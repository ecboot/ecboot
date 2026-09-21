// superadmin.go 种子超管过渡（特性 005 / research D4）：
// 启动时 admin_user 表无任何账号且 env ADMIN_USERNAME/ADMIN_PASSWORD 均配置时，
// 创建 is_super=1 种子账号（bcrypt）；env 缺省则跳过并输出警告（治理特性接入后由后台创建）。
package bootstrap

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	golangbcrypt "golang.org/x/crypto/bcrypt"
)

func EnsureSuperAdmin(ctx context.Context) {
	count, err := g.DB().Model("admin_user").Ctx(ctx).Count()
	if err != nil || count > 0 {
		return
	}
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	if username == "" || password == "" {
		g.Log().Warning(ctx, "admin_user 为空且未配置 ADMIN_USERNAME/ADMIN_PASSWORD——后台登录不可用")
		return
	}
	hash, err := golangbcrypt.GenerateFromPassword([]byte(password), golangbcrypt.DefaultCost)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	_, err = g.DB().Model("admin_user").Ctx(ctx).Data(g.Map{
		"username":      username,
		"password_hash": string(hash),
		"is_super":      1,
	}).Insert()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	g.Log().Infof(ctx, "seed super admin created: %s", username)
}
