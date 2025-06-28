package routes

import (
	"github.com/gin-gonic/gin"
	testRoutes "github.com/zhoudm1743/go-flow/app/admin/routes/test"
	httpCore "github.com/zhoudm1743/go-flow/core/http"
	"go.uber.org/fx"
)

// RouteRegistratorFunc 函数类型实现RouteRegistrator接口
type RouteRegistratorFunc func(*gin.Engine) error

func (f RouteRegistratorFunc) RegisterRoutes(engine *gin.Engine) error {
	return f(engine)
}

// RouteRegistratorResult 路由注册器结果结构
type RouteRegistratorResult struct {
	fx.Out
	Registrator httpCore.RouteRegistrator `group:"route_registrators"`
}

// NewAdminRouteRegistrator 创建admin路由注册器 - 简化为一个函数
func NewAdminRouteRegistrator(group httpCore.Group) RouteRegistratorResult {
	return RouteRegistratorResult{
		Registrator: RouteRegistratorFunc(func(engine *gin.Engine) error {
			return httpCore.RegisterModuleRoutes(engine, "admin", []httpCore.Group{
				group, // 🎉 终极简化！
			})
		}),
	}
}

// Module FX模块定义
var Module = fx.Options(
	fx.Provide(testRoutes.NewTestGroup),
	fx.Provide(NewAdminRouteRegistrator),
)
