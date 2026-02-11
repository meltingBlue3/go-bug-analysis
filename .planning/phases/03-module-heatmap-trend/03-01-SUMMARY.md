---
phase: 03-module-heatmap-trend
plan: 01
subsystem: analysis, ui
tags: [echarts, heatmap, trend, module-analysis, go]

# Dependency graph
requires:
  - phase: 02-core-kpi-bottleneck-analysis
    provides: "Analysis package pattern (types.go + domain files + analyze.go), dateutil, ECharts dashboard"
provides:
  - "Per-module bug statistics (total/active/activeRate) via ComputeModule"
  - "Module × Severity heatmap matrix for hotspot visualization"
  - "30-day module bug trend with 7d/30d toggle"
  - "Sortable module table with click-to-sort headers"
affects: [04-polish-export]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Section divider pattern for logical dashboard grouping"
    - "Click-to-sort table pattern with header indicator toggling"
    - "ECharts heatmap with reversed y-axis for top-down display"
    - "Trend line chart with date range slicing via days7 offset"

key-files:
  created:
    - "internal/analysis/module.go"
  modified:
    - "internal/analysis/types.go"
    - "internal/analysis/analyze.go"
    - "web/static/index.html"
    - "web/static/js/dashboard.js"
    - "web/static/css/style.css"

key-decisions:
  - "Top 15 modules for heatmap readability, top 10 for trend lines"
  - "Empty module field mapped to '未分类' for Chinese UI consistency"
  - "30-day trend window with days7=23 offset for 7-day slice"
  - "Reversed y-axis on heatmap so highest-count module appears at top"

patterns-established:
  - "Section divider: .section-divider + .section-title for visual grouping"
  - "Sortable table: data-sort attribute on th, sort-active/sort-desc classes"
  - "Heatmap pattern: flatten 2D array to [x, y, value] tuples for ECharts"

# Metrics
duration: 4min
completed: 2026-02-11
---

# Phase 3 Plan 1: Module Heatmap & Trend Summary

**Per-module bug statistics table with click-to-sort, Module × Severity ECharts heatmap with blue gradient, and 30-day trend line chart with 7d/30d toggle**

## Performance

- **Duration:** 4 min
- **Started:** 2026-02-11T12:54:44Z
- **Completed:** 2026-02-11T12:58:09Z
- **Tasks:** 2
- **Files modified:** 6

## Accomplishments
- Backend module analysis computing per-module total/active counts, severity heatmap matrix, and 30-day daily trend data
- Sortable module statistics table with active rate danger/warning styling
- Module × Severity heatmap visualizing bug concentration hotspots with blue gradient color scale
- Multi-line trend chart showing top 10 modules' daily bug creation with 7-day/30-day toggle

## Task Commits

Each task was committed atomically:

1. **Task 1: Backend — Module statistics, heatmap matrix, and trend computation** - `eba1f26` (feat)
2. **Task 2: Frontend — Module table, heatmap chart, and trend line chart** - `6bf8ecc` (feat)

## Files Created/Modified
- `internal/analysis/module.go` - ComputeModule with stats, heatmap, and trend sub-computations
- `internal/analysis/types.go` - ModuleData, ModuleStats, HeatmapData, TrendData, TrendSeries types
- `internal/analysis/analyze.go` - Wire ComputeModule into Analyze()
- `web/static/index.html` - Module section HTML (table, heatmap container, trend container with toggle)
- `web/static/js/dashboard.js` - renderModule, renderModuleTable, renderModuleHeatmap, renderModuleTrend, initTrendToggle
- `web/static/css/style.css` - Section divider, sortable table, rate styling, full-width chart card

## Decisions Made
- Top 15 modules for heatmap, top 10 for trend — balances readability and information density
- Empty module field → "未分类" consistent with workload "未指派" pattern
- Reversed y-axis on heatmap so top module (by total count) appears at the top
- days7 = 23 (30-7) as slice start index for the 7-day view
- Reused existing roundTo() from age.go for activeRate percentage rounding

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None — no external service configuration required.

## Next Phase Readiness
- Module analysis complete: MOD-01 (table), MOD-02 (heatmap), MOD-03 (trend) all satisfied
- All existing dashboard features (KPI, severity, age, workload) unaffected
- Ready for Phase 4 (Polish & Export) or additional Phase 3 plans

---
*Phase: 03-module-heatmap-trend*
*Completed: 2026-02-11*
