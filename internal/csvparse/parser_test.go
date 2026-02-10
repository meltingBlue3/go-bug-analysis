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

func TestParseEmptyFile(t *testing.T) {
	_, err := Parse(strings.NewReader(""))
	if err == nil {
		t.Fatal("Expected error for empty file")
	}
}
