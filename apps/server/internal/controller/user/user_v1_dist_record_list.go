package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// DistRecordList 我的佣金记录。
func (c *ControllerV1) DistRecordList(ctx context.Context, req *v1.DistRecordListReq) (res *v1.DistRecordListRes, err error) {
	out, err := user.NewDistributionLogic().Records(ctx, middleware.CtxUserIdFrom(ctx), req.Status, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.DistRecordListRes{List: make([]v1.DistRecordItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.DistRecordItem{
			OrderNo: it.OrderNo, Level: it.Level, Amount: it.Amount,
			Status: it.Status, SettleTime: it.SettleTime, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
