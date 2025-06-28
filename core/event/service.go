package event

import (
	"context"
	"fmt"
	"time"

	"github.com/zhoudm1743/go-flow/core/logger"
	"github.com/zhoudm1743/go-flow/pkg/response"
)

// EventService 事件服务接口
type EventService interface {
	// PublishEvent 发布即时事件
	PublishEvent(ctx context.Context, eventType string, payload interface{}, source string) error
	// PublishEventAsync 异步发布即时事件
	PublishEventAsync(ctx context.Context, eventType string, payload interface{}, source string) error
	// ScheduleDelayedEvent 调度延时事件
	ScheduleDelayedEvent(ctx context.Context, eventType string, payload interface{}, source string, delayUntil time.Time, options *DelayedEventOptions) (*DelayedEvent, error)
	// CancelDelayedEvent 取消延时事件
	CancelDelayedEvent(ctx context.Context, eventID string) error
	// GetDelayedEvent 获取延时事件详情
	GetDelayedEvent(ctx context.Context, eventID string) (*DelayedEvent, error)
	// GetDelayedEvents 获取延时事件列表
	GetDelayedEvents(ctx context.Context, filter *EventFilter) (*response.PageResult[*DelayedEvent], error)
	// GetEventRecords 获取事件记录
	GetEventRecords(ctx context.Context, filter *EventFilter) (*response.PageResult[*EventRecord], error)
	// Subscribe 订阅事件
	Subscribe(eventType string, handler EventHandler, options ...*SubscriptionOptions) error
	// Unsubscribe 取消订阅
	Unsubscribe(eventType string, handlerName string) error
	// GetMetrics 获取事件指标
	GetMetrics(ctx context.Context) (*EventMetrics, error)
}

// DelayedEventOptions 延时事件选项
type DelayedEventOptions struct {
	Priority   EventPriority `json:"priority"`
	MaxRetries int           `json:"max_retries"`
}

// DefaultDelayedEventOptions 默认延时事件选项
func DefaultDelayedEventOptions() *DelayedEventOptions {
	return &DelayedEventOptions{
		Priority:   EventPriorityNormal,
		MaxRetries: 3,
	}
}

// DefaultEventService 默认事件服务实现
type DefaultEventService struct {
	eventBus   EventBus
	eventStore EventStore
	scheduler  DelayedEventScheduler
	logger     logger.Logger
}

// NewDefaultEventService 创建默认事件服务
func NewDefaultEventService(
	eventBus EventBus,
	eventStore EventStore,
	scheduler DelayedEventScheduler,
	logger logger.Logger,
) *DefaultEventService {
	return &DefaultEventService{
		eventBus:   eventBus,
		eventStore: eventStore,
		scheduler:  scheduler,
		logger:     logger,
	}
}

// PublishEvent 发布即时事件
func (s *DefaultEventService) PublishEvent(ctx context.Context, eventType string, payload interface{}, source string) error {
	event := NewBaseEvent(eventType, payload, source)

	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.WithField("event_type", eventType).
			WithField("event_id", event.ID).
			WithError(err).Error("发布事件失败")
		return err
	}

	s.logger.WithField("event_type", eventType).
		WithField("event_id", event.ID).
		Debug("事件发布成功")

	return nil
}

// PublishEventAsync 异步发布即时事件
func (s *DefaultEventService) PublishEventAsync(ctx context.Context, eventType string, payload interface{}, source string) error {
	event := NewBaseEvent(eventType, payload, source)

	if err := s.eventBus.PublishAsync(ctx, event); err != nil {
		s.logger.WithField("event_type", eventType).
			WithField("event_id", event.ID).
			WithError(err).Error("异步发布事件失败")
		return err
	}

	s.logger.WithField("event_type", eventType).
		WithField("event_id", event.ID).
		Debug("事件异步发布成功")

	return nil
}

// ScheduleDelayedEvent 调度延时事件
func (s *DefaultEventService) ScheduleDelayedEvent(ctx context.Context, eventType string, payload interface{}, source string, delayUntil time.Time, options *DelayedEventOptions) (*DelayedEvent, error) {
	if options == nil {
		options = DefaultDelayedEventOptions()
	}

	// 检查延时时间是否合理
	if delayUntil.Before(time.Now()) {
		return nil, fmt.Errorf("延时时间不能是过去的时间")
	}

	// 创建延时事件
	delayedEvent := NewDelayedEvent(eventType, payload, source, delayUntil)
	delayedEvent.Priority = options.Priority
	delayedEvent.MaxRetries = options.MaxRetries

	// 调度事件
	if err := s.scheduler.ScheduleDelayedEvent(ctx, delayedEvent); err != nil {
		s.logger.WithField("event_id", delayedEvent.ID).
			WithField("event_type", eventType).
			WithField("delay_until", delayUntil).
			WithError(err).Error("调度延时事件失败")
		return nil, fmt.Errorf("调度延时事件失败: %w", err)
	}

	s.logger.WithField("event_id", delayedEvent.ID).
		WithField("event_type", eventType).
		WithField("delay_until", delayUntil).
		Info("延时事件调度成功")

	return delayedEvent, nil
}

