// Package config 定义 GVA 风格的分组配置，并支持 yaml + 环境变量两级覆盖。
package config

// 导入依赖包
import (
// 执行："os"
	"os"
// 执行："strconv"
	"strconv"
// 执行："time"
	"time"

// 执行："gopkg.in/yaml.v3"
	"gopkg.in/yaml.v3"
// 结束当前代码块
)

// System 服务基础配置

// 类型定义：type System struct {
type System struct {
// 执行：Addr string `yaml:"addr"`
	Addr string `yaml:"addr"`
// 执行：Mode string `yaml:"mode"` // debug | release | test
	Mode string `yaml:"mode"` // debug | release | test
	// AdminTarget GVA 管理端基地址：内置后台独立页面（/ops /warehouse /cs）通过
	// 本服务的 /gva/* 反向代理访问管理端接口（base/login 与 /shop/**），避免跨域。
	AdminTarget string `yaml:"admin-target"`
// 结束当前代码块
}

// Mysql 数据库配置

// 类型定义：type Mysql struct {
type Mysql struct {
// 执行：DSN          string `yaml:"dsn"`
	DSN          string `yaml:"dsn"`
// 执行：MaxIdleConns int    `yaml:"max-idle-conns"`
	MaxIdleConns int    `yaml:"max-idle-conns"`
// 执行：MaxOpenConns int    `yaml:"max-open-conns"`
	MaxOpenConns int    `yaml:"max-open-conns"`
// 执行：LogMode      bool   `yaml:"log-mode"`
	LogMode      bool   `yaml:"log-mode"`
// 结束当前代码块
}

// Jwt 令牌配置

// 类型定义：type Jwt struct {
type Jwt struct {
// 执行：Secret      string `yaml:"secret"`
	Secret      string `yaml:"secret"`
// 执行：ExpireHours int    `yaml:"expire-hours"`
	ExpireHours int    `yaml:"expire-hours"`
// 结束当前代码块
}

// Shop 业务参数
type Shop struct {
// 执行：OrderTimeoutMinutes int    `yaml:"order-timeout-minutes"`
	OrderTimeoutMinutes int    `yaml:"order-timeout-minutes"`
	// CloseExpiredIntervalSeconds 超时关单后台调度间隔（秒），<=0 时回落默认值；
	// 生产环境靠该定时任务关闭到期待付款订单并释放占用（PROD-10）。
	CloseExpiredIntervalSeconds int    `yaml:"close-expired-interval-seconds"`
// 执行：PayNotifySecret             string `yaml:"pay-notify-secret"`
	PayNotifySecret             string `yaml:"pay-notify-secret"`
// 执行：DemoPassword        string `yaml:"demo-password"`
	DemoPassword        string `yaml:"demo-password"`
// 执行：DemoMode            bool   `yaml:"demo-mode"` // 演示数据重置开关，生产必须为 false
	DemoMode            bool   `yaml:"demo-mode"` // 演示数据重置开关，生产必须为 false
	// PayChannel 支付渠道选择：simulated（默认，模拟渠道，验收口径不变）| alipay（支付宝沙箱真实渠道，
	// 演示扩展，见 docs/assumptions.md）。
	PayChannel string `yaml:"pay-channel"`
	// Alipay 支付宝渠道参数（沙箱凭据，无真实资金；密钥入配置已提请学生裁决）
	Alipay Alipay `yaml:"alipay"`
// 结束当前代码块
}

// Alipay 支付宝渠道配置：私钥/支付宝公钥为 PKCS8 PEM 内容（不含头尾行）。
// 网关地址由 SDK 按 production 开关自动选择（沙箱 openapi-sandbox.dl.alipaydev.com）。

// 类型定义：type Alipay struct {
type Alipay struct {
// 执行：AppID           string `yaml:"app-id"`
	AppID           string `yaml:"app-id"`
// 执行：PrivateKey      string `yaml:"private-key"`
	PrivateKey      string `yaml:"private-key"`
// 执行：AlipayPublicKey string `yaml:"alipay-public-key"`
	AlipayPublicKey string `yaml:"alipay-public-key"`
// 执行：NotifyURL       string `yaml:"notify-url"` // 公网可达的异步通知地址
	NotifyURL       string `yaml:"notify-url"` // 公网可达的异步通知地址
// 执行：ReturnURL       string `yaml:"return-url"` // 支付完成后跳回的买家页地址
	ReturnURL       string `yaml:"return-url"` // 支付完成后跳回的买家页地址
// 结束当前代码块
}

