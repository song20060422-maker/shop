// 包声明
package service

// 导入依赖包
import (
// 执行："strings"
	"strings"

// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/model/response"
	"shop/shop/model/response"
// 结束当前代码块
)

// ReportService 订单聚合视图与运营核查汇总

// 类型定义：type ReportService struct {
type ReportService struct {
// 执行：deps Deps
	deps Deps
// 结束当前代码块
}

// OrderDetail 聚合订单、收款、退款、发货、异常与操作轨迹，买家可据此知道钱和货去了哪里。

// 函数定义：func (s ReportService) OrderDetail(orderNo string, viewerID uint64, viewerRole string) (response.OrderDetail, ...
func (s ReportService) OrderDetail(orderNo string, viewerID uint64, viewerRole string) (response.OrderDetail, error) {
// 声明：var detail response.OrderDetail
	var detail response.OrderDetail
// 赋值/更新：orderNo = strings.TrimSpace(orderNo)
	orderNo = strings.TrimSpace(orderNo)
// 条件判断：if err := s.deps.DB.Where("order_no = ?", orderNo).First(&detail.Order).Error; err != nil {
	if err := s.deps.DB.Where("order_no = ?", orderNo).First(&detail.Order).Error; err != nil {
// 返回：return detail, ErrOrderNotFound
		return detail, ErrOrderNotFound
// 结束当前代码块
	}
// 条件判断：if viewerRole == model.RoleBuyer && detail.Order.BuyerID != viewerID {
	if viewerRole == model.RoleBuyer && detail.Order.BuyerID != viewerID {
// 返回：return response.OrderDetail{}, ErrForbidden
		return response.OrderDetail{}, ErrForbidden
// 结束当前代码块
	}

// 声明：var pay model.Payment
	var pay model.Payment
// 条件判断：if err := s.deps.DB.Where("order_no = ?", orderNo).First(&pay).Error; err == nil {
	if err := s.deps.DB.Where("order_no = ?", orderNo).First(&pay).Error; err == nil {
// 赋值/更新：detail.Payment = &pay
		detail.Payment = &pay
// 结束当前代码块
	}
// 声明：var refund model.Refund
	var refund model.Refund
// 条件判断：if err := s.deps.DB.Where("order_no = ?", orderNo).First(&refund).Error; err == nil {
	if err := s.deps.DB.Where("order_no = ?", orderNo).First(&refund).Error; err == nil {
// 赋值/更新：detail.Refund = &refund
		detail.Refund = &refund
// 结束当前代码块
	}
// 声明：var shipment model.Shipment
	var shipment model.Shipment
// 条件判断：if err := s.deps.DB.Where("order_no = ?", orderNo).First(&shipment).Error; err == nil {
	if err := s.deps.DB.Where("order_no = ?", orderNo).First(&shipment).Error; err == nil {
// 赋值/更新：detail.Shipment = &shipment
		detail.Shipment = &shipment
// 结束当前代码块
	}
// 声明：var product model.Product
	var product model.Product
// 条件判断：if err := s.deps.DB.Where("id = ?", detail.Order.ProductID).First(&product).Error; err == nil {
	if err := s.deps.DB.Where("id = ?", detail.Order.ProductID).First(&product).Error; err == nil {
// 赋值/更新：detail.ProductInfo = &product
		detail.ProductInfo = &product
// 结束当前代码块
	}
// 赋值/更新：_ = s.deps.DB.Where("order_no = ?", orderNo).Order("id asc").Find(&detail.Abnormal).Error
	_ = s.deps.DB.Where("order_no = ?", orderNo).Order("id asc").Find(&detail.Abnormal).Error
// 赋值/更新：_ = s.deps.DB.Where("order_no = ?", orderNo).Order("id asc").Find(&detail.Notices).Error
	_ = s.deps.DB.Where("order_no = ?", orderNo).Order("id asc").Find(&detail.Notices).Error
// 赋值/更新：_ = s.deps.DB.Where("order_no = ?", orderNo).Order("id asc").Find(&detail.Timeline).Error
	_ = s.deps.DB.Where("order_no = ?", orderNo).Order("id asc").Find(&detail.Timeline).Error

// 赋值/更新：detail.ShipableNow = detail.Order.Status == model.OrderStatusPaid
	detail.ShipableNow = detail.Order.Status == model.OrderStatusPaid
// 赋值/更新：detail.RefundPending = detail.Order.Status == model.OrderStatusRefundRequested ||
	detail.RefundPending = detail.Order.Status == model.OrderStatusRefundRequested ||
// 赋值/更新：detail.Order.Status == model.OrderStatusRefundFailed
		detail.Order.Status == model.OrderStatusRefundFailed
// 返回：return detail, nil
	return detail, nil
// 结束当前代码块
}

