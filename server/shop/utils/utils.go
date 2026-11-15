// Package utils 提供密码哈希、JWT 与业务单号生成等基础工具。
package utils

// 导入依赖包
import (
// 执行："crypto/hmac"
	"crypto/hmac"
// 执行："crypto/rand"
	"crypto/rand"
// 执行："crypto/sha256"
	"crypto/sha256"
// 执行："crypto/subtle"
	"crypto/subtle"
// 执行："encoding/base64"
	"encoding/base64"
// 执行："encoding/hex"
	"encoding/hex"
// 执行："encoding/json"
	"encoding/json"
// 执行："errors"
	"errors"
// 执行："fmt"
	"fmt"
// 执行："math/big"
	"math/big"
// 执行："strconv"
	"strconv"
// 执行："strings"
	"strings"
// 执行："time"
	"time"

// 执行："golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/pbkdf2"
// 结束当前代码块
)

// pbkdf2Iter 口令派生迭代次数（演示环境取 4096，兼顾登录耗时与爆破成本）。

// 声明：const pbkdf2Iter = 4096
const pbkdf2Iter = 4096

// HashPassword 生成 pbkdf2$iter$salt$hash 形式的密码摘要。

// 函数定义：func HashPassword(plain string) (string, error) {
func HashPassword(plain string) (string, error) {
// 赋值/更新：salt := make([]byte, 16)
	salt := make([]byte, 16)
// 条件判断：if _, err := rand.Read(salt); err != nil {
	if _, err := rand.Read(salt); err != nil {
// 返回：return "", err
		return "", err
// 结束当前代码块
	}
// 赋值/更新：sum := pbkdf2.Key([]byte(plain), salt, pbkdf2Iter, 32, sha256.New)
	sum := pbkdf2.Key([]byte(plain), salt, pbkdf2Iter, 32, sha256.New)
// 返回：return fmt.Sprintf("pbkdf2$%d$%s$%s", pbkdf2Iter, hex.EncodeToString(salt), hex.EncodeToString(sum)), nil
	return fmt.Sprintf("pbkdf2$%d$%s$%s", pbkdf2Iter, hex.EncodeToString(salt), hex.EncodeToString(sum)), nil
// 结束当前代码块
}

// VerifyPassword 校验明文密码与摘要是否匹配，比较过程使用常量时间。

// 函数定义：func VerifyPassword(plain, encoded string) bool {
func VerifyPassword(plain, encoded string) bool {
// 赋值/更新：parts := strings.Split(encoded, "$")
	parts := strings.Split(encoded, "$")
// 条件判断：if len(parts) != 4 || parts[0] != "pbkdf2" {
	if len(parts) != 4 || parts[0] != "pbkdf2" {
// 返回：return false
		return false
// 结束当前代码块
	}
// 赋值/更新：iter, err := strconv.Atoi(parts[1])
	iter, err := strconv.Atoi(parts[1])
// 条件判断：if err != nil || iter <= 0 {
	if err != nil || iter <= 0 {
// 返回：return false
		return false
// 结束当前代码块
	}
// 赋值/更新：salt, err := hex.DecodeString(parts[2])
	salt, err := hex.DecodeString(parts[2])
// 条件判断：if err != nil {
	if err != nil {
// 返回：return false
		return false
// 结束当前代码块
	}
// 赋值/更新：want, err := hex.DecodeString(parts[3])
	want, err := hex.DecodeString(parts[3])
// 条件判断：if err != nil {
	if err != nil {
// 返回：return false
		return false
// 结束当前代码块
	}
// 赋值/更新：got := pbkdf2.Key([]byte(plain), salt, iter, len(want), sha256.New)
	got := pbkdf2.Key([]byte(plain), salt, iter, len(want), sha256.New)
// 返回：return subtle.ConstantTimeCompare(got, want) == 1
	return subtle.ConstantTimeCompare(got, want) == 1
// 结束当前代码块
}

