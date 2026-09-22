package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminAfterSaleApprove 同意（仅退款→待退款; 退货退款→待寄回）
// 权限: aftersale:audit
func (c *ControllerV1) AdminAfterSaleApprove(ctx context.Context, req *v1.AdminAfterSaleApproveReq) (res *v1.AdminAfterSaleApproveRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermAfterSaleAudit); err != nil {
		return nil, err
	}
	if err = shop.NewAfterSaleLogic().Approve(ctx, req.AfterSaleNo, adminOperator(ctx)); err != nil {
		return nil, err
	}
	return &v1.AdminAfterSaleApproveRes{Success: true}, nil
}
