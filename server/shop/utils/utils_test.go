// 包声明
package utils

// 导入依赖包
import (
// 执行："strings"
	"strings"
// 执行："testing"
	"testing"
// 执行："time"
	"time"
// 结束当前代码块
)

// 函数定义：func TestPasswordHelpers(t *testing.T) {
func TestPasswordHelpers(t *testing.T) {
// 赋值/更新：hash, err := HashPassword("shop123")
	hash, err := HashPassword("shop123")
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 条件判断：if !VerifyPassword("shop123", hash) {
	if !VerifyPassword("shop123", hash) {
// 执行：t.Fatal("正确密码校验失败")
		t.Fatal("正确密码校验失败")
// 结束当前代码块
	}
// 条件判断：if VerifyPassword("wrong", hash) {
	if VerifyPassword("wrong", hash) {
// 执行：t.Fatal("错误密码不应通过")
		t.Fatal("错误密码不应通过")
// 结束当前代码块
	}
	// 非法编码格式
	for _, bad := range []string{"", "plain", "x$1$aa$bb", "pbkdf2$0$aa$bb", "pbkdf2$1$zz$bb"} {
// 条件判断：if VerifyPassword("shop123", bad) {
		if VerifyPassword("shop123", bad) {
// 执行：t.Fatalf("非法哈希不应通过: %s", bad)
			t.Fatalf("非法哈希不应通过: %s", bad)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestJWTAndSign(t *testing.T) {
func TestJWTAndSign(t *testing.T) {
// 赋值/更新：secret := "unit-secret"
	secret := "unit-secret"
// 赋值/更新：now := time.Date(2026, 11, 12, 10, 0, 0, 0, time.Local)
	now := time.Date(2026, 11, 12, 10, 0, 0, 0, time.Local)
// 赋值/更新：claims := Claims{UserID: 9, Username: "buyer1", Role: "buyer",
	claims := Claims{UserID: 9, Username: "buyer1", Role: "buyer",
// 执行：ExpireAt: now.Add(time.Hour).Unix(), IssuedAt: now.Unix()}
		ExpireAt: now.Add(time.Hour).Unix(), IssuedAt: now.Unix()}
// 赋值/更新：token, err := SignJWT(secret, claims)
	token, err := SignJWT(secret, claims)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 赋值/更新：got, err := ParseJWT(secret, token, now)
	got, err := ParseJWT(secret, token, now)
// 条件判断：if err != nil || got.UserID != 9 || got.Role != "buyer" {
	if err != nil || got.UserID != 9 || got.Role != "buyer" {
// 执行：t.Fatalf("令牌解析失败: %v %+v", err, got)
		t.Fatalf("令牌解析失败: %v %+v", err, got)
// 结束当前代码块
	}
	// 过期令牌
	if _, err := ParseJWT(secret, token, now.Add(2*time.Hour)); err == nil {
// 执行：t.Fatal("过期令牌应拒绝")
		t.Fatal("过期令牌应拒绝")
// 结束当前代码块
	}
	// 篡改签名
	parts := strings.Split(token, ".")
// 条件判断：if _, err := ParseJWT(secret, parts[0]+"."+parts[1]+".deadbeef", now); err == nil {
	if _, err := ParseJWT(secret, parts[0]+"."+parts[1]+".deadbeef", now); err == nil {
// 执行：t.Fatal("篡改签名应拒绝")
		t.Fatal("篡改签名应拒绝")
// 结束当前代码块
	}
	// 结构非法
	if _, err := ParseJWT(secret, "a.b", now); err == nil {
// 执行：t.Fatal("非法结构应拒绝")
		t.Fatal("非法结构应拒绝")
// 结束当前代码块
	}
// 条件判断：if _, err := ParseJWT(secret, "a.b.c", now); err == nil {
	if _, err := ParseJWT(secret, "a.b.c", now); err == nil {
// 执行：t.Fatal("非法载荷应拒绝")
		t.Fatal("非法载荷应拒绝")
// 结束当前代码块
	}

	// 支付签名对字段顺序敏感
	s1 := PaySign(secret, "SO1", "TXN1", 100)
// 赋值/更新：s2 := PaySign(secret, "SO1", "TXN1", 100)
	s2 := PaySign(secret, "SO1", "TXN1", 100)
// 条件判断：if s1 != s2 || len(s1) != 64 {
	if s1 != s2 || len(s1) != 64 {
// 执行：t.Fatal("支付签名不稳定")
		t.Fatal("支付签名不稳定")
// 结束当前代码块
	}
// 条件判断：if s1 == PaySign(secret, "SO2", "TXN1", 100) {
	if s1 == PaySign(secret, "SO2", "TXN1", 100) {
// 执行：t.Fatal("不同订单签名不应相同")
		t.Fatal("不同订单签名不应相同")
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestNumberAndFormatHelpers(t *testing.T) {
func TestNumberAndFormatHelpers(t *testing.T) {
// 赋值/更新：now := time.Date(2026, 11, 12, 10, 0, 0, 0, time.Local)
	now := time.Date(2026, 11, 12, 10, 0, 0, 0, time.Local)
// 赋值/更新：orderNo := NewOrderNo(now)
	orderNo := NewOrderNo(now)
// 条件判断：if !strings.HasPrefix(orderNo, "SH") || len(orderNo) != 20 {
	if !strings.HasPrefix(orderNo, "SH") || len(orderNo) != 20 {
// 执行：t.Fatalf("订单号格式错误: %s", orderNo)
		t.Fatalf("订单号格式错误: %s", orderNo)
// 结束当前代码块
	}
// 条件判断：if !strings.HasPrefix(NewRefundNo(now), "RF") {
	if !strings.HasPrefix(NewRefundNo(now), "RF") {
// 执行：t.Fatal("退款单号格式错误")
		t.Fatal("退款单号格式错误")
// 结束当前代码块
	}
// 条件判断：if Fen(0) == "" || Fen(12345) == "" || Fen(-100) == "" {
	if Fen(0) == "" || Fen(12345) == "" || Fen(-100) == "" {
// 执行：t.Fatal("金额格式化不应为空")
		t.Fatal("金额格式化不应为空")
// 结束当前代码块
	}
// 条件判断：if !strings.Contains(Fen(12345), "123") {
	if !strings.Contains(Fen(12345), "123") {
// 执行：t.Fatalf("金额格式化错误: %s", Fen(12345))
		t.Fatalf("金额格式化错误: %s", Fen(12345))
// 结束当前代码块
	}
// 结束当前代码块
}
