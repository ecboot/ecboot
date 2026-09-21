package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// LoginLogList 近 30 天登录记录（IP 脱敏）
func (c *ControllerV1) LoginLogList(ctx context.Context, req *v1.LoginLogListReq) (res *v1.LoginLogListRes, err error) {
	out, err := user.LoginLogs(ctx, middleware.CtxUserIdFrom(ctx),
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.LoginLogListRes{List: make([]v1.LoginLogItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.LoginLogItem{
			Channel: it.Channel, Status: it.Status, Ip: it.Ip,
			UserAgent: it.UserAgent, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
