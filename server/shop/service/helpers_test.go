// 包声明
package service

// 导入依赖包
import (
// 执行："errors"
	"errors"
// 执行："strconv"
	"strconv"
// 执行："testing"
	"testing"
// 执行："time"
	"time"

// 执行："shop/shop/model"
	"shop/shop/model"
// 结束当前代码块
)

// 函数定义：func TestBizErrorHelpers(t *testing.T) {
func TestBizErrorHelpers(t *testing.T) {
// 赋值/更新：e := NewBizError(1003, "可售数量不足")
	e := NewBizError(1003, "可售数量不足")
// 条件判断：if e.Error() != "可售数量不足" {
	if e.Error() != "可售数量不足" {
// 执行：t.Fatalf("错误文案错误: %s", e.Error())
		t.Fatalf("错误文案错误: %s", e.Error())
// 结束当前代码块
	}
// 条件判断：if !IsBiz(e) {
	if !IsBiz(e) {
// 执行：t.Fatal("BizError 应识别为业务错误")
		t.Fatal("BizError 应识别为业务错误")
// 结束当前代码块
	}
// 条件判断：if IsBiz(errors.New("plain")) {
	if IsBiz(errors.New("plain")) {
// 执行：t.Fatal("普通错误不应识别为业务错误")
		t.Fatal("普通错误不应识别为业务错误")
// 结束当前代码块
	}
// 赋值/更新：code, msg := BizCode(e)
	code, msg := BizCode(e)
// 条件判断：if code != 1003 || msg != "可售数量不足" {
	if code != 1003 || msg != "可售数量不足" {
// 执行：t.Fatalf("错误码解析错误: %d %s", code, msg)
		t.Fatalf("错误码解析错误: %d %s", code, msg)
// 结束当前代码块
	}
// 条件判断：if BizCodeMsg(e) != "可售数量不足" {
	if BizCodeMsg(e) != "可售数量不足" {
// 执行：t.Fatal("BizCodeMsg 错误")
		t.Fatal("BizCodeMsg 错误")
// 结束当前代码块
	}
	// 未知错误统一 9000，保留底层原因
	code, msg = BizCode(errors.New("db down"))
// 条件判断：if code != 9000 || msg == "" {
	if code != 9000 || msg == "" {
// 执行：t.Fatalf("未知错误应映射 9000: %d %s", code, msg)
		t.Fatalf("未知错误应映射 9000: %d %s", code, msg)
// 结束当前代码块
	}
	// nil 错误按未知错误处理（9000），提示非空但不泄露内部细节
	if code, _ := BizCode(nil); code != 9000 {
// 执行：t.Fatalf("nil 错误应映射 9000: %d", code)
		t.Fatalf("nil 错误应映射 9000: %d", code)
// 结束当前代码块
	}
// 条件判断：if BizCodeMsg(errors.New("")) == "" {
	if BizCodeMsg(errors.New("")) == "" {
// 执行：t.Fatal("未知错误应给出兜底提示")
		t.Fatal("未知错误应给出兜底提示")
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestStatusLabelsAndStateErrors(t *testing.T) {
func TestStatusLabelsAndStateErrors(t *testing.T) {
// 赋值/更新：cases := map[string]string{
	cases := map[string]string{
// 执行：model.OrderStatusPendingPayment:  "待付款",
		model.OrderStatusPendingPayment:  "待付款",
// 执行：model.OrderStatusCancelled:       "取消",
		model.OrderStatusCancelled:       "取消",
// 执行：model.OrderStatusClosed:          "关闭",
		model.OrderStatusClosed:          "关闭",
// 执行：model.OrderStatusPaid:            "付款",
		model.OrderStatusPaid:            "付款",
// 执行：model.OrderStatusShipped:         "发货",
		model.OrderStatusShipped:         "发货",
// 执行：model.OrderStatusRefundRequested: "退款申请",
		model.OrderStatusRefundRequested: "退款申请",
// 执行：model.OrderStatusRefundFailed:    "退款失败",
		model.OrderStatusRefundFailed:    "退款失败",
// 执行：model.OrderStatusRefunded:        "退款",
		model.OrderStatusRefunded:        "退款",
// 执行："unknown":                        "unknown",
		"unknown":                        "unknown",
// 结束当前代码块
	}
// 循环处理：for status, want := range cases {
	for status, want := range cases {
// 条件判断：if got := statusLabel(status); got != want {
		if got := statusLabel(status); got != want {
// 执行：t.Fatalf("状态文案错误 %s: %s", status, got)
			t.Fatalf("状态文案错误 %s: %s", status, got)
// 结束当前代码块
		}
// 结束当前代码块
	}

// 条件判断：if err := shipStateError(model.OrderStatusRefundRequested); !errors.Is(err, ErrRefundInProgress) {
	if err := shipStateError(model.OrderStatusRefundRequested); !errors.Is(err, ErrRefundInProgress) {
// 执行：t.Fatalf("退款申请中应禁止发货: %v", err)
		t.Fatalf("退款申请中应禁止发货: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := shipStateError(model.OrderStatusRefunded); !errors.Is(err, ErrRefunded) {
	if err := shipStateError(model.OrderStatusRefunded); !errors.Is(err, ErrRefunded) {
// 执行：t.Fatalf("已退款应禁止发货: %v", err)
		t.Fatalf("已退款应禁止发货: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := shipStateError(model.OrderStatusRefundFailed); err == nil {
	if err := shipStateError(model.OrderStatusRefundFailed); err == nil {
// 执行：t.Fatal("退款失败应禁止发货")
		t.Fatal("退款失败应禁止发货")
// 结束当前代码块
	}
// 条件判断：if err := shipStateError(model.OrderStatusPendingPayment); !errors.Is(err, ErrNotPaid) {
	if err := shipStateError(model.OrderStatusPendingPayment); !errors.Is(err, ErrNotPaid) {
// 执行：t.Fatalf("待付款应禁止发货: %v", err)
		t.Fatalf("待付款应禁止发货: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := shipStateError(model.OrderStatusClosed); !errors.Is(err, ErrOrderClosed) {
	if err := shipStateError(model.OrderStatusClosed); !errors.Is(err, ErrOrderClosed) {
// 执行：t.Fatalf("已关闭应禁止发货: %v", err)
		t.Fatalf("已关闭应禁止发货: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := shipStateError(model.OrderStatusShipped); !errors.Is(err, ErrAlreadyShipped) {
	if err := shipStateError(model.OrderStatusShipped); !errors.Is(err, ErrAlreadyShipped) {
// 执行：t.Fatalf("已发货应禁止重复发货: %v", err)
		t.Fatalf("已发货应禁止重复发货: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestStateErrorReadsCurrentStatus(t *testing.T) {
func TestStateErrorReadsCurrentStatus(t *testing.T) {
// 赋值/更新：g, _, clock := newEnv(t)
	g, _, clock := newEnv(t)
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")

// 赋值/更新：o, err := g.OrderService.Create(b1, pid, 1)
	o, err := g.OrderService.Create(b1, pid, 1)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
	// 待付款 → 状态不允许该操作
	if err := stateError(testDB, o.OrderNo); !errors.Is(err, ErrOrderNotPayable) {
// 执行：t.Fatalf("待付款订单状态错误: %v", err)
		t.Fatalf("待付款订单状态错误: %v", err)
// 结束当前代码块
	}
	// 未知订单
	if err := stateError(testDB, "SO-NOT-EXIST"); !errors.Is(err, ErrOrderNotFound) {
// 执行：t.Fatalf("未知订单应报错: %v", err)
		t.Fatalf("未知订单应报错: %v", err)
// 结束当前代码块
	}
	// 已关闭
	clock.Advance(11 * time.Minute)
// 条件判断：if _, err := g.OrderService.CloseExpired(); err != nil {
	if _, err := g.OrderService.CloseExpired(); err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if err := stateError(testDB, o.OrderNo); !errors.Is(err, ErrOrderClosed) {
	if err := stateError(testDB, o.OrderNo); !errors.Is(err, ErrOrderClosed) {
// 执行：t.Fatalf("已关闭订单状态错误: %v", err)
		t.Fatalf("已关闭订单状态错误: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestRefundFindByRef(t *testing.T) {
func TestRefundFindByRef(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：b1 := buyerID(t, "buyer1")
	b1 := buyerID(t, "buyer1")
// 赋值/更新：orderNo, _, _ := payOrder(t, g, b1, 1)
	orderNo, _, _ := payOrder(t, g, b1, 1)
// 赋值/更新：rf, err := g.RefundService.Apply(b1, orderNo, "测试")
	rf, err := g.RefundService.Apply(b1, orderNo, "测试")
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
	// 空引用
	if _, err := g.RefundService.FindByRef("  "); !errors.Is(err, ErrRefundNotFound) {
// 执行：t.Fatalf("空引用应报错: %v", err)
		t.Fatalf("空引用应报错: %v", err)
// 结束当前代码块
	}
	// 退款单号
	if got, err := g.RefundService.FindByRef(rf.RefundNo); err != nil || got.OrderNo != orderNo {
// 执行：t.Fatalf("按退款单号查询失败: %v %+v", err, got)
		t.Fatalf("按退款单号查询失败: %v %+v", err, got)
// 结束当前代码块
	}
	// 订单号
	if got, err := g.RefundService.FindByRef(orderNo); err != nil || got.RefundNo != rf.RefundNo {
// 执行：t.Fatalf("按订单号查询失败: %v %+v", err, got)
		t.Fatalf("按订单号查询失败: %v %+v", err, got)
// 结束当前代码块
	}
	// 数字主键
	if got, err := g.RefundService.FindByRef(strconv.FormatUint(rf.ID, 10)); err != nil || got.ID != rf.ID {
// 执行：t.Fatalf("按主键查询失败: %v %+v", err, got)
		t.Fatalf("按主键查询失败: %v %+v", err, got)
// 结束当前代码块
	}
	// 不存在
	if _, err := g.RefundService.FindByRef("999999"); !errors.Is(err, ErrRefundNotFound) {
// 执行：t.Fatalf("不存在的主键应报错: %v", err)
		t.Fatalf("不存在的主键应报错: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}
