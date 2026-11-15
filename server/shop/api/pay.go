// 包声明
package api

// 导入依赖包
import (
// 执行："encoding/json"
	"encoding/json"
// 执行："net/http"
	"net/http"
// 执行："net/url"
	"net/url"
// 执行："strings"
	"strings"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/middleware"
	"shop/shop/middleware"
// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/model/request"
	"shop/shop/model/request"
// 执行："shop/shop/model/response"
	"shop/shop/model/response"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// PayApi 模拟支付渠道与系统侧补偿接口

// 类型定义：type PayApi struct {
type PayApi struct {
// 执行：s *service.ServiceGroup
	s *service.ServiceGroup
// 结束当前代码块
}

// Notify 渠道支付成功通知：不依赖登录态，靠签名鉴权；任何分支都留痕可查。

// 函数定义：func (a PayApi) Notify(c *gin.Context) {
func (a PayApi) Notify(c *gin.Context) {
// 赋值/更新：raw, _ := c.GetRawData()
	raw, _ := c.GetRawData()
// 声明：var req request.PaymentNotify
	var req request.PaymentNotify
// 赋值/更新：_ = json.Unmarshal(raw, &req)
	_ = json.Unmarshal(raw, &req)

// 赋值/更新：res, err := a.s.PaymentService.Notify(req.OrderNo, req.ChannelTxnID, req.AmountCents, req.Sign, string(raw))
	res, err := a.s.PaymentService.Notify(req.OrderNo, req.ChannelTxnID, req.AmountCents, req.Sign, string(raw))
// 赋值/更新：code, msg := 0, "success"
	code, msg := 0, "success"
// 赋值/更新：status := http.StatusOK
	status := http.StatusOK
// 条件判断：if err != nil {
	if err != nil {
// 赋值/更新：code, msg = service.BizCode(err)
		code, msg = service.BizCode(err)
// 赋值/更新：status = statusOf(code)
		status = statusOf(code)
// 结束当前代码块
	}
// 执行：c.JSON(status, response.Body{Code: code, Data: res, Msg: msg})
	c.JSON(status, response.Body{Code: code, Data: res, Msg: msg})
// 结束当前代码块
}

// AlipayPay 生成支付宝收银台跳转地址（电脑网站支付，买家本人 + 待付款订单）。
// 回调地址按请求方实际访问域名动态生成（不写死配置）：内网穿透/域名变化零改动。

// 函数定义：func (a PayApi) AlipayPay(c *gin.Context) {
func (a PayApi) AlipayPay(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：orderNo := c.Param("orderNo")
	orderNo := c.Param("orderNo")
// 赋值/更新：payURL, err := a.s.PaymentService.AlipayPayURL(orderNo, claims.UserID, publicBaseURL(c))
	payURL, err := a.s.PaymentService.AlipayPayURL(orderNo, claims.UserID, publicBaseURL(c))
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, gin.H{"payURL": payURL})
	ok(c, gin.H{"payURL": payURL})
// 结束当前代码块
}

// publicBaseURL 推导当前请求对外的访问根地址：优先反向代理头（内网穿透会带），
// 否则用请求自身的 scheme + Host。例：经 http://xxx.cpolar.cn 访问则返回 http://xxx.cpolar.cn。

// 函数定义：func publicBaseURL(c *gin.Context) string {
func publicBaseURL(c *gin.Context) string {
// 赋值/更新：scheme := c.Request.Header.Get("X-Forwarded-Proto")
	scheme := c.Request.Header.Get("X-Forwarded-Proto")
// 条件判断：if i := strings.IndexByte(scheme, ','); i >= 0 {
	if i := strings.IndexByte(scheme, ','); i >= 0 {
// 赋值/更新：scheme = scheme[:i]
		scheme = scheme[:i]
// 结束当前代码块
	}
// 赋值/更新：scheme = strings.TrimSpace(scheme)
	scheme = strings.TrimSpace(scheme)
// 条件判断：if scheme == "" {
	if scheme == "" {
// 条件判断：if c.Request.TLS != nil {
		if c.Request.TLS != nil {
// 赋值/更新：scheme = "https"
			scheme = "https"
// 条件分支：} else {
		} else {
// 赋值/更新：scheme = "http"
			scheme = "http"
// 结束当前代码块
		}
// 结束当前代码块
	}
// 赋值/更新：host := c.Request.Header.Get("X-Forwarded-Host")
	host := c.Request.Header.Get("X-Forwarded-Host")
// 条件判断：if i := strings.IndexByte(host, ','); i >= 0 {
	if i := strings.IndexByte(host, ','); i >= 0 {
// 赋值/更新：host = host[:i]
		host = host[:i]
// 结束当前代码块
	}
// 赋值/更新：host = strings.TrimSpace(host)
	host = strings.TrimSpace(host)
// 条件判断：if host == "" {
	if host == "" {
// 赋值/更新：host = c.Request.Host
		host = c.Request.Host
// 结束当前代码块
	}
// 返回：return scheme + "://" + host
	return scheme + "://" + host
// 结束当前代码块
}

// AlipayNotify 支付宝异步通知回调：POST 表单，RSA2 验签；任何分支都留痕可查。
// 应答约定：受理成功返回纯文本 success，否则 fail（支付宝将按策略重发）。

