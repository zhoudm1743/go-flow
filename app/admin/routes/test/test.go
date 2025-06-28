package test

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-flow/app/admin/service/test"
	httpCore "github.com/zhoudm1743/go-flow/core/http"
)

type testRoutes struct {
	srv test.TestService
}

// NewTestGroup fx Provider函数，自动注入TestService并返回配置好的Group
func NewTestGroup(srv test.TestService) httpCore.Group {
	return httpCore.NewGroup("/test",
		func() interface{} {
			return &testRoutes{srv: srv}
		},
		regTest,
	)
}

// regTest 注册测试路由（内部函数）
func regTest(rg *httpCore.BaseGroup, instance interface{}) error {
	r := instance.(*testRoutes)

	// 使用新的可变参数API注册路由
	rg.GET("/test", r.test)

	// 演示多个处理器的用法
	rg.GET("/test-multi",
		func(c *gin.Context) {
			// 第一个处理器作为中间件
			c.Set("middleware_data", "from middleware")
			c.Next()
		},
		r.testMulti, // 使用实例方法作为主处理器
	)

	// POST 路由示例
	rg.POST("/test", r.testPost)

	return nil
}

func (r *testRoutes) test(c *gin.Context) {
	res := r.srv.Test()
	c.JSON(http.StatusOK, res)
}

func (r *testRoutes) testMulti(c *gin.Context) {
	middlewareData := c.GetString("middleware_data")
	c.JSON(http.StatusOK, gin.H{
		"message": "test with multiple handlers",
		"data":    middlewareData,
		"from":    "instance method",
	})
}

func (r *testRoutes) testPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "test post from instance method",
		"method":  "POST",
	})
}
