// 包声明
package api

// 导入依赖包
import (
// 执行："bytes"
	"bytes"
// 执行："encoding/json"
	"encoding/json"
// 执行："net/http"
	"net/http"
// 执行："net/http/httptest"
	"net/http/httptest"
// 执行："os"
	"os"
// 执行："strconv"
	"strconv"
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
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// 说明：本文件直测被契约兼容层委派遮蔽的原型 handler（同一路径已由 CompatApi 接管），
// 与 router 集成测试互补，保证原型实现同样被覆盖。
// 数据库可用 SHOP_TEST_DSN 覆盖（默认 shop_test，须 -p 1 串行执行）。

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

// 函数定义：func newGroup(t *testing.T) *service.ServiceGroup {
func newGroup(t *testing.T) *service.ServiceGroup {
// 执行：t.Helper()
	t.Helper()
// 执行：resetTables(t)
	resetTables(t)
// 赋值/更新：ch := service.NewSimulatedChannel(true)
	ch := service.NewSimulatedChannel(true)
// 返回：return service.NewServiceGroup(service.Deps{DB: testDB, Now: time.Now, Cfg: testCfg, Channel: ch})
	return service.NewServiceGroup(service.Deps{DB: testDB, Now: time.Now, Cfg: testCfg, Channel: ch})
// 结束当前代码块
}

// invoke 以带 body 与 claims 的测试上下文直调 handler，返回响应记录器。
func invoke(t *testing.T, method, path, body string, claims *utils.Claims, handler gin.HandlerFunc) *httptest.ResponseRecorder {
// 执行：t.Helper()
	t.Helper()
// 返回：return invokeParams(t, method, path, body, claims, handler, nil)
	return invokeParams(t, method, path, body, claims, handler, nil)
// 结束当前代码块
}

// 函数定义：func invokeParams(t *testing.T, method, path, body string, claims *utils.Claims, handler gin.HandlerFunc, para...
func invokeParams(t *testing.T, method, path, body string, claims *utils.Claims, handler gin.HandlerFunc, params gin.Params) *httptest.ResponseRecorder {
// 执行：t.Helper()
	t.Helper()
// 执行：gin.SetMode(gin.TestMode)
	gin.SetMode(gin.TestMode)
// 赋值/更新：w := httptest.NewRecorder()
	w := httptest.NewRecorder()
// 赋值/更新：c, _ := gin.CreateTestContext(w)
	c, _ := gin.CreateTestContext(w)
// 条件判断：if body != "" {
	if body != "" {
// 赋值/更新：c.Request = httptest.NewRequest(method, path, bytes.NewBufferString(body))
		c.Request = httptest.NewRequest(method, path, bytes.NewBufferString(body))
// 执行：c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Content-Type", "application/json")
// 条件分支：} else {
	} else {
// 赋值/更新：c.Request = httptest.NewRequest(method, path, nil)
		c.Request = httptest.NewRequest(method, path, nil)
// 结束当前代码块
	}
// 条件判断：if claims != nil {
	if claims != nil {
// 执行：c.Set("shopClaims", *claims)
		c.Set("shopClaims", *claims)
// 结束当前代码块
	}
// 赋值/更新：c.Params = params
	c.Params = params
// 执行：handler(c)
	handler(c)
// 返回：return w
	return w
// 结束当前代码块
}

// 函数定义：func userIDOf(t *testing.T, g *service.ServiceGroup, username string) uint64 {
func userIDOf(t *testing.T, g *service.ServiceGroup, username string) uint64 {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：login, err := g.AuthService.Login(username, "shop123")
	login, err := g.AuthService.Login(username, "shop123")
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("登录 %s 失败: %v", username, err)
		t.Fatalf("登录 %s 失败: %v", username, err)
// 结束当前代码块
	}
// 返回：return login.User.ID
	return login.User.ID
// 结束当前代码块
}

