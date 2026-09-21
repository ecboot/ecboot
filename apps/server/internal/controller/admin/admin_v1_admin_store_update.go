package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminStoreUpdate 修改门店
func (c *ControllerV1) AdminStoreUpdate(ctx context.Context, req *v1.AdminStoreUpdateReq) (res *v1.AdminStoreUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, "store:manage:update"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.AdminUpdate(ctx, id, storeInputFromReq(
		req.Name, req.ProvinceCode, req.CityCode, req.DistrictCode, req.DetailAddress,
		req.BusinessHours, req.ContactPhone, req.Longitude, req.Latitude, req.PickupEnabled, req.Status,
	)); err != nil {
		return nil, err
	}
	return &v1.AdminStoreUpdateRes{Success: true}, nil
}
