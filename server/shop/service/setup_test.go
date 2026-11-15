// 包声明
package service

// 导入依赖包
import (
// 执行："os"
	"os"
// 执行："testing"
	"testing"
// 执行："time"
	"time"

// 执行："gorm.io/gorm"
	"gorm.io/gorm"

// 执行："shop/shop/config"
	"shop/shop/config"
// 执行：shopdb "shop/shop/db"
	shopdb "shop/shop/db"
// 执行："shop/shop/migrations"
	"shop/shop/migrations"
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// 说明：service 层测试使用真实 MySQL（本机 3307），每个用例前清空业务数据并重放迁移/种子。
// 数据库可用 SHOP_TEST_DSN 覆盖，默认库名 shop_test（全项目唯一测试库；
// 因各测试包共用同一库，多包并发跑会互相清表，必须使用 go test -p 1 串行执行）。

// 声明：var testDB *gorm.DB
var testDB *gorm.DB
// 声明：var testCfg config.Config
var testCfg config.Config

// 函数定义：func TestMain(m *testing.M) {
func TestMain(m *testing.M) {
// 赋值/更新：dsn := os.Getenv("SHOP_TEST_DSN")
	dsn := os.Getenv("SHOP_TEST_DSN")
// 条件判断：if dsn == "" {
	if dsn == "" {
// 赋值/更新：dsn = "root:root@tcp(127.0.0.1:3307)/shop_test?charset=utf8mb4&parseTime=True&loc=Local"
		dsn = "root:root@tcp(127.0.0.1:3307)/shop_test?charset=utf8mb4&parseTime=True&loc=Local"
// 结束当前代码块
	}
// 赋值/更新：testCfg = config.Default()
	testCfg = config.Default()
// 赋值/更新：testCfg.Mysql.DSN = dsn
	testCfg.Mysql.DSN = dsn
// 赋值/更新：testCfg.Shop.PayNotifySecret = "test-secret"
	testCfg.Shop.PayNotifySecret = "test-secret"
// 赋值/更新：testCfg.Shop.OrderTimeoutMinutes = 10
	testCfg.Shop.OrderTimeoutMinutes = 10
// 赋值/更新：db, err := shopdb.Gorm(testCfg)
	db, err := shopdb.Gorm(testCfg)
// 条件判断：if err != nil {
	if err != nil {
// 流程控制：panic("连接测试库失败（请确认 MySQL 3307 可用）: " + err.Error())
		panic("连接测试库失败（请确认 MySQL 3307 可用）: " + err.Error())
// 结束当前代码块
	}
// 条件判断：if err := shopdb.Migrate(db, migrations.FS); err != nil {
	if err := shopdb.Migrate(db, migrations.FS); err != nil {
// 流程控制：panic("测试库迁移失败: " + err.Error())
		panic("测试库迁移失败: " + err.Error())
// 结束当前代码块
	}
// 条件判断：if err := shopdb.SeedAccounts(db, "shop123"); err != nil {
	if err := shopdb.SeedAccounts(db, "shop123"); err != nil {
// 流程控制：panic("种子账号失败: " + err.Error())
		panic("种子账号失败: " + err.Error())
// 结束当前代码块
	}
// 赋值/更新：testDB = db
	testDB = db
// 执行：os.Exit(m.Run())
	os.Exit(m.Run())
// 结束当前代码块
}

// 类型定义：type fakeClock struct{ t time.Time }
type fakeClock struct{ t time.Time }

// 函数定义：func (c *fakeClock) Now() time.Time          { return c.t }
func (c *fakeClock) Now() time.Time          { return c.t }
// 函数定义：func (c *fakeClock) Advance(d time.Duration) { c.t = c.t.Add(d) }
func (c *fakeClock) Advance(d time.Duration) { c.t = c.t.Add(d) }

// 函数定义：func resetDB(t *testing.T) {
func resetDB(t *testing.T) {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：tables := []string{
	tables := []string{
		// sessions 须最先清理：fk_sessions_user 引用 users，且共享库中其他包的
		// 登录测试会留下会话行，先删才能保证 users 可清。
		"sessions",
// 执行："abnormal_payments", "payments", "payment_notices", "refunds", "shipments",
		"abnormal_payments", "payments", "payment_notices", "refunds", "shipments",
// 执行："inventory_ledger", "operation_logs", "orders", "products", "users",
		"inventory_ledger", "operation_logs", "orders", "products", "users",
// 结束当前代码块
	}
// 循环处理：for _, tbl := range tables {
	for _, tbl := range tables {
// 条件判断：if err := testDB.Exec("DELETE FROM " + tbl).Error; err != nil {
		if err := testDB.Exec("DELETE FROM " + tbl).Error; err != nil {
// 执行：t.Fatalf("清理 %s 失败: %v", tbl, err)
			t.Fatalf("清理 %s 失败: %v", tbl, err)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if err := shopdb.Migrate(testDB, migrations.FS); err != nil {
	if err := shopdb.Migrate(testDB, migrations.FS); err != nil {
// 执行：t.Fatalf("重放迁移失败: %v", err)
		t.Fatalf("重放迁移失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := shopdb.SeedAccounts(testDB, "shop123"); err != nil {
	if err := shopdb.SeedAccounts(testDB, "shop123"); err != nil {
// 执行：t.Fatalf("重放种子账号失败: %v", err)
		t.Fatalf("重放种子账号失败: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func newEnv(t *testing.T) (*ServiceGroup, *SimulatedChannel, *fakeClock) {
func newEnv(t *testing.T) (*ServiceGroup, *SimulatedChannel, *fakeClock) {
// 执行：t.Helper()
	t.Helper()
// 执行：resetDB(t)
	resetDB(t)
// 赋值/更新：clock := &fakeClock{t: time.Date(2026, 11, 12, 10, 0, 0, 0, time.Local)}
	clock := &fakeClock{t: time.Date(2026, 11, 12, 10, 0, 0, 0, time.Local)}
// 赋值/更新：ch := NewSimulatedChannel(true)
	ch := NewSimulatedChannel(true)
// 赋值/更新：cfg := testCfg
	cfg := testCfg
// 赋值/更新：g := NewServiceGroup(Deps{DB: testDB, Now: clock.Now, Cfg: cfg, Channel: ch})
	g := NewServiceGroup(Deps{DB: testDB, Now: clock.Now, Cfg: cfg, Channel: ch})
// 返回：return g, ch, clock
	return g, ch, clock
// 结束当前代码块
}

// 函数定义：func buyerID(t *testing.T, username string) uint64 {
func buyerID(t *testing.T, username string) uint64 {
// 执行：t.Helper()
	t.Helper()
// 声明：var id uint64
	var id uint64
// 条件判断：if err := testDB.Table("users").Where("username = ?", username).Pluck("id", &id).Error; err != nil {
	if err := testDB.Table("users").Where("username = ?", username).Pluck("id", &id).Error; err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if id == 0 {
	if id == 0 {
// 执行：t.Fatalf("未找到账号 %s", username)
		t.Fatalf("未找到账号 %s", username)
// 结束当前代码块
	}
// 返回：return id
	return id
// 结束当前代码块
}

// 函数定义：func productID(t *testing.T, sku string) uint64 {
func productID(t *testing.T, sku string) uint64 {
// 执行：t.Helper()
	t.Helper()
// 声明：var id uint64
	var id uint64
// 条件判断：if err := testDB.Table("products").Where("sku = ?", sku).Pluck("id", &id).Error; err != nil {
	if err := testDB.Table("products").Where("sku = ?", sku).Pluck("id", &id).Error; err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if id == 0 {
	if id == 0 {
// 执行：t.Fatalf("未找到商品 %s", sku)
		t.Fatalf("未找到商品 %s", sku)
// 结束当前代码块
	}
// 返回：return id
	return id
// 结束当前代码块
}

// 函数定义：func available(t *testing.T, sku string) int {
func available(t *testing.T, sku string) int {
// 执行：t.Helper()
	t.Helper()
// 声明：var n int
	var n int
// 条件判断：if err := testDB.Table("products").Where("sku = ?", sku).Pluck("available", &n).Error; err != nil {
	if err := testDB.Table("products").Where("sku = ?", sku).Pluck("available", &n).Error; err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 返回：return n
	return n
// 结束当前代码块
}

// 函数定义：func payOrder(t *testing.T, g *ServiceGroup, buyer uint64, qty int) (orderNo, txn string, amount int64) {
func payOrder(t *testing.T, g *ServiceGroup, buyer uint64, qty int) (orderNo, txn string, amount int64) {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：o, err := g.OrderService.Create(buyer, pid, qty)
	o, err := g.OrderService.Create(buyer, pid, qty)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("下单失败: %v", err)
		t.Fatalf("下单失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：txn = "TXN-" + o.OrderNo
	txn = "TXN-" + o.OrderNo
// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, txn, o.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, o.OrderNo, txn, o.AmountCents)
// 赋值/更新：res, err := g.PaymentService.Notify(o.OrderNo, txn, o.AmountCents, sign, "{}")
	res, err := g.PaymentService.Notify(o.OrderNo, txn, o.AmountCents, sign, "{}")
// 条件判断：if err != nil || res.Result != "accepted" {
	if err != nil || res.Result != "accepted" {
// 执行：t.Fatalf("支付通知失败: %v %+v", err, res)
		t.Fatalf("支付通知失败: %v %+v", err, res)
// 结束当前代码块
	}
// 返回：return o.OrderNo, txn, o.AmountCents
	return o.OrderNo, txn, o.AmountCents
// 结束当前代码块
}
