package middleware

import (
	"context"
	"errors"
	"testing"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gdb"
	gredis "github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/consts"
	"ecboot/internal/errcode"
)

func init() {
	// 确定性测试配置（compose 基线, 与 service 层同源）
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			{Link: "mysql:myuser:secret@tcp(127.0.0.1:13306)/mydatabase"},
		},
	})
	gredis.SetConfig(&gredis.Config{Address: "127.0.0.1:6379", Db: 0})
}

func errCodeOf(err error) int {
	var ge *gerror.Error
	if errors.As(err, &ge) {
		return ge.Code().Code()
	}
	return -1
}

func seedMWAdmin(ctx context.Context, t *gtest.T, username string, isSuper int) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM `admin_user` WHERE username=?", username)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO `admin_user`(username,password_hash,real_name,is_super,status) VALUES(?,?,?, ?,1)",
		username, "$2a$10$AX.WGxKFiEoaIxKDzcbdGOIwtl555zjc1zLPFLDoz8mHu80DW1xg2", username, isSuper)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func cleanupMWAdmin(ctx context.Context, t *gtest.T, username string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM `admin_user` WHERE username=?", username)
}

// TestRequirePerm 权限拦截语义（SC-003, research D2）: 超管放行 / 无权 10005 / 未登录 10003。
func TestRequirePerm(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		superId := seedMWAdmin(ctx, t, "t_mw_super", 1)
		plainId := seedMWAdmin(ctx, t, "t_mw_plain", 0)
		defer cleanupMWAdmin(ctx, t, "t_mw_super")
		defer cleanupMWAdmin(ctx, t, "t_mw_plain")

		// 超管放行
		superCtx := context.WithValue(ctx, consts.CtxUserId, superId)
		t.AssertNil(RequirePerm(superCtx, "system:role:manage"))

		// 无权 → 10005
		plainCtx := context.WithValue(ctx, consts.CtxUserId, plainId)
		err := RequirePerm(plainCtx, "system:role:manage")
		t.Assert(errCodeOf(err), errcode.CodeForbidden)

		// 未登录（无主体）→ 10003
		err = RequirePerm(ctx, "system:role:manage")
		t.Assert(errCodeOf(err), errcode.CodeUnauthorized)
	})
}
