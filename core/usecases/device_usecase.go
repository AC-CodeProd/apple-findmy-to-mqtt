package usecases

import (
	"apple-findmy-to-mqtt/core/entities"
	"apple-findmy-to-mqtt/core/interfaces"
	"time"
)

type deviceUsecase struct {
	fileCacheReader interfaces.IFileCacheReader
}

func NewDeviceUsecase(fileCacheReader interfaces.IFileCacheReader) interfaces.IDeviceUsecase {
	return &deviceUsecase{
		fileCacheReader: fileCacheReader,
	}
}

func (du *deviceUsecase) GetDevicesCache() ([]entities.Device, error) {
	return du.fileCacheReader.ReadDevicesData()
}

func (du *deviceUsecase) HasDeviceMustBeUpdated(id, name string, lastUpdate time.Time) bool {
	return du.fileCacheReader.HasDeviceMustBeUpdated(id, name, lastUpdate)
}
