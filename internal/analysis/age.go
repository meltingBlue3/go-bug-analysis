package analysis

import (
	"math"
	"sort"
	"time"

	"go-bug-analysis/internal/csvparse"
)

// ComputeAge computes fix time statistics and backlog age ranking.
func ComputeAge(bugs []csvparse.Bug) *AgeData {
	return &AgeData{
		FixTime: computeFixTime(bugs),
		Backlog: computeBacklog(bugs),
	}
}

// computeFixTime calculates average, P50 (median), and distribution of fix times.
func computeFixTime(bugs []csvparse.Bug) *FixTimeStats {
	var durations []time.Duration

	for i := range bugs {
		b := &bugs[i]
		if b.ResolvedDate == "" || b.CreatedDate == "" {
			continue
		}

		created, okC := ParseDate(b.CreatedDate)
		resolved, okR := ParseDate(b.ResolvedDate)
		if !okC || !okR {
			continue
		}

		dur := resolved.Sub(created)
		if dur <= 0 {
			continue // skip negative or zero durations
		}

		durations = append(durations, dur)
	}

	if len(durations) == 0 {
		return nil
	}

	// Calculate average
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	avgHours := total.Hours() / float64(len(durations))

	// Calculate P50 (median)
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})
	var p50Hours float64
	n := len(durations)
	if n%2 == 1 {
		p50Hours = durations[n/2].Hours()
	} else {
		p50Hours = (durations[n/2-1].Hours() + durations[n/2].Hours()) / 2
	}

	// Distribution buckets
	buckets := []DistBucket{
		{Label: "0-1天", Count: 0},
		{Label: "2-3天", Count: 0},
		{Label: "4-7天", Count: 0},
		{Label: "7天以上", Count: 0},
	}
	for _, d := range durations {
		h := d.Hours()
		switch {
		case h <= 24:
			buckets[0].Count++
		case h <= 72:
			buckets[1].Count++
		case h <= 168:
			buckets[2].Count++
		default:
			buckets[3].Count++
		}
	}

	return &FixTimeStats{
		AvgHours:      roundTo(avgHours, 1),
		AvgDays:       roundTo(avgHours/24, 1),
		P50Hours:      roundTo(p50Hours, 1),
		P50Days:       roundTo(p50Hours/24, 1),
		Distribution:  buckets,
		TotalResolved: len(durations),
	}
}

// computeBacklog builds a list of active (unresolved) bugs sorted by age descending.
func computeBacklog(bugs []csvparse.Bug) []BacklogItem {
	now := time.Now()
	var items []BacklogItem

	for i := range bugs {
		b := &bugs[i]
		if b.Status != "激活" {
			continue
		}

		ageDays := 0
		if created, ok := ParseDate(b.CreatedDate); ok {
			ageDays = int(now.Sub(created).Hours() / 24)
			if ageDays < 0 {
				ageDays = 0
			}
		}

		items = append(items, BacklogItem{
			ID:          b.ID,
			Title:       b.Title,
			Severity:    b.Severity,
			Assignee:    b.Assignee,
			CreatedDate: b.CreatedDate,
			AgeDays:     ageDays,
		})
	}

	// Sort descending by age
	sort.Slice(items, func(i, j int) bool {
		return items[i].AgeDays > items[j].AgeDays
	})

	return items
}

// roundTo rounds a float64 to the specified number of decimal places.
func roundTo(val float64, places int) float64 {
	pow := math.Pow(10, float64(places))
	return math.Round(val*pow) / pow
}
