// GVA 管理端（后台）商城路由：挂在 PrivateGroup，走 GVA JWT + Casbin。
package router

// 导入依赖包
import (
// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/adminapi"
	"shop/shop/adminapi"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// InitShopAdminRouter 注册商城后台管理接口（运营/仓库/客服），前缀 /shop。
// 调用方需保证 group 已经过 GVA JWTAuth + CasbinHandler。

// 函数定义：func InitShopAdminRouter(group *gin.RouterGroup, g *service.ServiceGroup) {
func InitShopAdminRouter(group *gin.RouterGroup, g *service.ServiceGroup) {
// 赋值/更新：a := adminapi.NewApi(g)
	a := adminapi.NewApi(g)

// 赋值/更新：ops := group.Group("/shop/ops")
	ops := group.Group("/shop/ops")
// 执行：{
	{
// 执行：ops.GET("/summary", a.Summary)
		ops.GET("/summary", a.Summary)
// 执行：ops.GET("/products", a.Products)
		ops.GET("/products", a.Products)
// 执行：ops.PATCH("/products/:id", a.ProductUpdate)
		ops.PATCH("/products/:id", a.ProductUpdate)
// 执行：ops.GET("/orders", a.Orders)
		ops.GET("/orders", a.Orders)
// 执行：ops.GET("/orders/:orderNo", a.OrderTrace)
		ops.GET("/orders/:orderNo", a.OrderTrace)
// 执行：ops.GET("/notices", a.Notices)
		ops.GET("/notices", a.Notices)
// 执行：ops.GET("/abnormals", a.Abnormals)
		ops.GET("/abnormals", a.Abnormals)
// 执行：ops.GET("/ledger", a.Ledger)
		ops.GET("/ledger", a.Ledger)
// 执行：ops.GET("/trail", a.Trail)
		ops.GET("/trail", a.Trail)
// 结束当前代码块
	}

// 赋值/更新：wh := group.Group("/shop/warehouse")
	wh := group.Group("/shop/warehouse")
// 执行：{
	{
// 执行：wh.GET("/shipable", a.Shipable)
		wh.GET("/shipable", a.Shipable)
// 执行：wh.POST("/ship", a.Ship)
		wh.POST("/ship", a.Ship)
// 执行：wh.GET("/shipments", a.Shipments)
		wh.GET("/shipments", a.Shipments)
// 执行：wh.GET("/orders/:orderNo", a.WarehouseOrderTrace)
		wh.GET("/orders/:orderNo", a.WarehouseOrderTrace)
// 结束当前代码块
	}

// 赋值/更新：cs := group.Group("/shop/cs")
	cs := group.Group("/shop/cs")
// 执行：{
	{
// 执行：cs.GET("/refunds", a.Refunds)
		cs.GET("/refunds", a.Refunds)
// 执行：cs.POST("/refunds/audit", a.RefundAudit)
		cs.POST("/refunds/audit", a.RefundAudit)
// 执行：cs.POST("/refunds/retry", a.RefundRetry)
		cs.POST("/refunds/retry", a.RefundRetry)
// 执行：cs.GET("/abnormals", a.CsAbnormals)
		cs.GET("/abnormals", a.CsAbnormals)
// 执行：cs.POST("/abnormals/handle", a.HandleAbnormal)
		cs.POST("/abnormals/handle", a.HandleAbnormal)
// 执行：cs.GET("/notices", a.CsNotices)
		cs.GET("/notices", a.CsNotices)
// 执行：cs.GET("/orders/:orderNo", a.CsOrderTrace)
		cs.GET("/orders/:orderNo", a.CsOrderTrace)
// 执行：cs.POST("/channel/mode", a.ChannelMode)
		cs.POST("/channel/mode", a.ChannelMode)
// 结束当前代码块
	}
// 结束当前代码块
}
