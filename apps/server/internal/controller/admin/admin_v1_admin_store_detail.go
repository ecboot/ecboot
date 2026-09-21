package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/service/shop"
)

// AdminStoreDetail 门店详情
func (c *ControllerV1) AdminStoreDetail(ctx context.Context, req *v1.AdminStoreDetailReq) (res *v1.AdminStoreDetailRes, err error) {
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	d, err := shop.AdminDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.AdminStoreDetailRes{
		Id:            fmtID(d.Id),
		StoreNo:       d.StoreNo,
		Name:          d.Name,
		ProvinceCode:  d.ProvinceCode,
		CityCode:      d.CityCode,
		DistrictCode:  d.DistrictCode,
		DetailAddress: d.DetailAddress,
		Longitude:     d.Longitude,
		Latitude:      d.Latitude,
		BusinessHours: d.BusinessHours,
		ContactPhone:  d.ContactPhone,
		PickupEnabled: d.PickupEnabled,
		Status:        d.Status,
	}, nil
}
