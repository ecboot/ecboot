// Package system 平台治理域——RBAC（V10 五表）/ 审计（V10）/ 系统配置（V30）/ 看板 / 通知模板（V26）。
// 归属说明: 这些是平台级能力（不属 user/shop 业务域）; 管理渠道（api/admin）为主要消费方。
package system

import (
	"context"

	"ecboot/internal/model"
)

// IRBACLogic RBAC 权限（权限树: 菜单/按钮/接口统一, code=模块:资源:动作）。
type IRBACLogic interface {
	RoleList(ctx context.Context, page model.PageReq) (*model.PageResult[model.RoleItem], error)
	RoleCreate(ctx context.Context, in model.RoleInput) (int64, error)
	RoleUpdate(ctx context.Context, id int64, in model.RoleInput) error
	RoleDelete(ctx context.Context, id int64) error // 有账号引用禁删
	RoleDetailView(ctx context.Context, id int64) (*model.RoleDetailView, error)
	// PermissionTree 权限树（与 internal/consts/permission.go 及 000032 种子同源）。
	PermissionTree(ctx context.Context) ([]model.PermissionNode, error)
	// AssignPermissions 角色-权限全量替换（同事务删旧插新）。
	AssignPermissions(ctx context.Context, roleId int64, permissionIds []int64) error
	AdminUserList(ctx context.Context, status int, keyword string, page model.PageReq) (*model.PageResult[model.AdminUserItem], error)
	// AdminUserCreate 创建后台账号（密码 bcrypt; 审计留痕）。
	AdminUserCreate(ctx context.Context, in model.AdminUserInput) (int64, error)
	AdminUserUpdate(ctx context.Context, id int64, in model.AdminUserUpdateInput) error
	AdminUserDelete(ctx context.Context, id int64) error
	AssignRoles(ctx context.Context, adminId int64, roleIds []int64) error
	// HasPermission 权限校验（is_super 直通; 路由中间件调用）。
	HasPermission(ctx context.Context, adminId int64, permissionCode string) (bool, error)
}

// IAdminAuthLogic 后台账号认证（登录审计含失败尝试）。
type IAdminAuthLogic interface {
	Login(ctx context.Context, username, password, captchaKey, captchaCode, ip, userAgent string) (*model.AdminLoginResult, error)
	// ChangePassword 改密（旧密码校验）。
	ChangePassword(ctx context.Context, adminId int64, oldPassword, newPassword string) error
}
