---
phase: 02-core-kpi-bottleneck-analysis
plan: 01
subsystem: analysis-kpi-severity
tags: [go, analysis, kpi, severity, echarts, dashboard]

# Dependency graph
requires:
  - phase: 01-02
    provides: CSV parser, Bug struct, upload endpoint, AppState, frontend shell
provides:
  - Analysis package with date utils, KPI computation, severity distribution
  - GET /api/analysis endpoint returning KPI + severity JSON
  - Dashboard rendering module with KPI cards and ECharts severity chart
  - Toggle controls for pie/bar chart type and all/new-only data scope
  - Re-import navigation from dashboard back to upload view
affects: [02-02-backlog-age-workload, 02-03-module-trend, 03-01-advanced-charts]

# Tech tracking
tech-stack:
  added: [time.ParseInLocation, echarts-pie, echarts-bar, echarts-roseType]
  patterns: [analysis-package-pattern, dateutil-pattern, iife-dashboard-module, toggle-control-pattern]

key-files:
  created:
    - internal/analysis/types.go
    - internal/analysis/dateutil.go
    - internal/analysis/kpi.go
    - internal/analysis/severity.go
    - internal/analysis/analyze.go
    - web/static/js/dashboard.js
  modified:
    - internal/server/server.go
    - web/static/index.html
    - web/static/js/app.js
    - web/static/css/style.css

key-decisions:
  - "Date comparison via DateOnly() truncation — avoids time-of-day matching issues"
  - "ComputeKPI iterates bugs once for all KPI counts — O(n) single pass"
  - "Severity always returns all 4 levels (even if count=0) — charts render consistently"
  - "Unknown severity values mapped to level '0' (未设置) — graceful handling"
  - "ECharts roseType for pie chart — better visual differentiation of severity levels"
  - "Toggle controls use event delegation on parent toggle-group — clean event handling"

patterns-established:
  - "Analysis package pattern: types.go + domain-specific computation files + analyze.go entry point"
  - "Date utility pattern: ParseDate tries datetime then date-only format with local timezone"
  - "Dashboard module: IIFE extending window.BugAnalysis with renderDashboard entry point"
  - "Chart toggle: toggle-group with data-value buttons, updateChart() re-renders with setOption(opt, true)"
  - "KPI card pattern: kpi-card with kpi-label, kpi-value, kpi-sub; danger/warning CSS classes"

# Metrics
duration: 5min
completed: 2026-02-11
---

# Phase 02 Plan 01: KPI Cards & Severity Distribution Summary

**Analysis package with KPI computation, severity distribution, ECharts charts with toggle controls, and 5 KPI dashboard cards**

## Performance

- **Duration:** ~5 min
- **Started:** 2026-02-11
- **Completed:** 2026-02-11
- **Tasks:** 2
- **Files created:** 6, modified: 4

## Accomplishments
- Analysis package (`internal/analysis/`) with clean separation: types, date utilities, KPI computation, severity distribution, and top-level Analyze entry point
- `ParseDate` handles Zentao's two date formats ("2006-01-02 15:04:05" and "2006-01-02") with graceful handling of empty/"0000-00-00" strings
- `ComputeKPI` computes 7 KPIs in a single O(n) pass: todayNew, yesterdayNew, todayFixed, yesterdayFixed, total, active, pendingVerify
- `ComputeSeverity` produces severity distributions for all bugs and today's new bugs, always including all 4 severity levels
- GET /api/analysis endpoint returns combined KPI + severity JSON, with 400 error if no CSV uploaded
- 5 KPI cards: 今日新增 (with yesterday subtitle), 今日修复 (with yesterday subtitle), Bug 总数, 激活 Bug (red if > 0), 待验证 (warning if > 0)
- Severity chart renders as ECharts pie (roseType) by default, with pie/bar toggle and all/new-only data scope toggle
- Dashboard.js follows IIFE pattern, extends `window.BugAnalysis.renderDashboard`
- Re-import button in dashboard header returns to upload view and clears file input state

## Task Commits

Each task was committed atomically:

1. **Task 1: Analysis package — KPI, severity, API endpoint** - `5fdde92` (feat)
2. **Task 2: Frontend — KPI cards and severity chart** - `759a88e` (feat)

## Files Created/Modified
- `internal/analysis/types.go` - AnalysisResult, KPIData, SeverityItem, SeverityData structs
- `internal/analysis/dateutil.go` - ParseDate, DateOnly, Today, Yesterday date utilities
- `internal/analysis/kpi.go` - ComputeKPI with single-pass iteration over bugs
- `internal/analysis/severity.go` - ComputeSeverity, countSeverity, filterTodayNew helpers
- `internal/analysis/analyze.go` - Analyze() entry point composing KPI + severity
- `internal/server/server.go` - Added handleAnalysis handler, GET /api/analysis route
- `web/static/js/dashboard.js` - IIFE module: renderDashboard, renderKPI, renderSeverityChart, toggle controls
- `web/static/index.html` - 5 KPI cards with IDs, severity chart with toggle buttons, dashboard header, dashboard.js script
- `web/static/js/app.js` - renderDashboard() call after upload, initReimport() for re-import button
- `web/static/css/style.css` - KPI sub, warning/danger, dashboard header, reimport button, chart header, toggle group/button styles

## Decisions Made
- Date comparison uses `DateOnly()` truncation to midnight — avoids issues with time-of-day differences
- `ComputeKPI` uses a single O(n) loop for all 7 KPIs — efficient for large datasets
- Severity distribution always returns all 4 levels even with zero counts — ensures charts render consistently
- Unknown/missing severity values grouped under "未设置" (level "0") — graceful degradation
- ECharts `roseType: 'radius'` for pie chart — better visual differentiation between severity levels
- Toggle controls use click event delegation on parent toggle-group element — clean, minimal event listeners

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None — all code compiled cleanly, `go build ./...` and `go vet ./...` passed without issues.

## Next Plan Readiness
- Analysis package pattern established — 02-02 and 02-03 can add new analysis modules following the same structure
- Dashboard module pattern established — additional chart modules can extend `window.BugAnalysis`
- Toggle control pattern reusable for additional chart controls in later plans
- GET /api/analysis endpoint can be extended to include additional analysis results

---
*Phase: 02-core-kpi-bottleneck-analysis*
*Completed: 2026-02-11*
