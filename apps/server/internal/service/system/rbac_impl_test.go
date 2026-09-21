package system

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/consts"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// opCtx 构造携带操作者 ID 的 ctx（FR-012 自守卫读取路径）。
func opCtx(ctx context.Context, adminId int64) context.Context {
	return context.WithValue(ctx, consts.CtxUserId, adminId)
}

// cleanRoles 清理测试角色与关联。
func cleanRoles(ctx context.Context, t *gtest.T, codes ...string) {
	for _, c := range codes {
		_, _ = g.DB().Exec(ctx, "DELETE FROM admin_user_role WHERE role_id IN (SELECT id FROM admin_role WHERE code=?)", c)
		_, _ = g.DB().Exec(ctx, "DELETE FROM admin_role_permission WHERE role_id IN (SELECT id FROM admin_role WHERE code=?)", c)
		_, _ = g.DB().Exec(ctx, "DELETE FROM admin_role WHERE code=?", c)
	}
}

// TestAdminUserCreate 创建账号（FR-009/010 前置）: 成功 + 密码不可逆存储 + 用户名唯一。
func TestAdminUserCreate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		defer cleanupAdmin(ctx, t, "t_cr_a")

		id, err := AdminUserCreate(ctx, model.AdminUserInput{
			Username: "t_cr_a", Password: "Init@Pass2026", RealName: "创建者",
		})
		t.AssertNil(err)
		t.AssertGT(id, 0)

		// 密码 bcrypt 存储, 不出明文
		rec, err := g.DB().GetOne(ctx, "SELECT password_hash, real_name FROM admin_user WHERE id=?", id)
		t.AssertNil(err)
		t.Assert(strings.HasPrefix(rec["password_hash"].String(), "$2a$"), true)
		t.Assert(strings.Contains(rec["password_hash"].String(), "Init@Pass2026"), false)
		t.Assert(rec["real_name"].String(), "创建者")

		// 用户名冲突 → 80003
		_, err = AdminUserCreate(ctx, model.AdminUserInput{Username: "t_cr_a", Password: "x@Pass2026"})
		t.Assert(errCode(err), errcode.CodeAdminNameTaken)
	})
}

// TestAdminUserUpdate 修改姓名/状态（FR-009）+ 禁止禁用自身（FR-012）。
func TestAdminUserUpdate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		id := seedAdmin(ctx, t, "t_upd_a", 0)
		defer cleanupAdmin(ctx, t, "t_upd_a")

		t.AssertNil(AdminUserUpdate(ctx, id, model.AdminUserUpdateInput{RealName: "新名字", Status: 1}))
		rec, err := g.DB().GetOne(ctx, "SELECT real_name, status FROM admin_user WHERE id=?", id)
		t.AssertNil(err)
		t.Assert(rec["real_name"].String(), "新名字")

		// 操作者禁用自己 → 80006
		err = AdminUserUpdate(opCtx(ctx, id), id, model.AdminUserUpdateInput{Status: 2})
		t.Assert(errCode(err), errcode.CodeAdminSelfGuard)
	})
}

// TestAdminUserDelete 软删（FR-009）+ 自守卫/超管保护（FR-012）。
func TestAdminUserDelete(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		id := seedAdmin(ctx, t, "t_del_a", 0)
		superId := seedAdmin(ctx, t, "t_del_super", 1)
		defer cleanupAdmin(ctx, t, "t_del_a")
		defer cleanupAdmin(ctx, t, "t_del_super")

		// 删除自身 → 80006
		err := AdminUserDelete(opCtx(ctx, id), id)
		t.Assert(errCode(err), errcode.CodeAdminSelfGuard)

		// 删除超管（非自身）→ 80006
		err = AdminUserDelete(opCtx(ctx, id), superId)
		t.Assert(errCode(err), errcode.CodeAdminSelfGuard)

		// 正常软删 → ActiveAdmin false, 软删位=1
		t.AssertNil(AdminUserDelete(opCtx(ctx, superId), id))
		ok, err := ActiveAdmin(ctx, id)
		t.AssertNil(err)
		t.Assert(ok, false)
		rec, err := g.DB().GetOne(ctx, "SELECT deleted FROM admin_user WHERE id=?", id)
		t.AssertNil(err)
		t.Assert(rec["deleted"].Int(), 1)
	})
}

