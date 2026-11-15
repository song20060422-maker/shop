// 包声明
package router

// 导入依赖包
import (
// 执行："io"
	"io"
// 执行："net/http"
	"net/http"
// 执行："net/http/httptest"
	"net/http/httptest"
// 执行："strings"
	"strings"
// 执行："testing"
	"testing"
// 执行："time"
	"time"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"
// 执行："shop/shop/config"
	"shop/shop/config"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// newServerWithCfg 与 newServer 相同，但允许用例自定义配置（如 system.admin-target）。
func newServerWithCfg(t *testing.T, cfg config.Config) *gin.Engine {
// 执行：t.Helper()
	t.Helper()
// 执行：resetTables(t)
	resetTables(t)
// 执行：gin.SetMode(gin.TestMode)
	gin.SetMode(gin.TestMode)
// 赋值/更新：r := gin.New()
	r := gin.New()
// 赋值/更新：ch := service.NewSimulatedChannel(true)
	ch := service.NewSimulatedChannel(true)
// 赋值/更新：g := service.NewServiceGroup(service.Deps{DB: testDB, Now: time.Now, Cfg: cfg, Channel: ch})
	g := service.NewServiceGroup(service.Deps{DB: testDB, Now: time.Now, Cfg: cfg, Channel: ch})
// 执行：InitShopRouter(r.Group(""), g, cfg, time.Now)
	InitShopRouter(r.Group(""), g, cfg, time.Now)
// 返回：return r
	return r
// 结束当前代码块
}

// TestAdminPagesServed 买家端页随服务托管；根路径重定向；后台三页已迁移 GVA 管理端（404 提示去向）。
func TestAdminPagesServed(t *testing.T) {
// 赋值/更新：r, _ := newServer(t, true)
	r, _ := newServer(t, true)
	// 买家页：/shop、/shop/ 直接返回 HTML
	for _, p := range []string{"/shop", "/shop/"} {
// 赋值/更新：w := doJSON(r, "GET", p, "", "", nil)
		w := doJSON(r, "GET", p, "", "", nil)
// 条件判断：if w.Code != http.StatusOK {
		if w.Code != http.StatusOK {
// 执行：t.Fatalf("%s 应 200, got %d", p, w.Code)
			t.Fatalf("%s 应 200, got %d", p, w.Code)
// 结束当前代码块
		}
// 条件判断：if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
// 执行：t.Fatalf("%s Content-Type 应为 html: %s", p, ct)
			t.Fatalf("%s Content-Type 应为 html: %s", p, ct)
// 结束当前代码块
		}
// 条件判断：if !strings.Contains(w.Body.String(), "买家") {
		if !strings.Contains(w.Body.String(), "买家") {
// 执行：t.Fatalf("%s 页面缺少标记", p)
			t.Fatalf("%s 页面缺少标记", p)
// 结束当前代码块
		}
// 结束当前代码块
	}
	// 根路径与 /index.html 重定向到 /shop（历史入口兼容）
	for _, p := range []string{"/", "/index.html"} {
// 赋值/更新：w := doJSON(r, "GET", p, "", "", nil)
		w := doJSON(r, "GET", p, "", "", nil)
// 条件判断：if w.Code != http.StatusFound || w.Header().Get("Location") != "/shop" {
		if w.Code != http.StatusFound || w.Header().Get("Location") != "/shop" {
// 执行：t.Fatalf("%s 应 302 到 /shop, got %d %q", p, w.Code, w.Header().Get("Location"))
			t.Fatalf("%s 应 302 到 /shop, got %d %q", p, w.Code, w.Header().Get("Location"))
// 结束当前代码块
		}
// 结束当前代码块
	}
	// 后台三页已迁移：返回 404 + 明确去向提示，不做静默 404
	for _, c := range []struct{ path, marker string }{
// 执行：{"/ops", "运营核查"},
		{"/ops", "运营核查"},
// 执行：{"/warehouse", "仓库发货"},
		{"/warehouse", "仓库发货"},
// 执行：{"/cs", "客服退款"},
		{"/cs", "客服退款"},
// 结束当前代码块
	} {
// 赋值/更新：w := doJSON(r, "GET", c.path, "", "", nil)
		w := doJSON(r, "GET", c.path, "", "", nil)
// 条件判断：if w.Code != http.StatusNotFound {
		if w.Code != http.StatusNotFound {
// 执行：t.Fatalf("%s 应 404, got %d", c.path, w.Code)
			t.Fatalf("%s 应 404, got %d", c.path, w.Code)
// 结束当前代码块
		}
// 条件判断：if !strings.Contains(w.Body.String(), "GVA 管理端") || !strings.Contains(w.Body.String(), c.marker) {
		if !strings.Contains(w.Body.String(), "GVA 管理端") || !strings.Contains(w.Body.String(), c.marker) {
// 执行：t.Fatalf("%s 提示缺少去向说明: %s", c.path, w.Body.String())
			t.Fatalf("%s 提示缺少去向说明: %s", c.path, w.Body.String())
// 结束当前代码块
		}
// 结束当前代码块
	}
// 结束当前代码块
}

// TestGVAReverseProxy 后台独立页通过 /gva/* 访问管理端：剥掉 /gva 前缀并透传 x-token。
func TestGVAReverseProxy(t *testing.T) {
// 声明：var gotPath, gotToken, gotBody, gotMethod string
	var gotPath, gotToken, gotBody, gotMethod string
// 赋值/更新：admin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	admin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 赋值/更新：gotPath = req.URL.Path
		gotPath = req.URL.Path
// 赋值/更新：gotMethod = req.Method
		gotMethod = req.Method
// 赋值/更新：gotToken = req.Header.Get("x-token")
		gotToken = req.Header.Get("x-token")
// 赋值/更新：b, _ := io.ReadAll(req.Body)
		b, _ := io.ReadAll(req.Body)
// 赋值/更新：gotBody = string(b)
		gotBody = string(b)
// 执行：w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Type", "application/json")
// 赋值/更新：_, _ = w.Write([]byte(`{"code":0,"data":{"token":"T"},"msg":"ok"}`))
		_, _ = w.Write([]byte(`{"code":0,"data":{"token":"T"},"msg":"ok"}`))
// 结束当前代码块
	}))
