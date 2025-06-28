package database

import (
	"context"

	"go.uber.org/fx"
)

// Module fx模块
var Module = fx.Options(
	fx.Provide(NewDatabase),
	fx.Provide(NewMigrator),
	fx.Invoke(func(migrator *Migrator, lifecycle fx.Lifecycle) {
		lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return migrator.AutoMigrate()
			},
		})
	}),
)
