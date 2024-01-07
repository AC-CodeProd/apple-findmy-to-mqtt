package shared

import (
	"os"
	"path/filepath"
	"regexp"

	"go.uber.org/fx"
)

type IHelpers interface {
	FindNamedMatches(regex *regexp.Regexp, str string) map[string]string
}

type helpers struct {
}

type HelpersParams struct {
	fx.In
}

func NewHelpers(shp HelpersParams) IHelpers {
	return &helpers{}
}

func (h *helpers) FindNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)
	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}

func CreateDir(filePath string) error {
	dirPath := filepath.Dir(filePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
