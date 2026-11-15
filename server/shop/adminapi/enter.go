// Package adminapi 提供 GVA 管理端（后台）使用的商城运营/仓库/客服接口。
// 与买家端独立 Gin 服务（cmd/shop）共用同一套 service 层与数据库（gva 库）；
// 鉴权走 GVA JWT + Casbin（挂在 PrivateGroup），操作人以当前登录的 GVA 用户落 operation_logs。
package adminapi

// 导入依赖包
import (
// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"
// 执行："shop/model/common/response"
	"shop/model/common/response"
// 执行："shop/utils"
	"shop/utils"

// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// Api 商城后台管理接口（运营/仓库/客服）

// 类型定义：type Api struct {
type Api struct {
// 执行：s *service.ServiceGroup
	s *service.ServiceGroup
// 结束当前代码块
}

// NewApi 组装接口层

// 函数定义：func NewApi(g *service.ServiceGroup) *Api {
func NewApi(g *service.ServiceGroup) *Api {
// 返回：return &Api{s: g}
	return &Api{s: g}
// 结束当前代码块
}

// actor 返回操作人：GVA 登录用户 ID + 该接口对应的功能角色（ops/warehouse/cs）。
// service 层按角色做权限收敛，后台接口已在路由层经 GVA JWT + Casbin 鉴权。

// 函数定义：func actor(c *gin.Context, role string) (uint64, string) {
func actor(c *gin.Context, role string) (uint64, string) {
// 返回：return uint64(utils.GetUserID(c)), role
	return uint64(utils.GetUserID(c)), role
// 结束当前代码块
}

// ok 统一成功响应（GVA 结构 {code:0,data,msg}）

// 函数定义：func ok(c *gin.Context, data interface{}) {
func ok(c *gin.Context, data interface{}) {
// 执行：response.OkWithData(data, c)
	response.OkWithData(data, c)
// 结束当前代码块
}

// fail 业务错误统一转 GVA 失败响应（code:7），提示文案取业务错误信息

// 函数定义：func fail(c *gin.Context, err error) {
func fail(c *gin.Context, err error) {
// 赋值/更新：_, msg := service.BizCode(err)
	_, msg := service.BizCode(err)
// 执行：response.FailWithMessage(msg, c)
	response.FailWithMessage(msg, c)
// 结束当前代码块
}
