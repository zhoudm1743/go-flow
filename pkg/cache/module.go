package cache

import (
	"github.com/zhoudm1743/go-frame/pkg/config"
	"go.uber.org/fx"
)

// Module 缓存模块
var Module = fx.Options(
	fx.Provide(
		NewCacheProvider,
	),
)

// NewCacheProvider 根据配置选择并提供缓存实现
func NewCacheProvider(cfg *config.Config) fx.Option {
	switch cfg.Cache.Type {
	case "memory":
		return fx.Provide(NewMemoryCache)
	case "file":
		return fx.Provide(NewFileCache)
	default:
		return fx.Provide(NewRedisCache)
	}
}
