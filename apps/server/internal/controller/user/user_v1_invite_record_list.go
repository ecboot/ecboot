package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// InviteRecordList 邀请激励记录。
func (c *ControllerV1) InviteRecordList(ctx context.Context, req *v1.InviteRecordListReq) (res *v1.InviteRecordListRes, err error) {
	out, err := user.NewDistributionLogic().InviteRecords(ctx, middleware.CtxUserIdFrom(ctx), model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.InviteRecordListRes{List: make([]v1.InviteRecordItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.InviteRecordItem{
			NewUser: it.NewUser, RewardDesc: it.RewardDesc, Status: it.Status, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
