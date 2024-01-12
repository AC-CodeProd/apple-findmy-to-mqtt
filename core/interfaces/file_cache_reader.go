package interfaces

import (
	"apple-findmy-to-mqtt/core/entities"
	"time"
)

type IFileCacheReader interface {
	CalcAccuracy(horizontalAccuracy, verticalAccuracy float64) float64
	ConvertToDevice(data any) entities.Device
	GetSourceType(applePositionType string) string
	HasDeviceMustBeUpdated(id, name string, lastUpdate time.Time) bool
	ReadDevicesData() ([]entities.Device, error)
}
