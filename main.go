// Package main Go-Flow API服务
// @title Go-Flow API
// @version 1.0
// @description 这是一个基于Go和Gin的企业级后台管理系统API
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /api/v1
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer token for authentication
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/zhoudm1743/go-flow/boot"
	"github.com/zhoudm1743/go-flow/core/logger"

	_ "github.com/zhoudm1743/go-flow/docs" // swagger docs
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
