**1\. API Application (REST/GraphQL API)**
------------------------------------------

```
api/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── auth.go            # Auth endpoints
│   │   │   ├── users.go           # User CRUD endpoints
│   │   │   ├── posts.go           # Post CRUD endpoints
│   │   │   └── health.go          # Health check endpoint
│   │   ├── middleware/
│   │   │   ├── auth.go            # JWT auth middleware
│   │   │   ├── cors.go            # CORS middleware
│   │   │   ├── logging.go         # Request logging
│   │   │   └── ratelimit.go       # Rate limiting
│   │   ├── routes/
│   │   │   └── routes.go          # Route definitions
│   │   └── responses/
│   │       ├── error.go           # Error response formats
│   │       └── success.go         # Success response formats
│   ├── domain/
│   │   ├── user/
│   │   │   ├── model.go           # User domain model
│   │   │   ├── repository.go      # User repository interface
│   │   │   └── service.go         # User business logic
│   │   └── post/
│   │       ├── model.go           # Post domain model
│   │       ├── repository.go      # Post repository interface
│   │       └── service.go         # Post business logic
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── postgres/
│   │   │   │   ├── connection.go  # DB connection
│   │   │   │   ├── user_repo.go   # User repository impl
│   │   │   │   └── post_repo.go   # Post repository impl
│   │   │   └── redis/
│   │   │       └── cache.go       # Redis cache implementation
│   │   ├── external/
│   │   │   ├── email/
│   │   │   │   └── smtp.go        # Email service
│   │   │   └── storage/
│   │   │       └── s3.go          # File storage
│   │   └── auth/
│   │       ├── jwt.go             # JWT implementation
│   │       └── oauth.go           # OAuth implementation
│   ├── config/
│   │   └── config.go              # Configuration
│   └── pkg/
│       ├── validator/
│       │   └── validator.go       # Input validation
│       ├── logger/
│       │   └── logger.go          # Structured logging
│       └── errors/
│           └── errors.go          # Custom error types
├── api/
│   ├── openapi/
│   │   └── spec.yaml              # OpenAPI/Swagger spec
│   └── proto/                     # Protocol buffer files (if using gRPC)
│       ├── user.proto
│       └── post.proto
├── migrations/
│   ├── postgres/
│   │   ├── 001_create_users.up.sql
│   │   ├── 001_create_users.down.sql
│   │   ├── 002_create_posts.up.sql
│   │   └── 002_create_posts.down.sql
│   └── redis/
│       └── init.lua
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   ├── kubernetes/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── ingress.yaml
│   └── terraform/
│       ├── main.tf
│       └── variables.tf
├── scripts/
│   ├── build.sh
│   ├── test.sh
│   ├── migrate.sh
│   └── generate.sh                # Code generation
├── tests/
│   ├── integration/
│   │   ├── api_test.go
│   │   └── db_test.go
│   ├── unit/
│   │   ├── handlers_test.go
│   │   └── services_test.go
│   └── fixtures/
│       ├── users.json
│       └── posts.json
├── configs/
│   ├── local.yaml
│   ├── staging.yaml
│   └── production.yaml
├── docs/
│   ├── api.md                     # API documentation
│   └── deployment.md              # Deployment guide
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
├── Makefile
└── README.md
```

### **API Example Structure:**

go

```
// cmd/api/main.go
package main

import (
    "log"
    "net/http"

    "myapi/internal/api/routes"
    "myapi/internal/config"
    "myapi/internal/infrastructure/database/postgres"
)

func main() {
    cfg := config.Load()
    db := postgres.Connect(cfg.DatabaseURL)

    router := routes.Setup(db, cfg)

    log.Printf("API server starting on :%s", cfg.Port)
    log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
```

**2\. Web Application (Full-Stack with HTML/Templates)**
--------------------------------------------------------

```
webapp/
├── cmd/
│   └── webapp/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── handlers/
│   │   ├── auth.go                # Authentication handlers
│   │   ├── home.go                # Home page handlers
│   │   ├── user.go                # User-related handlers
│   │   └── middleware.go          # HTTP middleware
│   ├── models/
│   │   ├── user.go                # User model
│   │   ├── post.go                # Post model
│   │   └── database.go            # Database connection
│   ├── services/
│   │   ├── auth.go                # Auth business logic
│   │   ├── user.go                # User business logic
│   │   └── email.go               # Email service
│   └── utils/
│       ├── validator.go           # Input validation
│       ├── session.go             # Session management
│       └── helpers.go             # Helper functions
├── web/
│   ├── static/
│   │   ├── css/
│   │   │   ├── main.css
│   │   │   └── bootstrap.min.css
│   │   ├── js/
│   │   │   ├── main.js
│   │   │   └── htmx.min.js
│   │   └── images/
│   │       └── logo.png
│   └── templates/
│       ├── layout/
│       │   ├── base.html          # Base template
│       │   ├── header.html        # Header partial
│       │   └── footer.html        # Footer partial
│       ├── pages/
│       │   ├── home.html          # Home page
│       │   ├── login.html         # Login page
│       │   └── dashboard.html     # Dashboard page
│       └── components/
│           ├── form.html          # Reusable form components
│           └── table.html         # Reusable table components
├── migrations/
│   ├── 001_create_users.sql
│   └── 002_create_posts.sql
├── scripts/
│   ├── build.sh                   # Build script
│   └── deploy.sh                  # Deployment script
├── configs/
│   ├── app.yaml                   # App config
│   └── database.yaml              # Database config
├── go.mod
├── go.sum
├── .env                           # Environment variables
├── .gitignore
├── Dockerfile
├── docker-compose.yml
└── README.md
```

### **Web App Example Structure:**

go

```
// cmd/webapp/main.go
package main

import (
    "log"
    "net/http"

    "myapp/internal/config"
    "myapp/internal/handlers"
    "myapp/internal/models"
)

func main() {
    cfg := config.Load()
    db := models.InitDB(cfg.DatabaseURL)

    h := handlers.New(db)

    // Static files
    http.Handle("/static/", http.StripPrefix("/static/",
        http.FileServer(http.Dir("web/static/"))))

    // Routes
    http.HandleFunc("/", h.Home)
    http.HandleFunc("/login", h.Login)
    http.HandleFunc("/dashboard", h.Dashboard)

    log.Printf("Server starting on :%s", cfg.Port)
    log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
```

**3\. Microservice API Structure**
----------------------------------

```
user-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── user/
│   │   ├── handler.go             # HTTP handlers
│   │   ├── service.go             # Business logic
│   │   ├── repository.go          # Data access
│   │   └── model.go               # Domain models
│   ├── auth/
│   │   ├── middleware.go
│   │   └── jwt.go
│   └── health/
│       └── handler.go
├── pkg/
│   ├── database/
│   │   └── postgres.go
│   ├── logger/
│   │   └── zap.go
│   └── config/
│       └── config.go
├── api/
│   └── proto/
│       ├── user.proto
│       └── generated/
│           └── user.pb.go
├── deployments/
│   ├── Dockerfile
│   └── k8s/
│       ├── deployment.yaml
│       └── service.yaml
├── go.mod
├── go.sum
└── README.md
```