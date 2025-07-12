package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zhoudm1743/go-frame/pkg/log"
)

// ANSI颜色代码常量
const (
	greenColor   = "\033[32m"
	yellowColor  = "\033[33m"
	redColor     = "\033[31m"
	blueColor    = "\033[34m"
	magentaColor = "\033[35m"
	cyanColor    = "\033[36m"
	whiteColor   = "\033[37m"
	resetColor   = "\033[0m"
)

// getStatusColor 根据HTTP状态码返回对应的颜色
func getStatusColor(code int) string {
	switch {
	case code >= 500:
		return redColor
	case code >= 400:
		return yellowColor
	case code >= 300:
		return cyanColor
	default:
		return greenColor
	}
}

// getMethodColor 根据HTTP方法返回对应的颜色
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return blueColor
	case "POST":
		return greenColor
	case "PUT":
		return yellowColor
	case "DELETE":
		return redColor
	case "PATCH":
		return cyanColor
	case "HEAD":
		return magentaColor
	case "OPTIONS":
		return whiteColor
	default:
		return resetColor
	}
}

// LogrusLogger 使用logrus作为Gin的日志输出
func LogrusLogger(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 请求路径
		path := c.Request.URL.Path

		// 获取原始查询参数
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latency := endTime.Sub(startTime)

		// 如果执行时间大于1秒，则改用秒为单位显示
		if latency > time.Second {
			latency = latency.Round(time.Second / 10)
		}

		// 状态码
		statusCode := c.Writer.Status()

		// 请求方法
		method := c.Request.Method

		// 客户端IP
		clientIP := c.ClientIP()

		// 构建彩色日志消息
		statusColor := getStatusColor(statusCode)
		methodColor := getMethodColor(method)
		resetColor := "\033[0m"

		// 构建状态码和方法的彩色输出
		coloredStatus := fmt.Sprintf("%s%3d%s", statusColor, statusCode, resetColor)
		coloredMethod := fmt.Sprintf("%s%s%s", methodColor, method, resetColor)

		// 构建URL路径（包含查询参数）
		fullPath := path
		if query != "" {
			fullPath = path + "?" + query
		}

		// 构建完整日志消息
		msg := fmt.Sprintf("[GIN] %s | %s | %12v | %s",
			coloredStatus,
			coloredMethod,
			latency,
			fullPath,
		)

		// 构建日志字段
		fields := logrus.Fields{
			"status":    statusCode,
			"method":    method,
			"latency":   latency,
			"client_ip": clientIP,
			"path":      fullPath,
		}

		// 获取错误信息
		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		// 使用不同级别记录日志
		switch {
		case statusCode >= 500:
			logger.WithFields(fields).Error(msg)
		case statusCode >= 400:
			logger.WithFields(fields).Warn(msg)
		case statusCode >= 300:
			logger.WithFields(fields).Info(msg)
		default:
			logger.WithFields(fields).Info(msg)
		}
	}
}

// GinLogToLogrus 将Gin框架自身的日志输出重定向到logrus
func GinLogToLogrus(logger log.Logger) {
	// 完全禁用Gin的控制台颜色
	gin.DisableConsoleColor()

	// 创建一个自定义的Writer，将日志输出重定向到logrus
	ginLogger := &LogrusWriter{logger: logger}
	gin.DefaultWriter = ginLogger
	gin.DefaultErrorWriter = ginLogger

	// 禁止Gin在启动时打印路由
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		fields := logrus.Fields{
			"method":      httpMethod,
			"path":        absolutePath,
			"handler":     handlerName,
			"handlers_no": nuHandlers,
		}
		logger.WithFields(fields).Debug("注册路由")
	}
}

// LogrusWriter 实现了io.Writer接口，用于重定向Gin的日志输出到logrus
type LogrusWriter struct {
	logger log.Logger
}

// Write 实现io.Writer接口
func (w *LogrusWriter) Write(p []byte) (n int, err error) {
	text := string(p)
	text = strings.TrimSuffix(text, "\n") // 去除尾部换行符

	// 不同类型的日志采用不同的处理方式
	if strings.Contains(text, "[GIN-debug]") {
		// 提取真正的消息内容
		message := strings.TrimPrefix(text, "[GIN-debug] ")

		// 解析路由注册信息
		if strings.Contains(message, "-->") {
			// 提取HTTP方法、路径和处理器信息
			parts := strings.Split(message, "-->")
			if len(parts) >= 2 {
				methodPath := strings.TrimSpace(parts[0])
				handler := strings.TrimSpace(parts[1])

				// 进一步解析方法和路径
				fields := strings.Fields(methodPath)
				if len(fields) >= 2 {
					method := fields[0]
					path := fields[1]

					// 使用logrus字段记录
					w.logger.WithFields(logrus.Fields{
						"component": "gin",
						"method":    method,
						"path":      path,
						"handler":   handler,
					}).Debug("路由注册")
				}
			}
		} else if strings.Contains(message, "Listening and serving HTTP") {
			// 服务器启动信息
			w.logger.WithField("component", "gin").Info(message)
		} else {
			// 其他调试信息
			w.logger.WithField("component", "gin").Debug(message)
		}
	} else if strings.Contains(text, "[GIN]") {
		// GIN前缀的非调试日志
		message := strings.TrimPrefix(text, "[GIN] ")
		w.logger.WithField("component", "gin").Info(message)
	} else {
		// 其他日志
		w.logger.Info(text)
	}

	return len(p), nil
}
