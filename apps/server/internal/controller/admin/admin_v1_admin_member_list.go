package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// AdminMemberList 会员列表（手机号精确/状态/关键词筛选, 脱敏展示）。
func (c *ControllerV1) AdminMemberList(ctx context.Context, req *v1.AdminMemberListReq) (res *v1.AdminMemberListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermMemberRead); err != nil {
		return nil, err
	}
	out, err := user.NewMemberAdminLogic().AdminList(ctx, req.Phone, req.Status, req.Keyword, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminMemberListRes{List: make([]v1.AdminMemberItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminMemberItem{
			UserId: fmtID(it.UserId), Nickname: it.Nickname, Phone: it.Phone,
			Level: int(it.Level), Status: it.Status, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
