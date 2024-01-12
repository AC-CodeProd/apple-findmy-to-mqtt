package usecases

import (
	"apple-findmy-to-mqtt/core/entities"
	"apple-findmy-to-mqtt/core/interfaces"
	"math"
)

type knownLocationsUsecase struct {
	knownLocationFile interfaces.IKnownLocationFile
}

func NewKnownLocationsUsecase(knownLocationFile interfaces.IKnownLocationFile) interfaces.IKnownLocationsUsecase {
	return &knownLocationsUsecase{
		knownLocationFile: knownLocationFile,
	}
}

func (kluc *knownLocationsUsecase) GetLocationName(knownLocation entities.KnownLocation) string {
	knownLocations := kluc.knownLocationFile.GetAllLocations()
	for name, location := range knownLocations {
		tolerance := float64(location.Tolerance)
		if tolerance == 0 {
			tolerance = knownLocation.Tolerance
		}
		tolerance = getLatLngApprox(tolerance)
		if isClose(location.Latitude, knownLocation.Latitude, tolerance) && isClose(location.Longitude, knownLocation.Longitude, tolerance) {
			return name
		}
	}
	return "not_home"
}

func getLatLngApprox(meters float64) float64 {
	return meters / 111111
}

// isClose checks whether two float values are close within a certain tolerance.
func isClose(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
