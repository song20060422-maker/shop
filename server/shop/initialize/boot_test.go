// 包声明
package initialize

// 导入依赖包
import (
// 执行："os"
	"os"
// 执行："testing"
	"testing"

// 执行："github.com/spf13/viper"
	"github.com/spf13/viper"
// 执行："go.uber.org/zap"
	"go.uber.org/zap"

// 执行："shop/global"
	"shop/global"
// 执行：shopconfig "shop/shop/config"
	shopconfig "shop/shop/config"
// 执行：shopdb "shop/shop/db"
	shopdb "shop/shop/db"
// 结束当前代码块
)

// 共享测试库（shop_test）中可能残留其他包或历史运行的数据，
// Boot 依赖迁移+种子+密码补齐，测试前先按 FK 安全顺序清空业务表。
func resetTables(t *testing.T, dsn string) {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：cfg := shopconfig.Default()
	cfg := shopconfig.Default()
// 赋值/更新：cfg.Mysql.DSN = dsn
	cfg.Mysql.DSN = dsn
// 赋值/更新：db, err := shopdb.Gorm(cfg)
	db, err := shopdb.Gorm(cfg)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("复位连接失败: %v", err)
		t.Fatalf("复位连接失败: %v", err)
// 结束当前代码块
	}
// 循环处理：for _, tbl := range []string{
	for _, tbl := range []string{
// 执行："sessions", "abnormal_payments", "payments", "payment_notices", "refunds",
		"sessions", "abnormal_payments", "payments", "payment_notices", "refunds",
// 执行："shipments", "inventory_ledger", "operation_logs", "orders", "products", "users",
		"shipments", "inventory_ledger", "operation_logs", "orders", "products", "users",
// 结束当前代码块
	} {
// 条件判断：if err := db.Exec("DELETE FROM " + tbl).Error; err != nil {
		if err := db.Exec("DELETE FROM " + tbl).Error; err != nil {
// 执行：t.Fatalf("复位 %s 失败: %v", tbl, err)
			t.Fatalf("复位 %s 失败: %v", tbl, err)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if sqlDB, err := db.DB(); err == nil {
	if sqlDB, err := db.DB(); err == nil {
// 赋值/更新：_ = sqlDB.Close()
		_ = sqlDB.Close()
// 结束当前代码块
	}
// 结束当前代码块
}

// Boot 测试：从 viper 读取 shop.* 配置段，连库 -> 迁移 -> 种子 -> 装配 ServiceGroup。
// 数据库可用 SHOP_TEST_DSN 覆盖（默认 shop_test，须 -p 1 串行执行）。
func TestBoot(t *testing.T) {
// 赋值/更新：dsn := os.Getenv("SHOP_TEST_DSN")
	dsn := os.Getenv("SHOP_TEST_DSN")
// 条件判断：if dsn == "" {
	if dsn == "" {
// 赋值/更新：dsn = "root:root@tcp(127.0.0.1:3307)/shop_test?charset=utf8mb4&parseTime=True&loc=Local"
		dsn = "root:root@tcp(127.0.0.1:3307)/shop_test?charset=utf8mb4&parseTime=True&loc=Local"
// 结束当前代码块
	}
// 执行：resetTables(t, dsn)
	resetTables(t, dsn)
// 赋值/更新：vp := viper.New()
	vp := viper.New()
// 执行：vp.Set("shop.mysql-dsn", dsn)
	vp.Set("shop.mysql-dsn", dsn)
// 执行：vp.Set("shop.jwt-secret", "boot-secret")
	vp.Set("shop.jwt-secret", "boot-secret")
// 执行：vp.Set("shop.pay-notify-secret", "boot-pay")
	vp.Set("shop.pay-notify-secret", "boot-pay")
// 执行：vp.Set("shop.demo-password", "shop123")
	vp.Set("shop.demo-password", "shop123")
// 执行：vp.Set("shop.order-timeout-minutes", 10)
	vp.Set("shop.order-timeout-minutes", 10)
// 赋值/更新：global.GVA_VP = vp
	global.GVA_VP = vp
// 赋值/更新：global.GVA_LOG = zap.NewNop()
	global.GVA_LOG = zap.NewNop()

// 赋值/更新：group, cfg, err := Boot()
	group, cfg, err := Boot()
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("Boot 失败: %v", err)
		t.Fatalf("Boot 失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if group == nil || cfg.Jwt.Secret != "boot-secret" {
	if group == nil || cfg.Jwt.Secret != "boot-secret" {
// 执行：t.Fatalf("Boot 装配结果错误: %+v", cfg)
		t.Fatalf("Boot 装配结果错误: %+v", cfg)
// 结束当前代码块
	}
	// Shop 使用独立连接，严禁覆盖框架全局 GVA_DB（gva 主库）
	if global.GVA_DB != nil {
// 执行：t.Fatal("Boot 不应覆盖 global.GVA_DB")
		t.Fatal("Boot 不应覆盖 global.GVA_DB")
// 结束当前代码块
	}
	// 种子账号可登录
	if _, err := group.AuthService.Login("buyer1", "shop123"); err != nil {
// 执行：t.Fatalf("种子账号登录失败: %v", err)
		t.Fatalf("种子账号登录失败: %v", err)
// 结束当前代码块
	}
	// 二次启动幂等（迁移与种子可重复执行）
	if _, _, err := Boot(); err != nil {
// 执行：t.Fatalf("重复 Boot 失败: %v", err)
		t.Fatalf("重复 Boot 失败: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// 连库失败时 Boot 必须返回错误，不允许带病启动。
func TestBootFailsOnBadDSN(t *testing.T) {
// 赋值/更新：vp := viper.New()
	vp := viper.New()
// 执行：vp.Set("shop.mysql-dsn", "root:wrong@tcp(127.0.0.1:3399)/nope?charset=utf8mb4")
	vp.Set("shop.mysql-dsn", "root:wrong@tcp(127.0.0.1:3399)/nope?charset=utf8mb4")
// 赋值/更新：global.GVA_VP = vp
	global.GVA_VP = vp
// 赋值/更新：global.GVA_LOG = nil
	global.GVA_LOG = nil
// 条件判断：if _, _, err := Boot(); err == nil {
	if _, _, err := Boot(); err == nil {
// 执行：t.Fatal("错误 DSN 应启动失败")
		t.Fatal("错误 DSN 应启动失败")
// 结束当前代码块
	}
// 结束当前代码块
}
