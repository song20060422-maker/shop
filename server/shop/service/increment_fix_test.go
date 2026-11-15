// 包声明
package service

// 导入依赖包
import (
// 执行："errors"
	"errors"
// 执行："testing"
	"testing"
// 执行："time"
	"time"

// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/model/response"
	"shop/shop/model/response"
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// 增量续写配套测试：
// 1. 异常退款单（abnormal_late）渠道失败后按退款单号重试成功（R4）；
// 2. 汇总口径：运营补货后占用量仍守恒、异常收款与有效收款分离（R5）。

// TestAbnormalRefundRetry 渠道首次失败 -> 客服按退款单号重试成功：
// 异常记录转 handled、订单保持 closed、可售数量不恢复、无库存恢复流水。
func TestAbnormalRefundRetry(t *testing.T) {
// 赋值/更新：g, ch, clock := newEnv(t)
	g, ch, clock := newEnv(t)
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：ops := buyerID(t, "ops")
	ops := buyerID(t, "ops")

	// 下单 -> 超时关闭（释放占用）
	stockBeforeOrder := available(t, "SKU-A")
// 赋值/更新：o, err := g.OrderService.Create(b1, pid, 1)
	o, err := g.OrderService.Create(b1, pid, 1)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 执行：clock.Advance(11 * time.Minute)
	clock.Advance(11 * time.Minute)
// 条件判断：if n, _ := g.OrderService.CloseExpired(); n != 1 {
	if n, _ := g.OrderService.CloseExpired(); n != 1 {
// 执行：t.Fatalf("应关闭 1 笔订单，实际 %d", n)
		t.Fatalf("应关闭 1 笔订单，实际 %d", n)
// 结束当前代码块
	}
// 条件判断：if available(t, "SKU-A") != stockBeforeOrder {
	if available(t, "SKU-A") != stockBeforeOrder {
// 执行：t.Fatal("关闭订单应恰好释放一次占用")
		t.Fatal("关闭订单应恰好释放一次占用")
// 结束当前代码块
	}

	// 迟到支付到账 -> 异常收款
	txn := "TXN-LATE-" + o.OrderNo
// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, txn, o.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, txn, o.AmountCents)
// 赋值/更新：res, err := g.PaymentService.Notify(o.OrderNo, txn, o.AmountCents, sign, "{}")
	res, err := g.PaymentService.Notify(o.OrderNo, txn, o.AmountCents, sign, "{}")
// 条件判断：if err != nil || res.Result != model.NoticeAbnormal {
	if err != nil || res.Result != model.NoticeAbnormal {
// 执行：t.Fatalf("迟到支付应转异常: err=%v res=%+v", err, res)
		t.Fatalf("迟到支付应转异常: err=%v res=%+v", err, res)
// 结束当前代码块
	}

	// 渠道首次失败：异常记录保留 pending，退款单落 failed
	ch.SetRefundMode("fail")
// 赋值/更新：rf, err := g.RefundService.HandleAbnormal(ops, txn, "核实迟到支付，同意退款")
	rf, err := g.RefundService.HandleAbnormal(ops, txn, "核实迟到支付，同意退款")
// 条件判断：if !errors.Is(err, ErrChannelRefund) {
	if !errors.Is(err, ErrChannelRefund) {
// 执行：t.Fatalf("渠道失败应返回 ErrChannelRefund: %v", err)
		t.Fatalf("渠道失败应返回 ErrChannelRefund: %v", err)
// 结束当前代码块
	}
// 条件判断：if rf.ChannelStatus != model.ChannelFailed {
	if rf.ChannelStatus != model.ChannelFailed {
// 执行：t.Fatalf("退款单应落 failed: %s", rf.ChannelStatus)
		t.Fatalf("退款单应落 failed: %s", rf.ChannelStatus)
// 结束当前代码块
	}

	// 渠道恢复后按退款单号重试成功（修复点 R4：原逻辑此处必然 3010 回滚）
	ch.SetRefundMode("success")
// 赋值/更新：rf, err = g.RefundService.Retry(ops, rf.RefundNo)
	rf, err = g.RefundService.Retry(ops, rf.RefundNo)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("异常退款单重试应成功（R4 修复点）: %v", err)
		t.Fatalf("异常退款单重试应成功（R4 修复点）: %v", err)
// 结束当前代码块
	}
// 条件判断：if rf.ChannelStatus != model.ChannelSuccess {
	if rf.ChannelStatus != model.ChannelSuccess {
// 执行：t.Fatalf("重试后应为 success: %s", rf.ChannelStatus)
		t.Fatalf("重试后应为 success: %s", rf.ChannelStatus)
// 结束当前代码块
	}

	// 断言：异常记录已处理、订单保持关闭、库存未恢复、无 restore 流水
	var ab model.AbnormalPayment
// 条件判断：if err := testDB.Where("channel_txn_id = ?", txn).First(&ab).Error; err != nil {
	if err := testDB.Where("channel_txn_id = ?", txn).First(&ab).Error; err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if ab.Status != model.AbnormalHandled || ab.RefundID == nil || *ab.RefundID != rf.ID {
	if ab.Status != model.AbnormalHandled || ab.RefundID == nil || *ab.RefundID != rf.ID {
// 执行：t.Fatalf("异常记录应转 handled 且回链退款单: %+v", ab)
		t.Fatalf("异常记录应转 handled 且回链退款单: %+v", ab)
// 结束当前代码块
	}
// 声明：var o2 model.Order
	var o2 model.Order
// 条件判断：if err := testDB.Where("order_no = ?", o.OrderNo).First(&o2).Error; err != nil {
	if err := testDB.Where("order_no = ?", o.OrderNo).First(&o2).Error; err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if o2.Status != model.OrderStatusClosed {
	if o2.Status != model.OrderStatusClosed {
// 执行：t.Fatalf("异常退款不得改订单状态，应为 closed: %s", o2.Status)
		t.Fatalf("异常退款不得改订单状态，应为 closed: %s", o2.Status)
// 结束当前代码块
	}
// 条件判断：if got := available(t, "SKU-A"); got != stockBeforeOrder {
	if got := available(t, "SKU-A"); got != stockBeforeOrder {
// 执行：t.Fatalf("迟到支付退款禁止恢复库存: want %d got %d", stockBeforeOrder, got)
		t.Fatalf("迟到支付退款禁止恢复库存: want %d got %d", stockBeforeOrder, got)
// 结束当前代码块
	}
// 声明：var restoreCount int64
	var restoreCount int64
// 执行：testDB.Model(&model.InventoryLedger{}).
	testDB.Model(&model.InventoryLedger{}).
// 执行：Where("order_no = ? AND change_type = ?", o.OrderNo, model.LedgerRestore).Count(&restoreCount)
		Where("order_no = ? AND change_type = ?", o.OrderNo, model.LedgerRestore).Count(&restoreCount)
// 条件判断：if restoreCount != 0 {
	if restoreCount != 0 {
// 执行：t.Fatalf("异常退款不得产生 restore 流水，实际 %d 条", restoreCount)
		t.Fatalf("异常退款不得产生 restore 流水，实际 %d 条", restoreCount)
// 结束当前代码块
	}

	// 已成功的异常退款单不可再次重试（防重复退款）
	if _, err := g.RefundService.Retry(ops, rf.RefundNo); !errors.Is(err, ErrRefundNotRetryable) {
// 执行：t.Fatalf("已成功退款单重试应被闸门拒绝: %v", err)
		t.Fatalf("已成功退款单重试应被闸门拒绝: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// TestAbnormalRefundRetryFailsAgain 重试仍失败时异常记录保持 pending，可继续重试。
func TestAbnormalRefundRetryFailsAgain(t *testing.T) {
// 赋值/更新：g, ch, clock := newEnv(t)
	g, ch, clock := newEnv(t)
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：ops := buyerID(t, "ops")
	ops := buyerID(t, "ops")

// 赋值/更新：o, err := g.OrderService.Create(b1, pid, 1)
	o, err := g.OrderService.Create(b1, pid, 1)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 执行：clock.Advance(11 * time.Minute)
	clock.Advance(11 * time.Minute)
// 条件判断：if _, err := g.OrderService.CloseExpired(); err != nil {
	if _, err := g.OrderService.CloseExpired(); err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}

// 赋值/更新：txn := "TXN-LATE-2-" + o.OrderNo
	txn := "TXN-LATE-2-" + o.OrderNo
// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, txn, o.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, txn, o.AmountCents)
// 条件判断：if _, err := g.PaymentService.Notify(o.OrderNo, txn, o.AmountCents, sign, "{}"); err != nil {
	if _, err := g.PaymentService.Notify(o.OrderNo, txn, o.AmountCents, sign, "{}"); err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}

	// 两次失败：第一次 HandleAbnormal、第二次 Retry，异常记录始终 pending
	ch.SetRefundMode("fail")
// 赋值/更新：rf, err := g.RefundService.HandleAbnormal(ops, txn, "核实迟到支付")
	rf, err := g.RefundService.HandleAbnormal(ops, txn, "核实迟到支付")
// 条件判断：if !errors.Is(err, ErrChannelRefund) {
	if !errors.Is(err, ErrChannelRefund) {
// 执行：t.Fatalf("首次渠道失败: %v", err)
		t.Fatalf("首次渠道失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：rf, err = g.RefundService.Retry(ops, rf.RefundNo)
	rf, err = g.RefundService.Retry(ops, rf.RefundNo)
// 条件判断：if !errors.Is(err, ErrChannelRefund) {
	if !errors.Is(err, ErrChannelRefund) {
// 执行：t.Fatalf("重试仍失败应返回 ErrChannelRefund: %v", err)
		t.Fatalf("重试仍失败应返回 ErrChannelRefund: %v", err)
// 结束当前代码块
	}
	// retry_count 只统计 callChannel 路径；HandleAbnormal 首次失败不计数（既有口径）
	if rf.RetryCount < 1 {
// 执行：t.Fatalf("重试次数应累计: %d", rf.RetryCount)
		t.Fatalf("重试次数应累计: %d", rf.RetryCount)
// 结束当前代码块
	}
// 声明：var ab model.AbnormalPayment
	var ab model.AbnormalPayment
// 条件判断：if err := testDB.Where("channel_txn_id = ?", txn).First(&ab).Error; err != nil {
	if err := testDB.Where("channel_txn_id = ?", txn).First(&ab).Error; err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if ab.Status != model.AbnormalPending {
	if ab.Status != model.AbnormalPending {
// 执行：t.Fatalf("渠道未成功前异常记录必须保持 pending: %s", ab.Status)
		t.Fatalf("渠道未成功前异常记录必须保持 pending: %s", ab.Status)
// 结束当前代码块
	}

	// 渠道成功后重试闭环
	ch.SetRefundMode("success")
// 条件判断：if _, err := g.RefundService.Retry(ops, rf.RefundNo); err != nil {
	if _, err := g.RefundService.Retry(ops, rf.RefundNo); err != nil {
// 执行：t.Fatalf("二次重试应成功: %v", err)
		t.Fatalf("二次重试应成功: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// TestSummaryOccupiedAndAbnormalSeparated 汇总口径：
// ① 运营补货后占用量按订单状态统计不变；② 异常收款独立于有效收款。
func TestSummaryOccupiedAndAbnormalSeparated(t *testing.T) {
// 赋值/更新：g, _, clock := newEnv(t)
	g, _, clock := newEnv(t)
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：ops := buyerID(t, "ops1")
	ops := buyerID(t, "ops1")

	// 一笔待付款订单：占用 1 件
	if _, err := g.OrderService.Create(b1, pid, 1); err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
	// 一笔已支付订单：占用 1 件
	payOrder(t, g, b1, 1)
	// 再来一笔待付款订单（稍后超时关闭后迟到支付，制造异常收款）
	late, err := g.OrderService.Create(b1, pid, 1)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}

	// 运营把可售数量从 3 补到 8（OverwriteAvailable +5）
	if _, err := g.CatalogService.OverwriteAvailable(ops, model.RoleOps, pid, 8); err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}

// 赋值/更新：sum, err := g.ReportService.Summary()
	sum, err := g.ReportService.Summary()
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 声明：var stat *response.ProductStat
	var stat *response.ProductStat
// 循环处理：for i := range sum.Products {
	for i := range sum.Products {
// 条件判断：if sum.Products[i].ProductID == pid {
		if sum.Products[i].ProductID == pid {
// 赋值/更新：stat = &sum.Products[i]
			stat = &sum.Products[i]
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if stat == nil {
	if stat == nil {
// 执行：t.Fatal("汇总缺少目标商品")
		t.Fatal("汇总缺少目标商品")
// 结束当前代码块
	}
	// 占用量按状态直查：待付款 1 + 已支付 1 + 迟到单待付款 1 = 3（补货不虚增占用量）
	if stat.Occupied != 3 {
// 执行：t.Fatalf("补货后占用量应仍为 3（R5 修复点，原公式会算成 8）: %d", stat.Occupied)
		t.Fatalf("补货后占用量应仍为 3（R5 修复点，原公式会算成 8）: %d", stat.Occupied)
// 结束当前代码块
	}
	// 有效收款 1 单 10000 分，异常收款为 0
	if stat.PaidCents != 10000 {
// 执行：t.Fatalf("有效收款应为 10000: %d", stat.PaidCents)
		t.Fatalf("有效收款应为 10000: %d", stat.PaidCents)
// 结束当前代码块
	}
// 条件判断：if stat.AbnormalCents != 0 {
	if stat.AbnormalCents != 0 {
// 执行：t.Fatalf("异常收款应为 0: %d", stat.AbnormalCents)
		t.Fatalf("异常收款应为 0: %d", stat.AbnormalCents)
// 结束当前代码块
	}

	// 制造一笔迟到支付异常收款：两笔待付款单到期关闭后，对其中一笔发起迟到支付通知
	clock.Advance(11 * time.Minute)
// 条件判断：if n, _ := g.OrderService.CloseExpired(); n != 2 {
	if n, _ := g.OrderService.CloseExpired(); n != 2 {
// 执行：t.Fatalf("应关闭 2 笔待付款订单，实际 %d", n)
		t.Fatalf("应关闭 2 笔待付款订单，实际 %d", n)
// 结束当前代码块
	}
// 赋值/更新：lateTxn := "TXN-SUM-LATE-" + late.OrderNo
	lateTxn := "TXN-SUM-LATE-" + late.OrderNo
// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, late.OrderNo, lateTxn, late.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, late.OrderNo, lateTxn, late.AmountCents)
// 赋值/更新：res, err := g.PaymentService.Notify(late.OrderNo, lateTxn, late.AmountCents, sign, "{}")
	res, err := g.PaymentService.Notify(late.OrderNo, lateTxn, late.AmountCents, sign, "{}")
// 条件判断：if err != nil || res.Result != model.NoticeAbnormal {
	if err != nil || res.Result != model.NoticeAbnormal {
// 执行：t.Fatalf("应产生异常收款: err=%v res=%+v", err, res)
		t.Fatalf("应产生异常收款: err=%v res=%+v", err, res)
// 结束当前代码块
	}

// 赋值/更新：sum2, err := g.ReportService.Summary()
	sum2, err := g.ReportService.Summary()
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if sum2.TotalPaidCents != 10000 {
	if sum2.TotalPaidCents != 10000 {
// 执行：t.Fatalf("有效收款不得并入异常金额: %d", sum2.TotalPaidCents)
		t.Fatalf("有效收款不得并入异常金额: %d", sum2.TotalPaidCents)
// 结束当前代码块
	}
// 条件判断：if sum2.TotalAbnormalCents != late.AmountCents {
	if sum2.TotalAbnormalCents != late.AmountCents {
// 执行：t.Fatalf("异常收款应独立合计 %d（R5 修复点）: %d", late.AmountCents, sum2.TotalAbnormalCents)
		t.Fatalf("异常收款应独立合计 %d（R5 修复点）: %d", late.AmountCents, sum2.TotalAbnormalCents)
// 结束当前代码块
	}
// 声明：var lateStat *response.ProductStat
	var lateStat *response.ProductStat
// 循环处理：for i := range sum2.Products {
	for i := range sum2.Products {
// 条件判断：if sum2.Products[i].ProductID == pid {
		if sum2.Products[i].ProductID == pid {
// 赋值/更新：lateStat = &sum2.Products[i]
			lateStat = &sum2.Products[i]
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if lateStat == nil || lateStat.AbnormalCents != late.AmountCents {
	if lateStat == nil || lateStat.AbnormalCents != late.AmountCents {
// 执行：t.Fatalf("商品维度异常收款应为 %d", late.AmountCents)
		t.Fatalf("商品维度异常收款应为 %d", late.AmountCents)
// 结束当前代码块
	}
// 结束当前代码块
}
