// auth_impl.go 后台账号认证实现（接口契约见 rbac.go IAdminAuthLogic）。
// 规则: 登录成功/失败均写 admin_login_log（只追加, research D8 不锁定）;
// 密码 bcrypt（irreversible, FR-007）; 会话走 security 库 admin 渠道（research D1）。
package system

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/internal/dao"
)

// ActiveAdmin 账号存在且启用且未软删（FR-008: 中间件每请求校验, 禁用/软删后会话立即不可用）。
func ActiveAdmin(ctx context.Context, adminId int64) (bool, error) {
	cnt, err := dao.AdminUser.Ctx(ctx).
		Where(dao.AdminUser.Columns().Id, adminId).
		Where(dao.AdminUser.Columns().Status, 1).
		Where(dao.AdminUser.Columns().Deleted, 0).
		Count()
	if err != nil {
		return false, gerror.Wrap(err, "查询后台账号失败")
	}
	return cnt > 0, nil
}
