// rbac_impl.go RBAC 实现——后台账号管理 + 角色权限（接口契约见 rbac.go IRBACLogic）。
// 规则: 用户名/角色编码唯一; 软删; 被引用角色禁删（FR-014）; 操作者自守卫与超管保护（FR-012）;
// 账号-角色/角色-权限分配为全量替换（同事务删旧插新, 幂等, FR-011/016）。
package system

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/crypto/bcrypt"

	"ecboot/internal/consts"
	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
	"ecboot/internal/model/entity"
)

// currentAdminId 从 ctx 读取操作者管理员 ID（0=未知; consts 键无环, 见 consts/ctx.go）。
func currentAdminId(ctx context.Context) int64 {
	if v, ok := ctx.Value(consts.CtxUserId).(int64); ok {
		return v
	}
	return 0
}

// AdminUserDetail 账号详情（FR-010 字段同列表单行）。
func AdminUserDetail(ctx context.Context, id int64) (*model.AdminUserItem, error) {
	rec, err := dao.AdminUser.Ctx(ctx).
		Where(dao.AdminUser.Columns().Id, id).
		Where(dao.AdminUser.Columns().Deleted, 0).
		One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询后台账号失败")
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeAdminBadCredential, "账号不存在")
	}
	var admin entity.AdminUser
	if err = rec.Struct(&admin); err != nil {
		return nil, gerror.Wrap(err, "解析后台账号失败")
	}
	roleByAdmin, err := adminRoleMap(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	roles := roleByAdmin[id]
	if roles == nil {
		roles = []string{}
	}
	return &model.AdminUserItem{
		Id:            int64(admin.Id),
		Username:      admin.Username,
		RealName:      admin.RealName,
		IsSuper:       admin.IsSuper == 1,
		Roles:         roles,
		Status:        admin.Status,
		LastLoginTime: rfc3339(admin.LastLoginTime),
	}, nil
}

// AdminUserCreate 创建后台账号（FR-009）: 用户名唯一 + bcrypt 初始密码。
func AdminUserCreate(ctx context.Context, in model.AdminUserInput) (int64, error) {
	cnt, err := dao.AdminUser.Ctx(ctx).
		Where(dao.AdminUser.Columns().Username, in.Username).
		Count()
	if err != nil {
		return 0, gerror.Wrap(err, "查询后台账号失败")
	}
	if cnt > 0 {
		return 0, errcode.New(errcode.CodeAdminNameTaken, "登录名已存在")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, gerror.Wrap(err, "生成密码哈希失败")
	}
	res, err := dao.AdminUser.Ctx(ctx).Data(do.AdminUser{
		Username:     in.Username,
		PasswordHash: string(hash),
		RealName:     in.RealName,
	}).Insert()
	if err != nil {
		return 0, gerror.Wrap(err, "创建后台账号失败")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "读取新账号ID失败")
	}
	return id, nil
}

// AdminUserUpdate 修改姓名/状态（FR-009）: 操作者禁用自己被拒（FR-012）。
func AdminUserUpdate(ctx context.Context, id int64, in model.AdminUserUpdateInput) error {
	if in.Status == 2 && currentAdminId(ctx) == id {
		return errcode.New(errcode.CodeAdminSelfGuard, "禁止禁用自身账号")
	}
	data := do.AdminUser{}
	if in.RealName != "" {
		data.RealName = in.RealName
	}
	if in.Status > 0 {
		data.Status = in.Status
	}
	if _, err := dao.AdminUser.Ctx(ctx).
		Where(dao.AdminUser.Columns().Id, id).
		Where(dao.AdminUser.Columns().Deleted, 0).
		Data(data).Update(); err != nil {
		return gerror.Wrap(err, "修改后台账号失败")
	}
	return nil
}

// AdminUserDelete 软删账号（FR-009）: 自身与超管保护（FR-012）。
func AdminUserDelete(ctx context.Context, id int64) error {
	if op := currentAdminId(ctx); op == id {
		return errcode.New(errcode.CodeAdminSelfGuard, "禁止删除自身账号")
	}
	rec, err := dao.AdminUser.Ctx(ctx).
		Where(dao.AdminUser.Columns().Id, id).
		Where(dao.AdminUser.Columns().Deleted, 0).
		One()
	if err != nil {
		return gerror.Wrap(err, "查询后台账号失败")
	}
	if rec.IsEmpty() {
		return errcode.New(errcode.CodeAdminBadCredential, "账号不存在")
	}
	var admin entity.AdminUser
	if err = rec.Struct(&admin); err != nil {
		return gerror.Wrap(err, "解析后台账号失败")
	}
	if admin.IsSuper == 1 {
		return errcode.New(errcode.CodeAdminSelfGuard, "禁止删除超管账号")
	}
	_, err = dao.AdminUser.Ctx(ctx).
		Where(dao.AdminUser.Columns().Id, id).
		Data(do.AdminUser{Deleted: 1, Status: 2}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除后台账号失败")
	}
	return nil
}

