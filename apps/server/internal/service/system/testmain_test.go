package system

import (
	"context"
	"os"
	"time"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gdb"
	gredis "github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/testutil"
)

func init() {
	// 时区口径统一（009 评审 C2）: 与库内 UTC 墙钟一致（见 main.go）
	time.Local = time.UTC
	// 测试显式开启 mock（fail-closed 语义下的白盒开关）
	_ = os.Setenv("ECBOOT_MOCK", "true")
	// 确定性测试配置（compose 基线——与 user/shop 域测试同源, 不依赖本地 manifest 差异）
	_ = gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			{
				Link: testutil.DSN(),
			},
		},
	})
	gredis.SetConfig(&gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      0,
	})
}

// seedAdmin 建测试后台账号, 返回自增 ID（用后请 cleanupAdmin）。
func seedAdmin(ctx context.Context, t *gtest.T, username string, isSuper int) int64 {
	res, err := g.DB().Exec(ctx,
		"DELETE FROM `admin_user` WHERE username=?", username)
	t.AssertNil(err)
	_, _ = res.LastInsertId()
	res, err = g.DB().Exec(ctx,
		"INSERT INTO `admin_user`(username,password_hash,real_name,is_super,status) VALUES(?,?,?,?,1)",
		username, "$2a$10$AX.WGxKFiEoaIxKDzcbdGOIwtl555zjc1zLPFLDoz8mHu80DW1xg2", username, isSuper)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

// cleanupAdmin 清理测试后台账号及其登录审计行（含无账号的失败尝试行）。
func cleanupAdmin(ctx context.Context, t *gtest.T, username string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM `admin_user` WHERE username=?", username)
	_, _ = g.DB().Exec(ctx, "DELETE FROM `admin_login_log` WHERE username=?", username)
}
