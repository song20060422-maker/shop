// Package service 承载 Shop 交易系统的全部业务规则（类 GVA service 分层）。
// 并发互斥统一由「带前态条件的单条 UPDATE + 校验受影响行数」保证，
// 所有状态跃迁与库存变动都在数据库事务内完成。
package service

// 导入依赖包
import (
// 执行："errors"
	"errors"
// 执行："time"
	"time"

// 执行："gorm.io/gorm"
	"gorm.io/gorm"

// 执行："shop/shop/config"
	"shop/shop/config"
// 结束当前代码块
)

// Deps 服务依赖（可注入测试库、可控时钟与假渠道）

// 类型定义：type Deps struct {
type Deps struct {
// 执行：DB      *gorm.DB
	DB      *gorm.DB
// 执行：Now     func() time.Time
	Now     func() time.Time
// 执行：Cfg     config.Config
	Cfg     config.Config
// 执行：Channel ChannelClient
	Channel ChannelClient
// 结束当前代码块
}

// ServiceGroup 业务服务集合（GVA GroupApp 约定）

// 类型定义：type ServiceGroup struct {
type ServiceGroup struct {
// 执行：AuthService
	AuthService
// 执行：CatalogService
	CatalogService
// 执行：OrderService
	OrderService
// 执行：PaymentService
	PaymentService
// 执行：RefundService
	RefundService
// 执行：ShipmentService
	ShipmentService
// 执行：ReportService
	ReportService
// 执行：DemoService
	DemoService
// 执行：Channel ChannelClient
	Channel ChannelClient
// 结束当前代码块
}

// NewChannelFromConfig 按 shop.pay-channel 构建渠道：alipay（配置完整时）或 nil（回落模拟渠道，默认）。
// 返回 error 仅在显式要求支付宝但配置非法时发生——启动即失败，不静默降级。

// 函数定义：func NewChannelFromConfig(cfg config.Config) (ChannelClient, error) {
func NewChannelFromConfig(cfg config.Config) (ChannelClient, error) {
// 条件判断：if cfg.Shop.PayChannel != "alipay" {
	if cfg.Shop.PayChannel != "alipay" {
// 返回：return nil, nil
		return nil, nil
// 结束当前代码块
	}
// 返回：return NewAlipayChannel(cfg.Shop.Alipay)
	return NewAlipayChannel(cfg.Shop.Alipay)
// 结束当前代码块
}

// NewServiceGroup 组装全部业务服务

// 函数定义：func NewServiceGroup(d Deps) *ServiceGroup {
func NewServiceGroup(d Deps) *ServiceGroup {
// 条件判断：if d.Now == nil {
	if d.Now == nil {
// 赋值/更新：d.Now = time.Now
		d.Now = time.Now
// 结束当前代码块
	}
// 条件判断：if d.Channel == nil {
	if d.Channel == nil {
// 赋值/更新：d.Channel = NewSimulatedChannel(d.Cfg.AllowRefundByChannel())
		d.Channel = NewSimulatedChannel(d.Cfg.AllowRefundByChannel())
// 结束当前代码块
	}
// 赋值/更新：g := &ServiceGroup{Channel: d.Channel}
	g := &ServiceGroup{Channel: d.Channel}
// 赋值/更新：g.AuthService = AuthService{deps: d}
	g.AuthService = AuthService{deps: d}
// 赋值/更新：g.CatalogService = CatalogService{deps: d}
	g.CatalogService = CatalogService{deps: d}
// 赋值/更新：g.OrderService = OrderService{deps: d}
	g.OrderService = OrderService{deps: d}
// 赋值/更新：g.PaymentService = PaymentService{deps: d}
	g.PaymentService = PaymentService{deps: d}
// 赋值/更新：g.RefundService = RefundService{deps: d}
	g.RefundService = RefundService{deps: d}
// 赋值/更新：g.ShipmentService = ShipmentService{deps: d}
	g.ShipmentService = ShipmentService{deps: d}
// 赋值/更新：g.ReportService = ReportService{deps: d}
	g.ReportService = ReportService{deps: d}
// 赋值/更新：g.DemoService = DemoService{deps: d}
	g.DemoService = DemoService{deps: d}
// 返回：return g
	return g
// 结束当前代码块
}

