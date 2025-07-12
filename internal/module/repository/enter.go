package repository

import (
	"go.uber.org/fx"
)

var Repository = fx.Options(
	fx.Provide(
		NewDemoRepository,
	),
)
