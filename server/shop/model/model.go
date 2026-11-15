// Package model 定义 GORM 数据模型与业务枚举。
// 表结构由 migrations/*.sql 建立（保留 InnoDB 的 ENUM / CHECK / UNIQUE 约束），
// 模型仅做列映射，因此每张表都显式声明 TableName，避免 GORM 复数化推断偏差。
package model

// 导入依赖包
import "time"

// 订单状态机：pending_payment -> paid -> shipped / refund_requested -> refunded|refund_failed

// 声明：const (
const (
// 赋值/更新：OrderStatusPendingPayment  = "pending_payment"
	OrderStatusPendingPayment  = "pending_payment"
// 赋值/更新：OrderStatusCancelled       = "cancelled"
	OrderStatusCancelled       = "cancelled"
// 赋值/更新：OrderStatusClosed          = "closed"
	OrderStatusClosed          = "closed"
// 赋值/更新：OrderStatusPaid            = "paid"
	OrderStatusPaid            = "paid"
// 赋值/更新：OrderStatusShipped         = "shipped"
	OrderStatusShipped         = "shipped"
// 赋值/更新：OrderStatusRefundRequested = "refund_requested"
	OrderStatusRefundRequested = "refund_requested"
// 赋值/更新：OrderStatusRefundFailed    = "refund_failed"
	OrderStatusRefundFailed    = "refund_failed"
// 赋值/更新：OrderStatusRefunded        = "refunded"
	OrderStatusRefunded        = "refunded"
// 结束当前代码块
)

// 角色

// 声明：const (
const (
// 赋值/更新：RoleBuyer     = "buyer"
	RoleBuyer     = "buyer"
// 赋值/更新：RoleOps       = "ops"
	RoleOps       = "ops"
// 赋值/更新：RoleWarehouse = "warehouse"
	RoleWarehouse = "warehouse"
// 赋值/更新：RoleCS        = "cs"
	RoleCS        = "cs"
// 结束当前代码块
)

// 退款来源：售前退款 / 迟到支付退款

// 声明：const (
const (
// 赋值/更新：RefundSourcePresale      = "presale"
	RefundSourcePresale      = "presale"
// 赋值/更新：RefundSourceAbnormalLate = "abnormal_late"
	RefundSourceAbnormalLate = "abnormal_late"
// 结束当前代码块
)

// 退款审核与渠道状态

// 声明：const (
const (
// 赋值/更新：AuditPending  = "pending"
	AuditPending  = "pending"
// 赋值/更新：AuditApproved = "approved"
	AuditApproved = "approved"
// 赋值/更新：AuditRejected = "rejected"
	AuditRejected = "rejected"

// 赋值/更新：ChannelPending = "pending"
	ChannelPending = "pending"
// 赋值/更新：ChannelSuccess = "success"
	ChannelSuccess = "success"
// 赋值/更新：ChannelFailed  = "failed"
	ChannelFailed  = "failed"
// 结束当前代码块
)

// 支付通知处理结果

// 声明：const (
const (
// 赋值/更新：NoticeAccepted   = "accepted"
	NoticeAccepted   = "accepted"
// 赋值/更新：NoticeRejected   = "rejected"
	NoticeRejected   = "rejected"
// 赋值/更新：NoticeDuplicated = "duplicated"
	NoticeDuplicated = "duplicated"
// 赋值/更新：NoticeAbnormal   = "abnormal"
	NoticeAbnormal   = "abnormal"
// 结束当前代码块
)

// 库存流水类型

// 声明：const (
const (
// 赋值/更新：LedgerOccupy  = "occupy"
	LedgerOccupy  = "occupy"
// 赋值/更新：LedgerRelease = "release"
	LedgerRelease = "release"
// 赋值/更新：LedgerRestore = "restore"
	LedgerRestore = "restore"
// 赋值/更新：LedgerShipOut = "ship_out"
	LedgerShipOut = "ship_out"
// 结束当前代码块
)

