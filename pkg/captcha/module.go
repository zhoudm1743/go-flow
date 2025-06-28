package captcha

import "go.uber.org/fx"

var Module = fx.Module("captcha",
	fx.Provide(NewService),
)
