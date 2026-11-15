// Package request 存放 API 层入参结构体（类 GVA model/request 分层）。
package request

// Login 登录入参

// 类型定义：type Login struct {
type Login struct {
// 执行：Username string `json:"username" binding:"required"`
	Username string `json:"username" binding:"required"`
// 执行：Password string `json:"password" binding:"required"`
	Password string `json:"password" binding:"required"`
// 结束当前代码块
}

// CreateOrder 买家下单

// 类型定义：type CreateOrder struct {
type CreateOrder struct {
// 执行：ProductID uint64 `json:"productId" binding:"required"`
	ProductID uint64 `json:"productId" binding:"required"`
// 执行：Qty       int    `json:"qty" binding:"required"`
	Qty       int    `json:"qty" binding:"required"`
// 结束当前代码块
}

// PaymentNotify 模拟渠道支付通知

// 类型定义：type PaymentNotify struct {
type PaymentNotify struct {
// 执行：OrderNo      string `json:"orderNo" binding:"required"`
	OrderNo      string `json:"orderNo" binding:"required"`
// 执行：ChannelTxnID string `json:"channelTxnId" binding:"required"`
	ChannelTxnID string `json:"channelTxnId" binding:"required"`
// 执行：AmountCents  int64  `json:"amountCents" binding:"required"`
	AmountCents  int64  `json:"amountCents" binding:"required"`
// 执行：Sign         string `json:"sign" binding:"required"`
	Sign         string `json:"sign" binding:"required"`
// 结束当前代码块
}

// RefundApply 售前退款申请

// 类型定义：type RefundApply struct {
type RefundApply struct {
// 执行：OrderNo string `json:"orderNo" binding:"required"`
	OrderNo string `json:"orderNo" binding:"required"`
// 执行：Reason  string `json:"reason"`
	Reason  string `json:"reason"`
// 结束当前代码块
}

// RefundAudit 客服审核

// 类型定义：type RefundAudit struct {
type RefundAudit struct {
// 执行：OrderNo string `json:"orderNo" binding:"required"`
	OrderNo string `json:"orderNo" binding:"required"`
// 执行：Approve bool   `json:"approve"`
	Approve bool   `json:"approve"`
// 执行：Reason  string `json:"reason"`
	Reason  string `json:"reason"`
// 结束当前代码块
}

// RefundRetry 渠道失败后重试

// 类型定义：type RefundRetry struct {
type RefundRetry struct {
// 执行：RefundNo string `json:"refundNo" binding:"required"`
	RefundNo string `json:"refundNo" binding:"required"`
// 结束当前代码块
}

// Ship 仓库发货登记

// 类型定义：type Ship struct {
type Ship struct {
// 执行：OrderNo    string `json:"orderNo" binding:"required"`
	OrderNo    string `json:"orderNo" binding:"required"`
// 执行：Carrier    string `json:"carrier" binding:"required"`
	Carrier    string `json:"carrier" binding:"required"`
// 执行：TrackingNo string `json:"trackingNo" binding:"required"`
	TrackingNo string `json:"trackingNo" binding:"required"`
// 结束当前代码块
}

// AbnormalHandle 客服处理异常支付（发起退款）

// 类型定义：type AbnormalHandle struct {
type AbnormalHandle struct {
// 执行：ChannelTxnID string `json:"channelTxnId" binding:"required"`
	ChannelTxnID string `json:"channelTxnId" binding:"required"`
// 执行：Reason       string `json:"reason"`
	Reason       string `json:"reason"`
// 结束当前代码块
}

// ProductUpdate 运营维护商品

// 类型定义：type ProductUpdate struct {
type ProductUpdate struct {
// 执行：ProductID  uint64 `json:"productId" binding:"required"`
	ProductID  uint64 `json:"productId" binding:"required"`
// 执行：OnSale     *bool  `json:"onSale"`
	OnSale     *bool  `json:"onSale"`
// 执行：PriceCents *int64 `json:"priceCents"`
	PriceCents *int64 `json:"priceCents"`
// 结束当前代码块
}

// AdminProductUpdate GVA 管理端商品维护入参（按 URL 主键寻址）

// 类型定义：type AdminProductUpdate struct {
type AdminProductUpdate struct {
// 执行：OnSale     *bool  `json:"onSale"`
	OnSale     *bool  `json:"onSale"`
// 执行：PriceCents *int64 `json:"priceCents"`
	PriceCents *int64 `json:"priceCents"`
// 执行：Available  *int   `json:"available"`
	Available  *int   `json:"available"`
// 结束当前代码块
}

// ChannelMode 模拟渠道退款模式（演示用：首次失败可切换）

// 类型定义：type ChannelMode struct {
type ChannelMode struct {
// 执行：RefundMode string `json:"refundMode" binding:"required"` // success | fail
	RefundMode string `json:"refundMode" binding:"required"` // success | fail
// 结束当前代码块
}
