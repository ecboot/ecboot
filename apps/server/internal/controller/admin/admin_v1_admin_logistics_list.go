package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminLogisticsList 物流公司列表
func (c *ControllerV1) AdminLogisticsList(ctx context.Context, req *v1.AdminLogisticsListReq) (res *v1.AdminLogisticsListRes, err error) {
	out, err := shop.LogisticsList(ctx, req.Status, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminLogisticsListRes{List: make([]v1.AdminLogisticsItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminLogisticsItem{
			Id:           fmtID(it.Id),
			Code:         it.Code,
			Name:         it.Name,
			TrackingRule: it.TrackingRule,
			Status:       it.Status,
		})
	}
	return res, nil
}
