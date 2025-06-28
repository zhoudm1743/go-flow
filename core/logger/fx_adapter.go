package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx/fxevent"
)

// FxLogger fx 日志适配器
type FxLogger struct {
	logger *logrus.Logger
}

// NewFxLogger 创建 fx 日志适配器
func NewFxLogger(logger *logrus.Logger) *FxLogger {
	return &FxLogger{logger: logger}
}

// LogEvent 实现 fxevent.Logger 接口
func (l *FxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logger.WithFields(logrus.Fields{
			"caller": e.FunctionName,
			"callee": e.CallerName,
		}).Info("正在执行OnStart钩子")

	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logger.WithFields(logrus.Fields{
				"caller": e.FunctionName,
				"callee": e.CallerName,
				"error":  e.Err.Error(),
			}).Error("OnStart钩子执行失败")
		} else {
			l.logger.WithFields(logrus.Fields{
				"caller":  e.FunctionName,
				"callee":  e.CallerName,
				"runtime": e.Runtime.String(),
			}).Info("OnStart钩子执行完成")
		}

	case *fxevent.OnStopExecuting:
		l.logger.WithFields(logrus.Fields{
			"caller": e.FunctionName,
			"callee": e.CallerName,
		}).Info("正在执行OnStop钩子")

	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logger.WithFields(logrus.Fields{
				"caller": e.FunctionName,
				"callee": e.CallerName,
				"error":  e.Err.Error(),
			}).Error("OnStop钩子执行失败")
		} else {
			l.logger.WithFields(logrus.Fields{
				"caller":  e.FunctionName,
				"callee":  e.CallerName,
				"runtime": e.Runtime.String(),
			}).Info("OnStop钩子执行完成")
		}

	case *fxevent.Supplied:
		if e.Err != nil {
			l.logger.WithFields(logrus.Fields{
				"type":  e.TypeName,
				"error": e.Err.Error(),
			}).Error("依赖供应失败")
		} else {
			l.logger.WithField("type", e.TypeName).Debug("依赖供应成功")
		}

	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.logger.WithFields(logrus.Fields{
				"constructor": e.ConstructorName,
				"type":        rtype,
			}).Debug("依赖注入成功")
		}

	case *fxevent.Invoking:
		l.logger.WithField("function", e.FunctionName).Debug("正在调用函数")

	case *fxevent.Invoked:
		if e.Err != nil {
			l.logger.WithFields(logrus.Fields{
				"function": e.FunctionName,
				"error":    e.Err.Error(),
			}).Error("函数调用失败")
		} else {
			l.logger.WithFields(logrus.Fields{
				"function": e.FunctionName,
				"trace":    e.Trace,
			}).Debug("函数调用完成")
		}

	case *fxevent.Stopping:
		l.logger.WithField("signal", e.Signal.String()).Info("接收到信号")

	case *fxevent.Stopped:
		if e.Err != nil {
			l.logger.WithField("error", e.Err.Error()).Error("停止失败")
		} else {
			l.logger.Info("已停止")
		}

	case *fxevent.RollingBack:
		l.logger.WithField("error", e.StartErr.Error()).Error("启动失败，正在回滚")

	case *fxevent.RolledBack:
		if e.Err != nil {
			l.logger.WithField("error", e.Err.Error()).Error("回滚失败")
		} else {
			l.logger.Info("回滚完成")
		}

	case *fxevent.Started:
		if e.Err != nil {
			l.logger.WithField("error", e.Err.Error()).Error("启动失败")
		} else {
			l.logger.Info("启动成功")
		}

	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.logger.WithField("error", e.Err.Error()).Error("日志器初始化失败")
		} else {
			l.logger.WithField("constructor", e.ConstructorName).Debug("日志器初始化成功")
		}

	default:
		l.logger.WithField("event", fmt.Sprintf("%T", e)).Debug("未知的fx事件")
	}
}
