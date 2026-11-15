// Package web 内置买家端页面（随二进制发布，无需前端构建）：
//
// buyer.html   买家端页（原型图 live-prototype 对齐：商品橱窗 / 订单列 / 详情抽屉 / 验收场景）
// payment.html 模拟支付渠道收银台页（正常 / 重复 / 迟到 / 网络中断通知演示）
// assets/      原型样式、商品图与脚本（domain.js 为真实后端 API 适配层）
//
// 页面由独立 Gin 服务挂在 /shop，静态资源 /shop/assets/**。
// 运营/仓库/客服后台页面不在此处：统一由 GVA 管理端 Vue（web/src/view/shop/*.vue）
// 提供，接口走 GVA /shop/**（登录 + Casbin），避免两套后台页面重复维护。
package web

// 导入依赖包
import (
// 执行："embed"
	"embed"
// 执行："strings"
	"strings"
// 结束当前代码块
)

//go:embed buyer.html payment.html assets
var files embed.FS

// Buyer 返回买家端页面 HTML。

// 函数定义：func Buyer() ([]byte, error) {
func Buyer() ([]byte, error) {
// 返回：return files.ReadFile("buyer.html")
	return files.ReadFile("buyer.html")
// 结束当前代码块
}

// Payment 返回模拟收银台页面 HTML。

// 函数定义：func Payment() ([]byte, error) {
func Payment() ([]byte, error) {
// 返回：return files.ReadFile("payment.html")
	return files.ReadFile("payment.html")
// 结束当前代码块
}

// Asset 返回 /shop/assets/** 静态资源内容；name 为相对 assets 的文件名。
// 仅允许白名单扩展名，避免目录穿越与任意文件读取。

// 函数定义：func Asset(name string) ([]byte, string, bool) {
func Asset(name string) ([]byte, string, bool) {
// 条件判断：if strings.Contains(name, "..") || strings.Contains(name, "\\") || strings.HasPrefix(name, "/") {
	if strings.Contains(name, "..") || strings.Contains(name, "\\") || strings.HasPrefix(name, "/") {
// 返回：return nil, "", false
		return nil, "", false
// 结束当前代码块
	}
// 赋值/更新：data, err := files.ReadFile("assets/" + name)
	data, err := files.ReadFile("assets/" + name)
// 条件判断：if err != nil {
	if err != nil {
// 返回：return nil, "", false
		return nil, "", false
// 结束当前代码块
	}
// 赋值/更新：lower := strings.ToLower(name)
	lower := strings.ToLower(name)
// 声明：var contentType string
	var contentType string
// 多路分支：switch {
	switch {
// 匹配分支：case strings.HasSuffix(lower, ".css"):
	case strings.HasSuffix(lower, ".css"):
// 赋值/更新：contentType = "text/css; charset=utf-8"
		contentType = "text/css; charset=utf-8"
// 匹配分支：case strings.HasSuffix(lower, ".js"):
	case strings.HasSuffix(lower, ".js"):
// 赋值/更新：contentType = "application/javascript; charset=utf-8"
		contentType = "application/javascript; charset=utf-8"
// 匹配分支：case strings.HasSuffix(lower, ".png"):
	case strings.HasSuffix(lower, ".png"):
// 赋值/更新：contentType = "image/png"
		contentType = "image/png"
// 匹配分支：case strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"):
	case strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"):
// 赋值/更新：contentType = "image/jpeg"
		contentType = "image/jpeg"
// 匹配分支：case strings.HasSuffix(lower, ".svg"):
	case strings.HasSuffix(lower, ".svg"):
// 赋值/更新：contentType = "image/svg+xml"
		contentType = "image/svg+xml"
// 匹配分支：case strings.HasSuffix(lower, ".json"):
	case strings.HasSuffix(lower, ".json"):
// 赋值/更新：contentType = "application/json; charset=utf-8"
		contentType = "application/json; charset=utf-8"
// 默认分支
	default:
// 赋值/更新：contentType = "application/octet-stream"
		contentType = "application/octet-stream"
// 结束当前代码块
	}
// 返回：return data, contentType, true
	return data, contentType, true
// 结束当前代码块
}