// 异常支付处理状态

// 声明：const (
const (
// 赋值/更新：AbnormalPending = "pending"
	AbnormalPending = "pending"
// 赋值/更新：AbnormalHandled = "handled"
	AbnormalHandled = "handled"
// 结束当前代码块
)

// User 买家与员工账号

// 类型定义：type User struct {
type User struct {
// 执行：ID           uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
// 执行：Username     string    `gorm:"column:username" json:"username"`
	Username     string    `gorm:"column:username" json:"username"`
// 执行：PasswordHash string    `gorm:"column:password_hash" json:"-"`
	PasswordHash string    `gorm:"column:password_hash" json:"-"`
// 执行：Role         string    `gorm:"column:role" json:"role"`
	Role         string    `gorm:"column:role" json:"role"`
// 执行：DisplayName  string    `gorm:"column:display_name" json:"displayName"`
	DisplayName  string    `gorm:"column:display_name" json:"displayName"`
// 执行：CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
// 结束当前代码块
}

// TableName 鍥哄畾琛ㄥ悕 users

// 函数定义：func (User) TableName() string { return "users" }
func (User) TableName() string { return "users" }

// Product 商品与可售数量

// 类型定义：type Product struct {
type Product struct {
// 执行：ID           uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
// 执行：SKU          string    `gorm:"column:sku" json:"sku"`
	SKU          string    `gorm:"column:sku" json:"sku"`
// 执行：Name         string    `gorm:"column:name" json:"name"`
	Name         string    `gorm:"column:name" json:"name"`
// 执行：PriceCents   int64     `gorm:"column:price_cents" json:"priceCents"`
	PriceCents   int64     `gorm:"column:price_cents" json:"priceCents"`
// 执行：Available    int       `gorm:"column:available" json:"available"`
	Available    int       `gorm:"column:available" json:"available"`
// 执行：InitialStock int       `gorm:"column:initial_stock" json:"initialStock"`
	InitialStock int       `gorm:"column:initial_stock" json:"initialStock"`
// 执行：OnSale       bool      `gorm:"column:on_sale" json:"onSale"`
	OnSale       bool      `gorm:"column:on_sale" json:"onSale"`
// 执行：CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
// 执行：UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
// 结束当前代码块
}

// TableName 鍥哄畾琛ㄥ悕 products

// 函数定义：func (Product) TableName() string { return "products" }
func (Product) TableName() string { return "products" }

// Order 订单（含商品与价格快照）

// 类型定义：type Order struct {
type Order struct {
// 执行：ID          uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
// 执行：OrderNo     string     `gorm:"column:order_no" json:"orderNo"`
	OrderNo     string     `gorm:"column:order_no" json:"orderNo"`
// 执行：BuyerID     uint64     `gorm:"column:buyer_id" json:"buyerId"`
	BuyerID     uint64     `gorm:"column:buyer_id" json:"buyerId"`
// 执行：ProductID   uint64     `gorm:"column:product_id" json:"productId"`
	ProductID   uint64     `gorm:"column:product_id" json:"productId"`
// 执行：ProductName string     `gorm:"column:product_name" json:"productName"`
	ProductName string     `gorm:"column:product_name" json:"productName"`
// 执行：PriceCents  int64      `gorm:"column:price_cents" json:"priceCents"`
	PriceCents  int64      `gorm:"column:price_cents" json:"priceCents"`
// 执行：Qty         int        `gorm:"column:qty" json:"qty"`
	Qty         int        `gorm:"column:qty" json:"qty"`
// 执行：AmountCents int64      `gorm:"column:amount_cents" json:"amountCents"`
	AmountCents int64      `gorm:"column:amount_cents" json:"amountCents"`
// 执行：Status      string     `gorm:"column:status" json:"status"`
	Status      string     `gorm:"column:status" json:"status"`
// 执行：ExpireAt    time.Time  `gorm:"column:expire_at" json:"expireAt"`
	ExpireAt    time.Time  `gorm:"column:expire_at" json:"expireAt"`
// 执行：PaidAt      *time.Time `gorm:"column:paid_at" json:"paidAt,omitempty"`
	PaidAt      *time.Time `gorm:"column:paid_at" json:"paidAt,omitempty"`
// 执行：ShippedAt   *time.Time `gorm:"column:shipped_at" json:"shippedAt,omitempty"`
	ShippedAt   *time.Time `gorm:"column:shipped_at" json:"shippedAt,omitempty"`
// 执行：ClosedAt    *time.Time `gorm:"column:closed_at" json:"closedAt,omitempty"`
	ClosedAt    *time.Time `gorm:"column:closed_at" json:"closedAt,omitempty"`
// 执行：CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
// 执行：UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updatedAt"`
// 结束当前代码块
}

