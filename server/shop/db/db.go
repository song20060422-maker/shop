// Package db 提供 Shop 业务的基础设施：数据库连接、迁移执行与种子账号补齐。
// 独立成包以便测试与初始化层共用，避免测试包反向依赖业务装配层形成导入环。
package db

// 导入依赖包
import (
// 执行："fmt"
	"fmt"
// 执行："io/fs"
	"io/fs"
// 执行："sort"
	"sort"
// 执行："strings"
	"strings"
// 执行："time"
	"time"

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
// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// Gorm 建立连接：库不存在时自动创建，保证首次运行可直接启动。

// 函数定义：func Gorm(cfg config.Config) (*gorm.DB, error) {
func Gorm(cfg config.Config) (*gorm.DB, error) {
// 赋值/更新：newLogger := logger.Default.LogMode(logger.Silent)
	newLogger := logger.Default.LogMode(logger.Silent)
// 条件判断：if cfg.Mysql.LogMode {
	if cfg.Mysql.LogMode {
// 赋值/更新：newLogger = logger.Default.LogMode(logger.Info)
		newLogger = logger.Default.LogMode(logger.Info)
// 结束当前代码块
	}
// 赋值/更新：gdb, err := gorm.Open(mysql.Open(cfg.Mysql.DSN), &gorm.Config{Logger: newLogger})
	gdb, err := gorm.Open(mysql.Open(cfg.Mysql.DSN), &gorm.Config{Logger: newLogger})
// 条件判断：if err != nil {
	if err != nil {
// 条件判断：if !strings.Contains(err.Error(), "Unknown database") {
		if !strings.Contains(err.Error(), "Unknown database") {
// 返回：return nil, err
			return nil, err
// 结束当前代码块
		}
// 条件判断：if cerr := createDatabase(cfg.Mysql.DSN); cerr != nil {
		if cerr := createDatabase(cfg.Mysql.DSN); cerr != nil {
// 返回：return nil, cerr
			return nil, cerr
// 结束当前代码块
		}
// 赋值/更新：gdb, err = gorm.Open(mysql.Open(cfg.Mysql.DSN), &gorm.Config{Logger: newLogger})
		gdb, err = gorm.Open(mysql.Open(cfg.Mysql.DSN), &gorm.Config{Logger: newLogger})
// 条件判断：if err != nil {
		if err != nil {
// 返回：return nil, err
			return nil, err
// 结束当前代码块
		}
// 结束当前代码块
	}
// 赋值/更新：sqlDB, err := gdb.DB()
	sqlDB, err := gdb.DB()
// 条件判断：if err != nil {
	if err != nil {
// 返回：return nil, err
		return nil, err
// 结束当前代码块
	}
// 执行：sqlDB.SetMaxIdleConns(cfg.Mysql.MaxIdleConns)
	sqlDB.SetMaxIdleConns(cfg.Mysql.MaxIdleConns)
// 执行：sqlDB.SetMaxOpenConns(cfg.Mysql.MaxOpenConns)
	sqlDB.SetMaxOpenConns(cfg.Mysql.MaxOpenConns)
// 执行：sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxLifetime(time.Hour)
// 返回：return gdb, nil
	return gdb, nil
// 结束当前代码块
}

// createDatabase 首次启动时按 DSN 建库（不切库、不覆盖已有库）。

// 函数定义：func createDatabase(dsn string) error {
func createDatabase(dsn string) error {
// 赋值/更新：parsed, err := driver.ParseDSN(dsn)
	parsed, err := driver.ParseDSN(dsn)
// 条件判断：if err != nil {
	if err != nil {
// 返回：return err
		return err
// 结束当前代码块
	}
// 赋值/更新：name := parsed.DBName
	name := parsed.DBName
// 条件判断：if name == "" {
	if name == "" {
// 返回：return fmt.Errorf("DSN 未指定数据库名")
		return fmt.Errorf("DSN 未指定数据库名")
// 结束当前代码块
	}
// 赋值/更新：parsed.DBName = ""
	parsed.DBName = ""
// 赋值/更新：boot, err := gorm.Open(mysql.Open(parsed.FormatDSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Sile...
	boot, err := gorm.Open(mysql.Open(parsed.FormatDSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
// 条件判断：if err != nil {
	if err != nil {
// 返回：return err
		return err
// 结束当前代码块
	}
// 赋值/更新：stmt := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_...
	stmt := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", name)
// 条件判断：if err := boot.Exec(stmt).Error; err != nil {
	if err := boot.Exec(stmt).Error; err != nil {
// 返回：return err
		return err
// 结束当前代码块
	}
// 赋值/更新：conn, err := boot.DB()
	conn, err := boot.DB()
// 条件判断：if err == nil {
	if err == nil {
// 赋值/更新：_ = conn.Close()
		_ = conn.Close()
// 结束当前代码块
	}
// 返回：return nil
	return nil
// 结束当前代码块
}

// Migrate 按文件名顺序执行内置 migrations，脚本全部使用 IF NOT EXISTS 可重复执行。

// 函数定义：func Migrate(db *gorm.DB, fsys fs.FS) error {
func Migrate(db *gorm.DB, fsys fs.FS) error {
// 赋值/更新：names, err := fs.Glob(fsys, "*.sql")
	names, err := fs.Glob(fsys, "*.sql")
// 条件判断：if err != nil {
	if err != nil {
// 返回：return err
		return err
// 结束当前代码块
	}
// 条件判断：if len(names) == 0 {
	if len(names) == 0 {
// 返回：return fmt.Errorf("内置迁移脚本为空")
		return fmt.Errorf("内置迁移脚本为空")
// 结束当前代码块
	}
// 执行：sort.Strings(names)
	sort.Strings(names)
// 循环处理：for _, name := range names {
	for _, name := range names {
// 赋值/更新：raw, err := fs.ReadFile(fsys, name)
		raw, err := fs.ReadFile(fsys, name)
// 条件判断：if err != nil {
		if err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 循环处理：for _, stmt := range splitSQL(string(raw)) {
		for _, stmt := range splitSQL(string(raw)) {
// 条件判断：if err := db.Exec(stmt).Error; err != nil {
			if err := db.Exec(stmt).Error; err != nil {
// 返回：return fmt.Errorf("%s 执行失败: %w", name, err)
				return fmt.Errorf("%s 执行失败: %w", name, err)
// 结束当前代码块
			}
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return nil
	return nil
// 结束当前代码块
}

// splitSQL 去掉注释行后按分号切分语句。

// 函数定义：func splitSQL(content string) []string {
func splitSQL(content string) []string {
// 声明：var cleaned strings.Builder
	var cleaned strings.Builder
// 循环处理：for _, line := range strings.Split(content, "\n") {
	for _, line := range strings.Split(content, "\n") {
// 赋值/更新：trimmed := strings.TrimSpace(line)
		trimmed := strings.TrimSpace(line)
// 条件判断：if strings.HasPrefix(trimmed, "--") || trimmed == "" {
		if strings.HasPrefix(trimmed, "--") || trimmed == "" {
// 流程控制：continue
			continue
// 结束当前代码块
		}
// 执行：cleaned.WriteString(line)
		cleaned.WriteString(line)
// 执行：cleaned.WriteString("\n")
		cleaned.WriteString("\n")
// 结束当前代码块
	}
// 声明：var out []string
	var out []string
// 循环处理：for _, stmt := range strings.Split(cleaned.String(), ";") {
	for _, stmt := range strings.Split(cleaned.String(), ";") {
// 条件判断：if s := strings.TrimSpace(stmt); s != "" {
		if s := strings.TrimSpace(stmt); s != "" {
// 赋值/更新：out = append(out, s)
			out = append(out, s)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return out
	return out
// 结束当前代码块
}

// SeedAccounts 为预置账号补齐密码摘要（种子脚本不落明文，也不落固定哈希）。

// 函数定义：func SeedAccounts(db *gorm.DB, demoPassword string) error {
func SeedAccounts(db *gorm.DB, demoPassword string) error {
// 条件判断：if demoPassword == "" {
	if demoPassword == "" {
// 赋值/更新：demoPassword = "shop123"
		demoPassword = "shop123"
// 结束当前代码块
	}
// 声明：var users []model.User
	var users []model.User
// 条件判断：if err := db.Find(&users).Error; err != nil {
	if err := db.Find(&users).Error; err != nil {
// 返回：return err
		return err
// 结束当前代码块
	}
// 循环处理：for _, u := range users {
	for _, u := range users {
// 条件判断：if strings.TrimSpace(u.PasswordHash) != "" {
		if strings.TrimSpace(u.PasswordHash) != "" {
// 流程控制：continue
			continue
// 结束当前代码块
		}
// 赋值/更新：hash, err := utils.HashPassword(demoPassword)
		hash, err := utils.HashPassword(demoPassword)
// 条件判断：if err != nil {
		if err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 条件判断：if err := db.Model(&model.User{}).Where("id = ?", u.ID).Update("password_hash", hash).Error; err != nil {
		if err := db.Model(&model.User{}).Where("id = ?", u.ID).Update("password_hash", hash).Error; err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return nil
	return nil
// 结束当前代码块
}
