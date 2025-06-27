package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zhoudm1743/go-flow/core/config"
	"go.uber.org/fx"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 日志接口
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger
}

// LogrusWrapper logrus 包装器
type LogrusWrapper struct {
	*logrus.Entry
}

// WithField 添加字段
func (l *LogrusWrapper) WithField(key string, value interface{}) Logger {
	return &LogrusWrapper{l.Entry.WithField(key, value)}
}

// WithFields 添加多个字段
func (l *LogrusWrapper) WithFields(fields map[string]interface{}) Logger {
	return &LogrusWrapper{l.Entry.WithFields(fields)}
}

// WithError 添加错误字段
func (l *LogrusWrapper) WithError(err error) Logger {
	return &LogrusWrapper{l.Entry.WithError(err)}
}

// Module fx模块
var Module = fx.Options(
	fx.Provide(NewLogger),
)

// NewLogger 创建新的日志实例
func NewLogger(cfg *config.Config) (Logger, error) {
	logger := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		return nil, fmt.Errorf("无效的日志级别 '%s': %w", cfg.Log.Level, err)
	}
	logger.SetLevel(level)

	// 设置日志格式
	switch strings.ToLower(cfg.Log.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,  // 强制启用颜色
			DisableColors:   false, // 确保不禁用颜色
			PadLevelText:    true,  // 对齐级别文本
			DisableQuote:    true,  // 禁用字段值的引号
		})
	default:
		return nil, fmt.Errorf("不支持的日志格式 '%s'", cfg.Log.Format)
	}

	// 设置日志输出
	switch strings.ToLower(cfg.Log.Output) {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	case "file":
		// 使用 lumberjack 进行日志轮转
		lumberjackLogger := &lumberjack.Logger{
			Filename:   cfg.Log.Lumberjack.Filename,
			MaxSize:    cfg.Log.Lumberjack.MaxSize,
			MaxAge:     cfg.Log.Lumberjack.MaxAge,
			MaxBackups: cfg.Log.Lumberjack.MaxBackups,
			Compress:   cfg.Log.Lumberjack.Compress,
		}

		// 确保日志目录存在
		if err := os.MkdirAll(filepath.Dir(cfg.Log.Lumberjack.Filename), 0755); err != nil {
			logger.SetOutput(os.Stdout)
			logger.Warnf("无法创建日志目录，回退到标准输出: %v", err)
		} else {
			logger.SetOutput(lumberjackLogger)
		}
	case "both":
		// 同时输出到文件和标准输出，文件使用 lumberjack
		lumberjackLogger := &lumberjack.Logger{
			Filename:   cfg.Log.Lumberjack.Filename,
			MaxSize:    cfg.Log.Lumberjack.MaxSize,
			MaxAge:     cfg.Log.Lumberjack.MaxAge,
			MaxBackups: cfg.Log.Lumberjack.MaxBackups,
			Compress:   cfg.Log.Lumberjack.Compress,
		}

		if err := os.MkdirAll(filepath.Dir(cfg.Log.Lumberjack.Filename), 0755); err != nil {
			logger.SetOutput(os.Stdout)
			logger.Warnf("无法创建日志目录，仅输出到标准输出: %v", err)
		} else {
			multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)
			logger.SetOutput(multiWriter)
		}
	default:
		logger.SetOutput(os.Stdout)
		logger.Warnf("不支持的日志输出 '%s'，使用标准输出", cfg.Log.Output)
	}

	return &LogrusWrapper{logger.WithFields(logrus.Fields{})}, nil
}

// GetGlobalLogger 获取全局日志实例（用于非依赖注入的场景）
var globalLogger Logger

// SetGlobalLogger 设置全局日志实例
func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}

// Debug 全局调试日志
func Debug(args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debug(args...)
	}
}

// Debugf 全局调试日志格式化
func Debugf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debugf(format, args...)
	}
}

// Info 全局信息日志
func Info(args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Info(args...)
	}
}

// Infof 全局信息日志格式化
func Infof(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Infof(format, args...)
	}
}

// Warn 全局警告日志
func Warn(args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warn(args...)
	}
}

// Warnf 全局警告日志格式化
func Warnf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warnf(format, args...)
	}
}

// Error 全局错误日志
func Error(args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Error(args...)
	}
}

// Errorf 全局错误日志格式化
func Errorf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Errorf(format, args...)
	}
}

// Fatal 全局致命错误日志
func Fatal(args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Fatal(args...)
	}
}

// Fatalf 全局致命错误日志格式化
func Fatalf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Fatalf(format, args...)
	}
}
