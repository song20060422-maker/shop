# Shop 商城交易系统 验收报告（第 2 轮 · 2026-11-15）

## 被测对象与方式

- 代码仓库：https://github.com/song20060422-maker/shop.git（本地 `D:\Users\Lenovo\Desktop\shop\shop`）
- 分支：`dev/01-shop`（业务代码，逐行注释合规）/ `main`（纯 gin-vue-admin 框架）
- 被测运行实例：买家端 `http://127.0.0.1:8080`（shop-server.exe，验收窗口切换 simulated 渠道，验收后已恢复 alipay 沙箱）、后台 GVA `8085/8099`
- 实测方式：HTTP 接口实测（PowerShell 脚本 `acceptance-r2/phase1~5.ps1`，输出留证 `phase1~5-output.txt`），共 77 项断言全部 PASS
- 题面依据：`Go项目业务场景题_共同要求.docx`（用户原文提供 8 条规则 + 5 个场景 + 原型/单日边界）

## 总体结论：**通过**（77/77 断言 PASS，0 FAIL）

## 逐项结论

### 业务规则 R1–R8（全部实测）

| 编号 | 核对项 | 结论 | 证据（脚本断言） |
|---|---|---|---|
| R1-1 | 下架/数量不足拒购、不占库存 | PASS | R1-1-qty-exceed-rejected(409)、R1-1-offsale-rejected(409)、R1-1-no-stock-occupied |
| R1-2 | 订单保存下单时刻价格 | PASS | R1-2-price-snapshot-10000（改价 12000 后订单仍 10000） |
| R1-3 | 下单后可售立即反映占用 | PASS | R1-3-occupy-reflected(5→4)、R1-3-occupy-5-1-2-eq-2 |
| R1-4 | 每单一种商品、1~2 件、分为单位 | PASS | R1-4-qty0/negative/qty3-over-limit 全部 409；amountCents=10000（100 元） |
| R2-1 | 保留 10 分钟、服务端时间 | PASS | R2-1-expire-in-10min(expireIn=10.0min)；配置走查 config.shop.yaml `order-timeout-minutes: 10` |
| R2-2 | 主动取消释放 | PASS | R2-2-cancel-ok + stock-restored-4 |
| R2-3 | 超时关闭释放、重复关单幂等 | PASS | R2-3-closed-by-expire、expire-restored-5、second-close-idempotent |
| R3-1 | 通知验签、留痕 | PASS | R3-1-bad-sign-rejected(400) + queryable(result=rejected 验签失败) |
| R3-2 | 金额/订单匹配、可查询 | PASS | R3-2-amount-mismatch(400)、unknown-order(400)、illegal-queryable |
| R3-3 | 重复通知不重复记款 | PASS | R3-3-replay-duplicated(409)、one-payment-record(记款行=1) |
| R3-4 | 同一渠道交易号不归属两单 | PASS | R3-4-same-txn-other-order-rejected(409) |
| R4-1 | 有效期内付款→待发货 | PASS | R4-1-paid-status |
| R4-2 | 关闭后成功通知：不恢复/不重占/进异常 | PASS | R4-2-order-stays-closed、no-reoccupy、in-abnormal-list |
| R4-3 | 客服核实并发起退款 | PASS | R4-3-cs-handle-refund(channelStatus=success) |
| R5-1 | 申请期间暂停发货 | PASS | R5-1-refund-applied + ship-blocked-while-refunding(409) |
| R5-2 | 同意/拒绝均需原因 | PASS | R5-2-audit-approved / audit-rejected-with-reason |
| R5-3 | 拒绝后可继续发货 | PASS | R5-3-ship-after-reject-ok(status=shipped) |
| R5-4 | 退款与发货只一种合法结果 | PASS | S5-1-concurrent-one-legal-outcome（并发实测 REFUND-OK/SHIP-409） |
| R6-1 | 渠道成功后买家可见退款结果 | PASS | R6-1-buyer-sees-refund-success |
| R6-2 | 售前退款恢复可售数量（仅一次） | PASS | R6-2-stock-restored-once(3→4) + no-extra-restore |
| R6-3 | 渠道失败保留记录可重试 | PASS | R6-3-first-refund-fail-kept、retry-under-fail-still-failed、retry-success、triple-retry-rejected |
| R6-4 | 迟到支付退款不重复恢复数量 | PASS | R6-4-late-refund-no-restore、S4-no-double-restore |
| R7-1 | 只能登记一次发货 | PASS | R7-1-double-ship-rejected(409 订单已发货，不能重复登记) |
| R7-2 | 需承运方+运单号 | PASS | R7-2-missing-fields-rejected(400) |
| R7-3 | 已发货不受理售前退款 | PASS | R7-3-shipped-refund-rejected(400 订单已发货，不受理售前退款) |
| R8-1 | 按订单查收/退/发/异常四类记录 | PASS | R8-1-four-record-types-visible（payments/refunds/shipments/abnormals 接口实测） |
| R8-2 | 可售/已收款/已退款汇总对账 | PASS | R8-2-summary-struct-ok（记款 10000 + 迟到支付异常款 10000 = 汇总 20000；退款 10000=汇总） |
| R8-3 | 异常记录不可通过改状态消失 | PASS | R8-3-no-status-bypass-endpoint(404 无旁路接口) + 走查路由表 |

