// 包声明
package service

// 导入依赖包
import (
// 执行："time"
	"time"

// 执行："gorm.io/gorm"
	"gorm.io/gorm"

// 执行："shop/shop/model"
	"shop/shop/model"
// 结束当前代码块
)

// CatalogService 商品上下架、可售数量与库存流水

// 类型定义：type CatalogService struct {
type CatalogService struct {
// 执行：deps Deps
	deps Deps
// 结束当前代码块
}

// OnSale 买家可见的在售商品

// 函数定义：func (s CatalogService) OnSale() ([]model.Product, error) {
func (s CatalogService) OnSale() ([]model.Product, error) {
// 声明：var list []model.Product
	var list []model.Product
// 赋值/更新：err := s.deps.DB.Where("on_sale = ?", true).Order("id asc").Find(&list).Error
	err := s.deps.DB.Where("on_sale = ?", true).Order("id asc").Find(&list).Error
// 返回：return list, err
	return list, err
// 结束当前代码块
}

// All 运营可见的全部商品

// 函数定义：func (s CatalogService) All() ([]model.Product, error) {
func (s CatalogService) All() ([]model.Product, error) {
// 声明：var list []model.Product
	var list []model.Product
// 赋值/更新：err := s.deps.DB.Order("id asc").Find(&list).Error
	err := s.deps.DB.Order("id asc").Find(&list).Error
// 返回：return list, err
	return list, err
// 结束当前代码块
}

// Get 单个商品

// 函数定义：func (s CatalogService) Get(id uint64) (model.Product, error) {
func (s CatalogService) Get(id uint64) (model.Product, error) {
// 声明：var p model.Product
	var p model.Product
// 条件判断：if err := s.deps.DB.Where("id = ?", id).First(&p).Error; err != nil {
	if err := s.deps.DB.Where("id = ?", id).First(&p).Error; err != nil {
// 返回：return p, ErrProductNotFound
		return p, ErrProductNotFound
// 结束当前代码块
	}
// 返回：return p, nil
	return p, nil
// 结束当前代码块
}

// Update 运营维护上下架与售价；已存在订单的价格快照不受影响。

// 函数定义：func (s CatalogService) Update(actorID uint64, actorRole string, id uint64, onSale *bool, price *int64) (model...
func (s CatalogService) Update(actorID uint64, actorRole string, id uint64, onSale *bool, price *int64) (model.Product, error) {
// 条件判断：if actorRole != model.RoleOps {
	if actorRole != model.RoleOps {
// 返回：return model.Product{}, ErrForbidden
		return model.Product{}, ErrForbidden
// 结束当前代码块
	}
// 条件判断：if price != nil && *price <= 0 {
	if price != nil && *price <= 0 {
// 返回：return model.Product{}, NewBizError(1005, "售价必须大于 0")
		return model.Product{}, NewBizError(1005, "售价必须大于 0")
// 结束当前代码块
	}
// 赋值/更新：err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
	err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
// 声明：var p model.Product
		var p model.Product
// 条件判断：if err := tx.Where("id = ?", id).First(&p).Error; err != nil {
		if err := tx.Where("id = ?", id).First(&p).Error; err != nil {
// 返回：return ErrProductNotFound
			return ErrProductNotFound
// 结束当前代码块
		}
// 赋值/更新：updates := map[string]interface{}{}
		updates := map[string]interface{}{}
// 条件判断：if onSale != nil {
		if onSale != nil {
// 赋值/更新：updates["on_sale"] = *onSale
			updates["on_sale"] = *onSale
// 结束当前代码块
		}
// 条件判断：if price != nil {
		if price != nil {
// 赋值/更新：updates["price_cents"] = *price
			updates["price_cents"] = *price
// 结束当前代码块
		}
// 条件判断：if len(updates) > 0 {
		if len(updates) > 0 {
// 条件判断：if err := tx.Model(&model.Product{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			if err := tx.Model(&model.Product{}).Where("id = ?", id).Updates(updates).Error; err != nil {
// 返回：return err
				return err
// 结束当前代码块
			}
// 结束当前代码块
		}
