# Short URL Service

A URL shortening service built with Go, Gin, GORM, MySQL, and Redis. Features a command-line interface powered by Cobra and configuration management with Viper.

## Features

- **URL Shortening**: Create short URLs from long URLs with NanoID generation
- **Custom Keys**: Support for custom short URL keys with validation
- **Passkey Protection**: Optional passkey protection for URLs
- **Token Authentication**: UUID-based token authentication (configurable as mandatory)
- **Auto-revoke**: Support for URL expiration with semantic time duration
- **Redis Caching**: Fast URL lookups using Redis cache
- **Click Tracking**: Track click counts for each short URL
- **CLI Interface**: Command-line interface with Cobra
- **Configuration**: Flexible configuration with Viper (YAML, environment variables)
- **Multi-Database**: Support for MySQL, PostgreSQL, and SQLite
- **URL Validation**: Comprehensive URL format validation
- **Comprehensive Tests**: Full test coverage for all components

## Installation

```bash
go build -o shorturl
```

## Configuration

The application supports multiple configuration methods:

### 1. Configuration File (config.yaml)
```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  type: "mysql"  # mysql, postgres, sqlite
  host: "localhost"
  port: 3306
  user: "root"
  password: "password"
  database: "shorturl"
  # For SQLite: type: "sqlite", database: "./shorturl.db"
  # For PostgreSQL: type: "postgres", port: 5432

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

app:
  name: "Short URL Service"
  default_expire: "30d"           # e.g., 10s, 1h, 7d, 1y
  key_length: 6
  cache_duration: "7d"
  require_auth: false             # Mandatory token auth
  allow_custom_keys: true         # Allow custom short keys
  max_url_length: 2048
```

### 2. Environment Variables
```bash
export SHORTURL_SERVER_HOST="0.0.0.0"
export SHORTURL_SERVER_PORT=8080
export SHORTURL_DATABASE_HOST="localhost"
export SHORTURL_DATABASE_USER="root"
export SHORTURL_DATABASE_PASSWORD="password"
export SHORTURL_REDIS_HOST="localhost"
```

### 3. Command Line Flags
```bash
./shorturl serve --host 0.0.0.0 --port 8080 --db-host localhost
```

## Usage

### Available Commands

```bash
# Show help
./shorturl --help

# Run database migrations
./shorturl migrate

# Start the server
./shorturl serve

# Start server with custom config
./shorturl serve --config /path/to/config.yaml

# Start server with custom port
./shorturl serve --port 3000
```

### Command Line Options

```bash
Global Flags:
  --config string         config file (default is ./config.yaml)
  --db-host string        database host
  --db-name string        database name
  --db-password string    database password
  --db-user string        database user
  -H, --host string       server host
  -p, --port int          server port
  --redis-host string     redis host
  --redis-password string redis password
```

## API Endpoints

### URL Operations
- `POST /api/shorten` - Create a short URL
- `GET /:key` - Redirect to the original URL
- `GET /api/info/:key` - Get URL information
- `DELETE /api/urls/:key` - Revoke a URL
- `POST /api/auto-revoke` - Run auto-revoke for expired URLs

### Authentication
- `POST /api/auth/tokens` - Create an auth token
- `GET /api/auth/tokens` - List all tokens (requires auth)
- `DELETE /api/auth/tokens/:token` - Revoke a token (requires auth)

### Health Check
- `GET /health` - Health check endpoint

## Usage Examples

### Create a short URL
```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "long_url": "https://example.com/very/long/url",
    "custom_key": "mykey",
    "passkey": "secret123",
    "expires_in": "30d"
  }'
```

### Examples of expires_in formats:
- `"10s"` - 10 seconds
- `"5m"` - 5 minutes  
- `"2h"` - 2 hours
- `"7d"` - 7 days (converted to 168h)
- `"1y"` - 1 year (converted to 8760h)

### Access with passkey
```bash
curl "http://localhost:8080/mykey?passkey=secret123"
```

### Create auth token
```bash
curl -X POST http://localhost:8080/api/auth/tokens \
  -H "Content-Type: application/json" \
  -d '{"name": "My API Token"}'
```

### Use auth token
```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Authorization: Bearer your-token-here" \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://example.com"}'
```

## Quick Start

1. **Setup Database**: Start MySQL and Redis services
2. **Create Config**: Copy and modify `config.yaml` 
3. **Run Migrations**: `./shorturl migrate`
4. **Start Server**: `./shorturl serve`

The server will start on the configured host and port (default: 0.0.0.0:8080).

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run specific test packages:
```bash
go test ./utils           # Validation tests
go test ./services        # Service logic tests  
go test ./handlers        # HTTP handler tests
go test ./internal/config # Configuration tests
```

## Validation Features

### URL Validation
- Validates URL format and structure
- Ensures proper HTTP/HTTPS protocol
- Checks for valid host names
- Automatically adds HTTPS prefix if missing

### Custom Key Validation  
- Length validation (3-20 characters)
- Character restrictions (alphanumeric, hyphens, underscores)
- Reserved word protection (api, admin, www, etc.)

### Token Generation
- Short URL keys: Generated using NanoID (6 characters, URL-safe)
- Auth tokens: Generated using UUID v4 for security