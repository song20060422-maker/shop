// 包声明
package router

// 导入依赖包
import (
// 执行："net/http"
	"net/http"
// 执行："testing"
	"testing"
// 结束当前代码块
)

// TestDemoModeGatesTestHooks 安全整改回归：
// 测试钩子（渠道退款开关、强制过期、手动关单）与账号列表、演示重置
// 仅在 shop.demo-mode=true 时可用；生产配置下必须全部 403。
func TestDemoModeGatesTestHooks(t *testing.T) {
// 赋值/更新：hooks := []struct {
	hooks := []struct {
// 执行：method, path, body string
		method, path, body string
// 结束当前代码块
	}{
// 执行：{"POST", "/api/test/expire-orders", ""},
		{"POST", "/api/test/expire-orders", ""},
// 执行：{"POST", "/api/pay/channel-mode", `{"refund_result":"SUCCESS"}`},
		{"POST", "/api/pay/channel-mode", `{"refund_result":"SUCCESS"}`},
// 执行：{"POST", "/api/jobs/close-expired", ""},
		{"POST", "/api/jobs/close-expired", ""},
// 执行：{"GET", "/api/users", ""},
		{"GET", "/api/users", ""},
// 执行：{"POST", "/api/demo/reset", ""},
		{"POST", "/api/demo/reset", ""},
// 结束当前代码块
	}

// 赋值/更新：prod, _ := newServer(t, false) // 生产配置：demo-mode=false
	prod, _ := newServer(t, false) // 生产配置：demo-mode=false
// 循环处理：for _, h := range hooks {
	for _, h := range hooks {
// 赋值/更新：rr := doJSON(prod, h.method, h.path, h.body, "", nil)
		rr := doJSON(prod, h.method, h.path, h.body, "", nil)
// 条件判断：if rr.Code != http.StatusForbidden {
		if rr.Code != http.StatusForbidden {
// 执行：t.Fatalf("生产配置下 %s %s 应 403，实际 %d: %s", h.method, h.path, rr.Code, rr.Body.String())
			t.Fatalf("生产配置下 %s %s 应 403，实际 %d: %s", h.method, h.path, rr.Code, rr.Body.String())
// 结束当前代码块
		}
// 结束当前代码块
	}

// 赋值/更新：demo, _ := newServer(t, true) // 演示配置：demo-mode=true
	demo, _ := newServer(t, true) // 演示配置：demo-mode=true
// 循环处理：for _, h := range hooks {
	for _, h := range hooks {
// 赋值/更新：rr := doJSON(demo, h.method, h.path, h.body, "", nil)
		rr := doJSON(demo, h.method, h.path, h.body, "", nil)
// 条件判断：if rr.Code >= 400 {
		if rr.Code >= 400 {
// 执行：t.Fatalf("演示配置下 %s %s 应可用，实际 %d: %s", h.method, h.path, rr.Code, rr.Body.String())
			t.Fatalf("演示配置下 %s %s 应可用，实际 %d: %s", h.method, h.path, rr.Code, rr.Body.String())
// 结束当前代码块
		}
// 结束当前代码块
	}
// 结束当前代码块
}

// TestBuyerRoleGate 安全整改回归（I-02）：
// 买家区路由仅 buyer 角色可用，员工角色登录后调用必须 403。
func TestBuyerRoleGate(t *testing.T) {
// 赋值/更新：r, _ := newServer(t, true)
	r, _ := newServer(t, true)

// 赋值/更新：opsToken := login(t, r, "ops1")
	opsToken := login(t, r, "ops1")
// 赋值/更新：whToken := login(t, r, "wh1")
	whToken := login(t, r, "wh1")
// 赋值/更新：csToken := login(t, r, "cs1")
	csToken := login(t, r, "cs1")
// 赋值/更新：buyerToken := login(t, r, "buyer1")
	buyerToken := login(t, r, "buyer1")

// 循环处理：for _, tc := range []struct {
	for _, tc := range []struct {
// 执行：name, token string
		name, token string
// 结束当前代码块
	}{
// 执行：{"ops", opsToken},
		{"ops", opsToken},
// 执行：{"warehouse", whToken},
		{"warehouse", whToken},
// 执行：{"cs", csToken},
		{"cs", csToken},
// 结束当前代码块
	} {
// 赋值/更新：rr := doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":1}`, tc.token, nil)
		rr := doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":1}`, tc.token, nil)
// 条件判断：if rr.Code != http.StatusForbidden {
		if rr.Code != http.StatusForbidden {
// 执行：t.Fatalf("%s 角色下单应 403，实际 %d: %s", tc.name, rr.Code, rr.Body.String())
			t.Fatalf("%s 角色下单应 403，实际 %d: %s", tc.name, rr.Code, rr.Body.String())
// 结束当前代码块
		}
// 赋值/更新：rr = doJSON(r, "GET", "/api/orders", "", tc.token, nil)
		rr = doJSON(r, "GET", "/api/orders", "", tc.token, nil)
// 条件判断：if rr.Code != http.StatusForbidden {
		if rr.Code != http.StatusForbidden {
// 执行：t.Fatalf("%s 角色查买家订单应 403，实际 %d", tc.name, rr.Code)
			t.Fatalf("%s 角色查买家订单应 403，实际 %d", tc.name, rr.Code)
// 结束当前代码块
		}
// 结束当前代码块
	}

// 赋值/更新：rr := doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":1}`, buyerToken, nil)
	rr := doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":1}`, buyerToken, nil)
// 条件判断：if rr.Code != http.StatusOK {
	if rr.Code != http.StatusOK {
// 执行：t.Fatalf("买家下单应成功: %d %s", rr.Code, rr.Body.String())
		t.Fatalf("买家下单应成功: %d %s", rr.Code, rr.Body.String())
// 结束当前代码块
	}
// 结束当前代码块
}
