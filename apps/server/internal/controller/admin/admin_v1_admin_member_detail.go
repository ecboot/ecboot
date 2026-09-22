package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AdminMemberDetail 会员详情（资产概要+订单统计）。
func (c *ControllerV1) AdminMemberDetail(ctx context.Context, req *v1.AdminMemberDetailReq) (res *v1.AdminMemberDetailRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermMemberRead); err != nil {
		return nil, err
	}
	uid, err := parseID(req.UserId)
	if err != nil {
		return nil, err
	}
	it, err := user.NewMemberAdminLogic().AdminDetail(ctx, uid)
	if err != nil {
		return nil, err
	}
	assets := map[string]string{}
	for k, v := range it.Assets {
		assets[k] = v
	}
	stats := map[string]int64{}
	for k, v := range it.OrderStats {
		stats[k] = v
	}
	return &v1.AdminMemberDetailRes{
		UserId: fmtID(it.UserId), Nickname: it.Nickname, Phone: it.Phone,
		Gender: it.Gender, Level: int(it.Level), GrowthValue: it.GrowthValue, Status: it.Status,
		Assets: assets, OrderStats: stats, CreatedAt: it.CreatedAt,
	}, nil
}
