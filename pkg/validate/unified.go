package validate

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/zhoudm1743/go-flow/pkg/http/unified"
	"github.com/zhoudm1743/go-flow/pkg/response"
)

var UnifiedVerify = unifiedVerifyUtil{}

// unifiedVerifyUtil 统一参数验证工具类
type unifiedVerifyUtil struct{}

// VerifyJSON 验证JSON参数并返回中文错误信息
func (vu unifiedVerifyUtil) VerifyJSON(c unified.Context, obj any) (e error) {
	if err := c.BindJSON(obj); err != nil {
		// 如果是验证错误，返回中文翻译
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := TranslateError(validationErrors)
			e = response.ParamsValidError.MakeData(errors)
		} else {
			e = response.ParamsValidError.MakeData(err.Error())
		}
		return
	}
	return
}

// VerifyJSONWithValidator 使用自定义验证器验证JSON参数
func (vu unifiedVerifyUtil) VerifyJSONWithValidator(c unified.Context, obj any) (e error) {
	// 获取请求体
	req := c.GetRequest()
	if req == nil {
		// Fiber不支持直接获取request，需要从上下文获取body
		if fctx := c.FiberContext(); fctx != nil {
			if fc, ok := fctx.(*fiber.Ctx); ok {
				body := fc.Body()
				if err := json.Unmarshal(body, &obj); err != nil {
					e = response.ParamsValidError.MakeData(err.Error())
					return
				}
			}
		} else {
			e = response.ParamsValidError.MakeData("不支持的上下文类型")
			return
		}
	} else {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			e = response.ParamsValidError.MakeData(err.Error())
			return
		}

		if err := json.Unmarshal(body, &obj); err != nil {
			e = response.ParamsValidError.MakeData(err.Error())
			return
		}
	}

	// 使用validator进行验证
	if err := Validator.Struct(obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := TranslateError(validationErrors)
			e = response.ParamsValidError.MakeData(errors)
		} else {
			e = response.ParamsValidError.MakeData(err.Error())
		}
		return
	}

	return
}

// VerifyQuery 验证Query参数并返回中文错误信息
func (vu unifiedVerifyUtil) VerifyQuery(c unified.Context, obj any) (e error) {
	if err := c.BindQuery(obj); err != nil {
		// 如果是验证错误，返回中文翻译
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := TranslateError(validationErrors)
			e = response.ParamsValidError.MakeData(errors)
		} else {
			e = response.ParamsValidError.MakeData(err.Error())
		}
		return
	}
	return
}

// VerifyForm 验证Form参数
func (vu unifiedVerifyUtil) VerifyForm(c unified.Context, obj any) (e error) {
	if err := c.BindForm(obj); err != nil {
		// 如果是验证错误，返回中文翻译
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := TranslateError(validationErrors)
			e = response.ParamsValidError.MakeData(errors)
		} else {
			e = response.ParamsValidError.MakeData(err.Error())
		}
		return
	}
	return
}

// ValidateStruct 直接验证结构体
func (vu unifiedVerifyUtil) ValidateStruct(obj any) error {
	if Validator == nil {
		if err := InitValidator(); err != nil {
			return err
		}
	}

	if err := Validator.Struct(obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := TranslateError(validationErrors)
			return response.ParamsValidError.MakeData(errors)
		}
		return response.ParamsValidError.MakeData(err.Error())
	}
	return nil
}

// Verify 验证参数 使用validator进行验证
func (vu unifiedVerifyUtil) Verify(c unified.Context, obj ...any) (e error) {
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
			case strings.HasPrefix(contentType, "multipart/form-data") || contentType == "application/x-www-form-urlencoded":
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
