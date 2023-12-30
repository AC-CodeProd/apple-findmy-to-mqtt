package bootstrap

import (
	"apple-findmy-to-mqtt/infrastructure"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	infrastructure.Module,
)
