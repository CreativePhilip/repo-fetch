package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

const ConfigName = "rfetch.config.yml"

func CreateInitialConfiguration(provider string, tokenVar string, rootPath string) *RfetchConfig {
	return &RfetchConfig{
		Provider: provider,
		TokenVar: tokenVar,
		RootPath: rootPath,
	}
}

func ReadConfig(configPath string) (*RfetchConfig, error) {
	bytes, err := os.ReadFile(configPath)

	if err != nil {
		return nil, err
	}

	var config RfetchConfig

	err = yaml.Unmarshal(bytes, &config)

	if err != nil {
		return nil, err
	}

	dir, _ := path.Split(configPath)

	config.ConfigLocation = dir

	return &config, nil
}

func WalkFsAndLoadConfig(startPath string) (*RfetchConfig, error) {
	startPath = path.Clean(startPath)

	if startPath == "/" {
		return nil, CouldNotFindConfigError{}
	}

	candidateConfigPath := path.Join(startPath, ConfigName)

	if _, err := os.Stat(candidateConfigPath); errors.Is(err, os.ErrNotExist) {
		// Config does not exist, go directory higher and check again

		pathElements := strings.Split(startPath, "/")
		newPath := path.Join(pathElements[:len(pathElements)-1]...)

		return WalkFsAndLoadConfig("/" + newPath)
	}

	return ReadConfig(candidateConfigPath)
}