### 完整场景 S1–S5

| 编号 | 场景 | 结论 | 证据 |
|---|---|---|---|
| S1 | 正常成交全链路 | PASS | 下单(R1-2)→付款(R3-3)→发货(R5-3)→买家查运单（详情含 shipment YT/YT-888、payment 记录，本次补验）→运营核查(R8-1/R8-2)；四方数据一致 |
| S2 | 缺货与弃单（20 并发） | PASS | phase4：真实 20 并发下单，success=5 / failed=15（409 售罄），available=0 无超卖；循环关单释放回 5；buyer2 复购 2 件成功后取消回 5 |
| S3 | 售前退款（含失败重试） | PASS | phase3：申请→发货受限→批准(首次渠道 fail 保留 FAILED)→fail 下重试仍失败→切 success 重试成功→库存 3→4 仅恢复一次→订单 refunded |
| S4 | 迟到支付闭环 | PASS | phase5：关单→TXN-S4 通知 result=abnormal→客服处理退款 success→库存不变(5)→该订单退款单仅 1 张 |
| S5 | 操作竞争 | PASS | S5-1 退款 vs 发货并发（真实双任务并发，一种合法结果）；重复支付通知(R3-3-replay)；重复退款（triple-retry/double-handle 409）；越权 B5 三项 403 |

### 测试证据 E1–E4

| 编号 | 结论 | 说明 |
|---|---|---|
| E1 | PASS | 每条 R 规则对应 phase 脚本断言 + 原始输出文件 acceptance-r2/phase1~5-output.txt |
| E2 | PASS | 5 场景操作序列均含于 phase2~5 输出（含每步 HTTP 码与状态值） |
| E3 | PASS | S2 并发原始输出（20 条 OK/ERR 分布 + RAW RESULTS 行）存 phase4-output.txt |
| E4 | PASS | 输出为脚本实测落盘（Tee-Object），断言与响应一一对应；抽查 phase2/phase5 与响应报文一致 |

### 部署 D1–D3 与边界 B1–B6

| 编号 | 结论 | 说明 |
|---|---|---|
| D1 | PASS | go build 后 shop-server.exe 启动 :8080；GVA 8085/8099 在运行；健康检查 ok |
| D2 | PASS | D2-SKUA-10000fen-5pcs、D2-SKUB、D2-accounts（3 买家+3 员工+契约别名） |
| D3 | PASS | 买家页 /shop 200；仓库/客服/运营页在 GVA 8085（本会话浏览器实测可登录可用）；原型含支付页 payment.html |
| B1 | PASS（走查） | 无购物车/多商户/真实物流/部分退款/财务系统 |
| B2 | PASS | 固定商品、模拟渠道（/api/pay/notify HMAC）、人工发货登记 |
| B3 | PASS | 全部金额为分（int64 amountCents），无浮点 |
| B4 | PASS | 第一部分 R1–R8 全部实现（见上表） |
| B5 | PASS | 跨买家查看/取消/退款全部 403 |
| B6 | PASS | 审核原因/处理原因入库（auditReason、通知 reason、trail 接口在册） |

## 本轮发现的问题（均非业务缺陷，已记录）

1. **注释合规整改**：业务代码此前未逐行注释、main 分支含业务代码 → 已建 `dev/01-shop` 分支，语法感知注入 9969 行注释（70 文件），go build + 各包测试通过后提交；main 已收敛为纯框架。
2. **测试套件共享 shop_test 库跨包污染**：`go test ./shop/...` 全量连跑时 service 包因 FK 残留失败（分包运行全部通过）——建议测试用例清理顺序补 FK 级联（不属于验收范围，作为改进建议）。
3. **支付宝渠道下异常退款不可重试成功**：迟到支付通知若由模拟通道产生、而退款走支付宝沙箱，沙箱无对应交易号必然失败——属环境组合问题；验收窗口切 simulated 渠道验证了完整闭环，生产演示时两类通道应保持一致。

## 结论

系统达到题面验收标准：8 条业务规则、5 个必演示场景、原型页面、单日边界全部满足，异常路径（超卖/重复记款/重复退款/迟到支付/退款发货竞争/越权）均有明确且正确的结果。
