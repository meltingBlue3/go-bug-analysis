# Roadmap: 离线版禅道 Bug 分析工具

## Overview

从零搭建一个 Go 单文件离线 Bug 分析工具：先建好 Go 项目脚手架与 CSV 解析管线（Phase 1），再实现核心 KPI 看板与瓶颈/人员分析（Phase 2），然后补齐模块热力图与趋势（Phase 3），最后生成日报文本（Phase 4）。四个阶段完成后，用户导入禅道 CSV 即可在 30 秒内看到质量全貌并复制日报。

## Phases

- [x] **Phase 1: Foundation & CSV Pipeline** - Go 脚手架、嵌入资源、HTTP 服务与 CSV 解析管线 ✓ (2026-02-11)
- [x] **Phase 2: Core KPI & Bottleneck Analysis** - KPI 卡片、严重程度图表、修复时效、人员负载 ✓ (2026-02-11)
- [ ] **Phase 3: Module Heatmap & Trend** - 模块统计、热力图、趋势分析
- [ ] **Phase 4: Daily Report Generation** - 日报自动生成与剪贴板复制

## Phase Details

### Phase 1: Foundation & CSV Pipeline
**Goal**: Users can launch a single executable and import Zentao CSV data ready for analysis
**Depends on**: Nothing (first phase)
**Requirements**: INF-01, INF-02, INF-03, INF-04, CSV-01, CSV-02, CSV-03, CSV-04
**Success Criteria** (what must be TRUE):
  1. User double-clicks the executable, default browser opens automatically, and the app is ready within 2 seconds
  2. User can select a Zentao CSV file via file picker and see it accepted for processing
  3. System correctly parses both GBK and UTF-8 encoded CSV files without user intervention
  4. System displays clear Chinese error messages when required columns are missing or data format is invalid
**Plans**: 2 plans

Plans:
- [x] 01-01-PLAN.md — Go 项目脚手架：embed.FS、HTTP 服务器、自动打开浏览器、前端仪表盘外壳 (INF-01, INF-02, INF-03)
- [x] 01-02-PLAN.md — CSV 上传管线：编码检测、表头映射、数据验证、前端集成 (CSV-01..04, INF-04)

### Phase 2: Core KPI & Bottleneck Analysis
**Goal**: Users can see the full quality status picture — bug inventory, severity distribution, fix time bottlenecks, and personnel load
**Depends on**: Phase 1
**Requirements**: KPI-01, KPI-02, KPI-03, KPI-04, KPI-05, AGE-01, AGE-02, AGE-03, AGE-04, AGE-05, WRK-01, WRK-02
**Success Criteria** (what must be TRUE):
  1. User sees today's and yesterday's new/fixed bug counts as KPI cards immediately after CSV import
  2. User sees bug inventory (total, active, pending-verification) and can toggle severity chart between pie/bar and "new bugs only" vs "all bugs"
  3. User sees fix time stats (average, P50, distribution buckets) and a descending list of unresolved bugs ranked by backlog age
  4. User sees active and total bug distribution per assignee in chart form
**Plans**: 3 plans

Plans:
- [x] 02-01-PLAN.md — KPI 卡片（今日/昨日新增修复、库存）+ 严重程度分布图（饼图/柱状图切换、全部/仅新增筛选）(KPI-01..05)
- [x] 02-02-PLAN.md — 修复时效统计（均值/P50/分布桶）+ 未解决 Bug 积压排名表 (AGE-01..05)
- [x] 02-03-PLAN.md — 人员工作负载分布图（激活/总量按指派人）(WRK-01, WRK-02)

### Phase 3: Module Heatmap & Trend
**Goal**: Users can identify which modules have the worst bug concentration and track quality trends
**Depends on**: Phase 2
**Requirements**: MOD-01, MOD-02, MOD-03
**Success Criteria** (what must be TRUE):
  1. User sees per-module bug count breakdown (total and active) in a sortable view
  2. User sees a Module × Severity heatmap that visually highlights bug concentration hotspots
  3. User sees module bug trend lines for the last 7 and 30 days
**Plans**: 1 plan

Plans:
- [ ] 03-01: Module statistics table, ECharts heatmap, and trend line charts (MOD-01..03)

### Phase 4: Daily Report Generation
**Goal**: Users can generate and copy a ready-to-send daily quality report in seconds
**Depends on**: Phase 2, Phase 3
**Requirements**: RPT-01, RPT-02, RPT-03, RPT-04, RPT-05
**Success Criteria** (what must be TRUE):
  1. System auto-generates a structured daily report in both Markdown and plain text format after CSV import
  2. Report includes today's new/fixed/net-change counts, highlights severe bug risks (severity 1-2 with long backlog), and summarizes bottleneck modules and personnel
  3. User can one-click copy the generated report to clipboard and paste it into any chat/document tool
**Plans**: 1 plan

Plans:
- [ ] 04-01: Report template engine, Markdown/plain-text output, clipboard copy (RPT-01..05)

## Progress

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation & CSV Pipeline | 2/2 | Complete | 2026-02-11 |
| 2. Core KPI & Bottleneck Analysis | 3/3 | Complete | 2026-02-11 |
| 3. Module Heatmap & Trend | 0/1 | Not started | - |
| 4. Daily Report Generation | 0/1 | Not started | - |
