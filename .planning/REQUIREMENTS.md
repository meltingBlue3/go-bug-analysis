# Requirements: 离线版禅道 Bug 分析工具

**Defined:** 2026-02-11
**Core Value:** 让 Bug 积压和修复瓶颈一眼可见 — 导入 CSV 后 30 秒内看到质量状态全貌

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### CSV 导入

- [ ] **CSV-01**: User can select and upload a CSV file exported from Zentao via file picker
- [ ] **CSV-02**: System auto-detects GBK/UTF-8 encoding and converts to UTF-8 for processing
- [ ] **CSV-03**: System auto-recognizes Zentao CSV column headers and maps to internal fields
- [ ] **CSV-04**: System shows clear error messages when required fields are missing or data format is invalid

### 核心看板

- [ ] **KPI-01**: User can see today's and yesterday's new bug count
- [ ] **KPI-02**: User can see today's and yesterday's fixed bug count
- [ ] **KPI-03**: User can see bug inventory: total count, active count, pending-verification count
- [ ] **KPI-04**: User can see severity distribution chart with pie/bar chart toggle
- [ ] **KPI-05**: User can switch severity view between "new bugs only" and "all bugs"

### 瓶颈与时效分析

- [ ] **AGE-01**: User can see average fix time for resolved bugs
- [ ] **AGE-02**: User can see P50 (median) fix time for resolved bugs
- [ ] **AGE-03**: User can see fix time distribution breakdown (0-1 day / 2-3 days / 4-7 days / 7+ days)
- [ ] **AGE-04**: System calculates backlog age for each unresolved bug (days since creation)
- [ ] **AGE-05**: User can see unresolved bugs listed in descending order by backlog age

### 模块与质量分析

- [ ] **MOD-01**: User can see total bug count and active bug count per module
- [ ] **MOD-02**: User can see a heatmap of Module × Severity showing bug concentration
- [ ] **MOD-03**: User can see module bug trend over the last 7 and 30 days

### 人员与负载分析

- [ ] **WRK-01**: User can see distribution of currently active bugs by assignee
- [ ] **WRK-02**: User can see distribution of total bugs by assignee

### 日报辅助输出

- [ ] **RPT-01**: System auto-generates daily report text in Markdown and plain text format
- [ ] **RPT-02**: Daily report includes today's new / fixed / net change counts
- [ ] **RPT-03**: Daily report highlights severe bug risks (severity 1-2 with long backlog)
- [ ] **RPT-04**: Daily report includes bottleneck module and personnel overview
- [ ] **RPT-05**: User can one-click copy the generated report to clipboard

### 基础设施

- [ ] **INF-01**: Application compiles to a single executable binary via Go embed
- [ ] **INF-02**: Application auto-opens default browser on startup
- [ ] **INF-03**: Application starts in under 2 seconds
- [ ] **INF-04**: Application handles CSV files with 10,000–50,000 rows without degradation

## v2 Requirements

Deferred to future release. Tracked but not in current roadmap.

### 增强分析

- **ENH-01**: User can pick a custom reference date instead of system date
- **ENH-02**: User can filter analysis by product (所属产品)
- **ENH-03**: User can export charts as images

### 扩展导入

- **EXT-01**: Support drag-and-drop CSV import
- **EXT-02**: Support importing multiple CSV files and merging data

## Out of Scope

| Feature | Reason |
|---------|--------|
| 禅道数据库/API 直连 | 离线工具定位，仅 CSV 导入 |
| 多用户协作 | 单人本地工具 |
| 在线数据同步 | 无网络通信 |
| 英文/多语言 UI | 目标用户为中文团队 |
| 移动端适配 | 桌面浏览器使用场景 |
| 缺陷管理功能 | 定位为分析工具，不做流程管理 |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| CSV-01 | Phase 1 | Pending |
| CSV-02 | Phase 1 | Pending |
| CSV-03 | Phase 1 | Pending |
| CSV-04 | Phase 1 | Pending |
| KPI-01 | Phase 2 | Pending |
| KPI-02 | Phase 2 | Pending |
| KPI-03 | Phase 2 | Pending |
| KPI-04 | Phase 2 | Pending |
| KPI-05 | Phase 2 | Pending |
| AGE-01 | Phase 2 | Pending |
| AGE-02 | Phase 2 | Pending |
| AGE-03 | Phase 2 | Pending |
| AGE-04 | Phase 2 | Pending |
| AGE-05 | Phase 2 | Pending |
| MOD-01 | Phase 3 | Pending |
| MOD-02 | Phase 3 | Pending |
| MOD-03 | Phase 3 | Pending |
| WRK-01 | Phase 2 | Pending |
| WRK-02 | Phase 2 | Pending |
| RPT-01 | Phase 4 | Pending |
| RPT-02 | Phase 4 | Pending |
| RPT-03 | Phase 4 | Pending |
| RPT-04 | Phase 4 | Pending |
| RPT-05 | Phase 4 | Pending |
| INF-01 | Phase 1 | Pending |
| INF-02 | Phase 1 | Pending |
| INF-03 | Phase 1 | Pending |
| INF-04 | Phase 1 | Pending |

**Coverage:**
- v1 requirements: 28 total
- Mapped to phases: 28 ✓
- Unmapped: 0

---
*Requirements defined: 2026-02-11*
*Last updated: 2026-02-11 after roadmap creation*
