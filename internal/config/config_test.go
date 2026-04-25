package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set required env vars
	os.Setenv("PAPERLESS_BASE_URL", "http://example.com/api/")
	os.Setenv("PAPERLESS_API_TOKEN", "test-token")
	defer os.Unsetenv("PAPERLESS_BASE_URL")
	defer os.Unsetenv("PAPERLESS_API_TOKEN")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.BaseURL != "http://example.com/api/" {
		t.Errorf("expected BaseURL %q, got %q", "http://example.com/api/", cfg.BaseURL)
	}
	if cfg.APIToken != "test-token" {
		t.Errorf("expected APIToken %q, got %q", "test-token", cfg.APIToken)
	}
	if cfg.Transport != "stdio" {
		t.Errorf("expected default Transport %q, got %q", "stdio", cfg.Transport)
	}
}

func TestValidate(t *testing.T) {
	cfg := &Config{Transport: "invalid"}
	if err := cfg.Validate(); err == nil {
		t.Error("Validate() should have failed for invalid transport")
	}

	cfg.Transport = "stdio"
	if err := cfg.Validate(); err != nil {
		t.Errorf("Validate() failed for stdio: %v", err)
	}

	cfg.Transport = "http"
	if err := cfg.Validate(); err != nil {
		t.Errorf("Validate() failed for http: %v", err)
	}
}
