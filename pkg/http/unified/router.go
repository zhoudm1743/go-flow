package unified

import (
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

// HTTPMethod 定义HTTP方法类型
type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	PATCH   HTTPMethod = "PATCH"
	HEAD    HTTPMethod = "HEAD"
	OPTIONS HTTPMethod = "OPTIONS"
)

// Router 统一的路由器接口
type Router interface {
	// 基本路由方法
	GET(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router
	POST(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router
	PUT(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router
	DELETE(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router
	PATCH(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router
	OPTIONS(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router
	HEAD(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router

	// 静态文件
	Static(prefix, root string) Router

	// 路由组
	Group(prefix string, middlewares ...MiddlewareFunc) Router

	// 中间件
	Use(middlewares ...MiddlewareFunc) Router

	// 处理请求
	Handle(method HTTPMethod, path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router
}

// RouterAdapter 路由适配器接口
type RouterAdapter interface {
	// AdaptHandler 将统一处理函数适配为特定框架的处理函数
	AdaptHandler(handler HandlerFunc) interface{}

	// AdaptMiddleware 将统一中间件适配为特定框架的中间件
	AdaptMiddleware(middleware MiddlewareFunc) interface{}
}

// RouterRegister 路由注册器接口
type RouterRegister interface {
	// RegisterRoutes 注册路由到路由器
	RegisterRoutes(router Router)
}

// EngineType 定义引擎类型
type EngineType int

const (
	GinEngine EngineType = iota
	FiberEngine
)

// RouterImpl 路由器实现
type RouterImpl struct {
	engineType EngineType
	ginEngine  *gin.Engine
	fiberApp   *fiber.App
	prefix     string
	middleware []MiddlewareFunc
}

// NewRouter 创建新的路由器
func NewRouter(engineType EngineType, ginEngine *gin.Engine, fiberApp *fiber.App) Router {
	return &RouterImpl{
		engineType: engineType,
		ginEngine:  ginEngine,
		fiberApp:   fiberApp,
		prefix:     "",
		middleware: []MiddlewareFunc{},
	}
}

// Handle 实现Router接口
func (r *RouterImpl) Handle(method HTTPMethod, path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router {
	fullPath := r.prefix + path

	// 合并中间件
	allMiddlewares := append(r.middleware, middlewares...)

	// 应用中间件
	if len(allMiddlewares) > 0 {
		for _, m := range allMiddlewares {
			handler = m(handler)
		}
	}

	switch r.engineType {
	case GinEngine:
		if r.ginEngine != nil {
			r.ginEngine.Handle(string(method), fullPath, func(c *gin.Context) {
				ctx := NewGinContext(c)
				if err := handler(ctx); err != nil {
					c.Error(err)
				}
			})
		}
	case FiberEngine:
		if r.fiberApp != nil {
			r.fiberApp.Add(string(method), fullPath, func(c *fiber.Ctx) error {
				ctx := NewFiberContext(c)
				return handler(ctx)
			})
		}
	}

	return r
}

// GET 实现Router接口
func (r *RouterImpl) GET(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router {
	return r.Handle(GET, path, handler, middlewares...)
}

// POST 实现Router接口
func (r *RouterImpl) POST(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router {
	return r.Handle(POST, path, handler, middlewares...)
}

// PUT 实现Router接口
func (r *RouterImpl) PUT(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router {
	return r.Handle(PUT, path, handler, middlewares...)
}

// DELETE 实现Router接口
func (r *RouterImpl) DELETE(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router {
	return r.Handle(DELETE, path, handler, middlewares...)
}

// PATCH 实现Router接口
func (r *RouterImpl) PATCH(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router {
	return r.Handle(PATCH, path, handler, middlewares...)
}

// HEAD 实现Router接口
func (r *RouterImpl) HEAD(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router {
	return r.Handle(HEAD, path, handler, middlewares...)
}

// OPTIONS 实现Router接口
func (r *RouterImpl) OPTIONS(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) Router {
	return r.Handle(OPTIONS, path, handler, middlewares...)
}

// Static 实现Router接口
func (r *RouterImpl) Static(prefix, root string) Router {
	fullPrefix := r.prefix + prefix

	switch r.engineType {
	case GinEngine:
		if r.ginEngine != nil {
			r.ginEngine.Static(fullPrefix, root)
		}
	case FiberEngine:
		if r.fiberApp != nil {
			r.fiberApp.Static(fullPrefix, root)
		}
	}
	return r
}

// Use 实现Router接口
func (r *RouterImpl) Use(middlewares ...MiddlewareFunc) Router {
	r.middleware = append(r.middleware, middlewares...)
	return r
}

// Group 实现Router接口
func (r *RouterImpl) Group(prefix string, middlewares ...MiddlewareFunc) Router {
	group := &RouterImpl{
		engineType: r.engineType,
		ginEngine:  r.ginEngine,
		fiberApp:   r.fiberApp,
		prefix:     r.prefix + prefix,
		middleware: append(r.middleware, middlewares...),
	}
	return group
}
