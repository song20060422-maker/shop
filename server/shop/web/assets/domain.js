/**
 * Shop 商城买家端 —— 真实后端适配层
 * ============================================================
 * 依据原型图（原型图/live-prototype）对齐页面交互；本文件替代原型的
 * 本地 mock 数据层（window.ShopDomain），全部读写对接 Go 后端：
 *   - 同步读：基于本地缓存（refresh() 拉取后端真实数据）
 *   - 异步写：下单/取消/退款/支付通知等调用真实 REST 接口
 *   - 演示模式钩子：/api/demo/reset、/api/test/expire-orders、/api/pay/channel-mode
 * 接口面与原型 domain.js 完全一致，buyer.js / payment.js 仅做 async 适配。
 *
 * 后端契约（原生 camelCase + 兼容层 snake_case）：
 *   POST /api/login            -> data:{token,user,expireAt}
 *   GET  /api/products         -> data:[{id:SKU,db_id,price_cents,available,on_sale,...}]（兼容层）
 *   GET  /api/users            -> data:[{id,username,role,displayName}]
 *   POST /api/orders           -> data:{order:model.Order, availableAfter}（productId 为数据库主键）
 *   GET  /api/orders           -> data:[model.Order]
 *   GET  /api/orders/:no       -> data:OrderDetail{order,payment,refund,shipment,notices,...}
 *   POST /api/orders/:no/cancel            -> data:model.Order（重复取消 2004）
 *   POST /api/orders/:no/refund            -> data:{refund:model.Refund}
 *   POST /api/orders/:no/alipay-pay        -> data:{payURL}
 *   POST /api/pay/notify       -> data:NotifyResult{result,orderNo,channelTxnId,amountCents,abnormalId}
 *   POST /api/pay/channel-mode -> 兼容层（refund_result: success|fail）
 *   POST /api/test/expire-orders / /api/demo/reset -> 演示钩子
 *   GET  /api/ops/summary      -> data:TradeSummary（占用/出库跨买家口径）
 *   GET  /api/ops/notices      -> data:[PaymentNotice]
 *   GET  /api/ops/abnormals    -> data:[AbnormalPayment]
 */
