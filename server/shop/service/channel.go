// 包声明
package service

// 导入依赖包
import "sync"

// ChannelClient 支付渠道：受理退款并报告退款结果。
// 默认实现为可控内存模拟渠道，保证「首次失败、再次成功」可复现；
// 真实渠道（支付宝沙箱）见 alipay.go，经 shop.pay-channel=alipay 启用。

// 类型定义：type ChannelClient interface {
type ChannelClient interface {
	// Refund 受理退款：orderNo 用于真实渠道定位原交易（模拟渠道忽略）
	Refund(orderNo, refundNo string, amountCents int64) (bool, error)
// 执行：SetRefundMode(mode string)
	SetRefundMode(mode string)
// 执行：RefundMode() string
	RefundMode() string
// 执行：RefundCalls() int
	RefundCalls() int
// 结束当前代码块
}

// SimulatedChannel 模拟渠道实现

// 类型定义：type SimulatedChannel struct {
type SimulatedChannel struct {
// 执行：mu    sync.RWMutex
	mu    sync.RWMutex
// 执行：mode  string // success | fail
	mode  string // success | fail
// 执行：calls int
	calls int
// 结束当前代码块
}

// NewSimulatedChannel 创建模拟渠道；allow=true 时退款直接成功。

// 函数定义：func NewSimulatedChannel(allow bool) *SimulatedChannel {
func NewSimulatedChannel(allow bool) *SimulatedChannel {
// 赋值/更新：mode := "fail"
	mode := "fail"
// 条件判断：if allow {
	if allow {
// 赋值/更新：mode = "success"
		mode = "success"
// 结束当前代码块
	}
// 返回：return &SimulatedChannel{mode: mode}
	return &SimulatedChannel{mode: mode}
// 结束当前代码块
}

// Refund 受理退款（模拟渠道忽略订单号）

// 函数定义：func (c *SimulatedChannel) Refund(orderNo, refundNo string, amountCents int64) (bool, error) {
func (c *SimulatedChannel) Refund(orderNo, refundNo string, amountCents int64) (bool, error) {
// 执行：c.mu.Lock()
	c.mu.Lock()
// 异步/延迟执行：defer c.mu.Unlock()
	defer c.mu.Unlock()
// 执行：c.calls++
	c.calls++
// 返回：return c.mode == "success", nil
	return c.mode == "success", nil
// 结束当前代码块
}

// SetRefundMode 切换退款受理结果（演示用）

// 函数定义：func (c *SimulatedChannel) SetRefundMode(mode string) {
func (c *SimulatedChannel) SetRefundMode(mode string) {
// 执行：c.mu.Lock()
	c.mu.Lock()
// 异步/延迟执行：defer c.mu.Unlock()
	defer c.mu.Unlock()
// 条件判断：if mode == "fail" || mode == "success" {
	if mode == "fail" || mode == "success" {
// 赋值/更新：c.mode = mode
		c.mode = mode
// 结束当前代码块
	}
// 结束当前代码块
}

// RefundMode 当前退款受理结果

// 函数定义：func (c *SimulatedChannel) RefundMode() string {
func (c *SimulatedChannel) RefundMode() string {
// 执行：c.mu.RLock()
	c.mu.RLock()
// 异步/延迟执行：defer c.mu.RUnlock()
	defer c.mu.RUnlock()
// 返回：return c.mode
	return c.mode
// 结束当前代码块
}

// RefundCalls 累计受理次数（用于演示与测试断言）

// 函数定义：func (c *SimulatedChannel) RefundCalls() int {
func (c *SimulatedChannel) RefundCalls() int {
// 执行：c.mu.RLock()
	c.mu.RLock()
// 异步/延迟执行：defer c.mu.RUnlock()
	defer c.mu.RUnlock()
// 返回：return c.calls
	return c.calls
// 结束当前代码块
}