// ---------- JWT（HS256，标准库实现，避免额外依赖） ----------

// Claims JWT 载荷

// 类型定义：type Claims struct {
type Claims struct {
// 执行：UserID   uint64 `json:"sub"`
	UserID   uint64 `json:"sub"`
// 执行：Username string `json:"username"`
	Username string `json:"username"`
// 执行：Role     string `json:"role"`
	Role     string `json:"role"`
// 执行：ExpireAt int64  `json:"exp"`
	ExpireAt int64  `json:"exp"`
// 执行：IssuedAt int64  `json:"iat"`
	IssuedAt int64  `json:"iat"`
// 结束当前代码块
}

// ErrInvalidToken 令牌非法

// 声明：var ErrInvalidToken = errors.New("令牌无效或已过期")
var ErrInvalidToken = errors.New("令牌无效或已过期")

// b64 JWT 段编码：无填充的 URL 安全 Base64。

// 函数定义：func b64(buf []byte) string { return base64.RawURLEncoding.EncodeToString(buf) }
func b64(buf []byte) string { return base64.RawURLEncoding.EncodeToString(buf) }

// SignJWT 签发令牌

// 函数定义：func SignJWT(secret string, c Claims) (string, error) {
func SignJWT(secret string, c Claims) (string, error) {
// 赋值/更新：header := b64([]byte(`{"alg":"HS256","typ":"JWT"}`))
	header := b64([]byte(`{"alg":"HS256","typ":"JWT"}`))
// 赋值/更新：payloadJSON, err := json.Marshal(c)
	payloadJSON, err := json.Marshal(c)
// 条件判断：if err != nil {
	if err != nil {
// 返回：return "", err
		return "", err
// 结束当前代码块
	}
// 赋值/更新：payload := b64(payloadJSON)
	payload := b64(payloadJSON)
// 赋值/更新：signing := header + "." + payload
	signing := header + "." + payload
// 赋值/更新：mac := hmac.New(sha256.New, []byte(secret))
	mac := hmac.New(sha256.New, []byte(secret))
// 执行：mac.Write([]byte(signing))
	mac.Write([]byte(signing))
// 返回：return signing + "." + b64(mac.Sum(nil)), nil
	return signing + "." + b64(mac.Sum(nil)), nil
// 结束当前代码块
}

// ParseJWT 解析并校验令牌

// 函数定义：func ParseJWT(secret, token string, now time.Time) (Claims, error) {
func ParseJWT(secret, token string, now time.Time) (Claims, error) {
// 声明：var claims Claims
	var claims Claims
// 赋值/更新：parts := strings.Split(token, ".")
	parts := strings.Split(token, ".")
// 条件判断：if len(parts) != 3 {
	if len(parts) != 3 {
// 返回：return claims, ErrInvalidToken
		return claims, ErrInvalidToken
// 结束当前代码块
	}
// 赋值/更新：signing := parts[0] + "." + parts[1]
	signing := parts[0] + "." + parts[1]
// 赋值/更新：mac := hmac.New(sha256.New, []byte(secret))
	mac := hmac.New(sha256.New, []byte(secret))
// 执行：mac.Write([]byte(signing))
	mac.Write([]byte(signing))
// 赋值/更新：sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
// 条件判断：if err != nil || !hmac.Equal(sig, mac.Sum(nil)) {
	if err != nil || !hmac.Equal(sig, mac.Sum(nil)) {
// 返回：return claims, ErrInvalidToken
		return claims, ErrInvalidToken
// 结束当前代码块
	}
// 赋值/更新：body, err := base64.RawURLEncoding.DecodeString(parts[1])
	body, err := base64.RawURLEncoding.DecodeString(parts[1])
// 条件判断：if err != nil {
	if err != nil {
// 返回：return claims, ErrInvalidToken
		return claims, ErrInvalidToken
// 结束当前代码块
	}
// 条件判断：if err := json.Unmarshal(body, &claims); err != nil {
	if err := json.Unmarshal(body, &claims); err != nil {
// 返回：return claims, ErrInvalidToken
		return claims, ErrInvalidToken
// 结束当前代码块
	}
// 条件判断：if claims.ExpireAt <= now.Unix() {
	if claims.ExpireAt <= now.Unix() {
// 返回：return claims, ErrInvalidToken
		return claims, ErrInvalidToken
// 结束当前代码块
	}
// 返回：return claims, nil
	return claims, nil
// 结束当前代码块
}

