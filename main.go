package main

import (
	"github.com/gogf/gf/v2/os/gctx"

	"friberg/internal/cmd"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
