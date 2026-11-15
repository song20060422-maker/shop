// 包声明
package api

// 导入依赖包
import (
// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/middleware"
	"shop/shop/middleware"
// 执行："shop/shop/model/request"
	"shop/shop/model/request"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// AuthApi 登录与身份接口

// 类型定义：type AuthApi struct {
type AuthApi struct {
// 执行：s *service.ServiceGroup
	s *service.ServiceGroup
// 结束当前代码块
}

// Login 账号密码登录，返回 JWT 与用户信息（角色决定可见页面）

// 函数定义：func (a AuthApi) Login(c *gin.Context) {
func (a AuthApi) Login(c *gin.Context) {
// 声明：var req request.Login
	var req request.Login
// 条件判断：if err := c.ShouldBindJSON(&req); err != nil {
	if err := c.ShouldBindJSON(&req); err != nil {
// 执行：fail(c, service.ErrInvalidParam)
		fail(c, service.ErrInvalidParam)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：res, err := a.s.AuthService.Login(req.Username, req.Password)
	res, err := a.s.AuthService.Login(req.Username, req.Password)
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, res)
	ok(c, res)
// 结束当前代码块
}

// Me 当前登录用户

// 函数定义：func (a AuthApi) Me(c *gin.Context) {
func (a AuthApi) Me(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：user, err := a.s.AuthService.User(claims.UserID)
	user, err := a.s.AuthService.User(claims.UserID)
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, service.ErrLoginFailed)
		fail(c, service.ErrLoginFailed)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, user)
	ok(c, user)
// 结束当前代码块
}

// Users 账号列表（页面登录快捷入口用，不返回密码）

// 函数定义：func (a AuthApi) Users(c *gin.Context) {
func (a AuthApi) Users(c *gin.Context) {
// 赋值/更新：list, err := a.s.AuthService.Users(c.Query("role"))
	list, err := a.s.AuthService.Users(c.Query("role"))
// 条件判断：if err != nil {
	if err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, list)
	ok(c, list)
// 结束当前代码块
}
