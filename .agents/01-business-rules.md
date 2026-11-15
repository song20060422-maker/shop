# 01 · 业务规则与实现模式（核心）

> 动订单状态、库存、支付、退款、发货字段，或写/改 `shop/service` 业务流程前必读。与 `AGENTS.md` 冲突时以 `AGENTS.md` 为准。

## 技术栈与改动边界

- 框架 gin-vue-admin v2.9.2（Go 1.26 + Gin + GORM），模块名 `shop`；MySQL 8（3307，业务库 `shop`）；服务端口 8099；商城页 `server/shop/web/`（后端直出，非 Vue 构建链）。
- 分层：`api` 绑定/身份/错误码；`service` 事务/状态机/库存守恒/渠道/审计；`model` 表与 DTO（金额禁浮点）；`middleware` 角色门禁（只信 token）；`initialize` 装配/迁移/种子（幂等，失败即启动失败）；`router` 角色分组；`web` 内置演示页；`utils` 单号/HMAC/bcrypt（密钥不入库、不打日志）。
- 框架侧只允许动 `server/initialize/router_biz.go` 与 config 的 `shop` 段；不改 GVA `system`；不另起 main/端口/连接池；数据必须落 MySQL，禁止内存态冒充持久化；测试钩子受 `shop.demo-mode` 门禁，生产一律 403。

## 订单状态机

```
pending_payment → paid → shipped
        │  ↑(审核拒绝)
        ├→ cancelled / closed(超时)
        └→ refund_requested → refunded(渠道成功) / refund_failed(渠道失败,可重试)
```

跃迁只写在 `service/order.go`、`payment.go`、`refund.go`、`shipment.go`；完成判据必须可观察（状态、数字、记录、错误码），不接受"已实现"。

## 业务规则矩阵（BR-01 ~ BR-12）

| BR | 要点 | 落点 |
|---|---|---|
| BR-01/02 | 条件 UPDATE 扣减 + 快照 + 占用流水 | `order.go Create` |
| BR-03 | 10 分钟超时（服务端时间）、释放幂等 | `Cancel/CloseExpired` |
| BR-04~06 | HMAC 验签 → 金额 → 幂等（唯一键），拒绝可查 | `payment.go` |
| BR-07 | 有效付款推进 paid；迟到转 `abnormal_payments` | `payment.go Notify` |
| BR-08~10 | 退款申请互斥发货；渠道确认才 refunded；迟到退款不恢复库存 | `refund.go` |
| BR-11 | 发货唯一（`shipments` 唯一键） | `shipment.go Ship` |
| BR-12 | 核查汇总=明细之和；异常无"改单消失"通路 | `report.go` |

## 异常矩阵与不变量

- E-01~12（并发争抢、取消/超时竞争、重复通知、验签失败、迟到支付、退款×发货竞争、渠道失败重试、已发货退款、越权、重复发货、改单消除异常）均有明确拒绝/幂等/互斥行为；完整判据见 `shop_word/references/business-rules.md`（只读）。
- 不变量（每次改动复核）：INV-1 初始可售 = 可售 + 各态占用 + 已出库；INV-2 同一渠道交易号至多归属一单；INV-3 每单至多一次成功发货；INV-4 每单至多一次成功退款。

## 三条固定防御（审计 I-18 固化）

1. **渠道调用闸门**：`Retry`/`Audit(approve)`/`HandleAbnormal` 事务内 `SELECT ... FOR UPDATE` 行锁 + 入口状态校验后才调渠道；后到者被拒（`ErrRefundNotRetryable`）。
2. **错误分类**：仅唯一键冲突（errno 1062，`isDuplicateEntry`）判"重复"，其它 DB 错误原样上抛。
3. **RowsAffected 校验**：条件更新推进状态必须校验受影响行数，不为 1 即回滚。

## 实现模式（写 service 必读）

- **状态跃迁**：`UPDATE ... WHERE 状态=旧值 [AND 条件]` + `RowsAffected==1`；失败读最新态区分重复/过期/非法；禁止先 SELECT 再 UPDATE；禁止 handler 内改状态。
- **幂等**：幂等键落唯一索引（`payments.channel_txn_id`、`payments.order_no`、`shipments.order_no`），冲突按错误分类处理。
- **一次性资源闸门**：库存恢复用 `restored` 布尔闸门（`UPDATE ... WHERE restored=false` 且 RowsAffected==1 才释放）；迟到支付退款绝不走库存恢复。
- **审计留痕**：每次跃迁调 `logOp(tx, orderNo, actorID, role, action, detail)`，与业务同事务提交，禁止事后补写。
- **演示开关**：渠道退款结果可切 `success|fail` 并带 `RefundCalls()` 计数；超时可缩短但判断用服务端时间。
- **反模式（出现即打回）**：handler/service 混层写事务；用 panic/忽略错误代替业务码；内存 map 冒充持久化；为凑覆盖率给工具函数写用例而业务分支裸奔。
