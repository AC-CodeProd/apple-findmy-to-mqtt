package core

import (
	usecases "apple-findmy-to-mqtt/core/usecases"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(usecases.NewDeviceUsecase),
	fx.Provide(usecases.NewKnownLocationsUsecase),
)
