package validate

import (
	"go.uber.org/fx"
)

// Module validate模块
var Module = fx.Options(
	// 初始化验证器
	fx.Invoke(InitValidator),
)
