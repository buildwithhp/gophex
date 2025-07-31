# AGENTS.md - Development Guidelines

## Build/Test/Lint Commands
- `go build ./...` - Build all packages
- `go test ./...` - Run all tests
- `go test -v ./path/to/package` - Run single package tests
- `go test -run TestFunctionName` - Run specific test
- `go fmt ./...` - Format code
- `go vet ./...` - Static analysis
- `golangci-lint run` - Comprehensive linting (if available)

## Project Structure
This is a Go project following standard Go project layout patterns. See `projectstructures.md` for detailed architecture examples including API, web app, and microservice structures.

## Code Style Guidelines
- Follow standard Go conventions (gofmt, go vet)
- Use meaningful package names (lowercase, no underscores)
- Interfaces should be small and focused
- Error handling: always check errors, wrap with context using `fmt.Errorf`
- Naming: use camelCase for unexported, PascalCase for exported
- Imports: group standard library, third-party, then local packages
- Use `context.Context` for cancellation and timeouts
- Prefer composition over inheritance
- Write tests alongside code (package_test.go)
- Use dependency injection for testability