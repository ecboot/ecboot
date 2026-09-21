// devtools_impl.go 开发调试端点服务（research D4）。
// 规则: fail-closed——仅 ECBOOT_MOCK=true（开发/测试态）开放读取模拟短信;
// 生产环境默认拒绝且不泄露任何内容（FR-023）。
package system

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/errcode"
)

// MockLatestSms 返回最近一条模拟短信验证码（library/sms MockSender 落点: Redis mock:sms:{phone}）。
func MockLatestSms(ctx context.Context, phone string) (string, error) {
	if os.Getenv("ECBOOT_MOCK") != "true" {
		return "", errcode.New(errcode.CodeProdDenied, "调试端点仅限非生产环境")
	}
	v, err := g.Redis().Do(ctx, "GET", "mock:sms:"+phone)
	if err != nil {
		return "", gerror.Wrap(err, "查询模拟短信失败")
	}
	if v == nil || v.String() == "" {
		return "", errcode.New(errcode.CodeNotFound, "暂无该手机号的模拟短信")
	}
	return v.String(), nil
}
