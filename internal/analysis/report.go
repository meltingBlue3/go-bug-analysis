package analysis

import (
	"fmt"
	"strings"
)

// severityLabel maps severity level to Chinese label for reports.
var severityLabel = map[string]string{
	"1": "致命",
	"2": "严重",
	"3": "一般",
	"4": "轻微",
}

// truncateRunes truncates s to maxLen runes, appending "…" if truncated.
func truncateRunes(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "…"
}

// ComputeReport generates a daily quality report from the analysis result.
// It reads from already-computed fields (KPI, Age, Workload, Module) and
// produces both Markdown and plain text report strings.
func ComputeReport(result *AnalysisResult) *ReportData {
	date := Today().Format("2006-01-02")

	md := buildMarkdownReport(result, date)
	plain := buildPlainTextReport(result, date)

	return &ReportData{
		Markdown:  md,
		PlainText: plain,
		Date:      date,
	}
}

// formatNetChange formats the net change value with a +/- prefix.
func formatNetChange(todayNew, todayFixed int) string {
	net := todayNew - todayFixed
	if net > 0 {
		return fmt.Sprintf("+%d", net)
	}
	return fmt.Sprintf("%d", net)
}

// riskBugs filters backlog for severity 1-2 bugs with age > 7 days.
type riskBug struct {
	ID       string
	Title    string
	Severity string
	Assignee string
	AgeDays  int
}

func filterRiskBugs(result *AnalysisResult) []riskBug {
	if result.Age == nil || result.Age.Backlog == nil {
		return nil
	}

	var risks []riskBug
	for _, item := range result.Age.Backlog {
		if (item.Severity == "1" || item.Severity == "2") && item.AgeDays > 7 {
			risks = append(risks, riskBug{
				ID:       item.ID,
				Title:    item.Title,
				Severity: item.Severity,
				Assignee: item.Assignee,
				AgeDays:  item.AgeDays,
			})
		}
	}
	return risks
}

// buildMarkdownReport generates the Markdown format report.
func buildMarkdownReport(result *AnalysisResult, date string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# Bug 质量日报 — %s\n\n", date))

	// KPI Overview
	b.WriteString("## 📊 今日概览\n\n")
	if result.KPI != nil {
		kpi := result.KPI
		netChange := formatNetChange(kpi.TodayNew, kpi.TodayFixed)
		b.WriteString("| 指标 | 数值 |\n")
		b.WriteString("|------|------|\n")
		b.WriteString(fmt.Sprintf("| 今日新增 | %d |\n", kpi.TodayNew))
		b.WriteString(fmt.Sprintf("| 今日修复 | %d |\n", kpi.TodayFixed))
		b.WriteString(fmt.Sprintf("| 净增长 | %s |\n", netChange))
		b.WriteString(fmt.Sprintf("| 激活总数 | %d |\n", kpi.Active))
		b.WriteString(fmt.Sprintf("| 待验证 | %d |\n", kpi.PendingVerify))
	} else {
		b.WriteString("暂无数据\n")
	}
	b.WriteString("\n")

	// Risk Bugs
	b.WriteString("## ⚠️ 高风险 Bug 预警\n\n")
	b.WriteString("> 严重程度 1-2 级且积压超过 7 天\n\n")
	risks := filterRiskBugs(result)
	if len(risks) == 0 {
		b.WriteString("暂无高风险 Bug\n")
	} else {
		b.WriteString("| Bug编号 | 标题 | 严重程度 | 指派给 | 积压天数 |\n")
		b.WriteString("|---------|------|---------|-------|--------|\n")
		maxDisplay := 10
		displayCount := len(risks)
		if displayCount > maxDisplay {
			displayCount = maxDisplay
		}
		for i := 0; i < displayCount; i++ {
			r := risks[i]
			title := truncateRunes(r.Title, 30)
			sevLabel := severityLabel[r.Severity]
			if sevLabel == "" {
				sevLabel = r.Severity
			}
			assignee := r.Assignee
			if assignee == "" {
				assignee = "未指派"
			}
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %d天 |\n",
				r.ID, title, sevLabel, assignee, r.AgeDays))
		}
		if len(risks) > maxDisplay {
			b.WriteString(fmt.Sprintf("\n...及其他 %d 条\n", len(risks)-maxDisplay))
		}
	}
	b.WriteString("\n")

	// Bottleneck Modules Top 5
	b.WriteString("## 🔥 瓶颈模块 Top 5\n\n")
	if result.Module != nil && len(result.Module.Stats) > 0 {
		b.WriteString("| 模块 | 激活 Bug | 激活率 |\n")
		b.WriteString("|------|---------|-------|\n")
		// Filter modules with Active > 0, take top 5
		count := 0
		for _, ms := range result.Module.Stats {
			if ms.Active <= 0 {
				continue
			}
			b.WriteString(fmt.Sprintf("| %s | %d | %.1f%% |\n",
				ms.Name, ms.Active, ms.ActiveRate))
			count++
			if count >= 5 {
				break
			}
		}
		if count == 0 {
			b.WriteString("暂无数据\n")
		}
	} else {
		b.WriteString("暂无数据\n")
	}
	b.WriteString("\n")

	// Personnel Load Top 5
	b.WriteString("## 👥 人员负载 Top 5\n\n")
	if result.Workload != nil && len(result.Workload.ByActive) > 0 {
		b.WriteString("| 人员 | 激活 Bug 数 |\n")
		b.WriteString("|------|----------|\n")
		maxPersonnel := 5
		if len(result.Workload.ByActive) < maxPersonnel {
			maxPersonnel = len(result.Workload.ByActive)
		}
		for i := 0; i < maxPersonnel; i++ {
			a := result.Workload.ByActive[i]
			b.WriteString(fmt.Sprintf("| %s | %d |\n", a.Name, a.Count))
		}
	} else {
		b.WriteString("暂无数据\n")
	}

	return b.String()
}

