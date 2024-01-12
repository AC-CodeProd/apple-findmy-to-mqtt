package entities

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Device struct {
	ID            string
	Name          string
	ModelName     string
	BatteryStatus string
	SourceType    string
	Latitude      float64
	Longitude     float64
	Address       string
	GPSAccuracy   float64
	LastUpdate    time.Time
}

func NewDevice(address, batteryStatus string, gpsAccuracy float64, lastUpdate time.Time, latitude, longitude float64, modelName, name, sourceType string) *Device {
	return &Device{
		Address:       address,
		BatteryStatus: batteryStatus,
		GPSAccuracy:   gpsAccuracy,
		ID:            generateDeviceID(name),
		LastUpdate:    lastUpdate,
		Latitude:      latitude,
		Longitude:     longitude,
		ModelName:     modelName,
		Name:          name,
		SourceType:    sourceType,
	}
}

func generateDeviceID(name string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	name, _, _ = transform.String(t, name)

	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	reg, _ := regexp.Compile("[^a-zA-Z0-9_]+")
	name = reg.ReplaceAllString(name, "")

	name = strings.ToLower(name)

	return name
}
