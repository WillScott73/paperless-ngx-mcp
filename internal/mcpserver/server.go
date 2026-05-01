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
		server.WithTitle("Paperless-ngx"),
		server.WithDescription("MCP server for Paperless-ngx document search, retrieval, tagging, and document management."),
		server.WithInstructions(`Use this server for Paperless-ngx document workflows.
- Start with health, statistics, or search before get, upload, delete, or bulk-edit operations.
- Treat delete and bulk-edit tools as high-impact actions and confirm the target set before using them.
- Keep uploads and metadata edits explicit so the document archive stays organized.
- Prefer list and get tools to identify the right documents, tags, correspondents, or document types before making changes.`),
	)

	tools.Register(s, c)

	return s
}
