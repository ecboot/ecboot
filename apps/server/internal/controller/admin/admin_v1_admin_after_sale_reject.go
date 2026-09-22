package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminAfterSaleReject 拒绝售后（原因必填, 对买家可见）
// 权限: aftersale:audit
func (c *ControllerV1) AdminAfterSaleReject(ctx context.Context, req *v1.AdminAfterSaleRejectReq) (res *v1.AdminAfterSaleRejectRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermAfterSaleAudit); err != nil {
		return nil, err
	}
	if err = shop.NewAfterSaleLogic().Reject(ctx, req.AfterSaleNo, req.RejectReason, adminOperator(ctx)); err != nil {
		return nil, err
	}
	return &v1.AdminAfterSaleRejectRes{Success: true}, nil
}
