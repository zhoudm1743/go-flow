package demo

import (
	"fmt"

	"github.com/zhoudm1743/go-flow/internal/demo/controller"
	"github.com/zhoudm1743/go-flow/internal/demo/repository"
	"github.com/zhoudm1743/go-flow/internal/demo/service"
	"github.com/zhoudm1743/go-flow/pkg/http"
	"github.com/zhoudm1743/go-flow/pkg/log"
	"go.uber.org/fx"
)

// Module 示例模块
type Module struct {
	// 模块名称，可选，用于路由前缀
	ModuleName string
	// 模块唯一标识
	ModuleID string
	// 组件提供者集合
	providers []interface{}
}

// Name 模块名称
func (m *Module) Name() string {
	if m.ModuleName == "" {
		return "demo"
	}
	return m.ModuleName
}

// RoutePrefix 获取模块路由前缀
func (m *Module) RoutePrefix() string {
	if m.ModuleName == "" {
		return ""
	}
	return "/" + m.ModuleName
}

// RegisterProvider 注册组件提供者
func (m *Module) RegisterProvider(provider interface{}) *Module {
	m.providers = append(m.providers, provider)
	return m
}

// Providers 获取所有组件提供者
func (m *Module) Providers() []interface{} {
	// 默认提供的组件
	defaultProviders := []interface{}{
		// 仓库
		repository.NewUserRepository,
		// 服务
		service.NewUserService,
		// 控制器
		controller.NewUserController,
		// 路由器
		fx.Annotate(
			controller.NewUserRouter,
			fx.As(new(http.RouterRegister)),
			fx.ResultTags(`group:"demo_routers"`),
		),
	}

	// 合并默认和自定义组件
	return append(defaultProviders, m.providers...)
}

// Options 模块配置选项
func (m *Module) Options() fx.Option {
	// 将模块ID作为fx命名的一部分，避免多个模块实例冲突
	name := fmt.Sprintf("module_%s", m.ModuleID)

	return fx.Module(
		name,
		fx.Provide(m.Providers()...),
		fx.Invoke(
			// 注册路由
			fx.Annotate(
				func(server http.Server, logger log.Logger, registers []http.RouterRegister) {
					logger.Infof("注册模块: %s", m.Name())

					// 注册所有路由到服务器
					for _, register := range registers {
						server.RegisterRoutes(register)
					}

					logger.Infof("模块 %s 路由注册完成", m.Name())
				},
				fx.ParamTags("", "", `group:"demo_routers"`),
			),
		),
	)
}

// NewModuleWithName 创建带名称的示例模块
func NewModuleWithName(name string) *Module {
	return &Module{
		ModuleName: name,
		ModuleID:   name,
	}
}
