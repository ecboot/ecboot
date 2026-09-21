package admin

import (
	"context"
	"strconv"

	"ecboot/api/admin/v1"
	"ecboot/internal/service/system"
)

// AdminUserDetail 后台账号详情
func (c *ControllerV1) AdminUserDetail(ctx context.Context, req *v1.AdminUserDetailReq) (res *v1.AdminUserDetailRes, err error) {
	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	it, err := system.AdminUserDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.AdminUserDetailRes{
		Id:            strconv.FormatInt(it.Id, 10),
		Username:      it.Username,
		RealName:      it.RealName,
		IsSuper:       it.IsSuper,
		Roles:         it.Roles,
		Status:        it.Status,
		LastLoginTime: it.LastLoginTime,
	}, nil
}
