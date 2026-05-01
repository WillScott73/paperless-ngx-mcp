# Paperless-ngx MCP Server

A production-quality Model Context Protocol (MCP) server for Paperless-ngx, written in Go.

## Features
- Provides tools for an AI assistant to search and retrieve documents from your Paperless-ngx instance.
- Supports both **stdio** (for typical local desktop usage) and **HTTP/SSE** transports.
- Statically compiled single binary for easy deployment.

## Tools Available
- `paperless_health_check`: Check the health and connectivity of the Paperless-ngx API.
- `paperless_get_statistics`: Get system statistics including document counts and storage usage.
- `paperless_search_documents`: Search for documents in Paperless-ngx using a query string. Supports pagination with `page` and `page_size`.
- `paperless_get_document`: Get metadata and the full OCR text content for a specific document.
- `paperless_list_tags`: List all tags. Supports pagination.
- `paperless_list_correspondents`: List all correspondents. Supports pagination.
- `paperless_list_document_types`: List all document types. Supports pagination.

## Example Usage

Here are some things you can ask an AI assistant with this MCP:

- "What's the status of my document inbox?"
- "Find all documents related to 'energy bill' from 2024."
- "Show me my most recent documents (first page, 10 results)."
- "List all documents tagged with 'Work'."
- "What are all the correspondents I have on record?"

## Configuration

The server requires two environment variables:

- `PAPERLESS_BASE_URL`: The full URL to your Paperless-ngx API (e.g., `https://docs.example.com/api/`).
- `PAPERLESS_API_TOKEN`: Your Paperless API Token.

Optional variables:
- `TRANSPORT`: Set to `http` to use Server-Sent Events over HTTP. Defaults to `stdio`.
- `PORT`: The port to bind to in HTTP mode (default: `8080`).
- `HOST`: The host to bind to in HTTP mode (default: `127.0.0.1`).

## Usage

### Stdio Mode (Claude Desktop, Cursor, etc.)

Configure your client to execute the binary directly. Ensure the environment variables are passed to the process.

**Claude Desktop config (`claude_desktop_config.json`):**
```json
{
  "mcpServers": {
    "paperless": {
      "command": "/path/to/paperless-ngx-mcp",
      "env": {
        "PAPERLESS_BASE_URL": "https://docs.yourdomain.com/api/",
        "PAPERLESS_API_TOKEN": "your_api_token"
      }
    }
  }
}
```

### HTTP Mode

Run the server manually:
```bash
export PAPERLESS_BASE_URL="https://docs.yourdomain.com/api/"
export PAPERLESS_API_TOKEN="your_api_token"
export TRANSPORT="http"
export PORT="8080"

./paperless-ngx-mcp
```
The server will bind to `127.0.0.1:8080`.

## Building from source

Requirements: Go 1.25.5+

```bash
go build -trimpath -ldflags="-s -w" -o paperless-ngx-mcp ./cmd/paperless-ngx-mcp
```

## Releases

Tagged releases publish the raw binaries used by `mcpctl install-latest`:

- Linux: `paperless-ngx-mcp`
- Windows: `paperless-ngx-mcp.exe`

Latest Linux:

```text
https://github.com/WillScott73/paperless-ngx-mcp/releases/latest/download/paperless-ngx-mcp
```

Latest Windows:

```text
https://github.com/WillScott73/paperless-ngx-mcp/releases/latest/download/paperless-ngx-mcp.exe
```

## Security Notes
- Treat your `PAPERLESS_API_TOKEN` like a password.
- HTTP mode binds to localhost (`127.0.0.1`) by default for safety.
