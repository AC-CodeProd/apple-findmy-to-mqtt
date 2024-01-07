package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	envPath      string
	globalConfig *Config
	ENV_DEFAULT  = map[string]any{
		"DEBUG":                             true,
		"ENVIRONMENT":                       "development",
		"GO_ENV":                            "development",
		"KNOWN_LOCATIONS_DEFAULT_TOLERANCE": 70,
		"KNOWN_LOCATIONS_PATH":              "known_locations.json",
		"LOG_LEVEL":                         "info",
		"LOG_OUTPUT":                        "./logs/development.log",
		"MQTT_PORT":                         1883,
		"MQTT_CLIENT_ID":                    "apple_findmy_to_mqtt",
		"SCAN_TIMER":                        5,
		"TZ":                                "Europe/Paris",
	}
)

type Config struct {
	Environment                    string         `json:"environment"`
	ForceSync                      bool           `json:"force_sync"`
	KnownLocationsDefaultTolerance int            `json:"known_locations_default_tolerance"`
	KnownLocationsPath             string         `json:"known_locations_path"`
	Loggers                        []LoggerConfig `json:"loggers"`
	LogLevel                       string         `json:"log_level"`
	LogOutput                      string         `json:"log_output"`
	Mqtt                           Mqtt           `json:"mqtt"`
	ScanTimer                      int            `json:"scan_timer"`
	TZ                             string         `json:"tz"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type AliasConfig Config
	alias := &struct {
		ScanTimer                      string `json:"scan_timer"`
		ForceSync                      string `json:"force_sync"`
		KnownLocationsDefaultTolerance string `json:"known_locations_default_tolerance"`
		*AliasConfig
	}{
		AliasConfig: (*AliasConfig)(c),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	if alias.ScanTimer != "" {
		val := getEnvValue(strings.ToUpper(alias.ScanTimer))
		scanTimer, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return err
		}
		c.ScanTimer = int(scanTimer)
	}
	if alias.KnownLocationsDefaultTolerance != "" {
		val := getEnvValue(strings.ToUpper(alias.KnownLocationsDefaultTolerance))
		knownLocationsDefaultTolerance, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return err
		}
		c.KnownLocationsDefaultTolerance = int(knownLocationsDefaultTolerance)
	}
	if alias.ForceSync != "" {
		val := getEnvValue(strings.ToUpper(alias.ForceSync))
		boolValue, err := strconv.ParseBool(strings.Trim(val, "\""))
		if err != nil {
			return err
		}
		c.ForceSync = boolValue
	}

	return nil
}

type LoggerConfig struct {
	Directory    string        `json:"directory"`
	LayoutFormat string        `json:"layout_format"`
	Lef          string        `json:"lef"`
	Path         string        `json:"path"`
	Ropt         RotateOptions `json:"ropt"`
	Type         string        `json:"type"`
}
type Mqtt struct {
	Broker    string `json:"broker"`
	ClientID  string `json:"client_id"`
	HassTopic string `json:"hass_topic"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
	Topic     string `json:"topic"`
	Username  string `json:"username"`
}

func (m *Mqtt) UnmarshalJSON(data []byte) error {
	type AliasMqtt Mqtt
	alias := &struct {
		Port string `json:"port"`
		*AliasMqtt
	}{
		AliasMqtt: (*AliasMqtt)(m),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	if alias.Port != "" {
		val := getEnvValue(strings.ToUpper(alias.Port))
		port, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return err
		}
		m.Port = int(port)
	}

	return nil
}

type RotateOptions struct {
	Compress   bool `json:"compress"`
	MaxAge     int  `json:"max_age"`
	MaxBackups int  `json:"max_backups"`
	MaxSize    int  `json:"max_size"`
}

func SetupConfig(_envPath string) error {
	if len(_envPath) == 0 {
		return errors.New(`pass the path list`)
	}
	envPath = _envPath
	return nil
}

func newConfig(config *Config) error {
	currentDirectory, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	_dotEnvPath := filepath.Join(currentDirectory, envPath)
	if err := godotenv.Load(_dotEnvPath); err != nil {
		panic(err.Error())
	}

	file, err := os.ReadFile("config.json")
	if err != nil {
		panic(err.Error())
	}
	if err := json.Unmarshal(file, &config); err != nil {
		panic(err.Error())
	}

	replaceWithEnv(config)

	return nil
}

func GetConfig() Config {
	const names = "__config.go__ : GetConfig"
	if globalConfig == nil {
		globalConfig = &Config{}
		if err := newConfig(globalConfig); err != nil {
			panic(fmt.Sprintf("%s | %s", names, err))
		}
	}
	return *globalConfig
}

func replaceWithEnv(config *Config) {
	replaceField(reflect.ValueOf(config).Elem())
}

func replaceField(v reflect.Value) {
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			tag := t.Field(i).Tag.Get("json")
			if tag != "" && field.Kind() != reflect.Struct {
				replaceValue(field)
			} else if field.Kind() == reflect.Struct {
				replaceField(field)
			}
		}
	}
}

func getEnvValue(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if value, ok := ENV_DEFAULT[key]; ok {
		switch val := value.(type) {
		case string:
			return val
		case int:
			return strconv.Itoa(val)
		case float64:
			return strconv.FormatFloat(val, 'f', -1, 64)
		default:
			return fmt.Sprintf("%v", value)
		}
	}

	return key
}

func replaceValue(v reflect.Value) {
	if !v.IsValid() || !v.CanSet() {
		return
	}
	k := strings.ToUpper(v.String())
	val := getEnvValue(k)
	switch v.Kind() {
	case reflect.String:
		v.SetString(strings.Trim(val, "\""))
	case reflect.Int:
		if i, err := strconv.Atoi(strings.Trim(val, "\"")); err == nil {
			v.SetInt(int64(i))
		}
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(strings.Trim(val, "\"")); err == nil {
			v.SetBool(boolValue)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			newVal := reflect.New(val.Type()).Elem()
			newVal.Set(val)
			replaceValue(newVal)
			v.SetMapIndex(key, newVal)
		}
	}
}
