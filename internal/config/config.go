package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration.
type Config struct {
	BaseURL    string
	APIToken   string
	Transport  string
	Port       string
	Host       string
	MCPBaseURL string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		BaseURL:    os.Getenv("PAPERLESS_BASE_URL"),
		APIToken:   os.Getenv("PAPERLESS_API_TOKEN"),
		Transport:  os.Getenv("PAPERLESS_MCP_TRANSPORT"),
		Port:       os.Getenv("PAPERLESS_MCP_PORT"),
		Host:       os.Getenv("PAPERLESS_MCP_HOST"),
		MCPBaseURL: os.Getenv("PAPERLESS_MCP_BASE_URL"),
	}

	// Backward compatibility with the older env names documented in the README.
	if cfg.Transport == "" {
		cfg.Transport = os.Getenv("TRANSPORT")
	}
	if cfg.Port == "" {
		cfg.Port = os.Getenv("PORT")
	}
	if cfg.Host == "" {
		cfg.Host = os.Getenv("HOST")
	}

	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("PAPERLESS_BASE_URL is required")
	}
	if cfg.APIToken == "" {
		return nil, fmt.Errorf("PAPERLESS_API_TOKEN is required")
	}

	if cfg.Transport == "" {
		cfg.Transport = "stdio"
	}

	if cfg.Transport == "http" || cfg.Transport == "stateless" {
		if cfg.Port == "" {
			cfg.Port = "8080"
		}
		if cfg.Host == "" {
			cfg.Host = "127.0.0.1"
		}
		if cfg.MCPBaseURL == "" {
			cfg.MCPBaseURL = "http://" + cfg.Host + ":" + cfg.Port
		}
	}

	return cfg, nil
}

// Validate checks if the configuration is sound.
func (c *Config) Validate() error {
	switch c.Transport {
	case "stdio", "http", "stateless":
		return nil
	default:
		return fmt.Errorf("invalid transport: %s, must be stdio, http, or stateless", c.Transport)
	}
}
