package cron

import (
	"context"
	"encoding/json"
	"time"
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypeSystem TaskType = "system" // 系统任务
	TaskTypeHTTP   TaskType = "http"   // HTTP任务
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusActive  TaskStatus = "active"  // 激活状态
	TaskStatusPaused  TaskStatus = "paused"  // 暂停状态
	TaskStatusStopped TaskStatus = "stopped" // 停止状态
)

// ExecStatus 执行状态
type ExecStatus string

const (
	ExecStatusRunning ExecStatus = "running" // 执行中
	ExecStatusSuccess ExecStatus = "success" // 执行成功
	ExecStatusFailed  ExecStatus = "failed"  // 执行失败
	ExecStatusTimeout ExecStatus = "timeout" // 执行超时
)

// Task 任务接口
type Task interface {
	// GetID 获取任务ID
	GetID() string
	// GetName 获取任务名称
	GetName() string
	// GetType 获取任务类型
	GetType() TaskType
	// GetCron 获取cron表达式
	GetCron() string
	// GetStatus 获取任务状态
	GetStatus() TaskStatus
	// SetStatus 设置任务状态
	SetStatus(status TaskStatus)
	// Execute 执行任务
	Execute(ctx context.Context) *ExecResult
	// GetConfig 获取任务配置
	GetConfig() interface{}
	// Validate 验证任务配置
	Validate() error
}

// BaseTask 基础任务结构
type BaseTask struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Type        TaskType   `json:"type"`
	Cron        string     `json:"cron"`
	Status      TaskStatus `json:"status"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedBy   string     `json:"created_by"`
}

// GetID 获取任务ID
func (bt *BaseTask) GetID() string {
	return bt.ID
}

// GetName 获取任务名称
func (bt *BaseTask) GetName() string {
	return bt.Name
}

// GetType 获取任务类型
func (bt *BaseTask) GetType() TaskType {
	return bt.Type
}

// GetCron 获取cron表达式
func (bt *BaseTask) GetCron() string {
	return bt.Cron
}

// GetStatus 获取任务状态
func (bt *BaseTask) GetStatus() TaskStatus {
	return bt.Status
}

// SetStatus 设置任务状态
func (bt *BaseTask) SetStatus(status TaskStatus) {
	bt.Status = status
	bt.UpdatedAt = time.Now()
}

// ExecResult 执行结果
type ExecResult struct {
	TaskID    string        `json:"task_id"`
	Status    ExecStatus    `json:"status"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Output    string        `json:"output"`
	Error     string        `json:"error"`
	Retry     int           `json:"retry"`
}

// HTTPConfig HTTP任务配置
type HTTPConfig struct {
	URL          string            `json:"url"`
	Method       string            `json:"method"`
	Headers      map[string]string `json:"headers"`
	Body         string            `json:"body"`
	Timeout      time.Duration     `json:"timeout"`
	RetryCount   int               `json:"retry_count"`
	RetryDelay   time.Duration     `json:"retry_delay"`
	ExpectedCode int               `json:"expected_code"`
}

// SystemConfig 系统任务配置
type SystemConfig struct {
	HandlerName string                 `json:"handler_name"`
	Parameters  map[string]interface{} `json:"parameters"`
	Timeout     time.Duration          `json:"timeout"`
	RetryCount  int                    `json:"retry_count"`
	RetryDelay  time.Duration          `json:"retry_delay"`
}

// TaskInfo 任务信息（用于API返回）
type TaskInfo struct {
	BaseTask
	Config      interface{} `json:"config"`
	NextRunTime *time.Time  `json:"next_run_time"`
	LastRunTime *time.Time  `json:"last_run_time"`
	LastResult  *ExecResult `json:"last_result"`
	RunCount    int64       `json:"run_count"`
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Name        string      `json:"name" validate:"required,min=1,max=100"`
	Description string      `json:"description" validate:"max=500"`
	Cron        string      `json:"cron" validate:"required"`
	Type        TaskType    `json:"type" validate:"required,oneof=http"`
	Config      interface{} `json:"config" validate:"required"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Name        *string     `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string     `json:"description,omitempty" validate:"omitempty,max=500"`
	Cron        *string     `json:"cron,omitempty"`
	Status      *TaskStatus `json:"status,omitempty" validate:"omitempty,oneof=active paused stopped"`
	Config      interface{} `json:"config,omitempty"`
}

// TaskListRequest 任务列表请求
type TaskListRequest struct {
	Page     int        `form:"page" validate:"min=1"`
	PageSize int        `form:"page_size" validate:"min=1,max=100"`
	Type     TaskType   `form:"type,omitempty"`
	Status   TaskStatus `form:"status,omitempty"`
	Keyword  string     `form:"keyword,omitempty"`
}

// TaskExecHistoryRequest 任务执行历史请求
type TaskExecHistoryRequest struct {
	TaskID    string     `form:"task_id" validate:"required"`
	Page      int        `form:"page" validate:"min=1"`
	PageSize  int        `form:"page_size" validate:"min=1,max=100"`
	Status    ExecStatus `form:"status,omitempty"`
	StartTime *time.Time `form:"start_time,omitempty"`
	EndTime   *time.Time `form:"end_time,omitempty"`
}

// ToJSON 转换为JSON字符串
func (t *BaseTask) ToJSON() (string, error) {
	data, err := json.Marshal(t)
	return string(data), err
}

// FromJSON 从JSON字符串解析
func (t *BaseTask) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), t)
}
