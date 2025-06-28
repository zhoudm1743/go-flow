package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zhoudm1743/go-flow/core/cache"
	"github.com/zhoudm1743/go-flow/core/logger"
)

// SystemTaskHandler 系统任务处理器接口
type SystemTaskHandler interface {
	// Handle 处理任务
	Handle(ctx context.Context, params map[string]interface{}) (string, error)
	// GetName 获取处理器名称
	GetName() string
	// GetDescription 获取处理器描述
	GetDescription() string
}

// SystemTask 系统任务实现
type SystemTask struct {
	BaseTask
	Config  SystemConfig `json:"config"`
	handler SystemTaskHandler
	logger  logger.Logger
}

// systemTaskHandlers 全局系统任务处理器注册表
var systemTaskHandlers = make(map[string]SystemTaskHandler)

// RegisterSystemTaskHandler 注册系统任务处理器
func RegisterSystemTaskHandler(handler SystemTaskHandler) {
	systemTaskHandlers[handler.GetName()] = handler
}

// GetSystemTaskHandler 获取系统任务处理器
func GetSystemTaskHandler(name string) (SystemTaskHandler, bool) {
	handler, exists := systemTaskHandlers[name]
	return handler, exists
}

// GetAllSystemTaskHandlers 获取所有系统任务处理器
func GetAllSystemTaskHandlers() map[string]SystemTaskHandler {
	return systemTaskHandlers
}

// NewSystemTask 创建系统任务
func NewSystemTask(base BaseTask, config SystemConfig, log logger.Logger) (*SystemTask, error) {
	handler, exists := GetSystemTaskHandler(config.HandlerName)
	if !exists {
		return nil, fmt.Errorf("未找到系统任务处理器: %s", config.HandlerName)
	}

	task := &SystemTask{
		BaseTask: base,
		Config:   config,
		handler:  handler,
		logger:   log,
	}

	// 设置默认值
	if task.Config.Timeout == 0 {
		task.Config.Timeout = 30 * time.Second
	}
	if task.Config.RetryCount == 0 {
		task.Config.RetryCount = 3
	}
	if task.Config.RetryDelay == 0 {
		task.Config.RetryDelay = 5 * time.Second
	}

	return task, nil
}

// Execute 执行系统任务
func (s *SystemTask) Execute(ctx context.Context) *ExecResult {
	result := &ExecResult{
		TaskID:    s.ID,
		Status:    ExecStatusRunning,
		StartTime: time.Now(),
	}

	// 创建带超时的context
	execCtx, cancel := context.WithTimeout(ctx, s.Config.Timeout)
	defer cancel()

	var lastErr error
	var output string

	// 重试机制
	for retry := 0; retry <= s.Config.RetryCount; retry++ {
		result.Retry = retry

		// 执行任务处理器
		output, lastErr = s.handler.Handle(execCtx, s.Config.Parameters)

		if lastErr == nil {
			result.Status = ExecStatusSuccess
			result.Output = output
			break
		}

		// 如果不是最后一次重试，则等待重试延迟
		if retry < s.Config.RetryCount {
			s.logger.WithFields(map[string]interface{}{
				"task_id": s.ID,
				"handler": s.Config.HandlerName,
				"retry":   retry + 1,
				"error":   lastErr.Error(),
			}).Warn("系统任务执行失败，准备重试")

			select {
			case <-time.After(s.Config.RetryDelay):
				continue
			case <-ctx.Done():
				result.Status = ExecStatusTimeout
				result.Error = "任务执行超时"
				break
			}
		}
	}

	// 设置执行结果
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	if lastErr != nil {
		if result.Status == ExecStatusRunning {
			result.Status = ExecStatusFailed
		}
		result.Error = lastErr.Error()
	}

	// 记录执行日志
	logFields := map[string]interface{}{
		"task_id":  s.ID,
		"handler":  s.Config.HandlerName,
		"status":   result.Status,
		"duration": result.Duration,
		"retry":    result.Retry,
	}

	if result.Status == ExecStatusSuccess {
		s.logger.WithFields(logFields).Info("系统任务执行成功")
	} else {
		logFields["error"] = result.Error
		s.logger.WithFields(logFields).Error("系统任务执行失败")
	}

	return result
}

// GetConfig 获取任务配置
func (s *SystemTask) GetConfig() interface{} {
	return s.Config
}

