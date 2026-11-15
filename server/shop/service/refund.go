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
// 执行："gorm.io/gorm/clause"
	"gorm.io/gorm/clause"

// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// RefundService 售前退款与迟到支付退款

// 类型定义：type RefundService struct {
type RefundService struct {
// 执行：deps Deps
	deps Deps
// 结束当前代码块
}

// Apply 买家对「已付款未发货」订单申请整单退款；申请期间订单不在待发货列表，天然暂停发货。

// 函数定义：func (s RefundService) Apply(buyerID uint64, orderNo, reason string) (model.Refund, error) {
func (s RefundService) Apply(buyerID uint64, orderNo, reason string) (model.Refund, error) {
// 赋值/更新：now := s.deps.Now()
	now := s.deps.Now()
// 赋值/更新：orderNo = strings.TrimSpace(orderNo)
	orderNo = strings.TrimSpace(orderNo)
// 声明：var out model.Refund
	var out model.Refund
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
// 执行：Where("order_no = ? AND buyer_id = ? AND status = ?", orderNo, buyerID, model.OrderStatusPaid).
			Where("order_no = ? AND buyer_id = ? AND status = ?", orderNo, buyerID, model.OrderStatusPaid).
// 执行：Updates(map[string]interface{}{"status": model.OrderStatusRefundRequested, "updated_at": now})
			Updates(map[string]interface{}{"status": model.OrderStatusRefundRequested, "updated_at": now})
// 条件判断：if res.Error != nil {
		if res.Error != nil {
// 返回：return res.Error
			return res.Error
// 结束当前代码块
		}
// 条件判断：if res.RowsAffected != 1 {
		if res.RowsAffected != 1 {
// 多路分支：switch o.Status {
			switch o.Status {
// 匹配分支：case model.OrderStatusShipped:
			case model.OrderStatusShipped:
// 返回：return NewBizError(4010, "订单已发货，不受理售前退款")
				return NewBizError(4010, "订单已发货，不受理售前退款")
// 匹配分支：case model.OrderStatusRefunded:
			case model.OrderStatusRefunded:
// 返回：return ErrRefunded
				return ErrRefunded
// 匹配分支：case model.OrderStatusRefundRequested, model.OrderStatusRefundFailed:
			case model.OrderStatusRefundRequested, model.OrderStatusRefundFailed:
// 返回：return ErrRefundExists
				return ErrRefundExists
// 匹配分支：case model.OrderStatusPendingPayment:
			case model.OrderStatusPendingPayment:
// 返回：return ErrNotPaid
				return ErrNotPaid
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

// 声明：var existing model.Refund
		var existing model.Refund
// 赋值/更新：findErr := tx.Where("order_no = ?", orderNo).First(&existing).Error
		findErr := tx.Where("order_no = ?", orderNo).First(&existing).Error
// 条件判断：if findErr == nil {
		if findErr == nil {
			// 历史上被客服拒绝的单据可重新申请，复用同一张退款单（整单退款一单一张）
			if existing.ChannelStatus == model.ChannelSuccess {
// 返回：return ErrRefundExists
				return ErrRefundExists
// 结束当前代码块
			}
// 条件判断：if err := tx.Model(&model.Refund{}).Where("id = ?", existing.ID).
			if err := tx.Model(&model.Refund{}).Where("id = ?", existing.ID).
// 执行：Updates(map[string]interface{}{
				Updates(map[string]interface{}{
// 执行："audit_status": model.AuditPending,
					"audit_status": model.AuditPending,
// 执行："audit_reason": "",
					"audit_reason": "",
// 执行："audit_by":     nil,
					"audit_by":     nil,
// 执行："audit_at":     nil,
					"audit_at":     nil,
// 执行："updated_at":   now,
					"updated_at":   now,
// 结束当前代码块
				}).Error; err != nil {
// 返回：return err
				return err
// 结束当前代码块
			}
// 赋值/更新：out = existing
			out = existing
// 赋值/更新：out.AuditStatus = model.AuditPending
			out.AuditStatus = model.AuditPending
// 返回：return logOp(tx, &orderNo, &buyerID, model.RoleBuyer, "refund_apply", "重新申请售前退款")
			return logOp(tx, &orderNo, &buyerID, model.RoleBuyer, "refund_apply", "重新申请售前退款")
// 结束当前代码块
		}

// 赋值/更新：out = model.Refund{
		out = model.Refund{
// 执行：RefundNo:      s.uniqueRefundNo(tx, now),
			RefundNo:      s.uniqueRefundNo(tx, now),
// 执行：OrderNo:       orderNo,
			OrderNo:       orderNo,
// 执行：Source:        model.RefundSourcePresale,
			Source:        model.RefundSourcePresale,
// 执行：AmountCents:   o.AmountCents,
			AmountCents:   o.AmountCents,
// 执行：AuditStatus:   model.AuditPending,
			AuditStatus:   model.AuditPending,
// 执行：AuditReason:   reason,
			AuditReason:   reason,
// 执行：ChannelStatus: model.ChannelPending,
			ChannelStatus: model.ChannelPending,
// 执行：CreatedAt:     now,
			CreatedAt:     now,
// 执行：UpdatedAt:     now,
			UpdatedAt:     now,
// 结束当前代码块
		}
// 条件判断：if err := tx.Create(&out).Error; err != nil {
		if err := tx.Create(&out).Error; err != nil {
// 返回：return ErrRefundExists
			return ErrRefundExists
// 结束当前代码块
		}
// 返回：return logOp(tx, &orderNo, &buyerID, model.RoleBuyer, "refund_apply", "买家申请售前退款："+reason)
		return logOp(tx, &orderNo, &buyerID, model.RoleBuyer, "refund_apply", "买家申请售前退款："+reason)
// 结束当前代码块
	})
// 条件判断：if err != nil {
	if err != nil {
// 返回：return model.Refund{}, err
		return model.Refund{}, err
// 结束当前代码块
	}
// 返回：return out, nil
	return out, nil
// 结束当前代码块
}

