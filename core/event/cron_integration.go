package event

import (
	"context"
	"fmt"
	"time"

	"github.com/zhoudm1743/go-flow/core/cron"
	"github.com/zhoudm1743/go-flow/core/logger"
)

// CronEventIntegration Cron事件集成器
type CronEventIntegration struct {
	eventService EventService
	cronService  cron.Service
	logger       logger.Logger
}

// NewCronEventIntegration 创建Cron事件集成器
func NewCronEventIntegration(
	eventService EventService,
	cronService cron.Service,
	logger logger.Logger,
) *CronEventIntegration {
	return &CronEventIntegration{
		eventService: eventService,
		cronService:  cronService,
		logger:       logger,
	}
}

// Initialize 初始化集成
func (c *CronEventIntegration) Initialize(ctx context.Context) error {
	// 注册Cron任务执行事件处理器
	handlers := []EventHandler{
		NewCronTaskCreatedHandler(c.cronService, c.logger),
		NewCronTaskExecutedHandler(c.cronService, c.logger),
		NewCronTaskFailedHandler(c.cronService, c.logger),
		NewCronTaskCompletedHandler(c.cronService, c.logger),
	}

	for _, handler := range handlers {
		if err := c.eventService.Subscribe(handler.GetEventType(), handler); err != nil {
			return fmt.Errorf("订阅处理器 %s 失败: %w", handler.GetName(), err)
		}
	}

	c.logger.Info("Cron事件集成初始化成功")
	return nil
}

// PublishCronTaskEvent 发布Cron任务事件
func (c *CronEventIntegration) PublishCronTaskEvent(ctx context.Context, eventType string, taskInfo *cron.TaskInfo, result *cron.ExecResult) error {
	payload := &CronTaskPayload{
		TaskID:   taskInfo.ID,
		TaskName: taskInfo.Name,
		TaskType: string(taskInfo.Type),
		Status:   string(taskInfo.Status),
		Metadata: map[string]interface{}{
			"cron":        taskInfo.Cron,
			"description": taskInfo.Description,
			"created_at":  taskInfo.CreatedAt,
			"updated_at":  taskInfo.UpdatedAt,
		},
	}

	if result != nil {
		payload.ExecutedAt = result.StartTime
		payload.Duration = result.Duration
		payload.Result = result.Output
		payload.Error = result.Error
	}

	return c.eventService.PublishEventAsync(ctx, eventType, payload, "cron-service")
}

// CreateDelayedCronTask 创建延时Cron任务
func (c *CronEventIntegration) CreateDelayedCronTask(ctx context.Context, delayDuration time.Duration, taskRequest *cron.CreateTaskRequest) (*DelayedEvent, error) {
	payload := map[string]interface{}{
		"task_request": taskRequest,
		"action":       "create_task",
	}

	options := &DelayedEventOptions{
		Priority:   EventPriorityNormal,
		MaxRetries: 3,
	}

	return c.eventService.ScheduleDelayedEvent(ctx, "cron.task.delayed_create", payload, "event-integration", time.Now().Add(delayDuration), options)
}

// ScheduleTaskExecution 调度任务执行事件
func (c *CronEventIntegration) ScheduleTaskExecution(ctx context.Context, taskID string, executeAt time.Time) (*DelayedEvent, error) {
	payload := map[string]interface{}{
		"task_id": taskID,
		"action":  "execute",
	}

	options := &DelayedEventOptions{
		Priority:   EventPriorityHigh,
		MaxRetries: 2,
	}

	return c.eventService.ScheduleDelayedEvent(ctx, "cron.task.delayed_execute", payload, "event-integration", executeAt, options)
}

// CronTaskCreatedHandler Cron任务创建事件处理器
type CronTaskCreatedHandler struct {
	cronService cron.Service
	logger      logger.Logger
}

// NewCronTaskCreatedHandler 创建Cron任务创建事件处理器
func NewCronTaskCreatedHandler(cronService cron.Service, logger logger.Logger) *CronTaskCreatedHandler {
	return &CronTaskCreatedHandler{
		cronService: cronService,
		logger:      logger,
	}
}

// Handle 处理事件
func (h *CronTaskCreatedHandler) Handle(ctx context.Context, event Event) error {
	h.logger.WithField("event_type", event.GetType()).
		WithField("event_id", event.GetID()).
		WithField("source", event.GetSource()).
		Info("正在处理Cron任务创建事件")

	payload, ok := event.GetPayload().(*CronTaskPayload)
	if !ok {
		return fmt.Errorf("Cron任务创建事件的负载类型无效")
	}

	h.logger.WithField("task_id", payload.TaskID).
		WithField("task_name", payload.TaskName).
		Info("处理Cron任务创建事件")

	// 可以在这里执行一些后处理逻辑
	// 例如：发送通知、更新统计信息等

	return nil
}

// GetEventType 获取处理的事件类型
func (h *CronTaskCreatedHandler) GetEventType() string {
	return EventTypeCronTaskCreated
}

