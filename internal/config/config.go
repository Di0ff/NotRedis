package config

import (
	"NotRedis/internal/network"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Engine struct {
		Type string `yaml:"type"`
	} `yaml:"engine"`
	Network struct {
		Address        string `yaml:"address"`
		MaxConnections int    `yaml:"max_connections"`
		MaxMessageSize string `yaml:"max_message_size"`
		IdleTimeout    string `yaml:"idle_timeout"`
	} `yaml:"network"`
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

func Load(file string) (Config, error) {
	cfg := Config{
		Engine: struct {
			Type string `yaml:"type"`
		}{Type: "in_memory"},
		Network: struct {
			Address        string `yaml:"address"`
			MaxConnections int    `yaml:"max_connections"`
			MaxMessageSize string `yaml:"max_message_size"`
			IdleTimeout    string `yaml:"idle_timeout"`
		}{
			Address:        "127.0.0.1:3223",
			MaxConnections: 100,
			MaxMessageSize: "4KB",
			IdleTimeout:    "5m",
		},
		Logging: struct {
			Level  string `yaml:"level"`
			Output string `yaml:"output"`
		}{
			Level:  "info",
			Output: "/log/output.log",
		},
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return cfg, nil
	}

	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

func ToServerConfig(cfg Config) network.Config {
	maxMessageSize := 4096
	if cfg.Network.MaxMessageSize == "1KB" {
		maxMessageSize = 1024
	} else if cfg.Network.MaxMessageSize == "2KB" {
		maxMessageSize = 2048
	}

	idleTimeout, _ := time.ParseDuration(cfg.Network.IdleTimeout)
	if idleTimeout == 0 {
		idleTimeout = 5 * time.Minute
	}

	return network.Config{
		Address:        cfg.Network.Address,
		MaxConnections: cfg.Network.MaxConnections,
		MaxMessageSize: maxMessageSize,
		IdleTimeout:    idleTimeout,
	}
}

func SetupLogger(cfg Config) (*zap.Logger, error) {
	var zapCfg zap.Config
	switch cfg.Logging.Level {
	case "debug":
		zapCfg = zap.NewDevelopmentConfig()
	case "info":
		zapCfg = zap.NewProductionConfig()
	default:
		zapCfg = zap.NewProductionConfig()
	}
	zapCfg.OutputPaths = []string{cfg.Logging.Output}
	return zapCfg.Build()
}
