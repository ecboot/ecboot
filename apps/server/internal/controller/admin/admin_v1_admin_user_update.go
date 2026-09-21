package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/system"
)

// AdminUserUpdate 修改后台账号
func (c *ControllerV1) AdminUserUpdate(ctx context.Context, req *v1.AdminUserUpdateReq) (res *v1.AdminUserUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, "system:admin:manage"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	// 操作者 ID 经 ctx 由 service 层自守卫读取（consts.CtxUserId）
	if err = system.AdminUserUpdate(ctx, id, model.AdminUserUpdateInput{
		RealName: req.RealName,
		Status:   req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminUserUpdateRes{Success: true}, nil
}