// TableName 鍥哄畾琛ㄥ悕 orders

// 函数定义：func (Order) TableName() string { return "orders" }
func (Order) TableName() string { return "orders" }

// Payment 成功记款（幂等基石：order_no 与 channel_txn_id 双唯一）

// 类型定义：type Payment struct {
type Payment struct {
// 执行：ID           uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
// 执行：OrderNo      string    `gorm:"column:order_no" json:"orderNo"`
	OrderNo      string    `gorm:"column:order_no" json:"orderNo"`
// 执行：ChannelTxnID string    `gorm:"column:channel_txn_id" json:"channelTxnId"`
	ChannelTxnID string    `gorm:"column:channel_txn_id" json:"channelTxnId"`
// 执行：AmountCents  int64     `gorm:"column:amount_cents" json:"amountCents"`
	AmountCents  int64     `gorm:"column:amount_cents" json:"amountCents"`
// 执行：PaidAt       time.Time `gorm:"column:paid_at" json:"paidAt"`
	PaidAt       time.Time `gorm:"column:paid_at" json:"paidAt"`
// 执行：CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
// 结束当前代码块
}

// TableName 鍥哄畾琛ㄥ悕 payments

// 函数定义：func (Payment) TableName() string { return "payments" }
func (Payment) TableName() string { return "payments" }

// PaymentNotice 支付通知原始留痕（含被拒与重复）

// 类型定义：type PaymentNotice struct {
type PaymentNotice struct {
// 执行：ID           uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
// 执行：OrderNo      string    `gorm:"column:order_no" json:"orderNo"`
	OrderNo      string    `gorm:"column:order_no" json:"orderNo"`
// 执行：ChannelTxnID string    `gorm:"column:channel_txn_id" json:"channelTxnId"`
	ChannelTxnID string    `gorm:"column:channel_txn_id" json:"channelTxnId"`
// 执行：AmountCents  int64     `gorm:"column:amount_cents" json:"amountCents"`
	AmountCents  int64     `gorm:"column:amount_cents" json:"amountCents"`
// 执行：SignOK       bool      `gorm:"column:sign_ok" json:"signOk"`
	SignOK       bool      `gorm:"column:sign_ok" json:"signOk"`
// 执行：Result       string    `gorm:"column:result" json:"result"`
	Result       string    `gorm:"column:result" json:"result"`
// 执行：Reason       string    `gorm:"column:reason" json:"reason"`
	Reason       string    `gorm:"column:reason" json:"reason"`
// 执行：RawBody      *string   `gorm:"column:raw_body" json:"-"`
	RawBody      *string   `gorm:"column:raw_body" json:"-"`
// 执行：CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
// 结束当前代码块
}

// TableName 鍥哄畾琛ㄥ悕 payment_notices

// 函数定义：func (PaymentNotice) TableName() string { return "payment_notices" }
func (PaymentNotice) TableName() string { return "payment_notices" }

// Refund 退款单（整单退款，一单一张）

