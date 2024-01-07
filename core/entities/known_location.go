package entities

type KnownLocation struct {
	Latitude  float64
	Longitude float64
	Tolerance float64
}

type KnownLocationMap map[string]KnownLocation
