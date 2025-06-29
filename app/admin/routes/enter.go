package routes

import (
	"github.com/gin-gonic/gin"
	authRoutes "github.com/zhoudm1743/go-flow/app/admin/routes/auth"
	systemRoutes "github.com/zhoudm1743/go-flow/app/admin/routes/system"
	httpCore "github.com/zhoudm1743/go-flow/core/http"
	"go.uber.org/fx"
)

// RouteRegistratorFunc 函数类型实现RouteRegistrator接口
type RouteRegistratorFunc func(*gin.Engine) error

func (f RouteRegistratorFunc) RegisterRoutes(engine *gin.Engine) error {
	return f(engine)
}

// NewAdminRouteRegistrator 创建admin路由注册器，批量注册所有路由组
func NewAdminRouteRegistrator(params struct {
	fx.In
	Groups []httpCore.Group `group:"admin_route_groups"`
}) struct {
	fx.Out
	Registrator httpCore.RouteRegistrator `group:"route_registrators"`
} {
	return struct {
		fx.Out
		Registrator httpCore.RouteRegistrator `group:"route_registrators"`
	}{
		Registrator: RouteRegistratorFunc(func(engine *gin.Engine) error {
			return httpCore.RegisterModuleRoutes(engine, "admin", params.Groups)
		}),
	}
}

// Module FX模块定义
var Module = fx.Options(
	// 批量提供路由组
	fx.Provide(fx.Annotate(
		systemRoutes.NewAdminGroup,
		fx.ResultTags(`group:"admin_route_groups"`),
	)),
	fx.Provide(fx.Annotate(
		authRoutes.NewAuthGroup,
		fx.ResultTags(`group:"admin_route_groups"`),
	)),
	fx.Provide(NewAdminRouteRegistrator),
)
