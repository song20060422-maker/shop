// 包声明
package api

// 导入依赖包
import (
// 执行："strconv"
	"strconv"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/middleware"
	"shop/shop/middleware"
// 执行："shop/shop/model/request"
	"shop/shop/model/request"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// StaffApi 员工侧接口：运营核查、仓库发货、客服退款与异常处理

// 类型定义：type StaffApi struct {
type StaffApi struct {
// 执行：s *service.ServiceGroup
	s *service.ServiceGroup
// 结束当前代码块
}

// ---------- 运营 ----------

// ProductUpdate 上下架与调价（不影响已有订单金额快照）

// 函数定义：func (a StaffApi) ProductUpdate(c *gin.Context) {
func (a StaffApi) ProductUpdate(c *gin.Context) {
// 声明：var req request.ProductUpdate
	var req request.ProductUpdate
// 条件判断：if err := c.ShouldBindJSON(&req); err != nil {
	if err := c.ShouldBindJSON(&req); err != nil {
// 执行：fail(c, service.ErrInvalidParam)
		fail(c, service.ErrInvalidParam)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：p, err := a.s.CatalogService.Update(claims.UserID, claims.Role, req.ProductID, req.OnSale, req.PriceCents)
	p, err := a.s.CatalogService.Update(claims.UserID, claims.Role, req.ProductID, req.OnSale, req.PriceCents)
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, p)
	ok(c, p)
// 结束当前代码块
}

// Summary 交易汇总：可售数量、已收款、已退款与状态分布

// 函数定义：func (a StaffApi) Summary(c *gin.Context) {
func (a StaffApi) Summary(c *gin.Context) {
// 赋值/更新：sum, err := a.s.ReportService.Summary()
	sum, err := a.s.ReportService.Summary()
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, sum)
	ok(c, sum)
// 结束当前代码块
}

// Orders 按状态查询全量订单

