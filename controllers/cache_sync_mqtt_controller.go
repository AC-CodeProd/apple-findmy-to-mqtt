package controllers

import (
	"apple-findmy-to-mqtt/core/entities"
	"apple-findmy-to-mqtt/core/interfaces"
	"apple-findmy-to-mqtt/infrastructure/config"
	"apple-findmy-to-mqtt/infrastructure/logging"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/fx"
)

type DeviceConfig struct {
	UniqueID            string `json:"unique_id"`
	StateTopic          string `json:"state_topic"`
	JSONAttributesTopic string `json:"json_attributes_topic"`
	Device              struct {
		Identifiers  string `json:"identifiers"`
		Manufacturer string `json:"manufacturer"`
		Name         string `json:"name"`
		Mdl          string `json:"mdl"`
	} `json:"device"`
	SourceType     string `json:"source_type"`
	PayloadHome    string `json:"payload_home"`
	PayloadNotHome string `json:"payload_not_home"`
}

type DeviceAttributes struct {
	Latitude            float64   `json:"latitude"`
	Longitude           float64   `json:"longitude"`
	GPSAccuracy         float64   `json:"gps_accuracy"`
	Address             string    `json:"address"`
	BatteryStatus       string    `json:"batteryStatus"`
	LastUpdateTimestamp time.Time `json:"last_update_timestamp"`
	LastUpdate          string    `json:"last_update"`
	Provider            string    `json:"provider"`
}

type cacheSyncMQTTController struct {
	config                config.Config
	deviceUsecase         interfaces.IDeviceUsecase
	knownLocationsUsecase interfaces.IKnownLocationsUsecase
	logger                logging.Logger
	mqtt                  interfaces.IMQTTClient
}

type CacheSyncMQTTControllerParams struct {
	fx.In
	Config                config.Config
	DeviceUsecase         interfaces.IDeviceUsecase
	KnownLocationsUsecase interfaces.IKnownLocationsUsecase
	Logger                logging.Logger
	Mqtt                  interfaces.IMQTTClient
}

func NewCacheSyncMQTTController(p CacheSyncMQTTControllerParams) interfaces.ICacheSyncMQTTController {
	return &cacheSyncMQTTController{
		config:                p.Config,
		deviceUsecase:         p.DeviceUsecase,
		knownLocationsUsecase: p.KnownLocationsUsecase,
		logger:                p.Logger,
		mqtt:                  p.Mqtt,
	}
}

func (csmc *cacheSyncMQTTController) Process(forceSync bool) {
	const names = "__cache_sync_mqtt_controller.go__: Process"
	if err := csmc.mqtt.Connect(); err != nil {
		csmc.logger.Error(fmt.Sprintf("%s | %s", names, err.Error()))
		panic(err)
	}
	devices, err := csmc.deviceUsecase.GetDevicesCache()
	if err != nil {
		csmc.logger.Warn(fmt.Sprintf("%s | %s", names, err.Error()))
		return
	}
	csmc.logger.Info(fmt.Sprintf("%s | Processing %d devices", names, len(devices)))
	for _, device := range devices {
		if !forceSync && csmc.deviceUsecase.HasDeviceMustBeUpdated(device.ID, device.Name, device.LastUpdate) {
			continue
		}
		go csmc.processDevice(device)
	}
}

func (csmc *cacheSyncMQTTController) processDevice(device entities.Device) {
	const names = "__cache_sync_mqtt_controller.go__: processDevice"
	topics := []string{csmc.config.Mqtt.Topic, csmc.config.Mqtt.HassTopic}
	for _, topic := range topics {
		locationName := csmc.knownLocationsUsecase.GetLocationName(entities.KnownLocation{
			Latitude:  device.Latitude,
			Longitude: device.Longitude,
			Tolerance: float64(csmc.config.KnownLocationsDefaultTolerance),
		})
		deviceTopic := fmt.Sprintf("%s/%s/", topic, device.ID)
		if configJSON, attributesJSON, err := createDeviceConfigAndAttributes(device, deviceTopic); err != nil {
			csmc.logger.Error(fmt.Sprintf("%s | %s", names, err.Error()))
		} else {
			if err := csmc.mqtt.Publish(deviceTopic+"config", configJSON); err != nil {
				csmc.logger.Error(fmt.Sprintf("%s | %s", names, err.Error()))
			}
			if err := csmc.mqtt.Publish(deviceTopic+"attributes", attributesJSON); err != nil {
				csmc.logger.Error(fmt.Sprintf("%s | %s", names, err.Error()))
			}
			if err := csmc.mqtt.Publish(deviceTopic+"state", []byte(locationName)); err != nil {
				csmc.logger.Error(fmt.Sprintf("%s | %s", names, err.Error()))
			}
		}
	}
}

func createDeviceConfigAndAttributes(device entities.Device, topic string) (configJSON []byte, attributesJSON []byte, err error) {
	deviceHassTopic := fmt.Sprintf("%s/%s/", topic, device.ID)
	deviceConfig := DeviceConfig{
		UniqueID:            device.ID,
		StateTopic:          deviceHassTopic + "state",
		JSONAttributesTopic: deviceHassTopic + "attributes",
		SourceType:          device.SourceType,
		PayloadHome:         "home",
		PayloadNotHome:      "not_home",
	}
	deviceConfig.Device.Identifiers = device.ID
	deviceConfig.Device.Manufacturer = "Apple"
	deviceConfig.Device.Name = device.Name
	deviceConfig.Device.Mdl = device.ModelName

	deviceAttributes := DeviceAttributes{
		Latitude:            device.Latitude,
		Longitude:           device.Longitude,
		GPSAccuracy:         device.GPSAccuracy,
		Address:             device.Address,
		BatteryStatus:       device.BatteryStatus,
		LastUpdateTimestamp: device.LastUpdate,
		LastUpdate:          device.LastUpdate.Format(time.RFC3339),
		Provider:            "Apple FindMy To MQTT",
	}

	configJSON, err = json.Marshal(deviceConfig)
	if err != nil {
		return nil, nil, err
	}

	attributesJSON, err = json.Marshal(deviceAttributes)
	if err != nil {
		return nil, nil, err
	}

	return configJSON, attributesJSON, nil
}
