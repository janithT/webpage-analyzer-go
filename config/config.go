package config

import (
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	LinkTimeoutInMs int `yaml:"timeoutInMilliSec"`
	ServicePort     int `yaml:"servicePort"`
	ThreadCount     int `yaml:"ThreadCount"`
}

var (
	once     sync.Once
	instance *AppConfig
)

// Load conficgs to return
func loadConfig() *AppConfig {
	data, err := os.ReadFile("app.yaml")
	if err != nil {
		panic("Failed to read config file: " + err.Error())
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic("Failed to unmarshal YAML: " + err.Error())
	}

	// Provide a default if not set
	if cfg.LinkTimeoutInMs <= 0 {
		cfg.LinkTimeoutInMs = 3000 // default 3 seconds
	}

	return &cfg
}

// Return app configs.
func GetAppConfig() *AppConfig {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}

// GetLinkTimeout returns timeout as time.Duration
func (c *AppConfig) GetLinkTimeout() time.Duration {
	return time.Duration(c.LinkTimeoutInMs) * time.Millisecond
}