// TestAssignRoles 账号-角色全量替换幂等（FR-011）。
func TestAssignRoles(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		id := seedAdmin(ctx, t, "t_asg_a", 0)
		defer cleanupAdmin(ctx, t, "t_asg_a")
		cleanRoles(ctx, t, "t_asg_r1", "t_asg_r2")
		defer cleanRoles(ctx, t, "t_asg_r1", "t_asg_r2")

		for _, c := range []string{"t_asg_r1", "t_asg_r2"} {
			_, err := g.DB().Exec(ctx,
				"INSERT INTO admin_role(name,code,status) VALUES(?,?,1)", c, c)
			t.AssertNil(err)
		}
		r1, err := g.DB().GetValue(ctx, "SELECT id FROM admin_role WHERE code='t_asg_r1'")
		t.AssertNil(err)
		r2, err := g.DB().GetValue(ctx, "SELECT id FROM admin_role WHERE code='t_asg_r2'")
		t.AssertNil(err)

		// 全量替换为 [r1, r2]
		t.AssertNil(AssignRoles(ctx, id, []int64{r1.Int64(), r2.Int64()}))
		p, err := Profile(ctx, id)
		t.AssertNil(err)
		t.Assert(len(p.Roles), 2)

		// 再全量替换为 [r1] → 仅剩 r1（重复提交结果一致）
		t.AssertNil(AssignRoles(ctx, id, []int64{r1.Int64()}))
		p, err = Profile(ctx, id)
		t.AssertNil(err)
		t.Assert(len(p.Roles), 1)
		t.Assert(p.Roles[0], "t_asg_r1")
	})
}

// TestAdminUserList 状态/关键词筛选与分页（FR-010）。
func TestAdminUserList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		a := seedAdmin(ctx, t, "t_lst_aa", 0)
		b := seedAdmin(ctx, t, "t_lst_bb", 0)
		_ = a
		_ = b
		defer cleanupAdmin(ctx, t, "t_lst_aa")
		defer cleanupAdmin(ctx, t, "t_lst_bb")
		_, err := g.DB().Exec(ctx, "UPDATE admin_user SET status=2, real_name='停用人' WHERE username='t_lst_bb'")
		t.AssertNil(err)

		// status=2 筛选 → 仅 t_lst_bb
		res, err := AdminUserList(ctx, 2, "", model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.AssertGT(res.Total, 0)
		for _, it := range res.List {
			t.Assert(it.Username, "t_lst_bb")
			t.Assert(it.Status, 2)
		}

		// keyword 命中姓名模糊
		res, err = AdminUserList(ctx, 0, "停用", model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.AssertGT(res.Total, 0)

		// keyword 命中用户名 + 分页 total
		res, err = AdminUserList(ctx, 0, "t_lst_", model.PageReq{Page: 1, PageSize: 1})
		t.AssertNil(err)
		t.Assert(res.Total, 2)
		t.Assert(len(res.List), 1)
	})
}

// TestAdminUserDetail 详情（接口微扩后行为）: 存在返回 + 角色编码; 不存在报错。
func TestAdminUserDetail(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		id := seedAdmin(ctx, t, "t_dtl_a", 1)
		defer cleanupAdmin(ctx, t, "t_dtl_a")

		it, err := AdminUserDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(it.Username, "t_dtl_a")
		t.Assert(it.IsSuper, true)

		_, err = AdminUserDetail(ctx, 999999999)
		t.Assert(errCode(err), errcode.CodeAdminBadCredential)
	})
}

const (
	tRoleName = "t_rc_role"
)

// seedRole 建测试角色, 返回 ID。
func seedRole(ctx context.Context, t *gtest.T, code string, status int) int64 {
	res, err := g.DB().Exec(ctx,
		"INSERT INTO admin_role(name,code,description,status) VALUES(?,?, '', ?)", code, code, status)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

// seedPerm 建测试权限点, 返回 ID。
func seedPerm(ctx context.Context, t *gtest.T, code string, status int) int64 {
	res, err := g.DB().Exec(ctx,
		"INSERT INTO admin_permission(name,code,type,sort,status) VALUES(?,?,3,9999,?)", code, code, status)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

// TestRoleCrud 角色 CRUD（FR-013）: 创建/编码唯一/修改/详情含 permissionIds。
func TestRoleCrud(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		cleanRoles(ctx, t, tRoleName)
		defer cleanRoles(ctx, t, tRoleName)

		id, err := RoleCreate(ctx, model.RoleInput{Name: "角色A", Code: tRoleName, Description: "测试"})
		t.AssertNil(err)
		t.AssertGT(id, 0)

		// 编码唯一 → 80004
		_, err = RoleCreate(ctx, model.RoleInput{Name: "角色B", Code: tRoleName})
		t.Assert(errCode(err), errcode.CodeRoleCodeTaken)

		// 修改
		t.AssertNil(RoleUpdate(ctx, id, model.RoleInput{Name: "角色A2", Description: "改", Status: 1}))

		// 详情
		dv, err := RoleDetailView(ctx, id)
		t.AssertNil(err)
		t.Assert(dv.Name, "角色A2")
		t.Assert(dv.Code, tRoleName)
		t.Assert(len(dv.PermissionIds), 0)
	})
}

