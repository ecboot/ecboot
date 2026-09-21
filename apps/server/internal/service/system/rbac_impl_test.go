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
