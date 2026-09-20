package app

import (
	"context"

	_ "ecboot/internal/bootstrap"
	"ecboot/internal/routes"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http ecboot",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			routes.RouterGroup(s)
			s.Run()
			return nil
		},
	}
)
