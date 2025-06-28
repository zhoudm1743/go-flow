package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Event 事件接口
type Event interface {
	// GetID 获取事件ID
	GetID() string
	// GetType 获取事件类型
	GetType() string
	// GetPayload 获取事件负载数据
	GetPayload() interface{}
	// GetTimestamp 获取事件时间戳
	GetTimestamp() time.Time
	// GetSource 获取事件源
	GetSource() string
	// GetVersion 获取事件版本
	GetVersion() string
	// ToJSON 序列化为JSON
	ToJSON() (string, error)
}

// EventHandler 事件处理器接口
type EventHandler interface {
	// Handle 处理事件
	Handle(ctx context.Context, event Event) error
	// GetEventType 获取处理的事件类型
	GetEventType() string
	// GetName 获取处理器名称
	GetName() string
}

// EventStatus 事件状态
type EventStatus string

const (
	EventStatusPending   EventStatus = "pending"   // 待处理
	EventStatusDelivered EventStatus = "delivered" // 已投递
	EventStatusFailed    EventStatus = "failed"    // 处理失败
	EventStatusExpired   EventStatus = "expired"   // 已过期
)

// EventPriority 事件优先级
type EventPriority int

const (
	EventPriorityLow      EventPriority = 1
	EventPriorityNormal   EventPriority = 5
	EventPriorityHigh     EventPriority = 8
	EventPriorityCritical EventPriority = 10
)

// BaseEvent 基础事件实现
type BaseEvent struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
	Source    string      `json:"source"`
	Version   string      `json:"version"`
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(eventType string, payload interface{}, source string) *BaseEvent {
	return &BaseEvent{
		ID:        uuid.New().String(),
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now(),
		Source:    source,
		Version:   "1.0",
	}
}

// GetID 获取事件ID
func (e *BaseEvent) GetID() string {
	return e.ID
}

// GetType 获取事件类型
func (e *BaseEvent) GetType() string {
	return e.Type
}

// GetPayload 获取事件负载数据
func (e *BaseEvent) GetPayload() interface{} {
	return e.Payload
}

