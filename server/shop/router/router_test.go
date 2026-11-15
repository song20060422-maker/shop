// 包声明
package router

// 导入依赖包
import (
// 执行："bytes"
	"bytes"
// 执行："crypto/hmac"
	"crypto/hmac"
// 执行："crypto/sha256"
	"crypto/sha256"
// 执行："encoding/hex"
	"encoding/hex"
// 执行："encoding/json"
	"encoding/json"
// 执行："net/http"
	"net/http"
// 执行："net/http/httptest"
	"net/http/httptest"
// 执行："os"
	"os"
// 执行："testing"
	"testing"
// 执行："time"
	"time"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"
// 执行："gorm.io/gorm"
	"gorm.io/gorm"

// 执行："shop/shop/config"
	"shop/shop/config"
// 执行：shopdb "shop/shop/db"
	shopdb "shop/shop/db"
// 执行："shop/shop/migrations"
	"shop/shop/migrations"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// 说明：HTTP 集成测试使用真实 MySQL（本机 3307，测试库 shop_test），
// 覆盖 router 装配、JWT/角色门禁、契约兼容层与内置页面。
// 数据库可用 SHOP_TEST_DSN 覆盖（须 -p 1 串行执行）。

// 声明：var testDB *gorm.DB
var testDB *gorm.DB
// 声明：var testCfg config.Config
var testCfg config.Config

// 函数定义：func TestMain(m *testing.M) {
func TestMain(m *testing.M) {
// 赋值/更新：dsn := os.Getenv("SHOP_TEST_DSN")
	dsn := os.Getenv("SHOP_TEST_DSN")
// 条件判断：if dsn == "" {
	if dsn == "" {
// 赋值/更新：dsn = "root:root@tcp(127.0.0.1:3307)/shop_test?charset=utf8mb4&parseTime=True&loc=Local"
		dsn = "root:root@tcp(127.0.0.1:3307)/shop_test?charset=utf8mb4&parseTime=True&loc=Local"
// 结束当前代码块
	}
// 赋值/更新：testCfg = config.Default()
	testCfg = config.Default()
// 赋值/更新：testCfg.Mysql.DSN = dsn
	testCfg.Mysql.DSN = dsn
// 赋值/更新：testCfg.Shop.PayNotifySecret = "test-secret"
	testCfg.Shop.PayNotifySecret = "test-secret"
// 赋值/更新：testCfg.Shop.OrderTimeoutMinutes = 10
	testCfg.Shop.OrderTimeoutMinutes = 10
// 赋值/更新：db, err := shopdb.Gorm(testCfg)
	db, err := shopdb.Gorm(testCfg)
// 条件判断：if err != nil {
	if err != nil {
// 流程控制：panic("连接测试库失败（请确认 MySQL 3307 可用）: " + err.Error())
		panic("连接测试库失败（请确认 MySQL 3307 可用）: " + err.Error())
// 结束当前代码块
	}
// 条件判断：if err := shopdb.Migrate(db, migrations.FS); err != nil {
	if err := shopdb.Migrate(db, migrations.FS); err != nil {
// 流程控制：panic("测试库迁移失败: " + err.Error())
		panic("测试库迁移失败: " + err.Error())
// 结束当前代码块
	}
// 条件判断：if err := shopdb.SeedAccounts(db, "shop123"); err != nil {
	if err := shopdb.SeedAccounts(db, "shop123"); err != nil {
// 流程控制：panic("种子账号失败: " + err.Error())
		panic("种子账号失败: " + err.Error())
// 结束当前代码块
	}
// 赋值/更新：testDB = db
	testDB = db
// 执行：os.Exit(m.Run())
	os.Exit(m.Run())
// 结束当前代码块
}

// 函数定义：func resetTables(t *testing.T) {
func resetTables(t *testing.T) {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：tables := []string{
	tables := []string{
		// sessions 先清：fk_sessions_user 引用 users（共享测试库须防跨包残留）
		"sessions",
// 执行："abnormal_payments", "payments", "payment_notices", "refunds", "shipments",
		"abnormal_payments", "payments", "payment_notices", "refunds", "shipments",
// 执行："inventory_ledger", "operation_logs", "orders", "products", "users",
		"inventory_ledger", "operation_logs", "orders", "products", "users",
// 结束当前代码块
	}
// 循环处理：for _, tbl := range tables {
	for _, tbl := range tables {
// 条件判断：if err := testDB.Exec("DELETE FROM " + tbl).Error; err != nil {
		if err := testDB.Exec("DELETE FROM " + tbl).Error; err != nil {
// 执行：t.Fatalf("清理 %s 失败: %v", tbl, err)
			t.Fatalf("清理 %s 失败: %v", tbl, err)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if err := shopdb.Migrate(testDB, migrations.FS); err != nil {
	if err := shopdb.Migrate(testDB, migrations.FS); err != nil {
// 执行：t.Fatalf("重放迁移失败: %v", err)
		t.Fatalf("重放迁移失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := shopdb.SeedAccounts(testDB, "shop123"); err != nil {
	if err := shopdb.SeedAccounts(testDB, "shop123"); err != nil {
// 执行：t.Fatalf("重放种子账号失败: %v", err)
		t.Fatalf("重放种子账号失败: %v", err)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func newServer(t *testing.T, demoMode bool) (*gin.Engine, *service.SimulatedChannel) {
func newServer(t *testing.T, demoMode bool) (*gin.Engine, *service.SimulatedChannel) {
// 执行：t.Helper()
	t.Helper()
// 执行：resetTables(t)
	resetTables(t)
// 赋值/更新：cfg := testCfg
	cfg := testCfg
// 赋值/更新：cfg.Shop.DemoMode = demoMode
	cfg.Shop.DemoMode = demoMode
// 执行：gin.SetMode(gin.TestMode)
	gin.SetMode(gin.TestMode)
// 赋值/更新：r := gin.New()
	r := gin.New()
// 赋值/更新：ch := service.NewSimulatedChannel(true)
	ch := service.NewSimulatedChannel(true)
// 赋值/更新：g := service.NewServiceGroup(service.Deps{DB: testDB, Now: time.Now, Cfg: cfg, Channel: ch})
	g := service.NewServiceGroup(service.Deps{DB: testDB, Now: time.Now, Cfg: cfg, Channel: ch})
// 执行：InitShopRouter(r.Group(""), g, cfg, time.Now)
	InitShopRouter(r.Group(""), g, cfg, time.Now)
// 返回：return r, ch
	return r, ch
// 结束当前代码块
}

// 函数定义：func doJSON(r *gin.Engine, method, path, body, token string, headers map[string]string) *httptest.ResponseReco...
func doJSON(r *gin.Engine, method, path, body, token string, headers map[string]string) *httptest.ResponseRecorder {
// 声明：var req *http.Request
	var req *http.Request
// 条件判断：if body != "" {
	if body != "" {
// 赋值/更新：req = httptest.NewRequest(method, path, bytes.NewBufferString(body))
		req = httptest.NewRequest(method, path, bytes.NewBufferString(body))
// 执行：req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Type", "application/json")
// 条件分支：} else {
	} else {
// 赋值/更新：req = httptest.NewRequest(method, path, nil)
		req = httptest.NewRequest(method, path, nil)
// 结束当前代码块
	}
// 条件判断：if token != "" {
	if token != "" {
// 执行：req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Authorization", "Bearer "+token)
// 结束当前代码块
	}
// 循环处理：for k, v := range headers {
	for k, v := range headers {
// 执行：req.Header.Set(k, v)
		req.Header.Set(k, v)
// 结束当前代码块
	}
// 赋值/更新：rr := httptest.NewRecorder()
	rr := httptest.NewRecorder()
// 执行：r.ServeHTTP(rr, req)
	r.ServeHTTP(rr, req)
// 返回：return rr
	return rr
// 结束当前代码块
}

// 函数定义：func parseBody(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
func parseBody(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：out := map[string]interface{}{}
	out := map[string]interface{}{}
// 条件判断：if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
// 执行：t.Fatalf("响应不是 JSON: %v body=%s", err, rr.Body.String())
		t.Fatalf("响应不是 JSON: %v body=%s", err, rr.Body.String())
// 结束当前代码块
	}
// 返回：return out
	return out
// 结束当前代码块
}

// 函数定义：func login(t *testing.T, r *gin.Engine, username string) string {
func login(t *testing.T, r *gin.Engine, username string) string {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：rr := doJSON(r, "POST", "/api/login", `{"username":"`+username+`","password":"shop123"}`, "", nil)
	rr := doJSON(r, "POST", "/api/login", `{"username":"`+username+`","password":"shop123"}`, "", nil)
// 条件判断：if rr.Code != http.StatusOK {
	if rr.Code != http.StatusOK {
// 执行：t.Fatalf("登录 %s 失败: %d %s", username, rr.Code, rr.Body.String())
		t.Fatalf("登录 %s 失败: %d %s", username, rr.Code, rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：res := parseBody(t, rr)
	res := parseBody(t, rr)
// 赋值/更新：token, _ := res["token"].(string)
	token, _ := res["token"].(string)
// 条件判断：if token == "" {
	if token == "" {
// 执行：t.Fatalf("登录 %s 未返回 token: %s", username, rr.Body.String())
		t.Fatalf("登录 %s 未返回 token: %s", username, rr.Body.String())
// 结束当前代码块
	}
// 返回：return token
	return token
// 结束当前代码块
}

// 函数定义：func statusOK(t *testing.T, rr *httptest.ResponseRecorder, want int, ctx string) {
func statusOK(t *testing.T, rr *httptest.ResponseRecorder, want int, ctx string) {
// 执行：t.Helper()
	t.Helper()
// 条件判断：if rr.Code != want {
	if rr.Code != want {
// 执行：t.Fatalf("%s: 期望 %d 实际 %d body=%s", ctx, want, rr.Code, rr.Body.String())
		t.Fatalf("%s: 期望 %d 实际 %d body=%s", ctx, want, rr.Code, rr.Body.String())
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func paySignRaw(secret string, raw []byte) string {
func paySignRaw(secret string, raw []byte) string {
// 赋值/更新：mac := hmac.New(sha256.New, []byte(secret))
	mac := hmac.New(sha256.New, []byte(secret))
// 执行：mac.Write(raw)
	mac.Write(raw)
// 返回：return hex.EncodeToString(mac.Sum(nil))
	return hex.EncodeToString(mac.Sum(nil))
// 结束当前代码块
}

// compatPayNotify 走契约通道发送支付通知（HMAC 对原始请求体）。
func compatPayNotify(t *testing.T, r *gin.Engine, orderNo, txn string, amount int64, signHeader string) *httptest.ResponseRecorder {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：body := `{"order_no":"` + orderNo + `","channel_txn_no":"` + txn + `","amount":` + itoa(amount) + `}`
	body := `{"order_no":"` + orderNo + `","channel_txn_no":"` + txn + `","amount":` + itoa(amount) + `}`
// 赋值/更新：sig := paySignRaw(testCfg.Shop.PayNotifySecret, []byte(body))
	sig := paySignRaw(testCfg.Shop.PayNotifySecret, []byte(body))
// 条件判断：if signHeader == "bad" {
	if signHeader == "bad" {
// 赋值/更新：sig = "deadbeef"
		sig = "deadbeef"
// 结束当前代码块
	}
// 返回：return doJSON(r, "POST", "/api/pay/notify", body, "", map[string]string{"X-Channel-Signature": sig})
	return doJSON(r, "POST", "/api/pay/notify", body, "", map[string]string{"X-Channel-Signature": sig})
// 结束当前代码块
}

// 函数定义：func itoa(v int64) string {
func itoa(v int64) string {
// 返回：return json.Number(jsonInt(v)).String()
	return json.Number(jsonInt(v)).String()
// 结束当前代码块
}

// 函数定义：func jsonInt(v int64) string {
func jsonInt(v int64) string {
// 赋值/更新：b, _ := json.Marshal(v)
	b, _ := json.Marshal(v)
// 返回：return string(b)
	return string(b)
// 结束当前代码块
}

// 创建一笔已支付订单，返回订单号与金额。
func paidOrder(t *testing.T, r *gin.Engine, buyerToken, txn string) (string, int64) {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：rr := doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":2}`, buyerToken, nil)
	rr := doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":2}`, buyerToken, nil)
// 执行：statusOK(t, rr, http.StatusOK, "下单")
	statusOK(t, rr, http.StatusOK, "下单")
// 赋值/更新：res := parseBody(t, rr)
	res := parseBody(t, rr)
// 赋值/更新：orderNo, _ := res["order_no"].(string)
	orderNo, _ := res["order_no"].(string)
// 赋值/更新：amount := int64(numOf(res["amount"]))
	amount := int64(numOf(res["amount"]))
// 条件判断：if orderNo == "" || amount <= 0 {
	if orderNo == "" || amount <= 0 {
// 执行：t.Fatalf("下单响应缺少字段: %s", rr.Body.String())
		t.Fatalf("下单响应缺少字段: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：rr = compatPayNotify(t, r, orderNo, txn, amount, "")
	rr = compatPayNotify(t, r, orderNo, txn, amount, "")
// 执行：statusOK(t, rr, http.StatusOK, "支付通知")
	statusOK(t, rr, http.StatusOK, "支付通知")
// 赋值/更新：res = parseBody(t, rr)
	res = parseBody(t, rr)
// 条件判断：if res["status"] != "SUCCESS" {
	if res["status"] != "SUCCESS" {
// 执行：t.Fatalf("支付通知未受理: %s", rr.Body.String())
		t.Fatalf("支付通知未受理: %s", rr.Body.String())
// 结束当前代码块
	}
// 返回：return orderNo, amount
	return orderNo, amount
// 结束当前代码块
}

// 函数定义：func numOf(v interface{}) float64 {
func numOf(v interface{}) float64 {
// 条件判断：if f, ok := v.(float64); ok {
	if f, ok := v.(float64); ok {
// 返回：return f
		return f
// 结束当前代码块
	}
// 返回：return 0
	return 0
// 结束当前代码块
}

// 函数定义：func TestPublicPagesAndAuth(t *testing.T) {
func TestPublicPagesAndAuth(t *testing.T) {
// 赋值/更新：r, _ := newServer(t, true)
	r, _ := newServer(t, true)

// 赋值/更新：rr := doJSON(r, "GET", "/api/health", "", "", nil)
	rr := doJSON(r, "GET", "/api/health", "", "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "健康检查")
	statusOK(t, rr, http.StatusOK, "健康检查")
// 条件判断：if res := parseBody(t, rr); res["status"] != "ok" {
	if res := parseBody(t, rr); res["status"] != "ok" {
// 执行：t.Fatalf("健康检查内容错误: %s", rr.Body.String())
		t.Fatalf("健康检查内容错误: %s", rr.Body.String())
// 结束当前代码块
	}

// 赋值/更新：rr = doJSON(r, "GET", "/api/users", "", "", nil)
	rr = doJSON(r, "GET", "/api/users", "", "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "账号列表")
	statusOK(t, rr, http.StatusOK, "账号列表")

// 赋值/更新：rr = doJSON(r, "GET", "/api/products", "", "", nil)
	rr = doJSON(r, "GET", "/api/products", "", "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "商品列表")
	statusOK(t, rr, http.StatusOK, "商品列表")
// 赋值/更新：res := parseBody(t, rr)
	res := parseBody(t, rr)
// 条件判断：if list, ok := res["products"].([]interface{}); !ok || len(list) != 2 {
	if list, ok := res["products"].([]interface{}); !ok || len(list) != 2 {
// 执行：t.Fatalf("在售商品应为 2 个: %s", rr.Body.String())
		t.Fatalf("在售商品应为 2 个: %s", rr.Body.String())
// 结束当前代码块
	}

// 赋值/更新：rr = doJSON(r, "GET", "/shop", "", "", nil)
	rr = doJSON(r, "GET", "/shop", "", "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "内置页面")
	statusOK(t, rr, http.StatusOK, "内置页面")
// 条件判断：if ct := rr.Header().Get("Content-Type"); ct == "" {
	if ct := rr.Header().Get("Content-Type"); ct == "" {
// 执行：t.Fatal("内置页面缺少 Content-Type")
		t.Fatal("内置页面缺少 Content-Type")
// 结束当前代码块
	}

	// 登录：成功 / 密码错误 / 员工契约别名
	rr = doJSON(r, "POST", "/api/login", `{"username":"buyer1","password":"shop123"}`, "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "买家登录")
	statusOK(t, rr, http.StatusOK, "买家登录")
// 赋值/更新：rr = doJSON(r, "POST", "/api/login", `{"username":"buyer1","password":"wrong"}`, "", nil)
	rr = doJSON(r, "POST", "/api/login", `{"username":"buyer1","password":"wrong"}`, "", nil)
// 条件判断：if rr.Code < 400 {
	if rr.Code < 400 {
// 执行：t.Fatalf("错误密码应拒绝: %d", rr.Code)
		t.Fatalf("错误密码应拒绝: %d", rr.Code)
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "POST", "/api/login", `{"username":"ops","password":"shop123"}`, "", nil)
	rr = doJSON(r, "POST", "/api/login", `{"username":"ops","password":"shop123"}`, "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "契约员工别名登录")
	statusOK(t, rr, http.StatusOK, "契约员工别名登录")
// 条件判断：if res := parseBody(t, rr); res["role"] != "ops" {
	if res := parseBody(t, rr); res["role"] != "ops" {
// 执行：t.Fatalf("别名登录角色错误: %s", rr.Body.String())
		t.Fatalf("别名登录角色错误: %s", rr.Body.String())
// 结束当前代码块
	}

	// 未登录 / 伪造令牌
	rr = doJSON(r, "GET", "/api/me", "", "", nil)
// 执行：statusOK(t, rr, http.StatusUnauthorized, "未登录访问")
	statusOK(t, rr, http.StatusUnauthorized, "未登录访问")
// 赋值/更新：rr = doJSON(r, "GET", "/api/me", "", "not-a-token", nil)
	rr = doJSON(r, "GET", "/api/me", "", "not-a-token", nil)
// 执行：statusOK(t, rr, http.StatusUnauthorized, "伪造令牌")
	statusOK(t, rr, http.StatusUnauthorized, "伪造令牌")

// 赋值/更新：token := login(t, r, "buyer1")
	token := login(t, r, "buyer1")
// 赋值/更新：rr = doJSON(r, "GET", "/api/me", "", token, nil)
	rr = doJSON(r, "GET", "/api/me", "", token, nil)
// 执行：statusOK(t, rr, http.StatusOK, "当前用户")
	statusOK(t, rr, http.StatusOK, "当前用户")
// 条件判断：if me, ok := parseBody(t, rr)["data"].(map[string]interface{}); !ok || me["username"] != "buyer1" {
	if me, ok := parseBody(t, rr)["data"].(map[string]interface{}); !ok || me["username"] != "buyer1" {
// 执行：t.Fatalf("/api/me 用户错误: %s", rr.Body.String())
		t.Fatalf("/api/me 用户错误: %s", rr.Body.String())
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestBuyerOrderAndPayFlow(t *testing.T) {
func TestBuyerOrderAndPayFlow(t *testing.T) {
// 赋值/更新：r, _ := newServer(t, false)
	r, _ := newServer(t, false)
// 赋值/更新：token := login(t, r, "buyer1")
	token := login(t, r, "buyer1")

// 赋值/更新：rr := doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":2}`, token, nil)
	rr := doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":2}`, token, nil)
// 执行：statusOK(t, rr, http.StatusOK, "下单")
	statusOK(t, rr, http.StatusOK, "下单")
// 赋值/更新：res := parseBody(t, rr)
	res := parseBody(t, rr)
// 赋值/更新：orderNo, _ := res["order_no"].(string)
	orderNo, _ := res["order_no"].(string)
// 条件判断：if orderNo == "" || res["status"] != "WAIT_PAY" {
	if orderNo == "" || res["status"] != "WAIT_PAY" {
// 执行：t.Fatalf("下单响应错误: %s", rr.Body.String())
		t.Fatalf("下单响应错误: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：amount := int64(numOf(res["amount"]))
	amount := int64(numOf(res["amount"]))
// 条件判断：if list, ok := res["data"].(map[string]interface{}); !ok || list["order"] == nil {
	if list, ok := res["data"].(map[string]interface{}); !ok || list["order"] == nil {
// 执行：t.Fatalf("下单响应缺少 data 副本: %s", rr.Body.String())
		t.Fatalf("下单响应缺少 data 副本: %s", rr.Body.String())
// 结束当前代码块
	}

// 赋值/更新：rr = doJSON(r, "GET", "/api/orders", "", token, nil)
	rr = doJSON(r, "GET", "/api/orders", "", token, nil)
// 执行：statusOK(t, rr, http.StatusOK, "我的订单")
	statusOK(t, rr, http.StatusOK, "我的订单")
// 赋值/更新：rr = doJSON(r, "GET", "/api/orders/"+orderNo, "", token, nil)
	rr = doJSON(r, "GET", "/api/orders/"+orderNo, "", token, nil)
// 执行：statusOK(t, rr, http.StatusOK, "订单详情")
	statusOK(t, rr, http.StatusOK, "订单详情")
// 赋值/更新：rr = doJSON(r, "GET", "/api/orders/NO-SUCH-ORDER", "", token, nil)
	rr = doJSON(r, "GET", "/api/orders/NO-SUCH-ORDER", "", token, nil)
// 执行：statusOK(t, rr, http.StatusNotFound, "未知订单详情")
	statusOK(t, rr, http.StatusNotFound, "未知订单详情")
// 赋值/更新：rr = doJSON(r, "GET", "/api/orders/"+orderNo, "", login(t, r, "buyer2"), nil)
	rr = doJSON(r, "GET", "/api/orders/"+orderNo, "", login(t, r, "buyer2"), nil)
// 执行：statusOK(t, rr, http.StatusForbidden, "他人订单详情")
	statusOK(t, rr, http.StatusForbidden, "他人订单详情")

// 赋值/更新：rr = doJSON(r, "GET", "/api/pay/pending", "", "", nil)
	rr = doJSON(r, "GET", "/api/pay/pending", "", "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "渠道待支付列表")
	statusOK(t, rr, http.StatusOK, "渠道待支付列表")

	// 有效支付通知（契约通道）
	rr = compatPayNotify(t, r, orderNo, "TXN-1", amount, "")
// 执行：statusOK(t, rr, http.StatusOK, "支付通知")
	statusOK(t, rr, http.StatusOK, "支付通知")
	// 重复通知 → 409 且不重复记款
	rr = compatPayNotify(t, r, orderNo, "TXN-1", amount, "")
// 条件判断：if rr.Code != http.StatusConflict {
	if rr.Code != http.StatusConflict {
// 执行：t.Fatalf("重复通知应 409: %d %s", rr.Code, rr.Body.String())
		t.Fatalf("重复通知应 409: %d %s", rr.Code, rr.Body.String())
// 结束当前代码块
	}
	// 验签失败 → 4xx
	rr = compatPayNotify(t, r, orderNo, "TXN-2", amount, "bad")
// 条件判断：if rr.Code < 400 {
	if rr.Code < 400 {
// 执行：t.Fatalf("验签失败应拒绝: %d", rr.Code)
		t.Fatalf("验签失败应拒绝: %d", rr.Code)
// 结束当前代码块
	}
	// 渠道受理付款：返回交易号，可用于后续通知
	rr = doJSON(r, "POST", "/api/pay/channel-pay", `{"order_no":"`+orderNo+`"}`, "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "渠道受理付款")
	statusOK(t, rr, http.StatusOK, "渠道受理付款")
// 条件判断：if res := parseBody(t, rr); res["status"] != "SUCCESS" {
	if res := parseBody(t, rr); res["status"] != "SUCCESS" {
// 执行：t.Fatalf("渠道受理响应错误: %s", rr.Body.String())
		t.Fatalf("渠道受理响应错误: %s", rr.Body.String())
// 结束当前代码块
	}

	// 取消待付款订单
	rr = doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-B","quantity":1}`, token, nil)
// 赋值/更新：res = parseBody(t, rr)
	res = parseBody(t, rr)
// 赋值/更新：no2, _ := res["order_no"].(string)
	no2, _ := res["order_no"].(string)
// 赋值/更新：rr = doJSON(r, "POST", "/api/orders/"+no2+"/cancel", "", token, nil)
	rr = doJSON(r, "POST", "/api/orders/"+no2+"/cancel", "", token, nil)
// 执行：statusOK(t, rr, http.StatusOK, "取消订单")
	statusOK(t, rr, http.StatusOK, "取消订单")
// 赋值/更新：rr = doJSON(r, "POST", "/api/orders/"+no2+"/cancel", "", token, nil)
	rr = doJSON(r, "POST", "/api/orders/"+no2+"/cancel", "", token, nil)
// 条件判断：if rr.Code < 400 {
	if rr.Code < 400 {
// 执行：t.Fatalf("重复取消应拒绝: %d", rr.Code)
		t.Fatalf("重复取消应拒绝: %d", rr.Code)
// 结束当前代码块
	}

	// 契约独立退款申请路径（已付款订单）
	rr = doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-B","quantity":1}`, token, nil)
// 赋值/更新：res = parseBody(t, rr)
	res = parseBody(t, rr)
// 赋值/更新：no3, _ := res["order_no"].(string)
	no3, _ := res["order_no"].(string)
// 赋值/更新：amount3 := int64(numOf(res["amount"]))
	amount3 := int64(numOf(res["amount"]))
// 赋值/更新：rr = compatPayNotify(t, r, no3, "TXN-REFUND-APPLY", amount3, "")
	rr = compatPayNotify(t, r, no3, "TXN-REFUND-APPLY", amount3, "")
// 执行：statusOK(t, rr, http.StatusOK, "第三单支付")
	statusOK(t, rr, http.StatusOK, "第三单支付")
// 赋值/更新：rr = doJSON(r, "POST", "/api/orders/"+no3+"/refund-apply", `{"reason":"契约退款"}`, token, nil)
	rr = doJSON(r, "POST", "/api/orders/"+no3+"/refund-apply", `{"reason":"契约退款"}`, token, nil)
// 执行：statusOK(t, rr, http.StatusOK, "契约退款申请")
	statusOK(t, rr, http.StatusOK, "契约退款申请")
// 条件判断：if res := parseBody(t, rr); res["refund_id"] == nil || res["refund_id"] == "" {
	if res := parseBody(t, rr); res["refund_id"] == nil || res["refund_id"] == "" {
// 执行：t.Fatalf("退款申请缺少 refund_id: %s", rr.Body.String())
		t.Fatalf("退款申请缺少 refund_id: %s", rr.Body.String())
// 结束当前代码块
	}
	// 未登录申请退款 → 401
	rr = doJSON(r, "POST", "/api/orders/"+no3+"/refund-apply", `{"reason":"x"}`, "", nil)
// 执行：statusOK(t, rr, http.StatusUnauthorized, "未登录退款申请")
	statusOK(t, rr, http.StatusUnauthorized, "未登录退款申请")
// 结束当前代码块
}

// 函数定义：func TestWarehouseShipAndRoleGuard(t *testing.T) {
func TestWarehouseShipAndRoleGuard(t *testing.T) {
// 赋值/更新：r, _ := newServer(t, false)
	r, _ := newServer(t, false)
// 赋值/更新：buyer := login(t, r, "buyer1")
	buyer := login(t, r, "buyer1")
// 赋值/更新：wh := login(t, r, "wh1")
	wh := login(t, r, "wh1")
// 赋值/更新：orderNo, amount := paidOrder(t, r, buyer, "TXN-SHIP-1")
	orderNo, amount := paidOrder(t, r, buyer, "TXN-SHIP-1")

// 赋值/更新：rr := doJSON(r, "GET", "/api/warehouse/orders", "", wh, nil)
	rr := doJSON(r, "GET", "/api/warehouse/orders", "", wh, nil)
// 执行：statusOK(t, rr, http.StatusOK, "可发货列表(契约)")
	statusOK(t, rr, http.StatusOK, "可发货列表(契约)")
// 赋值/更新：rr = doJSON(r, "GET", "/api/warehouse/shipable", "", wh, nil)
	rr = doJSON(r, "GET", "/api/warehouse/shipable", "", wh, nil)
// 执行：statusOK(t, rr, http.StatusOK, "可发货列表(原型)")
	statusOK(t, rr, http.StatusOK, "可发货列表(原型)")

// 赋值/更新：rr = doJSON(r, "POST", "/api/warehouse/orders/"+orderNo+"/ship", `{"carrier":"顺丰","tracking_no":"SF100"}`, wh,...
	rr = doJSON(r, "POST", "/api/warehouse/orders/"+orderNo+"/ship", `{"carrier":"顺丰","tracking_no":"SF100"}`, wh, nil)
// 执行：statusOK(t, rr, http.StatusOK, "发货登记")
	statusOK(t, rr, http.StatusOK, "发货登记")
// 赋值/更新：res := parseBody(t, rr)
	res := parseBody(t, rr)
// 条件判断：if res["status"] != "SHIPPED" || res["tracking_no"] != "SF100" {
	if res["status"] != "SHIPPED" || res["tracking_no"] != "SF100" {
// 执行：t.Fatalf("发货响应错误: %s", rr.Body.String())
		t.Fatalf("发货响应错误: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "POST", "/api/warehouse/orders/"+orderNo+"/ship", `{"carrier":"顺丰","tracking_no":"SF101"}`, wh,...
	rr = doJSON(r, "POST", "/api/warehouse/orders/"+orderNo+"/ship", `{"carrier":"顺丰","tracking_no":"SF101"}`, wh, nil)
// 条件判断：if rr.Code < 400 {
	if rr.Code < 400 {
// 执行：t.Fatalf("重复发货应拒绝: %d", rr.Code)
		t.Fatalf("重复发货应拒绝: %d", rr.Code)
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "GET", "/api/warehouse/shipments", "", wh, nil)
	rr = doJSON(r, "GET", "/api/warehouse/shipments", "", wh, nil)
// 执行：statusOK(t, rr, http.StatusOK, "发货记录")
	statusOK(t, rr, http.StatusOK, "发货记录")

	// 原型发货入口（POST /api/warehouse/ship，request.Ship 结构）
	orderNo2, amount2 := paidOrder(t, r, buyer, "TXN-SHIP-2")
// 赋值/更新：rr = doJSON(r, "POST", "/api/warehouse/ship",
	rr = doJSON(r, "POST", "/api/warehouse/ship",
// 执行：`{"orderNo":"`+orderNo2+`","carrier":"SF","trackingNo":"SF300"}`, wh, nil)
		`{"orderNo":"`+orderNo2+`","carrier":"SF","trackingNo":"SF300"}`, wh, nil)
// 执行：statusOK(t, rr, http.StatusOK, "原型发货登记")
	statusOK(t, rr, http.StatusOK, "原型发货登记")
// 条件判断：if data, ok := parseBody(t, rr)["data"].(map[string]interface{}); !ok || data["status"] != "shipped" {
	if data, ok := parseBody(t, rr)["data"].(map[string]interface{}); !ok || data["status"] != "shipped" {
// 执行：t.Fatalf("原型发货状态错误: %s", rr.Body.String())
		t.Fatalf("原型发货状态错误: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：_ = amount2
	_ = amount2

	// 买家访问仓库接口 → 403
	rr = doJSON(r, "GET", "/api/warehouse/orders", "", buyer, nil)
// 执行：statusOK(t, rr, http.StatusForbidden, "买家访问仓库接口")
	statusOK(t, rr, http.StatusForbidden, "买家访问仓库接口")
// 赋值/更新：_ = amount
	_ = amount
// 结束当前代码块
}

// 函数定义：func TestOpsMaintainAndReport(t *testing.T) {
func TestOpsMaintainAndReport(t *testing.T) {
// 赋值/更新：r, _ := newServer(t, false)
	r, _ := newServer(t, false)
// 赋值/更新：buyer := login(t, r, "buyer1")
	buyer := login(t, r, "buyer1")
// 赋值/更新：ops := login(t, r, "ops1")
	ops := login(t, r, "ops1")
// 赋值/更新：orderNo, _ := paidOrder(t, r, buyer, "TXN-OPS-1")
	orderNo, _ := paidOrder(t, r, buyer, "TXN-OPS-1")

// 赋值/更新：rr := doJSON(r, "PATCH", "/api/ops/products/SKU-A", `{"price_cents":5200}`, ops, nil)
	rr := doJSON(r, "PATCH", "/api/ops/products/SKU-A", `{"price_cents":5200}`, ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "调价")
	statusOK(t, rr, http.StatusOK, "调价")
// 赋值/更新：rr = doJSON(r, "PATCH", "/api/ops/products/SKU-A", `{"available":9}`, ops, nil)
	rr = doJSON(r, "PATCH", "/api/ops/products/SKU-A", `{"available":9}`, ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "覆盖可售")
	statusOK(t, rr, http.StatusOK, "覆盖可售")
// 条件判断：if res := parseBody(t, rr); int(numOf(res["available"])) != 9 {
	if res := parseBody(t, rr); int(numOf(res["available"])) != 9 {
// 执行：t.Fatalf("覆盖可售失败: %s", rr.Body.String())
		t.Fatalf("覆盖可售失败: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "PATCH", "/api/ops/products/SKU-A", `{}`, ops, nil)
	rr = doJSON(r, "PATCH", "/api/ops/products/SKU-A", `{}`, ops, nil)
// 条件判断：if rr.Code < 400 {
	if rr.Code < 400 {
// 执行：t.Fatalf("空更新应拒绝: %d", rr.Code)
		t.Fatalf("空更新应拒绝: %d", rr.Code)
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "PATCH", "/api/ops/products/SKU-A", `{"price_cents":5200}`, buyer, nil)
	rr = doJSON(r, "PATCH", "/api/ops/products/SKU-A", `{"price_cents":5200}`, buyer, nil)
// 执行：statusOK(t, rr, http.StatusForbidden, "买家维护商品")
	statusOK(t, rr, http.StatusForbidden, "买家维护商品")

	// 原型维护入口：productId + priceCents（productId 从商品列表读取，避免依赖自增起始值）
	rr = doJSON(r, "GET", "/api/products", "", "", nil)
// 赋值/更新：plist := parseBody(t, rr)["products"].([]interface{})
	plist := parseBody(t, rr)["products"].([]interface{})
// 赋值/更新：pid := json.Number(jsonInt(int64(numOf(plist[0].(map[string]interface{})["productId"])))).String()
	pid := json.Number(jsonInt(int64(numOf(plist[0].(map[string]interface{})["productId"])))).String()
// 赋值/更新：rr = doJSON(r, "PATCH", "/api/ops/products", `{"productId":`+pid+`,"priceCents":5300}`, ops, nil)
	rr = doJSON(r, "PATCH", "/api/ops/products", `{"productId":`+pid+`,"priceCents":5300}`, ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "原型维护商品")
	statusOK(t, rr, http.StatusOK, "原型维护商品")

// 赋值/更新：rr = doJSON(r, "GET", "/api/ops/summary", "", ops, nil)
	rr = doJSON(r, "GET", "/api/ops/summary", "", ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "交易汇总")
	statusOK(t, rr, http.StatusOK, "交易汇总")
// 赋值/更新：rr = doJSON(r, "GET", "/api/ops/orders", "", ops, nil)
	rr = doJSON(r, "GET", "/api/ops/orders", "", ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "订单列表")
	statusOK(t, rr, http.StatusOK, "订单列表")
// 赋值/更新：rr = doJSON(r, "GET", "/api/ops/orders/"+orderNo, "", ops, nil)
	rr = doJSON(r, "GET", "/api/ops/orders/"+orderNo, "", ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "订单全链路")
	statusOK(t, rr, http.StatusOK, "订单全链路")
// 赋值/更新：rr = doJSON(r, "GET", "/api/ops/notices?orderNo="+orderNo, "", ops, nil)
	rr = doJSON(r, "GET", "/api/ops/notices?orderNo="+orderNo, "", ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "通知查询")
	statusOK(t, rr, http.StatusOK, "通知查询")
// 赋值/更新：rr = doJSON(r, "GET", "/api/ops/notices?result=rejected", "", ops, nil)
	rr = doJSON(r, "GET", "/api/ops/notices?result=rejected", "", ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "通知按结果过滤")
	statusOK(t, rr, http.StatusOK, "通知按结果过滤")
// 赋值/更新：rr = doJSON(r, "GET", "/api/ops/abnormals", "", ops, nil)
	rr = doJSON(r, "GET", "/api/ops/abnormals", "", ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "异常列表")
	statusOK(t, rr, http.StatusOK, "异常列表")
// 赋值/更新：rr = doJSON(r, "GET", "/api/ops/ledger?productId=1", "", ops, nil)
	rr = doJSON(r, "GET", "/api/ops/ledger?productId=1", "", ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "库存流水")
	statusOK(t, rr, http.StatusOK, "库存流水")
// 赋值/更新：rr = doJSON(r, "GET", "/api/ops/trail?orderNo="+orderNo, "", ops, nil)
	rr = doJSON(r, "GET", "/api/ops/trail?orderNo="+orderNo, "", ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "操作轨迹")
	statusOK(t, rr, http.StatusOK, "操作轨迹")
// 赋值/更新：rr = doJSON(r, "GET", "/api/ops/payments?order_no="+orderNo, "", ops, nil)
	rr = doJSON(r, "GET", "/api/ops/payments?order_no="+orderNo, "", ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "支付记录查询")
	statusOK(t, rr, http.StatusOK, "支付记录查询")
// 条件判断：if res := parseBody(t, rr); len(res["payments"].([]interface{})) == 0 {
	if res := parseBody(t, rr); len(res["payments"].([]interface{})) == 0 {
// 执行：t.Fatalf("支付记录为空: %s", rr.Body.String())
		t.Fatalf("支付记录为空: %s", rr.Body.String())
// 结束当前代码块
	}
	// 买家访问运营接口 → 403
	rr = doJSON(r, "GET", "/api/ops/summary", "", buyer, nil)
// 执行：statusOK(t, rr, http.StatusForbidden, "买家访问运营接口")
	statusOK(t, rr, http.StatusForbidden, "买家访问运营接口")
// 结束当前代码块
}

// 函数定义：func TestRefundFlow(t *testing.T) {
func TestRefundFlow(t *testing.T) {
// 赋值/更新：r, ch := newServer(t, true)
	r, ch := newServer(t, true)
// 赋值/更新：buyer := login(t, r, "buyer1")
	buyer := login(t, r, "buyer1")
// 赋值/更新：wh := login(t, r, "wh1")
	wh := login(t, r, "wh1")
// 赋值/更新：cs := login(t, r, "cs1")
	cs := login(t, r, "cs1")

	// 场景一：申请 → 发货被拒 → 客服拒绝 → 恢复发货
	orderNo, _ := paidOrder(t, r, buyer, "TXN-RF-1")
// 赋值/更新：rr := doJSON(r, "POST", "/api/orders/"+orderNo+"/refund", `{"reason":"不想要了"}`, buyer, nil)
	rr := doJSON(r, "POST", "/api/orders/"+orderNo+"/refund", `{"reason":"不想要了"}`, buyer, nil)
// 执行：statusOK(t, rr, http.StatusOK, "退款申请")
	statusOK(t, rr, http.StatusOK, "退款申请")
// 赋值/更新：rr = doJSON(r, "GET", "/api/orders/"+orderNo+"/refund", "", buyer, nil)
	rr = doJSON(r, "GET", "/api/orders/"+orderNo+"/refund", "", buyer, nil)
// 执行：statusOK(t, rr, http.StatusOK, "买家查退款")
	statusOK(t, rr, http.StatusOK, "买家查退款")
// 赋值/更新：rr = doJSON(r, "POST", "/api/warehouse/orders/"+orderNo+"/ship", `{"carrier":"SF","tracking_no":"SF200"}`, wh,...
	rr = doJSON(r, "POST", "/api/warehouse/orders/"+orderNo+"/ship", `{"carrier":"SF","tracking_no":"SF200"}`, wh, nil)
// 条件判断：if rr.Code < 400 {
	if rr.Code < 400 {
// 执行：t.Fatalf("退款申请期发货应被拒: %d", rr.Code)
		t.Fatalf("退款申请期发货应被拒: %d", rr.Code)
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "GET", "/api/cs/refunds?status=pending", "", cs, nil)
	rr = doJSON(r, "GET", "/api/cs/refunds?status=pending", "", cs, nil)
// 执行：statusOK(t, rr, http.StatusOK, "客服退款列表")
	statusOK(t, rr, http.StatusOK, "客服退款列表")
// 赋值/更新：res := parseBody(t, rr)
	res := parseBody(t, rr)
// 赋值/更新：list := res["refunds"].([]interface{})
	list := res["refunds"].([]interface{})
// 赋值/更新：first := list[0].(map[string]interface{})
	first := list[0].(map[string]interface{})
// 赋值/更新：refundID, _ := first["id"].(string)
	refundID, _ := first["id"].(string)
// 赋值/更新：rr = doJSON(r, "POST", "/api/cs/refunds/"+refundID+"/decide", `{"approve":false,"reason":"凭证不足"}`, cs, nil)
	rr = doJSON(r, "POST", "/api/cs/refunds/"+refundID+"/decide", `{"approve":false,"reason":"凭证不足"}`, cs, nil)
// 执行：statusOK(t, rr, http.StatusOK, "客服拒绝退款")
	statusOK(t, rr, http.StatusOK, "客服拒绝退款")
// 条件判断：if res := parseBody(t, rr); res["status"] != "REJECTED" {
	if res := parseBody(t, rr); res["status"] != "REJECTED" {
// 执行：t.Fatalf("拒绝后状态错误: %s", rr.Body.String())
		t.Fatalf("拒绝后状态错误: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "POST", "/api/warehouse/orders/"+orderNo+"/ship", `{"carrier":"SF","tracking_no":"SF201"}`, wh,...
	rr = doJSON(r, "POST", "/api/warehouse/orders/"+orderNo+"/ship", `{"carrier":"SF","tracking_no":"SF201"}`, wh, nil)
// 执行：statusOK(t, rr, http.StatusOK, "拒绝后可发货")
	statusOK(t, rr, http.StatusOK, "拒绝后可发货")

	// 场景二：申请 → 同意 → 渠道失败 → 重试成功
	orderNo2, _ := paidOrder(t, r, buyer, "TXN-RF-2")
// 赋值/更新：rr = doJSON(r, "POST", "/api/orders/"+orderNo2+"/refund", `{"reason":"换货"}`, buyer, nil)
	rr = doJSON(r, "POST", "/api/orders/"+orderNo2+"/refund", `{"reason":"换货"}`, buyer, nil)
// 执行：statusOK(t, rr, http.StatusOK, "第二笔退款申请")
	statusOK(t, rr, http.StatusOK, "第二笔退款申请")
// 赋值/更新：rr = doJSON(r, "POST", "/api/pay/channel-mode", `{"refund_result":"FAIL"}`, "", nil)
	rr = doJSON(r, "POST", "/api/pay/channel-mode", `{"refund_result":"FAIL"}`, "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "渠道切失败")
	statusOK(t, rr, http.StatusOK, "渠道切失败")
// 赋值/更新：rr = doJSON(r, "GET", "/api/cs/refunds?status=pending", "", cs, nil)
	rr = doJSON(r, "GET", "/api/cs/refunds?status=pending", "", cs, nil)
// 赋值/更新：list = parseBody(t, rr)["refunds"].([]interface{})
	list = parseBody(t, rr)["refunds"].([]interface{})
// 赋值/更新：rid2 := list[0].(map[string]interface{})["id"].(string)
	rid2 := list[0].(map[string]interface{})["id"].(string)
	// REF-03（增量续写 R2 修复后）：审核必须携带非空原因，空原因会被拒绝
	rr = doJSON(r, "POST", "/api/cs/refunds/"+rid2+"/decide", `{"approve":true}`, cs, nil)
// 条件判断：if rr.Code == http.StatusOK {
	if rr.Code == http.StatusOK {
// 执行：t.Fatalf("空原因审核应被拒绝（REF-03）: %s", rr.Body.String())
		t.Fatalf("空原因审核应被拒绝（REF-03）: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "POST", "/api/cs/refunds/"+rid2+"/decide", `{"approve":true,"reason":"凭证核验通过"}`, cs, nil)
	rr = doJSON(r, "POST", "/api/cs/refunds/"+rid2+"/decide", `{"approve":true,"reason":"凭证核验通过"}`, cs, nil)
// 执行：statusOK(t, rr, http.StatusOK, "同意后渠道失败仍 2xx")
	statusOK(t, rr, http.StatusOK, "同意后渠道失败仍 2xx")
// 条件判断：if res := parseBody(t, rr); res["status"] != "FAILED" {
	if res := parseBody(t, rr); res["status"] != "FAILED" {
// 执行：t.Fatalf("渠道失败状态错误: %s", rr.Body.String())
		t.Fatalf("渠道失败状态错误: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "POST", "/api/cs/refunds/"+rid2+"/retry", `{}`, cs, nil)
	rr = doJSON(r, "POST", "/api/cs/refunds/"+rid2+"/retry", `{}`, cs, nil)
// 条件判断：if rr.Code == http.StatusOK {
	if rr.Code == http.StatusOK {
// 执行：t.Fatalf("渠道仍失败时重试不应成功: %s", rr.Body.String())
		t.Fatalf("渠道仍失败时重试不应成功: %s", rr.Body.String())
// 结束当前代码块
	}
	// 原型重试入口（POST /api/cs/refunds/retry，request.RefundRetry 结构）同样返回冲突
	rr = doJSON(r, "POST", "/api/cs/refunds/retry", `{"refundNo":"`+rid2+`"}`, cs, nil)
// 条件判断：if rr.Code == http.StatusOK {
	if rr.Code == http.StatusOK {
// 执行：t.Fatalf("原型重试仍失败时不应成功: %s", rr.Body.String())
		t.Fatalf("原型重试仍失败时不应成功: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "POST", "/api/pay/channel-mode", `{"refund_result":"SUCCESS"}`, "", nil)
	rr = doJSON(r, "POST", "/api/pay/channel-mode", `{"refund_result":"SUCCESS"}`, "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "渠道切成功")
	statusOK(t, rr, http.StatusOK, "渠道切成功")
// 赋值/更新：rr = doJSON(r, "POST", "/api/cs/refunds/"+rid2+"/retry", `{}`, cs, nil)
	rr = doJSON(r, "POST", "/api/cs/refunds/"+rid2+"/retry", `{}`, cs, nil)
// 执行：statusOK(t, rr, http.StatusOK, "重试成功")
	statusOK(t, rr, http.StatusOK, "重试成功")
// 条件判断：if res := parseBody(t, rr); res["status"] != "REFUNDED" {
	if res := parseBody(t, rr); res["status"] != "REFUNDED" {
// 执行：t.Fatalf("重试后状态错误: %s", rr.Body.String())
		t.Fatalf("重试后状态错误: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "GET", "/api/orders/"+orderNo2+"/refund", "", buyer, nil)
	rr = doJSON(r, "GET", "/api/orders/"+orderNo2+"/refund", "", buyer, nil)
// 执行：statusOK(t, rr, http.StatusOK, "买家查退款终态")
	statusOK(t, rr, http.StatusOK, "买家查退款终态")
// 条件判断：if data, ok := parseBody(t, rr)["data"].(map[string]interface{}); !ok || data["channelStatus"] != "success" {
	if data, ok := parseBody(t, rr)["data"].(map[string]interface{}); !ok || data["channelStatus"] != "success" {
// 执行：t.Fatalf("买家可见退款终态错误: %s", rr.Body.String())
		t.Fatalf("买家可见退款终态错误: %s", rr.Body.String())
// 结束当前代码块
	}

	// 场景三：原型审核入口 RefundAudit（同单重复审核 → 409）
	orderNo3, _ := paidOrder(t, r, buyer, "TXN-RF-3")
// 赋值/更新：rr = doJSON(r, "POST", "/api/orders/"+orderNo3+"/refund", `{"reason":"测试"}`, buyer, nil)
	rr = doJSON(r, "POST", "/api/orders/"+orderNo3+"/refund", `{"reason":"测试"}`, buyer, nil)
// 执行：statusOK(t, rr, http.StatusOK, "第三笔退款申请")
	statusOK(t, rr, http.StatusOK, "第三笔退款申请")
// 赋值/更新：rr = doJSON(r, "POST", "/api/cs/refunds/audit", `{"orderNo":"`+orderNo3+`","approve":true,"reason":"ok"}`, cs,...
	rr = doJSON(r, "POST", "/api/cs/refunds/audit", `{"orderNo":"`+orderNo3+`","approve":true,"reason":"ok"}`, cs, nil)
// 执行：statusOK(t, rr, http.StatusOK, "原型审核同意")
	statusOK(t, rr, http.StatusOK, "原型审核同意")
// 赋值/更新：rr = doJSON(r, "POST", "/api/cs/refunds/audit", `{"orderNo":"`+orderNo3+`","approve":true,"reason":"ok"}`, cs,...
	rr = doJSON(r, "POST", "/api/cs/refunds/audit", `{"orderNo":"`+orderNo3+`","approve":true,"reason":"ok"}`, cs, nil)
// 条件判断：if rr.Code < 400 {
	if rr.Code < 400 {
// 执行：t.Fatalf("重复审核应拒绝: %d", rr.Code)
		t.Fatalf("重复审核应拒绝: %d", rr.Code)
// 结束当前代码块
	}

	// 渠道状态查询 + 非法模式
	rr = doJSON(r, "GET", "/api/pay/channel", "", "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "渠道状态")
	statusOK(t, rr, http.StatusOK, "渠道状态")
// 赋值/更新：rr = doJSON(r, "POST", "/api/pay/channel-mode", `{"refund_result":"BAD"}`, "", nil)
	rr = doJSON(r, "POST", "/api/pay/channel-mode", `{"refund_result":"BAD"}`, "", nil)
// 条件判断：if rr.Code < 400 {
	if rr.Code < 400 {
// 执行：t.Fatalf("非法渠道模式应拒绝: %d", rr.Code)
		t.Fatalf("非法渠道模式应拒绝: %d", rr.Code)
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "POST", "/api/cs/channel/mode", `{"refundMode":"fail"}`, cs, nil)
	rr = doJSON(r, "POST", "/api/cs/channel/mode", `{"refundMode":"fail"}`, cs, nil)
// 执行：statusOK(t, rr, http.StatusOK, "原型渠道模式切换")
	statusOK(t, rr, http.StatusOK, "原型渠道模式切换")
// 执行：ch.SetRefundMode("success")
	ch.SetRefundMode("success")
// 赋值/更新：_ = ch
	_ = ch
// 结束当前代码块
}

// 函数定义：func TestLatePaymentAbnormal(t *testing.T) {
func TestLatePaymentAbnormal(t *testing.T) {
// 赋值/更新：r, _ := newServer(t, true)
	r, _ := newServer(t, true)
// 赋值/更新：buyer := login(t, r, "buyer1")
	buyer := login(t, r, "buyer1")
// 赋值/更新：cs := login(t, r, "cs1")
	cs := login(t, r, "cs1")
// 赋值/更新：ops := login(t, r, "ops1")
	ops := login(t, r, "ops1")

// 赋值/更新：rr := doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":1}`, buyer, nil)
	rr := doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":1}`, buyer, nil)
// 赋值/更新：res := parseBody(t, rr)
	res := parseBody(t, rr)
// 赋值/更新：orderNo, _ := res["order_no"].(string)
	orderNo, _ := res["order_no"].(string)
// 赋值/更新：amount := int64(numOf(res["amount"]))
	amount := int64(numOf(res["amount"]))

// 赋值/更新：rr = doJSON(r, "POST", "/api/test/expire-orders", `{}`, "", nil)
	rr = doJSON(r, "POST", "/api/test/expire-orders", `{}`, "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "测试钩子超时关单")
	statusOK(t, rr, http.StatusOK, "测试钩子超时关单")
// 条件判断：if res := parseBody(t, rr); int(numOf(res["closed"])) != 1 {
	if res := parseBody(t, rr); int(numOf(res["closed"])) != 1 {
// 执行：t.Fatalf("应关闭 1 笔: %s", rr.Body.String())
		t.Fatalf("应关闭 1 笔: %s", rr.Body.String())
// 结束当前代码块
	}

	// 迟到通知 → ABNORMAL
	rr = compatPayNotify(t, r, orderNo, "TXN-LATE-1", amount, "")
// 执行：statusOK(t, rr, http.StatusOK, "迟到通知")
	statusOK(t, rr, http.StatusOK, "迟到通知")
// 条件判断：if res := parseBody(t, rr); res["status"] != "ABNORMAL" {
	if res := parseBody(t, rr); res["status"] != "ABNORMAL" {
// 执行：t.Fatalf("迟到通知应转异常: %s", rr.Body.String())
		t.Fatalf("迟到通知应转异常: %s", rr.Body.String())
// 结束当前代码块
	}

// 赋值/更新：rr = doJSON(r, "GET", "/api/cs/payment-exceptions", "", cs, nil)
	rr = doJSON(r, "GET", "/api/cs/payment-exceptions", "", cs, nil)
// 执行：statusOK(t, rr, http.StatusOK, "异常列表")
	statusOK(t, rr, http.StatusOK, "异常列表")
// 赋值/更新：list := parseBody(t, rr)["exceptions"].([]interface{})
	list := parseBody(t, rr)["exceptions"].([]interface{})
// 赋值/更新：ab := list[0].(map[string]interface{})
	ab := list[0].(map[string]interface{})
// 赋值/更新：abID := json.Number(jsonInt(int64(numOf(ab["id"])))).String()
	abID := json.Number(jsonInt(int64(numOf(ab["id"])))).String()
// 赋值/更新：rr = doJSON(r, "POST", "/api/cs/payment-exceptions/"+abID+"/refund", `{"reason":"核实迟到支付"}`, cs, nil)
	rr = doJSON(r, "POST", "/api/cs/payment-exceptions/"+abID+"/refund", `{"reason":"核实迟到支付"}`, cs, nil)
// 执行：statusOK(t, rr, http.StatusOK, "异常退款")
	statusOK(t, rr, http.StatusOK, "异常退款")
// 赋值/更新：rr = doJSON(r, "POST", "/api/cs/payment-exceptions/"+abID+"/refund", `{"reason":"再次处理"}`, cs, nil)
	rr = doJSON(r, "POST", "/api/cs/payment-exceptions/"+abID+"/refund", `{"reason":"再次处理"}`, cs, nil)
// 条件判断：if rr.Code < 400 {
	if rr.Code < 400 {
// 执行：t.Fatalf("已处理异常不应重复退款: %d", rr.Code)
		t.Fatalf("已处理异常不应重复退款: %d", rr.Code)
// 结束当前代码块
	}

// 赋值/更新：rr = doJSON(r, "GET", "/api/ops/abnormals?status=handled", "", ops, nil)
	rr = doJSON(r, "GET", "/api/ops/abnormals?status=handled", "", ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "运营查已处理异常")
	statusOK(t, rr, http.StatusOK, "运营查已处理异常")
// 赋值/更新：rr = doJSON(r, "POST", "/api/jobs/close-expired", `{}`, "", nil)
	rr = doJSON(r, "POST", "/api/jobs/close-expired", `{}`, "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "超时清扫任务")
	statusOK(t, rr, http.StatusOK, "超时清扫任务")

	// 原型异常处理入口 HandleAbnormal
	rr = doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-B","quantity":1}`, buyer, nil)
// 赋值/更新：res = parseBody(t, rr)
	res = parseBody(t, rr)
// 赋值/更新：no2, _ := res["order_no"].(string)
	no2, _ := res["order_no"].(string)
// 赋值/更新：amount2 := int64(numOf(res["amount"]))
	amount2 := int64(numOf(res["amount"]))
// 执行：doJSON(r, "POST", "/api/test/expire-orders", `{}`, "", nil)
	doJSON(r, "POST", "/api/test/expire-orders", `{}`, "", nil)
// 赋值/更新：rr = compatPayNotify(t, r, no2, "TXN-LATE-2", amount2, "")
	rr = compatPayNotify(t, r, no2, "TXN-LATE-2", amount2, "")
// 执行：statusOK(t, rr, http.StatusOK, "第二笔迟到通知")
	statusOK(t, rr, http.StatusOK, "第二笔迟到通知")
// 赋值/更新：rr = doJSON(r, "POST", "/api/cs/abnormals/handle", `{"channelTxnId":"TXN-LATE-2","reason":"原型入口处理"}`, cs, nil)
	rr = doJSON(r, "POST", "/api/cs/abnormals/handle", `{"channelTxnId":"TXN-LATE-2","reason":"原型入口处理"}`, cs, nil)
// 执行：statusOK(t, rr, http.StatusOK, "原型异常处理")
	statusOK(t, rr, http.StatusOK, "原型异常处理")

	// 演示模式关闭时禁止重置（负向断言已迁移至 TestDemoModeGatesTestHooks）
}

// 函数定义：func TestDemoResetEnabled(t *testing.T) {
func TestDemoResetEnabled(t *testing.T) {
// 赋值/更新：r, _ := newServer(t, true)
	r, _ := newServer(t, true)
// 赋值/更新：rr := doJSON(r, "POST", "/api/demo/reset", `{}`, "", nil)
	rr := doJSON(r, "POST", "/api/demo/reset", `{}`, "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "演示重置")
	statusOK(t, rr, http.StatusOK, "演示重置")
// 条件判断：if data, ok := parseBody(t, rr)["data"].(map[string]interface{}); !ok || data["reset"] != true {
	if data, ok := parseBody(t, rr)["data"].(map[string]interface{}); !ok || data["reset"] != true {
// 执行：t.Fatalf("重置响应错误: %s", rr.Body.String())
		t.Fatalf("重置响应错误: %s", rr.Body.String())
// 结束当前代码块
	}
// 结束当前代码块
}

// ---------- I-19：HTTP 层数值对账断言 ----------
//
// router 层此前以状态码断言为主，金额/库存数值对账依赖 service 层测试兜底。
// 本测试在 HTTP 响应上直接核对：订单金额快照、可售占用/释放、收退款汇总、库存流水守恒。

// 类型定义：type skuStat struct {
type skuStat struct {
// 执行：price     int64
	price     int64
// 执行：available int64
	available int64
// 执行：occupied  int64
	occupied  int64
// 执行：initial   int64
	initial   int64
// 执行：paid      int64
	paid      int64
// 执行：refunded  int64
	refunded  int64
// 结束当前代码块
}

// 函数定义：func summaryOf(t *testing.T, r *gin.Engine, opsToken, sku string) skuStat {
func summaryOf(t *testing.T, r *gin.Engine, opsToken, sku string) skuStat {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：rr := doJSON(r, "GET", "/api/ops/summary", "", opsToken, nil)
	rr := doJSON(r, "GET", "/api/ops/summary", "", opsToken, nil)
// 执行：statusOK(t, rr, http.StatusOK, "交易汇总")
	statusOK(t, rr, http.StatusOK, "交易汇总")
// 赋值/更新：res := parseBody(t, rr)
	res := parseBody(t, rr)
// 赋值/更新：list, ok := res["products"].([]interface{})
	list, ok := res["products"].([]interface{})
// 条件判断：if !ok {
	if !ok {
// 执行：t.Fatalf("汇总缺少 products: %s", rr.Body.String())
		t.Fatalf("汇总缺少 products: %s", rr.Body.String())
// 结束当前代码块
	}
// 循环处理：for _, it := range list {
	for _, it := range list {
// 赋值/更新：m := it.(map[string]interface{})
		m := it.(map[string]interface{})
// 条件判断：if m["sku"] == sku {
		if m["sku"] == sku {
// 返回：return skuStat{
			return skuStat{
// 执行：price:     int64(numOf(m["price_cents"])),
				price:     int64(numOf(m["price_cents"])),
// 执行：available: int64(numOf(m["available"])),
				available: int64(numOf(m["available"])),
// 执行：occupied:  int64(numOf(m["occupied"])),
				occupied:  int64(numOf(m["occupied"])),
// 执行：initial:   int64(numOf(m["initial_stock"])),
				initial:   int64(numOf(m["initial_stock"])),
// 执行：paid:      int64(numOf(m["paid_total"])),
				paid:      int64(numOf(m["paid_total"])),
// 执行：refunded:  int64(numOf(m["refunded_total"])),
				refunded:  int64(numOf(m["refunded_total"])),
// 结束当前代码块
			}
// 结束当前代码块
		}
// 结束当前代码块
	}
// 执行：t.Fatalf("汇总中未找到 %s: %s", sku, rr.Body.String())
	t.Fatalf("汇总中未找到 %s: %s", sku, rr.Body.String())
// 返回：return skuStat{}
	return skuStat{}
// 结束当前代码块
}

// 函数定义：func productDBID(t *testing.T, r *gin.Engine, sku string) int64 {
func productDBID(t *testing.T, r *gin.Engine, sku string) int64 {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：rr := doJSON(r, "GET", "/api/products", "", "", nil)
	rr := doJSON(r, "GET", "/api/products", "", "", nil)
// 执行：statusOK(t, rr, http.StatusOK, "商品列表")
	statusOK(t, rr, http.StatusOK, "商品列表")
// 循环处理：for _, it := range parseBody(t, rr)["products"].([]interface{}) {
	for _, it := range parseBody(t, rr)["products"].([]interface{}) {
// 赋值/更新：m := it.(map[string]interface{})
		m := it.(map[string]interface{})
// 条件判断：if m["sku"] == sku {
		if m["sku"] == sku {
// 返回：return int64(numOf(m["db_id"]))
			return int64(numOf(m["db_id"]))
// 结束当前代码块
		}
// 结束当前代码块
	}
// 执行：t.Fatalf("商品列表中未找到 %s", sku)
	t.Fatalf("商品列表中未找到 %s", sku)
// 返回：return 0
	return 0
// 结束当前代码块
}

// 函数定义：func TestHTTPNumericReconciliation(t *testing.T) {
func TestHTTPNumericReconciliation(t *testing.T) {
// 赋值/更新：r, _ := newServer(t, false)
	r, _ := newServer(t, false)
// 赋值/更新：buyer := login(t, r, "buyer1")
	buyer := login(t, r, "buyer1")
// 赋值/更新：ops := login(t, r, "ops1")
	ops := login(t, r, "ops1")
// 赋值/更新：cs := login(t, r, "cs1")
	cs := login(t, r, "cs1")

	// 1) HTTP 读取单价与初始可售
	before := summaryOf(t, r, ops, "SKU-A")
// 条件判断：if before.price <= 0 || before.available <= 0 {
	if before.price <= 0 || before.available <= 0 {
// 执行：t.Fatalf("商品基线异常: %+v", before)
		t.Fatalf("商品基线异常: %+v", before)
// 结束当前代码块
	}

	// 2) 下单 SKU-A×2：金额=单价×数量（快照），可售立即反映占用
	rr := doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":2}`, buyer, nil)
// 执行：statusOK(t, rr, http.StatusOK, "下单")
	statusOK(t, rr, http.StatusOK, "下单")
// 赋值/更新：res := parseBody(t, rr)
	res := parseBody(t, rr)
// 赋值/更新：orderNo, _ := res["order_no"].(string)
	orderNo, _ := res["order_no"].(string)
// 条件判断：if orderNo == "" {
	if orderNo == "" {
// 执行：t.Fatalf("缺少订单号: %s", rr.Body.String())
		t.Fatalf("缺少订单号: %s", rr.Body.String())
// 结束当前代码块
	}
// 条件判断：if amount := int64(numOf(res["amount"])); amount != before.price*2 {
	if amount := int64(numOf(res["amount"])); amount != before.price*2 {
// 执行：t.Fatalf("订单金额应=单价×数量 %d, 实际 %d", before.price*2, amount)
		t.Fatalf("订单金额应=单价×数量 %d, 实际 %d", before.price*2, amount)
// 结束当前代码块
	}
// 赋值/更新：afterOrder := summaryOf(t, r, ops, "SKU-A")
	afterOrder := summaryOf(t, r, ops, "SKU-A")
// 条件判断：if afterOrder.available != before.available-2 {
	if afterOrder.available != before.available-2 {
// 执行：t.Fatalf("下单后可售应 %d, 实际 %d", before.available-2, afterOrder.available)
		t.Fatalf("下单后可售应 %d, 实际 %d", before.available-2, afterOrder.available)
// 结束当前代码块
	}
// 条件判断：if afterOrder.occupied < 2 {
	if afterOrder.occupied < 2 {
// 执行：t.Fatalf("下单后占用应≥2, 实际 %d", afterOrder.occupied)
		t.Fatalf("下单后占用应≥2, 实际 %d", afterOrder.occupied)
// 结束当前代码块
	}
// 赋值/更新：amount := before.price * 2
	amount := before.price * 2

	// 3) 支付：已收款精确增加 amount，可售不变
	rr = compatPayNotify(t, r, orderNo, "TXN-RECON-1", amount, "")
// 执行：statusOK(t, rr, http.StatusOK, "支付通知")
	statusOK(t, rr, http.StatusOK, "支付通知")
// 条件判断：if res := parseBody(t, rr); res["status"] != "SUCCESS" {
	if res := parseBody(t, rr); res["status"] != "SUCCESS" {
// 执行：t.Fatalf("支付未受理: %s", rr.Body.String())
		t.Fatalf("支付未受理: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：afterPay := summaryOf(t, r, ops, "SKU-A")
	afterPay := summaryOf(t, r, ops, "SKU-A")
// 条件判断：if afterPay.paid != afterOrder.paid+amount {
	if afterPay.paid != afterOrder.paid+amount {
// 执行：t.Fatalf("支付后已收款应 %d, 实际 %d", afterOrder.paid+amount, afterPay.paid)
		t.Fatalf("支付后已收款应 %d, 实际 %d", afterOrder.paid+amount, afterPay.paid)
// 结束当前代码块
	}
// 条件判断：if afterPay.available != afterOrder.available {
	if afterPay.available != afterOrder.available {
// 执行：t.Fatalf("支付不应改变可售: %d → %d", afterOrder.available, afterPay.available)
		t.Fatalf("支付不应改变可售: %d → %d", afterOrder.available, afterPay.available)
// 结束当前代码块
	}

	// 4) 售前退款全额成功：可售恢复 2，已退款=amount，已收款不变
	rr = doJSON(r, "POST", "/api/orders/"+orderNo+"/refund", `{"reason":"对账测试"}`, buyer, nil)
// 执行：statusOK(t, rr, http.StatusOK, "退款申请")
	statusOK(t, rr, http.StatusOK, "退款申请")
// 赋值/更新：rr = doJSON(r, "POST", "/api/cs/refunds/audit", `{"orderNo":"`+orderNo+`","approve":true,"reason":"ok"}`, cs, ...
	rr = doJSON(r, "POST", "/api/cs/refunds/audit", `{"orderNo":"`+orderNo+`","approve":true,"reason":"ok"}`, cs, nil)
// 执行：statusOK(t, rr, http.StatusOK, "审核同意")
	statusOK(t, rr, http.StatusOK, "审核同意")
// 赋值/更新：afterRefund := summaryOf(t, r, ops, "SKU-A")
	afterRefund := summaryOf(t, r, ops, "SKU-A")
// 条件判断：if afterRefund.available != afterPay.available+2 {
	if afterRefund.available != afterPay.available+2 {
// 执行：t.Fatalf("退款成功后可售应恢复 %d, 实际 %d", afterPay.available+2, afterRefund.available)
		t.Fatalf("退款成功后可售应恢复 %d, 实际 %d", afterPay.available+2, afterRefund.available)
// 结束当前代码块
	}
// 条件判断：if afterRefund.refunded != afterPay.refunded+amount {
	if afterRefund.refunded != afterPay.refunded+amount {
// 执行：t.Fatalf("退款后已退款应 %d, 实际 %d", afterPay.refunded+amount, afterRefund.refunded)
		t.Fatalf("退款后已退款应 %d, 实际 %d", afterPay.refunded+amount, afterRefund.refunded)
// 结束当前代码块
	}
// 条件判断：if afterRefund.paid != afterPay.paid {
	if afterRefund.paid != afterPay.paid {
// 执行：t.Fatalf("退款不改变已收款: %d → %d", afterPay.paid, afterRefund.paid)
		t.Fatalf("退款不改变已收款: %d → %d", afterPay.paid, afterRefund.paid)
// 结束当前代码块
	}

	// 5) 取消待付款订单：金额快照=单价，取消后可售释放、收退款汇总不变
	rr = doJSON(r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":1}`, buyer, nil)
// 执行：statusOK(t, rr, http.StatusOK, "第二单")
	statusOK(t, rr, http.StatusOK, "第二单")
// 赋值/更新：res = parseBody(t, rr)
	res = parseBody(t, rr)
// 赋值/更新：no2, _ := res["order_no"].(string)
	no2, _ := res["order_no"].(string)
// 条件判断：if amt := int64(numOf(res["amount"])); amt != before.price {
	if amt := int64(numOf(res["amount"])); amt != before.price {
// 执行：t.Fatalf("第二单金额应=单价 %d, 实际 %d", before.price, amt)
		t.Fatalf("第二单金额应=单价 %d, 实际 %d", before.price, amt)
// 结束当前代码块
	}
// 赋值/更新：mid := summaryOf(t, r, ops, "SKU-A")
	mid := summaryOf(t, r, ops, "SKU-A")
// 条件判断：if mid.available != afterRefund.available-1 {
	if mid.available != afterRefund.available-1 {
// 执行：t.Fatalf("第二单下单后可售应 %d, 实际 %d", afterRefund.available-1, mid.available)
		t.Fatalf("第二单下单后可售应 %d, 实际 %d", afterRefund.available-1, mid.available)
// 结束当前代码块
	}
// 赋值/更新：rr = doJSON(r, "POST", "/api/orders/"+no2+"/cancel", "", buyer, nil)
	rr = doJSON(r, "POST", "/api/orders/"+no2+"/cancel", "", buyer, nil)
// 执行：statusOK(t, rr, http.StatusOK, "取消")
	statusOK(t, rr, http.StatusOK, "取消")
// 赋值/更新：afterCancel := summaryOf(t, r, ops, "SKU-A")
	afterCancel := summaryOf(t, r, ops, "SKU-A")
// 条件判断：if afterCancel.available != afterRefund.available {
	if afterCancel.available != afterRefund.available {
// 执行：t.Fatalf("取消后可售应回到 %d, 实际 %d", afterRefund.available, afterCancel.available)
		t.Fatalf("取消后可售应回到 %d, 实际 %d", afterRefund.available, afterCancel.available)
// 结束当前代码块
	}
// 条件判断：if afterCancel.paid != afterRefund.paid || afterCancel.refunded != afterRefund.refunded {
	if afterCancel.paid != afterRefund.paid || afterCancel.refunded != afterRefund.refunded {
// 执行：t.Fatalf("取消不应改变收退款汇总: %+v → %+v", afterRefund, afterCancel)
		t.Fatalf("取消不应改变收退款汇总: %+v → %+v", afterRefund, afterCancel)
// 结束当前代码块
	}

	// 6) 库存流水守恒（HTTP 层）：initial_stock + Σdelta == 当前可售，
	//    且最后一条流水的 available_after == 当前可售
	rr = doJSON(r, "GET", "/api/ops/ledger?productId="+itoa(productDBID(t, r, "SKU-A")), "", ops, nil)
// 执行：statusOK(t, rr, http.StatusOK, "库存流水")
	statusOK(t, rr, http.StatusOK, "库存流水")
// 赋值/更新：data, ok := parseBody(t, rr)["data"].([]interface{})
	data, ok := parseBody(t, rr)["data"].([]interface{})
// 条件判断：if !ok || len(data) == 0 {
	if !ok || len(data) == 0 {
// 执行：t.Fatalf("库存流水为空: %s", rr.Body.String())
		t.Fatalf("库存流水为空: %s", rr.Body.String())
// 结束当前代码块
	}
// 声明：var sumDelta int64
	var sumDelta int64
// 声明：var lastAfter int64
	var lastAfter int64
// 循环处理：for _, e := range data {
	for _, e := range data {
// 赋值/更新：m := e.(map[string]interface{})
		m := e.(map[string]interface{})
// 赋值/更新：sumDelta += int64(numOf(m["delta"]))
		sumDelta += int64(numOf(m["delta"]))
// 赋值/更新：lastAfter = int64(numOf(m["availableAfter"]))
		lastAfter = int64(numOf(m["availableAfter"]))
// 结束当前代码块
	}
// 条件判断：if want := afterCancel.initial + sumDelta; want != afterCancel.available {
	if want := afterCancel.initial + sumDelta; want != afterCancel.available {
// 执行：t.Fatalf("库存流水不守恒: 初始 %d + Σdelta %d = %d, 当前可售 %d",
		t.Fatalf("库存流水不守恒: 初始 %d + Σdelta %d = %d, 当前可售 %d",
// 执行：afterCancel.initial, sumDelta, want, afterCancel.available)
			afterCancel.initial, sumDelta, want, afterCancel.available)
// 结束当前代码块
	}
// 条件判断：if lastAfter != afterCancel.available {
	if lastAfter != afterCancel.available {
// 执行：t.Fatalf("流水末条 available_after %d 应等于当前可售 %d", lastAfter, afterCancel.available)
		t.Fatalf("流水末条 available_after %d 应等于当前可售 %d", lastAfter, afterCancel.available)
// 结束当前代码块
	}
// 结束当前代码块
}
