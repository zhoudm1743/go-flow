package facades

import (
	"github.com/sirupsen/logrus"
	"github.com/zhoudm1743/go-frame/pkg/log"
)

// Debug 记录调试级别日志
func (l *LogFacade) Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Info 记录信息级别日志
func (l *LogFacade) Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Warn 记录警告级别日志
func (l *LogFacade) Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Error 记录错误级别日志
func (l *LogFacade) Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Fatal 记录致命级别日志
func (l *LogFacade) Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Debugf 使用格式化字符串记录调试级别日志
func (l *LogFacade) Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Infof 使用格式化字符串记录信息级别日志
func (l *LogFacade) Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warnf 使用格式化字符串记录警告级别日志
func (l *LogFacade) Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Errorf 使用格式化字符串记录错误级别日志
func (l *LogFacade) Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatalf 使用格式化字符串记录致命级别日志
func (l *LogFacade) Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// WithField 创建带有单个字段的日志条目
func (l *LogFacade) WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

// WithFields 创建带有多个字段的日志条目
func (l *LogFacade) WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// Instance 获取原始日志实例
func (l *LogFacade) Instance() log.Logger {
	return GetLogger()
}