// 类型定义：type Refund struct {
type Refund struct {
// 执行：ID            uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ID            uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
// 执行：RefundNo      string     `gorm:"column:refund_no" json:"refundNo"`
	RefundNo      string     `gorm:"column:refund_no" json:"refundNo"`
// 执行：OrderNo       string     `gorm:"column:order_no" json:"orderNo"`
	OrderNo       string     `gorm:"column:order_no" json:"orderNo"`
// 执行：Source        string     `gorm:"column:source" json:"source"`
	Source        string     `gorm:"column:source" json:"source"`
// 执行：AmountCents   int64      `gorm:"column:amount_cents" json:"amountCents"`
	AmountCents   int64      `gorm:"column:amount_cents" json:"amountCents"`
// 执行：AuditStatus   string     `gorm:"column:audit_status" json:"auditStatus"`
	AuditStatus   string     `gorm:"column:audit_status" json:"auditStatus"`
// 执行：AuditReason   string     `gorm:"column:audit_reason" json:"auditReason"`
	AuditReason   string     `gorm:"column:audit_reason" json:"auditReason"`
// 执行：AuditBy       *uint64    `gorm:"column:audit_by" json:"auditBy,omitempty"`
	AuditBy       *uint64    `gorm:"column:audit_by" json:"auditBy,omitempty"`
// 执行：AuditAt       *time.Time `gorm:"column:audit_at" json:"auditAt,omitempty"`
	AuditAt       *time.Time `gorm:"column:audit_at" json:"auditAt,omitempty"`
// 执行：ChannelStatus string     `gorm:"column:channel_status" json:"channelStatus"`
	ChannelStatus string     `gorm:"column:channel_status" json:"channelStatus"`
// 执行：RetryCount    int        `gorm:"column:retry_count" json:"retryCount"`
	RetryCount    int        `gorm:"column:retry_count" json:"retryCount"`
// 执行：Restored      bool       `gorm:"column:restored" json:"restored"`
	Restored      bool       `gorm:"column:restored" json:"restored"`
// 执行：CreatedAt     time.Time  `gorm:"column:created_at" json:"createdAt"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"createdAt"`
// 执行：UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updatedAt"`
// 结束当前代码块
}

// TableName 鍥哄畾琛ㄥ悕 refunds

// 函数定义：func (Refund) TableName() string { return "refunds" }
func (Refund) TableName() string { return "refunds" }

// Shipment 发货登记

// 类型定义：type Shipment struct {
type Shipment struct {
// 执行：ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
// 执行：OrderNo    string    `gorm:"column:order_no" json:"orderNo"`
	OrderNo    string    `gorm:"column:order_no" json:"orderNo"`
// 执行：Carrier    string    `gorm:"column:carrier" json:"carrier"`
	Carrier    string    `gorm:"column:carrier" json:"carrier"`
// 执行：TrackingNo string    `gorm:"column:tracking_no" json:"trackingNo"`
	TrackingNo string    `gorm:"column:tracking_no" json:"trackingNo"`
// 执行：CreatedBy  *uint64   `gorm:"column:created_by" json:"createdBy,omitempty"`
	CreatedBy  *uint64   `gorm:"column:created_by" json:"createdBy,omitempty"`
// 执行：CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt"`
// 结束当前代码块
}

// TableName 鍥哄畾琛ㄥ悕 shipments

// 函数定义：func (Shipment) TableName() string { return "shipments" }
func (Shipment) TableName() string { return "shipments" }

// AbnormalPayment 迟到支付等异常收款待处理

