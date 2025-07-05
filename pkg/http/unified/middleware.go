package unified

import (
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

// HandlerFunc 统一的HTTP处理函数
type HandlerFunc func(Context) error

// MiddlewareFunc 统一的中间件函数
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Chain 中间件链
type Chain struct {
	middlewares []MiddlewareFunc
}

// NewChain 创建新的中间件链
func NewChain(middlewares ...MiddlewareFunc) *Chain {
	return &Chain{
		middlewares: middlewares,
	}
}

// Use 添加中间件到链中
func (c *Chain) Use(middlewares ...MiddlewareFunc) *Chain {
	c.middlewares = append(c.middlewares, middlewares...)
	return c
}

// Then 将处理函数包装在中间件链中
func (c *Chain) Then(handler HandlerFunc) HandlerFunc {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}
	return handler
}

// Compose 组合多个中间件为一个
func Compose(middlewares ...MiddlewareFunc) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// MiddlewareAdapter 中间件适配器接口
type MiddlewareAdapter interface {
	// AdaptGin 将统一中间件适配为Gin中间件
	AdaptGin(MiddlewareFunc) interface{}

	// AdaptFiber 将统一中间件适配为Fiber中间件
	AdaptFiber(MiddlewareFunc) interface{}
}

// DefaultMiddlewareAdapter 默认中间件适配器
type DefaultMiddlewareAdapter struct{}

// AdaptGin 将统一中间件适配为Gin中间件
func (a *DefaultMiddlewareAdapter) AdaptGin(middleware MiddlewareFunc) interface{} {
	return func(c interface{}) {
		ginCtx := c.(*gin.Context)
		ctx := NewGinContext(ginCtx)

		err := middleware(func(ctx Context) error {
			return nil
		})(ctx)

		if err != nil {
			ginCtx.Error(err)
			ginCtx.Abort()
		}
	}
}

// AdaptFiber 将统一中间件适配为Fiber中间件
func (a *DefaultMiddlewareAdapter) AdaptFiber(middleware MiddlewareFunc) interface{} {
	return func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c)

		return middleware(func(ctx Context) error {
			return nil
		})(ctx)
	}
}

// NewDefaultMiddlewareAdapter 创建默认中间件适配器
func NewDefaultMiddlewareAdapter() *DefaultMiddlewareAdapter {
	return &DefaultMiddlewareAdapter{}
}

// Adapt 将一个处理函数适配为中间件
func Adapt(handler HandlerFunc) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			if err := handler(c); err != nil {
				return err
			}
			return next(c)
		}
	}
}

// ToGinHandler 将HandlerFunc转换为Gin的处理函数
func ToGinHandler(handler HandlerFunc) func(*GinContext) {
	return func(c *GinContext) {
		if err := handler(c); err != nil {
			// 已经在handler中处理了错误
		}
	}
}

// ToFiberHandler 将HandlerFunc转换为Fiber的处理函数
func ToFiberHandler(handler HandlerFunc) func(*FiberContext) error {
	return func(c *FiberContext) error {
		return handler(c)
	}
}

// GinMiddlewareAdapter 将Gin中间件适配为统一的Middleware
func GinMiddlewareAdapter(ginMiddleware func(c *GinContext)) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			if gc, ok := c.(*GinContext); ok {
				done := make(chan bool)
				go func() {
					ginMiddleware(gc)
					done <- true
				}()
				<-done
				if gc.IsAborted() {
					return nil
				}
				return next(c)
			}
			// 非Gin上下文，跳过此中间件
			return next(c)
		}
	}
}

// FiberMiddlewareAdapter 将Fiber中间件适配为统一的Middleware
func FiberMiddlewareAdapter(fiberMiddleware func(c *FiberContext) error) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			if fc, ok := c.(*FiberContext); ok {
				if err := fiberMiddleware(fc); err != nil {
					return err
				}
				if fc.IsAborted() {
					return nil
				}
			}
			// 非Fiber上下文或未中止，继续执行
			return next(c)
		}
	}
}
