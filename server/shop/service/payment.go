// 包声明
package service

// 导入依赖包
import (
// 执行："context"
	"context"
// 执行："errors"
	"errors"
// 执行："net/url"
	"net/url"
// 执行："strings"
	"strings"
// 执行："time"
	"time"

// 执行："github.com/go-sql-driver/mysql"
	"github.com/go-sql-driver/mysql"
// 执行："github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/alipay/v3"
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

// isDuplicateEntry 判定是否 MySQL 唯一键冲突（errno 1062）

// 函数定义：func isDuplicateEntry(err error) bool {
func isDuplicateEntry(err error) bool {
// 声明：var me *mysql.MySQLError
	var me *mysql.MySQLError
// 返回：return errors.As(err, &me) && me.Number == 1062
	return errors.As(err, &me) && me.Number == 1062
// 结束当前代码块
}

// errDuplicatePayment 内部哨兵：触发事务回滚，外层据此判定重复通知。

// 声明：var errDuplicatePayment = errors.New("duplicate payment")
var errDuplicatePayment = errors.New("duplicate payment")

// NotifyResult 支付通知处理结果（渠道与运营都可查询）

// 类型定义：type NotifyResult struct {
type NotifyResult struct {
// 执行：Result       string `json:"result"` // accepted | rejected | duplicated | abnormal
	Result       string `json:"result"` // accepted | rejected | duplicated | abnormal
// 执行：OrderNo      string `json:"orderNo"`
	OrderNo      string `json:"orderNo"`
// 执行：ChannelTxnID string `json:"channelTxnId"`
	ChannelTxnID string `json:"channelTxnId"`
// 执行：AmountCents  int64  `json:"amountCents"`
	AmountCents  int64  `json:"amountCents"`
// 执行：Reason       string `json:"reason"`
	Reason       string `json:"reason"`
// 执行：OrderStatus  string `json:"orderStatus"`
	OrderStatus  string `json:"orderStatus"`
// 执行：NoticeID     uint64 `json:"noticeId"`
	NoticeID     uint64 `json:"noticeId"`
// 执行：AbnormalID   uint64 `json:"abnormalId,omitempty"`
	AbnormalID   uint64 `json:"abnormalId,omitempty"`
// 结束当前代码块
}

// PaymentService 支付通知受理与异常收款处理

// 类型定义：type PaymentService struct {
type PaymentService struct {
// 执行：deps Deps
	deps Deps
// 结束当前代码块
}

// Notify 受理模拟渠道支付通知：HMAC 验签 → 复用 applyPayment 推进。

// 函数定义：func (s PaymentService) Notify(orderNo, txn string, amount int64, sign, rawBody string) (NotifyResult, error) ...
func (s PaymentService) Notify(orderNo, txn string, amount int64, sign, rawBody string) (NotifyResult, error) {
// 赋值/更新：now := s.deps.Now()
	now := s.deps.Now()
// 赋值/更新：orderNo = strings.TrimSpace(orderNo)
	orderNo = strings.TrimSpace(orderNo)
// 赋值/更新：txn = strings.TrimSpace(txn)
	txn = strings.TrimSpace(txn)

	// 1. 参数与验签：验签失败必须留痕，且绝不记款
	if orderNo == "" || txn == "" || amount <= 0 || sign == "" {
// 赋值/更新：result := NotifyResult{Result: model.NoticeRejected, OrderNo: orderNo, ChannelTxnID: txn, AmountCents: amount,...
		result := NotifyResult{Result: model.NoticeRejected, OrderNo: orderNo, ChannelTxnID: txn, AmountCents: amount, Reason: "参数不合法"}
// 执行：s.saveNotice(result, false, rawBody, now)
		s.saveNotice(result, false, rawBody, now)
// 返回：return result, ErrInvalidParam
		return result, ErrInvalidParam
// 结束当前代码块
	}
// 条件判断：if utils.PaySign(s.deps.Cfg.Shop.PayNotifySecret, orderNo, txn, amount) != sign {
	if utils.PaySign(s.deps.Cfg.Shop.PayNotifySecret, orderNo, txn, amount) != sign {
// 赋值/更新：result := NotifyResult{Result: model.NoticeRejected, OrderNo: orderNo, ChannelTxnID: txn, AmountCents: amount,...
		result := NotifyResult{Result: model.NoticeRejected, OrderNo: orderNo, ChannelTxnID: txn, AmountCents: amount, Reason: "验签失败"}
// 执行：s.saveNotice(result, false, rawBody, now)
		s.saveNotice(result, false, rawBody, now)
// 返回：return result, ErrBadSign
		return result, ErrBadSign
// 结束当前代码块
	}
// 返回：return s.applyPayment(orderNo, txn, amount, rawBody, now)
	return s.applyPayment(orderNo, txn, amount, rawBody, now)
// 结束当前代码块
}

// NotifyAlipay 支付宝异步通知：RSA2 验签（渠道内实现）→ app_id/交易状态校验 → 复用同一套推进逻辑，
// 保证真实渠道与模拟渠道在订单状态机、幂等与异常收款上行为完全一致。

// 函数定义：func (s PaymentService) NotifyAlipay(form url.Values, rawBody string) (NotifyResult, error) {
func (s PaymentService) NotifyAlipay(form url.Values, rawBody string) (NotifyResult, error) {
// 赋值/更新：now := s.deps.Now()
	now := s.deps.Now()
// 赋值/更新：result := NotifyResult{Result: model.NoticeRejected,
	result := NotifyResult{Result: model.NoticeRejected,
// 执行：OrderNo: form.Get("out_trade_no"), ChannelTxnID: form.Get("trade_no")}
		OrderNo: form.Get("out_trade_no"), ChannelTxnID: form.Get("trade_no")}
// 赋值/更新：ch, isAli := s.deps.Channel.(NotifyVerifier)
	ch, isAli := s.deps.Channel.(NotifyVerifier)
// 条件判断：if !isAli {
	if !isAli {
// 赋值/更新：result.Reason = "支付宝渠道未启用"
		result.Reason = "支付宝渠道未启用"
// 执行：s.saveNotice(result, false, rawBody, now)
		s.saveNotice(result, false, rawBody, now)
// 返回：return result, ErrAlipayDisabled
		return result, ErrAlipayDisabled
// 结束当前代码块
	}
// 赋值/更新：payload, err := ch.VerifyNotify(context.Background(), form)
	payload, err := ch.VerifyNotify(context.Background(), form)
// 赋值/更新：result.OrderNo = payload.OrderNo
	result.OrderNo = payload.OrderNo
// 赋值/更新：result.ChannelTxnID = payload.Txn
	result.ChannelTxnID = payload.Txn
// 赋值/更新：result.AmountCents = payload.AmountCents
	result.AmountCents = payload.AmountCents
// 条件判断：if err != nil {
	if err != nil {
// 赋值/更新：result.Reason = "验签失败：" + err.Error()
		result.Reason = "验签失败：" + err.Error()
// 执行：s.saveNotice(result, false, rawBody, now)
		s.saveNotice(result, false, rawBody, now)
// 返回：return result, ErrBadSign
		return result, ErrBadSign
// 结束当前代码块
	}
// 条件判断：if alipay.TradeStatus(payload.TradeStatus) != alipay.TradeStatusSuccess && alipay.TradeStatus(payload.TradeSta...
	if alipay.TradeStatus(payload.TradeStatus) != alipay.TradeStatusSuccess && alipay.TradeStatus(payload.TradeStatus) != alipay.TradeStatusFinished {
// 赋值/更新：result.Reason = "交易未完成：" + payload.TradeStatus
		result.Reason = "交易未完成：" + payload.TradeStatus
// 执行：s.saveNotice(result, false, rawBody, now)
		s.saveNotice(result, false, rawBody, now)
// 返回：return result, ErrTradeNotPaid
		return result, ErrTradeNotPaid
// 结束当前代码块
	}
// 返回：return s.applyPayment(payload.OrderNo, payload.Txn, payload.AmountCents, rawBody, now)
	return s.applyPayment(payload.OrderNo, payload.Txn, payload.AmountCents, rawBody, now)
// 结束当前代码块
}

// applyPayment 验签通过后的共同路径：匹配订单与金额 → 幂等记款 / 迟到支付转异常。

// 函数定义：func (s PaymentService) applyPayment(orderNo, txn string, amount int64, rawBody string, now time.Time) (Notify...
func (s PaymentService) applyPayment(orderNo, txn string, amount int64, rawBody string, now time.Time) (NotifyResult, error) {
// 赋值/更新：result := NotifyResult{OrderNo: orderNo, ChannelTxnID: txn, AmountCents: amount}
	result := NotifyResult{OrderNo: orderNo, ChannelTxnID: txn, AmountCents: amount}

// 赋值/更新：err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
	err := s.deps.DB.Transaction(func(tx *gorm.DB) error {
// 声明：var o model.Order
		var o model.Order
// 条件判断：if err := tx.Where("order_no = ?", orderNo).First(&o).Error; err != nil {
		if err := tx.Where("order_no = ?", orderNo).First(&o).Error; err != nil {
// 赋值/更新：result.Result = model.NoticeRejected
			result.Result = model.NoticeRejected
// 赋值/更新：result.Reason = "未知订单"
			result.Reason = "未知订单"
// 返回：return s.noticeInTx(tx, &result, true, rawBody, now)
			return s.noticeInTx(tx, &result, true, rawBody, now)
// 结束当前代码块
		}
// 赋值/更新：result.OrderStatus = o.Status
		result.OrderStatus = o.Status
// 条件判断：if amount != o.AmountCents {
		if amount != o.AmountCents {
// 赋值/更新：result.Result = model.NoticeRejected
			result.Result = model.NoticeRejected
// 赋值/更新：result.Reason = "金额不符：订单应收 " + itoa64(o.AmountCents) + " 分，通知为 " + itoa64(amount) + " 分"
			result.Reason = "金额不符：订单应收 " + itoa64(o.AmountCents) + " 分，通知为 " + itoa64(amount) + " 分"
// 返回：return s.noticeInTx(tx, &result, true, rawBody, now)
			return s.noticeInTx(tx, &result, true, rawBody, now)
// 结束当前代码块
		}

		// 2. 订单已不在「待付款」：重复通知或迟到支付
		if o.Status != model.OrderStatusPendingPayment {
// 条件判断：if isClosedStatus(o.Status) {
			if isClosedStatus(o.Status) {
// 返回：return s.markAbnormal(tx, &result, amount, now, rawBody, "订单超时关闭后收到支付成功通知，需客服核实退款")
				return s.markAbnormal(tx, &result, amount, now, rawBody, "订单超时关闭后收到支付成功通知，需客服核实退款")
// 结束当前代码块
			}
// 赋值/更新：result.Result = model.NoticeDuplicated
			result.Result = model.NoticeDuplicated
// 赋值/更新：result.Reason = "订单已" + statusLabel(o.Status) + "，通知不重复记款"
			result.Reason = "订单已" + statusLabel(o.Status) + "，通知不重复记款"
// 返回：return s.noticeInTx(tx, &result, true, rawBody, now)
			return s.noticeInTx(tx, &result, true, rawBody, now)
// 结束当前代码块
		}

		// 3. 有效期内：条件推进订单状态，再记款；唯一键冲突即整体回滚
		res := tx.Model(&model.Order{}).
// 执行：Where("order_no = ? AND status = ? AND expire_at > ?", orderNo, model.OrderStatusPendingPayment, now).
			Where("order_no = ? AND status = ? AND expire_at > ?", orderNo, model.OrderStatusPendingPayment, now).
// 执行：Updates(map[string]interface{}{"status": model.OrderStatusPaid, "paid_at": now, "updated_at": now})
			Updates(map[string]interface{}{"status": model.OrderStatusPaid, "paid_at": now, "updated_at": now})
// 条件判断：if res.Error != nil {
		if res.Error != nil {
// 返回：return res.Error
			return res.Error
// 结束当前代码块
		}
// 条件判断：if res.RowsAffected != 1 {
		if res.RowsAffected != 1 {
// 声明：var latest model.Order
			var latest model.Order
// 条件判断：if err := tx.Where("order_no = ?", orderNo).First(&latest).Error; err != nil {
			if err := tx.Where("order_no = ?", orderNo).First(&latest).Error; err != nil {
// 返回：return err
				return err
// 结束当前代码块
			}
// 条件判断：if isClosedStatus(latest.Status) {
			if isClosedStatus(latest.Status) {
// 返回：return s.markAbnormal(tx, &result, amount, now, rawBody, "通知到达时订单已关闭，转入异常收款待处理")
				return s.markAbnormal(tx, &result, amount, now, rawBody, "通知到达时订单已关闭，转入异常收款待处理")
// 结束当前代码块
			}
// 赋值/更新：result.Result = model.NoticeDuplicated
			result.Result = model.NoticeDuplicated
// 赋值/更新：result.Reason = "订单已" + statusLabel(latest.Status) + "，通知不重复记款"
			result.Reason = "订单已" + statusLabel(latest.Status) + "，通知不重复记款"
// 返回：return s.noticeInTx(tx, &result, true, rawBody, now)
			return s.noticeInTx(tx, &result, true, rawBody, now)
// 结束当前代码块
		}
// 赋值/更新：pay := model.Payment{OrderNo: orderNo, ChannelTxnID: txn, AmountCents: amount, PaidAt: now, CreatedAt: now}
		pay := model.Payment{OrderNo: orderNo, ChannelTxnID: txn, AmountCents: amount, PaidAt: now, CreatedAt: now}
// 条件判断：if err := tx.Create(&pay).Error; err != nil {
		if err := tx.Create(&pay).Error; err != nil {
			// I-18 错误分类：只有唯一键冲突（重复交易号/一单多记款）才判重复支付；
			// 其它 DB 错误按真实错误上抛，不再一律误报为「重复支付」。
			if isDuplicateEntry(err) {
// 返回：return errDuplicatePayment
				return errDuplicatePayment
// 结束当前代码块
			}
// 返回：return err
			return err
// 结束当前代码块
		}
// 赋值/更新：result.Result = model.NoticeAccepted
		result.Result = model.NoticeAccepted
// 赋值/更新：result.OrderStatus = model.OrderStatusPaid
		result.OrderStatus = model.OrderStatusPaid
// 赋值/更新：result.Reason = "记款成功，订单进入待发货"
		result.Reason = "记款成功，订单进入待发货"
// 条件判断：if err := s.noticeInTx(tx, &result, true, rawBody, now); err != nil {
		if err := s.noticeInTx(tx, &result, true, rawBody, now); err != nil {
// 返回：return err
			return err
// 结束当前代码块
		}
// 返回：return logOp(tx, &orderNo, nil, "system", "payment_accepted",
		return logOp(tx, &orderNo, nil, "system", "payment_accepted",
// 执行："渠道记款："+txn+"，金额 "+itoa64(amount)+" 分")
			"渠道记款："+txn+"，金额 "+itoa64(amount)+" 分")
// 结束当前代码块
	})
// 条件判断：if err != nil && errors.Is(err, errDuplicatePayment) {
	if err != nil && errors.Is(err, errDuplicatePayment) {
// 赋值/更新：dup := NotifyResult{Result: model.NoticeDuplicated, OrderNo: orderNo, ChannelTxnID: txn,
		dup := NotifyResult{Result: model.NoticeDuplicated, OrderNo: orderNo, ChannelTxnID: txn,
// 执行：AmountCents: amount, Reason: "重复支付通知，未重复记款", OrderStatus: model.OrderStatusPaid}
			AmountCents: amount, Reason: "重复支付通知，未重复记款", OrderStatus: model.OrderStatusPaid}
// 执行：s.saveNotice(dup, true, rawBody, now)
		s.saveNotice(dup, true, rawBody, now)
// 返回：return dup, ErrDuplicatedPay
		return dup, ErrDuplicatedPay
// 结束当前代码块
	}
// 条件判断：if err != nil {
	if err != nil {
// 返回：return result, err
		return result, err
// 结束当前代码块
	}
// 条件判断：if result.Result == model.NoticeDuplicated {
	if result.Result == model.NoticeDuplicated {
// 返回：return result, ErrDuplicatedPay
		return result, ErrDuplicatedPay
// 结束当前代码块
	}
// 返回：return result, nil
	return result, nil
// 结束当前代码块
}

// AlipayPayURL 为待付款订单生成支付宝收银台跳转地址（电脑网站支付）：
// 校验本人订单与待付款状态，渠道未启用返回业务错误。

// 函数定义：func (s PaymentService) AlipayPayURL(orderNo string, buyerID uint64, publicBase string) (string, error) {
func (s PaymentService) AlipayPayURL(orderNo string, buyerID uint64, publicBase string) (string, error) {
// 声明：var o model.Order
	var o model.Order
// 条件判断：if err := s.deps.DB.Where("order_no = ?", strings.TrimSpace(orderNo)).First(&o).Error; err != nil {
	if err := s.deps.DB.Where("order_no = ?", strings.TrimSpace(orderNo)).First(&o).Error; err != nil {
// 返回：return "", ErrOrderNotFound
		return "", ErrOrderNotFound
// 结束当前代码块
	}
// 条件判断：if o.BuyerID != buyerID {
	if o.BuyerID != buyerID {
// 返回：return "", ErrForbidden
		return "", ErrForbidden
// 结束当前代码块
	}
// 条件判断：if o.Status != model.OrderStatusPendingPayment {
	if o.Status != model.OrderStatusPendingPayment {
// 返回：return "", ErrOrderNotPayable
		return "", ErrOrderNotPayable
// 结束当前代码块
	}
// 赋值/更新：ac, isAli := s.deps.Channel.(*AlipayChannel)
	ac, isAli := s.deps.Channel.(*AlipayChannel)
// 条件判断：if !isAli {
	if !isAli {
// 返回：return "", ErrAlipayDisabled
		return "", ErrAlipayDisabled
// 结束当前代码块
	}
	// 回调地址不写死：买家从哪个域名访问，通知就回哪个域名（config 显式配置仍优先生效）
	notifyURL, returnURL := "", ""
// 条件判断：if b := strings.TrimSpace(publicBase); b != "" {
	if b := strings.TrimSpace(publicBase); b != "" {
// 赋值/更新：notifyURL = b + "/api/pay/alipay/notify"
		notifyURL = b + "/api/pay/alipay/notify"
// 返回：returnURL = b + "/shop"
		returnURL = b + "/shop"
// 结束当前代码块
	}
// 返回：return ac.PagePayURL(o.OrderNo, o.AmountCents, "商城订单："+o.ProductName, notifyURL, returnURL)
	return ac.PagePayURL(o.OrderNo, o.AmountCents, "商城订单："+o.ProductName, notifyURL, returnURL)
// 结束当前代码块
}

// markAbnormal 迟到支付：登记异常收款，绝不恢复订单、绝不重新占用库存。

// 函数定义：func (s PaymentService) markAbnormal(tx *gorm.DB, r *NotifyResult, amount int64, now time.Time, rawBody, reaso...
func (s PaymentService) markAbnormal(tx *gorm.DB, r *NotifyResult, amount int64, now time.Time, rawBody, reason string) error {
// 赋值/更新：ab := model.AbnormalPayment{
	ab := model.AbnormalPayment{
// 执行：OrderNo:      r.OrderNo,
		OrderNo:      r.OrderNo,
// 执行：ChannelTxnID: r.ChannelTxnID,
		ChannelTxnID: r.ChannelTxnID,
// 执行：AmountCents:  amount,
		AmountCents:  amount,
// 执行：Reason:       reason,
		Reason:       reason,
// 执行：Status:       model.AbnormalPending,
		Status:       model.AbnormalPending,
// 执行：CreatedAt:    now,
		CreatedAt:    now,
// 结束当前代码块
	}
// 条件判断：if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&ab).Error; err != nil {
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&ab).Error; err != nil {
// 返回：return err
		return err
// 结束当前代码块
	}
// 声明：var saved model.AbnormalPayment
	var saved model.AbnormalPayment
// 条件判断：if err := tx.Where("channel_txn_id = ?", r.ChannelTxnID).First(&saved).Error; err == nil {
	if err := tx.Where("channel_txn_id = ?", r.ChannelTxnID).First(&saved).Error; err == nil {
// 赋值/更新：r.AbnormalID = saved.ID
		r.AbnormalID = saved.ID
// 结束当前代码块
	}
// 赋值/更新：r.Result = model.NoticeAbnormal
	r.Result = model.NoticeAbnormal
// 赋值/更新：r.Reason = reason
	r.Reason = reason
// 条件判断：if err := s.noticeInTx(tx, r, true, rawBody, now); err != nil {
	if err := s.noticeInTx(tx, r, true, rawBody, now); err != nil {
// 返回：return err
		return err
// 结束当前代码块
	}
// 返回：return logOp(tx, &r.OrderNo, nil, "system", "abnormal_payment", "迟到支付登记："+r.ChannelTxnID)
	return logOp(tx, &r.OrderNo, nil, "system", "abnormal_payment", "迟到支付登记："+r.ChannelTxnID)
// 结束当前代码块
}

