package common

import (
	"context"

	"ecboot/api/common/v1"
	"ecboot/internal/service/shop"
)

// StoreDetail 门店详情（歇业店可见）
func (c *ControllerV1) StoreDetail(ctx context.Context, req *v1.StoreDetailReq) (res *v1.StoreDetailRes, err error) {
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	d, err := shop.PublicDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.StoreDetailRes{
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
