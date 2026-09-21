package system

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// TestActiveAdmin 后台账号态校验（spec FR-008）：禁用/软删/不存在 → false。
func TestActiveAdmin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		id := seedAdmin(ctx, t, "t_act_admin", 0)
		defer cleanupAdmin(ctx, t, "t_act_admin")

		ok, err := ActiveAdmin(ctx, id)
		t.AssertNil(err)
		t.Assert(ok, true)

		// 禁用 → false
		_, err = g.DB().Exec(ctx, "UPDATE `admin_user` SET status=2 WHERE id=?", id)
		t.AssertNil(err)
		ok, err = ActiveAdmin(ctx, id)
		t.AssertNil(err)
		t.Assert(ok, false)

		// 恢复启用后软删 → false
		_, err = g.DB().Exec(ctx, "UPDATE `admin_user` SET status=1, deleted=1 WHERE id=?", id)
		t.AssertNil(err)
		ok, err = ActiveAdmin(ctx, id)
		t.AssertNil(err)
		t.Assert(ok, false)

		// 不存在 → false
		ok, err = ActiveAdmin(ctx, 999999999)
		t.AssertNil(err)
		t.Assert(ok, false)
	})
}
