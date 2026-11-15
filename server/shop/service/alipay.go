// 包声明
package service

// 导入依赖包
import (
// 执行："context"
	"context"
// 执行："fmt"
	"fmt"
// 执行："net/url"
	"net/url"
// 执行："strconv"
	"strconv"
// 执行："strings"
	"strings"
// 执行："sync"
	"sync"

// 执行："github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/alipay/v3"

// 执行："shop/shop/config"
	"shop/shop/config"
// 结束当前代码块
)

// 支付宝真实渠道（沙箱演示扩展）：电脑网站支付（alipay.trade.page.pay）+ 退款（alipay.trade.refund）。
// 启用条件：shop.pay-channel=alipay 且配置完整；模拟渠道仍是默认，保证验收口径不变。
// 通知推进与模拟渠道共用同一套状态机/幂等/异常路径（payment.go applyPayment）。

// NotifyVerifier 支付宝渠道的附加能力：校验异步通知签名并提取业务字段。

// 类型定义：type NotifyVerifier interface {
type NotifyVerifier interface {
// 执行：VerifyNotify(ctx context.Context, form url.Values) (NotifyPayload, error)
	VerifyNotify(ctx context.Context, form url.Values) (NotifyPayload, error)
// 结束当前代码块
}

// NotifyPayload 异步通知验签后提取的字段（已把元金额换算为分）

// 类型定义：type NotifyPayload struct {
type NotifyPayload struct {
// 执行：OrderNo     string // out_trade_no
	OrderNo     string // out_trade_no
// 执行：Txn         string // trade_no
	Txn         string // trade_no
// 执行：AmountCents int64  // total_amount 换算为分
	AmountCents int64  // total_amount 换算为分
// 执行：TradeStatus string // TRADE_SUCCESS | TRADE_FINISHED | ...
	TradeStatus string // TRADE_SUCCESS | TRADE_FINISHED | ...
// 执行：AppID       string // 通知中的 app_id，须与本渠道一致
	AppID       string // 通知中的 app_id，须与本渠道一致
// 结束当前代码块
}

// AlipayChannel 支付宝渠道实现（实现 ChannelClient + NotifyVerifier）。

// 类型定义：type AlipayChannel struct {
type AlipayChannel struct {
// 执行：client *alipay.Client
	client *alipay.Client
// 执行：cfg    config.Alipay
	cfg    config.Alipay
// 执行：mu     sync.RWMutex
	mu     sync.RWMutex
// 执行：calls  int
	calls  int
// 结束当前代码块
}

// NewAlipayChannel 构建支付宝渠道；沙箱网关由 SDK 按 production=false 自动选择。

// 函数定义：func NewAlipayChannel(cfg config.Alipay) (*AlipayChannel, error) {
func NewAlipayChannel(cfg config.Alipay) (*AlipayChannel, error) {
// 条件判断：if strings.TrimSpace(cfg.AppID) == "" || strings.TrimSpace(cfg.PrivateKey) == "" || strings.TrimSpace(cfg.Alip...
	if strings.TrimSpace(cfg.AppID) == "" || strings.TrimSpace(cfg.PrivateKey) == "" || strings.TrimSpace(cfg.AlipayPublicKey) == "" {
// 返回：return nil, fmt.Errorf("支付宝渠道配置不完整（app-id / private-key / alipay-public-key）")
		return nil, fmt.Errorf("支付宝渠道配置不完整（app-id / private-key / alipay-public-key）")
// 结束当前代码块
	}
// 赋值/更新：client, err := alipay.New(cfg.AppID, cfg.PrivateKey, false)
	client, err := alipay.New(cfg.AppID, cfg.PrivateKey, false)
// 条件判断：if err != nil {
	if err != nil {
// 返回：return nil, fmt.Errorf("初始化支付宝客户端失败: %w", err)
		return nil, fmt.Errorf("初始化支付宝客户端失败: %w", err)
// 结束当前代码块
	}
// 条件判断：if err := client.LoadAliPayPublicKey(cfg.AlipayPublicKey); err != nil {
	if err := client.LoadAliPayPublicKey(cfg.AlipayPublicKey); err != nil {
// 返回：return nil, fmt.Errorf("加载支付宝公钥失败: %w", err)
		return nil, fmt.Errorf("加载支付宝公钥失败: %w", err)
// 结束当前代码块
	}
// 返回：return &AlipayChannel{client: client, cfg: cfg}, nil
	return &AlipayChannel{client: client, cfg: cfg}, nil
// 结束当前代码块
}

// PagePayURL 生成电脑网站支付的收银台跳转地址（GET，直接 302/window.open 即可）。
// notifyURL / returnURL 由调用方按实际访问域名动态传入（数据不写死）；
// 传空则回落 config.shop.alipay 的静态配置（兼容显式指定场景）。

