package response

import (
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

// NoRoute 无路由的响应
func NoRoute(c *gin.Context) {
	Fail(c, Request404Error)
}

// NoMethod 无方法的响应
func NoMethod(c *gin.Context) {
	Fail(c, Request405Error)
}

// ErrDuplicateName 重复名称的响应
func ErrDuplicateName(c *gin.Context) {
	Fail(c, RequestErrDuplicateNameError)
}

// FiberNoRoute Fiber无路由的响应
func FiberNoRoute(c *fiber.Ctx) error {
	return FiberFail(c, Request404Error)
}

// FiberNoMethod Fiber无方法的响应
func FiberNoMethod(c *fiber.Ctx) error {
	return FiberFail(c, Request405Error)
}

// FiberErrDuplicateName Fiber重复名称的响应
func FiberErrDuplicateName(c *fiber.Ctx) error {
	return FiberFail(c, RequestErrDuplicateNameError)
}
