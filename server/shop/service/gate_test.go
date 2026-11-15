// 包声明
package service

// 导入依赖包
import (
// 执行："errors"
	"errors"
// 执行："sync"
	"sync"
// 执行："testing"
	"testing"
// 执行："time"
	"time"

// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// I-18 并发防御回归：渠道退款调用必须受闸门保护，绝不重复发起。

// TestConcurrentRetrySingleChannelCall 渠道失败的退款单并发重试：
// 无论多少协程同时 Retry，渠道受理恰好一次，其余请求被闸门拒绝。
func TestConcurrentRetrySingleChannelCall(t *testing.T) {
// 赋值/更新：g, ch, _ := newEnv(t)
	g, ch, _ := newEnv(t)
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：cs := buyerID(t, "cs1")
	cs := buyerID(t, "cs1")

// 赋值/更新：o, err := g.OrderService.Create(b1, pid, 1)
	o, err := g.OrderService.Create(b1, pid, 1)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("下单失败: %v", err)
		t.Fatalf("下单失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, "TXN-G1", o.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, "TXN-G1", o.AmountCents)
// 条件判断：if _, err := g.PaymentService.Notify(o.OrderNo, "TXN-G1", o.AmountCents, sign, "{}"); err != nil {
	if _, err := g.PaymentService.Notify(o.OrderNo, "TXN-G1", o.AmountCents, sign, "{}"); err != nil {
// 执行：t.Fatalf("支付通知失败: %v", err)
		t.Fatalf("支付通知失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.RefundService.Apply(b1, o.OrderNo, "不想要了"); err != nil {
	if _, err := g.RefundService.Apply(b1, o.OrderNo, "不想要了"); err != nil {
// 执行：t.Fatalf("申请退款失败: %v", err)
		t.Fatalf("申请退款失败: %v", err)
// 结束当前代码块
	}

// 执行：ch.SetRefundMode("fail")
	ch.SetRefundMode("fail")
// 条件判断：if _, err := g.RefundService.Audit(cs, o.OrderNo, true, "同意退款"); !errors.Is(err, ErrChannelRefund) {
	if _, err := g.RefundService.Audit(cs, o.OrderNo, true, "同意退款"); !errors.Is(err, ErrChannelRefund) {
// 执行：t.Fatalf("首次退款应渠道失败: %v", err)
		t.Fatalf("首次退款应渠道失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：firstCalls := ch.RefundCalls()
	firstCalls := ch.RefundCalls()
// 条件判断：if firstCalls != 1 {
	if firstCalls != 1 {
// 执行：t.Fatalf("首次退款应只调渠道 1 次: %d", firstCalls)
		t.Fatalf("首次退款应只调渠道 1 次: %d", firstCalls)
// 结束当前代码块
	}
// 赋值/更新：rfNo := mustRefundNo(t, g, o.OrderNo)
	rfNo := mustRefundNo(t, g, o.OrderNo)

	// 切换渠道成功后 8 协程并发重试
	ch.SetRefundMode("success")
// 声明：const workers = 8
	const workers = 8
// 赋值/更新：errs := make([]error, workers)
	errs := make([]error, workers)
// 赋值/更新：start := make(chan struct{})
	start := make(chan struct{})
// 声明：var wg sync.WaitGroup
	var wg sync.WaitGroup
// 循环处理：for i := 0; i < workers; i++ {
	for i := 0; i < workers; i++ {
// 执行：wg.Add(1)
		wg.Add(1)
// 异步/延迟执行：go func(i int) {
		go func(i int) {
// 异步/延迟执行：defer wg.Done()
			defer wg.Done()
// 执行：<-start
			<-start
// 赋值/更新：_, errs[i] = g.RefundService.Retry(cs, rfNo)
			_, errs[i] = g.RefundService.Retry(cs, rfNo)
// 结束当前代码块
		}(i)
// 结束当前代码块
	}
// 执行：close(start)
	close(start)
// 执行：wg.Wait()
	wg.Wait()

// 赋值/更新：okCount, gateCount := 0, 0
	okCount, gateCount := 0, 0
// 循环处理：for _, e := range errs {
	for _, e := range errs {
// 多路分支：switch {
		switch {
// 匹配分支：case e == nil:
		case e == nil:
// 执行：okCount++
			okCount++
// 匹配分支：case errors.Is(e, ErrRefundNotRetryable):
		case errors.Is(e, ErrRefundNotRetryable):
// 执行：gateCount++
			gateCount++
// 默认分支
		default:
// 执行：t.Fatalf("意外错误: %v", e)
			t.Fatalf("意外错误: %v", e)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if okCount != 1 || gateCount != workers-1 {
	if okCount != 1 || gateCount != workers-1 {
// 执行：t.Fatalf("应恰好 1 个重试成功、%d 个被闸门拒绝，实际 ok=%d gate=%d", workers-1, okCount, gateCount)
		t.Fatalf("应恰好 1 个重试成功、%d 个被闸门拒绝，实际 ok=%d gate=%d", workers-1, okCount, gateCount)
// 结束当前代码块
	}
// 条件判断：if ch.RefundCalls() != firstCalls+1 {
	if ch.RefundCalls() != firstCalls+1 {
// 执行：t.Fatalf("渠道受理应只增加 1 次（%d -> 应为 %d，实际 %d）", firstCalls, firstCalls+1, ch.RefundCalls())
		t.Fatalf("渠道受理应只增加 1 次（%d -> 应为 %d，实际 %d）", firstCalls, firstCalls+1, ch.RefundCalls())
// 结束当前代码块
	}

	// 订单已退款、售前退款恢复库存且只恢复一次
	o2, _ := g.OrderService.Get(o.OrderNo)
// 条件判断：if o2.Status != model.OrderStatusRefunded {
	if o2.Status != model.OrderStatusRefunded {
// 执行：t.Fatalf("订单应为已退款: %s", o2.Status)
		t.Fatalf("订单应为已退款: %s", o2.Status)
// 结束当前代码块
	}
// 声明：var restores int64
	var restores int64
// 执行：testDB.Table("inventory_ledger").Where("order_no = ? AND change_type = ?", o.OrderNo, model.LedgerRestore).Cou...
	testDB.Table("inventory_ledger").Where("order_no = ? AND change_type = ?", o.OrderNo, model.LedgerRestore).Count(&restores)
// 条件判断：if restores != 1 {
	if restores != 1 {
// 执行：t.Fatalf("库存恢复流水应恰好 1 条: %d", restores)
		t.Fatalf("库存恢复流水应恰好 1 条: %d", restores)
// 结束当前代码块
	}
// 结束当前代码块
}

// TestConcurrentAbnormalRefundSingleChannelCall 迟到支付异常退款并发处理：
// 同一交易号两个协程并发 HandleAbnormal，渠道受理恰好一次，后到者报"已处理"。
func TestConcurrentAbnormalRefundSingleChannelCall(t *testing.T) {
// 赋值/更新：g, ch, clock := newEnv(t)
	g, ch, clock := newEnv(t)
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：cs := buyerID(t, "cs1")
	cs := buyerID(t, "cs1")

// 赋值/更新：o, _ := g.OrderService.Create(b1, pid, 1)
	o, _ := g.OrderService.Create(b1, pid, 1)
// 执行：clock.Advance(11 * time.Minute)
	clock.Advance(11 * time.Minute)
// 条件判断：if _, err := g.OrderService.CloseExpired(); err != nil {
	if _, err := g.OrderService.CloseExpired(); err != nil {
// 执行：t.Fatalf("超时关单失败: %v", err)
		t.Fatalf("超时关单失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, "TXN-G2", o.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, "TXN-G2", o.AmountCents)
// 条件判断：if _, err := g.PaymentService.Notify(o.OrderNo, "TXN-G2", o.AmountCents, sign, "{}"); err != nil {
	if _, err := g.PaymentService.Notify(o.OrderNo, "TXN-G2", o.AmountCents, sign, "{}"); err != nil {
// 执行：t.Fatalf("迟到通知应受理为异常: %v", err)
		t.Fatalf("迟到通知应受理为异常: %v", err)
// 结束当前代码块
	}

// 执行：ch.SetRefundMode("success")
	ch.SetRefundMode("success")
// 赋值/更新：before := ch.RefundCalls()
	before := ch.RefundCalls()
// 赋值/更新：errs := make([]error, 2)
	errs := make([]error, 2)
// 声明：var wg sync.WaitGroup
	var wg sync.WaitGroup
// 赋值/更新：start := make(chan struct{})
	start := make(chan struct{})
// 循环处理：for i := 0; i < 2; i++ {
	for i := 0; i < 2; i++ {
// 执行：wg.Add(1)
		wg.Add(1)
// 异步/延迟执行：go func(i int) {
		go func(i int) {
// 异步/延迟执行：defer wg.Done()
			defer wg.Done()
// 执行：<-start
			<-start
// 赋值/更新：_, errs[i] = g.RefundService.HandleAbnormal(cs, "TXN-G2", "核实后退款")
			_, errs[i] = g.RefundService.HandleAbnormal(cs, "TXN-G2", "核实后退款")
// 结束当前代码块
		}(i)
// 结束当前代码块
	}
// 执行：close(start)
	close(start)
// 执行：wg.Wait()
	wg.Wait()

// 条件判断：if ch.RefundCalls() != before+1 {
	if ch.RefundCalls() != before+1 {
// 执行：t.Fatalf("渠道受理应只增加 1 次（%d -> 应为 %d，实际 %d）", before, before+1, ch.RefundCalls())
		t.Fatalf("渠道受理应只增加 1 次（%d -> 应为 %d，实际 %d）", before, before+1, ch.RefundCalls())
// 结束当前代码块
	}
// 赋值/更新：okCount := 0
	okCount := 0
// 循环处理：for _, e := range errs {
	for _, e := range errs {
// 条件判断：if e == nil {
		if e == nil {
// 执行：okCount++
			okCount++
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if okCount != 1 {
	if okCount != 1 {
// 执行：t.Fatalf("应恰好一个处理成功: %d (%v / %v)", okCount, errs[0], errs[1])
		t.Fatalf("应恰好一个处理成功: %d (%v / %v)", okCount, errs[0], errs[1])
// 结束当前代码块
	}
// 声明：var handled int64
	var handled int64
// 执行：testDB.Model(&model.AbnormalPayment{}).Where("channel_txn_id = ? AND status = ?", "TXN-G2", model.AbnormalHand...
	testDB.Model(&model.AbnormalPayment{}).Where("channel_txn_id = ? AND status = ?", "TXN-G2", model.AbnormalHandled).Count(&handled)
// 条件判断：if handled != 1 {
	if handled != 1 {
// 执行：t.Fatalf("异常记录应恰好标记已处理 1 条: %d", handled)
		t.Fatalf("异常记录应恰好标记已处理 1 条: %d", handled)
// 结束当前代码块
	}
// 结束当前代码块
}

// mustRefundNo 按订单号查询退款单号
func mustRefundNo(t *testing.T, g *ServiceGroup, orderNo string) string {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：rf, err := g.RefundService.Get(orderNo)
	rf, err := g.RefundService.Get(orderNo)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("查询退款单失败: %v", err)
		t.Fatalf("查询退款单失败: %v", err)
// 结束当前代码块
	}
// 返回：return rf.RefundNo
	return rf.RefundNo
// 结束当前代码块
}
