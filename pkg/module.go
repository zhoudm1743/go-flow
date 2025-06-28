package pkg

import (
	"go.uber.org/fx"

	"github.com/zhoudm1743/go-flow/pkg/captcha"
	"github.com/zhoudm1743/go-flow/pkg/types"
	"github.com/zhoudm1743/go-flow/pkg/validate"
)

// Module pkg模块 - 整合所有pkg下的子模块
var Module = fx.Options(
	// 敏感数据类型模块
	types.Module,

	// 验证器模块
	validate.Module,

	// 验证码模块
	captcha.Module,
)
