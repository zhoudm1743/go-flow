package util

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"github.com/zhoudm1743/go-flow/pkg/validate"
)

var VerifyUtil = verifyUtil{}

// verifyUtil 参数验证工具类
type verifyUtil struct{}

// VerifyJSON 验证JSON参数并返回中文错误信息
func (vu verifyUtil) VerifyJSON(c *gin.Context, obj any) (e error) {
	if err := c.ShouldBindBodyWith(obj, binding.JSON); err != nil {
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
func (vu verifyUtil) VerifyJSONWithValidator(c *gin.Context, obj any) (e error) {
	body, err := ioutil.ReadAll(c.Request.Body)
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

func (vu verifyUtil) VerifyJSONArray(c *gin.Context, obj any) (e error) {
	body, err := ioutil.ReadAll(c.Request.Body)
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

// VerifyBody 验证Body参数并返回中文错误信息
func (vu verifyUtil) VerifyBody(c *gin.Context, obj any) (e error) {
	if err := c.ShouldBind(obj); err != nil {
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

// VerifyHeader 验证Header参数并返回中文错误信息
func (vu verifyUtil) VerifyHeader(c *gin.Context, obj any) (e error) {
	if err := c.ShouldBindHeader(obj); err != nil {
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

// VerifyQuery 验证Query参数并返回中文错误信息
func (vu verifyUtil) VerifyQuery(c *gin.Context, obj any) (e error) {
	if err := c.ShouldBindQuery(obj); err != nil {
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

func (vu verifyUtil) VerifyFile(c *gin.Context, name string) (file *multipart.FileHeader, e error) {
	file, err := c.FormFile(name)
	if err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}
	return
}

// VerifyForm 验证Form参数
func (vu verifyUtil) VerifyForm(c *gin.Context, obj any) (e error) {
	formDecoder := form.NewDecoder()
	if err := formDecoder.Decode(obj, c.Request.Form); err != nil {
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

// VerifyPostForm 验证PostForm参数
func (vu verifyUtil) VerifyPostForm(c *gin.Context, obj any) (e error) {
	formDecoder := form.NewDecoder()
	if err := formDecoder.Decode(obj, c.Request.PostForm); err != nil {
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
