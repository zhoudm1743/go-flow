package event

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zhoudm1743/go-flow/core/cron"
	"github.com/zhoudm1743/go-flow/core/logger"
)

// DelayedEventScheduler 延时事件调度器接口
type DelayedEventScheduler interface {
	// ScheduleDelayedEvent 调度延时事件
	ScheduleDelayedEvent(ctx context.Context, event *DelayedEvent) error
	// CancelDelayedEvent 取消延时事件
	CancelDelayedEvent(ctx context.Context, eventID string) error
	// Start 启动调度器
	Start(ctx context.Context) error
	// Stop 停止调度器
	Stop(ctx context.Context) error
	// GetPendingCount 获取待处理事件数量
	GetPendingCount() int64
}

// CronDelayedEventScheduler 基于Cron的延时事件调度器
type CronDelayedEventScheduler struct {
	eventBus       EventBus           // 事件总线
	eventStore     EventStore         // 事件存储
	cronService    cron.Service       // Cron服务
	logger         logger.Logger      // 日志记录器
	options        *SchedulerOptions  // 调度器选项
	running        bool               // 运行状态
	scheduledTasks map[string]string  // 已调度任务映射 eventID -> cronTaskID
	mu             sync.RWMutex       // 读写锁
	ctx            context.Context    // 上下文
	cancel         context.CancelFunc // 取消函数
	wg             sync.WaitGroup     // 等待组
}

// SchedulerOptions 调度器选项
type SchedulerOptions struct {
	PollInterval    time.Duration `json:"poll_interval"`    // 轮询间隔
	BatchSize       int           `json:"batch_size"`       // 批处理大小
	MaxRetries      int           `json:"max_retries"`      // 最大重试次数
	RetryInterval   time.Duration `json:"retry_interval"`   // 重试间隔
	CleanupInterval time.Duration `json:"cleanup_interval"` // 清理间隔
	ExpireAfter     time.Duration `json:"expire_after"`     // 事件过期时间
}

// DefaultSchedulerOptions 默认调度器选项
func DefaultSchedulerOptions() *SchedulerOptions {
	return &SchedulerOptions{
		PollInterval:    30 * time.Second,
		BatchSize:       50,
		MaxRetries:      3,
		RetryInterval:   5 * time.Minute,
		CleanupInterval: 1 * time.Hour,
		ExpireAfter:     24 * time.Hour,
	}
}

// NewCronDelayedEventScheduler 创建基于Cron的延时事件调度器
func NewCronDelayedEventScheduler(
	eventBus EventBus,
	eventStore EventStore,
	cronService cron.Service,
	logger logger.Logger,
	options *SchedulerOptions,
) *CronDelayedEventScheduler {
	if options == nil {
		options = DefaultSchedulerOptions()
	}

	return &CronDelayedEventScheduler{
		eventBus:       eventBus,
		eventStore:     eventStore,
		cronService:    cronService,
		logger:         logger,
		options:        options,
		scheduledTasks: make(map[string]string),
	}
}

// Start 启动调度器
func (s *CronDelayedEventScheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("延时事件调度器已经在运行")
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	s.running = true

	// 启动事件轮询器
	s.wg.Add(1)
	go s.eventPoller()

	// 启动清理器
	s.wg.Add(1)
	go s.cleaner()

	s.logger.WithField("poll_interval", s.options.PollInterval).
		WithField("batch_size", s.options.BatchSize).
		Info("延时事件调度器启动成功")

	return nil
}

// Stop 停止调度器
func (s *CronDelayedEventScheduler) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.logger.Info("正在停止延时事件调度器...")

	// 取消上下文
	if s.cancel != nil {
		s.cancel()
	}

	// 等待所有goroutine退出
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("延时事件调度器优雅停止完成")
	case <-ctx.Done():
		s.logger.Warn("延时事件调度器停止超时")
		return ctx.Err()
	}

	s.running = false
	return nil
}

