package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// AdminDistRecordList 全局佣金记录。
func (c *ControllerV1) AdminDistRecordList(ctx context.Context, req *v1.AdminDistRecordListReq) (res *v1.AdminDistRecordListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDistributionRead); err != nil {
		return nil, err
	}
	out, err := user.NewDistributionAdminLogic().AdminRecordList(ctx, req.Status, req.OrderNo, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminDistRecordListRes{List: make([]v1.AdminDistRecordItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminDistRecordItem{
			OrderNo: it.OrderNo, Beneficiary: it.Beneficiary, Level: it.Level,
			BaseAmount: it.BaseAmount, Amount: it.Amount, Status: it.Status, SettleTime: it.SettleTime,
		})
	}
	return res, nil
}