// 函数定义：func (c *AlipayChannel) PagePayURL(orderNo string, amountCents int64, subject, notifyURL, returnURL string) (s...
func (c *AlipayChannel) PagePayURL(orderNo string, amountCents int64, subject, notifyURL, returnURL string) (string, error) {
// 声明：var p alipay.TradePagePay
	var p alipay.TradePagePay
// 赋值/更新：p.OutTradeNo = orderNo
	p.OutTradeNo = orderNo
// 赋值/更新：p.TotalAmount = centsToYuan(amountCents)
	p.TotalAmount = centsToYuan(amountCents)
// 赋值/更新：p.Subject = subject
	p.Subject = subject
// 赋值/更新：p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
// 赋值/更新：p.NotifyURL = notifyURL
	p.NotifyURL = notifyURL
// 条件判断：if strings.TrimSpace(p.NotifyURL) == "" {
	if strings.TrimSpace(p.NotifyURL) == "" {
// 赋值/更新：p.NotifyURL = c.cfg.NotifyURL
		p.NotifyURL = c.cfg.NotifyURL
// 结束当前代码块
	}
// 赋值/更新：p.ReturnURL = returnURL
	p.ReturnURL = returnURL
// 条件判断：if strings.TrimSpace(p.ReturnURL) == "" {
	if strings.TrimSpace(p.ReturnURL) == "" {
// 赋值/更新：p.ReturnURL = c.cfg.ReturnURL
		p.ReturnURL = c.cfg.ReturnURL
// 结束当前代码块
	}
// 赋值/更新：u, err := c.client.TradePagePay(p)
	u, err := c.client.TradePagePay(p)
// 条件判断：if err != nil {
	if err != nil {
// 返回：return "", fmt.Errorf("构造支付宝收银台地址失败: %w", err)
		return "", fmt.Errorf("构造支付宝收银台地址失败: %w", err)
// 结束当前代码块
	}
// 返回：return u.String(), nil
	return u.String(), nil
// 结束当前代码块
}

// VerifyNotify 校验异步通知签名（RSA2，支付宝公钥）并提取业务字段。

// 函数定义：func (c *AlipayChannel) VerifyNotify(ctx context.Context, form url.Values) (NotifyPayload, error) {
func (c *AlipayChannel) VerifyNotify(ctx context.Context, form url.Values) (NotifyPayload, error) {
// 赋值/更新：n, err := c.client.DecodeNotification(ctx, form)
	n, err := c.client.DecodeNotification(ctx, form)
// 条件判断：if err != nil {
	if err != nil {
// 返回：return NotifyPayload{}, err
		return NotifyPayload{}, err
// 结束当前代码块
	}
// 条件判断：if n.AppId != c.cfg.AppID {
	if n.AppId != c.cfg.AppID {
// 返回：return NotifyPayload{}, fmt.Errorf("通知 app_id 不匹配：期望 %s 实际 %s", c.cfg.AppID, n.AppId)
		return NotifyPayload{}, fmt.Errorf("通知 app_id 不匹配：期望 %s 实际 %s", c.cfg.AppID, n.AppId)
// 结束当前代码块
	}
// 赋值/更新：cents, err := yuanToCents(n.TotalAmount)
	cents, err := yuanToCents(n.TotalAmount)
// 条件判断：if err != nil {
	if err != nil {
// 返回：return NotifyPayload{}, fmt.Errorf("通知金额非法 %q: %w", n.TotalAmount, err)
		return NotifyPayload{}, fmt.Errorf("通知金额非法 %q: %w", n.TotalAmount, err)
// 结束当前代码块
	}
// 返回：return NotifyPayload{
	return NotifyPayload{
// 执行：OrderNo:     n.OutTradeNo,
		OrderNo:     n.OutTradeNo,
// 执行：Txn:         n.TradeNo,
		Txn:         n.TradeNo,
// 执行：AmountCents: cents,
		AmountCents: cents,
// 执行：TradeStatus: string(n.TradeStatus),
		TradeStatus: string(n.TradeStatus),
// 执行：AppID:       n.AppId,
		AppID:       n.AppId,
// 结束当前代码块
	}, nil
// 结束当前代码块
}

// Refund 调用 alipay.trade.refund（同步接口）：code=10000 视为受理成功（含重复请求的幂等返回）。

