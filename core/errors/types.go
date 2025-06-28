package errors

import (
	"fmt"
	"net/http"
	"runtime"
)

// ErrorCode 错误码类型
type ErrorCode int

// 定义错误码常量
const (
	// 成功
	SuccessCode ErrorCode = 200

	// 客户端错误 4xx
	BadRequest          ErrorCode = 400
	Unauthorized        ErrorCode = 401
	Forbidden           ErrorCode = 403
	NotFound            ErrorCode = 404
	MethodNotAllowed    ErrorCode = 405
	Conflict            ErrorCode = 409
	UnprocessableEntity ErrorCode = 422
	TooManyRequests     ErrorCode = 429

	// 服务端错误 5xx
	InternalServerError ErrorCode = 500
	NotImplemented      ErrorCode = 501
	BadGateway          ErrorCode = 502
	ServiceUnavailable  ErrorCode = 503
	GatewayTimeout      ErrorCode = 504

	// 业务错误码 1000+
	ValidationError      ErrorCode = 1001
	DuplicateError       ErrorCode = 1002
	DatabaseError        ErrorCode = 1003
	CacheError           ErrorCode = 1004
	ExternalServiceError ErrorCode = 1005
	ConfigurationError   ErrorCode = 1006
	PermissionDenied     ErrorCode = 1007
	ResourceNotFound     ErrorCode = 1008
	BusinessLogicError   ErrorCode = 1009
	ConcurrencyError     ErrorCode = 1010
)

// ErrorType 错误类型
type ErrorType string

const (
	ErrorTypeValidation  ErrorType = "validation"
	ErrorTypeBusiness    ErrorType = "business"
	ErrorTypeSystem      ErrorType = "system"
	ErrorTypeExternal    ErrorType = "external"
	ErrorTypePermission  ErrorType = "permission"
	ErrorTypeConcurrency ErrorType = "concurrency"
)

// AppError 应用错误接口
type AppError interface {
	error
	Code() ErrorCode
	Type() ErrorType
	Message() string
	Details() map[string]interface{}
	Stack() string
	HttpStatus() int
	WithDetails(details map[string]interface{}) AppError
	WithMessage(message string) AppError
	Unwrap() error
}

// BaseError 基础错误实现
type BaseError struct {
	code     ErrorCode
	errType  ErrorType
	message  string
	details  map[string]interface{}
	stack    string
	original error
}

// NewError 创建新错误
func NewError(code ErrorCode, errType ErrorType, message string) AppError {
	return &BaseError{
		code:    code,
		errType: errType,
		message: message,
		details: make(map[string]interface{}),
		stack:   getStack(),
	}
}

// WrapError 包装已有错误
func WrapError(err error, code ErrorCode, errType ErrorType, message string) AppError {
	return &BaseError{
		code:     code,
		errType:  errType,
		message:  message,
		details:  make(map[string]interface{}),
		stack:    getStack(),
		original: err,
	}
}

func (e *BaseError) Error() string {
	if e.original != nil {
		return fmt.Sprintf("%s: %v", e.message, e.original)
	}
	return e.message
}

func (e *BaseError) Code() ErrorCode {
	return e.code
}

func (e *BaseError) Type() ErrorType {
	return e.errType
}

func (e *BaseError) Message() string {
	return e.message
}

func (e *BaseError) Details() map[string]interface{} {
	return e.details
}

func (e *BaseError) Stack() string {
	return e.stack
}

