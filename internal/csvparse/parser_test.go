package csvparse

import (
	"os"
	"strings"
	"testing"
)

func TestParseGBKFile(t *testing.T) {
	f, err := os.Open("../../2225.csv")
	if err != nil {
		t.Skipf("Skipping: sample CSV not found: %v", err)
	}
	defer f.Close()

	result, err := Parse(f)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// 2225.csv has ~2225 bug records (multiline fields inflate line count)
	if result.TotalRows < 2000 {
		t.Errorf("Expected 2000+ bugs, got %d", result.TotalRows)
	}
	t.Logf("Total bugs parsed: %d", result.TotalRows)
	t.Logf("Columns recognized: %v", result.Columns)
	t.Logf("Warnings count: %d", len(result.Warnings))

	// First bug should have non-empty key fields
	if len(result.Bugs) > 0 {
		first := result.Bugs[0]
		t.Logf("First bug: ID=%s, Title=%s, Status=%s, Creator=%s", first.ID, first.Title, first.Status, first.Creator)
		if first.ID == "" {
			t.Error("First bug ID is empty")
		}
		if first.Title == "" {
			t.Error("First bug Title is empty")
		}
		if first.Status == "" {
			t.Error("First bug Status is empty")
		}
	}
}

func TestParseUTF8(t *testing.T) {
	csv := "Bug编号,Bug标题,严重程度,Bug状态,由谁创建,创建日期,指派给\n1001,测试Bug,3,激活,张三,2025-01-01,李四\n"
	result, err := Parse(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if result.TotalRows != 1 {
		t.Errorf("Expected 1 bug, got %d", result.TotalRows)
	}
	bug := result.Bugs[0]
	if bug.ID != "1001" || bug.Title != "测试Bug" || bug.Severity != "3" {
		t.Errorf("Unexpected bug data: %+v", bug)
	}
}

func TestParseMissingRequiredColumns(t *testing.T) {
	csv := "Bug编号,Bug标题\n1001,测试\n"
	_, err := Parse(strings.NewReader(csv))
	if err == nil {
		t.Fatal("Expected error for missing required columns")
	}
	if !strings.Contains(err.Error(), "缺少必填列") {
		t.Errorf("Error message should mention missing columns, got: %s", err.Error())
	}
	t.Logf("Got expected error: %s", err.Error())
}

func TestParseAssigneeClosedFallbackToResolver(t *testing.T) {
	csv := "Bug编号,Bug标题,严重程度,Bug状态,由谁创建,创建日期,指派给,解决者\n" +
		"2001,Bug一,3,已解决,张三,2025-06-01,Closed,王五\n" +
		"2002,Bug二,2,激活,张三,2025-06-02,李四,赵六\n"

	result, err := Parse(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if result.TotalRows != 2 {
		t.Fatalf("Expected 2 bugs, got %d", result.TotalRows)
	}

	// Row 1: Assignee was "Closed", Resolver is "王五" → Assignee should become "王五"
	if result.Bugs[0].Assignee != "王五" {
		t.Errorf("Bug 2001: expected Assignee=王五, got %q", result.Bugs[0].Assignee)
	}

	// Row 2: Assignee is "李四" (not "Closed") → should stay "李四", no fallback
	if result.Bugs[1].Assignee != "李四" {
		t.Errorf("Bug 2002: expected Assignee=李四, got %q", result.Bugs[1].Assignee)
	}
}

func TestParseEmptyFile(t *testing.T) {
	_, err := Parse(strings.NewReader(""))
	if err == nil {
		t.Fatal("Expected error for empty file")
	}
}
