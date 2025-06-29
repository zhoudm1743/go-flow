package captcha

import (
	"github.com/zhoudm1743/go-flow/core/cache"
	"go.uber.org/fx"
)

// ProvideCaptchaService 提供验证码服务的fx构造函数
func ProvideCaptchaService(cacheClient cache.Cache) *Service {
	return NewService(cacheClient, DefaultConfig)
}

var Module = fx.Module("captcha",
	fx.Provide(ProvideCaptchaService),
)
