package config

import (
	"NotRedis/internal/network"
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	testConfig := `
engine:
  type: "in_memory"
network:
  address: "127.0.0.1:9999"
  max_connections: 50
  max_message_size: "2KB"
  idle_timeout: "10m"
logging:
  level: "debug"
  output: "/tmp/test.log"
`
	tmpFile, err := os.CreateTemp("", "config_test.yaml")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(testConfig)); err != nil {
		t.Fatal("Failed to write to temp file:", err)
	}
	tmpFile.Close()

	tests := []struct {
		name    string
		file    string
		want    Config
		wantErr bool
	}{
		{
			name: "Load valid config",
			file: tmpFile.Name(),
			want: Config{
				Engine: struct {
					Type string `yaml:"type"`
				}{Type: "in_memory"},
				Network: struct {
					Address        string `yaml:"address"`
					MaxConnections int    `yaml:"max_connections"`
					MaxMessageSize string `yaml:"max_message_size"`
					IdleTimeout    string `yaml:"idle_timeout"`
				}{
					Address:        "127.0.0.1:9999",
					MaxConnections: 50,
					MaxMessageSize: "2KB",
					IdleTimeout:    "10m",
				},
				Logging: struct {
					Level  string `yaml:"level"`
					Output string `yaml:"output"`
				}{
					Level:  "debug",
					Output: "/tmp/test.log",
				},
			},
			wantErr: false,
		},
		{
			name: "Load with non-existent file (defaults)",
			file: "nonexistent.yaml",
			want: Config{
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
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.Engine.Type != tt.want.Engine.Type ||
				got.Network.Address != tt.want.Network.Address ||
				got.Network.MaxConnections != tt.want.Network.MaxConnections ||
				got.Network.MaxMessageSize != tt.want.Network.MaxMessageSize ||
				got.Network.IdleTimeout != tt.want.Network.IdleTimeout ||
				got.Logging.Level != tt.want.Logging.Level ||
				got.Logging.Output != tt.want.Logging.Output {
				t.Errorf("LoadConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestToServerConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want network.Config
	}{
		{
			name: "Default sizes",
			cfg: Config{
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
			},
			want: network.Config{
				Address:        "127.0.0.1:3223",
				MaxConnections: 100,
				MaxMessageSize: 4096,
				IdleTimeout:    5 * time.Minute,
			},
		},
		{
			name: "Custom sizes",
			cfg: Config{
				Network: struct {
					Address        string `yaml:"address"`
					MaxConnections int    `yaml:"max_connections"`
					MaxMessageSize string `yaml:"max_message_size"`
					IdleTimeout    string `yaml:"idle_timeout"`
				}{
					Address:        "127.0.0.1:9999",
					MaxConnections: 50,
					MaxMessageSize: "2KB",
					IdleTimeout:    "10m",
				},
			},
			want: network.Config{
				Address:        "127.0.0.1:9999",
				MaxConnections: 50,
				MaxMessageSize: 2048,
				IdleTimeout:    10 * time.Minute,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToServerConfig(tt.cfg)
			if got.Address != tt.want.Address ||
				got.MaxConnections != tt.want.MaxConnections ||
				got.MaxMessageSize != tt.want.MaxMessageSize ||
				got.IdleTimeout != tt.want.IdleTimeout {
				t.Errorf("ToServerConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
