package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminDashboardMember 会员看板（FR-2）: 新增/活跃（无窗口取近 30 天）/休眠（阈值读 system_config）。
func (c *ControllerV1) AdminDashboardMember(ctx context.Context, req *v1.AdminDashboardMemberReq) (res *v1.AdminDashboardMemberRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDashboardRead); err != nil {
		return nil, err
	}
	out, err := system.NewDashboardLogic().Member(ctx, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	return &v1.AdminDashboardMemberRes{
		NewCount: out.NewCount, ActiveCount: out.ActiveCount, DormantCount: out.DormantCount,
	}, nil
}
