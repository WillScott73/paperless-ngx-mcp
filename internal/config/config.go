package config

import (
	"errors"
	"fmt"
	"os"
)

// Config holds the application configuration.
type Config struct {
	BaseURL    string
	APIToken   string
	Transport  string // stdio or http
	Port       string // for http transport
	Host       string // for http transport
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		BaseURL:   os.Getenv("PAPERLESS_BASE_URL"),
		APIToken:  os.Getenv("PAPERLESS_API_TOKEN"),
		Transport: os.Getenv("TRANSPORT"),
		Port:      os.Getenv("PORT"),
		Host:      os.Getenv("HOST"),
	}

	if cfg.BaseURL == "" {
		return nil, errors.New("PAPERLESS_BASE_URL is required")
	}
	if cfg.APIToken == "" {
		return nil, errors.New("PAPERLESS_API_TOKEN is required")
	}

	if cfg.Transport == "" {
		cfg.Transport = "stdio" // Default
	}

	if cfg.Transport == "http" {
		if cfg.Port == "" {
			cfg.Port = "8080" // Default port for HTTP
		}
		if cfg.Host == "" {
			cfg.Host = "127.0.0.1" // Default host for HTTP
		}
	}

	return cfg, nil
}

// Validate checks if the configuration is sound.
func (c *Config) Validate() error {
	if c.Transport != "stdio" && c.Transport != "http" {
		return fmt.Errorf("invalid transport: %s, must be stdio or http", c.Transport)
	}
	return nil
}
