// dto_system.go system 域 DTO（自 internal/service/system 迁入 2026-09-21, 契约出入参统一归属 model）。
package model

type AuditQuery struct {
	AdminId   int64
	Username  string
	Module    string
	StartTime string
	EndTime   string
	PageReq
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

type AdminLoginLogItem struct {
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

type ConfigItem struct {
	Code        string `json:"code"`
	Value       string `json:"value"`
	ValueType   int    `json:"valueType"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int    `json:"status"`
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

type RoleItem struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

type RoleInput struct {
	Name        string
	Code        string
	Description string
	Status      int
}

type RoleDetailView struct {
	Id            int64   `json:"id"`
	Name          string  `json:"name"`
	Code          string  `json:"code"`
	Description   string  `json:"description"`
	Status        int     `json:"status"`
	PermissionIds []int64 `json:"permissionIds"`
}

type PermissionNode struct {
	Id       int64            `json:"id"`
	ParentId int64            `json:"parentId"`
	Name     string           `json:"name"`
	Code     string           `json:"code"`
	Type     int              `json:"type" dc:"1菜单 2按钮 3接口"`
	Sort     int              `json:"sort"`
	Status   int              `json:"status"`
	Children []PermissionNode `json:"children"`
}

type AdminUserItem struct {
	Id            int64    `json:"id"`
	Username      string   `json:"username"`
	RealName      string   `json:"realName"`
	IsSuper       bool     `json:"isSuper"`
	Roles         []string `json:"roles"`
	Status        int      `json:"status"`
	LastLoginTime string   `json:"lastLoginTime"`
}

type AdminUserInput struct {
	Username string
	Password string
	RealName string
}

type AdminUserUpdateInput struct {
	RealName string
	Status   int
}

type AdminLoginResult struct {
	AdminId      int64  `json:"adminId"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	RealName     string `json:"realName"`
	IsSuper      bool   `json:"isSuper"`
}

// AdminProfile 后台个人信息（007-admin-base research D6）。
type AdminProfile struct {
	Username string
	RealName string
	Roles    []string // 角色编码列表
}
