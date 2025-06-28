package types

import (
	"go.uber.org/fx"
)

// Module types模块
var Module = fx.Options(
	// 初始化敏感数据类型
	fx.Invoke(Init),
)
