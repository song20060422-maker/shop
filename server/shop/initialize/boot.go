// 包声明
package initialize

// 导入依赖包
import (
// 执行："time"
	"time"

// 执行："shop/global"
	"shop/global"
// 执行："shop/shop/config"
	"shop/shop/config"
// 执行：shopdb "shop/shop/db"
	shopdb "shop/shop/db"
// 执行："shop/shop/migrations"
	"shop/shop/migrations"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// Boot 初始化 Shop 业务模块：读取 shop 配置段 -> 连接数据库 -> 执行迁移 -> 补齐种子账号。
// 返回装配完成的 ServiceGroup，供路由层直接消费。
// 注意：Shop 使用独立连接（Deps.DB 注入），严禁覆盖 global.GVA_DB——
// 它是 Gin-Vue-Admin 框架主库（gva）的连接，覆盖会导致管理后台查询 shop 库失败。

// 函数定义：func Boot() (*service.ServiceGroup, config.Config, error) {
func Boot() (*service.ServiceGroup, config.Config, error) {
// 赋值/更新：cfg := config.FromViper(global.GVA_VP)
	cfg := config.FromViper(global.GVA_VP)

// 赋值/更新：db, err := shopdb.Gorm(cfg)
	db, err := shopdb.Gorm(cfg)
// 条件判断：if err != nil {
	if err != nil {
// 返回：return nil, cfg, err
		return nil, cfg, err
// 结束当前代码块
	}
// 条件判断：if err := shopdb.Migrate(db, migrations.FS); err != nil {
	if err := shopdb.Migrate(db, migrations.FS); err != nil {
// 返回：return nil, cfg, err
		return nil, cfg, err
// 结束当前代码块
	}
// 条件判断：if err := shopdb.SeedAccounts(db, cfg.Shop.DemoPassword); err != nil {
	if err := shopdb.SeedAccounts(db, cfg.Shop.DemoPassword); err != nil {
// 返回：return nil, cfg, err
		return nil, cfg, err
// 结束当前代码块
	}

// 赋值/更新：now := time.Now
	now := time.Now
	// 渠道装配：pay-channel=alipay 时管理端退款走支付宝沙箱真实渠道；失败即启动失败，不静默降级
	ch, err := service.NewChannelFromConfig(cfg)
// 条件判断：if err != nil {
	if err != nil {
// 返回：return nil, cfg, err
		return nil, cfg, err
// 结束当前代码块
	}
// 赋值/更新：group := service.NewServiceGroup(service.Deps{DB: db, Now: now, Cfg: cfg, Channel: ch})
	group := service.NewServiceGroup(service.Deps{DB: db, Now: now, Cfg: cfg, Channel: ch})
	// 超时关单后台调度：管理端进程同样兜底关闭到期待付款订单（PROD-10，增量续写）。
	// CloseExpired 幂等，与买家端调度器并发执行安全；stop 通道随进程生命周期存活。
	stopScheduler := make(chan struct{})
// 执行：group.StartCloseExpiredLoop(cfg.CloseExpiredInterval(), stopScheduler)
	group.StartCloseExpiredLoop(cfg.CloseExpiredInterval(), stopScheduler)
// 条件判断：if global.GVA_LOG != nil {
	if global.GVA_LOG != nil {
// 条件判断：if _, isAli := group.Channel.(*service.AlipayChannel); isAli {
		if _, isAli := group.Channel.(*service.AlipayChannel); isAli {
// 执行：global.GVA_LOG.Info("shop 业务模块初始化完成（支付渠道：支付宝沙箱）")
			global.GVA_LOG.Info("shop 业务模块初始化完成（支付渠道：支付宝沙箱）")
// 条件分支：} else {
		} else {
// 执行：global.GVA_LOG.Info("shop 业务模块初始化完成（支付渠道：模拟渠道）")
			global.GVA_LOG.Info("shop 业务模块初始化完成（支付渠道：模拟渠道）")
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return group, cfg, nil
	return group, cfg, nil
// 结束当前代码块
}
