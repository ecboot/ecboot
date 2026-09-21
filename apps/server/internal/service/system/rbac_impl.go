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
