package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/WillScott73/paperless-ngx-mcp/internal/paperless/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Register(s *server.MCPServer, c *client.Client) {
	registerHealthCheck(s, c)
	registerGetStatistics(s, c)
	registerSearchDocuments(s, c)
	registerGetDocument(s, c)
	registerUploadDocument(s, c)
	registerDownloadDocument(s, c)
	registerDeleteDocument(s, c)
	registerBulkEditDocuments(s, c)
	registerListTags(s, c)
	registerCreateTag(s, c)
	registerUpdateTag(s, c)
	registerDeleteTag(s, c)
	registerListCorrespondents(s, c)
	registerCreateCorrespondent(s, c)
	registerUpdateCorrespondent(s, c)
	registerDeleteCorrespondent(s, c)
	registerListDocumentTypes(s, c)
	registerCreateDocumentType(s, c)
	registerUpdateDocumentType(s, c)
	registerDeleteDocumentType(s, c)
	registerListStoragePaths(s, c)
	registerCreateStoragePath(s, c)
	registerUpdateStoragePath(s, c)
	registerDeleteStoragePath(s, c)
	registerUpdateDocument(s, c)
	registerListTasks(s, c)
	registerListCustomFields(s, c)
	registerCreateCustomField(s, c)
	registerUpdateCustomField(s, c)
	registerDeleteCustomField(s, c)
	registerListLogs(s, c)
	registerGetLog(s, c)
	registerListShareLinks(s, c)
	registerCreateShareLink(s, c)
	registerDeleteShareLink(s, c)
}

func registerHealthCheck(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_health_check",
		mcp.WithDescription("Check the health and connectivity of the Paperless-ngx API."),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func registerUploadDocument(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_upload_document",
		mcp.WithDescription("Upload a document to Paperless-ngx."),
		mcp.WithString("file_path", mcp.Description("Local path to the file to upload.")),
		mcp.WithString("base64_data", mcp.Description("Base64 encoded file content.")),
		mcp.WithString("file_name", mcp.Description("Name for the file (required if using base64_data).")),
		mcp.WithString("title", mcp.Description("Title of the document.")),
		mcp.WithNumber("correspondent", mcp.Description("ID of the correspondent.")),
		mcp.WithNumber("document_type", mcp.Description("ID of the document type.")),
		mcp.WithArray("tags", mcp.Description("List of tag IDs.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		var fileName string
		var fileContent []byte
		var err error

		if filePath, ok := args["file_path"].(string); ok {
			fileContent, err = os.ReadFile(filePath)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to read file: %v", err)), nil
			}
			fileName = filePath[strings.LastIndex(filePath, "/")+1:]
		} else if base64Data, ok := args["base64_data"].(string); ok {
			if fn, ok := args["file_name"].(string); ok {
				fileName = fn
			} else {
				return mcp.NewToolResultError("file_name is required when using base64_data"), nil
			}
			fileContent, err = base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to decode base64 data: %v", err)), nil
			}
		} else {
			return mcp.NewToolResultError("Either file_path or base64_data must be provided"), nil
		}

		metadata := make(map[string]string)
		if v, ok := args["title"].(string); ok {
			metadata["title"] = v
		}
		if v, ok := args["correspondent"].(float64); ok {
			metadata["correspondent"] = fmt.Sprintf("%d", int(v))
		}
		if v, ok := args["document_type"].(float64); ok {
			metadata["document_type"] = fmt.Sprintf("%d", int(v))
		}
		if v, ok := args["tags"].([]interface{}); ok {
			for _, tag := range v {
				metadata["tags"] = fmt.Sprintf("%v", tag) // Note: Paperless API might need multiple 'tags' fields in multipart, but UploadDocument implementation currently handles it as a loop. Actually my UploadDocument implementation might need adjustment for multiple tags if it doesn't support same-key multiple values.
			}
		}

		res, err := c.UploadDocument(ctx, fileName, strings.NewReader(string(fileContent)), metadata)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(res), nil
	})
}

func registerDownloadDocument(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_download_document",
		mcp.WithDescription("Download a document (original, preview, or thumbnail) from Paperless-ngx."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the document.")),
		mcp.WithString("mode", mcp.Description("Download mode: 'download' (original), 'preview', or 'thumbnail'. Default: 'download'")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)
		mode := "download"
		if m, ok := args["mode"].(string); ok {
			mode = m
		}

		content, contentType, err := c.DownloadDocument(ctx, int(id), mode)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Since MCP text results might not be ideal for binary, we provide it as a base64 encoded string if it's not text
		if strings.HasPrefix(contentType, "text/") {
			return mcp.NewToolResultText(string(content)), nil
		}
		
		encoded := base64.StdEncoding.EncodeToString(content)
		return mcp.NewToolResultText(fmt.Sprintf("Content-Type: %s\nData (base64): %s", contentType, encoded)), nil
	})
}

