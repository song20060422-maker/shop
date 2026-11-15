// 包声明
package api

// 导入依赖包
import (
// 执行："github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"

// 执行："shop/shop/service"
	"shop/shop/service"
// 结束当前代码块
)

// DemoApi 演示辅助接口（仅 shop.demo-mode=true 时由路由层放行）

// 类型定义：type DemoApi struct {
type DemoApi struct {
// 执行：s *service.ServiceGroup
	s *service.ServiceGroup
// 结束当前代码块
}

// Reset 重置演示业务数据：清空订单/收款/退款/发货/异常/流水，恢复商品初始可售数量。

// 函数定义：func (a DemoApi) Reset(c *gin.Context) {
func (a DemoApi) Reset(c *gin.Context) {
// 条件判断：if err := a.s.DemoService.Reset(); err != nil {
	if err := a.s.DemoService.Reset(); err != nil {
// 执行：fail(c, err)
		fail(c, err)
// 返回：return
		return
// 结束当前代码块
	}
// 执行：ok(c, gin.H{"reset": true})
	ok(c, gin.H{"reset": true})
// 结束当前代码块
}
