package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/system"
)

// AdminRiskRecordList 风控事件列表。
func (c *ControllerV1) AdminRiskRecordList(ctx context.Context, req *v1.AdminRiskRecordListReq) (res *v1.AdminRiskRecordListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermRiskRecordRead); err != nil {
		return nil, err
	}
	var uid int64
	if req.UserId != "" {
		if uid, err = parseID(req.UserId); err != nil {
			return nil, err
		}
	}
	out, err := system.NewRiskAdminLogic().AdminRecordList(ctx, uid, req.AppealStatus, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminRiskRecordListRes{List: make([]v1.AdminRiskRecordItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminRiskRecordItem{
			Id: fmtID(it.Id), UserId: fmtID(it.UserId), RuleName: it.RuleName,
			ObjectType: it.ObjectType, ObjectNo: it.ObjectNo, Action: it.Action,
			AppealStatus: it.AppealStatus, Remark: it.Remark, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
