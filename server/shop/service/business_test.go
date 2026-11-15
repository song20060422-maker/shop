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

// 函数定义：func TestCreateOrderSnapshotAndOccupy(t *testing.T) {
func TestCreateOrderSnapshotAndOccupy(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")

// 赋值/更新：o, err := g.OrderService.Create(b1, pid, 2)
	o, err := g.OrderService.Create(b1, pid, 2)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("下单失败: %v", err)
		t.Fatalf("下单失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if o.PriceCents != 10000 || o.AmountCents != 20000 || o.Qty != 2 {
	if o.PriceCents != 10000 || o.AmountCents != 20000 || o.Qty != 2 {
// 执行：t.Fatalf("价格快照错误: %+v", o)
		t.Fatalf("价格快照错误: %+v", o)
// 结束当前代码块
	}
// 条件判断：if got := available(t, "SKU-A"); got != 3 {
	if got := available(t, "SKU-A"); got != 3 {
// 执行：t.Fatalf("可售数量应反映占用: %d", got)
		t.Fatalf("可售数量应反映占用: %d", got)
// 结束当前代码块
	}

	// 调价不影响已有订单金额
	newPrice := int64(12000)
// 条件判断：if _, err := g.CatalogService.Update(buyerID(t, "ops1"), "ops", pid, nil, &newPrice); err != nil {
	if _, err := g.CatalogService.Update(buyerID(t, "ops1"), "ops", pid, nil, &newPrice); err != nil {
// 执行：t.Fatalf("调价失败: %v", err)
		t.Fatalf("调价失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：old, _ := g.OrderService.Get(o.OrderNo)
	old, _ := g.OrderService.Get(o.OrderNo)
// 条件判断：if old.AmountCents != 20000 || old.PriceCents != 10000 {
	if old.AmountCents != 20000 || old.PriceCents != 10000 {
// 执行：t.Fatalf("已有订单金额被调价改变: %+v", old)
		t.Fatalf("已有订单金额被调价改变: %+v", old)
// 结束当前代码块
	}
// 赋值/更新：o2, _ := g.OrderService.Create(b1, pid, 1)
	o2, _ := g.OrderService.Create(b1, pid, 1)
// 条件判断：if o2.PriceCents != 12000 {
	if o2.PriceCents != 12000 {
// 执行：t.Fatalf("新订单应使用新价格: %d", o2.PriceCents)
		t.Fatalf("新订单应使用新价格: %d", o2.PriceCents)
// 结束当前代码块
	}

// 声明：var ledgers int64
	var ledgers int64
// 执行：testDB.Model(&model.InventoryLedger{}).Where("order_no IN ?", []string{o.OrderNo, o2.OrderNo}).Count(&ledgers)
	testDB.Model(&model.InventoryLedger{}).Where("order_no IN ?", []string{o.OrderNo, o2.OrderNo}).Count(&ledgers)
// 条件判断：if ledgers != 2 {
	if ledgers != 2 {
// 执行：t.Fatalf("应记录 2 条占用流水: %d", ledgers)
		t.Fatalf("应记录 2 条占用流水: %d", ledgers)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestCreateOrderValidation(t *testing.T) {
func TestCreateOrderValidation(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")

// 条件判断：if _, err := g.OrderService.Create(b1, pid, 0); !errors.Is(err, ErrQtyInvalid) {
	if _, err := g.OrderService.Create(b1, pid, 0); !errors.Is(err, ErrQtyInvalid) {
// 执行：t.Fatalf("数量 0 应拒绝: %v", err)
		t.Fatalf("数量 0 应拒绝: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.OrderService.Create(b1, pid, 3); !errors.Is(err, ErrQtyInvalid) {
	if _, err := g.OrderService.Create(b1, pid, 3); !errors.Is(err, ErrQtyInvalid) {
// 执行：t.Fatalf("数量 3 应拒绝: %v", err)
		t.Fatalf("数量 3 应拒绝: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.OrderService.Create(b1, 99999, 1); !errors.Is(err, ErrProductNotFound) {
	if _, err := g.OrderService.Create(b1, 99999, 1); !errors.Is(err, ErrProductNotFound) {
// 执行：t.Fatalf("未知商品应拒绝: %v", err)
		t.Fatalf("未知商品应拒绝: %v", err)
// 结束当前代码块
	}
// 赋值/更新：off := false
	off := false
// 条件判断：if _, err := g.CatalogService.Update(buyerID(t, "ops1"), "ops", pid, &off, nil); err != nil {
	if _, err := g.CatalogService.Update(buyerID(t, "ops1"), "ops", pid, &off, nil); err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.OrderService.Create(b1, pid, 1); !errors.Is(err, ErrProductOffSale) {
	if _, err := g.OrderService.Create(b1, pid, 1); !errors.Is(err, ErrProductOffSale) {
// 执行：t.Fatalf("下架商品应拒绝: %v", err)
		t.Fatalf("下架商品应拒绝: %v", err)
// 结束当前代码块
	}
// 赋值/更新：on := true
	on := true
// 条件判断：if _, err := g.CatalogService.Update(buyerID(t, "ops1"), "ops", pid, &on, nil); err != nil {
	if _, err := g.CatalogService.Update(buyerID(t, "ops1"), "ops", pid, &on, nil); err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
	// 库存 5 件：2+2+2 的第三单必须失败
	if _, err := g.OrderService.Create(b1, pid, 2); err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.OrderService.Create(b1, pid, 2); err != nil {
	if _, err := g.OrderService.Create(b1, pid, 2); err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.OrderService.Create(b1, pid, 2); !errors.Is(err, ErrInsufficientStock) {
	if _, err := g.OrderService.Create(b1, pid, 2); !errors.Is(err, ErrInsufficientStock) {
// 执行：t.Fatalf("库存不足应拒绝: %v", err)
		t.Fatalf("库存不足应拒绝: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestConcurrentCreateNoOversell(t *testing.T) {
func TestConcurrentCreateNoOversell(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：buyers := []uint64{buyerID(t, "buyer1"), buyerID(t, "buyer2"), buyerID(t, "buyer3")}
	buyers := []uint64{buyerID(t, "buyer1"), buyerID(t, "buyer2"), buyerID(t, "buyer3")}

// 声明：var wg sync.WaitGroup
	var wg sync.WaitGroup
// 声明：var mu sync.Mutex
	var mu sync.Mutex
// 声明：var ok, insufficient, other int
	var ok, insufficient, other int
// 循环处理：for i := 0; i < 20; i++ {
	for i := 0; i < 20; i++ {
// 执行：wg.Add(1)
		wg.Add(1)
// 异步/延迟执行：go func(i int) {
		go func(i int) {
// 异步/延迟执行：defer wg.Done()
			defer wg.Done()
// 赋值/更新：_, err := g.OrderService.Create(buyers[i%3], pid, 1)
			_, err := g.OrderService.Create(buyers[i%3], pid, 1)
// 执行：mu.Lock()
			mu.Lock()
// 异步/延迟执行：defer mu.Unlock()
			defer mu.Unlock()
// 多路分支：switch {
			switch {
// 匹配分支：case err == nil:
			case err == nil:
// 执行：ok++
				ok++
// 匹配分支：case errors.Is(err, ErrInsufficientStock):
			case errors.Is(err, ErrInsufficientStock):
// 执行：insufficient++
				insufficient++
// 默认分支
			default:
// 执行：other++
				other++
// 执行：t.Logf("意外错误: %v", err)
				t.Logf("意外错误: %v", err)
// 结束当前代码块
			}
// 结束当前代码块
		}(i)
// 结束当前代码块
	}
// 执行：wg.Wait()
	wg.Wait()
// 条件判断：if ok != 5 || insufficient != 15 || other != 0 {
	if ok != 5 || insufficient != 15 || other != 0 {
// 执行：t.Fatalf("并发争抢结果错误: 成功=%d 不足=%d 其他=%d", ok, insufficient, other)
		t.Fatalf("并发争抢结果错误: 成功=%d 不足=%d 其他=%d", ok, insufficient, other)
// 结束当前代码块
	}
// 条件判断：if got := available(t, "SKU-A"); got != 0 {
	if got := available(t, "SKU-A"); got != 0 {
// 执行：t.Fatalf("售罄后应为 0: %d", got)
		t.Fatalf("售罄后应为 0: %d", got)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestCancelAndTimeoutIdempotent(t *testing.T) {
func TestCancelAndTimeoutIdempotent(t *testing.T) {
// 赋值/更新：g, _, clock := newEnv(t)
	g, _, clock := newEnv(t)
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")

// 赋值/更新：o, _ := g.OrderService.Create(b1, pid, 2)
	o, _ := g.OrderService.Create(b1, pid, 2)
// 赋值/更新：cancelled, err := g.OrderService.Cancel(b1, o.OrderNo)
	cancelled, err := g.OrderService.Cancel(b1, o.OrderNo)
// 条件判断：if err != nil || cancelled.Status != model.OrderStatusCancelled {
	if err != nil || cancelled.Status != model.OrderStatusCancelled {
// 执行：t.Fatalf("取消失败: %v %+v", err, cancelled)
		t.Fatalf("取消失败: %v %+v", err, cancelled)
// 结束当前代码块
	}
// 条件判断：if got := available(t, "SKU-A"); got != 5 {
	if got := available(t, "SKU-A"); got != 5 {
// 执行：t.Fatalf("取消后应释放为 5: %d", got)
		t.Fatalf("取消后应释放为 5: %d", got)
// 结束当前代码块
	}
// 条件判断：if _, err := g.OrderService.Cancel(b1, o.OrderNo); err == nil {
	if _, err := g.OrderService.Cancel(b1, o.OrderNo); err == nil {
// 执行：t.Fatal("重复取消应报错")
		t.Fatal("重复取消应报错")
// 结束当前代码块
	}
// 条件判断：if got := available(t, "SKU-A"); got != 5 {
	if got := available(t, "SKU-A"); got != 5 {
// 执行：t.Fatalf("重复取消不得重复释放: %d", got)
		t.Fatalf("重复取消不得重复释放: %d", got)
// 结束当前代码块
	}

	// 超时关闭使用服务端时间
	o2, _ := g.OrderService.Create(b1, pid, 2)
// 执行：clock.Advance(11 * time.Minute)
	clock.Advance(11 * time.Minute)
// 赋值/更新：n, err := g.OrderService.CloseExpired()
	n, err := g.OrderService.CloseExpired()
// 条件判断：if err != nil || n != 1 {
	if err != nil || n != 1 {
// 执行：t.Fatalf("超时关闭失败: n=%d err=%v", n, err)
		t.Fatalf("超时关闭失败: n=%d err=%v", n, err)
// 结束当前代码块
	}
// 赋值/更新：got, _ := g.OrderService.Get(o2.OrderNo)
	got, _ := g.OrderService.Get(o2.OrderNo)
// 条件判断：if got.Status != model.OrderStatusClosed || got.ClosedAt == nil {
	if got.Status != model.OrderStatusClosed || got.ClosedAt == nil {
// 执行：t.Fatalf("超时订单状态错误: %+v", got)
		t.Fatalf("超时订单状态错误: %+v", got)
// 结束当前代码块
	}
// 条件判断：if available(t, "SKU-A") != 5 {
	if available(t, "SKU-A") != 5 {
// 执行：t.Fatal("超时关闭应释放占用")
		t.Fatal("超时关闭应释放占用")
// 结束当前代码块
	}
// 条件判断：if n, _ := g.OrderService.CloseExpired(); n != 0 {
	if n, _ := g.OrderService.CloseExpired(); n != 0 {
// 执行：t.Fatalf("重复清扫不得重复关闭: %d", n)
		t.Fatalf("重复清扫不得重复关闭: %d", n)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestNotifyRules(t *testing.T) {
func TestNotifyRules(t *testing.T) {
// 赋值/更新：g, _, clock := newEnv(t)
	g, _, clock := newEnv(t)
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：o, _ := g.OrderService.Create(b1, pid, 1)
	o, _ := g.OrderService.Create(b1, pid, 1)

// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, "TXN-1", o.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, "TXN-1", o.AmountCents)
// 赋值/更新：res, err := g.PaymentService.Notify(o.OrderNo, "TXN-1", o.AmountCents, sign, "{}")
	res, err := g.PaymentService.Notify(o.OrderNo, "TXN-1", o.AmountCents, sign, "{}")
// 条件判断：if err != nil || res.Result != model.NoticeAccepted {
	if err != nil || res.Result != model.NoticeAccepted {
// 执行：t.Fatalf("有效通知应受理: %v %+v", err, res)
		t.Fatalf("有效通知应受理: %v %+v", err, res)
// 结束当前代码块
	}
// 赋值/更新：paid, _ := g.OrderService.Get(o.OrderNo)
	paid, _ := g.OrderService.Get(o.OrderNo)
// 条件判断：if paid.Status != model.OrderStatusPaid || paid.PaidAt == nil {
	if paid.Status != model.OrderStatusPaid || paid.PaidAt == nil {
// 执行：t.Fatalf("订单未进入待发货: %+v", paid)
		t.Fatalf("订单未进入待发货: %+v", paid)
// 结束当前代码块
	}

	// 重复通知：同一交易号，不重复记款
	res, err = g.PaymentService.Notify(o.OrderNo, "TXN-1", o.AmountCents, sign, "{}")
// 条件判断：if !errors.Is(err, ErrDuplicatedPay) || res.Result != model.NoticeDuplicated {
	if !errors.Is(err, ErrDuplicatedPay) || res.Result != model.NoticeDuplicated {
// 执行：t.Fatalf("重复通知应拒绝记款: %v %+v", err, res)
		t.Fatalf("重复通知应拒绝记款: %v %+v", err, res)
// 结束当前代码块
	}
// 声明：var payCount int64
	var payCount int64
// 执行：testDB.Model(&model.Payment{}).Where("order_no = ?", o.OrderNo).Count(&payCount)
	testDB.Model(&model.Payment{}).Where("order_no = ?", o.OrderNo).Count(&payCount)
// 条件判断：if payCount != 1 {
	if payCount != 1 {
// 执行：t.Fatalf("重复通知不得重复记款: %d", payCount)
		t.Fatalf("重复通知不得重复记款: %d", payCount)
// 结束当前代码块
	}

	// 金额不符
	o2, _ := g.OrderService.Create(b1, pid, 1)
// 赋值/更新：badAmount := utils.PaySign(testCfg.Shop.PayNotifySecret, o2.OrderNo, "TXN-2", o2.AmountCents+1)
	badAmount := utils.PaySign(testCfg.Shop.PayNotifySecret, o2.OrderNo, "TXN-2", o2.AmountCents+1)
// 赋值/更新：res, _ = g.PaymentService.Notify(o2.OrderNo, "TXN-2", o2.AmountCents+1, badAmount, "{}")
	res, _ = g.PaymentService.Notify(o2.OrderNo, "TXN-2", o2.AmountCents+1, badAmount, "{}")
// 条件判断：if res.Result != model.NoticeRejected {
	if res.Result != model.NoticeRejected {
// 执行：t.Fatalf("金额不符应拒绝: %+v", res)
		t.Fatalf("金额不符应拒绝: %+v", res)
// 结束当前代码块
	}

	// 验签失败（且必须留痕可查）
	res, err = g.PaymentService.Notify(o2.OrderNo, "TXN-3", o2.AmountCents, "deadbeef", "{}")
// 条件判断：if !errors.Is(err, ErrBadSign) || res.Result != model.NoticeRejected {
	if !errors.Is(err, ErrBadSign) || res.Result != model.NoticeRejected {
// 执行：t.Fatalf("验签失败应拒绝: %v %+v", err, res)
		t.Fatalf("验签失败应拒绝: %v %+v", err, res)
// 结束当前代码块
	}
// 赋值/更新：notices, _ := g.PaymentService.Notices(o2.OrderNo)
	notices, _ := g.PaymentService.Notices(o2.OrderNo)
// 赋值/更新：foundBadSign := false
	foundBadSign := false
// 循环处理：for _, n := range notices {
	for _, n := range notices {
// 条件判断：if !n.SignOK && n.Result == model.NoticeRejected {
		if !n.SignOK && n.Result == model.NoticeRejected {
// 赋值/更新：foundBadSign = true
			foundBadSign = true
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if !foundBadSign {
	if !foundBadSign {
// 执行：t.Fatal("验签失败应可查询处理结果")
		t.Fatal("验签失败应可查询处理结果")
// 结束当前代码块
	}

	// 未知订单
	unknownSign := utils.PaySign(testCfg.Shop.PayNotifySecret, "SO-NOPE", "TXN-4", 10000)
// 赋值/更新：res, _ = g.PaymentService.Notify("SO-NOPE", "TXN-4", 10000, unknownSign, "{}")
	res, _ = g.PaymentService.Notify("SO-NOPE", "TXN-4", 10000, unknownSign, "{}")
// 条件判断：if res.Result != model.NoticeRejected {
	if res.Result != model.NoticeRejected {
// 执行：t.Fatalf("未知订单应拒绝: %+v", res)
		t.Fatalf("未知订单应拒绝: %+v", res)
// 结束当前代码块
	}

	// 同一渠道交易号不得归属两个订单
	o3, _ := g.OrderService.Create(b1, pid, 1)
// 赋值/更新：signDup := utils.PaySign(testCfg.Shop.PayNotifySecret, o3.OrderNo, "TXN-1", o3.AmountCents)
	signDup := utils.PaySign(testCfg.Shop.PayNotifySecret, o3.OrderNo, "TXN-1", o3.AmountCents)
// 赋值/更新：res, err = g.PaymentService.Notify(o3.OrderNo, "TXN-1", o3.AmountCents, signDup, "{}")
	res, err = g.PaymentService.Notify(o3.OrderNo, "TXN-1", o3.AmountCents, signDup, "{}")
// 条件判断：if !errors.Is(err, ErrDuplicatedPay) || res.Result != model.NoticeDuplicated {
	if !errors.Is(err, ErrDuplicatedPay) || res.Result != model.NoticeDuplicated {
// 执行：t.Fatalf("交易号跨订单应拒绝: %v %+v", err, res)
		t.Fatalf("交易号跨订单应拒绝: %v %+v", err, res)
// 结束当前代码块
	}
// 赋值/更新：o3db, _ := g.OrderService.Get(o3.OrderNo)
	o3db, _ := g.OrderService.Get(o3.OrderNo)
// 条件判断：if o3db.Status != model.OrderStatusPendingPayment {
	if o3db.Status != model.OrderStatusPendingPayment {
// 执行：t.Fatalf("冲突通知不得改变另一订单状态: %s", o3db.Status)
		t.Fatalf("冲突通知不得改变另一订单状态: %s", o3db.Status)
// 结束当前代码块
	}

	// 有效期外到达 → 异常收款，不恢复订单/库存
	// 注：此处 o2/o3/o4 三笔待付款单均已超时，清扫器应一并关闭（正确行为，原断言 1 有误）
	o4, _ := g.OrderService.Create(b1, pid, 1)
// 执行：clock.Advance(11 * time.Minute)
	clock.Advance(11 * time.Minute)
// 条件判断：if n, err := g.OrderService.CloseExpired(); err != nil || n != 3 {
	if n, err := g.OrderService.CloseExpired(); err != nil || n != 3 {
// 执行：t.Fatalf("清扫器应关闭全部超时订单: %d %v", n, err)
		t.Fatalf("清扫器应关闭全部超时订单: %d %v", n, err)
// 结束当前代码块
	}
// 赋值/更新：sign4 := utils.PaySign(testCfg.Shop.PayNotifySecret, o4.OrderNo, "TXN-5", o4.AmountCents)
	sign4 := utils.PaySign(testCfg.Shop.PayNotifySecret, o4.OrderNo, "TXN-5", o4.AmountCents)
// 赋值/更新：res, err = g.PaymentService.Notify(o4.OrderNo, "TXN-5", o4.AmountCents, sign4, "{}")
	res, err = g.PaymentService.Notify(o4.OrderNo, "TXN-5", o4.AmountCents, sign4, "{}")
// 条件判断：if err != nil || res.Result != model.NoticeAbnormal {
	if err != nil || res.Result != model.NoticeAbnormal {
// 执行：t.Fatalf("迟到支付应转异常: %v %+v", err, res)
		t.Fatalf("迟到支付应转异常: %v %+v", err, res)
// 结束当前代码块
	}
// 赋值/更新：late, _ := g.OrderService.Get(o4.OrderNo)
	late, _ := g.OrderService.Get(o4.OrderNo)
// 条件判断：if late.Status != model.OrderStatusClosed {
	if late.Status != model.OrderStatusClosed {
// 执行：t.Fatalf("迟到通知不得恢复已关闭订单: %s", late.Status)
		t.Fatalf("迟到通知不得恢复已关闭订单: %s", late.Status)
// 结束当前代码块
	}
// 赋值/更新：abns, _ := g.PaymentService.AbnormalList(model.AbnormalPending)
	abns, _ := g.PaymentService.AbnormalList(model.AbnormalPending)
// 条件判断：if len(abns) == 0 || abns[0].ChannelTxnID != "TXN-5" {
	if len(abns) == 0 || abns[0].ChannelTxnID != "TXN-5" {
// 执行：t.Fatalf("异常记录未落库: %+v", abns)
		t.Fatalf("异常记录未落库: %+v", abns)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestRefundApplyPauseShipAndReject(t *testing.T) {
func TestRefundApplyPauseShipAndReject(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：orderNo, _, _ := payOrder(t, g, b1, 2)
	orderNo, _, _ := payOrder(t, g, b1, 2)

// 条件判断：if _, err := g.RefundService.Apply(buyerID(t, "buyer2"), orderNo, "x"); !errors.Is(err, ErrForbidden) {
	if _, err := g.RefundService.Apply(buyerID(t, "buyer2"), orderNo, "x"); !errors.Is(err, ErrForbidden) {
// 执行：t.Fatalf("他人订单应拒绝: %v", err)
		t.Fatalf("他人订单应拒绝: %v", err)
// 结束当前代码块
	}
// 赋值/更新：refund, err := g.RefundService.Apply(b1, orderNo, "不想要了")
	refund, err := g.RefundService.Apply(b1, orderNo, "不想要了")
// 条件判断：if err != nil || refund.AuditStatus != model.AuditPending {
	if err != nil || refund.AuditStatus != model.AuditPending {
// 执行：t.Fatalf("申请退款失败: %v %+v", err, refund)
		t.Fatalf("申请退款失败: %v %+v", err, refund)
// 结束当前代码块
	}
// 条件判断：if _, err := g.ShipmentService.Ship(buyerID(t, "wh1"), orderNo, "SF", "SF1"); !errors.Is(err, ErrRefundInProgr...
	if _, err := g.ShipmentService.Ship(buyerID(t, "wh1"), orderNo, "SF", "SF1"); !errors.Is(err, ErrRefundInProgress) {
// 执行：t.Fatalf("退款申请期间发货应拒绝: %v", err)
		t.Fatalf("退款申请期间发货应拒绝: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.RefundService.Apply(b1, orderNo, "again"); !errors.Is(err, ErrRefundExists) {
	if _, err := g.RefundService.Apply(b1, orderNo, "again"); !errors.Is(err, ErrRefundExists) {
// 执行：t.Fatalf("重复申请应拒绝: %v", err)
		t.Fatalf("重复申请应拒绝: %v", err)
// 结束当前代码块
	}

// 赋值/更新：rejected, err := g.RefundService.Audit(buyerID(t, "cs1"), orderNo, false, "发货在即")
	rejected, err := g.RefundService.Audit(buyerID(t, "cs1"), orderNo, false, "发货在即")
// 条件判断：if err != nil || rejected.AuditStatus != model.AuditRejected {
	if err != nil || rejected.AuditStatus != model.AuditRejected {
// 执行：t.Fatalf("拒绝退款失败: %v %+v", err, rejected)
		t.Fatalf("拒绝退款失败: %v %+v", err, rejected)
// 结束当前代码块
	}
// 赋值/更新：paid, _ := g.OrderService.Get(orderNo)
	paid, _ := g.OrderService.Get(orderNo)
// 条件判断：if paid.Status != model.OrderStatusPaid {
	if paid.Status != model.OrderStatusPaid {
// 执行：t.Fatalf("拒绝后订单应回到待发货: %s", paid.Status)
		t.Fatalf("拒绝后订单应回到待发货: %s", paid.Status)
// 结束当前代码块
	}
// 条件判断：if _, err := g.ShipmentService.Ship(buyerID(t, "wh1"), orderNo, "SF", "SF2"); err != nil {
	if _, err := g.ShipmentService.Ship(buyerID(t, "wh1"), orderNo, "SF", "SF2"); err != nil {
// 执行：t.Fatalf("拒绝后应允许发货: %v", err)
		t.Fatalf("拒绝后应允许发货: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.RefundService.Apply(b1, orderNo, "售后"); err == nil {
	if _, err := g.RefundService.Apply(b1, orderNo, "售后"); err == nil {
// 执行：t.Fatal("已发货订单不应受理售前退款")
		t.Fatal("已发货订单不应受理售前退款")
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestRefundChannelFailThenRetrySuccess(t *testing.T) {
func TestRefundChannelFailThenRetrySuccess(t *testing.T) {
// 赋值/更新：g, ch, _ := newEnv(t)
	g, ch, _ := newEnv(t)
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：cs := buyerID(t, "cs1")
	cs := buyerID(t, "cs1")
// 赋值/更新：orderNo, _, amount := payOrder(t, g, b1, 2)
	orderNo, _, amount := payOrder(t, g, b1, 2)
// 赋值/更新：before := available(t, "SKU-A")
	before := available(t, "SKU-A")

// 赋值/更新：refund, err := g.RefundService.Apply(b1, orderNo, "整单退款")
	refund, err := g.RefundService.Apply(b1, orderNo, "整单退款")
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 执行：ch.SetRefundMode("fail")
	ch.SetRefundMode("fail")
// 条件判断：if _, err := g.RefundService.Audit(cs, orderNo, true, "同意"); !errors.Is(err, ErrChannelRefund) {
	if _, err := g.RefundService.Audit(cs, orderNo, true, "同意"); !errors.Is(err, ErrChannelRefund) {
// 执行：t.Fatalf("首次退款应失败: %v", err)
		t.Fatalf("首次退款应失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：r, _ := g.RefundService.Get(orderNo)
	r, _ := g.RefundService.Get(orderNo)
// 条件判断：if r.ChannelStatus != model.ChannelFailed || r.RetryCount != 1 {
	if r.ChannelStatus != model.ChannelFailed || r.RetryCount != 1 {
// 执行：t.Fatalf("失败应保留待处理记录: %+v", r)
		t.Fatalf("失败应保留待处理记录: %+v", r)
// 结束当前代码块
	}
// 条件判断：if available(t, "SKU-A") != before {
	if available(t, "SKU-A") != before {
// 执行：t.Fatal("退款失败不得恢复库存")
		t.Fatal("退款失败不得恢复库存")
// 结束当前代码块
	}
// 赋值/更新：order, _ := g.OrderService.Get(orderNo)
	order, _ := g.OrderService.Get(orderNo)
// 条件判断：if order.Status != model.OrderStatusRefundFailed {
	if order.Status != model.OrderStatusRefundFailed {
// 执行：t.Fatalf("退款失败订单状态错误: %s", order.Status)
		t.Fatalf("退款失败订单状态错误: %s", order.Status)
// 结束当前代码块
	}

// 执行：ch.SetRefundMode("success")
	ch.SetRefundMode("success")
// 赋值/更新：retried, err := g.RefundService.Retry(cs, refund.RefundNo)
	retried, err := g.RefundService.Retry(cs, refund.RefundNo)
// 条件判断：if err != nil || retried.ChannelStatus != model.ChannelSuccess {
	if err != nil || retried.ChannelStatus != model.ChannelSuccess {
// 执行：t.Fatalf("重试退款应成功: %v %+v", err, retried)
		t.Fatalf("重试退款应成功: %v %+v", err, retried)
// 结束当前代码块
	}
// 赋值/更新：order, _ = g.OrderService.Get(orderNo)
	order, _ = g.OrderService.Get(orderNo)
// 条件判断：if order.Status != model.OrderStatusRefunded {
	if order.Status != model.OrderStatusRefunded {
// 执行：t.Fatalf("退款成功后订单状态错误: %s", order.Status)
		t.Fatalf("退款成功后订单状态错误: %s", order.Status)
// 结束当前代码块
	}
// 条件判断：if available(t, "SKU-A") != before+2 {
	if available(t, "SKU-A") != before+2 {
// 执行：t.Fatalf("售前退款应恢复 2 件: %d", available(t, "SKU-A"))
		t.Fatalf("售前退款应恢复 2 件: %d", available(t, "SKU-A"))
// 结束当前代码块
	}
// 条件判断：if _, err := g.RefundService.Retry(cs, refund.RefundNo); !errors.Is(err, ErrRefundNotRetryable) {
	if _, err := g.RefundService.Retry(cs, refund.RefundNo); !errors.Is(err, ErrRefundNotRetryable) {
// 执行：t.Fatalf("成功后退款不可重试: %v", err)
		t.Fatalf("成功后退款不可重试: %v", err)
// 结束当前代码块
	}
	// 重复退款不会重复恢复：直接再次调用 callChannel 等价的 Audit/Retry 已被状态拦住
	var restored int64
// 执行：testDB.Model(&model.Refund{}).Where("order_no = ? AND restored = ?", orderNo, true).Count(&restored)
	testDB.Model(&model.Refund{}).Where("order_no = ? AND restored = ?", orderNo, true).Count(&restored)
// 条件判断：if restored != 1 {
	if restored != 1 {
// 执行：t.Fatalf("restored 闸门应只生效一次: %d", restored)
		t.Fatalf("restored 闸门应只生效一次: %d", restored)
// 结束当前代码块
	}
// 条件判断：if amount != 20000 {
	if amount != 20000 {
// 执行：t.Fatalf("退款金额应等于订单金额: %d", amount)
		t.Fatalf("退款金额应等于订单金额: %d", amount)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestShipOnlyOnceAndStateErrors(t *testing.T) {
func TestShipOnlyOnceAndStateErrors(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：wh := buyerID(t, "wh1")
	wh := buyerID(t, "wh1")
// 赋值/更新：orderNo, _, _ := payOrder(t, g, b1, 1)
	orderNo, _, _ := payOrder(t, g, b1, 1)

// 条件判断：if _, err := g.ShipmentService.Ship(wh, orderNo, "", "SF"); !errors.Is(err, ErrInvalidParam) {
	if _, err := g.ShipmentService.Ship(wh, orderNo, "", "SF"); !errors.Is(err, ErrInvalidParam) {
// 执行：t.Fatalf("缺少承运方应拒绝: %v", err)
		t.Fatalf("缺少承运方应拒绝: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.ShipmentService.Ship(wh, orderNo, "SF", "SF100"); err != nil {
	if _, err := g.ShipmentService.Ship(wh, orderNo, "SF", "SF100"); err != nil {
// 执行：t.Fatalf("发货失败: %v", err)
		t.Fatalf("发货失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.ShipmentService.Ship(wh, orderNo, "SF", "SF101"); !errors.Is(err, ErrAlreadyShipped) {
	if _, err := g.ShipmentService.Ship(wh, orderNo, "SF", "SF101"); !errors.Is(err, ErrAlreadyShipped) {
// 执行：t.Fatalf("重复发货应拒绝: %v", err)
		t.Fatalf("重复发货应拒绝: %v", err)
// 结束当前代码块
	}
// 赋值/更新：sh, err := g.ShipmentService.Get(orderNo)
	sh, err := g.ShipmentService.Get(orderNo)
// 条件判断：if err != nil || sh == nil || sh.TrackingNo != "SF100" {
	if err != nil || sh == nil || sh.TrackingNo != "SF100" {
// 执行：t.Fatalf("运单查询失败: %v %+v", err, sh)
		t.Fatalf("运单查询失败: %v %+v", err, sh)
// 结束当前代码块
	}
// 条件判断：if list, _ := g.ShipmentService.Shipable(); len(list) != 0 {
	if list, _ := g.ShipmentService.Shipable(); len(list) != 0 {
// 执行：t.Fatalf("已发货订单不应出现在待发货列表: %d", len(list))
		t.Fatalf("已发货订单不应出现在待发货列表: %d", len(list))
// 结束当前代码块
	}

	// 未付款订单发货
	o2, _ := g.OrderService.Create(b1, productID(t, "SKU-A"), 1)
// 条件判断：if _, err := g.ShipmentService.Ship(wh, o2.OrderNo, "SF", "SF200"); !errors.Is(err, ErrNotPaid) {
	if _, err := g.ShipmentService.Ship(wh, o2.OrderNo, "SF", "SF200"); !errors.Is(err, ErrNotPaid) {
// 执行：t.Fatalf("未付款订单发货应拒绝: %v", err)
		t.Fatalf("未付款订单发货应拒绝: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestLatePaymentAbnormalRefund(t *testing.T) {
func TestLatePaymentAbnormalRefund(t *testing.T) {
// 赋值/更新：g, ch, clock := newEnv(t)
	g, ch, clock := newEnv(t)
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：cs := buyerID(t, "cs1")
	cs := buyerID(t, "cs1")
// 赋值/更新：o, _ := g.OrderService.Create(b1, productID(t, "SKU-A"), 2)
	o, _ := g.OrderService.Create(b1, productID(t, "SKU-A"), 2)
// 执行：clock.Advance(11 * time.Minute)
	clock.Advance(11 * time.Minute)
// 条件判断：if n, err := g.OrderService.CloseExpired(); err != nil || n != 1 {
	if n, err := g.OrderService.CloseExpired(); err != nil || n != 1 {
// 执行：t.Fatalf("超时关闭失败: %d %v", n, err)
		t.Fatalf("超时关闭失败: %d %v", n, err)
// 结束当前代码块
	}
// 赋值/更新：stockAfterClose := available(t, "SKU-A")
	stockAfterClose := available(t, "SKU-A")

// 赋值/更新：txn := "TXN-LATE-" + o.OrderNo
	txn := "TXN-LATE-" + o.OrderNo
// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, txn, o.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, txn, o.AmountCents)
// 赋值/更新：res, err := g.PaymentService.Notify(o.OrderNo, txn, o.AmountCents, sign, "{}")
	res, err := g.PaymentService.Notify(o.OrderNo, txn, o.AmountCents, sign, "{}")
// 条件判断：if err != nil || res.Result != model.NoticeAbnormal {
	if err != nil || res.Result != model.NoticeAbnormal {
// 执行：t.Fatalf("迟到支付应转异常: %v %+v", err, res)
		t.Fatalf("迟到支付应转异常: %v %+v", err, res)
// 结束当前代码块
	}
// 赋值/更新：abns, _ := g.PaymentService.AbnormalList("")
	abns, _ := g.PaymentService.AbnormalList("")
// 条件判断：if len(abns) != 1 {
	if len(abns) != 1 {
// 执行：t.Fatalf("异常记录应为 1 条: %d", len(abns))
		t.Fatalf("异常记录应为 1 条: %d", len(abns))
// 结束当前代码块
	}

// 执行：ch.SetRefundMode("fail")
	ch.SetRefundMode("fail")
// 条件判断：if _, err := g.RefundService.HandleAbnormal(cs, txn, "核实迟到支付"); !errors.Is(err, ErrChannelRefund) {
	if _, err := g.RefundService.HandleAbnormal(cs, txn, "核实迟到支付"); !errors.Is(err, ErrChannelRefund) {
// 执行：t.Fatalf("渠道失败应返回待重试: %v", err)
		t.Fatalf("渠道失败应返回待重试: %v", err)
// 结束当前代码块
	}
// 条件判断：if available(t, "SKU-A") != stockAfterClose {
	if available(t, "SKU-A") != stockAfterClose {
// 执行：t.Fatal("异常退款失败不得改变库存")
		t.Fatal("异常退款失败不得改变库存")
// 结束当前代码块
	}

// 执行：ch.SetRefundMode("success")
	ch.SetRefundMode("success")
// 赋值/更新：refund, err := g.RefundService.HandleAbnormal(cs, txn, "重试退款")
	refund, err := g.RefundService.HandleAbnormal(cs, txn, "重试退款")
// 条件判断：if err != nil || refund.ChannelStatus != model.ChannelSuccess {
	if err != nil || refund.ChannelStatus != model.ChannelSuccess {
// 执行：t.Fatalf("异常退款重试失败: %v %+v", err, refund)
		t.Fatalf("异常退款重试失败: %v %+v", err, refund)
// 结束当前代码块
	}
// 条件判断：if available(t, "SKU-A") != stockAfterClose {
	if available(t, "SKU-A") != stockAfterClose {
// 执行：t.Fatalf("迟到支付退款不得恢复已释放数量: %d", available(t, "SKU-A"))
		t.Fatalf("迟到支付退款不得恢复已释放数量: %d", available(t, "SKU-A"))
// 结束当前代码块
	}
// 赋值/更新：abns, _ = g.PaymentService.AbnormalList("")
	abns, _ = g.PaymentService.AbnormalList("")
// 条件判断：if abns[0].Status != model.AbnormalHandled {
	if abns[0].Status != model.AbnormalHandled {
// 执行：t.Fatalf("异常记录应标记已处理: %s", abns[0].Status)
		t.Fatalf("异常记录应标记已处理: %s", abns[0].Status)
// 结束当前代码块
	}
// 条件判断：if _, err := g.RefundService.HandleAbnormal(cs, txn, "再次处理"); err == nil {
	if _, err := g.RefundService.HandleAbnormal(cs, txn, "再次处理"); err == nil {
// 执行：t.Fatal("已处理异常不应重复退款")
		t.Fatal("已处理异常不应重复退款")
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestConcurrentRefundAndShipOnlyOneResult(t *testing.T) {
func TestConcurrentRefundAndShipOnlyOneResult(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：wh := buyerID(t, "wh1")
	wh := buyerID(t, "wh1")
// 赋值/更新：orderNo, _, _ := payOrder(t, g, b1, 2)
	orderNo, _, _ := payOrder(t, g, b1, 2)
// 赋值/更新：before := available(t, "SKU-A")
	before := available(t, "SKU-A")

// 声明：var wg sync.WaitGroup
	var wg sync.WaitGroup
// 声明：var shipErr, refundErr error
	var shipErr, refundErr error
// 执行：wg.Add(2)
	wg.Add(2)
// 异步/延迟执行：go func() { defer wg.Done(); _, shipErr = g.ShipmentService.Ship(wh, orderNo, "SF", "SF-RACE") }()
	go func() { defer wg.Done(); _, shipErr = g.ShipmentService.Ship(wh, orderNo, "SF", "SF-RACE") }()
// 异步/延迟执行：go func() { defer wg.Done(); _, refundErr = g.RefundService.Apply(b1, orderNo, "并发退款") }()
	go func() { defer wg.Done(); _, refundErr = g.RefundService.Apply(b1, orderNo, "并发退款") }()
// 执行：wg.Wait()
	wg.Wait()

// 赋值/更新：order, _ := g.OrderService.Get(orderNo)
	order, _ := g.OrderService.Get(orderNo)
// 多路分支：switch {
	switch {
// 匹配分支：case shipErr == nil && refundErr != nil:
	case shipErr == nil && refundErr != nil:
// 条件判断：if order.Status != model.OrderStatusShipped {
		if order.Status != model.OrderStatusShipped {
// 执行：t.Fatalf("发货成功则应为已发货: %s", order.Status)
			t.Fatalf("发货成功则应为已发货: %s", order.Status)
// 结束当前代码块
		}
// 条件判断：if available(t, "SKU-A") != before {
		if available(t, "SKU-A") != before {
// 执行：t.Fatal("发货路径不应恢复库存")
			t.Fatal("发货路径不应恢复库存")
// 结束当前代码块
		}
// 匹配分支：case refundErr == nil && shipErr != nil:
	case refundErr == nil && shipErr != nil:
// 条件判断：if order.Status != model.OrderStatusRefundRequested {
		if order.Status != model.OrderStatusRefundRequested {
// 执行：t.Fatalf("退款申请成功则应为退款申请中: %s", order.Status)
			t.Fatalf("退款申请成功则应为退款申请中: %s", order.Status)
// 结束当前代码块
		}
// 默认分支
	default:
// 执行：t.Fatalf("并发必须恰好一种成功: ship=%v refund=%v", shipErr, refundErr)
		t.Fatalf("并发必须恰好一种成功: ship=%v refund=%v", shipErr, refundErr)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestReportDetailTrailLedgerAndSummary(t *testing.T) {
func TestReportDetailTrailLedgerAndSummary(t *testing.T) {
// 赋值/更新：g, _, clock := newEnv(t)
	g, _, clock := newEnv(t)
// 赋值/更新：b1, b2 := buyerID(t, "buyer1"), buyerID(t, "buyer2")
	b1, b2 := buyerID(t, "buyer1"), buyerID(t, "buyer2")
// 赋值/更新：ops := buyerID(t, "ops1")
	ops := buyerID(t, "ops1")
// 赋值/更新：orderNo, _, _ := payOrder(t, g, b1, 2)
	orderNo, _, _ := payOrder(t, g, b1, 2)

	// 越权视图
	if _, err := g.ReportService.OrderDetail(orderNo, b2, "buyer"); !errors.Is(err, ErrForbidden) {
// 执行：t.Fatalf("他人订单详情应拒绝: %v", err)
		t.Fatalf("他人订单详情应拒绝: %v", err)
// 结束当前代码块
	}
// 赋值/更新：detail, err := g.ReportService.OrderDetail(orderNo, b1, "buyer")
	detail, err := g.ReportService.OrderDetail(orderNo, b1, "buyer")
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if detail.Payment == nil || detail.ProductInfo == nil || !detail.ShipableNow {
	if detail.Payment == nil || detail.ProductInfo == nil || !detail.ShipableNow {
// 执行：t.Fatalf("订单详情聚合字段缺失: %+v", detail)
		t.Fatalf("订单详情聚合字段缺失: %+v", detail)
// 结束当前代码块
	}
// 条件判断：if len(detail.Timeline) == 0 {
	if len(detail.Timeline) == 0 {
// 执行：t.Fatal("订单轨迹应可查")
		t.Fatal("订单轨迹应可查")
// 结束当前代码块
	}

	// 迟到支付收款 10000，验证运营汇总口径
	o2, _ := g.OrderService.Create(b2, productID(t, "SKU-A"), 1)
// 执行：clock.Advance(11 * time.Minute)
	clock.Advance(11 * time.Minute)
// 执行：g.OrderService.CloseExpired()
	g.OrderService.CloseExpired()
// 赋值/更新：txn := "TXN-LATE-" + o2.OrderNo
	txn := "TXN-LATE-" + o2.OrderNo
// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o2.OrderNo, txn, o2.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o2.OrderNo, txn, o2.AmountCents)
// 条件判断：if _, err := g.PaymentService.Notify(o2.OrderNo, txn, o2.AmountCents, sign, "{}"); err != nil {
	if _, err := g.PaymentService.Notify(o2.OrderNo, txn, o2.AmountCents, sign, "{}"); err != nil {
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
	// 增量续写 R5 修复后：有效收款与异常渠道收款分开显示（OPS-03/BR-14），
	// 迟到支付 10000 分不再并入 TotalPaidCents，改入 TotalAbnormalCents。
	if sum.TotalPaidCents != 20000 {
// 执行：t.Fatalf("有效收款应不含迟到支付（20000）: %d", sum.TotalPaidCents)
		t.Fatalf("有效收款应不含迟到支付（20000）: %d", sum.TotalPaidCents)
// 结束当前代码块
	}
// 条件判断：if sum.TotalAbnormalCents != 10000 {
	if sum.TotalAbnormalCents != 10000 {
// 执行：t.Fatalf("异常收款应独立合计（10000）: %d", sum.TotalAbnormalCents)
		t.Fatalf("异常收款应独立合计（10000）: %d", sum.TotalAbnormalCents)
// 结束当前代码块
	}
// 条件判断：if sum.TotalRefundedCents != 0 {
	if sum.TotalRefundedCents != 0 {
// 执行：t.Fatalf("未退款时应为 0: %d", sum.TotalRefundedCents)
		t.Fatalf("未退款时应为 0: %d", sum.TotalRefundedCents)
// 结束当前代码块
	}
// 声明：var prodPaid int64
	var prodPaid int64
// 声明：var prodAbnormal int64
	var prodAbnormal int64
// 循环处理：for _, p := range sum.Products {
	for _, p := range sum.Products {
// 条件判断：if p.SKU == "SKU-A" {
		if p.SKU == "SKU-A" {
// 赋值/更新：prodPaid = p.PaidCents
			prodPaid = p.PaidCents
// 赋值/更新：prodAbnormal = p.AbnormalCents
			prodAbnormal = p.AbnormalCents
// 条件判断：if p.Available != 3 {
			if p.Available != 3 {
// 执行：t.Fatalf("商品可售数量错误: %d", p.Available)
				t.Fatalf("商品可售数量错误: %d", p.Available)
// 结束当前代码块
			}
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if prodPaid != 20000 {
	if prodPaid != 20000 {
// 执行：t.Fatalf("商品维度有效收款应为 20000: %d", prodPaid)
		t.Fatalf("商品维度有效收款应为 20000: %d", prodPaid)
// 结束当前代码块
	}
// 条件判断：if prodAbnormal != 10000 {
	if prodAbnormal != 10000 {
// 执行：t.Fatalf("商品维度异常收款应为 10000: %d", prodAbnormal)
		t.Fatalf("商品维度异常收款应为 10000: %d", prodAbnormal)
// 结束当前代码块
	}
// 条件判断：if sum.OrderCount != 2 {
	if sum.OrderCount != 2 {
// 执行：t.Fatalf("订单总数错误: %d", sum.OrderCount)
		t.Fatalf("订单总数错误: %d", sum.OrderCount)
// 结束当前代码块
	}

// 赋值/更新：trail, err := g.ReportService.Trail(orderNo)
	trail, err := g.ReportService.Trail(orderNo)
// 条件判断：if err != nil || len(trail) == 0 {
	if err != nil || len(trail) == 0 {
// 执行：t.Fatalf("操作轨迹应可查: %v %d", err, len(trail))
		t.Fatalf("操作轨迹应可查: %v %d", err, len(trail))
// 结束当前代码块
	}
// 赋值/更新：ledger, err := g.ReportService.Ledger(productID(t, "SKU-A"))
	ledger, err := g.ReportService.Ledger(productID(t, "SKU-A"))
// 条件判断：if err != nil || len(ledger) == 0 {
	if err != nil || len(ledger) == 0 {
// 执行：t.Fatalf("库存流水应可查: %v %d", err, len(ledger))
		t.Fatalf("库存流水应可查: %v %d", err, len(ledger))
// 结束当前代码块
	}
// 条件判断：if _, err := g.ReportService.OrderDetail("SO-NOPE", ops, "ops"); !errors.Is(err, ErrOrderNotFound) {
	if _, err := g.ReportService.OrderDetail("SO-NOPE", ops, "ops"); !errors.Is(err, ErrOrderNotFound) {
// 执行：t.Fatalf("未知订单详情应报错: %v", err)
		t.Fatalf("未知订单详情应报错: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestCatalogAndAuthAndDemoReset(t *testing.T) {
func TestCatalogAndAuthAndDemoReset(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：ops := buyerID(t, "ops1")
	ops := buyerID(t, "ops1")
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")

// 条件判断：if _, err := g.CatalogService.Update(buyerID(t, "buyer1"), "buyer", pid, nil, nil); !errors.Is(err, ErrForbidd...
	if _, err := g.CatalogService.Update(buyerID(t, "buyer1"), "buyer", pid, nil, nil); !errors.Is(err, ErrForbidden) {
// 执行：t.Fatalf("买家不能维护商品: %v", err)
		t.Fatalf("买家不能维护商品: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.CatalogService.Update(ops, "ops", 99999, nil, nil); !errors.Is(err, ErrProductNotFound) {
	if _, err := g.CatalogService.Update(ops, "ops", 99999, nil, nil); !errors.Is(err, ErrProductNotFound) {
// 执行：t.Fatalf("未知商品应报错: %v", err)
		t.Fatalf("未知商品应报错: %v", err)
// 结束当前代码块
	}
// 条件判断：if list, err := g.CatalogService.OnSale(); err != nil || len(list) != 2 {
	if list, err := g.CatalogService.OnSale(); err != nil || len(list) != 2 {
// 执行：t.Fatalf("在售商品应为 2 个: %v %d", err, len(list))
		t.Fatalf("在售商品应为 2 个: %v %d", err, len(list))
// 结束当前代码块
	}
// 条件判断：if list, err := g.CatalogService.All(); err != nil || len(list) != 2 {
	if list, err := g.CatalogService.All(); err != nil || len(list) != 2 {
// 执行：t.Fatalf("全部商品应为 2 个: %v %d", err, len(list))
		t.Fatalf("全部商品应为 2 个: %v %d", err, len(list))
// 结束当前代码块
	}

	// 登录
	if _, err := g.AuthService.Login("buyer1", "wrong"); !errors.Is(err, ErrLoginFailed) {
// 执行：t.Fatalf("错误密码应拒绝: %v", err)
		t.Fatalf("错误密码应拒绝: %v", err)
// 结束当前代码块
	}
// 赋值/更新：login, err := g.AuthService.Login("buyer1", "shop123")
	login, err := g.AuthService.Login("buyer1", "shop123")
// 条件判断：if err != nil || login.Token == "" || login.User.Role != "buyer" {
	if err != nil || login.Token == "" || login.User.Role != "buyer" {
// 执行：t.Fatalf("登录失败: %v %+v", err, login)
		t.Fatalf("登录失败: %v %+v", err, login)
// 结束当前代码块
	}
// 条件判断：if u, err := g.AuthService.User(login.User.ID); err != nil || u.Username != "buyer1" {
	if u, err := g.AuthService.User(login.User.ID); err != nil || u.Username != "buyer1" {
// 执行：t.Fatalf("用户查询失败: %v", err)
		t.Fatalf("用户查询失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if users, err := g.AuthService.Users("buyer"); err != nil || len(users) != 3 {
	if users, err := g.AuthService.Users("buyer"); err != nil || len(users) != 3 {
// 执行：t.Fatalf("买家账号应为 3 个: %v %d", err, len(users))
		t.Fatalf("买家账号应为 3 个: %v %d", err, len(users))
// 结束当前代码块
	}
// 条件判断：if users, err := g.AuthService.Users(""); err != nil || len(users) < 6 {
	if users, err := g.AuthService.Users(""); err != nil || len(users) < 6 {
// 执行：t.Fatalf("全部账号应至少 6 个: %v %d", err, len(users))
		t.Fatalf("全部账号应至少 6 个: %v %d", err, len(users))
// 结束当前代码块
	}

	// 演示重置
	if _, err := g.OrderService.Create(buyerID(t, "buyer1"), pid, 2); err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if err := g.DemoService.Reset(); err != nil {
	if err := g.DemoService.Reset(); err != nil {
// 执行：t.Fatalf("重置失败: %v", err)
		t.Fatalf("重置失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if available(t, "SKU-A") != 5 {
	if available(t, "SKU-A") != 5 {
// 执行：t.Fatal("重置应恢复初始库存")
		t.Fatal("重置应恢复初始库存")
// 结束当前代码块
	}
// 条件判断：if list, _ := g.OrderService.ListByStatus("", 0); len(list) != 0 {
	if list, _ := g.OrderService.ListByStatus("", 0); len(list) != 0 {
// 执行：t.Fatalf("重置应清空订单: %d", len(list))
		t.Fatalf("重置应清空订单: %d", len(list))
// 结束当前代码块
	}
// 条件判断：if users, _ := g.AuthService.Users(""); len(users) < 6 {
	if users, _ := g.AuthService.Users(""); len(users) < 6 {
// 执行：t.Fatal("重置应保留账号")
		t.Fatal("重置应保留账号")
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestSimulatedChannel(t *testing.T) {
func TestSimulatedChannel(t *testing.T) {
// 赋值/更新：ch := NewSimulatedChannel(false)
	ch := NewSimulatedChannel(false)
// 条件判断：if ch.RefundMode() != "fail" {
	if ch.RefundMode() != "fail" {
// 执行：t.Fatalf("初始模式错误: %s", ch.RefundMode())
		t.Fatalf("初始模式错误: %s", ch.RefundMode())
// 结束当前代码块
	}
// 条件判断：if ok, _ := ch.Refund("SO-1", "RF1", 100); ok {
	if ok, _ := ch.Refund("SO-1", "RF1", 100); ok {
// 执行：t.Fatal("fail 模式应返回失败")
		t.Fatal("fail 模式应返回失败")
// 结束当前代码块
	}
// 执行：ch.SetRefundMode("success")
	ch.SetRefundMode("success")
// 条件判断：if ok, _ := ch.Refund("SO-1", "RF1", 100); !ok {
	if ok, _ := ch.Refund("SO-1", "RF1", 100); !ok {
// 执行：t.Fatal("success 模式应返回成功")
		t.Fatal("success 模式应返回成功")
// 结束当前代码块
	}
// 执行：ch.SetRefundMode("非法值")
	ch.SetRefundMode("非法值")
// 条件判断：if ch.RefundMode() != "success" {
	if ch.RefundMode() != "success" {
// 执行：t.Fatal("非法模式不应生效")
		t.Fatal("非法模式不应生效")
// 结束当前代码块
	}
// 条件判断：if ch.RefundCalls() != 2 {
	if ch.RefundCalls() != 2 {
// 执行：t.Fatalf("调用次数错误: %d", ch.RefundCalls())
		t.Fatalf("调用次数错误: %d", ch.RefundCalls())
// 结束当前代码块
	}
// 条件判断：if NewSimulatedChannel(true).RefundMode() != "success" {
	if NewSimulatedChannel(true).RefundMode() != "success" {
// 执行：t.Fatal("允许退款时应为 success")
		t.Fatal("允许退款时应为 success")
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestOrderErrorsAndHelpers(t *testing.T) {
func TestOrderErrorsAndHelpers(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 条件判断：if _, err := g.OrderService.Get("SO-NOPE"); !errors.Is(err, ErrOrderNotFound) {
	if _, err := g.OrderService.Get("SO-NOPE"); !errors.Is(err, ErrOrderNotFound) {
// 执行：t.Fatalf("未知订单应报错: %v", err)
		t.Fatalf("未知订单应报错: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.OrderService.Cancel(b1, "SO-NOPE"); !errors.Is(err, ErrOrderNotFound) {
	if _, err := g.OrderService.Cancel(b1, "SO-NOPE"); !errors.Is(err, ErrOrderNotFound) {
// 执行：t.Fatalf("取消未知订单应报错: %v", err)
		t.Fatalf("取消未知订单应报错: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.RefundService.Apply(b1, "SO-NOPE", "x"); !errors.Is(err, ErrOrderNotFound) {
	if _, err := g.RefundService.Apply(b1, "SO-NOPE", "x"); !errors.Is(err, ErrOrderNotFound) {
// 执行：t.Fatalf("未知订单退款应报错: %v", err)
		t.Fatalf("未知订单退款应报错: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.RefundService.Audit(buyerID(t, "cs1"), "SO-NOPE", true, "x"); !errors.Is(err, ErrRefundNotFound...
	if _, err := g.RefundService.Audit(buyerID(t, "cs1"), "SO-NOPE", true, "x"); !errors.Is(err, ErrRefundNotFound) {
// 执行：t.Fatalf("未知退款审核应报错: %v", err)
		t.Fatalf("未知退款审核应报错: %v", err)
// 结束当前代码块
	}
	// 客服审核/异常处理必须写明原因（空或纯空白一律拒绝）
	if _, err := g.RefundService.Audit(buyerID(t, "cs1"), "SO-NOPE", true, "   "); !errors.Is(err, ErrReasonRequired) {
// 执行：t.Fatalf("空原因审核应报错: %v", err)
		t.Fatalf("空原因审核应报错: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.RefundService.HandleAbnormal(buyerID(t, "cs1"), "TXN-NOPE", ""); !errors.Is(err, ErrReasonRequi...
	if _, err := g.RefundService.HandleAbnormal(buyerID(t, "cs1"), "TXN-NOPE", ""); !errors.Is(err, ErrReasonRequired) {
// 执行：t.Fatalf("空原因异常处理应报错: %v", err)
		t.Fatalf("空原因异常处理应报错: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.RefundService.Retry(buyerID(t, "cs1"), "RF-NOPE"); !errors.Is(err, ErrRefundNotFound) {
	if _, err := g.RefundService.Retry(buyerID(t, "cs1"), "RF-NOPE"); !errors.Is(err, ErrRefundNotFound) {
// 执行：t.Fatalf("未知退款重试应报错: %v", err)
		t.Fatalf("未知退款重试应报错: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.RefundService.HandleAbnormal(buyerID(t, "cs1"), "TXN-NOPE", "x"); !errors.Is(err, ErrAbnormalNo...
	if _, err := g.RefundService.HandleAbnormal(buyerID(t, "cs1"), "TXN-NOPE", "x"); !errors.Is(err, ErrAbnormalNotFound) {
// 执行：t.Fatalf("未知异常应报错: %v", err)
		t.Fatalf("未知异常应报错: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.ShipmentService.Ship(buyerID(t, "wh1"), "SO-NOPE", "SF", "1"); !errors.Is(err, ErrOrderNotFound...
	if _, err := g.ShipmentService.Ship(buyerID(t, "wh1"), "SO-NOPE", "SF", "1"); !errors.Is(err, ErrOrderNotFound) {
// 执行：t.Fatalf("未知订单发货应报错: %v", err)
		t.Fatalf("未知订单发货应报错: %v", err)
// 结束当前代码块
	}
// 赋值/更新：sh, err := g.ShipmentService.Get("SO-NOPE")
	sh, err := g.ShipmentService.Get("SO-NOPE")
// 条件判断：if err != nil || sh != nil {
	if err != nil || sh != nil {
// 执行：t.Fatalf("无发货记录应返回 nil: %v %+v", err, sh)
		t.Fatalf("无发货记录应返回 nil: %v %+v", err, sh)
// 结束当前代码块
	}
// 条件判断：if list, err := g.ShipmentService.List(); err != nil || len(list) != 0 {
	if list, err := g.ShipmentService.List(); err != nil || len(list) != 0 {
// 执行：t.Fatalf("发货列表应为空: %v %d", err, len(list))
		t.Fatalf("发货列表应为空: %v %d", err, len(list))
// 结束当前代码块
	}
// 条件判断：if list, err := g.RefundService.List(model.AuditPending); err != nil || len(list) != 0 {
	if list, err := g.RefundService.List(model.AuditPending); err != nil || len(list) != 0 {
// 执行：t.Fatalf("退款列表应为空: %v %d", err, len(list))
		t.Fatalf("退款列表应为空: %v %d", err, len(list))
// 结束当前代码块
	}
// 条件判断：if list, err := g.PaymentService.NoticesByResult(model.NoticeAccepted, 10); err != nil || len(list) != 0 {
	if list, err := g.PaymentService.NoticesByResult(model.NoticeAccepted, 10); err != nil || len(list) != 0 {
// 执行：t.Fatalf("通知过滤应为空: %v %d", err, len(list))
		t.Fatalf("通知过滤应为空: %v %d", err, len(list))
// 结束当前代码块
	}
// 条件判断：if list, err := g.OrderService.ListByBuyer(b1, ""); err != nil || len(list) != 0 {
	if list, err := g.OrderService.ListByBuyer(b1, ""); err != nil || len(list) != 0 {
// 执行：t.Fatalf("买家订单应为空: %v %d", err, len(list))
		t.Fatalf("买家订单应为空: %v %d", err, len(list))
// 结束当前代码块
	}

	// 状态错误翻译
	o, _ := g.OrderService.Create(b1, productID(t, "SKU-A"), 1)
// 条件判断：if _, err := g.RefundService.Apply(b1, o.OrderNo, "未付款"); !errors.Is(err, ErrNotPaid) {
	if _, err := g.RefundService.Apply(b1, o.OrderNo, "未付款"); !errors.Is(err, ErrNotPaid) {
// 执行：t.Fatalf("未付款申请退款应报未支付: %v", err)
		t.Fatalf("未付款申请退款应报未支付: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.ShipmentService.Ship(buyerID(t, "wh1"), o.OrderNo, "SF", "1"); !errors.Is(err, ErrNotPaid) {
	if _, err := g.ShipmentService.Ship(buyerID(t, "wh1"), o.OrderNo, "SF", "1"); !errors.Is(err, ErrNotPaid) {
// 执行：t.Fatalf("未付款发货应报未支付: %v", err)
		t.Fatalf("未付款发货应报未支付: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.OrderService.Cancel(buyerID(t, "buyer2"), o.OrderNo); !errors.Is(err, ErrForbidden) {
	if _, err := g.OrderService.Cancel(buyerID(t, "buyer2"), o.OrderNo); !errors.Is(err, ErrForbidden) {
// 执行：t.Fatalf("他人取消应拒绝: %v", err)
		t.Fatalf("他人取消应拒绝: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}
