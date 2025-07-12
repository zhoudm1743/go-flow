package controller

import (
	"go.uber.org/fx"
)

// Router 路由模块选项
var Router = fx.Options(
	DemoModule,
	// 可以添加更多模块
)
