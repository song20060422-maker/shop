// 包声明
package service

// 导入依赖包
import (
// 执行："errors"
	"errors"
// 执行："testing"
	"testing"
// 结束当前代码块
)

// TestUtilityStringBranches 覆盖 itoa 正/负/零与 boolLabel 两个分支。
func TestUtilityStringBranches(t *testing.T) {
// 赋值/更新：cases := []struct {
	cases := []struct {
// 执行：in   int
		in   int
// 执行：want string
		want string
// 结束当前代码块
	}{
// 执行：{0, "0"},
		{0, "0"},
// 执行：{7, "7"},
		{7, "7"},
// 执行：{-5, "-5"},
		{-5, "-5"},
// 执行：{123456, "123456"},
		{123456, "123456"},
// 执行：{-987654, "-987654"},
		{-987654, "-987654"},
// 结束当前代码块
	}
// 循环处理：for _, c := range cases {
	for _, c := range cases {
// 条件判断：if got := itoa(c.in); got != c.want {
		if got := itoa(c.in); got != c.want {
// 执行：t.Fatalf("itoa(%d) = %s, want %s", c.in, got, c.want)
			t.Fatalf("itoa(%d) = %s, want %s", c.in, got, c.want)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if boolLabel(true) != "是" || boolLabel(false) != "否" {
	if boolLabel(true) != "是" || boolLabel(false) != "否" {
// 执行：t.Fatal("boolLabel 分支错误")
		t.Fatal("boolLabel 分支错误")
// 结束当前代码块
	}
// 结束当前代码块
}

// TestCatalogGetNotFound 覆盖商品查询的未命中分支。
func TestCatalogGetNotFound(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：_, err := g.CatalogService.Get(999999)
	_, err := g.CatalogService.Get(999999)
// 条件判断：if !errors.Is(err, ErrProductNotFound) {
	if !errors.Is(err, ErrProductNotFound) {
// 执行：t.Fatalf("应返回商品不存在: %v", err)
		t.Fatalf("应返回商品不存在: %v", err)
// 结束当前代码块
	}
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：p, err := g.CatalogService.Get(pid)
	p, err := g.CatalogService.Get(pid)
// 条件判断：if err != nil || p.SKU != "SKU-A" {
	if err != nil || p.SKU != "SKU-A" {
// 执行：t.Fatalf("命中查询失败: %v %+v", err, p)
		t.Fatalf("命中查询失败: %v %+v", err, p)
// 结束当前代码块
	}
// 结束当前代码块
}

// TestListByBuyerEmpty 覆盖买家订单列表空结果分支。
func TestListByBuyerEmpty(t *testing.T) {
// 赋值/更新：g, _, _ := newEnv(t)
	g, _, _ := newEnv(t)
// 赋值/更新：list, err := g.OrderService.ListByBuyer(buyerID(t, "buyer2"), "")
	list, err := g.OrderService.ListByBuyer(buyerID(t, "buyer2"), "")
// 条件判断：if err != nil || len(list) != 0 {
	if err != nil || len(list) != 0 {
// 执行：t.Fatalf("新买家列表应为空: %v %d", err, len(list))
		t.Fatalf("新买家列表应为空: %v %d", err, len(list))
// 结束当前代码块
	}
// 结束当前代码块
}