// 函数定义：func (a StaffApi) Orders(c *gin.Context) {
func (a StaffApi) Orders(c *gin.Context) {
// 赋值/更新：list, err := a.s.OrderService.ListByStatus(c.Query("status"), 0)
	list, err := a.s.OrderService.ListByStatus(c.Query("status"), 0)
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

// OrderTrace 订单全链路视图（运营/客服/仓库共用）

// 函数定义：func (a StaffApi) OrderTrace(c *gin.Context) {
func (a StaffApi) OrderTrace(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：detail, err := a.s.ReportService.OrderDetail(c.Param("orderNo"), claims.UserID, claims.Role)
	detail, err := a.s.ReportService.OrderDetail(c.Param("orderNo"), claims.UserID, claims.Role)
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

// Notices 支付通知处理结果查询（按 result 过滤）

// 函数定义：func (a StaffApi) Notices(c *gin.Context) {
func (a StaffApi) Notices(c *gin.Context) {
// 赋值/更新：no := c.Query("orderNo")
	no := c.Query("orderNo")
// 条件判断：if no != "" {
	if no != "" {
// 赋值/更新：list, err := a.s.PaymentService.Notices(no)
		list, err := a.s.PaymentService.Notices(no)
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
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
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

// Abnormals 异常收款列表

// 函数定义：func (a StaffApi) Abnormals(c *gin.Context) {
func (a StaffApi) Abnormals(c *gin.Context) {
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

// Ledger 库存流水（核对占用/释放/恢复守恒）

// 函数定义：func (a StaffApi) Ledger(c *gin.Context) {
func (a StaffApi) Ledger(c *gin.Context) {
// 赋值/更新：id, _ := strconv.ParseUint(c.Query("productId"), 10, 64)
	id, _ := strconv.ParseUint(c.Query("productId"), 10, 64)
// 赋值/更新：list, err := a.s.ReportService.Ledger(id)
	list, err := a.s.ReportService.Ledger(id)
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

// Trail 操作留痕

// 函数定义：func (a StaffApi) Trail(c *gin.Context) {
func (a StaffApi) Trail(c *gin.Context) {
// 赋值/更新：list, err := a.s.ReportService.Trail(c.Query("orderNo"))
	list, err := a.s.ReportService.Trail(c.Query("orderNo"))
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

// ---------- 仓库 ----------

// Shipable 待发货列表

// 函数定义：func (a StaffApi) Shipable(c *gin.Context) {
func (a StaffApi) Shipable(c *gin.Context) {
// 赋值/更新：list, err := a.s.ShipmentService.Shipable()
	list, err := a.s.ShipmentService.Shipable()
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

// Ship 发货登记

// 函数定义：func (a StaffApi) Ship(c *gin.Context) {
func (a StaffApi) Ship(c *gin.Context) {
// 声明：var req request.Ship
	var req request.Ship
// 条件判断：if err := c.ShouldBindJSON(&req); err != nil {
	if err := c.ShouldBindJSON(&req); err != nil {
// 执行：fail(c, service.ErrInvalidParam)
		fail(c, service.ErrInvalidParam)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：order, err := a.s.ShipmentService.Ship(claims.UserID, req.OrderNo, req.Carrier, req.TrackingNo)
	order, err := a.s.ShipmentService.Ship(claims.UserID, req.OrderNo, req.Carrier, req.TrackingNo)
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, order)
	ok(c, order)
// 结束当前代码块
}

// Shipments 发货记录

// 函数定义：func (a StaffApi) Shipments(c *gin.Context) {
func (a StaffApi) Shipments(c *gin.Context) {
// 赋值/更新：list, err := a.s.ShipmentService.List()
	list, err := a.s.ShipmentService.List()
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

// ---------- 客服 ----------

// Refunds 退款单列表

// 函数定义：func (a StaffApi) Refunds(c *gin.Context) {
func (a StaffApi) Refunds(c *gin.Context) {
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

// RefundAudit 审核售前退款

// 函数定义：func (a StaffApi) RefundAudit(c *gin.Context) {
func (a StaffApi) RefundAudit(c *gin.Context) {
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
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：refund, err := a.s.RefundService.Audit(claims.UserID, req.OrderNo, req.Approve, req.Reason)
	refund, err := a.s.RefundService.Audit(claims.UserID, req.OrderNo, req.Approve, req.Reason)
// 条件判断：if err != nil {
	if err != nil {
		// 审核已生效但渠道失败时，把退款单一起返回，方便客服重试
		c.JSON(409, gin.H{"code": 3004, "msg": service.BizCodeMsg(err), "data": refund})
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, refund)
	ok(c, refund)
// 结束当前代码块
}

// RefundRetry 渠道失败后重试

// 函数定义：func (a StaffApi) RefundRetry(c *gin.Context) {
func (a StaffApi) RefundRetry(c *gin.Context) {
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
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：refund, err := a.s.RefundService.Retry(claims.UserID, req.RefundNo)
	refund, err := a.s.RefundService.Retry(claims.UserID, req.RefundNo)
// 条件判断：if err != nil {
	if err != nil {
// 执行：c.JSON(409, gin.H{"code": 3004, "msg": service.BizCodeMsg(err), "data": refund})
		c.JSON(409, gin.H{"code": 3004, "msg": service.BizCodeMsg(err), "data": refund})
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, refund)
	ok(c, refund)
// 结束当前代码块
}

// HandleAbnormal 处理迟到支付：核实后发起退款

// 函数定义：func (a StaffApi) HandleAbnormal(c *gin.Context) {
func (a StaffApi) HandleAbnormal(c *gin.Context) {
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
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：refund, err := a.s.RefundService.HandleAbnormal(claims.UserID, req.ChannelTxnID, req.Reason)
	refund, err := a.s.RefundService.HandleAbnormal(claims.UserID, req.ChannelTxnID, req.Reason)
// 条件判断：if err != nil {
	if err != nil {
// 执行：c.JSON(409, gin.H{"code": 3004, "msg": service.BizCodeMsg(err), "data": refund})
		c.JSON(409, gin.H{"code": 3004, "msg": service.BizCodeMsg(err), "data": refund})
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, refund)
	ok(c, refund)
// 结束当前代码块
}
