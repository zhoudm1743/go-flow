package event

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"

	"github.com/zhoudm1743/go-flow/core/cache"
	"github.com/zhoudm1743/go-flow/core/logger"
)

// Module 事件模块
var Module = fx.Module("event",
	// 提供事件总线
	fx.Provide(
		fx.Annotate(
			NewDefaultEventBus,
			fx.As(new(EventBus)),
		),
		func() *EventBusOptions {
			return DefaultEventBusOptions()
		},
	),

	// 提供事件存储
	fx.Provide(
		func(cache cache.Cache, logger logger.Logger) EventStore {
			return NewRedisEventStore(cache.GetClient(), logger)
		},
	),

	// 提供延时事件调度器
	fx.Provide(
		fx.Annotate(
			NewCronDelayedEventScheduler,
			fx.As(new(DelayedEventScheduler)),
		),
		func() *SchedulerOptions {
			return DefaultSchedulerOptions()
		},
	),

	// 提供事件服务
	fx.Provide(
		fx.Annotate(
			NewDefaultEventService,
			fx.As(new(EventService)),
		),
	),

	// 提供Cron集成
	fx.Provide(
		NewCronEventIntegration,
	),

	// 提供系统事件处理器
	fx.Provide(
		NewSystemEventHandlers,
	),

	// 生命周期管理
	fx.Invoke(RegisterLifecycle),

	// 初始化系统事件处理器
	fx.Invoke(InitializeSystemEventHandlers),
)

// SystemEventHandlers 系统事件处理器集合
type SystemEventHandlers struct {
	UserHandler   *UserEventHandler
	SystemHandler *SystemEventHandler
	DataHandler   *DataEventHandler
}

// NewSystemEventHandlers 创建系统事件处理器
func NewSystemEventHandlers(logger logger.Logger) *SystemEventHandlers {
	return &SystemEventHandlers{
		UserHandler:   NewUserEventHandler(logger),
		SystemHandler: NewSystemEventHandler(logger),
		DataHandler:   NewDataEventHandler(logger),
	}
}

// RegisterLifecycle 注册生命周期
func RegisterLifecycle(
	lc fx.Lifecycle,
	eventBus EventBus,
	scheduler DelayedEventScheduler,
	cronIntegration *CronEventIntegration,
	logger logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("正在启动事件系统...")

			// 启动事件总线
			if err := eventBus.Start(ctx); err != nil {
				return err
			}

			// 启动延时事件调度器
			if err := scheduler.Start(ctx); err != nil {
				return err
			}

			// 初始化Cron集成
			if err := cronIntegration.Initialize(ctx); err != nil {
				return err
			}

			logger.Info("事件系统启动成功")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("正在停止事件系统...")

			// 创建停止超时上下文
			stopCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			// 停止延时事件调度器
			if err := scheduler.Stop(stopCtx); err != nil {
				logger.WithError(err).Warn("停止延时事件调度器失败")
			}

			// 停止事件总线
			if err := eventBus.Stop(stopCtx); err != nil {
				logger.WithError(err).Warn("停止事件总线失败")
			}

			logger.Info("事件系统已停止")
			return nil
		},
	})
}

// InitializeSystemEventHandlers 初始化系统事件处理器
func InitializeSystemEventHandlers(
	eventService EventService,
	handlers *SystemEventHandlers,
	logger logger.Logger,
) error {
	logger.Info("正在初始化系统事件处理器...")

	// 订阅用户事件
	userHandlers := []EventHandler{
		handlers.UserHandler,
	}

	for _, handler := range userHandlers {
		if err := eventService.Subscribe(handler.GetEventType(), handler); err != nil {
			return err
		}
	}

	// 订阅系统事件
	systemHandlers := []EventHandler{
		handlers.SystemHandler,
	}

	for _, handler := range systemHandlers {
		if err := eventService.Subscribe(handler.GetEventType(), handler); err != nil {
			return err
		}
	}

	// 订阅数据事件
	dataHandlers := []EventHandler{
		handlers.DataHandler,
	}

	for _, handler := range dataHandlers {
		if err := eventService.Subscribe(handler.GetEventType(), handler); err != nil {
			return err
		}
	}

	logger.Info("系统事件处理器初始化成功")
	return nil
}

// UserEventHandler 用户事件处理器
type UserEventHandler struct {
	logger logger.Logger
}

// NewUserEventHandler 创建用户事件处理器
func NewUserEventHandler(logger logger.Logger) *UserEventHandler {
	return &UserEventHandler{
		logger: logger,
	}
}

// Handle 处理事件
func (h *UserEventHandler) Handle(ctx context.Context, event Event) error {
	h.logger.WithField("event_type", event.GetType()).
		WithField("event_id", event.GetID()).
		WithField("source", event.GetSource()).
		Info("正在处理用户事件")

	// 处理用户相关事件的逻辑
	switch event.GetType() {
	case EventTypeUserLogin:
		return h.handleUserLogin(ctx, event)
	case EventTypeUserLogout:
		return h.handleUserLogout(ctx, event)
	case EventTypeUserRegistered:
		return h.handleUserRegistered(ctx, event)
	default:
		return fmt.Errorf("未知的用户事件类型: %s", event.GetType())
	}
}

// handleUserLogin 处理用户登录事件
func (h *UserEventHandler) handleUserLogin(ctx context.Context, event Event) error {
	// 实现用户登录事件处理逻辑
	h.logger.WithField("event_id", event.GetID()).Info("用户登录事件处理完成")
	return nil
}

