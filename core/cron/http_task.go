package cron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/zhoudm1743/go-flow/core/logger"
)

// HTTPTask HTTP任务实现
type HTTPTask struct {
	BaseTask
	Config HTTPConfig `json:"config"`
	logger logger.Logger
}

// NewHTTPTask 创建HTTP任务
func NewHTTPTask(base BaseTask, config HTTPConfig, log logger.Logger) *HTTPTask {
	task := &HTTPTask{
		BaseTask: base,
		Config:   config,
		logger:   log,
	}

	// 设置默认值
	if task.Config.Method == "" {
		task.Config.Method = "GET"
	}
	if task.Config.Timeout == 0 {
		task.Config.Timeout = 30 * time.Second
	}
	if task.Config.RetryCount == 0 {
		task.Config.RetryCount = 3
	}
	if task.Config.RetryDelay == 0 {
		task.Config.RetryDelay = 5 * time.Second
	}
	if task.Config.ExpectedCode == 0 {
		task.Config.ExpectedCode = 200
	}

	return task
}

// Execute 执行HTTP任务
func (h *HTTPTask) Execute(ctx context.Context) *ExecResult {
	result := &ExecResult{
		TaskID:    h.ID,
		Status:    ExecStatusRunning,
		StartTime: time.Now(),
	}

	// 创建带超时的context
	execCtx, cancel := context.WithTimeout(ctx, h.Config.Timeout)
	defer cancel()

	var lastErr error
	var response *http.Response

	// 重试机制
	for retry := 0; retry <= h.Config.RetryCount; retry++ {
		result.Retry = retry

		// 执行HTTP请求
		response, lastErr = h.executeHTTPRequest(execCtx)

		if lastErr == nil {
			// 检查状态码
			if response.StatusCode == h.Config.ExpectedCode {
				result.Status = ExecStatusSuccess
				break
			} else {
				lastErr = fmt.Errorf("意外的状态码: %d, 期望: %d", response.StatusCode, h.Config.ExpectedCode)
			}
		}

		// 如果不是最后一次重试，则等待重试延迟
		if retry < h.Config.RetryCount {
			h.logger.WithFields(map[string]interface{}{
				"task_id": h.ID,
				"retry":   retry + 1,
				"error":   lastErr.Error(),
			}).Warn("HTTP任务执行失败，准备重试")

			select {
			case <-time.After(h.Config.RetryDelay):
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

	// 读取响应内容
	if response != nil {
		if body, err := ioutil.ReadAll(response.Body); err == nil {
			result.Output = string(body)
			response.Body.Close()
		}
	}

	// 记录执行日志
	logFields := map[string]interface{}{
		"task_id":  h.ID,
		"status":   result.Status,
		"duration": result.Duration,
		"retry":    result.Retry,
	}

	if result.Status == ExecStatusSuccess {
		h.logger.WithFields(logFields).Info("HTTP任务执行成功")
	} else {
		logFields["error"] = result.Error
		h.logger.WithFields(logFields).Error("HTTP任务执行失败")
	}

	return result
}

// executeHTTPRequest 执行HTTP请求
func (h *HTTPTask) executeHTTPRequest(ctx context.Context) (*http.Response, error) {
	// 创建请求
	var body *bytes.Buffer
	if h.Config.Body != "" {
		body = bytes.NewBufferString(h.Config.Body)
	}

	req, err := http.NewRequestWithContext(ctx, h.Config.Method, h.Config.URL, body)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	for key, value := range h.Config.Headers {
		req.Header.Set(key, value)
	}

	// 如果没有设置Content-Type且有Body，则设置默认的Content-Type
	if h.Config.Body != "" && req.Header.Get("Content-Type") == "" {
		if h.isJSON(h.Config.Body) {
			req.Header.Set("Content-Type", "application/json")
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: h.Config.Timeout,
	}

	// 执行请求
	return client.Do(req)
}

// isJSON 检查字符串是否为JSON格式
func (h *HTTPTask) isJSON(str string) bool {
	str = strings.TrimSpace(str)
	return (strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}")) ||
		(strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]"))
}

// GetConfig 获取任务配置
func (h *HTTPTask) GetConfig() interface{} {
	return h.Config
}

// Validate 验证任务配置
func (h *HTTPTask) Validate() error {
	if h.Config.URL == "" {
		return fmt.Errorf("URL不能为空")
	}

	if !strings.HasPrefix(h.Config.URL, "http://") && !strings.HasPrefix(h.Config.URL, "https://") {
		return fmt.Errorf("URL必须以http://或https://开头")
	}

	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	method := strings.ToUpper(h.Config.Method)
	valid := false
	for _, validMethod := range validMethods {
		if method == validMethod {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("不支持的HTTP方法: %s", h.Config.Method)
	}

	if h.Config.Timeout < 0 {
		return fmt.Errorf("超时时间不能为负数")
	}

	if h.Config.RetryCount < 0 {
		return fmt.Errorf("重试次数不能为负数")
	}

	if h.Config.RetryDelay < 0 {
		return fmt.Errorf("重试延迟不能为负数")
	}

	if h.Config.ExpectedCode < 100 || h.Config.ExpectedCode > 599 {
		return fmt.Errorf("期望状态码必须在100-599之间")
	}

	return nil
}

// ToJSON 转换为JSON
func (h *HTTPTask) ToJSON() (string, error) {
	data, err := json.Marshal(h)
	return string(data), err
}

// FromJSON 从JSON解析
func (h *HTTPTask) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), h)
}

// UpdateConfig 更新配置
func (h *HTTPTask) UpdateConfig(config interface{}) error {
	configData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("配置序列化失败: %w", err)
	}

	var httpConfig HTTPConfig
	if err := json.Unmarshal(configData, &httpConfig); err != nil {
		return fmt.Errorf("配置反序列化失败: %w", err)
	}

	h.Config = httpConfig
	return h.Validate()
}
