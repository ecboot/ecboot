package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// ReviewExtra 追加评价（一次; 主评 90 天内）
func (c *ControllerV1) ReviewExtra(ctx context.Context, req *v1.ReviewExtraReq) (res *v1.ReviewExtraRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	reviewId, err := parseID(req.ReviewId)
	if err != nil {
		return nil, err
	}
	if err = shop.NewReviewLogic().Extra(ctx, userId, reviewId, req.Content, req.Images); err != nil {
		return nil, err
	}
	return &v1.ReviewExtraRes{Success: true}, nil
}