// buildPlainTextReport generates the plain text format report.
func buildPlainTextReport(result *AnalysisResult, date string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("【Bug 质量日报】%s\n\n", date))

	// KPI Overview
	b.WriteString("▎今日概览\n")
	if result.KPI != nil {
		kpi := result.KPI
		netChange := formatNetChange(kpi.TodayNew, kpi.TodayFixed)
		b.WriteString(fmt.Sprintf("  今日新增: %d | 今日修复: %d | 净增长: %s\n",
			kpi.TodayNew, kpi.TodayFixed, netChange))
		b.WriteString(fmt.Sprintf("  激活总数: %d | 待验证: %d\n",
			kpi.Active, kpi.PendingVerify))
	} else {
		b.WriteString("  暂无数据\n")
	}
	b.WriteString("\n")

	// Risk Bugs
	b.WriteString("▎高风险 Bug 预警 (严重程度1-2级, 积压>7天)\n")
	risks := filterRiskBugs(result)
	if len(risks) == 0 {
		b.WriteString("  暂无高风险 Bug\n")
	} else {
		maxDisplay := 10
		displayCount := len(risks)
		if displayCount > maxDisplay {
			displayCount = maxDisplay
		}
		for i := 0; i < displayCount; i++ {
			r := risks[i]
			title := truncateRunes(r.Title, 30)
			sevLabel := severityLabel[r.Severity]
			if sevLabel == "" {
				sevLabel = r.Severity
			}
			assignee := r.Assignee
			if assignee == "" {
				assignee = "未指派"
			}
			b.WriteString(fmt.Sprintf("  #%s %s  %s | %s | %d天\n",
				r.ID, title, sevLabel, assignee, r.AgeDays))
		}
		if len(risks) > maxDisplay {
			b.WriteString(fmt.Sprintf("  ...及其他 %d 条\n", len(risks)-maxDisplay))
		}
	}
	b.WriteString("\n")

	// Bottleneck Modules Top 5
	b.WriteString("▎瓶颈模块 Top 5\n")
	if result.Module != nil && len(result.Module.Stats) > 0 {
		count := 0
		for _, ms := range result.Module.Stats {
			if ms.Active <= 0 {
				continue
			}
			b.WriteString(fmt.Sprintf("  %s: 激活 %d 个 (%.1f%%)\n",
				ms.Name, ms.Active, ms.ActiveRate))
			count++
			if count >= 5 {
				break
			}
		}
		if count == 0 {
			b.WriteString("  暂无数据\n")
		}
	} else {
		b.WriteString("  暂无数据\n")
	}
	b.WriteString("\n")

	// Personnel Load Top 5
	b.WriteString("▎人员负载 Top 5\n")
	if result.Workload != nil && len(result.Workload.ByActive) > 0 {
		maxPersonnel := 5
		if len(result.Workload.ByActive) < maxPersonnel {
			maxPersonnel = len(result.Workload.ByActive)
		}
		for i := 0; i < maxPersonnel; i++ {
			a := result.Workload.ByActive[i]
			b.WriteString(fmt.Sprintf("  %s: 激活 %d 个\n", a.Name, a.Count))
		}
	} else {
		b.WriteString("  暂无数据\n")
	}

	return b.String()
}
