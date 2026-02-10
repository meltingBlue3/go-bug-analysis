package analysis

// AnalysisResult is the top-level result returned by Analyze().
type AnalysisResult struct {
	KPI      *KPIData      `json:"kpi"`
	Severity *SeverityData `json:"severity"`
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
