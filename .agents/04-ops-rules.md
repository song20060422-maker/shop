# 04 · 配置部署 / 冲突裁决 / 缺陷修复

> 改配置或部署、需求矛盾、修缺陷时必读。

## 一、配置与部署

| 项 | 值 | 说明 |
|---|---|---|
| `addr` | ":8099" | 监听端口 |
| `demo-mode` | true | 演示钩子可用；生产必须 false |
| `mysql` | 3307 / `shop` | 业务库连接 |
| `pay-notify-secret` | 演示密钥 | HMAC 验签，不入代码、不打日志 |

坑：配置初始化会重写 config.yaml，其它段也会出现 `addr:`——解析端口必须限定 `shop` 段（见 `scripts/run_acceptance.py read_base_url`）。

- CI `.github/workflows/shop-ci.yml`（lint/build/test/coverage）；容器 `deploy/shop/`（Dockerfile + compose + 冒烟）；未真实部署前 DELIVERY/README 保持"未部署"标注。
- 本地运行：`cd server; go build -o bin\shop-server.exe .; .\bin\shop-server.exe`（8099）；验收 `py -3 shop\scripts\run_acceptance.py --no-build --no-start`。演示账号 `ops1/wh1/cs1/buyer1~3`，口令 `shop123`，商城页 `http://127.0.0.1:8099/shop`。

## 二、冲突与例外

- 需求材料只读；发现矛盾停下来提请学生裁决并记入 `docs/assumptions.md`，不自行改需求，未确认的假设不落地。
- 优先级：`AGENTS.md` > `.agents/01`~`04` 约束文件 > `versions/` 历史版本。
- 与交付实现冲突：以仓库中可复核的代码、profile、CI 运行记录为准，文档随后同步修正（不得改文档掩盖事实）。
- 只读对象（永不修改）：`.agents/versions/`（新增 vN + 更新 README 索引）、`shop_word/archive/`、检查者脚本（`shop-trade-check/scripts/acceptance_check.py` 等）、审计报告、GVA `system` 模块、仓库根 `commit-convention.md`、`C:\Users\Lenovo\Desktop\shop_word` 同名旧副本。
- Windows GBK 控制台：Python 入口强制 `sys.stdout.reconfigure(encoding="utf-8")`，子进程注入 `PYTHONIOENCODING=utf-8`（I-24）。

## 三、缺陷修复流程

1. **登记**：`docs/review_records.md`（问题/复现/期望实际/违反的 BR/版本哈希）。
2. **回归先行**：先写失败回归测试并跑出真实失败。
3. **根因定位**：到行级（`log/zap-error.log`、慢查、并发时序）。
4. **最小改动**：修正（不顺手重构）。
5. **复测**：回归转绿 + `go test ./shop/... -count=1` + 受影响验收套件，exit 0。
6. **回填留痕**并登记 `docs/events.md`。
7. **提交**：`[fix] [版本] [编号]`（引用审计编号如 I-18）。

教训库（已固化）：断言写业务语义不写魔法数字；委派 handler 必须复核角色；覆盖率只用单命令；验收脚本端口解析限定 `shop` 段；渠道重试用事务行锁闸门 + 调用计数；仅唯一键冲突判重复支付；入口强制 UTF-8。
修不了的如实标注"未解决+原因"；不为让测试变绿改测试语义。
