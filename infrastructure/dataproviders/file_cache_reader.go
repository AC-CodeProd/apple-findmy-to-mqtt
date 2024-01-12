package dataproviders

import (
	"apple-findmy-to-mqtt/core/entities"
	"apple-findmy-to-mqtt/core/interfaces"
	"apple-findmy-to-mqtt/infrastructure/config"
	"apple-findmy-to-mqtt/infrastructure/logging"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/fx"
)

type FindMyData struct {
	Address       FindMyDataAddress  `json:"address"`
	BatteryStatus string             `json:"batteryStatus"`
	DeviceClass   string             `json:"deviceClass"`
	Location      FindMyDataLocation `json:"location"`
	ModelName     string             `json:"modelDisplayName"`
	Name          string             `json:"name"`
}

type FindMyDataAddress struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Locality    string `json:"locality"`
	FullAddress string `json:"mapItemFullAddress"`
}

type FindMyDataLocation struct {
	HorizontalAccuracy float64 `json:"horizontalAccuracy"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	PositionType       string  `json:"positionType"`
	TimeStamp          int64   `json:"timeStamp"`
	VerticalAccuracy   float64 `json:"verticalAccuracy"`
}

type FindMyDevice struct {
	FindMyData
}

type deviceUpdateInfo struct {
	LastUpdateTime time.Time
}

var deviceUpdates = make(map[string]deviceUpdateInfo)

type FileCacheReaderParams struct {
	fx.In
	Config config.Config
	Logger logging.Logger
}

type fileCacheReader struct {
	logger logging.Logger
}

func NewFileCacheReader(fcrp FileCacheReaderParams) interfaces.IFileCacheReader {
	return &fileCacheReader{
		logger: fcrp.Logger,
	}
}

func (fcr *fileCacheReader) CalcAccuracy(horizontalAccuracy, verticalAccuracy float64) float64 {
	return math.Sqrt(math.Pow(horizontalAccuracy, 2) + math.Pow(verticalAccuracy, 2))
}

func (fcr *fileCacheReader) ConvertToDevice(data any) entities.Device {
	findMyDevice := data.(FindMyDevice)
	timestamp := findMyDevice.Location.TimeStamp
	lastUpdate := time.Unix(timestamp/1000, (timestamp%1000)*1000000)
	sourceType := fcr.GetSourceType(findMyDevice.Location.PositionType)
	gpsAccuracy := fcr.CalcAccuracy(findMyDevice.Location.HorizontalAccuracy, findMyDevice.Location.VerticalAccuracy)

	return *entities.NewDevice(
		findMyDevice.Address.FullAddress,
		findMyDevice.BatteryStatus,
		gpsAccuracy,
		lastUpdate,
		findMyDevice.Location.Latitude,
		findMyDevice.Location.Longitude,
		findMyDevice.ModelName,
		findMyDevice.Name,
		sourceType,
	)
}

func (fcr *fileCacheReader) GetSourceType(applePositionType string) string {
	switch applePositionType {
	case "crowdsourced", "safeLocation":
		return "gps"
	case "Wifi":
		return "router"
	default:
		return "gps"
	}
}

func (fcr *fileCacheReader) HasDeviceMustBeUpdated(id, name string, lastUpdate time.Time) bool {
	updatesIdentifier := fmt.Sprintf("%s (%s)", name, id)
	lastUpdateInfo, exists := deviceUpdates[updatesIdentifier]
	updated := exists && lastUpdateInfo.LastUpdateTime.Equal(lastUpdate)

	deviceUpdates[updatesIdentifier] = deviceUpdateInfo{
		LastUpdateTime: lastUpdate,
	}
	return updated
}

func (fcr *fileCacheReader) ReadDevicesData() ([]entities.Device, error) {
	const names = "__file_cache_reader.go__: ReadDevicesData"
	usr, err := user.Current()

	if err != nil {
		return nil, fmt.Errorf("%s | error getting current user: %w", names, err)
	}

	filePathDevices := filepath.Join(usr.HomeDir, "Library/Caches/com.apple.findmy.fmipcore/Devices.data")
	filePathItems := filepath.Join(usr.HomeDir, "Library/Caches/com.apple.findmy.fmipcore/Items.data")

	var wg sync.WaitGroup
	wg.Add(2)

	var devicesData, itemsData []FindMyDevice
	var errDevices, errItems error

	go func() {
		defer wg.Done()
		devicesData, errDevices = readAndUnmarshalData(filePathDevices)
		if errDevices != nil {
			fcr.logger.Warn(fmt.Sprintf("%s | %s", names, errDevices.Error()))
		}
	}()
	go func() {
		defer wg.Done()
		itemsData, errItems = readAndUnmarshalData(filePathItems)
		if errItems != nil {
			fcr.logger.Warn(fmt.Sprintf("%s | %s", names, errItems.Error()))
		}
	}()

	wg.Wait()

	data := append(devicesData, itemsData...)
	devices := make([]entities.Device, len(data))
	for i, findMyDevice := range data {
		devices[i] = fcr.ConvertToDevice(findMyDevice)
	}
	return devices, nil
}

func readData(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func readAndUnmarshalData(filePath string) ([]FindMyDevice, error) {
	data, err := readData(filePath)
	if err != nil {
		return nil, err
	}
	var findMyDevices []FindMyDevice
	if err := json.Unmarshal(data, &findMyDevices); err != nil {
		return nil, err
	}

	return findMyDevices, nil
}
