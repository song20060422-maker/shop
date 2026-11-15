// 包声明
package middleware

// 导入依赖包
import (
// 执行："net/http"
	"net/http"
// 执行："net/http/httptest"
	"net/http/httptest"
// 执行："testing"
	"testing"
// 执行："time"
	"time"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// 函数定义：func newCtx(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
func newCtx(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
// 执行：gin.SetMode(gin.TestMode)
	gin.SetMode(gin.TestMode)
// 赋值/更新：w := httptest.NewRecorder()
	w := httptest.NewRecorder()
// 赋值/更新：c, _ := gin.CreateTestContext(w)
	c, _ := gin.CreateTestContext(w)
// 赋值/更新：c.Request = httptest.NewRequest(method, path, nil)
	c.Request = httptest.NewRequest(method, path, nil)
// 返回：return c, w
	return c, w
// 结束当前代码块
}

// 函数定义：func TestJWTBranches(t *testing.T) {
func TestJWTBranches(t *testing.T) {
// 赋值/更新：now := func() time.Time { return time.Date(2026, 11, 12, 10, 0, 0, 0, time.Local) }
	now := func() time.Time { return time.Date(2026, 11, 12, 10, 0, 0, 0, time.Local) }
// 赋值/更新：secret := "unit-secret"
	secret := "unit-secret"
// 赋值/更新：token, err := utils.SignJWT(secret, utils.Claims{
	token, err := utils.SignJWT(secret, utils.Claims{
// 执行：UserID: 7, Username: "buyer1", Role: model.RoleBuyer,
		UserID: 7, Username: "buyer1", Role: model.RoleBuyer,
// 执行：ExpireAt: now().Add(time.Hour).Unix(), IssuedAt: now().Unix(),
		ExpireAt: now().Add(time.Hour).Unix(), IssuedAt: now().Unix(),
// 结束当前代码块
	})
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}

	// 缺少令牌
	c, w := newCtx("GET", "/api/me")
// 执行：JWT(secret, now)(c)
	JWT(secret, now)(c)
// 条件判断：if w.Code != http.StatusUnauthorized {
	if w.Code != http.StatusUnauthorized {
// 执行：t.Fatalf("无令牌应 401: %d", w.Code)
		t.Fatalf("无令牌应 401: %d", w.Code)
// 结束当前代码块
	}

	// 伪造令牌
	c, w = newCtx("GET", "/api/me")
// 执行：c.Request.Header.Set("Authorization", "Bearer forged")
	c.Request.Header.Set("Authorization", "Bearer forged")
// 执行：JWT(secret, now)(c)
	JWT(secret, now)(c)
// 条件判断：if w.Code != http.StatusUnauthorized {
	if w.Code != http.StatusUnauthorized {
// 执行：t.Fatalf("伪造令牌应 401: %d", w.Code)
		t.Fatalf("伪造令牌应 401: %d", w.Code)
// 结束当前代码块
	}

	// Bearer 令牌
	c, w = newCtx("GET", "/api/me")
// 执行：c.Request.Header.Set("Authorization", "Bearer "+token)
	c.Request.Header.Set("Authorization", "Bearer "+token)
// 执行：JWT(secret, now)(c)
	JWT(secret, now)(c)
// 条件判断：if claims := ClaimsOf(c); claims.UserID != 7 || claims.Role != model.RoleBuyer {
	if claims := ClaimsOf(c); claims.UserID != 7 || claims.Role != model.RoleBuyer {
// 执行：t.Fatalf("令牌载荷错误: %+v", claims)
		t.Fatalf("令牌载荷错误: %+v", claims)
// 结束当前代码块
	}

	// Cookie 令牌
	c, w = newCtx("GET", "/api/me")
// 执行：c.Request.AddCookie(&http.Cookie{Name: "token", Value: token})
	c.Request.AddCookie(&http.Cookie{Name: "token", Value: token})
// 执行：JWT(secret, now)(c)
	JWT(secret, now)(c)
// 条件判断：if claims := ClaimsOf(c); claims.Username != "buyer1" {
	if claims := ClaimsOf(c); claims.Username != "buyer1" {
// 执行：t.Fatalf("Cookie 令牌未生效: %+v", claims)
		t.Fatalf("Cookie 令牌未生效: %+v", claims)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestRequireRoleAndHelpers(t *testing.T) {
func TestRequireRoleAndHelpers(t *testing.T) {
	// 无 claims 时拒绝
	c, w := newCtx("GET", "/api/ops/summary")
// 执行：RequireRole(model.RoleOps)(c)
	RequireRole(model.RoleOps)(c)
// 条件判断：if w.Code != http.StatusForbidden {
	if w.Code != http.StatusForbidden {
// 执行：t.Fatalf("无角色应 403: %d", w.Code)
		t.Fatalf("无角色应 403: %d", w.Code)
// 结束当前代码块
	}

	// 角色不符
	c, w = newCtx("GET", "/api/ops/summary")
// 执行：c.Set(claimsKey, utils.Claims{UserID: 1, Role: model.RoleBuyer})
	c.Set(claimsKey, utils.Claims{UserID: 1, Role: model.RoleBuyer})
// 执行：RequireRole(model.RoleOps)(c)
	RequireRole(model.RoleOps)(c)
// 条件判断：if w.Code != http.StatusForbidden {
	if w.Code != http.StatusForbidden {
// 执行：t.Fatalf("角色不符应 403: %d", w.Code)
		t.Fatalf("角色不符应 403: %d", w.Code)
// 结束当前代码块
	}

	// 角色命中（使用新的记录器，避免复用上一次的 403 结果）
	c, w2 := newCtx("GET", "/api/ops/summary")
// 执行：c.Set(claimsKey, utils.Claims{UserID: 1, Role: model.RoleOps})
	c.Set(claimsKey, utils.Claims{UserID: 1, Role: model.RoleOps})
// 执行：RequireRole(model.RoleOps)(c)
	RequireRole(model.RoleOps)(c)
// 条件判断：if w2.Code != http.StatusOK {
	if w2.Code != http.StatusOK {
// 执行：t.Fatalf("角色命中不应被拦截: %d", w2.Code)
		t.Fatalf("角色命中不应被拦截: %d", w2.Code)
// 结束当前代码块
	}

	// IsBuyer 正例
	c, _ = newCtx("GET", "/api/me")
// 执行：c.Set(claimsKey, utils.Claims{UserID: 1, Role: model.RoleBuyer})
	c.Set(claimsKey, utils.Claims{UserID: 1, Role: model.RoleBuyer})
// 条件判断：if !IsBuyer(c) {
	if !IsBuyer(c) {
// 执行：t.Fatal("买家判定失败")
		t.Fatal("买家判定失败")
// 结束当前代码块
	}

	// 空上下文 ClaimsOf 返回零值
	c, _ = newCtx("GET", "/api/me")
// 条件判断：if claims := ClaimsOf(c); claims.UserID != 0 {
	if claims := ClaimsOf(c); claims.UserID != 0 {
// 执行：t.Fatalf("空 claims 应返回零值: %+v", claims)
		t.Fatalf("空 claims 应返回零值: %+v", claims)
// 结束当前代码块
	}
// 结束当前代码块
}
