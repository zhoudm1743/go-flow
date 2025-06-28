package domain

import (
	"fmt"
	"time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventID() string
	EventType() string
	AggregateID() string
	AggregateType() string
	EventVersion() int
	OccurredOn() time.Time
	EventData() map[string]interface{}
}

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	eventID       string                 `json:"event_id"`
	eventType     string                 `json:"event_type"`
	aggregateID   string                 `json:"aggregate_id"`
	aggregateType string                 `json:"aggregate_type"`
	eventVersion  int                    `json:"event_version"`
	occurredOn    time.Time              `json:"occurred_on"`
	eventData     map[string]interface{} `json:"event_data"`
}

// NewBaseDomainEvent 创建基础领域事件
func NewBaseDomainEvent(eventType, aggregateID, aggregateType string, eventData map[string]interface{}) *BaseDomainEvent {
	return &BaseDomainEvent{
		eventID:       generateEventID(),
		eventType:     eventType,
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		eventVersion:  1,
		occurredOn:    time.Now(),
		eventData:     eventData,
	}
}

// EventID 获取事件ID
func (e *BaseDomainEvent) EventID() string {
	return e.eventID
}

// EventType 获取事件类型
func (e *BaseDomainEvent) EventType() string {
	return e.eventType
}

// AggregateID 获取聚合ID
func (e *BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

// AggregateType 获取聚合类型
func (e *BaseDomainEvent) AggregateType() string {
	return e.aggregateType
}

// EventVersion 获取事件版本
func (e *BaseDomainEvent) EventVersion() int {
	return e.eventVersion
}

// OccurredOn 获取发生时间
func (e *BaseDomainEvent) OccurredOn() time.Time {
	return e.occurredOn
}

// EventData 获取事件数据
func (e *BaseDomainEvent) EventData() map[string]interface{} {
	return e.eventData
}

// SetEventData 设置事件数据
func (e *BaseDomainEvent) SetEventData(data map[string]interface{}) {
	e.eventData = data
}

// DomainEventHandler 领域事件处理器接口
type DomainEventHandler interface {
	Handle(event DomainEvent) error
	CanHandle(eventType string) bool
}

// DomainEventDispatcher 领域事件分发器接口
type DomainEventDispatcher interface {
	Register(eventType string, handler DomainEventHandler)
	Dispatch(event DomainEvent) error
	DispatchAll(events []DomainEvent) error
}

// SimpleDomainEventDispatcher 简单领域事件分发器
type SimpleDomainEventDispatcher struct {
	handlers map[string][]DomainEventHandler
}

// NewSimpleDomainEventDispatcher 创建简单领域事件分发器
func NewSimpleDomainEventDispatcher() *SimpleDomainEventDispatcher {
	return &SimpleDomainEventDispatcher{
		handlers: make(map[string][]DomainEventHandler),
	}
}

// Register 注册事件处理器
func (d *SimpleDomainEventDispatcher) Register(eventType string, handler DomainEventHandler) {
	if d.handlers[eventType] == nil {
		d.handlers[eventType] = make([]DomainEventHandler, 0)
	}
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Dispatch 分发单个事件
func (d *SimpleDomainEventDispatcher) Dispatch(event DomainEvent) error {
	handlers, exists := d.handlers[event.EventType()]
	if !exists {
		return nil // 没有处理器，忽略事件
	}

	for _, handler := range handlers {
		if handler.CanHandle(event.EventType()) {
			if err := handler.Handle(event); err != nil {
				return fmt.Errorf("failed to handle event %s: %w", event.EventType(), err)
			}
		}
	}

	return nil
}

// DispatchAll 分发多个事件
func (d *SimpleDomainEventDispatcher) DispatchAll(events []DomainEvent) error {
	for _, event := range events {
		if err := d.Dispatch(event); err != nil {
			return err
		}
	}
	return nil
}

// AsyncDomainEventDispatcher 异步领域事件分发器
type AsyncDomainEventDispatcher struct {
	SimpleDomainEventDispatcher
	eventChan chan DomainEvent
	workers   int
}

// NewAsyncDomainEventDispatcher 创建异步领域事件分发器
func NewAsyncDomainEventDispatcher(workers int, bufferSize int) *AsyncDomainEventDispatcher {
	dispatcher := &AsyncDomainEventDispatcher{
		SimpleDomainEventDispatcher: *NewSimpleDomainEventDispatcher(),
		eventChan:                   make(chan DomainEvent, bufferSize),
		workers:                     workers,
	}

	// 启动工作goroutine
	for i := 0; i < workers; i++ {
		go dispatcher.worker()
	}

	return dispatcher
}

// worker 工作goroutine
func (d *AsyncDomainEventDispatcher) worker() {
	for event := range d.eventChan {
		// 异步处理事件，错误只记录，不阻塞
		if err := d.SimpleDomainEventDispatcher.Dispatch(event); err != nil {
			// 这里应该记录日志，但由于没有logger依赖，暂时忽略
			fmt.Printf("Error handling event %s: %v\n", event.EventType(), err)
		}
	}
}

// Dispatch 异步分发事件
func (d *AsyncDomainEventDispatcher) Dispatch(event DomainEvent) error {
	select {
	case d.eventChan <- event:
		return nil
	default:
		return fmt.Errorf("event channel is full, cannot dispatch event %s", event.EventType())
	}
}

// DispatchAll 异步分发多个事件
func (d *AsyncDomainEventDispatcher) DispatchAll(events []DomainEvent) error {
	for _, event := range events {
		if err := d.Dispatch(event); err != nil {
			return err
		}
	}
	return nil
}

// Close 关闭分发器
func (d *AsyncDomainEventDispatcher) Close() {
	close(d.eventChan)
}

// IntegrationEvent 集成事件接口
type IntegrationEvent interface {
	DomainEvent
	ExternalSystemID() string
	MessageType() string
}

// BaseIntegrationEvent 基础集成事件
type BaseIntegrationEvent struct {
	BaseDomainEvent
	externalSystemID string `json:"external_system_id"`
	messageType      string `json:"message_type"`
}

// NewBaseIntegrationEvent 创建基础集成事件
func NewBaseIntegrationEvent(eventType, aggregateID, aggregateType, externalSystemID, messageType string, eventData map[string]interface{}) *BaseIntegrationEvent {
	return &BaseIntegrationEvent{
		BaseDomainEvent:  *NewBaseDomainEvent(eventType, aggregateID, aggregateType, eventData),
		externalSystemID: externalSystemID,
		messageType:      messageType,
	}
}

// ExternalSystemID 获取外部系统ID
func (e *BaseIntegrationEvent) ExternalSystemID() string {
	return e.externalSystemID
}

// MessageType 获取消息类型
func (e *BaseIntegrationEvent) MessageType() string {
	return e.messageType
}

// 生成事件ID的辅助函数
func generateEventID() string {
	// 简单的ID生成，实际项目中可能使用UUID
	return fmt.Sprintf("event_%d", time.Now().UnixNano())
}

// 常用事件类型常量
const (
	EventTypeEntityCreated = "entity.created"
	EventTypeEntityUpdated = "entity.updated"
	EventTypeEntityDeleted = "entity.deleted"
)

// EntityCreatedEvent 实体创建事件
type EntityCreatedEvent struct {
	BaseDomainEvent
}

// NewEntityCreatedEvent 创建实体创建事件
func NewEntityCreatedEvent(aggregateID, aggregateType string, data map[string]interface{}) *EntityCreatedEvent {
	return &EntityCreatedEvent{
		BaseDomainEvent: *NewBaseDomainEvent(EventTypeEntityCreated, aggregateID, aggregateType, data),
	}
}

// EntityUpdatedEvent 实体更新事件
type EntityUpdatedEvent struct {
	BaseDomainEvent
	Changes map[string]interface{} `json:"changes"`
}

// NewEntityUpdatedEvent 创建实体更新事件
func NewEntityUpdatedEvent(aggregateID, aggregateType string, changes map[string]interface{}) *EntityUpdatedEvent {
	eventData := map[string]interface{}{
		"changes": changes,
	}
	return &EntityUpdatedEvent{
		BaseDomainEvent: *NewBaseDomainEvent(EventTypeEntityUpdated, aggregateID, aggregateType, eventData),
		Changes:         changes,
	}
}

// EntityDeletedEvent 实体删除事件
type EntityDeletedEvent struct {
	BaseDomainEvent
}

// NewEntityDeletedEvent 创建实体删除事件
func NewEntityDeletedEvent(aggregateID, aggregateType string) *EntityDeletedEvent {
	return &EntityDeletedEvent{
		BaseDomainEvent: *NewBaseDomainEvent(EventTypeEntityDeleted, aggregateID, aggregateType, map[string]interface{}{}),
	}
}