// Audit 客服审核：同意则调用渠道退款，拒绝则订单回到可发货状态。
// 题面要求客服同意/拒绝都必须写明原因，空原因直接拒绝，保证处理记录可追溯。

// 函数定义：func (s RefundService) Audit(csID uint64, orderNo string, approve bool, reason string) (model.Refund, error) {
func (s RefundService) Audit(csID uint64, orderNo string, approve bool, reason string) (model.Refund, error) {
// 赋值/更新：now := s.deps.Now()
	now := s.deps.Now()
// 赋值/更新：orderNo = strings.TrimSpace(orderNo)
	orderNo = strings.TrimSpace(orderNo)
// 赋值/更新：reason = strings.TrimSpace(reason)
	reason = strings.TrimSpace(reason)
// 条件判断：if reason == "" {
	if reason == "" {
// 返回：return model.Refund{}, ErrReasonRequired
		return model.Refund{}, ErrReasonRequired
// 结束当前代码块
	}
// 声明：var refund model.Refund
	var refund model.Refund
// 赋值/更新：err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
	err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
// 条件判断：if err := tx.Where("order_no = ?", orderNo).First(&refund).Error; err != nil {
		if err := tx.Where("order_no = ?", orderNo).First(&refund).Error; err != nil {
// 返回：return ErrRefundNotFound
			return ErrRefundNotFound
// 结束当前代码块
		}
// 条件判断：if refund.AuditStatus != model.AuditPending {
		if refund.AuditStatus != model.AuditPending {
// 返回：return ErrRefundAudited
			return ErrRefundAudited
// 结束当前代码块
		}
// 赋值/更新：status := model.AuditRejected
		status := model.AuditRejected
// 条件判断：if approve {
		if approve {
// 赋值/更新：status = model.AuditApproved
			status = model.AuditApproved
// 结束当前代码块
		}
// 赋值/更新：res := tx.Model(&model.Refund{}).
		res := tx.Model(&model.Refund{}).
// 执行：Where("order_no = ? AND audit_status = ?", orderNo, model.AuditPending).
			Where("order_no = ? AND audit_status = ?", orderNo, model.AuditPending).
// 执行：Updates(map[string]interface{}{
			Updates(map[string]interface{}{
// 执行："audit_status": status,
				"audit_status": status,
// 执行："audit_reason": reason,
				"audit_reason": reason,
// 执行："audit_by":     csID,
				"audit_by":     csID,
// 执行："audit_at":     now,
				"audit_at":     now,
// 执行："updated_at":   now,
				"updated_at":   now,
// 结束当前代码块
			})
// 条件判断：if res.Error != nil {
		if res.Error != nil {
// 返回：return res.Error
			return res.Error
// 结束当前代码块
		}
// 条件判断：if res.RowsAffected != 1 {
		if res.RowsAffected != 1 {
// 返回：return ErrRefundAudited
			return ErrRefundAudited
// 结束当前代码块
		}
// 条件判断：if !approve {
		if !approve {
			// 拒绝后允许仓库继续发货
			if err := tx.Model(&model.Order{}).
// 执行：Where("order_no = ? AND status = ?", orderNo, model.OrderStatusRefundRequested).
				Where("order_no = ? AND status = ?", orderNo, model.OrderStatusRefundRequested).
// 执行：Updates(map[string]interface{}{"status": model.OrderStatusPaid, "updated_at": now}).Error; err != nil {
				Updates(map[string]interface{}{"status": model.OrderStatusPaid, "updated_at": now}).Error; err != nil {
// 返回：return err
				return err
// 结束当前代码块
			}
// 返回：return logOp(tx, &orderNo, &csID, model.RoleCS, "refund_reject", "客服拒绝退款："+reason)
			return logOp(tx, &orderNo, &csID, model.RoleCS, "refund_reject", "客服拒绝退款："+reason)
// 结束当前代码块
		}
// 返回：return logOp(tx, &orderNo, &csID, model.RoleCS, "refund_approve", "客服同意退款："+reason)
		return logOp(tx, &orderNo, &csID, model.RoleCS, "refund_approve", "客服同意退款："+reason)
// 结束当前代码块
	})
