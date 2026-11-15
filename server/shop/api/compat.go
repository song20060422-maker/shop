// 包声明
package api

// 导入依赖包
import (
// 执行："crypto/hmac"
	"crypto/hmac"
// 执行："crypto/sha256"
	"crypto/sha256"
// 执行："encoding/hex"
	"encoding/hex"
// 执行："encoding/json"
	"encoding/json"
// 执行："fmt"
	"fmt"
// 执行："net/http"
	"net/http"
// 执行："strconv"
	"strconv"
// 执行："strings"
	"strings"
// 执行："time"
	"time"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/config"
	"shop/shop/config"
// 执行："shop/shop/middleware"
	"shop/shop/middleware"
// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/service"
	"shop/shop/service"
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// compatStaffAlias 接口契约账号名 -> 原型种子账号名（仅登录兜底使用）。

// 声明：var compatStaffAlias = map[string]string{
var compatStaffAlias = map[string]string{
// 执行："ops":       "ops1",
	"ops":       "ops1",
// 执行："warehouse": "wh1",
	"warehouse": "wh1",
// 执行："cs":        "cs1",
	"cs":        "cs1",
// 结束当前代码块
}

// CompatApi 接口契约兼容层（api-contract.md 的调用约定）。
//
// 设计原则：
//  1. 完全复用既有 service 层实现，不复制业务逻辑、不直接操作数据库；
//  2. 输出「契约扁平字段 + 原型 data 副本」双轨结构：验收脚本读顶层字段，
//     内置演示页面继续读 data，既有 /api 原型页面的调用方式不被破坏；
//  3. 对状态、金额单位、字段命名做契约侧映射（分、WAIT_PAY/TO_SHIP 等）。

// 类型定义：type CompatApi struct {
type CompatApi struct {
// 执行：s   *service.ServiceGroup
	s   *service.ServiceGroup
// 执行：cfg config.Config
	cfg config.Config
// 执行：now func() time.Time
	now func() time.Time
// 结束当前代码块
}

// NewCompatApi 构造兼容层处理器。

// 函数定义：func NewCompatApi(g *service.ServiceGroup, cfg config.Config, now func() time.Time) *CompatApi {
func NewCompatApi(g *service.ServiceGroup, cfg config.Config, now func() time.Time) *CompatApi {
// 条件判断：if now == nil {
	if now == nil {
// 赋值/更新：now = time.Now
		now = time.Now
// 结束当前代码块
	}
// 返回：return &CompatApi{s: g, cfg: cfg, now: now}
	return &CompatApi{s: g, cfg: cfg, now: now}
// 结束当前代码块
}

// PublicGuard 公开区委派：命中与原型同名但契约语义不同的路径时改走兼容实现。

// 函数定义：func (a *CompatApi) PublicGuard() gin.HandlerFunc {
func (a *CompatApi) PublicGuard() gin.HandlerFunc {
// 返回：return func(c *gin.Context) {
	return func(c *gin.Context) {
// 条件判断：if a.delegate(c) {
		if a.delegate(c) {
// 执行：c.Abort()
			c.Abort()
// 返回：return
			return
// 结束当前代码块
		}
// 执行：c.Next()
		c.Next()
// 结束当前代码块
	}
// 结束当前代码块
}

// AuthGuard 登录区委派：必须挂在 JWT 中间件之后，保证 claims 已就绪。

// 函数定义：func (a *CompatApi) AuthGuard() gin.HandlerFunc {
func (a *CompatApi) AuthGuard() gin.HandlerFunc {
// 返回：return func(c *gin.Context) {
	return func(c *gin.Context) {
// 条件判断：if a.delegate(c) {
		if a.delegate(c) {
// 执行：c.Abort()
			c.Abort()
// 返回：return
			return
// 结束当前代码块
		}
// 执行：c.Next()
		c.Next()
// 结束当前代码块
	}
// 结束当前代码块
}

// delegate 按「方法 + 注册路径」判断是否需要交给契约兼容实现。

// 函数定义：func (a *CompatApi) delegate(c *gin.Context) bool {
func (a *CompatApi) delegate(c *gin.Context) bool {
// 赋值/更新：full := c.FullPath()
	full := c.FullPath()
// 赋值/更新：match := func(target string) bool { return full == target || strings.HasSuffix(full, target) }
	match := func(target string) bool { return full == target || strings.HasSuffix(full, target) }
// 多路分支：switch c.Request.Method {
	switch c.Request.Method {
// 匹配分支：case http.MethodPost:
	case http.MethodPost:
// 多路分支：switch {
		switch {
// 匹配分支：case match("/api/login"):
		case match("/api/login"):
// 执行：a.Login(c)
			a.Login(c)
// 返回：return true
			return true
// 匹配分支：case match("/api/pay/notify"):
		case match("/api/pay/notify"):
// 执行：a.PayNotify(c)
			a.PayNotify(c)
// 返回：return true
			return true
// 匹配分支：case match("/api/orders"):
		case match("/api/orders"):
// 执行：a.CreateOrder(c)
			a.CreateOrder(c)
// 返回：return true
			return true
// 匹配分支：case match("/api/orders/:orderNo/cancel"):
		case match("/api/orders/:orderNo/cancel"):
// 执行：a.CancelOrder(c)
			a.CancelOrder(c)
// 返回：return true
			return true
// 结束当前代码块
		}
// 匹配分支：case http.MethodGet:
	case http.MethodGet:
// 多路分支：switch {
		switch {
// 匹配分支：case match("/api/products"):
		case match("/api/products"):
// 执行：a.Products(c)
			a.Products(c)
// 返回：return true
			return true
// 匹配分支：case match("/api/orders/:orderNo"):
		case match("/api/orders/:orderNo"):
// 执行：a.OrderDetail(c)
			a.OrderDetail(c)
// 返回：return true
			return true
// 匹配分支：case match("/api/ops/summary"):
		case match("/api/ops/summary"):
// 执行：a.OpsSummary(c)
			a.OpsSummary(c)
// 返回：return true
			return true
// 匹配分支：case match("/api/cs/refunds"):
		case match("/api/cs/refunds"):
// 执行：a.CSRefunds(c)
			a.CSRefunds(c)
// 返回：return true
			return true
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return false
	return false
// 结束当前代码块
}

// ---------------- 状态与视图映射 ----------------

// orderStatusLabel 内部订单状态 → 契约状态（WAIT_PAY / TO_SHIP / SHIPPED / REFUND_PENDING / REFUNDED / CLOSED）。

// 函数定义：func orderStatusLabel(st string) string {
func orderStatusLabel(st string) string {
// 多路分支：switch st {
	switch st {
// 匹配分支：case model.OrderStatusPendingPayment:
	case model.OrderStatusPendingPayment:
// 返回：return "WAIT_PAY"
		return "WAIT_PAY"
// 匹配分支：case model.OrderStatusPaid:
	case model.OrderStatusPaid:
// 返回：return "TO_SHIP"
		return "TO_SHIP"
// 匹配分支：case model.OrderStatusShipped:
	case model.OrderStatusShipped:
// 返回：return "SHIPPED"
		return "SHIPPED"
// 匹配分支：case model.OrderStatusRefundRequested, model.OrderStatusRefundFailed:
	case model.OrderStatusRefundRequested, model.OrderStatusRefundFailed:
// 返回：return "REFUND_PENDING"
		return "REFUND_PENDING"
// 匹配分支：case model.OrderStatusRefunded:
	case model.OrderStatusRefunded:
// 返回：return "REFUNDED"
		return "REFUNDED"
// 默认分支
	default:
// 返回：return "CLOSED"
		return "CLOSED"
// 结束当前代码块
	}
// 结束当前代码块
}

// refundStatusLabel 退款单内部状态 → 契约状态（PENDING / REFUNDING / FAILED / REFUNDED / REJECTED）。

// 函数定义：func refundStatusLabel(r model.Refund) string {
func refundStatusLabel(r model.Refund) string {
// 多路分支：switch {
	switch {
// 匹配分支：case r.AuditStatus == model.AuditRejected:
	case r.AuditStatus == model.AuditRejected:
// 返回：return "REJECTED"
		return "REJECTED"
// 匹配分支：case r.ChannelStatus == model.ChannelSuccess:
	case r.ChannelStatus == model.ChannelSuccess:
// 返回：return "REFUNDED"
		return "REFUNDED"
// 匹配分支：case r.ChannelStatus == model.ChannelFailed:
	case r.ChannelStatus == model.ChannelFailed:
// 返回：return "FAILED"
		return "FAILED"
// 匹配分支：case r.AuditStatus == model.AuditApproved:
	case r.AuditStatus == model.AuditApproved:
// 返回：return "REFUNDING"
		return "REFUNDING"
// 默认分支
	default:
// 返回：return "PENDING"
		return "PENDING"
// 结束当前代码块
	}
// 结束当前代码块
}

// refundFailReason 渠道失败原因（契约 refund.fail_reason）。

// 函数定义：func refundFailReason(r model.Refund) string {
func refundFailReason(r model.Refund) string {
// 条件判断：if r.ChannelStatus != model.ChannelFailed {
	if r.ChannelStatus != model.ChannelFailed {
// 返回：return ""
		return ""
// 结束当前代码块
	}
// 条件判断：if r.AuditReason != "" {
	if r.AuditReason != "" {
// 返回：return "渠道退款未受理：" + r.AuditReason
		return "渠道退款未受理：" + r.AuditReason
// 结束当前代码块
	}
// 返回：return "渠道退款未受理，可重试发起"
	return "渠道退款未受理，可重试发起"
// 结束当前代码块
}

// ts 统一时间输出格式（RFC3339）。旧版验收脚本会尝试解析 ISO8601。

// 函数定义：func ts(t time.Time) string {
func ts(t time.Time) string {
// 条件判断：if t.IsZero() {
	if t.IsZero() {
// 返回：return ""
		return ""
// 结束当前代码块
	}
// 返回：return t.Format(time.RFC3339)
	return t.Format(time.RFC3339)
// 结束当前代码块
}

// tsP 指针时间输出，nil 时返回空字符串。

// 函数定义：func tsP(t *time.Time) string {
func tsP(t *time.Time) string {
// 条件判断：if t == nil {
	if t == nil {
// 返回：return ""
		return ""
// 结束当前代码块
	}
// 返回：return ts(*t)
	return ts(*t)
// 结束当前代码块
}

// productView 商品视图：id 使用 SKU（契约约定），同时保留原有 camelCase 字段。

// 函数定义：func productView(p model.Product) gin.H {
func productView(p model.Product) gin.H {
// 返回：return gin.H{
	return gin.H{
// 执行："id":            p.SKU,
		"id":            p.SKU,
// 执行："sku":           p.SKU,
		"sku":           p.SKU,
// 执行："product_id":    p.SKU,
		"product_id":    p.SKU,
// 执行："productId":     p.ID,
		"productId":     p.ID,
// 执行："db_id":         p.ID,
		"db_id":         p.ID,
// 执行："name":          p.Name,
		"name":          p.Name,
// 执行："price_cents":   p.PriceCents,
		"price_cents":   p.PriceCents,
// 执行："priceCents":    p.PriceCents,
		"priceCents":    p.PriceCents,
// 执行："available":     p.Available,
		"available":     p.Available,
// 执行："initial_stock": p.InitialStock,
		"initial_stock": p.InitialStock,
// 执行："initialStock":  p.InitialStock,
		"initialStock":  p.InitialStock,
// 执行："on_sale":       p.OnSale,
		"on_sale":       p.OnSale,
// 执行："onSale":        p.OnSale,
		"onSale":        p.OnSale,
// 结束当前代码块
	}
// 结束当前代码块
}

// refundView 退款单视图（契约 {id,status,reason,fail_reason}）。

// 函数定义：func refundView(r model.Refund) gin.H {
func refundView(r model.Refund) gin.H {
// 赋值/更新：reason := r.AuditReason
	reason := r.AuditReason
// 条件判断：if reason == "" {
	if reason == "" {
// 赋值/更新：reason = r.Source
		reason = r.Source
// 结束当前代码块
	}
// 返回：return gin.H{
	return gin.H{
// 执行："id":             r.RefundNo,
		"id":             r.RefundNo,
// 执行："refund_id":      r.RefundNo,
		"refund_id":      r.RefundNo,
// 执行："refund_no":      r.RefundNo,
		"refund_no":      r.RefundNo,
// 执行："refundNo":       r.RefundNo,
		"refundNo":       r.RefundNo,
// 执行："order_no":       r.OrderNo,
		"order_no":       r.OrderNo,
// 执行："orderNo":        r.OrderNo,
		"orderNo":        r.OrderNo,
// 执行："status":         refundStatusLabel(r),
		"status":         refundStatusLabel(r),
// 执行："statusLabel":    refundStatusLabel(r),
		"statusLabel":    refundStatusLabel(r),
// 执行："amount":         r.AmountCents,
		"amount":         r.AmountCents,
// 执行："amount_cents":   r.AmountCents,
		"amount_cents":   r.AmountCents,
// 执行："amountCents":    r.AmountCents,
		"amountCents":    r.AmountCents,
// 执行："reason":         reason,
		"reason":         reason,
// 执行："apply_reason":   r.Source,
		"apply_reason":   r.Source,
// 执行："audit_reason":   r.AuditReason,
		"audit_reason":   r.AuditReason,
// 执行："fail_reason":    refundFailReason(r),
		"fail_reason":    refundFailReason(r),
// 执行："channel_status": r.ChannelStatus,
		"channel_status": r.ChannelStatus,
// 执行："audit_status":   r.AuditStatus,
		"audit_status":   r.AuditStatus,
// 执行："restored":       r.Restored,
		"restored":       r.Restored,
// 执行："retry_count":    r.RetryCount,
		"retry_count":    r.RetryCount,
// 执行："retryCount":     r.RetryCount,
		"retryCount":     r.RetryCount,
// 执行："created_at":     ts(r.CreatedAt),
		"created_at":     ts(r.CreatedAt),
// 执行："createdAt":      ts(r.CreatedAt),
		"createdAt":      ts(r.CreatedAt),
// 结束当前代码块
	}
// 结束当前代码块
}

// orderView 订单视图：同时提供契约 snake_case 字段与原型 camelCase 字段。

// 函数定义：func (a *CompatApi) orderView(o model.Order, refund *model.Refund, ship *model.Shipment) gin.H {
func (a *CompatApi) orderView(o model.Order, refund *model.Refund, ship *model.Shipment) gin.H {
// 赋值/更新：sku := ""
	sku := ""
// 条件判断：if p, err := a.s.CatalogService.Get(o.ProductID); err == nil {
	if p, err := a.s.CatalogService.Get(o.ProductID); err == nil {
// 赋值/更新：sku = p.SKU
		sku = p.SKU
// 结束当前代码块
	}
// 赋值/更新：view := gin.H{
	view := gin.H{
// 执行："order_id":     o.ID,
		"order_id":     o.ID,
// 执行："orderId":      o.ID,
		"orderId":      o.ID,
// 执行："id":           o.ID,
		"id":           o.ID,
// 执行："order_no":     o.OrderNo,
		"order_no":     o.OrderNo,
// 执行："orderNo":      o.OrderNo,
		"orderNo":      o.OrderNo,
// 执行："product_id":   sku,
		"product_id":   sku,
// 执行："productId":    o.ProductID,
		"productId":    o.ProductID,
// 执行："product_name": o.ProductName,
		"product_name": o.ProductName,
// 执行："productName":  o.ProductName,
		"productName":  o.ProductName,
// 执行："price_cents":  o.PriceCents,
		"price_cents":  o.PriceCents,
// 执行："priceCents":   o.PriceCents,
		"priceCents":   o.PriceCents,
// 执行："quantity":     o.Qty,
		"quantity":     o.Qty,
// 执行："qty":          o.Qty,
		"qty":          o.Qty,
// 执行："amount":       o.AmountCents,
		"amount":       o.AmountCents,
// 执行："amount_cents": o.AmountCents,
		"amount_cents": o.AmountCents,
// 执行："amountCents":  o.AmountCents,
		"amountCents":  o.AmountCents,
// 执行："status":       orderStatusLabel(o.Status),
		"status":       orderStatusLabel(o.Status),
// 执行："statusLabel":  orderStatusLabel(o.Status),
		"statusLabel":  orderStatusLabel(o.Status),
// 执行："raw_status":   o.Status,
		"raw_status":   o.Status,
// 执行："buyer_id":     o.BuyerID,
		"buyer_id":     o.BuyerID,
// 执行："buyerId":      o.BuyerID,
		"buyerId":      o.BuyerID,
// 执行："expire_at":    ts(o.ExpireAt),
		"expire_at":    ts(o.ExpireAt),
// 执行："expireAt":     ts(o.ExpireAt),
		"expireAt":     ts(o.ExpireAt),
// 执行："created_at":   ts(o.CreatedAt),
		"created_at":   ts(o.CreatedAt),
// 执行："createdAt":    ts(o.CreatedAt),
		"createdAt":    ts(o.CreatedAt),
// 执行："paid_at":      tsP(o.PaidAt),
		"paid_at":      tsP(o.PaidAt),
// 执行："paidAt":       tsP(o.PaidAt),
		"paidAt":       tsP(o.PaidAt),
// 执行："closed_at":    tsP(o.ClosedAt),
		"closed_at":    tsP(o.ClosedAt),
// 执行："carrier":      "",
		"carrier":      "",
// 执行："tracking_no":  "",
		"tracking_no":  "",
// 执行："trackingNo":   "",
		"trackingNo":   "",
// 结束当前代码块
	}
// 条件判断：if ship != nil {
	if ship != nil {
// 赋值/更新：view["carrier"] = ship.Carrier
		view["carrier"] = ship.Carrier
// 赋值/更新：view["tracking_no"] = ship.TrackingNo
		view["tracking_no"] = ship.TrackingNo
// 赋值/更新：view["trackingNo"] = ship.TrackingNo
		view["trackingNo"] = ship.TrackingNo
// 结束当前代码块
	}
// 条件判断：if refund != nil {
	if refund != nil {
// 赋值/更新：view["refund"] = refundView(*refund)
		view["refund"] = refundView(*refund)
// 条件分支：} else {
	} else {
// 赋值/更新：view["refund"] = nil
		view["refund"] = nil
// 结束当前代码块
	}
// 返回：return view
	return view
// 结束当前代码块
}

// ---------------- 请求解析辅助 ----------------

// 函数定义：func bindAny(c *gin.Context) map[string]interface{} {
func bindAny(c *gin.Context) map[string]interface{} {
// 赋值/更新：body := map[string]interface{}{}
	body := map[string]interface{}{}
// 赋值/更新：raw, err := c.GetRawData()
	raw, err := c.GetRawData()
// 条件判断：if err == nil && len(raw) > 0 {
	if err == nil && len(raw) > 0 {
// 赋值/更新：_ = json.Unmarshal(raw, &body)
		_ = json.Unmarshal(raw, &body)
// 结束当前代码块
	}
// 返回：return body
	return body
// 结束当前代码块
}

// strOf 依次取第一个非空字段，兼容契约 snake_case 与原型 camelCase。

// 函数定义：func strOf(m map[string]interface{}, keys ...string) string {
func strOf(m map[string]interface{}, keys ...string) string {
// 循环处理：for _, k := range keys {
	for _, k := range keys {
// 赋值/更新：v, ok := m[k]
		v, ok := m[k]
// 条件判断：if !ok || v == nil {
		if !ok || v == nil {
// 流程控制：continue
			continue
// 结束当前代码块
		}
// 多路分支：switch t := v.(type) {
		switch t := v.(type) {
// 匹配分支：case string:
		case string:
// 条件判断：if s := strings.TrimSpace(t); s != "" {
			if s := strings.TrimSpace(t); s != "" {
// 返回：return s
				return s
// 结束当前代码块
			}
// 匹配分支：case float64:
		case float64:
// 条件判断：if t == float64(int64(t)) {
			if t == float64(int64(t)) {
// 返回：return strconv.FormatInt(int64(t), 10)
				return strconv.FormatInt(int64(t), 10)
// 结束当前代码块
			}
// 返回：return strconv.FormatFloat(t, 'f', -1, 64)
			return strconv.FormatFloat(t, 'f', -1, 64)
// 匹配分支：case json.Number:
		case json.Number:
// 返回：return t.String()
			return t.String()
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return ""
	return ""
// 结束当前代码块
}

// intOf 取整数值字段。

// 函数定义：func intOf(m map[string]interface{}, keys ...string) (int64, bool) {
func intOf(m map[string]interface{}, keys ...string) (int64, bool) {
// 循环处理：for _, k := range keys {
	for _, k := range keys {
// 赋值/更新：v, ok := m[k]
		v, ok := m[k]
// 条件判断：if !ok || v == nil {
		if !ok || v == nil {
// 流程控制：continue
			continue
// 结束当前代码块
		}
// 多路分支：switch t := v.(type) {
		switch t := v.(type) {
// 匹配分支：case float64:
		case float64:
// 返回：return int64(t), true
			return int64(t), true
// 匹配分支：case int64:
		case int64:
// 返回：return t, true
			return t, true
// 匹配分支：case json.Number:
		case json.Number:
// 条件判断：if n, err := t.Int64(); err == nil {
			if n, err := t.Int64(); err == nil {
// 返回：return n, true
				return n, true
// 结束当前代码块
			}
// 匹配分支：case string:
		case string:
// 条件判断：if n, err := strconv.ParseInt(strings.TrimSpace(t), 10, 64); err == nil {
			if n, err := strconv.ParseInt(strings.TrimSpace(t), 10, 64); err == nil {
// 返回：return n, true
				return n, true
// 结束当前代码块
			}
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return 0, false
	return 0, false
// 结束当前代码块
}

// boolOf 取布尔值字段。

// 函数定义：func boolOf(m map[string]interface{}, keys ...string) (bool, bool) {
func boolOf(m map[string]interface{}, keys ...string) (bool, bool) {
// 循环处理：for _, k := range keys {
	for _, k := range keys {
// 条件判断：if v, ok := m[k]; ok && v != nil {
		if v, ok := m[k]; ok && v != nil {
// 条件判断：if b, o := v.(bool); o {
			if b, o := v.(bool); o {
// 返回：return b, true
				return b, true
// 结束当前代码块
			}
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return false, false
	return false, false
// 结束当前代码块
}

// failJSON 契约要求的业务失败响应：4xx + {"code","message"}。

// 函数定义：func (a *CompatApi) failJSON(c *gin.Context, err error) {
func (a *CompatApi) failJSON(c *gin.Context, err error) {
// 赋值/更新：code, msg := service.BizCode(err)
	code, msg := service.BizCode(err)
// 条件判断：if msg == "" {
	if msg == "" {
// 赋值/更新：msg = err.Error()
		msg = err.Error()
// 结束当前代码块
	}
	// 未知内部错误（9000）保留底层原因，便于现场排查与验收取证
	if code == 9000 && err != nil {
// 赋值/更新：msg = "系统繁忙，请稍后重试（底层原因：" + err.Error() + "）"
		msg = "系统繁忙，请稍后重试（底层原因：" + err.Error() + "）"
// 结束当前代码块
	}
// 执行：c.JSON(statusOf(code), gin.H{
	c.JSON(statusOf(code), gin.H{
// 执行："code":    code,
		"code":    code,
// 执行："message": msg,
		"message": msg,
// 执行："msg":     msg,
		"msg":     msg,
// 结束当前代码块
	})
// 结束当前代码块
}

// okJSON 统一输出：顶层契约字段 + data 副本（原型页面继续读 data）。

// 函数定义：func (a *CompatApi) okJSON(c *gin.Context, flat gin.H, data gin.H) {
func (a *CompatApi) okJSON(c *gin.Context, flat gin.H, data gin.H) {
// 赋值/更新：out := gin.H{"code": 0, "msg": "ok"}
	out := gin.H{"code": 0, "msg": "ok"}
// 循环处理：for k, v := range flat {
	for k, v := range flat {
// 赋值/更新：out[k] = v
		out[k] = v
// 结束当前代码块
	}
// 条件判断：if data == nil {
	if data == nil {
// 赋值/更新：data = flat
		data = flat
// 结束当前代码块
	}
// 赋值/更新：out["data"] = data
	out["data"] = data
// 执行：c.JSON(http.StatusOK, out)
	c.JSON(http.StatusOK, out)
// 结束当前代码块
}

// resolveProduct 商品双寻址：SKU 优先，其次数字主键。

// 函数定义：func (a *CompatApi) resolveProduct(ref string) (model.Product, error) {
func (a *CompatApi) resolveProduct(ref string) (model.Product, error) {
// 条件判断：if p, err := a.s.CatalogService.ProductBySKU(ref); err == nil {
	if p, err := a.s.CatalogService.ProductBySKU(ref); err == nil {
// 返回：return p, nil
		return p, nil
// 结束当前代码块
	}
// 条件判断：if id, err := strconv.ParseUint(strings.TrimSpace(ref), 10, 64); err == nil && id > 0 {
	if id, err := strconv.ParseUint(strings.TrimSpace(ref), 10, 64); err == nil && id > 0 {
// 返回：return a.s.CatalogService.Get(id)
		return a.s.CatalogService.Get(id)
// 结束当前代码块
	}
// 返回：return model.Product{}, service.ErrProductNotFound
	return model.Product{}, service.ErrProductNotFound
// 结束当前代码块
}

// ---------------- 公开区接口 ----------------

// Login POST /api/login -> {token, role}

// 函数定义：func (a *CompatApi) Login(c *gin.Context) {
func (a *CompatApi) Login(c *gin.Context) {
// 赋值/更新：body := bindAny(c)
	body := bindAny(c)
// 赋值/更新：username := strOf(body, "username")
	username := strOf(body, "username")
// 赋值/更新：password := strOf(body, "password")
	password := strOf(body, "password")
// 条件判断：if username == "" || password == "" {
	if username == "" || password == "" {
// 执行：a.failJSON(c, service.NewBizError(1000, "用户名或密码不能为空"))
		a.failJSON(c, service.NewBizError(1000, "用户名或密码不能为空"))
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：res, err := a.s.AuthService.Login(username, password)
	res, err := a.s.AuthService.Login(username, password)
// 条件判断：if err != nil {
	if err != nil {
		// 接口契约以 ops/warehouse/cs 作为员工账号名，原型种子账号为 ops1/wh1/cs1；
		// 仅在原型账号登录失败时按别名重试，保证契约账号可直接登录且不改变原 /api/login 行为。
		if alias, ok := compatStaffAlias[strings.ToLower(username)]; ok {
// 赋值/更新：res, err = a.s.AuthService.Login(alias, password)
			res, err = a.s.AuthService.Login(alias, password)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：user := gin.H{
	user := gin.H{
// 执行："id":           res.User.ID,
		"id":           res.User.ID,
// 执行："username":     res.User.Username,
		"username":     res.User.Username,
// 执行："role":         res.User.Role,
		"role":         res.User.Role,
// 执行："display_name": res.User.DisplayName,
		"display_name": res.User.DisplayName,
// 执行："displayName":  res.User.DisplayName,
		"displayName":  res.User.DisplayName,
// 结束当前代码块
	}
// 执行：a.okJSON(c, gin.H{
	a.okJSON(c, gin.H{
// 执行："token":     res.Token,
		"token":     res.Token,
// 执行："role":      res.User.Role,
		"role":      res.User.Role,
// 执行："user":      user,
		"user":      user,
// 执行："expire_at": res.ExpireAt,
		"expire_at": res.ExpireAt,
// 结束当前代码块
	}, gin.H{
// 执行："token":    res.Token,
		"token":    res.Token,
// 执行："user":     user,
		"user":     user,
// 执行："expireAt": res.ExpireAt,
		"expireAt": res.ExpireAt,
// 结束当前代码块
	})
// 结束当前代码块
}

// Products GET /api/products -> [{id: sku, price_cents, available, on_sale}]

// 函数定义：func (a *CompatApi) Products(c *gin.Context) {
func (a *CompatApi) Products(c *gin.Context) {
// 赋值/更新：list, err := a.s.CatalogService.OnSale()
	list, err := a.s.CatalogService.OnSale()
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：items := make([]gin.H, 0, len(list))
	items := make([]gin.H, 0, len(list))
// 循环处理：for _, p := range list {
	for _, p := range list {
// 赋值/更新：items = append(items, productView(p))
		items = append(items, productView(p))
// 结束当前代码块
	}
// 执行：c.JSON(http.StatusOK, gin.H{
	c.JSON(http.StatusOK, gin.H{
// 执行："code":     0,
		"code":     0,
// 执行："msg":      "ok",
		"msg":      "ok",
// 执行："data":     items,
		"data":     items,
// 执行："products": items,
		"products": items,
// 结束当前代码块
	})
// 结束当前代码块
}

// ---------------- 买家接口 ----------------

// CreateOrder POST /api/orders -> {order_id, order_no, amount, status, expire_at}

// 函数定义：func (a *CompatApi) CreateOrder(c *gin.Context) {
func (a *CompatApi) CreateOrder(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
	// 兼容层委派路径绕过了路由组上的 RequireRole，必须在此复核角色（I-02 越权防护）
	if claims.Role != model.RoleBuyer {
// 执行：a.failJSON(c, service.NewBizError(2002, "当前角色无权执行该操作"))
		a.failJSON(c, service.NewBizError(2002, "当前角色无权执行该操作"))
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：body := bindAny(c)
	body := bindAny(c)
// 赋值/更新：ref := strOf(body, "product_id", "productId", "sku")
	ref := strOf(body, "product_id", "productId", "sku")
// 条件判断：if ref == "" {
	if ref == "" {
// 执行：a.failJSON(c, service.NewBizError(1000, "product_id 不能为空"))
		a.failJSON(c, service.NewBizError(1000, "product_id 不能为空"))
// 返回：return
		return
// 结束当前代码块
	}
	// 数量：缺省 1 件；显式传入时必须是 1~2 的整数（题面：每单 1-2 件），
	// 0 / 负数 / 超过 2 一律拒绝，不能静默兜底成 1 件。
	qty := 1
// 条件判断：if _, hasQty := body["quantity"]; hasQty {
	if _, hasQty := body["quantity"]; hasQty {
// 赋值/更新：n, ok := intOf(body, "quantity")
		n, ok := intOf(body, "quantity")
// 条件判断：if !ok || n < 1 || n > 2 {
		if !ok || n < 1 || n > 2 {
// 执行：a.failJSON(c, service.ErrQtyInvalid)
			a.failJSON(c, service.ErrQtyInvalid)
// 返回：return
			return
// 结束当前代码块
		}
// 赋值/更新：qty = int(n)
		qty = int(n)
// 条件分支：} else if _, hasQty := body["qty"]; hasQty {
	} else if _, hasQty := body["qty"]; hasQty {
// 赋值/更新：n, ok := intOf(body, "qty")
		n, ok := intOf(body, "qty")
// 条件判断：if !ok || n < 1 || n > 2 {
		if !ok || n < 1 || n > 2 {
// 执行：a.failJSON(c, service.ErrQtyInvalid)
			a.failJSON(c, service.ErrQtyInvalid)
// 返回：return
			return
// 结束当前代码块
		}
// 赋值/更新：qty = int(n)
		qty = int(n)
// 结束当前代码块
	}
// 赋值/更新：p, err := a.resolveProduct(ref)
	p, err := a.resolveProduct(ref)
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：o, err := a.s.OrderService.Create(claims.UserID, p.ID, qty)
	o, err := a.s.OrderService.Create(claims.UserID, p.ID, qty)
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：view := a.orderView(o, nil, nil)
	view := a.orderView(o, nil, nil)
// 赋值/更新：avail, _ := a.s.CatalogService.CurrentAvailable(p.ID)
	avail, _ := a.s.CatalogService.CurrentAvailable(p.ID)
// 执行：a.okJSON(c, view, gin.H{
	a.okJSON(c, view, gin.H{
// 执行："order":          view,
		"order":          view,
// 执行："availableAfter": avail,
		"availableAfter": avail,
// 结束当前代码块
	})
// 结束当前代码块
}

// OrderDetail GET /api/orders/{id} -> 订单详情（{id} 兼容数字主键与订单号）

// 函数定义：func (a *CompatApi) OrderDetail(c *gin.Context) {
func (a *CompatApi) OrderDetail(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：o, err := a.s.OrderService.FindByRef(c.Param("orderNo"))
	o, err := a.s.OrderService.FindByRef(c.Param("orderNo"))
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：detail, err := a.s.ReportService.OrderDetail(o.OrderNo, claims.UserID, claims.Role)
	detail, err := a.s.ReportService.OrderDetail(o.OrderNo, claims.UserID, claims.Role)
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：view := a.orderView(detail.Order, detail.Refund, detail.Shipment)
	view := a.orderView(detail.Order, detail.Refund, detail.Shipment)
// 赋值/更新：data := gin.H{
	data := gin.H{
// 执行："order":         view,
		"order":         view,
// 执行："shipableNow":   detail.ShipableNow,
		"shipableNow":   detail.ShipableNow,
// 执行："refundPending": detail.RefundPending,
		"refundPending": detail.RefundPending,
// 结束当前代码块
	}
// 条件判断：if detail.Payment != nil {
	if detail.Payment != nil {
// 赋值/更新：pay := gin.H{
		pay := gin.H{
// 执行："order_no":       detail.Payment.OrderNo,
			"order_no":       detail.Payment.OrderNo,
// 执行："channel_txn_no": detail.Payment.ChannelTxnID,
			"channel_txn_no": detail.Payment.ChannelTxnID,
// 执行："channelTxnNo":   detail.Payment.ChannelTxnID,
			"channelTxnNo":   detail.Payment.ChannelTxnID,
// 执行："amount":         detail.Payment.AmountCents,
			"amount":         detail.Payment.AmountCents,
// 执行："amountCents":    detail.Payment.AmountCents,
			"amountCents":    detail.Payment.AmountCents,
// 执行："paid_at":        ts(detail.Payment.PaidAt),
			"paid_at":        ts(detail.Payment.PaidAt),
// 执行："paidAt":         ts(detail.Payment.PaidAt),
			"paidAt":         ts(detail.Payment.PaidAt),
// 结束当前代码块
		}
// 赋值/更新：view["payment"] = pay
		view["payment"] = pay
// 赋值/更新：data["payment"] = pay
		data["payment"] = pay
// 结束当前代码块
	}
// 条件判断：if detail.Refund != nil {
	if detail.Refund != nil {
// 赋值/更新：data["refund"] = refundView(*detail.Refund)
		data["refund"] = refundView(*detail.Refund)
// 结束当前代码块
	}
// 条件判断：if detail.Shipment != nil {
	if detail.Shipment != nil {
// 赋值/更新：ship := gin.H{
		ship := gin.H{
// 执行："carrier":     detail.Shipment.Carrier,
			"carrier":     detail.Shipment.Carrier,
// 执行："tracking_no": detail.Shipment.TrackingNo,
			"tracking_no": detail.Shipment.TrackingNo,
// 执行："trackingNo":  detail.Shipment.TrackingNo,
			"trackingNo":  detail.Shipment.TrackingNo,
// 执行："created_at":  ts(detail.Shipment.CreatedAt),
			"created_at":  ts(detail.Shipment.CreatedAt),
// 结束当前代码块
		}
// 赋值/更新：view["shipment"] = ship
		view["shipment"] = ship
// 赋值/更新：data["shipment"] = ship
		data["shipment"] = ship
// 结束当前代码块
	}
// 执行：a.okJSON(c, view, data)
	a.okJSON(c, view, data)
// 结束当前代码块
}

// CancelOrder POST /api/orders/{id}/cancel -> 待付款可取消（重复取消 4xx）

// 函数定义：func (a *CompatApi) CancelOrder(c *gin.Context) {
func (a *CompatApi) CancelOrder(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
	// 兼容层委派路径绕过了路由组上的 RequireRole，必须在此复核角色（I-02 越权防护）
	if claims.Role != model.RoleBuyer {
// 执行：a.failJSON(c, service.NewBizError(2002, "当前角色无权执行该操作"))
		a.failJSON(c, service.NewBizError(2002, "当前角色无权执行该操作"))
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：o, err := a.s.OrderService.FindByRef(c.Param("orderNo"))
	o, err := a.s.OrderService.FindByRef(c.Param("orderNo"))
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：updated, err := a.s.OrderService.Cancel(claims.UserID, o.OrderNo)
	updated, err := a.s.OrderService.Cancel(claims.UserID, o.OrderNo)
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：view := a.orderView(updated, nil, nil)
	view := a.orderView(updated, nil, nil)
// 执行：a.okJSON(c, view, gin.H{"order": view})
	a.okJSON(c, view, gin.H{"order": view})
// 结束当前代码块
}

// RefundApply POST /api/orders/{id}/refund-apply -> {refund_id}（仅 TO_SHIP 可申请）

// 函数定义：func (a *CompatApi) RefundApply(c *gin.Context) {
func (a *CompatApi) RefundApply(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
	// 兼容层委派/直接注册路径均绕过 RequireRole，必须在此复核角色（I-02 越权防护）
	if claims.Role != model.RoleBuyer {
// 执行：a.failJSON(c, service.NewBizError(2002, "当前角色无权执行该操作"))
		a.failJSON(c, service.NewBizError(2002, "当前角色无权执行该操作"))
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：o, err := a.s.OrderService.FindByRef(c.Param("orderNo"))
	o, err := a.s.OrderService.FindByRef(c.Param("orderNo"))
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：body := bindAny(c)
	body := bindAny(c)
// 赋值/更新：reason := strOf(body, "reason", "audit_reason")
	reason := strOf(body, "reason", "audit_reason")
// 条件判断：if reason == "" {
	if reason == "" {
// 赋值/更新：reason = "买家申请退款"
		reason = "买家申请退款"
// 结束当前代码块
	}
// 赋值/更新：rf, err := a.s.RefundService.Apply(claims.UserID, o.OrderNo, reason)
	rf, err := a.s.RefundService.Apply(claims.UserID, o.OrderNo, reason)
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：latest, _ := a.s.OrderService.Get(o.OrderNo)
	latest, _ := a.s.OrderService.Get(o.OrderNo)
// 赋值/更新：view := a.orderView(latest, &rf, nil)
	view := a.orderView(latest, &rf, nil)
// 执行：a.okJSON(c, gin.H{
	a.okJSON(c, gin.H{
// 执行："refund_id":    rf.RefundNo,
		"refund_id":    rf.RefundNo,
// 执行："id":           rf.RefundNo,
		"id":           rf.RefundNo,
// 执行："status":       refundStatusLabel(rf),
		"status":       refundStatusLabel(rf),
// 执行："order_no":     o.OrderNo,
		"order_no":     o.OrderNo,
// 执行："order_status": orderStatusLabel(latest.Status),
		"order_status": orderStatusLabel(latest.Status),
// 执行："refund":       refundView(rf),
		"refund":       refundView(rf),
// 结束当前代码块
	}, gin.H{
// 执行："refund":        refundView(rf),
		"refund":        refundView(rf),
// 执行："order":         view,
		"order":         view,
// 执行："refundPending": true,
		"refundPending": true,
// 结束当前代码块
	})
// 结束当前代码块
}

// ---------------- 运营接口 ----------------

// OpsProductUpdate PATCH /api/ops/products/{id} -> available 直接覆盖

// 函数定义：func (a *CompatApi) OpsProductUpdate(c *gin.Context) {
func (a *CompatApi) OpsProductUpdate(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：p, err := a.resolveProduct(c.Param("id"))
	p, err := a.resolveProduct(c.Param("id"))
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：body := bindAny(c)
	body := bindAny(c)
// 声明：var price *int64
	var price *int64
// 条件判断：if n, ok := intOf(body, "price_cents", "priceCents"); ok {
	if n, ok := intOf(body, "price_cents", "priceCents"); ok {
// 赋值/更新：price = &n
		price = &n
// 结束当前代码块
	}
// 声明：var onSale *bool
	var onSale *bool
// 条件判断：if b, ok := boolOf(body, "on_sale", "onSale"); ok {
	if b, ok := boolOf(body, "on_sale", "onSale"); ok {
// 赋值/更新：onSale = &b
		onSale = &b
// 结束当前代码块
	}
// 声明：var available *int
	var available *int
// 条件判断：if n, ok := intOf(body, "available"); ok {
	if n, ok := intOf(body, "available"); ok {
// 赋值/更新：v := int(n)
		v := int(n)
// 赋值/更新：available = &v
		available = &v
// 结束当前代码块
	}
// 条件判断：if price == nil && onSale == nil && available == nil {
	if price == nil && onSale == nil && available == nil {
// 执行：a.failJSON(c, service.NewBizError(1000, "未提供可更新字段"))
		a.failJSON(c, service.NewBizError(1000, "未提供可更新字段"))
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：updated := p
	updated := p
// 条件判断：if price != nil || onSale != nil {
	if price != nil || onSale != nil {
// 赋值/更新：updated, err = a.s.CatalogService.Update(claims.UserID, claims.Role, p.ID, onSale, price)
		updated, err = a.s.CatalogService.Update(claims.UserID, claims.Role, p.ID, onSale, price)
// 条件判断：if err != nil {
		if err != nil {
// 执行：a.failJSON(c, err)
			a.failJSON(c, err)
// 返回：return
			return
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if available != nil {
	if available != nil {
// 赋值/更新：updated, err = a.s.CatalogService.OverwriteAvailable(claims.UserID, claims.Role, p.ID, *available)
		updated, err = a.s.CatalogService.OverwriteAvailable(claims.UserID, claims.Role, p.ID, *available)
// 条件判断：if err != nil {
		if err != nil {
// 执行：a.failJSON(c, err)
			a.failJSON(c, err)
// 返回：return
			return
// 结束当前代码块
		}
// 结束当前代码块
	}
// 赋值/更新：view := productView(updated)
	view := productView(updated)
// 执行：a.okJSON(c, view, view)
	a.okJSON(c, view, view)
// 结束当前代码块
}

// OpsSummary GET /api/ops/summary -> {products:[{id,available}], paid_total, refunded_total}（分）

// 函数定义：func (a *CompatApi) OpsSummary(c *gin.Context) {
func (a *CompatApi) OpsSummary(c *gin.Context) {
	// 兼容层委派路径绕过了路由组上的 RequireRole，必须在此复核角色（E-10 越权防护）
	if claims := middleware.ClaimsOf(c); claims.Role != model.RoleOps {
// 执行：a.failJSON(c, service.NewBizError(2002, "当前角色无权执行该操作"))
		a.failJSON(c, service.NewBizError(2002, "当前角色无权执行该操作"))
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：sum, err := a.s.ReportService.Summary()
	sum, err := a.s.ReportService.Summary()
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：items := make([]gin.H, 0, len(sum.Products))
	items := make([]gin.H, 0, len(sum.Products))
// 循环处理：for _, it := range sum.Products {
	for _, it := range sum.Products {
// 赋值/更新：items = append(items, gin.H{
		items = append(items, gin.H{
// 执行："id":             it.SKU,
			"id":             it.SKU,
// 执行："sku":            it.SKU,
			"sku":            it.SKU,
// 执行："product_id":     it.SKU,
			"product_id":     it.SKU,
// 执行："db_id":          it.ProductID,
			"db_id":          it.ProductID,
// 执行："name":           it.Name,
			"name":           it.Name,
// 执行："price_cents":    it.PriceCents,
			"price_cents":    it.PriceCents,
// 执行："priceCents":     it.PriceCents,
			"priceCents":     it.PriceCents,
// 执行："available":      it.Available,
			"available":      it.Available,
// 执行："initial_stock":  it.InitialStock,
			"initial_stock":  it.InitialStock,
// 执行："initialStock":   it.InitialStock,
			"initialStock":   it.InitialStock,
// 执行："on_sale":        it.OnSale,
			"on_sale":        it.OnSale,
// 执行："onSale":         it.OnSale,
			"onSale":         it.OnSale,
// 执行："occupied":       it.Occupied,
			"occupied":       it.Occupied,
// 执行："shipped_qty":    it.ShippedQty,
			"shipped_qty":    it.ShippedQty,
// 执行："shippedQty":     it.ShippedQty,
			"shippedQty":     it.ShippedQty,
// 执行："paid_total":     it.PaidCents,
			"paid_total":     it.PaidCents,
// 执行："paidCents":      it.PaidCents,
			"paidCents":      it.PaidCents,
// 执行："abnormal_cents": it.AbnormalCents,
			"abnormal_cents": it.AbnormalCents,
// 执行："abnormalCents":  it.AbnormalCents,
			"abnormalCents":  it.AbnormalCents,
// 执行："refunded_total": it.RefundedCents,
			"refunded_total": it.RefundedCents,
// 执行："refundedCents":  it.RefundedCents,
			"refundedCents":  it.RefundedCents,
// 结束当前代码块
		})
// 结束当前代码块
	}
// 执行：c.JSON(http.StatusOK, gin.H{
	c.JSON(http.StatusOK, gin.H{
// 执行："code":               0,
		"code":               0,
// 执行："msg":                "ok",
		"msg":                "ok",
// 执行："products":           items,
		"products":           items,
// 执行："paid_total":         sum.TotalPaidCents,
		"paid_total":         sum.TotalPaidCents,
// 执行："abnormal_total":     sum.TotalAbnormalCents,
		"abnormal_total":     sum.TotalAbnormalCents,
// 执行："refunded_total":     sum.TotalRefundedCents,
		"refunded_total":     sum.TotalRefundedCents,
// 执行："order_count":        sum.OrderCount,
		"order_count":        sum.OrderCount,
// 执行："status_count":       sum.StatusCount,
		"status_count":       sum.StatusCount,
// 执行："data": gin.H{
		"data": gin.H{
// 执行："products":           items,
			"products":           items,
// 执行："totalPaidCents":     sum.TotalPaidCents,
			"totalPaidCents":     sum.TotalPaidCents,
// 执行："totalAbnormalCents": sum.TotalAbnormalCents,
			"totalAbnormalCents": sum.TotalAbnormalCents,
// 执行："totalRefundedCents": sum.TotalRefundedCents,
			"totalRefundedCents": sum.TotalRefundedCents,
// 执行："orderCount":         sum.OrderCount,
			"orderCount":         sum.OrderCount,
// 执行："statusCount":        sum.StatusCount,
			"statusCount":        sum.StatusCount,
// 结束当前代码块
		},
// 结束当前代码块
	})
// 结束当前代码块
}

// noticeStatusLabel 通知留痕结果 → 对外状态（成功记款单独以 SUCCESS 出现）。

// 函数定义：func noticeStatusLabel(n model.PaymentNotice) string {
func noticeStatusLabel(n model.PaymentNotice) string {
// 多路分支：switch n.Result {
	switch n.Result {
// 匹配分支：case model.NoticeAccepted:
	case model.NoticeAccepted:
// 返回：return "ACCEPTED"
		return "ACCEPTED"
// 匹配分支：case model.NoticeRejected:
	case model.NoticeRejected:
// 返回：return "REJECTED"
		return "REJECTED"
// 匹配分支：case model.NoticeDuplicated:
	case model.NoticeDuplicated:
// 返回：return "DUPLICATED"
		return "DUPLICATED"
// 匹配分支：case model.NoticeAbnormal:
	case model.NoticeAbnormal:
// 返回：return "ABNORMAL"
		return "ABNORMAL"
// 默认分支
	default:
// 返回：return strings.ToUpper(n.Result)
		return strings.ToUpper(n.Result)
// 结束当前代码块
	}
// 结束当前代码块
}

// OpsPayments GET /api/ops/payments?order_no=X -> 该订单全部支付/通知记录（含失败记录）

// 函数定义：func (a *CompatApi) OpsPayments(c *gin.Context) {
func (a *CompatApi) OpsPayments(c *gin.Context) {
// 赋值/更新：orderNo := strings.TrimSpace(c.Query("order_no"))
	orderNo := strings.TrimSpace(c.Query("order_no"))
// 条件判断：if orderNo == "" {
	if orderNo == "" {
// 赋值/更新：orderNo = strings.TrimSpace(c.Query("orderNo"))
		orderNo = strings.TrimSpace(c.Query("orderNo"))
// 结束当前代码块
	}
// 赋值/更新：items := make([]gin.H, 0, 4)
	items := make([]gin.H, 0, 4)
// 条件判断：if orderNo != "" {
	if orderNo != "" {
// 赋值/更新：payments, err := a.s.PaymentService.PaymentsByOrder(orderNo)
		payments, err := a.s.PaymentService.PaymentsByOrder(orderNo)
// 条件判断：if err != nil {
		if err != nil {
// 执行：a.failJSON(c, err)
			a.failJSON(c, err)
// 返回：return
			return
// 结束当前代码块
		}
// 循环处理：for _, p := range payments {
		for _, p := range payments {
// 赋值/更新：items = append(items, gin.H{
			items = append(items, gin.H{
// 执行："order_no":       p.OrderNo,
				"order_no":       p.OrderNo,
// 执行："orderNo":        p.OrderNo,
				"orderNo":        p.OrderNo,
// 执行："channel_txn_no": p.ChannelTxnID,
				"channel_txn_no": p.ChannelTxnID,
// 执行："channelTxnNo":   p.ChannelTxnID,
				"channelTxnNo":   p.ChannelTxnID,
// 执行："amount":         p.AmountCents,
				"amount":         p.AmountCents,
// 执行："amount_cents":   p.AmountCents,
				"amount_cents":   p.AmountCents,
// 执行："amountCents":    p.AmountCents,
				"amountCents":    p.AmountCents,
// 执行："status":         "SUCCESS",
				"status":         "SUCCESS",
// 执行："reason":         "渠道收款成功并已记款",
				"reason":         "渠道收款成功并已记款",
// 执行："source":         "payment",
				"source":         "payment",
// 执行："sign_ok":        true,
				"sign_ok":        true,
// 执行："paid_at":        ts(p.PaidAt),
				"paid_at":        ts(p.PaidAt),
// 执行："created_at":     ts(p.CreatedAt),
				"created_at":     ts(p.CreatedAt),
// 结束当前代码块
			})
// 结束当前代码块
		}
// 赋值/更新：notices, err := a.s.PaymentService.Notices(orderNo)
		notices, err := a.s.PaymentService.Notices(orderNo)
// 条件判断：if err != nil {
		if err != nil {
// 执行：a.failJSON(c, err)
			a.failJSON(c, err)
// 返回：return
			return
// 结束当前代码块
		}
// 循环处理：for _, n := range notices {
		for _, n := range notices {
// 赋值/更新：items = append(items, gin.H{
			items = append(items, gin.H{
// 执行："order_no":       n.OrderNo,
				"order_no":       n.OrderNo,
// 执行："orderNo":        n.OrderNo,
				"orderNo":        n.OrderNo,
// 执行："channel_txn_no": n.ChannelTxnID,
				"channel_txn_no": n.ChannelTxnID,
// 执行："channelTxnNo":   n.ChannelTxnID,
				"channelTxnNo":   n.ChannelTxnID,
// 执行："amount":         n.AmountCents,
				"amount":         n.AmountCents,
// 执行："amount_cents":   n.AmountCents,
				"amount_cents":   n.AmountCents,
// 执行："amountCents":    n.AmountCents,
				"amountCents":    n.AmountCents,
// 执行："status":         noticeStatusLabel(n),
				"status":         noticeStatusLabel(n),
// 执行："result":         n.Result,
				"result":         n.Result,
// 执行："reason":         n.Reason,
				"reason":         n.Reason,
// 执行："source":         "notice",
				"source":         "notice",
// 执行："sign_ok":        n.SignOK,
				"sign_ok":        n.SignOK,
// 执行："created_at":     ts(n.CreatedAt),
				"created_at":     ts(n.CreatedAt),
// 结束当前代码块
			})
// 结束当前代码块
		}
// 结束当前代码块
	}
// 执行：c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": items, "payments": items})
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": items, "payments": items})
// 结束当前代码块
}

// ---------------- 支付渠道模拟（脚本侧） ----------------

// PayChannelPay POST /api/pay/channel-pay -> {channel_txn_no, amount}

// 函数定义：func (a *CompatApi) PayChannelPay(c *gin.Context) {
func (a *CompatApi) PayChannelPay(c *gin.Context) {
// 赋值/更新：body := bindAny(c)
	body := bindAny(c)
// 赋值/更新：ref := strOf(body, "order_no", "orderNo")
	ref := strOf(body, "order_no", "orderNo")
// 条件判断：if ref == "" {
	if ref == "" {
// 执行：a.failJSON(c, service.NewBizError(1000, "order_no 不能为空"))
		a.failJSON(c, service.NewBizError(1000, "order_no 不能为空"))
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：o, err := a.s.OrderService.FindByRef(ref)
	o, err := a.s.OrderService.FindByRef(ref)
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：txn := fmt.Sprintf("TXN%s%06d", strings.TrimPrefix(o.OrderNo, "SO"), a.now().UnixNano()%1000000)
	txn := fmt.Sprintf("TXN%s%06d", strings.TrimPrefix(o.OrderNo, "SO"), a.now().UnixNano()%1000000)
// 执行：a.okJSON(c, gin.H{
	a.okJSON(c, gin.H{
// 执行："channel_txn_no": txn,
		"channel_txn_no": txn,
// 执行："channelTxnNo":   txn,
		"channelTxnNo":   txn,
// 执行："order_no":       o.OrderNo,
		"order_no":       o.OrderNo,
// 执行："amount":         o.AmountCents,
		"amount":         o.AmountCents,
// 执行："status":         "SUCCESS",
		"status":         "SUCCESS",
// 结束当前代码块
	}, gin.H{
// 执行："channelTxnNo": txn,
		"channelTxnNo": txn,
// 执行："orderNo":      o.OrderNo,
		"orderNo":      o.OrderNo,
// 执行："amountCents":  o.AmountCents,
		"amountCents":  o.AmountCents,
// 结束当前代码块
	})
// 结束当前代码块
}

// hmacHex 计算 HMAC-SHA256 十六进制小写摘要（渠道回调验签）。

// 函数定义：func hmacHex(secret string, raw []byte) string {
func hmacHex(secret string, raw []byte) string {
// 赋值/更新：mac := hmac.New(sha256.New, []byte(secret))
	mac := hmac.New(sha256.New, []byte(secret))
// 执行：mac.Write(raw)
	mac.Write(raw)
// 返回：return hex.EncodeToString(mac.Sum(nil))
	return hex.EncodeToString(mac.Sum(nil))
// 结束当前代码块
}

// PayNotify POST /api/pay/notify：按原始请求体字节验签，失败分支 4xx 且留痕可查。

// 函数定义：func (a *CompatApi) PayNotify(c *gin.Context) {
func (a *CompatApi) PayNotify(c *gin.Context) {
// 赋值/更新：raw, _ := c.GetRawData()
	raw, _ := c.GetRawData()
// 赋值/更新：body := map[string]interface{}{}
	body := map[string]interface{}{}
// 条件判断：if len(raw) > 0 {
	if len(raw) > 0 {
// 赋值/更新：_ = json.Unmarshal(raw, &body)
		_ = json.Unmarshal(raw, &body)
// 结束当前代码块
	}
// 赋值/更新：orderNo := strOf(body, "order_no", "orderNo")
	orderNo := strOf(body, "order_no", "orderNo")
// 赋值/更新：txnID := strOf(body, "channel_txn_no", "channelTxnNo", "channelTxnId", "channel_txn_id")
	txnID := strOf(body, "channel_txn_no", "channelTxnNo", "channelTxnId", "channel_txn_id")
// 赋值/更新：amount, _ := intOf(body, "amount", "amount_cents", "amountCents")
	amount, _ := intOf(body, "amount", "amount_cents", "amountCents")
// 赋值/更新：sig := strings.TrimSpace(c.GetHeader("X-Channel-Signature"))
	sig := strings.TrimSpace(c.GetHeader("X-Channel-Signature"))
// 赋值/更新：bodySign := strings.TrimSpace(strOf(body, "sign", "signature"))
	bodySign := strings.TrimSpace(strOf(body, "sign", "signature"))

// 赋值/更新：secret := a.cfg.Shop.PayNotifySecret
	secret := a.cfg.Shop.PayNotifySecret
// 赋值/更新：sign := "invalid-signature"
	sign := "invalid-signature"
// 多路分支：switch {
	switch {
// 匹配分支：case sig != "":
	case sig != "":
		// 契约通道：X-Channel-Signature = HMAC-SHA256(原始请求体)
		expect := hmacHex(secret, raw)
// 条件判断：if hmac.Equal([]byte(strings.ToLower(expect)), []byte(strings.ToLower(sig))) {
		if hmac.Equal([]byte(strings.ToLower(expect)), []byte(strings.ToLower(sig))) {
			// 验签通过：交给 service 按同一算法复核（内部使用同一密钥与字段顺序）
			sign = utils.PaySign(secret, orderNo, txnID, amount)
// 结束当前代码块
		}
// 匹配分支：case bodySign != "":
	case bodySign != "":
		// 内置演示页通道：请求体 sign = HMAC-SHA256("订单号|交易号|金额(分)")，
		// 与原生 /api/pay/notify 语义一致，避免演示页只能拿到"参数不合法"。
		if expect := utils.PaySign(secret, orderNo, txnID, amount); hmac.Equal([]byte(strings.ToLower(expect)), []byte(strings.ToLower(bodySign))) {
// 赋值/更新：sign = expect
			sign = expect
// 结束当前代码块
		}
// 默认分支
	default:
		// 既无签名头也无 body 签名：走"验签失败"分支并留痕
		sign = "invalid-signature"
// 结束当前代码块
	}

// 赋值/更新：res, err := a.s.PaymentService.Notify(orderNo, txnID, amount, sign, string(raw))
	res, err := a.s.PaymentService.Notify(orderNo, txnID, amount, sign, string(raw))
// 赋值/更新：payload := gin.H{
	payload := gin.H{
// 执行："order_no":       res.OrderNo,
		"order_no":       res.OrderNo,
// 执行："channel_txn_no": res.ChannelTxnID,
		"channel_txn_no": res.ChannelTxnID,
// 执行："amount":         res.AmountCents,
		"amount":         res.AmountCents,
// 执行："result":         res.Result,
		"result":         res.Result,
// 执行："reason":         res.Reason,
		"reason":         res.Reason,
// 执行："order_status":   orderStatusLabel(res.OrderStatus),
		"order_status":   orderStatusLabel(res.OrderStatus),
// 结束当前代码块
	}
// 多路分支：switch res.Result {
	switch res.Result {
// 匹配分支：case model.NoticeAccepted:
	case model.NoticeAccepted:
// 执行：c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "status": "SUCCESS", "data": payload, "notify": payload})
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "status": "SUCCESS", "data": payload, "notify": payload})
// 匹配分支：case model.NoticeAbnormal:
	case model.NoticeAbnormal:
		// 迟到支付：订单不复活，转入异常收款列表
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "abnormal", "status": "ABNORMAL", "data": payload, "notify": payload})
// 匹配分支：case model.NoticeDuplicated:
	case model.NoticeDuplicated:
// 赋值/更新：payload["code"] = 3003
		payload["code"] = 3003
// 赋值/更新：payload["message"] = res.Reason
		payload["message"] = res.Reason
// 执行：c.JSON(http.StatusConflict, payload)
		c.JSON(http.StatusConflict, payload)
// 默认分支
	default:
// 条件判断：if err != nil {
		if err != nil {
// 执行：a.failJSON(c, err)
			a.failJSON(c, err)
// 返回：return
			return
// 结束当前代码块
		}
// 赋值/更新：payload["code"] = 3002
		payload["code"] = 3002
// 赋值/更新：payload["message"] = res.Reason
		payload["message"] = res.Reason
// 执行：c.JSON(http.StatusBadRequest, payload)
		c.JSON(http.StatusBadRequest, payload)
// 结束当前代码块
	}
// 结束当前代码块
}

// PayChannelMode POST /api/pay/channel-mode -> 设定下一次渠道退款结果

// 函数定义：func (a *CompatApi) PayChannelMode(c *gin.Context) {
func (a *CompatApi) PayChannelMode(c *gin.Context) {
// 赋值/更新：body := bindAny(c)
	body := bindAny(c)
// 赋值/更新：raw := strOf(body, "refund_result", "refundResult", "mode", "refund_mode")
	raw := strOf(body, "refund_result", "refundResult", "mode", "refund_mode")
// 条件判断：if raw == "" {
	if raw == "" {
// 执行：a.failJSON(c, service.NewBizError(1000, "refund_result 不能为空"))
		a.failJSON(c, service.NewBizError(1000, "refund_result 不能为空"))
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：mode := "fail"
	mode := "fail"
// 多路分支：switch strings.ToUpper(strings.TrimSpace(raw)) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
// 匹配分支：case "SUCCESS", "OK", "PASS":
	case "SUCCESS", "OK", "PASS":
// 赋值/更新：mode = "success"
		mode = "success"
// 匹配分支：case "FAIL", "FAILED", "ERROR":
	case "FAIL", "FAILED", "ERROR":
// 赋值/更新：mode = "fail"
		mode = "fail"
// 默认分支
	default:
// 执行：a.failJSON(c, service.NewBizError(1000, "refund_result 仅支持 SUCCESS / FAIL"))
		a.failJSON(c, service.NewBizError(1000, "refund_result 仅支持 SUCCESS / FAIL"))
// 返回：return
		return
// 结束当前代码块
	}
// 执行：a.s.Channel.SetRefundMode(mode)
	a.s.Channel.SetRefundMode(mode)
// 赋值/更新：label := "FAIL"
	label := "FAIL"
// 条件判断：if mode == "success" {
	if mode == "success" {
// 赋值/更新：label = "SUCCESS"
		label = "SUCCESS"
// 结束当前代码块
	}
// 执行：a.okJSON(c, gin.H{
	a.okJSON(c, gin.H{
// 执行："refund_result": label,
		"refund_result": label,
// 执行："mode":          mode,
		"mode":          mode,
// 结束当前代码块
	}, gin.H{"refundResult": label, "refundMode": mode})
// 结束当前代码块
}

// ---------------- 仓库接口 ----------------

// WarehouseOrders GET /api/warehouse/orders -> 可发货（TO_SHIP）订单列表

// 函数定义：func (a *CompatApi) WarehouseOrders(c *gin.Context) {
func (a *CompatApi) WarehouseOrders(c *gin.Context) {
// 赋值/更新：list, err := a.s.ShipmentService.Shipable()
	list, err := a.s.ShipmentService.Shipable()
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：items := make([]gin.H, 0, len(list))
	items := make([]gin.H, 0, len(list))
// 循环处理：for _, o := range list {
	for _, o := range list {
// 赋值/更新：items = append(items, a.orderView(o, nil, nil))
		items = append(items, a.orderView(o, nil, nil))
// 结束当前代码块
	}
// 执行：c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": items, "orders": items})
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": items, "orders": items})
// 结束当前代码块
}

// WarehouseShip POST /api/warehouse/orders/{id}/ship -> 仅 TO_SHIP 可发货且幂等

// 函数定义：func (a *CompatApi) WarehouseShip(c *gin.Context) {
func (a *CompatApi) WarehouseShip(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：o, err := a.s.OrderService.FindByRef(c.Param("orderNo"))
	o, err := a.s.OrderService.FindByRef(c.Param("orderNo"))
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：body := bindAny(c)
	body := bindAny(c)
// 赋值/更新：carrier := strOf(body, "carrier", "express")
	carrier := strOf(body, "carrier", "express")
// 赋值/更新：trackingNo := strOf(body, "tracking_no", "trackingNo")
	trackingNo := strOf(body, "tracking_no", "trackingNo")
	// SHIP-02：承运方与运单号去空白后不能为空，空参必须拒绝且不产生发货记录
	// （增量续写【修改原因=R1】：删除原来的"顺丰"/自动单号兜底，兜底会绕过必填校验）。
	if carrier == "" || trackingNo == "" {
// 执行：a.failJSON(c, service.NewBizError(1000, "承运方和运单号不能为空"))
		a.failJSON(c, service.NewBizError(1000, "承运方和运单号不能为空"))
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：shipped, err := a.s.ShipmentService.Ship(claims.UserID, o.OrderNo, carrier, trackingNo)
	shipped, err := a.s.ShipmentService.Ship(claims.UserID, o.OrderNo, carrier, trackingNo)
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 声明：var ship *model.Shipment
	var ship *model.Shipment
// 条件判断：if s, e := a.s.ShipmentService.Get(o.OrderNo); e == nil {
	if s, e := a.s.ShipmentService.Get(o.OrderNo); e == nil {
// 赋值/更新：ship = s
		ship = s
// 结束当前代码块
	}
// 赋值/更新：view := a.orderView(shipped, nil, ship)
	view := a.orderView(shipped, nil, ship)
// 执行：a.okJSON(c, view, gin.H{"order": view})
	a.okJSON(c, view, gin.H{"order": view})
// 结束当前代码块
}

// ---------------- 客服接口 ----------------

// CSRefunds GET /api/cs/refunds -> 退款单列表

// 函数定义：func (a *CompatApi) CSRefunds(c *gin.Context) {
func (a *CompatApi) CSRefunds(c *gin.Context) {
	// 兼容层委派路径绕过了路由组上的 RequireRole，必须在此复核角色（E-10 越权防护）
	if claims := middleware.ClaimsOf(c); claims.Role != model.RoleCS {
// 执行：a.failJSON(c, service.NewBizError(2002, "当前角色无权执行该操作"))
		a.failJSON(c, service.NewBizError(2002, "当前角色无权执行该操作"))
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：status := strings.ToLower(strings.TrimSpace(c.Query("status")))
	status := strings.ToLower(strings.TrimSpace(c.Query("status")))
// 赋值/更新：list, err := a.s.RefundService.List(status)
	list, err := a.s.RefundService.List(status)
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：items := make([]gin.H, 0, len(list))
	items := make([]gin.H, 0, len(list))
// 循环处理：for _, r := range list {
	for _, r := range list {
// 赋值/更新：items = append(items, refundView(r))
		items = append(items, refundView(r))
// 结束当前代码块
	}
// 执行：c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": items, "refunds": items})
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": items, "refunds": items})
// 结束当前代码块
}

// CSRefundDecide POST /api/cs/refunds/{id}/decide -> 审核并触发渠道退款

// 函数定义：func (a *CompatApi) CSRefundDecide(c *gin.Context) {
func (a *CompatApi) CSRefundDecide(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：r, err := a.s.RefundService.FindByRef(c.Param("id"))
	r, err := a.s.RefundService.FindByRef(c.Param("id"))
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：body := bindAny(c)
	body := bindAny(c)
// 赋值/更新：approve, _ := boolOf(body, "approve")
	approve, _ := boolOf(body, "approve")
	// REF-03：客服同意/拒绝都必须填写非空原因，空原因拒绝且不改变审核状态
	// （增量续写【修改原因=R2】：删除自动填充默认文案的兜底）。
	reason := strOf(body, "reason", "audit_reason")
// 条件判断：if reason == "" {
	if reason == "" {
// 执行：a.failJSON(c, service.ErrReasonRequired)
		a.failJSON(c, service.ErrReasonRequired)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：updated, err := a.s.RefundService.Audit(claims.UserID, r.OrderNo, approve, reason)
	updated, err := a.s.RefundService.Audit(claims.UserID, r.OrderNo, approve, reason)
// 条件判断：if err != nil {
	if err != nil {
		// 渠道退款失败属于业务可预期结果：状态已落库（FAILED），契约要求 2xx 且可后续查询/重试
		if code, _ := service.BizCode(err); code == 3004 {
// 赋值/更新：view := refundView(updated)
			view := refundView(updated)
// 赋值/更新：latest, _ := a.s.OrderService.Get(r.OrderNo)
			latest, _ := a.s.OrderService.Get(r.OrderNo)
// 执行：a.okJSON(c, gin.H{
			a.okJSON(c, gin.H{
// 执行："refund_id":      view["id"],
				"refund_id":      view["id"],
// 执行："status":         view["status"],
				"status":         view["status"],
// 执行："channel_result": "FAIL",
				"channel_result": "FAIL",
// 执行："order_no":       r.OrderNo,
				"order_no":       r.OrderNo,
// 执行："order_status":   orderStatusLabel(latest.Status),
				"order_status":   orderStatusLabel(latest.Status),
// 执行："reason":         view["reason"],
				"reason":         view["reason"],
// 执行："fail_reason":    view["fail_reason"],
				"fail_reason":    view["fail_reason"],
// 执行："refund":         view,
				"refund":         view,
// 结束当前代码块
			}, gin.H{"refund": view, "channelResult": "FAIL"})
// 返回：return
			return
// 结束当前代码块
		}
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：view := refundView(updated)
	view := refundView(updated)
// 赋值/更新：latest, _ := a.s.OrderService.Get(r.OrderNo)
	latest, _ := a.s.OrderService.Get(r.OrderNo)
// 执行：a.okJSON(c, gin.H{
	a.okJSON(c, gin.H{
// 执行："refund_id":      view["id"],
		"refund_id":      view["id"],
// 执行："status":         view["status"],
		"status":         view["status"],
// 执行："channel_result": "SUCCESS",
		"channel_result": "SUCCESS",
// 执行："order_no":       r.OrderNo,
		"order_no":       r.OrderNo,
// 执行："order_status":   orderStatusLabel(latest.Status),
		"order_status":   orderStatusLabel(latest.Status),
// 执行："reason":         view["reason"],
		"reason":         view["reason"],
// 执行："refund":         view,
		"refund":         view,
// 结束当前代码块
	}, gin.H{"refund": view, "channelResult": "SUCCESS"})
// 结束当前代码块
}

// CSRefundRetry POST /api/cs/refunds/{id}/retry -> 对 FAILED 退款单重试

// 函数定义：func (a *CompatApi) CSRefundRetry(c *gin.Context) {
func (a *CompatApi) CSRefundRetry(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：r, err := a.s.RefundService.FindByRef(c.Param("id"))
	r, err := a.s.RefundService.FindByRef(c.Param("id"))
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：updated, err := a.s.RefundService.Retry(claims.UserID, r.RefundNo)
	updated, err := a.s.RefundService.Retry(claims.UserID, r.RefundNo)
// 条件判断：if err != nil {
	if err != nil {
// 条件判断：if code, _ := service.BizCode(err); code == 3004 {
		if code, _ := service.BizCode(err); code == 3004 {
// 赋值/更新：view := refundView(updated)
			view := refundView(updated)
// 赋值/更新：view["code"] = 3004
			view["code"] = 3004
// 赋值/更新：view["message"] = "渠道退款仍失败"
			view["message"] = "渠道退款仍失败"
// 执行：c.JSON(http.StatusConflict, view)
			c.JSON(http.StatusConflict, view)
// 返回：return
			return
// 结束当前代码块
		}
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：view := refundView(updated)
	view := refundView(updated)
// 赋值/更新：latest, _ := a.s.OrderService.Get(r.OrderNo)
	latest, _ := a.s.OrderService.Get(r.OrderNo)
// 执行：a.okJSON(c, gin.H{
	a.okJSON(c, gin.H{
// 执行："refund_id":    view["id"],
		"refund_id":    view["id"],
// 执行："status":       view["status"],
		"status":       view["status"],
// 执行："order_no":     r.OrderNo,
		"order_no":     r.OrderNo,
// 执行："order_status": orderStatusLabel(latest.Status),
		"order_status": orderStatusLabel(latest.Status),
// 执行："retry_count":  view["retry_count"],
		"retry_count":  view["retry_count"],
// 执行："refund":       view,
		"refund":       view,
// 结束当前代码块
	}, gin.H{"refund": view})
// 结束当前代码块
}

// abnormalStatusLabel 异常收款状态（契约 payment-exceptions）。

// 函数定义：func abnormalStatusLabel(status string) string {
func abnormalStatusLabel(status string) string {
// 条件判断：if status == model.AbnormalHandled {
	if status == model.AbnormalHandled {
// 返回：return "HANDLED"
		return "HANDLED"
// 结束当前代码块
	}
// 返回：return "PENDING"
	return "PENDING"
// 结束当前代码块
}

// CSExceptions GET /api/cs/payment-exceptions -> 异常支付列表

// 函数定义：func (a *CompatApi) CSExceptions(c *gin.Context) {
func (a *CompatApi) CSExceptions(c *gin.Context) {
// 赋值/更新：status := strings.ToLower(strings.TrimSpace(c.Query("status")))
	status := strings.ToLower(strings.TrimSpace(c.Query("status")))
// 赋值/更新：list, err := a.s.PaymentService.AbnormalList(status)
	list, err := a.s.PaymentService.AbnormalList(status)
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：items := make([]gin.H, 0, len(list))
	items := make([]gin.H, 0, len(list))
// 循环处理：for _, ab := range list {
	for _, ab := range list {
// 赋值/更新：items = append(items, gin.H{
		items = append(items, gin.H{
// 执行："id":             ab.ID,
			"id":             ab.ID,
// 执行："abnormal_id":    ab.ID,
			"abnormal_id":    ab.ID,
// 执行："order_no":       ab.OrderNo,
			"order_no":       ab.OrderNo,
// 执行："orderNo":        ab.OrderNo,
			"orderNo":        ab.OrderNo,
// 执行："channel_txn_no": ab.ChannelTxnID,
			"channel_txn_no": ab.ChannelTxnID,
// 执行："channelTxnNo":   ab.ChannelTxnID,
			"channelTxnNo":   ab.ChannelTxnID,
// 执行："amount":         ab.AmountCents,
			"amount":         ab.AmountCents,
// 执行："amount_cents":   ab.AmountCents,
			"amount_cents":   ab.AmountCents,
// 执行："amountCents":    ab.AmountCents,
			"amountCents":    ab.AmountCents,
// 执行："reason":         ab.Reason,
			"reason":         ab.Reason,
// 执行："status":         abnormalStatusLabel(ab.Status),
			"status":         abnormalStatusLabel(ab.Status),
// 执行："raw_status":     ab.Status,
			"raw_status":     ab.Status,
// 赋值/更新："handled":        ab.Status == model.AbnormalHandled,
			"handled":        ab.Status == model.AbnormalHandled,
// 执行："refund_id":      ab.RefundID,
			"refund_id":      ab.RefundID,
// 执行："created_at":     ts(ab.CreatedAt),
			"created_at":     ts(ab.CreatedAt),
// 执行："createdAt":      ts(ab.CreatedAt),
			"createdAt":      ts(ab.CreatedAt),
// 结束当前代码块
		})
// 结束当前代码块
	}
// 执行：c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": items, "exceptions": items})
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": items, "exceptions": items})
// 结束当前代码块
}

// CSExceptionRefund POST /api/cs/payment-exceptions/{id}/refund -> 对迟到支付发起退款（幂等）

// 函数定义：func (a *CompatApi) CSExceptionRefund(c *gin.Context) {
func (a *CompatApi) CSExceptionRefund(c *gin.Context) {
// 赋值/更新：claims := middleware.ClaimsOf(c)
	claims := middleware.ClaimsOf(c)
// 赋值/更新：idText := strings.TrimSpace(c.Param("id"))
	idText := strings.TrimSpace(c.Param("id"))
// 赋值/更新：id, err := strconv.ParseUint(idText, 10, 64)
	id, err := strconv.ParseUint(idText, 10, 64)
// 条件判断：if err != nil || id == 0 {
	if err != nil || id == 0 {
// 执行：a.failJSON(c, service.NewBizError(1000, "异常记录 id 不合法"))
		a.failJSON(c, service.NewBizError(1000, "异常记录 id 不合法"))
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：ab, err := a.s.PaymentService.AbnormalByID(id)
	ab, err := a.s.PaymentService.AbnormalByID(id)
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：body := bindAny(c)
	body := bindAny(c)
	// ABN-04：客服处理迟到支付必须填写非空原因，无原因不能处理
	// （增量续写【修改原因=R3】：删除自动填充默认原因的兜底）。
	reason := strOf(body, "reason")
// 条件判断：if reason == "" {
	if reason == "" {
// 执行：a.failJSON(c, service.ErrReasonRequired)
		a.failJSON(c, service.ErrReasonRequired)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：rf, err := a.s.RefundService.HandleAbnormal(claims.UserID, ab.ChannelTxnID, reason)
	rf, err := a.s.RefundService.HandleAbnormal(claims.UserID, ab.ChannelTxnID, reason)
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：view := refundView(rf)
	view := refundView(rf)
// 执行：a.okJSON(c, gin.H{
	a.okJSON(c, gin.H{
// 执行："refund_id":   rf.RefundNo,
		"refund_id":   rf.RefundNo,
// 执行："id":          rf.RefundNo,
		"id":          rf.RefundNo,
// 执行："status":      refundStatusLabel(rf),
		"status":      refundStatusLabel(rf),
// 执行："order_no":    ab.OrderNo,
		"order_no":    ab.OrderNo,
// 执行："abnormal_id": ab.ID,
		"abnormal_id": ab.ID,
// 执行："refund":      view,
		"refund":      view,
// 结束当前代码块
	}, gin.H{"refund": view, "abnormalId": ab.ID})
// 结束当前代码块
}

// ---------------- 测试钩子 ----------------

// diffNos 求 before - after（被关闭的订单号）。

// 函数定义：func diffNos(before, after []string) []string {
func diffNos(before, after []string) []string {
// 赋值/更新：keep := make(map[string]bool, len(after))
	keep := make(map[string]bool, len(after))
// 循环处理：for _, n := range after {
	for _, n := range after {
// 赋值/更新：keep[n] = true
		keep[n] = true
// 结束当前代码块
	}
// 赋值/更新：out := make([]string, 0, len(before))
	out := make([]string, 0, len(before))
// 循环处理：for _, n := range before {
	for _, n := range before {
// 条件判断：if !keep[n] {
		if !keep[n] {
// 赋值/更新：out = append(out, n)
			out = append(out, n)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return out
	return out
// 结束当前代码块
}

// TestExpireOrders POST /api/test/expire-orders -> 立即关闭超时未付订单并返回 {expired:[...]}
//
// 语义：先按服务端时间关闭真实超时订单；若本次没有订单到期（脚本刚创建订单的场景），
// 测试钩子会让最新一笔待付款订单立即到期后再关闭，等价于「时间快进」到该订单到期时刻。

// 函数定义：func (a *CompatApi) TestExpireOrders(c *gin.Context) {
func (a *CompatApi) TestExpireOrders(c *gin.Context) {
// 赋值/更新：before, err := a.s.OrderService.PendingOrderNos()
	before, err := a.s.OrderService.PendingOrderNos()
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 条件判断：if _, err := a.s.OrderService.CloseExpired(); err != nil {
	if _, err := a.s.OrderService.CloseExpired(); err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：after, err := a.s.OrderService.PendingOrderNos()
	after, err := a.s.OrderService.PendingOrderNos()
// 条件判断：if err != nil {
	if err != nil {
// 执行：a.failJSON(c, err)
		a.failJSON(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 赋值/更新：expired := diffNos(before, after)
	expired := diffNos(before, after)
// 条件判断：if len(expired) == 0 && len(before) > 0 {
	if len(expired) == 0 && len(before) > 0 {
// 条件判断：if err := a.s.OrderService.ForceExpire(before[0], a.now()); err != nil {
		if err := a.s.OrderService.ForceExpire(before[0], a.now()); err != nil {
// 执行：a.failJSON(c, err)
			a.failJSON(c, err)
// 返回：return
			return
// 结束当前代码块
		}
// 条件判断：if _, err := a.s.OrderService.CloseExpired(); err != nil {
		if _, err := a.s.OrderService.CloseExpired(); err != nil {
// 执行：a.failJSON(c, err)
			a.failJSON(c, err)
// 返回：return
			return
// 结束当前代码块
		}
// 赋值/更新：after2, err := a.s.OrderService.PendingOrderNos()
		after2, err := a.s.OrderService.PendingOrderNos()
// 条件判断：if err != nil {
		if err != nil {
// 执行：a.failJSON(c, err)
			a.failJSON(c, err)
// 返回：return
			return
// 结束当前代码块
		}
// 赋值/更新：expired = diffNos(before, after2)
		expired = diffNos(before, after2)
// 结束当前代码块
	}
// 执行：a.okJSON(c, gin.H{
	a.okJSON(c, gin.H{
// 执行："expired": expired,
		"expired": expired,
// 执行："closed":  len(expired),
		"closed":  len(expired),
// 结束当前代码块
	}, gin.H{"expired": expired, "closed": len(expired)})
// 结束当前代码块
}
