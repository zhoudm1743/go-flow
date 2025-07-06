package main

import (
	"os"

	"github.com/zhoudm1743/go-flow/cmd"
	"github.com/zhoudm1743/go-flow/pkg/core"
	"github.com/zhoudm1743/go-flow/pkg/http"
)

func main() {
	// 检查是否是命令行模式
	if len(os.Args) > 1 && os.Args[1] == "gen" {
		cmd.Execute()
		return
	}

	// 创建应用
	app := core.NewApp("go-flow")

	// 添加HTTP模块
	app.WithOptions(http.UnifiedModule)

	// 启动应用
	app.Run()
}
