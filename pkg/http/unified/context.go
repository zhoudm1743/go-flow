package unified

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// Context 统一的HTTP上下文接口
type Context interface {
	// 基础信息
	Method() string
	Path() string
	Host() string
	URL() *url.URL
	ClientIP() string
	GetHeader(key string) string
	SetHeader(key, value string)

	// 请求相关
	Query(key string) string
	QueryDefault(key, defaultValue string) string
	QueryMap() map[string]string
	Param(key string) string
	ParamInt(key string) (int, error)
	ParamUint(key string) (uint, error)

	// 请求体解析
	BindJSON(obj interface{}) error
	BindQuery(obj interface{}) error
	BindForm(obj interface{}) error
	FormFile(name string) (*multipart.FileHeader, error)
	FormValue(name string) string

	// 响应相关
	Status(code int) Context
	JSON(code int, obj interface{}) error
	String(code int, format string, values ...interface{}) error
	HTML(code int, html string) error
	Redirect(code int, url string) error
	File(filepath string) error
	Stream(contentType string, r io.Reader) error

	// 上下文数据
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	MustGet(key string) interface{}

	// 框架原生上下文
	GinContext() interface{}
	FiberContext() interface{}
	GetRequest() *http.Request
	GetResponse() http.ResponseWriter

	// 错误处理
	Error(err error) error
	HasErrors() bool
	Errors() []error

	// 其他
	Next()
	IsAborted() bool
	Abort()
	AbortWithStatus(code int)
	AbortWithJSON(code int, obj interface{}) error
}
