package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/zhoudm1743/go-flow/core/cache"
	"github.com/zhoudm1743/go-flow/core/logger"
	"github.com/zhoudm1743/go-flow/pkg/response"
)

// Service cron服务接口
type Service interface {
	// 任务管理
	CreateTask(ctx context.Context, req *CreateTaskRequest) (*TaskInfo, error)
	UpdateTask(ctx context.Context, taskID string, req *UpdateTaskRequest) (*TaskInfo, error)
	DeleteTask(ctx context.Context, taskID string) error
	GetTask(ctx context.Context, taskID string) (*TaskInfo, error)
	ListTasks(ctx context.Context, req *TaskListRequest) (*response.PageResult[*TaskInfo], error)

	// 任务控制
	StartTask(ctx context.Context, taskID string) error
	StopTask(ctx context.Context, taskID string) error
	PauseTask(ctx context.Context, taskID string) error
	ExecuteTaskNow(ctx context.Context, taskID string) error

	// 执行历史
	GetTaskExecHistory(ctx context.Context, req *TaskExecHistoryRequest) (*response.PageResult[*ExecResult], error)

	// 系统任务处理器
	GetSystemTaskHandlers(ctx context.Context) (map[string]SystemTaskHandler, error)

	// 调度器状态
	GetSchedulerStatus(ctx context.Context) (*SchedulerStatus, error)
	StartScheduler(ctx context.Context) error
	StopScheduler(ctx context.Context) error
}

// CronService cron服务实现
type CronService struct {
	repository Repository
	scheduler  Scheduler
	cache      cache.Cache
	logger     logger.Logger
}

// NewCronService 创建cron服务
func NewCronService(repo Repository, scheduler Scheduler, cache cache.Cache, log logger.Logger) Service {
	return &CronService{
		repository: repo,
		scheduler:  scheduler,
		cache:      cache,
		logger:     log,
	}
}

// CreateTask 创建任务
func (cs *CronService) CreateTask(ctx context.Context, req *CreateTaskRequest) (*TaskInfo, error) {
	// 验证cron表达式
	if _, err := cron.ParseStandard(req.Cron); err != nil {
		return nil, fmt.Errorf("无效的cron表达式: %w", err)
	}

	// 生成任务ID
	taskID := uuid.New().String()

	// 创建基础任务信息
	baseTask := BaseTask{
		ID:          taskID,
		Name:        req.Name,
		Type:        req.Type,
		Cron:        req.Cron,
		Status:      TaskStatusActive,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   "system", // TODO: 从上下文获取用户信息
	}

	// 根据类型创建任务
	var task Task

	switch req.Type {
	case TaskTypeHTTP:
		httpConfig, err := cs.parseHTTPConfig(req.Config)
		if err != nil {
			return nil, fmt.Errorf("解析HTTP配置失败: %w", err)
		}
		task = NewHTTPTask(baseTask, httpConfig, cs.logger)

	case TaskTypeSystem:
		return nil, fmt.Errorf("系统任务不支持通过API创建")

	default:
		return nil, fmt.Errorf("不支持的任务类型: %s", req.Type)
	}

	// 验证任务配置
	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("任务配置验证失败: %w", err)
	}

	// 保存任务到仓库
	if err := cs.repository.CreateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("保存任务失败: %w", err)
	}

	// 添加到调度器
	if err := cs.scheduler.AddTask(task); err != nil {
		cs.logger.WithField("task_id", taskID).Warn("添加任务到调度器失败")
	}

	// 返回任务信息
	return cs.GetTask(ctx, taskID)
}