// 类型定义：type AbnormalPayment struct {
type AbnormalPayment struct {
// 执行：ID           uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ID           uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
// 执行：OrderNo      string     `gorm:"column:order_no" json:"orderNo"`
	OrderNo      string     `gorm:"column:order_no" json:"orderNo"`
// 执行：ChannelTxnID string     `gorm:"column:channel_txn_id" json:"channelTxnId"`
	ChannelTxnID string     `gorm:"column:channel_txn_id" json:"channelTxnId"`
// 执行：AmountCents  int64      `gorm:"column:amount_cents" json:"amountCents"`
	AmountCents  int64      `gorm:"column:amount_cents" json:"amountCents"`
// 执行：Reason       string     `gorm:"column:reason" json:"reason"`
	Reason       string     `gorm:"column:reason" json:"reason"`
// 执行：Status       string     `gorm:"column:status" json:"status"`
	Status       string     `gorm:"column:status" json:"status"`
// 执行：HandledBy    *uint64    `gorm:"column:handled_by" json:"handledBy,omitempty"`
	HandledBy    *uint64    `gorm:"column:handled_by" json:"handledBy,omitempty"`
// 执行：HandledAt    *time.Time `gorm:"column:handled_at" json:"handledAt,omitempty"`
	HandledAt    *time.Time `gorm:"column:handled_at" json:"handledAt,omitempty"`
// 执行：RefundID     *uint64    `gorm:"column:refund_id" json:"refundId,omitempty"`
	RefundID     *uint64    `gorm:"column:refund_id" json:"refundId,omitempty"`
// 执行：CreatedAt    time.Time  `gorm:"column:created_at" json:"createdAt"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"createdAt"`
// 结束当前代码块
}

// TableName 鍥哄畾琛ㄥ悕 abnormal_payments

// 函数定义：func (AbnormalPayment) TableName() string { return "abnormal_payments" }
func (AbnormalPayment) TableName() string { return "abnormal_payments" }

// InventoryLedger 库存变动流水（守恒校验依据）

// 类型定义：type InventoryLedger struct {
type InventoryLedger struct {
// 执行：ID             uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
// 执行：ProductID      uint64    `gorm:"column:product_id" json:"productId"`
	ProductID      uint64    `gorm:"column:product_id" json:"productId"`
// 执行：OrderNo        string    `gorm:"column:order_no" json:"orderNo"`
	OrderNo        string    `gorm:"column:order_no" json:"orderNo"`
// 执行：ChangeType     string    `gorm:"column:change_type" json:"changeType"`
	ChangeType     string    `gorm:"column:change_type" json:"changeType"`
// 执行：Delta          int       `gorm:"column:delta" json:"delta"`
	Delta          int       `gorm:"column:delta" json:"delta"`
// 执行：AvailableAfter int       `gorm:"column:available_after" json:"availableAfter"`
	AvailableAfter int       `gorm:"column:available_after" json:"availableAfter"`
// 执行：CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
// 结束当前代码块
}

// TableName 鍥哄畾琛ㄥ悕 inventory_ledger

// 函数定义：func (InventoryLedger) TableName() string { return "inventory_ledger" }
func (InventoryLedger) TableName() string { return "inventory_ledger" }

// OperationLog 操作留痕（追加写，不提供修改与删除接口）

// 类型定义：type OperationLog struct {
type OperationLog struct {
// 执行：ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
// 执行：OrderNo   *string   `gorm:"column:order_no" json:"orderNo,omitempty"`
	OrderNo   *string   `gorm:"column:order_no" json:"orderNo,omitempty"`
// 执行：ActorID   *uint64   `gorm:"column:actor_id" json:"actorId,omitempty"`
	ActorID   *uint64   `gorm:"column:actor_id" json:"actorId,omitempty"`
// 执行：ActorRole string    `gorm:"column:actor_role" json:"actorRole"`
	ActorRole string    `gorm:"column:actor_role" json:"actorRole"`
// 执行：Action    string    `gorm:"column:action" json:"action"`
	Action    string    `gorm:"column:action" json:"action"`
// 执行：Detail    string    `gorm:"column:detail" json:"detail"`
	Detail    string    `gorm:"column:detail" json:"detail"`
// 执行：CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
// 结束当前代码块
}

// TableName 鍥哄畾琛ㄥ悕 operation_logs

// 函数定义：func (OperationLog) TableName() string { return "operation_logs" }
func (OperationLog) TableName() string { return "operation_logs" }
