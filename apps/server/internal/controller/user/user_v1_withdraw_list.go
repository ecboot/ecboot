package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// WithdrawList 我的提现列表。
func (c *ControllerV1) WithdrawList(ctx context.Context, req *v1.WithdrawListReq) (res *v1.WithdrawListRes, err error) {
	out, err := user.NewDistributionLogic().WithdrawList(ctx, middleware.CtxUserIdFrom(ctx), req.Status, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.WithdrawListRes{List: make([]v1.WithdrawItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.WithdrawItem{
			WithdrawNo: it.WithdrawNo, Amount: it.Amount, Status: it.Status, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
