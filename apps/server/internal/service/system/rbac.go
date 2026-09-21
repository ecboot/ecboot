// Package system 平台治理域——RBAC（V10 五表）/ 审计（V10）/ 系统配置（V30）/ 看板 / 通知模板（V26）。
// 归属说明: 这些是平台级能力（不属 user/shop 业务域）; 管理渠道（api/admin）为主要消费方。
package system

import "context"

// IRBACLogic RBAC 权限（权限树: 菜单/按钮/接口统一, code=模块:资源:动作）。
type IRBACLogic interface {
	RoleList(ctx context.Context, page PageQuery) (*PageResult[RoleItem], error)
	RoleCreate(ctx context.Context, in RoleInput) (int64, error)
	RoleUpdate(ctx context.Context, id int64, in RoleInput) error
	RoleDelete(ctx context.Context, id int64) error // 有账号引用禁删
	RoleDetail(ctx context.Context, id int64) (*RoleDetail, error)
	// PermissionTree 权限树（与 internal/consts/permission.go 及 000032 种子同源）。
	PermissionTree(ctx context.Context) ([]PermissionNode, error)
	// AssignPermissions 角色-权限全量替换（同事务删旧插新）。
	AssignPermissions(ctx context.Context, roleId int64, permissionIds []int64) error
	AdminUserList(ctx context.Context, status int, keyword string, page PageQuery) (*PageResult[AdminUserItem], error)
	// AdminUserCreate 创建后台账号（密码 bcrypt; 审计留痕）。
	AdminUserCreate(ctx context.Context, in AdminUserInput) (int64, error)
	AdminUserUpdate(ctx context.Context, id int64, in AdminUserUpdateInput) error
	AdminUserDelete(ctx context.Context, id int64) error
	AssignRoles(ctx context.Context, adminId int64, roleIds []int64) error
	// HasPermission 权限校验（is_super 直通; 路由中间件调用）。
	HasPermission(ctx context.Context, adminId int64, permissionCode string) (bool, error)
}

type RoleItem struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

type RoleInput struct {
	Name        string
	Code        string
	Description string
	Status      int
}

type RoleDetail struct {
	Id            int64   `json:"id"`
	Name          string  `json:"name"`
	Code          string  `json:"code"`
	Description   string  `json:"description"`
	Status        int     `json:"status"`
	PermissionIds []int64 `json:"permissionIds"`
}

type PermissionNode struct {
	Id       int64            `json:"id"`
	ParentId int64            `json:"parentId"`
	Name     string           `json:"name"`
	Code     string           `json:"code"`
	Type     int              `json:"type" dc:"1菜单 2按钮 3接口"`
	Sort     int              `json:"sort"`
	Status   int              `json:"status"`
	Children []PermissionNode `json:"children"`
}

type AdminUserItem struct {
	Id            int64    `json:"id"`
	Username      string   `json:"username"`
	RealName      string   `json:"realName"`
	IsSuper       bool     `json:"isSuper"`
	Roles         []string `json:"roles"`
	Status        int      `json:"status"`
	LastLoginTime string   `json:"lastLoginTime"`
}

type AdminUserInput struct {
	Username string
	Password string
	RealName string
}

type AdminUserUpdateInput struct {
	RealName string
	Status   int
}

// IAdminAuthLogic 后台账号认证（登录审计含失败尝试）。
type IAdminAuthLogic interface {
	Login(ctx context.Context, username, password, captchaKey, captchaCode, ip, userAgent string) (*AdminLoginResult, error)
	// ChangePassword 改密（旧密码校验）。
	ChangePassword(ctx context.Context, adminId int64, oldPassword, newPassword string) error
}

type AdminLoginResult struct {
	AdminId      int64  `json:"adminId"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	RealName     string `json:"realName"`
	IsSuper      bool   `json:"isSuper"`
}
