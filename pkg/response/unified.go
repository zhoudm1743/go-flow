package response

import (
	"net/http"

	"github.com/zhoudm1743/go-flow/pkg/http/unified"
)

// UnifiedResponse 统一响应结构体
type UnifiedResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

// UnifiedResult 统一响应
func UnifiedResult(c unified.Context, resp RespType, data interface{}) error {
	if data == nil {
		data = resp.data
	}

	if resp != Success {
		c.Error(resp)
	}

	return c.JSON(http.StatusOK, UnifiedResponse{
		Code: resp.code,
		Msg:  resp.msg,
		Data: data,
	})
}

// UnifiedOk 正常响应
func UnifiedOk(c unified.Context) error {
	return UnifiedResult(c, Success, []string{})
}

// UnifiedOkWithMsg 正常响应附带msg
func UnifiedOkWithMsg(c unified.Context, msg string) error {
	resp := Success
	resp.msg = msg
	return UnifiedResult(c, resp, []string{})
}

// UnifiedOkWithData 正常响应附带data
func UnifiedOkWithData(c unified.Context, data interface{}) error {
	return UnifiedResult(c, Success, data)
}

// UnifiedFail 错误响应
func UnifiedFail(c unified.Context, resp RespType) error {
	c.Status(http.StatusBadRequest)
	return UnifiedResult(c, resp, []string{})
}

// UnifiedFailWithMsg 错误响应附带msg
func UnifiedFailWithMsg(c unified.Context, resp RespType, msg string) error {
	resp.msg = msg
	c.Status(http.StatusBadRequest)
	return UnifiedResult(c, resp, []string{})
}

// UnifiedFailWithData 错误响应附带data
func UnifiedFailWithData(c unified.Context, resp RespType, data interface{}) error {
	c.Status(http.StatusBadRequest)
	return UnifiedResult(c, resp, data)
}

// UnifiedIsFailWithResp 判断是否出现错误，并追加错误返回信息
func UnifiedIsFailWithResp(c unified.Context, err error) bool {
	if err == nil {
		return false
	}
	switch v := err.(type) {
	// 自定义类型
	case RespType:
		data := v.Data()
		if data == nil {
			data = []string{}
		}
		UnifiedFailWithData(c, v, data)
	// 其他类型
	default:
		UnifiedFailWithMsg(c, SystemError, err.Error())
	}
	return true
}

// UnifiedCheckAndResp 判断是否出现错误，并返回对应响应
func UnifiedCheckAndResp(c unified.Context, err error) error {
	if UnifiedIsFailWithResp(c, err) {
		return nil
	}
	return UnifiedOk(c)
}

// UnifiedCheckAndRespWithData 判断是否出现错误，并返回对应响应（带data数据）
func UnifiedCheckAndRespWithData(c unified.Context, data interface{}, err error) error {
	if UnifiedIsFailWithResp(c, err) {
		return nil
	}
	return UnifiedOkWithData(c, data)
}
