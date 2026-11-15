package service

import (
	"log"
	"time"
)

// StartCloseExpiredLoop 超时关单后台调度（新增于增量续写【修改原因=R6】）：
// 生产环境没有 demo-only 手动钩子，必须由定时任务关闭到期待付款订单并释放占用（PROD-10）。
// CloseExpired 内部为「条件更新 + 行数校验」，与买家取消/支付通知并发时天然幂等，
// 重复触发不会二次释放库存，因此调度器只管节奏，不做额外加锁。
// stop 在服务优雅退出时关闭；interval <= 0 由 Config.CloseExpiredInterval 兜底默认值。
//
// 买家端（cmd/shop）与 GVA 管理端（initialize.Boot）启动时各拉起一个实例即可，
// 两个实例并发执行也安全（同一订单只会被其中一个成功关闭）。

func (g *ServiceGroup) StartCloseExpiredLoop(interval time.Duration, stop <-chan struct{}) {
	if interval <= 0 {
		interval = g.OrderService.deps.Cfg.CloseExpiredInterval()
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			// 启动即执行一次，补上停机期间到期的订单
			g.closeExpiredQuiet()
			select {
			case <-stop:
				return
			case <-ticker.C:
				g.closeExpiredQuiet()
			}
		}
	}()
}

// closeExpiredQuiet 调度轮次：关闭失败只记日志，不允许 panic 拖垮调度循环。

func (g *ServiceGroup) closeExpiredQuiet() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("shop: 超时关单调度异常（已忽略，下一轮重试）: %v", r)
		}
	}()
	n, err := g.OrderService.CloseExpired()
	if err != nil {
		log.Printf("shop: 超时关单执行失败（下一轮重试）: %v", err)
		return
	}
	if n > 0 {
		log.Printf("shop: 超时关单调度关闭 %d 笔到期订单并释放占用", n)
	}
}
