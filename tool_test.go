package main_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WillScott73/paperless-ngx-mcp/internal/paperless/client"
	"github.com/WillScott73/paperless-ngx-mcp/internal/paperless/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func TestAllTools(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	c := client.New(ts.URL, "dummy-token")
	s := server.NewMCPServer("test", "1.0.0")
	tools.Register(s, c)

	ctx := context.Background()

	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{"paperless_health_check", nil},
		{"paperless_get_statistics", nil},
		{"paperless_search_documents", map[string]interface{}{"query": "test", "page": 1.0, "page_size": 25.0}},
		{"paperless_get_document", map[string]interface{}{"id": 123.0}},
		{"paperless_list_tags", map[string]interface{}{"page": 1.0, "page_size": 25.0}},
		{"paperless_list_correspondents", map[string]interface{}{"page": 1.0, "page_size": 25.0}},
		{"paperless_list_document_types", map[string]interface{}{"page": 1.0, "page_size": 25.0}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := mcp.CallToolRequest{}
			req.Params.Name = tc.name
			if tc.args != nil {
				req.Params.Arguments = tc.args
			}

			// We have to invoke the tool callback directly, but mark3labs server doesn't expose it easily.
			// Let's just find the tool in the server and call it.
			// Actually, mcp-go's server doesn't expose the router.
			// Let's use the provided handlers if possible, or just realize that the user wants to test
			// what is WRONG with the tool they used earlier.
		})
	}
}
