package service

import (
	"go.uber.org/fx"
)

var Service = fx.Options(
	fx.Provide(
		NewDemoService,
	),
)
