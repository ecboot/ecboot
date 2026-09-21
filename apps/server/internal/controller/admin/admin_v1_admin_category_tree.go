package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminCategoryTree 分类树（含禁用）
func (c *ControllerV1) AdminCategoryTree(ctx context.Context, req *v1.AdminCategoryTreeReq) (res *v1.AdminCategoryTreeRes, err error) {
	tree, err := shop.NewProductLogic().AdminCategoryTree(ctx)
	if err != nil {
		return nil, err
	}
	var conv func(nodes []model.AdminCategoryNode) []v1.AdminCategoryNode
	conv = func(nodes []model.AdminCategoryNode) []v1.AdminCategoryNode {
		out := make([]v1.AdminCategoryNode, 0, len(nodes))
		for _, n := range nodes {
			out = append(out, v1.AdminCategoryNode{
				Id: fmtID(n.Id), ParentId: fmtID(n.ParentId), Name: n.Name, Icon: n.Icon,
				Level: n.Level, Sort: n.Sort, Status: n.Status, Children: conv(n.Children),
			})
		}
		return out
	}
	return &v1.AdminCategoryTreeRes{Tree: conv(tree)}, nil
}
