package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/WillScott73/paperless-ngx-mcp/internal/paperless/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Register(s *server.MCPServer, c *client.Client) {
	registerHealthCheck(s, c)
	registerGetStatistics(s, c)
	registerSearchDocuments(s, c)
	registerGetDocument(s, c)
	registerListTags(s, c)
	registerListCorrespondents(s, c)
	registerListDocumentTypes(s, c)
}

func registerHealthCheck(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_health_check",
		mcp.WithDescription("Check the health and connectivity of the Paperless-ngx API."),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Try to list correspondents as a simple connectivity check
		var res interface{}
		err := c.Get(ctx, "correspondents/?page_size=1", &res)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Health check failed: %v", err)), nil
		}
		return mcp.NewToolResultText("Paperless-ngx API is reachable and authenticated."), nil
	})
}

func registerGetStatistics(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_get_statistics",
		mcp.WithDescription("Get system statistics including document counts, inbox status, and storage usage."),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := c.GetStatistics(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerSearchDocuments(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_search_documents",
		mcp.WithDescription("Search for documents in Paperless-ngx."),
		mcp.WithString("query", mcp.Required(), mcp.Description("The search query.")),
		mcp.WithNumber("page", mcp.Description("Page number (default: 1)")),
		mcp.WithNumber("page_size", mcp.Description("Results per page (default: 25, max: 100)")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		query, _ := args["query"].(string)
		
		page := 1
		if p, ok := args["page"].(float64); ok && p > 0 {
			page = int(p)
		}

		pageSize := 25
		if ps, ok := args["page_size"].(float64); ok && ps > 0 {
			pageSize = int(ps)
			if pageSize > 100 { pageSize = 100 }
		}

		res, err := c.SearchDocuments(ctx, query, page, pageSize)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerGetDocument(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_get_document",
		mcp.WithDescription("Get metadata and OCR text content for a specific document in Paperless-ngx."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the document.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		res, err := c.GetDocument(ctx, int(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerListTags(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_list_tags",
		mcp.WithDescription("List all tags in Paperless-ngx."),
		mcp.WithNumber("page", mcp.Description("Page number (default: 1)")),
		mcp.WithNumber("page_size", mcp.Description("Results per page (default: 25, max: 100)")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		page := 1
		if p, ok := args["page"].(float64); ok && p > 0 {
			page = int(p)
		}

		pageSize := 25
		if ps, ok := args["page_size"].(float64); ok && ps > 0 {
			pageSize = int(ps)
			if pageSize > 100 { pageSize = 100 }
		}

		res, err := c.ListTags(ctx, page, pageSize)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerListCorrespondents(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_list_correspondents",
		mcp.WithDescription("List all correspondents in Paperless-ngx."),
		mcp.WithNumber("page", mcp.Description("Page number (default: 1)")),
		mcp.WithNumber("page_size", mcp.Description("Results per page (default: 25, max: 100)")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		page := 1
		if p, ok := args["page"].(float64); ok && p > 0 {
			page = int(p)
		}

		pageSize := 25
		if ps, ok := args["page_size"].(float64); ok && ps > 0 {
			pageSize = int(ps)
			if pageSize > 100 { pageSize = 100 }
		}

		res, err := c.ListCorrespondents(ctx, page, pageSize)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerListDocumentTypes(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_list_document_types",
		mcp.WithDescription("List all document types in Paperless-ngx."),
		mcp.WithNumber("page", mcp.Description("Page number (default: 1)")),
		mcp.WithNumber("page_size", mcp.Description("Results per page (default: 25, max: 100)")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		page := 1
		if p, ok := args["page"].(float64); ok && p > 0 {
			page = int(p)
		}

		pageSize := 25
		if ps, ok := args["page_size"].(float64); ok && ps > 0 {
			pageSize = int(ps)
			if pageSize > 100 { pageSize = 100 }
		}

		res, err := c.ListDocumentTypes(ctx, page, pageSize)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}


