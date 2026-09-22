package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/system"
)

// AdminOperationLogList 操作审计日志。
func (c *ControllerV1) AdminOperationLogList(ctx context.Context, req *v1.AdminOperationLogListReq) (res *v1.AdminOperationLogListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermSystemAuditRead); err != nil {
		return nil, err
	}
	var adminId int64
	if req.AdminId != "" {
		if adminId, err = parseID(req.AdminId); err != nil {
			return nil, err
		}
	}
	out, err := system.NewAuditLogic().OperationLogs(ctx, model.AuditQuery{
		AdminId: adminId, Module: req.Module,
		StartTime: req.StartTime, EndTime: req.EndTime,
		PageReq: model.PageReq{Page: req.Page, PageSize: req.PageSize},
	})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminOperationLogListRes{List: make([]v1.AdminOperationLogItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminOperationLogItem{
			Id: fmtID(it.Id), Username: it.Username, Module: it.Module, Operation: it.Operation,
			Method: it.Method, RequestUri: it.RequestUri, ResultStatus: it.ResultStatus,
			Ip: it.Ip, CostMs: it.CostMs, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
