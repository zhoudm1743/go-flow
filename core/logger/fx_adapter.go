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
		}).Info("OnStart hook executing")

	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logger.WithFields(logrus.Fields{
				"caller": e.FunctionName,
				"callee": e.CallerName,
				"error":  e.Err.Error(),
			}).Error("OnStart hook failed")
		} else {
			l.logger.WithFields(logrus.Fields{
				"caller":  e.FunctionName,
				"callee":  e.CallerName,
				"runtime": e.Runtime.String(),
			}).Info("OnStart hook executed")
		}

	case *fxevent.OnStopExecuting:
		l.logger.WithFields(logrus.Fields{
			"caller": e.FunctionName,
			"callee": e.CallerName,
		}).Info("OnStop hook executing")

	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logger.WithFields(logrus.Fields{
				"caller": e.FunctionName,
				"callee": e.CallerName,
				"error":  e.Err.Error(),
			}).Error("OnStop hook failed")
		} else {
			l.logger.WithFields(logrus.Fields{
				"caller":  e.FunctionName,
				"callee":  e.CallerName,
				"runtime": e.Runtime.String(),
			}).Info("OnStop hook executed")
		}

	case *fxevent.Supplied:
		if e.Err != nil {
			l.logger.WithFields(logrus.Fields{
				"type":  e.TypeName,
				"error": e.Err.Error(),
			}).Error("Supply failed")
		} else {
			l.logger.WithField("type", e.TypeName).Debug("Supplied")
		}

	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.logger.WithFields(logrus.Fields{
				"constructor": e.ConstructorName,
				"type":        rtype,
			}).Debug("Provided")
		}

	case *fxevent.Invoking:
		l.logger.WithField("function", e.FunctionName).Debug("Invoking")

	case *fxevent.Invoked:
		if e.Err != nil {
			l.logger.WithFields(logrus.Fields{
				"function": e.FunctionName,
				"error":    e.Err.Error(),
			}).Error("Invoke failed")
		} else {
			l.logger.WithFields(logrus.Fields{
				"function": e.FunctionName,
				"trace":    e.Trace,
			}).Debug("Invoked")
		}

	case *fxevent.Stopping:
		l.logger.WithField("signal", e.Signal.String()).Info("Received signal")

	case *fxevent.Stopped:
		if e.Err != nil {
			l.logger.WithField("error", e.Err.Error()).Error("Stop failed")
		} else {
			l.logger.Info("Stopped")
		}

	case *fxevent.RollingBack:
		l.logger.WithField("error", e.StartErr.Error()).Error("Start failed, rolling back")

	case *fxevent.RolledBack:
		if e.Err != nil {
			l.logger.WithField("error", e.Err.Error()).Error("Rollback failed")
		} else {
			l.logger.Info("Rolled back")
		}

	case *fxevent.Started:
		if e.Err != nil {
			l.logger.WithField("error", e.Err.Error()).Error("Start failed")
		} else {
			l.logger.Info("Started")
		}

	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.logger.WithField("error", e.Err.Error()).Error("Logger initialization failed")
		} else {
			l.logger.WithField("constructor", e.ConstructorName).Debug("Logger initialized")
		}

	default:
		l.logger.WithField("event", fmt.Sprintf("%T", e)).Debug("Unknown fx event")
	}
}
