package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/zhoudm1743/go-flow/pkg/config"
	"github.com/zhoudm1743/go-flow/pkg/http/middleware"
	ctx "github.com/zhoudm1743/go-flow/pkg/http/unified"
	"github.com/zhoudm1743/go-flow/pkg/log"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"go.uber.org/fx"
)

// Server 统一的HTTP服务器接口
type Server interface {
	// 启动服务器
	Start() error

	// 关闭服务器
	Shutdown(ctx context.Context) error

	// 获取路由器
	Router() ctx.Router

	// 添加全局中间件
	Use(middlewares ...ctx.MiddlewareFunc) Server

	// 注册路由
	RegisterRoutes(register RouterRegister) Server

	// 设置错误处理器
	SetErrorHandler(handler ctx.HandlerFunc) Server

	// 设置404处理器
	SetNotFoundHandler(handler ctx.HandlerFunc) Server

	// 设置405处理器
	SetMethodNotAllowedHandler(handler ctx.HandlerFunc) Server
}

// RouterRegister 路由注册接口
type RouterRegister interface {
	// RegisterRoutes 注册路由到路由器
	RegisterRoutes(router ctx.Router)
}

// ServerConfig 服务器配置
type ServerConfig struct {
	// 服务器地址
	Addr string

	// 引擎类型: "gin" 或 "fiber"
	Engine string

	// 是否启用调试模式
	Debug bool

	// 读超时
	ReadTimeout time.Duration

	// 写超时
	WriteTimeout time.Duration

	// 请求体大小限制
	BodyLimit int

	// 是否启用CORS
	EnableCORS bool

	// 是否启用请求日志
	EnableRequestLog bool

	// 是否启用恢复中间件
	EnableRecover bool
}

// UnifiedServer 统一的HTTP服务器实现
type UnifiedServer struct {
	config     *ServerConfig
	router     ctx.Router
	ginEngine  *gin.Engine
	ginServer  *http.Server
	fiberApp   *fiber.App
	logger     log.Logger
	middleware []ctx.MiddlewareFunc
}

// NewUnifiedServer 创建统一的HTTP服务器
func NewUnifiedServer(config *ServerConfig, logger log.Logger) *UnifiedServer {
	server := &UnifiedServer{
		config:     config,
		logger:     logger,
		middleware: []ctx.MiddlewareFunc{},
	}

	// 根据配置创建引擎
	switch config.Engine {
	case "fiber":
		server.initFiber()
	default:
		server.initGin()
	}

	return server
}

