package main_test

import (
	"testing"

	"github.com/WillScott73/paperless-ngx-mcp/internal/paperless/client"
	"github.com/WillScott73/paperless-ngx-mcp/internal/paperless/tools"
	"github.com/mark3labs/mcp-go/server"
)

func TestToolRegistration(t *testing.T) {
	s := server.NewMCPServer("paperless-ngx-mcp", "1.0.0")
	tools.Register(s, client.New("http://example.com/api/", "dummy-token"))

	got := s.ListTools()
	if len(got) == 0 {
		t.Fatal("expected registered tools, got none")
	}

	want := []string{
		"paperless_health_check",
		"paperless_get_statistics",
		"paperless_search_documents",
		"paperless_get_document",
		"paperless_upload_document",
		"paperless_download_document",
		"paperless_delete_document",
		"paperless_bulk_edit_documents",
		"paperless_list_tags",
		"paperless_list_correspondents",
		"paperless_list_document_types",
		"paperless_list_tasks",
		"paperless_list_logs",
	}

	for _, name := range want {
		if _, ok := got[name]; !ok {
			t.Fatalf("expected tool %q to be registered", name)
		}
	}
}
