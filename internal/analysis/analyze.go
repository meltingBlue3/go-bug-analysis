package analysis

import (
	"go-bug-analysis/internal/csvparse"
)

// Analyze runs all analysis computations on the parsed bug data.
func Analyze(bugs []csvparse.Bug) *AnalysisResult {
	return &AnalysisResult{
		KPI:      ComputeKPI(bugs),
		Severity: ComputeSeverity(bugs),
		Age:      ComputeAge(bugs),
	}
}
