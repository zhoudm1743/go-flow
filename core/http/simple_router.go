package http

import (
	"github.com/gin-gonic/gin"
)

// RouterBuilder 路由构建器 - 使用链式调用简化路由定义
type RouterBuilder struct {
	engine *gin.Engine
	group  *gin.RouterGroup
}

// NewRouter 创建新的路由构建器
func NewRouter(engine *gin.Engine) *RouterBuilder {
	return &RouterBuilder{
		engine: engine,
		group:  &engine.RouterGroup,
	}
}

// Group 创建路由组 - 支持链式调用
func (r *RouterBuilder) Group(prefix string, middlewares ...gin.HandlerFunc) *RouterBuilder {
	group := r.group.Group(prefix)
	for _, middleware := range middlewares {
		group.Use(middleware)
	}
	return &RouterBuilder{
		engine: r.engine,
		group:  group,
	}
}

// Use 添加中间件
func (r *RouterBuilder) Use(middlewares ...gin.HandlerFunc) *RouterBuilder {
	r.group.Use(middlewares...)
	return r
}

// GET 注册GET路由
func (r *RouterBuilder) GET(path string, handlers ...gin.HandlerFunc) *RouterBuilder {
	r.group.GET(path, handlers...)
	return r
}

// POST 注册POST路由
func (r *RouterBuilder) POST(path string, handlers ...gin.HandlerFunc) *RouterBuilder {
	r.group.POST(path, handlers...)
	return r
}

// PUT 注册PUT路由
func (r *RouterBuilder) PUT(path string, handlers ...gin.HandlerFunc) *RouterBuilder {
	r.group.PUT(path, handlers...)
	return r
}

// DELETE 注册DELETE路由
func (r *RouterBuilder) DELETE(path string, handlers ...gin.HandlerFunc) *RouterBuilder {
	r.group.DELETE(path, handlers...)
	return r
}

// PATCH 注册PATCH路由
func (r *RouterBuilder) PATCH(path string, handlers ...gin.HandlerFunc) *RouterBuilder {
	r.group.PATCH(path, handlers...)
	return r
}

// SimpleRouteHandler 简化的路由处理器接口
type SimpleRouteHandler interface {
	RegisterRoutes(router *RouterBuilder)
}

// ServiceRouteHandler 带服务的路由处理器接口
type ServiceRouteHandler interface {
	RegisterRoutes(router *RouterBuilder, service interface{})
}

// SimpleRouteRegistrar 简化的路由注册器
type SimpleRouteRegistrar struct {
	handler SimpleRouteHandler
}

// ServiceRouteRegistrar 带服务的路由注册器
type ServiceRouteRegistrar struct {
	handler ServiceRouteHandler
	service interface{}
}

// NewSimpleRouteRegistrar 创建简化路由注册器
func NewSimpleRouteRegistrar(handler SimpleRouteHandler) *SimpleRouteRegistrar {
	return &SimpleRouteRegistrar{
		handler: handler,
	}
}

// NewServiceRouteRegistrar 创建带服务的路由注册器
func NewServiceRouteRegistrar(handler ServiceRouteHandler, service interface{}) *ServiceRouteRegistrar {
	return &ServiceRouteRegistrar{
		handler: handler,
		service: service,
	}
}

// RegisterRoutes 实现RouteRegistrator接口
func (r *SimpleRouteRegistrar) RegisterRoutes(engine *gin.Engine) error {
	router := NewRouter(engine)
	r.handler.RegisterRoutes(router)
	return nil
}

// RegisterRoutes 实现RouteRegistrator接口
func (r *ServiceRouteRegistrar) RegisterRoutes(engine *gin.Engine) error {
	router := NewRouter(engine)
	r.handler.RegisterRoutes(router, r.service)
	return nil
}

// ControllerBase 控制器基类 - 提供常用的响应方法
type ControllerBase struct{}

// Success 成功响应
func (cb *ControllerBase) Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

// Error 错误响应
func (cb *ControllerBase) Error(c *gin.Context, code int, message string) {
	c.JSON(200, gin.H{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}

// BadRequest 参数错误响应
func (cb *ControllerBase) BadRequest(c *gin.Context, message string) {
	cb.Error(c, 400, message)
}

// NotFound 未找到响应
func (cb *ControllerBase) NotFound(c *gin.Context, message string) {
	cb.Error(c, 404, message)
}

// InternalError 内部错误响应
func (cb *ControllerBase) InternalError(c *gin.Context, message string) {
	cb.Error(c, 500, message)
}

// RouteGroup 路由组配置函数类型
type RouteGroup func(*RouterBuilder)

// RegisterRouteGroups 批量注册路由组
func RegisterRouteGroups(engine *gin.Engine, groups ...RouteGroup) {
	router := NewRouter(engine)
	for _, group := range groups {
		group(router)
	}
}

// RouteRegistratorFunc 函数类型实现RouteRegistrator接口
type RouteRegistratorFunc func(*gin.Engine) error

// RegisterRoutes 实现RouteRegistrator接口
func (f RouteRegistratorFunc) RegisterRoutes(engine *gin.Engine) error {
	return f(engine)
}
