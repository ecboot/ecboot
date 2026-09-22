package common

import (
	"context"

	"ecboot/api/common/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// ShareReport 分享行为上报（归因来源记录; 携带会员凭证则归因到人——017 批次 11 接线,
// 原桩未接导致归因链断头（评审 C2））。
func (c *ControllerV1) ShareReport(ctx context.Context, req *v1.ShareReportReq) (res *v1.ShareReportRes, err error) {
	var spuId int64
	if req.SpuId != "" {
		if spuId, err = parseID(req.SpuId); err != nil {
			return nil, err
		}
	}
	if err = user.NewDistributionLogic().ShareReport(ctx,
		middleware.CtxUserIdFrom(ctx), spuId, req.Channel, req.SceneValue); err != nil {
		return nil, err
	}
	return &v1.ShareReportRes{Success: true}, nil
}
