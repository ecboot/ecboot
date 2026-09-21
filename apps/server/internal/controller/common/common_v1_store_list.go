package common

import (
	"context"

	"ecboot/api/common/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// StoreList 门店列表（区县筛选或经纬度附近检索, 距离排序; 仅营业）
func (c *ControllerV1) StoreList(ctx context.Context, req *v1.StoreListReq) (res *v1.StoreListRes, err error) {
	out, err := shop.PublicList(ctx, model.StoreQuery{
		DistrictCode: req.DistrictCode,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		RadiusKm:     req.RadiusKm,
		PageReq:      model.PageReq{Page: req.Page, PageSize: req.PageSize},
	})
	if err != nil {
		return nil, err
	}
	res = &v1.StoreListRes{List: make([]v1.StoreItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.StoreItem{
			Id:            fmtID(it.Id),
			StoreNo:       it.StoreNo,
			Name:          it.Name,
			DetailAddress: it.DetailAddress,
			DistanceM:     it.DistanceM,
			BusinessHours: it.BusinessHours,
			PickupEnabled: it.PickupEnabled,
			Status:        it.Status,
		})
	}
	return res, nil
}
