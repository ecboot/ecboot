package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// AdminDistributorList 推广员列表。
func (c *ControllerV1) AdminDistributorList(ctx context.Context, req *v1.AdminDistributorListReq) (res *v1.AdminDistributorListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDistributionRead); err != nil {
		return nil, err
	}
	out, err := user.NewDistributionAdminLogic().AdminDistributorList(ctx, req.Status, req.Keyword, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminDistributorListRes{List: make([]v1.AdminDistributorItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminDistributorItem{
			Id: fmtID(it.Id), UserId: fmtID(it.UserId), Nickname: it.Nickname,
			Level: 0, Status: it.Status, ApplyTime: it.ApplyTime, AuditTime: it.AuditTime,
		})
	}
	return res, nil
}
