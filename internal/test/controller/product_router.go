package controller

import (
	"github.com/zhoudm1743/go-flow/pkg/http/unified"
)

// ProductRouter 产品路由
type ProductRouter struct {
	controller *ProductController
}

// NewProductRouter 创建产品路由
func NewProductRouter(controller *ProductController) *ProductRouter {
	return &ProductRouter{
		controller: controller,
	}
}

// RegisterRoutes 注册路由
func (r *ProductRouter) RegisterRoutes(router unified.Router) {
	group := router.Group("/api/products")

	group.GET("", r.controller.List)
	group.GET("/:id", r.controller.Get)
	group.POST("", r.controller.Create)
	group.PUT("/:id", r.controller.Update)
	group.DELETE("/:id", r.controller.Delete)
	group.GET("/page", r.controller.Page)
}
