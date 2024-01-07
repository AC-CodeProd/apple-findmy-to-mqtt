package bootstrap

import (
	"apple-findmy-to-mqtt/controllers"
	"apple-findmy-to-mqtt/core"
	"apple-findmy-to-mqtt/infrastructure"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	controllers.Module,
	core.Module,
	infrastructure.Module,
)