// 赋值/更新：detail := "运营维护商品"
		detail := "运营维护商品"
// 条件判断：if onSale != nil {
		if onSale != nil {
// 条件判断：if *onSale {
			if *onSale {
// 赋值/更新：detail += "，上架"
				detail += "，上架"
// 条件分支：} else {
			} else {
// 赋值/更新：detail += "，下架"
				detail += "，下架"
// 结束当前代码块
			}
// 结束当前代码块
		}
// 条件判断：if price != nil {
		if price != nil {
// 赋值/更新：detail += "，调整售价"
			detail += "，调整售价"
// 结束当前代码块
		}
// 返回：return logOp(tx, nil, &actorID, actorRole, "product_update", detail)
		return logOp(tx, nil, &actorID, actorRole, "product_update", detail)
// 结束当前代码块
	})
// 条件判断：if err != nil {
	if err != nil {
// 返回：return model.Product{}, err
		return model.Product{}, err
// 结束当前代码块
	}
// 返回：return s.Get(id)
	return s.Get(id)
// 结束当前代码块
}

// appendLedger 记录库存流水：先抢占/归还可售数量，再读回实时值写入 available_after。

// 函数定义：func appendLedger(tx *gorm.DB, productID uint64, orderNo, changeType string, delta int, now time.Time) error {
func appendLedger(tx *gorm.DB, productID uint64, orderNo, changeType string, delta int, now time.Time) error {
// 声明：var available int
	var available int
// 条件判断：if err := tx.Model(&model.Product{}).Select("available").Where("id = ?", productID).Scan(&available).Error; er...
	if err := tx.Model(&model.Product{}).Select("available").Where("id = ?", productID).Scan(&available).Error; err != nil {
// 返回：return err
		return err
// 结束当前代码块
	}
// 返回：return tx.Create(&model.InventoryLedger{
	return tx.Create(&model.InventoryLedger{
// 执行：ProductID:      productID,
		ProductID:      productID,
// 执行：OrderNo:        orderNo,
		OrderNo:        orderNo,
// 执行：ChangeType:     changeType,
		ChangeType:     changeType,
// 执行：Delta:          delta,
		Delta:          delta,
// 执行：AvailableAfter: available,
		AvailableAfter: available,
// 执行：CreatedAt:      now,
		CreatedAt:      now,
// 结束当前代码块
	}).Error
// 结束当前代码块
}

// logOp 追加操作留痕（只增不改）

// 函数定义：func logOp(tx *gorm.DB, orderNo *string, actorID *uint64, actorRole, action, detail string) error {
func logOp(tx *gorm.DB, orderNo *string, actorID *uint64, actorRole, action, detail string) error {
// 返回：return tx.Create(&model.OperationLog{
	return tx.Create(&model.OperationLog{
// 执行：OrderNo:   orderNo,
		OrderNo:   orderNo,
// 执行：ActorID:   actorID,
		ActorID:   actorID,
// 执行：ActorRole: actorRole,
		ActorRole: actorRole,
// 执行：Action:    action,
		Action:    action,
// 执行：Detail:    detail,
		Detail:    detail,
// 结束当前代码块
	}).Error
// 结束当前代码块
}

// CurrentAvailable 读取商品当前可售数量（页面展示与演示断言使用）

// 函数定义：func (s CatalogService) CurrentAvailable(id uint64) (int, error) {
func (s CatalogService) CurrentAvailable(id uint64) (int, error) {
// 声明：var available int
	var available int
// 条件判断：if err := s.deps.DB.Model(&model.Product{}).Select("available").Where("id = ?", id).Scan(&available).Error; er...
	if err := s.deps.DB.Model(&model.Product{}).Select("available").Where("id = ?", id).Scan(&available).Error; err != nil {
// 返回：return 0, err
		return 0, err
// 结束当前代码块
	}
// 返回：return available, nil
	return available, nil
// 结束当前代码块
}
