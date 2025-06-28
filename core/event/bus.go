package event

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zhoudm1743/go-flow/core/logger"
)

// EventBus 事件总线接口
type EventBus interface {
	// Publish 发布事件
	Publish(ctx context.Context, event Event) error
	// PublishAsync 异步发布事件
	PublishAsync(ctx context.Context, event Event) error
	// Subscribe 订阅事件
	Subscribe(eventType string, handler EventHandler, options ...*SubscriptionOptions) error
	// Unsubscribe 取消订阅
	Unsubscribe(eventType string, handlerName string) error
	// Start 启动事件总线
	Start(ctx context.Context) error
	// Stop 停止事件总线
	Stop(ctx context.Context) error
	// GetMetrics 获取事件指标
	GetMetrics() *EventMetrics
}

// Subscription 订阅信息
type Subscription struct {
	Handler EventHandler
	Options *SubscriptionOptions
}

// DefaultEventBus 默认事件总线实现
type DefaultEventBus struct {
	subscriptions map[string][]*Subscription // 订阅信息 eventType -> subscriptions
	eventChan     chan Event                 // 事件通道
	workers       []*Worker                  // 工作协程
	options       *EventBusOptions           // 配置选项
	logger        logger.Logger              // 日志记录器
	metrics       *EventMetrics              // 事件指标
	running       bool                       // 运行状态
	mu            sync.RWMutex               // 读写锁
	ctx           context.Context            // 上下文
	cancel        context.CancelFunc         // 取消函数
	wg            sync.WaitGroup             // 等待组
}

// Worker 工作协程
type Worker struct {
	id      int
	eventCh chan Event
	bus     *DefaultEventBus
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewDefaultEventBus 创建默认事件总线
func NewDefaultEventBus(logger logger.Logger, options *EventBusOptions) *DefaultEventBus {
	if options == nil {
		options = DefaultEventBusOptions()
	}

	return &DefaultEventBus{
		subscriptions: make(map[string][]*Subscription),
		eventChan:     make(chan Event, options.BufferSize),
		options:       options,
		logger:        logger,
		metrics:       &EventMetrics{},
	}
}

// Start 启动事件总线
func (bus *DefaultEventBus) Start(ctx context.Context) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if bus.running {
		return fmt.Errorf("事件总线已经在运行")
	}

	bus.ctx, bus.cancel = context.WithCancel(ctx)
	bus.running = true

	// 启动工作协程
	bus.workers = make([]*Worker, bus.options.WorkerCount)
	for i := 0; i < bus.options.WorkerCount; i++ {
		worker := &Worker{
			id:      i,
			eventCh: make(chan Event, bus.options.BufferSize/bus.options.WorkerCount),
			bus:     bus,
		}
		worker.ctx, worker.cancel = context.WithCancel(bus.ctx)
		bus.workers[i] = worker

		bus.wg.Add(1)
		go worker.run()
	}

	// 启动事件分发器
	bus.wg.Add(1)
	go bus.dispatcher()

	bus.logger.WithField("worker_count", bus.options.WorkerCount).
		WithField("buffer_size", bus.options.BufferSize).
		Info("事件总线启动成功")

	return nil
}

// Stop 停止事件总线
func (bus *DefaultEventBus) Stop(ctx context.Context) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if !bus.running {
		return nil
	}

	bus.logger.Info("正在停止事件总线...")

	// 取消上下文
	if bus.cancel != nil {
		bus.cancel()
	}

	// 关闭事件通道
	close(bus.eventChan)

	// 等待所有工作协程退出
	done := make(chan struct{})
	go func() {
		bus.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		bus.logger.Info("事件总线优雅停止完成")
	case <-ctx.Done():
		bus.logger.Warn("事件总线停止超时")
		return ctx.Err()
	}

	bus.running = false
	return nil
}

