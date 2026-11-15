/**
 * Shop 收银台 —— 支付处理（依据原型 payment.js 对齐，数据层换成真实后端）
 * 原型中同步的 mock 写操作改为 async/await：
 *   - normal：向 /api/pay/notify 发起带 HMAC 验签的渠道通知
 *   - duplicate：同一渠道交易号重复通知，验证幂等去重
 *   - late：先触发超时关单，再补发通知，验证异常收款路径
 *   - network：本地模拟网络中断（订单保持待付款，可重试）
 */
(function () {
  "use strict";

  const D = window.ShopDomain;
  const root = document.querySelector("#checkout-content");
  const params = new URLSearchParams(window.location.search);
  const orderId = params.get("orderId") || "";
  const buyerId = params.get("buyerId") || "";
  const method = params.get("method") === "alipay" ? "alipay" : "wechat";
  const requestedScenario = params.get("scenario");
  const allowedScenarios = new Set(["network", "duplicate", "late"]);
  const methodInfo =
    method === "alipay"
      ? { label: "支付宝", initial: "支", tone: "alipay" }
      : { label: "微信支付", initial: "微", tone: "wechat" };

  let simulationMode = allowedScenarios.has(requestedScenario) ? requestedScenario : "normal";
  let processing = false;
  let result = null;

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

  function formatTime(value) {
    if (!value) return "—";
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return "—";
    return new Intl.DateTimeFormat("zh-CN", {
      timeZone: "Asia/Shanghai",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false
    }).format(date);
  }

  function remainingText(expireAt) {
    const remaining = Math.max(0, new Date(expireAt).getTime() - D.now());
    const seconds = Math.ceil(remaining / 1000);
    return `${String(Math.floor(seconds / 60)).padStart(2, "0")}:${String(
      seconds % 60
    ).padStart(2, "0")}`;
  }

  function currentOrder() {
    return orderId ? D.getOrder(orderId) : null;
  }

  function buyerReturnUrl(paymentState, returnedOrderId) {
    const query = new URLSearchParams();
    if (paymentState) query.set("payment", paymentState);
    if (returnedOrderId) query.set("order", returnedOrderId);
    const suffix = query.toString();
    return `/shop${suffix ? `?${suffix}` : ""}`;
  }

  function goToBuyer(paymentState, returnedOrderId) {
    window.location.href = buyerReturnUrl(paymentState, returnedOrderId);
  }

  function renderError(title, message) {
    root.innerHTML = `
      <div class="checkout-error">
        <span class="checkout-result-mark">${icon("warning")}</span>
        <h1>${escapeHtml(title)}</h1>
        <p>${escapeHtml(message)}</p>
        <button class="btn btn-primary" type="button" data-action="back-buyer">
          ${icon("chevron-left")}返回商城
        </button>
      </div>
    `;
  }

  function renderCheckout() {
    const order = currentOrder();
    if (!order) {
      renderError("订单不存在", "当前支付订单无法查询，请返回商城订单列表后重试。");
      return;
    }
    if (order.buyerId !== buyerId) {
      renderError("无权支付该订单", "当前支付账号与订单归属不一致，系统已拒绝本次操作。");
      return;
    }
    if (order.status !== D.ORDER_STATUS.PENDING) {
      if (order.status === D.ORDER_STATUS.PAID) {
        renderResult({
          tone: "success",
          title: "支付成功",
          message: "该订单已完成支付，商城订单状态已更新为待发货。",
          amountCents: order.amountCents,
          rows: [
            ["订单编号", order.id],
            ["支付渠道", methodInfo.label],
            ["渠道交易号", order.paymentTxnNo || "已由渠道确认"],
            ["订单状态", "待发货"]
          ],
          primaryState: "success"
        });
        return;
      }
      renderError(
        "订单当前不可支付",
        `订单状态为${D.STATUS_LABEL[order.status] || order.status}，商城不会对失效订单重新收款。`
      );
      return;
    }

    const productImage = D.getProduct(order.productId)?.image || order.productImage || "";

    root.innerHTML = `
      <div class="checkout-body">
        <section class="checkout-amount">
          <span class="checkout-amount-label">待支付金额</span>
          <strong>${money(order.amountCents)}</strong>
          <div class="checkout-amount-meta">
            <span>剩余支付时间</span>
            <strong data-payment-countdown="${escapeHtml(order.expireAt)}">${remainingText(
              order.expireAt
            )}</strong>
          </div>
        </section>

        <section class="checkout-section">
          <div class="checkout-section-head">
            <strong>订单商品</strong>
            <span>${escapeHtml(order.id)}</span>
          </div>
          <div class="checkout-product">
            <img src="${escapeHtml(productImage)}" alt="${escapeHtml(order.productName)}" />
            <div class="checkout-product-copy">
              <strong>${escapeHtml(order.productName)}</strong>
              <span>${escapeHtml(order.productSpec)} · 数量 × ${order.qty}</span>
            </div>
            <span class="checkout-product-total">${money(order.amountCents)}</span>
          </div>
        </section>

        <section class="checkout-section">
          <div class="checkout-section-head">
            <strong>支付方式</strong>
            <span>渠道已选择</span>
          </div>
          <div class="checkout-channel">
            <span class="payment-brand ${methodInfo.tone}">${methodInfo.initial}</span>
            <div class="checkout-channel-copy">
              <strong>${methodInfo.label}</strong>
              <span>支付结果将返回 Shop 商城</span>
            </div>
            <span class="badge badge-paid">已选择</span>
          </div>
        </section>

      </div>

      <footer class="checkout-paybar">
        <button class="btn btn-primary btn-block" type="button" data-action="pay-now">
          ${icon("lock")}确认支付 ${money(order.amountCents)}
        </button>
        <button class="checkout-cancel" type="button" data-action="back-buyer">
          取消支付，返回商城
        </button>
      </footer>
    `;
  }

  function renderProcessing() {
    root.innerHTML = `
      <div class="checkout-processing" aria-live="polite">
        <span class="checkout-spinner" aria-hidden="true"></span>
        <strong>正在向${methodInfo.label}发起支付</strong>
        <p>请稍候，商城将在渠道返回结果后更新订单，不会提前显示支付成功。</p>
      </div>
    `;
  }

  function renderResult(nextResult) {
    result = nextResult;
    const state = nextResult.primaryState || null;
    root.innerHTML = `
      <div class="checkout-result">
        <div class="checkout-result-main">
          <span class="checkout-result-mark ${nextResult.tone}">
            ${icon(nextResult.tone === "success" ? "check" : nextResult.tone === "warning" ? "clock" : "warning")}
          </span>
          <h1 id="payment-result-title">${escapeHtml(nextResult.title)}</h1>
          <p class="checkout-result-message">${escapeHtml(nextResult.message)}</p>
          <strong class="checkout-result-amount">${money(nextResult.amountCents)}</strong>
          <div class="checkout-result-list">
            ${nextResult.rows
              .map(
                ([label, value]) => `
                  <div class="checkout-result-row">
                    <span>${escapeHtml(label)}</span>
                    <strong>${escapeHtml(value)}</strong>
                  </div>
                `
              )
              .join("")}
          </div>
        </div>
        <footer class="checkout-result-actions">
          <button
            class="btn btn-primary btn-block"
            type="button"
            data-action="return-buyer"
            data-payment-state="${escapeHtml(state || "")}"
          >
            ${icon("receipt")}返回商城查看订单
          </button>
          ${
            nextResult.canRetry
              ? `<button class="btn btn-secondary btn-block" type="button" data-action="retry-payment">${icon(
                  "refresh"
                )}重新支付</button>`
              : ""
          }
        </footer>
      </div>
    `;
  }

  function paymentMode() {
    if (simulationMode === "normal") return method === "alipay" ? "ali" : "success";
    return simulationMode === "network" ? "network_failure" : simulationMode;
  }

  async function executePayment() {
    if (processing) return;
    const order = currentOrder();
    if (!order) {
      renderError("订单不存在", "当前支付订单无法查询，请返回商城后重试。");
      return;
    }
    processing = true;
    renderProcessing();

    window.setTimeout(async () => {
      try {
        if (simulationMode === "late") {
          const timedOut = await D.timeoutOrder(order.id);
          const txnNo = `${method === "alipay" ? "ALI" : "WX"}_PAY_LATE_${Date.now()}`;
          const response = await D.receivePayment({
            orderId: order.id,
            txnNo,
            amountCents: order.amountCents,
            authOk: true,
            signatureOk: true,
            source: "checkout-late-simulation"
          });
          processing = false;
          renderResult({
            tone: "warning",
            title: "收款到账，订单已关闭",
            message:
              "渠道完成收款时订单已超过有效期。商城不会恢复订单或重新占用库存，该笔收款已进入客服异常支付待处理列表。",
            amountCents: order.amountCents,
            rows: [
              ["订单编号", order.id],
              ["支付渠道", methodInfo.label],
              ["服务端关单", timedOut.ok ? "已完成" : timedOut.message],
              ["异常记录", (response.data && response.data.anomaly && response.data.anomaly.id) || "已进入客服处理"],
              ["库存结果", "保持已释放，不再次占用"]
            ],
            primaryState: "late"
          });
          return;
        }

        if (simulationMode === "duplicate") {
          const first = await D.simulatePayment(
            order.id,
            method === "alipay" ? "ali" : "success",
            buyerId
          );
          if (!first.ok) {
            processing = false;
            renderResult({
              tone: "danger",
              title: "支付处理失败",
              message: first.message,
              amountCents: order.amountCents,
              rows: [
                ["订单编号", order.id],
                ["订单状态", D.STATUS_LABEL[D.getOrder(order.id).status]]
              ],
              primaryState: "failed"
            });
            return;
          }
          const txnNo = first.data.payment.txnNo;
          await D.receivePayment({
            orderId: order.id,
            txnNo,
            amountCents: order.amountCents,
            authOk: true,
            signatureOk: true,
            source: "network-retry"
          });
          processing = false;
          renderResult({
            tone: "success",
            title: "支付成功",
            message:
              "渠道首次通知已形成有效收款；网络重试携带同一渠道交易号，商城只记录一笔支付。",
            amountCents: order.amountCents,
            rows: [
              ["订单编号", order.id],
              ["支付渠道", methodInfo.label],
              ["渠道交易号", txnNo],
              ["有效收款", "1 笔"],
              ["重复通知", "已幂等去重"],
              ["订单状态", "待发货"]
            ],
            primaryState: "success"
          });
          return;
        }

        const response = await D.simulatePayment(order.id, paymentMode(), buyerId);
        processing = false;
        if (simulationMode === "network") {
          const latest = D.getOrder(order.id);
          renderResult({
            tone: "warning",
            title: "支付结果待确认",
            message:
              "网络在渠道回调前中断，商城暂时无法确认收款结果。订单仍保持待付款，库存继续占用。",
            amountCents: order.amountCents,
            rows: [
              ["订单编号", order.id],
              ["支付渠道", methodInfo.label],
              ["订单状态", D.STATUS_LABEL[latest.status]],
              ["库存结果", "继续占用"],
              ["建议操作", "返回订单后重新查询或支付"]
            ],
            primaryState: response.code === "network_unknown" ? "uncertain" : "failed",
            canRetry: latest.status === D.ORDER_STATUS.PENDING
          });
          return;
        }

        if (!response.ok) {
          renderResult({
            tone: "danger",
            title: "支付处理失败",
            message: response.message,
            amountCents: order.amountCents,
            rows: [
              ["订单编号", order.id],
              ["支付渠道", methodInfo.label],
              ["订单状态", D.STATUS_LABEL[D.getOrder(order.id).status]]
            ],
            primaryState: "failed"
          });
          return;
        }

        renderResult({
          tone: "success",
          title: "支付成功",
          message: "渠道回调已验签，支付金额与订单一致。商城订单已更新为待发货。",
          amountCents: order.amountCents,
          rows: [
            ["订单编号", order.id],
            ["支付渠道", methodInfo.label],
            ["渠道交易号", response.data.payment.txnNo],
            ["实付金额", money(response.data.payment.amountCents)],
            ["订单状态", "待发货"]
          ],
          primaryState: "success"
        });
      } catch (error) {
        processing = false;
        renderResult({
          tone: "danger",
          title: "支付处理异常",
          message: (error && error.message) || "支付通道异常，请返回商城重试。",
          amountCents: order.amountCents,
          rows: [
            ["订单编号", order.id],
            ["支付渠道", methodInfo.label]
          ],
          primaryState: "failed"
        });
      }
    }, 950);
  }

  root.addEventListener("click", (event) => {
    const actionButton = event.target.closest("[data-action]");
    if (!actionButton) return;
    const action = actionButton.dataset.action;
    if (action === "back-buyer") {
      goToBuyer("", "");
    } else if (action === "pay-now") {
      executePayment();
    } else if (action === "retry-payment") {
      simulationMode = "normal";
      result = null;
      renderCheckout();
    } else if (action === "return-buyer") {
      goToBuyer(actionButton.dataset.paymentState, orderId);
    }
  });

  window.setInterval(() => {
    const countdown = document.querySelector("[data-payment-countdown]");
    if (!countdown) return;
    const order = currentOrder();
    if (!order) return;
    countdown.textContent = remainingText(order.expireAt);
    if (remainingText(order.expireAt) === "00:00") {
      void D.tickTimeouts().then(() => renderCheckout());
    }
  }, 1000);

  // 真实后端：先同步一次订单数据，再渲染收银台
  D.refresh().then(() => renderCheckout());
})();
