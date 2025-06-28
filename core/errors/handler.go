package errors

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-flow/core/logger"
)

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Type      ErrorType              `json:"type"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Timestamp int64                  `json:"timestamp"`
	Path      string                 `json:"path"`
}

// ErrorHandler 错误处理器
type ErrorHandler struct {
	logger logger.Logger
	debug  bool
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler(log logger.Logger, debug bool) *ErrorHandler {
	return &ErrorHandler{
		logger: log,
		debug:  debug,
	}
}

// Middleware 错误处理中间件
func (eh *ErrorHandler) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if len(c.Errors) > 0 {
				eh.handleErrors(c)
			}
		}()
		c.Next()
	}
}

// HandleError 处理单个错误
func (eh *ErrorHandler) HandleError(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
		eh.handleErrors(c)
		c.Abort()
	}
}

// handleErrors 处理所有错误
func (eh *ErrorHandler) handleErrors(c *gin.Context) {
	if len(c.Errors) == 0 {
		return
	}

	// 取最后一个错误作为主要错误
	lastError := c.Errors.Last()
	err := lastError.Err

	var appErr AppError
	var ok bool

	// 检查是否为AppError
	if appErr, ok = err.(AppError); !ok {
		// 不是AppError，包装为系统错误
		appErr = WrapError(err, InternalServerError, ErrorTypeSystem, "Internal server error")
	}

	// 创建错误响应
	response := ErrorResponse{
		Code:      appErr.Code(),
		Message:   appErr.Message(),
		Type:      appErr.Type(),
		Details:   appErr.Details(),
		RequestID: eh.getRequestID(c),
		Timestamp: time.Now().Unix(),
		Path:      c.Request.URL.Path,
	}

	// 记录错误日志
	eh.logError(c, appErr)

	// 如果是debug模式，添加堆栈信息
	if eh.debug {
		response.Details["stack"] = appErr.Stack()
		response.Details["all_errors"] = c.Errors.Errors()
	}

	// 设置HTTP状态码并返回JSON响应
	httpStatus := appErr.HttpStatus()
	c.JSON(httpStatus, response)
}

// logError 记录错误日志
func (eh *ErrorHandler) logError(c *gin.Context, err AppError) {
	fields := map[string]interface{}{
		"error_code":  err.Code(),
		"error_type":  err.Type(),
		"request_id":  eh.getRequestID(c),
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"user_agent":  c.Request.UserAgent(),
		"remote_addr": c.ClientIP(),
	}

	// 添加错误详情
	for k, v := range err.Details() {
		fields["detail_"+k] = v
	}

	switch err.Type() {
	case ErrorTypeValidation, ErrorTypeBusiness:
		eh.logger.WithFields(fields).Warn(err.Message())
	case ErrorTypePermission:
		eh.logger.WithFields(fields).Warn(fmt.Sprintf("Permission denied: %s", err.Message()))
	case ErrorTypeSystem, ErrorTypeExternal:
		eh.logger.WithFields(fields).Error(err.Error())
		if eh.debug {
			eh.logger.WithFields(fields).Error(err.Stack())
		}
	default:
		eh.logger.WithFields(fields).Error(err.Error())
	}
}

// getRequestID 获取请求ID
func (eh *ErrorHandler) getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// ========== 便捷的响应方法 ==========

// SendSuccess 成功响应
func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":      SuccessCode,
		"message":   "操作成功",
		"data":      data,
		"timestamp": time.Now().Unix(),
	})
}

// SendSuccessWithMessage 带消息的成功响应
func SendSuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":      SuccessCode,
		"message":   message,
		"data":      data,
		"timestamp": time.Now().Unix(),
	})
}

// SendPagedSuccess 分页成功响应
func SendPagedSuccess(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, gin.H{
		"code":    SuccessCode,
		"message": "操作成功",
		"data": gin.H{
			"items":     data,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
		"timestamp": time.Now().Unix(),
	})
}

// SendError 错误响应
func SendError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

// SendValidationError 验证错误响应
func SendValidationError(c *gin.Context, message string, details map[string]interface{}) {
	err := NewValidationError(message)
	if details != nil {
		err = err.WithDetails(details)
	}
	SendError(c, err)
}

// SendBusinessError 业务错误响应
func SendBusinessError(c *gin.Context, message string) {
	SendError(c, NewBusinessError(message))
}

// SendPermissionError 权限错误响应
func SendPermissionError(c *gin.Context, message string) {
	SendError(c, NewPermissionError(message))
}

// SendNotFoundError 资源未找到错误响应
func SendNotFoundError(c *gin.Context, resource string) {
	SendError(c, NewNotFoundError(resource))
}

// SendSystemError 系统错误响应
func SendSystemError(c *gin.Context, message string) {
	SendError(c, NewSystemError(message))
}

// ========== 增强的ControllerBase ==========

// ControllerBase 控制器基类 - 集成统一错误处理
type ControllerBase struct{}

// Success 成功响应
func (cb *ControllerBase) Success(c *gin.Context, data interface{}) {
	SendSuccess(c, data)
}

// SuccessWithMessage 带消息的成功响应
func (cb *ControllerBase) SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	SendSuccessWithMessage(c, message, data)
}

// PagedSuccess 分页成功响应
func (cb *ControllerBase) PagedSuccess(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	SendPagedSuccess(c, data, total, page, pageSize)
}

// Error 通用错误响应
func (cb *ControllerBase) Error(c *gin.Context, err error) {
	SendError(c, err)
}

// ValidationError 验证错误
func (cb *ControllerBase) ValidationError(c *gin.Context, message string, details map[string]interface{}) {
	SendValidationError(c, message, details)
}

// BusinessError 业务错误
func (cb *ControllerBase) BusinessError(c *gin.Context, message string) {
	SendBusinessError(c, message)
}

// PermissionError 权限错误
func (cb *ControllerBase) PermissionError(c *gin.Context, message string) {
	SendPermissionError(c, message)
}

// NotFoundError 资源未找到错误
func (cb *ControllerBase) NotFoundError(c *gin.Context, resource string) {
	SendNotFoundError(c, resource)
}

// SystemError 系统错误
func (cb *ControllerBase) SystemError(c *gin.Context, message string) {
	SendSystemError(c, message)
}

// WithValidation 参数验证辅助方法
func (cb *ControllerBase) WithValidation(c *gin.Context, req interface{}, handler func()) {
	if err := c.ShouldBindJSON(req); err != nil {
		cb.ValidationError(c, "参数验证失败", map[string]interface{}{
			"validation_error": err.Error(),
		})
		return
	}
	handler()
}

// WithAuth 权限验证辅助方法
func (cb *ControllerBase) WithAuth(c *gin.Context, requiredRole string, handler func()) {
	userRole, exists := c.Get("user_role")
	if !exists {
		cb.PermissionError(c, "用户未认证")
		return
	}

	if userRole != requiredRole && userRole != "admin" {
		cb.PermissionError(c, "权限不足")
		return
	}

	handler()
}

// ========== Context增强 ==========

// WithContext 为Context添加增强功能
func WithContext(c *gin.Context) *EnhancedContext {
	return &EnhancedContext{Context: c}
}

// EnhancedContext 增强的Context
type EnhancedContext struct {
	*gin.Context
}

// MustBind 必须绑定成功，否则返回验证错误
func (ec *EnhancedContext) MustBind(req interface{}) bool {
	if err := ec.ShouldBindJSON(req); err != nil {
		SendValidationError(ec.Context, "参数验证失败", map[string]interface{}{
			"error": err.Error(),
		})
		return false
	}
	return true
}

// GetUserID 获取用户ID
func (ec *EnhancedContext) GetUserID() (string, bool) {
	userID, exists := ec.Get("user_id")
	if !exists {
		return "", false
	}
	if id, ok := userID.(string); ok {
		return id, true
	}
	return "", false
}

// MustGetUserID 必须获取用户ID，否则返回权限错误
func (ec *EnhancedContext) MustGetUserID() (string, bool) {
	userID, exists := ec.GetUserID()
	if !exists {
		SendPermissionError(ec.Context, "用户未认证")
		return "", false
	}
	return userID, true
}
