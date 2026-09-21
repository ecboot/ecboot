package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// PointLogList 积分流水（类型筛选 + 分页）
func (c *ControllerV1) PointLogList(ctx context.Context, req *v1.PointLogListReq) (res *v1.PointLogListRes, err error) {
	out, err := user.PointLogs(ctx, middleware.CtxUserIdFrom(ctx), req.BizType,
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.PointLogListRes{List: make([]v1.PointLogItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.PointLogItem{
			BizType: it.BizType, Points: it.Points, BalanceAfter: it.BalanceAfter,
			OrderNo: it.OrderNo, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
