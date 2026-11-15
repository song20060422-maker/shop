# AGENTS.md

> 本文件是 AI 会话的**入口说明**（只讲"读什么、守什么底线"）。
> **详细规则拆分在 `.agents/` 目录**，按主题分文件，会话按需读取；冲突时以本文件为准。

## 会话启动顺序

1. 本文件（入口说明）
2. `.agents/00-index.md`（约束文件目录）
3. 按任务从 `.agents/` 按需读取分文件：
   - 开发角色：`.agents/agent-dev.md`（v6.0 全文任务书）→ 按 `01`~`04` 分文件查详细规则
   - 测试角色：`.agents/agent-test.md`（v2.0 任务书）
   - 验收角色：`.agents/agent-acceptance.md`（v2.0 任务书）
   - 提交前：`.agents/02-workflow-git.md` + 仓库根 `commit-convention.md`

## 详细规则文件清单（`.agents/`）

| 文件 | 内容 |
|---|---|
| `01-business-rules.md` | 技术栈边界、状态机、BR/E/INV 矩阵、并发防御、实现模式 |
| `02-workflow-git.md` | 工作项、提交规范、编号核验、门禁、交接七段 |
| `03-quality-security.md` | 覆盖率、race、自检清单、安全红线 |
| `04-ops-rules.md` | 配置部署、冲突裁决、缺陷修复 |
| `README.md` / `versions/` | 版本索引、历史归档（只读） |

## 硬性底线（摘要，细则在各分文件）

- 状态跃迁只用"条件 UPDATE + RowsAffected 校验"，禁止先查后写。
- 数据落 MySQL，禁止内存态冒充持久化；异常记录只增不改。
- 覆盖率核心包 ≥85%，单命令 profile，`-p 1` 串行，并发路径 `-race`。
- 提交 `[类型] [版本号] [编号] 中文描述`，feat/fix 编号必填飞书 `#数字` 并经 meegle 核验；分支 `dev/01-shop`，禁直推 `main`。
- 不编造执行结果；未执行的命令不写"已验证"。
- 开发/测试/验收三角色必须独立会话，测试与验收读原始题面与被测版本代码。

## 留痕与裁决

- 会话真实加载了哪些约束文件，登记 `.agents/README.md` 的"加载记录"。
- 修改约束：先归档 `versions/<name>.vN.md`，再更新生效文件与索引，禁止覆盖历史。
- 工作项与事件：`docs/dev_log.md`（T 编号）、`docs/events.md`、`docs/review_records.md`（缺陷）。
- 冲突裁决：原始题面 > 可审计事实 > 文档；细则见 `.agents/08-conflict-rules.md`。
- 文档层分工：`AGENT.MD` 管 `aiDoc/` 文档层；本文件管 `.agents/` 约束层。
