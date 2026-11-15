// 独立 Gin 服务：Shop 商城买家端 + 模拟支付渠道（非后台管理页面）。
//
// 与 GVA 管理端（server/main.go，默认 8099）共用同一个 gva 数据库：
//   - 本服务承载买家与渠道侧页面/接口（含接口契约路径，供验收脚本使用）
//   - 运营/仓库/客服后台页面与接口由 GVA 管理端提供（gva 登录 + Casbin）
//
// 启动：go run ./cmd/shop -c config.shop.yaml
package main

// 导入依赖包
import (
// 执行："flag"
	"flag"
// 执行："log"
	"log"
// 执行："time"
	"time"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行：shopconfig "shop/shop/config"
	shopconfig "shop/shop/config"
// 执行：shopdb "shop/shop/db"
	shopdb "shop/shop/db"
// 执行："shop/shop/migrations"
	"shop/shop/migrations"
// 执行：shoprouter "shop/shop/router"
	shoprouter "shop/shop/router"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// 函数定义：func main() {
func main() {
// 赋值/更新：cfgPath := flag.String("c", "", "配置文件路径（缺省用内置默认值，可用 SHOP_DSN/SHOP_ADDR/SHOP_MODE 等环境变量覆盖）")
	cfgPath := flag.String("c", "", "配置文件路径（缺省用内置默认值，可用 SHOP_DSN/SHOP_ADDR/SHOP_MODE 等环境变量覆盖）")
// 执行：flag.Parse()
	flag.Parse()

// 赋值/更新：cfg, err := shopconfig.Load(*cfgPath)
	cfg, err := shopconfig.Load(*cfgPath)
// 条件判断：if err != nil {
	if err != nil {
// 执行：log.Fatalf("读取配置失败: %v", err)
		log.Fatalf("读取配置失败: %v", err)
// 结束当前代码块
	}

// 赋值/更新：db, err := shopdb.Gorm(cfg)
	db, err := shopdb.Gorm(cfg)
// 条件判断：if err != nil {
	if err != nil {
// 执行：log.Fatalf("连接数据库失败: %v", err)
		log.Fatalf("连接数据库失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := shopdb.Migrate(db, migrations.FS); err != nil {
	if err := shopdb.Migrate(db, migrations.FS); err != nil {
// 执行：log.Fatalf("执行迁移失败: %v", err)
		log.Fatalf("执行迁移失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := shopdb.SeedAccounts(db, cfg.Shop.DemoPassword); err != nil {
	if err := shopdb.SeedAccounts(db, cfg.Shop.DemoPassword); err != nil {
// 执行：log.Fatalf("种子账号失败: %v", err)
		log.Fatalf("种子账号失败: %v", err)
// 结束当前代码块
	}

	// 渠道装配：pay-channel=alipay 且配置完整时启用支付宝沙箱渠道，否则保持模拟渠道（默认）
	ch, err := service.NewChannelFromConfig(cfg)
// 条件判断：if err != nil {
	if err != nil {
// 执行：log.Fatalf("支付宝渠道初始化失败: %v", err)
		log.Fatalf("支付宝渠道初始化失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：g := service.NewServiceGroup(service.Deps{DB: db, Now: time.Now, Cfg: cfg, Channel: ch})
	g := service.NewServiceGroup(service.Deps{DB: db, Now: time.Now, Cfg: cfg, Channel: ch})
// 条件判断：if _, isAli := g.Channel.(*service.AlipayChannel); isAli {
	if _, isAli := g.Channel.(*service.AlipayChannel); isAli {
// 执行：log.Printf("支付渠道：支付宝沙箱（%s），异步通知 %s", cfg.Shop.Alipay.AppID, cfg.Shop.Alipay.NotifyURL)
		log.Printf("支付渠道：支付宝沙箱（%s），异步通知 %s", cfg.Shop.Alipay.AppID, cfg.Shop.Alipay.NotifyURL)
// 条件分支：} else {
	} else {
// 执行：log.Print("支付渠道：模拟渠道（shop.pay-channel=alipay 可切换）")
		log.Print("支付渠道：模拟渠道（shop.pay-channel=alipay 可切换）")
// 结束当前代码块
	}

// 条件判断：if cfg.System.Mode == "release" {
	if cfg.System.Mode == "release" {
// 执行：gin.SetMode(gin.ReleaseMode)
		gin.SetMode(gin.ReleaseMode)
// 结束当前代码块
	}
	// 超时关单后台调度：生产环境兜底关闭到期待付款订单并释放占用（PROD-10，增量续写）
	stopScheduler := make(chan struct{})
// 异步/延迟执行：defer close(stopScheduler)
	defer close(stopScheduler)
// 执行：g.StartCloseExpiredLoop(cfg.CloseExpiredInterval(), stopScheduler)
	g.StartCloseExpiredLoop(cfg.CloseExpiredInterval(), stopScheduler)

// 赋值/更新：r := gin.New()
	r := gin.New()
// 执行：r.Use(gin.Logger(), gin.Recovery())
	r.Use(gin.Logger(), gin.Recovery())
// 执行：shoprouter.InitShopRouter(r.Group(""), g, cfg, time.Now)
	shoprouter.InitShopRouter(r.Group(""), g, cfg, time.Now)

// 执行：log.Printf("Shop 买家端独立服务启动：http://127.0.0.1%s/shop", cfg.System.Addr)
	log.Printf("Shop 买家端独立服务启动：http://127.0.0.1%s/shop", cfg.System.Addr)
// 条件判断：if err := r.Run(cfg.System.Addr); err != nil {
	if err := r.Run(cfg.System.Addr); err != nil {
// 执行：log.Fatalf("服务启动失败: %v", err)
		log.Fatalf("服务启动失败: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}
