// Package router 按 GVA 分层约定注册 Shop 业务路由，挂载在 GVA 业务路由扩展点。
package router

// 导入依赖包
import (
// 执行："errors"
	"errors"
// 执行："log"
	"log"
// 执行："net/http"
	"net/http"
// 执行："net/http/httputil"
	"net/http/httputil"
// 执行："net/url"
	"net/url"
// 执行："strings"
	"strings"
// 执行："time"
	"time"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/api"
	"shop/shop/api"
// 执行："shop/shop/config"
	"shop/shop/config"
// 执行："shop/shop/middleware"
	"shop/shop/middleware"
// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/model/response"
	"shop/shop/model/response"
// 执行："shop/shop/service"
	"shop/shop/service"
// 执行："shop/shop/web"
	"shop/shop/web"
// 结束当前代码块
)

// InitShopRouter 装配买家端 Gin 服务路由：公开区、登录区、买家/运营/仓库/客服接口分区，
// 页面仅托管买家端 /shop（后台页面由 GVA 管理端 Vue 提供）。

// 函数定义：func InitShopRouter(group *gin.RouterGroup, g *service.ServiceGroup, cfg config.Config, now func() time.Time) ...
func InitShopRouter(group *gin.RouterGroup, g *service.ServiceGroup, cfg config.Config, now func() time.Time) {
// 赋值/更新：ag := api.NewApiGroup(g)
	ag := api.NewApiGroup(g)
	// 接口契约兼容层：复用同一套 service，仅做入参/出参契约适配
	ca := api.NewCompatApi(g, cfg, now)

	// PublicGuard 挂在公开区路由链最前：命中契约语义不同的同名路径时委派给兼容实现
	pub := group.Group("/api", ca.PublicGuard())
// 执行：{
	{
// 执行：pub.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok", "module": "shop"}) })
		pub.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok", "module": "shop"}) })
// 执行：pub.POST("/login", ag.Login)
		pub.POST("/login", ag.Login)
		// 账号列表：演示模式下供演示页选择账号；生产环境关闭，防止枚举账号
		pub.GET("/users", func(c *gin.Context) {
// 条件判断：if !cfg.Shop.DemoMode {
			if !cfg.Shop.DemoMode {
// 执行：c.JSON(http.StatusForbidden, response.Fail(2002, "演示模式未开启，禁止查询账号列表"))
				c.JSON(http.StatusForbidden, response.Fail(2002, "演示模式未开启，禁止查询账号列表"))
// 返回：return
				return
// 结束当前代码块
			}
// 执行：ag.Users(c)
			ag.Users(c)
// 结束当前代码块
		})
// 执行：pub.GET("/products", ag.Products)
		pub.GET("/products", ag.Products)
		// 模拟渠道：待支付订单查询、支付通知回调（靠签名鉴权，不需要登录态）
		pub.GET("/pay/pending", ag.PendingOrders)
// 执行：pub.POST("/pay/notify", ag.Notify)
		pub.POST("/pay/notify", ag.Notify)
// 执行：pub.POST("/pay/alipay/notify", ag.AlipayNotify)
		pub.POST("/pay/alipay/notify", ag.AlipayNotify)
// 执行：pub.GET("/pay/channel", ag.ChannelState)
		pub.GET("/pay/channel", ag.ChannelState)
		// 超时关单测试钩子：仅演示模式可用（生产由定时任务触发）
		pub.POST("/jobs/close-expired", func(c *gin.Context) {
// 条件判断：if !cfg.Shop.DemoMode {
			if !cfg.Shop.DemoMode {
// 执行：c.JSON(http.StatusForbidden, response.Fail(2002, "演示模式未开启，禁止手动触发关单"))
				c.JSON(http.StatusForbidden, response.Fail(2002, "演示模式未开启，禁止手动触发关单"))
// 返回：return
				return
// 结束当前代码块
			}
// 执行：ag.CloseExpired(c)
			ag.CloseExpired(c)
// 结束当前代码块
		})
		// 演示数据重置：仅 shop.demo-mode=true 时可用
		pub.POST("/demo/reset", func(c *gin.Context) {
// 条件判断：if !cfg.Shop.DemoMode {
			if !cfg.Shop.DemoMode {
// 执行：c.JSON(http.StatusForbidden, response.Fail(2002, "演示模式未开启，禁止重置"))
				c.JSON(http.StatusForbidden, response.Fail(2002, "演示模式未开启，禁止重置"))
// 返回：return
				return
// 结束当前代码块
			}
// 执行：ag.DemoApi.Reset(c)
			ag.DemoApi.Reset(c)
// 结束当前代码块
		})
// 结束当前代码块
	}

	// 内置页面只保留买家端（随服务发布，无需前端构建）：/shop 买家 + 模拟支付渠道页。
	// 运营/仓库/客服后台页面统一由 GVA 管理端 Vue（web/src/view/shop/*.vue）提供，
	// 接口走 GVA /shop/**（gva 登录 + Casbin），本服务不再重复内置后台页。
	page := func(name string) gin.HandlerFunc {
// 返回：return func(c *gin.Context) {
		return func(c *gin.Context) {
// 声明：var (
			var (
// 执行：data []byte
				data []byte
// 执行：err  error
				err  error
// 结束当前代码块
			)
// 多路分支：switch name {
			switch name {
// 匹配分支：case "buyer":
			case "buyer":
// 赋值/更新：data, err = web.Buyer()
				data, err = web.Buyer()
// 匹配分支：case "payment":
			case "payment":
				// 模拟收银台页（原型对齐）：微信模拟渠道通知演示入口
				data, err = web.Payment()
// 默认分支
			default:
// 赋值/更新：err = errors.New("未知页面")
				err = errors.New("未知页面")
// 结束当前代码块
			}
// 条件判断：if err != nil {
			if err != nil {
// 执行：c.String(http.StatusInternalServerError, "页面加载失败")
				c.String(http.StatusInternalServerError, "页面加载失败")
// 返回：return
				return
// 结束当前代码块
			}
// 执行：c.Data(http.StatusOK, "text/html; charset=utf-8", data)
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
// 结束当前代码块
		}
// 结束当前代码块
	}
	// 根路径与 /index.html 统一跳转买家页（历史入口兼容）
	group.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/shop") })
// 执行：group.GET("/index.html", func(c *gin.Context) { c.Redirect(http.StatusFound, "/shop") })
	group.GET("/index.html", func(c *gin.Context) { c.Redirect(http.StatusFound, "/shop") })
// 执行：group.GET("/shop", page("buyer"))
	group.GET("/shop", page("buyer"))
// 执行：group.GET("/shop/", page("buyer"))
	group.GET("/shop/", page("buyer"))
	// 模拟收银台页与页面静态资源（原型对齐，随二进制发布）
	group.GET("/shop/payment", page("payment"))
// 执行：group.GET("/shop/assets/:name", func(c *gin.Context) {
	group.GET("/shop/assets/:name", func(c *gin.Context) {
// 赋值/更新：data, contentType, ok := web.Asset(c.Param("name"))
		data, contentType, ok := web.Asset(c.Param("name"))
// 条件判断：if !ok {
		if !ok {
// 执行：c.String(http.StatusNotFound, "资源不存在")
			c.String(http.StatusNotFound, "资源不存在")
// 返回：return
			return
// 结束当前代码块
		}
// 执行：c.Data(http.StatusOK, contentType, data)
		c.Data(http.StatusOK, contentType, data)
// 结束当前代码块
	})
	// 后台三页入口已迁移到 GVA 管理端：命中旧路径时提示去向，避免静默 404
	adminHint := func(page string) gin.HandlerFunc {
// 返回：return func(c *gin.Context) {
		return func(c *gin.Context) {
// 执行：c.String(http.StatusNotFound,
			c.String(http.StatusNotFound,
// 执行："后台管理页面已迁移至 GVA 管理端（Vue）：请在 http://127.0.0.1:8080 登录后打开「商城管理 → %s」；买家端仅保留 /shop", page)
				"后台管理页面已迁移至 GVA 管理端（Vue）：请在 http://127.0.0.1:8080 登录后打开「商城管理 → %s」；买家端仅保留 /shop", page)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 执行：group.GET("/ops", adminHint("运营核查"))
	group.GET("/ops", adminHint("运营核查"))
// 执行：group.GET("/warehouse", adminHint("仓库发货"))
	group.GET("/warehouse", adminHint("仓库发货"))
// 执行：group.GET("/cs", adminHint("客服退款"))
	group.GET("/cs", adminHint("客服退款"))

	// /gva/* 反向代理到 GVA 管理端（历史兼容保留）：/gva/base/login、/gva/shop/** 等
	// 可经本服务同源访问，目标地址取 system.admin-target。正式后台入口为 GVA Vue（:8080）。
	adminURL, err := url.Parse(cfg.System.AdminTarget)
// 条件判断：if err != nil || adminURL.Scheme == "" || adminURL.Host == "" {
	if err != nil || adminURL.Scheme == "" || adminURL.Host == "" {
// 执行：log.Printf("shop: system.admin-target 非法 %q，/gva 反向代理未启用", cfg.System.AdminTarget)
		log.Printf("shop: system.admin-target 非法 %q，/gva 反向代理未启用", cfg.System.AdminTarget)
// 条件分支：} else {
	} else {
		// 剥掉 /gva 前缀再转发：/gva/base/login -> {admin-target}/base/login
		proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
// 赋值/更新：req.URL.Scheme = adminURL.Scheme
			req.URL.Scheme = adminURL.Scheme
// 赋值/更新：req.URL.Host = adminURL.Host
			req.URL.Host = adminURL.Host
// 赋值/更新：req.URL.Path = strings.TrimPrefix(req.URL.Path, "/gva")
			req.URL.Path = strings.TrimPrefix(req.URL.Path, "/gva")
// 条件判断：if req.URL.Path == "" {
			if req.URL.Path == "" {
// 赋值/更新：req.URL.Path = "/"
				req.URL.Path = "/"
// 结束当前代码块
			}
// 赋值/更新：req.Host = adminURL.Host
			req.Host = adminURL.Host
// 结束当前代码块
		}}
// 执行：group.Any("/gva/*path", func(c *gin.Context) {
		group.Any("/gva/*path", func(c *gin.Context) {
// 执行：proxy.ServeHTTP(c.Writer, c.Request)
			proxy.ServeHTTP(c.Writer, c.Request)
// 结束当前代码块
		})
// 结束当前代码块
	}

	// AuthGuard 必须在 JWT 之后执行：委派契约实现时需要 claims 判断越权
	authed := group.Group("/api", middleware.JWT(cfg.Jwt.Secret, now), ca.AuthGuard())
// 执行：{
	{
// 执行：authed.GET("/me", ag.Me)
		authed.GET("/me", ag.Me)

		// 买家：下单、查询、取消、售前退款（契约规定买家区仅 buyer 角色可用）
		buyer := authed.Group("", middleware.RequireRole(model.RoleBuyer))
// 执行：{
		{
// 执行：buyer.POST("/orders", ag.CreateOrder)
			buyer.POST("/orders", ag.CreateOrder)
// 执行：buyer.GET("/orders", ag.MyOrders)
			buyer.GET("/orders", ag.MyOrders)
// 执行：buyer.GET("/orders/:orderNo", ag.OrderDetail)
			buyer.GET("/orders/:orderNo", ag.OrderDetail)
// 执行：buyer.POST("/orders/:orderNo/cancel", ag.Cancel)
			buyer.POST("/orders/:orderNo/cancel", ag.Cancel)
// 执行：buyer.POST("/orders/:orderNo/refund", ag.RefundApply)
			buyer.POST("/orders/:orderNo/refund", ag.RefundApply)
// 执行：buyer.GET("/orders/:orderNo/refund", ag.MyRefund)
			buyer.GET("/orders/:orderNo/refund", ag.MyRefund)
// 执行：buyer.POST("/orders/:orderNo/alipay-pay", ag.AlipayPay)
			buyer.POST("/orders/:orderNo/alipay-pay", ag.AlipayPay)
// 结束当前代码块
		}

		// 运营：上下架、核查、查询
		ops := authed.Group("/ops", middleware.RequireRole(model.RoleOps))
// 执行：{
		{
// 执行：ops.PATCH("/products", ag.ProductUpdate)
			ops.PATCH("/products", ag.ProductUpdate)
// 执行：ops.GET("/summary", ag.Summary)
			ops.GET("/summary", ag.Summary)
// 执行：ops.GET("/orders", ag.Orders)
			ops.GET("/orders", ag.Orders)
// 执行：ops.GET("/orders/:orderNo", ag.OrderTrace)
			ops.GET("/orders/:orderNo", ag.OrderTrace)
// 执行：ops.GET("/notices", ag.Notices)
			ops.GET("/notices", ag.Notices)
// 执行：ops.GET("/abnormals", ag.Abnormals)
			ops.GET("/abnormals", ag.Abnormals)
// 执行：ops.GET("/ledger", ag.Ledger)
			ops.GET("/ledger", ag.Ledger)
// 执行：ops.GET("/trail", ag.Trail)
			ops.GET("/trail", ag.Trail)
// 结束当前代码块
		}

		// 仓库：待发货与发货登记
		wh := authed.Group("/warehouse", middleware.RequireRole(model.RoleWarehouse))
// 执行：{
		{
// 执行：wh.GET("/shipable", ag.Shipable)
			wh.GET("/shipable", ag.Shipable)
// 执行：wh.POST("/ship", ag.Ship)
			wh.POST("/ship", ag.Ship)
// 执行：wh.GET("/shipments", ag.Shipments)
			wh.GET("/shipments", ag.Shipments)
// 执行：wh.GET("/orders/:orderNo", ag.OrderTrace)
			wh.GET("/orders/:orderNo", ag.OrderTrace)
// 结束当前代码块
		}

		// 客服：退款审核、渠道重试、异常收款处理
		cs := authed.Group("/cs", middleware.RequireRole(model.RoleCS))
// 执行：{
		{
// 执行：cs.GET("/refunds", ag.Refunds)
			cs.GET("/refunds", ag.Refunds)
// 执行：cs.POST("/refunds/audit", ag.RefundAudit)
			cs.POST("/refunds/audit", ag.RefundAudit)
// 执行：cs.POST("/refunds/retry", ag.RefundRetry)
			cs.POST("/refunds/retry", ag.RefundRetry)
// 执行：cs.GET("/abnormals", ag.Abnormals)
			cs.GET("/abnormals", ag.Abnormals)
// 执行：cs.POST("/abnormals/handle", ag.HandleAbnormal)
			cs.POST("/abnormals/handle", ag.HandleAbnormal)
// 执行：cs.GET("/notices", ag.Notices)
			cs.GET("/notices", ag.Notices)
// 执行：cs.GET("/orders/:orderNo", ag.OrderTrace)
			cs.GET("/orders/:orderNo", ag.OrderTrace)
// 执行：cs.POST("/channel/mode", ag.ChannelMode)
			cs.POST("/channel/mode", ag.ChannelMode)
// 结束当前代码块
		}
// 结束当前代码块
	}

	// 契约新增路径（不与原型重名的部分）
	InitCompatRouter(pub, authed, ca, cfg)
// 结束当前代码块
}