// GetName 获取处理器名称
func (h *CronTaskCreatedHandler) GetName() string {
	return "CronTaskCreatedHandler"
}

// CronTaskExecutedHandler Cron任务执行事件处理器
type CronTaskExecutedHandler struct {
	cronService cron.Service
	logger      logger.Logger
}

// NewCronTaskExecutedHandler 创建Cron任务执行事件处理器
func NewCronTaskExecutedHandler(cronService cron.Service, logger logger.Logger) *CronTaskExecutedHandler {
	return &CronTaskExecutedHandler{
		cronService: cronService,
		logger:      logger,
	}
}

// Handle 处理事件
func (h *CronTaskExecutedHandler) Handle(ctx context.Context, event Event) error {
	h.logger.WithField("event_type", event.GetType()).
		WithField("event_id", event.GetID()).
		WithField("source", event.GetSource()).
		Info("正在处理Cron任务执行事件")

	payload, ok := event.GetPayload().(*CronTaskPayload)
	if !ok {
		return fmt.Errorf("Cron任务执行事件的负载类型无效")
	}

	h.logger.WithField("task_id", payload.TaskID).
		WithField("task_name", payload.TaskName).
		WithField("duration", payload.Duration).
		Info("处理Cron任务执行事件")

	// 可以在这里执行一些监控和统计逻辑
	// 例如：记录执行时间、更新性能指标等

	return nil
}

// GetEventType 获取处理的事件类型
func (h *CronTaskExecutedHandler) GetEventType() string {
	return EventTypeCronTaskExecuted
}

// GetName 获取处理器名称
func (h *CronTaskExecutedHandler) GetName() string {
	return "CronTaskExecutedHandler"
}

// CronTaskFailedHandler Cron任务失败事件处理器
type CronTaskFailedHandler struct {
	cronService cron.Service
	logger      logger.Logger
}

// NewCronTaskFailedHandler 创建Cron任务失败事件处理器
func NewCronTaskFailedHandler(cronService cron.Service, logger logger.Logger) *CronTaskFailedHandler {
	return &CronTaskFailedHandler{
		cronService: cronService,
		logger:      logger,
	}
}

// Handle 处理事件
func (h *CronTaskFailedHandler) Handle(ctx context.Context, event Event) error {
	h.logger.WithField("event_type", event.GetType()).
		WithField("event_id", event.GetID()).
		WithField("source", event.GetSource()).
		Warn("正在处理Cron任务失败事件")

	payload, ok := event.GetPayload().(*CronTaskPayload)
	if !ok {
		return fmt.Errorf("Cron任务失败事件的负载类型无效")
	}

	h.logger.WithField("task_id", payload.TaskID).
		WithField("task_name", payload.TaskName).
		WithField("error", payload.Error).
		Warn("处理Cron任务失败事件")

	// 可以在这里执行失败处理逻辑
	// 例如：发送告警、记录失败统计、触发重试等

	// 如果是重要任务失败，可以调度一个延时重试事件
	if payload.Metadata != nil {
		if priority, exists := payload.Metadata["priority"]; exists && priority == "high" {
			// 调度5分钟后重试
			// 这里可以创建一个延时事件来重新执行任务
		}
	}

	return nil
}

// GetEventType 获取处理的事件类型
func (h *CronTaskFailedHandler) GetEventType() string {
	return EventTypeCronTaskFailed
}

// GetName 获取处理器名称
func (h *CronTaskFailedHandler) GetName() string {
	return "CronTaskFailedHandler"
}

// CronTaskCompletedHandler Cron任务完成事件处理器
type CronTaskCompletedHandler struct {
	cronService cron.Service
	logger      logger.Logger
}

// NewCronTaskCompletedHandler 创建Cron任务完成事件处理器
func NewCronTaskCompletedHandler(cronService cron.Service, logger logger.Logger) *CronTaskCompletedHandler {
	return &CronTaskCompletedHandler{
		cronService: cronService,
		logger:      logger,
	}
}

// Handle 处理事件
func (h *CronTaskCompletedHandler) Handle(ctx context.Context, event Event) error {
	h.logger.WithField("event_type", event.GetType()).
		WithField("event_id", event.GetID()).
		WithField("source", event.GetSource()).
		Debug("正在处理Cron任务完成事件")

	payload, ok := event.GetPayload().(*CronTaskPayload)
	if !ok {
		return fmt.Errorf("Cron任务完成事件的负载类型无效")
	}

	h.logger.WithField("task_id", payload.TaskID).
		WithField("task_name", payload.TaskName).
		WithField("status", payload.Status).
		Debug("处理Cron任务完成事件")

	// 可以在这里执行完成后的清理和统计逻辑
	// 例如：更新任务状态、清理临时数据等

	return nil
}

// GetEventType 获取处理的事件类型
func (h *CronTaskCompletedHandler) GetEventType() string {
	return EventTypeCronTaskCompleted
}

