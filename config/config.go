package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	ServicePort int `yaml:"servicePort"`
}

var (
	once     sync.Once
	instance *AppConfig
)

func loadConfig() *AppConfig {
	data, err := os.ReadFile("app.yaml")
	if err != nil {
		panic("Failed to read config file: " + err.Error())
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic("Failed to unmarshal YAML: " + err.Error())
	}

	return &cfg
}

func GetAppConfig() *AppConfig {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}