// ScheduleDelayedEvent 调度延时事件
func (s *CronDelayedEventScheduler) ScheduleDelayedEvent(ctx context.Context, event *DelayedEvent) error {
	// 存储事件
	if err := s.eventStore.StoreDelayedEvent(ctx, event); err != nil {
		return fmt.Errorf("存储延时事件失败: %w", err)
	}

	// 如果延时时间很短，直接由轮询器处理
	if event.DelayUntil.Sub(time.Now()) <= s.options.PollInterval {
		s.logger.WithField("event_id", event.ID).
			WithField("delay_until", event.DelayUntil).
			Debug("事件延时时间较短，将由轮询器处理")
		return nil
	}

	// 创建Cron任务调度事件
	cronExpression := s.generateCronExpression(event.DelayUntil)
	cronTask := &cron.CreateTaskRequest{
		Name:        fmt.Sprintf("delayed_event_%s", event.ID),
		Type:        cron.TaskTypeSystem,
		Cron:        cronExpression,
		Description: fmt.Sprintf("Delayed event trigger: %s", event.Type),
		Config: &cron.SystemConfig{
			HandlerName: "DelayedEventTrigger",
			Parameters: map[string]interface{}{
				"event_id": event.ID,
			},
		},
	}

	cronTaskResp, err := s.cronService.CreateTask(ctx, cronTask)
	if err != nil {
		return fmt.Errorf("为延时事件创建定时任务失败: %w", err)
	}

	// 记录任务映射
	s.mu.Lock()
	s.scheduledTasks[event.ID] = cronTaskResp.ID
	s.mu.Unlock()

	s.logger.WithField("event_id", event.ID).
		WithField("cron_task_id", cronTaskResp.ID).
		WithField("delay_until", event.DelayUntil).
		Debug("延时事件已通过Cron任务调度")

	return nil
}

// CancelDelayedEvent 取消延时事件
func (s *CronDelayedEventScheduler) CancelDelayedEvent(ctx context.Context, eventID string) error {
	// 从存储中删除事件
	if err := s.eventStore.DeleteDelayedEvent(ctx, eventID); err != nil {
		s.logger.WithField("event_id", eventID).
			WithError(err).Warn("从存储中删除延时事件失败")
	}

	// 取消对应的Cron任务
	s.mu.Lock()
	cronTaskID, exists := s.scheduledTasks[eventID]
	if exists {
		delete(s.scheduledTasks, eventID)
	}
	s.mu.Unlock()

	if exists {
		if err := s.cronService.DeleteTask(ctx, cronTaskID); err != nil {
			s.logger.WithField("event_id", eventID).
				WithField("cron_task_id", cronTaskID).
				WithError(err).Warn("删除延时事件的Cron任务失败")
		} else {
			s.logger.WithField("event_id", eventID).
				WithField("cron_task_id", cronTaskID).
				Debug("已取消延时事件的Cron任务")
		}
	}

	return nil
}

// GetPendingCount 获取待处理事件数量
func (s *CronDelayedEventScheduler) GetPendingCount() int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metrics, err := s.eventStore.GetMetrics(ctx)
	if err != nil {
		s.logger.WithError(err).Warn("获取事件指标失败")
		return 0
	}

	return metrics.PendingEvents
}

// eventPoller 事件轮询器
func (s *CronDelayedEventScheduler) eventPoller() {
	defer s.wg.Done()

	s.logger.Debug("事件轮询器启动")
	ticker := time.NewTicker(s.options.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.processReadyEvents()
		case <-s.ctx.Done():
			s.logger.Debug("事件轮询器停止")
			return
		}
	}
}

// processReadyEvents 处理准备就绪的事件
func (s *CronDelayedEventScheduler) processReadyEvents() {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	events, err := s.eventStore.GetReadyEvents(ctx, s.options.BatchSize)
	if err != nil {
		s.logger.WithError(err).Error("获取准备就绪的事件失败")
		return
	}

	if len(events) == 0 {
		return
	}

	s.logger.WithField("count", len(events)).Debug("正在处理准备就绪的事件")

	for _, event := range events {
		if err := s.processDelayedEvent(ctx, event); err != nil {
			s.logger.WithField("event_id", event.ID).
				WithField("event_type", event.Type).
				WithError(err).Error("处理延时事件失败")
		}
	}
}