func (e *BaseError) HttpStatus() int {
	switch e.code {
	case SuccessCode:
		return http.StatusOK
	case BadRequest, ValidationError:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden, PermissionDenied:
		return http.StatusForbidden
	case NotFound, ResourceNotFound:
		return http.StatusNotFound
	case MethodNotAllowed:
		return http.StatusMethodNotAllowed
	case Conflict, DuplicateError:
		return http.StatusConflict
	case UnprocessableEntity:
		return http.StatusUnprocessableEntity
	case TooManyRequests:
		return http.StatusTooManyRequests
	case InternalServerError, DatabaseError, CacheError:
		return http.StatusInternalServerError
	case NotImplemented:
		return http.StatusNotImplemented
	case BadGateway, ExternalServiceError:
		return http.StatusBadGateway
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	case GatewayTimeout:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

func (e *BaseError) WithDetails(details map[string]interface{}) AppError {
	newErr := *e
	for k, v := range details {
		newErr.details[k] = v
	}
	return &newErr
}

func (e *BaseError) WithMessage(message string) AppError {
	newErr := *e
	newErr.message = message
	return &newErr
}

func (e *BaseError) Unwrap() error {
	return e.original
}

// 获取调用栈
func getStack() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// 便捷的错误创建函数

// NewValidationError 创建验证错误
func NewValidationError(message string) AppError {
	return NewError(ValidationError, ErrorTypeValidation, message)
}

// NewBusinessError 创建业务错误
func NewBusinessError(message string) AppError {
	return NewError(BusinessLogicError, ErrorTypeBusiness, message)
}

// NewPermissionError 创建权限错误
func NewPermissionError(message string) AppError {
	return NewError(PermissionDenied, ErrorTypePermission, message)
}

// NewNotFoundError 创建资源未找到错误
func NewNotFoundError(resource string) AppError {
	return NewError(ResourceNotFound, ErrorTypeBusiness, fmt.Sprintf("%s not found", resource))
}

// NewDuplicateError 创建重复资源错误
func NewDuplicateError(resource string) AppError {
	return NewError(DuplicateError, ErrorTypeBusiness, fmt.Sprintf("%s already exists", resource))
}

// NewSystemError 创建系统错误
func NewSystemError(message string) AppError {
	return NewError(InternalServerError, ErrorTypeSystem, message)
}

// NewDatabaseError 创建数据库错误
func NewDatabaseError(err error) AppError {
	return WrapError(err, DatabaseError, ErrorTypeSystem, "Database operation failed")
}

// NewCacheError 创建缓存错误
func NewCacheError(err error) AppError {
	return WrapError(err, CacheError, ErrorTypeSystem, "Cache operation failed")
}

// NewExternalServiceError 创建外部服务错误
func NewExternalServiceError(service string, err error) AppError {
	return WrapError(err, ExternalServiceError, ErrorTypeExternal,
		fmt.Sprintf("External service %s failed", service))
}

// NewConcurrencyError 创建并发错误
func NewConcurrencyError(operation string) AppError {
	return NewError(ConcurrencyError, ErrorTypeConcurrency,
		fmt.Sprintf("Concurrency conflict in %s", operation))
}

// IsErrorType 检查错误类型
func IsErrorType(err error, errType ErrorType) bool {
	if appErr, ok := err.(AppError); ok {
		return appErr.Type() == errType
	}
	return false
}

// IsErrorCode 检查错误码
func IsErrorCode(err error, code ErrorCode) bool {
	if appErr, ok := err.(AppError); ok {
		return appErr.Code() == code
	}
	return false
}

// GetErrorCode 获取错误码，如果不是AppError返回InternalServerError
func GetErrorCode(err error) ErrorCode {
	if err == nil {
		return SuccessCode
	}
	if appErr, ok := err.(AppError); ok {
		return appErr.Code()
	}
	return InternalServerError
}

// GetHttpStatus 获取HTTP状态码
func GetHttpStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if appErr, ok := err.(AppError); ok {
		return appErr.HttpStatus()
	}
	return http.StatusInternalServerError
}

// 常用错误实例
var (
	ErrValidationFailed   = NewValidationError("Validation failed")
	ErrUnauthorized       = NewError(Unauthorized, ErrorTypePermission, "Unauthorized")
	ErrForbidden          = NewError(Forbidden, ErrorTypePermission, "Forbidden")
	ErrNotFound           = NewError(NotFound, ErrorTypeBusiness, "Resource not found")
	ErrMethodNotAllowed   = NewError(MethodNotAllowed, ErrorTypeSystem, "Method not allowed")
	ErrTooManyRequests    = NewError(TooManyRequests, ErrorTypeSystem, "Too many requests")
	ErrInternalServer     = NewError(InternalServerError, ErrorTypeSystem, "Internal server error")
	ErrNotImplemented     = NewError(NotImplemented, ErrorTypeSystem, "Not implemented")
	ErrServiceUnavailable = NewError(ServiceUnavailable, ErrorTypeSystem, "Service unavailable")
)
