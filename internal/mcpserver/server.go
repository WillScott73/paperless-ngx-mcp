package mcpserver

import (
	"github.com/WillScott73/paperless-ngx-mcp/internal/paperless/client"
	"github.com/WillScott73/paperless-ngx-mcp/internal/paperless/tools"
	"github.com/mark3labs/mcp-go/server"
)

// New creates and configures the MCP server instance.
func New(c *client.Client) *server.MCPServer {
	s := server.NewMCPServer(
		"paperless-ngx-mcp",
		"1.0.0",
	)

	tools.Register(s, c)

	return s
}
