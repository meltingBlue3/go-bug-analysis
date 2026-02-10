package analysis

import (
	"go-bug-analysis/internal/csvparse"
	"sort"
	"strings"
)

// ComputeWorkload computes per-assignee bug counts for active and total bugs.
func ComputeWorkload(bugs []csvparse.Bug) *WorkloadData {
	activeCount := make(map[string]int)
	totalCount := make(map[string]int)

	for _, b := range bugs {
		assignee := strings.TrimSpace(b.Assignee)
		if assignee == "" {
			assignee = "未指派"
		}

		totalCount[assignee]++

		if b.Status == "激活" {
			activeCount[assignee]++
		}
	}

	// Convert maps to sorted slices
	byActive := mapToSortedSlice(activeCount)
	byTotal := mapToSortedSlice(totalCount)

	return &WorkloadData{
		ByActive: byActive,
		ByTotal:  byTotal,
	}
}

// mapToSortedSlice converts a name→count map to a slice sorted descending by count.
func mapToSortedSlice(m map[string]int) []AssigneeStats {
	result := make([]AssigneeStats, 0, len(m))
	for name, count := range m {
		if count > 0 {
			result = append(result, AssigneeStats{Name: name, Count: count})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Name < result[j].Name
	})
	return result
}
