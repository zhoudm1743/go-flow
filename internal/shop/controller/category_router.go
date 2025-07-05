package controller

import (
	"github.com/zhoudm1743/go-flow/pkg/http/unified"
)

// CategoryRouter 分类路由
type CategoryRouter struct {
	controller *CategoryController
}

// NewCategoryRouter 创建分类路由
func NewCategoryRouter(controller *CategoryController) *CategoryRouter {
	return &CategoryRouter{
		controller: controller,
	}
}

// RegisterRoutes 注册路由
func (r *CategoryRouter) RegisterRoutes(router unified.Router) {
	group := router.Group("/api/categorys")

	group.GET("", r.controller.List)
	group.GET("/:id", r.controller.Get)
	group.POST("", r.controller.Create)
	group.PUT("/:id", r.controller.Update)
	group.DELETE("/:id", r.controller.Delete)
	group.GET("/page", r.controller.Page)
}
