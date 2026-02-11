package csvparse

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// Parse reads a CSV stream (GBK or UTF-8), detects encoding, maps Zentao
// Chinese headers to Bug struct fields, validates required columns, and
// returns the parsed result. The entire content is read into memory for
// encoding detection (~1.7MB for 50k rows is acceptable).
func Parse(r io.Reader) (*ParseResult, error) {
	// 1. Read all bytes for encoding detection
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("文件内容为空")
	}

	// 2. Detect and convert encoding to UTF-8
	utf8Data, err := detectAndConvert(data)
	if err != nil {
		return nil, err
	}

	// 3. Parse CSV
	csvReader := csv.NewReader(strings.NewReader(string(utf8Data)))
	csvReader.LazyQuotes = true
	csvReader.FieldsPerRecord = -1 // Allow variable column counts

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CSV 格式错误: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV 文件至少需要包含表头行和一行数据")
	}

	// 4. Map headers: column index → field name
	headerMapping, columns := mapHeaders(records[0])

	// 5. Validate required columns exist
	if err := validateHeaders(headerMapping); err != nil {
		return nil, err
	}

	// 6. Parse data rows into Bug structs
	dataRows := records[1:]
	bugs := make([]Bug, 0, len(dataRows))
	var warnings []string

	for i, row := range dataRows {
		// Skip completely empty rows
		if isEmptyRow(row) {
			continue
		}

		bug := rowToBug(row, headerMapping)

		// 禅道导出 CSV 中，已解决的 Bug 的"指派给"列可能变为 "Closed"，此时用"解决者"列填充
		if bug.Assignee == "Closed" && bug.Resolver != "" {
			bug.Assignee = bug.Resolver
		}

		// Warn about missing required fields on specific rows
		if bug.ID == "" {
			warnings = append(warnings, fmt.Sprintf("第 %d 行: Bug编号 字段为空", i+2))
		}

		bugs = append(bugs, bug)
	}

	return &ParseResult{
		Bugs:      bugs,
		TotalRows: len(bugs),
		Warnings:  warnings,
		Columns:   columns,
	}, nil
}

// detectAndConvert detects the encoding of raw bytes and converts to UTF-8.
// Detection order: UTF-8 BOM → valid UTF-8 → assume GBK.
func detectAndConvert(data []byte) ([]byte, error) {
	// Check for UTF-8 BOM (0xEF 0xBB 0xBF)
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return data[3:], nil // Strip BOM and return
	}

	// Check if data is valid UTF-8
	if utf8.Valid(data) {
		return data, nil
	}

	// Assume GBK and attempt conversion
	utf8Data, err := simplifiedchinese.GBK.NewDecoder().Bytes(data)
	if err != nil {
		return nil, fmt.Errorf("文件编码无法识别，请确认文件为 GBK 或 UTF-8 编码")
	}

	return utf8Data, nil
}

// mapHeaders maps the CSV header row to Bug field names using HeaderMap.
// Returns: columnIndex→fieldName mapping, and the list of recognized column names.
func mapHeaders(row []string) (map[int]string, []string) {
	mapping := make(map[int]string)
	var columns []string

	for i, cell := range row {
		// Clean: strip BOM remnants, quotes, and whitespace
		cell = strings.TrimSpace(cell)
		cell = strings.Trim(cell, "\ufeff")
		cell = strings.Trim(cell, "\"")
		cell = strings.TrimSpace(cell)

		if fieldName, ok := HeaderMap[cell]; ok {
			mapping[i] = fieldName
			columns = append(columns, cell)
		}
	}

	return mapping, columns
}

// validateHeaders checks that all required columns are present in the mapping.
func validateHeaders(mapping map[int]string) error {
	// Build a set of mapped field names
	fieldSet := make(map[string]bool)
	for _, fieldName := range mapping {
		fieldSet[fieldName] = true
	}

	var missing []string
	for _, header := range RequiredHeaders {
		fieldName := HeaderMap[header]
		if !fieldSet[fieldName] {
			missing = append(missing, header)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("CSV 缺少必填列：%s", strings.Join(missing, ", "))
	}
	return nil
}

// rowToBug converts a CSV row to a Bug struct using the header mapping.
func rowToBug(row []string, mapping map[int]string) Bug {
	var bug Bug

	for colIdx, fieldName := range mapping {
		// Bounds check — Zentao exports occasionally have short rows
		if colIdx >= len(row) {
			continue
		}
		value := strings.TrimSpace(row[colIdx])

		switch fieldName {
		case "ID":
			bug.ID = value
		case "Product":
			bug.Product = value
		case "Module":
			bug.Module = value
		case "Title":
			bug.Title = value
		case "Severity":
			bug.Severity = value
		case "Priority":
			bug.Priority = value
		case "BugType":
			bug.BugType = value
		case "Status":
			bug.Status = value
		case "Creator":
			bug.Creator = value
		case "CreatedDate":
			bug.CreatedDate = value
		case "Assignee":
			bug.Assignee = value
		case "AssignedDate":
			bug.AssignedDate = value
		case "Resolver":
			bug.Resolver = value
		case "Resolution":
			bug.Resolution = value
		case "ResolvedDate":
			bug.ResolvedDate = value
		case "Closer":
			bug.Closer = value
		case "ClosedDate":
			bug.ClosedDate = value
		case "ActivationCount":
			bug.ActivationCount = value
		case "Deadline":
			bug.Deadline = value
		case "AffectedVersion":
			bug.AffectedVersion = value
		case "ResolvedVersion":
			bug.ResolvedVersion = value
		case "Keywords":
			bug.Keywords = value
		}
	}

	return bug
}

// isEmptyRow checks if a CSV row has no non-empty values.
func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}
