package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zhoudm1743/go-flow/core/cache"
	"github.com/zhoudm1743/go-flow/core/logger"
)

// Repository 任务仓库接口
type Repository interface {
	// 任务管理
	CreateTask(ctx context.Context, task Task) error
	UpdateTask(ctx context.Context, task Task) error
	DeleteTask(ctx context.Context, taskID string) error
	GetTask(ctx context.Context, taskID string) (Task, error)
	ListTasks(ctx context.Context, req *TaskListRequest) ([]*TaskInfo, int64, error)

	// 任务状态管理
	SetTaskStatus(ctx context.Context, taskID string, status TaskStatus) error
	GetTaskStatus(ctx context.Context, taskID string) (TaskStatus, error)

	// 执行历史管理
	SaveExecResult(ctx context.Context, result *ExecResult) error
	GetExecHistory(ctx context.Context, req *TaskExecHistoryRequest) ([]*ExecResult, int64, error)
	GetLastExecResult(ctx context.Context, taskID string) (*ExecResult, error)

	// 统计信息
	GetTaskStats(ctx context.Context, taskID string) (*TaskStats, error)
	IncrTaskRunCount(ctx context.Context, taskID string) error
	SetTaskLastRunTime(ctx context.Context, taskID string, t time.Time) error
}

// TaskStats 任务统计信息
type TaskStats struct {
	TaskID       string    `json:"task_id"`
	RunCount     int64     `json:"run_count"`
	SuccessCount int64     `json:"success_count"`
	FailedCount  int64     `json:"failed_count"`
	LastRunTime  time.Time `json:"last_run_time"`
	AvgDuration  float64   `json:"avg_duration"`
}

// RedisRepository Redis任务仓库实现
type RedisRepository struct {
	cache  cache.Cache
	logger logger.Logger
}

// NewRedisRepository 创建Redis任务仓库
func NewRedisRepository(c cache.Cache, log logger.Logger) Repository {
	return &RedisRepository{
		cache:  c,
		logger: log,
	}
}

// Redis键定义
const (
	TaskKeyPrefix        = "cron:task:"        // 任务信息键前缀
	TaskListKey          = "cron:tasks"        // 任务列表键
	TaskStatsKeyPrefix   = "cron:stats:"       // 任务统计键前缀
	ExecHistoryKeyPrefix = "cron:history:"     // 执行历史键前缀
	ExecHistoryListKey   = "cron:history_list" // 执行历史列表键
)

// CreateTask 创建任务
func (r *RedisRepository) CreateTask(ctx context.Context, task Task) error {
	taskKey := TaskKeyPrefix + task.GetID()

	// 将任务序列化为JSON
	taskData, err := r.serializeTask(task)
	if err != nil {
		return fmt.Errorf("序列化任务失败: %w", err)
	}

	// 保存任务信息
	if err := r.cache.SetCtx(ctx, taskKey, taskData, 0); err != nil {
		return fmt.Errorf("保存任务失败: %w", err)
	}

	// 添加到任务列表
	if _, err := r.cache.SAddCtx(ctx, TaskListKey, task.GetID()); err != nil {
		return fmt.Errorf("添加到任务列表失败: %w", err)
	}

	// 初始化任务统计
	statsKey := TaskStatsKeyPrefix + task.GetID()
	stats := &TaskStats{
		TaskID:       task.GetID(),
		RunCount:     0,
		SuccessCount: 0,
		FailedCount:  0,
	}

	statsData, _ := json.Marshal(stats)
	if err := r.cache.SetCtx(ctx, statsKey, string(statsData), 0); err != nil {
		r.logger.WithField("task_id", task.GetID()).Warn("初始化任务统计失败")
	}

	return nil
}