// 函数定义：func (a PayApi) AlipayNotify(c *gin.Context) {
func (a PayApi) AlipayNotify(c *gin.Context) {
// 赋值/更新：raw, err := c.GetRawData()
	raw, err := c.GetRawData()
// 条件判断：if err != nil {
	if err != nil {
// 执行：c.String(http.StatusBadRequest, "fail")
		c.String(http.StatusBadRequest, "fail")
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：form, perr := url.ParseQuery(string(raw))
	form, perr := url.ParseQuery(string(raw))
// 条件判断：if perr != nil {
	if perr != nil {
// 执行：c.String(http.StatusBadRequest, "fail")
		c.String(http.StatusBadRequest, "fail")
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：res, err := a.s.PaymentService.NotifyAlipay(form, string(raw))
	res, err := a.s.PaymentService.NotifyAlipay(form, string(raw))
// 条件判断：if err == nil && res.Result == model.NoticeAccepted {
	if err == nil && res.Result == model.NoticeAccepted {
// 执行：c.String(http.StatusOK, "success")
		c.String(http.StatusOK, "success")
// 返回：return
		return
// 结束当前代码块
	}
// 执行：c.String(http.StatusOK, "fail")
	c.String(http.StatusOK, "fail")
// 结束当前代码块
}

// CloseExpired 服务端时间触发的超时关闭（演示可手动触发，也可由定时任务调用）

// 函数定义：func (a PayApi) CloseExpired(c *gin.Context) {
func (a PayApi) CloseExpired(c *gin.Context) {
// 赋值/更新：n, err := a.s.OrderService.CloseExpired()
	n, err := a.s.OrderService.CloseExpired()
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, gin.H{"closed": n})
	ok(c, gin.H{"closed": n})
// 结束当前代码块
}

// ChannelMode 切换模拟渠道退款结果（success | fail），用于复现「首次失败、再次成功」

// 函数定义：func (a PayApi) ChannelMode(c *gin.Context) {
func (a PayApi) ChannelMode(c *gin.Context) {
// 声明：var req request.ChannelMode
	var req request.ChannelMode
// 条件判断：if err := c.ShouldBindJSON(&req); err != nil {
	if err := c.ShouldBindJSON(&req); err != nil {
// 执行：fail(c, service.ErrInvalidParam)
		fail(c, service.ErrInvalidParam)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：sc, isSim := a.s.Channel.(*service.SimulatedChannel)
	sc, isSim := a.s.Channel.(*service.SimulatedChannel)
// 条件判断：if !isSim {
	if !isSim {
// 执行：fail(c, service.NewBizError(9001, "当前渠道不支持切换退款模式"))
		fail(c, service.NewBizError(9001, "当前渠道不支持切换退款模式"))
// 返回：return
		return
// 结束当前代码块
	}
// 执行：sc.SetRefundMode(req.RefundMode)
	sc.SetRefundMode(req.RefundMode)
// 执行：ok(c, gin.H{"refundMode": sc.RefundMode(), "refundCalls": sc.RefundCalls()})
	ok(c, gin.H{"refundMode": sc.RefundMode(), "refundCalls": sc.RefundCalls()})
// 结束当前代码块
}

// PendingOrders 模拟渠道侧可见的待支付订单（演示页面据此发起支付/重发/迟到通知）

// 函数定义：func (a PayApi) PendingOrders(c *gin.Context) {
func (a PayApi) PendingOrders(c *gin.Context) {
// 赋值/更新：list, err := a.s.OrderService.ListByStatus("pending_payment", 50)
	list, err := a.s.OrderService.ListByStatus("pending_payment", 50)
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：out := make([]gin.H, 0, len(list))
	out := make([]gin.H, 0, len(list))
// 循环处理：for _, o := range list {
	for _, o := range list {
// 赋值/更新：out = append(out, gin.H{
		out = append(out, gin.H{
// 执行："orderNo":     o.OrderNo,
			"orderNo":     o.OrderNo,
// 执行："productName": o.ProductName,
			"productName": o.ProductName,
// 执行："amountCents": o.AmountCents,
			"amountCents": o.AmountCents,
// 执行："expireAt":    o.ExpireAt,
			"expireAt":    o.ExpireAt,
// 执行："createdAt":   o.CreatedAt,
			"createdAt":   o.CreatedAt,
// 结束当前代码块
		})
// 结束当前代码块
	}
// 执行：ok(c, out)
	ok(c, out)
// 结束当前代码块
}

// ChannelState 查看渠道当前模式（演示与验收脚本用）

// 函数定义：func (a PayApi) ChannelState(c *gin.Context) {
func (a PayApi) ChannelState(c *gin.Context) {
// 赋值/更新：sc, isSim := a.s.Channel.(*service.SimulatedChannel)
	sc, isSim := a.s.Channel.(*service.SimulatedChannel)
// 条件判断：if !isSim {
	if !isSim {
// 执行：ok(c, gin.H{"refundMode": "external"})
		ok(c, gin.H{"refundMode": "external"})
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, gin.H{"refundMode": sc.RefundMode(), "refundCalls": sc.RefundCalls()})
	ok(c, gin.H{"refundMode": sc.RefundMode(), "refundCalls": sc.RefundCalls()})
// 结束当前代码块
}
