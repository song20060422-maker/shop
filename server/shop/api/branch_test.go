// 包声明
package api

// 导入依赖包
import (
// 执行："encoding/json"
	"encoding/json"
// 执行："errors"
	"errors"
// 执行："strconv"
	"strconv"
// 执行："strings"
	"strings"
// 执行："testing"
	"testing"
// 执行："time"
	"time"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/model"
	"shop/shop/model"
// 执行："shop/shop/service"
	"shop/shop/service"
// 执行："shop/shop/utils"
	"shop/shop/utils"
// 结束当前代码块
)

// externalChannel 非模拟渠道实现：用于覆盖「渠道不支持切换退款模式」分支。
type externalChannel struct{}

// 函数定义：func (externalChannel) Refund(string, string, int64) (bool, error) { return true, nil }
func (externalChannel) Refund(string, string, int64) (bool, error) { return true, nil }
// 函数定义：func (externalChannel) SetRefundMode(string)                       {}
func (externalChannel) SetRefundMode(string)                       {}
// 函数定义：func (externalChannel) RefundMode() string                         { return "external" }
func (externalChannel) RefundMode() string                         { return "external" }
// 函数定义：func (externalChannel) RefundCalls() int                           { return 0 }
func (externalChannel) RefundCalls() int                           { return 0 }

// 函数定义：func newGroupWithChannel(t *testing.T, ch service.ChannelClient) *service.ServiceGroup {
func newGroupWithChannel(t *testing.T, ch service.ChannelClient) *service.ServiceGroup {
// 执行：t.Helper()
	t.Helper()
// 执行：resetTables(t)
	resetTables(t)
// 返回：return service.NewServiceGroup(service.Deps{DB: testDB, Now: time.Now, Cfg: testCfg, Channel: ch})
	return service.NewServiceGroup(service.Deps{DB: testDB, Now: time.Now, Cfg: testCfg, Channel: ch})
// 结束当前代码块
}

