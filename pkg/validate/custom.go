package validate

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 自定义验证规则
func RegisterCustomValidations(v *validator.Validate) {
	// 中国手机号验证
	v.RegisterValidation("phone", validatePhone)
	// 中国身份证验证
	v.RegisterValidation("idcard", validateIDCard)
	// 中文姓名验证
	v.RegisterValidation("chinese_name", validateChineseName)
	// 强密码验证
	v.RegisterValidation("strong_password", validateStrongPassword)
	// 字段相等
	v.RegisterValidation("eqfield", validateEqField)
}

// 验证中国手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return false
	}
	// 中国手机号正则
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}

// 验证中国身份证号
func validateIDCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()
	if idCard == "" {
		return false
	}
	// 18位身份证正则
	idCardRegex := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`)
	return idCardRegex.MatchString(idCard)
}

// 验证中文姓名
func validateChineseName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if name == "" {
		return false
	}
	// 中文姓名正则（2-4个中文字符）
	nameRegex := regexp.MustCompile(`^[\u4e00-\u9fa5]{2,4}$`)
	return nameRegex.MatchString(name)
}

// 验证强密码（至少8位，包含大小写字母、数字和特殊字符）
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	// 检查是否包含大写字母
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// 检查是否包含小写字母
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// 检查是否包含数字
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	// 检查是否包含特殊字符
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// 验证字段相等
func validateEqField(fl validator.FieldLevel) bool {
	field := fl.Field()
	otherFieldName := fl.Param()
	otherField := fl.Parent().FieldByName(otherFieldName)
	return field.String() == otherField.String()
}
