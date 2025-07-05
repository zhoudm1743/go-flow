package unified

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GinContext 是Gin上下文的适配器
type GinContext struct {
	ctx *gin.Context
}

// NewGinContext 创建一个Gin上下文适配器
func NewGinContext(ctx *gin.Context) *GinContext {
	return &GinContext{
		ctx: ctx,
	}
}

// Method 实现Context接口
func (c *GinContext) Method() string {
	return c.ctx.Request.Method
}

// Path 实现Context接口
func (c *GinContext) Path() string {
	return c.ctx.FullPath()
}

// Host 实现Context接口
func (c *GinContext) Host() string {
	return c.ctx.Request.Host
}

// URL 实现Context接口
func (c *GinContext) URL() *url.URL {
	return c.ctx.Request.URL
}

// ClientIP 实现Context接口
func (c *GinContext) ClientIP() string {
	return c.ctx.ClientIP()
}

// GetHeader 实现Context接口
func (c *GinContext) GetHeader(key string) string {
	return c.ctx.GetHeader(key)
}

// SetHeader 实现Context接口
func (c *GinContext) SetHeader(key, value string) {
	c.ctx.Header(key, value)
}

// Query 实现Context接口
func (c *GinContext) Query(key string) string {
	return c.ctx.Query(key)
}

// QueryDefault 实现Context接口
func (c *GinContext) QueryDefault(key, defaultValue string) string {
	return c.ctx.DefaultQuery(key, defaultValue)
}

// QueryMap 实现Context接口
func (c *GinContext) QueryMap() map[string]string {
	result := make(map[string]string)
	for k, v := range c.ctx.Request.URL.Query() {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}

// Param 实现Context接口
func (c *GinContext) Param(key string) string {
	return c.ctx.Param(key)
}

// ParamInt 实现Context接口
func (c *GinContext) ParamInt(key string) (int, error) {
	return strconv.Atoi(c.ctx.Param(key))
}

// ParamUint 实现Context接口
func (c *GinContext) ParamUint(key string) (uint, error) {
	val, err := strconv.ParseUint(c.ctx.Param(key), 10, 64)
	return uint(val), err
}

// BindJSON 实现Context接口
func (c *GinContext) BindJSON(obj interface{}) error {
	return c.ctx.ShouldBindJSON(obj)
}

// BindQuery 实现Context接口
func (c *GinContext) BindQuery(obj interface{}) error {
	return c.ctx.ShouldBindQuery(obj)
}

// BindForm 实现Context接口
func (c *GinContext) BindForm(obj interface{}) error {
	return c.ctx.ShouldBind(obj)
}

// FormFile 实现Context接口
func (c *GinContext) FormFile(name string) (*multipart.FileHeader, error) {
	return c.ctx.FormFile(name)
}

// FormValue 实现Context接口
func (c *GinContext) FormValue(name string) string {
	return c.ctx.PostForm(name)
}

// Status 实现Context接口
func (c *GinContext) Status(code int) Context {
	c.ctx.Status(code)
	return c
}

// JSON 实现Context接口
func (c *GinContext) JSON(code int, obj interface{}) error {
	c.ctx.JSON(code, obj)
	return nil
}

// String 实现Context接口
func (c *GinContext) String(code int, format string, values ...interface{}) error {
	c.ctx.String(code, format, values...)
	return nil
}

// HTML 实现Context接口
func (c *GinContext) HTML(code int, html string) error {
	c.ctx.Data(code, "text/html", []byte(html))
	return nil
}

// Redirect 实现Context接口
func (c *GinContext) Redirect(code int, url string) error {
	c.ctx.Redirect(code, url)
	return nil
}

// File 实现Context接口
func (c *GinContext) File(filepath string) error {
	c.ctx.File(filepath)
	return nil
}

// Stream 实现Context接口
func (c *GinContext) Stream(contentType string, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	c.ctx.Data(http.StatusOK, contentType, data)
	return nil
}

// Set 实现Context接口
func (c *GinContext) Set(key string, value interface{}) {
	c.ctx.Set(key, value)
}

// Get 实现Context接口
func (c *GinContext) Get(key string) (interface{}, bool) {
	return c.ctx.Get(key)
}

// MustGet 实现Context接口
func (c *GinContext) MustGet(key string) interface{} {
	return c.ctx.MustGet(key)
}

// GinContext 实现Context接口
func (c *GinContext) GinContext() interface{} {
	return c.ctx
}

// FiberContext 实现Context接口
func (c *GinContext) FiberContext() interface{} {
	return nil
}

// GetRequest 实现Context接口
func (c *GinContext) GetRequest() *http.Request {
	return c.ctx.Request
}

// GetResponse 实现Context接口
func (c *GinContext) GetResponse() http.ResponseWriter {
	return c.ctx.Writer
}

// Error 实现Context接口
func (c *GinContext) Error(err error) error {
	c.ctx.Error(err)
	return err
}

// HasErrors 实现Context接口
func (c *GinContext) HasErrors() bool {
	return len(c.ctx.Errors) > 0
}

// Errors 实现Context接口
func (c *GinContext) Errors() []error {
	result := make([]error, len(c.ctx.Errors))
	for i, err := range c.ctx.Errors {
		result[i] = err.Err
	}
	return result
}

// Next 实现Context接口
func (c *GinContext) Next() {
	c.ctx.Next()
}

// IsAborted 实现Context接口
func (c *GinContext) IsAborted() bool {
	return c.ctx.IsAborted()
}

// Abort 实现Context接口
func (c *GinContext) Abort() {
	c.ctx.Abort()
}

// AbortWithStatus 实现Context接口
func (c *GinContext) AbortWithStatus(code int) {
	c.ctx.AbortWithStatus(code)
}

// AbortWithJSON 实现Context接口
func (c *GinContext) AbortWithJSON(code int, obj interface{}) error {
	c.ctx.AbortWithStatusJSON(code, obj)
	return nil
}
