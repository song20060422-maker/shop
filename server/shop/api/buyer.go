// 包声明
package api

// 导入依赖包
import (
// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/middleware"
	"shop/shop/middleware"
// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/model/request"
	"shop/shop/model/request"
// 执行："shop/shop/service"
	"shop/shop/service"
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// BuyerApi 买家侧接口：浏览、下单、查询、取消、售前退款

// 类型定义：type BuyerApi struct {
type BuyerApi struct {
// 执行：s *service.ServiceGroup
	s *service.ServiceGroup
// 结束当前代码块
}

// Products 在售商品（含可售数量与占用情况）

// 函数定义：func (a BuyerApi) Products(c *gin.Context) {
func (a BuyerApi) Products(c *gin.Context) {
// 赋值/更新：list, err := a.s.CatalogService.OnSale()
	list, err := a.s.CatalogService.OnSale()
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
// 循环处理：for _, p := range list {
	for _, p := range list {
// 赋值/更新：out = append(out, gin.H{
		out = append(out, gin.H{
// 执行："id":         p.ID,
			"id":         p.ID,
// 执行："sku":        p.SKU,
			"sku":        p.SKU,
// 执行："name":       p.Name,
			"name":       p.Name,
// 执行："priceCents": p.PriceCents,
			"priceCents": p.PriceCents,
// 执行："priceText":  utils.Fen(p.PriceCents),
			"priceText":  utils.Fen(p.PriceCents),
// 执行："available":  p.Available,
			"available":  p.Available,
// 执行："onSale":     p.OnSale,
			"onSale":     p.OnSale,
// 结束当前代码块
		})
// 结束当前代码块
	}
// 执行：ok(c, out)
	ok(c, out)
// 结束当前代码块
}

// CreateOrder 下单：数量 1～2 件，成功即可售数量已反映占用

// 函数定义：func (a BuyerApi) CreateOrder(c *gin.Context) {
func (a BuyerApi) CreateOrder(c *gin.Context) {
// 声明：var req request.CreateOrder
	var req request.CreateOrder
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
// 赋值/更新：order, err := a.s.OrderService.Create(claims.UserID, req.ProductID, req.Qty)
	order, err := a.s.OrderService.Create(claims.UserID, req.ProductID, req.Qty)
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：available, _ := a.s.CatalogService.CurrentAvailable(req.ProductID)
	available, _ := a.s.CatalogService.CurrentAvailable(req.ProductID)
// 执行：ok(c, gin.H{"order": order, "availableAfter": available, "amountText": utils.Fen(order.AmountCents)})
	ok(c, gin.H{"order": order, "availableAfter": available, "amountText": utils.Fen(order.AmountCents)})
// 结束当前代码块
}

// MyOrders 我的订单

// 函数定义：func (a BuyerApi) MyOrders(c *gin.Context) {
func (a BuyerApi) MyOrders(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：list, err := a.s.OrderService.ListByBuyer(claims.UserID, c.Query("status"))
	list, err := a.s.OrderService.ListByBuyer(claims.UserID, c.Query("status"))
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

// OrderDetail 订单详情：可见收款、退款、发货、通知与操作轨迹

// 函数定义：func (a BuyerApi) OrderDetail(c *gin.Context) {
func (a BuyerApi) OrderDetail(c *gin.Context) {
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

// Cancel 取消待付款订单并释放占用

// 函数定义：func (a BuyerApi) Cancel(c *gin.Context) {
func (a BuyerApi) Cancel(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：order, err := a.s.OrderService.Cancel(claims.UserID, c.Param("orderNo"))
	order, err := a.s.OrderService.Cancel(claims.UserID, c.Param("orderNo"))
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

// RefundApply 售前退款申请（已付款未发货）

// 函数定义：func (a BuyerApi) RefundApply(c *gin.Context) {
func (a BuyerApi) RefundApply(c *gin.Context) {
// 声明：var req request.RefundApply
	var req request.RefundApply
// 条件判断：if err := c.ShouldBindJSON(&req); err != nil || req.OrderNo == "" {
	if err := c.ShouldBindJSON(&req); err != nil || req.OrderNo == "" {
// 赋值/更新：req.OrderNo = c.Param("orderNo")
		req.OrderNo = c.Param("orderNo")
// 结束当前代码块
	}
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：refund, err := a.s.RefundService.Apply(claims.UserID, req.OrderNo, req.Reason)
	refund, err := a.s.RefundService.Apply(claims.UserID, req.OrderNo, req.Reason)
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, gin.H{"refund": refund, "refundPending": refund.AuditStatus == model.AuditPending})
	ok(c, gin.H{"refund": refund, "refundPending": refund.AuditStatus == model.AuditPending})
// 结束当前代码块
}

// MyRefund 查询本人订单退款结果

// 函数定义：func (a BuyerApi) MyRefund(c *gin.Context) {
func (a BuyerApi) MyRefund(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：orderNo := c.Param("orderNo")
	orderNo := c.Param("orderNo")
// 赋值/更新：order, err := a.s.OrderService.Get(orderNo)
	order, err := a.s.OrderService.Get(orderNo)
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 条件判断：if order.BuyerID != claims.UserID {
	if order.BuyerID != claims.UserID {
// 执行：fail(c, service.ErrForbidden)
		fail(c, service.ErrForbidden)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：refund, err := a.s.RefundService.Get(orderNo)
	refund, err := a.s.RefundService.Get(orderNo)
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, refund)
	ok(c, refund)
// 结束当前代码块
}
