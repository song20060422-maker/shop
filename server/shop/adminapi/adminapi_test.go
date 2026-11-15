// 包声明
package adminapi_test

// 说明：GVA 管理端商城后台接口（运营/仓库/客服）测试。
// 鉴权由 gva 中间件负责（PrivateGroup + JWT + Casbin），这里用中间件注入 claims
// 模拟已登录的管理员，覆盖 adminapi 三个分组的成功/失败分支与 service 装配。
// 使用真实 MySQL（shop_test，须 -p 1 串行）。

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
// 执行："testing"
	"testing"
// 执行："time"
	"time"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"
// 执行："gorm.io/gorm"
	"gorm.io/gorm"

// 执行：systemReq "shop/model/system/request"
	systemReq "shop/model/system/request"
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
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// 声明：var (
var (
// 执行：testDB  *gorm.DB
	testDB  *gorm.DB
// 执行：testCfg config.Config
	testCfg config.Config
// 结束当前代码块
)

// adminID 模拟已登录的 GVA 管理员用户
const adminID = 888

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

// newAdminServer 组装带「已登录管理员」的引擎：claims 注入等价于 gva JWT 中间件产物。
// 返回可调时钟指针，便于构造「订单已超时关闭后收到支付通知」的迟到支付场景。
func newAdminServer(t *testing.T) (*gin.Engine, *service.ServiceGroup, *time.Time) {
// 执行：t.Helper()
	t.Helper()
// 执行：resetTables(t)
	resetTables(t)
// 执行：gin.SetMode(gin.TestMode)
	gin.SetMode(gin.TestMode)
// 赋值/更新：r := gin.New()
	r := gin.New()
// 赋值/更新：ch := service.NewSimulatedChannel(true)
	ch := service.NewSimulatedChannel(true)
// 赋值/更新：clock := time.Now()
	clock := time.Now()
// 赋值/更新：g := service.NewServiceGroup(service.Deps{DB: testDB, Now: func() time.Time { return clock }, Cfg: testCfg, Ch...
	g := service.NewServiceGroup(service.Deps{DB: testDB, Now: func() time.Time { return clock }, Cfg: testCfg, Channel: ch})
// 赋值/更新：group := r.Group("", func(c *gin.Context) {
	group := r.Group("", func(c *gin.Context) {
// 执行：c.Set("claims", &systemReq.CustomClaims{BaseClaims: systemReq.BaseClaims{ID: adminID}})
		c.Set("claims", &systemReq.CustomClaims{BaseClaims: systemReq.BaseClaims{ID: adminID}})
// 结束当前代码块
	})
// 执行：shoprouter.InitShopAdminRouter(group, g)
	shoprouter.InitShopAdminRouter(group, g)
// 返回：return r, g, &clock
	return r, g, &clock
// 结束当前代码块
}

// 函数定义：func doJSON(t *testing.T, r *gin.Engine, method, path, body string) map[string]interface{} {
func doJSON(t *testing.T, r *gin.Engine, method, path, body string) map[string]interface{} {
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
// 赋值/更新：rr := httptest.NewRecorder()
	rr := httptest.NewRecorder()
// 执行：r.ServeHTTP(rr, req)
	r.ServeHTTP(rr, req)
// 条件判断：if rr.Code != http.StatusOK {
	if rr.Code != http.StatusOK {
// 执行：t.Fatalf("%s %s: HTTP %d body=%s", method, path, rr.Code, rr.Body.String())
		t.Fatalf("%s %s: HTTP %d body=%s", method, path, rr.Code, rr.Body.String())
// 结束当前代码块
	}
// 赋值/更新：out := map[string]interface{}{}
	out := map[string]interface{}{}
// 条件判断：if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
// 执行：t.Fatalf("%s %s: 响应不是 JSON: %v", method, path, err)
		t.Fatalf("%s %s: 响应不是 JSON: %v", method, path, err)
// 结束当前代码块
	}
// 返回：return out
	return out
// 结束当前代码块
}