// GetName 获取处理器名称
func (h *CronTaskCompletedHandler) GetName() string {
	return "CronTaskCompletedHandler"
}

// DelayedCronTaskHandler 延时Cron任务处理器
type DelayedCronTaskHandler struct {
	cronService cron.Service
	logger      logger.Logger
}

// NewDelayedCronTaskHandler 创建延时Cron任务处理器
func NewDelayedCronTaskHandler(cronService cron.Service, logger logger.Logger) *DelayedCronTaskHandler {
	return &DelayedCronTaskHandler{
		cronService: cronService,
		logger:      logger,
	}
}

// Handle 处理事件
func (h *DelayedCronTaskHandler) Handle(ctx context.Context, event Event) error {
	payload := event.GetPayload().(map[string]interface{})

	action, ok := payload["action"].(string)
	if !ok {
		return fmt.Errorf("延时Cron任务事件中的操作无效")
	}

	switch action {
	case "create_task":
		return h.handleCreateTask(ctx, payload)
	case "execute":
		return h.handleExecuteTask(ctx, payload)
	default:
		return fmt.Errorf("未知操作: %s", action)
	}
}

// handleCreateTask 处理创建任务
func (h *DelayedCronTaskHandler) handleCreateTask(ctx context.Context, payload map[string]interface{}) error {
	taskRequestData, ok := payload["task_request"]
	if !ok {
		return fmt.Errorf("负载中未找到task_request")
	}

	// 这里需要反序列化taskRequest
	// 由于类型转换的复杂性，实际使用时可能需要更复杂的处理
	h.logger.WithField("task_request", taskRequestData).Info("正在处理延时Cron任务创建")

	// 实际创建任务的逻辑
	// taskInfo, err := h.cronService.CreateTask(ctx, taskRequest)
	// if err != nil {
	//     return fmt.Errorf("failed to create delayed cron task: %w", err)
	// }

	return nil
}

// handleExecuteTask 处理执行任务
func (h *DelayedCronTaskHandler) handleExecuteTask(ctx context.Context, payload map[string]interface{}) error {
	taskID, ok := payload["task_id"].(string)
	if !ok {
		return fmt.Errorf("负载中未找到task_id")
	}

	h.logger.WithField("task_id", taskID).Info("正在处理延时Cron任务执行")

	// 执行任务
	err := h.cronService.ExecuteTaskNow(ctx, taskID)
	if err != nil {
		return fmt.Errorf("执行延时Cron任务失败: %w", err)
	}

	h.logger.WithField("task_id", taskID).Info("延时Cron任务执行成功")

	return nil
}

// GetEventType 获取处理的事件类型
func (h *DelayedCronTaskHandler) GetEventType() string {
	return "cron.task.delayed_create"
}

// GetName 获取处理器名称
func (h *DelayedCronTaskHandler) GetName() string {
	return "DelayedCronTaskHandler"
}

// CronTaskStatusWatcher Cron任务状态监控器
type CronTaskStatusWatcher struct {
	cronService  cron.Service
	eventService EventService
	logger       logger.Logger
}

// NewCronTaskStatusWatcher 创建Cron任务状态监控器
func NewCronTaskStatusWatcher(
	cronService cron.Service,
	eventService EventService,
	logger logger.Logger,
) *CronTaskStatusWatcher {
	return &CronTaskStatusWatcher{
		cronService:  cronService,
		eventService: eventService,
		logger:       logger,
	}
}

// WatchTaskStatus 监控任务状态变化
func (w *CronTaskStatusWatcher) WatchTaskStatus(ctx context.Context, taskID string) error {
	// 定期检查任务状态并发布相应事件
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	var lastStatus cron.TaskStatus

	for {
		select {
		case <-ticker.C:
			taskInfo, err := w.cronService.GetTask(ctx, taskID)
			if err != nil {
				w.logger.WithField("task_id", taskID).
					WithError(err).Error("获取任务信息失败，无法进行状态监控")
				continue
			}

			// 检查状态是否变化
			if taskInfo.Status != lastStatus {
				eventType := ""
				switch taskInfo.Status {
				case cron.TaskStatusActive:
					eventType = EventTypeCronTaskUpdated
				case cron.TaskStatusPaused:
					eventType = EventTypeCronTaskUpdated
				case cron.TaskStatusStopped:
					eventType = EventTypeCronTaskDeleted
				}

				if eventType != "" {
					payload := &CronTaskPayload{
						TaskID:   taskInfo.ID,
						TaskName: taskInfo.Name,
						TaskType: string(taskInfo.Type),
						Status:   string(taskInfo.Status),
						Metadata: map[string]interface{}{
							"previous_status": string(lastStatus),
							"status_changed":  true,
						},
					}

					if err := w.eventService.PublishEventAsync(ctx, eventType, payload, "status-watcher"); err != nil {
						w.logger.WithField("task_id", taskID).
							WithField("event_type", eventType).
							WithError(err).Error("发布任务状态变化事件失败")
					}
				}

				lastStatus = taskInfo.Status
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