// handleUserLogout 处理用户登出事件
func (h *UserEventHandler) handleUserLogout(ctx context.Context, event Event) error {
	// 实现用户登出事件处理逻辑
	h.logger.WithField("event_id", event.GetID()).Info("用户登出事件处理完成")
	return nil
}

// handleUserRegistered 处理用户注册事件
func (h *UserEventHandler) handleUserRegistered(ctx context.Context, event Event) error {
	// 实现用户注册事件处理逻辑
	h.logger.WithField("event_id", event.GetID()).Info("用户注册事件处理完成")
	return nil
}

// GetEventType 获取处理的事件类型
func (h *UserEventHandler) GetEventType() string {
	return EventTypeUserLogin // 实际实现中可能需要订阅多个事件类型
}

// GetName 获取处理器名称
func (h *UserEventHandler) GetName() string {
	return "UserEventHandler"
}

// SystemEventHandler 系统事件处理器
type SystemEventHandler struct {
	logger logger.Logger
}

// NewSystemEventHandler 创建系统事件处理器
func NewSystemEventHandler(logger logger.Logger) *SystemEventHandler {
	return &SystemEventHandler{
		logger: logger,
	}
}

// Handle 处理事件
func (h *SystemEventHandler) Handle(ctx context.Context, event Event) error {
	h.logger.WithField("event_type", event.GetType()).
		WithField("event_id", event.GetID()).
		WithField("source", event.GetSource()).
		Info("正在处理系统事件")

	// 处理系统相关事件的逻辑
	switch event.GetType() {
	case EventTypeSystemStartup:
		return h.handleSystemStartup(ctx, event)
	case EventTypeSystemShutdown:
		return h.handleSystemShutdown(ctx, event)
	case EventTypeSystemError:
		return h.handleSystemError(ctx, event)
	case EventTypeSystemWarning:
		return h.handleSystemWarning(ctx, event)
	default:
		return fmt.Errorf("未知的系统事件类型: %s", event.GetType())
	}
}

// handleSystemStartup 处理系统启动事件
func (h *SystemEventHandler) handleSystemStartup(ctx context.Context, event Event) error {
	// 实现系统启动事件处理逻辑
	h.logger.WithField("event_id", event.GetID()).Info("系统启动事件处理完成")
	return nil
}

// handleSystemShutdown 处理系统关闭事件
func (h *SystemEventHandler) handleSystemShutdown(ctx context.Context, event Event) error {
	// 实现系统关闭事件处理逻辑
	h.logger.WithField("event_id", event.GetID()).Info("系统关闭事件处理完成")
	return nil
}

// handleSystemError 处理系统错误事件
func (h *SystemEventHandler) handleSystemError(ctx context.Context, event Event) error {
	// 实现系统错误事件处理逻辑
	h.logger.WithField("event_id", event.GetID()).Error("系统错误事件处理完成")
	return nil
}

// handleSystemWarning 处理系统警告事件
func (h *SystemEventHandler) handleSystemWarning(ctx context.Context, event Event) error {
	// 实现系统警告事件处理逻辑
	h.logger.WithField("event_id", event.GetID()).Warn("系统警告事件处理完成")
	return nil
}

// GetEventType 获取处理的事件类型
func (h *SystemEventHandler) GetEventType() string {
	return EventTypeSystemStartup // 实际实现中可能需要订阅多个事件类型
}

// GetName 获取处理器名称
func (h *SystemEventHandler) GetName() string {
	return "SystemEventHandler"
}

// DataEventHandler 数据事件处理器
type DataEventHandler struct {
	logger logger.Logger
}

// NewDataEventHandler 创建数据事件处理器
func NewDataEventHandler(logger logger.Logger) *DataEventHandler {
	return &DataEventHandler{
		logger: logger,
	}
}

// Handle 处理事件
func (h *DataEventHandler) Handle(ctx context.Context, event Event) error {
	h.logger.WithField("event_type", event.GetType()).
		WithField("event_id", event.GetID()).
		WithField("source", event.GetSource()).
		Info("正在处理数据事件")

	// 处理数据相关事件的逻辑
	switch event.GetType() {
	case EventTypeDataCreated:
		return h.handleDataCreated(ctx, event)
	case EventTypeDataUpdated:
		return h.handleDataUpdated(ctx, event)
	case EventTypeDataDeleted:
		return h.handleDataDeleted(ctx, event)
	default:
		return fmt.Errorf("未知的数据事件类型: %s", event.GetType())
	}
}

// handleDataCreated 处理数据创建事件
func (h *DataEventHandler) handleDataCreated(ctx context.Context, event Event) error {
	// 实现数据创建事件处理逻辑
	h.logger.WithField("event_id", event.GetID()).Info("数据创建事件处理完成")
	return nil
}

// handleDataUpdated 处理数据更新事件
func (h *DataEventHandler) handleDataUpdated(ctx context.Context, event Event) error {
	// 实现数据更新事件处理逻辑
	h.logger.WithField("event_id", event.GetID()).Info("数据更新事件处理完成")
	return nil
}

// handleDataDeleted 处理数据删除事件
func (h *DataEventHandler) handleDataDeleted(ctx context.Context, event Event) error {
	// 实现数据删除事件处理逻辑
	h.logger.WithField("event_id", event.GetID()).Info("数据删除事件处理完成")
	return nil
}

// GetEventType 获取处理的事件类型
func (h *DataEventHandler) GetEventType() string {
	return EventTypeDataCreated // 实际实现中可能需要订阅多个事件类型
}

// GetName 获取处理器名称
func (h *DataEventHandler) GetName() string {
	return "DataEventHandler"
}
