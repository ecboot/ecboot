package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/system"
)

// AdminLoginLogList 后台登录审计（路径 /admin/admin-login-logs; 含失败尝试）。
func (c *ControllerV1) AdminLoginLogList(ctx context.Context, req *v1.AdminLoginLogListReq) (res *v1.AdminLoginLogListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermSystemAuditRead); err != nil {
		return nil, err
	}
	out, err := system.NewAuditLogic().LoginLogs(ctx, model.AuditQuery{
		Username: req.Username,
		PageReq:  model.PageReq{Page: req.Page, PageSize: req.PageSize},
	})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminLoginLogListRes{List: make([]v1.AdminLoginLogItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminLoginLogItem{
			Username: it.Username, AdminId: fmtID(it.AdminId), LoginStatus: it.LoginStatus,
			Ip: it.Ip, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
