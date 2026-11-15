// 包声明
package service

// 导入依赖包
import (
// 执行："log"
	"log"
// 执行："time"
	"time"
// 结束当前代码块
)

// StartCloseExpiredLoop 超时关单后台调度（新增于增量续写【修改原因=R6】）：
// 生产环境没有 demo-only 手动钩子，必须由定时任务关闭到期待付款订单并释放占用（PROD-10）。
// CloseExpired 内部为「条件更新 + 行数校验」，与买家取消/支付通知并发时天然幂等，
// 重复触发不会二次释放库存，因此调度器只管节奏，不做额外加锁。
// stop 在服务优雅退出时关闭；interval <= 0 由 Config.CloseExpiredInterval 兜底默认值。
//
// 买家端（cmd/shop）与 GVA 管理端（initialize.Boot）启动时各拉起一个实例即可，
// 两个实例并发执行也安全（同一订单只会被其中一个成功关闭）。

// 函数定义：func (g *ServiceGroup) StartCloseExpiredLoop(interval time.Duration, stop <-chan struct{}) {
func (g *ServiceGroup) StartCloseExpiredLoop(interval time.Duration, stop <-chan struct{}) {
// 条件判断：if interval <= 0 {
	if interval <= 0 {
// 赋值/更新：interval = g.OrderService.deps.Cfg.CloseExpiredInterval()
		interval = g.OrderService.deps.Cfg.CloseExpiredInterval()
// 结束当前代码块
	}
// 异步/延迟执行：go func() {
	go func() {
// 赋值/更新：ticker := time.NewTicker(interval)
		ticker := time.NewTicker(interval)
// 异步/延迟执行：defer ticker.Stop()
		defer ticker.Stop()
// 循环处理：for {
		for {
			// 启动即执行一次，补上停机期间到期的订单
			g.closeExpiredQuiet()
// 多路分支：select {
			select {
// 匹配分支：case <-stop:
			case <-stop:
// 返回：return
				return
// 匹配分支：case <-ticker.C:
			case <-ticker.C:
// 执行：g.closeExpiredQuiet()
				g.closeExpiredQuiet()
// 结束当前代码块
			}
// 结束当前代码块
		}
// 结束当前代码块
	}()
// 结束当前代码块
}

// closeExpiredQuiet 调度轮次：关闭失败只记日志，不允许 panic 拖垮调度循环。

// 函数定义：func (g *ServiceGroup) closeExpiredQuiet() {
func (g *ServiceGroup) closeExpiredQuiet() {
// 异步/延迟执行：defer func() {
	defer func() {
// 条件判断：if r := recover(); r != nil {
		if r := recover(); r != nil {
// 执行：log.Printf("shop: 超时关单调度异常（已忽略，下一轮重试）: %v", r)
			log.Printf("shop: 超时关单调度异常（已忽略，下一轮重试）: %v", r)
// 结束当前代码块
		}
// 结束当前代码块
	}()
// 赋值/更新：n, err := g.OrderService.CloseExpired()
	n, err := g.OrderService.CloseExpired()
// 条件判断：if err != nil {
	if err != nil {
// 执行：log.Printf("shop: 超时关单执行失败（下一轮重试）: %v", err)
		log.Printf("shop: 超时关单执行失败（下一轮重试）: %v", err)
// 返回：return
		return
// 结束当前代码块
	}
// 条件判断：if n > 0 {
	if n > 0 {
// 执行：log.Printf("shop: 超时关单调度关闭 %d 笔到期订单并释放占用", n)
		log.Printf("shop: 超时关单调度关闭 %d 笔到期订单并释放占用", n)
// 结束当前代码块
	}
// 结束当前代码块
}
