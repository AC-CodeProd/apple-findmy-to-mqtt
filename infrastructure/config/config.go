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
)

type Config struct {
	Environment string         `json:"environment"`
	LogOutput   string         `json:"log_output"`
	LogLevel    string         `json:"log_level"`
	Timer       int            `json:"timer"`
	TZ          string         `json:"tz"`
	Loggers     []LoggerConfig `json:"loggers"`
}

type LoggerConfig struct {
	Path         string        `json:"path"`
	Type         string        `json:"type"`
	LayoutFormat string        `json:"layout_format"`
	Directory    string        `json:"directory"`
	Ropt         RotateOptions `json:"ropt"`
	Lef          string        `json:"lef"`
}

type RotateOptions struct {
	MaxSize    int  `json:"max_size"`
	MaxAge     int  `json:"max_age"`
	MaxBackups int  `json:"max_backups"`
	Compress   bool `json:"compress"`
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

func replaceValue(v reflect.Value) {
	if !v.IsValid() || !v.CanSet() {
		return
	}

	switch v.Kind() {
	case reflect.String:
		envKey := strings.ToUpper(v.String())
		if envValue, exists := os.LookupEnv(envKey); exists {
			v.SetString(strings.Trim(envValue, "\""))
		}
	case reflect.Int:
		envKey := strings.ToUpper(v.String())
		if envValue, exists := os.LookupEnv(envKey); exists {
			if intValue, err := strconv.Atoi(strings.Trim(envValue, "\"")); err == nil {
				v.SetInt(int64(intValue))
			}
		}
	case reflect.Bool:
		envKey := strings.ToUpper(v.String())
		if envValue, exists := os.LookupEnv(envKey); exists {
			if boolValue, err := strconv.ParseBool(strings.Trim(envValue, "\"")); err == nil {
				v.SetBool(boolValue)
			}
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
