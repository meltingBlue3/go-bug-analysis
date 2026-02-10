package csvparse

// Bug represents a single bug record parsed from a Zentao CSV export.
// All fields are strings — date parsing is deferred to analysis modules.
type Bug struct {
	ID              string `json:"id"`              // Bug编号
	Product         string `json:"product"`         // 所属产品
	Module          string `json:"module"`          // 所属模块
	Title           string `json:"title"`           // Bug标题
	Severity        string `json:"severity"`        // 严重程度 (1-4)
	Priority        string `json:"priority"`        // 优先级
	BugType         string `json:"bugType"`         // Bug类型
	Status          string `json:"status"`          // Bug状态 (激活/已解决/已关闭)
	Creator         string `json:"creator"`         // 由谁创建
	CreatedDate     string `json:"createdDate"`     // 创建日期
	Assignee        string `json:"assignee"`        // 指派给
	AssignedDate    string `json:"assignedDate"`    // 指派日期
	Resolver        string `json:"resolver"`        // 解决者
	Resolution      string `json:"resolution"`      // 解决方案
	ResolvedDate    string `json:"resolvedDate"`    // 解决日期
	Closer          string `json:"closer"`          // 由谁关闭
	ClosedDate      string `json:"closedDate"`      // 关闭日期
	ActivationCount string `json:"activationCount"` // 激活次数
	Deadline        string `json:"deadline"`        // 截止日期
	AffectedVersion string `json:"affectedVersion"` // 影响版本
	ResolvedVersion string `json:"resolvedVersion"` // 解决版本
	Keywords        string `json:"keywords"`        // 关键词
}

// HeaderMap maps Chinese column headers from Zentao CSV exports to Bug field names.
// Unrecognized columns are silently ignored during parsing.
var HeaderMap = map[string]string{
	"Bug编号": "ID",
	"所属产品": "Product",
	"所属模块": "Module",
	"Bug标题": "Title",
	"严重程度": "Severity",
	"优先级":  "Priority",
	"Bug类型": "BugType",
	"Bug状态": "Status",
	"由谁创建": "Creator",
	"创建日期": "CreatedDate",
	"指派给":  "Assignee",
	"指派日期": "AssignedDate",
	"解决者":  "Resolver",
	"解决方案": "Resolution",
	"解决日期": "ResolvedDate",
	"由谁关闭": "Closer",
	"关闭日期": "ClosedDate",
	"激活次数": "ActivationCount",
	"截止日期": "Deadline",
	"影响版本": "AffectedVersion",
	"解决版本": "ResolvedVersion",
	"关键词":  "Keywords",
}

// RequiredHeaders lists the Chinese column names that must be present in the CSV.
// Missing any of these causes a parse error with a descriptive Chinese message.
var RequiredHeaders = []string{
	"Bug编号",
	"Bug标题",
	"严重程度",
	"Bug状态",
	"由谁创建",
	"创建日期",
	"指派给",
}

// ParseResult holds the output of a successful CSV parse operation.
type ParseResult struct {
	Bugs      []Bug    `json:"bugs"`      // Parsed bug records
	TotalRows int      `json:"totalRows"` // Total data rows processed
	Warnings  []string `json:"warnings"`  // Non-fatal warnings (e.g., empty fields)
	Columns   []string `json:"columns"`   // Recognized column names found in header
}