// ---------- 业务单号 ----------

// 声明：const noCharset = "0123456789"
const noCharset = "0123456789"

// randDigits 生成 n 位密码学随机数字串，用于订单号/退款单号的随机尾段。

// 函数定义：func randDigits(n int) string {
func randDigits(n int) string {
// 赋值/更新：buf := make([]byte, 0, n)
	buf := make([]byte, 0, n)
// 循环处理：for i := 0; i < n; i++ {
	for i := 0; i < n; i++ {
// 赋值/更新：v, err := rand.Int(rand.Reader, big.NewInt(10))
		v, err := rand.Int(rand.Reader, big.NewInt(10))
// 条件判断：if err != nil {
		if err != nil {
// 赋值/更新：buf = append(buf, noCharset[0])
			buf = append(buf, noCharset[0])
// 流程控制：continue
			continue
// 结束当前代码块
		}
// 赋值/更新：buf = append(buf, noCharset[v.Int64()])
		buf = append(buf, noCharset[v.Int64()])
// 结束当前代码块
	}
// 返回：return string(buf)
	return string(buf)
// 结束当前代码块
}

// NewOrderNo 生成 20 位订单号：SH + yyMMddHHmmss + 6 位随机数字。

// 函数定义：func NewOrderNo(now time.Time) string {
func NewOrderNo(now time.Time) string {
// 返回：return "SH" + now.Format("060102150405") + randDigits(6)
	return "SH" + now.Format("060102150405") + randDigits(6)
// 结束当前代码块
}

// NewRefundNo 生成 20 位退款单号：RF + yyMMddHHmmss + 6 位随机数字。

// 函数定义：func NewRefundNo(now time.Time) string {
func NewRefundNo(now time.Time) string {
// 返回：return "RF" + now.Format("060102150405") + randDigits(6)
	return "RF" + now.Format("060102150405") + randDigits(6)
// 结束当前代码块
}

// PaySign 模拟渠道通知签名：HMAC-SHA256(secret, orderNo|txn|amount)。

// 函数定义：func PaySign(secret, orderNo, txn string, amount int64) string {
func PaySign(secret, orderNo, txn string, amount int64) string {
// 赋值/更新：mac := hmac.New(sha256.New, []byte(secret))
	mac := hmac.New(sha256.New, []byte(secret))
// 执行：mac.Write([]byte(fmt.Sprintf("%s|%s|%d", orderNo, txn, amount)))
	mac.Write([]byte(fmt.Sprintf("%s|%s|%d", orderNo, txn, amount)))
// 返回：return hex.EncodeToString(mac.Sum(nil))
	return hex.EncodeToString(mac.Sum(nil))
// 结束当前代码块
}

// Fen 把分转换为人民币字符串，便于页面展示。

// 函数定义：func Fen(cents int64) string {
func Fen(cents int64) string {
// 赋值/更新：sign := ""
	sign := ""
// 条件判断：if cents < 0 {
	if cents < 0 {
// 赋值/更新：sign = "-"
		sign = "-"
// 赋值/更新：cents = -cents
		cents = -cents
// 结束当前代码块
	}
// 返回：return fmt.Sprintf("%s%d.%02d", sign, cents/100, cents%100)
	return fmt.Sprintf("%s%d.%02d", sign, cents/100, cents%100)
// 结束当前代码块
}
