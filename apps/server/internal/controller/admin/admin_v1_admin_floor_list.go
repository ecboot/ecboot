package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminFloorList 楼层列表
func (c *ControllerV1) AdminFloorList(ctx context.Context, req *v1.AdminFloorListReq) (res *v1.AdminFloorListRes, err error) {
	out, err := shop.FloorList(ctx, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminFloorListRes{List: make([]v1.AdminFloorItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminFloorItem{
			Id: fmtID(it.Id), FloorType: it.FloorType, Title: it.Title,
			Config: it.Config, Sort: it.Sort, Status: it.Status,
		})
	}
	return res, nil
}
