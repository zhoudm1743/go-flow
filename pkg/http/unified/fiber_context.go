package unified

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// FiberContext 是Fiber上下文的适配器
type FiberContext struct {
	ctx *fiber.Ctx
}

// NewFiberContext 创建一个Fiber上下文适配器
func NewFiberContext(ctx *fiber.Ctx) *FiberContext {
	return &FiberContext{
		ctx: ctx,
	}
}

// Method 实现Context接口
func (c *FiberContext) Method() string {
	return c.ctx.Method()
}

// Path 实现Context接口
func (c *FiberContext) Path() string {
	return c.ctx.Route().Path
}

// Host 实现Context接口
func (c *FiberContext) Host() string {
	return c.ctx.Hostname()
}

// URL 实现Context接口
func (c *FiberContext) URL() *url.URL {
	url, _ := url.Parse(string(c.ctx.Request().URI().FullURI()))
	return url
}

// ClientIP 实现Context接口
func (c *FiberContext) ClientIP() string {
	return c.ctx.IP()
}

// GetHeader 实现Context接口
func (c *FiberContext) GetHeader(key string) string {
	return c.ctx.Get(key)
}

// SetHeader 实现Context接口
func (c *FiberContext) SetHeader(key, value string) {
	c.ctx.Set(key, value)
}

// Query 实现Context接口
func (c *FiberContext) Query(key string) string {
	return c.ctx.Query(key)
}

// QueryDefault 实现Context接口
func (c *FiberContext) QueryDefault(key, defaultValue string) string {
	return c.ctx.Query(key, defaultValue)
}

// QueryMap 实现Context接口
func (c *FiberContext) QueryMap() map[string]string {
	result := make(map[string]string)
	c.ctx.QueryParser(&result)
	return result
}

// Param 实现Context接口
func (c *FiberContext) Param(key string) string {
	return c.ctx.Params(key)
}

// ParamInt 实现Context接口
func (c *FiberContext) ParamInt(key string) (int, error) {
	return c.ctx.ParamsInt(key)
}

// ParamUint 实现Context接口
func (c *FiberContext) ParamUint(key string) (uint, error) {
	val, err := strconv.ParseUint(c.ctx.Params(key), 10, 64)
	return uint(val), err
}

// BindJSON 实现Context接口
func (c *FiberContext) BindJSON(obj interface{}) error {
	return c.ctx.BodyParser(obj)
}

// BindQuery 实现Context接口
func (c *FiberContext) BindQuery(obj interface{}) error {
	return c.ctx.QueryParser(obj)
}

// BindForm 实现Context接口
func (c *FiberContext) BindForm(obj interface{}) error {
	return c.ctx.BodyParser(obj)
}

// FormFile 实现Context接口
func (c *FiberContext) FormFile(name string) (*multipart.FileHeader, error) {
	return c.ctx.FormFile(name)
}

// FormValue 实现Context接口
func (c *FiberContext) FormValue(name string) string {
	return c.ctx.FormValue(name)
}

// Status 实现Context接口
func (c *FiberContext) Status(code int) Context {
	c.ctx.Status(code)
	return c
}

// JSON 实现Context接口
func (c *FiberContext) JSON(code int, obj interface{}) error {
	return c.ctx.Status(code).JSON(obj)
}

// String 实现Context接口
func (c *FiberContext) String(code int, format string, values ...interface{}) error {
	return c.ctx.Status(code).SendString(fmt.Sprintf(format, values...))
}

// HTML 实现Context接口
func (c *FiberContext) HTML(code int, html string) error {
	c.ctx.Set("Content-Type", "text/html")
	return c.ctx.Status(code).SendString(html)
}

// Redirect 实现Context接口
func (c *FiberContext) Redirect(code int, url string) error {
	return c.ctx.Status(code).Redirect(url)
}

// File 实现Context接口
func (c *FiberContext) File(filepath string) error {
	return c.ctx.SendFile(filepath)
}

// Stream 实现Context接口
func (c *FiberContext) Stream(contentType string, r io.Reader) error {
	c.ctx.Set("Content-Type", contentType)
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return c.ctx.Send(data)
}

// Set 实现Context接口
func (c *FiberContext) Set(key string, value interface{}) {
	c.ctx.Locals(key, value)
}

// Get 实现Context接口
func (c *FiberContext) Get(key string) (interface{}, bool) {
	val := c.ctx.Locals(key)
	if val == nil {
		return nil, false
	}
	return val, true
}

// MustGet 实现Context接口
func (c *FiberContext) MustGet(key string) interface{} {
	val := c.ctx.Locals(key)
	if val == nil {
		panic("key " + key + " not found")
	}
	return val
}

// GinContext 实现Context接口
func (c *FiberContext) GinContext() interface{} {
	return nil
}

// FiberContext 实现Context接口
func (c *FiberContext) FiberContext() interface{} {
	return c.ctx
}

// GetRequest 实现Context接口
func (c *FiberContext) GetRequest() *http.Request {
	// Fiber不直接暴露http.Request
	// 这里返回nil，如果需要原始请求可以使用FiberContext()获取原始ctx
	return nil
}

// GetResponse 实现Context接口
func (c *FiberContext) GetResponse() http.ResponseWriter {
	// Fiber不直接暴露http.ResponseWriter
	// 这里返回nil，如果需要原始响应可以使用FiberContext()获取原始ctx
	return nil
}

// 用于存储错误的本地变量名
const errorsKey = "_errors"

// Error 实现Context接口
func (c *FiberContext) Error(err error) error {
	if err == nil {
		return nil
	}

	var errs []error
	if v := c.ctx.Locals(errorsKey); v != nil {
		errs = v.([]error)
	}
	errs = append(errs, err)
	c.ctx.Locals(errorsKey, errs)
	return err
}

// HasErrors 实现Context接口
func (c *FiberContext) HasErrors() bool {
	return len(c.Errors()) > 0
}

// Errors 实现Context接口
func (c *FiberContext) Errors() []error {
	if v := c.ctx.Locals(errorsKey); v != nil {
		return v.([]error)
	}
	return nil
}

// Next 实现Context接口
func (c *FiberContext) Next() {
	if err := c.ctx.Next(); err != nil {
		c.Error(err)
	}
}

// 用于存储中止状态的本地变量名
const abortedKey = "_aborted"

// IsAborted 实现Context接口
func (c *FiberContext) IsAborted() bool {
	return c.ctx.Locals(abortedKey) != nil
}

// Abort 实现Context接口
func (c *FiberContext) Abort() {
	c.ctx.Locals(abortedKey, true)
}

// AbortWithStatus 实现Context接口
func (c *FiberContext) AbortWithStatus(code int) {
	c.Abort()
	c.ctx.Status(code)
}

// AbortWithJSON 实现Context接口
func (c *FiberContext) AbortWithJSON(code int, obj interface{}) error {
	c.Abort()
	return c.ctx.Status(code).JSON(obj)
}
