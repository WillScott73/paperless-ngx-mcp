package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sort"

	"github.com/WillScott73/paperless-ngx-mcp/internal/mcpserver"
	"github.com/WillScott73/paperless-ngx-mcp/internal/paperless/client"
)

func main() {
	baseURL := flag.String("url", "https://mcp.theblueonion.app/paperless", "Base URL for the API")
	title := flag.String("title", "Paperless MCP", "Title for the API")
	flag.Parse()

	c := client.New("https://paperless.local", "token")
	s := mcpserver.New(c)

	tools := s.ListTools()
	
	api := map[string]any{
		"openapi": "3.0.0",
		"info": map[string]any{
			"title":   *title,
			"version": "1.0.0",
		},
		"servers": []map[string]any{
			{"url": *baseURL},
		},
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]any{
					"type":         "http",
					"scheme":       "bearer",
				},
			},
		},
		"security": []map[string]any{
			{"bearerAuth": []string{}},
		},
		"paths": make(map[string]any),
	}

	paths := api["paths"].(map[string]any)
	var toolNames []string
	for name := range tools {
		toolNames = append(toolNames, name)
	}
	sort.Strings(toolNames)

	for _, name := range toolNames {
		t := tools[name].Tool
		paths["/"+name] = map[string]any{
			"post": map[string]any{
				"operationId": name,
				"summary":     t.Description,
				"requestBody": map[string]any{
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"jsonrpc": map[string]any{"type": "string", "example": "2.0"},
									"id":      map[string]any{"type": "string", "example": "1"},
									"method":  map[string]any{"type": "string", "example": "tools/call"},
									"params":  map[string]any{
										"type": "object",
										"properties": map[string]any{
											"name":      map[string]any{"type": "string", "example": name},
											"arguments": t.InputSchema,
										},
										"required": []string{"name", "arguments"},
									},
								},
								"required": []string{"jsonrpc", "id", "method", "params"},
							},
						},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Successful response",
					},
				},
			},
		}
	}

	out, _ := json.MarshalIndent(api, "", "  ")
	fmt.Println(string(out))
}
