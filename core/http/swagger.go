package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwaggerRegistrator swagger路由注册器
type SwaggerRegistrator struct{}

// NewSwaggerRegistrator 创建swagger路由注册器
func NewSwaggerRegistrator() RouteRegistrator {
	return &SwaggerRegistrator{}
}

// RegisterRoutes 注册swagger路由
func (s *SwaggerRegistrator) RegisterRoutes(engine *gin.Engine) error {
	// 注册swagger路由
	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 添加重定向，方便访问
	engine.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/docs/index.html")
	})

	return nil
}
