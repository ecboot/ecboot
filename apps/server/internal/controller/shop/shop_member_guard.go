package shop

import (
	"context"

	"ecboot/internal/errcode"
	"ecboot/internal/middleware"
)

// requireMember 会员端点未登录防御（评审 I8）: 会话 userId<=0 一律 10003。
// 背景: OrderDetail/Cancel 以 userId=0 表示"后台视角"（管理面复用）——若某天 C 端端点
// 漏挂鉴权或白名单误放宽, 缺此防御会静默变成"任意订单可见/可取消"。此处显式收口。
func requireMember(ctx context.Context) (int64, error) {
	userId := middleware.CtxUserIdFrom(ctx)
	if userId <= 0 {
		return 0, errcode.New(errcode.CodeUnauthorized, "未登录或凭证失效")
	}
	return userId, nil
}
