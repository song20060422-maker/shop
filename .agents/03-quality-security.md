# 03 · 质量与安全（覆盖率 / race / 红线）

> 每次自测、声称覆盖率之前必读。

## 覆盖率门禁

- 自编代码整体语句覆盖 ≥85%，核心包 `service/api/middleware/config/initialize/router` 各 ≥85%。
- profile 必须单命令生成（禁止拼接）：`go test -p 1 ./shop/... "-coverprofile=coverage-single.out" "-coverpkg=./shop/..." -covermode=atomic -count=1`。
- 共库测试一律 `-p 1` 串行；覆盖不足只能补真实业务路径，工具函数凑数视为低价值（审计 I-23）。
- 基线（2026-11-13 拆分后）：并集 88.7%，service 85.5%，api 81.3%，adminapi 81.4%。

## 测试要求

- 单元 + HTTP 集成并存，覆盖正常/边界/非法状态/重复/权限/并发（20 并发、退款×发货、重复通知、并发重试）。
- 并发路径必须 `-race`（本地无 cgo 时如实标注）。
- router 测试不得只断言状态码，库存与金额要有对账断言（I-19）。
- 买家区 6 路由：路由层 `RequireRole(RoleBuyer)` **且**委派 handler 内复核（AuthGuard 委派先于组门禁，仅路由层拦不住）。

## 提交前自检清单（逐项打勾）

- [ ] gofmt 空、vet/build/test 退出 0
- [ ] 跃迁全走"条件 UPDATE + RowsAffected"
- [ ] 并发路径有闸门与计数断言
- [ ] DB 错误仅唯一键判重复
- [ ] 覆盖率单命令复跑不低于基线
- [ ] `demo-mode` 门禁仍在
- [ ] 提交编号在 dev_log 存在且类型匹配
- [ ] `docs/events.md` 已登记
- [ ] INV-1~4 有数字

## 红线

- 不编造执行结果/覆盖率/部署记录；未执行的命令不写"已验证"。
- 不删失败测试、不弱化断言；不为让测试变绿改测试语义。
- 密钥令牌个人信息不进代码日志。
- 不扩大范围（购物车、多商户、真实物流/支付）。
- 模拟渠道与框架能力不得包装成自研。