// AssignRoles 账号-角色全量替换（FR-011）: 同事务删旧插新, 幂等。
func AssignRoles(ctx context.Context, adminId int64, roleIds []int64) error {
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.AdminUserRole.Ctx(ctx).
			Where(dao.AdminUserRole.Columns().AdminId, adminId).
			Delete()
		if err != nil {
			return err
		}
		if len(roleIds) == 0 {
			return nil
		}
		rows := make([]do.AdminUserRole, 0, len(roleIds))
		for _, rid := range roleIds {
			rows = append(rows, do.AdminUserRole{AdminId: adminId, RoleId: rid})
		}
		_, err = dao.AdminUserRole.Ctx(ctx).Data(rows).Insert()
		return err
	})
	if err != nil {
		return gerror.Wrap(err, "分配角色失败")
	}
	return nil
}

// AdminUserList 账号列表（FR-010）: 状态/关键词筛选 + 分页 + 角色编码与最后登录时间。
func AdminUserList(ctx context.Context, status int, keyword string, page model.PageReq) (*model.PageResult[model.AdminUserItem], error) {
	page = page.Normalized()
	m := dao.AdminUser.Ctx(ctx).Where(dao.AdminUser.Columns().Deleted, 0)
	if status > 0 {
		m = m.Where(dao.AdminUser.Columns().Status, status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		u, n := dao.AdminUser.Columns().Username, dao.AdminUser.Columns().RealName
		m = m.Where(u+" LIKE ? OR "+n+" LIKE ?", like, like)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计后台账号失败")
	}
	recs, err := m.Page(page.Page, page.PageSize).OrderDesc(dao.AdminUser.Columns().Id).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询后台账号失败")
	}
	var admins []entity.AdminUser
	if err = recs.Structs(&admins); err != nil {
		return nil, gerror.Wrap(err, "解析后台账号失败")
	}

	// 批量装配角色编码: 本页账号 → 关联 → 角色编码映射
	ids := make([]int64, 0, len(admins))
	for _, a := range admins {
		ids = append(ids, int64(a.Id))
	}
	roleByAdmin, err := adminRoleMap(ctx, ids)
	if err != nil {
		return nil, err
	}
	list := make([]model.AdminUserItem, 0, len(admins))
	for _, a := range admins {
		item := model.AdminUserItem{
			Id:            int64(a.Id),
			Username:      a.Username,
			RealName:      a.RealName,
			IsSuper:       a.IsSuper == 1,
			Roles:         roleByAdmin[int64(a.Id)],
			Status:        a.Status,
			LastLoginTime: rfc3339(a.LastLoginTime),
		}
		if item.Roles == nil {
			item.Roles = []string{}
		}
		list = append(list, item)
	}
	return &model.PageResult[model.AdminUserItem]{List: list, Total: int64(total)}, nil
}

// adminRoleMap 批量查询一组账号的角色编码映射（启用未删角色）。
func adminRoleMap(ctx context.Context, adminIds []int64) (map[int64][]string, error) {
	out := make(map[int64][]string, len(adminIds))
	if len(adminIds) == 0 {
		return out, nil
	}
	links, err := dao.AdminUserRole.Ctx(ctx).
		WhereIn(dao.AdminUserRole.Columns().AdminId, adminIds).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询账号角色失败")
	}
	roleIds := make([]int64, 0, len(links))
	linksByRole := map[int64][]int64{} // roleId → adminIds
	for _, l := range links {
		rid := l["role_id"].Int64()
		aid := l["admin_id"].Int64()
		if _, seen := linksByRole[rid]; !seen {
			roleIds = append(roleIds, rid)
		}
		linksByRole[rid] = append(linksByRole[rid], aid)
	}
	if len(roleIds) == 0 {
		return out, nil
	}
	roleRecs, err := dao.AdminRole.Ctx(ctx).
		Fields(dao.AdminRole.Columns().Id, dao.AdminRole.Columns().Code).
		Where(dao.AdminRole.Columns().Deleted, 0).
		Where(dao.AdminRole.Columns().Status, 1).
		WhereIn(dao.AdminRole.Columns().Id, roleIds).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询角色编码失败")
	}
	for _, r := range roleRecs {
		code := r["code"].String()
		for _, aid := range linksByRole[r["id"].Int64()] {
			out[aid] = append(out[aid], code)
		}
	}
	return out, nil
}

