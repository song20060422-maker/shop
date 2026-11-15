// 包声明
package service

// 导入依赖包
import (
// 执行："gorm.io/gorm"
	"gorm.io/gorm"

// 执行："shop/shop/model"
	"shop/shop/model"
// 结束当前代码块
)

// DemoService 演示数据重置：仅清理业务数据并恢复商品初始可售数量，保留账号。
// 由 shop.demo-mode 配置开关保护，避免生产误用。

// 类型定义：type DemoService struct {
type DemoService struct {
// 执行：deps Deps
	deps Deps
// 结束当前代码块
}

// Reset 按外键依赖顺序清理业务表，并把商品可售数量恢复到初始库存。

// 函数定义：func (s DemoService) Reset() error {
func (s DemoService) Reset() error {
// 返回：return s.deps.DB.Transaction(func(tx *gorm.DB) error {
	return s.deps.DB.Transaction(func(tx *gorm.DB) error {
		// 先删子表，再删订单，避免外键约束失败
		tables := []string{
// 执行："payments",
			"payments",
// 执行："payment_notices",
			"payment_notices",
// 执行："refunds",
			"refunds",
// 执行："shipments",
			"shipments",
// 执行："abnormal_payments",
			"abnormal_payments",
// 执行："inventory_ledger",
			"inventory_ledger",
// 执行："operation_logs",
			"operation_logs",
// 执行："orders",
			"orders",
// 结束当前代码块
		}
// 循环处理：for _, t := range tables {
		for _, t := range tables {
// 条件判断：if err := tx.Exec("DELETE FROM " + t).Error; err != nil {
			if err := tx.Exec("DELETE FROM " + t).Error; err != nil {
// 返回：return err
				return err
// 结束当前代码块
			}
// 结束当前代码块
		}
// 条件判断：if err := tx.Model(&model.Product{}).Where("1 = 1").
		if err := tx.Model(&model.Product{}).Where("1 = 1").
// 执行：Updates(map[string]interface{}{
			Updates(map[string]interface{}{
// 执行："available": gorm.Expr("initial_stock"),
				"available": gorm.Expr("initial_stock"),
// 执行："on_sale":   true,
				"on_sale":   true,
// 结束当前代码块
			}).Error; err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 返回：return nil
		return nil
// 结束当前代码块
	})
// 结束当前代码块
}
