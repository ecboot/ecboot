package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminStoreCreate 新增门店
func (c *ControllerV1) AdminStoreCreate(ctx context.Context, req *v1.AdminStoreCreateReq) (res *v1.AdminStoreCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, "store:manage:create"); err != nil {
		return nil, err
	}
	id, err := shop.AdminCreate(ctx, storeInputFromReq(
		req.Name, req.ProvinceCode, req.CityCode, req.DistrictCode, req.DetailAddress,
		req.BusinessHours, req.ContactPhone, req.Longitude, req.Latitude, req.PickupEnabled, 1))
	if err != nil {
		return nil, err
	}
	return &v1.AdminStoreCreateRes{Id: fmtID(id)}, nil
}
