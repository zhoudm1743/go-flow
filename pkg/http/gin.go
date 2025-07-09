package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-frame/pkg/config"
	"github.com/zhoudm1743/go-frame/pkg/http/middleware"
	ctx "github.com/zhoudm1743/go-frame/pkg/http/unified"
	"github.com/zhoudm1743/go-frame/pkg/log"
	"go.uber.org/fx"
)

// EngineParams HTTP引擎参数
type EngineParams struct {
	fx.In
	Config *config.Config
	Logger log.Logger
}

// NewGinEngine 创建Gin引擎
func NewGinEngine(p EngineParams) *gin.Engine {
	// 设置运行模式
	switch p.Config.App.Mode {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 将Gin的默认日志输出重定向到logrus
	middleware.GinLogToLogrus(p.Logger)

	engine := gin.New()

	// 配置中间件
	engine.Use(
		// 恢复中间件
		gin.Recovery(),
		// 使用美观的logrus日志中间件
		middleware.LogrusLogger(p.Logger),
	)

	return engine
}

// ServerParams HTTP服务参数
type ServerParams struct {
	fx.In
	Config *config.Config
	Engine *gin.Engine
	Logger log.Logger
	// 修改：使用tag指定注入的路由注册器
	Routers []RouterRegister `group:"http_routers"`
}

// NewHTTPServer 创建HTTP服务
func NewHTTPServer(p ServerParams) *http.Server {
	// 注册所有路由
	for _, router := range p.Routers {
		// 创建统一路由器
		unifiedRouter := ctx.NewRouter(ctx.GinEngine, p.Engine, nil)
		router.RegisterRoutes(unifiedRouter)
	}

	// 创建服务器
	addr := fmt.Sprintf("%s:%d", p.Config.HTTP.Host, p.Config.HTTP.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        p.Engine,
		ReadTimeout:    p.Config.HTTP.ReadTimeout,
		WriteTimeout:   p.Config.HTTP.WriteTimeout,
		MaxHeaderBytes: p.Config.HTTP.MaxHeaderBytes,
	}

	return srv
}

// StartHTTPServer 启动HTTP服务
func StartHTTPServer(lc fx.Lifecycle, server *http.Server, logger log.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 非阻塞方式启动服务器
			go func() {
				logger.Infof("HTTP服务启动在 %s", server.Addr)
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Errorf("HTTP服务启动失败: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("正在关闭HTTP服务...")

			// 创建一个用于关闭的上下文，设置超时时间
			stopCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := server.Shutdown(stopCtx); err != nil {
				logger.Errorf("HTTP服务关闭出错: %v", err)
				return err
			}

			logger.Info("HTTP服务已成功关闭")
			return nil
		},
	})
}

// Module 提供HTTP模块
var Module = fx.Options(
	fx.Provide(NewGinEngine),
	// 修改：直接使用ServerParams中的Routers
	fx.Provide(NewHTTPServer),
	fx.Invoke(StartHTTPServer),
)
