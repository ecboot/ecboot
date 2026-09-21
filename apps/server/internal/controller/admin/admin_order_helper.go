package admin

import (
	"context"
	"fmt"

	"ecboot/internal/middleware"
)

// adminOperator 操作者标识（表注释约定 admin:{id}）。
func adminOperator(ctx context.Context) string {
	return fmt.Sprintf("admin:%d", middleware.CtxUserIdFrom(ctx))
}