// Channel 模拟支付渠道配置

// 类型定义：type Channel struct {
type Channel struct {
// 执行：RefundMode string `yaml:"refund-mode"` // success | fail
	RefundMode string `yaml:"refund-mode"` // success | fail
// 结束当前代码块
}

// Config 聚合配置

// 类型定义：type Config struct {
type Config struct {
// 执行：System  System  `yaml:"system"`
	System  System  `yaml:"system"`
// 执行：Mysql   Mysql   `yaml:"mysql"`
	Mysql   Mysql   `yaml:"mysql"`
// 执行：Jwt     Jwt     `yaml:"jwt"`
	Jwt     Jwt     `yaml:"jwt"`
// 执行：Shop    Shop    `yaml:"shop"`
	Shop    Shop    `yaml:"shop"`
// 执行：Channel Channel `yaml:"channel"`
	Channel Channel `yaml:"channel"`
// 结束当前代码块
}

// Default 返回可立即运行的默认配置（本地 shop-mysql:3307）。

// 函数定义：func Default() Config {
func Default() Config {
// 返回：return Config{
	return Config{
// 执行：System: System{Addr: ":8090", Mode: "debug", AdminTarget: "http://127.0.0.1:8099"},
		System: System{Addr: ":8090", Mode: "debug", AdminTarget: "http://127.0.0.1:8099"},
// 执行：Mysql: Mysql{
		Mysql: Mysql{
// 执行：DSN:          "root:root@tcp(127.0.0.1:3307)/shop?charset=utf8mb4&parseTime=True&loc=Local",
			DSN:          "root:root@tcp(127.0.0.1:3307)/shop?charset=utf8mb4&parseTime=True&loc=Local",
// 执行：MaxIdleConns: 5,
			MaxIdleConns: 5,
// 执行：MaxOpenConns: 20,
			MaxOpenConns: 20,
// 执行：LogMode:      false,
			LogMode:      false,
// 结束当前代码块
		},
// 执行：Jwt:     Jwt{Secret: "shop-gva-dev-secret", ExpireHours: 12},
		Jwt:     Jwt{Secret: "shop-gva-dev-secret", ExpireHours: 12},
// 执行：Shop:    Shop{OrderTimeoutMinutes: 10, PayNotifySecret: "channel-shared-secret", DemoPassword: "shop123", PayC...
		Shop:    Shop{OrderTimeoutMinutes: 10, PayNotifySecret: "channel-shared-secret", DemoPassword: "shop123", PayChannel: "simulated"},
// 执行：Channel: Channel{RefundMode: "success"},
		Channel: Channel{RefundMode: "success"},
// 结束当前代码块
	}
// 结束当前代码块
}

// Load 读取 yaml 配置（文件不存在时使用默认值），随后应用环境变量覆盖。

