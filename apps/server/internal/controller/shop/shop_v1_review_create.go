package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// ReviewCreate 提交评价（一项一评; 已完成订单项）
func (c *ControllerV1) ReviewCreate(ctx context.Context, req *v1.ReviewCreateReq) (res *v1.ReviewCreateRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	itemId, err := parseID(req.OrderItemId)
	if err != nil {
		return nil, err
	}
	id, err := shop.NewReviewLogic().Create(ctx, userId, model.ReviewCreateInput{
		OrderItemId: itemId,
		Score:       req.Score,
		Content:     req.Content,
		Images:      req.Images,
		IsAnonymous: req.IsAnonymous,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ReviewCreateRes{ReviewId: fmtID(id)}, nil
}