// rfc3339 时间出参格式（契约约定 RFC3339; nil 时间返回空串）。
func rfc3339(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02T15:04:05Z07:00")
}

// RoleList 角色分页列表（FR-013）。
func RoleList(ctx context.Context, page model.PageReq) (*model.PageResult[model.RoleItem], error) {
	page = page.Normalized()
	m := dao.AdminRole.Ctx(ctx).Where(dao.AdminRole.Columns().Deleted, 0)
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计角色失败")
	}
	recs, err := m.Page(page.Page, page.PageSize).OrderAsc(dao.AdminRole.Columns().Id).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询角色失败")
	}
	var roles []entity.AdminRole
	if err = recs.Structs(&roles); err != nil {
		return nil, gerror.Wrap(err, "解析角色失败")
	}
	list := make([]model.RoleItem, 0, len(roles))
	for _, r := range roles {
		list = append(list, model.RoleItem{
			Id:          int64(r.Id),
			Name:        r.Name,
			Code:        r.Code,
			Description: r.Description,
			Status:      r.Status,
		})
	}
	return &model.PageResult[model.RoleItem]{List: list, Total: int64(total)}, nil
}

// RoleCreate 创建角色（FR-013）: 编码唯一。
func RoleCreate(ctx context.Context, in model.RoleInput) (int64, error) {
	cnt, err := dao.AdminRole.Ctx(ctx).
		Where(dao.AdminRole.Columns().Code, in.Code).
		Count()
	if err != nil {
		return 0, gerror.Wrap(err, "查询角色失败")
	}
	if cnt > 0 {
		return 0, errcode.New(errcode.CodeRoleCodeTaken, "角色编码已存在")
	}
	status := in.Status
	if status <= 0 {
		status = 1
	}
	res, err := dao.AdminRole.Ctx(ctx).Data(do.AdminRole{
		Name:        in.Name,
		Code:        in.Code,
		Description: in.Description,
		Status:      status,
	}).Insert()
	if err != nil {
		return 0, gerror.Wrap(err, "创建角色失败")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "读取新角色ID失败")
	}
	return id, nil
}

// RoleUpdate 修改角色（FR-013; 编码为身份不可改）。
func RoleUpdate(ctx context.Context, id int64, in model.RoleInput) error {
	data := do.AdminRole{}
	if in.Name != "" {
		data.Name = in.Name
	}
	if in.Description != "" {
		data.Description = in.Description
	}
	if in.Status > 0 {
		data.Status = in.Status
	}
	_, err := dao.AdminRole.Ctx(ctx).
		Where(dao.AdminRole.Columns().Id, id).
		Where(dao.AdminRole.Columns().Deleted, 0).
		Data(data).Update()
	if err != nil {
		return gerror.Wrap(err, "修改角色失败")
	}
	return nil
}

// RoleDelete 软删角色（FR-013/014）: 被账号引用禁删。
func RoleDelete(ctx context.Context, id int64) error {
	cnt, err := dao.AdminUserRole.Ctx(ctx).
		Where(dao.AdminUserRole.Columns().RoleId, id).
		Count()
	if err != nil {
		return gerror.Wrap(err, "查询角色引用失败")
	}
	if cnt > 0 {
		return errcode.New(errcode.CodeRoleInUse, "角色被账号引用, 禁止删除")
	}
	_, err = dao.AdminRole.Ctx(ctx).
		Where(dao.AdminRole.Columns().Id, id).
		Where(dao.AdminRole.Columns().Deleted, 0).
		Data(do.AdminRole{Deleted: 1, Status: 0}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除角色失败")
	}
	return nil
}