// UpdateTask 更新任务
func (r *RedisRepository) UpdateTask(ctx context.Context, task Task) error {
	taskKey := TaskKeyPrefix + task.GetID()

	// 检查任务是否存在
	exists, err := r.cache.ExistsCtx(ctx, taskKey)
	if err != nil {
		return fmt.Errorf("检查任务存在性失败: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("任务不存在: %s", task.GetID())
	}

	// 将任务序列化为JSON
	taskData, err := r.serializeTask(task)
	if err != nil {
		return fmt.Errorf("序列化任务失败: %w", err)
	}

	// 更新任务信息
	if err := r.cache.SetCtx(ctx, taskKey, taskData, 0); err != nil {
		return fmt.Errorf("更新任务失败: %w", err)
	}

	return nil
}

// DeleteTask 删除任务
func (r *RedisRepository) DeleteTask(ctx context.Context, taskID string) error {
	taskKey := TaskKeyPrefix + taskID

	// 删除任务信息
	if _, err := r.cache.DelCtx(ctx, taskKey); err != nil {
		return fmt.Errorf("删除任务失败: %w", err)
	}

	// 从任务列表中移除
	if _, err := r.cache.SRemCtx(ctx, TaskListKey, taskID); err != nil {
		r.logger.WithField("task_id", taskID).Warn("从任务列表移除失败")
	}

	// 删除任务统计
	statsKey := TaskStatsKeyPrefix + taskID
	if _, err := r.cache.DelCtx(ctx, statsKey); err != nil {
		r.logger.WithField("task_id", taskID).Warn("删除任务统计失败")
	}

	return nil
}

// GetTask 获取任务
func (r *RedisRepository) GetTask(ctx context.Context, taskID string) (Task, error) {
	taskKey := TaskKeyPrefix + taskID

	taskData, err := r.cache.GetCtx(ctx, taskKey)
	if err != nil {
		return nil, fmt.Errorf("获取任务失败: %w", err)
	}

	return r.deserializeTask(taskData)
}

// ListTasks 列出任务
func (r *RedisRepository) ListTasks(ctx context.Context, req *TaskListRequest) ([]*TaskInfo, int64, error) {
	// 获取所有任务ID
	taskIDs, err := r.cache.SMembersCtx(ctx, TaskListKey)
	if err != nil {
		return nil, 0, fmt.Errorf("获取任务列表失败: %w", err)
	}

	var tasks []*TaskInfo

	// 获取每个任务的详细信息
	for _, taskID := range taskIDs {
		task, err := r.GetTask(ctx, taskID)
		if err != nil {
			r.logger.WithField("task_id", taskID).Warn("获取任务信息失败")
			continue
		}

		// 应用过滤条件
		if req.Type != "" && task.GetType() != req.Type {
			continue
		}
		if req.Status != "" && task.GetStatus() != req.Status {
			continue
		}
		if req.Keyword != "" && !strings.Contains(task.GetName(), req.Keyword) {
			continue
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

		// 获取任务统计信息
		if stats, err := r.GetTaskStats(ctx, taskID); err == nil {
			taskInfo.RunCount = stats.RunCount
			if !stats.LastRunTime.IsZero() {
				taskInfo.LastRunTime = &stats.LastRunTime
			}
		}

		// 获取最后执行结果
		if lastResult, err := r.GetLastExecResult(ctx, taskID); err == nil {
			taskInfo.LastResult = lastResult
		}

		tasks = append(tasks, taskInfo)
	}

	// 分页处理
	total := int64(len(tasks))

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	if start >= len(tasks) {
		return []*TaskInfo{}, total, nil
	}
	if end > len(tasks) {
		end = len(tasks)
	}

	return tasks[start:end], total, nil
}

// SetTaskStatus 设置任务状态
func (r *RedisRepository) SetTaskStatus(ctx context.Context, taskID string, status TaskStatus) error {
	task, err := r.GetTask(ctx, taskID)
	if err != nil {
		return err
	}

	task.SetStatus(status)
	return r.UpdateTask(ctx, task)
}

// GetTaskStatus 获取任务状态
func (r *RedisRepository) GetTaskStatus(ctx context.Context, taskID string) (TaskStatus, error) {
	task, err := r.GetTask(ctx, taskID)
	if err != nil {
		return "", err
	}

	return task.GetStatus(), nil
}

// SaveExecResult 保存执行结果
func (r *RedisRepository) SaveExecResult(ctx context.Context, result *ExecResult) error {
	// 生成执行结果ID
	resultID := fmt.Sprintf("%s_%d", result.TaskID, time.Now().UnixNano())
	resultKey := ExecHistoryKeyPrefix + resultID

	// 序列化执行结果
	resultData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("序列化执行结果失败: %w", err)
	}

	// 保存执行结果
	if err := r.cache.SetCtx(ctx, resultKey, string(resultData), 30*24*time.Hour); err != nil {
		return fmt.Errorf("保存执行结果失败: %w", err)
	}

	// 添加到执行历史列表（按时间排序）
	score := float64(result.StartTime.Unix())
	if _, err := r.cache.ZAddCtx(ctx, ExecHistoryListKey, redis.Z{Score: score, Member: resultID}); err != nil {
		r.logger.WithField("result_id", resultID).Warn("添加到执行历史列表失败")
	}

	// 更新任务统计
	if err := r.updateTaskStats(ctx, result); err != nil {
		r.logger.WithField("task_id", result.TaskID).Warn("更新任务统计失败")
	}

	return nil
}

// GetExecHistory 获取执行历史
func (r *RedisRepository) GetExecHistory(ctx context.Context, req *TaskExecHistoryRequest) ([]*ExecResult, int64, error) {
	// 获取所有执行结果ID（按时间倒序）
	resultIDs, err := r.cache.ZRangeCtx(ctx, ExecHistoryListKey, 0, -1)
	if err != nil {
		return nil, 0, fmt.Errorf("获取执行历史列表失败: %w", err)
	}

	var results []*ExecResult

	// 获取每个执行结果的详细信息
	for i := len(resultIDs) - 1; i >= 0; i-- { // 倒序遍历
		resultID := resultIDs[i]
		resultKey := ExecHistoryKeyPrefix + resultID

		resultData, err := r.cache.GetCtx(ctx, resultKey)
		if err != nil {
			continue
		}

		var result ExecResult
		if err := json.Unmarshal([]byte(resultData), &result); err != nil {
			continue
		}

		// 应用过滤条件
		if result.TaskID != req.TaskID {
			continue
		}
		if req.Status != "" && result.Status != req.Status {
			continue
		}
		if req.StartTime != nil && result.StartTime.Before(*req.StartTime) {
			continue
		}
		if req.EndTime != nil && result.StartTime.After(*req.EndTime) {
			continue
		}

		results = append(results, &result)
	}

	// 分页处理
	total := int64(len(results))

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	if start >= len(results) {
		return []*ExecResult{}, total, nil
	}
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], total, nil
}