// Summary 运营核查：商品可售数量、已收款、已退款与订单状态分布。

// 函数定义：func (s ReportService) Summary() (response.TradeSummary, error) {
func (s ReportService) Summary() (response.TradeSummary, error) {
// 赋值/更新：out := response.TradeSummary{StatusCount: map[string]int64{}}
	out := response.TradeSummary{StatusCount: map[string]int64{}}

// 声明：var products []model.Product
	var products []model.Product
// 条件判断：if err := s.deps.DB.Order("id asc").Find(&products).Error; err != nil {
	if err := s.deps.DB.Order("id asc").Find(&products).Error; err != nil {
// 返回：return out, err
		return out, err
// 结束当前代码块
	}

// 类型定义：type aggRow struct {
	type aggRow struct {
// 执行：ProductID uint64
		ProductID uint64
// 执行：Qty       int64
		Qty       int64
// 执行：Cents     int64
		Cents     int64
// 结束当前代码块
	}
// 声明：var shipRows []aggRow
	var shipRows []aggRow
// 条件判断：if err := s.deps.DB.Model(&model.Order{}).
	if err := s.deps.DB.Model(&model.Order{}).
// 执行：Select("product_id AS product_id, COALESCE(SUM(qty),0) AS qty").
		Select("product_id AS product_id, COALESCE(SUM(qty),0) AS qty").
// 执行：Where("status = ?", model.OrderStatusShipped).
		Where("status = ?", model.OrderStatusShipped).
// 执行：Group("product_id").Scan(&shipRows).Error; err != nil {
		Group("product_id").Scan(&shipRows).Error; err != nil {
// 返回：return out, err
		return out, err
// 结束当前代码块
	}
	// 占用量口径改为按订单状态直查（增量续写【修改原因=R5】）：
	// 原公式 initial_stock-available-shipped 在运营补货/减库后失衡；
	// 状态说明 4.2 定义占用量 = 尚未关闭且未发货订单的数量合计。
	var occRows []aggRow
// 条件判断：if err := s.deps.DB.Model(&model.Order{}).
	if err := s.deps.DB.Model(&model.Order{}).
// 执行：Select("product_id AS product_id, COALESCE(SUM(qty),0) AS qty").
		Select("product_id AS product_id, COALESCE(SUM(qty),0) AS qty").
// 执行：Where("status IN ?", []string{
		Where("status IN ?", []string{
// 执行：model.OrderStatusPendingPayment, model.OrderStatusPaid,
			model.OrderStatusPendingPayment, model.OrderStatusPaid,
// 执行：model.OrderStatusRefundRequested, model.OrderStatusRefundFailed,
			model.OrderStatusRefundRequested, model.OrderStatusRefundFailed,
// 结束当前代码块
		}).
// 执行：Group("product_id").Scan(&occRows).Error; err != nil {
		Group("product_id").Scan(&occRows).Error; err != nil {
// 返回：return out, err
		return out, err
// 结束当前代码块
	}
// 声明：var paidRows []aggRow
	var paidRows []aggRow
// 条件判断：if err := s.deps.DB.Table("payments").
	if err := s.deps.DB.Table("payments").
// 执行：Select("orders.product_id AS product_id, COALESCE(SUM(payments.amount_cents),0) AS cents").
		Select("orders.product_id AS product_id, COALESCE(SUM(payments.amount_cents),0) AS cents").
// 执行：Joins("JOIN orders ON orders.order_no = payments.order_no").
		Joins("JOIN orders ON orders.order_no = payments.order_no").
// 执行：Group("orders.product_id").Scan(&paidRows).Error; err != nil {
		Group("orders.product_id").Scan(&paidRows).Error; err != nil {
// 返回：return out, err
		return out, err
// 结束当前代码块
	}
// 声明：var refundRows []aggRow
	var refundRows []aggRow
// 条件判断：if err := s.deps.DB.Table("refunds").
	if err := s.deps.DB.Table("refunds").
// 执行：Select("orders.product_id AS product_id, COALESCE(SUM(refunds.amount_cents),0) AS cents").
		Select("orders.product_id AS product_id, COALESCE(SUM(refunds.amount_cents),0) AS cents").
// 执行：Joins("JOIN orders ON orders.order_no = refunds.order_no").
		Joins("JOIN orders ON orders.order_no = refunds.order_no").
// 执行：Where("refunds.channel_status = ?", model.ChannelSuccess).
		Where("refunds.channel_status = ?", model.ChannelSuccess).
// 执行：Group("orders.product_id").Scan(&refundRows).Error; err != nil {
		Group("orders.product_id").Scan(&refundRows).Error; err != nil {
// 返回：return out, err
		return out, err
// 结束当前代码块
	}
	// 迟到支付：渠道已收款但订单已关闭，不会写入 payments；金额按异常收款记录聚合，
	// 独立列示不并入有效收款，避免"渠道已收款而商城报表查不到"的盲区（题面规则 8）。
	var abnormalRows []aggRow
// 条件判断：if err := s.deps.DB.Table("abnormal_payments").
	if err := s.deps.DB.Table("abnormal_payments").
// 执行：Select("orders.product_id AS product_id, COALESCE(SUM(abnormal_payments.amount_cents),0) AS cents").
		Select("orders.product_id AS product_id, COALESCE(SUM(abnormal_payments.amount_cents),0) AS cents").
// 执行：Joins("JOIN orders ON orders.order_no = abnormal_payments.order_no").
		Joins("JOIN orders ON orders.order_no = abnormal_payments.order_no").
// 执行：Group("orders.product_id").Scan(&abnormalRows).Error; err != nil {
		Group("orders.product_id").Scan(&abnormalRows).Error; err != nil {
// 返回：return out, err
		return out, err
// 结束当前代码块
	}
// 声明：var abnormalTotal int64
	var abnormalTotal int64
// 条件判断：if err := s.deps.DB.Model(&model.AbnormalPayment{}).
	if err := s.deps.DB.Model(&model.AbnormalPayment{}).
// 执行：Select("COALESCE(SUM(amount_cents),0)").Scan(&abnormalTotal).Error; err != nil {
		Select("COALESCE(SUM(amount_cents),0)").Scan(&abnormalTotal).Error; err != nil {
// 返回：return out, err
		return out, err
// 结束当前代码块
	}

// 赋值/更新：shippedOf := func(id uint64) int64 {
	shippedOf := func(id uint64) int64 {
// 循环处理：for _, r := range shipRows {
		for _, r := range shipRows {
// 条件判断：if r.ProductID == id {
			if r.ProductID == id {
// 返回：return r.Qty
				return r.Qty
// 结束当前代码块
			}
// 结束当前代码块
		}
// 返回：return 0
		return 0
// 结束当前代码块
	}
// 赋值/更新：occupiedOf := func(id uint64) int64 {
	occupiedOf := func(id uint64) int64 {
// 循环处理：for _, r := range occRows {
		for _, r := range occRows {
// 条件判断：if r.ProductID == id {
			if r.ProductID == id {
// 返回：return r.Qty
				return r.Qty
// 结束当前代码块
			}
// 结束当前代码块
		}
// 返回：return 0
		return 0
// 结束当前代码块
	}
	// 有效收款只统计成功记款 payments；异常渠道收款单独列示（增量续写【修改原因=R5】，
	// OPS-03/BR-14：不得把迟到支付等异常金额并入有效收款，避免隐藏）。
	paidOf := func(id uint64) int64 {
// 循环处理：for _, r := range paidRows {
		for _, r := range paidRows {
// 条件判断：if r.ProductID == id {
			if r.ProductID == id {
// 返回：return r.Cents
				return r.Cents
// 结束当前代码块
			}
// 结束当前代码块
		}
// 返回：return 0
		return 0
// 结束当前代码块
	}
// 赋值/更新：abnormalOf := func(id uint64) int64 {
	abnormalOf := func(id uint64) int64 {
// 循环处理：for _, r := range abnormalRows {
		for _, r := range abnormalRows {
// 条件判断：if r.ProductID == id {
			if r.ProductID == id {
// 返回：return r.Cents
				return r.Cents
// 结束当前代码块
			}
// 结束当前代码块
		}
// 返回：return 0
		return 0
// 结束当前代码块
	}
// 赋值/更新：refundedOf := func(id uint64) int64 {
	refundedOf := func(id uint64) int64 {
// 循环处理：for _, r := range refundRows {
		for _, r := range refundRows {
// 条件判断：if r.ProductID == id {
			if r.ProductID == id {
// 返回：return r.Cents
				return r.Cents
// 结束当前代码块
			}
// 结束当前代码块
		}
// 返回：return 0
		return 0
// 结束当前代码块
	}

// 循环处理：for _, p := range products {
	for _, p := range products {
// 赋值/更新：shipped := shippedOf(p.ID)
		shipped := shippedOf(p.ID)
// 赋值/更新：occupied := int(occupiedOf(p.ID))
		occupied := int(occupiedOf(p.ID))
// 条件判断：if occupied < 0 {
		if occupied < 0 {
// 赋值/更新：occupied = 0
			occupied = 0
// 结束当前代码块
		}
// 赋值/更新：out.Products = append(out.Products, response.ProductStat{
		out.Products = append(out.Products, response.ProductStat{
// 执行：ProductID:     p.ID,
			ProductID:     p.ID,
// 执行：SKU:           p.SKU,
			SKU:           p.SKU,
// 执行：Name:          p.Name,
			Name:          p.Name,
// 执行：PriceCents:    p.PriceCents,
			PriceCents:    p.PriceCents,
// 执行：Available:     p.Available,
			Available:     p.Available,
// 执行：InitialStock:  p.InitialStock,
			InitialStock:  p.InitialStock,
// 执行：OnSale:        p.OnSale,
			OnSale:        p.OnSale,
// 执行：Occupied:      occupied,
			Occupied:      occupied,
// 执行：ShippedQty:    int(shipped),
			ShippedQty:    int(shipped),
// 执行：PaidCents:     paidOf(p.ID),
			PaidCents:     paidOf(p.ID),
// 执行：AbnormalCents: abnormalOf(p.ID),
			AbnormalCents: abnormalOf(p.ID),
// 执行：RefundedCents: refundedOf(p.ID),
			RefundedCents: refundedOf(p.ID),
// 结束当前代码块
		})
// 结束当前代码块
	}

// 条件判断：if err := s.deps.DB.Model(&model.Payment{}).
	if err := s.deps.DB.Model(&model.Payment{}).
// 执行：Select("COALESCE(SUM(amount_cents),0)").Scan(&out.TotalPaidCents).Error; err != nil {
		Select("COALESCE(SUM(amount_cents),0)").Scan(&out.TotalPaidCents).Error; err != nil {
// 返回：return out, err
		return out, err
// 结束当前代码块
	}
	// 异常渠道收款（迟到支付等）单独合计，不并入有效收款（增量续写【修改原因=R5】）
	out.TotalAbnormalCents = abnormalTotal
// 条件判断：if err := s.deps.DB.Model(&model.Refund{}).
	if err := s.deps.DB.Model(&model.Refund{}).
// 执行：Select("COALESCE(SUM(amount_cents),0)").
		Select("COALESCE(SUM(amount_cents),0)").
// 执行：Where("channel_status = ?", model.ChannelSuccess).Scan(&out.TotalRefundedCents).Error; err != nil {
		Where("channel_status = ?", model.ChannelSuccess).Scan(&out.TotalRefundedCents).Error; err != nil {
// 返回：return out, err
		return out, err
// 结束当前代码块
	}
// 条件判断：if err := s.deps.DB.Model(&model.Order{}).Count(&out.OrderCount).Error; err != nil {
	if err := s.deps.DB.Model(&model.Order{}).Count(&out.OrderCount).Error; err != nil {
// 返回：return out, err
		return out, err
// 结束当前代码块
	}

// 类型定义：type statusRow struct {
	type statusRow struct {
// 执行：Status string
		Status string
// 执行：Total  int64
		Total  int64
// 结束当前代码块
	}
// 声明：var rows []statusRow
	var rows []statusRow
// 条件判断：if err := s.deps.DB.Model(&model.Order{}).Select("status AS status, COUNT(*) AS total").
	if err := s.deps.DB.Model(&model.Order{}).Select("status AS status, COUNT(*) AS total").
// 执行：Group("status").Scan(&rows).Error; err != nil {
		Group("status").Scan(&rows).Error; err != nil {
// 返回：return out, err
		return out, err
// 结束当前代码块
	}
// 循环处理：for _, r := range rows {
	for _, r := range rows {
// 赋值/更新：out.StatusCount[r.Status] = r.Total
		out.StatusCount[r.Status] = r.Total
// 结束当前代码块
	}
// 返回：return out, nil
	return out, nil
// 结束当前代码块
}

