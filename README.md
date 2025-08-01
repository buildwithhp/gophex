# 🚀 Gophex - Go Project Generator

A powerful CLI tool that generates production-ready Go project structures following industry best practices and clean architecture principles.

## ✨ Features

- 🏗️ **Clean Architecture** - Domain-driven design with proper separation of concerns
- 🔐 **Security First** - JWT authentication, CORS, rate limiting, and request logging
- 📊 **Database Ready** - PostgreSQL integration with repository pattern
- 🚀 **High Performance** - Redis caching and optimized middleware
- 🧪 **Testing Structure** - Comprehensive unit and integration test organization
- 📝 **API Documentation** - OpenAPI/Swagger specification ready
- 🐳 **DevOps Ready** - Docker, Kubernetes, and deployment configurations
- 🎯 **Interactive CLI** - User-friendly prompts with Survey library
- 📦 **Single Binary** - Templates embedded using Go's embed filesystem

## 🛠️ Installation

### Install from GitHub

```bash
go install github.com/buildwithhp/gophex@latest
```

### Build from Source

```bash
git clone https://github.com/buildwithhp/gophex.git
cd gophex
go build -o gophex main.go
```

## 🎯 Usage

### Interactive Mode

Simply run gophex and follow the interactive prompts:

```bash
gophex
```

### Command Line Mode

```bash
# Generate a new project
gophex generate

# Show version
gophex version

# Show help
gophex help
```

## 📋 Project Types

### 🌐 REST API (Production-Ready)

Generates a comprehensive REST API with clean architecture:

```bash
gophex generate
# Select: api - REST API with clean architecture
# Enter project name: myapi
```

**Generated Structure:**
```
myapi/
├── cmd/api/                    # Application entry point
│   └── main.go                 # Graceful shutdown, DB connections
├── internal/
│   ├── api/                    # HTTP layer
│   │   ├── handlers/           # Request handlers
│   │   │   ├── auth.go         # JWT authentication
│   │   │   ├── users.go        # User CRUD operations
│   │   │   ├── posts.go        # Post management
│   │   │   └── health.go       # Health checks
│   │   ├── middleware/         # HTTP middleware
│   │   │   ├── auth.go         # JWT validation
│   │   │   ├── cors.go         # CORS handling
│   │   │   ├── logging.go      # Request logging
│   │   │   └── ratelimit.go    # Rate limiting
│   │   ├── routes/             # Route definitions
│   │   │   └── routes.go       # API routing setup
│   │   └── responses/          # Response formatting
│   │       ├── error.go        # Error responses
│   │       └── success.go      # Success responses
│   ├── domain/                 # Business logic
│   │   ├── user/               # User domain
│   │   │   ├── model.go        # User entity
│   │   │   ├── repository.go   # Repository interface
│   │   │   └── service.go      # Business logic
│   │   └── post/               # Post domain
│   │       ├── model.go        # Post entity
│   │       ├── repository.go   # Repository interface
│   │       └── service.go      # Business logic
│   ├── infrastructure/         # External dependencies
│   │   ├── database/           # Database implementations
│   │   │   ├── postgres/       # PostgreSQL repositories
│   │   │   └── redis/          # Redis caching
│   │   ├── external/           # External services
│   │   │   ├── email/          # Email service
│   │   │   └── storage/        # File storage
│   │   └── auth/               # Authentication
│   │       ├── jwt.go          # JWT implementation
│   │       └── oauth.go        # OAuth integration
│   ├── config/                 # Configuration management
│   │   └── config.go           # App configuration
│   └── pkg/                    # Shared utilities
│       ├── validator/          # Input validation
│       ├── logger/             # Structured logging
│       └── errors/             # Custom error types
├── api/                        # API specifications
│   ├── openapi/                # OpenAPI/Swagger specs
│   └── proto/                  # Protocol buffers
├── migrations/                 # Database migrations
│   ├── postgres/               # PostgreSQL migrations
│   └── redis/                  # Redis scripts
├── deployments/                # Deployment configurations
│   ├── docker/                 # Docker setup
│   ├── kubernetes/             # K8s manifests
│   └── terraform/              # Infrastructure as code
├── scripts/                    # Build and utility scripts
├── tests/                      # Test files
│   ├── integration/            # Integration tests
│   ├── unit/                   # Unit tests
│   └── fixtures/               # Test data
├── configs/                    # Environment configs
├── docs/                       # Documentation
├── go.mod                      # Go modules
└── README.md                   # Project documentation
```

