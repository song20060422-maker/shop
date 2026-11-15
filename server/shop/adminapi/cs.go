// 包声明
package adminapi

// 导入依赖包
import (
// 执行："strconv"
	"strconv"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/model/common/response"
	"shop/model/common/response"
// 执行："shop/shop/model/request"
	"shop/shop/model/request"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// ---------- 客服（GVA 管理端） ----------

// Refunds 退款单列表

// 函数定义：func (a *Api) Refunds(c *gin.Context) {
func (a *Api) Refunds(c *gin.Context) {
// 赋值/更新：list, err := a.s.RefundService.List(c.Query("status"))
	list, err := a.s.RefundService.List(c.Query("status"))
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, list)
	ok(c, list)
// 结束当前代码块
}

// RefundAudit 审核售前退款（同意→渠道退款；拒绝→回到待发货）

// 函数定义：func (a *Api) RefundAudit(c *gin.Context) {
func (a *Api) RefundAudit(c *gin.Context) {
// 声明：var req request.RefundAudit
	var req request.RefundAudit
// 条件判断：if err := c.ShouldBindJSON(&req); err != nil {
	if err := c.ShouldBindJSON(&req); err != nil {
// 执行：fail(c, service.ErrInvalidParam)
		fail(c, service.ErrInvalidParam)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：actorID, _ := actor(c, "cs")
	actorID, _ := actor(c, "cs")
// 赋值/更新：refund, err := a.s.RefundService.Audit(actorID, req.OrderNo, req.Approve, req.Reason)
	refund, err := a.s.RefundService.Audit(actorID, req.OrderNo, req.Approve, req.Reason)
// 条件判断：if err != nil {
	if err != nil {
		// 审核已生效但渠道失败时，把退款单一起返回，方便客服重试
		response.FailWithDetailed(refund, service.BizCodeMsg(err), c)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, refund)
	ok(c, refund)
// 结束当前代码块
}

// RefundRetry 渠道失败后重试

// 函数定义：func (a *Api) RefundRetry(c *gin.Context) {
func (a *Api) RefundRetry(c *gin.Context) {
// 声明：var req request.RefundRetry
	var req request.RefundRetry
// 条件判断：if err := c.ShouldBindJSON(&req); err != nil {
	if err := c.ShouldBindJSON(&req); err != nil {
// 执行：fail(c, service.ErrInvalidParam)
		fail(c, service.ErrInvalidParam)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：actorID, _ := actor(c, "cs")
	actorID, _ := actor(c, "cs")
// 赋值/更新：refund, err := a.s.RefundService.Retry(actorID, req.RefundNo)
	refund, err := a.s.RefundService.Retry(actorID, req.RefundNo)
// 条件判断：if err != nil {
	if err != nil {
// 执行：response.FailWithDetailed(refund, service.BizCodeMsg(err), c)
		response.FailWithDetailed(refund, service.BizCodeMsg(err), c)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, refund)
	ok(c, refund)
// 结束当前代码块
}

// Abnormals 异常收款列表（客服）

// 函数定义：func (a *Api) CsAbnormals(c *gin.Context) {
func (a *Api) CsAbnormals(c *gin.Context) {
// 赋值/更新：list, err := a.s.PaymentService.AbnormalList(c.Query("status"))
	list, err := a.s.PaymentService.AbnormalList(c.Query("status"))
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, list)
	ok(c, list)
// 结束当前代码块
}

// HandleAbnormal 处理迟到支付：核实后发起退款

// 函数定义：func (a *Api) HandleAbnormal(c *gin.Context) {
func (a *Api) HandleAbnormal(c *gin.Context) {
// 声明：var req request.AbnormalHandle
	var req request.AbnormalHandle
// 条件判断：if err := c.ShouldBindJSON(&req); err != nil {
	if err := c.ShouldBindJSON(&req); err != nil {
// 执行：fail(c, service.ErrInvalidParam)
		fail(c, service.ErrInvalidParam)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：actorID, _ := actor(c, "cs")
	actorID, _ := actor(c, "cs")
// 赋值/更新：refund, err := a.s.RefundService.HandleAbnormal(actorID, req.ChannelTxnID, req.Reason)
	refund, err := a.s.RefundService.HandleAbnormal(actorID, req.ChannelTxnID, req.Reason)
// 条件判断：if err != nil {
	if err != nil {
// 执行：response.FailWithDetailed(refund, service.BizCodeMsg(err), c)
		response.FailWithDetailed(refund, service.BizCodeMsg(err), c)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, refund)
	ok(c, refund)
// 结束当前代码块
}

// CsNotices 支付通知处理记录（客服）

// 函数定义：func (a *Api) CsNotices(c *gin.Context) {
func (a *Api) CsNotices(c *gin.Context) {
// 赋值/更新：limit := 100
	limit := 100
// 条件判断：if v, err := strconv.Atoi(c.Query("limit")); err == nil && v > 0 {
	if v, err := strconv.Atoi(c.Query("limit")); err == nil && v > 0 {
// 赋值/更新：limit = v
		limit = v
// 结束当前代码块
	}
// 赋值/更新：list, err := a.s.PaymentService.NoticesByResult(c.Query("result"), limit)
	list, err := a.s.PaymentService.NoticesByResult(c.Query("result"), limit)
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, list)
	ok(c, list)
// 结束当前代码块
}

// CsOrderTrace 订单全链路视图（客服视角）

// 函数定义：func (a *Api) CsOrderTrace(c *gin.Context) {
func (a *Api) CsOrderTrace(c *gin.Context) {
// 赋值/更新：actorID, role := actor(c, "cs")
	actorID, role := actor(c, "cs")
// 赋值/更新：detail, err := a.s.ReportService.OrderDetail(c.Param("orderNo"), actorID, role)
	detail, err := a.s.ReportService.OrderDetail(c.Param("orderNo"), actorID, role)
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, detail)
	ok(c, detail)
// 结束当前代码块
}

// ChannelMode 切换模拟渠道退款结果（success | fail），用于复现「首次失败、再次成功」

// 函数定义：func (a *Api) ChannelMode(c *gin.Context) {
func (a *Api) ChannelMode(c *gin.Context) {
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
