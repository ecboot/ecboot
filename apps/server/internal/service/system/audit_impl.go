// audit_impl.go IAuditLogic 实现（018 批次 12——审计查询面; 写入侧由批次 01 AOP 产出）。
// 只追加日志的只读查询: 筛选分页, 敏感字段按既有脱敏口径（批次 01 登录日志不落明文密码）。
package system

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
)

// AuditLogicImpl IAuditLogic 实现。
type AuditLogicImpl struct{}

func NewAuditLogic() *AuditLogicImpl { return &AuditLogicImpl{} }

var _ IAuditLogic = (*AuditLogicImpl)(nil)

// OperationLogs 操作日志分页（模块/时间/操作者筛选）。
func (i *AuditLogicImpl) OperationLogs(ctx context.Context, q model.AuditQuery) (*model.PageResult[model.OperationLogItem], error) {
	q.PageReq = q.Normalized()
	m := g.DB().Model("admin_operation_log").Ctx(ctx)
	if q.AdminId > 0 {
		m = m.Where("admin_id", q.AdminId)
	}
	if q.Module != "" {
		m = m.Where("module", q.Module)
	}
	if q.StartTime != "" {
		m = m.Where("created_at >= ?", q.StartTime)
	}
	if q.EndTime != "" {
		m = m.Where("created_at <= ?", q.EndTime)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计操作日志失败")
	}
	recs, err := m.OrderDesc("id").Page(q.Page, q.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询操作日志失败")
	}
	list := make([]model.OperationLogItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.OperationLogItem{
			Id: r["id"].Int64(), Username: r["username"].String(),
			Module: r["module"].String(), Operation: r["operation"].String(),
			Method: r["method"].String(), RequestUri: r["request_uri"].String(),
			ResultStatus: r["result_status"].Int(), Ip: r["ip"].String(),
			CostMs: r["cost_ms"].Int(), CreatedAt: r["created_at"].String(),
		})
	}
	return &model.PageResult[model.OperationLogItem]{List: list, Total: int64(total)}, nil
}

// LoginLogs 登录日志分页（操作者/时间筛选; 含失败尝试——防暴力破解分析依据）。
func (i *AuditLogicImpl) LoginLogs(ctx context.Context, q model.AuditQuery) (*model.PageResult[model.AdminLoginLogItem], error) {
	q.PageReq = q.Normalized()
	m := g.DB().Model("admin_login_log").Ctx(ctx)
	if q.AdminId > 0 {
		m = m.Where("admin_id", q.AdminId)
	}
	if q.Username != "" {
		m = m.Where("username", q.Username)
	}
	if q.StartTime != "" {
		m = m.Where("created_at >= ?", q.StartTime)
	}
	if q.EndTime != "" {
		m = m.Where("created_at <= ?", q.EndTime)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计登录日志失败")
	}
	recs, err := m.OrderDesc("id").Page(q.Page, q.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询登录日志失败")
	}
	list := make([]model.AdminLoginLogItem, 0, len(recs))
	for _, r := range recs {
		adminId := int64(0)
		if !r["admin_id"].IsNil() {
			adminId = r["admin_id"].Int64()
		}
		list = append(list, model.AdminLoginLogItem{
			Username: r["username"].String(), AdminId: adminId,
			LoginStatus: r["login_status"].Int(), Ip: r["ip"].String(),
			CreatedAt: r["created_at"].String(),
		})
	}
	return &model.PageResult[model.AdminLoginLogItem]{List: list, Total: int64(total)}, nil
}

// Write 操作审计写入（AOP 切面调用; 异步语义由调用方决定, 本方法同步落库不告警阻断）。
// 批次 01 既有产出路径——本批补齐接口实现（此前 IAuditLogic 无实现者）。
func (i *AuditLogicImpl) Write(ctx context.Context, e model.OperationLogEntry) error {
	params := ""
	if e.RequestParams != nil {
		if b, err := json.Marshal(e.RequestParams); err == nil {
			params = string(b)
		}
	}
	_, err := g.DB().Model("admin_operation_log").Ctx(ctx).Data(g.Map{
		"admin_id": e.AdminId, "username": e.Username,
		"module": e.Module, "operation": e.Operation,
		"method": e.Method, "request_uri": e.RequestUri,
		"request_params": params, "result_status": e.ResultStatus,
		"error_msg": e.ErrorMsg, "ip": e.Ip,
	}).Insert()
	return gerror.Wrap(err, "写入操作审计失败")
}
