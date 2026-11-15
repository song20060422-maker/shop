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
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// 契约兼容层 service 扩展方法的真实业务路径测试（替代工具函数补洞）。

// TestProductBySKU 契约以 SKU 作为商品寻址：命中/空串/未命中三分支。
func TestProductBySKU(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：p, err := g.CatalogService.ProductBySKU("SKU-A")
	p, err := g.CatalogService.ProductBySKU("SKU-A")
// 条件判断：if err != nil || p.SKU != "SKU-A" {
	if err != nil || p.SKU != "SKU-A" {
// 执行：t.Fatalf("按 SKU 查询失败: %v %+v", err, p)
		t.Fatalf("按 SKU 查询失败: %v %+v", err, p)
// 结束当前代码块
	}
// 条件判断：if _, err := g.CatalogService.ProductBySKU(""); !errors.Is(err, ErrProductNotFound) {
	if _, err := g.CatalogService.ProductBySKU(""); !errors.Is(err, ErrProductNotFound) {
// 执行：t.Fatalf("空 SKU 应返回未找到: %v", err)
		t.Fatalf("空 SKU 应返回未找到: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.CatalogService.ProductBySKU("NOPE"); !errors.Is(err, ErrProductNotFound) {
	if _, err := g.CatalogService.ProductBySKU("NOPE"); !errors.Is(err, ErrProductNotFound) {
// 执行：t.Fatalf("未知 SKU 应返回未找到: %v", err)
		t.Fatalf("未知 SKU 应返回未找到: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// TestOverwriteAvailable 运营覆盖可售数量：生效、追加流水、拒绝负数、无变化只记日志。
func TestOverwriteAvailable(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：ops := buyerID(t, "ops1")
	ops := buyerID(t, "ops1")

// 赋值/更新：p, err := g.CatalogService.OverwriteAvailable(ops, "ops", pid, 10)
	p, err := g.CatalogService.OverwriteAvailable(ops, "ops", pid, 10)
// 条件判断：if err != nil || p.Available != 10 {
	if err != nil || p.Available != 10 {
// 执行：t.Fatalf("覆盖可售数量失败: %v %d", err, p.Available)
		t.Fatalf("覆盖可售数量失败: %v %d", err, p.Available)
// 结束当前代码块
	}
// 声明：var led int64
	var led int64
// 执行：testDB.Table("inventory_ledger").Where("product_id = ? AND change_type = ?", pid, "restore").Count(&led)
	testDB.Table("inventory_ledger").Where("product_id = ? AND change_type = ?", pid, "restore").Count(&led)
// 条件判断：if led != 1 {
	if led != 1 {
// 执行：t.Fatalf("覆盖应追加 1 条库存流水: %d", led)
		t.Fatalf("覆盖应追加 1 条库存流水: %d", led)
// 结束当前代码块
	}

// 条件判断：if _, err := g.CatalogService.OverwriteAvailable(ops, "ops", pid, -1); err == nil {
	if _, err := g.CatalogService.OverwriteAvailable(ops, "ops", pid, -1); err == nil {
// 执行：t.Fatal("负数应被拒绝")
		t.Fatal("负数应被拒绝")
// 结束当前代码块
	}
// 条件判断：if _, err := g.CatalogService.OverwriteAvailable(ops, "ops", pid, 10); err != nil {
	if _, err := g.CatalogService.OverwriteAvailable(ops, "ops", pid, 10); err != nil {
// 执行：t.Fatalf("同值覆盖（无变化）不应报错: %v", err)
		t.Fatalf("同值覆盖（无变化）不应报错: %v", err)
// 结束当前代码块
	}
// 声明：var led2 int64
	var led2 int64
// 执行：testDB.Table("inventory_ledger").Where("product_id = ? AND change_type = ?", pid, "restore").Count(&led2)
	testDB.Table("inventory_ledger").Where("product_id = ? AND change_type = ?", pid, "restore").Count(&led2)
// 条件判断：if led2 != 1 {
	if led2 != 1 {
// 执行：t.Fatalf("无变化覆盖不得追加流水: %d", led2)
		t.Fatalf("无变化覆盖不得追加流水: %d", led2)
// 结束当前代码块
	}
// 条件判断：if _, err := g.CatalogService.OverwriteAvailable(ops, "ops", 999999, 5); !errors.Is(err, ErrProductNotFound) {
	if _, err := g.CatalogService.OverwriteAvailable(ops, "ops", 999999, 5); !errors.Is(err, ErrProductNotFound) {
// 执行：t.Fatalf("未知商品应返回未找到: %v", err)
		t.Fatalf("未知商品应返回未找到: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// TestFindByRefDualAddressing 订单双寻址（order_no / 数字主键）与未命中。
func TestFindByRefDualAddressing(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：o, _ := g.OrderService.Create(b1, pid, 1)
	o, _ := g.OrderService.Create(b1, pid, 1)

// 赋值/更新：byNo, err := g.OrderService.FindByRef(o.OrderNo)
	byNo, err := g.OrderService.FindByRef(o.OrderNo)
// 条件判断：if err != nil || byNo.OrderNo != o.OrderNo {
	if err != nil || byNo.OrderNo != o.OrderNo {
// 执行：t.Fatalf("按订单号寻址失败: %v", err)
		t.Fatalf("按订单号寻址失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：byID, err := g.OrderService.FindByRef(itoa(int(o.ID)))
	byID, err := g.OrderService.FindByRef(itoa(int(o.ID)))
// 条件判断：if err != nil || byID.OrderNo != o.OrderNo {
	if err != nil || byID.OrderNo != o.OrderNo {
// 执行：t.Fatalf("按主键寻址失败: %v", err)
		t.Fatalf("按主键寻址失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.OrderService.FindByRef(""); !errors.Is(err, ErrOrderNotFound) {
	if _, err := g.OrderService.FindByRef(""); !errors.Is(err, ErrOrderNotFound) {
// 执行：t.Fatalf("空引用应返回未找到: %v", err)
		t.Fatalf("空引用应返回未找到: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.OrderService.FindByRef("NOPE"); !errors.Is(err, ErrOrderNotFound) {
	if _, err := g.OrderService.FindByRef("NOPE"); !errors.Is(err, ErrOrderNotFound) {
// 执行：t.Fatalf("未知引用应返回未找到: %v", err)
		t.Fatalf("未知引用应返回未找到: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// TestPendingAndForceExpire 待付款清单与强制过期钩子。
func TestPendingAndForceExpire(t *testing.T) {
// 赋值/更新：g, _, clock := newEnv(t)
	g, _, clock := newEnv(t)
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：o, _ := g.OrderService.Create(b1, pid, 1)
	o, _ := g.OrderService.Create(b1, pid, 1)

// 赋值/更新：nos, err := g.OrderService.PendingOrderNos()
	nos, err := g.OrderService.PendingOrderNos()
// 条件判断：if err != nil || len(nos) != 1 || nos[0] != o.OrderNo {
	if err != nil || len(nos) != 1 || nos[0] != o.OrderNo {
// 执行：t.Fatalf("待付款清单应含新订单: %v %v", nos, err)
		t.Fatalf("待付款清单应含新订单: %v %v", nos, err)
// 结束当前代码块
	}

// 执行：clock.Advance(11 * time.Minute)
	clock.Advance(11 * time.Minute)
// 条件判断：if err := g.OrderService.ForceExpire(o.OrderNo, clock.Now()); err != nil {
	if err := g.OrderService.ForceExpire(o.OrderNo, clock.Now()); err != nil {
// 执行：t.Fatalf("强制过期失败: %v", err)
		t.Fatalf("强制过期失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := g.OrderService.ForceExpire("", clock.Now()); !errors.Is(err, ErrOrderNotFound) {
	if err := g.OrderService.ForceExpire("", clock.Now()); !errors.Is(err, ErrOrderNotFound) {
// 执行：t.Fatalf("空单号应返回未找到: %v", err)
		t.Fatalf("空单号应返回未找到: %v", err)
// 结束当前代码块
	}
// 条件判断：if n, _ := g.OrderService.CloseExpired(); n != 1 {
	if n, _ := g.OrderService.CloseExpired(); n != 1 {
// 执行：t.Fatalf("强制过期后应可被清扫: %d", n)
		t.Fatalf("强制过期后应可被清扫: %d", n)
// 结束当前代码块
	}
// 条件判断：if nos2, _ := g.OrderService.PendingOrderNos(); len(nos2) != 0 {
	if nos2, _ := g.OrderService.PendingOrderNos(); len(nos2) != 0 {
// 执行：t.Fatalf("清扫后待付款清单应为空: %v", nos2)
		t.Fatalf("清扫后待付款清单应为空: %v", nos2)
// 结束当前代码块
	}
// 结束当前代码块
}

// TestPaymentsAndAbnormalQueries 收款记录与异常记录查询路径。
func TestPaymentsAndAbnormalQueries(t *testing.T) {
// 赋值/更新：g, _, clock := newEnv(t)
	g, _, clock := newEnv(t)
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：o, _ := g.OrderService.Create(b1, pid, 1)
	o, _ := g.OrderService.Create(b1, pid, 1)

	// 正常记款 → PaymentsByOrder 1 条
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, "TXN-C1", o.AmountCents)
// 条件判断：if _, err := g.PaymentService.Notify(o.OrderNo, "TXN-C1", o.AmountCents, sign, "{}"); err != nil {
	if _, err := g.PaymentService.Notify(o.OrderNo, "TXN-C1", o.AmountCents, sign, "{}"); err != nil {
// 执行：t.Fatalf("记款失败: %v", err)
		t.Fatalf("记款失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：pays, err := g.PaymentService.PaymentsByOrder(o.OrderNo)
	pays, err := g.PaymentService.PaymentsByOrder(o.OrderNo)
// 条件判断：if err != nil || len(pays) != 1 || pays[0].ChannelTxnID != "TXN-C1" {
	if err != nil || len(pays) != 1 || pays[0].ChannelTxnID != "TXN-C1" {
// 执行：t.Fatalf("收款记录查询错误: %v %d", err, len(pays))
		t.Fatalf("收款记录查询错误: %v %d", err, len(pays))
// 结束当前代码块
	}

	// 迟到支付 → 异常记录，AbnormalByID 命中/未命中
	o2, _ := g.OrderService.Create(b1, pid, 1)
// 执行：clock.Advance(11 * time.Minute)
	clock.Advance(11 * time.Minute)
// 条件判断：if _, err := g.OrderService.CloseExpired(); err != nil {
	if _, err := g.OrderService.CloseExpired(); err != nil {
// 执行：t.Fatalf("超时关单失败: %v", err)
		t.Fatalf("超时关单失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：sign2 := utils.PaySign(testCfg.Shop.PayNotifySecret, o2.OrderNo, "TXN-C2", o2.AmountCents)
	sign2 := utils.PaySign(testCfg.Shop.PayNotifySecret, o2.OrderNo, "TXN-C2", o2.AmountCents)
// 条件判断：if _, err := g.PaymentService.Notify(o2.OrderNo, "TXN-C2", o2.AmountCents, sign2, "{}"); err != nil {
	if _, err := g.PaymentService.Notify(o2.OrderNo, "TXN-C2", o2.AmountCents, sign2, "{}"); err != nil {
// 执行：t.Fatalf("迟到通知失败: %v", err)
		t.Fatalf("迟到通知失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：ab, err := g.PaymentService.AbnormalList("")
	ab, err := g.PaymentService.AbnormalList("")
// 条件判断：if err != nil || len(ab) != 1 {
	if err != nil || len(ab) != 1 {
// 执行：t.Fatalf("异常列表应有 1 条: %v %d", err, len(ab))
		t.Fatalf("异常列表应有 1 条: %v %d", err, len(ab))
// 结束当前代码块
	}
// 赋值/更新：got, err := g.PaymentService.AbnormalByID(ab[0].ID)
	got, err := g.PaymentService.AbnormalByID(ab[0].ID)
// 条件判断：if err != nil || got.ChannelTxnID != "TXN-C2" {
	if err != nil || got.ChannelTxnID != "TXN-C2" {
// 执行：t.Fatalf("按主键查异常失败: %v %+v", err, got)
		t.Fatalf("按主键查异常失败: %v %+v", err, got)
// 结束当前代码块
	}
// 条件判断：if _, err := g.PaymentService.AbnormalByID(999999); !errors.Is(err, ErrAbnormalNotFound) {
	if _, err := g.PaymentService.AbnormalByID(999999); !errors.Is(err, ErrAbnormalNotFound) {
// 执行：t.Fatalf("未知异常主键应返回未找到: %v", err)
		t.Fatalf("未知异常主键应返回未找到: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// TestApplyRefundGuards 退款申请的归属与状态守卫。
func TestApplyRefundGuards(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：b1, b2 := buyerID(t, "buyer1"), buyerID(t, "buyer2")
	b1, b2 := buyerID(t, "buyer1"), buyerID(t, "buyer2")
// 赋值/更新：cs := buyerID(t, "cs1")
	cs := buyerID(t, "cs1")

// 赋值/更新：o, _ := g.OrderService.Create(b1, pid, 1)
	o, _ := g.OrderService.Create(b1, pid, 1)
// 条件判断：if _, err := g.RefundService.Apply(b1, o.OrderNo, "未付款退款"); err == nil {
	if _, err := g.RefundService.Apply(b1, o.OrderNo, "未付款退款"); err == nil {
// 执行：t.Fatal("待付款订单不得申请售前退款")
		t.Fatal("待付款订单不得申请售前退款")
// 结束当前代码块
	}

// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, "TXN-C3", o.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, "TXN-C3", o.AmountCents)
// 条件判断：if _, err := g.PaymentService.Notify(o.OrderNo, "TXN-C3", o.AmountCents, sign, "{}"); err != nil {
	if _, err := g.PaymentService.Notify(o.OrderNo, "TXN-C3", o.AmountCents, sign, "{}"); err != nil {
// 执行：t.Fatalf("记款失败: %v", err)
		t.Fatalf("记款失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.RefundService.Apply(b2, o.OrderNo, "他人订单"); !errors.Is(err, ErrForbidden) {
	if _, err := g.RefundService.Apply(b2, o.OrderNo, "他人订单"); !errors.Is(err, ErrForbidden) {
// 执行：t.Fatalf("非归属买家应被拒绝: %v", err)
		t.Fatalf("非归属买家应被拒绝: %v", err)
// 结束当前代码块
	}

	// 拒绝路径：拒绝后订单回可发货，可正常发货
	if _, err := g.RefundService.Apply(b1, o.OrderNo, "无理由"); err != nil {
// 执行：t.Fatalf("合法退款申请失败: %v", err)
		t.Fatalf("合法退款申请失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if _, err := g.RefundService.Audit(cs, o.OrderNo, false, "证据不足"); err != nil {
	if _, err := g.RefundService.Audit(cs, o.OrderNo, false, "证据不足"); err != nil {
// 执行：t.Fatalf("审核拒绝失败: %v", err)
		t.Fatalf("审核拒绝失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：o3, _ := g.OrderService.Get(o.OrderNo)
	o3, _ := g.OrderService.Get(o.OrderNo)
// 条件判断：if o3.Status != model.OrderStatusPaid {
	if o3.Status != model.OrderStatusPaid {
// 执行：t.Fatalf("拒绝后订单应回待发货: %s", o3.Status)
		t.Fatalf("拒绝后订单应回待发货: %s", o3.Status)
// 结束当前代码块
	}
	// 被拒绝的退款单允许重新申请，复用同一张退款单（整单退款一单一张）
	rf2, err := g.RefundService.Apply(b1, o.OrderNo, "补充理由后再次申请")
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("重新申请应成功: %v", err)
		t.Fatalf("重新申请应成功: %v", err)
// 结束当前代码块
	}
// 赋值/更新：rfOld, _ := g.RefundService.Get(o.OrderNo)
	rfOld, _ := g.RefundService.Get(o.OrderNo)
// 条件判断：if rf2.ID != rfOld.ID || rf2.AuditStatus != model.AuditPending {
	if rf2.ID != rfOld.ID || rf2.AuditStatus != model.AuditPending {
// 执行：t.Fatalf("应复用原退款单并回到待审核: %+v %+v", rf2, rfOld)
		t.Fatalf("应复用原退款单并回到待审核: %+v %+v", rf2, rfOld)
// 结束当前代码块
	}
// 结束当前代码块
}
