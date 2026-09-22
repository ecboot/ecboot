package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// AfterSaleCancel 撤销售后申请（仅"尚未进入资金环节": 待审核/待寄回/待退款）
func (c *ControllerV1) AfterSaleCancel(ctx context.Context, req *v1.AfterSaleCancelReq) (res *v1.AfterSaleCancelRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	if err = shop.NewAfterSaleLogic().Cancel(ctx, userId, req.AfterSaleNo); err != nil {
		return nil, err
	}
	return &v1.AfterSaleCancelRes{Success: true}, nil
}
