// 包声明
package service

// 导入依赖包
import (
// 执行："strings"
	"strings"

// 执行："gorm.io/gorm"
	"gorm.io/gorm"

// 执行："shop/shop/model"
	"shop/shop/model"
// 结束当前代码块
)

// ShipmentService 仓库发货：只为合法待发货订单登记一次。

// 类型定义：type ShipmentService struct {
type ShipmentService struct {
// 执行：deps Deps
	deps Deps
// 结束当前代码块
}

// Shipable 待发货列表：仅 status=paid（退款申请/已退款/未付款都不在其中）

// 函数定义：func (s ShipmentService) Shipable() ([]model.Order, error) {
func (s ShipmentService) Shipable() ([]model.Order, error) {
// 声明：var list []model.Order
	var list []model.Order
// 赋值/更新：err := s.deps.DB.Where("status = ?", model.OrderStatusPaid).Order("id asc").Find(&list).Error
	err := s.deps.DB.Where("status = ?", model.OrderStatusPaid).Order("id asc").Find(&list).Error
// 返回：return list, err
	return list, err
// 结束当前代码块
}

// Ship 登记发货：状态条件更新 + 唯一键双保险，确保同一订单只发一次。

// 函数定义：func (s ShipmentService) Ship(whID uint64, orderNo, carrier, trackingNo string) (model.Order, error) {
func (s ShipmentService) Ship(whID uint64, orderNo, carrier, trackingNo string) (model.Order, error) {
// 赋值/更新：orderNo = strings.TrimSpace(orderNo)
	orderNo = strings.TrimSpace(orderNo)
// 赋值/更新：carrier = strings.TrimSpace(carrier)
	carrier = strings.TrimSpace(carrier)
// 赋值/更新：trackingNo = strings.TrimSpace(trackingNo)
	trackingNo = strings.TrimSpace(trackingNo)
// 条件判断：if orderNo == "" || carrier == "" || trackingNo == "" {
	if orderNo == "" || carrier == "" || trackingNo == "" {
// 返回：return model.Order{}, ErrInvalidParam
		return model.Order{}, ErrInvalidParam
// 结束当前代码块
	}
// 赋值/更新：now := s.deps.Now()
	now := s.deps.Now()
// 声明：var out model.Order
	var out model.Order
// 赋值/更新：err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
	err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
// 声明：var o model.Order
		var o model.Order
// 条件判断：if err := tx.Where("order_no = ?", orderNo).First(&o).Error; err != nil {
		if err := tx.Where("order_no = ?", orderNo).First(&o).Error; err != nil {
// 返回：return ErrOrderNotFound
			return ErrOrderNotFound
// 结束当前代码块
		}
// 赋值/更新：res := tx.Model(&model.Order{}).
		res := tx.Model(&model.Order{}).
// 执行：Where("order_no = ? AND status = ?", orderNo, model.OrderStatusPaid).
			Where("order_no = ? AND status = ?", orderNo, model.OrderStatusPaid).
// 执行：Updates(map[string]interface{}{"status": model.OrderStatusShipped, "shipped_at": now, "updated_at": now})
			Updates(map[string]interface{}{"status": model.OrderStatusShipped, "shipped_at": now, "updated_at": now})
// 条件判断：if res.Error != nil {
		if res.Error != nil {
// 返回：return res.Error
			return res.Error
// 结束当前代码块
		}
// 条件判断：if res.RowsAffected != 1 {
		if res.RowsAffected != 1 {
// 返回：return shipStateError(o.Status)
			return shipStateError(o.Status)
// 结束当前代码块
		}
// 赋值/更新：shipment := model.Shipment{
		shipment := model.Shipment{
// 执行：OrderNo:    orderNo,
			OrderNo:    orderNo,
// 执行：Carrier:    carrier,
			Carrier:    carrier,
// 执行：TrackingNo: trackingNo,
			TrackingNo: trackingNo,
// 执行：CreatedBy:  &whID,
			CreatedBy:  &whID,
// 执行：CreatedAt:  now,
			CreatedAt:  now,
// 结束当前代码块
		}
// 条件判断：if err := tx.Create(&shipment).Error; err != nil {
		if err := tx.Create(&shipment).Error; err != nil {
// 返回：return ErrAlreadyShipped
			return ErrAlreadyShipped
// 结束当前代码块
		}
// 条件判断：if err := appendLedger(tx, o.ProductID, orderNo, model.LedgerShipOut, 0, now); err != nil {
		if err := appendLedger(tx, o.ProductID, orderNo, model.LedgerShipOut, 0, now); err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 条件判断：if err := logOp(tx, &orderNo, &whID, model.RoleWarehouse, "shipment_create",
		if err := logOp(tx, &orderNo, &whID, model.RoleWarehouse, "shipment_create",
// 赋值/更新："发货登记："+carrier+" "+trackingNo); err != nil {
			"发货登记："+carrier+" "+trackingNo); err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 返回：return tx.Where("order_no = ?", orderNo).First(&out).Error
		return tx.Where("order_no = ?", orderNo).First(&out).Error
// 结束当前代码块
	})
