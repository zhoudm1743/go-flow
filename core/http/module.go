package http

import (
	"context"

	"github.com/zhoudm1743/go-flow/core/config"
	"github.com/zhoudm1743/go-flow/core/logger"
	"go.uber.org/fx"
)

// RouteRegistratorParams 路由注册器参数结构
type RouteRegistratorParams struct {
	fx.In
	Registrators []RouteRegistrator `group:"route_registrators"`
}

// Module HTTP 服务的 fx 模块
var Module = fx.Options(
	fx.Provide(func(cfg *config.Config, log logger.Logger, params RouteRegistratorParams) *Service {
		return NewService(cfg, log, params.Registrators)
	}),
	fx.Invoke(registerLifecycle),
)

// registerLifecycle 注册 HTTP 服务的生命周期钩子
func registerLifecycle(lc fx.Lifecycle, service *Service) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return service.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return service.Stop(ctx)
		},
	})
}
