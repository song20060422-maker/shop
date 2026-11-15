// 包声明
package config

// 导入依赖包
import (
// 执行："os"
	"os"
// 执行："path/filepath"
	"path/filepath"
// 执行："testing"
	"testing"
// 执行："time"
	"time"

// 执行："github.com/spf13/viper"
	"github.com/spf13/viper"
// 结束当前代码块
)

// 函数定义：func TestDefaultAndTimeout(t *testing.T) {
func TestDefaultAndTimeout(t *testing.T) {
// 赋值/更新：cfg := Default()
	cfg := Default()
// 条件判断：if cfg.Mysql.DSN == "" || cfg.Jwt.Secret == "" || cfg.Shop.PayNotifySecret == "" {
	if cfg.Mysql.DSN == "" || cfg.Jwt.Secret == "" || cfg.Shop.PayNotifySecret == "" {
// 执行：t.Fatalf("默认配置缺少关键项: %+v", cfg)
		t.Fatalf("默认配置缺少关键项: %+v", cfg)
// 结束当前代码块
	}
// 条件判断：if cfg.OrderTimeout() != 10*time.Minute {
	if cfg.OrderTimeout() != 10*time.Minute {
// 执行：t.Fatalf("默认超时应为 10 分钟: %s", cfg.OrderTimeout())
		t.Fatalf("默认超时应为 10 分钟: %s", cfg.OrderTimeout())
// 结束当前代码块
	}
// 赋值/更新：cfg.Shop.OrderTimeoutMinutes = 3
	cfg.Shop.OrderTimeoutMinutes = 3
// 条件判断：if cfg.OrderTimeout() != 3*time.Minute {
	if cfg.OrderTimeout() != 3*time.Minute {
// 执行：t.Fatalf("自定义超时错误: %s", cfg.OrderTimeout())
		t.Fatalf("自定义超时错误: %s", cfg.OrderTimeout())
// 结束当前代码块
	}
// 赋值/更新：cfg.Channel.RefundMode = "fail"
	cfg.Channel.RefundMode = "fail"
// 条件判断：if cfg.AllowRefundByChannel() {
	if cfg.AllowRefundByChannel() {
// 执行：t.Fatal("渠道 fail 模式不应受理退款")
		t.Fatal("渠道 fail 模式不应受理退款")
// 结束当前代码块
	}
// 赋值/更新：cfg.Channel.RefundMode = "success"
	cfg.Channel.RefundMode = "success"
// 条件判断：if !cfg.AllowRefundByChannel() {
	if !cfg.AllowRefundByChannel() {
// 执行：t.Fatal("渠道 success 模式应受理退款")
		t.Fatal("渠道 success 模式应受理退款")
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestLoadFromFileAndEnv(t *testing.T) {
func TestLoadFromFileAndEnv(t *testing.T) {
// 赋值/更新：dir := t.TempDir()
	dir := t.TempDir()
// 赋值/更新：path := filepath.Join(dir, "shop.yaml")
	path := filepath.Join(dir, "shop.yaml")
// 赋值/更新：content := `
	content := `
system:
  addr: ":9099"
  mode: release
mysql:
  dsn: "root:root@tcp(127.0.0.1:3307)/shop_yaml?charset=utf8mb4"
jwt:
  secret: "yaml-secret"
  expire-hours: 6
shop:
  order-timeout-minutes: 2
  pay-notify-secret: "yaml-pay"
  demo-password: "yaml-pass"
  demo-mode: true
channel:
  refund-mode: fail
`
// 条件判断：if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 赋值/更新：cfg, err := Load(path)
	cfg, err := Load(path)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("读取配置失败: %v", err)
		t.Fatalf("读取配置失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if cfg.System.Addr != ":9099" || cfg.Jwt.Secret != "yaml-secret" || !cfg.Shop.DemoMode {
	if cfg.System.Addr != ":9099" || cfg.Jwt.Secret != "yaml-secret" || !cfg.Shop.DemoMode {
// 执行：t.Fatalf("配置文件未生效: %+v", cfg)
		t.Fatalf("配置文件未生效: %+v", cfg)
// 结束当前代码块
	}
// 条件判断：if cfg.AllowRefundByChannel() {
	if cfg.AllowRefundByChannel() {
// 执行：t.Fatal("配置渠道失败模式未生效")
		t.Fatal("配置渠道失败模式未生效")
// 结束当前代码块
	}

	// 环境变量覆盖
	t.Setenv("SHOP_DSN", "root:root@tcp(127.0.0.1:3307)/shop_env")
// 执行：t.Setenv("SHOP_ADDR", ":9191")
	t.Setenv("SHOP_ADDR", ":9191")
// 执行：t.Setenv("SHOP_MODE", "test")
	t.Setenv("SHOP_MODE", "test")
// 执行：t.Setenv("SHOP_JWT_SECRET", "env-secret")
	t.Setenv("SHOP_JWT_SECRET", "env-secret")
// 执行：t.Setenv("SHOP_PAY_SECRET", "env-pay")
	t.Setenv("SHOP_PAY_SECRET", "env-pay")
// 执行：t.Setenv("SHOP_ORDER_TIMEOUT_MINUTES", "5")
	t.Setenv("SHOP_ORDER_TIMEOUT_MINUTES", "5")
// 赋值/更新：cfg, err = Load(path)
	cfg, err = Load(path)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if cfg.System.Addr != ":9191" || cfg.System.Mode != "test" ||
	if cfg.System.Addr != ":9191" || cfg.System.Mode != "test" ||
// 赋值/更新：cfg.Jwt.Secret != "env-secret" || cfg.Shop.PayNotifySecret != "env-pay" ||
		cfg.Jwt.Secret != "env-secret" || cfg.Shop.PayNotifySecret != "env-pay" ||
// 赋值/更新：cfg.Shop.OrderTimeoutMinutes != 5 || cfg.Mysql.DSN != "root:root@tcp(127.0.0.1:3307)/shop_env" {
		cfg.Shop.OrderTimeoutMinutes != 5 || cfg.Mysql.DSN != "root:root@tcp(127.0.0.1:3307)/shop_env" {
// 执行：t.Fatalf("环境变量覆盖未生效: %+v", cfg)
		t.Fatalf("环境变量覆盖未生效: %+v", cfg)
// 结束当前代码块
	}

	// 非法 yaml
	badPath := filepath.Join(dir, "bad.yaml")
// 条件判断：if err := os.WriteFile(badPath, []byte("system: ["), 0o600); err != nil {
	if err := os.WriteFile(badPath, []byte("system: ["), 0o600); err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if _, err := Load(badPath); err == nil {
	if _, err := Load(badPath); err == nil {
// 执行：t.Fatal("非法 yaml 应报错")
		t.Fatal("非法 yaml 应报错")
// 结束当前代码块
	}

	// 文件不存在回落到默认值（先清空上面设置的环境变量）
	t.Setenv("SHOP_ADDR", "")
// 执行：t.Setenv("SHOP_DSN", "")
	t.Setenv("SHOP_DSN", "")
// 赋值/更新：cfg, err = Load(filepath.Join(dir, "missing.yaml"))
	cfg, err = Load(filepath.Join(dir, "missing.yaml"))
// 条件判断：if err != nil || cfg.System.Addr != ":8090" {
	if err != nil || cfg.System.Addr != ":8090" {
// 执行：t.Fatalf("缺省配置错误: %v %+v", err, cfg)
		t.Fatalf("缺省配置错误: %v %+v", err, cfg)
// 结束当前代码块
	}
	// 空路径直接使用默认值
	if _, err := Load(""); err != nil {
// 执行：t.Fatalf("空路径不应报错: %v", err)
		t.Fatalf("空路径不应报错: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestFromViper(t *testing.T) {
func TestFromViper(t *testing.T) {
// 条件判断：if cfg := FromViper(nil); cfg.Jwt.Secret != Default().Jwt.Secret {
	if cfg := FromViper(nil); cfg.Jwt.Secret != Default().Jwt.Secret {
// 执行：t.Fatal("nil viper 应返回默认配置")
		t.Fatal("nil viper 应返回默认配置")
// 结束当前代码块
	}
// 赋值/更新：vp := viper.New()
	vp := viper.New()
// 执行：vp.Set("shop.mysql-dsn", "root:root@tcp(127.0.0.1:3307)/shop_vp")
	vp.Set("shop.mysql-dsn", "root:root@tcp(127.0.0.1:3307)/shop_vp")
// 执行：vp.Set("shop.jwt-secret", "vp-secret")
	vp.Set("shop.jwt-secret", "vp-secret")
// 执行：vp.Set("shop.pay-notify-secret", "vp-pay")
	vp.Set("shop.pay-notify-secret", "vp-pay")
// 执行：vp.Set("shop.demo-password", "vp-pass")
	vp.Set("shop.demo-password", "vp-pass")
// 执行：vp.Set("shop.refund-mode", "fail")
	vp.Set("shop.refund-mode", "fail")
// 执行：vp.Set("shop.order-timeout-minutes", 7)
	vp.Set("shop.order-timeout-minutes", 7)
// 执行：vp.Set("shop.jwt-expire-hours", 8)
	vp.Set("shop.jwt-expire-hours", 8)
// 执行：vp.Set("shop.mode", "release")
	vp.Set("shop.mode", "release")
// 执行：vp.Set("shop.demo-mode", true)
	vp.Set("shop.demo-mode", true)
// 执行：vp.Set("shop.max-open-conns", 30)
	vp.Set("shop.max-open-conns", 30)
// 赋值/更新：cfg := FromViper(vp)
	cfg := FromViper(vp)
// 条件判断：if cfg.Mysql.DSN != "root:root@tcp(127.0.0.1:3307)/shop_vp" ||
	if cfg.Mysql.DSN != "root:root@tcp(127.0.0.1:3307)/shop_vp" ||
// 赋值/更新：cfg.Jwt.Secret != "vp-secret" || cfg.Shop.PayNotifySecret != "vp-pay" ||
		cfg.Jwt.Secret != "vp-secret" || cfg.Shop.PayNotifySecret != "vp-pay" ||
// 赋值/更新：cfg.Shop.DemoPassword != "vp-pass" || cfg.Channel.RefundMode != "fail" ||
		cfg.Shop.DemoPassword != "vp-pass" || cfg.Channel.RefundMode != "fail" ||
// 赋值/更新：cfg.Shop.OrderTimeoutMinutes != 7 || cfg.Jwt.ExpireHours != 8 ||
		cfg.Shop.OrderTimeoutMinutes != 7 || cfg.Jwt.ExpireHours != 8 ||
// 赋值/更新：cfg.System.Mode != "release" || !cfg.Shop.DemoMode || cfg.Mysql.MaxOpenConns != 30 {
		cfg.System.Mode != "release" || !cfg.Shop.DemoMode || cfg.Mysql.MaxOpenConns != 30 {
// 执行：t.Fatalf("viper 配置映射错误: %+v", cfg)
		t.Fatalf("viper 配置映射错误: %+v", cfg)
// 结束当前代码块
	}
// 结束当前代码块
}
