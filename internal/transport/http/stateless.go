package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

func ServeStateless(s *server.MCPServer, port string) error {
	addr := ":" + port
	if port != "" && port[0] == ':' {
		addr = port
	}

	fmt.Fprintf(os.Stderr, "Starting Paperless MCP server on Stateless HTTP at %s\n", addr)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-mcp-protocol-version")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method == http.MethodGet {
			if r.URL.Path == "/.well-known/oauth-authorization-server" {
				http.NotFound(w, r)
				return
			}
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Paperless MCP Server Operational (Stateless Mode)"))
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var rawMessage json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&rawMessage); err != nil {
			http.Error(w, "Parse error", http.StatusBadRequest)
			return
		}

		result := s.HandleMessage(r.Context(), rawMessage)
		if result == nil {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	})

	return http.ListenAndServe(addr, handler)
}
