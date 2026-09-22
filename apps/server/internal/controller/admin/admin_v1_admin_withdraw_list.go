package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// AdminWithdrawList 提现列表。
func (c *ControllerV1) AdminWithdrawList(ctx context.Context, req *v1.AdminWithdrawListReq) (res *v1.AdminWithdrawListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDistributionWithdrawRead); err != nil {
		return nil, err
	}
	out, err := user.NewDistributionAdminLogic().AdminWithdrawList(ctx, req.Status, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminWithdrawListRes{List: make([]v1.AdminWithdrawItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminWithdrawItem{
			WithdrawNo: it.WithdrawNo, UserId: "", Amount: it.Amount,
			Status: it.Status, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