// noticeInTx 事务内写通知留痕

// 函数定义：func (s PaymentService) noticeInTx(tx *gorm.DB, r *NotifyResult, signOK bool, rawBody string, now time.Time) e...
func (s PaymentService) noticeInTx(tx *gorm.DB, r *NotifyResult, signOK bool, rawBody string, now time.Time) error {
// 赋值/更新：body := rawBody
	body := rawBody
// 赋值/更新：notice := model.PaymentNotice{
	notice := model.PaymentNotice{
// 执行：OrderNo:      r.OrderNo,
		OrderNo:      r.OrderNo,
// 执行：ChannelTxnID: r.ChannelTxnID,
		ChannelTxnID: r.ChannelTxnID,
// 执行：AmountCents:  r.AmountCents,
		AmountCents:  r.AmountCents,
// 执行：SignOK:       signOK,
		SignOK:       signOK,
// 执行：Result:       r.Result,
		Result:       r.Result,
// 执行：Reason:       r.Reason,
		Reason:       r.Reason,
// 执行：RawBody:      &body,
		RawBody:      &body,
// 执行：CreatedAt:    now,
		CreatedAt:    now,
// 结束当前代码块
	}
// 条件判断：if err := tx.Create(&notice).Error; err != nil {
	if err := tx.Create(&notice).Error; err != nil {
// 返回：return err
		return err
// 结束当前代码块
	}
// 赋值/更新：r.NoticeID = notice.ID
	r.NoticeID = notice.ID
// 返回：return nil
	return nil
// 结束当前代码块
}

