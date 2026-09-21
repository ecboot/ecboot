package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminStoreList 门店列表
func (c *ControllerV1) AdminStoreList(ctx context.Context, req *v1.AdminStoreListReq) (res *v1.AdminStoreListRes, err error) {
	out, err := shop.AdminList(ctx, req.Status, req.Keyword, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminStoreListRes{List: make([]v1.AdminStoreItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminStoreItem{
			Id:            fmtID(it.Id),
			StoreNo:       it.StoreNo,
			Name:          it.Name,
			ProvinceCode:  it.ProvinceCode,
			CityCode:      it.CityCode,
			DistrictCode:  it.DistrictCode,
			DetailAddress: it.DetailAddress,
			BusinessHours: it.BusinessHours,
			ContactPhone:  it.ContactPhone,
			PickupEnabled: it.PickupEnabled,
			Status:        it.Status,
		})
	}
	return res, nil
}
