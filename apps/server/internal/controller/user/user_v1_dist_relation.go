package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// DistRelation 我的邀请关系（上级+下级分页, 脱敏）。
func (c *ControllerV1) DistRelation(ctx context.Context, req *v1.DistRelationReq) (res *v1.DistRelationRes, err error) {
	out, err := user.NewDistributionLogic().Relations(ctx, middleware.CtxUserIdFrom(ctx), model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.DistRelationRes{Inviter: out.Inviter, Invitees: make([]v1.DistRelationInvitee, 0, len(out.Invitees)), Total: out.Total}
	for _, it := range out.Invitees {
		res.Invitees = append(res.Invitees, v1.DistRelationInvitee{
			UserId: fmtID(it.UserId), Nickname: it.Nickname, BindTime: it.BindTime,
		})
	}
	return res, nil
}
