package infrastructure

import (
	"apple-findmy-to-mqtt/infrastructure/config"
	"apple-findmy-to-mqtt/infrastructure/logging"
	"apple-findmy-to-mqtt/infrastructure/shared"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(config.GetConfig),
	fx.Provide(logging.GetLogger),
	fx.Provide(shared.NewHelpers),
)
