package admin

import (
	"context"
	"strconv"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/system"
)

// AdminUserCreate 创建后台账号
func (c *ControllerV1) AdminUserCreate(ctx context.Context, req *v1.AdminUserCreateReq) (res *v1.AdminUserCreateRes, err error) {
	id, err := system.AdminUserCreate(ctx, model.AdminUserInput{
		Username: req.Username,
		Password: req.Password,
		RealName: req.RealName,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminUserCreateRes{Id: strconv.FormatInt(id, 10)}, nil
}
