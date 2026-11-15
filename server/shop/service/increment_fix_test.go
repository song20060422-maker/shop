package service

import (
	"errors"
	"testing"
	"time"

	"shop/shop/model"
	"shop/shop/model/response"
	"shop/shop/utils"
)

// 增量续写配套测试：
// 1. 异常退款单（abnormal_late）渠道失败后按退款单号重试成功（R4）；
// 2. 汇总口径：运营补货后占用量仍守恒、异常收款与有效收款分离（R5）。

// TestAbnormalRefundRetry 渠道首次失败 -> 客服按退款单号重试成功：
// 异常记录转 handled、订单保持 closed、可售数量不恢复、无库存恢复流水。
func TestAbnormalRefundRetry(t *testing.T) {
	g, ch, clock := newEnv(t)
	b1 := buyerID(t, "buyer1")
	pid := productID(t, "SKU-A")
	ops := buyerID(t, "ops")

	// 下单 -> 超时关闭（释放占用）
	stockBeforeOrder := available(t, "SKU-A")
	o, err := g.OrderService.Create(b1, pid, 1)
	if err != nil {
		t.Fatal(err)
	}
	clock.Advance(11 * time.Minute)
	if n, _ := g.OrderService.CloseExpired(); n != 1 {
		t.Fatalf("应关闭 1 笔订单，实际 %d", n)
	}
	if available(t, "SKU-A") != stockBeforeOrder {
		t.Fatal("关闭订单应恰好释放一次占用")
	}

	// 迟到支付到账 -> 异常收款
	txn := "TXN-LATE-" + o.OrderNo
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, txn, o.AmountCents)
	res, err := g.PaymentService.Notify(o.OrderNo, txn, o.AmountCents, sign, "{}")
	if err != nil || res.Result != model.NoticeAbnormal {
		t.Fatalf("迟到支付应转异常: err=%v res=%+v", err, res)
	}

	// 渠道首次失败：异常记录保留 pending，退款单落 failed
	ch.SetRefundMode("fail")
	rf, err := g.RefundService.HandleAbnormal(ops, txn, "核实迟到支付，同意退款")
	if !errors.Is(err, ErrChannelRefund) {
		t.Fatalf("渠道失败应返回 ErrChannelRefund: %v", err)
	}
	if rf.ChannelStatus != model.ChannelFailed {
		t.Fatalf("退款单应落 failed: %s", rf.ChannelStatus)
	}

	// 渠道恢复后按退款单号重试成功（修复点 R4：原逻辑此处必然 3010 回滚）
	ch.SetRefundMode("success")
	rf, err = g.RefundService.Retry(ops, rf.RefundNo)
	if err != nil {
		t.Fatalf("异常退款单重试应成功（R4 修复点）: %v", err)
	}
	if rf.ChannelStatus != model.ChannelSuccess {
		t.Fatalf("重试后应为 success: %s", rf.ChannelStatus)
	}

	// 断言：异常记录已处理、订单保持关闭、库存未恢复、无 restore 流水
	var ab model.AbnormalPayment
	if err := testDB.Where("channel_txn_id = ?", txn).First(&ab).Error; err != nil {
		t.Fatal(err)
	}
	if ab.Status != model.AbnormalHandled || ab.RefundID == nil || *ab.RefundID != rf.ID {
		t.Fatalf("异常记录应转 handled 且回链退款单: %+v", ab)
	}
	var o2 model.Order
	if err := testDB.Where("order_no = ?", o.OrderNo).First(&o2).Error; err != nil {
		t.Fatal(err)
	}
	if o2.Status != model.OrderStatusClosed {
		t.Fatalf("异常退款不得改订单状态，应为 closed: %s", o2.Status)
	}
	if got := available(t, "SKU-A"); got != stockBeforeOrder {
		t.Fatalf("迟到支付退款禁止恢复库存: want %d got %d", stockBeforeOrder, got)
	}
	var restoreCount int64
	testDB.Model(&model.InventoryLedger{}).
		Where("order_no = ? AND change_type = ?", o.OrderNo, model.LedgerRestore).Count(&restoreCount)
	if restoreCount != 0 {
		t.Fatalf("异常退款不得产生 restore 流水，实际 %d 条", restoreCount)
	}

	// 已成功的异常退款单不可再次重试（防重复退款）
	if _, err := g.RefundService.Retry(ops, rf.RefundNo); !errors.Is(err, ErrRefundNotRetryable) {
		t.Fatalf("已成功退款单重试应被闸门拒绝: %v", err)
	}
}

