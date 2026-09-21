package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/system"
)

// AdminRoleList 角色列表
func (c *ControllerV1) AdminRoleList(ctx context.Context, req *v1.AdminRoleListReq) (res *v1.AdminRoleListRes, err error) {
	out, err := system.RoleList(ctx, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminRoleListRes{List: make([]v1.AdminRoleItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminRoleItem{
			Id:          fmtID(it.Id),
			Name:        it.Name,
			Code:        it.Code,
			Description: it.Description,
			Status:      it.Status,
		})
	}
	return res, nil
}