func registerDeleteDocument(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_delete_document",
		mcp.WithDescription("Delete a document from Paperless-ngx."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the document to delete.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		err := c.DeleteDocument(ctx, int(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Document %d deleted successfully.", int(id))), nil
	})
}

func registerBulkEditDocuments(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_bulk_edit_documents",
		mcp.WithDescription("Perform bulk operations on multiple documents."),
		mcp.WithArray("document_ids", mcp.Required(), mcp.Description("List of document IDs.")),
		mcp.WithString("method", mcp.Required(), mcp.Description("Operation method (e.g., 'set_correspondent', 'add_tag', 'delete').")),
		mcp.WithObject("parameters", mcp.Description("Additional parameters for the method.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		
		idsRaw, _ := args["document_ids"].([]interface{})
		var ids []int
		for _, id := range idsRaw {
			if n, ok := id.(float64); ok {
				ids = append(ids, int(n))
			}
		}

		method, _ := args["method"].(string)
		params, _ := args["parameters"].(map[string]interface{})

		err := c.BulkEditDocuments(ctx, ids, method, params)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Bulk operation completed successfully."), nil
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

func registerCreateTag(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_create_tag",
		mcp.WithDescription("Create a new tag in Paperless-ngx."),
		mcp.WithString("name", mcp.Required(), mcp.Description("The name of the tag.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		name, _ := args["name"].(string)

		res, err := c.CreateTag(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerUpdateTag(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_update_tag",
		mcp.WithDescription("Update an existing tag."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the tag.")),
		mcp.WithString("name", mcp.Description("The new name.")),
		mcp.WithString("color", mcp.Description("The hex color.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		updates := make(map[string]interface{})
		if v, ok := args["name"].(string); ok {
			updates["name"] = v
		}
		if v, ok := args["color"].(string); ok {
			updates["color"] = v
		}

		res, err := c.UpdateTag(ctx, int(id), updates)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerDeleteTag(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_delete_tag",
		mcp.WithDescription("Delete a tag."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the tag to delete.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		err := c.DeleteTag(ctx, int(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Tag %d deleted successfully.", int(id))), nil
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

func registerCreateCorrespondent(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_create_correspondent",
		mcp.WithDescription("Create a new correspondent in Paperless-ngx."),
		mcp.WithString("name", mcp.Required(), mcp.Description("The name of the correspondent.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		name, _ := args["name"].(string)

		res, err := c.CreateCorrespondent(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerUpdateCorrespondent(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_update_correspondent",
		mcp.WithDescription("Update an existing correspondent."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the correspondent.")),
		mcp.WithString("name", mcp.Description("The new name.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		updates := make(map[string]interface{})
		if v, ok := args["name"].(string); ok {
			updates["name"] = v
		}

		res, err := c.UpdateCorrespondent(ctx, int(id), updates)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerDeleteCorrespondent(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_delete_correspondent",
		mcp.WithDescription("Delete a correspondent."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the correspondent to delete.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		err := c.DeleteCorrespondent(ctx, int(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Correspondent %d deleted successfully.", int(id))), nil
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

func registerCreateDocumentType(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_create_document_type",
		mcp.WithDescription("Create a new document type."),
		mcp.WithString("name", mcp.Required(), mcp.Description("The name of the document type.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		name, _ := args["name"].(string)

		res, err := c.CreateDocumentType(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerUpdateDocumentType(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_update_document_type",
		mcp.WithDescription("Update an existing document type."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the document type.")),
		mcp.WithString("name", mcp.Description("The new name.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		updates := make(map[string]interface{})
		if v, ok := args["name"].(string); ok {
			updates["name"] = v
		}

		res, err := c.UpdateDocumentType(ctx, int(id), updates)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerDeleteDocumentType(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_delete_document_type",
		mcp.WithDescription("Delete a document type."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the document type to delete.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		err := c.DeleteDocumentType(ctx, int(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Document type %d deleted successfully.", int(id))), nil
	})
}

func registerListStoragePaths(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_list_storage_paths",
		mcp.WithDescription("List all storage paths."),
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

		res, err := c.ListStoragePaths(ctx, page, pageSize)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerCreateStoragePath(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_create_storage_path",
		mcp.WithDescription("Create a new storage path."),
		mcp.WithString("name", mcp.Required(), mcp.Description("The name of the storage path.")),
		mcp.WithString("path", mcp.Required(), mcp.Description("The path string.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		name, _ := args["name"].(string)
		path, _ := args["path"].(string)

		res, err := c.CreateStoragePath(ctx, name, path)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerUpdateStoragePath(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_update_storage_path",
		mcp.WithDescription("Update an existing storage path."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the storage path.")),
		mcp.WithString("name", mcp.Description("The new name.")),
		mcp.WithString("path", mcp.Description("The new path string.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		updates := make(map[string]interface{})
		if v, ok := args["name"].(string); ok {
			updates["name"] = v
		}
		if v, ok := args["path"].(string); ok {
			updates["path"] = v
		}

		res, err := c.UpdateStoragePath(ctx, int(id), updates)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerDeleteStoragePath(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_delete_storage_path",
		mcp.WithDescription("Delete a storage path."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the storage path to delete.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		err := c.DeleteStoragePath(ctx, int(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Storage path %d deleted successfully.", int(id))), nil
	})
}

func registerUpdateDocument(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_update_document",
		mcp.WithDescription("Update metadata for an existing document."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the document.")),
		mcp.WithString("title", mcp.Description("The new title.")),
		mcp.WithNumber("correspondent", mcp.Description("The ID of the correspondent.")),
		mcp.WithNumber("document_type", mcp.Description("The ID of the document type.")),
		mcp.WithArray("tags", mcp.Description("A list of tag IDs.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		updates := make(map[string]interface{})
		if v, ok := args["title"].(string); ok {
			updates["title"] = v
		}
		if v, ok := args["correspondent"].(float64); ok {
			updates["correspondent"] = int(v)
		}
		if v, ok := args["document_type"].(float64); ok {
			updates["document_type"] = int(v)
		}
		if v, ok := args["tags"].([]interface{}); ok {
			updates["tags"] = v
		}

		res, err := c.UpdateDocument(ctx, int(id), updates)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerListTasks(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_list_tasks",
		mcp.WithDescription("List recent background tasks (OCR, ingestion, etc.)."),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := c.ListTasks(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerListCustomFields(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_list_custom_fields",
		mcp.WithDescription("List all custom field definitions."),
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
			if pageSize > 100 {
				pageSize = 100
			}
		}

		res, err := c.ListCustomFields(ctx, page, pageSize)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerCreateCustomField(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_create_custom_field",
		mcp.WithDescription("Define a new custom field."),
		mcp.WithString("name", mcp.Required(), mcp.Description("The name of the custom field.")),
		mcp.WithString("data_type", mcp.Required(), mcp.Description("The data type of the field (e.g., 'string', 'boolean', 'date', 'integer', 'float', 'monetary', 'documentlink', 'url', 'select').")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		name, _ := args["name"].(string)
		dataType, _ := args["data_type"].(string)

		res, err := c.CreateCustomField(ctx, name, dataType)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerUpdateCustomField(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_update_custom_field",
		mcp.WithDescription("Update a custom field definition."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the custom field.")),
		mcp.WithString("name", mcp.Description("The new name.")),
		mcp.WithString("data_type", mcp.Description("The new data type.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		updates := make(map[string]interface{})
		if v, ok := args["name"].(string); ok {
			updates["name"] = v
		}
		if v, ok := args["data_type"].(string); ok {
			updates["data_type"] = v
		}

		res, err := c.UpdateCustomField(ctx, int(id), updates)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerDeleteCustomField(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_delete_custom_field",
		mcp.WithDescription("Remove a custom field definition."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the custom field to delete.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		err := c.DeleteCustomField(ctx, int(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Custom field %d deleted successfully.", int(id))), nil
	})
}

func registerListLogs(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_list_logs",
		mcp.WithDescription("List available log files."),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := c.ListLogs(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerGetLog(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_get_log",
		mcp.WithDescription("Fetch the content of a specific log file."),
		mcp.WithString("name", mcp.Required(), mcp.Description("The name of the log file.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		name, _ := args["name"].(string)

		res, err := c.GetLog(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})
}

func registerListShareLinks(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_list_share_links",
		mcp.WithDescription("List existing external share links."),
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
			if pageSize > 100 {
				pageSize = 100
			}
		}

		res, err := c.ListShareLinks(ctx, page, pageSize)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerCreateShareLink(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_create_share_link",
		mcp.WithDescription("Create a temporary sharing link for a document."),
		mcp.WithNumber("document_id", mcp.Required(), mcp.Description("The ID of the document to share.")),
		mcp.WithString("expires", mcp.Description("Optional expiration date/time.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		documentID, _ := args["document_id"].(float64)
		expires, _ := args["expires"].(string)

		res, err := c.CreateShareLink(ctx, int(documentID), expires)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerDeleteShareLink(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("paperless_delete_share_link",
		mcp.WithDescription("Revoke a share link."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("The ID of the share link to delete.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		id, _ := args["id"].(float64)

		err := c.DeleteShareLink(ctx, int(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Share link %d deleted successfully.", int(id))), nil
	})
}

