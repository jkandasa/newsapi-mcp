# newsapi-mcp

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server for [newsapi.org](https://newsapi.org/), written in Go. Supports both stdio and HTTP/SSE transports.

## Prerequisites

- A free API key from [newsapi.org](https://newsapi.org/register)
- Go 1.26+ (to build from source)

## Build

```bash
# current platform
make build

# all platforms (linux amd64/arm64/arm32, windows amd64)
make build-all
```

`make build-all` places binaries in `dist/`. The version is injected automatically from `git describe`.

## Usage

Check the version:

```bash
./newsapi-mcp -version
```

Set your API key via the environment variable:

```bash
export NEWSAPI_KEY=your_api_key_here
```

**stdio** (default — for MCP clients like Claude Desktop):

```bash
./newsapi-mcp
```

**HTTP/SSE:**

```bash
./newsapi-mcp -transport http -addr :8080
```

**HTTPS** (provide your own certificate and key):

```bash
./newsapi-mcp -transport http -addr :8443 -cert /path/to/cert.pem -key /path/to/key.pem
```

### Flags

| Flag         | Default | Description                                            |
| ------------ | ------- | ------------------------------------------------------ |
| `-version`   | —       | Print version and exit                                 |
| `-transport` | `stdio` | Transport type: `stdio` or `http`                      |
| `-addr`      | `:8080` | Listen address (HTTP only)                             |
| `-cert`      | —       | TLS certificate file (enables HTTPS; requires `-key`)  |
| `-key`       | —       | TLS private key file (enables HTTPS; requires `-cert`) |

## Tools

| Tool                | Description                                                               |
| ------------------- | ------------------------------------------------------------------------- |
| `get_top_headlines` | Breaking news headlines filtered by country, category, source, or keyword |
| `search_everything` | Full-text search across 150,000+ sources going back 5 years               |
| `get_sources`       | List available news sources and their IDs                                 |

## Claude Desktop setup

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "newsapi": {
      "command": "/path/to/newsapi-mcp",
      "env": {
        "NEWSAPI_KEY": "your_api_key_here"
      }
    }
  }
}
```

## API reference

See [newsapi.org](https://newsapi.org/) for full API documentation, endpoint details, and plan limits.

## License

MIT — see [LICENSE](LICENSE). This license covers the MCP server code only; use of the News API is subject to [newsapi.org's terms](https://newsapi.org/terms).
