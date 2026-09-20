package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"

	"ecboot/api/user/v1"
	"ecboot/internal/dao"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// Me 当前用户信息
func (c *ControllerV1) Me(ctx context.Context, req *v1.MeReq) (res *v1.MeRes, err error) {
	userId := middleware.CtxUserIdFrom(ctx)
	if userId == 0 {
		return nil, gerror.NewCode(gcode.New(10003, "未登录或凭证失效", nil))
	}
	record, err := dao.User.Ctx(ctx).Where(dao.User.Columns().Id, userId).One()
	if err != nil || record.IsEmpty() {
		return nil, gerror.NewCode(gcode.New(10003, "未登录或凭证失效", nil))
	}
	// 密文解密后脱敏（后台展示与接口契约: 138****1234）
	phone := ""
	if pc := user.PhoneCipherFor(ctx); pc != nil {
		if plain, dErr := pc.Decrypt(record["phone"].String()); dErr == nil && len(plain) >= 7 {
			phone = plain[:3] + "****" + plain[len(plain)-4:]
		}
	}
	return &v1.MeRes{
		UserId:          gconv.String(userId),
		Nickname:        record["nickname"].String(),
		Avatar:          record["avatar"].String(),
		Phone:           phone,
		Level:           record["level"].Int(),
		RegisterChannel: record["register_channel"].Int(),
	}, nil
}
