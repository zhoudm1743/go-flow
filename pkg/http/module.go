package http

import (
	"github.com/zhoudm1743/go-frame/pkg/config"
	"go.uber.org/fx"
)

// UnifiedModule 提供统一的HTTP模块
var UnifiedModule = fx.Options(
	fx.Provide(NewUnifiedHTTPServer),
	fx.Invoke(StartUnifiedHTTPServer),
)

// ConfigurableModule 创建可配置的HTTP模块
func ConfigurableModule(configFunc func(*config.Config)) fx.Option {
	return fx.Options(
		// 先配置
		fx.Invoke(func(config *config.Config) {
			if configFunc != nil {
				configFunc(config)
			}
		}),
		// 再创建服务器
		fx.Provide(NewUnifiedHTTPServer),
		fx.Invoke(StartUnifiedHTTPServer),
	)
}

// WithEngine 配置使用指定的HTTP引擎
func WithEngine(engine string) func(*config.Config) {
	return func(config *config.Config) {
		config.HTTP.Engine = engine
	}
}

// GinEngine 使用Gin引擎
var GinEngine = WithEngine("gin")

// FiberEngine 使用Fiber引擎
var FiberEngine = WithEngine("fiber")
