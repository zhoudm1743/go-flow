package module

import (
	"github.com/zhoudm1743/go-frame/internal/module/controller"
	"github.com/zhoudm1743/go-frame/internal/module/repository"
	"github.com/zhoudm1743/go-frame/internal/module/service"
	"go.uber.org/fx"
)

// Module 模块定义
type Module struct {
	moduleID string
	name     string
}

// NewModule 创建模块
func NewModule() *Module {
	return &Module{
		moduleID: "module",
		name:     "Module",
	}
}

// Name 模块名称
func (m *Module) Name() string {
	return m.name
}

// Options 模块配置选项
func (m *Module) Options() fx.Option {
	return fx.Module(
		m.moduleID,
		// 注册各层组件
		repository.Repository,
		service.Service,
		controller.Router,
	)
}
