package boot

import (
	"context"

	"github.com/zhoudm1743/go-flow/app/admin"
	"github.com/zhoudm1743/go-flow/core/cache"
	"github.com/zhoudm1743/go-flow/core/config"
	"github.com/zhoudm1743/go-flow/core/cron"
	"github.com/zhoudm1743/go-flow/core/database"
	"github.com/zhoudm1743/go-flow/core/domain"
	"github.com/zhoudm1743/go-flow/core/http"
	"github.com/zhoudm1743/go-flow/core/logger"
	"github.com/zhoudm1743/go-flow/pkg"
	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	logger.Module,
	database.Module,
	cache.Module,
	cron.Module,
	domain.Module,
	pkg.Module,
	http.Module,

	// 应用模块
	admin.Module,

	// 启动
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	cfg *config.Config,
	log logger.Logger,
	db database.Database,
	cacheClient cache.Cache,
) {
	// 设置全局日志实例
	logger.SetGlobalLogger(log)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.WithField("level", cfg.Log.Level).Info("日志服务已初始化")

			// 测试数据库连接
			if err := db.Ping(); err != nil {
				log.WithError(err).Error("数据库连接测试失败")
				return err
			}
			log.Info("数据库连接测试成功")

			// 测试 Redis 连接
			if err := cacheClient.PingCtx(ctx); err != nil {
				log.WithError(err).Error("Redis 连接测试失败")
				return err
			}
			log.Info("Redis 连接测试成功")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("应用停止")

			// 关闭数据库连接
			if err := db.Close(); err != nil {
				log.WithError(err).Error("关闭数据库连接失败")
			} else {
				log.Info("数据库连接已关闭")
			}

			// 关闭 Redis 连接
			if err := cacheClient.Close(); err != nil {
				log.WithError(err).Error("关闭 Redis 连接失败")
			} else {
				log.Info("Redis 连接已关闭")
			}

			return nil
		},
	})
}
