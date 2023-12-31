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
		"DEBUG":       true,
		"ENVIRONMENT": "development",
		"GO_ENV":      "development",
		"LOG_LEVEL":   "info",
		"LOG_OUTPUT":  "./logs/development.log",
		"SCAN_TIMER":  15,
		"TZ":          "Europe/Paris",
	}
)

type Config struct {
	Environment string         `json:"environment"`
	Loggers     []LoggerConfig `json:"loggers"`
	LogLevel    string         `json:"log_level"`
	LogOutput   string         `json:"log_output"`
	ScanTimer   int            `json:"scan_timer"`
	TZ          string         `json:"tz"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type AliasConfig Config
	alias := &struct {
		ScanTimer string `json:"scan_timer"`
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
