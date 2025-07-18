package config

import (
	"path"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/makhkets/7.17.25/pkg/utils"
)

type Config struct {
	App    App    `json:"app"`
	Filter Filter `json:"filter"`
}

type App struct {
	Port    int    `json:"port"`
	Address string `json:"address"`
}

type Filter struct {
	MaxFiles          int      `json:"max_files"`
	NotAllowedExtensions []string `json:"not_allowed_extensions"`
}

var config Config

func MustLoad(filename string) *Config {
	configFile := path.Join(utils.FindDirectoryName("config"), filename)

	if err := cleanenv.ReadConfig(configFile, &config); err != nil {
		panic(err)
	}

	return &config
}