**Features Included:**
- ✅ JWT Authentication system
- ✅ User management (CRUD)
- ✅ Post management system
- ✅ Security middleware (CORS, rate limiting, logging)
- ✅ PostgreSQL integration
- ✅ Redis caching support
- ✅ Input validation
- ✅ Structured logging
- ✅ Error handling
- ✅ Health check endpoints
- ✅ Clean architecture pattern
- ✅ Repository pattern
- ✅ Dependency injection
- ✅ Graceful shutdown
- ✅ Configuration management

### 🌍 Web Application

```bash
gophex generate
# Select: webapp - Web application with templates
```

Generates a web application with HTML templates and static file serving.

### 🔧 Microservice

```bash
gophex generate
# Select: microservice - Microservice with gRPC support
```

Creates a lightweight microservice with health checks and service endpoints.

### 💻 CLI Tool

```bash
gophex generate
# Select: cli - Command-line tool
```

Generates a CLI application using Cobra framework.

## 🏗️ Architecture Principles

### Clean Architecture

- **Domain Layer**: Business logic and entities
- **Infrastructure Layer**: External dependencies (database, APIs)
- **API Layer**: HTTP handlers and middleware
- **Dependency Inversion**: Interfaces define contracts

### Best Practices

- ✅ **Separation of Concerns** - Each layer has a single responsibility
- ✅ **Dependency Injection** - Loose coupling between components
- ✅ **Repository Pattern** - Abstract data access
- ✅ **Service Layer** - Business logic encapsulation
- ✅ **Middleware Pattern** - Cross-cutting concerns
- ✅ **Error Handling** - Consistent error responses
- ✅ **Configuration Management** - Environment-based config
- ✅ **Structured Logging** - Observability and debugging
- ✅ **Input Validation** - Data integrity and security
- ✅ **Testing Structure** - Unit and integration tests

## 🔧 Technology Stack

### Core Dependencies

- **HTTP Router**: Gorilla Mux
- **Database**: PostgreSQL with lib/pq driver
- **Caching**: Redis
- **Authentication**: JWT with golang-jwt/jwt
- **Password Hashing**: bcrypt
- **Configuration**: YAML-based config
- **Logging**: Structured logging
- **Validation**: Custom validation package

### Development Tools

- **Testing**: Go's built-in testing framework
- **Documentation**: OpenAPI/Swagger
- **Containerization**: Docker
- **Orchestration**: Kubernetes
- **Infrastructure**: Terraform
- **CI/CD**: GitHub Actions ready

## 🚀 Getting Started with Generated Projects

### API Project

1. **Generate the project:**
   ```bash
   gophex generate
   # Select API template and enter project name
   ```

2. **Navigate to project:**
   ```bash
   cd your-project-name
   ```

3. **Install dependencies:**
   ```bash
   go mod tidy
   ```

4. **Set up environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

5. **Run migrations:**
   ```bash
   ./scripts/migrate.sh up
   ```

6. **Start the server:**
   ```bash
   go run cmd/api/main.go
   ```

7. **Test the API:**
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

### API Endpoints

Generated APIs include these endpoints:

- `GET /api/v1/health` - Health check
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/users` - List users (protected)
- `GET /api/v1/users/{id}` - Get user (protected)
- `PUT /api/v1/users/{id}` - Update user (protected)
- `DELETE /api/v1/users/{id}` - Delete user (protected)
- `GET /api/v1/posts` - List posts (public)
- `GET /api/v1/posts/{id}` - Get post (public)
- `POST /api/v1/posts` - Create post (protected)
- `PUT /api/v1/posts/{id}` - Update post (protected)
- `DELETE /api/v1/posts/{id}` - Delete post (protected)

## 🧪 Testing

Generated projects include comprehensive testing structure:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests
go test ./tests/integration/...

# Run unit tests
go test ./internal/...
```

## 📦 Template System

Gophex uses Go's embedded filesystem for templates:

- **Embedded Templates**: All templates are embedded in the binary
- **No External Dependencies**: Single binary distribution
- **Fast Generation**: Templates loaded from memory
- **Version Controlled**: Templates are part of the codebase
- **Easy Maintenance**: Each template is a separate file

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Adding New Templates

1. Create template files in `internal/templates/{type}/`
2. Use `.tmpl` extension for template files
3. Use `{{.ProjectName}}` and `{{.ModuleName}}` variables
4. Templates are automatically discovered by the embedded filesystem

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by modern Go project layouts
- Built with Go's powerful standard library
- Uses industry-standard packages and patterns
- Follows clean architecture principles

## 📞 Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/buildwithhp/gophex/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/buildwithhp/gophex/discussions)
- 📧 **Email**: [Contact Us](mailto:support@gophex.dev)

---

**Made with ❤️ for the Go community**