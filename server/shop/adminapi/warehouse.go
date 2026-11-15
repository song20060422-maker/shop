// 包声明
package adminapi

// 导入依赖包
import (
// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/model/request"
	"shop/shop/model/request"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// ---------- 仓库（GVA 管理端） ----------

// Shipable 待发货列表

// 函数定义：func (a *Api) Shipable(c *gin.Context) {
func (a *Api) Shipable(c *gin.Context) {
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

// Ship 发货登记（一单只能登记一次）

// 函数定义：func (a *Api) Ship(c *gin.Context) {
func (a *Api) Ship(c *gin.Context) {
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
// 赋值/更新：actorID, _ := actor(c, "warehouse")
	actorID, _ := actor(c, "warehouse")
// 赋值/更新：order, err := a.s.ShipmentService.Ship(actorID, req.OrderNo, req.Carrier, req.TrackingNo)
	order, err := a.s.ShipmentService.Ship(actorID, req.OrderNo, req.Carrier, req.TrackingNo)
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

// 函数定义：func (a *Api) Shipments(c *gin.Context) {
func (a *Api) Shipments(c *gin.Context) {
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

// WarehouseOrderTrace 订单全链路视图（仓库视角）

// 函数定义：func (a *Api) WarehouseOrderTrace(c *gin.Context) {
func (a *Api) WarehouseOrderTrace(c *gin.Context) {
// 赋值/更新：actorID, role := actor(c, "warehouse")
	actorID, role := actor(c, "warehouse")
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
