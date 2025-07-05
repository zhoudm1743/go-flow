package cache

import (
	"go.uber.org/fx"
)

// Module 缓存模块
var Module = fx.Options(
	fx.Provide(
		// 使用内存缓存实现，而不是真实的Redis
		// 在生产环境中可以替换为 NewRedisCache
		// NewMockCache,
		NewRedisCache,
	),
)
