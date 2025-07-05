package test

import (
	"fmt"

	"github.com/zhoudm1743/go-flow/internal/test/controller"
	"github.com/zhoudm1743/go-flow/internal/test/repository"
	"github.com/zhoudm1743/go-flow/internal/test/service"
	"github.com/zhoudm1743/go-flow/pkg/http"
	"github.com/zhoudm1743/go-flow/pkg/http/unified"
	"github.com/zhoudm1743/go-flow/pkg/log"
	"go.uber.org/fx"
)

// Router 路由接口
type Router interface {
	RegisterRoutes(router unified.Router)
}

// Module 模块定义
type Module struct {
	moduleID string
	name     string
}

// NewModule 创建模块
func NewModule() *Module {
	return &Module{
		moduleID: "test",
		name:     "Test",
	}
}

// Name 模块名称
func (m *Module) Name() string {
	return m.name
}

// RoutePrefix 获取模块路由前缀
func (m *Module) RoutePrefix() string {
	return "/test"
}

// Options 模块配置选项
func (m *Module) Options() fx.Option {
	// 将模块ID作为fx命名的一部分，避免多个模块实例冲突
	name := fmt.Sprintf("module_%s", m.moduleID)

	return fx.Module(
		name,
		// 提供所有依赖
		fx.Provide(
			repository.NewProductRepository,
			service.NewProductService,
			controller.NewProductController,
			controller.NewProductRouter,
			// 注册路由组
			fx.Annotate(
				func(router *controller.ProductRouter) Router {
					return router
				},
				fx.ResultTags(fmt.Sprintf(`group:"%s_routers"`, m.moduleID)),
			),
		),
		// 注册路由
		fx.Invoke(
			fx.Annotate(
				func(server http.Server, logger log.Logger, routers []Router) {
					logger.Infof("注册模块: %s", m.Name())
					for _, router := range routers {
						router.RegisterRoutes(server.Router())
					}
					logger.Infof("模块 %s 路由注册完成", m.moduleID)
				},
				fx.ParamTags(``, ``, fmt.Sprintf(`group:"%s_routers"`, m.moduleID)),
			),
		),
	)
}
