package analysis

// AnalysisResult is the top-level result returned by Analyze().
type AnalysisResult struct {
	KPI      *KPIData      `json:"kpi"`
	Severity *SeverityData `json:"severity"`
	Age      *AgeData      `json:"age,omitempty"`
	Workload *WorkloadData `json:"workload,omitempty"`
}

// WorkloadData holds per-assignee bug distribution.
type WorkloadData struct {
	ByActive []AssigneeStats `json:"byActive"` // sorted descending by count
	ByTotal  []AssigneeStats `json:"byTotal"`  // sorted descending by count
}

// AssigneeStats represents one assignee with their bug count.
type AssigneeStats struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// AgeData holds fix time statistics and backlog age ranking.
type AgeData struct {
	FixTime *FixTimeStats `json:"fixTime"`
	Backlog []BacklogItem `json:"backlog"` // sorted descending by AgeDays
}

// FixTimeStats holds average/median fix times and distribution buckets.
type FixTimeStats struct {
	AvgHours      float64      `json:"avgHours"`      // average fix time in hours
	AvgDays       float64      `json:"avgDays"`       // average fix time in days (avgHours/24)
	P50Hours      float64      `json:"p50Hours"`      // median fix time in hours
	P50Days       float64      `json:"p50Days"`       // median fix time in days
	Distribution  []DistBucket `json:"distribution"`  // 4 buckets
	TotalResolved int          `json:"totalResolved"` // how many bugs had valid fix times
}

// DistBucket represents one bucket in the fix time distribution.
type DistBucket struct {
	Label string `json:"label"` // "0-1天", "2-3天", "4-7天", "7天以上"
	Count int    `json:"count"`
}

// BacklogItem represents one unresolved bug in the backlog ranking.
type BacklogItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Severity    string `json:"severity"`
	Assignee    string `json:"assignee"`
	CreatedDate string `json:"createdDate"`
	AgeDays     int    `json:"ageDays"` // days since creation
}

// KPIData holds bug count KPIs relative to today / yesterday.
type KPIData struct {
	TodayNew       int `json:"todayNew"`
	YesterdayNew   int `json:"yesterdayNew"`
	TodayFixed     int `json:"todayFixed"`
	YesterdayFixed int `json:"yesterdayFixed"`
	Total          int `json:"total"`
	Active         int `json:"active"`
	PendingVerify  int `json:"pendingVerify"`
}

// SeverityItem represents one severity level with its count.
type SeverityItem struct {
	Level string `json:"level"` // "1","2","3","4"
	Label string `json:"label"` // "致命","严重","一般","轻微"
	Count int    `json:"count"`
}

// SeverityData holds severity distributions for all bugs and today's new bugs.
type SeverityData struct {
	All     []SeverityItem `json:"all"`
	NewOnly []SeverityItem `json:"newOnly"` // today's new bugs only
}