// 函数定义：func (c *AlipayChannel) Refund(orderNo, refundNo string, amountCents int64) (bool, error) {
func (c *AlipayChannel) Refund(orderNo, refundNo string, amountCents int64) (bool, error) {
// 执行：c.mu.Lock()
	c.mu.Lock()
// 执行：c.calls++
	c.calls++
// 执行：c.mu.Unlock()
	c.mu.Unlock()
// 声明：var p alipay.TradeRefund
	var p alipay.TradeRefund
// 赋值/更新：p.OutTradeNo = orderNo
	p.OutTradeNo = orderNo
// 赋值/更新：p.OutRequestNo = refundNo
	p.OutRequestNo = refundNo
// 赋值/更新：p.RefundAmount = centsToYuan(amountCents)
	p.RefundAmount = centsToYuan(amountCents)
// 赋值/更新：p.RefundReason = "商城整单退款"
	p.RefundReason = "商城整单退款"
// 赋值/更新：rsp, err := c.client.TradeRefund(context.Background(), p)
	rsp, err := c.client.TradeRefund(context.Background(), p)
// 条件判断：if err != nil {
	if err != nil {
// 返回：return false, err
		return false, err
// 结束当前代码块
	}
// 条件判断：if rsp.Code != alipay.CodeSuccess {
	if rsp.Code != alipay.CodeSuccess {
// 返回：return false, fmt.Errorf("支付宝退款受理失败: code=%s msg=%s sub=%s(%s)", rsp.Code, rsp.Msg, rsp.SubCode, rsp.SubMsg)
		return false, fmt.Errorf("支付宝退款受理失败: code=%s msg=%s sub=%s(%s)", rsp.Code, rsp.Msg, rsp.SubCode, rsp.SubMsg)
// 结束当前代码块
	}
// 返回：return true, nil
	return true, nil
// 结束当前代码块
}

// SetRefundMode 真实渠道不支持切换退款结果，忽略（保持 ChannelClient 兼容）。

// 函数定义：func (c *AlipayChannel) SetRefundMode(mode string) {}
func (c *AlipayChannel) SetRefundMode(mode string) {}

// RefundMode 渠道标识（真实渠道）

// 函数定义：func (c *AlipayChannel) RefundMode() string { return "real-alipay" }
func (c *AlipayChannel) RefundMode() string { return "real-alipay" }

// RefundCalls 累计退款调用次数（并发闸门测试断言仍可用）

// 函数定义：func (c *AlipayChannel) RefundCalls() int {
func (c *AlipayChannel) RefundCalls() int {
// 执行：c.mu.RLock()
	c.mu.RLock()
// 异步/延迟执行：defer c.mu.RUnlock()
	defer c.mu.RUnlock()
// 返回：return c.calls
	return c.calls
// 结束当前代码块
}

// centsToYuan 分转元字符串：纯整数运算，避免浮点误差（两位小数）。

// 函数定义：func centsToYuan(cents int64) string {
func centsToYuan(cents int64) string {
// 条件判断：if cents < 0 {
	if cents < 0 {
// 赋值/更新：cents = 0
		cents = 0
// 结束当前代码块
	}
// 返回：return fmt.Sprintf("%d.%02d", cents/100, cents%100)
	return fmt.Sprintf("%d.%02d", cents/100, cents%100)
// 结束当前代码块
}

// yuanToCents 元字符串转分：要求至多两位小数（支付宝 total_amount 格式）。

// 函数定义：func yuanToCents(s string) (int64, error) {
func yuanToCents(s string) (int64, error) {
// 赋值/更新：s = strings.TrimSpace(s)
	s = strings.TrimSpace(s)
// 赋值/更新：parts := strings.SplitN(s, ".", 2)
	parts := strings.SplitN(s, ".", 2)
// 赋值/更新：whole, err := strconv.ParseInt(parts[0], 10, 64)
	whole, err := strconv.ParseInt(parts[0], 10, 64)
// 条件判断：if err != nil || whole < 0 {
	if err != nil || whole < 0 {
// 返回：return 0, fmt.Errorf("金额整数部分非法")
		return 0, fmt.Errorf("金额整数部分非法")
// 结束当前代码块
	}
// 声明：var frac int64
	var frac int64
// 条件判断：if len(parts) == 2 {
	if len(parts) == 2 {
// 赋值/更新：f := parts[1]
		f := parts[1]
// 条件判断：if len(f) > 2 {
		if len(f) > 2 {
// 返回：return 0, fmt.Errorf("金额小数位超过两位")
			return 0, fmt.Errorf("金额小数位超过两位")
// 结束当前代码块
		}
// 循环处理：for len(f) < 2 {
		for len(f) < 2 {
// 赋值/更新：f += "0"
			f += "0"
// 结束当前代码块
		}
// 赋值/更新：frac, err = strconv.ParseInt(f, 10, 64)
		frac, err = strconv.ParseInt(f, 10, 64)
// 条件判断：if err != nil {
		if err != nil {
// 返回：return 0, fmt.Errorf("金额小数部分非法")
			return 0, fmt.Errorf("金额小数部分非法")
// 结束当前代码块
		}
// 结束当前代码块
	}
// 返回：return whole*100 + frac, nil
	return whole*100 + frac, nil
// 结束当前代码块
}
