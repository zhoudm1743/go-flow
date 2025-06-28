package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-flow/core/config"
	"github.com/zhoudm1743/go-flow/core/logger"
)

// RouteRegistrator 路由注册器接口
type RouteRegistrator interface {
	RegisterRoutes(engine *gin.Engine) error
}

// Service HTTP 服务
type Service struct {
	server       *http.Server
	engine       *gin.Engine
	config       *config.Config
	logger       logger.Logger
	registrators []RouteRegistrator // 路由注册器列表
}

// NewService 创建新的 HTTP 服务实例
func NewService(cfg *config.Config, log logger.Logger, registrators []RouteRegistrator) *Service {
	// 根据环境设置 Gin 模式
	switch cfg.App.Env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 创建 Gin 引擎
	engine := gin.New()

	// 设置 Gin 的日志输出为我们的 logger
	gin.DefaultWriter = NewGinLoggerWriter(log)
	gin.DefaultErrorWriter = NewGinLoggerWriter(log)

	// 添加中间件
	engine.Use(GinLoggerMiddleware(log))
	engine.Use(GinRecoveryMiddleware(log))

	// 根据配置添加 CORS 中间件
	if cfg.App.HTTP.EnableCORS {
		engine.Use(CORSMiddleware())
	}

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      engine,
		ReadTimeout:  time.Duration(cfg.App.HTTP.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.App.HTTP.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.App.HTTP.IdleTimeout) * time.Second,
	}

	service := &Service{
		server:       server,
		engine:       engine,
		config:       cfg,
		logger:       log,
		registrators: registrators,
	}

	// 初始化全局路由注册器
	InitGlobalRegistrar(engine)

	// 注册路由
	service.registerRoutes()

	return service
}

// registerRoutes 注册路由
func (s *Service) registerRoutes() {
	// 健康检查路由
	s.engine.GET("/health", s.healthCheck)
	s.engine.GET("/ping", s.ping)

	// API 路由组
	api := s.engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// 示例路由
			v1.GET("/status", s.getStatus)
		}
	}

	// 注册所有模块的路由
	for _, registrator := range s.registrators {
		if err := registrator.RegisterRoutes(s.engine); err != nil {
			s.logger.Errorf("路由注册失败: %v", err)
		}
	}
}

// Start 启动 HTTP 服务
func (s *Service) Start(ctx context.Context) error {
	s.logger.Infof("HTTP 服务正在启动，监听端口: %d", s.config.App.Port)

	// 在 goroutine 中启动服务器
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Errorf("HTTP 服务启动失败: %v", err)
		}
	}()

	s.logger.Infof("HTTP 服务已启动，地址: http://localhost:%d", s.config.App.Port)
	return nil
}

// Stop 停止 HTTP 服务
func (s *Service) Stop(ctx context.Context) error {
	s.logger.Info("正在停止 HTTP 服务...")

	// 创建一个带超时的上下文
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		s.logger.Errorf("HTTP 服务停止失败: %v", err)
		return err
	}

	s.logger.Info("HTTP 服务已停止")
	return nil
}

// GetEngine 获取 Gin 引擎（用于外部注册路由）
func (s *Service) GetEngine() *gin.Engine {
	return s.engine
}

// healthCheck 健康检查处理器
func (s *Service) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   s.config.App.Name,
		"version":   s.config.App.Version,
	})
}

// ping 简单的 ping 处理器
func (s *Service) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// getStatus 获取状态信息
func (s *Service) getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app": gin.H{
			"name":    s.config.App.Name,
			"version": s.config.App.Version,
			"env":     s.config.App.Env,
		},
		"server": gin.H{
			"port":   s.config.App.Port,
			"uptime": time.Now().Format(time.RFC3339),
		},
	})
}
