package middleware

import (
	"fmt"
	"io"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/zhoudm1743/go-flow/pkg/log"
)

// FiberLogrusConfig 配置
type FiberLogrusConfig struct {
	Logger log.Logger
}

// FiberLogrusLogger 使用logrus作为Fiber的日志输出
func FiberLogrusLogger(logger log.Logger) fiber.Handler {
	// 创建自定义中间件
	return func(c *fiber.Ctx) error {
		// 开始时间
		startTime := time.Now()

		// 请求路径
		path := c.Path()

		// 获取原始查询参数
		query := string(c.Request().URI().QueryString())

		// 处理请求
		err := c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latency := endTime.Sub(startTime)

		// 如果执行时间大于1秒，则改用秒为单位显示
		if latency > time.Second {
			latency = latency.Round(time.Second / 10)
		}

		// 状态码
		statusCode := c.Response().StatusCode()

		// 请求方法
		method := c.Method()

		// 客户端IP
		clientIP := c.IP()

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
		msg := fmt.Sprintf("[FIBER] %s | %s | %12v | %s",
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

		return err
	}
}

// FiberLogrusOutput 创建一个Fiber日志输出器
func FiberLogrusOutput(logger log.Logger) io.Writer {
	return &fiberLogrusWriter{logger: logger}
}

// fiberLogrusWriter 实现了io.Writer接口，用于重定向Fiber日志到logrus
type fiberLogrusWriter struct {
	logger log.Logger
}

// Write 实现io.Writer接口
func (w *fiberLogrusWriter) Write(p []byte) (n int, err error) {
	text := string(p)
	w.logger.Debug(fmt.Sprintf("[FIBER] %s", text))
	return len(p), nil
}
