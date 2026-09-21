package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminLogisticsUpdate 修改物流公司（编码不可改——入参无该字段）
func (c *ControllerV1) AdminLogisticsUpdate(ctx context.Context, req *v1.AdminLogisticsUpdateReq) (res *v1.AdminLogisticsUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, "logistics:company:manage"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.LogisticsUpdate(ctx, id, model.LogisticsCompanyInput{
		Name: req.Name, TrackingRule: req.TrackingRule, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminLogisticsUpdateRes{Success: true}, nil
}