// saveNotice 独立写通知留痕（验签失败等不进入主事务的场景）

// 函数定义：func (s PaymentService) saveNotice(r NotifyResult, signOK bool, rawBody string, now time.Time) {
func (s PaymentService) saveNotice(r NotifyResult, signOK bool, rawBody string, now time.Time) {
// 赋值/更新：body := rawBody
	body := rawBody
// 赋值/更新：_ = s.deps.DB.Create(&model.PaymentNotice{
	_ = s.deps.DB.Create(&model.PaymentNotice{
// 执行：OrderNo:      r.OrderNo,
		OrderNo:      r.OrderNo,
// 执行：ChannelTxnID: r.ChannelTxnID,
		ChannelTxnID: r.ChannelTxnID,
// 执行：AmountCents:  r.AmountCents,
		AmountCents:  r.AmountCents,
// 执行：SignOK:       signOK,
		SignOK:       signOK,
// 执行：Result:       r.Result,
		Result:       r.Result,
// 执行：Reason:       r.Reason,
		Reason:       r.Reason,
// 执行：RawBody:      &body,
		RawBody:      &body,
// 执行：CreatedAt:    now,
		CreatedAt:    now,
// 结束当前代码块
	}).Error
// 结束当前代码块
}

// Notices 按订单查询通知处理结果（需求：金额不符/未知订单/验签失败须可查询）