// 条件判断：if err != nil {
	if err != nil {
// 返回：return model.Refund{}, err
		return model.Refund{}, err
// 结束当前代码块
	}
// 条件判断：if !approve {
	if !approve {
// 赋值/更新：refund.AuditStatus = model.AuditRejected
		refund.AuditStatus = model.AuditRejected
// 返回：return refund, nil
		return refund, nil
// 结束当前代码块
	}
// 返回：return s.callChannel(refund, csID, model.ChannelPending, now)
	return s.callChannel(refund, csID, model.ChannelPending, now)
// 结束当前代码块
}

// HandleAbnormal 客服处理迟到支付：发起异常退款，且不恢复已释放的数量。
// I-18 并发防御：整段（锁定异常记录 -> 建/复用退款单 -> 调渠道 -> 推进）在同一事务内完成，
// 事务内行锁串行化同一交易号的并发处理，杜绝并发重复发起渠道退款。
// 模拟渠道为同步内存实现，事务内调用安全；真实渠道应改为闸门状态机+异步回调。

// 函数定义：func (s RefundService) HandleAbnormal(csID uint64, txnID, reason string) (model.Refund, error) {
func (s RefundService) HandleAbnormal(csID uint64, txnID, reason string) (model.Refund, error) {
// 赋值/更新：now := s.deps.Now()
	now := s.deps.Now()
// 赋值/更新：txnID = strings.TrimSpace(txnID)
	txnID = strings.TrimSpace(txnID)
// 赋值/更新：reason = strings.TrimSpace(reason)
	reason = strings.TrimSpace(reason)
// 条件判断：if reason == "" {
	if reason == "" {
// 返回：return model.Refund{}, ErrReasonRequired
		return model.Refund{}, ErrReasonRequired
// 结束当前代码块
	}
// 声明：var refund model.Refund
	var refund model.Refund
// 声明：var abnormalID uint64
	var abnormalID uint64
// 赋值/更新：ok := false
	ok := false
// 赋值/更新：err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
	err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
		// I-18：行锁读取异常记录，并发处理同一交易号时后到者将阻塞至前者提交
		var ab model.AbnormalPayment
// 条件判断：if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
// 执行：Where("channel_txn_id = ?", txnID).First(&ab).Error; err != nil {
			Where("channel_txn_id = ?", txnID).First(&ab).Error; err != nil {
// 返回：return ErrAbnormalNotFound
			return ErrAbnormalNotFound
// 结束当前代码块
		}
// 条件判断：if ab.Status == model.AbnormalHandled {
		if ab.Status == model.AbnormalHandled {
// 返回：return NewBizError(3006, "该异常收款已处理")
			return NewBizError(3006, "该异常收款已处理")
// 结束当前代码块
		}
// 声明：var o model.Order
		var o model.Order
// 条件判断：if err := tx.Where("order_no = ?", ab.OrderNo).First(&o).Error; err != nil {
		if err := tx.Where("order_no = ?", ab.OrderNo).First(&o).Error; err != nil {
// 返回：return ErrOrderNotFound
			return ErrOrderNotFound
// 结束当前代码块
		}
// 声明：var existing model.Refund
		var existing model.Refund
// 赋值/更新：findErr := tx.Where("order_no = ?", ab.OrderNo).First(&existing).Error
		findErr := tx.Where("order_no = ?", ab.OrderNo).First(&existing).Error
// 条件判断：if findErr == nil {
		if findErr == nil {
			// 整单退款一单一张：上次渠道失败留下的退款单直接复用重试，不得重复建单
			if existing.ChannelStatus == model.ChannelSuccess {
// 返回：return ErrRefundExists
				return ErrRefundExists
// 结束当前代码块
			}
// 条件判断：if err := tx.Model(&model.Refund{}).Where("id = ?", existing.ID).
			if err := tx.Model(&model.Refund{}).Where("id = ?", existing.ID).
// 执行：Updates(map[string]interface{}{
				Updates(map[string]interface{}{
// 执行："audit_status":   model.AuditApproved,
					"audit_status":   model.AuditApproved,
// 执行："audit_reason":   reason,
					"audit_reason":   reason,
// 执行："audit_by":       csID,
					"audit_by":       csID,
// 执行："audit_at":       now,
					"audit_at":       now,
// 执行："channel_status": model.ChannelPending,
					"channel_status": model.ChannelPending,
// 执行："updated_at":     now,
					"updated_at":     now,
// 结束当前代码块
				}).Error; err != nil {
// 返回：return err
				return err
// 结束当前代码块
			}
// 赋值/更新：refund = existing
			refund = existing
// 赋值/更新：refund.AuditStatus = model.AuditApproved
			refund.AuditStatus = model.AuditApproved
// 赋值/更新：refund.AuditReason = reason
			refund.AuditReason = reason
// 赋值/更新：refund.ChannelStatus = model.ChannelPending
			refund.ChannelStatus = model.ChannelPending
// 条件分支：} else {
		} else {
// 赋值/更新：refund = model.Refund{
			refund = model.Refund{
// 执行：RefundNo:      s.uniqueRefundNo(tx, now),
				RefundNo:      s.uniqueRefundNo(tx, now),
// 执行：OrderNo:       ab.OrderNo,
				OrderNo:       ab.OrderNo,
// 执行：Source:        model.RefundSourceAbnormalLate,
				Source:        model.RefundSourceAbnormalLate,
// 执行：AmountCents:   ab.AmountCents,
				AmountCents:   ab.AmountCents,
// 执行：AuditStatus:   model.AuditApproved, // 客服已核实才发起，无需二次审核
				AuditStatus:   model.AuditApproved, // 客服已核实才发起，无需二次审核
// 执行：AuditReason:   reason,
				AuditReason:   reason,
// 执行：AuditBy:       &csID,
				AuditBy:       &csID,
// 执行：AuditAt:       &now,
				AuditAt:       &now,
// 执行：ChannelStatus: model.ChannelPending,
				ChannelStatus: model.ChannelPending,
// 执行：CreatedAt:     now,
				CreatedAt:     now,
// 执行：UpdatedAt:     now,
				UpdatedAt:     now,
// 结束当前代码块
			}
// 条件判断：if err := tx.Create(&refund).Error; err != nil {
			if err := tx.Create(&refund).Error; err != nil {
// 返回：return ErrRefundExists
				return ErrRefundExists
// 结束当前代码块
			}
// 结束当前代码块
		}
// 赋值/更新：abnormalID = ab.ID
		abnormalID = ab.ID
		// I-18：行锁持有期间调用渠道——同一异常记录的并发处理在事务层互斥，
		// 渠道只会被真实受理一次；失败时异常记录保留待处理（可重试）。
		var cerr error
// 赋值/更新：ok, cerr = s.deps.Channel.Refund(refund.OrderNo, refund.RefundNo, refund.AmountCents)
		ok, cerr = s.deps.Channel.Refund(refund.OrderNo, refund.RefundNo, refund.AmountCents)
// 条件判断：if cerr != nil {
		if cerr != nil {
// 赋值/更新：ok = false
			ok = false
// 结束当前代码块
		}
// 赋值/更新：channelStatus := model.ChannelFailed
		channelStatus := model.ChannelFailed
// 条件判断：if ok {
		if ok {
// 赋值/更新：channelStatus = model.ChannelSuccess
			channelStatus = model.ChannelSuccess
// 结束当前代码块
		}
// 条件判断：if err := tx.Model(&model.Refund{}).Where("refund_no = ?", refund.RefundNo).
		if err := tx.Model(&model.Refund{}).Where("refund_no = ?", refund.RefundNo).
// 执行：Updates(map[string]interface{}{"channel_status": channelStatus, "updated_at": now}).Error; err != nil {
			Updates(map[string]interface{}{"channel_status": channelStatus, "updated_at": now}).Error; err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 条件判断：if !ok {
		if !ok {
// 返回：return logOp(tx, &refund.OrderNo, &csID, model.RoleCS, "abnormal_refund_failed",
			return logOp(tx, &refund.OrderNo, &csID, model.RoleCS, "abnormal_refund_failed",
// 执行："渠道退款失败，异常记录保留待处理")
				"渠道退款失败，异常记录保留待处理")
// 结束当前代码块
		}
		// 迟到支付场景：订单早已关闭、数量早已释放，这里绝不再恢复库存
		if err := tx.Model(&model.AbnormalPayment{}).Where("id = ?", abnormalID).
// 执行：Updates(map[string]interface{}{
			Updates(map[string]interface{}{
// 执行："status":     model.AbnormalHandled,
				"status":     model.AbnormalHandled,
// 执行："handled_by": csID,
				"handled_by": csID,
// 执行："handled_at": now,
				"handled_at": now,
// 执行："refund_id":  refund.ID,
				"refund_id":  refund.ID,
// 结束当前代码块
			}).Error; err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 返回：return logOp(tx, &refund.OrderNo, &csID, model.RoleCS, "abnormal_refund_success",
		return logOp(tx, &refund.OrderNo, &csID, model.RoleCS, "abnormal_refund_success",
// 执行："异常收款退款成功，未重复释放库存")
			"异常收款退款成功，未重复释放库存")
// 结束当前代码块
	})
// 条件判断：if err != nil {
	if err != nil {
// 返回：return model.Refund{}, err
		return model.Refund{}, err
// 结束当前代码块
	}
// 赋值/更新：refund.ChannelStatus = model.ChannelSuccess
	refund.ChannelStatus = model.ChannelSuccess
// 条件判断：if !ok {
	if !ok {
// 赋值/更新：refund.ChannelStatus = model.ChannelFailed
		refund.ChannelStatus = model.ChannelFailed
// 返回：return refund, ErrChannelRefund
		return refund, ErrChannelRefund
// 结束当前代码块
	}
// 返回：return refund, nil
	return refund, nil
// 结束当前代码块
}

