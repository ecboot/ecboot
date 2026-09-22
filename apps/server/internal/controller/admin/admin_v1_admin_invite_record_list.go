package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// AdminInviteRecordList 全局邀请激励记录。
func (c *ControllerV1) AdminInviteRecordList(ctx context.Context, req *v1.AdminInviteRecordListReq) (res *v1.AdminInviteRecordListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDistributionRead); err != nil {
		return nil, err
	}
	out, err := user.NewDistributionAdminLogic().AdminInviteRecords(ctx, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminInviteRecordListRes{List: make([]v1.AdminInviteRecordItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminInviteRecordItem{
			Inviter: it.Inviter, NewUser: it.NewUser, RewardDesc: it.RewardDesc,
			Trigger: it.Trigger, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
