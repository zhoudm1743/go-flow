package http

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

// RouteHandler 路由处理器函数类型
type RouteHandler gin.HandlerFunc

// BaseGroup 基础路由组，提供路由注册方法
type BaseGroup struct {
	group *gin.RouterGroup
}

// GET 注册 GET 路由
func (bg *BaseGroup) GET(path string, handlers ...RouteHandler) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = gin.HandlerFunc(handler)
	}
	bg.group.GET(path, ginHandlers...)
}

// POST 注册 POST 路由
func (bg *BaseGroup) POST(path string, handlers ...RouteHandler) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = gin.HandlerFunc(handler)
	}
	bg.group.POST(path, ginHandlers...)
}

// PUT 注册 PUT 路由
func (bg *BaseGroup) PUT(path string, handlers ...RouteHandler) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = gin.HandlerFunc(handler)
	}
	bg.group.PUT(path, ginHandlers...)
}

// DELETE 注册 DELETE 路由
func (bg *BaseGroup) DELETE(path string, handlers ...RouteHandler) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = gin.HandlerFunc(handler)
	}
	bg.group.DELETE(path, ginHandlers...)
}

// PATCH 注册 PATCH 路由
func (bg *BaseGroup) PATCH(path string, handlers ...RouteHandler) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = gin.HandlerFunc(handler)
	}
	bg.group.PATCH(path, ginHandlers...)
}

// HEAD 注册 HEAD 路由
func (bg *BaseGroup) HEAD(path string, handlers ...RouteHandler) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = gin.HandlerFunc(handler)
	}
	bg.group.HEAD(path, ginHandlers...)
}

// OPTIONS 注册 OPTIONS 路由
func (bg *BaseGroup) OPTIONS(path string, handlers ...RouteHandler) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = gin.HandlerFunc(handler)
	}
	bg.group.OPTIONS(path, ginHandlers...)
}

// Group 路由组定义
type Group struct {
	Prefix      string            // 路由前缀
	Constructor interface{}       // 构造函数
	Register    interface{}       // 注册函数
	Middleware  []gin.HandlerFunc // 中间件
}

// RouteRegistrar 路由注册器
type RouteRegistrar struct {
	engine       *gin.Engine
	dependencies map[reflect.Type]interface{} // 依赖注入容器
}

// NewRouteRegistrar 创建路由注册器
func NewRouteRegistrar(engine *gin.Engine) *RouteRegistrar {
	return &RouteRegistrar{
		engine:       engine,
		dependencies: make(map[reflect.Type]interface{}),
	}
}

// AddDependency 添加依赖
func (r *RouteRegistrar) AddDependency(dep interface{}) {
	r.dependencies[reflect.TypeOf(dep)] = dep
}

// RegisterModuleGroups 注册模块路由组
func (r *RouteRegistrar) RegisterModuleGroups(module string, groups []Group) error {
	moduleGroup := r.engine.Group("/" + module)

	for _, group := range groups {
		if err := r.registerGroup(moduleGroup, group); err != nil {
			return err
		}
	}

	return nil
}

// registerGroup 注册单个路由组
func (r *RouteRegistrar) registerGroup(parent *gin.RouterGroup, group Group) error {
	// 创建路由组
	routeGroup := parent.Group(group.Prefix)

	// 应用中间件
	for _, middleware := range group.Middleware {
		routeGroup.Use(middleware)
	}

	// 创建 BaseGroup
	baseGroup := &BaseGroup{group: routeGroup}

	// 如果有构造函数，先创建实例
	var instance interface{}
	if group.Constructor != nil {
		// 这里简化处理，实际项目中可能需要依赖注入容器
		// 暂时假设构造函数不需要参数或者使用默认参数
		if constructor, ok := group.Constructor.(func() interface{}); ok {
			instance = constructor()
		}
	}

	// 如果有注册函数，调用它
	if group.Register != nil {
		if instance != nil {
			// 尝试调用带实例参数的注册函数
			if regFunc, ok := group.Register.(func(*BaseGroup, interface{}) error); ok {
				return regFunc(baseGroup, instance)
			}
		}
		// 回退到原来的无参数版本
		if regFunc, ok := group.Register.(func(*BaseGroup) error); ok {
			return regFunc(baseGroup)
		}
	}

	return nil
}

// Reg 用于在注册函数中执行路由注册 - 支持带实例参数的版本
func RegWithInstance(regFunc func(interface{})) func(*BaseGroup, interface{}) error {
	return func(bg *BaseGroup, instance interface{}) error {
		regFunc(instance)
		return nil
	}
}

// Reg 用于在注册函数中执行路由注册 - 原版本保持兼容
func Reg(regFunc func()) error {
	regFunc()
	return nil
}

// RegisterModuleRoutes 便捷函数，用于注册模块路由
func RegisterModuleRoutes(engine *gin.Engine, module string, groups []Group) error {
	registrar := NewRouteRegistrar(engine)
	return registrar.RegisterModuleGroups(module, groups)
}

// NewGroup 创建新的路由组（工厂函数）
// 支持按照test.go的模式声明路由组
func NewGroup(prefix string, constructor interface{}, register interface{}, middleware ...gin.HandlerFunc) Group {
	return Group{
		Prefix:      prefix,
		Constructor: constructor,
		Register:    register,
		Middleware:  middleware,
	}
}

// GlobalRegistrar 全局路由注册器
var GlobalRegistrar *RouteRegistrar

// InitGlobalRegistrar 初始化全局路由注册器
func InitGlobalRegistrar(engine *gin.Engine) {
	GlobalRegistrar = NewRouteRegistrar(engine)
}

// RegisterGlobalRoutes 全局路由注册函数
func RegisterGlobalRoutes(module string, groups ...Group) error {
	if GlobalRegistrar == nil {
		panic("GlobalRegistrar not initialized. Call InitGlobalRegistrar first.")
	}
	return GlobalRegistrar.RegisterModuleGroups(module, groups)
}
