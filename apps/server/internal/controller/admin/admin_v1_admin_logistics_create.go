package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminLogisticsCreate 新增物流公司
func (c *ControllerV1) AdminLogisticsCreate(ctx context.Context, req *v1.AdminLogisticsCreateReq) (res *v1.AdminLogisticsCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, "logistics:company:manage"); err != nil {
		return nil, err
	}
	id, err := shop.LogisticsCreate(ctx, model.LogisticsCompanyInput{
		Code: req.Code, Name: req.Name, TrackingRule: req.TrackingRule,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminLogisticsCreateRes{Id: fmtID(id)}, nil
}