// Publish 发布事件（同步）
func (bus *DefaultEventBus) Publish(ctx context.Context, event Event) error {
	bus.mu.RLock()
	if !bus.running {
		bus.mu.RUnlock()
		return fmt.Errorf("事件总线未运行")
	}
	bus.mu.RUnlock()

	// 更新指标
	bus.metrics.TotalEvents++

	// 获取订阅者
	subscriptions := bus.getSubscriptions(event.GetType())
	if len(subscriptions) == 0 {
		bus.logger.WithField("event_type", event.GetType()).
			WithField("event_id", event.GetID()).
			Debug("该事件类型没有订阅者")
		return nil
	}

	// 同步处理事件
	var lastErr error
	for _, sub := range subscriptions {
		if sub.Options.Async {
			// 异步处理
			select {
			case bus.eventChan <- event:
				bus.metrics.DeliveredEvents++
			case <-ctx.Done():
				return ctx.Err()
			default:
				bus.metrics.FailedEvents++
				lastErr = fmt.Errorf("事件通道已满")
				bus.logger.WithField("event_type", event.GetType()).
					WithField("event_id", event.GetID()).
					Warn("事件通道已满，发布事件失败")
			}
		} else {
			// 同步处理
			if err := bus.handleEventSync(ctx, event, sub); err != nil {
				bus.metrics.FailedEvents++
				lastErr = err
				bus.logger.WithField("event_type", event.GetType()).
					WithField("event_id", event.GetID()).
					WithField("handler", sub.Handler.GetName()).
					WithError(err).Error("同步处理事件失败")
			} else {
				bus.metrics.DeliveredEvents++
			}
		}
	}

	return lastErr
}

// PublishAsync 异步发布事件
func (bus *DefaultEventBus) PublishAsync(ctx context.Context, event Event) error {
	bus.mu.RLock()
	if !bus.running {
		bus.mu.RUnlock()
		return fmt.Errorf("事件总线未运行")
	}
	bus.mu.RUnlock()

	// 更新指标
	bus.metrics.TotalEvents++

	// 异步发送到事件通道
	select {
	case bus.eventChan <- event:
		bus.metrics.PendingEvents++
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		bus.metrics.FailedEvents++
		return fmt.Errorf("事件通道已满")
	}
}

// Subscribe 订阅事件
func (bus *DefaultEventBus) Subscribe(eventType string, handler EventHandler, options ...*SubscriptionOptions) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	// 使用默认选项
	opts := DefaultSubscriptionOptions()
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	}

	// 检查是否已经订阅
	subscriptions := bus.subscriptions[eventType]
	for _, sub := range subscriptions {
		if sub.Handler.GetName() == handler.GetName() {
			return fmt.Errorf("处理器 %s 已订阅事件类型 %s",
				handler.GetName(), eventType)
		}
	}

	// 添加订阅
	subscription := &Subscription{
		Handler: handler,
		Options: opts,
	}
	bus.subscriptions[eventType] = append(bus.subscriptions[eventType], subscription)

	bus.logger.WithField("event_type", eventType).
		WithField("handler", handler.GetName()).
		WithField("async", opts.Async).
		Info("事件处理器订阅成功")

	return nil
}

// Unsubscribe 取消订阅
func (bus *DefaultEventBus) Unsubscribe(eventType string, handlerName string) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	subscriptions := bus.subscriptions[eventType]
	for i, sub := range subscriptions {
		if sub.Handler.GetName() == handlerName {
			// 删除订阅
			bus.subscriptions[eventType] = append(subscriptions[:i], subscriptions[i+1:]...)

			bus.logger.WithField("event_type", eventType).
				WithField("handler", handlerName).
				Info("事件处理器取消订阅成功")

			return nil
		}
	}

	return fmt.Errorf("事件类型 %s 的处理器 %s 未找到", eventType, handlerName)
}

