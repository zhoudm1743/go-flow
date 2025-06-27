package main

import (
	"github.com/sirupsen/logrus"
	"github.com/zhoudm1743/go-flow/boot"
	"github.com/zhoudm1743/go-flow/core/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	fx.New(
		fx.WithLogger(func() fxevent.Logger {
			fxLogger := logrus.New()
			fxLogger.SetLevel(logrus.InfoLevel)
			fxLogger.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05",
				ForceColors:     true,
				DisableColors:   false,
				PadLevelText:    true,
				DisableQuote:    true,
			})
			return logger.NewFxLogger(fxLogger)
		}),
		boot.Module,
	).Run()
}
