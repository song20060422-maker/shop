// Package api 承载 HTTP 接口层（类 GVA api 分层），只做参数绑定、权限与响应封装，
// 业务规则全部收敛在 service 层。
package api

// 导入依赖包
import (
// 执行："net/http"
	"net/http"

// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/model/response"
	"shop/shop/model/response"
// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// ApiGroup 全部接口处理器

// 类型定义：type ApiGroup struct {
type ApiGroup struct {
// 执行：AuthApi
	AuthApi
// 执行：BuyerApi
	BuyerApi
// 执行：StaffApi
	StaffApi
// 执行：PayApi
	PayApi
// 执行：DemoApi
	DemoApi
// 结束当前代码块
}

// NewApiGroup 组装接口层

// 函数定义：func NewApiGroup(g *service.ServiceGroup) *ApiGroup {
func NewApiGroup(g *service.ServiceGroup) *ApiGroup {
// 返回：return &ApiGroup{
	return &ApiGroup{
// 执行：AuthApi:  AuthApi{s: g},
		AuthApi:  AuthApi{s: g},
// 执行：BuyerApi: BuyerApi{s: g},
		BuyerApi: BuyerApi{s: g},
// 执行：StaffApi: StaffApi{s: g},
		StaffApi: StaffApi{s: g},
// 执行：PayApi:   PayApi{s: g},
		PayApi:   PayApi{s: g},
// 执行：DemoApi:  DemoApi{s: g},
		DemoApi:  DemoApi{s: g},
// 结束当前代码块
	}
// 结束当前代码块
}

// fail 把业务错误翻译为 HTTP 状态 + 稳定错误码，便于三类使用者看到明确原因。

// 函数定义：func fail(c *gin.Context, err error) {
func fail(c *gin.Context, err error) {
// 赋值/更新：code, msg := service.BizCode(err)
	code, msg := service.BizCode(err)
// 执行：c.JSON(statusOf(code), response.Fail(code, msg))
	c.JSON(statusOf(code), response.Fail(code, msg))
// 结束当前代码块
}

// ok 统一成功响应包装：业务码 0 + data。

// 函数定义：func ok(c *gin.Context, data interface{}) {
func ok(c *gin.Context, data interface{}) {
// 执行：c.JSON(http.StatusOK, response.OK(data))
	c.JSON(http.StatusOK, response.OK(data))
// 结束当前代码块
}

// statusOf 业务错误码 → HTTP 状态码映射（契约中 2002 越权为 403）。

// 函数定义：func statusOf(code int) int {
func statusOf(code int) int {
// 多路分支：switch code {
	switch code {
// 匹配分支：case 1000, 3001, 3002, 4010:
	case 1000, 3001, 3002, 4010:
// 返回：return http.StatusBadRequest
		return http.StatusBadRequest
// 匹配分支：case 1001, 2001, 3005, 4001:
	case 1001, 2001, 3005, 4001:
// 返回：return http.StatusNotFound
		return http.StatusNotFound
// 匹配分支：case 2002:
	case 2002:
// 返回：return http.StatusForbidden
		return http.StatusForbidden
// 匹配分支：case 9000:
	case 9000:
// 返回：return http.StatusInternalServerError
		return http.StatusInternalServerError
// 默认分支
	default:
		// 其余为业务状态冲突：订单已关闭、重复退款、已发货等，语义上是 409
		return http.StatusConflict
// 结束当前代码块
	}
// 结束当前代码块
}