// GetTimestamp 获取事件时间戳
func (e *BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetSource 获取事件源
func (e *BaseEvent) GetSource() string {
	return e.Source
}

// GetVersion 获取事件版本
func (e *BaseEvent) GetVersion() string {
	return e.Version
}

// ToJSON 序列化为JSON
func (e *BaseEvent) ToJSON() (string, error) {
	data, err := json.Marshal(e)
	return string(data), err
}

// DelayedEvent 延时事件
type DelayedEvent struct {
	BaseEvent
	DelayUntil time.Time     `json:"delay_until"` // 延时到什么时候
	Priority   EventPriority `json:"priority"`    // 事件优先级
	MaxRetries int           `json:"max_retries"` // 最大重试次数
	RetryCount int           `json:"retry_count"` // 当前重试次数
	Status     EventStatus   `json:"status"`      // 事件状态
	CreatedAt  time.Time     `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time     `json:"updated_at"`  // 更新时间
}

// NewDelayedEvent 创建延时事件
func NewDelayedEvent(eventType string, payload interface{}, source string, delayUntil time.Time) *DelayedEvent {
	baseEvent := NewBaseEvent(eventType, payload, source)

	return &DelayedEvent{
		BaseEvent:  *baseEvent,
		DelayUntil: delayUntil,
		Priority:   EventPriorityNormal,
		MaxRetries: 3,
		RetryCount: 0,
		Status:     EventStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// IsReady 检查事件是否已准备好处理
func (d *DelayedEvent) IsReady() bool {
	return time.Now().After(d.DelayUntil) || time.Now().Equal(d.DelayUntil)
}

// CanRetry 检查是否可以重试
func (d *DelayedEvent) CanRetry() bool {
	return d.RetryCount < d.MaxRetries
}

// IncrRetry 增加重试次数
func (d *DelayedEvent) IncrRetry() {
	d.RetryCount++
	d.UpdatedAt = time.Now()
}

// SetStatus 设置事件状态
func (d *DelayedEvent) SetStatus(status EventStatus) {
	d.Status = status
	d.UpdatedAt = time.Now()
}

// EventRecord 事件记录（用于持久化）
type EventRecord struct {
	ID          string        `json:"id"`
	Type        string        `json:"type"`
	Payload     string        `json:"payload"` // JSON字符串
	Source      string        `json:"source"`
	DelayUntil  time.Time     `json:"delay_until"`
	Priority    EventPriority `json:"priority"`
	MaxRetries  int           `json:"max_retries"`
	RetryCount  int           `json:"retry_count"`
	Status      EventStatus   `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	ProcessedAt *time.Time    `json:"processed_at,omitempty"`
	Error       string        `json:"error,omitempty"`
}

// EventFilter 事件过滤器
type EventFilter struct {
	Types     []string      `json:"types,omitempty"`
	Sources   []string      `json:"sources,omitempty"`
	Status    EventStatus   `json:"status,omitempty"`
	Priority  EventPriority `json:"priority,omitempty"`
	StartTime *time.Time    `json:"start_time,omitempty"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
	Page      int           `json:"page"`
	PageSize  int           `json:"page_size"`
}

// SubscriptionOptions 订阅选项
type SubscriptionOptions struct {
	Async      bool          `json:"async"`       // 是否异步处理
	MaxRetries int           `json:"max_retries"` // 最大重试次数
	RetryDelay time.Duration `json:"retry_delay"` // 重试延迟
	Timeout    time.Duration `json:"timeout"`     // 处理超时
	Persistent bool          `json:"persistent"`  // 是否持久化订阅
}

// DefaultSubscriptionOptions 默认订阅选项
func DefaultSubscriptionOptions() *SubscriptionOptions {
	return &SubscriptionOptions{
		Async:      true,
		MaxRetries: 3,
		RetryDelay: 5 * time.Second,
		Timeout:    30 * time.Second,
		Persistent: false,
	}
}

// EventBusOptions 事件总线选项
type EventBusOptions struct {
	BufferSize    int           `json:"buffer_size"`    // 事件缓冲区大小
	WorkerCount   int           `json:"worker_count"`   // 工作协程数量
	BatchSize     int           `json:"batch_size"`     // 批处理大小
	FlushInterval time.Duration `json:"flush_interval"` // 刷新间隔
}

// DefaultEventBusOptions 默认事件总线选项
func DefaultEventBusOptions() *EventBusOptions {
	return &EventBusOptions{
		BufferSize:    1000,
		WorkerCount:   10,
		BatchSize:     50,
		FlushInterval: 1 * time.Second,
	}
}

// EventMetrics 事件指标
type EventMetrics struct {
	TotalEvents     int64 `json:"total_events"`
	PendingEvents   int64 `json:"pending_events"`
	DeliveredEvents int64 `json:"delivered_events"`
	FailedEvents    int64 `json:"failed_events"`
	ExpiredEvents   int64 `json:"expired_events"`
}

// SystemEvents 系统内置事件类型
const (
	// Cron相关事件
	EventTypeCronTaskCreated   = "cron.task.created"
	EventTypeCronTaskUpdated   = "cron.task.updated"
	EventTypeCronTaskDeleted   = "cron.task.deleted"
	EventTypeCronTaskExecuted  = "cron.task.executed"
	EventTypeCronTaskFailed    = "cron.task.failed"
	EventTypeCronTaskCompleted = "cron.task.completed"

	// 系统事件
	EventTypeSystemStartup  = "system.startup"
	EventTypeSystemShutdown = "system.shutdown"
	EventTypeSystemError    = "system.error"
	EventTypeSystemWarning  = "system.warning"

	// 用户事件
	EventTypeUserLogin      = "user.login"
	EventTypeUserLogout     = "user.logout"
	EventTypeUserRegistered = "user.registered"

	// 数据事件
	EventTypeDataCreated = "data.created"
	EventTypeDataUpdated = "data.updated"
	EventTypeDataDeleted = "data.deleted"
)

// CronTaskPayload Cron任务事件负载
type CronTaskPayload struct {
	TaskID     string                 `json:"task_id"`
	TaskName   string                 `json:"task_name"`
	TaskType   string                 `json:"task_type"`
	Status     string                 `json:"status"`
	ExecutedAt time.Time              `json:"executed_at,omitempty"`
	Duration   time.Duration          `json:"duration,omitempty"`
	Result     string                 `json:"result,omitempty"`
	Error      string                 `json:"error,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}
