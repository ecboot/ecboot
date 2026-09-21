package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/service/shop"
)

// AdminLogisticsDetail 物流公司详情
func (c *ControllerV1) AdminLogisticsDetail(ctx context.Context, req *v1.AdminLogisticsDetailReq) (res *v1.AdminLogisticsDetailRes, err error) {
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	d, err := shop.LogisticsDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.AdminLogisticsDetailRes{
		Id:           fmtID(d.Id),
		Code:         d.Code,
		Name:         d.Name,
		TrackingRule: d.TrackingRule,
		Status:       d.Status,
	}, nil
}
