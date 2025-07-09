package http

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/zhoudm1743/go-frame/pkg/config"
	"github.com/zhoudm1743/go-frame/pkg/log"
	"go.uber.org/fx"
)

// FiberRouterRegister Fiber路由注册接口
type FiberRouterRegister interface {
	Register(app *fiber.App)
}

// FiberEngineParams Fiber引擎参数
type FiberEngineParams struct {
	fx.In
	Config *config.Config
	Logger log.Logger
}

// NewFiberApp 创建Fiber应用
func NewFiberApp(p FiberEngineParams) *fiber.App {
	// 创建Fiber配置
	config := fiber.Config{
		AppName:      p.Config.App.Name,
		ReadTimeout:  p.Config.HTTP.ReadTimeout,
		WriteTimeout: p.Config.HTTP.WriteTimeout,
	}

	// 根据运行模式配置
	switch p.Config.App.Mode {
	case "prod":
		// 生产环境配置
		config.Prefork = true
	case "test":
		// 测试环境配置
		config.Prefork = false
	default:
		// 开发环境配置
		config.Prefork = false
	}

	// 创建Fiber应用
	app := fiber.New(config)

	// 添加中间件
	app.Use(
		// 恢复中间件
		recover.New(),
		// 日志中间件
		logger.New(logger.Config{
			Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
			TimeFormat: "2006-01-02 15:04:05",
			Output:     NewFiberLoggerOutput(p.Logger),
		}),
	)

	return app
}

// FiberLoggerOutput 日志输出适配器
type FiberLoggerOutput struct {
	logger log.Logger
}

// NewFiberLoggerOutput 创建日志输出适配器
func NewFiberLoggerOutput(logger log.Logger) *FiberLoggerOutput {
	return &FiberLoggerOutput{logger: logger}
}

// Write 实现io.Writer接口
func (o *FiberLoggerOutput) Write(p []byte) (n int, err error) {
	o.logger.Info(string(p))
	return len(p), nil
}

// FiberServerParams Fiber服务参数
type FiberServerParams struct {
	fx.In
	Config  *config.Config
	App     *fiber.App
	Logger  log.Logger
	Routers []FiberRouterRegister `group:"fiber_routers"`
}

// NewFiberServer 创建Fiber服务
func NewFiberServer(p FiberServerParams) *fiber.App {
	// 注册所有路由
	for _, router := range p.Routers {
		router.Register(p.App)
	}

	return p.App
}

// StartFiberServer 启动Fiber服务
func StartFiberServer(lc fx.Lifecycle, app *fiber.App, config *config.Config, logger log.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 非阻塞方式启动服务器
			go func() {
				addr := fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port)
				logger.Infof("Fiber HTTP服务启动在 %s", addr)
				if err := app.Listen(addr); err != nil {
					logger.Errorf("Fiber HTTP服务启动失败: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("正在关闭Fiber HTTP服务...")

			// 设置关闭超时时间
			if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
				logger.Errorf("Fiber HTTP服务关闭出错: %v", err)
				return err
			}

			logger.Info("Fiber HTTP服务已成功关闭")
			return nil
		},
	})
}

// FiberModule 提供Fiber HTTP模块
var FiberModule = fx.Options(
	fx.Provide(NewFiberApp),
	fx.Provide(NewFiberServer),
	fx.Invoke(StartFiberServer),
)