// RoleDetailView 角色详情（FR-013）: 含已分配权限 ID 集合。
func RoleDetailView(ctx context.Context, id int64) (*model.RoleDetailView, error) {
	rec, err := dao.AdminRole.Ctx(ctx).
		Where(dao.AdminRole.Columns().Id, id).
		Where(dao.AdminRole.Columns().Deleted, 0).
		One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询角色失败")
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeActivityNotFound, "角色不存在")
	}
	var role entity.AdminRole
	if err = rec.Struct(&role); err != nil {
		return nil, gerror.Wrap(err, "解析角色失败")
	}
	permIds, err := dao.AdminRolePermission.Ctx(ctx).
		Fields(dao.AdminRolePermission.Columns().PermissionId).
		Where(dao.AdminRolePermission.Columns().RoleId, id).
		Array()
	if err != nil {
		return nil, gerror.Wrap(err, "查询角色权限失败")
	}
	ids := make([]int64, 0, len(permIds))
	for _, v := range permIds {
		ids = append(ids, v.Int64())
	}
	return &model.RoleDetailView{
		Id:            int64(role.Id),
		Name:          role.Name,
		Code:          role.Code,
		Description:   role.Description,
		Status:        role.Status,
		PermissionIds: ids,
	}, nil
}

// AssignPermissions 角色-权限全量替换（FR-016）: 同事务删旧插新, 幂等。
func AssignPermissions(ctx context.Context, roleId int64, permissionIds []int64) error {
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.AdminRolePermission.Ctx(ctx).
			Where(dao.AdminRolePermission.Columns().RoleId, roleId).
			Delete()
		if err != nil {
			return err
		}
		if len(permissionIds) == 0 {
			return nil
		}
		rows := make([]do.AdminRolePermission, 0, len(permissionIds))
		for _, pid := range permissionIds {
			rows = append(rows, do.AdminRolePermission{RoleId: roleId, PermissionId: pid})
		}
		_, err = dao.AdminRolePermission.Ctx(ctx).Data(rows).Insert()
		return err
	})
	if err != nil {
		return gerror.Wrap(err, "分配角色权限失败")
	}
	return nil
}

// PermissionTree 权限树（FR-015）: 菜单/按钮/接口统一树, 与 000032 种子同源;
// 当前种子为 type=3 扁平登记（parent=0）, 树构建按 parent_id 通用嵌套。
func PermissionTree(ctx context.Context) ([]model.PermissionNode, error) {
	recs, err := dao.AdminPermission.Ctx(ctx).
		Where(dao.AdminPermission.Columns().Deleted, 0).
		Where(dao.AdminPermission.Columns().Status, 1).
		OrderAsc(dao.AdminPermission.Columns().Sort).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询权限失败")
	}
	var perms []entity.AdminPermission
	if err = recs.Structs(&perms); err != nil {
		return nil, gerror.Wrap(err, "解析权限失败")
	}
	byParent := map[int64][]model.PermissionNode{}
	for _, p := range perms {
		n := model.PermissionNode{
			Id:       int64(p.Id),
			ParentId: int64(p.ParentId),
			Name:     p.Name,
			Code:     p.Code,
			Type:     p.Type,
			Sort:     p.Sort,
			Status:   p.Status,
			Children: []model.PermissionNode{},
		}
		byParent[n.ParentId] = append(byParent[n.ParentId], n)
	}
	// 父子挂接（深度按种子层级, 当前为 1~2 层）
	var build func(parentId int64) []model.PermissionNode
	build = func(parentId int64) []model.PermissionNode {
		out := byParent[parentId]
		for i := range out {
			out[i].Children = build(out[i].Id)
		}
		return out
	}
	return build(0), nil
}

// HasPermission 权限判定（FR-017/018）: is_super 直通;
// 数据流 admin_user_role → admin_role_permission → admin_permission(code, 启用未删)。
func HasPermission(ctx context.Context, adminId int64, code string) (bool, error) {
	rec, err := dao.AdminUser.Ctx(ctx).
		Fields(dao.AdminUser.Columns().IsSuper).
		Where(dao.AdminUser.Columns().Id, adminId).
		Where(dao.AdminUser.Columns().Deleted, 0).
		One()
	if err != nil {
		return false, gerror.Wrap(err, "查询后台账号失败")
	}
	if rec.IsEmpty() {
		return false, nil
	}
	if rec["is_super"].Int() == 1 {
		return true, nil
	}
	cnt, err := dao.AdminUserRole.Ctx(ctx).As("ur").
		InnerJoin(dao.AdminRolePermission.Table()+" rp", "rp.role_id=ur.role_id").
		InnerJoin(dao.AdminPermission.Table()+" p", "p.id=rp.permission_id").
		Where("ur.admin_id", adminId).
		Where("p.code", code).
		Where("p.status", 1).
		Where("p.deleted", 0).
		Count()
	if err != nil {
		return false, gerror.Wrap(err, "权限判定查询失败")
	}
	return cnt > 0, nil
}
