package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// BargainLaunch 发起砍价（会员; 发起即首刀）
func (c *ControllerV1) BargainLaunch(ctx context.Context, req *v1.BargainLaunchReq) (res *v1.BargainLaunchRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	itemId, err := parseID(req.BargainItemId)
	if err != nil {
		return nil, err
	}
	out, err := shop.NewBargainLogic().Launch(ctx, userId, itemId)
	if err != nil {
		return nil, err
	}
	return &v1.BargainLaunchRes{RecordId: fmtID(out.RecordId), CurrentPrice: out.CurrentPrice}, nil
}
