// 包声明
package db

// 导入依赖包
import (
// 执行："errors"
	"errors"
// 执行："fmt"
	"fmt"
// 执行："io/fs"
	"io/fs"
// 执行："os"
	"os"
// 执行："testing"
	"testing"
// 执行："testing/fstest"
	"testing/fstest"

// 执行：driver "github.com/go-sql-driver/mysql"
	driver "github.com/go-sql-driver/mysql"
// 执行："gorm.io/driver/mysql"
	"gorm.io/driver/mysql"
// 执行："gorm.io/gorm"
	"gorm.io/gorm"
// 执行："gorm.io/gorm/logger"
	"gorm.io/gorm/logger"

// 执行："shop/shop/config"
	"shop/shop/config"
// 执行："shop/shop/migrations"
	"shop/shop/migrations"
// 结束当前代码块
)

// dsnOnTestHost 返回指向当前测试库主机端口的 DSN：
// 设置了 SHOP_TEST_DSN 时复用其主机端口与口令（容器/远程 MySQL 场景），
// 否则回落到本机 3307 默认值。
func dsnOnTestHost(dbname string) string {
// 条件判断：if dsn := os.Getenv("SHOP_TEST_DSN"); dsn != "" {
	if dsn := os.Getenv("SHOP_TEST_DSN"); dsn != "" {
// 条件判断：if parsed, err := driver.ParseDSN(dsn); err == nil && parsed.Addr != "" {
		if parsed, err := driver.ParseDSN(dsn); err == nil && parsed.Addr != "" {
// 赋值/更新：user := parsed.User
			user := parsed.User
// 条件判断：if user == "" {
			if user == "" {
// 赋值/更新：user = "root"
				user = "root"
// 结束当前代码块
			}
// 赋值/更新：passwd := parsed.Passwd
			passwd := parsed.Passwd
// 返回：return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, passwd, parsed.Addr, dbn...
			return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, passwd, parsed.Addr, dbname)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return "root:root@tcp(127.0.0.1:3307)/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	return "root:root@tcp(127.0.0.1:3307)/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
// 结束当前代码块
}

// 函数定义：func testConfig(t *testing.T) config.Config {
func testConfig(t *testing.T) config.Config {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：cfg := config.Default()
	cfg := config.Default()
// 条件判断：if dsn := os.Getenv("SHOP_TEST_DSN"); dsn != "" {
	if dsn := os.Getenv("SHOP_TEST_DSN"); dsn != "" {
// 赋值/更新：cfg.Mysql.DSN = dsn
		cfg.Mysql.DSN = dsn
// 条件分支：} else {
	} else {
// 赋值/更新：cfg.Mysql.DSN = dsnOnTestHost("shop_test")
		cfg.Mysql.DSN = dsnOnTestHost("shop_test")
// 结束当前代码块
	}
// 赋值/更新：cfg.Mysql.MaxIdleConns = 2
	cfg.Mysql.MaxIdleConns = 2
// 赋值/更新：cfg.Mysql.MaxOpenConns = 4
	cfg.Mysql.MaxOpenConns = 4
// 返回：return cfg
	return cfg
// 结束当前代码块
}

// 函数定义：func TestGormCreateDatabaseAndMigrate(t *testing.T) {
func TestGormCreateDatabaseAndMigrate(t *testing.T) {
// 赋值/更新：cfg := testConfig(t)
	cfg := testConfig(t)
	// 库不存在时自动创建（覆盖 createDatabase 分支）。
	// 使用一次性临时库名，测试结束后 DROP，保持 MySQL 实例整洁。
	cfg.Mysql.DSN = dsnOnTestHost("shop_test_autocreate_tmp")
// 赋值/更新：gdb, err := Gorm(cfg)
	gdb, err := Gorm(cfg)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("自动建库连接失败: %v", err)
		t.Fatalf("自动建库连接失败: %v", err)
// 结束当前代码块
	}
// 执行：t.Cleanup(func() {
	t.Cleanup(func() {
// 赋值/更新：parsed, perr := driver.ParseDSN(cfg.Mysql.DSN)
		parsed, perr := driver.ParseDSN(cfg.Mysql.DSN)
// 条件判断：if perr != nil {
		if perr != nil {
// 返回：return
			return
// 结束当前代码块
		}
// 赋值/更新：parsed.DBName = ""
		parsed.DBName = ""
// 赋值/更新：boot, berr := gorm.Open(mysql.Open(parsed.FormatDSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Sil...
		boot, berr := gorm.Open(mysql.Open(parsed.FormatDSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
// 条件判断：if berr == nil {
		if berr == nil {
// 赋值/更新：_ = boot.Exec("DROP DATABASE IF EXISTS `shop_test_autocreate_tmp`").Error
			_ = boot.Exec("DROP DATABASE IF EXISTS `shop_test_autocreate_tmp`").Error
// 条件判断：if conn, cerr := boot.DB(); cerr == nil {
			if conn, cerr := boot.DB(); cerr == nil {
// 赋值/更新：_ = conn.Close()
				_ = conn.Close()
// 结束当前代码块
			}
// 结束当前代码块
		}
// 结束当前代码块
	})
// 赋值/更新：sqlDB, err := gdb.DB()
	sqlDB, err := gdb.DB()
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 异步/延迟执行：defer func() { _ = sqlDB.Close() }()
	defer func() { _ = sqlDB.Close() }()

	// 空迁移集报错
	if err := Migrate(gdb, fstest.MapFS{}); err == nil {
// 执行：t.Fatal("空迁移集应报错")
		t.Fatal("空迁移集应报错")
// 结束当前代码块
	}
	// 非法 SQL 报错
	bad := fstest.MapFS{"001_bad.sql": &fstest.MapFile{Data: []byte("CREATE TABLE ;")}}
// 条件判断：if err := Migrate(gdb, bad); err == nil {
	if err := Migrate(gdb, bad); err == nil {
// 执行：t.Fatal("非法 SQL 应报错")
		t.Fatal("非法 SQL 应报错")
// 结束当前代码块
	}
	// 正常迁移与种子补齐
	cfg.Mysql.DSN = dsnOnTestHost("shop_test")
// 赋值/更新：gdb2, err := Gorm(cfg)
	gdb2, err := Gorm(cfg)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("连接测试库失败: %v", err)
		t.Fatalf("连接测试库失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := Migrate(gdb2, migrations.FS); err != nil {
	if err := Migrate(gdb2, migrations.FS); err != nil {
// 执行：t.Fatalf("迁移失败: %v", err)
		t.Fatalf("迁移失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := SeedAccounts(gdb2, "shop123"); err != nil {
	if err := SeedAccounts(gdb2, "shop123"); err != nil {
// 执行：t.Fatalf("种子账号失败: %v", err)
		t.Fatalf("种子账号失败: %v", err)
// 结束当前代码块
	}
	// 幂等：重复补齐不报错
	if err := SeedAccounts(gdb2, ""); err != nil {
// 执行：t.Fatalf("重复种子账号失败: %v", err)
		t.Fatalf("重复种子账号失败: %v", err)
// 结束当前代码块
	}
// 声明：var count int64
	var count int64
// 条件判断：if err := gdb2.Table("users").Where("password_hash <> ''").Count(&count).Error; err != nil {
	if err := gdb2.Table("users").Where("password_hash <> ''").Count(&count).Error; err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if count == 0 {
	if count == 0 {
// 执行：t.Fatal("种子账号未补齐密码")
		t.Fatal("种子账号未补齐密码")
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestSplitSQL(t *testing.T) {
func TestSplitSQL(t *testing.T) {
// 赋值/更新：stmts := splitSQL("-- 注释\nCREATE TABLE a (id INT);\n\nINSERT INTO a VALUES (1);\n")
	stmts := splitSQL("-- 注释\nCREATE TABLE a (id INT);\n\nINSERT INTO a VALUES (1);\n")
// 条件判断：if len(stmts) != 2 {
	if len(stmts) != 2 {
// 执行：t.Fatalf("语句切分错误: %d %+v", len(stmts), stmts)
		t.Fatalf("语句切分错误: %d %+v", len(stmts), stmts)
// 结束当前代码块
	}
// 结束当前代码块
}

// errFS 读取即失败，用于覆盖迁移读取错误分支。
type errFS struct{}

// 函数定义：func (errFS) Open(string) (fs.File, error) { return nil, errors.New("read failure") }
func (errFS) Open(string) (fs.File, error) { return nil, errors.New("read failure") }

// 函数定义：func TestGormAndMigrateErrorBranches(t *testing.T) {
func TestGormAndMigrateErrorBranches(t *testing.T) {
	// 非法 DSN：非「库不存在」类错误直接返回
	cfg := testConfig(t)
// 赋值/更新：cfg.Mysql.DSN = "not-a-valid-dsn"
	cfg.Mysql.DSN = "not-a-valid-dsn"
// 条件判断：if _, err := Gorm(cfg); err == nil {
	if _, err := Gorm(cfg); err == nil {
// 执行：t.Fatal("非法 DSN 应报错")
		t.Fatal("非法 DSN 应报错")
// 结束当前代码块
	}
	// DSN 未指定库名：自动建库内部必须报错（直接调内部函数，避免连接默认库成功）
	if err := createDatabase(dsnOnTestHost("")); err == nil {
// 执行：t.Fatal("未指定库名应报错")
		t.Fatal("未指定库名应报错")
// 结束当前代码块
	}

	// 迁移文件读取失败
	cfg = testConfig(t)
// 赋值/更新：gdb, err := Gorm(cfg)
	gdb, err := Gorm(cfg)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("连接测试库失败: %v", err)
		t.Fatalf("连接测试库失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := Migrate(gdb, errFS{}); err == nil {
	if err := Migrate(gdb, errFS{}); err == nil {
// 执行：t.Fatal("文件读取失败应报错")
		t.Fatal("文件读取失败应报错")
// 结束当前代码块
	}

	// 连接关闭后种子补齐必须报错（不允许静默成功）
	sqlDB, err := gdb.DB()
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 赋值/更新：_ = sqlDB.Close()
	_ = sqlDB.Close()
// 条件判断：if err := SeedAccounts(gdb, "shop123"); err == nil {
	if err := SeedAccounts(gdb, "shop123"); err == nil {
// 执行：t.Fatal("连接关闭后应报错")
		t.Fatal("连接关闭后应报错")
// 结束当前代码块
	}
// 结束当前代码块
}
