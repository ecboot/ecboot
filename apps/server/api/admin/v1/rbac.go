package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 角色 CRUD
	AdminRoleListReq struct {
		g.Meta `path:"/roles" method:"GET" summary:"角色列表"`
		PageReq
	}
	AdminRoleItem struct {
		Id          string `json:"id"`
		Name        string `json:"name" dc:"角色名称"`
		Code        string `json:"code" dc:"角色编码"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}
	AdminRoleListRes struct {
		PageRes
		List []AdminRoleItem `json:"list"`
	}

	AdminRoleCreateReq struct {
		g.Meta      `path:"/roles" method:"POST" summary:"新增角色"`
		Name        string `json:"name" v:"required" dc:"名称"`
		Code        string `json:"code" v:"required" dc:"编码"`
		Description string `json:"description" dc:"描述"`
	}
	AdminRoleCreateRes struct {
		Id string `json:"id"`
	}

	AdminRoleUpdateReq struct {
		g.Meta      `path:"/roles/{id}" method:"PUT" summary:"修改角色"`
		Id          string `json:"id" v:"required" dc:"角色ID"`
		Name        string `json:"name" dc:"名称"`
		Description string `json:"description" dc:"描述"`
		Status      int    `json:"status" dc:"状态"`
	}
	AdminRoleUpdateRes struct {
		Success bool `json:"success"`
	}

	AdminRoleDeleteReq struct {
		g.Meta `path:"/roles/{id}" method:"DELETE" summary:"删除角色(软删)"`
		Id     string `json:"id" v:"required" dc:"角色ID"`
	}
	AdminRoleDeleteRes struct {
		Success bool `json:"success"`
	}

	// 权限树（菜单/按钮/接口统一树）
	AdminPermissionTreeReq struct {
		g.Meta `path:"/permissions" method:"GET" summary:"权限树"`
	}
	AdminPermissionNode struct {
		Id       string                `json:"id"`
		ParentId string                `json:"parentId" dc:"父ID,0为根"`
		Name     string                `json:"name"`
		Code     string                `json:"code" dc:"权限编码"`
		Type     int                   `json:"type" dc:"1菜单 2按钮 3接口"`
		Sort     int                   `json:"sort"`
		Status   int                   `json:"status"`
		Children []AdminPermissionNode `json:"children"`
	}
	AdminPermissionTreeRes struct {
		Tree []AdminPermissionNode `json:"tree"`
	}

	// 角色-权限全量替换
	AdminRoleAssignPermReq struct {
		g.Meta        `path:"/roles/{id}/permissions" method:"PUT" summary:"角色权限分配"`
		Id            string   `json:"id" v:"required" dc:"角色ID"`
		PermissionIds []string `json:"permissionIds" dc:"权限ID集合(全量替换)"`
	}
	AdminRoleAssignPermRes struct {
		Success bool `json:"success"`
	}

	// 后台账号列表
	AdminUserListReq struct {
		g.Meta  `path:"/admin-users" method:"GET" summary:"后台账号列表"`
		Status  int    `json:"status" dc:"状态筛选"`
		Keyword string `json:"keyword" dc:"用户名/姓名"`
		PageReq
	}
	AdminUserItem struct {
		Id            string   `json:"id"`
		Username      string   `json:"username"`
		RealName      string   `json:"realName"`
		IsSuper       bool     `json:"isSuper"`
		Roles         []string `json:"roles" dc:"角色编码"`
		Status        int      `json:"status"`
		LastLoginTime string   `json:"lastLoginTime"`
	}
	AdminUserListRes struct {
		PageRes
		List []AdminUserItem `json:"list"`
	}

	// 创建后台账号
	AdminUserCreateReq struct {
		g.Meta   `path:"/admin-users" method:"POST" summary:"创建后台账号"`
		Username string `json:"username" v:"required" dc:"登录名"`
		Password string `json:"password" v:"required" dc:"初始密码"`
		RealName string `json:"realName" dc:"姓名"`
	}
	AdminUserCreateRes struct {
		Id string `json:"id"`
	}

	AdminUserUpdateReq struct {
		g.Meta   `path:"/admin-users/{id}" method:"PUT" summary:"修改后台账号"`
		Id       string `json:"id" v:"required" dc:"账号ID"`
		RealName string `json:"realName" dc:"姓名"`
		Status   int    `json:"status" dc:"1正常 2禁用"`
	}
	AdminUserUpdateRes struct {
		Success bool `json:"success"`
	}

	AdminUserDeleteReq struct {
		g.Meta `path:"/admin-users/{id}" method:"DELETE" summary:"删除后台账号(软删)"`
		Id     string `json:"id" v:"required" dc:"账号ID"`
	}
	AdminUserDeleteRes struct {
		Success bool `json:"success"`
	}

	// 账号-角色分配
	AdminUserAssignRolesReq struct {
		g.Meta  `path:"/admin-users/{id}/roles" method:"PUT" summary:"账号角色分配"`
		Id      string   `json:"id" v:"required" dc:"账号ID"`
		RoleIds []string `json:"roleIds" dc:"角色ID集合(全量替换)"`
	}
	AdminUserAssignRolesRes struct {
		Success bool `json:"success"`
	}
)
