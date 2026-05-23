package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	DefaultConfigDir  = ".envoy"
	DefaultConfigFile = "config.json"
)

// Config holds the global CLI configuration.
type Config struct {
	RemoteURL   string            `json:"remote_url"`
	ActiveEnv   string            `json:"active_env"`
	AuthToken   string            `json:"auth_token,omitempty"`
	Aliases     map[string]string `json:"aliases,omitempty"`
}

// Load reads the config from the default location (~/.envoy/config.json).
func Load() (*Config, error) {
	path, err := defaultConfigPath()
	if err != nil {
		return nil, err
	}
	return LoadFrom(path)
}

// LoadFrom reads the config from the given file path.
func LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{Aliases: make(map[string]string)}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Aliases == nil {
		cfg.Aliases = make(map[string]string)
	}
	return &cfg, nil
}

// Save writes the config to the default location.
func (c *Config) Save() error {
	path, err := defaultConfigPath()
	if err != nil {
		return err
	}
	return c.SaveTo(path)
}

// SaveTo writes the config to the given file path.
func (c *Config) SaveTo(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func defaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, DefaultConfigDir, DefaultConfigFile), nil
}