// Trail 订单全链路操作留痕（只读，异常不能被改状态抹掉）

// 函数定义：func (s ReportService) Trail(orderNo string) ([]model.OperationLog, error) {
func (s ReportService) Trail(orderNo string) ([]model.OperationLog, error) {
// 声明：var list []model.OperationLog
	var list []model.OperationLog
// 赋值/更新：q := s.deps.DB
	q := s.deps.DB
// 条件判断：if strings.TrimSpace(orderNo) != "" {
	if strings.TrimSpace(orderNo) != "" {
// 赋值/更新：q = q.Where("order_no = ?", strings.TrimSpace(orderNo))
		q = q.Where("order_no = ?", strings.TrimSpace(orderNo))
// 结束当前代码块
	}
// 赋值/更新：err := q.Order("id asc").Limit(500).Find(&list).Error
	err := q.Order("id asc").Limit(500).Find(&list).Error
// 返回：return list, err
	return list, err
// 结束当前代码块
}

// Ledger 库存流水（核对占用/释放/恢复/出库是否守恒）

// 函数定义：func (s ReportService) Ledger(productID uint64) ([]model.InventoryLedger, error) {
func (s ReportService) Ledger(productID uint64) ([]model.InventoryLedger, error) {
// 声明：var list []model.InventoryLedger
	var list []model.InventoryLedger
// 赋值/更新：q := s.deps.DB
	q := s.deps.DB
// 条件判断：if productID > 0 {
	if productID > 0 {
// 赋值/更新：q = q.Where("product_id = ?", productID)
		q = q.Where("product_id = ?", productID)
// 结束当前代码块
	}
// 赋值/更新：err := q.Order("id asc").Limit(500).Find(&list).Error
	err := q.Order("id asc").Limit(500).Find(&list).Error
// 返回：return list, err
	return list, err
// 结束当前代码块
}
