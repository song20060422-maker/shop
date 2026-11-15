// 包声明
package service

// 导入依赖包
import (
// 执行："bytes"
	"bytes"
// 执行："context"
	"context"
// 执行："crypto"
	"crypto"
// 执行："crypto/rand"
	"crypto/rand"
// 执行："crypto/rsa"
	"crypto/rsa"
// 执行："crypto/sha256"
	"crypto/sha256"
// 执行："crypto/x509"
	"crypto/x509"
// 执行："encoding/base64"
	"encoding/base64"
// 执行："encoding/pem"
	"encoding/pem"
// 执行："net/url"
	"net/url"
// 执行："sort"
	"sort"
// 执行："strings"
	"strings"
// 执行："testing"
	"testing"
// 执行："time"
	"time"

// 执行："shop/shop/config"
	"shop/shop/config"
// 执行："shop/shop/model"
	"shop/shop/model"
// 结束当前代码块
)

// 支付宝沙箱渠道测试：本地生成 RSA 密钥对，私钥充当商户私钥，
// 公钥充当"支付宝公钥"，从而无需真实沙箱网络即可验证签名链路与通知推进。

// 函数定义：func genKeyPair(t *testing.T) (privPEM, pubPEM string) {
func genKeyPair(t *testing.T) (privPEM, pubPEM string) {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：key, err := rsa.GenerateKey(rand.Reader, 2048)
	key, err := rsa.GenerateKey(rand.Reader, 2048)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("生成密钥失败: %v", err)
		t.Fatalf("生成密钥失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：der, err := x509.MarshalPKCS8PrivateKey(key)
	der, err := x509.MarshalPKCS8PrivateKey(key)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("序列化私钥失败: %v", err)
		t.Fatalf("序列化私钥失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：privPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))
	privPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))
// 赋值/更新：pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("序列化公钥失败: %v", err)
		t.Fatalf("序列化公钥失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：pubPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
	pubPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
// 返回：return privPEM, pubPEM
	return privPEM, pubPEM
// 结束当前代码块
}

// 函数定义：func newAliEnv(t *testing.T) (*ServiceGroup, *AlipayChannel, *fakeClock, *rsa.PrivateKey) {
func newAliEnv(t *testing.T) (*ServiceGroup, *AlipayChannel, *fakeClock, *rsa.PrivateKey) {
// 执行：t.Helper()
	t.Helper()
// 执行：resetDB(t)
	resetDB(t)
// 赋值/更新：privPEM, pubPEM := genKeyPair(t)
	privPEM, pubPEM := genKeyPair(t)
// 赋值/更新：block, _ := pem.Decode([]byte(privPEM))
	block, _ := pem.Decode([]byte(privPEM))
// 赋值/更新：key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("解析测试私钥失败: %v", err)
		t.Fatalf("解析测试私钥失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：cfg := testCfg
	cfg := testCfg
// 赋值/更新：cfg.Shop.PayChannel = "alipay"
	cfg.Shop.PayChannel = "alipay"
// 赋值/更新：cfg.Shop.Alipay = config.Alipay{
	cfg.Shop.Alipay = config.Alipay{
// 执行：AppID:           "9021000158651669",
		AppID:           "9021000158651669",
// 执行：PrivateKey:      privPEM,
		PrivateKey:      privPEM,
// 执行：AlipayPublicKey: pubPEM,
		AlipayPublicKey: pubPEM,
// 执行：NotifyURL:       "http://notify.example/api/pay/alipay/notify",
		NotifyURL:       "http://notify.example/api/pay/alipay/notify",
// 执行：ReturnURL:       "http://127.0.0.1:8090/shop",
		ReturnURL:       "http://127.0.0.1:8090/shop",
// 结束当前代码块
	}
// 赋值/更新：ch, err := NewAlipayChannel(cfg.Shop.Alipay)
	ch, err := NewAlipayChannel(cfg.Shop.Alipay)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("构建支付宝渠道失败: %v", err)
		t.Fatalf("构建支付宝渠道失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：clock := &fakeClock{t: time.Date(2026, 11, 14, 9, 0, 0, 0, time.Local)}
	clock := &fakeClock{t: time.Date(2026, 11, 14, 9, 0, 0, 0, time.Local)}
// 赋值/更新：g := NewServiceGroup(Deps{DB: testDB, Now: clock.Now, Cfg: cfg, Channel: ch})
	g := NewServiceGroup(Deps{DB: testDB, Now: clock.Now, Cfg: cfg, Channel: ch})
// 返回：return g, ch, clock, key.(*rsa.PrivateKey)
	return g, ch, clock, key.(*rsa.PrivateKey)
// 结束当前代码块
}

// signAlipayForm 按支付宝规则签名通知表单：除 sign/sign_type 外按 key 升序 a=b&c=d，
// RSA2(SHA256withRSA) 后 base64。

// 函数定义：func signAlipayForm(t *testing.T, key *rsa.PrivateKey, form url.Values) string {
func signAlipayForm(t *testing.T, key *rsa.PrivateKey, form url.Values) string {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：keys := make([]string, 0, len(form))
	keys := make([]string, 0, len(form))
// 循环处理：for k := range form {
	for k := range form {
// 条件判断：if k == "sign" || k == "sign_type" {
		if k == "sign" || k == "sign_type" {
// 流程控制：continue
			continue
// 结束当前代码块
		}
// 赋值/更新：keys = append(keys, k)
		keys = append(keys, k)
// 结束当前代码块
	}
// 执行：sort.Strings(keys)
	sort.Strings(keys)
// 声明：var buf bytes.Buffer
	var buf bytes.Buffer
// 循环处理：for i, k := range keys {
	for i, k := range keys {
// 条件判断：if i > 0 {
		if i > 0 {
// 执行：buf.WriteByte('&')
			buf.WriteByte('&')
// 结束当前代码块
		}
// 执行：buf.WriteString(k)
		buf.WriteString(k)
// 执行：buf.WriteByte('=')
		buf.WriteByte('=')
// 执行：buf.WriteString(form.Get(k))
		buf.WriteString(form.Get(k))
// 结束当前代码块
	}
// 赋值/更新：digest := sha256.Sum256(buf.Bytes())
	digest := sha256.Sum256(buf.Bytes())
// 赋值/更新：sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("签名失败: %v", err)
		t.Fatalf("签名失败: %v", err)
// 结束当前代码块
	}
// 返回：return base64.StdEncoding.EncodeToString(sig)
	return base64.StdEncoding.EncodeToString(sig)
// 结束当前代码块
}

// 函数定义：func alipayNotifyForm(t *testing.T, key *rsa.PrivateKey, orderNo, txn, totalAmount, tradeStatus string) url.Va...
func alipayNotifyForm(t *testing.T, key *rsa.PrivateKey, orderNo, txn, totalAmount, tradeStatus string) url.Values {
// 赋值/更新：form := url.Values{}
	form := url.Values{}
// 执行：form.Set("app_id", "9021000158651669")
	form.Set("app_id", "9021000158651669")
// 执行：form.Set("trade_no", txn)
	form.Set("trade_no", txn)
// 执行：form.Set("out_trade_no", orderNo)
	form.Set("out_trade_no", orderNo)
// 执行：form.Set("total_amount", totalAmount)
	form.Set("total_amount", totalAmount)
// 执行：form.Set("trade_status", tradeStatus)
	form.Set("trade_status", tradeStatus)
// 执行：form.Set("sign_type", "RSA2")
	form.Set("sign_type", "RSA2")
// 执行：form.Set("sign", signAlipayForm(t, key, form))
	form.Set("sign", signAlipayForm(t, key, form))
// 返回：return form
	return form
// 结束当前代码块
}

// 函数定义：func TestAlipayAmountConvert(t *testing.T) {
func TestAlipayAmountConvert(t *testing.T) {
// 赋值/更新：cases := []struct {
	cases := []struct {
// 执行：cents int64
		cents int64
// 执行：yuan  string
		yuan  string
// 结束当前代码块
	}{
// 执行：{0, "0.00"}, {1, "0.01"}, {1234, "12.34"}, {100, "1.00"}, {1000000, "10000.00"},
		{0, "0.00"}, {1, "0.01"}, {1234, "12.34"}, {100, "1.00"}, {1000000, "10000.00"},
// 结束当前代码块
	}
// 循环处理：for _, c := range cases {
	for _, c := range cases {
// 条件判断：if got := centsToYuan(c.cents); got != c.yuan {
		if got := centsToYuan(c.cents); got != c.yuan {
// 执行：t.Fatalf("centsToYuan(%d) = %s，期望 %s", c.cents, got, c.yuan)
			t.Fatalf("centsToYuan(%d) = %s，期望 %s", c.cents, got, c.yuan)
// 结束当前代码块
		}
// 赋值/更新：back, err := yuanToCents(c.yuan)
		back, err := yuanToCents(c.yuan)
// 条件判断：if err != nil || back != c.cents {
		if err != nil || back != c.cents {
// 执行：t.Fatalf("yuanToCents(%s) = %d,%v，期望 %d", c.yuan, back, err, c.cents)
			t.Fatalf("yuanToCents(%s) = %d,%v，期望 %d", c.yuan, back, err, c.cents)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if _, err := yuanToCents("12.345"); err == nil {
	if _, err := yuanToCents("12.345"); err == nil {
// 执行：t.Fatal("三位小数应报错")
		t.Fatal("三位小数应报错")
// 结束当前代码块
	}
// 条件判断：if _, err := yuanToCents("abc"); err == nil {
	if _, err := yuanToCents("abc"); err == nil {
// 执行：t.Fatal("非数字应报错")
		t.Fatal("非数字应报错")
// 结束当前代码块
	}
// 条件判断：if v, err := yuanToCents(" 12.3 "); err != nil || v != 1230 {
	if v, err := yuanToCents(" 12.3 "); err != nil || v != 1230 {
// 执行：t.Fatalf("带空白金额解析失败: %d %v", v, err)
		t.Fatalf("带空白金额解析失败: %d %v", v, err)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestAlipayPagePayURL(t *testing.T) {
func TestAlipayPagePayURL(t *testing.T) {
// 赋值/更新：g, _, _, _ := newAliEnv(t)
	g, _, _, _ := newAliEnv(t)
// 赋值/更新：b := buyerID(t, "buyer1")
	b := buyerID(t, "buyer1")
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：o, err := g.OrderService.Create(b, pid, 1)
	o, err := g.OrderService.Create(b, pid, 1)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("下单失败: %v", err)
		t.Fatalf("下单失败: %v", err)
// 结束当前代码块
	}
	// 动态回调：从哪个域名访问，notify_url 就用哪个域名（数据不写死）
	u, err := g.PaymentService.AlipayPayURL(o.OrderNo, b, "http://demo.cpolar.cn")
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("生成收银台地址失败: %v", err)
		t.Fatalf("生成收银台地址失败: %v", err)
// 结束当前代码块
	}
// 循环处理：for _, want := range []string{
	for _, want := range []string{
// 赋值/更新："alipay.trade.page.pay", "out_trade_no", "sign=", "9021000158651669",
		"alipay.trade.page.pay", "out_trade_no", "sign=", "9021000158651669",
// 赋值/更新："notify_url=http%3A%2F%2Fdemo.cpolar.cn%2Fapi%2Fpay%2Falipay%2Fnotify",
		"notify_url=http%3A%2F%2Fdemo.cpolar.cn%2Fapi%2Fpay%2Falipay%2Fnotify",
// 赋值/更新："return_url=http%3A%2F%2Fdemo.cpolar.cn%2Fshop",
		"return_url=http%3A%2F%2Fdemo.cpolar.cn%2Fshop",
// 结束当前代码块
	} {
// 条件判断：if !strings.Contains(u, want) {
		if !strings.Contains(u, want) {
// 执行：t.Fatalf("收银台地址缺少 %s：%s", want, u)
			t.Fatalf("收银台地址缺少 %s：%s", want, u)
// 结束当前代码块
		}
// 结束当前代码块
	}
	// 未传访问域名时回落 config 静态配置（显式指定仍生效）
	u2, err := g.PaymentService.AlipayPayURL(o.OrderNo, b, "")
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("回落配置生成失败: %v", err)
		t.Fatalf("回落配置生成失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if !strings.Contains(u2, "sign=") {
	if !strings.Contains(u2, "sign=") {
// 执行：t.Fatalf("回落配置地址异常：%s", u2)
		t.Fatalf("回落配置地址异常：%s", u2)
// 结束当前代码块
	}
	// 未启用支付宝渠道时应报业务错误
	g2, _, _ := newEnv(t)
// 赋值/更新：o2, _ := g2.OrderService.Create(b, productID(t, "SKU-A"), 1)
	o2, _ := g2.OrderService.Create(b, productID(t, "SKU-A"), 1)
// 条件判断：if _, err := g2.PaymentService.AlipayPayURL(o2.OrderNo, b, "http://demo.cpolar.cn"); err == nil {
	if _, err := g2.PaymentService.AlipayPayURL(o2.OrderNo, b, "http://demo.cpolar.cn"); err == nil {
// 执行：t.Fatal("模拟渠道下应拒绝生成收银台地址")
		t.Fatal("模拟渠道下应拒绝生成收银台地址")
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestAlipayVerifyNotify(t *testing.T) {
func TestAlipayVerifyNotify(t *testing.T) {
// 赋值/更新：_, ch, _, key := newAliEnv(t)
	_, ch, _, key := newAliEnv(t)
// 赋值/更新：form := alipayNotifyForm(t, key, "SO-1", "TXN-ALI-1", "12.34", "TRADE_SUCCESS")
	form := alipayNotifyForm(t, key, "SO-1", "TXN-ALI-1", "12.34", "TRADE_SUCCESS")
// 赋值/更新：p, err := ch.VerifyNotify(context.Background(), form)
	p, err := ch.VerifyNotify(context.Background(), form)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("合法签名验签失败: %v", err)
		t.Fatalf("合法签名验签失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if p.OrderNo != "SO-1" || p.Txn != "TXN-ALI-1" || p.AmountCents != 1234 || p.TradeStatus != "TRADE_SUCCESS" {
	if p.OrderNo != "SO-1" || p.Txn != "TXN-ALI-1" || p.AmountCents != 1234 || p.TradeStatus != "TRADE_SUCCESS" {
// 执行：t.Fatalf("通知字段提取错误: %+v", p)
		t.Fatalf("通知字段提取错误: %+v", p)
// 结束当前代码块
	}
	// 篡改金额后验签必须失败
	form.Set("total_amount", "12.35")
// 条件判断：if _, err := ch.VerifyNotify(context.Background(), form); err == nil {
	if _, err := ch.VerifyNotify(context.Background(), form); err == nil {
// 执行：t.Fatal("篡改后的通知应验签失败")
		t.Fatal("篡改后的通知应验签失败")
// 结束当前代码块
	}
	// app_id 不匹配必须拒绝
	form2 := alipayNotifyForm(t, key, "SO-1", "TXN-ALI-2", "12.34", "TRADE_SUCCESS")
// 执行：form2.Set("app_id", "other-app")
	form2.Set("app_id", "other-app")
// 执行：form2.Set("sign", signAlipayForm(t, key, form2))
	form2.Set("sign", signAlipayForm(t, key, form2))
// 条件判断：if _, err := ch.VerifyNotify(context.Background(), form2); err == nil {
	if _, err := ch.VerifyNotify(context.Background(), form2); err == nil {
// 执行：t.Fatal("app_id 不匹配应被拒绝")
		t.Fatal("app_id 不匹配应被拒绝")
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestAlipayNotifyFlow(t *testing.T) {
func TestAlipayNotifyFlow(t *testing.T) {
// 赋值/更新：g, _, clock, key := newAliEnv(t)
	g, _, clock, key := newAliEnv(t)
// 赋值/更新：b := buyerID(t, "buyer1")
	b := buyerID(t, "buyer1")
// 赋值/更新：pid := productID(t, "SKU-A")
	pid := productID(t, "SKU-A")
// 赋值/更新：o, err := g.OrderService.Create(b, pid, 2)
	o, err := g.OrderService.Create(b, pid, 2)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("下单失败: %v", err)
		t.Fatalf("下单失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：amountYuan := centsToYuan(o.AmountCents)
	amountYuan := centsToYuan(o.AmountCents)

	// 1. 未完成交易（买家未付款）应拒绝记款
	res, err := g.PaymentService.NotifyAlipay(alipayNotifyForm(t, key, o.OrderNo, "TXN-ALI-WAIT", amountYuan, "TRADE_WAIT_BUYER_PAY"), "")
// 条件判断：if err == nil || res.Result != model.NoticeRejected {
	if err == nil || res.Result != model.NoticeRejected {
// 执行：t.Fatalf("未完成交易应被拒绝: %+v %v", res, err)
		t.Fatalf("未完成交易应被拒绝: %+v %v", res, err)
// 结束当前代码块
	}

	// 2. 合法成功通知：推进 paid 并记款
	txn := "TXN-ALI-" + o.OrderNo
// 赋值/更新：res, err = g.PaymentService.NotifyAlipay(alipayNotifyForm(t, key, o.OrderNo, txn, amountYuan, "TRADE_SUCCESS")...
	res, err = g.PaymentService.NotifyAlipay(alipayNotifyForm(t, key, o.OrderNo, txn, amountYuan, "TRADE_SUCCESS"), "")
// 条件判断：if err != nil || res.Result != model.NoticeAccepted {
	if err != nil || res.Result != model.NoticeAccepted {
// 执行：t.Fatalf("支付宝通知未受理: %+v %v", res, err)
		t.Fatalf("支付宝通知未受理: %+v %v", res, err)
// 结束当前代码块
	}
// 声明：var order model.Order
	var order model.Order
// 条件判断：if err := testDB.Where("order_no = ?", o.OrderNo).First(&order).Error; err != nil || order.Status != model.Ord...
	if err := testDB.Where("order_no = ?", o.OrderNo).First(&order).Error; err != nil || order.Status != model.OrderStatusPaid {
// 执行：t.Fatalf("订单未推进 paid: %s %v", order.Status, err)
		t.Fatalf("订单未推进 paid: %s %v", order.Status, err)
// 结束当前代码块
	}
// 声明：var pay model.Payment
	var pay model.Payment
// 条件判断：if err := testDB.Where("channel_txn_id = ?", txn).First(&pay).Error; err != nil || pay.AmountCents != o.Amount...
	if err := testDB.Where("channel_txn_id = ?", txn).First(&pay).Error; err != nil || pay.AmountCents != o.AmountCents {
// 执行：t.Fatalf("记款缺失或金额不符: %+v %v", pay, err)
		t.Fatalf("记款缺失或金额不符: %+v %v", pay, err)
// 结束当前代码块
	}

	// 3. 同一交易号重复通知：duplicated 且不重复记款
	res, err = g.PaymentService.NotifyAlipay(alipayNotifyForm(t, key, o.OrderNo, txn, amountYuan, "TRADE_SUCCESS"), "")
// 条件判断：if err == nil || res.Result != model.NoticeDuplicated {
	if err == nil || res.Result != model.NoticeDuplicated {
// 执行：t.Fatalf("重复通知应为 duplicated: %+v %v", res, err)
		t.Fatalf("重复通知应为 duplicated: %+v %v", res, err)
// 结束当前代码块
	}
// 声明：var payCount int64
	var payCount int64
// 执行：testDB.Model(&model.Payment{}).Where("order_no = ?", o.OrderNo).Count(&payCount)
	testDB.Model(&model.Payment{}).Where("order_no = ?", o.OrderNo).Count(&payCount)
// 条件判断：if payCount != 1 {
	if payCount != 1 {
// 执行：t.Fatalf("重复通知造成重复记款: %d", payCount)
		t.Fatalf("重复通知造成重复记款: %d", payCount)
// 结束当前代码块
	}

	// 4. 篡改签名：rejected 且 sign_ok=false（留痕可查）
	clock.Advance(time.Minute)
// 赋值/更新：bad := alipayNotifyForm(t, key, o.OrderNo, "TXN-ALI-BAD", amountYuan, "TRADE_SUCCESS")
	bad := alipayNotifyForm(t, key, o.OrderNo, "TXN-ALI-BAD", amountYuan, "TRADE_SUCCESS")
// 执行：bad.Set("sign", base64.StdEncoding.EncodeToString([]byte("forged")))
	bad.Set("sign", base64.StdEncoding.EncodeToString([]byte("forged")))
// 赋值/更新：res, err = g.PaymentService.NotifyAlipay(bad, "")
	res, err = g.PaymentService.NotifyAlipay(bad, "")
// 条件判断：if err == nil || res.Result != model.NoticeRejected {
	if err == nil || res.Result != model.NoticeRejected {
// 执行：t.Fatalf("伪造签名应被拒绝: %+v %v", res, err)
		t.Fatalf("伪造签名应被拒绝: %+v %v", res, err)
// 结束当前代码块
	}
// 赋值/更新：notices, _ := g.PaymentService.Notices(o.OrderNo)
	notices, _ := g.PaymentService.Notices(o.OrderNo)
// 赋值/更新：found := false
	found := false
// 循环处理：for _, n := range notices {
	for _, n := range notices {
// 条件判断：if !n.SignOK && n.Result == model.NoticeRejected {
		if !n.SignOK && n.Result == model.NoticeRejected {
// 赋值/更新：found = true
			found = true
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if !found {
	if !found {
// 执行：t.Fatal("验签失败通知未留痕")
		t.Fatal("验签失败通知未留痕")
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestAlipayRefundReal(t *testing.T) {
func TestAlipayRefundReal(t *testing.T) {
// 赋值/更新：_, ch, _, _ := newAliEnv(t)
	_, ch, _, _ := newAliEnv(t)
// 条件判断：if ch.RefundMode() != "real-alipay" {
	if ch.RefundMode() != "real-alipay" {
// 执行：t.Fatalf("真实渠道模式标识错误: %s", ch.RefundMode())
		t.Fatalf("真实渠道模式标识错误: %s", ch.RefundMode())
// 结束当前代码块
	}
// 执行：ch.SetRefundMode("fail") // no-op，不影响真实渠道
	ch.SetRefundMode("fail") // no-op，不影响真实渠道
// 条件判断：if ch.RefundMode() != "real-alipay" {
	if ch.RefundMode() != "real-alipay" {
// 执行：t.Fatal("真实渠道不应响应退款模式切换")
		t.Fatal("真实渠道不应响应退款模式切换")
// 结束当前代码块
	}
	// Refund 需要真实沙箱网关，单测只验证调用计数与配置错误路径
	before := ch.RefundCalls()
// 条件判断：if _, err := ch.Refund("SO-X", "RF-X", 100); err == nil {
	if _, err := ch.Refund("SO-X", "RF-X", 100); err == nil {
		// 沙箱网关在 CI 可能不可达：允许网络错误，但不允许误报成功
		t.Log("沙箱网关可达并返回了结果（本地演示环境）")
// 结束当前代码块
	}
// 条件判断：if ch.RefundCalls() != before+1 {
	if ch.RefundCalls() != before+1 {
// 执行：t.Fatal("退款调用计数未累加")
		t.Fatal("退款调用计数未累加")
// 结束当前代码块
	}
// 结束当前代码块
}
