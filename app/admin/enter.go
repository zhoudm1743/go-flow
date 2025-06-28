package admin

import (
	"github.com/zhoudm1743/go-flow/app/admin/routes"
	"github.com/zhoudm1743/go-flow/app/admin/service"
	"go.uber.org/fx"
)

var Module = fx.Module("admin",
	routes.Module,
	service.Module,
)
