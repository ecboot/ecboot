package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// CategoryTree 商品分类树（禁用不返回; 公开）
func (c *ControllerV1) CategoryTree(ctx context.Context, req *v1.CategoryTreeReq) (res *v1.CategoryTreeRes, err error) {
	tree, err := shop.NewProductLogic().CategoryTree(ctx)
	if err != nil {
		return nil, err
	}
	var conv func(nodes []model.CategoryNode) []v1.CategoryNode
	conv = func(nodes []model.CategoryNode) []v1.CategoryNode {
		out := make([]v1.CategoryNode, 0, len(nodes))
		for _, n := range nodes {
			out = append(out, v1.CategoryNode{
				Id: fmtID(n.Id), Name: n.Name, Icon: n.Icon, Children: conv(n.Children),
			})
		}
		return out
	}
	return &v1.CategoryTreeRes{Tree: conv(tree)}, nil
}
