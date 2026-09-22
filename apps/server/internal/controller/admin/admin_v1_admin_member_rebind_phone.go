package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AdminMemberRebindPhone 改绑手机号（唯一键兜底并发）。
func (c *ControllerV1) AdminMemberRebindPhone(ctx context.Context, req *v1.AdminMemberRebindPhoneReq) (res *v1.AdminMemberRebindPhoneRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermMemberUpdate); err != nil {
		return nil, err
	}
	uid, err := parseID(req.UserId)
	if err != nil {
		return nil, err
	}
	if err = user.NewMemberAdminLogic().RebindPhone(ctx, uid, req.NewPhone); err != nil {
		return nil, err
	}
	return &v1.AdminMemberRebindPhoneRes{Success: true}, nil
}