// Retry 渠道明确失败后由客服重试

// 函数定义：func (s RefundService) Retry(csID uint64, refundNo string) (model.Refund, error) {
func (s RefundService) Retry(csID uint64, refundNo string) (model.Refund, error) {
// 赋值/更新：now := s.deps.Now()
	now := s.deps.Now()
// 声明：var refund model.Refund
	var refund model.Refund
// 条件判断：if err := s.deps.DB.Where("refund_no = ?", strings.TrimSpace(refundNo)).First(&refund).Error; err != nil {
	if err := s.deps.DB.Where("refund_no = ?", strings.TrimSpace(refundNo)).First(&refund).Error; err != nil {
// 返回：return refund, ErrRefundNotFound
		return refund, ErrRefundNotFound
// 结束当前代码块
	}
// 条件判断：if refund.AuditStatus != model.AuditApproved || refund.ChannelStatus != model.ChannelFailed {
	if refund.AuditStatus != model.AuditApproved || refund.ChannelStatus != model.ChannelFailed {
// 返回：return refund, ErrRefundNotRetryable
		return refund, ErrRefundNotRetryable
// 结束当前代码块
	}
// 返回：return s.callChannel(refund, csID, model.ChannelFailed, now)
	return s.callChannel(refund, csID, model.ChannelFailed, now)
// 结束当前代码块
}

