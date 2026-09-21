package admin

import (
	"context"
	"strconv"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/system"
)

// AdminUserList 后台账号列表
func (c *ControllerV1) AdminUserList(ctx context.Context, req *v1.AdminUserListReq) (res *v1.AdminUserListRes, err error) {
	out, err := system.AdminUserList(ctx, req.Status, req.Keyword, model.PageReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminUserListRes{List: make([]v1.AdminUserItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminUserItem{
			Id:            strconv.FormatInt(it.Id, 10),
			Username:      it.Username,
			RealName:      it.RealName,
			IsSuper:       it.IsSuper,
			Roles:         it.Roles,
			Status:        it.Status,
			LastLoginTime: it.LastLoginTime,
		})
	}
	return res, nil
}