// 函数定义：func TestCompatErrorBranches(t *testing.T) {
func TestCompatErrorBranches(t *testing.T) {
// 赋值/更新：g := newGroup(t)
	g := newGroup(t)
// 赋值/更新：buyer := &utils.Claims{UserID: userIDOf(t, g, "buyer1"), Role: model.RoleBuyer}
	buyer := &utils.Claims{UserID: userIDOf(t, g, "buyer1"), Role: model.RoleBuyer}
// 赋值/更新：ops := &utils.Claims{UserID: userIDOf(t, g, "ops1"), Role: model.RoleOps}
	ops := &utils.Claims{UserID: userIDOf(t, g, "ops1"), Role: model.RoleOps}
// 赋值/更新：cs := &utils.Claims{UserID: userIDOf(t, g, "cs1"), Role: model.RoleCS}
	cs := &utils.Claims{UserID: userIDOf(t, g, "cs1"), Role: model.RoleCS}
// 赋值/更新：a := NewCompatApi(g, testCfg, nil) // now 为 nil 时回落到 time.Now
	a := NewCompatApi(g, testCfg, nil) // now 为 nil 时回落到 time.Now

	// 登录参数缺失 / 密码错误
	w := invoke(t, "POST", "/api/login", `{}`, nil, a.Login)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("空账号应拒绝: %d", w.Code)
		t.Fatalf("空账号应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/login", `{"username":"buyer1","password":"bad"}`, nil, a.Login)
	w = invoke(t, "POST", "/api/login", `{"username":"buyer1","password":"bad"}`, nil, a.Login)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("错误密码应拒绝: %d", w.Code)
		t.Fatalf("错误密码应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 下单参数缺失 / 商品不存在 / 数量非法
	w = invoke(t, "POST", "/api/orders", `{}`, buyer, a.CreateOrder)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("缺商品应拒绝: %d", w.Code)
		t.Fatalf("缺商品应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/orders", `{"product_id":"NO-SKU"}`, buyer, a.CreateOrder)
	w = invoke(t, "POST", "/api/orders", `{"product_id":"NO-SKU"}`, buyer, a.CreateOrder)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知商品应拒绝: %d", w.Code)
		t.Fatalf("未知商品应拒绝: %d", w.Code)
// 结束当前代码块
	}
	// 显式数量 0 属越界（题面：每单 1-2 件），必须拒绝而不是兜底成 1 件
	w = invoke(t, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":0}`, buyer, a.CreateOrder)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("数量 0 应拒绝: %d %s", w.Code, w.Body.String())
		t.Fatalf("数量 0 应拒绝: %d %s", w.Code, w.Body.String())
// 结束当前代码块
	}
	// 完全未传数量字段时按 1 件处理
	w = invoke(t, "POST", "/api/orders", `{"product_id":"SKU-A"}`, buyer, a.CreateOrder)
// 条件判断：if w.Code != 200 {
	if w.Code != 200 {
// 执行：t.Fatalf("未传数量应按 1 件下单: %d %s", w.Code, w.Body.String())
		t.Fatalf("未传数量应按 1 件下单: %d %s", w.Code, w.Body.String())
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":3}`, buyer, a.CreateOrder)
	w = invoke(t, "POST", "/api/orders", `{"product_id":"SKU-A","quantity":3}`, buyer, a.CreateOrder)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("数量 3 应拒绝: %d", w.Code)
		t.Fatalf("数量 3 应拒绝: %d", w.Code)
// 结束当前代码块
	}
	// 数字主键寻址分支
	list, _ := g.CatalogService.OnSale()
// 赋值/更新：pid := strconv.FormatUint(list[0].ID, 10)
	pid := strconv.FormatUint(list[0].ID, 10)
// 赋值/更新：w = invoke(t, "POST", "/api/orders", `{"product_id":"`+pid+`","quantity":1}`, buyer, a.CreateOrder)
	w = invoke(t, "POST", "/api/orders", `{"product_id":"`+pid+`","quantity":1}`, buyer, a.CreateOrder)
// 条件判断：if w.Code != 200 {
	if w.Code != 200 {
// 执行：t.Fatalf("数字主键下单失败: %d %s", w.Code, w.Body.String())
		t.Fatalf("数字主键下单失败: %d %s", w.Code, w.Body.String())
// 结束当前代码块
	}

	// 订单详情未知 / 取消未知
	w = invokeParams(t, "GET", "/api/orders/:orderNo", "", buyer, a.OrderDetail, gin.Params{{Key: "orderNo", Value: "SO-NOPE"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知订单详情应拒绝: %d", w.Code)
		t.Fatalf("未知订单详情应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invokeParams(t, "POST", "/api/orders/:orderNo/cancel", "", buyer, a.CancelOrder, gin.Params{{Key: "orderNo...
	w = invokeParams(t, "POST", "/api/orders/:orderNo/cancel", "", buyer, a.CancelOrder, gin.Params{{Key: "orderNo", Value: "SO-NOPE"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知订单取消应拒绝: %d", w.Code)
		t.Fatalf("未知订单取消应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invokeParams(t, "POST", "/api/orders/:orderNo/refund-apply", "", buyer, a.RefundApply, gin.Params{{Key: "o...
	w = invokeParams(t, "POST", "/api/orders/:orderNo/refund-apply", "", buyer, a.RefundApply, gin.Params{{Key: "orderNo", Value: "SO-NOPE"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知订单退款应拒绝: %d", w.Code)
		t.Fatalf("未知订单退款应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 运营维护：未知商品 / 非法价格 / 空字段
	w = invokeParams(t, "PATCH", "/api/ops/products/:id", `{"price_cents":100}`, ops, a.OpsProductUpdate, gin.Params{{Key: "id", Value: "NO-SKU"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知商品应拒绝: %d", w.Code)
		t.Fatalf("未知商品应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invokeParams(t, "PATCH", "/api/ops/products/:id", `{"price_cents":0}`, ops, a.OpsProductUpdate, gin.Params...
	w = invokeParams(t, "PATCH", "/api/ops/products/:id", `{"price_cents":0}`, ops, a.OpsProductUpdate, gin.Params{{Key: "id", Value: "SKU-A"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("非法价格应拒绝: %d", w.Code)
		t.Fatalf("非法价格应拒绝: %d", w.Code)
// 结束当前代码块
	}
	// 支付记录：空订单号返回空列表
	w = invoke(t, "GET", "/api/ops/payments", "", ops, a.OpsPayments)
// 条件判断：if w.Code != 200 {
	if w.Code != 200 {
// 执行：t.Fatalf("空条件查询应返回空列表: %d", w.Code)
		t.Fatalf("空条件查询应返回空列表: %d", w.Code)
// 结束当前代码块
	}

	// 渠道受理：缺订单号 / 未知订单
	w = invoke(t, "POST", "/api/pay/channel-pay", `{}`, nil, a.PayChannelPay)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("缺订单号应拒绝: %d", w.Code)
		t.Fatalf("缺订单号应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/pay/channel-pay", `{"order_no":"SO-NOPE"}`, nil, a.PayChannelPay)
	w = invoke(t, "POST", "/api/pay/channel-pay", `{"order_no":"SO-NOPE"}`, nil, a.PayChannelPay)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知订单应拒绝: %d", w.Code)
		t.Fatalf("未知订单应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 支付通知：空体 / 金额不符 / 未知订单
	w = invoke(t, "POST", "/api/pay/notify", `{}`, nil, a.PayNotify)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("空通知应拒绝: %d", w.Code)
		t.Fatalf("空通知应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 渠道模式：缺参数 / 非法值
	w = invoke(t, "POST", "/api/pay/channel-mode", `{}`, nil, a.PayChannelMode)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("缺参数应拒绝: %d", w.Code)
		t.Fatalf("缺参数应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/pay/channel-mode", `{"refund_result":"WHAT"}`, nil, a.PayChannelMode)
	w = invoke(t, "POST", "/api/pay/channel-mode", `{"refund_result":"WHAT"}`, nil, a.PayChannelMode)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("非法模式应拒绝: %d", w.Code)
		t.Fatalf("非法模式应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 仓库发货：未知订单（兼容层）
	w = invokeParams(t, "POST", "/api/warehouse/orders/:orderNo/ship", `{}`, cs, a.WarehouseShip, gin.Params{{Key: "orderNo", Value: "SO-NOPE"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知订单发货应拒绝: %d", w.Code)
		t.Fatalf("未知订单发货应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 客服：非法 id / 不存在的异常记录
	w = invokeParams(t, "POST", "/api/cs/payment-exceptions/:id/refund", `{}`, cs, a.CSExceptionRefund, gin.Params{{Key: "id", Value: "abc"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("非法异常 id 应拒绝: %d", w.Code)
		t.Fatalf("非法异常 id 应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invokeParams(t, "POST", "/api/cs/payment-exceptions/:id/refund", `{}`, cs, a.CSExceptionRefund, gin.Params...
	w = invokeParams(t, "POST", "/api/cs/payment-exceptions/:id/refund", `{}`, cs, a.CSExceptionRefund, gin.Params{{Key: "id", Value: "99999"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("不存在异常记录应拒绝: %d", w.Code)
		t.Fatalf("不存在异常记录应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invokeParams(t, "POST", "/api/cs/refunds/:id/decide", `{}`, cs, a.CSRefundDecide, gin.Params{{Key: "id", V...
	w = invokeParams(t, "POST", "/api/cs/refunds/:id/decide", `{}`, cs, a.CSRefundDecide, gin.Params{{Key: "id", Value: "RF-NOPE"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("不存在退款单应拒绝: %d", w.Code)
		t.Fatalf("不存在退款单应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invokeParams(t, "POST", "/api/cs/refunds/:id/retry", `{}`, cs, a.CSRefundRetry, gin.Params{{Key: "id", Val...
	w = invokeParams(t, "POST", "/api/cs/refunds/:id/retry", `{}`, cs, a.CSRefundRetry, gin.Params{{Key: "id", Value: "RF-NOPE"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("不存在退款单重试应拒绝: %d", w.Code)
		t.Fatalf("不存在退款单重试应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 商城页面渲染失败分支（web.Index 正常返回，此处仅验证正常页可读）
	w = invoke(t, "GET", "/api/products", "", nil, a.Products)
// 条件判断：if w.Code != 200 {
	if w.Code != 200 {
// 执行：t.Fatalf("兼容商品列表失败: %d", w.Code)
		t.Fatalf("兼容商品列表失败: %d", w.Code)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestPrototypeErrorBranches(t *testing.T) {
func TestPrototypeErrorBranches(t *testing.T) {
// 赋值/更新：g := newGroup(t)
	g := newGroup(t)
// 赋值/更新：buyer := &utils.Claims{UserID: userIDOf(t, g, "buyer1"), Role: model.RoleBuyer}
	buyer := &utils.Claims{UserID: userIDOf(t, g, "buyer1"), Role: model.RoleBuyer}
// 赋值/更新：cs := &utils.Claims{UserID: userIDOf(t, g, "cs1"), Role: model.RoleCS}
	cs := &utils.Claims{UserID: userIDOf(t, g, "cs1"), Role: model.RoleCS}
// 赋值/更新：wh := &utils.Claims{UserID: userIDOf(t, g, "wh1"), Role: model.RoleWarehouse}
	wh := &utils.Claims{UserID: userIDOf(t, g, "wh1"), Role: model.RoleWarehouse}

	// 当前用户不存在（伪造 claims）
	w := invoke(t, "GET", "/api/me", "", &utils.Claims{UserID: 999999, Role: "buyer"}, AuthApi{g}.Me)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("不存在用户应拒绝: %d", w.Code)
		t.Fatalf("不存在用户应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 买家退款查询：未知订单 / 他人订单 / 无退款单
	w = invokeParams(t, "GET", "/api/orders/:orderNo/refund", "", buyer, BuyerApi{g}.MyRefund, gin.Params{{Key: "orderNo", Value: "SO-NOPE"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知订单退款查询应拒绝: %d", w.Code)
		t.Fatalf("未知订单退款查询应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：o, err := g.OrderService.Create(buyer.UserID, productRowOf(t, g, "SKU-A"), 1)
	o, err := g.OrderService.Create(buyer.UserID, productRowOf(t, g, "SKU-A"), 1)
// 条件判断：if err != nil {
	if err != nil {
// 执行：t.Fatal(err)
		t.Fatal(err)
// 结束当前代码块
	}
// 赋值/更新：other := &utils.Claims{UserID: userIDOf(t, g, "buyer2"), Role: model.RoleBuyer}
	other := &utils.Claims{UserID: userIDOf(t, g, "buyer2"), Role: model.RoleBuyer}
// 赋值/更新：w = invokeParams(t, "GET", "/api/orders/:orderNo/refund", "", other, BuyerApi{g}.MyRefund, gin.Params{{Key: "o...
	w = invokeParams(t, "GET", "/api/orders/:orderNo/refund", "", other, BuyerApi{g}.MyRefund, gin.Params{{Key: "orderNo", Value: o.OrderNo}})
// 条件判断：if w.Code != 403 {
	if w.Code != 403 {
// 执行：t.Fatalf("他人退款查询应 403: %d", w.Code)
		t.Fatalf("他人退款查询应 403: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invokeParams(t, "GET", "/api/orders/:orderNo/refund", "", buyer, BuyerApi{g}.MyRefund, gin.Params{{Key: "o...
	w = invokeParams(t, "GET", "/api/orders/:orderNo/refund", "", buyer, BuyerApi{g}.MyRefund, gin.Params{{Key: "orderNo", Value: o.OrderNo}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("无退款单应拒绝: %d", w.Code)
		t.Fatalf("无退款单应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 原型退款申请：非法入参落到 path 参数后再校验订单
	w = invokeParams(t, "POST", "/api/orders/:orderNo/refund", `{bad json`, buyer, BuyerApi{g}.RefundApply, gin.Params{{Key: "orderNo", Value: "SO-NOPE"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知订单退款申请应拒绝: %d", w.Code)
		t.Fatalf("未知订单退款申请应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 原型发货：非法入参 / 未知订单
	w = invoke(t, "POST", "/api/warehouse/ship", `{bad json`, wh, StaffApi{g}.Ship)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("非法发货入参应拒绝: %d", w.Code)
		t.Fatalf("非法发货入参应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/warehouse/ship", `{"orderNo":"SO-NOPE","carrier":"SF","trackingNo":"X"}`, wh, Staf...
	w = invoke(t, "POST", "/api/warehouse/ship", `{"orderNo":"SO-NOPE","carrier":"SF","trackingNo":"X"}`, wh, StaffApi{g}.Ship)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知订单发货应拒绝: %d", w.Code)
		t.Fatalf("未知订单发货应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 原型运营维护：非法入参 / 未知商品
	w = invoke(t, "PATCH", "/api/ops/products", `{bad json`, cs, StaffApi{g}.ProductUpdate)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("非法入参应拒绝: %d", w.Code)
		t.Fatalf("非法入参应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "PATCH", "/api/ops/products", `{"productId":99999,"priceCents":100}`, cs, StaffApi{g}.ProductUpd...
	w = invoke(t, "PATCH", "/api/ops/products", `{"productId":99999,"priceCents":100}`, cs, StaffApi{g}.ProductUpdate)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知商品应拒绝: %d", w.Code)
		t.Fatalf("未知商品应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 订单全链路：未知订单
	w = invokeParams(t, "GET", "/api/ops/orders/:orderNo", "", cs, StaffApi{g}.OrderTrace, gin.Params{{Key: "orderNo", Value: "SO-NOPE"}})
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知订单全链路应拒绝: %d", w.Code)
		t.Fatalf("未知订单全链路应拒绝: %d", w.Code)
// 结束当前代码块
	}

	// 原型退款审核 / 重试 / 异常处理：非法入参
	w = invoke(t, "POST", "/api/cs/refunds/audit", `{bad json`, cs, StaffApi{g}.RefundAudit)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("非法审核入参应拒绝: %d", w.Code)
		t.Fatalf("非法审核入参应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/cs/refunds/audit", `{"orderNo":"SO-NOPE","approve":true}`, cs, StaffApi{g}.RefundA...
	w = invoke(t, "POST", "/api/cs/refunds/audit", `{"orderNo":"SO-NOPE","approve":true}`, cs, StaffApi{g}.RefundAudit)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知订单审核应拒绝: %d", w.Code)
		t.Fatalf("未知订单审核应拒绝: %d", w.Code)
// 结束当前代码块
	}
	// 客服同意/拒绝必须写明原因（题面要求，空原因不得通过审核）
	w = invoke(t, "POST", "/api/cs/refunds/audit", `{"orderNo":"SO-NOPE","approve":true,"reason":"   "}`, cs, StaffApi{g}.RefundAudit)
// 条件判断：if w.Code != 409 || !strings.Contains(w.Body.String(), "必须填写原因") {
	if w.Code != 409 || !strings.Contains(w.Body.String(), "必须填写原因") {
// 执行：t.Fatalf("空原因审核应被拒: %d %s", w.Code, w.Body.String())
		t.Fatalf("空原因审核应被拒: %d %s", w.Code, w.Body.String())
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/cs/refunds/retry", `{bad json`, cs, StaffApi{g}.RefundRetry)
	w = invoke(t, "POST", "/api/cs/refunds/retry", `{bad json`, cs, StaffApi{g}.RefundRetry)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("非法重试入参应拒绝: %d", w.Code)
		t.Fatalf("非法重试入参应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/cs/refunds/retry", `{"refundNo":"RF-NOPE"}`, cs, StaffApi{g}.RefundRetry)
	w = invoke(t, "POST", "/api/cs/refunds/retry", `{"refundNo":"RF-NOPE"}`, cs, StaffApi{g}.RefundRetry)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知退款单重试应拒绝: %d", w.Code)
		t.Fatalf("未知退款单重试应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/cs/abnormals/handle", `{bad json`, cs, StaffApi{g}.HandleAbnormal)
	w = invoke(t, "POST", "/api/cs/abnormals/handle", `{bad json`, cs, StaffApi{g}.HandleAbnormal)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("非法异常处理入参应拒绝: %d", w.Code)
		t.Fatalf("非法异常处理入参应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 赋值/更新：w = invoke(t, "POST", "/api/cs/abnormals/handle", `{"channelTxnId":"TXN-NOPE"}`, cs, StaffApi{g}.HandleAbnorma...
	w = invoke(t, "POST", "/api/cs/abnormals/handle", `{"channelTxnId":"TXN-NOPE"}`, cs, StaffApi{g}.HandleAbnormal)
// 条件判断：if w.Code < 400 {
	if w.Code < 400 {
// 执行：t.Fatalf("未知异常应拒绝: %d", w.Code)
		t.Fatalf("未知异常应拒绝: %d", w.Code)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestExternalChannelBranches(t *testing.T) {
func TestExternalChannelBranches(t *testing.T) {
// 赋值/更新：g := newGroupWithChannel(t, externalChannel{})
	g := newGroupWithChannel(t, externalChannel{})
// 赋值/更新：a := NewCompatApi(g, testCfg, nil)
	a := NewCompatApi(g, testCfg, nil)
// 赋值/更新：claims := &utils.Claims{UserID: 1, Role: model.RoleOps}
	claims := &utils.Claims{UserID: 1, Role: model.RoleOps}

	// 契约通道仅做契约侧映射，外部渠道调用不报错
	w := invoke(t, "POST", "/api/pay/channel-mode", `{"refund_result":"SUCCESS"}`, nil, a.PayChannelMode)
// 条件判断：if w.Code != 200 {
	if w.Code != 200 {
// 执行：t.Fatalf("契约通道切换应成功: %d %s", w.Code, w.Body.String())
		t.Fatalf("契约通道切换应成功: %d %s", w.Code, w.Body.String())
// 结束当前代码块
	}
	// 原型 ChannelMode 对非模拟渠道拒绝（9001 → 409）
	w = invoke(t, "POST", "/api/cs/channel/mode", `{"refundMode":"fail"}`, claims, PayApi{g}.ChannelMode)
// 条件判断：if w.Code != 409 {
	if w.Code != 409 {
// 执行：t.Fatalf("原型外部渠道切换应失败: %d", w.Code)
		t.Fatalf("原型外部渠道切换应失败: %d", w.Code)
// 结束当前代码块
	}
	// 原型 ChannelState 返回 external
	w = invoke(t, "GET", "/api/pay/channel", "", nil, PayApi{g}.ChannelState)
// 条件判断：if w.Code != 200 {
	if w.Code != 200 {
// 执行：t.Fatalf("渠道状态查询失败: %d", w.Code)
		t.Fatalf("渠道状态查询失败: %d", w.Code)
// 结束当前代码块
	}
	// 兼容层 OpsSummary 角色不符（买家）
	buyer := &utils.Claims{UserID: 2, Role: model.RoleBuyer}
// 赋值/更新：w = invoke(t, "GET", "/api/ops/summary", "", buyer, a.OpsSummary)
	w = invoke(t, "GET", "/api/ops/summary", "", buyer, a.OpsSummary)
// 条件判断：if w.Code != 403 {
	if w.Code != 403 {
// 执行：t.Fatalf("买家查汇总应 403: %d", w.Code)
		t.Fatalf("买家查汇总应 403: %d", w.Code)
// 结束当前代码块
	}
	// 兼容层 CSRefunds 角色不符
	w = invoke(t, "GET", "/api/cs/refunds", "", buyer, a.CSRefunds)
// 条件判断：if w.Code != 403 {
	if w.Code != 403 {
// 执行：t.Fatalf("买家查退款列表应 403: %d", w.Code)
		t.Fatalf("买家查退款列表应 403: %d", w.Code)
// 结束当前代码块
	}
// 结束当前代码块
}

// 函数定义：func TestCompatViewHelpers(t *testing.T) {
func TestCompatViewHelpers(t *testing.T) {
	// 状态映射全分支
	refunds := []model.Refund{
// 执行：{AuditStatus: model.AuditRejected},
		{AuditStatus: model.AuditRejected},
// 执行：{ChannelStatus: model.ChannelSuccess},
		{ChannelStatus: model.ChannelSuccess},
// 执行：{ChannelStatus: model.ChannelFailed},
		{ChannelStatus: model.ChannelFailed},
// 执行：{AuditStatus: model.AuditApproved},
		{AuditStatus: model.AuditApproved},
// 执行：{},
		{},
// 结束当前代码块
	}
// 赋值/更新：wants := []string{"REJECTED", "REFUNDED", "FAILED", "REFUNDING", "PENDING"}
	wants := []string{"REJECTED", "REFUNDED", "FAILED", "REFUNDING", "PENDING"}
// 循环处理：for i, r := range refunds {
	for i, r := range refunds {
// 条件判断：if got := refundStatusLabel(r); got != wants[i] {
		if got := refundStatusLabel(r); got != wants[i] {
// 执行：t.Fatalf("退款状态映射错误 %d: %s", i, got)
			t.Fatalf("退款状态映射错误 %d: %s", i, got)
// 结束当前代码块
		}
// 结束当前代码块
	}
	// 失败原因：有审核原因 / 无审核原因 / 非失败
	if got := refundFailReason(model.Refund{ChannelStatus: model.ChannelFailed, AuditReason: "x"}); got != "渠道退款未受理：x" {
// 执行：t.Fatalf("失败原因错误: %s", got)
		t.Fatalf("失败原因错误: %s", got)
// 结束当前代码块
	}
// 条件判断：if got := refundFailReason(model.Refund{ChannelStatus: model.ChannelFailed}); got == "" {
	if got := refundFailReason(model.Refund{ChannelStatus: model.ChannelFailed}); got == "" {
// 执行：t.Fatal("失败原因兜底文案为空")
		t.Fatal("失败原因兜底文案为空")
// 结束当前代码块
	}
// 条件判断：if got := refundFailReason(model.Refund{ChannelStatus: model.ChannelSuccess}); got != "" {
	if got := refundFailReason(model.Refund{ChannelStatus: model.ChannelSuccess}); got != "" {
// 执行：t.Fatalf("成功退款不应有失败原因: %s", got)
		t.Fatalf("成功退款不应有失败原因: %s", got)
// 结束当前代码块
	}
	// 时间输出：零值 / 指针
	if got := ts(time.Time{}); got != "" {
// 执行：t.Fatalf("零值时间应为空: %s", got)
		t.Fatalf("零值时间应为空: %s", got)
// 结束当前代码块
	}
// 条件判断：if got := tsP(nil); got != "" {
	if got := tsP(nil); got != "" {
// 执行：t.Fatalf("nil 时间指针应为空: %s", got)
		t.Fatalf("nil 时间指针应为空: %s", got)
// 结束当前代码块
	}
	// 退款视图：无审核原因时回落到来源
	view := refundView(model.Refund{RefundNo: "RF1", Source: model.RefundSourcePresale})
// 条件判断：if view["reason"] != model.RefundSourcePresale {
	if view["reason"] != model.RefundSourcePresale {
// 执行：t.Fatalf("退款视图原因回落错误: %v", view["reason"])
		t.Fatalf("退款视图原因回落错误: %v", view["reason"])
// 结束当前代码块
	}
	// 订单状态映射默认分支
	if orderStatusLabel("odd") != "CLOSED" {
// 执行：t.Fatal("未知订单状态应映射 CLOSED")
		t.Fatal("未知订单状态应映射 CLOSED")
// 结束当前代码块
	}
	// 通知状态映射全分支
	cases := map[string]string{
// 执行：model.NoticeAccepted: "ACCEPTED", model.NoticeRejected: "REJECTED",
		model.NoticeAccepted: "ACCEPTED", model.NoticeRejected: "REJECTED",
// 执行：model.NoticeDuplicated: "DUPLICATED", model.NoticeAbnormal: "ABNORMAL",
		model.NoticeDuplicated: "DUPLICATED", model.NoticeAbnormal: "ABNORMAL",
// 执行："weird": "WEIRD",
		"weird": "WEIRD",
// 结束当前代码块
	}
// 循环处理：for in, want := range cases {
	for in, want := range cases {
// 条件判断：if got := noticeStatusLabel(model.PaymentNotice{Result: in}); got != want {
		if got := noticeStatusLabel(model.PaymentNotice{Result: in}); got != want {
// 执行：t.Fatalf("通知状态映射错误 %s: %s", in, got)
			t.Fatalf("通知状态映射错误 %s: %s", in, got)
// 结束当前代码块
		}
// 结束当前代码块
	}
	// 异常状态映射
	if abnormalStatusLabel(model.AbnormalHandled) != "HANDLED" || abnormalStatusLabel(model.AbnormalPending) != "PENDING" {
// 执行：t.Fatal("异常状态映射错误")
		t.Fatal("异常状态映射错误")
// 结束当前代码块
	}

	// 解析辅助：字符串 / 数字 / 布尔 / json.Number
	m := map[string]interface{}{
// 执行："s": " 42 ", "f": float64(7), "frac": float64(1.5),
		"s": " 42 ", "f": float64(7), "frac": float64(1.5),
// 执行："n": json.Number("9"), "b": true,
		"n": json.Number("9"), "b": true,
// 结束当前代码块
	}
// 条件判断：if strOf(m, "s") != "42" || strOf(m, "f") != "7" || strOf(m, "frac") != "1.5" || strOf(m, "n") != "9" {
	if strOf(m, "s") != "42" || strOf(m, "f") != "7" || strOf(m, "frac") != "1.5" || strOf(m, "n") != "9" {
// 执行：t.Fatalf("strOf 分支错误: %s %s %s %s", strOf(m, "s"), strOf(m, "f"), strOf(m, "frac"), strOf(m, "n"))
		t.Fatalf("strOf 分支错误: %s %s %s %s", strOf(m, "s"), strOf(m, "f"), strOf(m, "frac"), strOf(m, "n"))
// 结束当前代码块
	}
// 条件判断：if strOf(m, "none") != "" {
	if strOf(m, "none") != "" {
// 执行：t.Fatal("缺省字段应为空串")
		t.Fatal("缺省字段应为空串")
// 结束当前代码块
	}
// 条件判断：if n, ok := intOf(m, "f"); !ok || n != 7 {
	if n, ok := intOf(m, "f"); !ok || n != 7 {
// 执行：t.Fatal("intOf float 分支错误")
		t.Fatal("intOf float 分支错误")
// 结束当前代码块
	}
// 条件判断：if n, ok := intOf(map[string]interface{}{"n": json.Number("11")}, "n"); !ok || n != 11 {
	if n, ok := intOf(map[string]interface{}{"n": json.Number("11")}, "n"); !ok || n != 11 {
// 执行：t.Fatal("intOf json.Number 分支错误")
		t.Fatal("intOf json.Number 分支错误")
// 结束当前代码块
	}
// 条件判断：if n, ok := intOf(map[string]interface{}{"s": "13"}, "s"); !ok || n != 13 {
	if n, ok := intOf(map[string]interface{}{"s": "13"}, "s"); !ok || n != 13 {
// 执行：t.Fatal("intOf string 分支错误")
		t.Fatal("intOf string 分支错误")
// 结束当前代码块
	}
// 条件判断：if _, ok := intOf(map[string]interface{}{"s": "bad"}, "s"); ok {
	if _, ok := intOf(map[string]interface{}{"s": "bad"}, "s"); ok {
// 执行：t.Fatal("非法数字应解析失败")
		t.Fatal("非法数字应解析失败")
// 结束当前代码块
	}
// 条件判断：if b, ok := boolOf(m, "b"); !ok || !b {
	if b, ok := boolOf(m, "b"); !ok || !b {
// 执行：t.Fatal("boolOf 错误")
		t.Fatal("boolOf 错误")
// 结束当前代码块
	}
// 条件判断：if _, ok := boolOf(m, "s"); ok {
	if _, ok := boolOf(m, "s"); ok {
// 执行：t.Fatal("非布尔字段应解析失败")
		t.Fatal("非布尔字段应解析失败")
// 结束当前代码块
	}
	// 失败响应：未知错误（9000）不暴露内部实现细节
	a := NewCompatApi(nil, testCfg, nil)
// 赋值/更新：w := invoke(t, "GET", "/api/x", "", nil, func(c *gin.Context) { a.failJSON(c, errors.New("internal detail")) }...
	w := invoke(t, "GET", "/api/x", "", nil, func(c *gin.Context) { a.failJSON(c, errors.New("internal detail")) })
// 条件判断：if w.Code != 500 {
	if w.Code != 500 {
// 执行：t.Fatalf("未知错误应 500: %d", w.Code)
		t.Fatalf("未知错误应 500: %d", w.Code)
// 结束当前代码块
	}
// 结束当前代码块
}
