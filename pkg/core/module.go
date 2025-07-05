package core

import (
	"go.uber.org/fx"
)

// Module 模块接口
type Module interface {
	// Name 模块名称
	Name() string

	// Options 模块配置选项
	Options() fx.Option

	// RoutePrefix 模块路由前缀，可选实现
	// 如果返回空字符串则不使用前缀
	RoutePrefix() string
}

// RegisterRoutes 注册路由接口
// 应用模块如果需要注册路由，需要实现此接口
type RegisterRoutes interface {
	// Routes 获取路由
	Routes() interface{}
}
