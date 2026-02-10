---
phase: 02-core-kpi-bottleneck-analysis
plan: 02
subsystem: age-fixtime-backlog
tags: [go, analysis, fix-time, backlog, echarts, dashboard]

# Dependency graph
requires:
  - phase: 02-01
    provides: Analysis package, dashboard module, ECharts patterns, date utils
provides:
  - Fix time computation (avg, P50, distribution)
  - Backlog age ranking for active bugs
  - Fix time stats cards and horizontal bar chart
  - Scrollable backlog table with severity badges and age color coding
affects: [02-03-module-trend, 03-01-advanced-charts]

# Tech tracking
tech-stack:
  added: [echarts-horizontal-bar, sort.Slice, math.Round]
  patterns: [age-analysis-pattern, stat-box-pattern, data-table-pattern, severity-badge-pattern]

key-files:
  created:
    - internal/analysis/age.go
  modified:
    - internal/analysis/types.go
    - internal/analysis/analyze.go
    - web/static/index.html
    - web/static/js/dashboard.js
    - web/static/css/style.css

key-decisions:
  - "Fix time only counts positive durations (resolved after creation) — skip edge cases"
  - "P50 computed as true median: middle element for odd, average of two middles for even"
  - "Distribution uses 4 buckets: 0-1d, 2-3d, 4-7d, 7+d — matches common triage thresholds"
  - "Backlog includes all active bugs even with unparseable dates (ageDays=0) — no data loss"
  - "roundTo helper for consistent decimal display (1 decimal place)"
  - "Horizontal bar chart for distribution — better readability for category comparison"

patterns-established:
  - "Age analysis pattern: ComputeAge orchestrates computeFixTime + computeBacklog"
  - "Stat box pattern: stat-box with stat-label, stat-value, stat-unit for metric display"
  - "Data table pattern: table-scroll for overflow, sticky headers, hover highlighting"
  - "Severity badge pattern: .severity-badge.s1-s4 with matching colors"
  - "Age color coding: age-danger (>30d red), age-warning (>14d orange)"

# Metrics
duration: 5min
completed: 2026-02-11
---

# Phase 02 Plan 02: Fix Time Stats & Backlog Age Ranking Summary

**Fix time computation with avg/P50 stats, distribution chart, and scrollable backlog age ranking table**

## Performance

- **Duration:** ~5 min
- **Started:** 2026-02-11
- **Completed:** 2026-02-11
- **Tasks:** 2
- **Files created:** 1, modified: 5

## Accomplishments
- `ComputeAge` function in `age.go` orchestrates fix time stats and backlog ranking
- `computeFixTime` calculates average and P50 (median) fix times from resolved bugs, with 4-bucket distribution (0-1d, 2-3d, 4-7d, 7+d)
- `computeBacklog` filters active (激活) bugs, computes age in days since creation, sorts descending by age
- Edge case handling: skips empty/unparseable dates, negative durations; includes active bugs with bad dates at ageDays=0
- Fix time stats section with two stat boxes (average + P50) with smart formatting (hours if <24h, days otherwise)
- Horizontal bar chart for fix time distribution using ECharts with green→red color gradient
- Scrollable backlog table with severity badges (color-coded s1-s4), age color coding (red >30d, orange >14d)
- All new types: AgeData, FixTimeStats, DistBucket, BacklogItem added to types.go
- AnalysisResult extended with Age field — API response now includes age analysis data

## Task Commits

Each task was committed atomically:

1. **Task 1: Backend — Fix time computation and backlog age ranking** - `6b97784` (feat)
2. **Task 2: Frontend — Fix time stats display, distribution chart, and backlog table** - `09a326f` (feat)

## Files Created/Modified
- `internal/analysis/age.go` - ComputeAge, computeFixTime, computeBacklog, roundTo helper
- `internal/analysis/types.go` - AgeData, FixTimeStats, DistBucket, BacklogItem structs; Age field on AnalysisResult
- `internal/analysis/analyze.go` - Added ComputeAge(bugs) call in Analyze()
- `web/static/index.html` - Fix time stats section (2 cards), backlog table with headers
- `web/static/js/dashboard.js` - renderAge, renderFixTimeStats (stat boxes + ECharts horizontal bar), renderBacklogTable (severity badges, age coloring)
- `web/static/css/style.css` - fix-time-summary grid, stat-box, data-table, table-scroll, severity-badge (s1-s4), age-danger/warning, backlog-count styles

## Decisions Made
- Fix time uses only positive durations (resolved > created) — avoids data entry errors
- P50 is true median: middle for odd count, average of two middles for even count
- Distribution buckets match common bug triage thresholds (1d, 3d, 7d boundaries)
- Active bugs with unparseable dates still appear in backlog (ageDays=0) — no silent data loss
- `roundTo` helper for consistent 1-decimal-place display
- Horizontal bar chart chosen over vertical for better label readability in distribution

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None — all code compiled cleanly, `go build ./...` and `go vet ./...` passed without issues.

## Next Plan Readiness
- Age analysis pattern established — can be extended with additional time-based metrics
- Data table pattern reusable for module-level bug lists in 02-03
- Severity badge and age color coding patterns reusable across dashboard
- API response structure extensible for remaining analysis modules

---
*Phase: 02-core-kpi-bottleneck-analysis*
*Completed: 2026-02-11*
