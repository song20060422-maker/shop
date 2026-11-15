# 00 · 约束文件目录

> 详细规则拆分为 4 个核心文件，会话按需读取对应文件（不必全读）。
> 所有条款都是本项目实际约束；与 `AGENTS.md` 冲突时以 `AGENTS.md` 为准。

| 文件 | 主题 | 何时必读 |
|---|---|---|
| `01-business-rules.md` | 技术栈边界、状态机、BR-01~12、E-01~12、INV-1~4、并发防御、实现模式 | 动业务字段 / 写 service 流程 |
| `02-workflow-git.md` | 工作项、分支、提交规范、编号核验、门禁、证据、交接七段 | 每次提交前 |
| `03-quality-security.md` | 覆盖率 ≥85%、单命令 profile、race、自检清单、红线 | 每次自测 / 声称覆盖率 |
| `04-ops-rules.md` | 配置部署、冲突裁决、只读对象、缺陷修复流程 | 改配置 / 矛盾 / 修缺陷 |

## 角色任务书（独立会话整体下发）

| 文件 | 角色 |
|---|---|
| `agent-dev.md` | 开发 Agent（v6.0 全文） |
| `agent-test.md` | 独立测试 Agent（v2.0） |
| `agent-acceptance.md` | 独立验收 Agent（v2.0） |

历史版本与加载记录见 `README.md` 与 `versions/`（只读）。提交规范来源副本见仓库根 `commit-convention.md`。
