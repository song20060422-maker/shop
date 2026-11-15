// 包声明
package service

// 导入依赖包
import (
// 执行："strconv"
	"strconv"
// 执行："strings"
	"strings"
// 执行："time"
	"time"

// 执行："gorm.io/gorm"
	"gorm.io/gorm"

// 执行："shop/shop/model"
	"shop/shop/model"
// 结束当前代码块
)

// 本文件为「接口契约兼容层」补齐的 service 扩展方法。
// 只新增方法，不修改任何既有方法签名与行为，保证原有 /api 原型接口语义不变。

// ProductBySKU 按 SKU 查询商品：契约层以 SKU 作为对外商品 id。

// 函数定义：func (s CatalogService) ProductBySKU(sku string) (model.Product, error) {
func (s CatalogService) ProductBySKU(sku string) (model.Product, error) {
// 声明：var p model.Product
	var p model.Product
// 赋值/更新：sku = strings.TrimSpace(sku)
	sku = strings.TrimSpace(sku)
// 条件判断：if sku == "" {
	if sku == "" {
// 返回：return p, ErrProductNotFound
		return p, ErrProductNotFound
// 结束当前代码块
	}
// 条件判断：if err := s.deps.DB.Where("sku = ?", sku).First(&p).Error; err != nil {
	if err := s.deps.DB.Where("sku = ?", sku).First(&p).Error; err != nil {
// 返回：return p, ErrProductNotFound
		return p, ErrProductNotFound
// 结束当前代码块
	}
// 返回：return p, nil
	return p, nil
// 结束当前代码块
}

// OverwriteAvailable 运营直接覆盖可售数量（契约约定：available 为绝对值，而非增量）。
// 与占用/释放路径一致地追加库存流水，便于运营核对数量来源。

// 函数定义：func (s CatalogService) OverwriteAvailable(actorID uint64, actorRole string, productID uint64, available int) ...
func (s CatalogService) OverwriteAvailable(actorID uint64, actorRole string, productID uint64, available int) (model.Product, error) {
// 条件判断：if available < 0 {
	if available < 0 {
// 返回：return model.Product{}, NewBizError(1006, "可售数量不能小于 0")
		return model.Product{}, NewBizError(1006, "可售数量不能小于 0")
// 结束当前代码块
	}
// 赋值/更新：now := s.deps.Now()
	now := s.deps.Now()
// 赋值/更新：err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
	err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
// 声明：var p model.Product
		var p model.Product
// 条件判断：if err := tx.Where("id = ?", productID).First(&p).Error; err != nil {
		if err := tx.Where("id = ?", productID).First(&p).Error; err != nil {
// 返回：return ErrProductNotFound
			return ErrProductNotFound
// 结束当前代码块
		}
// 条件判断：if p.Available == available {
		if p.Available == available {
// 返回：return logOp(tx, nil, &actorID, actorRole, "product_stock_overwrite",
			return logOp(tx, nil, &actorID, actorRole, "product_stock_overwrite",
// 执行："运营覆盖可售数量："+strconv.Itoa(available)+"（无变化）")
				"运营覆盖可售数量："+strconv.Itoa(available)+"（无变化）")
// 结束当前代码块
		}
// 条件判断：if err := tx.Model(&model.Product{}).Where("id = ?", productID).
		if err := tx.Model(&model.Product{}).Where("id = ?", productID).
// 执行：Updates(map[string]interface{}{"available": available, "updated_at": now}).Error; err != nil {
			Updates(map[string]interface{}{"available": available, "updated_at": now}).Error; err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 条件判断：if err := appendLedger(tx, productID, "", "restore", available-p.Available, now); err != nil {
		if err := appendLedger(tx, productID, "", "restore", available-p.Available, now); err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 返回：return logOp(tx, nil, &actorID, actorRole, "product_stock_overwrite",
		return logOp(tx, nil, &actorID, actorRole, "product_stock_overwrite",
// 执行："运营覆盖可售数量："+strconv.Itoa(p.Available)+" -> "+strconv.Itoa(available))
			"运营覆盖可售数量："+strconv.Itoa(p.Available)+" -> "+strconv.Itoa(available))
// 结束当前代码块
	})
// 条件判断：if err != nil {
	if err != nil {
// 返回：return model.Product{}, err
		return model.Product{}, err
// 结束当前代码块
	}
// 返回：return s.Get(productID)
	return s.Get(productID)
// 结束当前代码块
}

// FindByRef 订单双寻址：先按 order_no，再按数字主键（契约 {id} 同时兼容两者）。

