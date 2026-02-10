package analysis

import (
	"go-bug-analysis/internal/csvparse"
)

// ComputeKPI iterates bugs once to compute all KPI counts.
func ComputeKPI(bugs []csvparse.Bug) *KPIData {
	today := Today()
	yesterday := Yesterday()

	kpi := &KPIData{
		Total: len(bugs),
	}

	for i := range bugs {
		b := &bugs[i]

		// Today / yesterday new bugs (by CreatedDate)
		if ct, ok := ParseDate(b.CreatedDate); ok {
			cd := DateOnly(ct)
			if cd.Equal(today) {
				kpi.TodayNew++
			} else if cd.Equal(yesterday) {
				kpi.YesterdayNew++
			}
		}

		// Today / yesterday fixed bugs (by ResolvedDate)
		if rt, ok := ParseDate(b.ResolvedDate); ok {
			rd := DateOnly(rt)
			if rd.Equal(today) {
				kpi.TodayFixed++
			} else if rd.Equal(yesterday) {
				kpi.YesterdayFixed++
			}
		}

		// Status-based counts
		switch b.Status {
		case "激活":
			kpi.Active++
		case "已解决":
			kpi.PendingVerify++
		}
	}

	return kpi
}