// processDelayedEvent 处理延时事件
func (s *CronDelayedEventScheduler) processDelayedEvent(ctx context.Context, delayedEvent *DelayedEvent) error {
	// 创建基础事件
	baseEvent := &BaseEvent{
		ID:        delayedEvent.ID,
		Type:      delayedEvent.Type,
		Source:    delayedEvent.Source,
		Payload:   delayedEvent.Payload,
		Timestamp: time.Now(),
	}

	// 发布事件
	if err := s.eventBus.PublishAsync(ctx, baseEvent); err != nil {
		// 处理重试逻辑
		if delayedEvent.RetryCount < delayedEvent.MaxRetries {
			// 更新重试次数和下次重试时间
			delayedEvent.RetryCount++
			delayedEvent.DelayUntil = time.Now().Add(s.options.RetryInterval)
			delayedEvent.Status = EventStatusPending

			if err := s.eventStore.UpdateDelayedEvent(ctx, delayedEvent); err != nil {
				s.logger.WithField("event_id", delayedEvent.ID).
					WithError(err).Error("更新延时事件重试状态失败")
				return err
			}

			s.logger.WithField("event_id", delayedEvent.ID).
				WithField("retry_count", delayedEvent.RetryCount).
				WithField("max_retries", delayedEvent.MaxRetries).
				WithField("next_retry", delayedEvent.DelayUntil).
				Warn("延时事件处理失败，已安排重试")

			return nil
		} else {
			// 超过最大重试次数，标记为失败
			delayedEvent.Status = EventStatusFailed
			if err := s.eventStore.UpdateDelayedEvent(ctx, delayedEvent); err != nil {
				s.logger.WithField("event_id", delayedEvent.ID).
					WithError(err).Error("更新延时事件失败状态失败")
			}

			return fmt.Errorf("延时事件在 %d 次重试后失败: %w", delayedEvent.MaxRetries, err)
		}
	}

	// 发布成功，标记为已发送
	delayedEvent.Status = EventStatusDelivered
	if err := s.eventStore.UpdateDelayedEvent(ctx, delayedEvent); err != nil {
		s.logger.WithField("event_id", delayedEvent.ID).
			WithError(err).Warn("更新延时事件状态为已发送失败")
	}

	// 从调度任务映射中移除
	s.mu.Lock()
	delete(s.scheduledTasks, delayedEvent.ID)
	s.mu.Unlock()

	return nil
}

// cleaner 清理器
func (s *CronDelayedEventScheduler) cleaner() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.options.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanup()
		case <-s.ctx.Done():
			return
		}
	}
}

// cleanup 清理过期事件
func (s *CronDelayedEventScheduler) cleanup() {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	expireBefore := time.Now().Add(-s.options.ExpireAfter)
	count, err := s.eventStore.Cleanup(ctx, expireBefore)
	if err != nil {
		s.logger.WithError(err).Warn("清理过期事件失败")
		return
	}

	if count > 0 {
		s.logger.WithField("count", count).Info("清理过期事件完成")
	}
}

// generateCronExpression 为指定时间生成Cron表达式
func (s *CronDelayedEventScheduler) generateCronExpression(targetTime time.Time) string {
	// 生成一次性执行的Cron表达式（秒级精度）
	return fmt.Sprintf("%d %d %d %d %d ? %d",
		targetTime.Second(),
		targetTime.Minute(),
		targetTime.Hour(),
		targetTime.Day(),
		int(targetTime.Month()),
		targetTime.Year())
}

// DelayedEventTriggerHandler 延时事件触发处理器
type DelayedEventTriggerHandler struct {
	scheduler DelayedEventScheduler
	logger    logger.Logger
}

// NewDelayedEventTriggerHandler 创建延时事件触发处理器
func NewDelayedEventTriggerHandler(scheduler DelayedEventScheduler, logger logger.Logger) *DelayedEventTriggerHandler {
	return &DelayedEventTriggerHandler{
		scheduler: scheduler,
		logger:    logger,
	}
}

// Handle 处理Cron任务触发
func (h *DelayedEventTriggerHandler) Handle(ctx context.Context, params map[string]interface{}) (string, error) {
	eventID, ok := params["event_id"].(string)
	if !ok {
		return "", fmt.Errorf("event_id参数是必需的且必须是字符串")
	}

	h.logger.WithField("event_id", eventID).Debug("Cron任务触发延时事件")

	// 延时事件的实际处理由轮询器完成，这里只是确保触发
	// 可以向事件总线发送一个内部事件来立即处理特定的延时事件

	return fmt.Sprintf("延时事件 %s 已触发", eventID), nil
}

// GetName 获取处理器名称
func (h *DelayedEventTriggerHandler) GetName() string {
	return "DelayedEventTrigger"
}

// GetDescription 获取处理器描述
func (h *DelayedEventTriggerHandler) GetDescription() string {
	return "Trigger delayed event execution"
}