// GetLastExecResult 获取最后执行结果
func (r *RedisRepository) GetLastExecResult(ctx context.Context, taskID string) (*ExecResult, error) {
	req := &TaskExecHistoryRequest{
		TaskID:   taskID,
		Page:     1,
		PageSize: 1,
	}

	results, _, err := r.GetExecHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("未找到执行结果")
	}

	return results[0], nil
}

// GetTaskStats 获取任务统计信息
func (r *RedisRepository) GetTaskStats(ctx context.Context, taskID string) (*TaskStats, error) {
	statsKey := TaskStatsKeyPrefix + taskID

	statsData, err := r.cache.GetCtx(ctx, statsKey)
	if err != nil {
		return nil, fmt.Errorf("获取任务统计失败: %w", err)
	}

	var stats TaskStats
	if err := json.Unmarshal([]byte(statsData), &stats); err != nil {
		return nil, fmt.Errorf("反序列化任务统计失败: %w", err)
	}

	return &stats, nil
}

// IncrTaskRunCount 增加任务执行次数
func (r *RedisRepository) IncrTaskRunCount(ctx context.Context, taskID string) error {
	statsKey := TaskStatsKeyPrefix + taskID

	stats, err := r.GetTaskStats(ctx, taskID)
	if err != nil {
		// 如果统计不存在，则创建
		stats = &TaskStats{TaskID: taskID}
	}

	stats.RunCount++

	statsData, _ := json.Marshal(stats)
	return r.cache.SetCtx(ctx, statsKey, string(statsData), 0)
}

