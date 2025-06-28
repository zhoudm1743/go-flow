package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zhoudm1743/go-flow/core/logger"
)

// Scheduler 调度器接口
type Scheduler interface {
	// 启动调度器
	Start(ctx context.Context) error
	// 停止调度器
	Stop() error
	// 添加任务到调度器
	AddTask(task Task) error
	// 从调度器移除任务
	RemoveTask(taskID string) error
	// 更新任务
	UpdateTask(task Task) error
	// 立即执行任务
	ExecuteTaskNow(taskID string) error
	// 获取任务下次执行时间
	GetNextRunTime(taskID string) (*time.Time, error)
	// 获取调度器状态
	GetStatus() SchedulerStatus
}

// SchedulerStatus 调度器状态
type SchedulerStatus struct {
	Running     bool      `json:"running"`
	StartTime   time.Time `json:"start_time"`
	TaskCount   int       `json:"task_count"`
	ActiveTasks int       `json:"active_tasks"`
}

// CronScheduler cron调度器实现
type CronScheduler struct {
	cron       *cron.Cron
	repository Repository
	logger     logger.Logger
	running    bool
	startTime  time.Time
	tasks      map[string]cron.EntryID // 任务ID到cron条目ID的映射
	mutex      sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewCronScheduler 创建cron调度器
func NewCronScheduler(repo Repository, log logger.Logger) Scheduler {
	// 创建cron实例，使用秒级精度
	c := cron.New(cron.WithSeconds())

	return &CronScheduler{
		cron:       c,
		repository: repo,
		logger:     log,
		tasks:      make(map[string]cron.EntryID),
	}
}

// Start 启动调度器
func (cs *CronScheduler) Start(ctx context.Context) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	if cs.running {
		return fmt.Errorf("调度器已经在运行")
	}

	cs.ctx, cs.cancel = context.WithCancel(ctx)
	cs.startTime = time.Now()
	cs.running = true

	// 启动cron调度器
	cs.cron.Start()

	// 加载现有任务
	if err := cs.loadExistingTasks(); err != nil {
		cs.logger.WithError(err).Error("加载现有任务失败")
	}

	cs.logger.Info("Cron调度器启动成功")
	return nil
}

// Stop 停止调度器
func (cs *CronScheduler) Stop() error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	if !cs.running {
		return nil
	}

	// 停止cron调度器
	cronCtx := cs.cron.Stop()

	// 等待所有任务执行完成
	select {
	case <-cronCtx.Done():
		cs.logger.Info("所有任务执行完成")
	case <-time.After(30 * time.Second):
		cs.logger.Warn("等待任务完成超时")
	}

	if cs.cancel != nil {
		cs.cancel()
	}

	cs.running = false
	cs.tasks = make(map[string]cron.EntryID)

	cs.logger.Info("Cron调度器停止成功")
	return nil
}

// AddTask 添加任务到调度器
func (cs *CronScheduler) AddTask(task Task) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	if !cs.running {
		return fmt.Errorf("调度器未运行")
	}

	// 检查任务状态
	if task.GetStatus() != TaskStatusActive {
		cs.logger.WithFields(map[string]interface{}{
			"task_id": task.GetID(),
			"status":  task.GetStatus(),
		}).Debug("跳过非激活任务")
		return nil
	}

	// 验证cron表达式
	schedule, err := cron.ParseStandard(task.GetCron())
	if err != nil {
		return fmt.Errorf("无效的cron表达式 '%s': %w", task.GetCron(), err)
	}

	// 创建任务执行函数
	jobFunc := cs.createJobFunc(task)

	// 添加到cron调度器
	entryID := cs.cron.Schedule(schedule, jobFunc)

	// 记录任务映射
	cs.tasks[task.GetID()] = entryID

	cs.logger.WithFields(map[string]interface{}{
		"task_id":  task.GetID(),
		"cron":     task.GetCron(),
		"entry_id": entryID,
	}).Info("任务添加到调度器成功")

	return nil
}

// RemoveTask 从调度器移除任务
func (cs *CronScheduler) RemoveTask(taskID string) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	entryID, exists := cs.tasks[taskID]
	if !exists {
		return fmt.Errorf("任务不存在: %s", taskID)
	}

	// 从cron调度器移除
	cs.cron.Remove(entryID)

	// 删除映射记录
	delete(cs.tasks, taskID)

	cs.logger.WithField("task_id", taskID).Info("任务从调度器移除成功")
	return nil
}

// UpdateTask 更新任务
func (cs *CronScheduler) UpdateTask(task Task) error {
	// 先移除旧任务
	if err := cs.RemoveTask(task.GetID()); err != nil {
		// 如果任务不存在，忽略错误
		cs.logger.WithField("task_id", task.GetID()).Debug("移除任务时未找到，可能是新任务")
	}

	// 添加新任务
	return cs.AddTask(task)
}

