# URL Shortener API

A high-performance URL shortening service built with Go, featuring user authentication, analytics, and comprehensive API documentation.

## 🚀 Features

- **URL Shortening**: Create short URLs for long links
- **User Authentication**: JWT-based authentication with refresh tokens
- **Guest Mode**: Anonymous URL shortening without registration
- **Analytics**: Track click counts and URL statistics
- **Health Checks**: Comprehensive health monitoring endpoints
- **Swagger Documentation**: Interactive API documentation
- **Docker Support**: Containerized deployment with Docker Compose

## 📁 Project Structure

```
url-shortener-app/
├── cmd/                     # Application entry points
├── internal/                # Private application code
├── configs/                 # Configuration management
├── migrations/              # Database migrations
├── deploy/                  # Infrastructure configurations
├── docker-compose.yml       # Development environment
├── docker-compose.prod.yml  # Production environment
└── Dockerfile.*             # Container definitions
```

## 🛠️ Development Setup

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### Quick Start

1. **Clone and setup environment**
   ```bash
   git clone <repository-url>
   cd url-shortener-app
   cp .env.example .env
   ```

2. **Start dependencies**
   ```bash
   docker-compose up -d
   ```

3. **Run migrations**
   ```bash
   docker-compose --profile migrate up migrate
   ```

4. **Generate SSL keys**
   ```bash
   mkdir -p configs/ssl
   openssl ecparam -name prime256v1 -genkey -noout -out configs/ssl/private.pem
   openssl ec -in configs/ssl/private.pem -pubout -out configs/ssl/public.pem
   ```

5. **Start the API**
   ```bash
   go run cmd/api/main.go
   ```

6. **Access Swagger Documentation**
   ```
   http://localhost:8080/swagger/index.html
   ```

## 🐳 Docker Deployment

### Development
```bash
docker-compose up
```

### Production
```bash
docker-compose -f docker-compose.prod.yml up --build
```

## 📊 API Endpoints

### Health
- `GET /health` - Service health check
- `GET /health/ready` - Readiness probe
- `GET /health/live` - Liveness probe

### Authentication
- `POST /api/v1/user/signup` - User registration
- `POST /api/v1/user/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token

### User Management
- `GET /api/v1/user/me` - Get current user profile

### URL Operations
- `POST /api/v1/url/shorten` - Shorten URL (authenticated)
- `POST /api/v1/url/shorten/guest` - Shorten URL (guest)
- `DELETE /api/v1/url/{id}` - Delete URL
- `GET /{code}` - Redirect to original URL

## 🔧 Development Commands

### Database Operations
```bash
# Generate SQLC code
sqlc generate

# Create new migration
sudo atlas migrate diff migration_name \
--dir "file://migrations" \
--to "file://internal/sqlc/schema.sql" \
--dev-url "docker://postgres/16/url-shortener-db" \
--format '{{ sql . " " }}'

# Apply migrations
sudo atlas migrate apply \
--url "postgres://url-shortener-db-user:url-shortener-db-pass@localhost:5432/url-shortener-db?sslmode=disable" \
--dir "file://migrations"
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/service/...
```

### Code Generation
```bash
# Generate Swagger documentation
swag init -g cmd/api/main.go -o internal/docs

# Tidy dependencies
go mod tidy
```

## 🏗️ Architecture

- **Clean Architecture**: Separation of concerns with clear boundaries
- **Dependency Injection**: Loose coupling between components  
- **Repository Pattern**: Data access abstraction
- **Middleware**: Cross-cutting concerns (auth, logging, CORS)
- **Health Checks**: Kubernetes-ready health endpoints

## 🔐 Security

- **JWT Authentication**: Secure token-based authentication
- **ECDSA Signing**: Cryptographic signing for tokens
- **CORS Protection**: Configurable cross-origin policies
- **Input Validation**: Request validation and sanitization

## 📈 Monitoring

- **Health Endpoints**: Ready for Kubernetes probes
- **Structured Logging**: JSON-formatted logs
- **Metrics Ready**: Prepared for Prometheus integration

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## 📄 License

This project is licensed under the Apache 2.0 License.