// 函数定义：func (s PaymentService) Notices(orderNo string) ([]model.PaymentNotice, error) {
func (s PaymentService) Notices(orderNo string) ([]model.PaymentNotice, error) {
// 声明：var list []model.PaymentNotice
	var list []model.PaymentNotice
// 赋值/更新：err := s.deps.DB.Where("order_no = ?", strings.TrimSpace(orderNo)).Order("id asc").Find(&list).Error
	err := s.deps.DB.Where("order_no = ?", strings.TrimSpace(orderNo)).Order("id asc").Find(&list).Error
// 返回：return list, err
	return list, err
// 结束当前代码块
}

// NoticesByResult 运营按处理结果过滤（可为空）

// 函数定义：func (s PaymentService) NoticesByResult(result string, limit int) ([]model.PaymentNotice, error) {
func (s PaymentService) NoticesByResult(result string, limit int) ([]model.PaymentNotice, error) {
// 声明：var list []model.PaymentNotice
	var list []model.PaymentNotice
// 赋值/更新：q := s.deps.DB
	q := s.deps.DB
// 条件判断：if result != "" {
	if result != "" {
// 赋值/更新：q = q.Where("result = ?", result)
		q = q.Where("result = ?", result)
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

// AbnormalList 异常收款待处理列表

// 函数定义：func (s PaymentService) AbnormalList(status string) ([]model.AbnormalPayment, error) {
func (s PaymentService) AbnormalList(status string) ([]model.AbnormalPayment, error) {
// 声明：var list []model.AbnormalPayment
	var list []model.AbnormalPayment
// 赋值/更新：q := s.deps.DB
	q := s.deps.DB
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

// isClosedStatus 判断订单是否已终结（关闭或取消）：迟到支付据此转异常收款而不是恢复订单。

// 函数定义：func isClosedStatus(status string) bool {
func isClosedStatus(status string) bool {
// 返回：return status == model.OrderStatusClosed || status == model.OrderStatusCancelled
	return status == model.OrderStatusClosed || status == model.OrderStatusCancelled
// 结束当前代码块
}

// statusLabel 订单状态中文名，用于通知留痕与审计描述。

// 函数定义：func statusLabel(status string) string {
func statusLabel(status string) string {
// 多路分支：switch status {
	switch status {
// 匹配分支：case model.OrderStatusPendingPayment:
	case model.OrderStatusPendingPayment:
// 返回：return "待付款"
		return "待付款"
// 匹配分支：case model.OrderStatusCancelled:
	case model.OrderStatusCancelled:
// 返回：return "取消"
		return "取消"
// 匹配分支：case model.OrderStatusClosed:
	case model.OrderStatusClosed:
// 返回：return "关闭"
		return "关闭"
// 匹配分支：case model.OrderStatusPaid:
	case model.OrderStatusPaid:
// 返回：return "付款"
		return "付款"
// 匹配分支：case model.OrderStatusShipped:
	case model.OrderStatusShipped:
// 返回：return "发货"
		return "发货"
// 匹配分支：case model.OrderStatusRefundRequested:
	case model.OrderStatusRefundRequested:
// 返回：return "退款申请"
		return "退款申请"
// 匹配分支：case model.OrderStatusRefundFailed:
	case model.OrderStatusRefundFailed:
// 返回：return "退款失败"
		return "退款失败"
// 匹配分支：case model.OrderStatusRefunded:
	case model.OrderStatusRefunded:
// 返回：return "退款"
		return "退款"
// 默认分支
	default:
// 返回：return status
		return status
// 结束当前代码块
	}
// 结束当前代码块
}

// itoa64 int64 版本的单号拼接辅助。

// 函数定义：func itoa64(v int64) string { return itoa(int(v)) }
func itoa64(v int64) string { return itoa(int(v)) }