// ExecuteTaskNow 立即执行任务
func (cs *CronScheduler) ExecuteTaskNow(taskID string) error {
	task, err := cs.repository.GetTask(cs.ctx, taskID)
	if err != nil {
		return fmt.Errorf("获取任务失败: %w", err)
	}

	// 在新的goroutine中执行任务
	go func() {
		cs.executeTask(task)
	}()

	return nil
}

// GetNextRunTime 获取任务下次执行时间
func (cs *CronScheduler) GetNextRunTime(taskID string) (*time.Time, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	entryID, exists := cs.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("任务不在调度器中: %s", taskID)
	}

	entry := cs.cron.Entry(entryID)
	if entry.ID == 0 {
		return nil, fmt.Errorf("任务条目不存在: %s", taskID)
	}

	nextTime := entry.Next
	return &nextTime, nil
}

// GetStatus 获取调度器状态
func (cs *CronScheduler) GetStatus() SchedulerStatus {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	activeTasks := 0
	for taskID := range cs.tasks {
		if task, err := cs.repository.GetTask(cs.ctx, taskID); err == nil {
			if task.GetStatus() == TaskStatusActive {
				activeTasks++
			}
		}
	}

	return SchedulerStatus{
		Running:     cs.running,
		StartTime:   cs.startTime,
		TaskCount:   len(cs.tasks),
		ActiveTasks: activeTasks,
	}
}

// loadExistingTasks 加载现有任务
func (cs *CronScheduler) loadExistingTasks() error {
	req := &TaskListRequest{
		Page:     1,
		PageSize: 1000, // 加载所有任务
	}

	tasks, _, err := cs.repository.ListTasks(cs.ctx, req)
	if err != nil {
		return fmt.Errorf("获取任务列表失败: %w", err)
	}

	var loadedCount int
	for _, taskInfo := range tasks {
		task, err := cs.repository.GetTask(cs.ctx, taskInfo.ID)
		if err != nil {
			cs.logger.WithField("task_id", taskInfo.ID).Warn("获取任务详情失败")
			continue
		}

		if err := cs.AddTask(task); err != nil {
			cs.logger.WithFields(map[string]interface{}{
				"task_id": task.GetID(),
				"error":   err.Error(),
			}).Warn("添加任务到调度器失败")
			continue
		}

		loadedCount++
	}

	cs.logger.WithField("count", loadedCount).Info("现有任务加载完成")
	return nil
}

// createJobFunc 创建任务执行函数
func (cs *CronScheduler) createJobFunc(task Task) cron.Job {
	return cron.FuncJob(func() {
		cs.executeTask(task)
	})
}

// executeTask 执行任务
func (cs *CronScheduler) executeTask(task Task) {
	taskID := task.GetID()

	// 检查任务状态
	currentStatus, err := cs.repository.GetTaskStatus(cs.ctx, taskID)
	if err != nil {
		cs.logger.WithField("task_id", taskID).Error("获取任务状态失败")
		return
	}

	if currentStatus != TaskStatusActive {
		cs.logger.WithFields(map[string]interface{}{
			"task_id": taskID,
			"status":  currentStatus,
		}).Debug("跳过非激活任务执行")
		return
	}

	cs.logger.WithField("task_id", taskID).Info("开始执行任务")

	// 更新任务执行次数
	if err := cs.repository.IncrTaskRunCount(cs.ctx, taskID); err != nil {
		cs.logger.WithField("task_id", taskID).Warn("更新任务执行次数失败")
	}

	// 设置最后执行时间
	if err := cs.repository.SetTaskLastRunTime(cs.ctx, taskID, time.Now()); err != nil {
		cs.logger.WithField("task_id", taskID).Warn("设置最后执行时间失败")
	}

	// 创建执行上下文
	execCtx, cancel := context.WithCancel(cs.ctx)
	defer cancel()

	// 执行任务
	result := task.Execute(execCtx)

	// 保存执行结果
	if err := cs.repository.SaveExecResult(cs.ctx, result); err != nil {
		cs.logger.WithFields(map[string]interface{}{
			"task_id": taskID,
			"error":   err.Error(),
		}).Error("保存执行结果失败")
	}

	// 记录执行完成日志
	logFields := map[string]interface{}{
		"task_id":  taskID,
		"status":   result.Status,
		"duration": result.Duration,
	}

	if result.Status == ExecStatusSuccess {
		cs.logger.WithFields(logFields).Info("任务执行完成")
	} else {
		logFields["error"] = result.Error
		cs.logger.WithFields(logFields).Error("任务执行失败")
	}
}
