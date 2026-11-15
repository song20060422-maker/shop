// 包声明
package adminapi

// 导入依赖包
import (
// 执行："strconv"
	"strconv"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/model/request"
	"shop/shop/model/request"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// ---------- 运营（GVA 管理端） ----------

// Summary 交易汇总：可售数量、已收款、已退款与状态分布

// 函数定义：func (a *Api) Summary(c *gin.Context) {
func (a *Api) Summary(c *gin.Context) {
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

// Products 全部商品（含下架）

// 函数定义：func (a *Api) Products(c *gin.Context) {
func (a *Api) Products(c *gin.Context) {
// 赋值/更新：list, err := a.s.CatalogService.All()
	list, err := a.s.CatalogService.All()
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

// ProductUpdate 运营维护商品：上下架/调价/直接覆盖可售数量

// 函数定义：func (a *Api) ProductUpdate(c *gin.Context) {
func (a *Api) ProductUpdate(c *gin.Context) {
// 赋值/更新：id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 条件判断：if err != nil || id == 0 {
	if err != nil || id == 0 {
// 执行：fail(c, service.ErrInvalidParam)
		fail(c, service.ErrInvalidParam)
// 返回：return
		return
// 结束当前代码块
	}
// 声明：var req request.AdminProductUpdate
	var req request.AdminProductUpdate
// 条件判断：if err := c.ShouldBindJSON(&req); err != nil {
	if err := c.ShouldBindJSON(&req); err != nil {
// 执行：fail(c, service.ErrInvalidParam)
		fail(c, service.ErrInvalidParam)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：actorID, role := actor(c, "ops")
	actorID, role := actor(c, "ops")
// 条件判断：if req.Available != nil {
	if req.Available != nil {
// 赋值/更新：p, err := a.s.CatalogService.OverwriteAvailable(actorID, role, id, *req.Available)
		p, err := a.s.CatalogService.OverwriteAvailable(actorID, role, id, *req.Available)
// 条件判断：if err != nil {
		if err != nil {
// 执行：fail(c, err)
			fail(c, err)
// 返回：return
			return
// 结束当前代码块
		}
// 执行：c.Set("shopProduct", p)
		c.Set("shopProduct", p)
// 结束当前代码块
	}
// 条件判断：if req.OnSale != nil || req.PriceCents != nil {
	if req.OnSale != nil || req.PriceCents != nil {
// 赋值/更新：p, err := a.s.CatalogService.Update(actorID, role, id, req.OnSale, req.PriceCents)
		p, err := a.s.CatalogService.Update(actorID, role, id, req.OnSale, req.PriceCents)
// 条件判断：if err != nil {
		if err != nil {
// 执行：fail(c, err)
			fail(c, err)
// 返回：return
			return
// 结束当前代码块
		}
// 执行：c.Set("shopProduct", p)
		c.Set("shopProduct", p)
// 结束当前代码块
	}
// 赋值/更新：pv, existed := c.Get("shopProduct")
	pv, existed := c.Get("shopProduct")
// 条件判断：if !existed {
	if !existed {
// 执行：fail(c, service.ErrInvalidParam)
		fail(c, service.ErrInvalidParam)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, pv)
	ok(c, pv)
// 结束当前代码块
}

// Orders 按状态查询全量订单

// 函数定义：func (a *Api) Orders(c *gin.Context) {
func (a *Api) Orders(c *gin.Context) {
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

// OrderTrace 订单全链路核查视图

// 函数定义：func (a *Api) OrderTrace(c *gin.Context) {
func (a *Api) OrderTrace(c *gin.Context) {
// 赋值/更新：actorID, role := actor(c, "ops")
	actorID, role := actor(c, "ops")
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

// Notices 支付通知处理结果查询（按 orderNo / result 过滤）

// 函数定义：func (a *Api) Notices(c *gin.Context) {
func (a *Api) Notices(c *gin.Context) {
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

// 函数定义：func (a *Api) Abnormals(c *gin.Context) {
func (a *Api) Abnormals(c *gin.Context) {
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

// 函数定义：func (a *Api) Ledger(c *gin.Context) {
func (a *Api) Ledger(c *gin.Context) {
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

// 函数定义：func (a *Api) Trail(c *gin.Context) {
func (a *Api) Trail(c *gin.Context) {
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
