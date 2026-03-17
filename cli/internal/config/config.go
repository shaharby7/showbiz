package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	configDir       = ".showbiz"
	configFile      = "config.yaml"
	credentialsFile = "credentials.json"

	DefaultAPIURL = "https://api.showbiz.dev"
)

// Config holds the CLI configuration.
type Config struct {
	APIURL string `yaml:"api_url" json:"api_url"`
	Org    string `yaml:"org,omitempty" json:"org,omitempty"`
}

// Credentials holds the stored authentication tokens.
type Credentials struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Dir returns the showbiz config directory path.
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, configDir), nil
}

// Load reads the config file from disk. Returns defaults if the file does not exist.
func Load() (*Config, error) {
	cfg := &Config{
		APIURL: DefaultAPIURL,
	}

	dir, err := Dir()
	if err != nil {
		return cfg, nil
	}

	data, err := os.ReadFile(filepath.Join(dir, configFile))
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.APIURL == "" {
		cfg.APIURL = DefaultAPIURL
	}

	return cfg, nil
}

// Save writes the config to disk.
func (c *Config) Save() error {
	dir, err := Dir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(filepath.Join(dir, configFile), data, 0600)
}

// Get retrieves a config value by key.
func (c *Config) Get(key string) (string, error) {
	switch key {
	case "api_url":
		return c.APIURL, nil
	case "org":
		return c.Org, nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}

// Set sets a config value by key.
func (c *Config) Set(key, value string) error {
	switch key {
	case "api_url":
		c.APIURL = value
	case "org":
		c.Org = value
	default:
		return fmt.Errorf("unknown config key: %s (valid keys: api_url, org)", key)
	}
	return nil
}

// LoadCredentials reads the stored credentials from disk.
func LoadCredentials() (*Credentials, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(dir, credentialsFile))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return &creds, nil
}

// SaveCredentials writes credentials to disk.
func SaveCredentials(creds *Credentials) error {
	dir, err := Dir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	return os.WriteFile(filepath.Join(dir, credentialsFile), data, 0600)
}

// ClearCredentials removes stored credentials.
func ClearCredentials() error {
	dir, err := Dir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, credentialsFile)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove credentials: %w", err)
	}
	return nil
}

// ResolveAPIURL returns the API URL from the environment variable, config, or default.
func ResolveAPIURL(cfg *Config) string {
	if env := os.Getenv("SHOWBIZ_API_URL"); env != "" {
		return env
	}
	if cfg != nil && cfg.APIURL != "" {
		return cfg.APIURL
	}
	return DefaultAPIURL
}