// TestAbnormalRefundRetryFailsAgain 重试仍失败时异常记录保持 pending，可继续重试。
func TestAbnormalRefundRetryFailsAgain(t *testing.T) {
	g, ch, clock := newEnv(t)
	b1 := buyerID(t, "buyer1")
	pid := productID(t, "SKU-A")
	ops := buyerID(t, "ops")

	o, err := g.OrderService.Create(b1, pid, 1)
	if err != nil {
		t.Fatal(err)
	}
	clock.Advance(11 * time.Minute)
	if _, err := g.OrderService.CloseExpired(); err != nil {
		t.Fatal(err)
	}

	txn := "TXN-LATE-2-" + o.OrderNo
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, txn, o.AmountCents)
	if _, err := g.PaymentService.Notify(o.OrderNo, txn, o.AmountCents, sign, "{}"); err != nil {
		t.Fatal(err)
	}

	// 两次失败：第一次 HandleAbnormal、第二次 Retry，异常记录始终 pending
	ch.SetRefundMode("fail")
	rf, err := g.RefundService.HandleAbnormal(ops, txn, "核实迟到支付")
	if !errors.Is(err, ErrChannelRefund) {
		t.Fatalf("首次渠道失败: %v", err)
	}
	rf, err = g.RefundService.Retry(ops, rf.RefundNo)
	if !errors.Is(err, ErrChannelRefund) {
		t.Fatalf("重试仍失败应返回 ErrChannelRefund: %v", err)
	}
	// retry_count 只统计 callChannel 路径；HandleAbnormal 首次失败不计数（既有口径）
	if rf.RetryCount < 1 {
		t.Fatalf("重试次数应累计: %d", rf.RetryCount)
	}
	var ab model.AbnormalPayment
	if err := testDB.Where("channel_txn_id = ?", txn).First(&ab).Error; err != nil {
		t.Fatal(err)
	}
	if ab.Status != model.AbnormalPending {
		t.Fatalf("渠道未成功前异常记录必须保持 pending: %s", ab.Status)
	}

	// 渠道成功后重试闭环
	ch.SetRefundMode("success")
	if _, err := g.RefundService.Retry(ops, rf.RefundNo); err != nil {
		t.Fatalf("二次重试应成功: %v", err)
	}
}

// TestSummaryOccupiedAndAbnormalSeparated 汇总口径：
// ① 运营补货后占用量按订单状态统计不变；② 异常收款独立于有效收款。
func TestSummaryOccupiedAndAbnormalSeparated(t *testing.T) {
	g, _, clock := newEnv(t)
	b1 := buyerID(t, "buyer1")
	pid := productID(t, "SKU-A")
	ops := buyerID(t, "ops1")

	// 一笔待付款订单：占用 1 件
	if _, err := g.OrderService.Create(b1, pid, 1); err != nil {
		t.Fatal(err)
	}
	// 一笔已支付订单：占用 1 件
	payOrder(t, g, b1, 1)
	// 再来一笔待付款订单（稍后超时关闭后迟到支付，制造异常收款）
	late, err := g.OrderService.Create(b1, pid, 1)
	if err != nil {
		t.Fatal(err)
	}

	// 运营把可售数量从 3 补到 8（OverwriteAvailable +5）
	if _, err := g.CatalogService.OverwriteAvailable(ops, model.RoleOps, pid, 8); err != nil {
		t.Fatal(err)
	}

	sum, err := g.ReportService.Summary()
	if err != nil {
		t.Fatal(err)
	}
	var stat *response.ProductStat
	for i := range sum.Products {
		if sum.Products[i].ProductID == pid {
			stat = &sum.Products[i]
		}
	}
	if stat == nil {
		t.Fatal("汇总缺少目标商品")
	}
	// 占用量按状态直查：待付款 1 + 已支付 1 + 迟到单待付款 1 = 3（补货不虚增占用量）
	if stat.Occupied != 3 {
		t.Fatalf("补货后占用量应仍为 3（R5 修复点，原公式会算成 8）: %d", stat.Occupied)
	}
	// 有效收款 1 单 10000 分，异常收款为 0
	if stat.PaidCents != 10000 {
		t.Fatalf("有效收款应为 10000: %d", stat.PaidCents)
	}
	if stat.AbnormalCents != 0 {
		t.Fatalf("异常收款应为 0: %d", stat.AbnormalCents)
	}

	// 制造一笔迟到支付异常收款：两笔待付款单到期关闭后，对其中一笔发起迟到支付通知
	clock.Advance(11 * time.Minute)
	if n, _ := g.OrderService.CloseExpired(); n != 2 {
		t.Fatalf("应关闭 2 笔待付款订单，实际 %d", n)
	}
	lateTxn := "TXN-SUM-LATE-" + late.OrderNo
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, late.OrderNo, lateTxn, late.AmountCents)
	res, err := g.PaymentService.Notify(late.OrderNo, lateTxn, late.AmountCents, sign, "{}")
	if err != nil || res.Result != model.NoticeAbnormal {
		t.Fatalf("应产生异常收款: err=%v res=%+v", err, res)
	}

	sum2, err := g.ReportService.Summary()
	if err != nil {
		t.Fatal(err)
	}
	if sum2.TotalPaidCents != 10000 {
		t.Fatalf("有效收款不得并入异常金额: %d", sum2.TotalPaidCents)
	}
	if sum2.TotalAbnormalCents != late.AmountCents {
		t.Fatalf("异常收款应独立合计 %d（R5 修复点）: %d", late.AmountCents, sum2.TotalAbnormalCents)
	}
	var lateStat *response.ProductStat
	for i := range sum2.Products {
		if sum2.Products[i].ProductID == pid {
			lateStat = &sum2.Products[i]
		}
	}
	if lateStat == nil || lateStat.AbnormalCents != late.AmountCents {
		t.Fatalf("商品维度异常收款应为 %d", late.AmountCents)
	}
}
