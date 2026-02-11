package analysis

import (
	"go-bug-analysis/internal/csvparse"
	"sort"
	"strings"
	"time"
)

// ComputeModule orchestrates module-level analysis: stats, heatmap, and trend.
func ComputeModule(bugs []csvparse.Bug) *ModuleData {
	stats := computeModuleStats(bugs)
	heatmap := computeHeatmap(bugs, stats)
	trend := computeTrend(bugs, stats)
	return &ModuleData{
		Stats:   stats,
		Heatmap: heatmap,
		Trend:   trend,
	}
}

// computeModuleStats counts total and active bugs per module.
func computeModuleStats(bugs []csvparse.Bug) []ModuleStats {
	type counts struct {
		total  int
		active int
	}
	m := make(map[string]*counts)

	for _, b := range bugs {
		mod := strings.TrimSpace(b.Module)
		if mod == "" {
			mod = "未分类"
		}
		c, ok := m[mod]
		if !ok {
			c = &counts{}
			m[mod] = c
		}
		c.total++
		if b.Status == "激活" {
			c.active++
		}
	}

	stats := make([]ModuleStats, 0, len(m))
	for name, c := range m {
		rate := 0.0
		if c.total > 0 {
			rate = roundTo(float64(c.active)/float64(c.total)*100, 1)
		}
		stats = append(stats, ModuleStats{
			Name:       name,
			Total:      c.total,
			Active:     c.active,
			ActiveRate: rate,
		})
	}

	// Sort descending by Total; ties broken alphabetically
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Total != stats[j].Total {
			return stats[i].Total > stats[j].Total
		}
		return stats[i].Name < stats[j].Name
	})

	return stats
}

// severityIndex maps severity level string to column index.
// "1"→0 (致命), "2"→1 (严重), "3"→2 (一般), "4"→3 (轻微).
var severityIndex = map[string]int{
	"1": 0,
	"2": 1,
	"3": 2,
	"4": 3,
}

var severityLabels = []string{"致命", "严重", "一般", "轻微"}

// computeHeatmap builds the Module × Severity heatmap matrix using top 15 modules.
func computeHeatmap(bugs []csvparse.Bug, stats []ModuleStats) *HeatmapData {
	maxModules := 15
	if len(stats) < maxModules {
		maxModules = len(stats)
	}

	topModules := make([]string, maxModules)
	moduleIdx := make(map[string]int, maxModules)
	for i := 0; i < maxModules; i++ {
		topModules[i] = stats[i].Name
		moduleIdx[stats[i].Name] = i
	}

	// Initialize matrix [module][severity] = 0
	data := make([][]int, maxModules)
	for i := range data {
		data[i] = make([]int, 4) // 4 severity levels
	}

	maxValue := 0
	for _, b := range bugs {
		mod := strings.TrimSpace(b.Module)
		if mod == "" {
			mod = "未分类"
		}
		mIdx, ok := moduleIdx[mod]
		if !ok {
			continue
		}
		sIdx, ok := severityIndex[strings.TrimSpace(b.Severity)]
		if !ok {
			continue
		}
		data[mIdx][sIdx]++
		if data[mIdx][sIdx] > maxValue {
			maxValue = data[mIdx][sIdx]
		}
	}

	return &HeatmapData{
		Modules:    topModules,
		Severities: severityLabels,
		Data:       data,
		MaxValue:   maxValue,
	}
}

// computeTrend computes daily bug creation counts per module for the last 30 days.
func computeTrend(bugs []csvparse.Bug, stats []ModuleStats) *TrendData {
	today := Today()
	const totalDays = 30

	// Generate 30 date slots (today-29 to today)
	dateSlots := make([]time.Time, totalDays)
	dateLabels := make([]string, totalDays)
	for i := 0; i < totalDays; i++ {
		d := today.AddDate(0, 0, i-(totalDays-1))
		dateSlots[i] = d
		dateLabels[i] = d.Format("01-02")
	}

	rangeStart := dateSlots[0]

	// Pick top 10 modules for trend lines
	maxTrend := 10
	if len(stats) < maxTrend {
		maxTrend = len(stats)
	}
	topNames := make(map[string]int, maxTrend) // name → index in series
	for i := 0; i < maxTrend; i++ {
		topNames[stats[i].Name] = i
	}

	// Initialize counts: [moduleSeriesIdx][dateIdx]
	counts := make([][]int, maxTrend)
	for i := range counts {
		counts[i] = make([]int, totalDays)
	}

	for _, b := range bugs {
		mod := strings.TrimSpace(b.Module)
		if mod == "" {
			mod = "未分类"
		}
		seriesIdx, ok := topNames[mod]
		if !ok {
			continue
		}
		created, ok := ParseDate(b.CreatedDate)
		if !ok {
			continue
		}
		createdDay := DateOnly(created)
		if createdDay.Before(rangeStart) || createdDay.After(today) {
			continue
		}
		// Calculate date index
		dayOffset := int(createdDay.Sub(rangeStart).Hours() / 24)
		if dayOffset >= 0 && dayOffset < totalDays {
			counts[seriesIdx][dayOffset]++
		}
	}

	// Build series
	series := make([]TrendSeries, maxTrend)
	for i := 0; i < maxTrend; i++ {
		series[i] = TrendSeries{
			Name:   stats[i].Name,
			Counts: counts[i],
		}
	}

	return &TrendData{
		Dates:  dateLabels,
		Series: series,
		Days7:  totalDays - 7, // 23: index where last-7-day window starts
	}
}