(function () {
  "use strict";

  var API_BASE = "";
  var DEMO_PASSWORD = "shop123"; // 演示口令：与部署说明一致，仅演示环境使用

  var REAL_TO_PROTO = {
    pending_payment: "PENDING",
    paid: "PAID",
    shipped: "SHIPPED",
    refund_requested: "REFUND_REQUESTED",
    refund_failed: "REFUND_FAILED",
    refunded: "REFUNDED",
    cancelled: "CANCELLED",
    closed: "CLOSED"
  };

  // 兼容层契约状态字（下单/详情接口的 status 字段）兜底映射
  var COMPAT_TO_PROTO = {
    WAIT_PAY: "PENDING",
    TO_SHIP: "PAID",
    SHIPPED: "SHIPPED",
    REFUND_PENDING: "REFUND_REQUESTED",
    REFUNDED: "REFUNDED",
    CANCELLED: "CANCELLED",
    CLOSED: "CLOSED"
  };

  var ORDER_STATUS = {
    PENDING: "PENDING",
    PAID: "PAID",
    SHIPPED: "SHIPPED",
    REFUND_REQUESTED: "REFUND_REQUESTED",
    REFUND_FAILED: "REFUND_FAILED",
    REFUNDED: "REFUNDED",
    CANCELLED: "CANCELLED",
    CLOSED: "CLOSED"
  };

  var STATUS_LABEL = {
    PENDING: "待付款",
    PAID: "待发货",
    SHIPPED: "已发货",
    REFUND_REQUESTED: "退款申请中",
    REFUND_FAILED: "退款渠道失败",
    REFUNDED: "已退款",
    CANCELLED: "已取消",
    CLOSED: "已关闭"
  };

  var STATUS_TONE = {
    PENDING: "pending",
    PAID: "paid",
    SHIPPED: "shipped",
    REFUND_REQUESTED: "refund",
    REFUND_FAILED: "refund",
    REFUNDED: "refunded",
    CANCELLED: "closed",
    CLOSED: "closed"
  };

  var NOTICE_LABEL = {
    ACCEPTED: "已收账",
    DUPLICATED: "重复通知",
    ABNORMAL: "异常收款",
    REJECTED: "已拒绝"
  };

  var REFUND_SOURCE_LABEL = {
    PRESALE: "售前退款",
    ABNORMAL_LATE: "迟到支付退款"
  };

  // 演示用渠道密钥（与 server config.shop.yaml 的 pay-notify-secret 一致；
  // 可被 localStorage 覆盖，仅用于模拟支付渠道发起通知）
  var CHANNEL_SECRET = localStorage.getItem("shop.channelSecret") || "channel-shared-secret";

  var cache = {
    products: [],
    buyers: [],
    orders: [],
    payments: [], // 有效收款（来自订单详情 payment）
    notifications: [], // 全量渠道通知留痕（ops 视图）
    refunds: [],
    anomalies: [],
    shipments: [],
    channelRefundMode: "success",
    lastRefreshAt: 0
  };

  var tokens = {};
  var usersById = {};
  var productsByDbId = {};
  var listeners = [];
  var refreshPromise = null;

  function clone(value) {
    return JSON.parse(JSON.stringify(value));
  }

  function notifyListeners() {
    listeners.forEach(function (fn) {
      try { fn(); } catch (e) { console.error(e); }
    });
  }

  function money(cents) {
    return Math.round(Number(cents || 0));
  }

  function iso(value) {
    if (!value) return "";
    var d = new Date(value);
    return Number.isNaN(d.getTime()) ? String(value) : d.toISOString();
  }

  function productImage(sku) {
    if (sku === "SKU-B") return "/shop/assets/product-b.png";
    return "/shop/assets/product-a-visual.png";
  }

  function usernameOf(buyerNumericId) {
    var user = usersById[String(buyerNumericId)];
    return user ? user.username : String(buyerNumericId);
  }

  function skuOf(productDbId) {
    var product = productsByDbId[String(productDbId)];
    return product ? product.id : String(productDbId);
  }

  // ---- 低层请求 -----------------------------------------------------------
  function request(method, path, body, token) {
    var headers = { "Content-Type": "application/json" };
    if (token) headers["Authorization"] = "Bearer " + token;
    return fetch(API_BASE + path, {
      method: method,
      headers: headers,
      body: body === undefined ? undefined : JSON.stringify(body)
    })
      .then(function (resp) {
        return resp
          .json()
          .catch(function () { return {}; })
          .then(function (payload) {
            return { httpOk: resp.ok, status: resp.status, payload: payload };
          });
      })
      .catch(function () {
        return { httpOk: false, status: 0, payload: { message: "网络异常，请稍后重试" } };
      });
  }

  function tokenFor(username) {
    return tokens[username] || null;
  }

  function login(username, password) {
    return request("POST", "/api/login", { username: username, password: password || DEMO_PASSWORD }).then(
      function (resp) {
        var data = resp.payload.data || {};
        if (resp.httpOk && data.token) {
          tokens[username] = data.token;
          if (data.user && data.user.id !== undefined) usersById[String(data.user.id)] = data.user;
          return { ok: true, username: username, token: data.token };
        }
        return { ok: false, message: resp.payload.msg || resp.payload.message || "登录失败" };
      }
    );
  }

  function ensureBuyerToken(buyerId, password) {
    if (tokenFor(buyerId)) return Promise.resolve({ ok: true });
    return login(buyerId, password);
  }

  function ensureStaffToken(username, password) {
    if (tokenFor(username)) return Promise.resolve({ ok: true });
    return login(username, password);
  }

  // ---- 数据映射（后端视图 → 原型形态） ------------------------------------
  function mapProduct(p) {
    // 兼容层 productView：id=SKU，db_id=数据库主键
    var product = {
      id: p.sku || p.id,
      sku: p.sku || p.id,
      name: p.name,
      subtitle: (p.sku || p.id) === "SKU-A" ? "标准款 · 上架演示商品" : "配件款 · 上架演示商品",
      spec: (p.sku || p.id) === "SKU-A" ? "标配款" : "配件款",
      priceCents: money(p.price_cents !== undefined ? p.price_cents : p.priceCents),
      image: productImage(p.sku || p.id),
      onSale: !!(p.on_sale !== undefined ? p.on_sale : p.onSale),
      dbId: p.db_id || p.productId,
      initialStock: p.initial_stock !== undefined ? p.initial_stock : p.initialStock || 0,
      available: p.available || 0
    };
    productsByDbId[String(product.dbId)] = product;
    return product;
  }

  function mapBuyer(u) {
    usersById[String(u.id)] = u;
    return {
      id: u.username,
      name: u.displayName || u.username,
      phone: "账号 " + u.username,
      role: u.role
    };
  }

  function mapOrder(o) {
    // 原生 model.Order（camelCase）或兼容层 orderView（raw_status 可靠）均可
    var rawStatus = o.raw_status || o.status;
    var sku = typeof o.productId === "string" ? o.productId : skuOf(o.productId);
    if (sku === "undefined") sku = o.product_id || "";
    var product = cache.products.find(function (p) { return p.id === sku; });
    var order = {
      id: o.orderNo,
      buyerId: typeof o.buyerId === "string" ? o.buyerId : usernameOf(o.buyer_id !== undefined ? o.buyer_id : o.buyerId),
      productId: sku,
      productName: o.productName || o.product_name || (product ? product.name : sku),
      productSpec: product ? product.spec : "",
      productImage: productImage(sku),
      unitPriceCents: money(o.priceCents !== undefined ? o.priceCents : o.price_cents),
      qty: o.qty !== undefined ? o.qty : o.quantity,
      amountCents: money(o.amountCents !== undefined ? o.amountCents : o.amount_cents),
      status: REAL_TO_PROTO[rawStatus] || COMPAT_TO_PROTO[o.status] || o.status,
      createdAt: iso(o.createdAt || o.created_at),
      expireAt: iso(o.expireAt || o.expire_at),
      paidAt: iso(o.paidAt || o.paid_at),
      closedAt: iso(o.closedAt || o.closed_at),
      closeReason: rawStatus === "cancelled" ? "买家取消" : rawStatus === "closed" ? "超时未支付" : "",
      paymentTxnNo: o.paymentTxnNo || "",
      carrier: o.carrier || "",
      trackingNo: o.trackingNo || o.tracking_no || "",
      shippedAt: iso(o.shippedAt)
    };
    if (order.status === "SHIPPED" && order.trackingNo) {
      upsertShipment({
        orderId: order.id,
        carrier: order.carrier,
        trackingNo: order.trackingNo,
        shippedAt: order.shippedAt
      });
    }
    return order;
  }

  function mapRefund(r) {
    // model.Refund（原生 camelCase）
    var source = (r.source || "presale").toLowerCase() === "abnormal_late" ? "ABNORMAL_LATE" : "PRESALE";
    var channelStatus = (r.channelStatus || "").toUpperCase();
    var stockRestored = typeof r.restored === "boolean" ? r.restored : source === "PRESALE" && channelStatus === "SUCCESS";
    return {
      id: r.refundNo,
      orderId: r.orderNo,
      source: source,
      requestedAt: iso(r.createdAt),
      requestReason: r.source === "presale" ? (r.auditReason || "买家申请退款") : "迟到支付异常退款",
      reviewHistory: [
        {
          by: "cs1",
          decision: (r.auditStatus || "").toUpperCase(),
          reason: r.auditReason || "",
          at: iso(r.auditAt || r.createdAt)
        }
      ],
      status: (r.auditStatus || "").toUpperCase(),
      channelStatus: channelStatus,
      attempts: r.retryCount || 0,
      refundAmountCents: money(r.amountCents),
      stockRestored: stockRestored,
      lastError: channelStatus === "FAILED" ? "渠道明确失败，可重试" : ""
    };
  }

  function mapShipment(s) {
    return {
      orderId: s.orderNo,
      carrier: s.carrier,
      trackingNo: s.trackingNo,
      shippedAt: iso(s.createdAt)
    };
  }

  function upsertShipment(shipment) {
    if (!shipment.orderId || !shipment.trackingNo) return;
    var index = cache.shipments.findIndex(function (s) { return s.orderId === shipment.orderId; });
    if (index >= 0) cache.shipments[index] = shipment;
    else cache.shipments.push(shipment);
  }

  function mapNotice(n) {
    // PaymentNotice（原生 camelCase）：全量通知留痕
    return {
      id: n.id ? "N" + n.id : "",
      orderId: n.orderNo,
      txnNo: n.channelTxnId,
      amountCents: money(n.amountCents),
      status: (n.result || "").toUpperCase(),
      signatureOk: n.signOk !== undefined ? !!n.signOk : true,
      createdAt: iso(n.createdAt),
      note: n.reason || ""
    };
  }

  function mapAnomaly(a) {
    // AbnormalPayment（原生 camelCase）
    return {
      id: a.id,
      txnNo: a.channelTxnId,
      orderId: a.orderNo,
      buyerId: "",
      amountCents: money(a.amountCents),
      reason: a.reason || "",
      status: (a.status || "").toUpperCase(),
      createdAt: iso(a.createdAt),
      handledBy: a.handledBy || "",
      handledAt: iso(a.handledAt),
      refundId: a.refundId || "",
      refundStatus: ""
    };
  }

  // ---- 同步读接口（buyer.js / payment.js 使用） ---------------------------
  function getState() {
    return {
      products: clone(cache.products),
      buyers: clone(cache.buyers),
      orders: clone(cache.orders),
      payments: clone(cache.payments),
      notifications: clone(cache.notifications),
      refunds: clone(cache.refunds),
      anomalies: clone(cache.anomalies),
      shipments: clone(cache.shipments),
      channelRefundMode: cache.channelRefundMode,
      lastRefreshAt: cache.lastRefreshAt,
      activeBuyer: null,
      activeStaff: null,
      ledger: []
    };
  }

  function products() {
    return clone(cache.products);
  }

  function buyers() {
    return clone(cache.buyers);
  }

  function getProduct(productId) {
    return cache.products.find(function (p) { return p.id === productId; }) || null;
  }

  function getBuyer(buyerId) {
    return cache.buyers.find(function (b) { return b.id === buyerId; }) || null;
  }

  function getOrder(orderId) {
    return cache.orders.find(function (o) { return o.id === orderId; }) || null;
  }

  function getRefundForOrder(orderId) {
    return cache.refunds.find(function (r) { return r.orderId === orderId; }) || null;
  }

  function getShipment(orderId) {
    return cache.shipments.find(function (s) { return s.orderId === orderId; }) || null;
  }

  function productMetrics(productId) {
    var product = getProduct(productId);
    var available = product ? product.available : 0;
    // 占用/出库优先取运营核查汇总（跨买家口径准确；未登录 ops 时回退本地订单推算）
    if (product && typeof product.occupied === "number") {
      return {
        productId: productId,
        available: available,
        occupied: product.occupied,
        shipped: product.shipped || 0,
        total: product.initialStock
      };
    }
    var occupied = 0;
    var shipped = 0;
    cache.orders.forEach(function (order) {
      if (order.productId !== productId) return;
      if (["PENDING", "PAID", "REFUND_REQUESTED", "REFUND_FAILED"].indexOf(order.status) >= 0) {
        occupied += order.qty;
      }
      if (order.status === "SHIPPED") shipped += order.qty;
    });
    return {
      productId: productId,
      available: available,
      occupied: occupied,
      shipped: shipped,
      total: (product ? product.initialStock : 0) || available + occupied + shipped
    };
  }

  function financialSummary() {
    var validCents = 0;
    var refundCents = 0;
    var abnormalCents = 0;
    var abnormalCount = 0;
    cache.payments.forEach(function (n) {
      if (n.status === "ACCEPTED") validCents += n.amountCents;
    });
    cache.anomalies.forEach(function (a) {
      abnormalCents += a.amountCents;
      abnormalCount += 1;
    });
    cache.refunds.forEach(function (r) {
      if (r.channelStatus === "SUCCESS") refundCents += r.refundAmountCents;
    });
    return {
      validCents: validCents,
      refundCents: refundCents,
      abnormalCents: abnormalCents,
      abnormalCount: abnormalCount,
      netCents: validCents - refundCents
    };
  }

  function inventoryInvariant() {
    var rows = cache.products.map(function (product) {
      var metrics = productMetrics(product.id);
      var expected = product.initialStock;
      var calculated = metrics.available + metrics.occupied + metrics.shipped;
      return {
        productId: product.id,
        expected: expected,
        calculated: calculated,
        ok: expected === calculated
      };
    });
    return { ok: rows.every(function (r) { return r.ok; }), products: rows };
  }

  function paymentInvariant() {
    var seen = {};
    var duplicates = [];
    cache.payments.forEach(function (n) {
      if (n.status !== "ACCEPTED") return;
      if (seen[n.orderId]) duplicates.push(n.orderId);
      seen[n.orderId] = true;
    });
    return { ok: duplicates.length === 0, duplicates: duplicates };
  }

  function validateInvariants() {
    var inventory = inventoryInvariant();
    var payments = paymentInvariant();
    var issues = [];
    if (!inventory.ok) {
      inventory.products.forEach(function (row) {
        if (!row.ok) issues.push(row.productId + " 库存不守恒（应为 " + row.expected + "，实测 " + row.calculated + "）");
      });
    }
    if (!payments.ok) issues.push("存在重复有效收款：" + payments.duplicates.join("、"));
    return { ok: issues.length === 0, issues: issues };
  }

  function now() {
    return Date.now();
  }

  function subscribe(listener) {
    listeners.push(listener);
    return function () {
      listeners = listeners.filter(function (fn) { return fn !== listener; });
    };
  }

  // ---- 数据刷新（异步；完成后触发订阅者重渲染） ---------------------------
  function fetchProductSnapshot(opsToken) {
    if (!opsToken) return Promise.resolve();
    return request("GET", "/api/ops/summary", undefined, opsToken).then(function (resp) {
      if (!resp.httpOk) return;
      var data = resp.payload.data || {};
      (data.products || []).forEach(function (row) {
        var product = cache.products.find(function (p) { return p.sku === row.sku || p.id === row.sku; });
        if (product) {
          product.available = row.available;
          product.initialStock = row.initialStock || product.initialStock;
          product.occupied = row.occupied || 0;
          product.shipped = row.shippedQty || 0;
        }
      });
    });
  }

  function upsertOrder(mapped) {
    var index = cache.orders.findIndex(function (o) { return o.id === mapped.id; });
    if (index >= 0) cache.orders[index] = mapped;
    else cache.orders.push(mapped);
  }

  function fetchBuyerOrders(buyerId, token) {
    if (!token) return Promise.resolve();
    return request("GET", "/api/orders", undefined, token).then(function (resp) {
      if (!resp.httpOk) return;
      (resp.payload.data || []).forEach(function (o) {
        upsertOrder(mapOrder(o));
      });
    });
  }

  function fetchBuyerOrdersAll() {
    var entries = Object.keys(tokens)
      .filter(function (name) {
        return !!cache.buyers.find(function (b) { return b.id === name; });
      })
      .map(function (name) { return { buyerId: name, token: tokens[name] }; });
    return entries.reduce(function (chain, entry) {
      return chain.then(function () { return fetchBuyerOrders(entry.buyerId, entry.token); });
    }, Promise.resolve());
  }

  // 订单全链路视图（ops token）：补充有效收款 / 退款 / 运单 / 通知留痕（列表视图不含）
  function fetchOrderDetail(orderId) {
    var opsToken = tokenFor("ops1");
    if (!opsToken) return Promise.resolve();
    if (!getOrder(orderId)) return Promise.resolve();
    return request("GET", "/api/ops/orders/" + encodeURIComponent(orderId), undefined, opsToken).then(function (resp) {
      if (!resp.httpOk) return;
      var d = resp.payload.data || {};
      if (d.order) upsertOrder(mapOrder(d.order));
      if (d.payment) {
        var exists = cache.payments.some(function (p) { return p.orderId === orderId; });
        if (!exists) {
          cache.payments.push({
            id: "P" + d.payment.id,
            orderId: d.payment.orderNo,
            txnNo: d.payment.channelTxnId,
            amountCents: money(d.payment.amountCents),
            status: "ACCEPTED",
            createdAt: iso(d.payment.paidAt || d.payment.createdAt),
            note: ""
          });
        }
      }
      if (d.refund) {
        var mappedRefund = mapRefund(d.refund);
        var rIndex = cache.refunds.findIndex(function (r) { return r.orderId === mappedRefund.orderId; });
        if (rIndex >= 0) cache.refunds[rIndex] = mappedRefund;
        else cache.refunds.push(mappedRefund);
      }
      if (d.shipment) {
        var shipment = mapShipment(d.shipment);
        upsertShipment(shipment);
        var shipped = getOrder(orderId);
        if (shipped) {
          shipped.carrier = d.shipment.carrier;
          shipped.trackingNo = d.shipment.trackingNo;
          shipped.shippedAt = iso(d.shipment.createdAt);
        }
      }
      if (d.notices && d.notices.length) {
        d.notices.forEach(function (n) {
          var mapped = mapNotice(n);
          var nIndex = cache.notifications.findIndex(function (x) { return x.id === mapped.id; });
          if (nIndex < 0) cache.notifications.push(mapped);
        });
      }
    });
  }

  function fetchAllOrderDetails() {
    var ids = cache.orders.slice(0, 40).map(function (order) { return order.id; });
    return ids.reduce(function (chain, orderId) {
      return chain.then(function () { return fetchOrderDetail(orderId); });
    }, Promise.resolve());
  }

  function fetchOpsViews(opsToken) {
    if (!opsToken) return Promise.resolve();
    return request("GET", "/api/ops/notices?limit=80", undefined, opsToken)
      .then(function (resp) {
        if (resp.httpOk) {
          cache.notifications = (resp.payload.data || []).map(mapNotice);
        }
        return request("GET", "/api/ops/abnormals", undefined, opsToken);
      })
      .then(function (resp) {
        if (resp.httpOk) {
          cache.anomalies = (resp.payload.data || []).map(mapAnomaly);
        }
        return request("GET", "/api/pay/channel", undefined, null);
      })
      .then(function (resp) {
        if (resp.httpOk) {
          var data = resp.payload.data || {};
          cache.channelRefundMode = data.refundMode || "success";
        }
      })
      .catch(function () {});
  }

  function refresh() {
    if (refreshPromise) return refreshPromise;
    refreshPromise = request("GET", "/api/products")
      .then(function (resp) {
        if (resp.httpOk) {
          cache.products = (resp.payload.data || resp.payload.products || []).map(mapProduct);
        }
        return request("GET", "/api/users");
      })
      .then(function (resp) {
        if (resp.httpOk) {
          cache.buyers = (resp.payload.data || [])
            .filter(function (u) { return u.role === "buyer"; })
            .map(mapBuyer);
        }
        return fetchBuyerOrdersAll();
      })
      .then(function () {
        return fetchAllOrderDetails();
      })
      .then(function () {
        return fetchProductSnapshot(tokenFor("ops1"));
      })
      .then(function () {
        return fetchOpsViews(tokenFor("ops1"));
      })
      .then(function () {
        cache.lastRefreshAt = Date.now();
        notifyListeners();
      })
      .catch(function () {
        notifyListeners();
      })
      .then(function () {
        refreshPromise = null;
      });
    return refreshPromise;
  }

  // ---- 异步写接口（对接真实后端） ----------------------------------------
  function createOrder(buyerId, productId, qty) {
    return ensureBuyerToken(buyerId).then(function (auth) {
      if (!auth.ok) return { ok: false, code: "unauthorized", message: auth.message };
      var product = getProduct(productId);
      if (!product) return { ok: false, code: "not_found", message: "商品不存在" };
      return request("POST", "/api/orders", { productId: product.dbId, qty: qty }, tokenFor(buyerId)).then(
        function (resp) {
          var payload = resp.payload;
          if (!resp.httpOk) {
            return { ok: false, code: "order_failed", message: payload.msg || payload.message || "下单失败" };
          }
          var order = mapOrder(payload.data.order);
          order.buyerId = buyerId;
          upsertOrder(order);
          return refresh().then(function () {
            return { ok: true, code: "created", data: { order: getOrder(order.id) || order } };
          });
        }
      );
    });
  }

  function cancelOrder(orderId, buyerId) {
    var order = getOrder(orderId);
    if (order && order.status === "CANCELLED") {
      return Promise.resolve({ ok: true, code: "duplicated", message: "订单已取消，重复请求幂等处理", data: { order: order } });
    }
    return ensureBuyerToken(buyerId).then(function (auth) {
      if (!auth.ok) return { ok: false, code: "unauthorized", message: auth.message };
      return request("POST", "/api/orders/" + encodeURIComponent(orderId) + "/cancel", {}, tokenFor(buyerId)).then(
        function (resp) {
          if (!resp.httpOk) {
            // 2003（状态不允许）/2004（已关闭）：缓存中订单已为取消态时按幂等重复请求处理
            var code = resp.payload.code;
            var cached = getOrder(orderId);
            if ((code === 2003 || code === 2004) && cached && cached.status === "CANCELLED") {
              return refresh().then(function () {
                return { ok: true, code: "duplicated", message: "订单已取消，重复请求幂等处理", data: { order: getOrder(orderId) } };
              });
            }
            return { ok: false, code: "cancel_failed", message: resp.payload.msg || resp.payload.message || "取消失败" };
          }
          if (resp.payload.data) upsertOrder(mapOrder(resp.payload.data));
          return refresh().then(function () {
            return { ok: true, code: "cancelled", data: { order: getOrder(orderId) } };
          });
        }
      );
    });
  }

  function applyRefund(orderId, buyerId, reason) {
    return ensureBuyerToken(buyerId).then(function (auth) {
      if (!auth.ok) return { ok: false, code: "unauthorized", message: auth.message };
      return request(
        "POST",
        "/api/orders/" + encodeURIComponent(orderId) + "/refund",
        { orderNo: orderId, reason: reason },
        tokenFor(buyerId)
      ).then(function (resp) {
        if (!resp.httpOk) {
          return { ok: false, code: "refund_failed", message: resp.payload.msg || resp.payload.message || "申请失败" };
        }
        return refresh().then(function () {
          return { ok: true, code: "refund_requested", data: { refund: getRefundForOrder(orderId) } };
        });
      });
    });
  }

  // 支付宝真实渠道：请求支付链接，由真实支付宝沙箱完成收款，
  // 渠道异步通知（notify URL 动态取请求域名）推进订单状态
  function alipayPay(orderId, buyerId) {
    return ensureBuyerToken(buyerId).then(function (auth) {
      if (!auth.ok) return { ok: false, message: auth.message };
      return request("POST", "/api/orders/" + encodeURIComponent(orderId) + "/alipay-pay", {}, tokenFor(buyerId)).then(
        function (resp) {
          if (!resp.httpOk) {
            return { ok: false, code: "alipay_failed", message: resp.payload.msg || resp.payload.message || "支付宝下单失败" };
          }
          return refresh().then(function () {
            return { ok: true, code: "alipay_created", payURL: (resp.payload.data && resp.payload.data.payURL) || "", data: resp.payload.data };
          });
        }
      );
    });
  }

  // 渠道通知签名：HMAC-SHA256(secret, "订单号|交易号|金额(分)")，与 utils.PaySign 一致
  function hmacSign(orderNo, txnNo, amountCents) {
    var encoder = new TextEncoder();
    return crypto.subtle
      .importKey("raw", encoder.encode(CHANNEL_SECRET), { name: "HMAC", hash: "SHA-256" }, false, ["sign"])
      .then(function (key) {
        return crypto.subtle.sign("HMAC", key, encoder.encode(orderNo + "|" + txnNo + "|" + money(amountCents)));
      })
      .then(function (sig) {
        return Array.from(new Uint8Array(sig))
          .map(function (b) { return ("0" + b.toString(16)).slice(-2); })
          .join("");
      });
  }

  function sendNotify(orderNo, txnNo, amountCents, sign) {
    return request("POST", "/api/pay/notify", {
      orderNo: orderNo,
      channelTxnId: txnNo,
      amountCents: money(amountCents),
      sign: sign
    }).then(function (resp) {
      var data = resp.payload.data || {};
      var result = (data.result || "").toLowerCase();
      var okResult = resp.httpOk && result === "accepted";
      return refresh().then(function () {
        return {
          httpOk: resp.httpOk,
          ok: okResult,
          code: result || (resp.httpOk ? "unknown" : "network_error"),
          message: data.reason || resp.payload.msg || "通知已受理",
          result: result.toUpperCase(),
          data: {
            payment: { orderId: orderNo, txnNo: txnNo, amountCents: money(amountCents) },
            anomaly: data.abnormalId ? { id: data.abnormalId } : null
          }
        };
      });
    });
  }

  function simulatePayment(orderId, mode, buyerId) {
    var order = getOrder(orderId);
    if (!order) return Promise.resolve({ ok: false, code: "not_found", message: "订单不存在" });
    if (mode === "network_failure") {
      return Promise.resolve({
        ok: false,
        code: "network_unknown",
        message: "网络中断，支付结果暂时未知，订单仍保持待付款。"
      });
    }
    var txnNo = "WX_PAY_" + Date.now() + "_" + Math.floor(Math.random() * 1e6);
    var amount = order.amountCents;
    if (mode === "amount_mismatch") amount += 1;
    return hmacSign(orderId, txnNo, amount).then(function (sign) {
      return sendNotify(orderId, txnNo, amount, sign);
    });
  }

  function receivePayment(payload) {
    var orderId = payload.orderId;
    var order = getOrder(orderId);
    var amount = payload.amountCents || (order ? order.amountCents : 0);
    return hmacSign(orderId, payload.txnNo, amount).then(function (sign) {
      return sendNotify(orderId, payload.txnNo, amount, sign);
    });
  }

  function timeoutOrder(orderId) {
    // 演示钩子：真实关闭到期待付款订单；若目标单仍未到期且为最新待付款单则强制到期
    return request("POST", "/api/test/expire-orders").then(function (resp) {
      return refresh().then(function () {
        var order = getOrder(orderId);
        return {
          ok: !!order && ["CLOSED", "CANCELLED"].indexOf(order.status) >= 0,
          message: order && order.status === "CLOSED" ? "订单已超时关闭，占用已释放" : "订单尚未关闭",
          data: { expired: (resp.payload && resp.payload.expired) || [] }
        };
      });
    });
  }

  function advanceClock(ms) {
    return request("POST", "/api/test/expire-orders").then(function (resp) {
      return refresh().then(function () {
        return { ok: resp.httpOk, advanced: ms, expired: (resp.payload && resp.payload.expired) || [] };
      });
    });
  }

  function tickTimeouts() {
    // 服务端已内置超时关单调度；客户端仅刷新状态
    return Promise.resolve([]);
  }

  function sendTestNotification(payload) {
    var orderId = payload.orderId;
    var order = getOrder(orderId);
    var amount = payload.amountCents || (order ? order.amountCents : 0) + (payload.amountDelta || 0);
    var txnNo = payload.txnNo || "WX_TEST_" + Date.now();
    if (payload.signatureOk === false) {
      return sendNotify(orderId, txnNo, amount, "deadbeef").then(function (resp) {
        return { ok: !resp.ok, code: resp.code, message: resp.message, result: resp.result, data: resp.data };
      });
    }
    return hmacSign(orderId, txnNo, amount).then(function (sign) {
      return sendNotify(orderId, txnNo, amount, sign);
    });
  }

  function resetState() {
    return request("POST", "/api/demo/reset").then(function (resp) {
      cache.orders = [];
      cache.payments = [];
      cache.notifications = [];
      cache.refunds = [];
      cache.anomalies = [];
      cache.shipments = [];
      tokens = {};
      return refresh().then(function () {
        return { ok: resp.httpOk, message: resp.httpOk ? "已重置" : "重置失败（需要演示模式）" };
      });
    });
  }

  // ---- 演示员工动作（场景编排 / 买家页复用） ------------------------------
  function setChannelRefundMode(mode) {
    return request("POST", "/api/pay/channel-mode", { refund_result: mode }).then(function (resp) {
      cache.channelRefundMode = mode;
      return { ok: resp.httpOk, message: resp.payload.msg || "" };
    });
  }

  function shipOrder(orderId, carrier, trackingNo) {
    return ensureStaffToken("wh1").then(function (auth) {
      if (!auth.ok) return { ok: false, message: auth.message };
      return request(
        "POST",
        "/api/warehouse/orders/" + encodeURIComponent(orderId) + "/ship",
        { carrier: carrier || "顺丰", trackingNo: trackingNo || "SF" + Date.now() },
        tokens["wh1"]
      ).then(function (resp) {
        return refresh().then(function () {
          return {
            ok: resp.httpOk,
            code: resp.payload.code || "",
            message: resp.payload.msg || (resp.httpOk ? "发货成功" : "发货被拒绝"),
            data: { shipment: getShipment(orderId) }
          };
        });
      });
    });
  }

  function reviewRefund(orderId, decision, reason) {
    return ensureStaffToken("cs1").then(function (auth) {
      if (!auth.ok) return { ok: false, message: auth.message };
      var refund = getRefundForOrder(orderId);
      if (!refund) return { ok: false, message: "退款单不存在" };
      return request(
        "POST",
        "/api/cs/refunds/" + encodeURIComponent(refund.id) + "/decide",
        { approve: decision === "approve", reason: reason || "客服审核" },
        tokens["cs1"]
      ).then(function (resp) {
        return refresh().then(function () {
          return {
            ok: resp.httpOk,
            message: resp.payload.msg || (resp.httpOk ? "审核完成" : "审核被拒绝"),
            data: { refund: getRefundForOrder(orderId) }
          };
        });
      });
    });
  }

  function retryRefund(orderId) {
    return ensureStaffToken("cs1").then(function (auth) {
      if (!auth.ok) return { ok: false, message: auth.message };
      var refund = getRefundForOrder(orderId);
      if (!refund) return { ok: false, message: "退款单不存在" };
      return request("POST", "/api/cs/refunds/" + encodeURIComponent(refund.id) + "/retry", {}, tokens["cs1"]).then(
        function (resp) {
          return refresh().then(function () {
            return {
              ok: resp.httpOk,
              message: resp.payload.msg || (resp.httpOk ? "重试成功，退款到账" : "重试失败"),
              data: { refund: getRefundForOrder(orderId) }
            };
          });
        }
      );
    });
  }

  function handleAnomaly(txnNo, note, channelMode) {
    return ensureStaffToken("cs1").then(function (auth) {
      if (!auth.ok) return { ok: false, message: auth.message };
      var anomaly = cache.anomalies.find(function (a) { return a.txnNo === txnNo && a.status === "PENDING"; });
      if (!anomaly) return { ok: false, message: "待处理异常不存在" };
      var modePromise = channelMode ? setChannelRefundMode(channelMode) : Promise.resolve({ ok: true });
      return modePromise.then(function () {
        return request(
          "POST",
          "/api/cs/payment-exceptions/" + anomaly.id + "/refund",
          { reason: note || "客服处理迟到支付" },
          tokens["cs1"]
        );
      }).then(function (resp) {
        return refresh().then(function () {
          return { ok: resp.httpOk, message: resp.payload.msg || (resp.httpOk ? "处理完成" : "处理失败"), data: resp.payload.data };
        });
      });
    });
  }

  function setProductSale(productId, onSale) {
    return ensureStaffToken("ops1").then(function (auth) {
      if (!auth.ok) return { ok: false, message: auth.message };
      return request("PATCH", "/api/ops/products/" + encodeURIComponent(productId), { onSale: onSale }, tokens["ops1"]).then(
        function (resp) {
          return refresh().then(function () { return { ok: resp.httpOk, message: resp.payload.msg || "" }; });
        }
      );
    });
  }

  function updatePrice(productId, newPriceCents) {
    return ensureStaffToken("ops1").then(function (auth) {
      if (!auth.ok) return { ok: false, message: auth.message };
      return request(
        "PATCH",
        "/api/ops/products/" + encodeURIComponent(productId),
        { priceCents: money(newPriceCents) },
        tokens["ops1"]
      ).then(function (resp) {
        return refresh().then(function () { return { ok: resp.httpOk, message: resp.payload.msg || "" }; });
      });
    });
  }

  function adjustStock(productId, newAvailable) {
    return ensureStaffToken("ops1").then(function (auth) {
      if (!auth.ok) return { ok: false, message: auth.message };
      return request(
        "PATCH",
        "/api/ops/products/" + encodeURIComponent(productId),
        { available: newAvailable },
        tokens["ops1"]
      ).then(function (resp) {
        return refresh().then(function () { return { ok: resp.httpOk, message: resp.payload.msg || "" }; });
      });
    });
  }

  function logout(username) {
    delete tokens[username];
    cache.orders = cache.orders.filter(function (order) { return order.buyerId !== username; });
    notifyListeners();
    return Promise.resolve({ ok: true });
  }

  // ---- 验收场景编排（对应原型 runScenario，全部走真实接口） ----------------
  function scenarioLog(steps, text, tone) {
    steps.push({ text: text, tone: tone || "" });
  }

  function runScenario(name) {
    var steps = [];
    var password = DEMO_PASSWORD;
    var orderNo = null;

    return resetState()
      .then(function () { return login("buyer1", password); })
      .then(function () { return login("buyer2", password); })
      .then(function () { return login("wh1", password); })
      .then(function () { return login("cs1", password); })
      .then(function () { return login("ops1", password); })
      .then(function () { return refresh(); })
      .then(function () {
        if (name === "normal") {
          // 正常成交：下单 -> 支付成功 -> 仓库发货（占用转出库）
          return createOrder("buyer1", "SKU-A", 1).then(function (r1) {
            if (!r1.ok) throw new Error(r1.message);
            orderNo = r1.data.order.id;
            scenarioLog(steps, "buyer1 下单 SKU-A ×1，占用 1 件", "ok");
            return simulatePayment(orderNo, "success", "buyer1");
          }).then(function (r2) {
            if (!r2.ok) throw new Error(r2.message);
            scenarioLog(steps, "模拟渠道通知验签通过，订单进入待发货", "ok");
            return shipOrder(orderNo, "顺丰", "SF" + Date.now());
          }).then(function (r3) {
            scenarioLog(steps, r3.ok ? "仓库发货成功，出库 1 件" : "发货被拒绝：" + r3.message, r3.ok ? "ok" : "bad");
            scenarioLog(steps, "库存守恒：可售 4 + 占用 0 + 出库 1 = 5", "ok");
            return { ok: r3.ok, orderId: orderNo, steps: steps };
          });
        }
        if (name === "stock") {
          // 防超卖：20 个并发请求抢 4 件可售库存，仅 4 单成功
          var tasks = [];
          for (var i = 0; i < 20; i += 1) {
            tasks.push(createOrder(i % 2 === 0 ? "buyer1" : "buyer2", "SKU-A", 1));
          }
          return Promise.all(tasks).then(function (results) {
            var success = results.filter(function (r) { return r.ok; }).length;
            var failed = results.length - success;
            scenarioLog(steps, "并发下单 20 笔：成功 " + success + "，失败 " + failed + "（可售仅 4）", "ok");
            scenarioLog(steps, "可售数量归零，未出现负库存或超卖", success === 4 ? "ok" : "bad");
            return { ok: success === 4, orderId: null, steps: steps, success: success, failed: failed };
          });
        }
        if (name === "refund") {
          // 售前退款：申请 -> 发货拦截 -> 审核（渠道失败）-> 重试成功
          return createOrder("buyer1", "SKU-A", 1).then(function (r1) {
            if (!r1.ok) throw new Error(r1.message);
            orderNo = r1.data.order.id;
            return simulatePayment(orderNo, "success", "buyer1");
          }).then(function (r2) {
            if (!r2.ok) throw new Error(r2.message);
            scenarioLog(steps, "订单已支付，申请售前退款", "ok");
            return applyRefund(orderNo, "buyer1", "七天无理由退货");
          }).then(function (r3) {
            if (!r3.ok) throw new Error(r3.message);
            return shipOrder(orderNo, "顺丰", "SF" + Date.now());
          }).then(function (r4) {
            scenarioLog(steps, r4.ok ? "异常：已退款申请仍发货成功" : "发货被拒绝（退款申请中不可发货）", r4.ok ? "bad" : "ok");
            return setChannelRefundMode("fail");
          }).then(function () {
            return reviewRefund(orderNo, "approve", "凭证核验通过，同意退款");
          }).then(function (r6) {
            scenarioLog(steps, "客服同意退款，渠道返回明确失败", "ok");
            return setChannelRefundMode("success");
          }).then(function () {
            return retryRefund(orderNo);
          }).then(function (r8) {
            scenarioLog(steps, r8.ok ? "渠道重试成功，退款到账，库存恢复 1 件" : "重试失败：" + r8.message, r8.ok ? "ok" : "bad");
            return { ok: r8.ok, orderId: orderNo, steps: steps };
          });
        }
        if (name === "late") {
          // 迟到支付：超时关单 -> 渠道补发通知 -> 异常收款 -> 客服退款（不恢复库存）
          return createOrder("buyer1", "SKU-A", 1).then(function (r1) {
            if (!r1.ok) throw new Error(r1.message);
            orderNo = r1.data.order.id;
            return timeoutOrder(orderNo);
          }).then(function (r2) {
            scenarioLog(steps, r2.ok ? "订单超时关闭，库存占用已释放" : "关单结果：" + r2.message, r2.ok ? "ok" : "ok");
            return sendTestNotification({ orderId: orderNo, amountDelta: 0, source: "late" });
          }).then(function (r4) {
            scenarioLog(steps, "迟到通知到账：" + (r4.result === "ABNORMAL" ? "进入异常收款列表，订单未复活" : r4.message), r4.result === "ABNORMAL" ? "ok" : "bad");
            var anomaly = cache.anomalies.find(function (a) { return a.orderId === orderNo; });
            return handleAnomaly(anomaly ? anomaly.txnNo : null, "核实迟到支付，同意退款", "success");
          }).then(function (r5) {
            scenarioLog(steps, r5.ok ? "客服退款成功；库存保持已释放，不恢复" : "异常退款失败：" + r5.message, r5.ok ? "ok" : "bad");
            var product = getProduct("SKU-A");
            scenarioLog(steps, "库存核对：可售 " + (product ? product.available : "?") + "（未恢复）", "ok");
            return { ok: r5.ok, orderId: orderNo, steps: steps };
          });
        }
        if (name === "network") {
          // 网络中断与重复通知：同一渠道交易号只记一笔
          return createOrder("buyer1", "SKU-A", 1).then(function (r1) {
            if (!r1.ok) throw new Error(r1.message);
            orderNo = r1.data.order.id;
            var txnNo = "WX_PAY_NET_" + Date.now();
            return hmacSign(orderNo, txnNo, r1.data.order.amountCents).then(function (sign) {
              return sendNotify(orderNo, txnNo, r1.data.order.amountCents, sign);
            });
          }).then(function (r2) {
            if (!r2.ok) throw new Error(r2.message);
            scenarioLog(steps, "首次通知已收账", "ok");
            var payment = cache.payments.find(function (n) { return n.orderId === orderNo; }) || {};
            return sendNotify(orderNo, payment.txnNo, 10000, "resend");
          }).then(function (r3) {
            scenarioLog(steps, r3.result === "DUPLICATED" ? "重复通知被幂等去重，仅一笔有效收款" : "重复通知结果：" + r3.result, r3.result === "DUPLICATED" ? "ok" : "bad");
            return { ok: r3.result === "DUPLICATED", orderId: orderNo, steps: steps };
          });
        }
        if (name === "race") {
          // 退款与发货竞争：已支付订单上只能成功其一
          return createOrder("buyer1", "SKU-A", 1).then(function (r1) {
            if (!r1.ok) throw new Error(r1.message);
            orderNo = r1.data.order.id;
            return simulatePayment(orderNo, "success", "buyer1");
          }).then(function (r2) {
            if (!r2.ok) throw new Error(r2.message);
            return applyRefund(orderNo, "buyer1", "竞争演示：先申请退款");
          }).then(function (r3) {
            if (!r3.ok) throw new Error(r3.message);
            return shipOrder(orderNo, "顺丰", "SF" + Date.now());
          }).then(function (r4) {
            scenarioLog(steps, r4.ok ? "异常：退款申请期间发货成功" : "发货被拒绝（退款申请中）", r4.ok ? "bad" : "ok");
            return reviewRefund(orderNo, "reject", "凭证缺失，拒绝退款");
          }).then(function (r5) {
            scenarioLog(steps, r5.ok ? "客服拒绝退款，订单回到待发货" : "拒绝失败", "ok");
            return shipOrder(orderNo, "顺丰", "SF" + Date.now());
          }).then(function (r6) {
            scenarioLog(steps, r6.ok ? "拒绝后发货成功，互斥正确" : "发货仍被拒绝", r6.ok ? "ok" : "bad");
            return { ok: r6.ok, orderId: orderNo, steps: steps };
          });
        }
        if (name === "auth") {
          // 越权拦截：buyer2 操作 buyer1 订单 / 伪造通知被验签拒绝
          return createOrder("buyer1", "SKU-A", 1).then(function (r1) {
            if (!r1.ok) throw new Error(r1.message);
            orderNo = r1.data.order.id;
            return cancelOrder(orderNo, "buyer2");
          }).then(function (r2) {
            scenarioLog(steps, r2.ok ? "异常：越权取消成功" : "buyer2 取消 buyer1 订单被 403 拒绝", r2.ok ? "bad" : "ok");
            return sendTestNotification({ orderId: orderNo, amountDelta: 0, signatureOk: false });
          }).then(function (r3) {
            scenarioLog(steps, r3.ok ? "异常：伪造签名通知被接受" : "伪造签名通知被拒绝并留痕", r3.ok ? "bad" : "ok");
            return sendTestNotification({ orderId: orderNo, amountDelta: 1 });
          }).then(function (r4) {
            scenarioLog(steps, r4.ok ? "异常：金额不符通知被接受" : "金额不符通知被拒绝并留痕", r4.ok ? "bad" : "ok");
            var allBad = steps.some(function (s) { return s.tone === "bad"; });
            return { ok: !allBad, orderId: orderNo, steps: steps };
          });
        }
        scenarioLog(steps, "未知场景：" + name, "bad");
        return { ok: false, steps: steps };
      })
      .catch(function (error) {
        scenarioLog(steps, "场景执行中断：" + (error && error.message ? error.message : error), "bad");
        return { ok: false, orderId: orderNo, steps: steps };
      });
  }

  // ---- 原型兼容入口 -------------------------------------------------------
  function seedShowcase() {
    // 真实数据由后端迁移种子提供（SKU-A/SKU-B），无需本地造数
    return refresh();
  }

  window.ShopDomain = {
    ORDER_STATUS: ORDER_STATUS,
    STATUS_LABEL: STATUS_LABEL,
    STATUS_TONE: STATUS_TONE,
    NOTICE_LABEL: NOTICE_LABEL,
    REFUND_SOURCE_LABEL: REFUND_SOURCE_LABEL,
    DEMO_PASSWORD: DEMO_PASSWORD,
    now: now,
    getState: getState,
    products: products,
    buyers: buyers,
    getProduct: getProduct,
    getBuyer: getBuyer,
    getOrder: getOrder,
    getRefundForOrder: getRefundForOrder,
    getShipment: getShipment,
    productMetrics: productMetrics,
    financialSummary: financialSummary,
    inventoryInvariant: inventoryInvariant,
    validateInvariants: validateInvariants,
    subscribe: subscribe,
    refresh: refresh,
    login: login,
    tokenFor: tokenFor,
    createOrder: createOrder,
    cancelOrder: cancelOrder,
    applyRefund: applyRefund,
    alipayPay: alipayPay,
    logout: logout,
    simulatePayment: simulatePayment,
    receivePayment: receivePayment,
    sendNotify: sendNotify,
    sendTestNotification: sendTestNotification,
    timeoutOrder: timeoutOrder,
    advanceClock: advanceClock,
    tickTimeouts: tickTimeouts,
    resetState: resetState,
    setChannelRefundMode: setChannelRefundMode,
    shipOrder: shipOrder,
    reviewRefund: reviewRefund,
    retryRefund: retryRefund,
    handleAnomaly: handleAnomaly,
    setProductSale: setProductSale,
    updatePrice: updatePrice,
    adjustStock: adjustStock,
    runScenario: runScenario,
    seedShowcase: seedShowcase
  };
})();