// GetMetrics 获取事件指标
func (bus *DefaultEventBus) GetMetrics() *EventMetrics {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	// 返回副本
	return &EventMetrics{
		TotalEvents:     bus.metrics.TotalEvents,
		PendingEvents:   bus.metrics.PendingEvents,
		DeliveredEvents: bus.metrics.DeliveredEvents,
		FailedEvents:    bus.metrics.FailedEvents,
		ExpiredEvents:   bus.metrics.ExpiredEvents,
	}
}

// getSubscriptions 获取事件类型的订阅者
func (bus *DefaultEventBus) getSubscriptions(eventType string) []*Subscription {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	subscriptions := bus.subscriptions[eventType]
	if len(subscriptions) == 0 {
		return nil
	}

	// 返回副本
	result := make([]*Subscription, len(subscriptions))
	copy(result, subscriptions)
	return result
}

// dispatcher 事件分发器
func (bus *DefaultEventBus) dispatcher() {
	defer bus.wg.Done()

	bus.logger.Debug("事件分发器启动")

	for {
		select {
		case event, ok := <-bus.eventChan:
			if !ok {
				bus.logger.Debug("事件通道已关闭，分发器退出")
				return
			}

			// 计算哈希值选择工作协程
			hash := hash(event.GetType() + event.GetID())
			workerIndex := hash % len(bus.workers)
			worker := bus.workers[workerIndex]

			// 发送到对应的工作协程
			select {
			case worker.eventCh <- event:
				bus.metrics.PendingEvents--
			default:
				// 工作协程通道满了，丢弃事件
				bus.metrics.FailedEvents++
				bus.logger.WithField("event_type", event.GetType()).
					WithField("event_id", event.GetID()).
					WithField("worker_id", workerIndex).
					Warn("工作协程通道已满，丢弃事件")
			}

		case <-bus.ctx.Done():
			bus.logger.Debug("上下文已取消，分发器退出")
			return
		}
	}
}

// handleEventSync 同步处理事件
func (bus *DefaultEventBus) handleEventSync(ctx context.Context, event Event, sub *Subscription) error {
	// 创建超时上下文
	timeout := sub.Options.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second // 默认超时时间
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 使用 goroutine 处理事件以支持超时
	errCh := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errCh <- fmt.Errorf("事件处理器发生恐慌: %v", r)
			}
		}()
		errCh <- sub.Handler.Handle(timeoutCtx, event)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			bus.logger.WithField("event_type", event.GetType()).
				WithField("event_id", event.GetID()).
				WithField("handler", sub.Handler.GetName()).
				WithError(err).Warn("事件处理器执行失败")
		}
		return err
	case <-timeoutCtx.Done():
		return fmt.Errorf("事件处理器超时")
	}
}

// run 工作协程运行
func (w *Worker) run() {
	defer w.bus.wg.Done()

	w.bus.logger.WithField("worker_id", w.id).Debug("工作协程启动")

	for {
		select {
		case event, ok := <-w.eventCh:
			if !ok {
				w.bus.logger.WithField("worker_id", w.id).Debug("工作协程通道已关闭")
				return
			}

			w.handleEvent(event)

		case <-w.ctx.Done():
			w.bus.logger.WithField("worker_id", w.id).Debug("工作协程上下文已取消")
			return
		}
	}
}

// handleEvent 处理事件
func (w *Worker) handleEvent(event Event) {
	subscriptions := w.bus.getSubscriptions(event.GetType())
	if len(subscriptions) == 0 {
		return
	}

	for _, sub := range subscriptions {
		if !sub.Options.Async {
			continue // 跳过同步处理的订阅
		}

		if err := w.bus.handleEventSync(w.ctx, event, sub); err != nil {
			w.bus.metrics.FailedEvents++
		} else {
			w.bus.metrics.DeliveredEvents++
		}
	}
}

// hash 简单哈希函数
func hash(s string) int {
	h := 0
	for _, c := range s {
		h = 31*h + int(c)
	}
	if h < 0 {
		h = -h
	}
	return h
}
