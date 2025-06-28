package domain

import (
	"go.uber.org/fx"

	"github.com/zhoudm1743/go-flow/core/database"
	"github.com/zhoudm1743/go-flow/core/logger"
)

// Module DDD领域模块
var Module = fx.Options(
	// 提供领域事件分发器
	fx.Provide(
		fx.Annotate(
			NewDefaultAsyncDomainEventDispatcher,
			fx.As(new(DomainEventDispatcher)),
		),
	),

	// 提供工作单元
	fx.Provide(
		func(db database.Database, dispatcher DomainEventDispatcher) UnitOfWork {
			return NewSimpleUnitOfWork(db.GetDB(), dispatcher)
		},
	),

	// 提供领域服务基础设施
	fx.Provide(
		NewDomainServiceRegistry,
	),

	// 注册内置领域事件处理器
	fx.Provide(
		NewSystemDomainEventHandlers,
	),

	// 初始化领域模块
	fx.Invoke(InitDomainModule),
)

// NewDefaultAsyncDomainEventDispatcher 创建默认异步领域事件分发器
func NewDefaultAsyncDomainEventDispatcher() DomainEventDispatcher {
	// 创建异步分发器，10个工作协程，1000个事件缓冲
	return NewAsyncDomainEventDispatcher(10, 1000)
}

// DomainServiceRegistry 领域服务注册表
type DomainServiceRegistry struct {
	services map[string]DomainService
	logger   logger.Logger
}

// NewDomainServiceRegistry 创建领域服务注册表
func NewDomainServiceRegistry(log logger.Logger) *DomainServiceRegistry {
	return &DomainServiceRegistry{
		services: make(map[string]DomainService),
		logger:   log,
	}
}

// Register 注册领域服务
func (r *DomainServiceRegistry) Register(service DomainService) {
	r.services[service.Name()] = service
	r.logger.WithField("service_name", service.Name()).Info("领域服务注册成功")
}

// Get 获取领域服务
func (r *DomainServiceRegistry) Get(name string) (DomainService, bool) {
	service, exists := r.services[name]
	return service, exists
}

// GetAll 获取所有领域服务
func (r *DomainServiceRegistry) GetAll() map[string]DomainService {
	return r.services
}

// SystemDomainEventHandlers 系统领域事件处理器
type SystemDomainEventHandlers struct {
	logger logger.Logger
}

// NewSystemDomainEventHandlers 创建系统领域事件处理器
func NewSystemDomainEventHandlers(log logger.Logger) *SystemDomainEventHandlers {
	return &SystemDomainEventHandlers{
		logger: log,
	}
}

// EntityCreatedHandler 实体创建事件处理器
func (h *SystemDomainEventHandlers) EntityCreatedHandler() DomainEventHandler {
	return &entityCreatedHandler{logger: h.logger}
}

// EntityUpdatedHandler 实体更新事件处理器
func (h *SystemDomainEventHandlers) EntityUpdatedHandler() DomainEventHandler {
	return &entityUpdatedHandler{logger: h.logger}
}

// EntityDeletedHandler 实体删除事件处理器
func (h *SystemDomainEventHandlers) EntityDeletedHandler() DomainEventHandler {
	return &entityDeletedHandler{logger: h.logger}
}

// 实体创建事件处理器
type entityCreatedHandler struct {
	logger logger.Logger
}

func (h *entityCreatedHandler) Handle(event DomainEvent) error {
	h.logger.WithFields(map[string]interface{}{
		"event_id":       event.EventID(),
		"event_type":     event.EventType(),
		"aggregate_id":   event.AggregateID(),
		"aggregate_type": event.AggregateType(),
		"occurred_on":    event.OccurredOn(),
	}).Info("处理实体创建事件")
	return nil
}

func (h *entityCreatedHandler) CanHandle(eventType string) bool {
	return eventType == EventTypeEntityCreated
}

// 实体更新事件处理器
type entityUpdatedHandler struct {
	logger logger.Logger
}

func (h *entityUpdatedHandler) Handle(event DomainEvent) error {
	h.logger.WithFields(map[string]interface{}{
		"event_id":       event.EventID(),
		"event_type":     event.EventType(),
		"aggregate_id":   event.AggregateID(),
		"aggregate_type": event.AggregateType(),
		"occurred_on":    event.OccurredOn(),
	}).Info("处理实体更新事件")
	return nil
}

func (h *entityUpdatedHandler) CanHandle(eventType string) bool {
	return eventType == EventTypeEntityUpdated
}

// 实体删除事件处理器
type entityDeletedHandler struct {
	logger logger.Logger
}

func (h *entityDeletedHandler) Handle(event DomainEvent) error {
	h.logger.WithFields(map[string]interface{}{
		"event_id":       event.EventID(),
		"event_type":     event.EventType(),
		"aggregate_id":   event.AggregateID(),
		"aggregate_type": event.AggregateType(),
		"occurred_on":    event.OccurredOn(),
	}).Info("处理实体删除事件")
	return nil
}

func (h *entityDeletedHandler) CanHandle(eventType string) bool {
	return eventType == EventTypeEntityDeleted
}

// InitDomainModule 初始化领域模块
func InitDomainModule(
	dispatcher DomainEventDispatcher,
	handlers *SystemDomainEventHandlers,
	registry *DomainServiceRegistry,
	log logger.Logger,
) {
	// 注册系统领域事件处理器
	dispatcher.Register(EventTypeEntityCreated, handlers.EntityCreatedHandler())
	dispatcher.Register(EventTypeEntityUpdated, handlers.EntityUpdatedHandler())
	dispatcher.Register(EventTypeEntityDeleted, handlers.EntityDeletedHandler())

	log.Info("DDD领域模块初始化完成", map[string]interface{}{
		"event_handlers": []string{
			EventTypeEntityCreated,
			EventTypeEntityUpdated,
			EventTypeEntityDeleted,
		},
	})
}
