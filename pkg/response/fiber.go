package response

import (
	"github.com/gofiber/fiber/v2"
)

// FiberResponse 定义Fiber响应结构体
type FiberResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// FiberResult 返回Fiber响应结果
func FiberResult(c *fiber.Ctx, resp RespType, data interface{}) error {
	if data == nil {
		data = resp.data
	}

	return c.JSON(FiberResponse{
		Code: resp.code,
		Msg:  resp.msg,
		Data: data,
	})
}

// FiberOk Fiber正常响应
func FiberOk(c *fiber.Ctx) error {
	return FiberResult(c, Success, []string{})
}

// FiberOkWithMsg Fiber正常响应附带msg
func FiberOkWithMsg(c *fiber.Ctx, msg string) error {
	resp := Success
	resp.msg = msg
	return FiberResult(c, resp, []string{})
}

// FiberOkWithData Fiber正常响应附带data
func FiberOkWithData(c *fiber.Ctx, data interface{}) error {
	return FiberResult(c, Success, data)
}

// FiberFail Fiber错误响应
func FiberFail(c *fiber.Ctx, resp RespType) error {
	c.Status(fiber.StatusBadRequest)
	return FiberResult(c, resp, []string{})
}

// FiberFailWithMsg Fiber错误响应附带msg
func FiberFailWithMsg(c *fiber.Ctx, resp RespType, msg string) error {
	resp.msg = msg
	c.Status(fiber.StatusBadRequest)
	return FiberResult(c, resp, []string{})
}

// FiberFailWithData Fiber错误响应附带data
func FiberFailWithData(c *fiber.Ctx, resp RespType, data interface{}) error {
	c.Status(fiber.StatusBadRequest)
	return FiberResult(c, resp, data)
}

// FiberCheckAndResp Fiber检查错误并响应
func FiberCheckAndResp(c *fiber.Ctx, err error) error {
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return FiberResult(c, SystemError, []string{})
	}
	return FiberOk(c)
}

// FiberCheckAndRespWithData Fiber检查错误并响应数据
func FiberCheckAndRespWithData(c *fiber.Ctx, data interface{}, err error) error {
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return FiberResult(c, SystemError, []string{})
	}
	return FiberOkWithData(c, data)
}
