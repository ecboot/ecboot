package main

import (
	"github.com/gogf/gf/v2/os/gctx"

	cmd "ecboot/internal/app"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
