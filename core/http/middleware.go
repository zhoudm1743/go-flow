package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		// 设置允许的响应头
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 设置源
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 对于预检请求，直接返回
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 继续处理请求
		c.Next()
	}
}

// RequestIDMiddleware 请求 ID 中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// 可以使用 UUID 或其他方式生成请求 ID
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// generateRequestID 生成请求 ID（简单实现）
func generateRequestID() string {
	// 这里使用时间戳作为简单的请求 ID
	// 在生产环境中应该使用更好的 UUID 生成器
	return "req_" + fmt.Sprintf("%d", time.Now().UnixNano())
}

// AuthMiddleware 认证中间件
func AuthMiddleware(c *gin.Context) {
	// 简单的认证检查示例
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		c.Abort()
		return
	}

	// 这里可以添加实际的token验证逻辑
	// 例如：验证JWT token，检查用户权限等

	c.Next()
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware(c *gin.Context) {
	// 检查管理员权限
	userRole := c.GetHeader("X-User-Role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "需要管理员权限",
		})
		c.Abort()
		return
	}

	c.Next()
}
