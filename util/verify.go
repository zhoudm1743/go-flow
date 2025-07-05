package util

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zhoudm1743/go-flow/pkg/http/unified"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"github.com/zhoudm1743/go-flow/pkg/validate"
)

var VerifyUtil = verifyUtil{}

// verifyUtil 统一参数验证工具类
type verifyUtil struct{}

// VerifyJSON 验证JSON参数并返回中文错误信息
func (vu verifyUtil) VerifyJSON(c unified.Context, obj any) (e error) {
	// 使用统一上下文绑定JSON
	if err := c.BindJSON(obj); err != nil {
		// 如果是验证错误，返回中文翻译
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := validate.TranslateError(validationErrors)
			e = response.ParamsValidError.MakeData(errors)
		} else {
			e = response.ParamsValidError.MakeData(err.Error())
		}
		return
	}
	return
}

// VerifyJSONWithValidator 使用自定义验证器验证JSON参数
func (vu verifyUtil) VerifyJSONWithValidator(c unified.Context, obj any) (e error) {
	// 读取请求体
	body, err := ioutil.ReadAll(c.GetRequest().Body)
	if err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}

	if err := json.Unmarshal(body, &obj); err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}

	// 使用validator进行验证
	if err := validate.Validator.Struct(obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := validate.TranslateError(validationErrors)
			e = response.ParamsValidError.MakeData(errors)
		} else {
			e = response.ParamsValidError.MakeData(err.Error())
		}
		return
	}

	return
}

// VerifyJSONArray 验证JSON数组
func (vu verifyUtil) VerifyJSONArray(c unified.Context, obj any) (e error) {
	body, err := ioutil.ReadAll(c.GetRequest().Body)
	if err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}
	return
}

// VerifyQuery 验证Query参数并返回中文错误信息
func (vu verifyUtil) VerifyQuery(c unified.Context, obj any) (e error) {
	if err := c.BindQuery(obj); err != nil {
		// 如果是验证错误，返回中文翻译
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := validate.TranslateError(validationErrors)
			e = response.ParamsValidError.MakeData(errors)
		} else {
			e = response.ParamsValidError.MakeData(err.Error())
		}
		return
	}
	return
}

// VerifyForm 验证Form参数
func (vu verifyUtil) VerifyForm(c unified.Context, obj any) (e error) {
	if err := c.BindForm(obj); err != nil {
		// 如果是验证错误，返回中文翻译
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := validate.TranslateError(validationErrors)
			e = response.ParamsValidError.MakeData(errors)
		} else {
			e = response.ParamsValidError.MakeData(err.Error())
		}
		return
	}
	return
}

// VerifyFile 验证文件参数
func (vu verifyUtil) VerifyFile(c unified.Context, name string) (file *multipart.FileHeader, e error) {
	// 统一抽象层暂无直接支持，需要根据具体实现处理
	// 尝试从Gin上下文获取文件
	if ginCtx, ok := c.(*unified.GinContext); ok {
		if ctx, ok := ginCtx.GinContext().(*gin.Context); ok {
			file, err := ctx.FormFile(name)
			if err != nil {
				e = response.ParamsValidError.MakeData(err.Error())
				return nil, e
			}
			return file, nil
		}
	}

	e = response.ParamsValidError.MakeData("不支持的上下文类型")
	return nil, e
}

// ValidateStruct 直接验证结构体
func (vu verifyUtil) ValidateStruct(obj any) error {
	if err := validate.Validator.Struct(obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := validate.TranslateError(validationErrors)
			return response.ParamsValidError.MakeData(errors)
		}
		return response.ParamsValidError.MakeData(err.Error())
	}
	return nil
}

// Verify 验证参数 使用validator进行验证
func (vu verifyUtil) Verify(c unified.Context, obj ...any) (e error) {
	method := c.Method()
	contentType := c.GetHeader("Content-Type")

	// 如果没有传入任何对象，直接返回
	if len(obj) == 0 {
		return nil
	}

	// 验证每个对象，返回第一个错误
	for _, o := range obj {
		var err error

		switch method {
		case "GET", "DELETE":
			// GET和DELETE请求通常使用Query参数
			err = vu.VerifyQuery(c, o)
		case "POST", "PUT", "PATCH":
			// POST、PUT、PATCH请求根据Content-Type选择验证方式
			switch {
			case contentType == "application/json" || contentType == "application/json; charset=utf-8":
				err = vu.VerifyJSON(c, o)
			case strings.HasPrefix(contentType, "multipart/form-data"):
				err = vu.VerifyForm(c, o)
			default:
				// 默认尝试JSON验证
				err = vu.VerifyJSON(c, o)
			}
		default:
			// 其他请求方法默认使用JSON验证
			err = vu.VerifyJSON(c, o)
		}

		// 如果验证失败，立即返回错误
		if err != nil {
			return err
		}
	}

	return nil
}