// 异步/延迟执行：defer admin.Close()
	defer admin.Close()

// 赋值/更新：cfg := testCfg
	cfg := testCfg
// 赋值/更新：cfg.System.AdminTarget = admin.URL
	cfg.System.AdminTarget = admin.URL
// 赋值/更新：r := newServerWithCfg(t, cfg)
	r := newServerWithCfg(t, cfg)
	// ReverseProxy 需要真实 ResponseWriter（httptest.ResponseRecorder 缺 CloseNotifier/Flusher），
	// 因此这里起真实监听再请求。
	srv := httptest.NewServer(r)
// 异步/延迟执行：defer srv.Close()
	defer srv.Close()

// 赋值/更新：resp, err := http.Post(srv.URL+"/gva/base/login", "application/json", strings.NewReader(`{"username":"admin","...
	resp, err := http.Post(srv.URL+"/gva/base/login", "application/json", strings.NewReader(`{"username":"admin","password":"123456"}`))
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("反代登录请求失败: %v", err)
		t.Fatalf("反代登录请求失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：body, _ := io.ReadAll(resp.Body)
	body, _ := io.ReadAll(resp.Body)
// 执行：resp.Body.Close()
	resp.Body.Close()
// 条件判断：if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), `"code":0`) {
	if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), `"code":0`) {
// 执行：t.Fatalf("反代登录应透传响应: %d %s", resp.StatusCode, string(body))
		t.Fatalf("反代登录应透传响应: %d %s", resp.StatusCode, string(body))
// 结束当前代码块
	}
// 条件判断：if gotPath != "/base/login" || gotMethod != "POST" {
	if gotPath != "/base/login" || gotMethod != "POST" {
// 执行：t.Fatalf("/gva 前缀未剥离: path=%s method=%s", gotPath, gotMethod)
		t.Fatalf("/gva 前缀未剥离: path=%s method=%s", gotPath, gotMethod)
// 结束当前代码块
	}
// 条件判断：if !strings.Contains(gotBody, "admin") {
	if !strings.Contains(gotBody, "admin") {
// 执行：t.Fatalf("请求体未透传: %s", gotBody)
		t.Fatalf("请求体未透传: %s", gotBody)
// 结束当前代码块
	}

// 赋值/更新：req, _ := http.NewRequest("GET", srv.URL+"/gva/shop/ops/summary", nil)
	req, _ := http.NewRequest("GET", srv.URL+"/gva/shop/ops/summary", nil)
// 执行：req.Header.Set("x-token", "TOKEN-X")
	req.Header.Set("x-token", "TOKEN-X")
// 赋值/更新：resp2, err := http.DefaultClient.Do(req)
	resp2, err := http.DefaultClient.Do(req)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("反代查询失败: %v", err)
		t.Fatalf("反代查询失败: %v", err)
// 结束当前代码块
	}
// 执行：resp2.Body.Close()
	resp2.Body.Close()
// 条件判断：if resp2.StatusCode != http.StatusOK {
	if resp2.StatusCode != http.StatusOK {
// 执行：t.Fatalf("反代查询应 200: %d", resp2.StatusCode)
		t.Fatalf("反代查询应 200: %d", resp2.StatusCode)
// 结束当前代码块
	}
// 条件判断：if gotPath != "/shop/ops/summary" {
	if gotPath != "/shop/ops/summary" {
// 执行：t.Fatalf("管理端接口路径错误: %s", gotPath)
		t.Fatalf("管理端接口路径错误: %s", gotPath)
// 结束当前代码块
	}
// 条件判断：if gotToken != "TOKEN-X" {
	if gotToken != "TOKEN-X" {
// 执行：t.Fatalf("x-token 未透传: %q", gotToken)
		t.Fatalf("x-token 未透传: %q", gotToken)
// 结束当前代码块
	}
// 结束当前代码块
}

// TestGVAReverseProxyInvalidTarget 目标地址非法时不注册反代，服务不崩溃。
func TestGVAReverseProxyInvalidTarget(t *testing.T) {
// 赋值/更新：cfg := testCfg
	cfg := testCfg
// 赋值/更新：cfg.System.AdminTarget = "://bad-target"
	cfg.System.AdminTarget = "://bad-target"
// 赋值/更新：r := newServerWithCfg(t, cfg)
	r := newServerWithCfg(t, cfg)

// 条件判断：if w := doJSON(r, "GET", "/shop", "", "", nil); w.Code != http.StatusOK {
	if w := doJSON(r, "GET", "/shop", "", "", nil); w.Code != http.StatusOK {
// 执行：t.Fatalf("买家页仍应可访问: %d", w.Code)
		t.Fatalf("买家页仍应可访问: %d", w.Code)
// 结束当前代码块
	}
// 条件判断：if w := doJSON(r, "GET", "/gva/health", "", "", nil); w.Code == http.StatusOK {
	if w := doJSON(r, "GET", "/gva/health", "", "", nil); w.Code == http.StatusOK {
// 执行：t.Fatalf("非法目标不应注册反代: %d", w.Code)
		t.Fatalf("非法目标不应注册反代: %d", w.Code)
// 结束当前代码块
	}
// 结束当前代码块
}