// SetTaskLastRunTime 设置任务最后执行时间
func (r *RedisRepository) SetTaskLastRunTime(ctx context.Context, taskID string, t time.Time) error {
	statsKey := TaskStatsKeyPrefix + taskID

	stats, err := r.GetTaskStats(ctx, taskID)
	if err != nil {
		stats = &TaskStats{TaskID: taskID}
	}

	stats.LastRunTime = t

	statsData, _ := json.Marshal(stats)
	return r.cache.SetCtx(ctx, statsKey, string(statsData), 0)
}

// updateTaskStats 更新任务统计
func (r *RedisRepository) updateTaskStats(ctx context.Context, result *ExecResult) error {
	statsKey := TaskStatsKeyPrefix + result.TaskID

	stats, err := r.GetTaskStats(ctx, result.TaskID)
	if err != nil {
		stats = &TaskStats{TaskID: result.TaskID}
	}

	stats.RunCount++
	if result.Status == ExecStatusSuccess {
		stats.SuccessCount++
	} else {
		stats.FailedCount++
	}
	stats.LastRunTime = result.StartTime

	// 计算平均执行时间
	if stats.RunCount > 0 {
		stats.AvgDuration = (stats.AvgDuration*float64(stats.RunCount-1) + float64(result.Duration.Milliseconds())) / float64(stats.RunCount)
	}

	statsData, _ := json.Marshal(stats)
	return r.cache.SetCtx(ctx, statsKey, string(statsData), 0)
}

// serializeTask 序列化任务
func (r *RedisRepository) serializeTask(task Task) (string, error) {
	taskData := map[string]interface{}{
		"id":         task.GetID(),
		"name":       task.GetName(),
		"type":       task.GetType(),
		"cron":       task.GetCron(),
		"status":     task.GetStatus(),
		"config":     task.GetConfig(),
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	data, err := json.Marshal(taskData)
	return string(data), err
}

// deserializeTask 反序列化任务
func (r *RedisRepository) deserializeTask(data string) (Task, error) {
	var taskData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &taskData); err != nil {
		return nil, fmt.Errorf("反序列化任务数据失败: %w", err)
	}

	taskType, _ := taskData["type"].(string)

	// 根据类型创建相应的任务
	switch TaskType(taskType) {
	case TaskTypeHTTP:
		return r.deserializeHTTPTask(taskData)
	case TaskTypeSystem:
		return r.deserializeSystemTask(taskData)
	default:
		return nil, fmt.Errorf("不支持的任务类型: %s", taskType)
	}
}

// deserializeHTTPTask 反序列化HTTP任务
func (r *RedisRepository) deserializeHTTPTask(taskData map[string]interface{}) (Task, error) {
	baseTask := BaseTask{
		ID:     taskData["id"].(string),
		Name:   taskData["name"].(string),
		Type:   TaskType(taskData["type"].(string)),
		Cron:   taskData["cron"].(string),
		Status: TaskStatus(taskData["status"].(string)),
	}

	configData, _ := json.Marshal(taskData["config"])
	var config HTTPConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("反序列化HTTP配置失败: %w", err)
	}

	return NewHTTPTask(baseTask, config, r.logger), nil
}

// deserializeSystemTask 反序列化系统任务
func (r *RedisRepository) deserializeSystemTask(taskData map[string]interface{}) (Task, error) {
	baseTask := BaseTask{
		ID:     taskData["id"].(string),
		Name:   taskData["name"].(string),
		Type:   TaskType(taskData["type"].(string)),
		Cron:   taskData["cron"].(string),
		Status: TaskStatus(taskData["status"].(string)),
	}

	configData, _ := json.Marshal(taskData["config"])
	var config SystemConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("反序列化系统配置失败: %w", err)
	}

	return NewSystemTask(baseTask, config, r.logger)
}
