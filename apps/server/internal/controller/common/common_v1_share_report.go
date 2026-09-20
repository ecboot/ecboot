package common

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/common/v1"
)

// ShareReport 分享行为上报（归因来源记录; 游客可报, 携带会员凭证则归因到人）
func (c *ControllerV1) ShareReport(ctx context.Context, req *v1.ShareReportReq) (res *v1.ShareReportRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
