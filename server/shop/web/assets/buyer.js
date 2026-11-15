/**
 * Shop 商城买家端 —— 页面交互层
 * 与原型（原型图/live-prototype/assets/buyer.js）保持一致的 DOM 结构与交互；
 * 数据层 window.ShopDomain 为真实后端适配器（web/assets/domain.js）。
 * 与原型的差异仅在：mock 的同步写操作改为 async/await、买家切换时自动用
 * 演示口令登录、支付宝走真实渠道支付链接、员工工作台指向 GVA 管理端。
 */
(function () {
  "use strict";

  const D = window.ShopDomain;
  const $ = (selector, root) => (root || document).querySelector(selector);
  const $$ = (selector, root) => Array.from((root || document).querySelectorAll(selector));
  const buyerStorageKey = "shop-active-buyer";
  const minute = 60 * 1000;

  let selectedBuyerId = readPreference();
  let activeOrderFilter = "all";
  let activeMobileView = "market";
  let drawerOpen = false;
  let modalOpen = false;
  let activeOrderDetailId = null;
  let quantities = {};
  let lastExpiryCheck = 0;

  function readPreference() {
    try {
      return localStorage.getItem(buyerStorageKey) || "";
    } catch (error) {
      return "";
    }
  }

  function savePreference(value) {
    try {
      localStorage.setItem(buyerStorageKey, value);
    } catch (error) {
      // Preference is non-critical when local storage is blocked.
    }
  }

  function escapeHtml(value) {
    return String(value == null ? "" : value)
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#039;");
  }

  function icon(name) {
    return `<svg class="icon" aria-hidden="true"><use href="#i-${name}"></use></svg>`;
  }

  function money(cents) {
    return `¥${(Number(cents || 0) / 100).toFixed(2)}`;
  }

  function formatTime(value, withYear) {
    if (!value) return "—";
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return "—";
    return new Intl.DateTimeFormat("zh-CN", {
      timeZone: "Asia/Shanghai",
      year: withYear ? "numeric" : undefined,
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      second: withYear ? undefined : "2-digit",
      hour12: false
    }).format(date);
  }

  function remainingParts(isoTime) {
    const remaining = Math.max(0, new Date(isoTime).getTime() - D.now());
    const totalSeconds = Math.ceil(remaining / 1000);
    return {
      remaining,
      text: `${String(Math.floor(totalSeconds / 60)).padStart(2, "0")}:${String(
        totalSeconds % 60
      ).padStart(2, "0")}`
    };
  }

  function getSelectedBuyer() {
    return D.getBuyer(selectedBuyerId) || D.buyers()[0];
  }

  function buyerOrders() {
    return D.getState()
      .orders.filter((order) => order.buyerId === selectedBuyerId)
      .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
  }

  function orderTone(order) {
    if (order.status === D.ORDER_STATUS.PENDING) return "pending";
    if (order.status === D.ORDER_STATUS.PAID || order.status === D.ORDER_STATUS.SHIPPED) {
      return order.status === D.ORDER_STATUS.SHIPPED ? "shipped" : "paid";
    }
    if (
      order.status === D.ORDER_STATUS.REFUND_REQUESTED ||
      order.status === D.ORDER_STATUS.REFUNDED
    ) {
      return "refund";
    }
    if (order.status === D.ORDER_STATUS.REFUND_FAILED) return "danger";
    return "neutral";
  }

  function badgeClass(status) {
    const mapping = {
      PENDING: "badge-pending",
      PAID: "badge-paid",
      SHIPPED: "badge-shipped",
      REFUND_REQUESTED: "badge-refund",
      REFUND_FAILED: "badge-danger",
      REFUNDED: "badge-refund",
      CANCELLED: "badge-neutral",
      CLOSED: "badge-neutral"
    };
    return mapping[status] || "badge-neutral";
  }

  function statusBadge(order) {
    const label =
      order.status === D.ORDER_STATUS.REFUND_REQUESTED
        ? "退款审核中"
        : D.STATUS_LABEL[order.status] || order.status;
    return `<span class="badge ${badgeClass(order.status)}">${escapeHtml(label)}</span>`;
  }

  function productImage(productId) {
    const product = D.getProduct(productId);
    return product ? product.image : "/shop/assets/product-a-visual.png";
  }

  function orderStatusNote(order) {
    const refund = D.getRefundForOrder(order.id);
    const shipment = D.getShipment(order.id);
    if (order.status === D.ORDER_STATUS.PENDING) {
      return {
        className: "",
        text: `订单已占用库存，请在 ${remainingParts(order.expireAt).text} 内完成付款。`
      };
    }
    if (order.status === D.ORDER_STATUS.PAID) {
      return { className: "success", text: "支付已确认。仓库可查询并登记发货，售前退款申请会立即冻结发货。" };
    }
    if (order.status === D.ORDER_STATUS.SHIPPED) {
      return {
        className: "success",
        text: shipment
          ? `${shipment.carrier} ${shipment.trackingNo} · ${formatTime(shipment.shippedAt)} 已出库`
          : "仓库已登记发货。"
      };
    }
    if (order.status === D.ORDER_STATUS.REFUND_REQUESTED) {
      return {
        className: "refund",
        text: refund
          ? "退款申请已提交，仓库暂停发货；当前等待客服审核或渠道处理。"
          : "退款申请已提交，仓库暂停发货。"
      };
    }
    if (order.status === D.ORDER_STATUS.REFUND_FAILED) {
      return {
        className: "danger",
        text: "渠道明确退款失败，记录保留待客服重试；库存尚未恢复，仓库不可发货。"
      };
    }
    if (order.status === D.ORDER_STATUS.REFUNDED) {
      const restored = refund && refund.stockRestored ? "库存已恢复一次" : "异常支付退款不恢复库存";
      return { className: "refund", text: `退款已原路退回成功；${restored}。` };
    }
    return {
      className: "",
      text: order.closeReason || "订单已结束，历史交易记录仍可查询。"
    };
  }

  function orderActions(order) {
    const actions = [];
    if (order.status === D.ORDER_STATUS.PENDING) {
      actions.push(
        `<button class="btn btn-secondary btn-sm" type="button" data-action="cancel-order" data-order-id="${order.id}">取消订单</button>`,
        `<button class="btn btn-primary btn-sm" type="button" data-action="pay-order" data-order-id="${order.id}">${icon("card")}立即支付</button>`
      );
    } else if (order.status === D.ORDER_STATUS.PAID) {
      actions.push(
        `<button class="btn btn-secondary btn-sm" type="button" data-action="apply-refund" data-order-id="${order.id}">${icon("refund")}申请退款</button>`
      );
    } else if (order.status === D.ORDER_STATUS.SHIPPED) {
      actions.push(
        `<button class="btn btn-secondary btn-sm" type="button" data-action="view-tracking" data-order-id="${order.id}">${icon("truck")}查看运单</button>`
      );
    } else if (order.status === D.ORDER_STATUS.CLOSED || order.status === D.ORDER_STATUS.CANCELLED) {
      actions.push(
        `<button class="btn btn-secondary btn-sm" type="button" data-action="buy-again" data-product-id="${order.productId}">${icon("rotate")}重新购买</button>`
      );
    }
    actions.push(
      `<button class="btn btn-ghost btn-sm" type="button" data-action="view-order" data-order-id="${order.id}">订单详情${icon("chevron-right")}</button>`
    );
    return actions.join("");
  }

  function renderIdentity() {
    const select = $("#buyer-identity");
    const buyers = D.buyers();
    if (!buyers.length) {
      select.innerHTML = `<option value="">加载买家账号…</option>`;
      return;
    }
    if (!buyers.some((buyer) => buyer.id === selectedBuyerId)) selectedBuyerId = buyers[0].id;
    select.innerHTML = buyers
      .map(
        (buyer) =>
          `<option value="${buyer.id}" ${buyer.id === selectedBuyerId ? "selected" : ""}>${escapeHtml(
            buyer.name
          )} · ${escapeHtml(buyer.phone)}</option>`
      )
      .join("");
  }

  function renderProducts() {
    const products = D.products();
    if (!products.length) {
      $("#product-grid").innerHTML = `<p class="empty-state">暂无在售商品，请确认服务端种子数据已初始化。</p>`;
      return;
    }
    const hero = D.productMetrics(products[0].id);
    $("#hero-stock").textContent = hero.available;
    $("#hero-occupied").textContent = hero.occupied;

    $("#product-grid").innerHTML = products
      .map((product) => {
        const metrics = D.productMetrics(product.id);
        const qty = quantities[product.id] || 1;
        const disabled = !product.onSale || metrics.available < qty;
        const stockText = product.onSale
          ? metrics.available > 0
            ? `仅剩 ${metrics.available} 件可售`
            : "当前已售罄"
          : "商品已下架";
        const buyLabel = !product.onSale
          ? "暂不可购买"
          : metrics.available < qty
            ? "库存不足"
            : "立即购买";
        return `
          <article class="product-card" data-product-card="${product.id}">
            <div class="product-visual">
              <img src="${escapeHtml(product.image)}" alt="${escapeHtml(product.name)}" />
              <span class="badge no-dot product-stock">${icon("package")}${escapeHtml(stockText)}</span>
              <span class="product-live">${metrics.occupied} 件交易占用中</span>
            </div>
            <div class="product-body">
              <h3>${escapeHtml(product.name)}</h3>
              <p class="product-subtitle">${escapeHtml(product.subtitle)} · ${escapeHtml(product.spec)}</p>
              <div class="product-inventory" aria-label="${escapeHtml(product.name)}实时库存">
                <span><small>可售</small><strong>${metrics.available}</strong></span>
                <span><small>订单占用</small><strong>${metrics.occupied}</strong></span>
                <span><small>已发货</small><strong>${metrics.shipped}</strong></span>
              </div>
              <div class="product-price-row">
                <div class="product-price">${money(product.priceCents)}<small>${product.priceCents} 分</small></div>
                <div class="qty-stepper" aria-label="${escapeHtml(product.name)}购买数量">
                  <button type="button" data-qty="minus" data-product-id="${product.id}" aria-label="减少数量" ${
                    qty <= 1 ? "disabled" : ""
                  }>−</button>
                  <span>${qty}</span>
                  <button type="button" data-qty="plus" data-product-id="${product.id}" aria-label="增加数量" ${
                    qty >= 2 ? "disabled" : ""
                  }>+</button>
                </div>
              </div>
              <button class="btn btn-primary btn-block product-buy" type="button" data-buy-product="${
                product.id
              }" ${disabled ? "disabled" : ""}>
                ${icon("bag")}${buyLabel}${disabled ? "" : ` · ${money(product.priceCents * qty)}`}
              </button>
            </div>
          </article>
        `;
      })
      .join("");
  }

  function renderMiniKpis(orders) {
    const awaitingShipment = orders.filter((order) => order.status === D.ORDER_STATUS.PAID).length;
    const refundOrders = orders.filter((order) =>
      [
        D.ORDER_STATUS.REFUND_REQUESTED,
        D.ORDER_STATUS.REFUND_FAILED,
        D.ORDER_STATUS.REFUNDED
      ].includes(order.status)
    ).length;
    $("#buyer-mini-kpis").innerHTML = `
      <div class="mini-kpi"><strong>${orders.length}</strong><span>全部订单</span></div>
      <div class="mini-kpi"><strong>${awaitingShipment}</strong><span>等待发货</span></div>
      <div class="mini-kpi"><strong>${refundOrders}</strong><span>退款记录</span></div>
    `;
  }

  function filteredOrders(orders) {
    const mapping = {
      all: () => true,
      pending: (order) => order.status === D.ORDER_STATUS.PENDING,
      paid: (order) => order.status === D.ORDER_STATUS.PAID,
      refund: (order) =>
        [
          D.ORDER_STATUS.REFUND_REQUESTED,
          D.ORDER_STATUS.REFUND_FAILED,
          D.ORDER_STATUS.REFUNDED
        ].includes(order.status),
      history: (order) =>
        [
          D.ORDER_STATUS.SHIPPED,
          D.ORDER_STATUS.CANCELLED,
          D.ORDER_STATUS.CLOSED
        ].includes(order.status)
    };
    return orders.filter(mapping[activeOrderFilter] || mapping.all);
  }

  function renderOrderTabs(orders) {
    const tabs = [
      { key: "all", label: "全部", count: orders.length },
      {
        key: "pending",
        label: "待付款",
        count: orders.filter((order) => order.status === D.ORDER_STATUS.PENDING).length
      },
      {
        key: "paid",
        label: "待发货",
        count: orders.filter((order) => order.status === D.ORDER_STATUS.PAID).length
      },
      {
        key: "refund",
        label: "退款",
        count: orders.filter((order) =>
          [
            D.ORDER_STATUS.REFUND_REQUESTED,
            D.ORDER_STATUS.REFUND_FAILED,
            D.ORDER_STATUS.REFUNDED
          ].includes(order.status)
        ).length
      },
      {
        key: "history",
        label: "历史",
        count: orders.filter((order) =>
          [D.ORDER_STATUS.SHIPPED, D.ORDER_STATUS.CANCELLED, D.ORDER_STATUS.CLOSED].includes(
            order.status
          )
        ).length
      }
    ];
    $("#order-tabs").innerHTML = tabs
      .map(
        (tab) =>
          `<button type="button" role="tab" aria-selected="${tab.key === activeOrderFilter}" class="${
            tab.key === activeOrderFilter ? "active" : ""
          }" data-order-filter="${tab.key}">${tab.label} ${tab.count}</button>`
      )
      .join("");
  }

  function renderOrders() {
    const buyer = getSelectedBuyer();
    if (buyer) $("#buyer-order-subtitle").textContent = `${buyer.name} · ${buyer.phone}`;
    const orders = buyerOrders();
    renderMiniKpis(orders);
    renderOrderTabs(orders);
    const visible = filteredOrders(orders);
    const list = $("#order-list");
    if (!visible.length) {
      list.innerHTML = `
        <div class="empty-state">
          ${icon("receipt")}
          <strong>${orders.length ? "当前筛选下没有订单" : "还没有订单"}</strong>
          <p>${orders.length ? "切换筛选条件查看其他交易记录。" : "从左侧选择商品下单，订单和支付结果会在这里持续更新。"}</p>
        </div>
      `;
      return;
    }
    list.innerHTML = visible
      .map((order) => {
        const note = orderStatusNote(order);
        const countdown =
          order.status === D.ORDER_STATUS.PENDING
            ? `<div class="countdown-bar">
                <span>${icon("clock")} 支付剩余</span>
                <strong data-countdown-order="${order.id}" data-expire="${order.expireAt}">${remainingParts(
                  order.expireAt
                ).text}</strong>
              </div>`
            : "";
        return `
          <article class="order-card" data-tone="${orderTone(order)}" data-order-card="${order.id}">
            <div class="order-card-head">
              <div>
                <div class="order-id">${escapeHtml(order.id)} ${statusBadge(order)}</div>
                <time datetime="${escapeHtml(order.createdAt)}">${formatTime(order.createdAt, true)}</time>
              </div>
              <span class="mono muted">${escapeHtml(order.paymentTxnNo || "等待渠道流水")}</span>
            </div>
            ${countdown}
            <div class="order-product">
              <img class="order-thumb" src="${escapeHtml(productImage(order.productId))}" alt="${escapeHtml(
                order.productName
              )}" />
              <div class="order-product-name">
                <strong>${escapeHtml(order.productName)}</strong>
                <span>${escapeHtml(order.productSpec)} · ${money(order.unitPriceCents)} / 件</span>
              </div>
              <div class="order-amount">${money(order.amountCents)}<span>× ${order.qty}</span></div>
            </div>
            <div class="order-status-note ${note.className}">${escapeHtml(note.text)}</div>
            <div class="order-actions">${orderActions(order)}</div>
          </article>
        `;
      })
      .join("");
  }

  function render() {
    renderIdentity();
    renderProducts();
    renderOrders();
    if (drawerOpen) renderDrawer();
    document.body.classList.toggle("buyer-view-market", activeMobileView === "market");
    document.body.classList.toggle("buyer-view-orders", activeMobileView === "orders");
    $$(".buyer-mobile-nav button[data-view]").forEach((button) => {
      button.classList.toggle("active", button.dataset.view === activeMobileView);
    });
  }

  function toast(title, message, tone) {
    const node = document.createElement("div");
    node.className = `toast ${tone || ""}`;
    node.innerHTML = `
      ${icon(tone === "error" ? "warning" : tone === "success" ? "check" : "info")}
      <div><strong>${escapeHtml(title)}</strong>${escapeHtml(message || "")}</div>
    `;
    $("#toast-root").appendChild(node);
    window.setTimeout(() => {
      node.style.opacity = "0";
      node.style.transform = "translateY(-5px)";
      window.setTimeout(() => node.remove(), 180);
    }, 4200);
  }

  function openModal(options) {
    const root = $("#modal-root");
    modalOpen = true;
    activeOrderDetailId = options.orderDetailId || null;
    document.body.style.overflow = "hidden";
    root.innerHTML = `
      <div class="modal-backdrop" data-modal-backdrop>
        <section class="modal ${options.wide ? "modal-wide" : ""} ${
          options.compact ? "modal-compact" : ""
        }" role="dialog" aria-modal="true" aria-labelledby="modal-title">
          <header class="modal-head">
            <div>
              <h2 id="modal-title">${escapeHtml(options.title)}</h2>
              ${options.subtitle ? `<p>${escapeHtml(options.subtitle)}</p>` : ""}
            </div>
            <button class="icon-btn" type="button" data-close-modal aria-label="关闭弹层">
              ${icon("x")}
            </button>
          </header>
          <div class="modal-body">${options.content || ""}</div>
          ${
            options.footer
              ? `<footer class="modal-foot">${options.footer}</footer>`
              : ""
          }
        </section>
      </div>
    `;
    window.requestAnimationFrame(() => {
      const focusTarget =
        root.querySelector("input:not([type=hidden]), select, textarea") ||
        root.querySelector("button");
      if (focusTarget) focusTarget.focus({ preventScroll: true });
    });
  }

  function closeModal() {
    modalOpen = false;
    activeOrderDetailId = null;
    $("#modal-root").innerHTML = "";
    document.body.style.overflow = "";
  }

  function openPaymentMethodModal(orderId) {
    const order = D.getOrder(orderId);
    if (!order) {
      toast("订单不存在", "请刷新订单后重试。", "error");
      return;
    }
    if (order.status !== D.ORDER_STATUS.PENDING) {
      toast("订单当前不可支付", `订单状态为 ${D.STATUS_LABEL[order.status]}。`, "warning");
      render();
      return;
    }
    openModal({
      title: "选择支付方式",
      subtitle: `${order.id} · ${order.productName} × ${order.qty}`,
      content: `
        <div class="payment-summary">
          <div class="payment-summary-top">
            <span>应付金额</span>
            <strong>${money(order.amountCents)}</strong>
          </div>
          <p>订单库存已占用，支付有效期至 ${formatTime(order.expireAt, true)}。</p>
        </div>
        <div class="payment-method-list">
          <button class="payment-method-card" type="button" data-action="choose-payment-method" data-method="wechat" data-order-id="${
            order.id
          }">
            <span class="payment-brand wechat">微</span>
            <span class="payment-method-copy">
              <strong>微信支付（模拟渠道）</strong>
              <small>进入收银台页，可演示正常 / 重复 / 迟到 / 网络中断通知</small>
            </span>
            ${icon("chevron-right")}
          </button>
          <button class="payment-method-card" type="button" data-action="choose-payment-method" data-method="alipay" data-order-id="${
            order.id
          }">
            <span class="payment-brand alipay">支</span>
            <span class="payment-method-copy">
              <strong>支付宝（真实渠道）</strong>
              <small>跳转支付宝沙箱收银页，回调由渠道异步通知</small>
            </span>
            ${icon("chevron-right")}
          </button>
        </div>
        <div class="payment-channel-note">
          ${icon("shield")}
          <span>支付结果以渠道返回为准；未完成支付前订单保持待付款。</span>
        </div>
      `,
      footer: `<button class="btn btn-secondary" type="button" data-close-modal>暂不支付</button>`
    });
  }

  function openChannelPaymentPage(orderId, method) {
    const order = D.getOrder(orderId);
    if (!order) {
      toast("订单不存在", "请刷新订单后重试。", "error");
      return;
    }
    if (order.status !== D.ORDER_STATUS.PENDING) {
      toast("订单当前不可支付", `订单状态为 ${D.STATUS_LABEL[order.status]}。`, "warning");
      closeModal();
      render();
      return;
    }
    if (method === "alipay") {
      // 支付宝真实渠道：先请求支付链接，再跳转沙箱收银页
      D.alipayPay(orderId, selectedBuyerId).then((result) => {
        if (!result.ok) {
          toast("支付宝下单失败", result.message || "请检查沙箱配置。", "error");
          return;
        }
        if (result.payURL) {
          window.open(result.payURL, "_blank", "noopener");
          toast("已跳转支付宝", "支付结果由支付宝异步通知，返回本页即可查看订单状态。", "success");
        } else {
          toast("未获取到支付链接", "请检查支付宝沙箱配置，可改用微信模拟渠道。", "warning");
        }
      });
      return;
    }
    const params = new URLSearchParams({
      orderId: order.id,
      buyerId: selectedBuyerId,
      method: "wechat"
    });
    window.location.href = `/shop/payment?${params.toString()}`;
  }

  function openOrderConfirmModal(productId, qty) {
    const product = D.getProduct(productId);
    if (!product) {
      toast("商品不存在", "请刷新商品列表后重试。", "error");
      return;
    }
    const metrics = D.productMetrics(product.id);
    if (!product.onSale || metrics.available < qty) {
      toast(
        "当前无法下单",
        !product.onSale ? "商品已下架。" : `库存仅剩 ${metrics.available} 件，请调整购买数量。`,
        "warning"
      );
      render();
      return;
    }
    const buyer = getSelectedBuyer();
    const amountCents = product.priceCents * qty;
    openModal({
      title: "确认订单",
      subtitle: "请核对商品、数量与应付金额，确认后才会占用库存",
      content: `
        <div class="order-confirm-product">
          <img src="${escapeHtml(product.image)}" alt="${escapeHtml(product.name)}" />
          <div class="order-confirm-product-copy">
            <strong>${escapeHtml(product.name)}</strong>
            <span>${escapeHtml(product.spec)} · 标准商品</span>
            <span>购买数量 × ${qty}</span>
          </div>
          <strong class="order-confirm-price">${money(amountCents)}</strong>
        </div>
        <dl class="order-confirm-breakdown">
          <div><dt>商品金额（${qty} 件）</dt><dd>${money(amountCents)}</dd></div>
          <div><dt>运费 / 优惠</dt><dd>免运费 / 无优惠</dd></div>
          <div class="order-confirm-total"><dt>应付金额</dt><dd>${money(amountCents)}</dd></div>
        </dl>
        <div class="order-confirm-meta">
          <div><span>下单账号</span><strong>${escapeHtml(buyer.name)} · ${escapeHtml(buyer.phone)}</strong></div>
          <div><span>当前可售</span><strong>${metrics.available} 件</strong></div>
        </div>
        <div class="notice-box warning">
          ${icon("clock")}
          <div>确认后将占用 ${qty} 件库存，订单保留 10 分钟。取消、超时关闭或退款成功后再释放；重复关闭不会重复释放。</div>
        </div>
      `,
      footer: `
        <button class="btn btn-secondary" type="button" data-close-modal>返回选购</button>
        <button class="btn btn-primary" type="button" data-action="confirm-create-order" data-product-id="${product.id}" data-order-qty="${qty}">
          ${icon("check")}确认提交订单
        </button>
      `
    });
  }

  async function createConfirmedOrder(productId, qty) {
    const response = await D.createOrder(selectedBuyerId, productId, qty);
    if (!response.ok) {
      closeModal();
      toast("无法创建订单", response.message, "error");
      render();
      return;
    }
    closeModal();
    toast("订单已提交", `${response.data.order.id} 已占用库存，请选择支付方式。`, "success");
    activeOrderFilter = "pending";
    render();
    openPaymentMethodModal(response.data.order.id);
  }

  function openCancelModal(orderId) {
    const order = D.getOrder(orderId);
    if (!order) return;
    openModal({
      title: "确认取消订单",
      compact: true,
      subtitle: `${order.id} · ${money(order.amountCents)}`,
      content: `
        <div class="notice-box warning">
          ${icon("warning")}
          <div>取消成功后会立即释放本单占用的 ${order.qty} 件库存。重复关闭不会重复释放。</div>
        </div>
      `,
      footer: `
        <button class="btn btn-secondary" type="button" data-close-modal>继续支付</button>
        <button class="btn btn-danger" type="button" data-action="confirm-cancel-order" data-order-id="${order.id}">确认取消</button>
      `
    });
  }

  async function cancelOrder(orderId) {
    const response = await D.cancelOrder(orderId, selectedBuyerId);
    closeModal();
    if (response.ok) {
      toast(
        "订单已取消",
        response.code === "duplicated" ? "订单此前已取消，未重复释放库存。" : "占用库存已释放。",
        "success"
      );
      activeOrderFilter = "all";
      render();
    } else {
      toast("无法取消订单", response.message, "error");
    }
  }

  function openRefundModal(orderId) {
    const order = D.getOrder(orderId);
    if (!order) return;
    if (order.status !== D.ORDER_STATUS.PAID) {
      toast("当前不可申请退款", "售前退款仅适用于已付款且未发货订单。", "warning");
      render();
      return;
    }
    openModal({
      title: "申请售前整单退款",
      subtitle: `${order.id} · 退款后仓库将立即暂停发货`,
      content: `
        <div class="refund-hero">
          <span>原路退回金额</span>
          <strong>${money(order.amountCents)}</strong>
        </div>
        <div class="notice-box warning" style="margin-bottom: 16px">
          ${icon("warning")}
          <div>申请提交后仓库端口令会转为冻结。客服拒绝后可恢复发货；同意后由模拟渠道退款。</div>
        </div>
        <div class="form-grid">
          <div class="choice-list" role="radiogroup" aria-label="退款原因">
            <label class="choice"><input type="radio" name="refund-reason" value="规格或型号选错" checked />规格或型号选错</label>
            <label class="choice"><input type="radio" name="refund-reason" value="计划有变，不再需要" />计划有变，不再需要</label>
            <label class="choice"><input type="radio" name="refund-reason" value="信息填错，准备重新下单" />信息填错，准备重新下单</label>
          </div>
          <div class="field">
            <label for="refund-note">补充说明</label>
            <textarea id="refund-note" placeholder="可补充退款原因，客服审核时会看到此内容"></textarea>
            <span class="field-hint">整单退款，不支持部分退款。退款期间不可发货。</span>
          </div>
        </div>
      `,
      footer: `
        <button class="btn btn-secondary" type="button" data-close-modal>暂不申请</button>
        <button class="btn btn-danger" type="button" data-action="confirm-refund" data-order-id="${order.id}">确认申请退款</button>
      `
    });
  }

  async function confirmRefund(orderId) {
    const selected = $('input[name="refund-reason"]:checked')?.value || "买家申请售前退款";
    const note = $("#refund-note")?.value.trim();
    const response = await D.applyRefund(orderId, selectedBuyerId, note ? `${selected}；${note}` : selected);
    closeModal();
    if (response.ok) {
      toast("退款申请已提交", "仓库发货已冻结，等待客服审核和渠道退款结果。", "success");
      activeOrderFilter = "refund";
      render();
    } else {
      toast("退款申请未通过", response.message, "error");
      render();
    }
  }

  function buildOrderTimeline(order) {
    const events = [];
    events.push({ at: order.createdAt, title: "订单创建", text: `占用 ${order.qty} 件库存，金额快照 ${money(order.amountCents)}` });
    if (order.paidAt) {
      events.push({
        at: order.paidAt,
        title: "支付成功",
        text: `${order.paymentTxnNo || "渠道交易号"} · ${money(order.amountCents)}`
      });
    }
    const refund = D.getRefundForOrder(order.id);
    if (refund) {
      events.push({
        at: refund.requestedAt,
        title: refund.source === "PRESALE" ? "售前退款申请" : "异常支付退款",
        text: refund.requestReason
      });
      (refund.reviewHistory || []).forEach((item) => {
        events.push({
          at: item.at,
          title:
            item.decision === "REJECTED"
              ? "客服拒绝退款"
              : item.decision === "APPROVED"
                ? "客服同意退款"
                : "退款处理记录",
          text: item.reason || ""
        });
      });
      if (refund.channelStatus === "SUCCESS") {
        events.push({ at: "", title: "渠道退款成功", text: `退款金额 ${money(refund.refundAmountCents)}` });
      } else if (refund.channelStatus === "FAILED") {
        events.push({
          at: "",
          title: "渠道退款失败",
          text: `已尝试 ${refund.attempts} 次 · ${refund.lastError || "待客服重试"}`
        });
      }
    }
    const shipment = D.getShipment(order.id);
    if (shipment) {
      events.push({
        at: shipment.shippedAt,
        title: "仓库登记发货",
        text: `${shipment.carrier} · ${shipment.trackingNo}`
      });
    }
    if (order.closedAt) {
      events.push({ at: order.closedAt, title: D.STATUS_LABEL[order.status], text: order.closeReason });
    }
    return events
      .filter((item) => item.at)
      .sort((a, b) => new Date(a.at) - new Date(b.at))
      .map(
        (item) => `
          <div class="timeline-item">
            <strong>${escapeHtml(item.title)}</strong>
            <p>${escapeHtml(item.text)}</p>
            <time>${formatTime(item.at, true)}</time>
          </div>
        `
      )
      .join("");
  }

  function openOrderDetail(orderId) {
    const order = D.getOrder(orderId);
    if (!order || order.buyerId !== selectedBuyerId) {
      toast("无法查看订单", "订单不存在，或不属于当前买家。", "error");
      return;
    }
    const payment = D.getState().payments.find(
      (item) => item.orderId === order.id && item.status === "ACCEPTED"
    );
    const refund = D.getRefundForOrder(order.id);
    const shipment = D.getShipment(order.id);
    openModal({
      title: `订单 ${order.id}`,
      subtitle: "状态、收款、库存、退款与发货使用同一交易快照",
      wide: true,
      orderDetailId: order.id,
      content: `
        <div style="margin-bottom: 18px">${statusBadge(order)}</div>
        <dl class="detail-list">
          <div class="detail-row"><dt>商品快照</dt><dd>${escapeHtml(order.productName)} · ${
            escapeHtml(order.productSpec)
          } × ${order.qty}</dd></div>
          <div class="detail-row"><dt>订单金额</dt><dd><strong>${money(
            order.amountCents
          )}</strong>（${order.unitPriceCents} 分/件 × ${order.qty}）</dd></div>
          <div class="detail-row"><dt>渠道交易号</dt><dd class="mono">${escapeHtml(
            order.paymentTxnNo || (payment && payment.txnNo) || "尚未形成有效收款"
          )}</dd></div>
          <div class="detail-row"><dt>退款单</dt><dd>${
            refund
              ? `${refund.id} · ${refund.source === "PRESALE" ? "售前整单退款" : "异常支付退款"} · 渠道${
                  refund.channelStatus === "SUCCESS"
                    ? "成功"
                    : refund.channelStatus === "FAILED"
                      ? "失败待重试"
                      : "处理中"
                }`
              : "无"
          }</dd></div>
          <div class="detail-row"><dt>库存处理</dt><dd>${
            refund
              ? refund.stockRestored
                ? "售前退款成功，库存已恢复"
                : "异常支付退款，不恢复库存"
              : order.status === "CANCELLED" || order.status === "CLOSED"
                ? "订单关闭，占用已释放"
                : "占用中"
          }</dd></div>
          <div class="detail-row"><dt>发货信息</dt><dd>${
            shipment
              ? `${escapeHtml(shipment.carrier)} / ${escapeHtml(shipment.trackingNo)} · ${formatTime(
                  shipment.shippedAt,
                  true
                )}`
              : "尚未发货"
          }</dd></div>
          ${
            order.closeReason
              ? `<div class="detail-row"><dt>关闭原因</dt><dd>${escapeHtml(order.closeReason)}</dd></div>`
              : ""
          }
        </dl>
        <div class="divider" style="margin: 18px 0"></div>
        <h3 style="margin-bottom: 14px; font-size: 14px">业务时间线</h3>
        <div class="timeline">${buildOrderTimeline(order)}</div>
      `,
      footer: `<button class="btn btn-primary" type="button" data-close-modal>关闭</button>`
    });
  }

  function openTrackingModal(orderId) {
    const order = D.getOrder(orderId);
    const shipment = D.getShipment(orderId);
    if (!order || order.buyerId !== selectedBuyerId || !shipment) {
      toast("暂无运单", "该订单还没有有效的发货登记。", "warning");
      return;
    }
    openModal({
      title: "物流运单",
      compact: true,
      content: `
        <div class="tracking-box">
          <div class="tracking-box-head">
            <div class="tracking-icon">${icon("truck")}</div>
            <div>
              <strong>${escapeHtml(shipment.carrier)}</strong>
              <span>已于 ${formatTime(shipment.shippedAt, true)} 登记出库</span>
            </div>
          </div>
          <div class="tracking-number">
            <span class="mono">${escapeHtml(shipment.trackingNo)}</span>
            <button class="btn btn-secondary btn-sm" type="button" data-copy="${escapeHtml(
              shipment.trackingNo
            )}">${icon("copy")}复制单号</button>
          </div>
        </div>
        <div class="notice-box">${icon(
          "info"
        )}<div>本题仅登记承运方与运单号，不对接真实物流轨迹。</div></div>
      `,
      footer: `<button class="btn btn-primary" type="button" data-close-modal>知道了</button>`
    });
  }

  async function copyText(text) {
    try {
      await navigator.clipboard.writeText(text);
      toast("已复制", text, "success");
    } catch (error) {
      const textarea = document.createElement("textarea");
      textarea.value = text;
      textarea.style.position = "fixed";
      textarea.style.opacity = "0";
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand("copy");
      textarea.remove();
      toast("已复制", text, "success");
    }
  }

  function renderDrawer() {
    const root = $("#drawer-root");
    if (!drawerOpen) {
      root.innerHTML = "";
      return;
    }
    const scenarios = [
      { id: "normal", title: "正常成交", text: "下单、付款、发货闭环" },
      { id: "stock", title: "20 并发抢库存", text: "超卖拦截与可售归零" },
      { id: "refund", title: "售前退款", text: "冻结发货、失败重试、库存回补" },
      { id: "late", title: "迟到支付", text: "关闭后收款转客服退款，不恢复库存" },
      { id: "network", title: "网络异常", text: "重复通知幂等去重" },
      { id: "race", title: "操作竞争", text: "退款与发货只能一方成功" },
      { id: "auth", title: "越权访问", text: "跨买家、伪造签名、金额不符全拦截" }
    ];
    root.innerHTML = `
      <div class="drawer-scrim" data-action="close-demo"></div>
      <aside class="demo-drawer" aria-label="演示场景">
        <header class="drawer-head">
          <div>
            <h2>验收场景台</h2>
            <p>一键场景会先重置服务端演示数据，再通过真实接口执行跨角色流程。</p>
          </div>
          <button class="icon-btn" type="button" data-action="close-demo" aria-label="关闭演示台">${icon(
            "x"
          )}</button>
        </header>
        <div class="drawer-body">
          <section class="drawer-section">
            <div class="drawer-section-head"><h3>完整业务场景</h3><span class="badge badge-paid no-dot">7 组</span></div>
            <div class="scenario-grid">
              ${scenarios
                .map(
                  (scenario) => `
                    <button class="scenario-button" type="button" data-action="run-scenario" data-scenario="${scenario.id}">
                      <strong>${scenario.title}</strong>
                      <span>${scenario.text}</span>
                    </button>
                  `
                )
                .join("")}
            </div>
            <p class="scenario-note">场景会重置后端演示数据并重新登录演示账号；完成后可在本页订单及 GVA 管理端继续查看。</p>
          </section>
          <section class="drawer-section">
            <div class="drawer-section-head"><h3>超时与异常通知</h3></div>
            <div class="channel-test-row">
              <button class="btn btn-secondary" type="button" data-action="advance-time">${icon(
                "clock"
              )}推进服务端时间 11 分钟并关闭待付单</button>
              <button class="btn btn-secondary" type="button" data-action="test-notice" data-mode="unknown">${icon(
                "search"
              )}发送未知订单成功通知</button>
              <button class="btn btn-secondary" type="button" data-action="test-notice" data-mode="amount">${icon(
                "warning"
              )}发送金额不符通知</button>
              <button class="btn btn-secondary" type="button" data-action="test-notice" data-mode="signature">${icon(
                "shield"
              )}发送验签失败通知</button>
            </div>
          </section>
          <section class="drawer-section">
            <div class="drawer-section-head"><h3>数据一致性</h3></div>
            <div id="buyer-invariant-summary"></div>
            <button class="btn btn-primary btn-block" type="button" data-action="validate-invariants" style="margin-top: 9px">
              ${icon("shield")}校验库存、支付、退款与发货约束
            </button>
          </section>
        </div>
      </aside>
    `;
    renderBuyerInvariantSummary();
  }

  function renderBuyerInvariantSummary() {
    const summary = D.financialSummary();
    const inventory = D.inventoryInvariant();
    const target = $("#buyer-invariant-summary");
    if (!target) return;
    target.innerHTML = `
      <div class="notice-box ${inventory.ok ? "success" : "danger"}">
        ${icon(inventory.ok ? "check" : "warning")}
        <div>
          库存守恒${inventory.ok ? "正常" : "异常"}；有效收款 ${money(
            summary.validCents
          )}，已退款 ${money(summary.refundCents)}，异常待处理 ${summary.abnormalCount} 笔。
        </div>
      </div>
    `;
  }

  function openDemo() {
    drawerOpen = true;
    renderDrawer();
  }

  function closeDemo() {
    drawerOpen = false;
    renderDrawer();
  }

  async function runScenario(name) {
    const result = await D.runScenario(name);
    if (!result.ok) {
      closeDemo();
      openModal({
        title: "场景存在未通过断言",
        subtitle: "请检查以下步骤输出",
        content: `
          <div class="notice-box danger">
            ${icon("warning")}
            <div>场景执行未完全通过，请查看步骤明细。</div>
          </div>
          <div class="timeline">${renderScenarioSteps(result)}</div>
        `,
        footer: `<button class="btn btn-primary" type="button" data-close-modal>关闭</button>`
      });
      return;
    }
    activeOrderFilter = "all";
    activeMobileView = "orders";
    closeDemo();
    render();
    openModal({
      title: `场景「${scenarioTitle(name)}」执行完成`,
      subtitle: "验收场景已通过真实接口从固定商品与账号重新执行",
      content: `
        <div class="notice-box success">
          ${icon("check")}
          <div>全部步骤断言通过；订单、库存流水与通知留痕可在本页与 GVA 管理端核查。</div>
        </div>
        <div class="timeline">${renderScenarioSteps(result)}</div>
        <p class="field-hint" style="margin: 13px 0 0">可切换到员工端继续查看仓库、客服和运营结果。</p>
      `,
      footer: `
        <button class="btn btn-secondary" type="button" data-close-modal>留在买家端</button>
        <a class="btn btn-primary" href="http://127.0.0.1:8099" target="_blank" rel="noreferrer">打开 GVA 管理端</a>
      `
    });
  }

  function scenarioTitle(name) {
    const mapping = {
      normal: "正常成交",
      stock: "防超卖",
      refund: "售前退款",
      late: "迟到支付",
      network: "网络与重复通知",
      race: "退款发货互斥",
      auth: "越权拦截"
    };
    return mapping[name] || name;
  }

  function renderScenarioSteps(result) {
    if (!result || !result.steps) return "";
    return result.steps
      .map(
        (step) => `
          <div class="timeline-item">
            <strong>${icon(step.tone === "bad" ? "warning" : "check")} ${escapeHtml(step.text)}</strong>
          </div>
        `
      )
      .join("");
  }

  async function runNoticeTest(mode) {
    const state = D.getState();
    const pending = state.orders.find((order) => order.status === D.ORDER_STATUS.PENDING);
    let payload;
    if (mode === "unknown") {
      payload = {
        orderId: "ORD-VOID-9999",
        txnNo: `WX_UNKNOWN_${Date.now()}`,
        amountCents: 10000,
        signatureOk: true
      };
    } else {
      if (!pending) {
        toast("没有待付款订单", "请先创建一笔待付款订单，再测试金额或验签异常。", "warning");
        return;
      }
      payload =
        mode === "amount"
          ? {
              orderId: pending.id,
              txnNo: `WX_AMOUNT_BAD_${Date.now()}`,
              amountCents: 1,
              signatureOk: true
            }
          : {
              orderId: pending.id,
              txnNo: `WX_SIGN_BAD_${Date.now()}`,
              amountCents: pending.amountCents,
              signatureOk: false
            };
    }
    const response = await D.sendTestNotification(payload);
    closeDemo();
    openModal({
      title: "渠道通知已留痕",
      compact: true,
      content: `
        <div class="notice-box ${response.ok ? "success" : "danger"}">
          ${icon(response.ok ? "check" : "warning")}
          <div>${escapeHtml(response.message)}。该结果可在 GVA 管理端运营核查与渠道通知记录中查询。</div>
        </div>
      `,
      footer: `<button class="btn btn-primary" type="button" data-close-modal>知道了</button>`
    });
  }

  async function advanceTime() {
    const pendingCount = D.getState().orders.filter(
      (order) => order.status === D.ORDER_STATUS.PENDING
    ).length;
    const response = await D.advanceClock(11 * minute);
    closeDemo();
    render();
    toast(
      "服务端时间已推进",
      response.expired && response.expired.length
        ? `${response.expired.length} 笔待付款订单已超时关闭并释放占用。`
        : pendingCount
          ? "待付款订单已由服务端超时关闭。"
          : "没有待付款订单需要关闭。",
      "warning"
    );
  }

  function validateInvariants() {
    const validation = D.validateInvariants();
    const inventory = D.inventoryInvariant();
    openModal({
      title: validation.ok ? "业务约束全部通过" : "发现业务约束异常",
      subtitle: "检查库存守恒、支付幂等、退款幂等、退款发货互斥及异常退款库存",
      content: validation.ok
        ? `
          <div class="notice-box success">
            ${icon("check")}
            <div>库存守恒、有效支付唯一、发货唯一、退款与发货互斥，以及异常退款不恢复库存均通过。</div>
          </div>
          <div class="invariant-visual" style="margin-top: 14px">
            ${inventory.products
              .map((item) => {
                const metrics = D.productMetrics(item.productId);
                return `
                  <div class="invariant-product">
                    <div class="invariant-product-head">
                      <span>${escapeHtml((D.getProduct(item.productId) || {}).name || item.productId)}</span>
                      <span class="badge badge-paid">守恒</span>
                    </div>
                    <div class="invariant-equation">可售 <strong>${metrics.available}</strong> + 占用 <strong>${
                      metrics.occupied
                    }</strong> + 已发货 <strong>${metrics.shipped}</strong> = 总量 <strong>${
                      metrics.total
                    }</strong></div>
                  </div>
                `;
              })
              .join("")}
          </div>
        `
        : `
          <div class="notice-box danger">
            ${icon("warning")}
            <div>${validation.issues.map(escapeHtml).join("<br>")}</div>
          </div>
        `,
      footer: `<button class="btn btn-primary" type="button" data-close-modal>关闭</button>`
    });
  }

  function setMobileView(view) {
    activeMobileView = view;
    render();
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  async function switchBuyer(buyerId) {
    selectedBuyerId = buyerId;
    savePreference(selectedBuyerId);
    activeOrderFilter = "all";
    if (!D.tokenFor(buyerId)) {
      const result = await D.login(buyerId, D.DEMO_PASSWORD);
      if (!result.ok) {
        render();
        toast("登录失败", result.message || "请检查演示口令。", "error");
        return;
      }
      await D.refresh();
    }
    render();
    toast("已切换买家", `当前身份：${getSelectedBuyer().name}`, "success");
  }

  function bindEvents() {
    $("#buyer-identity").addEventListener("change", (event) => {
      void switchBuyer(event.target.value);
    });

    document.addEventListener("click", (event) => {
      const closeModalTrigger = event.target.closest("[data-close-modal]");
      if (closeModalTrigger) {
        closeModal();
        return;
      }
      if (event.target.matches("[data-modal-backdrop]")) {
        closeModal();
        return;
      }

      const filter = event.target.closest("[data-order-filter]");
      if (filter) {
        activeOrderFilter = filter.dataset.orderFilter;
        renderOrders();
        return;
      }

      const qtyButton = event.target.closest("[data-qty]");
      if (qtyButton) {
        const productId = qtyButton.dataset.productId;
        quantities[productId] = Math.max(
          1,
          Math.min(2, (quantities[productId] || 1) + (qtyButton.dataset.qty === "plus" ? 1 : -1))
        );
        renderProducts();
        return;
      }

      const buyButton = event.target.closest("[data-buy-product]");
      if (buyButton) {
        const productId = buyButton.dataset.buyProduct;
        openOrderConfirmModal(productId, quantities[productId] || 1);
        return;
      }

      const viewButton = event.target.closest("[data-view]");
      if (viewButton) {
        setMobileView(viewButton.dataset.view);
        return;
      }

      const copyButton = event.target.closest("[data-copy]");
      if (copyButton) {
        void copyText(copyButton.dataset.copy);
        return;
      }

      const actionButton = event.target.closest("[data-action]");
      if (!actionButton) return;
      const action = actionButton.dataset.action;
      const orderId = actionButton.dataset.orderId;
      if (action === "refresh-orders") {
        void D.refresh().then(() => {
          render();
          toast("订单已刷新", "已读取服务端最新交易状态。", "success");
        });
      } else if (action === "confirm-create-order") {
        void createConfirmedOrder(actionButton.dataset.productId, Number(actionButton.dataset.orderQty || 1));
      } else if (action === "pay-order") {
        openPaymentMethodModal(orderId);
      } else if (action === "choose-payment-method") {
        openChannelPaymentPage(orderId, actionButton.dataset.method);
      } else if (action === "cancel-order") {
        openCancelModal(orderId);
      } else if (action === "confirm-cancel-order") {
        void cancelOrder(orderId);
      } else if (action === "apply-refund") {
        openRefundModal(orderId);
      } else if (action === "confirm-refund") {
        void confirmRefund(orderId);
      } else if (action === "view-order") {
        openOrderDetail(orderId);
      } else if (action === "view-tracking") {
        openTrackingModal(orderId);
      } else if (action === "buy-again") {
        const product = D.getProduct(actionButton.dataset.productId);
        setMobileView("market");
        if (product) {
          const card = $(`[data-product-card="${product.id}"]`);
          if (card) card.scrollIntoView({ behavior: "smooth", block: "center" });
          toast("可重新购买", `${product.name} 当前可售 ${D.productMetrics(product.id).available} 件。`, "success");
        }
      } else if (action === "open-demo") {
        openDemo();
      } else if (action === "close-demo") {
        closeDemo();
      } else if (action === "run-scenario") {
        void runScenario(actionButton.dataset.scenario);
      } else if (action === "advance-time") {
        void advanceTime();
      } else if (action === "test-notice") {
        void runNoticeTest(actionButton.dataset.mode);
      } else if (action === "validate-invariants") {
        validateInvariants();
      }
    });

    document.addEventListener("keydown", (event) => {
      if (event.key === "Escape") {
        if (modalOpen) closeModal();
        else if (drawerOpen) closeDemo();
      }
    });

    D.subscribe(() => {
      const previousShipment =
        modalOpen && activeOrderDetailId ? D.getShipment(activeOrderDetailId) : null;
      render();
      if (!modalOpen || !activeOrderDetailId) return;
      const latestOrder = D.getOrder(activeOrderDetailId);
      if (!latestOrder || latestOrder.buyerId !== selectedBuyerId) {
        closeModal();
        return;
      }
      const currentShipment = D.getShipment(activeOrderDetailId);
      if (!previousShipment && currentShipment) {
        toast(
          "仓库已发货",
          `${currentShipment.carrier} ${currentShipment.trackingNo} 已写入当前订单。`,
          "success"
        );
      }
      openOrderDetail(activeOrderDetailId);
    });
  }

  function updateCountdowns() {
    $$("[data-countdown-order]").forEach((node) => {
      node.textContent = remainingParts(node.dataset.expire).text;
    });
    if (Date.now() - lastExpiryCheck < 1000) return;
    lastExpiryCheck = Date.now();
    // 服务端已内置超时关单调度；这里仅周期性同步状态
    void D.refresh();
  }

  async function init() {
    await D.refresh();
    const buyers = D.buyers();
    if (!buyers.some((buyer) => buyer.id === selectedBuyerId)) {
      selectedBuyerId = buyers.length ? buyers[0].id : "";
    }
    if (selectedBuyerId && !D.tokenFor(selectedBuyerId)) {
      await D.login(selectedBuyerId, D.DEMO_PASSWORD);
      await D.refresh();
    }
    const returnParams = new URLSearchParams(window.location.search);
    const paymentResult = returnParams.get("payment");
    const returnedOrderId = returnParams.get("order");
    if (paymentResult) {
      activeMobileView = "orders";
      activeOrderFilter = paymentResult === "success" ? "paid" : "all";
    }
    render();
    bindEvents();
    if (paymentResult) {
      const returnedOrder = returnedOrderId ? D.getOrder(returnedOrderId) : null;
      if (paymentResult === "success") {
        toast(
          "支付已确认",
          returnedOrder
            ? `${returnedOrder.id} 已进入待发货，仓库现在可以查询到该订单。`
            : "支付结果已返回，订单状态已更新。",
          "success"
        );
      } else if (paymentResult === "late") {
        toast(
          "收款到账，订单已关闭",
          "渠道完成收款时订单已超时关闭；该笔收款已进入客服异常支付待处理列表。",
          "warning"
        );
      } else {
        toast("支付未完成", "订单仍保持待付款，可返回订单重新发起支付。", "warning");
      }
      window.history.replaceState({}, "", "/shop");
    }
    window.setInterval(updateCountdowns, 1000);
    // 周期同步服务端数据（仓库/客服在 GVA 管理端的操作会反映到本页）
    window.setInterval(() => {
      if (!drawerOpen && !modalOpen) void D.refresh();
    }, 5000);
  }

  init();
})();