// BizError 业务错误：携带稳定错误码，API 层直接映射 HTTP 状态与提示。

// 类型定义：type BizError struct {
type BizError struct {
// 执行：Code int    `json:"code"`
	Code int    `json:"code"`
// 执行：Msg  string `json:"msg"`
	Msg  string `json:"msg"`
// 结束当前代码块
}

// Error 实现 error 接口，msg 即对外提示文案。

// 函数定义：func (e BizError) Error() string { return e.Msg }
func (e BizError) Error() string { return e.Msg }

// NewBizError 构造业务错误

// 函数定义：func NewBizError(code int, msg string) BizError { return BizError{Code: code, Msg: msg} }
func NewBizError(code int, msg string) BizError { return BizError{Code: code, Msg: msg} }

// 业务错误清单（对外可解释，禁止用「操作失败」这类含糊提示掩盖真实原因）

// 声明：var (
var (
// 赋值/更新：ErrInvalidParam      = NewBizError(1000, "参数不合法")
	ErrInvalidParam      = NewBizError(1000, "参数不合法")
// 赋值/更新：ErrProductNotFound   = NewBizError(1001, "商品不存在")
	ErrProductNotFound   = NewBizError(1001, "商品不存在")
// 赋值/更新：ErrProductOffSale    = NewBizError(1002, "商品已下架，无法购买")
	ErrProductOffSale    = NewBizError(1002, "商品已下架，无法购买")
// 赋值/更新：ErrInsufficientStock = NewBizError(1003, "可售数量不足")
	ErrInsufficientStock = NewBizError(1003, "可售数量不足")
// 赋值/更新：ErrQtyInvalid        = NewBizError(1004, "购买数量只能为 1～2 件")
	ErrQtyInvalid        = NewBizError(1004, "购买数量只能为 1～2 件")
// 赋值/更新：ErrReasonRequired    = NewBizError(1005, "客服处理必须填写原因")
	ErrReasonRequired    = NewBizError(1005, "客服处理必须填写原因")

// 赋值/更新：ErrOrderNotFound   = NewBizError(2001, "订单不存在")
	ErrOrderNotFound   = NewBizError(2001, "订单不存在")
// 赋值/更新：ErrForbidden       = NewBizError(2002, "无权访问他人订单")
	ErrForbidden       = NewBizError(2002, "无权访问他人订单")
// 赋值/更新：ErrOrderNotPayable = NewBizError(2003, "订单当前状态不允许该操作")
	ErrOrderNotPayable = NewBizError(2003, "订单当前状态不允许该操作")
// 赋值/更新：ErrOrderClosed     = NewBizError(2004, "订单已关闭，不能恢复")
	ErrOrderClosed     = NewBizError(2004, "订单已关闭，不能恢复")
// 赋值/更新：ErrAlreadyPaid     = NewBizError(2005, "订单已支付，请勿重复支付")
	ErrAlreadyPaid     = NewBizError(2005, "订单已支付，请勿重复支付")
// 赋值/更新：ErrNotPaid         = NewBizError(2006, "订单未支付，不能发货")
	ErrNotPaid         = NewBizError(2006, "订单未支付，不能发货")
// 赋值/更新：ErrAlreadyShipped  = NewBizError(2007, "订单已发货，不能重复登记")
	ErrAlreadyShipped  = NewBizError(2007, "订单已发货，不能重复登记")
// 赋值/更新：ErrNotShippable    = NewBizError(2008, "订单当前状态不可发货")
	ErrNotShippable    = NewBizError(2008, "订单当前状态不可发货")

// 赋值/更新：ErrBadSign          = NewBizError(3001, "通知验签失败")
	ErrBadSign          = NewBizError(3001, "通知验签失败")
// 赋值/更新：ErrAmountMismatch   = NewBizError(3002, "通知金额与订单金额不一致")
	ErrAmountMismatch   = NewBizError(3002, "通知金额与订单金额不一致")
// 赋值/更新：ErrDuplicatedPay    = NewBizError(3003, "重复支付通知，未重复记款")
	ErrDuplicatedPay    = NewBizError(3003, "重复支付通知，未重复记款")
// 赋值/更新：ErrChannelRefund    = NewBizError(3004, "渠道退款失败，已保留待处理记录")
	ErrChannelRefund    = NewBizError(3004, "渠道退款失败，已保留待处理记录")
// 赋值/更新：ErrAbnormalNotFound = NewBizError(3005, "异常收款记录不存在或已处理")
	ErrAbnormalNotFound = NewBizError(3005, "异常收款记录不存在或已处理")
// 赋值/更新：ErrAlipayDisabled   = NewBizError(3007, "支付宝渠道未启用")
	ErrAlipayDisabled   = NewBizError(3007, "支付宝渠道未启用")
// 赋值/更新：ErrTradeNotPaid     = NewBizError(3008, "交易未完成，拒绝记款")
	ErrTradeNotPaid     = NewBizError(3008, "交易未完成，拒绝记款")

// 赋值/更新：ErrRefundNotFound     = NewBizError(4001, "退款单不存在")
	ErrRefundNotFound     = NewBizError(4001, "退款单不存在")
// 赋值/更新：ErrRefundExists       = NewBizError(4002, "该订单已退款或已在退款流程中")
	ErrRefundExists       = NewBizError(4002, "该订单已退款或已在退款流程中")
// 赋值/更新：ErrRefundAudited      = NewBizError(4003, "退款单已审核，不能重复审核")
	ErrRefundAudited      = NewBizError(4003, "退款单已审核，不能重复审核")
// 赋值/更新：ErrRefundInProgress   = NewBizError(4004, "退款处理中，暂停发货")
	ErrRefundInProgress   = NewBizError(4004, "退款处理中，暂停发货")
// 赋值/更新：ErrRefunded           = NewBizError(4005, "订单已退款，不能发货")
	ErrRefunded           = NewBizError(4005, "订单已退款，不能发货")
// 赋值/更新：ErrRefundNotRetryable = NewBizError(4006, "退款单当前状态不可重试")
	ErrRefundNotRetryable = NewBizError(4006, "退款单当前状态不可重试")

// 赋值/更新：ErrLoginFailed = NewBizError(5001, "账号或密码错误")
	ErrLoginFailed = NewBizError(5001, "账号或密码错误")
// 结束当前代码块
)

// IsBiz 判断是否为业务错误（其余视为系统错误）

// 函数定义：func IsBiz(err error) bool {
func IsBiz(err error) bool {
// 声明：var be BizError
	var be BizError
// 返回：return errors.As(err, &be)
	return errors.As(err, &be)
// 结束当前代码块
}

// BizCode 提取业务错误码

// 函数定义：func BizCode(err error) (int, string) {
func BizCode(err error) (int, string) {
// 声明：var be BizError
	var be BizError
// 条件判断：if errors.As(err, &be) {
	if errors.As(err, &be) {
// 返回：return be.Code, be.Msg
		return be.Code, be.Msg
// 结束当前代码块
	}
// 返回：return 9000, "系统繁忙，请稍后重试"
	return 9000, "系统繁忙，请稍后重试"
// 结束当前代码块
}

// BizCodeMsg 仅取提示文案（接口层回传业务失败时使用）

// 函数定义：func BizCodeMsg(err error) string {
func BizCodeMsg(err error) string {
// 赋值/更新：_, msg := BizCode(err)
	_, msg := BizCode(err)
// 返回：return msg
	return msg
// 结束当前代码块
}