// CancelDelayedEvent 取消延时事件
func (s *DefaultEventService) CancelDelayedEvent(ctx context.Context, eventID string) error {
	if err := s.scheduler.CancelDelayedEvent(ctx, eventID); err != nil {
		s.logger.WithField("event_id", eventID).
			WithError(err).Error("取消延时事件失败")
		return fmt.Errorf("取消延时事件失败: %w", err)
	}

	s.logger.WithField("event_id", eventID).Info("延时事件取消成功")

	return nil
}

// GetDelayedEvent 获取延时事件详情
func (s *DefaultEventService) GetDelayedEvent(ctx context.Context, eventID string) (*DelayedEvent, error) {
	event, err := s.eventStore.GetDelayedEvent(ctx, eventID)
	if err != nil {
		s.logger.WithField("event_id", eventID).
			WithError(err).Error("获取延时事件失败")
		return nil, fmt.Errorf("获取延时事件失败: %w", err)
	}

	return event, nil
}

// GetDelayedEvents 获取延时事件列表
func (s *DefaultEventService) GetDelayedEvents(ctx context.Context, filter *EventFilter) (*response.PageResult[*DelayedEvent], error) {
	// 设置默认分页参数
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	result, err := s.eventStore.GetDelayedEvents(ctx, filter)
	if err != nil {
		s.logger.WithField("page", filter.Page).
			WithField("page_size", filter.PageSize).
			WithError(err).Error("获取延时事件列表失败")
		return nil, fmt.Errorf("获取延时事件列表失败: %w", err)
	}

	return result, nil
}

// GetEventRecords 获取事件记录
func (s *DefaultEventService) GetEventRecords(ctx context.Context, filter *EventFilter) (*response.PageResult[*EventRecord], error) {
	// 设置默认分页参数
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	result, err := s.eventStore.GetEventRecords(ctx, filter)
	if err != nil {
		s.logger.WithField("page", filter.Page).
			WithField("page_size", filter.PageSize).
			WithError(err).Error("获取事件记录列表失败")
		return nil, fmt.Errorf("获取事件记录失败: %w", err)
	}

	return result, nil
}

// Subscribe 订阅事件
func (s *DefaultEventService) Subscribe(eventType string, handler EventHandler, options ...*SubscriptionOptions) error {
	if err := s.eventBus.Subscribe(eventType, handler, options...); err != nil {
		s.logger.WithField("event_type", eventType).
			WithField("handler", handler.GetName()).
			WithError(err).Error("订阅事件失败")
		return fmt.Errorf("订阅事件失败: %w", err)
	}

	s.logger.WithField("event_type", eventType).
		WithField("handler", handler.GetName()).
		Info("成功订阅事件")

	return nil
}

// Unsubscribe 取消订阅
func (s *DefaultEventService) Unsubscribe(eventType string, handlerName string) error {
	if err := s.eventBus.Unsubscribe(eventType, handlerName); err != nil {
		s.logger.WithField("event_type", eventType).
			WithField("handler", handlerName).
			WithError(err).Error("取消订阅事件失败")
		return fmt.Errorf("取消订阅事件失败: %w", err)
	}

	s.logger.WithField("event_type", eventType).
		WithField("handler", handlerName).
		Info("成功取消订阅事件")

	return nil
}

// GetMetrics 获取事件指标
func (s *DefaultEventService) GetMetrics(ctx context.Context) (*EventMetrics, error) {
	// 获取事件总线指标
	busMetrics := s.eventBus.GetMetrics()

	// 获取存储指标
	storeMetrics, err := s.eventStore.GetMetrics(ctx)
	if err != nil {
		s.logger.WithError(err).Error("获取存储指标失败")
		// 返回总线指标，忽略存储指标错误
		return busMetrics, nil
	}

	// 合并指标
	combinedMetrics := &EventMetrics{
		TotalEvents:     busMetrics.TotalEvents + storeMetrics.TotalEvents,
		PendingEvents:   busMetrics.PendingEvents + storeMetrics.PendingEvents,
		DeliveredEvents: busMetrics.DeliveredEvents + storeMetrics.DeliveredEvents,
		FailedEvents:    busMetrics.FailedEvents + storeMetrics.FailedEvents,
		ExpiredEvents:   storeMetrics.ExpiredEvents,
	}

	return combinedMetrics, nil
}
