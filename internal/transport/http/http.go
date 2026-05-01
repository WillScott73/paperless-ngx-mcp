package http

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

// Serve runs the MCP server on HTTP using Server-Sent Events (SSE).
func Serve(s *server.MCPServer, host, port, baseURL string) error {
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Fprintf(os.Stderr, "Starting Paperless-ngx MCP server on HTTP at %s\n", addr)
	
	sseServer := server.NewSSEServer(s, server.WithBaseURL(baseURL))
	return sseServer.Start(addr)
}
