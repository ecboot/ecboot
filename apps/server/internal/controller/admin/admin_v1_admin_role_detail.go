package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/service/system"
)

// AdminRoleDetail 角色详情
func (c *ControllerV1) AdminRoleDetail(ctx context.Context, req *v1.AdminRoleDetailReq) (res *v1.AdminRoleDetailRes, err error) {
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	dv, err := system.RoleDetailView(ctx, id)
	if err != nil {
		return nil, err
	}
	permIds := make([]string, 0, len(dv.PermissionIds))
	for _, pid := range dv.PermissionIds {
		permIds = append(permIds, fmtID(pid))
	}
	return &v1.AdminRoleDetailRes{
		Id:            fmtID(dv.Id),
		Name:          dv.Name,
		Code:          dv.Code,
		Description:   dv.Description,
		Status:        dv.Status,
		PermissionIds: permIds,
	}, nil
}