// 条件判断：if err != nil {
	if err != nil {
// 返回：return model.Order{}, err
		return model.Order{}, err
// 结束当前代码块
	}
// 返回：return out, nil
	return out, nil
// 结束当前代码块
}

// Get 查询发货记录（买家可查运单）

// 函数定义：func (s ShipmentService) Get(orderNo string) (*model.Shipment, error) {
func (s ShipmentService) Get(orderNo string) (*model.Shipment, error) {
// 声明：var sh model.Shipment
	var sh model.Shipment
// 条件判断：if err := s.deps.DB.Where("order_no = ?", strings.TrimSpace(orderNo)).First(&sh).Error; err != nil {
	if err := s.deps.DB.Where("order_no = ?", strings.TrimSpace(orderNo)).First(&sh).Error; err != nil {
// 返回：return nil, nil
		return nil, nil
// 结束当前代码块
	}
// 返回：return &sh, nil
	return &sh, nil
// 结束当前代码块
}

// List 全部发货记录（运营核查）

// 函数定义：func (s ShipmentService) List() ([]model.Shipment, error) {
func (s ShipmentService) List() ([]model.Shipment, error) {
// 声明：var list []model.Shipment
	var list []model.Shipment
// 赋值/更新：err := s.deps.DB.Order("id desc").Limit(200).Find(&list).Error
	err := s.deps.DB.Order("id desc").Limit(200).Find(&list).Error
// 返回：return list, err
	return list, err
// 结束当前代码块
}

// shipStateError 把订单当前状态翻译为仓库看得懂的原因。

// 函数定义：func shipStateError(status string) error {
func shipStateError(status string) error {
// 多路分支：switch status {
	switch status {
// 匹配分支：case model.OrderStatusRefundRequested:
	case model.OrderStatusRefundRequested:
// 返回：return ErrRefundInProgress
		return ErrRefundInProgress
// 匹配分支：case model.OrderStatusRefunded:
	case model.OrderStatusRefunded:
// 返回：return ErrRefunded
		return ErrRefunded
// 匹配分支：case model.OrderStatusRefundFailed:
	case model.OrderStatusRefundFailed:
// 返回：return NewBizError(4011, "退款渠道处理失败，订单暂不可发货")
		return NewBizError(4011, "退款渠道处理失败，订单暂不可发货")
// 匹配分支：case model.OrderStatusPendingPayment:
	case model.OrderStatusPendingPayment:
// 返回：return ErrNotPaid
		return ErrNotPaid
// 匹配分支：case model.OrderStatusClosed, model.OrderStatusCancelled:
	case model.OrderStatusClosed, model.OrderStatusCancelled:
// 返回：return ErrOrderClosed
		return ErrOrderClosed
// 匹配分支：case model.OrderStatusShipped:
	case model.OrderStatusShipped:
// 返回：return ErrAlreadyShipped
		return ErrAlreadyShipped
// 默认分支
	default:
// 返回：return ErrNotShippable
		return ErrNotShippable
// 结束当前代码块
	}
// 结束当前代码块
}
