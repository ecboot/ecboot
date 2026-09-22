package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// AfterSaleLogistics 填写寄回单号（退货退款; 仅"待买家寄回"可填）
func (c *ControllerV1) AfterSaleLogistics(ctx context.Context, req *v1.AfterSaleLogisticsReq) (res *v1.AfterSaleLogisticsRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	if err = shop.NewAfterSaleLogic().SubmitReturn(ctx, userId, req.AfterSaleNo, req.ReturnLogisticsNo); err != nil {
		return nil, err
	}
	return &v1.AfterSaleLogisticsRes{Success: true}, nil
}
