package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/zhoudm1743/go-frame/pkg/response"
	"go.uber.org/fx"
)

// 全局变量
var (
	// Validator 全局验证器
	Validator *validator.Validate
)

// Module 验证模块
var Module = fx.Options(
	// 初始化验证器
	fx.Invoke(InitValidator),
)

// ConfigurableModule 可配置的验证模块
func ConfigurableModule(opts ...Option) fx.Option {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	return fx.Options(
		fx.Provide(
			NewConfigurableValidator(options),
		),
		fx.Invoke(func(v *validator.Validate) {
			// 设置全局验证器
			Validator = v
		}),
	)
}

// Options 验证器配置选项
type Options struct {
	RegisterFuncs []func(*validator.Validate)
}

// Option 验证器配置函数
type Option func(*Options)

// WithCustomValidations 添加自定义验证
func WithCustomValidations(registerFunc func(*validator.Validate)) Option {
	return func(o *Options) {
		o.RegisterFuncs = append(o.RegisterFuncs, registerFunc)
	}
}

// NewValidator 创建验证器
func NewValidator() *validator.Validate {
	v := validator.New()
	RegisterCustomValidations(v)
	return v
}

// NewConfigurableValidator 创建可配置的验证器
func NewConfigurableValidator(options *Options) func() *validator.Validate {
	return func() *validator.Validate {
		v := validator.New()

		// 注册内置的自定义验证
		RegisterCustomValidations(v)

		// 注册用户自定义的验证
		for _, registerFunc := range options.RegisterFuncs {
			registerFunc(v)
		}

		return v
	}
}

// ValidateStruct 验证结构体
func ValidateStruct(obj interface{}) error {
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
