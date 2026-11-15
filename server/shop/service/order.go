// 包声明
package service

// 导入依赖包
import (
// 执行："strings"
	"strings"
// 执行："time"
	"time"

// 执行："gorm.io/gorm"
	"gorm.io/gorm"

// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// OrderService 下单、取消、超时关闭与查询

// 类型定义：type OrderService struct {
type OrderService struct {
// 执行：deps Deps
	deps Deps
// 结束当前代码块
}

// Create 下单：条件扣减可售数量 → 订单快照 → 库存流水，全部在一个事务内。

// 函数定义：func (s OrderService) Create(buyerID, productID uint64, qty int) (model.Order, error) {
func (s OrderService) Create(buyerID, productID uint64, qty int) (model.Order, error) {
// 声明：var order model.Order
	var order model.Order
// 条件判断：if qty < 1 || qty > 2 {
	if qty < 1 || qty > 2 {
// 返回：return order, ErrQtyInvalid
		return order, ErrQtyInvalid
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
		// 关键并发闸门：仅当在售且可售数量足够时才扣减，靠 InnoDB 行锁保证互斥。
		res := tx.Model(&model.Product{}).
// 执行：Where("id = ? AND on_sale = ? AND available >= ?", productID, true, qty).
			Where("id = ? AND on_sale = ? AND available >= ?", productID, true, qty).
// 执行：Update("available", gorm.Expr("available - ?", qty))
			Update("available", gorm.Expr("available - ?", qty))
// 条件判断：if res.Error != nil {
		if res.Error != nil {
// 返回：return res.Error
			return res.Error
// 结束当前代码块
		}
// 条件判断：if res.RowsAffected != 1 {
		if res.RowsAffected != 1 {
// 条件判断：if !p.OnSale {
			if !p.OnSale {
// 返回：return ErrProductOffSale
				return ErrProductOffSale
// 结束当前代码块
			}
// 返回：return ErrInsufficientStock
			return ErrInsufficientStock
// 结束当前代码块
		}
// 赋值/更新：order = model.Order{
		order = model.Order{
// 执行：OrderNo:     s.uniqueOrderNo(tx, now),
			OrderNo:     s.uniqueOrderNo(tx, now),
// 执行：BuyerID:     buyerID,
			BuyerID:     buyerID,
// 执行：ProductID:   productID,
			ProductID:   productID,
// 执行：ProductName: p.Name,
			ProductName: p.Name,
// 执行：PriceCents:  p.PriceCents,
			PriceCents:  p.PriceCents,
// 执行：Qty:         qty,
			Qty:         qty,
// 执行：AmountCents: p.PriceCents * int64(qty),
			AmountCents: p.PriceCents * int64(qty),
// 执行：Status:      model.OrderStatusPendingPayment,
			Status:      model.OrderStatusPendingPayment,
// 执行：ExpireAt:    now.Add(s.deps.Cfg.OrderTimeout()),
			ExpireAt:    now.Add(s.deps.Cfg.OrderTimeout()),
// 执行：CreatedAt:   now,
			CreatedAt:   now,
// 执行：UpdatedAt:   now,
			UpdatedAt:   now,
// 结束当前代码块
		}
// 条件判断：if err := tx.Create(&order).Error; err != nil {
		if err := tx.Create(&order).Error; err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 条件判断：if err := appendLedger(tx, productID, order.OrderNo, model.LedgerOccupy, -qty, now); err != nil {
		if err := appendLedger(tx, productID, order.OrderNo, model.LedgerOccupy, -qty, now); err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 返回：return logOp(tx, &order.OrderNo, &buyerID, model.RoleBuyer, "order_create",
		return logOp(tx, &order.OrderNo, &buyerID, model.RoleBuyer, "order_create",
// 执行："下单占用:"+order.ProductName+" x"+itoa(qty))
			"下单占用:"+order.ProductName+" x"+itoa(qty))
// 结束当前代码块
	})
// 条件判断：if err != nil {
	if err != nil {
// 返回：return model.Order{}, err
		return model.Order{}, err
// 结束当前代码块
	}
// 返回：return order, nil
	return order, nil
// 结束当前代码块
}

// Cancel 买家取消待付款订单：状态条件更新成功才释放占用，重复关闭不会多释放。

// 函数定义：func (s OrderService) Cancel(buyerID uint64, orderNo string) (model.Order, error) {
func (s OrderService) Cancel(buyerID uint64, orderNo string) (model.Order, error) {
// 赋值/更新：now := s.deps.Now()
	now := s.deps.Now()
// 赋值/更新：orderNo = strings.TrimSpace(orderNo)
	orderNo = strings.TrimSpace(orderNo)
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
// 条件判断：if o.BuyerID != buyerID {
		if o.BuyerID != buyerID {
// 返回：return ErrForbidden
			return ErrForbidden
// 结束当前代码块
		}
// 赋值/更新：res := tx.Model(&model.Order{}).
		res := tx.Model(&model.Order{}).
// 执行：Where("order_no = ? AND buyer_id = ? AND status = ?", orderNo, buyerID, model.OrderStatusPendingPayment).
			Where("order_no = ? AND buyer_id = ? AND status = ?", orderNo, buyerID, model.OrderStatusPendingPayment).
// 执行：Updates(map[string]interface{}{"status": model.OrderStatusCancelled, "closed_at": now, "updated_at": now})
			Updates(map[string]interface{}{"status": model.OrderStatusCancelled, "closed_at": now, "updated_at": now})
// 条件判断：if res.Error != nil {
		if res.Error != nil {
// 返回：return res.Error
			return res.Error
// 结束当前代码块
		}
// 条件判断：if res.RowsAffected != 1 {
		if res.RowsAffected != 1 {
// 返回：return stateError(tx, orderNo)
			return stateError(tx, orderNo)
// 结束当前代码块
		}
// 条件判断：if err := releaseStock(tx, o, model.LedgerRelease, now); err != nil {
		if err := releaseStock(tx, o, model.LedgerRelease, now); err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 条件判断：if err := logOp(tx, &orderNo, &buyerID, model.RoleBuyer, "order_cancel", "买家主动取消并释放占用"); err != nil {
		if err := logOp(tx, &orderNo, &buyerID, model.RoleBuyer, "order_cancel", "买家主动取消并释放占用"); err != nil {
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

// CloseExpired 超时关闭：以服务端时间为准，逐单条件更新后释放占用。

// 函数定义：func (s OrderService) CloseExpired() (int, error) {
func (s OrderService) CloseExpired() (int, error) {
// 赋值/更新：now := s.deps.Now()
	now := s.deps.Now()
// 声明：var candidates []model.Order
	var candidates []model.Order
// 条件判断：if err := s.deps.DB.Where("status = ? AND expire_at <= ?", model.OrderStatusPendingPayment, now).
	if err := s.deps.DB.Where("status = ? AND expire_at <= ?", model.OrderStatusPendingPayment, now).
// 执行：Order("id asc").Limit(200).Find(&candidates).Error; err != nil {
		Order("id asc").Limit(200).Find(&candidates).Error; err != nil {
// 返回：return 0, err
		return 0, err
// 结束当前代码块
	}
// 赋值/更新：closed := 0
	closed := 0
// 循环处理：for _, c := range candidates {
	for _, c := range candidates {
// 赋值/更新：err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
		err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
// 赋值/更新：res := tx.Model(&model.Order{}).
			res := tx.Model(&model.Order{}).
// 执行：Where("order_no = ? AND status = ? AND expire_at <= ?", c.OrderNo, model.OrderStatusPendingPayment, now).
				Where("order_no = ? AND status = ? AND expire_at <= ?", c.OrderNo, model.OrderStatusPendingPayment, now).
// 执行：Updates(map[string]interface{}{"status": model.OrderStatusClosed, "closed_at": now, "updated_at": now})
				Updates(map[string]interface{}{"status": model.OrderStatusClosed, "closed_at": now, "updated_at": now})
// 条件判断：if res.Error != nil {
			if res.Error != nil {
// 返回：return res.Error
				return res.Error
// 结束当前代码块
			}
// 条件判断：if res.RowsAffected != 1 {
			if res.RowsAffected != 1 {
// 返回：return nil // 已被支付或已取消，跳过，不释放数量
				return nil // 已被支付或已取消，跳过，不释放数量
// 结束当前代码块
			}
// 条件判断：if err := releaseStock(tx, c, model.LedgerRelease, now); err != nil {
			if err := releaseStock(tx, c, model.LedgerRelease, now); err != nil {
// 返回：return err
				return err
// 结束当前代码块
			}
// 返回：return logOp(tx, &c.OrderNo, nil, "system", "order_timeout_close", "服务端超时关闭并释放占用")
			return logOp(tx, &c.OrderNo, nil, "system", "order_timeout_close", "服务端超时关闭并释放占用")
// 结束当前代码块
		})
// 条件判断：if err != nil {
		if err != nil {
// 返回：return closed, err
			return closed, err
// 结束当前代码块
		}
// 执行：closed++
		closed++
// 结束当前代码块
	}
// 返回：return closed, nil
	return closed, nil
// 结束当前代码块
}

// Get 按单号读取订单（不做视图权限判断，调用方负责）

// 函数定义：func (s OrderService) Get(orderNo string) (model.Order, error) {
func (s OrderService) Get(orderNo string) (model.Order, error) {
// 声明：var o model.Order
	var o model.Order
// 条件判断：if err := s.deps.DB.Where("order_no = ?", strings.TrimSpace(orderNo)).First(&o).Error; err != nil {
	if err := s.deps.DB.Where("order_no = ?", strings.TrimSpace(orderNo)).First(&o).Error; err != nil {
// 返回：return o, ErrOrderNotFound
		return o, ErrOrderNotFound
// 结束当前代码块
	}
// 返回：return o, nil
	return o, nil
// 结束当前代码块
}

// ListByBuyer 买家自己的订单

// 函数定义：func (s OrderService) ListByBuyer(buyerID uint64, status string) ([]model.Order, error) {
func (s OrderService) ListByBuyer(buyerID uint64, status string) ([]model.Order, error) {
// 声明：var list []model.Order
	var list []model.Order
// 赋值/更新：q := s.deps.DB.Where("buyer_id = ?", buyerID)
	q := s.deps.DB.Where("buyer_id = ?", buyerID)
// 条件判断：if status != "" {
	if status != "" {
// 赋值/更新：q = q.Where("status = ?", status)
		q = q.Where("status = ?", status)
// 结束当前代码块
	}
// 赋值/更新：err := q.Order("id desc").Find(&list).Error
	err := q.Order("id desc").Find(&list).Error
// 返回：return list, err
	return list, err
// 结束当前代码块
}

// ListByStatus 员工视角按状态查询订单

// 函数定义：func (s OrderService) ListByStatus(status string, limit int) ([]model.Order, error) {
func (s OrderService) ListByStatus(status string, limit int) ([]model.Order, error) {
// 声明：var list []model.Order
	var list []model.Order
// 赋值/更新：q := s.deps.DB
	q := s.deps.DB
// 条件判断：if status != "" {
	if status != "" {
// 赋值/更新：q = q.Where("status = ?", status)
		q = q.Where("status = ?", status)
// 结束当前代码块
	}
// 条件判断：if limit <= 0 {
	if limit <= 0 {
// 赋值/更新：limit = 100
		limit = 100
// 结束当前代码块
	}
// 赋值/更新：err := q.Order("id desc").Limit(limit).Find(&list).Error
	err := q.Order("id desc").Limit(limit).Find(&list).Error
// 返回：return list, err
	return list, err
// 结束当前代码块
}

// releaseStock 释放/归还可售数量：以该单库存流水净额为守恒依据，杜绝重复释放造成的虚增。
// 说明：运营允许把可售数量覆盖到高于 initial_stock（如活动补货），因此不能用 initial_stock 作为上界，
// 而是校验「该订单仍有数量处于占用状态」（净额为负）。

// 函数定义：func releaseStock(tx *gorm.DB, o model.Order, changeType string, now time.Time) error {
func releaseStock(tx *gorm.DB, o model.Order, changeType string, now time.Time) error {
// 声明：var net int
	var net int
// 条件判断：if err := tx.Model(&model.InventoryLedger{}).
	if err := tx.Model(&model.InventoryLedger{}).
// 执行：Select("COALESCE(SUM(delta), 0)").Where("order_no = ?", o.OrderNo).Scan(&net).Error; err != nil {
		Select("COALESCE(SUM(delta), 0)").Where("order_no = ?", o.OrderNo).Scan(&net).Error; err != nil {
// 返回：return err
		return err
// 结束当前代码块
	}
// 条件判断：if net > -o.Qty {
	if net > -o.Qty {
		// 该单占用已归还过，重复释放会虚增可售数量，直接回滚
		return NewBizError(2009, "库存归还不合法：可能已重复释放")
// 结束当前代码块
	}
// 条件判断：if err := tx.Model(&model.Product{}).
	if err := tx.Model(&model.Product{}).
// 执行：Where("id = ?", o.ProductID).
		Where("id = ?", o.ProductID).
// 执行：Update("available", gorm.Expr("available + ?", o.Qty)).Error; err != nil {
		Update("available", gorm.Expr("available + ?", o.Qty)).Error; err != nil {
// 返回：return err
		return err
// 结束当前代码块
	}
// 返回：return appendLedger(tx, o.ProductID, o.OrderNo, changeType, o.Qty, now)
	return appendLedger(tx, o.ProductID, o.OrderNo, changeType, o.Qty, now)
// 结束当前代码块
}

// stateError 根据订单当前状态给出准确原因，避免笼统报错。

// 函数定义：func stateError(tx *gorm.DB, orderNo string) error {
func stateError(tx *gorm.DB, orderNo string) error {
// 声明：var o model.Order
	var o model.Order
// 条件判断：if err := tx.Where("order_no = ?", orderNo).First(&o).Error; err != nil {
	if err := tx.Where("order_no = ?", orderNo).First(&o).Error; err != nil {
// 返回：return ErrOrderNotFound
		return ErrOrderNotFound
// 结束当前代码块
	}
// 多路分支：switch o.Status {
	switch o.Status {
// 匹配分支：case model.OrderStatusPaid, model.OrderStatusShipped, model.OrderStatusRefundRequested,
	case model.OrderStatusPaid, model.OrderStatusShipped, model.OrderStatusRefundRequested,
// 执行：model.OrderStatusRefundFailed, model.OrderStatusRefunded:
		model.OrderStatusRefundFailed, model.OrderStatusRefunded:
// 返回：return ErrAlreadyPaid
		return ErrAlreadyPaid
// 匹配分支：case model.OrderStatusClosed, model.OrderStatusCancelled:
	case model.OrderStatusClosed, model.OrderStatusCancelled:
// 返回：return ErrOrderClosed
		return ErrOrderClosed
// 默认分支
	default:
// 返回：return ErrOrderNotPayable
		return ErrOrderNotPayable
// 结束当前代码块
	}
// 结束当前代码块
}

// uniqueOrderNo 生成未占用的订单号（极端碰撞下重试）

// 函数定义：func (s OrderService) uniqueOrderNo(tx *gorm.DB, now time.Time) string {
func (s OrderService) uniqueOrderNo(tx *gorm.DB, now time.Time) string {
// 循环处理：for i := 0; i < 5; i++ {
	for i := 0; i < 5; i++ {
// 赋值/更新：no := utils.NewOrderNo(now)
		no := utils.NewOrderNo(now)
// 声明：var count int64
		var count int64
// 执行：tx.Model(&model.Order{}).Where("order_no = ?", no).Count(&count)
		tx.Model(&model.Order{}).Where("order_no = ?", no).Count(&count)
// 条件判断：if count == 0 {
		if count == 0 {
// 返回：return no
			return no
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return utils.NewOrderNo(now.Add(time.Duration(1) * time.Nanosecond))
	return utils.NewOrderNo(now.Add(time.Duration(1) * time.Nanosecond))
// 结束当前代码块
}

// itoa 整数转字符串（避免引入 fmt，错误码拼接路径高频调用）。

// 函数定义：func itoa(v int) string {
func itoa(v int) string {
// 条件判断：if v == 0 {
	if v == 0 {
// 返回：return "0"
		return "0"
// 结束当前代码块
	}
// 赋值/更新：neg := v < 0
	neg := v < 0
// 条件判断：if neg {
	if neg {
// 赋值/更新：v = -v
		v = -v
// 结束当前代码块
	}
// 声明：var buf [20]byte
	var buf [20]byte
// 赋值/更新：i := len(buf)
	i := len(buf)
// 循环处理：for v > 0 {
	for v > 0 {
// 执行：i--
		i--
// 赋值/更新：buf[i] = byte('0' + v%10)
		buf[i] = byte('0' + v%10)
// 赋值/更新：v /= 10
		v /= 10
// 结束当前代码块
	}
// 条件判断：if neg {
	if neg {
// 执行：i--
		i--
// 赋值/更新：buf[i] = '-'
		buf[i] = '-'
// 结束当前代码块
	}
// 返回：return string(buf[i:])
	return string(buf[i:])
// 结束当前代码块
}
