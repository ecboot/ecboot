package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// DistAccountLogList 账户流水。
func (c *ControllerV1) DistAccountLogList(ctx context.Context, req *v1.DistAccountLogListReq) (res *v1.DistAccountLogListRes, err error) {
	out, err := user.NewDistributionLogic().AccountLogs(ctx, middleware.CtxUserIdFrom(ctx), req.BizType, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.DistAccountLogListRes{List: make([]v1.DistAccountLogItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.DistAccountLogItem{
			BizType: it.BizType, Amount: it.Amount, BalanceAfter: it.BalanceAfter,
			BizNo: it.BizNo, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
