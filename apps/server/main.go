package main

import (
	"time"

	"github.com/gogf/gf/v2/os/gctx"

	cmd "ecboot/internal/app"
)

func main() {
	// 时区口径统一（009 评审 C2）: 库内 DATETIME 存 UTC 墙钟, 而 gf 读取时按 time.Local 解析——
	// 二者必须一致, 否则强类型时间列读回偏移（读早 8h, 回显再提交每轮再漂 8h）。
	// 应用进程固定 UTC: 对外 RFC3339 出参带 Z, 前端按本地时区渲染; 库会话保持 UTC（默认）。
	time.Local = time.UTC
	cmd.Main.Run(gctx.GetInitCtx())
}
