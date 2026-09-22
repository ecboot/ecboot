package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// BargainCut 帮砍一刀（会员; 一人一刀）
func (c *ControllerV1) BargainCut(ctx context.Context, req *v1.BargainCutReq) (res *v1.BargainCutRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	recordId, err := parseID(req.RecordId)
	if err != nil {
		return nil, err
	}
	out, err := shop.NewBargainLogic().Cut(ctx, userId, recordId)
	if err != nil {
		return nil, err
	}
	return &v1.BargainCutRes{CutAmount: out.CutAmount, CurrentPrice: out.CurrentPrice, FloorReached: out.FloorReached}, nil
}
