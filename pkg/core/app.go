package core

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zhoudm1743/go-frame/pkg/cache"
	"github.com/zhoudm1743/go-frame/pkg/config"
	"github.com/zhoudm1743/go-frame/pkg/database"
	"github.com/zhoudm1743/go-frame/pkg/facades"
	"github.com/zhoudm1743/go-frame/pkg/log"
	"go.uber.org/fx"
)

// App 应用结构体
type App struct {
	opts   []fx.Option
	name   string
	logger log.Logger
}

// NewApp 创建新应用
func NewApp(name string) *App {
	app := &App{
		name: name,
		opts: []fx.Option{},
	}

	// 添加基础模块
	app.opts = append(app.opts,
		config.Module,
		log.Module,
		database.Module,
		// 使用统一HTTP框架，但不添加默认模块，让应用自己选择
		cache.Module,
	)

	// 添加门面模块（在基础模块之后，确保服务已注册）
	app.opts = append(app.opts,
		facades.Module,
	)

	return app
}

// AddModule 添加应用模块
func (a *App) AddModule(module Module) *App {
	a.opts = append(a.opts, module.Options())
	return a
}

// AddModules 添加多个应用模块
func (a *App) AddModules(modules ...Module) *App {
	for _, m := range modules {
		a.AddModule(m)
	}
	return a
}

// WithOptions 添加fx选项
func (a *App) WithOptions(opts ...fx.Option) *App {
	a.opts = append(a.opts, opts...)
	return a
}

// Run 运行应用
func (a *App) Run() {
	a.RunWithOptions(true)
}

// RunWithOptions 使用选项运行应用
func (a *App) RunWithOptions(blocking bool) {
	// 创建FX应用
	fxApp := fx.New(
		fx.Options(a.opts...),
		fx.NopLogger, // 禁用FX默认日志
		fx.Invoke(func(logger log.Logger) {
			// 保存日志记录器供后续使用
			a.logger = logger
			logger.Infof("应用 %s 启动中...", a.name)
		}),
	)

	// 启动应用，非阻塞方式
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := fxApp.Start(startCtx); err != nil {
		fmt.Fprintf(os.Stderr, "应用启动错误: %v\n", err)
		fmt.Fprintf(os.Stderr, "完整错误信息: %+v\n", err)
		os.Exit(1)
	}

	a.logger.Info("应用已启动")

	if blocking {
		// 等待中断信号优雅地关闭应用
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		a.logger.Info("正在关闭应用...")

		// 关闭应用
		stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := fxApp.Stop(stopCtx); err != nil {
			a.logger.Errorf("应用关闭错误: %v", err)
		}

		a.logger.Info("应用已优雅关闭")
	}
}
