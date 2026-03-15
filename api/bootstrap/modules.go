package bootstrap

import (
	"github.com/kickplate/api/handler/handlers"
	"github.com/kickplate/api/handler/routes"
	"github.com/kickplate/api/lib"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
	routes.Module,
	handlers.Module,
)