// UpdateTask 更新任务
func (cs *CronService) UpdateTask(ctx context.Context, taskID string, req *UpdateTaskRequest) (*TaskInfo, error) {
	// 获取现有任务
	task, err := cs.repository.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("获取任务失败: %w", err)
	}

	// 更新基础信息
	if req.Name != nil {
		// 这里需要修改任务的名称，但Task接口没有提供修改方法
		// 实际实现中可能需要扩展接口或使用类型断言
	}

	if req.Description != nil {
		// 同样需要处理描述字段
	}

	if req.Cron != nil {
		// 验证新的cron表达式
		if _, err := cron.ParseStandard(*req.Cron); err != nil {
			return nil, fmt.Errorf("无效的cron表达式: %w", err)
		}
		// 需要更新cron表达式
	}

	if req.Status != nil {
		task.SetStatus(*req.Status)
	}

	if req.Config != nil {
		// 根据任务类型更新配置
		switch task.GetType() {
		case TaskTypeHTTP:
			if httpTask, ok := task.(*HTTPTask); ok {
				if err := httpTask.UpdateConfig(req.Config); err != nil {
					return nil, fmt.Errorf("更新HTTP配置失败: %w", err)
				}
			}
		case TaskTypeSystem:
			if systemTask, ok := task.(*SystemTask); ok {
				if err := systemTask.UpdateConfig(req.Config); err != nil {
					return nil, fmt.Errorf("更新系统配置失败: %w", err)
				}
			}
		}
	}

	// 验证更新后的任务配置
	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("任务配置验证失败: %w", err)
	}

	// 更新任务到仓库
	if err := cs.repository.UpdateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("更新任务失败: %w", err)
	}

	// 更新调度器中的任务
	if err := cs.scheduler.UpdateTask(task); err != nil {
		cs.logger.WithField("task_id", taskID).Warn("更新调度器中的任务失败")
	}

	return cs.GetTask(ctx, taskID)
}

// DeleteTask 删除任务
func (cs *CronService) DeleteTask(ctx context.Context, taskID string) error {
	// 从调度器移除任务
	if err := cs.scheduler.RemoveTask(taskID); err != nil {
		cs.logger.WithField("task_id", taskID).Warn("从调度器移除任务失败")
	}

	// 从仓库删除任务
	if err := cs.repository.DeleteTask(ctx, taskID); err != nil {
		return fmt.Errorf("删除任务失败: %w", err)
	}

	cs.logger.WithField("task_id", taskID).Info("任务删除成功")
	return nil
}

// GetTask 获取任务详情
func (cs *CronService) GetTask(ctx context.Context, taskID string) (*TaskInfo, error) {
	task, err := cs.repository.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("获取任务失败: %w", err)
	}

	taskInfo := &TaskInfo{
		BaseTask: BaseTask{
			ID:          task.GetID(),
			Name:        task.GetName(),
			Type:        task.GetType(),
			Cron:        task.GetCron(),
			Status:      task.GetStatus(),
			Description: "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Config: task.GetConfig(),
	}

	// 获取下次执行时间
	if nextTime, err := cs.scheduler.GetNextRunTime(taskID); err == nil {
		taskInfo.NextRunTime = nextTime
	}

	// 获取统计信息
	if stats, err := cs.repository.GetTaskStats(ctx, taskID); err == nil {
		taskInfo.RunCount = stats.RunCount
		if !stats.LastRunTime.IsZero() {
			taskInfo.LastRunTime = &stats.LastRunTime
		}
	}

	// 获取最后执行结果
	if lastResult, err := cs.repository.GetLastExecResult(ctx, taskID); err == nil {
		taskInfo.LastResult = lastResult
	}

	return taskInfo, nil
}

// ListTasks 获取任务列表
func (cs *CronService) ListTasks(ctx context.Context, req *TaskListRequest) (*response.PageResult[*TaskInfo], error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	tasks, total, err := cs.repository.ListTasks(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("获取任务列表失败: %w", err)
	}

	// 补充调度器信息
	for _, task := range tasks {
		if nextTime, err := cs.scheduler.GetNextRunTime(task.ID); err == nil {
			task.NextRunTime = nextTime
		}
	}

	return response.NewPageResult(tasks, total, int64(req.Page), int64(req.PageSize)), nil
}

