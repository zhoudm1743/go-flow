package user

import (
	"fmt"

	userController "github.com/zhoudm1743/go-flow/internal/user/controller/user"
	"github.com/zhoudm1743/go-flow/internal/user/repository"
	userService "github.com/zhoudm1743/go-flow/internal/user/service/user"
	"github.com/zhoudm1743/go-flow/pkg/http"
	"github.com/zhoudm1743/go-flow/pkg/http/unified"
	"github.com/zhoudm1743/go-flow/pkg/log"
	"go.uber.org/fx"
)

// Router 统一路由器接口
type Router interface {
	// RegisterRoutes 注册路由到路由器
	RegisterRoutes(router unified.Router)
}

// Module User模块
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
		return "user"
	}
	return "user-" + m.ModuleName
}

// RoutePrefix 获取模块路由前缀
func (m *Module) RoutePrefix() string {
	if m.ModuleName == "" {
		return "/user"
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
	// 这里添加默认的提供者
	defaultProviders := []interface{}{
		// Account 相关组件
		repository.NewAccountRepository,
		userService.NewAccountService,
		userController.NewAccountController,
		fx.Annotate(
			userController.NewAccountRouter,
			fx.As(new(Router)),
			fx.As(new(http.RouterRegister)),
			fx.ResultTags(`group:"user_routers" group:"http_routers"`),
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
			fx.Annotate(
				func(server http.Server, logger log.Logger, routers []Router) {
					logger.Infof("注册模块: %s", m.Name())
					for _, router := range routers {
						router.RegisterRoutes(server.Router())
					}
				},
				fx.ParamTags(``, ``, `group:"user_routers"`),
			),
		),
	)
}

// NewModule 创建User模块
func NewModule() *Module {
	return &Module{
		ModuleID: "user_default",
	}
}

// NewModuleWithName 创建带名称的User模块
func NewModuleWithName(name string) *Module {
	return &Module{
		ModuleName: name,
		ModuleID:   "user_" + name,
	}
}