// 函数定义：func Load(path string) (Config, error) {
func Load(path string) (Config, error) {
// 赋值/更新：cfg := Default()
	cfg := Default()
// 条件判断：if path != "" {
	if path != "" {
// 条件判断：if raw, err := os.ReadFile(path); err == nil {
		if raw, err := os.ReadFile(path); err == nil {
// 条件判断：if err := yaml.Unmarshal(raw, &cfg); err != nil {
			if err := yaml.Unmarshal(raw, &cfg); err != nil {
// 返回：return cfg, err
				return cfg, err
// 结束当前代码块
			}
// 条件分支：} else if !os.IsNotExist(err) {
		} else if !os.IsNotExist(err) {
// 返回：return cfg, err
			return cfg, err
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if v := os.Getenv("SHOP_DSN"); v != "" {
	if v := os.Getenv("SHOP_DSN"); v != "" {
// 赋值/更新：cfg.Mysql.DSN = v
		cfg.Mysql.DSN = v
// 结束当前代码块
	}
// 条件判断：if v := os.Getenv("SHOP_ADDR"); v != "" {
	if v := os.Getenv("SHOP_ADDR"); v != "" {
// 赋值/更新：cfg.System.Addr = v
		cfg.System.Addr = v
// 结束当前代码块
	}
// 条件判断：if v := os.Getenv("SHOP_MODE"); v != "" {
	if v := os.Getenv("SHOP_MODE"); v != "" {
// 赋值/更新：cfg.System.Mode = v
		cfg.System.Mode = v
// 结束当前代码块
	}
// 条件判断：if v := os.Getenv("SHOP_ADMIN_TARGET"); v != "" {
	if v := os.Getenv("SHOP_ADMIN_TARGET"); v != "" {
// 赋值/更新：cfg.System.AdminTarget = v
		cfg.System.AdminTarget = v
// 结束当前代码块
	}
// 条件判断：if v := os.Getenv("SHOP_JWT_SECRET"); v != "" {
	if v := os.Getenv("SHOP_JWT_SECRET"); v != "" {
// 赋值/更新：cfg.Jwt.Secret = v
		cfg.Jwt.Secret = v
// 结束当前代码块
	}
// 条件判断：if v := os.Getenv("SHOP_PAY_SECRET"); v != "" {
	if v := os.Getenv("SHOP_PAY_SECRET"); v != "" {
// 赋值/更新：cfg.Shop.PayNotifySecret = v
		cfg.Shop.PayNotifySecret = v
// 结束当前代码块
	}
// 条件判断：if v := os.Getenv("SHOP_PAY_CHANNEL"); v != "" {
	if v := os.Getenv("SHOP_PAY_CHANNEL"); v != "" {
// 赋值/更新：cfg.Shop.PayChannel = v
		cfg.Shop.PayChannel = v
// 结束当前代码块
	}
// 条件判断：if v := os.Getenv("SHOP_ORDER_TIMEOUT_MINUTES"); v != "" {
	if v := os.Getenv("SHOP_ORDER_TIMEOUT_MINUTES"); v != "" {
// 条件判断：if n, err := strconv.Atoi(v); err == nil && n > 0 {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
// 赋值/更新：cfg.Shop.OrderTimeoutMinutes = n
			cfg.Shop.OrderTimeoutMinutes = n
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if v := os.Getenv("SHOP_CLOSE_EXPIRED_SECONDS"); v != "" {
	if v := os.Getenv("SHOP_CLOSE_EXPIRED_SECONDS"); v != "" {
// 条件判断：if n, err := strconv.Atoi(v); err == nil && n > 0 {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
// 赋值/更新：cfg.Shop.CloseExpiredIntervalSeconds = n
			cfg.Shop.CloseExpiredIntervalSeconds = n
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return cfg, nil
	return cfg, nil
// 结束当前代码块
}

// OrderTimeout 待付款订单保留时长，缺省 10 分钟。

// 函数定义：func (c Config) OrderTimeout() time.Duration {
func (c Config) OrderTimeout() time.Duration {
// 条件判断：if c.Shop.OrderTimeoutMinutes <= 0 {
	if c.Shop.OrderTimeoutMinutes <= 0 {
// 返回：return 10 * time.Minute
		return 10 * time.Minute
// 结束当前代码块
	}
// 返回：return time.Duration(c.Shop.OrderTimeoutMinutes) * time.Minute
	return time.Duration(c.Shop.OrderTimeoutMinutes) * time.Minute
// 结束当前代码块
}

// CloseExpiredInterval 超时关单调度间隔，缺省 30 秒。

// 函数定义：func (c Config) CloseExpiredInterval() time.Duration {
func (c Config) CloseExpiredInterval() time.Duration {
// 条件判断：if c.Shop.CloseExpiredIntervalSeconds <= 0 {
	if c.Shop.CloseExpiredIntervalSeconds <= 0 {
// 返回：return 30 * time.Second
		return 30 * time.Second
// 结束当前代码块
	}
// 返回：return time.Duration(c.Shop.CloseExpiredIntervalSeconds) * time.Second
	return time.Duration(c.Shop.CloseExpiredIntervalSeconds) * time.Second
// 结束当前代码块
}

// AllowRefundByChannel 渠道当前是否受理退款成功。

// 函数定义：func (c Config) AllowRefundByChannel() bool { return c.Channel.RefundMode != "fail" }
func (c Config) AllowRefundByChannel() bool { return c.Channel.RefundMode != "fail" }
