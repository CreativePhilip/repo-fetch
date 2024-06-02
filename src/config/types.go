package config

import (
	"os"
	"path"
)

type RfetchConfig struct {
	Provider       string // gitlab github etc
	TokenVar       string
	RootPath       string
	GroupWhitelist []string

	ConfigLocation string
}

func (c RfetchConfig) GetToken() string {
	value, _ := os.LookupEnv(c.TokenVar)

	return value
}

func (c RfetchConfig) GetRootPathAbs() string {
	return path.Join(c.ConfigLocation, c.RootPath)
}