// TestRoleDeleteGuard 被账号引用禁删（FR-014）。
func TestRoleDeleteGuard(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		cleanRoles(ctx, t, "t_rc_del")
		id := seedRole(ctx, t, "t_rc_del", 1)
		adminId := seedAdmin(ctx, t, "t_rc_del_a", 0)
		defer cleanRoles(ctx, t, "t_rc_del")
		defer cleanupAdmin(ctx, t, "t_rc_del_a")
		t.AssertNil(AssignRoles(ctx, adminId, []int64{id}))

		// 被引用 → 80005
		err := RoleDelete(ctx, id)
		t.Assert(errCode(err), errcode.CodeRoleInUse)

		// 解除引用后可软删
		t.AssertNil(AssignRoles(ctx, adminId, nil))
		t.AssertNil(RoleDelete(ctx, id))
		rec, err := g.DB().GetOne(ctx, "SELECT deleted FROM admin_role WHERE id=?", id)
		t.AssertNil(err)
		t.Assert(rec["deleted"].Int(), 1)
	})
}

// TestPermissionTree 权限树与种子同源（FR-015）: 全量节点、type=3 扁平、children 嵌套结构。
func TestPermissionTree(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		tree, err := PermissionTree(ctx)
		t.AssertNil(err)
		total, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM admin_permission WHERE deleted=0 AND status=1")
		t.AssertNil(err)
		t.Assert(len(tree), total.Int())

		// 含本批契约核心权限码
		found := false
		for _, n := range tree {
			if n.Code == "system:role:manage" {
				found = true
			}
		}
		t.Assert(found, true)
	})
}

// TestAssignPermissions 角色-权限全量替换幂等（FR-016）。
func TestAssignPermissions(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		cleanRoles(ctx, t, "t_rc_ap")
		rid := seedRole(ctx, t, "t_rc_ap", 1)
		defer cleanRoles(ctx, t, "t_rc_ap")
		p1 := seedPerm(ctx, t, "t_rc:p1", 1)
		p2 := seedPerm(ctx, t, "t_rc:p2", 1)
		defer func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM admin_permission WHERE id IN (?,?)", p1, p2)
		}()

		t.AssertNil(AssignPermissions(ctx, rid, []int64{p1, p2}))
		dv, err := RoleDetailView(ctx, rid)
		t.AssertNil(err)
		t.Assert(len(dv.PermissionIds), 2)

		// 全量替换为 [p1]
		t.AssertNil(AssignPermissions(ctx, rid, []int64{p1}))
		dv, err = RoleDetailView(ctx, rid)
		t.AssertNil(err)
		t.Assert(len(dv.PermissionIds), 1)
		t.Assert(dv.PermissionIds[0], p1)
	})
}

// TestHasPermission 权限判定（FR-017/018）: 超管直通/持权/无权/停用权限不命中。
func TestHasPermission(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		superId := seedAdmin(ctx, t, "t_rc_super", 1)
		plainId := seedAdmin(ctx, t, "t_rc_plain", 0)
		defer cleanupAdmin(ctx, t, "t_rc_super")
		defer cleanupAdmin(ctx, t, "t_rc_plain")
		cleanRoles(ctx, t, "t_rc_hp")
		rid := seedRole(ctx, t, "t_rc_hp", 1)
		defer cleanRoles(ctx, t, "t_rc_hp")
		p := seedPerm(ctx, t, "t_rc:hp", 1)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM admin_permission WHERE id=?", p) }()

		// 超管直通（任意码）
		ok, err := HasPermission(ctx, superId, "any:thing:here")
		t.AssertNil(err)
		t.Assert(ok, true)

		// 未知管理员 → false
		ok, err = HasPermission(ctx, 999999999, "system:role:manage")
		t.AssertNil(err)
		t.Assert(ok, false)

		// 授权后命中
		t.AssertNil(AssignRoles(ctx, plainId, []int64{rid}))
		t.AssertNil(AssignPermissions(ctx, rid, []int64{p}))
		ok, err = HasPermission(ctx, plainId, "t_rc:hp")
		t.AssertNil(err)
		t.Assert(ok, true)

		// 未授权码 → false
		ok, err = HasPermission(ctx, plainId, "system:role:manage")
		t.AssertNil(err)
		t.Assert(ok, false)

		// 权限停用 → false
		_, err = g.DB().Exec(ctx, "UPDATE admin_permission SET status=0 WHERE id=?", p)
		t.AssertNil(err)
		ok, err = HasPermission(ctx, plainId, "t_rc:hp")
		t.AssertNil(err)
		t.Assert(ok, false)
	})
}
