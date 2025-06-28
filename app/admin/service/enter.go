package service

import (
	"github.com/zhoudm1743/go-flow/app/admin/service/system"
	"github.com/zhoudm1743/go-flow/app/admin/service/test"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(test.NewTestService),
	fx.Provide(system.NewAdminService),
)
