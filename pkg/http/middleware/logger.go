package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-frame/pkg/log"
)

// Logger 创建日志中间件
func Logger(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latency := endTime.Sub(startTime)

		// 请求方法
		reqMethod := c.Request.Method

		// 请求路由
		reqURI := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logMsg := fmt.Sprintf("[GIN] %s | %3d | %13v | %15s | %s",
			reqMethod,
			statusCode,
			latency,
			clientIP,
			reqURI,
		)

		// 根据状态码确定日志级别
		switch {
		case statusCode >= 500:
			logger.Errorf(logMsg)
		case statusCode >= 400:
			logger.Warnf(logMsg)
		default:
			logger.Infof(logMsg)
		}
	}
}
