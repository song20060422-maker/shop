// 包声明
package service

// 导入依赖包
import (
// 执行："strings"
	"strings"
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

// AuthService 账号与令牌

// 类型定义：type AuthService struct {
type AuthService struct {
// 执行：deps Deps
	deps Deps
// 结束当前代码块
}

// Login 校验账号密码并签发 JWT

// 函数定义：func (s AuthService) Login(username, password string) (response.LoginResult, error) {
func (s AuthService) Login(username, password string) (response.LoginResult, error) {
// 声明：var out response.LoginResult
	var out response.LoginResult
// 赋值/更新：username = strings.TrimSpace(username)
	username = strings.TrimSpace(username)
// 条件判断：if username == "" || password == "" {
	if username == "" || password == "" {
// 返回：return out, ErrLoginFailed
		return out, ErrLoginFailed
// 结束当前代码块
	}
// 声明：var user model.User
	var user model.User
// 条件判断：if err := s.deps.DB.Where("username = ?", username).First(&user).Error; err != nil {
	if err := s.deps.DB.Where("username = ?", username).First(&user).Error; err != nil {
// 返回：return out, ErrLoginFailed
		return out, ErrLoginFailed
// 结束当前代码块
	}
// 条件判断：if user.PasswordHash == "" || !utils.VerifyPassword(password, user.PasswordHash) {
	if user.PasswordHash == "" || !utils.VerifyPassword(password, user.PasswordHash) {
// 返回：return out, ErrLoginFailed
		return out, ErrLoginFailed
// 结束当前代码块
	}
// 赋值/更新：now := s.deps.Now()
	now := s.deps.Now()
// 赋值/更新：expire := now.Add(time.Duration(s.deps.Cfg.Jwt.ExpireHours) * time.Hour)
	expire := now.Add(time.Duration(s.deps.Cfg.Jwt.ExpireHours) * time.Hour)
// 条件判断：if s.deps.Cfg.Jwt.ExpireHours <= 0 {
	if s.deps.Cfg.Jwt.ExpireHours <= 0 {
// 赋值/更新：expire = now.Add(12 * time.Hour)
		expire = now.Add(12 * time.Hour)
// 结束当前代码块
	}
// 赋值/更新：token, err := utils.SignJWT(s.deps.Cfg.Jwt.Secret, utils.Claims{
	token, err := utils.SignJWT(s.deps.Cfg.Jwt.Secret, utils.Claims{
// 执行：UserID:   user.ID,
		UserID:   user.ID,
// 执行：Username: user.Username,
		Username: user.Username,
// 执行：Role:     user.Role,
		Role:     user.Role,
// 执行：IssuedAt: now.Unix(),
		IssuedAt: now.Unix(),
// 执行：ExpireAt: expire.Unix(),
		ExpireAt: expire.Unix(),
// 结束当前代码块
	})
// 条件判断：if err != nil {
	if err != nil {
// 返回：return out, err
		return out, err
// 结束当前代码块
	}
// 返回：return response.LoginResult{Token: token, User: user, ExpireAt: expire.Format("2006-01-02 15:04:05")}, nil
	return response.LoginResult{Token: token, User: user, ExpireAt: expire.Format("2006-01-02 15:04:05")}, nil
// 结束当前代码块
}

// User 按 ID 查账号

// 函数定义：func (s AuthService) User(id uint64) (model.User, error) {
func (s AuthService) User(id uint64) (model.User, error) {
// 声明：var user model.User
	var user model.User
// 条件判断：if err := s.deps.DB.Where("id = ?", id).First(&user).Error; err != nil {
	if err := s.deps.DB.Where("id = ?", id).First(&user).Error; err != nil {
// 返回：return user, ErrLoginFailed
		return user, ErrLoginFailed
// 结束当前代码块
	}
// 返回：return user, nil
	return user, nil
// 结束当前代码块
}

// Users 按角色列出账号（运营核查页可查看角色分布）

// 函数定义：func (s AuthService) Users(role string) ([]model.User, error) {
func (s AuthService) Users(role string) ([]model.User, error) {
// 声明：var users []model.User
	var users []model.User
// 赋值/更新：q := s.deps.DB
	q := s.deps.DB
// 条件判断：if role != "" {
	if role != "" {
// 赋值/更新：q = q.Where("role = ?", role)
		q = q.Where("role = ?", role)
// 结束当前代码块
	}
// 赋值/更新：err := q.Order("id asc").Find(&users).Error
	err := q.Order("id asc").Find(&users).Error
// 返回：return users, err
	return users, err
// 结束当前代码块
}
