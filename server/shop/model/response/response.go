// Package response 存放统一响应体与对外 DTO。
package response

// 导入依赖包
import "shop/shop/model"

// Body 统一响应体（GVA 约定 code=0 成功）

// 类型定义：type Body struct {
type Body struct {
// 执行：Code int         `json:"code"`
	Code int         `json:"code"`
// 执行：Data interface{} `json:"data"`
	Data interface{} `json:"data"`
// 执行：Msg  string      `json:"msg"`
	Msg  string      `json:"msg"`
// 结束当前代码块
}

// OK 成功响应

// 函数定义：func OK(data interface{}) Body { return Body{Code: 0, Data: data, Msg: "success"} }
func OK(data interface{}) Body { return Body{Code: 0, Data: data, Msg: "success"} }

// Fail 失败响应

// 函数定义：func Fail(code int, msg string) Body { return Body{Code: code, Data: nil, Msg: msg} }
func Fail(code int, msg string) Body { return Body{Code: code, Data: nil, Msg: msg} }

// LoginResult 登录结果

// 类型定义：type LoginResult struct {
type LoginResult struct {
// 执行：Token    string     `json:"token"`
	Token    string     `json:"token"`
// 执行：User     model.User `json:"user"`
	User     model.User `json:"user"`
// 执行：ExpireAt string     `json:"expireAt"`
	ExpireAt string     `json:"expireAt"`
// 结束当前代码块
}

// OrderDetail 订单聚合视图：买家能看懂钱与订单去向

// 类型定义：type OrderDetail struct {
type OrderDetail struct {
// 执行：Order         model.Order             `json:"order"`
	Order         model.Order             `json:"order"`
// 执行：Payment       *model.Payment          `json:"payment,omitempty"`
	Payment       *model.Payment          `json:"payment,omitempty"`
// 执行：Refund        *model.Refund           `json:"refund,omitempty"`
	Refund        *model.Refund           `json:"refund,omitempty"`
// 执行：Shipment      *model.Shipment         `json:"shipment,omitempty"`
	Shipment      *model.Shipment         `json:"shipment,omitempty"`
// 执行：Abnormal      []model.AbnormalPayment `json:"abnormal,omitempty"`
	Abnormal      []model.AbnormalPayment `json:"abnormal,omitempty"`
// 执行：Notices       []model.PaymentNotice   `json:"notices,omitempty"`
	Notices       []model.PaymentNotice   `json:"notices,omitempty"`
// 执行：Timeline      []model.OperationLog    `json:"timeline,omitempty"`
	Timeline      []model.OperationLog    `json:"timeline,omitempty"`
// 执行：ProductInfo   *model.Product          `json:"product,omitempty"`
	ProductInfo   *model.Product          `json:"product,omitempty"`
// 执行：ShipableNow   bool                    `json:"shipableNow"`
	ShipableNow   bool                    `json:"shipableNow"`
// 执行：RefundPending bool                    `json:"refundPending"`
	RefundPending bool                    `json:"refundPending"`
// 结束当前代码块
}

// ProductStat 商品维度汇总

// 类型定义：type ProductStat struct {
type ProductStat struct {
// 执行：ProductID     uint64 `json:"productId"`
	ProductID     uint64 `json:"productId"`
// 执行：SKU           string `json:"sku"`
	SKU           string `json:"sku"`
// 执行：Name          string `json:"name"`
	Name          string `json:"name"`
// 执行：PriceCents    int64  `json:"priceCents"`
	PriceCents    int64  `json:"priceCents"`
// 执行：Available     int    `json:"available"`
	Available     int    `json:"available"`
// 执行：InitialStock  int    `json:"initialStock"`
	InitialStock  int    `json:"initialStock"`
// 执行：OnSale        bool   `json:"onSale"`
	OnSale        bool   `json:"onSale"`
// 执行：Occupied      int    `json:"occupied"`
	Occupied      int    `json:"occupied"`
// 执行：ShippedQty    int    `json:"shippedQty"`
	ShippedQty    int    `json:"shippedQty"`
// 执行：PaidCents     int64  `json:"paidCents"`
	PaidCents     int64  `json:"paidCents"`
// 执行：AbnormalCents int64  `json:"abnormalCents"` // 异常渠道收款，与有效收款分开显示（OPS-03）
	AbnormalCents int64  `json:"abnormalCents"` // 异常渠道收款，与有效收款分开显示（OPS-03）
// 执行：RefundedCents int64  `json:"refundedCents"`
	RefundedCents int64  `json:"refundedCents"`
// 结束当前代码块
}

// TradeSummary 交易汇总（运营核查用）

// 类型定义：type TradeSummary struct {
type TradeSummary struct {
// 执行：Products           []ProductStat    `json:"products"`
	Products           []ProductStat    `json:"products"`
// 执行：TotalPaidCents     int64            `json:"totalPaidCents"`
	TotalPaidCents     int64            `json:"totalPaidCents"`
// 执行：TotalAbnormalCents int64            `json:"totalAbnormalCents"` // 异常渠道收款合计（BR-14：不得并入有效收款）
	TotalAbnormalCents int64            `json:"totalAbnormalCents"` // 异常渠道收款合计（BR-14：不得并入有效收款）
// 执行：TotalRefundedCents int64            `json:"totalRefundedCents"`
	TotalRefundedCents int64            `json:"totalRefundedCents"`
// 执行：OrderCount         int64            `json:"orderCount"`
	OrderCount         int64            `json:"orderCount"`
// 执行：StatusCount        map[string]int64 `json:"statusCount"`
	StatusCount        map[string]int64 `json:"statusCount"`
// 结束当前代码块
}
