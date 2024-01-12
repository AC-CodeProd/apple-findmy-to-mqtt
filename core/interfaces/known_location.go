package interfaces

import "apple-findmy-to-mqtt/core/entities"

type IKnownLocationFile interface {
	LoadLocationsFromFile(filePath string) (entities.KnownLocationMap, error)
	GetAllLocations() entities.KnownLocationMap
}

type IKnownLocationsUsecase interface {
	GetLocationName(knownLocation entities.KnownLocation) string
}
