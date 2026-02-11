# 离线版禅道 Bug 分析工具 (go-bug-analysis)

## What This Is

一款面向测试/质量团队的离线 Bug 分析工具。用户从禅道导出 CSV，导入工具后即可获得 Bug 状态总览（KPI 卡片、严重程度分布）、瓶颈分析（修复时效、积压排名）、模块质量热力图与趋势、人员负载分布，以及自动生成日报文本。Go 单文件编译（~9MB），启动即用，无需数据库或网络连接。

## Core Value

**让 Bug 积压和修复瓶颈一眼可见** — 测试人员导入 CSV 后，30 秒内看到当前质量状态全貌，不再依赖手工统计。

## Requirements

### Validated

- ✓ CSV 导入：支持禅道导出的 CSV 文件，自动检测 GBK/UTF-8 编码 — v1.0
- ✓ 核心看板：新增/修复 Bug 数、Bug 库存、严重程度分布（饼图/柱状图切换） — v1.0
- ✓ 瓶颈与时效分析：修复时长统计（均值/P50/分布）、未修复 Bug 积压排名 — v1.0
- ✓ 模块与质量分析：模块 Bug 总数/激活数、模块×严重程度热力图、模块趋势（7/30天） — v1.0
- ✓ 人员与负载分析：激活/总量 Bug 按指派人分布 — v1.0
- ✓ 日报辅助输出：自动生成 Markdown/纯文本日报，支持一键复制 — v1.0
- ✓ 单文件分发：Go embed 打包前端资源，编译为单个可执行文件 — v1.0
- ✓ 启动自动打开浏览器 — v1.0

### Active

(None — next milestone requirements to be defined via `/gsd-new-milestone`)

### Out of Scope

- 多用户协作 — 定位为本地单人工具
- 禅道数据库/API 直连 — 仅通过 CSV 离线分析
- 在线数据同步 — 离线工具，不涉及网络通信
- 多语言 UI — 仅中文，目标用户为中文测试团队
- 多产品筛选 — 一次导入分析全量数据，不区分产品
- 移动端适配 — 桌面浏览器使用场景

## Context

**Current State:** Shipped v1.0 MVP with ~3,245 LOC (Go 1,345 + JS 1,050 + CSS 615 + HTML 235).
**Tech stack:** Go 1.25, net/http + embed.FS, native HTML/JS (IIFE), ECharts 5.x, golang.org/x/text.
**Data:** 禅道 CSV 导出文件，GBK/UTF-8 编码，典型 1–5 万条记录。
**Known tech debt:**
- All Bug fields stored as strings (date parsing deferred to analysis modules)
- No comprehensive test suite beyond CSV parser tests (4 tests)
- No CI/CD pipeline

## Constraints

- **Tech Stack**: Go 1.20+, net/http + embed.FS, 原生 HTML/JS, ECharts — 无框架依赖，单文件分发
- **Performance**: 启动时间 < 2s，5 万条 CSV 解析可接受
- **Distribution**: 单文件可执行程序（exe / binary），不依赖外部运行时
- **Offline**: 完全离线运行，ECharts 等前端资源内嵌

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Go + embed.FS 而非 Electron | 单文件小体积、启动快、无需 Node 运行时 | ✓ Good — ~9MB binary, <2s startup |
| ECharts 而非 D3/Chart.js | 开箱即用、热力图支持好、中文生态成熟 | ✓ Good — heatmap, pie, bar, line all used |
| 原生 HTML/JS 而非 Vue/React | 减少构建复杂度、embed 打包简单 | ✓ Good — IIFE pattern works well, ~1050 LOC JS |
| GBK/UTF-8 自动检测 | 禅道导出默认 GBK，降低用户使用门槛 | ✓ Good — transparent to users |
| 热力图维度：模块 × 严重程度 | 快速定位哪些模块有最严重的 Bug 集中 | ✓ Good — clear hotspot visualization |
| 启动时自动打开浏览器 | 降低使用门槛，运行即见界面 | ✓ Good — cross-platform (Win/Mac/Linux) |
| All Bug fields as string | 避免过早日期解析，分析模块按需解析 | ⚠️ Revisit — adds parsing overhead in analysis |
| strings.Builder for reports | 无需模板库，简洁直接 | ✓ Good — lightweight, no dependency |
| IIFE + window.BugAnalysis API | 跨模块通信无需 ES modules | ✓ Good — simple, compatible |

---
*Last updated: 2026-02-11 after v1.0 milestone*
