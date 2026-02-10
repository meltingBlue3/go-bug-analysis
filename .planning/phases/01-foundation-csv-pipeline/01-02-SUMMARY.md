---
phase: 01-foundation-csv-pipeline
plan: 02
subsystem: csv-pipeline
tags: [go, csv, gbk, encoding, upload, multipart, x-text]

# Dependency graph
requires:
  - phase: 01-01
    provides: Go scaffold with embed.FS, HTTP server, frontend shell with BugAnalysis API
provides:
  - CSV parser with automatic GBK/UTF-8 encoding detection
  - Bug struct with 22 fields mapped from Zentao Chinese headers
  - POST /api/upload endpoint with multipart file handling
  - GET /api/data endpoint returning full parsed data
  - Frontend upload flow with progress, success, and error feedback
  - Thread-safe AppState for in-memory data storage
affects: [02-01-severity-trend, 02-02-module-assignee, 03-01-advanced-charts]

# Tech tracking
tech-stack:
  added: [golang.org/x/text, encoding/csv, encoding/json, sync.RWMutex]
  patterns: [gbk-detection, header-mapping, appstate-pattern, json-api-response]

key-files:
  created:
    - internal/csvparse/types.go
    - internal/csvparse/parser.go
    - internal/csvparse/parser_test.go
  modified:
    - internal/server/server.go
    - main.go
    - web/static/js/app.js
    - go.mod
    - go.sum

key-decisions:
  - "All Bug fields as string type — defers date parsing to analysis modules"
  - "One-shot io.ReadAll for encoding detection — acceptable for 1.7MB files"
  - "LazyQuotes + FieldsPerRecord=-1 for Zentao CSV quirks"
  - "AppState with sync.RWMutex for thread-safe result storage"
  - "JSON API with Chinese error messages for consistent UX"

patterns-established:
  - "Encoding detection: UTF-8 BOM → utf8.Valid → assume GBK fallback"
  - "Header mapping: Chinese header → field name via HeaderMap lookup table"
  - "API response pattern: writeJSON helper with SetEscapeHTML(false)"
  - "Upload flow: FormData fetch → JSON response → showStatus/showDashboard"

# Metrics
duration: 7min
completed: 2026-02-11
---

# Phase 01 Plan 02: CSV Upload & Parsing Pipeline Summary

**GBK/UTF-8 auto-detecting CSV parser with Zentao header mapping, multipart upload endpoint, and frontend upload UX with Chinese error messages**

## Performance

- **Duration:** ~7 min
- **Started:** 2026-02-10T20:58:38Z
- **Completed:** 2026-02-10T21:06:00Z
- **Tasks:** 2
- **Files created:** 3, modified: 5

## Accomplishments
- CSV parser correctly detects GBK encoding and converts to UTF-8, parsing 2225 bug records from the sample Zentao export
- All 22 target columns recognized and mapped to Bug struct fields; unrecognized columns silently ignored
- POST /api/upload handles multipart file upload, validates .csv extension, returns JSON summary with bug count, columns, sample bug
- GET /api/data returns full parsed data (~1.5MB JSON) for Phase 2 analysis modules
- Frontend upload flow: file selection → "正在解析..." status → success summary with sample bug → auto-switch to dashboard
- Required column validation returns Chinese error messages naming the specific missing columns
- Test suite: 4 tests covering GBK file parsing, UTF-8 parsing, missing required columns, and empty file

## Task Commits

Each task was committed atomically:

1. **Task 1: CSV data model and parser** - `2243ba0` (feat)
2. **Task 2: Upload endpoint and frontend integration** - `105e980` (feat)

## Files Created/Modified
- `internal/csvparse/types.go` - Bug struct (22 fields), HeaderMap, RequiredHeaders, ParseResult
- `internal/csvparse/parser.go` - Parse(), detectAndConvert(), mapHeaders(), validateHeaders(), rowToBug()
- `internal/csvparse/parser_test.go` - 4 tests: GBK file, UTF-8, missing columns, empty file
- `internal/server/server.go` - AppState, handleUpload, handleData, writeJSON, New() updated signature
- `main.go` - Creates AppState, passes to server.New()
- `web/static/js/app.js` - uploadCSV() with fetch, setUploading(), showStatusHTML(), error handling
- `go.mod` - Added golang.org/x/text dependency
- `go.sum` - Module checksums

## Decisions Made
- All Bug fields kept as `string` — avoids premature date parsing; analysis modules can parse dates as needed
- Single io.ReadAll for encoding detection — memory-safe for expected file sizes (< 100MB limit enforced)
- LazyQuotes=true and FieldsPerRecord=-1 for Zentao's non-standard CSV formatting (unbalanced quotes, variable column counts)
- AppState uses sync.RWMutex rather than channels — simpler for read-heavy access pattern
- JSON encoder with SetEscapeHTML(false) — preserves Chinese characters in JSON output without HTML entity encoding
- File contains 2225 bug records (not 50K) — the 50K line count reflects multiline "重现步骤" fields within CSV records

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered
- CSV file has 2225 bug records, not ~51000 as initially estimated. The 50624 line count in the raw file comes from multiline "重现步骤" (reproduction steps) fields. Go's `csv.ReadAll()` handles this correctly via quoted field parsing.

## User Setup Required
None — no external service configuration required.

## Next Phase Readiness
- CSV data pipeline complete: upload → parse → store in memory → serve via API
- ParseResult with Bug slice ready for Phase 2 analysis (severity distribution, module analysis, trend charts)
- GET /api/data endpoint provides full data access for analysis endpoints
- Frontend BugAnalysis.data stores summary for chart modules to use
- ECharts library already loaded (from Plan 01) — ready for chart rendering

---
*Phase: 01-foundation-csv-pipeline*
*Completed: 2026-02-11*