// wantOK 断言 GVA 成功响应（code=0）
func wantOK(t *testing.T, res map[string]interface{}, ctx string) map[string]interface{} {
// 执行：t.Helper()
	t.Helper()
// 条件判断：if code, _ := res["code"].(float64); code != 0 {
	if code, _ := res["code"].(float64); code != 0 {
// 执行：t.Fatalf("%s: 期望 code=0 实际 %v msg=%v", ctx, res["code"], res["msg"])
		t.Fatalf("%s: 期望 code=0 实际 %v msg=%v", ctx, res["code"], res["msg"])
// 结束当前代码块
	}
// 赋值/更新：data, _ := res["data"].(map[string]interface{})
	data, _ := res["data"].(map[string]interface{})
// 返回：return data
	return data
// 结束当前代码块
}

// wantFail 断言业务失败响应（code=7，GVA FailWithMessage/FailWithDetailed）
func wantFail(t *testing.T, res map[string]interface{}, ctx string) {
// 执行：t.Helper()
	t.Helper()
// 条件判断：if code, _ := res["code"].(float64); code == 0 {
	if code, _ := res["code"].(float64); code == 0 {
// 执行：t.Fatalf("%s: 期望业务失败，实际成功 data=%v", ctx, res["data"])
		t.Fatalf("%s: 期望业务失败，实际成功 data=%v", ctx, res["data"])
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func firstProductID(t *testing.T, g *service.ServiceGroup) uint64 {
func firstProductID(t *testing.T, g *service.ServiceGroup) uint64 {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：list, err := g.CatalogService.OnSale()
	list, err := g.CatalogService.OnSale()
// 条件判断：if err != nil || len(list) == 0 {
	if err != nil || len(list) == 0 {
// 执行：t.Fatalf("取商品失败: %v", err)
		t.Fatalf("取商品失败: %v", err)
// 结束当前代码块
	}
// 返回：return list[0].ID
	return list[0].ID
// 结束当前代码块
}

// buyerID 取种子买家账号 ID（清表重放种子后自增 ID 会变化，不能写死）
func buyerID(t *testing.T, g *service.ServiceGroup) uint64 {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：list, err := g.AuthService.Users("buyer")
	list, err := g.AuthService.Users("buyer")
// 条件判断：if err != nil || len(list) == 0 {
	if err != nil || len(list) == 0 {
// 执行：t.Fatalf("取买家账号失败: %v", err)
		t.Fatalf("取买家账号失败: %v", err)
// 结束当前代码块
	}
// 返回：return list[0].ID
	return list[0].ID
// 结束当前代码块
}

// paidOrderNo 创建一笔已支付订单，返回订单号、金额、商品 ID
func paidOrderNo(t *testing.T, g *service.ServiceGroup, txn string) (string, int64, uint64) {
// 执行：t.Helper()
	t.Helper()
// 赋值/更新：pid := firstProductID(t, g)
	pid := firstProductID(t, g)
// 赋值/更新：order, err := g.OrderService.Create(buyerID(t, g), pid, 1)
	order, err := g.OrderService.Create(buyerID(t, g), pid, 1)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("下单失败: %v", err)
		t.Fatalf("下单失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, order.OrderNo, txn, order.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, order.OrderNo, txn, order.AmountCents)
// 条件判断：if _, err := g.PaymentService.Notify(order.OrderNo, txn, order.AmountCents, sign, ""); err != nil {
	if _, err := g.PaymentService.Notify(order.OrderNo, txn, order.AmountCents, sign, ""); err != nil {
// 执行：t.Fatalf("支付通知失败: %v", err)
		t.Fatalf("支付通知失败: %v", err)
// 结束当前代码块
	}
// 返回：return order.OrderNo, order.AmountCents, pid
	return order.OrderNo, order.AmountCents, pid
// 结束当前代码块
}

// ---------- 运营 ----------

// 函数定义：func TestOpsHandlers(t *testing.T) {
func TestOpsHandlers(t *testing.T) {
// 赋值/更新：r, g, _ := newAdminServer(t)
	r, g, _ := newAdminServer(t)

// 赋值/更新：orderNo, amount, pid := paidOrderNo(t, g, "TXN-OPS-1")
	orderNo, amount, pid := paidOrderNo(t, g, "TXN-OPS-1")

// 条件判断：if data := wantOK(t, doJSON(t, r, "GET", "/shop/ops/summary", ""), "汇总"); data == nil {
	if data := wantOK(t, doJSON(t, r, "GET", "/shop/ops/summary", ""), "汇总"); data == nil {
// 执行：t.Fatal("汇总应返回数据")
		t.Fatal("汇总应返回数据")
// 结束当前代码块
	}
// 条件判断：if data := wantOK(t, doJSON(t, r, "GET", "/shop/ops/products", ""), "商品列表"); data == nil {
	if data := wantOK(t, doJSON(t, r, "GET", "/shop/ops/products", ""), "商品列表"); data == nil {
// 执行：t.Logf("商品列表返回非对象（数组）：%v", data)
		t.Logf("商品列表返回非对象（数组）：%v", data)
// 结束当前代码块
	}
// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/ops/orders?status=paid", ""), "订单列表(按状态)")
	wantOK(t, doJSON(t, r, "GET", "/shop/ops/orders?status=paid", ""), "订单列表(按状态)")
// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/ops/orders", ""), "订单列表(全量)")
	wantOK(t, doJSON(t, r, "GET", "/shop/ops/orders", ""), "订单列表(全量)")

	// 商品维护：调价 + 覆盖可售 + 上下架
	wantOK(t, doJSON(t, r, "PATCH", "/shop/ops/products/"+itoa(pid), `{"priceCents":12300,"available":9}`), "商品调价/库存")
// 执行：wantOK(t, doJSON(t, r, "PATCH", "/shop/ops/products/"+itoa(pid), `{"onSale":false}`), "商品下架")
	wantOK(t, doJSON(t, r, "PATCH", "/shop/ops/products/"+itoa(pid), `{"onSale":false}`), "商品下架")
// 执行：wantOK(t, doJSON(t, r, "PATCH", "/shop/ops/products/"+itoa(pid), `{"onSale":true,"priceCents":10000}`), "商品上架"...
	wantOK(t, doJSON(t, r, "PATCH", "/shop/ops/products/"+itoa(pid), `{"onSale":true,"priceCents":10000}`), "商品上架")

	// 失败分支：非法主键 / 非法 JSON / 不存在商品 / 空入参
	wantFail(t, doJSON(t, r, "PATCH", "/shop/ops/products/abc", `{}`), "非法主键")
// 执行：wantFail(t, doJSON(t, r, "PATCH", "/shop/ops/products/"+itoa(pid), `{`), "非法 JSON")
	wantFail(t, doJSON(t, r, "PATCH", "/shop/ops/products/"+itoa(pid), `{`), "非法 JSON")
// 执行：wantFail(t, doJSON(t, r, "PATCH", "/shop/ops/products/999999", `{"onSale":true}`), "商品不存在")
	wantFail(t, doJSON(t, r, "PATCH", "/shop/ops/products/999999", `{"onSale":true}`), "商品不存在")
// 执行：wantFail(t, doJSON(t, r, "PATCH", "/shop/ops/products/"+itoa(pid), `{}`), "空入参")
	wantFail(t, doJSON(t, r, "PATCH", "/shop/ops/products/"+itoa(pid), `{}`), "空入参")

	// 订单核查 / 留痕 / 流水 / 通知 / 异常
	wantOK(t, doJSON(t, r, "GET", "/shop/ops/orders/"+orderNo, ""), "订单核查")
// 执行：wantFail(t, doJSON(t, r, "GET", "/shop/ops/orders/SH-NOT-EXIST", ""), "订单核查-不存在")
	wantFail(t, doJSON(t, r, "GET", "/shop/ops/orders/SH-NOT-EXIST", ""), "订单核查-不存在")
// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/ops/trail?orderNo="+orderNo, ""), "操作留痕")
	wantOK(t, doJSON(t, r, "GET", "/shop/ops/trail?orderNo="+orderNo, ""), "操作留痕")
// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/ops/ledger?productId="+itoa(pid), ""), "库存流水")
	wantOK(t, doJSON(t, r, "GET", "/shop/ops/ledger?productId="+itoa(pid), ""), "库存流水")
// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/ops/notices?orderNo="+orderNo, ""), "通知按单号")
	wantOK(t, doJSON(t, r, "GET", "/shop/ops/notices?orderNo="+orderNo, ""), "通知按单号")
// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/ops/notices?result=SUCCESS&limit=5", ""), "通知按结果")
	wantOK(t, doJSON(t, r, "GET", "/shop/ops/notices?result=SUCCESS&limit=5", ""), "通知按结果")
// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/ops/abnormals?status=pending", ""), "异常收款列表")
	wantOK(t, doJSON(t, r, "GET", "/shop/ops/abnormals?status=pending", ""), "异常收款列表")

// 赋值/更新：_ = amount
	_ = amount
// 结束当前代码块
}

// ---------- 仓库 ----------

// 函数定义：func TestWarehouseHandlers(t *testing.T) {
func TestWarehouseHandlers(t *testing.T) {
// 赋值/更新：r, g, _ := newAdminServer(t)
	r, g, _ := newAdminServer(t)

// 赋值/更新：orderNo, _, _ := paidOrderNo(t, g, "TXN-WH-1")
	orderNo, _, _ := paidOrderNo(t, g, "TXN-WH-1")

// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/warehouse/shipable", ""), "待发货列表")
	wantOK(t, doJSON(t, r, "GET", "/shop/warehouse/shipable", ""), "待发货列表")

// 执行：wantFail(t, doJSON(t, r, "POST", "/shop/warehouse/ship", `{`), "发货-非法 JSON")
	wantFail(t, doJSON(t, r, "POST", "/shop/warehouse/ship", `{`), "发货-非法 JSON")
// 执行：wantFail(t, doJSON(t, r, "POST", "/shop/warehouse/ship", `{"orderNo":"SH-NOT-EXIST","carrier":"SF","trackingNo...
	wantFail(t, doJSON(t, r, "POST", "/shop/warehouse/ship", `{"orderNo":"SH-NOT-EXIST","carrier":"SF","trackingNo":"1"}`), "发货-订单不存在")

// 执行：wantOK(t, doJSON(t, r, "POST", "/shop/warehouse/ship",
	wantOK(t, doJSON(t, r, "POST", "/shop/warehouse/ship",
// 执行：`{"orderNo":"`+orderNo+`","carrier":"顺丰","trackingNo":"SF123456789"}`), "发货登记")
		`{"orderNo":"`+orderNo+`","carrier":"顺丰","trackingNo":"SF123456789"}`), "发货登记")
// 执行：wantFail(t, doJSON(t, r, "POST", "/shop/warehouse/ship",
	wantFail(t, doJSON(t, r, "POST", "/shop/warehouse/ship",
// 执行：`{"orderNo":"`+orderNo+`","carrier":"顺丰","trackingNo":"SF123456789"}`), "重复发货")
		`{"orderNo":"`+orderNo+`","carrier":"顺丰","trackingNo":"SF123456789"}`), "重复发货")

// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/warehouse/shipments", ""), "发货记录")
	wantOK(t, doJSON(t, r, "GET", "/shop/warehouse/shipments", ""), "发货记录")
// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/warehouse/orders/"+orderNo, ""), "订单详情（仓库）")
	wantOK(t, doJSON(t, r, "GET", "/shop/warehouse/orders/"+orderNo, ""), "订单详情（仓库）")
// 执行：wantFail(t, doJSON(t, r, "GET", "/shop/warehouse/orders/SH-NOT-EXIST", ""), "订单详情-不存在")
	wantFail(t, doJSON(t, r, "GET", "/shop/warehouse/orders/SH-NOT-EXIST", ""), "订单详情-不存在")
// 结束当前代码块
}

// ---------- 客服 ----------

// 函数定义：func TestCsHandlers(t *testing.T) {
func TestCsHandlers(t *testing.T) {
// 赋值/更新：r, g, clk := newAdminServer(t)
	r, g, clk := newAdminServer(t)

	// 渠道开关：失败模式 -> 审核退款失败可重试 -> 成功模式重试成功
	wantFail(t, doJSON(t, r, "POST", "/shop/cs/channel/mode", `{`), "渠道模式-非法 JSON")
// 执行：wantOK(t, doJSON(t, r, "POST", "/shop/cs/channel/mode", `{"refundMode":"fail"}`), "渠道模式-失败")
	wantOK(t, doJSON(t, r, "POST", "/shop/cs/channel/mode", `{"refundMode":"fail"}`), "渠道模式-失败")

// 赋值/更新：orderNo, _, _ := paidOrderNo(t, g, "TXN-CS-1")
	orderNo, _, _ := paidOrderNo(t, g, "TXN-CS-1")
// 赋值/更新：refund, err := g.RefundService.Apply(buyerID(t, g), orderNo, "买家申请退款")
	refund, err := g.RefundService.Apply(buyerID(t, g), orderNo, "买家申请退款")
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("退款申请失败: %v", err)
		t.Fatalf("退款申请失败: %v", err)
// 结束当前代码块
	}

// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/cs/refunds?status=pending", ""), "退款单列表")
	wantOK(t, doJSON(t, r, "GET", "/shop/cs/refunds?status=pending", ""), "退款单列表")
// 执行：wantFail(t, doJSON(t, r, "POST", "/shop/cs/refunds/audit", `{`), "审核-非法 JSON")
	wantFail(t, doJSON(t, r, "POST", "/shop/cs/refunds/audit", `{`), "审核-非法 JSON")
	// 渠道失败：审核已生效但渠道调用失败，返回 code!=0 且带退款单
	wantFail(t, doJSON(t, r, "POST", "/shop/cs/refunds/audit",
// 执行：`{"orderNo":"`+orderNo+`","approve":true,"reason":"同意整单退款"}`), "审核-渠道失败")
		`{"orderNo":"`+orderNo+`","approve":true,"reason":"同意整单退款"}`), "审核-渠道失败")

// 执行：wantOK(t, doJSON(t, r, "POST", "/shop/cs/channel/mode", `{"refundMode":"success"}`), "渠道模式-成功")
	wantOK(t, doJSON(t, r, "POST", "/shop/cs/channel/mode", `{"refundMode":"success"}`), "渠道模式-成功")
// 执行：wantOK(t, doJSON(t, r, "POST", "/shop/cs/refunds/retry", `{"refundNo":"`+refund.RefundNo+`"}`), "退款重试")
	wantOK(t, doJSON(t, r, "POST", "/shop/cs/refunds/retry", `{"refundNo":"`+refund.RefundNo+`"}`), "退款重试")
// 执行：wantFail(t, doJSON(t, r, "POST", "/shop/cs/refunds/retry", `{"refundNo":"RF-NOT-EXIST"}`), "重试-退款单不存在")
	wantFail(t, doJSON(t, r, "POST", "/shop/cs/refunds/retry", `{"refundNo":"RF-NOT-EXIST"}`), "重试-退款单不存在")

	// 售前拒绝分支：新订单 -> 申请 -> 拒绝（订单回到待发货）
	orderNo2, _, _ := paidOrderNo(t, g, "TXN-CS-2")
// 条件判断：if _, err := g.RefundService.Apply(buyerID(t, g), orderNo2, "临时申请"); err != nil {
	if _, err := g.RefundService.Apply(buyerID(t, g), orderNo2, "临时申请"); err != nil {
// 执行：t.Fatalf("退款申请失败: %v", err)
		t.Fatalf("退款申请失败: %v", err)
// 结束当前代码块
	}
// 执行：wantOK(t, doJSON(t, r, "POST", "/shop/cs/refunds/audit",
	wantOK(t, doJSON(t, r, "POST", "/shop/cs/refunds/audit",
// 执行：`{"orderNo":"`+orderNo2+`","approve":false,"reason":"发货在即"}`), "审核-拒绝")
		`{"orderNo":"`+orderNo2+`","approve":false,"reason":"发货在即"}`), "审核-拒绝")

	// 迟到支付 -> 异常收款 -> 客服发起退款（推进时钟 + 服务端关单构造真实场景）
	orderNo3, err := g.OrderService.Create(buyerID(t, g), firstProductID(t, g), 1)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("下单失败: %v", err)
		t.Fatalf("下单失败: %v", err)
// 结束当前代码块
	}
	*clk = clk.Add(11 * time.Minute)
// 条件判断：if _, err := g.OrderService.CloseExpired(); err != nil {
	if _, err := g.OrderService.CloseExpired(); err != nil {
// 执行：t.Fatalf("超时关单失败: %v", err)
		t.Fatalf("超时关单失败: %v", err)
// 结束当前代码块
	}
// 赋值/更新：txn := "TXN-LATE-1"
	txn := "TXN-LATE-1"
// 赋值/更新：sign := utils.PaySign(testCfg.Shop.PayNotifySecret, orderNo3.OrderNo, txn, orderNo3.AmountCents)
	sign := utils.PaySign(testCfg.Shop.PayNotifySecret, orderNo3.OrderNo, txn, orderNo3.AmountCents)
// 赋值/更新：res, err := g.PaymentService.Notify(orderNo3.OrderNo, txn, orderNo3.AmountCents, sign, "")
	res, err := g.PaymentService.Notify(orderNo3.OrderNo, txn, orderNo3.AmountCents, sign, "")
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatalf("迟到支付通知失败: %v", err)
		t.Fatalf("迟到支付通知失败: %v", err)
// 结束当前代码块
	}
// 条件判断：if res.Result != "abnormal" { // model.NoticeAbnormal
	if res.Result != "abnormal" { // model.NoticeAbnormal
// 执行：t.Fatalf("迟到支付应转异常收款，实际 %s（%s）", res.Result, res.Reason)
		t.Fatalf("迟到支付应转异常收款，实际 %s（%s）", res.Result, res.Reason)
// 结束当前代码块
	}

// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/cs/abnormals?status=pending", ""), "异常收款列表（客服）")
	wantOK(t, doJSON(t, r, "GET", "/shop/cs/abnormals?status=pending", ""), "异常收款列表（客服）")
// 执行：wantFail(t, doJSON(t, r, "POST", "/shop/cs/abnormals/handle", `{`), "异常处理-非法 JSON")
	wantFail(t, doJSON(t, r, "POST", "/shop/cs/abnormals/handle", `{`), "异常处理-非法 JSON")
// 执行：wantFail(t, doJSON(t, r, "POST", "/shop/cs/abnormals/handle",
	wantFail(t, doJSON(t, r, "POST", "/shop/cs/abnormals/handle",
// 执行：`{"channelTxnId":"TXN-NOT-EXIST","reason":"x"}`), "异常处理-交易号不存在")
		`{"channelTxnId":"TXN-NOT-EXIST","reason":"x"}`), "异常处理-交易号不存在")
// 执行：wantOK(t, doJSON(t, r, "POST", "/shop/cs/abnormals/handle",
	wantOK(t, doJSON(t, r, "POST", "/shop/cs/abnormals/handle",
// 执行：`{"channelTxnId":"`+txn+`","reason":"核实后原路退回"}`), "异常处理-发起退款")
		`{"channelTxnId":"`+txn+`","reason":"核实后原路退回"}`), "异常处理-发起退款")

// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/cs/notices?limit=10", ""), "通知记录（客服）")
	wantOK(t, doJSON(t, r, "GET", "/shop/cs/notices?limit=10", ""), "通知记录（客服）")
// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/cs/notices?result=ABNORMAL&limit=10", ""), "通知记录-按结果")
	wantOK(t, doJSON(t, r, "GET", "/shop/cs/notices?result=ABNORMAL&limit=10", ""), "通知记录-按结果")
// 执行：wantOK(t, doJSON(t, r, "GET", "/shop/cs/orders/"+orderNo, ""), "订单详情（客服）")
	wantOK(t, doJSON(t, r, "GET", "/shop/cs/orders/"+orderNo, ""), "订单详情（客服）")
// 执行：wantFail(t, doJSON(t, r, "GET", "/shop/cs/orders/SH-NOT-EXIST", ""), "订单详情-不存在")
	wantFail(t, doJSON(t, r, "GET", "/shop/cs/orders/SH-NOT-EXIST", ""), "订单详情-不存在")
// 结束当前代码块
}

// 函数定义：func itoa(v uint64) string {
func itoa(v uint64) string {
// 赋值/更新：b, _ := json.Marshal(v)
	b, _ := json.Marshal(v)
// 返回：return string(b)
	return string(b)
// 结束当前代码块
}
