package cron

import (
	"context"

	"github.com/zhoudm1743/go-flow/core/logger"
	"go.uber.org/fx"
)

// Module cron模块
var Module = fx.Options(
	// 提供组件
	fx.Provide(
		NewRedisRepository,
		NewCronScheduler,
		NewCronService,

		// 内置系统任务处理器
		NewLogHandler,
		NewCacheCleanHandler,
		NewDatabaseMaintenanceHandler,
	),

	// 启动调度器
	fx.Invoke(InitCronModule),
)

// InitCronModule 初始化cron模块
func InitCronModule(
	lifecycle fx.Lifecycle,
	scheduler Scheduler,
	service Service,
	log logger.Logger,
	logHandler *LogHandler,
	cacheCleanHandler *CacheCleanHandler,
	dbMaintenanceHandler *DatabaseMaintenanceHandler,
) {
	// 注册内置系统任务处理器
	RegisterSystemTaskHandler(logHandler)
	RegisterSystemTaskHandler(cacheCleanHandler)
	RegisterSystemTaskHandler(dbMaintenanceHandler)

	log.Info("已注册内置系统任务处理器", map[string]interface{}{
		"handlers": []string{
			logHandler.GetName(),
			cacheCleanHandler.GetName(),
			dbMaintenanceHandler.GetName(),
		},
	})

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 启动调度器
			if err := scheduler.Start(ctx); err != nil {
				log.WithError(err).Error("启动cron调度器失败")
				return err
			}

			log.Info("Cron模块启动成功")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// 停止调度器
			if err := scheduler.Stop(); err != nil {
				log.WithError(err).Error("停止cron调度器失败")
				return err
			}

			log.Info("Cron模块停止成功")
			return nil
		},
	})
}

// CreateSystemTasks 创建内置系统任务
func CreateSystemTasks(ctx context.Context, service Service, log logger.Logger) error {
	// 示例：创建一些预定义的系统任务

	// 1. 缓存清理任务 - 每天凌晨2点执行
	cacheCleanTask := &SystemTask{
		BaseTask: BaseTask{
			ID:          "system-cache-clean",
			Name:        "缓存清理任务",
			Type:        TaskTypeSystem,
			Cron:        "0 0 2 * * *", // 每天凌晨2点
			Status:      TaskStatusActive,
			Description: "定期清理过期缓存",
			CreatedBy:   "system",
		},
		Config: SystemConfig{
			HandlerName: "cache_clean",
			Parameters: map[string]interface{}{
				"pattern": "temp:*", // 只清理临时缓存
			},
			Timeout:    300, // 5分钟超时
			RetryCount: 2,
		},
		logger: log,
	}

	// 设置处理器
	if handler, exists := GetSystemTaskHandler("cache_clean"); exists {
		cacheCleanTask.handler = handler
	}

	// 2. 系统日志记录任务 - 每小时执行
	logTask := &SystemTask{
		BaseTask: BaseTask{
			ID:          "system-log-record",
			Name:        "系统日志记录",
			Type:        TaskTypeSystem,
			Cron:        "0 0 * * * *", // 每小时
			Status:      TaskStatusActive,
			Description: "记录系统运行状态日志",
			CreatedBy:   "system",
		},
		Config: SystemConfig{
			HandlerName: "log",
			Parameters: map[string]interface{}{
				"level":   "info",
				"message": "系统运行正常，定时检查完成",
			},
			Timeout:    60, // 1分钟超时
			RetryCount: 1,
		},
		logger: log,
	}

	// 设置处理器
	if handler, exists := GetSystemTaskHandler("log"); exists {
		logTask.handler = handler
	}

	// 注册任务到仓库（这里需要直接操作仓库，因为系统任务不通过API创建）
	log.Info("系统任务已准备就绪，可通过管理界面手动创建")

	return nil
}
