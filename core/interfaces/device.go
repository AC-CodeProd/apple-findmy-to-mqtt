package interfaces

import (
	"apple-findmy-to-mqtt/core/entities"
	"time"
)

type IDeviceUsecase interface {
	GetDevicesCache() ([]entities.Device, error)
	HasDeviceMustBeUpdated(id, name string, lastUpdate time.Time) bool
}