// StartTask 启动任务
func (cs *CronService) StartTask(ctx context.Context, taskID string) error {
	if err := cs.repository.SetTaskStatus(ctx, taskID, TaskStatusActive); err != nil {
		return fmt.Errorf("设置任务状态失败: %w", err)
	}

	// 重新添加到调度器
	task, err := cs.repository.GetTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("获取任务失败: %w", err)
	}

	if err := cs.scheduler.AddTask(task); err != nil {
		return fmt.Errorf("添加任务到调度器失败: %w", err)
	}

	cs.logger.WithField("task_id", taskID).Info("任务启动成功")
	return nil
}

// StopTask 停止任务
func (cs *CronService) StopTask(ctx context.Context, taskID string) error {
	if err := cs.repository.SetTaskStatus(ctx, taskID, TaskStatusStopped); err != nil {
		return fmt.Errorf("设置任务状态失败: %w", err)
	}

	// 从调度器移除
	if err := cs.scheduler.RemoveTask(taskID); err != nil {
		cs.logger.WithField("task_id", taskID).Warn("从调度器移除任务失败")
	}

	cs.logger.WithField("task_id", taskID).Info("任务停止成功")
	return nil
}

// PauseTask 暂停任务
func (cs *CronService) PauseTask(ctx context.Context, taskID string) error {
	if err := cs.repository.SetTaskStatus(ctx, taskID, TaskStatusPaused); err != nil {
		return fmt.Errorf("设置任务状态失败: %w", err)
	}

	// 从调度器移除
	if err := cs.scheduler.RemoveTask(taskID); err != nil {
		cs.logger.WithField("task_id", taskID).Warn("从调度器移除任务失败")
	}

	cs.logger.WithField("task_id", taskID).Info("任务暂停成功")
	return nil
}

// ExecuteTaskNow 立即执行任务
func (cs *CronService) ExecuteTaskNow(ctx context.Context, taskID string) error {
	if err := cs.scheduler.ExecuteTaskNow(taskID); err != nil {
		return fmt.Errorf("执行任务失败: %w", err)
	}

	cs.logger.WithField("task_id", taskID).Info("任务手动执行成功")
	return nil
}

// GetTaskExecHistory 获取任务执行历史
func (cs *CronService) GetTaskExecHistory(ctx context.Context, req *TaskExecHistoryRequest) (*response.PageResult[*ExecResult], error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	results, total, err := cs.repository.GetExecHistory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("获取执行历史失败: %w", err)
	}

	return response.NewPageResult(results, total, int64(req.Page), int64(req.PageSize)), nil
}

// GetSystemTaskHandlers 获取系统任务处理器
func (cs *CronService) GetSystemTaskHandlers(ctx context.Context) (map[string]SystemTaskHandler, error) {
	return GetAllSystemTaskHandlers(), nil
}

// GetSchedulerStatus 获取调度器状态
func (cs *CronService) GetSchedulerStatus(ctx context.Context) (*SchedulerStatus, error) {
	status := cs.scheduler.GetStatus()
	return &status, nil
}

// StartScheduler 启动调度器
func (cs *CronService) StartScheduler(ctx context.Context) error {
	if err := cs.scheduler.Start(ctx); err != nil {
		return fmt.Errorf("启动调度器失败: %w", err)
	}

	cs.logger.Info("调度器启动成功")
	return nil
}

// StopScheduler 停止调度器
func (cs *CronService) StopScheduler(ctx context.Context) error {
	if err := cs.scheduler.Stop(); err != nil {
		return fmt.Errorf("停止调度器失败: %w", err)
	}

	cs.logger.Info("调度器停止成功")
	return nil
}

// parseHTTPConfig 解析HTTP配置
func (cs *CronService) parseHTTPConfig(config interface{}) (HTTPConfig, error) {
	var httpConfig HTTPConfig

	configData, err := json.Marshal(config)
	if err != nil {
		return httpConfig, fmt.Errorf("配置序列化失败: %w", err)
	}

	if err := json.Unmarshal(configData, &httpConfig); err != nil {
		return httpConfig, fmt.Errorf("配置反序列化失败: %w", err)
	}

	return httpConfig, nil
}