// 初始化Gin引擎
func (s *UnifiedServer) initGin() {
	// 设置模式
	if s.config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	s.ginEngine = gin.New()

	// 使用恢复中间件
	if s.config.EnableRecover {
		s.ginEngine.Use(gin.Recovery())
	}

	// 使用日志中间件
	if s.config.EnableRequestLog {
		s.ginEngine.Use(middleware.LogrusLogger(s.logger))
	}

	// 设置404处理器
	s.ginEngine.NoRoute(response.NoRoute)

	// 设置405处理器
	s.ginEngine.NoMethod(response.NoMethod)

	// 创建HTTP服务器
	s.ginServer = &http.Server{
		Addr:           s.config.Addr,
		Handler:        s.ginEngine,
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// 创建统一路由器
	s.router = ctx.NewRouter(ctx.GinEngine, s.ginEngine, nil)

	// 记录启动信息
	s.logger.Info("[GIN] 服务器初始化完成")
}

// 初始化Fiber引擎
func (s *UnifiedServer) initFiber() {
	// 创建Fiber应用
	s.fiberApp = fiber.New(fiber.Config{
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		BodyLimit:    s.config.BodyLimit,
		// 将Fiber的日志输出重定向到logrus
		DisableStartupMessage: true, // 禁用默认的启动消息，由logrus处理
	})

	// 使用恢复中间件
	if s.config.EnableRecover {
		s.fiberApp.Use(recover.New())
	}

	// 使用日志中间件
	if s.config.EnableRequestLog {
		s.fiberApp.Use(middleware.FiberLogrusLogger(s.logger))
	}

	// 设置404处理器 - 使用标准方式
	s.fiberApp.Use(func(c *fiber.Ctx) error {
		// 检查路由是否存在
		// 注意：这是一个简化的实现，可能不能完全检测所有404情况
		// 但对于大多数情况应该足够了
		if c.Route() == nil {
			return response.FiberNoRoute(c)
		}
		return c.Next()
	})

	// 创建统一路由器
	s.router = ctx.NewRouter(ctx.FiberEngine, nil, s.fiberApp)

	// 记录启动信息
	s.logger.Info("[FIBER] 服务器初始化完成")
}

// Router 实现Server接口
func (s *UnifiedServer) Router() ctx.Router {
	return s.router
}

// Use 实现Server接口
func (s *UnifiedServer) Use(middlewares ...ctx.MiddlewareFunc) Server {
	s.middleware = append(s.middleware, middlewares...)

	// 应用中间件到路由器
	s.router.Use(middlewares...)

	return s
}

// RegisterRoutes 实现Server接口
func (s *UnifiedServer) RegisterRoutes(register RouterRegister) Server {
	register.RegisterRoutes(s.router)
	return s
}

// SetErrorHandler 实现Server接口
func (s *UnifiedServer) SetErrorHandler(handler ctx.HandlerFunc) Server {
	switch s.config.Engine {
	case "fiber":
		if s.fiberApp != nil {
			// 创建一个新的Fiber应用，保留原有配置但更新错误处理器
			config := s.fiberApp.Config()
			config.ErrorHandler = func(c *fiber.Ctx, err error) error {
				ctx := ctx.NewFiberContext(c)
				return handler(ctx)
			}

			// 由于无法直接修改Config，我们只能记录这个情况
			s.logger.Warn("Fiber不支持动态修改错误处理器，请在创建应用时设置")
		}
	default:
		// Gin不支持全局错误处理器，需要在每个路由中处理
	}
	return s
}

// SetNotFoundHandler 实现Server接口
func (s *UnifiedServer) SetNotFoundHandler(handler ctx.HandlerFunc) Server {
	switch s.config.Engine {
	case "fiber":
		if s.fiberApp != nil {
			s.fiberApp.Use(func(c *fiber.Ctx) error {
				// 检查路由是否存在
				if c.Route() == nil {
					ctx := ctx.NewFiberContext(c)
					return handler(ctx)
				}
				return c.Next()
			})
		}
	default:
		if s.ginEngine != nil {
			s.ginEngine.NoRoute(func(c *gin.Context) {
				ctx := ctx.NewGinContext(c)
				if err := handler(ctx); err != nil {
					c.Error(err)
				}
			})
		}
	}
	return s
}

// SetMethodNotAllowedHandler 实现Server接口
func (s *UnifiedServer) SetMethodNotAllowedHandler(handler ctx.HandlerFunc) Server {
	switch s.config.Engine {
	case "fiber":
		// Fiber不直接支持方法不允许处理器，需要自定义
	default:
		if s.ginEngine != nil {
			s.ginEngine.NoMethod(func(c *gin.Context) {
				ctx := ctx.NewGinContext(c)
				if err := handler(ctx); err != nil {
					c.Error(err)
				}
			})
		}
	}
	return s
}

// Start 实现Server接口
func (s *UnifiedServer) Start() error {
	s.logger.Infof("HTTP服务启动在 %s", s.config.Addr)

	// 根据引擎类型启动服务器
	switch s.config.Engine {
	case "fiber":
		return s.fiberApp.Listen(s.config.Addr)
	default:
		return s.ginServer.ListenAndServe()
	}
}

// Shutdown 实现Server接口
func (s *UnifiedServer) Shutdown(ctx context.Context) error {
	s.logger.Info("正在关闭HTTP服务...")

	// 根据引擎类型关闭服务器
	var err error
	switch s.config.Engine {
	case "fiber":
		err = s.fiberApp.Shutdown()
	default:
		err = s.ginServer.Shutdown(ctx)
	}

	if err != nil {
		s.logger.Errorf("HTTP服务关闭出错: %v", err)
		return err
	}

	s.logger.Info("HTTP服务已成功关闭")
	return nil
}

// UnifiedServerParams 统一服务器参数
type UnifiedServerParams struct {
	fx.In
	Config *config.Config
	Logger log.Logger
}

// NewUnifiedHTTPServer 创建统一的HTTP服务器
func NewUnifiedHTTPServer(p UnifiedServerParams) Server {
	// 从配置文件获取HTTP引擎类型
	engineType := p.Config.HTTP.Engine
	if engineType == "" {
		engineType = "gin" // 默认使用Gin
	}

	// 创建服务器配置
	serverConfig := &ServerConfig{
		Engine:           engineType,
		Addr:             fmt.Sprintf("%s:%d", p.Config.HTTP.Host, p.Config.HTTP.Port),
		Debug:            p.Config.App.Mode == "dev",
		ReadTimeout:      p.Config.HTTP.ReadTimeout,
		WriteTimeout:     p.Config.HTTP.WriteTimeout,
		BodyLimit:        p.Config.HTTP.MaxBodySize,
		EnableCORS:       true, // 默认启用CORS
		EnableRequestLog: true,
		EnableRecover:    true,
	}

	// 创建服务器
	return NewUnifiedServer(serverConfig, p.Logger)
}

// StartUnifiedHTTPServer 启动统一的HTTP服务器
func StartUnifiedHTTPServer(lc fx.Lifecycle, server Server, logger log.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 非阻塞方式启动服务器
			go func() {
				if err := server.Start(); err != nil {
					logger.Errorf("HTTP服务启动失败: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// 创建一个用于关闭的上下文，设置超时时间
			stopCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			return server.Shutdown(stopCtx)
		},
	})
}
