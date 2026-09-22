package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AdminMemberDisable 禁用/启用会员（审计 reason 由操作日志面承载）。
func (c *ControllerV1) AdminMemberDisable(ctx context.Context, req *v1.AdminMemberDisableReq) (res *v1.AdminMemberDisableRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermMemberUpdate); err != nil {
		return nil, err
	}
	uid, err := parseID(req.UserId)
	if err != nil {
		return nil, err
	}
	if err = user.NewMemberAdminLogic().Disable(ctx, uid, req.Disable, req.Reason); err != nil {
		return nil, err
	}
	return &v1.AdminMemberDisableRes{Success: true}, nil
}
