package analysis

import (
	"go-bug-analysis/internal/csvparse"
)

// Analyze runs all analysis computations on the parsed bug data.
func Analyze(bugs []csvparse.Bug) *AnalysisResult {
	result := &AnalysisResult{
		KPI:      ComputeKPI(bugs),
		Severity: ComputeSeverity(bugs),
		Age:      ComputeAge(bugs),
		Workload: ComputeWorkload(bugs),
		Module:   ComputeModule(bugs),
	}
	result.Report = ComputeReport(result)
	return result
}