// callChannel 调用模拟渠道并根据结果推进订单、退款单与库存。
// I-18 并发防御：事务内行锁读取退款单并校验入口渠道状态（pending=首次发起，failed=重试），
// 校验通过才真正调渠道——并发 Retry/审核到达同一退款单时，后到者在行锁上阻塞，
// 解锁后读到的渠道状态已变更为 success/failed，闸门拒绝其再次发起渠道退款。
// 模拟渠道为同步内存实现，事务内调用安全；真实渠道应改为闸门状态机+异步回调。

// 函数定义：func (s RefundService) callChannel(refund model.Refund, csID uint64, expectedChannel string, now time.Time) (m...
func (s RefundService) callChannel(refund model.Refund, csID uint64, expectedChannel string, now time.Time) (model.Refund, error) {
// 赋值/更新：ok := false
	ok := false
// 声明：var result model.Refund
	var result model.Refund
// 赋值/更新：err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
	err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
// 条件判断：if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
// 执行：Where("refund_no = ?", refund.RefundNo).First(&result).Error; err != nil {
			Where("refund_no = ?", refund.RefundNo).First(&result).Error; err != nil {
// 返回：return ErrRefundNotFound
			return ErrRefundNotFound
// 结束当前代码块
		}
		// 闸门：审核已通过且渠道状态等于入口期望，否则拒绝发起（含并发竞争失败方）
		if result.AuditStatus != model.AuditApproved || result.ChannelStatus != expectedChannel {
// 返回：return ErrRefundNotRetryable
			return ErrRefundNotRetryable
// 结束当前代码块
		}
// 声明：var cerr error
		var cerr error
// 赋值/更新：ok, cerr = s.deps.Channel.Refund(result.OrderNo, result.RefundNo, result.AmountCents)
		ok, cerr = s.deps.Channel.Refund(result.OrderNo, result.RefundNo, result.AmountCents)
// 条件判断：if cerr != nil {
		if cerr != nil {
// 赋值/更新：ok = false
			ok = false
// 结束当前代码块
		}
// 赋值/更新：retry := result.RetryCount + 1
		retry := result.RetryCount + 1
// 条件判断：if !ok {
		if !ok {
// 条件判断：if err := tx.Model(&model.Refund{}).Where("refund_no = ?", refund.RefundNo).
			if err := tx.Model(&model.Refund{}).Where("refund_no = ?", refund.RefundNo).
// 执行：Updates(map[string]interface{}{
				Updates(map[string]interface{}{
// 执行："channel_status": model.ChannelFailed,
					"channel_status": model.ChannelFailed,
// 执行："retry_count":    retry,
					"retry_count":    retry,
// 执行："updated_at":     now,
					"updated_at":     now,
// 结束当前代码块
				}).Error; err != nil {
// 返回：return err
				return err
// 结束当前代码块
			}
// 条件判断：if err := tx.Model(&model.Order{}).
			if err := tx.Model(&model.Order{}).
// 执行：Where("order_no = ? AND status = ?", refund.OrderNo, model.OrderStatusRefundRequested).
				Where("order_no = ? AND status = ?", refund.OrderNo, model.OrderStatusRefundRequested).
// 执行：Updates(map[string]interface{}{"status": model.OrderStatusRefundFailed, "updated_at": now}).Error; err != nil ...
				Updates(map[string]interface{}{"status": model.OrderStatusRefundFailed, "updated_at": now}).Error; err != nil {
// 返回：return err
				return err
// 结束当前代码块
			}
// 返回：return logOp(tx, &refund.OrderNo, &csID, model.RoleCS, "refund_channel_failed",
			return logOp(tx, &refund.OrderNo, &csID, model.RoleCS, "refund_channel_failed",
// 执行："渠道退款失败，保留待处理记录等待重试")
				"渠道退款失败，保留待处理记录等待重试")
// 结束当前代码块
		}

		// 渠道成功：按退款来源分流推进（增量续写【修改原因=R4】：
		// abnormal_late 对应订单早已关闭，原逻辑要求订单处于可退款状态导致
		// 异常退款单重试必然 RowsAffected=0 整体回滚，违反 ABN-06 可重试要求）。
		if result.Source == model.RefundSourceAbnormalLate {
			// 迟到支付场景：不改订单状态（订单保持 closed/cancelled）、
			// 绝不恢复库存（红线 2/3）；异常记录 pending -> handled，
			// RowsAffected==1 闸门保证同一异常记录不会被二次推进。
			res := tx.Model(&model.AbnormalPayment{}).
// 执行：Where("order_no = ? AND status = ?", refund.OrderNo, model.AbnormalPending).
				Where("order_no = ? AND status = ?", refund.OrderNo, model.AbnormalPending).
// 执行：Updates(map[string]interface{}{
				Updates(map[string]interface{}{
// 执行："status":     model.AbnormalHandled,
					"status":     model.AbnormalHandled,
// 执行："handled_by": csID,
					"handled_by": csID,
// 执行："handled_at": now,
					"handled_at": now,
// 执行："refund_id":  result.ID,
					"refund_id":  result.ID,
// 结束当前代码块
				})
// 条件判断：if res.Error != nil {
			if res.Error != nil {
// 返回：return res.Error
				return res.Error
// 结束当前代码块
			}
// 条件判断：if res.RowsAffected != 1 {
			if res.RowsAffected != 1 {
// 返回：return NewBizError(3010, "异常收款不在待处理状态，退款推进被拒绝")
				return NewBizError(3010, "异常收款不在待处理状态，退款推进被拒绝")
// 结束当前代码块
			}
// 条件分支：} else {
		} else {
			// 售前退款：订单进入已退款；库存恢复受 restored 闸门保护，只做一次。
			// I-18：RowsAffected 必须为 1——订单不在可退款状态（重复推进/并发竞争）时整体回滚，
			// 退款单不落成功、库存不恢复。
			res := tx.Model(&model.Order{}).
// 执行：Where("order_no = ? AND status IN ?", refund.OrderNo,
				Where("order_no = ? AND status IN ?", refund.OrderNo,
// 执行：[]string{model.OrderStatusRefundRequested, model.OrderStatusRefundFailed}).
					[]string{model.OrderStatusRefundRequested, model.OrderStatusRefundFailed}).
// 执行：Updates(map[string]interface{}{"status": model.OrderStatusRefunded, "updated_at": now})
				Updates(map[string]interface{}{"status": model.OrderStatusRefunded, "updated_at": now})
// 条件判断：if res.Error != nil {
			if res.Error != nil {
// 返回：return res.Error
				return res.Error
// 结束当前代码块
			}
// 条件判断：if res.RowsAffected != 1 {
			if res.RowsAffected != 1 {
// 返回：return NewBizError(3010, "订单不在可退款状态，退款推进被拒绝")
				return NewBizError(3010, "订单不在可退款状态，退款推进被拒绝")
// 结束当前代码块
			}
// 结束当前代码块
		}
// 赋值/更新：updates := map[string]interface{}{
		updates := map[string]interface{}{
// 执行："channel_status": model.ChannelSuccess,
			"channel_status": model.ChannelSuccess,
// 执行："retry_count":    retry,
			"retry_count":    retry,
// 执行："updated_at":     now,
			"updated_at":     now,
// 结束当前代码块
		}
// 赋值/更新：restore := false
		restore := false
// 条件判断：if result.Source == model.RefundSourcePresale && !result.Restored {
		if result.Source == model.RefundSourcePresale && !result.Restored {
// 赋值/更新：gate := tx.Model(&model.Refund{}).
			gate := tx.Model(&model.Refund{}).
// 执行：Where("refund_no = ? AND restored = ?", refund.RefundNo, false).
				Where("refund_no = ? AND restored = ?", refund.RefundNo, false).
// 执行：Update("restored", true)
				Update("restored", true)
// 条件判断：if gate.Error != nil {
			if gate.Error != nil {
// 返回：return gate.Error
				return gate.Error
// 结束当前代码块
			}
// 赋值/更新：restore = gate.RowsAffected == 1
			restore = gate.RowsAffected == 1
// 条件判断：if restore {
			if restore {
// 赋值/更新：updates["restored"] = true
				updates["restored"] = true
// 结束当前代码块
			}
// 结束当前代码块
		}
// 条件判断：if err := tx.Model(&model.Refund{}).Where("refund_no = ?", refund.RefundNo).Updates(updates).Error; err != nil...
		if err := tx.Model(&model.Refund{}).Where("refund_no = ?", refund.RefundNo).Updates(updates).Error; err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 条件判断：if restore {
		if restore {
// 声明：var o model.Order
			var o model.Order
// 条件判断：if err := tx.Where("order_no = ?", refund.OrderNo).First(&o).Error; err != nil {
			if err := tx.Where("order_no = ?", refund.OrderNo).First(&o).Error; err != nil {
// 返回：return err
				return err
// 结束当前代码块
			}
// 条件判断：if err := releaseStock(tx, o, model.LedgerRestore, now); err != nil {
			if err := releaseStock(tx, o, model.LedgerRestore, now); err != nil {
// 返回：return err
				return err
// 结束当前代码块
			}
// 结束当前代码块
		}
// 返回：return logOp(tx, &refund.OrderNo, &csID, model.RoleCS, "refund_success",
		return logOp(tx, &refund.OrderNo, &csID, model.RoleCS, "refund_success",
// 赋值/更新："渠道退款成功；库存恢复="+boolLabel(restore))
			"渠道退款成功；库存恢复="+boolLabel(restore))
// 结束当前代码块
	})
