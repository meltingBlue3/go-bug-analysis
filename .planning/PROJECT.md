# 离线版禅道 Bug 分析工具 (go-bug-analysis)

## What This Is

一款面向测试/质量团队的离线 Bug 分析工具。用户从禅道导出 CSV，导入工具后即可获得 Bug 状态总览、瓶颈分析、模块质量热力图、人员负载分布及自动生成日报文本。Go 单文件编译，启动即用，无需数据库或网络连接。

## Core Value

**让 Bug 积压和修复瓶颈一眼可见** — 测试人员导入 CSV 后，30 秒内看到当前质量状态全貌，不再依赖手工统计。

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] CSV 导入：支持禅道导出的 CSV 文件，自动检测 GBK/UTF-8 编码
- [ ] 核心看板：新增/修复 Bug 数、Bug 库存（总数/激活/待验证）、严重程度分布（饼图/柱状图切换）
- [ ] 瓶颈与时效分析：已修复 Bug 平均修复时长、P50 修复时长、修复耗时分布；未修复 Bug 积压时长倒序列表
- [ ] 模块与质量分析：各模块 Bug 总数/激活数、模块×严重程度热力图、模块 Bug 趋势（近 7/30 天）
- [ ] 人员与负载分析：激活 Bug 指派人分布、总 Bug 指派人分布
- [ ] 日报辅助输出：自动生成 Markdown/纯文本日报（新增/修复/净变化、严重 Bug 风险、瓶颈概览），支持一键复制
- [ ] 单文件分发：Go embed 打包前端资源，编译为单个可执行文件
- [ ] 启动自动打开浏览器

### Out of Scope

- 多用户协作 — 定位为本地单人工具
- 禅道数据库/API 直连 — 仅通过 CSV 离线分析
- 在线数据同步 — 离线工具，不涉及网络通信
- 多语言 UI — 仅中文，目标用户为中文测试团队
- 多产品筛选 — 一次导入分析全量数据，不区分产品

## Context

- **数据源**：禅道(ZenTao) Bug 管理系统的 CSV 导出文件
- **CSV 字段**：Bug编号、所属产品、所属模块、Bug标题、严重程度、优先级、Bug类型、Bug状态、由谁创建、创建日期、指派给、指派日期、解决者、解决方案、解决日期、由谁关闭、关闭日期 等 36 个字段
- **CSV 编码**：禅道导出通常为 GBK 编码，需自动检测并转换为 UTF-8
- **数据规模**：典型 1–5 万条记录，样本数据约 51,000 行
- **日期判定**："今日"使用系统日期，"今日新增"以 CSV 中创建日期字段匹配当天日期
- **Bug 状态值**：激活、已解决、已关闭 等（来自禅道状态体系）
- **严重程度值**：1(致命)、2(严重)、3(一般)、4(轻微)（数字型）

## Constraints

- **Tech Stack**: Go 1.20+, net/http + embed.FS, 原生 HTML/JS, ECharts — 无框架依赖，单文件分发
- **Performance**: 启动时间 < 2s，5 万条 CSV 解析可接受
- **Distribution**: 单文件可执行程序（exe / binary），不依赖外部运行时
- **Offline**: 完全离线运行，ECharts 等前端资源内嵌

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Go + embed.FS 而非 Electron | 单文件小体积、启动快、无需 Node 运行时 | — Pending |
| ECharts 而非 D3/Chart.js | 开箱即用、热力图支持好、中文生态成熟 | — Pending |
| 原生 HTML/JS 而非 Vue/React | 减少构建复杂度、embed 打包简单 | — Pending |
| GBK/UTF-8 自动检测 | 禅道导出默认 GBK，降低用户使用门槛 | — Pending |
| 热力图维度：模块 × 严重程度 | 快速定位哪些模块有最严重的 Bug 集中 | — Pending |
| 启动时自动打开浏览器 | 降低使用门槛，运行即见界面 | — Pending |

---
*Last updated: 2026-02-11 after initialization*
