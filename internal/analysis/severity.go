package analysis

import (
	"go-bug-analysis/internal/csvparse"
)

// severityLevels defines the ordered severity levels with their labels.
var severityLevels = []struct {
	Level string
	Label string
}{
	{"1", "致命"},
	{"2", "严重"},
	{"3", "一般"},
	{"4", "轻微"},
}

// severityLabelMap maps severity level strings to Chinese labels.
var severityLabelMap = map[string]string{
	"1": "致命",
	"2": "严重",
	"3": "一般",
	"4": "轻微",
}

// ComputeSeverity computes the severity distribution for all bugs and today's new bugs.
func ComputeSeverity(bugs []csvparse.Bug) *SeverityData {
	todayNew := filterTodayNew(bugs)

	return &SeverityData{
		All:     countSeverity(bugs),
		NewOnly: countSeverity(todayNew),
	}
}

// countSeverity counts bugs per severity level, always returning all 4 levels.
func countSeverity(bugs []csvparse.Bug) []SeverityItem {
	counts := make(map[string]int)
	for i := range bugs {
		level := bugs[i].Severity
		if _, ok := severityLabelMap[level]; !ok {
			level = "0" // unknown
		}
		counts[level]++
	}

	items := make([]SeverityItem, 0, 5)
	for _, sl := range severityLevels {
		items = append(items, SeverityItem{
			Level: sl.Level,
			Label: sl.Label,
			Count: counts[sl.Level],
		})
	}

	// Add unknown/unset if any bugs had unrecognized severity
	if counts["0"] > 0 {
		items = append(items, SeverityItem{
			Level: "0",
			Label: "未设置",
			Count: counts["0"],
		})
	}

	return items
}

// filterTodayNew returns bugs whose CreatedDate matches today.
func filterTodayNew(bugs []csvparse.Bug) []csvparse.Bug {
	today := Today()
	var result []csvparse.Bug
	for i := range bugs {
		if ct, ok := ParseDate(bugs[i].CreatedDate); ok {
			if DateOnly(ct).Equal(today) {
				result = append(result, bugs[i])
			}
		}
	}
	return result
}
