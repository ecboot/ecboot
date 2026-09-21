// ops.go 治理域运营面——审计（V10 双日志）/ 系统配置（V30）/ 看板 / 通知模板（V26）。
// 规则: 审计只追加（含失败尝试）; 配置是覆盖层（停用/缺失回退代码默认——不阻断业务）;
// 看板为只读统计（口径随实现细化）。
package system

import "context"

// IAuditLogic 审计查询。
type IAuditLogic interface {
	OperationLogs(ctx context.Context, q AuditQuery) (*PageResult[OperationLogItem], error)
	LoginLogs(ctx context.Context, q AuditQuery) (*PageResult[LoginLogItem], error)
	// Write 操作审计写入（AOP 切面调用: 写操作成功/失败均记; 异步不阻断业务）。
	Write(ctx context.Context, e OperationLogEntry) error
}

type AuditQuery struct {
	AdminId   int64
	Username  string
	Module    string
	StartTime string
	EndTime   string
	PageQuery
}

type OperationLogItem struct {
	Id           int64  `json:"id"`
	Username     string `json:"username"`
	Module       string `json:"module"`
	Operation    string `json:"operation"`
	Method       string `json:"method"`
	RequestUri   string `json:"requestUri"`
	ResultStatus int    `json:"resultStatus"`
	Ip           string `json:"ip"`
	CostMs       int    `json:"costMs"`
	CreatedAt    string `json:"createdAt"`
}

type LoginLogItem struct {
	Username    string `json:"username"`
	AdminId     int64  `json:"adminId" dc:"成功时回填"`
	LoginStatus int    `json:"loginStatus" dc:"1成功 2失败"`
	Ip          string `json:"ip"`
	CreatedAt   string `json:"createdAt"`
}

type OperationLogEntry struct {
	AdminId       int64
	Username      string
	Module        string
	Operation     string
	Method        string
	RequestUri    string
	RequestParams map[string]any `json:"-" dc:"脱敏后"`
	ResultStatus  int
	ErrorMsg      string
	Ip            string
	CostMs        int
}

// IConfigLogic 系统配置。
type IConfigLogic interface {
	List(ctx context.Context) ([]ConfigItem, error)
	// Update 修改配置（代码默认值兜底语义不变; 变更 30 分钟缓存窗口内生效）。
	Update(ctx context.Context, code, value string, status int) error
}

type ConfigItem struct {
	Code        string `json:"code"`
	Value       string `json:"value"`
	ValueType   int    `json:"valueType"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

// IDashboardLogic 运营看板（只读统计）。
type IDashboardLogic interface {
	Trade(ctx context.Context, startTime, endTime string) (*TradeDashboard, error)
	Member(ctx context.Context, startTime, endTime string) (*MemberDashboard, error)
	Product(ctx context.Context) (*ProductDashboard, error)
}

type TradeDashboard struct {
	OrderCount     int64  `json:"orderCount"`
	SalesAmount    string `json:"salesAmount"`
	RefundAmount   string `json:"refundAmount"`
	PendingDeliver int64  `json:"pendingDeliver"`
}

type MemberDashboard struct {
	NewCount     int64 `json:"newCount"`
	ActiveCount  int64 `json:"activeCount"`
	DormantCount int64 `json:"dormantCount" dc:"≥90天口径"`
}

type ProductDashboard struct {
	OnSaleCount   int64 `json:"onSaleCount"`
	LowStockCount int64 `json:"lowStockCount"`
	PendingReview int64 `json:"pendingReview"`
}

// INotifyTemplateLogic 通知模板管理（渠道现实建模: 平台侧模板ID+参数契约; 站内信全文自控）。
type INotifyTemplateLogic interface {
	List(ctx context.Context, channel int, page PageQuery) (*PageResult[NotifyTemplateItem], error)
	Create(ctx context.Context, in NotifyTemplateInput) (int64, error)
	Update(ctx context.Context, id int64, in NotifyTemplateInput) error
	Delete(ctx context.Context, id int64) error
}

type NotifyTemplateInput struct {
	Code               string
	Channel            int
	ExternalTemplateId string
	Title              string
	ContentTemplate    string
	ParamsSchema       map[string]any
	Status             int
}

type NotifyTemplateItem struct {
	Id                 int64  `json:"id"`
	Code               string `json:"code"`
	Channel            int    `json:"channel"`
	ExternalTemplateId string `json:"externalTemplateId"`
	Title              string `json:"title"`
	ContentTemplate    string `json:"contentTemplate"`
	Status             int    `json:"status"`
}
