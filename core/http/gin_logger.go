package http

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-flow/core/logger"
)

// GinLoggerWriter 实现 io.Writer 接口，将 Gin 的日志写入到我们的 logger
type GinLoggerWriter struct {
	logger logger.Logger
}

// Write 实现 io.Writer 接口
func (w *GinLoggerWriter) Write(p []byte) (n int, err error) {
	// 移除末尾的换行符
	msg := string(p)
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}
	w.logger.Info(msg)
	return len(p), nil
}

// NewGinLoggerWriter 创建新的 Gin 日志写入器
func NewGinLoggerWriter(log logger.Logger) io.Writer {
	return &GinLoggerWriter{logger: log}
}

// GinLoggerMiddleware 创建 Gin 日志中间件
func GinLoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 记录请求信息
		log.WithFields(map[string]interface{}{
			"client_ip":   param.ClientIP,
			"timestamp":   param.TimeStamp.Format(time.RFC3339),
			"method":      param.Method,
			"path":        param.Path,
			"protocol":    param.Request.Proto,
			"status_code": param.StatusCode,
			"latency":     param.Latency,
			"user_agent":  param.Request.UserAgent(),
			"data_length": param.BodySize,
		}).Infof("%s %s %d %v",
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)

		// 如果有错误消息，单独记录
		if param.ErrorMessage != "" {
			log.WithField("error", param.ErrorMessage).Error("请求处理错误")
		}

		// 返回空字符串，因为我们已经通过自定义 logger 记录了
		return ""
	})
}

// GinRecoveryMiddleware 创建 Gin 恢复中间件
func GinRecoveryMiddleware(log logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.WithFields(map[string]interface{}{
			"panic":     recovered,
			"path":      c.Request.URL.Path,
			"method":    c.Request.Method,
			"client_ip": c.ClientIP(),
		}).Error("服务器内部错误，请求处理时发生恐慌")

		c.JSON(500, gin.H{
			"error":   "服务器内部错误",
			"message": "请求处理失败",
		})
	})
}
