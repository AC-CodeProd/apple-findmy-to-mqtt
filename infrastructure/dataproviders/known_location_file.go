package dataproviders

import (
	"apple-findmy-to-mqtt/core/entities"
	"apple-findmy-to-mqtt/core/interfaces"
	"apple-findmy-to-mqtt/infrastructure/config"
	"apple-findmy-to-mqtt/infrastructure/logging"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/fx"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Tolerance float64 `json:"tolerance"`
}

type LocationMap map[string]Location

type KnownLocationFileParams struct {
	fx.In
	Config config.Config
	Logger logging.Logger
}
type knownLocationFile struct {
	config    config.Config
	logger    logging.Logger
	locations entities.KnownLocationMap
}

func NewKnownLocationFile(klfp KnownLocationFileParams) interfaces.IKnownLocationFile {
	klf := &knownLocationFile{
		logger: klfp.Logger,
		config: klfp.Config,
	}
	locations, _ := klf.LoadLocationsFromFile(klfp.Config.KnownLocationsPath)
	klf.locations = locations
	return klf
}

func (klfp *knownLocationFile) LoadLocationsFromFile(filePath string) (entities.KnownLocationMap, error) {
	var locations LocationMap

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &locations)
	if err != nil {
		return nil, err
	}

	return klfp.ConvertToKnownLocationMap(locations)
}

func (klf *knownLocationFile) ConvertToKnownLocationMap(data any) (entities.KnownLocationMap, error) {
	locationMap := data.(LocationMap)
	locationMap, ok := data.(LocationMap)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected LocationMap")
	}
	knownLocationMap := make(entities.KnownLocationMap)
	for key, location := range locationMap {
		knownLocationMap[key] = entities.KnownLocation{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
			Tolerance: location.Tolerance,
		}
	}

	return knownLocationMap, nil
}

func (klf *knownLocationFile) GetAllLocations() entities.KnownLocationMap {
	return klf.locations
}