// 函数定义：func productRowOf(t *testing.T, g *service.ServiceGroup, sku string) uint64 {
func productRowOf(t *testing.T, g *service.ServiceGroup, sku string) uint64 {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：list, err := g.CatalogService.OnSale()
	list, err := g.CatalogService.OnSale()
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 循环处理：for _, p := range list {
	for _, p := range list {
// 条件判断：if p.SKU == sku {
		if p.SKU == sku {
// 返回：return p.ID
			return p.ID
// 结束当前代码块
		}
// 结束当前代码块
	}
// 执行：t.Fatalf("未找到商品 %s", sku)
	t.Fatalf("未找到商品 %s", sku)
// 返回：return 0
	return 0
// 结束当前代码块
}

// 函数定义：func codeOf(t *testing.T, w *httptest.ResponseRecorder) float64 {
func codeOf(t *testing.T, w *httptest.ResponseRecorder) float64 {
// 执行：t.Helper()
	t.Helper()
// 声明：var out map[string]interface{}
	var out map[string]interface{}
// 条件判断：if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
// 执行：t.Fatalf("响应不是 JSON: %s", w.Body.String())
		t.Fatalf("响应不是 JSON: %s", w.Body.String())
// 结束当前代码块
	}
// 条件判断：if f, ok := out["code"].(float64); ok {
	if f, ok := out["code"].(float64); ok {
// 返回：return f
		return f
// 结束当前代码块
	}
// 执行：t.Fatalf("响应缺少 code: %s", w.Body.String())
	t.Fatalf("响应缺少 code: %s", w.Body.String())
// 返回：return -1
	return -1
// 结束当前代码块
}

