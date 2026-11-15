// Package middleware 提供 JWT 鉴权与角色校验中间件。
package middleware

// 导入依赖包
import (
// 执行："net/http"
	"net/http"
// 执行："strings"
	"strings"
// 执行："time"
	"time"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/model/response"
	"shop/shop/model/response"
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// claimsKey gin.Context 中存放已验签 Claims 的键名，RequireRole 从这里取角色。

// 声明：const claimsKey = "shopClaims"
const claimsKey = "shopClaims"

// JWT 校验令牌，支持 Authorization: Bearer 与 token Cookie（页面演示免传头）。

// 函数定义：func JWT(secret string, now func() time.Time) gin.HandlerFunc {
func JWT(secret string, now func() time.Time) gin.HandlerFunc {
// 返回：return func(c *gin.Context) {
	return func(c *gin.Context) {
// 赋值/更新：token := ""
		token := ""
// 条件判断：if h := c.GetHeader("Authorization"); strings.HasPrefix(h, "Bearer ") {
		if h := c.GetHeader("Authorization"); strings.HasPrefix(h, "Bearer ") {
// 赋值/更新：token = strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
			token = strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
// 结束当前代码块
		}
// 条件判断：if token == "" {
		if token == "" {
// 条件判断：if ck, err := c.Cookie("token"); err == nil {
			if ck, err := c.Cookie("token"); err == nil {
// 赋值/更新：token = ck
				token = ck
// 结束当前代码块
			}
// 结束当前代码块
		}
// 条件判断：if token == "" {
		if token == "" {
// 执行：c.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(401, "请先登录"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(401, "请先登录"))
// 返回：return
			return
// 结束当前代码块
		}
// 赋值/更新：claims, err := utils.ParseJWT(secret, token, now())
		claims, err := utils.ParseJWT(secret, token, now())
// 条件判断：if err != nil {
		if err != nil {
// 执行：c.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(401, "登录状态已失效，请重新登录"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(401, "登录状态已失效，请重新登录"))
// 返回：return
			return
// 结束当前代码块
		}
// 执行：c.Set(claimsKey, claims)
		c.Set(claimsKey, claims)
// 执行：c.Next()
		c.Next()
// 结束当前代码块
	}
// 结束当前代码块
}

// RequireRole 限制角色访问；买家越权访问员工接口时返回明确的 403。

// 函数定义：func RequireRole(roles ...string) gin.HandlerFunc {
func RequireRole(roles ...string) gin.HandlerFunc {
// 赋值/更新：allow := make(map[string]bool, len(roles))
	allow := make(map[string]bool, len(roles))
// 循环处理：for _, r := range roles {
	for _, r := range roles {
// 赋值/更新：allow[r] = true
		allow[r] = true
// 结束当前代码块
	}
// 返回：return func(c *gin.Context) {
	return func(c *gin.Context) {
// 赋值/更新：claims := ClaimsOf(c)
		claims := ClaimsOf(c)
// 条件判断：if !allow[claims.Role] {
		if !allow[claims.Role] {
// 执行：c.AbortWithStatusJSON(http.StatusForbidden, response.Fail(2002, "当前角色无权执行该操作"))
			c.AbortWithStatusJSON(http.StatusForbidden, response.Fail(2002, "当前角色无权执行该操作"))
// 返回：return
			return
// 结束当前代码块
		}
// 执行：c.Next()
		c.Next()
// 结束当前代码块
	}
// 结束当前代码块
}

// ClaimsOf 读取当前请求的令牌载荷

// 函数定义：func ClaimsOf(c *gin.Context) utils.Claims {
func ClaimsOf(c *gin.Context) utils.Claims {
// 条件判断：if v, ok := c.Get(claimsKey); ok {
	if v, ok := c.Get(claimsKey); ok {
// 条件判断：if claims, ok := v.(utils.Claims); ok {
		if claims, ok := v.(utils.Claims); ok {
// 返回：return claims
			return claims
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return utils.Claims{}
	return utils.Claims{}
// 结束当前代码块
}

// IsBuyer 便于业务层判定

// 函数定义：func IsBuyer(c *gin.Context) bool { return ClaimsOf(c).Role == model.RoleBuyer }
func IsBuyer(c *gin.Context) bool { return ClaimsOf(c).Role == model.RoleBuyer }
