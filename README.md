# Paperless-ngx MCP Server

A production-quality Model Context Protocol (MCP) server for Paperless-ngx, written in Go.

## Features
- Provides tools for an AI assistant to search and retrieve documents from your Paperless-ngx instance.
- Supports both **stdio** (for typical local desktop usage) and **HTTP/SSE** transports.
- Statically compiled single binary for easy deployment.

## Tools Available
- `paperless_health_check`: Check API connectivity.
- `paperless_search_documents`: Search for documents with a query.

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

Requirements: Go 1.25+

```bash
go build -o paperless-ngx-mcp ./cmd/paperless-ngx-mcp
```

## Security Notes
- Treat your `PAPERLESS_API_TOKEN` like a password.
- HTTP mode binds to localhost (`127.0.0.1`) by default for safety.
