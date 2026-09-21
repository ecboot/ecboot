package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/system"
)

// AdminPermissionTree 权限树（菜单/按钮/接口统一树）
func (c *ControllerV1) AdminPermissionTree(ctx context.Context, req *v1.AdminPermissionTreeReq) (res *v1.AdminPermissionTreeRes, err error) {
	tree, err := system.PermissionTree(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.AdminPermissionTreeRes{Tree: make([]v1.AdminPermissionNode, 0, len(tree))}
	var conv func(nodes []model.PermissionNode) []v1.AdminPermissionNode
	conv = func(nodes []model.PermissionNode) []v1.AdminPermissionNode {
		out := make([]v1.AdminPermissionNode, 0, len(nodes))
		for _, n := range nodes {
			out = append(out, v1.AdminPermissionNode{
				Id:       fmtID(n.Id),
				ParentId: fmtID(n.ParentId),
				Name:     n.Name,
				Code:     n.Code,
				Type:     n.Type,
				Sort:     n.Sort,
				Status:   n.Status,
				Children: conv(n.Children),
			})
		}
		return out
	}
	res.Tree = conv(tree)
	return res, nil
}
