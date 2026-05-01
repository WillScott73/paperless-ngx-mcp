package main

import (
	"fmt"
	"os"

	"github.com/WillScott73/paperless-ngx-mcp/internal/config"
	"github.com/WillScott73/paperless-ngx-mcp/internal/mcpserver"
	"github.com/WillScott73/paperless-ngx-mcp/internal/paperless/client"
	"github.com/WillScott73/paperless-ngx-mcp/internal/transport/http"
	"github.com/WillScott73/paperless-ngx-mcp/internal/transport/stdio"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration Error: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Validation Error: %v\n", err)
		os.Exit(1)
	}

	// Initialize API client
	pClient := client.New(cfg.BaseURL, cfg.APIToken)

	// Create MCP server
	srv := mcpserver.New(pClient)

	// Start server on requested transport
	switch cfg.Transport {
	case "stdio":
		if err := stdio.Serve(srv); err != nil {
			fmt.Fprintf(os.Stderr, "Stdio server error: %v\n", err)
			os.Exit(1)
		}
	case "stateless":
		if err := http.ServeStateless(srv, cfg.Port); err != nil {
			fmt.Fprintf(os.Stderr, "Stateless server error: %v\\n", err)
			os.Exit(1)
		}
	case "http":
		if err := http.Serve(srv, cfg.Host, cfg.Port, cfg.MCPBaseURL); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP server error: %v\n", err)
			os.Exit(1)
		}
	}
}