// 条件判断：if err != nil {
	if err != nil {
// 返回：return model.Refund{}, err
		return model.Refund{}, err
// 结束当前代码块
	}
// 条件判断：if err := s.deps.DB.Where("refund_no = ?", refund.RefundNo).First(&result).Error; err != nil {
	if err := s.deps.DB.Where("refund_no = ?", refund.RefundNo).First(&result).Error; err != nil {
// 返回：return model.Refund{}, err
		return model.Refund{}, err
// 结束当前代码块
	}
// 条件判断：if !ok {
	if !ok {
// 返回：return result, ErrChannelRefund
		return result, ErrChannelRefund
// 结束当前代码块
	}
// 返回：return result, nil
	return result, nil
// 结束当前代码块
}

// Get 按单号查询退款单

// 函数定义：func (s RefundService) Get(orderNo string) (model.Refund, error) {
func (s RefundService) Get(orderNo string) (model.Refund, error) {
// 声明：var r model.Refund
	var r model.Refund
// 条件判断：if err := s.deps.DB.Where("order_no = ?", strings.TrimSpace(orderNo)).First(&r).Error; err != nil {
	if err := s.deps.DB.Where("order_no = ?", strings.TrimSpace(orderNo)).First(&r).Error; err != nil {
// 返回：return r, ErrRefundNotFound
		return r, ErrRefundNotFound
// 结束当前代码块
	}
// 返回：return r, nil
	return r, nil
// 结束当前代码块
}

