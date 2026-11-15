// 包声明
package api_test

// 说明：接口契约（api-contract.md）全链路 HTTP 测试，挂在 api 包（外部测试包）内，
// 使 api 层「原型处理器 + 契约兼容层 + 员工分区」在 api 包覆盖统计中真实计数。
// 与 router 包的集成测试互补：router 关注路由装配与门禁，这里关注处理器行为覆盖率。
// 使用真实 MySQL（shop_test，须 -p 1 串行）。

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
// 执行："sync"
	"sync"
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
// 执行：shoprouter "shop/shop/router"
	shoprouter "shop/shop/router"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// 声明：var (
var (
// 执行：extOnce sync.Once
	extOnce sync.Once
// 执行：extDB   *gorm.DB
	extDB   *gorm.DB
// 执行：extCfg  config.Config
	extCfg  config.Config
// 执行：extErr  error
	extErr  error
// 结束当前代码块
)

// 函数定义：func extSetup(t *testing.T) (*gin.Engine, *service.SimulatedChannel) {
func extSetup(t *testing.T) (*gin.Engine, *service.SimulatedChannel) {
// 执行：t.Helper()
	t.Helper()
// 执行：extOnce.Do(func() {
	extOnce.Do(func() {
// 赋值/更新：dsn := os.Getenv("SHOP_TEST_DSN")
		dsn := os.Getenv("SHOP_TEST_DSN")
// 条件判断：if dsn == "" {
		if dsn == "" {
// 赋值/更新：dsn = "root:root@tcp(127.0.0.1:3307)/shop_test?charset=utf8mb4&parseTime=True&loc=Local"
			dsn = "root:root@tcp(127.0.0.1:3307)/shop_test?charset=utf8mb4&parseTime=True&loc=Local"
// 结束当前代码块
		}
// 赋值/更新：extCfg = config.Default()
		extCfg = config.Default()
// 赋值/更新：extCfg.Mysql.DSN = dsn
		extCfg.Mysql.DSN = dsn
// 赋值/更新：extCfg.Shop.PayNotifySecret = "test-secret"
		extCfg.Shop.PayNotifySecret = "test-secret"
// 赋值/更新：extCfg.Shop.OrderTimeoutMinutes = 10
		extCfg.Shop.OrderTimeoutMinutes = 10
		// 演示模式开启：账号列表、强制过期、演示重置等演示钩子才可用
		extCfg.Shop.DemoMode = true
// 赋值/更新：db, err := shopdb.Gorm(extCfg)
		db, err := shopdb.Gorm(extCfg)
// 条件判断：if err != nil {
		if err != nil {
// 赋值/更新：extErr = err
			extErr = err
// 返回：return
			return
// 结束当前代码块
		}
// 条件判断：if err := shopdb.Migrate(db, migrations.FS); err != nil {
		if err := shopdb.Migrate(db, migrations.FS); err != nil {
// 赋值/更新：extErr = err
			extErr = err
// 返回：return
			return
// 结束当前代码块
		}
// 条件判断：if err := shopdb.SeedAccounts(db, "shop123"); err != nil {
		if err := shopdb.SeedAccounts(db, "shop123"); err != nil {
// 赋值/更新：extErr = err
			extErr = err
// 返回：return
			return
// 结束当前代码块
		}
// 赋值/更新：extDB = db
		extDB = db
// 结束当前代码块
	})
// 条件判断：if extErr != nil {
	if extErr != nil {
// 执行：t.Fatalf("初始化测试库失败（需 MySQL 3307 可用）: %v", extErr)
		t.Fatalf("初始化测试库失败（需 MySQL 3307 可用）: %v", extErr)
// 结束当前代码块
	}

// 赋值/更新：tables := []string{
	tables := []string{
// 执行："sessions",
		"sessions",
// 执行："abnormal_payments", "payments", "payment_notices", "refunds", "shipments",
		"abnormal_payments", "payments", "payment_notices", "refunds", "shipments",
// 执行："inventory_ledger", "operation_logs", "orders", "products", "users",
		"inventory_ledger", "operation_logs", "orders", "products", "users",
// 结束当前代码块
	}
// 循环处理：for _, tbl := range tables {
	for _, tbl := range tables {
// 条件判断：if err := extDB.Exec("DELETE FROM " + tbl).Error; err != nil {
		if err := extDB.Exec("DELETE FROM " + tbl).Error; err != nil {
// 执行：t.Fatalf("清理 %s 失败: %v", tbl, err)
			t.Fatalf("清理 %s 失败: %v", tbl, err)
// 结束当前代码块
		}
// 结束当前代码块
	}
// 条件判断：if err := shopdb.Migrate(extDB, migrations.FS); err != nil {
	if err := shopdb.Migrate(extDB, migrations.FS); err != nil {
// 执行：t.Fatalf("重放迁移失败: %v", err)
		t.Fatalf("重放迁移失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if err := shopdb.SeedAccounts(extDB, "shop123"); err != nil {
	if err := shopdb.SeedAccounts(extDB, "shop123"); err != nil {
// 执行：t.Fatalf("重放种子账号失败: %v", err)
		t.Fatalf("重放种子账号失败: %v", err)
// 结束当前代码块
	}

// 执行：gin.SetMode(gin.TestMode)
	gin.SetMode(gin.TestMode)
// 赋值/更新：cfg := extCfg
	cfg := extCfg
// 赋值/更新：ch := service.NewSimulatedChannel(true)
	ch := service.NewSimulatedChannel(true)
// 赋值/更新：g := service.NewServiceGroup(service.Deps{DB: extDB, Now: time.Now, Cfg: cfg, Channel: ch})
	g := service.NewServiceGroup(service.Deps{DB: extDB, Now: time.Now, Cfg: cfg, Channel: ch})
// 赋值/更新：r := gin.New()
	r := gin.New()
// 执行：shoprouter.InitShopRouter(r.Group(""), g, cfg, time.Now)
	shoprouter.InitShopRouter(r.Group(""), g, cfg, time.Now)
// 返回：return r, ch
	return r, ch
// 结束当前代码块
}

// 函数定义：func extDo(t *testing.T, r *gin.Engine, method, path, body, token string, headers map[string]string) *httptest...
func extDo(t *testing.T, r *gin.Engine, method, path, body, token string, headers map[string]string) *httptest.ResponseRecorder {
// 执行：t.Helper()
	t.Helper()
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

// 函数定义：func extStatus(t *testing.T, rr *httptest.ResponseRecorder, want int, ctx string) {
func extStatus(t *testing.T, rr *httptest.ResponseRecorder, want int, ctx string) {
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

// 函数定义：func extJSON(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
func extJSON(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
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

// 函数定义：func extLogin(t *testing.T, r *gin.Engine, username string) string {
func extLogin(t *testing.T, r *gin.Engine, username string) string {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：rr := extDo(t, r, "POST", "/api/login", `{"username":"`+username+`","password":"shop123"}`, "", nil)
	rr := extDo(t, r, "POST", "/api/login", `{"username":"`+username+`","password":"shop123"}`, "", nil)
// 执行：extStatus(t, rr, http.StatusOK, "登录 "+username)
	extStatus(t, rr, http.StatusOK, "登录 "+username)
// 赋值/更新：token, _ := extJSON(t, rr)["token"].(string)
	token, _ := extJSON(t, rr)["token"].(string)
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

// extSign 契约通道签名：对原始请求体做 HMAC-SHA256（X-Channel-Signature）。
func extSign(secret string, raw []byte) string {
// 赋值/更新：mac := hmac.New(sha256.New, []byte(secret))
	mac := hmac.New(sha256.New, []byte(secret))
// 执行：mac.Write(raw)
	mac.Write(raw)
// 返回：return hex.EncodeToString(mac.Sum(nil))
	return hex.EncodeToString(mac.Sum(nil))
// 结束当前代码块
}

// 函数定义：func jsonNum(v int64) string {
func jsonNum(v int64) string {
// 赋值/更新：b, _ := json.Marshal(v)
	b, _ := json.Marshal(v)
// 返回：return string(b)
	return string(b)
// 结束当前代码块
}

// 走契约通道创建并支付一笔订单，返回订单号与金额。
func extPaidOrder(t *testing.T, r *gin.Engine, buyerToken, txn string) (string, int64) {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：rr := extDo(t, r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":1}`, buyerToken, nil)
	rr := extDo(t, r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":1}`, buyerToken, nil)
// 执行：extStatus(t, rr, http.StatusOK, "下单")
	extStatus(t, rr, http.StatusOK, "下单")
// 赋值/更新：res := extJSON(t, rr)
	res := extJSON(t, rr)
// 赋值/更新：orderNo, _ := res["order_no"].(string)
	orderNo, _ := res["order_no"].(string)
// 赋值/更新：amount := int64(0)
	amount := int64(0)
// 条件判断：if f, ok := res["amount"].(float64); ok {
	if f, ok := res["amount"].(float64); ok {
// 赋值/更新：amount = int64(f)
		amount = int64(f)
// 结束当前代码块
	}
// 条件判断：if orderNo == "" || amount <= 0 {
	if orderNo == "" || amount <= 0 {
// 执行：t.Fatalf("下单响应缺少字段: %s", rr.Body.String())
		t.Fatalf("下单响应缺少字段: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：body := `{"order_no":"` + orderNo + `","channel_txn_no":"` + txn + `","amount":` + jsonNum(amount) + `}`
	body := `{"order_no":"` + orderNo + `","channel_txn_no":"` + txn + `","amount":` + jsonNum(amount) + `}`
// 赋值/更新：rr = extDo(t, r, "POST", "/api/pay/notify", body, "", map[string]string{
	rr = extDo(t, r, "POST", "/api/pay/notify", body, "", map[string]string{
// 执行："X-Channel-Signature": extSign(extCfg.Shop.PayNotifySecret, []byte(body)),
		"X-Channel-Signature": extSign(extCfg.Shop.PayNotifySecret, []byte(body)),
// 结束当前代码块
	})
// 执行：extStatus(t, rr, http.StatusOK, "支付通知")
	extStatus(t, rr, http.StatusOK, "支付通知")
// 返回：return orderNo, amount
	return orderNo, amount
// 结束当前代码块
}

// TestContractStaffFlows 覆盖运营 / 仓库 / 客服三个分区的契约接口与越权门禁。
func TestContractStaffFlows(t *testing.T) {
// 赋值/更新：r, _ := extSetup(t)
	r, _ := extSetup(t)

// 赋值/更新：buyer := extLogin(t, r, "buyer1")
	buyer := extLogin(t, r, "buyer1")
// 赋值/更新：ops := extLogin(t, r, "ops1")
	ops := extLogin(t, r, "ops1")
// 赋值/更新：wh := extLogin(t, r, "wh1")
	wh := extLogin(t, r, "wh1")
// 赋值/更新：cs := extLogin(t, r, "cs1")
	cs := extLogin(t, r, "cs1")

	// 公开区：商品、账号列表、健康检查、买家端页面
	extStatus(t, extDo(t, r, "GET", "/api/products", "", "", nil), http.StatusOK, "商品列表")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/users", "", "", nil), http.StatusOK, "账号列表")
	extStatus(t, extDo(t, r, "GET", "/api/users", "", "", nil), http.StatusOK, "账号列表")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/health", "", "", nil), http.StatusOK, "健康检查")
	extStatus(t, extDo(t, r, "GET", "/api/health", "", "", nil), http.StatusOK, "健康检查")
// 执行：extStatus(t, extDo(t, r, "GET", "/shop", "", "", nil), http.StatusOK, "买家端页面")
	extStatus(t, extDo(t, r, "GET", "/shop", "", "", nil), http.StatusOK, "买家端页面")

	// 未登录访问分区接口 -> 401
	extStatus(t, extDo(t, r, "GET", "/api/ops/summary", "", "", nil), http.StatusUnauthorized, "未登录运营接口")

// 赋值/更新：orderNo, amount := extPaidOrder(t, r, buyer, "TXN-API-1")
	orderNo, amount := extPaidOrder(t, r, buyer, "TXN-API-1")

	// 买家区：我的订单 / 详情 / 退款申请（未发货前可申请）
	extStatus(t, extDo(t, r, "GET", "/api/orders", "", buyer, nil), http.StatusOK, "我的订单")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/orders/"+orderNo, "", buyer, nil), http.StatusOK, "订单详情（买家）")
	extStatus(t, extDo(t, r, "GET", "/api/orders/"+orderNo, "", buyer, nil), http.StatusOK, "订单详情（买家）")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/me", "", buyer, nil), http.StatusOK, "当前登录信息")
	extStatus(t, extDo(t, r, "GET", "/api/me", "", buyer, nil), http.StatusOK, "当前登录信息")
	// 买家越权访问运营接口 -> 403
	extStatus(t, extDo(t, r, "GET", "/api/ops/summary", "", buyer, nil), http.StatusForbidden, "买家访问运营接口")

	// 运营区：汇总 / 订单 / 核查 / 通知 / 异常 / 流水 / 留痕 / 支付记录 / 商品维护
	extStatus(t, extDo(t, r, "GET", "/api/ops/summary", "", ops, nil), http.StatusOK, "运营汇总")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/ops/orders?status=paid", "", ops, nil), http.StatusOK, "运营订单列表")
	extStatus(t, extDo(t, r, "GET", "/api/ops/orders?status=paid", "", ops, nil), http.StatusOK, "运营订单列表")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/ops/orders/"+orderNo, "", ops, nil), http.StatusOK, "运营订单核查")
	extStatus(t, extDo(t, r, "GET", "/api/ops/orders/"+orderNo, "", ops, nil), http.StatusOK, "运营订单核查")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/ops/notices?orderNo="+orderNo, "", ops, nil), http.StatusOK, "通知记录（按单）")
	extStatus(t, extDo(t, r, "GET", "/api/ops/notices?orderNo="+orderNo, "", ops, nil), http.StatusOK, "通知记录（按单）")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/ops/notices?limit=10", "", ops, nil), http.StatusOK, "通知记录（按结果）")
	extStatus(t, extDo(t, r, "GET", "/api/ops/notices?limit=10", "", ops, nil), http.StatusOK, "通知记录（按结果）")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/ops/abnormals", "", ops, nil), http.StatusOK, "异常收款")
	extStatus(t, extDo(t, r, "GET", "/api/ops/abnormals", "", ops, nil), http.StatusOK, "异常收款")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/ops/ledger?productId=1", "", ops, nil), http.StatusOK, "库存流水")
	extStatus(t, extDo(t, r, "GET", "/api/ops/ledger?productId=1", "", ops, nil), http.StatusOK, "库存流水")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/ops/trail?orderNo="+orderNo, "", ops, nil), http.StatusOK, "操作留痕")
	extStatus(t, extDo(t, r, "GET", "/api/ops/trail?orderNo="+orderNo, "", ops, nil), http.StatusOK, "操作留痕")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/ops/payments?orderNo="+orderNo, "", ops, nil), http.StatusOK, "支付记录")
	extStatus(t, extDo(t, r, "GET", "/api/ops/payments?orderNo="+orderNo, "", ops, nil), http.StatusOK, "支付记录")
	// 原型商品维护按数字主键寻址：契约视图用 db_id 暴露真实主键（种子自增 ID 不固定）
	products := extJSON(t, extDo(t, r, "GET", "/api/products", "", "", nil))["data"].([]interface{})
// 赋值/更新：firstProduct := products[0].(map[string]interface{})
	firstProduct := products[0].(map[string]interface{})
// 赋值/更新：pid := uint64(firstProduct["db_id"].(float64))
	pid := uint64(firstProduct["db_id"].(float64))
// 执行：extStatus(t, extDo(t, r, "PATCH", "/api/ops/products", `{"productId":`+jsonNum(int64(pid))+`,"priceCents":1000...
	extStatus(t, extDo(t, r, "PATCH", "/api/ops/products", `{"productId":`+jsonNum(int64(pid))+`,"priceCents":10000}`, ops, nil), http.StatusOK, "商品调价")
// 执行：extStatus(t, extDo(t, r, "PATCH", "/api/ops/products/SKU-A", `{"available":5}`, ops, nil), http.StatusOK, "覆盖可...
	extStatus(t, extDo(t, r, "PATCH", "/api/ops/products/SKU-A", `{"available":5}`, ops, nil), http.StatusOK, "覆盖可售")
// 执行：extStatus(t, extDo(t, r, "PATCH", "/api/ops/products/SKU-A", `{"price_cents":10000,"on_sale":true}`, ops, nil)...
	extStatus(t, extDo(t, r, "PATCH", "/api/ops/products/SKU-A", `{"price_cents":10000,"on_sale":true}`, ops, nil), http.StatusOK, "按 SKU 调价")
// 执行：extStatus(t, extDo(t, r, "PATCH", "/api/ops/products/SKU-A", `{}`, ops, nil), http.StatusBadRequest, "空入参")
	extStatus(t, extDo(t, r, "PATCH", "/api/ops/products/SKU-A", `{}`, ops, nil), http.StatusBadRequest, "空入参")

	// 仓库区：待发货 / 发货 / 重复发货 / 记录 / 详情
	extStatus(t, extDo(t, r, "GET", "/api/warehouse/shipable", "", wh, nil), http.StatusOK, "待发货列表")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/warehouse/orders", "", wh, nil), http.StatusOK, "待发货列表（契约）")
	extStatus(t, extDo(t, r, "GET", "/api/warehouse/orders", "", wh, nil), http.StatusOK, "待发货列表（契约）")
// 执行：extStatus(t, extDo(t, r, "POST", "/api/warehouse/ship", `{"orderNo":"`+orderNo+`","carrier":"顺丰","trackingNo":...
	extStatus(t, extDo(t, r, "POST", "/api/warehouse/ship", `{"orderNo":"`+orderNo+`","carrier":"顺丰","trackingNo":"SF-1"}`, wh, nil), http.StatusOK, "发货登记")
// 执行：extStatus(t, extDo(t, r, "POST", "/api/warehouse/ship", `{"orderNo":"`+orderNo+`","carrier":"顺丰","trackingNo":...
	extStatus(t, extDo(t, r, "POST", "/api/warehouse/ship", `{"orderNo":"`+orderNo+`","carrier":"顺丰","trackingNo":"SF-2"}`, wh, nil), http.StatusConflict, "重复发货")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/warehouse/shipments", "", wh, nil), http.StatusOK, "发货记录")
	extStatus(t, extDo(t, r, "GET", "/api/warehouse/shipments", "", wh, nil), http.StatusOK, "发货记录")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/warehouse/orders/"+orderNo, "", wh, nil), http.StatusOK, "订单详情（仓库）")
	extStatus(t, extDo(t, r, "GET", "/api/warehouse/orders/"+orderNo, "", wh, nil), http.StatusOK, "订单详情（仓库）")
	// 仓库越权访问客服接口 -> 403
	extStatus(t, extDo(t, r, "GET", "/api/cs/refunds", "", wh, nil), http.StatusForbidden, "仓库访问客服接口")

	// 客服区：退款列表 / 审核（售前拒绝）/ 退款重试 / 异常收款 / 通知 / 订单详情
	orderNo2, _ := extPaidOrder(t, r, buyer, "TXN-API-2")
// 赋值/更新：rr := extDo(t, r, "POST", "/api/orders/"+orderNo2+"/refund", `{"reason":"不想要了"}`, buyer, nil)
	rr := extDo(t, r, "POST", "/api/orders/"+orderNo2+"/refund", `{"reason":"不想要了"}`, buyer, nil)
// 执行：extStatus(t, rr, http.StatusOK, "退款申请")
	extStatus(t, rr, http.StatusOK, "退款申请")
// 赋值/更新：refundID := ""
	refundID := ""
// 条件判断：if refund, ok := extJSON(t, rr)["refund"].(map[string]interface{}); ok {
	if refund, ok := extJSON(t, rr)["refund"].(map[string]interface{}); ok {
// 赋值/更新：refundID, _ = refund["refund_no"].(string)
		refundID, _ = refund["refund_no"].(string)
// 结束当前代码块
	}

// 执行：extStatus(t, extDo(t, r, "GET", "/api/cs/refunds", "", cs, nil), http.StatusOK, "退款单列表")
	extStatus(t, extDo(t, r, "GET", "/api/cs/refunds", "", cs, nil), http.StatusOK, "退款单列表")
// 执行：extStatus(t, extDo(t, r, "POST", "/api/cs/refunds/audit", `{"orderNo":"`+orderNo2+`","approve":false,"reason":...
	extStatus(t, extDo(t, r, "POST", "/api/cs/refunds/audit", `{"orderNo":"`+orderNo2+`","approve":false,"reason":"即将发货"}`, cs, nil), http.StatusOK, "退款审核-拒绝")
// 条件判断：if refundID != "" {
	if refundID != "" {
// 执行：extStatus(t, extDo(t, r, "POST", "/api/cs/refunds/"+refundID+"/retry", `{}`, cs, nil), http.StatusConflict, "退...
		extStatus(t, extDo(t, r, "POST", "/api/cs/refunds/"+refundID+"/retry", `{}`, cs, nil), http.StatusConflict, "退款重试-非失败态")
// 结束当前代码块
	}
// 执行：extStatus(t, extDo(t, r, "GET", "/api/cs/payment-exceptions?status=pending", "", cs, nil), http.StatusOK, "异常支...
	extStatus(t, extDo(t, r, "GET", "/api/cs/payment-exceptions?status=pending", "", cs, nil), http.StatusOK, "异常支付列表")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/cs/notices?limit=10", "", cs, nil), http.StatusOK, "通知记录（客服）")
	extStatus(t, extDo(t, r, "GET", "/api/cs/notices?limit=10", "", cs, nil), http.StatusOK, "通知记录（客服）")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/cs/orders/"+orderNo, "", cs, nil), http.StatusOK, "订单详情（客服）")
	extStatus(t, extDo(t, r, "GET", "/api/cs/orders/"+orderNo, "", cs, nil), http.StatusOK, "订单详情（客服）")
// 执行：extStatus(t, extDo(t, r, "POST", "/api/cs/channel/mode", `{"refundMode":"success"}`, cs, nil), http.StatusOK, ...
	extStatus(t, extDo(t, r, "POST", "/api/cs/channel/mode", `{"refundMode":"success"}`, cs, nil), http.StatusOK, "渠道退款模式")

	// 渠道侧：待支付订单、渠道状态、渠道代付、强制过期
	rr = extDo(t, r, "POST", "/api/orders", `{"product_id":"SKU-B","quantity":1}`, buyer, nil)
// 执行：extStatus(t, rr, http.StatusOK, "下单（待付款）")
	extStatus(t, rr, http.StatusOK, "下单（待付款）")
// 赋值/更新：pending := extJSON(t, rr)
	pending := extJSON(t, rr)
// 赋值/更新：pendingNo, _ := pending["order_no"].(string)
	pendingNo, _ := pending["order_no"].(string)
// 执行：extStatus(t, extDo(t, r, "GET", "/api/pay/pending", "", "", nil), http.StatusOK, "渠道待支付订单")
	extStatus(t, extDo(t, r, "GET", "/api/pay/pending", "", "", nil), http.StatusOK, "渠道待支付订单")
// 执行：extStatus(t, extDo(t, r, "GET", "/api/pay/channel", "", "", nil), http.StatusOK, "渠道状态")
	extStatus(t, extDo(t, r, "GET", "/api/pay/channel", "", "", nil), http.StatusOK, "渠道状态")
// 执行：extStatus(t, extDo(t, r, "POST", "/api/pay/channel-mode", `{"refund_mode":"success"}`, "", nil), http.StatusOK...
	extStatus(t, extDo(t, r, "POST", "/api/pay/channel-mode", `{"refund_mode":"success"}`, "", nil), http.StatusOK, "渠道模式（契约）")
// 赋值/更新：rr = extDo(t, r, "POST", "/api/pay/channel-pay", `{"order_no":"`+pendingNo+`","channel_txn_no":"TXN-CH-1","amo...
	rr = extDo(t, r, "POST", "/api/pay/channel-pay", `{"order_no":"`+pendingNo+`","channel_txn_no":"TXN-CH-1","amount":`+jsonNum(amount)+`}`, "", nil)
// 条件判断：if rr.Code != http.StatusOK && rr.Code != http.StatusConflict && rr.Code != http.StatusBadRequest {
	if rr.Code != http.StatusOK && rr.Code != http.StatusConflict && rr.Code != http.StatusBadRequest {
// 执行：t.Fatalf("渠道代付异常状态: %d %s", rr.Code, rr.Body.String())
		t.Fatalf("渠道代付异常状态: %d %s", rr.Code, rr.Body.String())
// 结束当前代码块
	}
// 执行：extStatus(t, extDo(t, r, "POST", "/api/test/expire-orders", `{"order_nos":["`+pendingNo+`"]}`, "", nil), http....
	extStatus(t, extDo(t, r, "POST", "/api/test/expire-orders", `{"order_nos":["`+pendingNo+`"]}`, "", nil), http.StatusOK, "强制过期钩子")
// 执行：extStatus(t, extDo(t, r, "POST", "/api/jobs/close-expired", "", "", nil), http.StatusOK, "超时关单")
	extStatus(t, extDo(t, r, "POST", "/api/jobs/close-expired", "", "", nil), http.StatusOK, "超时关单")
// 执行：extStatus(t, extDo(t, r, "POST", "/api/demo/reset", "", "", nil), http.StatusOK, "演示数据重置")
	extStatus(t, extDo(t, r, "POST", "/api/demo/reset", "", "", nil), http.StatusOK, "演示数据重置")
// 结束当前代码块
}

// TestContractLateAndRefundFlows 覆盖迟到支付转异常、客服处理异常退款、渠道失败与重试。
func TestContractLateAndRefundFlows(t *testing.T) {
// 赋值/更新：r, _ := extSetup(t)
	r, _ := extSetup(t)

// 赋值/更新：buyer := extLogin(t, r, "buyer1")
	buyer := extLogin(t, r, "buyer1")
// 赋值/更新：cs := extLogin(t, r, "cs1")
	cs := extLogin(t, r, "cs1")

	// 1. 退款：渠道首次失败（fail）-> 审核返回 409 + 退款单 -> 切成功 -> 重试成功
	extStatus(t, extDo(t, r, "POST", "/api/pay/channel-mode", `{"refund_mode":"fail"}`, "", nil), http.StatusOK, "渠道置失败")
// 赋值/更新：orderNo, amount := extPaidOrder(t, r, buyer, "TXN-RF-1")
	orderNo, amount := extPaidOrder(t, r, buyer, "TXN-RF-1")
// 赋值/更新：rr := extDo(t, r, "POST", "/api/orders/"+orderNo+"/refund-apply", `{"reason":"退货退款"}`, buyer, nil)
	rr := extDo(t, r, "POST", "/api/orders/"+orderNo+"/refund-apply", `{"reason":"退货退款"}`, buyer, nil)
// 执行：extStatus(t, rr, http.StatusOK, "契约退款申请")
	extStatus(t, rr, http.StatusOK, "契约退款申请")
// 赋值/更新：refundNo := ""
	refundNo := ""
// 条件判断：if refund, ok := extJSON(t, rr)["refund"].(map[string]interface{}); ok {
	if refund, ok := extJSON(t, rr)["refund"].(map[string]interface{}); ok {
// 赋值/更新：refundNo, _ = refund["refund_no"].(string)
		refundNo, _ = refund["refund_no"].(string)
// 结束当前代码块
	}
// 条件判断：if refundNo == "" {
	if refundNo == "" {
// 执行：t.Fatalf("退款申请未返回退款单号: %s", rr.Body.String())
		t.Fatalf("退款申请未返回退款单号: %s", rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：rr = extDo(t, r, "POST", "/api/cs/refunds/"+refundNo+"/decide", `{"approve":true,"reason":"同意"}`, cs, nil)
	rr = extDo(t, r, "POST", "/api/cs/refunds/"+refundNo+"/decide", `{"approve":true,"reason":"同意"}`, cs, nil)
// 执行：extStatus(t, rr, http.StatusOK, "审核-渠道失败仍受理并留痕")
	extStatus(t, rr, http.StatusOK, "审核-渠道失败仍受理并留痕")
// 条件判断：if st, _ := extJSON(t, rr)["status"].(string); st != "FAILED" {
	if st, _ := extJSON(t, rr)["status"].(string); st != "FAILED" {
// 执行：t.Fatalf("渠道失败退款单状态应为 FAILED，实际 %v body=%s", st, rr.Body.String())
		t.Fatalf("渠道失败退款单状态应为 FAILED，实际 %v body=%s", st, rr.Body.String())
// 结束当前代码块
	}
// 执行：extStatus(t, extDo(t, r, "POST", "/api/pay/channel-mode", `{"refund_mode":"success"}`, "", nil), http.StatusOK...
	extStatus(t, extDo(t, r, "POST", "/api/pay/channel-mode", `{"refund_mode":"success"}`, "", nil), http.StatusOK, "渠道置成功")
// 执行：extStatus(t, extDo(t, r, "POST", "/api/cs/refunds/"+refundNo+"/retry", `{}`, cs, nil), http.StatusOK, "退款重试成功"...
	extStatus(t, extDo(t, r, "POST", "/api/cs/refunds/"+refundNo+"/retry", `{}`, cs, nil), http.StatusOK, "退款重试成功")

	// 2. 迟到支付：待付款订单被强制过期 -> 支付通知转异常收款 -> 客服发起退款
	rr = extDo(t, r, "POST", "/api/orders", `{"product_id":"SKU-B","quantity":1}`, buyer, nil)
// 执行：extStatus(t, rr, http.StatusOK, "下单（迟到支付用例）")
	extStatus(t, rr, http.StatusOK, "下单（迟到支付用例）")
// 赋值/更新：late := extJSON(t, rr)
	late := extJSON(t, rr)
// 赋值/更新：lateNo, _ := late["order_no"].(string)
	lateNo, _ := late["order_no"].(string)
// 赋值/更新：lateAmount := int64(0)
	lateAmount := int64(0)
// 条件判断：if f, ok := late["amount"].(float64); ok {
	if f, ok := late["amount"].(float64); ok {
// 赋值/更新：lateAmount = int64(f)
		lateAmount = int64(f)
// 结束当前代码块
	}
// 执行：extStatus(t, extDo(t, r, "POST", "/api/test/expire-orders", `{"order_nos":["`+lateNo+`"]}`, "", nil), http.Sta...
	extStatus(t, extDo(t, r, "POST", "/api/test/expire-orders", `{"order_nos":["`+lateNo+`"]}`, "", nil), http.StatusOK, "强制过期")
// 执行：extStatus(t, extDo(t, r, "POST", "/api/jobs/close-expired", "", "", nil), http.StatusOK, "关单")
	extStatus(t, extDo(t, r, "POST", "/api/jobs/close-expired", "", "", nil), http.StatusOK, "关单")

// 赋值/更新：body := `{"order_no":"` + lateNo + `","channel_txn_no":"TXN-LATE-9","amount":` + jsonNum(lateAmount) + `}`
	body := `{"order_no":"` + lateNo + `","channel_txn_no":"TXN-LATE-9","amount":` + jsonNum(lateAmount) + `}`
// 赋值/更新：rr = extDo(t, r, "POST", "/api/pay/notify", body, "", map[string]string{
	rr = extDo(t, r, "POST", "/api/pay/notify", body, "", map[string]string{
// 执行："X-Channel-Signature": extSign(extCfg.Shop.PayNotifySecret, []byte(body)),
		"X-Channel-Signature": extSign(extCfg.Shop.PayNotifySecret, []byte(body)),
// 结束当前代码块
	})
// 执行：extStatus(t, rr, http.StatusOK, "迟到支付通知")
	extStatus(t, rr, http.StatusOK, "迟到支付通知")
// 赋值/更新：res := extJSON(t, rr)
	res := extJSON(t, rr)
// 条件判断：if status, _ := res["status"].(string); status != "ABNORMAL" && status != "DUPLICATED" {
	if status, _ := res["status"].(string); status != "ABNORMAL" && status != "DUPLICATED" {
// 执行：t.Fatalf("迟到支付应转异常收款，实际 %v body=%s", res["status"], rr.Body.String())
		t.Fatalf("迟到支付应转异常收款，实际 %v body=%s", res["status"], rr.Body.String())
// 结束当前代码块
	}

	// 3. 验签失败 / 金额不符 / 未知订单：均被拒且留痕
	extStatus(t, extDo(t, r, "POST", "/api/pay/notify", body, "", map[string]string{"X-Channel-Signature": "deadbeef"}), http.StatusBadRequest, "验签失败")
// 赋值/更新：unknownBody := `{"order_no":"SH-NOT-EXIST","channel_txn_no":"X-9","amount":100}`
	unknownBody := `{"order_no":"SH-NOT-EXIST","channel_txn_no":"X-9","amount":100}`
// 执行：extStatus(t, extDo(t, r, "POST", "/api/pay/notify", unknownBody, "", map[string]string{
	extStatus(t, extDo(t, r, "POST", "/api/pay/notify", unknownBody, "", map[string]string{
// 执行："X-Channel-Signature": extSign(extCfg.Shop.PayNotifySecret, []byte(unknownBody)),
		"X-Channel-Signature": extSign(extCfg.Shop.PayNotifySecret, []byte(unknownBody)),
// 结束当前代码块
	}), http.StatusBadRequest, "未知订单通知被拒")

	// 4. 内置演示页通道：camelCase 字段 + 请求体 sign（HMAC "订单号|交易号|金额(分)"）
	//    回归：曾只识别 channel_txn_no 与 X-Channel-Signature 头，页面点"发送成功通知"报"参数不合法"。
	rr = extDo(t, r, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":1}`, buyer, nil)
// 执行：extStatus(t, rr, http.StatusOK, "演示页下单")
	extStatus(t, rr, http.StatusOK, "演示页下单")
// 赋值/更新：demoRes := extJSON(t, rr)
	demoRes := extJSON(t, rr)
// 赋值/更新：demoOrder, _ := demoRes["order_no"].(string)
	demoOrder, _ := demoRes["order_no"].(string)
// 赋值/更新：demoAmount := int64(0)
	demoAmount := int64(0)
// 条件判断：if f, ok := demoRes["amount"].(float64); ok {
	if f, ok := demoRes["amount"].(float64); ok {
// 赋值/更新：demoAmount = int64(f)
		demoAmount = int64(f)
// 结束当前代码块
	}
// 赋值/更新：demoTxn := "TXN-DEMO-" + demoOrder
	demoTxn := "TXN-DEMO-" + demoOrder
// 赋值/更新：demoSign := extPayBodySign(extCfg.Shop.PayNotifySecret, demoOrder, demoTxn, demoAmount)
	demoSign := extPayBodySign(extCfg.Shop.PayNotifySecret, demoOrder, demoTxn, demoAmount)
// 赋值/更新：demoBody := `{"orderNo":"` + demoOrder + `","channelTxnId":"` + demoTxn + `","amountCents":` + jsonNum(demoAmo...
	demoBody := `{"orderNo":"` + demoOrder + `","channelTxnId":"` + demoTxn + `","amountCents":` + jsonNum(demoAmount) + `,"sign":"` + demoSign + `"}`
// 执行：extStatus(t, extDo(t, r, "POST", "/api/pay/notify", demoBody, "", nil), http.StatusOK, "演示页通道通知")
	extStatus(t, extDo(t, r, "POST", "/api/pay/notify", demoBody, "", nil), http.StatusOK, "演示页通道通知")
	// 签名错误同样被拒：放宽字段别名不应弱化验签
	badDemo := `{"orderNo":"` + demoOrder + `","channelTxnId":"TXN-DEMO-BAD","amountCents":` + jsonNum(demoAmount) + `,"sign":"deadbeef"}`
// 执行：extStatus(t, extDo(t, r, "POST", "/api/pay/notify", badDemo, "", nil), http.StatusBadRequest, "演示页签名错误被拒")
	extStatus(t, extDo(t, r, "POST", "/api/pay/notify", badDemo, "", nil), http.StatusBadRequest, "演示页签名错误被拒")

// 赋值/更新：_ = amount
	_ = amount
// 结束当前代码块
}

// extPayBodySign 演示页通道签名：HMAC-SHA256("订单号|交易号|金额(分)")。
func extPayBodySign(secret, orderNo, txn string, amount int64) string {
// 返回：return extSign(secret, []byte(orderNo+"|"+txn+"|"+jsonNum(amount)))
	return extSign(secret, []byte(orderNo+"|"+txn+"|"+jsonNum(amount)))
// 结束当前代码块
}
