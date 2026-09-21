package admin

import (
	"context"
	"fmt"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminInventoryAdjust 库存调整（留痕 inventory_log; operator 按表注释约定 admin:{id}）
func (c *ControllerV1) AdminInventoryAdjust(ctx context.Context, req *v1.AdminInventoryAdjustReq) (res *v1.AdminInventoryAdjustRes, err error) {
	if err = middleware.RequirePerm(ctx, "inventory:adjust"); err != nil {
		return nil, err
	}
	skuId, err := parseID(req.SkuId)
	if err != nil {
		return nil, err
	}
	operator := fmt.Sprintf("admin:%d", middleware.CtxUserIdFrom(ctx))
	totalAfter, err := shop.InventoryAdjust(ctx, skuId, req.Delta, operator, req.Remark)
	if err != nil {
		return nil, err
	}
	return &v1.AdminInventoryAdjustRes{TotalAfter: totalAfter}, nil
}
