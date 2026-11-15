// 包声明
package config

// 导入依赖包
import (
// 执行："strconv"
	"strconv"
// 执行："strings"
	"strings"

// 执行："github.com/spf13/viper"
	"github.com/spf13/viper"
// 结束当前代码块
)

// FromViper 从 GVA 全局 viper 读取 shop.* 配置段，未配置项回落到默认值。
// 这样 Shop 业务与 GVA 自身的 mysql/jwt/system 配置解耦，互不覆盖。

// 函数定义：func FromViper(vp *viper.Viper) Config {
func FromViper(vp *viper.Viper) Config {
// 赋值/更新：cfg := Default()
	cfg := Default()
// 条件判断：if vp == nil {
	if vp == nil {
// 返回：return cfg
		return cfg
// 结束当前代码块
	}
// 条件判断：if v := vp.GetString("shop.mysql-dsn"); strings.TrimSpace(v) != "" {
	if v := vp.GetString("shop.mysql-dsn"); strings.TrimSpace(v) != "" {
// 赋值/更新：cfg.Mysql.DSN = v
		cfg.Mysql.DSN = v
// 结束当前代码块
	}
// 条件判断：if v := vp.GetString("shop.jwt-secret"); strings.TrimSpace(v) != "" {
	if v := vp.GetString("shop.jwt-secret"); strings.TrimSpace(v) != "" {
// 赋值/更新：cfg.Jwt.Secret = v
		cfg.Jwt.Secret = v
// 结束当前代码块
	}
// 条件判断：if v := vp.GetString("shop.pay-notify-secret"); strings.TrimSpace(v) != "" {
	if v := vp.GetString("shop.pay-notify-secret"); strings.TrimSpace(v) != "" {
// 赋值/更新：cfg.Shop.PayNotifySecret = v
		cfg.Shop.PayNotifySecret = v
// 结束当前代码块
	}
// 条件判断：if v := vp.GetString("shop.demo-password"); strings.TrimSpace(v) != "" {
	if v := vp.GetString("shop.demo-password"); strings.TrimSpace(v) != "" {
// 赋值/更新：cfg.Shop.DemoPassword = v
		cfg.Shop.DemoPassword = v
// 结束当前代码块
	}
// 条件判断：if v := vp.GetString("shop.refund-mode"); strings.TrimSpace(v) != "" {
	if v := vp.GetString("shop.refund-mode"); strings.TrimSpace(v) != "" {
// 赋值/更新：cfg.Channel.RefundMode = v
		cfg.Channel.RefundMode = v
// 结束当前代码块
	}
// 条件判断：if v := vp.GetString("shop.pay-channel"); strings.TrimSpace(v) != "" {
	if v := vp.GetString("shop.pay-channel"); strings.TrimSpace(v) != "" {
// 赋值/更新：cfg.Shop.PayChannel = v
		cfg.Shop.PayChannel = v
// 结束当前代码块
	}
	// 支付宝渠道（沙箱演示扩展）：管理端只用到退款能力，密钥等配置与买家端共用同一 yaml 段。
	cfg.Shop.Alipay.AppID = vp.GetString("shop.alipay.app-id")
// 赋值/更新：cfg.Shop.Alipay.PrivateKey = vp.GetString("shop.alipay.private-key")
	cfg.Shop.Alipay.PrivateKey = vp.GetString("shop.alipay.private-key")
// 赋值/更新：cfg.Shop.Alipay.AlipayPublicKey = vp.GetString("shop.alipay.alipay-public-key")
	cfg.Shop.Alipay.AlipayPublicKey = vp.GetString("shop.alipay.alipay-public-key")
// 赋值/更新：cfg.Shop.Alipay.NotifyURL = vp.GetString("shop.alipay.notify-url")
	cfg.Shop.Alipay.NotifyURL = vp.GetString("shop.alipay.notify-url")
// 赋值/更新：cfg.Shop.Alipay.ReturnURL = vp.GetString("shop.alipay.return-url")
	cfg.Shop.Alipay.ReturnURL = vp.GetString("shop.alipay.return-url")
// 条件判断：if n := vp.GetInt("shop.order-timeout-minutes"); n > 0 {
	if n := vp.GetInt("shop.order-timeout-minutes"); n > 0 {
// 赋值/更新：cfg.Shop.OrderTimeoutMinutes = n
		cfg.Shop.OrderTimeoutMinutes = n
// 结束当前代码块
	}
// 条件判断：if n := vp.GetInt("shop.jwt-expire-hours"); n > 0 {
	if n := vp.GetInt("shop.jwt-expire-hours"); n > 0 {
// 赋值/更新：cfg.Jwt.ExpireHours = n
		cfg.Jwt.ExpireHours = n
// 结束当前代码块
	}
// 条件判断：if v := vp.GetString("shop.mode"); strings.TrimSpace(v) != "" {
	if v := vp.GetString("shop.mode"); strings.TrimSpace(v) != "" {
// 赋值/更新：cfg.System.Mode = v
		cfg.System.Mode = v
// 结束当前代码块
	}
// 赋值/更新：cfg.Shop.DemoMode = vp.GetBool("shop.demo-mode")
	cfg.Shop.DemoMode = vp.GetBool("shop.demo-mode")
// 条件判断：if v := strconv.Itoa(vp.GetInt("shop.max-open-conns")); v != "0" {
	if v := strconv.Itoa(vp.GetInt("shop.max-open-conns")); v != "0" {
// 赋值/更新：cfg.Mysql.MaxOpenConns = vp.GetInt("shop.max-open-conns")
		cfg.Mysql.MaxOpenConns = vp.GetInt("shop.max-open-conns")
// 结束当前代码块
	}
// 返回：return cfg
	return cfg
// 结束当前代码块
}
