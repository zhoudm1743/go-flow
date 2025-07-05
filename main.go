package main

import (
	"os"

	"github.com/zhoudm1743/go-flow/cmd"
	"github.com/zhoudm1743/go-flow/internal/demo"
	"github.com/zhoudm1743/go-flow/internal/shop"
	"github.com/zhoudm1743/go-flow/internal/test"
	"github.com/zhoudm1743/go-flow/internal/user"
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

	// 添加模块
	app.AddModule(demo.NewModuleWithName("demo")) // 使用/demo前缀
	app.AddModule(user.NewModuleWithName("v1"))   // 使用/v1前缀注册user模块
	app.AddModule(test.NewModule())               // 添加新生成的test模块
	app.AddModule(shop.NewModule())               // 添加新生成的shop模块

	// 启动应用
	app.Run()
}