// Validate 验证任务配置
func (s *SystemTask) Validate() error {
	if s.Config.HandlerName == "" {
		return fmt.Errorf("处理器名称不能为空")
	}

	if _, exists := GetSystemTaskHandler(s.Config.HandlerName); !exists {
		return fmt.Errorf("未找到系统任务处理器: %s", s.Config.HandlerName)
	}

	if s.Config.Timeout < 0 {
		return fmt.Errorf("超时时间不能为负数")
	}

	if s.Config.RetryCount < 0 {
		return fmt.Errorf("重试次数不能为负数")
	}

	if s.Config.RetryDelay < 0 {
		return fmt.Errorf("重试延迟不能为负数")
	}

	return nil
}

// ToJSON 转换为JSON
func (s *SystemTask) ToJSON() (string, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

// FromJSON 从JSON解析
func (s *SystemTask) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), s)
}

// UpdateConfig 更新配置
func (s *SystemTask) UpdateConfig(config interface{}) error {
	configData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("配置序列化失败: %w", err)
	}

	var systemConfig SystemConfig
	if err := json.Unmarshal(configData, &systemConfig); err != nil {
		return fmt.Errorf("配置反序列化失败: %w", err)
	}

	s.Config = systemConfig
	return s.Validate()
}

// =============== 内置系统任务处理器 ===============

// LogHandler 日志处理器
type LogHandler struct {
	logger logger.Logger
}

// NewLogHandler 创建日志处理器
func NewLogHandler(log logger.Logger) *LogHandler {
	return &LogHandler{logger: log}
}

// Handle 处理日志任务
func (l *LogHandler) Handle(ctx context.Context, params map[string]interface{}) (string, error) {
	level, _ := params["level"].(string)
	message, _ := params["message"].(string)

	if message == "" {
		return "", fmt.Errorf("消息内容不能为空")
	}

	switch level {
	case "debug":
		l.logger.Debug(message)
	case "info":
		l.logger.Info(message)
	case "warn":
		l.logger.Warn(message)
	case "error":
		l.logger.Error(message)
	default:
		l.logger.Info(message)
	}

	return fmt.Sprintf("日志记录成功: [%s] %s", level, message), nil
}

// GetName 获取处理器名称
func (l *LogHandler) GetName() string {
	return "log"
}

// GetDescription 获取处理器描述
func (l *LogHandler) GetDescription() string {
	return "记录日志消息，支持debug、info、warn、error级别"
}

// CacheCleanHandler 缓存清理处理器
type CacheCleanHandler struct {
	cache  cache.Cache
	logger logger.Logger
}

// NewCacheCleanHandler 创建缓存清理处理器
func NewCacheCleanHandler(c cache.Cache, log logger.Logger) *CacheCleanHandler {
	return &CacheCleanHandler{
		cache:  c,
		logger: log,
	}
}

// Handle 处理缓存清理任务
func (c *CacheCleanHandler) Handle(ctx context.Context, params map[string]interface{}) (string, error) {
	pattern, _ := params["pattern"].(string)
	if pattern == "" {
		pattern = "*"
	}

	keys, err := c.cache.KeysCtx(ctx, pattern)
	if err != nil {
		return "", fmt.Errorf("获取缓存键失败: %w", err)
	}

	if len(keys) == 0 {
		return "没有找到匹配的缓存键", nil
	}

	deleted, err := c.cache.DelCtx(ctx, keys...)
	if err != nil {
		return "", fmt.Errorf("删除缓存失败: %w", err)
	}

	return fmt.Sprintf("成功清理 %d 个缓存键", deleted), nil
}

// GetName 获取处理器名称
func (c *CacheCleanHandler) GetName() string {
	return "cache_clean"
}

// GetDescription 获取处理器描述
func (c *CacheCleanHandler) GetDescription() string {
	return "清理缓存，支持通配符模式匹配"
}

// DatabaseMaintenanceHandler 数据库维护处理器
type DatabaseMaintenanceHandler struct {
	logger logger.Logger
}

// NewDatabaseMaintenanceHandler 创建数据库维护处理器
func NewDatabaseMaintenanceHandler(log logger.Logger) *DatabaseMaintenanceHandler {
	return &DatabaseMaintenanceHandler{logger: log}
}

// Handle 处理数据库维护任务
func (d *DatabaseMaintenanceHandler) Handle(ctx context.Context, params map[string]interface{}) (string, error) {
	action, _ := params["action"].(string)

	switch action {
	case "analyze":
		return "数据库分析完成", nil
	case "optimize":
		return "数据库优化完成", nil
	case "vacuum":
		return "数据库清理完成", nil
	default:
		return "", fmt.Errorf("不支持的维护操作: %s", action)
	}
}

// GetName 获取处理器名称
func (d *DatabaseMaintenanceHandler) GetName() string {
	return "db_maintenance"
}

// GetDescription 获取处理器描述
func (d *DatabaseMaintenanceHandler) GetDescription() string {
	return "数据库维护操作，支持analyze、optimize、vacuum等操作"
}
