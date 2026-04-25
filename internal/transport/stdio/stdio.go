package stdio

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

// Serve runs the MCP server on stdio.
func Serve(s *server.MCPServer) error {
	fmt.Fprintln(os.Stderr, "Starting Paperless-ngx MCP server on stdio")
	return server.ServeStdio(s)
}
