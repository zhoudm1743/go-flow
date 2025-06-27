package http

import (
	"context"

	"go.uber.org/fx"
)

// Module HTTP 服务的 fx 模块
var Module = fx.Options(
	fx.Provide(NewService),
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
