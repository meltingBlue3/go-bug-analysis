---
phase: quick-1
plan: 1
subsystem: csvparse
tags: [bugfix, parser, zentao]
dependency-graph:
  requires: []
  provides: [assignee-closed-fallback]
  affects: [internal/csvparse]
tech-stack:
  added: []
  patterns: [post-processing-after-rowToBug]
key-files:
  created: []
  modified:
    - internal/csvparse/parser.go
    - internal/csvparse/parser_test.go
decisions:
  - Fallback placed in Parse() loop after rowToBug(), not inside rowToBug() itself
metrics:
  duration: ~2 min
  completed: 2026-02-11
---

# Quick Task 1: Fix Assignee "Closed" Fallback to Resolver Summary

**One-liner:** When Zentao CSV exports set Assignee to "Closed" for resolved bugs, fall back to the Resolver column value.

## What Was Done

### Task 1: Add Assignee fallback logic in Parse()

Added post-processing in the `Parse()` loop immediately after `rowToBug()` call. When `bug.Assignee == "Closed"` and `bug.Resolver != ""`, the Assignee is replaced with the Resolver value. This handles Zentao's behavior of setting the Assignee column to "Closed" for resolved bugs.

**Commit:** `4fb7c2d` — `fix(quick-1): fallback Assignee from Closed to Resolver in Parse()`

### Task 2: Add test for Closed→Resolver fallback

Added `TestParseAssigneeClosedFallbackToResolver` in `parser_test.go` with two cases:
1. **Fallback case:** Assignee="Closed", Resolver="王五" → asserts Assignee becomes "王五"
2. **No-fallback case:** Assignee="李四", Resolver="赵六" → asserts Assignee stays "李四"

**Commit:** `09ac01d` — `test(quick-1): add TestParseAssigneeClosedFallbackToResolver`

## Verification

- `go build ./...` — passes
- `go test ./internal/csvparse/...` — all tests pass (including new fallback test)

## Deviations from Plan

None — plan executed exactly as written.

## Self-Check: PASSED

- [x] `internal/csvparse/parser.go` modified with fallback logic
- [x] `internal/csvparse/parser_test.go` modified with new test
- [x] Commit `4fb7c2d` exists
- [x] Commit `09ac01d` exists
- [x] All tests pass