// 函数定义：func (s OrderService) FindByRef(ref string) (model.Order, error) {
func (s OrderService) FindByRef(ref string) (model.Order, error) {
// 赋值/更新：ref = strings.TrimSpace(ref)
	ref = strings.TrimSpace(ref)
// 条件判断：if ref == "" {
	if ref == "" {
// 返回：return model.Order{}, ErrOrderNotFound
		return model.Order{}, ErrOrderNotFound
// 结束当前代码块
	}
// 声明：var o model.Order
	var o model.Order
// 条件判断：if err := s.deps.DB.Where("order_no = ?", ref).First(&o).Error; err == nil {
	if err := s.deps.DB.Where("order_no = ?", ref).First(&o).Error; err == nil {
// 返回：return o, nil
		return o, nil
// 结束当前代码块
	}
// 条件判断：if id, err := strconv.ParseUint(ref, 10, 64); err == nil && id > 0 {
	if id, err := strconv.ParseUint(ref, 10, 64); err == nil && id > 0 {
// 条件判断：if err := s.deps.DB.Where("id = ?", id).First(&o).Error; err == nil {
		if err := s.deps.DB.Where("id = ?", id).First(&o).Error; err == nil {
// 返回：return o, nil
			return o, nil
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return model.Order{}, ErrOrderNotFound
	return model.Order{}, ErrOrderNotFound
// 结束当前代码块
}

// PendingOrderNos 待付款订单号列表，按创建时间倒序（最新在前）。

// 函数定义：func (s OrderService) PendingOrderNos() ([]string, error) {
func (s OrderService) PendingOrderNos() ([]string, error) {
// 赋值/更新：nos := make([]string, 0, 8)
	nos := make([]string, 0, 8)
// 赋值/更新：err := s.deps.DB.Model(&model.Order{}).
	err := s.deps.DB.Model(&model.Order{}).
// 执行：Where("status = ?", model.OrderStatusPendingPayment).
		Where("status = ?", model.OrderStatusPendingPayment).
// 执行：Order("id desc").Pluck("order_no", &nos).Error
		Order("id desc").Pluck("order_no", &nos).Error
// 返回：return nos, err
	return nos, err
// 结束当前代码块
}

// ForceExpire 把指定订单的到期时间提前到 at 之前（测试钩子：让最新一笔待付款订单立即到期）。

// 函数定义：func (s OrderService) ForceExpire(orderNo string, at time.Time) error {
func (s OrderService) ForceExpire(orderNo string, at time.Time) error {
// 赋值/更新：orderNo = strings.TrimSpace(orderNo)
	orderNo = strings.TrimSpace(orderNo)
// 条件判断：if orderNo == "" {
	if orderNo == "" {
// 返回：return ErrOrderNotFound
		return ErrOrderNotFound
// 结束当前代码块
	}
// 返回：return s.deps.DB.Model(&model.Order{}).
	return s.deps.DB.Model(&model.Order{}).
// 执行：Where("order_no = ? AND status = ?", orderNo, model.OrderStatusPendingPayment).
		Where("order_no = ? AND status = ?", orderNo, model.OrderStatusPendingPayment).
// 执行：Updates(map[string]interface{}{
		Updates(map[string]interface{}{
// 执行："expire_at":  at.Add(-2 * time.Second),
			"expire_at":  at.Add(-2 * time.Second),
// 执行："updated_at": s.deps.Now(),
			"updated_at": s.deps.Now(),
// 结束当前代码块
		}).Error
// 结束当前代码块
}

// FindByRef 退款单三寻址：refund_no / order_no / 数字主键（契约 {id} 传退款单号）。

// 函数定义：func (s RefundService) FindByRef(ref string) (model.Refund, error) {
func (s RefundService) FindByRef(ref string) (model.Refund, error) {
// 赋值/更新：ref = strings.TrimSpace(ref)
	ref = strings.TrimSpace(ref)
// 条件判断：if ref == "" {
	if ref == "" {
// 返回：return model.Refund{}, ErrRefundNotFound
		return model.Refund{}, ErrRefundNotFound
// 结束当前代码块
	}
// 声明：var r model.Refund
	var r model.Refund
// 条件判断：if err := s.deps.DB.Where("refund_no = ?", ref).First(&r).Error; err == nil {
	if err := s.deps.DB.Where("refund_no = ?", ref).First(&r).Error; err == nil {
// 返回：return r, nil
		return r, nil
// 结束当前代码块
	}
// 条件判断：if err := s.deps.DB.Where("order_no = ?", ref).Order("id desc").First(&r).Error; err == nil {
	if err := s.deps.DB.Where("order_no = ?", ref).Order("id desc").First(&r).Error; err == nil {
// 返回：return r, nil
		return r, nil
// 结束当前代码块
	}
// 条件判断：if id, err := strconv.ParseUint(ref, 10, 64); err == nil && id > 0 {
	if id, err := strconv.ParseUint(ref, 10, 64); err == nil && id > 0 {
// 条件判断：if err := s.deps.DB.Where("id = ?", id).First(&r).Error; err == nil {
		if err := s.deps.DB.Where("id = ?", id).First(&r).Error; err == nil {
// 返回：return r, nil
			return r, nil
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return model.Refund{}, ErrRefundNotFound
	return model.Refund{}, ErrRefundNotFound
// 结束当前代码块
}

// PaymentsByOrder 查询订单的成功收款记录（契约 /api/ops/payments 的记款部分）。

// 函数定义：func (s PaymentService) PaymentsByOrder(orderNo string) ([]model.Payment, error) {
func (s PaymentService) PaymentsByOrder(orderNo string) ([]model.Payment, error) {
// 赋值/更新：list := make([]model.Payment, 0, 2)
	list := make([]model.Payment, 0, 2)
// 赋值/更新：err := s.deps.DB.Where("order_no = ?", strings.TrimSpace(orderNo)).
	err := s.deps.DB.Where("order_no = ?", strings.TrimSpace(orderNo)).
// 执行：Order("id asc").Find(&list).Error
		Order("id asc").Find(&list).Error
// 返回：return list, err
	return list, err
// 结束当前代码块
}

// AbnormalByID 按主键查询异常收款记录（契约 /api/cs/payment-exceptions/{id}/refund 使用）。

// 函数定义：func (s PaymentService) AbnormalByID(id uint64) (model.AbnormalPayment, error) {
func (s PaymentService) AbnormalByID(id uint64) (model.AbnormalPayment, error) {
// 声明：var ab model.AbnormalPayment
	var ab model.AbnormalPayment
// 条件判断：if err := s.deps.DB.Where("id = ?", id).First(&ab).Error; err != nil {
	if err := s.deps.DB.Where("id = ?", id).First(&ab).Error; err != nil {
// 返回：return ab, ErrAbnormalNotFound
		return ab, ErrAbnormalNotFound
// 结束当前代码块
	}
// 返回：return ab, nil
	return ab, nil
// 结束当前代码块
}
