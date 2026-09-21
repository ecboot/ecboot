// ops.go 治理域运营面——审计（V10 双日志）/ 系统配置（V30）/ 看板 / 通知模板（V26）。
// 规则: 审计只追加（含失败尝试）; 配置是覆盖层（停用/缺失回退代码默认——不阻断业务）;
// 看板为只读统计（口径随实现细化）。
package system

import (
	"context"

	"ecboot/internal/model"
)

// IAuditLogic 审计查询。
type IAuditLogic interface {
	OperationLogs(ctx context.Context, q model.AuditQuery) (*model.PageResult[model.OperationLogItem], error)
	LoginLogs(ctx context.Context, q model.AuditQuery) (*model.PageResult[model.AdminLoginLogItem], error)
	// Write 操作审计写入（AOP 切面调用: 写操作成功/失败均记; 异步不阻断业务）。
	Write(ctx context.Context, e model.OperationLogEntry) error
}

// IConfigLogic 系统配置。
type IConfigLogic interface {
	List(ctx context.Context) ([]model.ConfigItem, error)
	// Update 修改配置（代码默认值兜底语义不变; 变更 30 分钟缓存窗口内生效）。
	Update(ctx context.Context, code, value string, status int) error
}

// IDashboardLogic 运营看板（只读统计）。
type IDashboardLogic interface {
	Trade(ctx context.Context, startTime, endTime string) (*model.TradeDashboard, error)
	Member(ctx context.Context, startTime, endTime string) (*model.MemberDashboard, error)
	Product(ctx context.Context) (*model.ProductDashboard, error)
}

// INotifyTemplateLogic 通知模板管理（渠道现实建模: 平台侧模板ID+参数契约; 站内信全文自控）。
type INotifyTemplateLogic interface {
	List(ctx context.Context, channel int, page model.PageReq) (*model.PageResult[model.NotifyTemplateItem], error)
	Create(ctx context.Context, in model.NotifyTemplateInput) (int64, error)
	Update(ctx context.Context, id int64, in model.NotifyTemplateInput) error
	Delete(ctx context.Context, id int64) error
}