// List 客服工作台的退款单列表

// 函数定义：func (s RefundService) List(status string) ([]model.Refund, error) {
func (s RefundService) List(status string) ([]model.Refund, error) {
// 声明：var list []model.Refund
	var list []model.Refund
// 赋值/更新：q := s.deps.DB
	q := s.deps.DB
// 条件判断：if status != "" {
	if status != "" {
// 赋值/更新：q = q.Where("audit_status = ? OR channel_status = ?", status, status)
		q = q.Where("audit_status = ? OR channel_status = ?", status, status)
// 结束当前代码块
	}
// 赋值/更新：err := q.Order("id desc").Limit(100).Find(&list).Error
	err := q.Order("id desc").Limit(100).Find(&list).Error
// 返回：return list, err
	return list, err
// 结束当前代码块
}

// uniqueRefundNo 生成退款单号：事务内先查重，冲突换时间片重试，兜底用纳秒偏移。

// 函数定义：func (s RefundService) uniqueRefundNo(tx *gorm.DB, now time.Time) string {
func (s RefundService) uniqueRefundNo(tx *gorm.DB, now time.Time) string {
// 循环处理：for i := 0; i < 5; i++ {
	for i := 0; i < 5; i++ {
// 赋值/更新：no := utils.NewRefundNo(now)
		no := utils.NewRefundNo(now)
// 声明：var count int64
		var count int64
// 执行：tx.Model(&model.Refund{}).Where("refund_no = ?", no).Count(&count)
		tx.Model(&model.Refund{}).Where("refund_no = ?", no).Count(&count)
// 条件判断：if count == 0 {
		if count == 0 {
// 返回：return no
			return no
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return utils.NewRefundNo(now.Add(time.Nanosecond))
	return utils.NewRefundNo(now.Add(time.Nanosecond))
// 结束当前代码块
}

// boolLabel 布尔转中文，用于审计描述中的"是否恢复库存"等字段。

// 函数定义：func boolLabel(v bool) string {
func boolLabel(v bool) string {
// 条件判断：if v {
	if v {
// 返回：return "是"
		return "是"
// 结束当前代码块
	}
// 返回：return "否"
	return "否"
// 结束当前代码块
}
