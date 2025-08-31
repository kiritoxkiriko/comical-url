# Short URL Service

A high-performance URL shortening service built with Go, featuring Redis caching, multiple database backends, and optional authentication.

## Features

- 🚀 Fast URL shortening and redirection
- 🔒 Optional passkey protection for URLs
- ⏰ Configurable URL expiration
- 🗄️ Multiple database backends (MySQL, PostgreSQL, SQLite)
- ⚡ Redis caching for high performance
- 🔐 Token-based authentication
- 🎯 Custom short keys support
- 📊 URL analytics and management
- 🐳 Docker support with docker-compose
- 📖 OpenAPI/Swagger documentation

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd comical-url

# Start the services
make docker-run
# or
docker-compose -f deployments/docker/docker-compose.yaml up
```

### Manual Installation

1. **Prerequisites**
   - Go 1.21+
   - MySQL/PostgreSQL/SQLite
   - Redis (optional, for caching)

2. **Installation**
   ```bash
   # Build the application
   make build
   # or
   go build -o shorturl

   # Copy and configure
   cp configs/config.example.yaml configs/config.yaml
   # Edit configs/config.yaml with your settings

   # Run database migration
   ./shorturl migrate

   # Start the server
   ./shorturl serve
   ```

## Configuration

Configuration is managed through YAML files and environment variables. See `configs/config.example.yaml` for all available options.

### Environment Variables

All configuration options can be overridden using environment variables with the `SHORTURL_` prefix:

```bash
export SHORTURL_SERVER_PORT=8080
export SHORTURL_DATABASE_TYPE=mysql
export SHORTURL_DATABASE_HOST=localhost
export SHORTURL_REDIS_HOST=localhost
```

## API Usage

### Basic URL Shortening

```bash
# Shorten a URL
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://example.com"}'

# Response
{
  "short_key": "abc123",
  "short_url": "http://localhost:8080/abc123",
  "long_url": "https://example.com"
}

# Access the shortened URL
curl http://localhost:8080/abc123
# → Redirects to https://example.com
```

### Advanced Features

```bash
# Custom short key
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "long_url": "https://example.com",
    "custom_key": "my-link",
    "expires_in": "7d"
  }'

# Protected URL with passkey
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "long_url": "https://example.com",
    "passkey": "secret123"
  }'

# Access protected URL
curl "http://localhost:8080/abc123?passkey=secret123"
```

### Authentication

```bash
# Create auth token
curl -X POST http://localhost:8080/api/auth/tokens

# Use token in requests
curl -X POST http://localhost:8080/api/shorten \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://example.com"}'
```

## Development

### Make Commands

```bash
make help          # Show all available commands
make build         # Build the application
make test          # Run tests
make test-coverage # Run tests with coverage
make run           # Build and run
make docker-build  # Build Docker image
make docker-run    # Run with docker-compose
make clean         # Clean build artifacts
```

### Project Structure

```
├── cmd/                    # Application entrypoints
├── internal/              # Private application code
│   ├── handlers/          # HTTP handlers
│   ├── services/          # Business logic
│   ├── models/           # Data models
│   ├── middleware/       # HTTP middleware
│   ├── utils/            # Utility functions
│   └── config/           # Configuration management
├── api/                  # OpenAPI specifications
├── configs/              # Configuration files
├── deployments/          # Docker and deployment files
│   └── docker/
├── scripts/              # Build and utility scripts
├── docs/                 # Documentation
└── test/                # Test data and utilities
```

## Database Support

### MySQL
```yaml
database:
  type: "mysql"
  host: "localhost"
  port: 3306
  user: "root"
  password: "password"
  database: "shorturl"
```

### PostgreSQL
```yaml
database:
  type: "postgres"
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  database: "shorturl"
```

### SQLite
```yaml
database:
  type: "sqlite"
  database: "./shorturl.db"
```

## Deployment

### Docker

```bash
# Build image
make docker-build

# Run with MySQL
docker-compose -f deployments/docker/docker-compose.yaml up

# Run with PostgreSQL
docker-compose -f deployments/docker/docker-compose.yaml --profile postgres up
```

### Manual Deployment

1. Build for target platform:
   ```bash
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o shorturl
   ```

2. Deploy binary with config file

3. Set up database and run migrations:
   ```bash
   ./shorturl migrate
   ```

4. Start the service:
   ```bash
   ./shorturl serve
   ```

## API Documentation

Full API documentation is available in OpenAPI format at `api/shorturl.yaml`. You can use tools like Swagger UI to view and interact with the API.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.