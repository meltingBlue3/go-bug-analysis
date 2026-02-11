---
phase: 1-fix-assignee-closed-fallback-to-resolver
plan: 1
type: execute
wave: 1
depends_on: []
files_modified: [internal/csvparse/parser.go, internal/csvparse/parser_test.go]
autonomous: true

must_haves:
  truths:
    - "When Assignee is 'Closed' and Resolver is set, Assignee shows Resolver value"
  artifacts:
    - path: internal/csvparse/parser.go
      provides: "Assignee fallback logic in Parse()"
      pattern: "bug.Assignee == \"Closed\""
    - path: internal/csvparse/parser_test.go
      provides: "Test for Closed→Resolver fallback"
      pattern: "TestParseAssigneeClosedFallbackToResolver"
  key_links:
    - from: internal/csvparse/parser.go
      to: "Parse() loop"
      via: "post-processing after rowToBug"
      pattern: "rowToBug"
---

<objective>
Add fallback logic: when Zentao sets Assignee to "Closed" for resolved bugs, use the Resolver column value instead.

Purpose: Zentao CSV exports sometimes put "Closed" in the 指派给 (Assignee) column when bugs are resolved. The 解决者 (Resolver) column holds the actual assignee. Use Resolver as fallback for downstream analysis.

Output: Updated parser.go with fallback logic; new test in parser_test.go.
</objective>

<context>
@internal/csvparse/parser.go
@internal/csvparse/types.go
@internal/csvparse/parser_test.go
</context>

<tasks>

<task type="auto">
  <name>Task 1: Add Assignee fallback logic in Parse()</name>
  <files>internal/csvparse/parser.go</files>
  <action>
    In Parse(), immediately after line 65 (`bug := rowToBug(row, headerMapping)`), add post-processing:

    if bug.Assignee == "Closed" && bug.Resolver != "" {
        bug.Assignee = bug.Resolver
    }

    Place this before the ID empty check (before line 67). Do not modify rowToBug.
  </action>
  <verify>go build ./internal/csvparse/...</verify>
  <done>Assignee "Closed" with non-empty Resolver is replaced by Resolver when bugs are appended.</done>
</task>

<task type="auto">
  <name>Task 2: Add test for Closed→Resolver fallback</name>
  <files>internal/csvparse/parser_test.go</files>
  <action>
    Add TestParseAssigneeClosedFallbackToResolver that:
    1. Uses a minimal CSV with required columns plus 解决者 (Resolver)
    2. One row has 指派给=Closed and 解决者=王五
    3. Parse and assert result.Bugs[0].Assignee == "王五"
    4. Optionally add a second row: 指派给=李四, 解决者=赵六 — assert Assignee stays "李四" (no fallback)
  </action>
  <verify>go test ./internal/csvparse/... -run TestParseAssigneeClosedFallbackToResolver -v</verify>
  <done>Test passes; fallback and non-fallback cases both correct.</done>
</task>

</tasks>

<verification>
- go build ./...
- go test ./internal/csvparse/...
</verification>

<success_criteria>
- Assignee "Closed" with non-empty Resolver is replaced by Resolver
- All existing tests pass
- New test covers fallback and non-fallback cases
</success_criteria>
