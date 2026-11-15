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

import (
	"embed"
	"strings"
)

//go:embed buyer.html payment.html assets
var files embed.FS

// Buyer 返回买家端页面 HTML。

func Buyer() ([]byte, error) {
	return files.ReadFile("buyer.html")
}

// Payment 返回模拟收银台页面 HTML。

func Payment() ([]byte, error) {
	return files.ReadFile("payment.html")
}

// Asset 返回 /shop/assets/** 静态资源内容；name 为相对 assets 的文件名。
// 仅允许白名单扩展名，避免目录穿越与任意文件读取。

func Asset(name string) ([]byte, string, bool) {
	if strings.Contains(name, "..") || strings.Contains(name, "\\") || strings.HasPrefix(name, "/") {
		return nil, "", false
	}
	data, err := files.ReadFile("assets/" + name)
	if err != nil {
		return nil, "", false
	}
	lower := strings.ToLower(name)
	var contentType string
	switch {
	case strings.HasSuffix(lower, ".css"):
		contentType = "text/css; charset=utf-8"
	case strings.HasSuffix(lower, ".js"):
		contentType = "application/javascript; charset=utf-8"
	case strings.HasSuffix(lower, ".png"):
		contentType = "image/png"
	case strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"):
		contentType = "image/jpeg"
	case strings.HasSuffix(lower, ".svg"):
		contentType = "image/svg+xml"
	case strings.HasSuffix(lower, ".json"):
		contentType = "application/json; charset=utf-8"
	default:
		contentType = "application/octet-stream"
	}
	return data, contentType, true
}