// 函数定义：func TestShadowedPrototypeHandlers(t *testing.T) {
func TestShadowedPrototypeHandlers(t *testing.T) {
// 赋值/更新：g := newGroup(t)
	g := newGroup(t)
// 赋值/更新：buyer := userIDOf(t, g, "buyer1")
	buyer := userIDOf(t, g, "buyer1")
// 赋值/更新：pid := productRowOf(t, g, "SKU-A")
	pid := productRowOf(t, g, "SKU-A")

	// 原型登录（被 CompatApi.Login 遮蔽的原型实现）
	w := invoke(t, "POST", "/api/login", `{"username":"buyer1","password":"shop123"}`, nil, AuthApi{g}.Login)
// 条件判断：if w.Code != http.StatusOK || codeOf(t, w) != 0 {
	if w.Code != http.StatusOK || codeOf(t, w) != 0 {
// 执行：t.Fatalf("原型登录失败: %d %s", w.Code, w.Body.String())
		t.Fatalf("原型登录失败: %d %s", w.Code, w.Body.String())
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/login", `{"username":"buyer1","password":"wrong"}`, nil, AuthApi{g}.Login)
	w = invoke(t, "POST", "/api/login", `{"username":"buyer1","password":"wrong"}`, nil, AuthApi{g}.Login)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("错误密码应拒绝: %d", w.Code)
		t.Fatalf("错误密码应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/login", `{bad json`, nil, AuthApi{g}.Login)
	w = invoke(t, "POST", "/api/login", `{bad json`, nil, AuthApi{g}.Login)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("非法入参应拒绝: %d", w.Code)
		t.Fatalf("非法入参应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 原型商品列表
	w = invoke(t, "GET", "/api/products", "", nil, BuyerApi{g}.Products)
// 条件判断：if w.Code != http.StatusOK {
	if w.Code != http.StatusOK {
// 执行：t.Fatalf("原型商品列表失败: %d", w.Code)
		t.Fatalf("原型商品列表失败: %d", w.Code)
// 结束当前代码块
	}

	// 原型下单
	claims := &utils.Claims{UserID: buyer, Role: "buyer"}
// 赋值/更新：w = invoke(t, "POST", "/api/orders", `{"productId":`+strconv.FormatUint(pid, 10)+`,"qty":2}`, claims, BuyerApi...
	w = invoke(t, "POST", "/api/orders", `{"productId":`+strconv.FormatUint(pid, 10)+`,"qty":2}`, claims, BuyerApi{g}.CreateOrder)
// 条件判断：if w.Code != http.StatusOK || codeOf(t, w) != 0 {
	if w.Code != http.StatusOK || codeOf(t, w) != 0 {
// 执行：t.Fatalf("原型下单失败: %d %s", w.Code, w.Body.String())
		t.Fatalf("原型下单失败: %d %s", w.Code, w.Body.String())
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/orders", `{"productId":`+strconv.FormatUint(pid, 10)+`,"qty":5}`, claims, BuyerApi...
	w = invoke(t, "POST", "/api/orders", `{"productId":`+strconv.FormatUint(pid, 10)+`,"qty":5}`, claims, BuyerApi{g}.CreateOrder)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("超量下单应拒绝: %d", w.Code)
		t.Fatalf("超量下单应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/orders", `{bad json`, claims, BuyerApi{g}.CreateOrder)
	w = invoke(t, "POST", "/api/orders", `{bad json`, claims, BuyerApi{g}.CreateOrder)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("非法入参应拒绝: %d", w.Code)
		t.Fatalf("非法入参应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 原型订单详情（本人可见，他人 403）
	list, _ := g.OrderService.ListByBuyer(buyer, "")
// 条件判断：if len(list) == 0 {
	if len(list) == 0 {
// 执行：t.Fatal("下单失败")
		t.Fatal("下单失败")
// 结束当前代码块
	}
// 赋值/更新：orderNo := list[0].OrderNo
	orderNo := list[0].OrderNo
// 赋值/更新：w = invokeParams(t, "GET", "/api/orders/:orderNo", "", claims, BuyerApi{g}.OrderDetail, gin.Params{{Key: "orde...
	w = invokeParams(t, "GET", "/api/orders/:orderNo", "", claims, BuyerApi{g}.OrderDetail, gin.Params{{Key: "orderNo", Value: orderNo}})
// 条件判断：if w.Code != http.StatusOK {
	if w.Code != http.StatusOK {
// 执行：t.Fatalf("原型订单详情失败: %d", w.Code)
		t.Fatalf("原型订单详情失败: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：other := &utils.Claims{UserID: userIDOf(t, g, "buyer2"), Role: "buyer"}
	other := &utils.Claims{UserID: userIDOf(t, g, "buyer2"), Role: "buyer"}
// 赋值/更新：w = invokeParams(t, "GET", "/api/orders/:orderNo", "", other, BuyerApi{g}.OrderDetail, gin.Params{{Key: "order...
	w = invokeParams(t, "GET", "/api/orders/:orderNo", "", other, BuyerApi{g}.OrderDetail, gin.Params{{Key: "orderNo", Value: orderNo}})
// 条件判断：if w.Code != http.StatusForbidden {
	if w.Code != http.StatusForbidden {
// 执行：t.Fatalf("他人订单详情应 403: %d", w.Code)
		t.Fatalf("他人订单详情应 403: %d", w.Code)
// 结束当前代码块
	}

	// 原型取消
	w = invokeParams(t, "POST", "/api/orders/:orderNo/cancel", "", claims, BuyerApi{g}.Cancel, gin.Params{{Key: "orderNo", Value: orderNo}})
// 条件判断：if w.Code != http.StatusOK {
	if w.Code != http.StatusOK {
// 执行：t.Fatalf("原型取消失败: %d", w.Code)
		t.Fatalf("原型取消失败: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invokeParams(t, "POST", "/api/orders/:orderNo/cancel", "", claims, BuyerApi{g}.Cancel, gin.Params{{Key: "o...
	w = invokeParams(t, "POST", "/api/orders/:orderNo/cancel", "", claims, BuyerApi{g}.Cancel, gin.Params{{Key: "orderNo", Value: orderNo}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("重复取消应拒绝: %d", w.Code)
		t.Fatalf("重复取消应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invokeParams(t, "POST", "/api/orders/:orderNo/cancel", "", other, BuyerApi{g}.Cancel, gin.Params{{Key: "or...
	w = invokeParams(t, "POST", "/api/orders/:orderNo/cancel", "", other, BuyerApi{g}.Cancel, gin.Params{{Key: "orderNo", Value: orderNo}})
// 条件判断：if w.Code != http.StatusForbidden {
	if w.Code != http.StatusForbidden {
// 执行：t.Fatalf("他人取消应 403: %d", w.Code)
		t.Fatalf("他人取消应 403: %d", w.Code)
// 结束当前代码块
	}

	// 原型支付通知（成功 / 重复 409 / 验签失败拒绝）
	order2, err := g.OrderService.Create(buyer, pid, 1)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, order2.OrderNo, "TXN-API-1", order2.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, order2.OrderNo, "TXN-API-1", order2.AmountCents)
// 赋值/更新：body := `{"orderNo":"` + order2.OrderNo + `","channelTxnId":"TXN-API-1","amountCents":` +
	body := `{"orderNo":"` + order2.OrderNo + `","channelTxnId":"TXN-API-1","amountCents":` +
// 执行：strconv.FormatInt(order2.AmountCents, 10) + `,"sign":"` + sign + `"}`
		strconv.FormatInt(order2.AmountCents, 10) + `,"sign":"` + sign + `"}`
// 赋值/更新：w = invoke(t, "POST", "/api/pay/notify", body, nil, PayApi{g}.Notify)
	w = invoke(t, "POST", "/api/pay/notify", body, nil, PayApi{g}.Notify)
// 条件判断：if w.Code != http.StatusOK {
	if w.Code != http.StatusOK {
// 执行：t.Fatalf("原型通知失败: %d %s", w.Code, w.Body.String())
		t.Fatalf("原型通知失败: %d %s", w.Code, w.Body.String())
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/pay/notify", body, nil, PayApi{g}.Notify)
	w = invoke(t, "POST", "/api/pay/notify", body, nil, PayApi{g}.Notify)
// 条件判断：if w.Code != http.StatusConflict {
	if w.Code != http.StatusConflict {
// 执行：t.Fatalf("重复通知应 409: %d", w.Code)
		t.Fatalf("重复通知应 409: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/pay/notify", `{"orderNo":"x","channelTxnId":"y","amountCents":1,"sign":"bad"}`, ni...
	w = invoke(t, "POST", "/api/pay/notify", `{"orderNo":"x","channelTxnId":"y","amountCents":1,"sign":"bad"}`, nil, PayApi{g}.Notify)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("验签失败应拒绝: %d", w.Code)
		t.Fatalf("验签失败应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 原型运营汇总与退款列表
	opsClaims := &utils.Claims{UserID: userIDOf(t, g, "ops1"), Role: "ops"}
// 赋值/更新：w = invoke(t, "GET", "/api/ops/summary", "", opsClaims, StaffApi{g}.Summary)
	w = invoke(t, "GET", "/api/ops/summary", "", opsClaims, StaffApi{g}.Summary)
// 条件判断：if w.Code != http.StatusOK {
	if w.Code != http.StatusOK {
// 执行：t.Fatalf("原型汇总失败: %d", w.Code)
		t.Fatalf("原型汇总失败: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "GET", "/api/cs/refunds", "", opsClaims, StaffApi{g}.Refunds)
	w = invoke(t, "GET", "/api/cs/refunds", "", opsClaims, StaffApi{g}.Refunds)
// 条件判断：if w.Code != http.StatusOK {
	if w.Code != http.StatusOK {
// 执行：t.Fatalf("原型退款列表失败: %d", w.Code)
		t.Fatalf("原型退款列表失败: %d", w.Code)
// 结束当前代码块
	}

	// 原型渠道模式切换（SimulatedChannel 支持；非法入参拒绝）
	w = invoke(t, "POST", "/api/cs/channel/mode", `{"refundMode":"success"}`, opsClaims, PayApi{g}.ChannelMode)
// 条件判断：if w.Code != http.StatusOK {
	if w.Code != http.StatusOK {
// 执行：t.Fatalf("原型渠道模式切换失败: %d", w.Code)
		t.Fatalf("原型渠道模式切换失败: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/cs/channel/mode", `{bad json`, opsClaims, PayApi{g}.ChannelMode)
	w = invoke(t, "POST", "/api/cs/channel/mode", `{bad json`, opsClaims, PayApi{g}.ChannelMode)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("非法入参应拒绝: %d", w.Code)
		t.Fatalf("非法入参应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 结束当前代码块
}
