package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/buildwithhp/gophex/internal/generator"
)

// ProjectConfiguration represents the complete project configuration
type ProjectConfiguration struct {
	Name           string
	Type           string
	Framework      string
	DatabaseConfig *generator.DatabaseConfig
	RedisConfig    *generator.RedisConfig
	Path           string
	Features       []ProjectFeature
}

// ProjectFeature represents a feature that can be enabled in the project
type ProjectFeature struct {
	Name        string
	Description string
	Enabled     bool
	Educational string
}

// RunEnhancedProjectWizard runs the enhanced educational project generation wizard
func RunEnhancedProjectWizard() error {
	clearScreen()
	fmt.Println("🎓 Enhanced Project Generation Wizard")
	fmt.Println("Learn Go project architecture by building step-by-step!")
	fmt.Println()

	config := &ProjectConfiguration{}

	// Step 1: Project Architecture Overview
	if err := showProjectArchitectureOverview(); err != nil {
		return err
	}

	// Step 2: Project Type Selection with Education
	if err := selectProjectTypeWithEducation(config); err != nil {
		return err
	}

	// Step 3: Project Naming and Structure
	if err := configureProjectBasics(config); err != nil {
		return err
	}

	// Step 4: Framework Selection (if applicable)
	if config.Type == "api" {
		if err := selectFrameworkWithEducation(config); err != nil {
			return err
		}
	}

	// Step 5: Database Architecture Design
	if config.Type == "api" || config.Type == "webapp" {
		if err := designDatabaseArchitecture(config); err != nil {
			return err
		}
	}

	// Step 6: Feature Selection and Configuration
	if err := configureProjectFeatures(config); err != nil {
		return err
	}

	// Step 7: Project Structure Visualization
	if err := visualizeProjectStructure(config); err != nil {
		return err
	}

	// Step 8: Generate and Explain
	if err := generateProjectWithExplanation(config); err != nil {
		return err
	}

	return nil
}

// showProjectArchitectureOverview provides an overview of Go project architectures
func showProjectArchitectureOverview() error {
	fmt.Println("📚 Go Project Architecture Overview")
	fmt.Println("Let's explore different Go project types and their architectures:")
	fmt.Println()

	architectures := []struct {
		Type        string
		Description string
		UseCase     string
		Structure   string
		Examples    string
	}{
		{
			Type:        "API (REST/GraphQL)",
			Description: "Backend services with HTTP endpoints",
			UseCase:     "Microservices, web backends, mobile app APIs",
			Structure:   "Clean Architecture with domain-driven design",
			Examples:    "E-commerce API, user management service, payment gateway",
		},
		{
			Type:        "Web Application",
			Description: "Full-stack web applications with server-side rendering",
			UseCase:     "Traditional web apps, admin dashboards, content sites",
			Structure:   "MVC pattern with template rendering",
			Examples:    "Blog platform, CMS, admin panel, documentation site",
		},
		{
			Type:        "CLI Tool",
			Description: "Command-line applications and utilities",
			UseCase:     "Developer tools, system utilities, automation scripts",
			Structure:   "Command pattern with subcommands",
			Examples:    "Git, Docker CLI, kubectl, custom build tools",
		},
		{
			Type:        "Microservice",
			Description: "Distributed service with gRPC and health checks",
			UseCase:     "Service mesh, distributed systems, cloud-native apps",
			Structure:   "Hexagonal architecture with ports and adapters",
			Examples:    "User service, notification service, analytics service",
		},
	}

	for i, arch := range architectures {
		fmt.Printf("%d. %s\n", i+1, arch.Type)
		fmt.Printf("   📖 %s\n", arch.Description)
		fmt.Printf("   🎯 Use Case: %s\n", arch.UseCase)
		fmt.Printf("   🏗️  Architecture: %s\n", arch.Structure)
		fmt.Printf("   💡 Examples: %s\n\n", arch.Examples)
	}

	var proceed string
	proceedPrompt := &survey.Select{
		Message: "Ready to choose your project architecture?",
		Options: []string{
			"Yes - Let's start building",
			"Tell me more about Clean Architecture",
			"Explain the differences between these types",
			"Quit",
		},
	}

	if err := survey.AskOne(proceedPrompt, &proceed); err != nil {
		return err
	}

	if proceed == "Quit" {
		return nil
	}

	if strings.HasPrefix(proceed, "Tell me more") {
		clearScreen()
		return explainCleanArchitecture()
	}

	if strings.HasPrefix(proceed, "Explain the differences") {
		clearScreen()
		return explainProjectTypeDifferences()
	}

	clearScreen()
	return nil
}

// explainCleanArchitecture provides detailed Clean Architecture explanation
func explainCleanArchitecture() error {
	fmt.Println("\n🏛️  Clean Architecture Deep Dive")
	fmt.Println("Clean Architecture is a software design philosophy that separates concerns.")
	fmt.Println()

	fmt.Println("🎯 Core Principles:")
	principles := []string{
		"Independence: Business rules don't depend on frameworks, databases, or UI",
		"Testability: Business logic can be tested without external dependencies",
		"Flexibility: Easy to change databases, frameworks, or external services",
		"Maintainability: Clear separation makes code easier to understand and modify",
	}

	for _, principle := range principles {
		fmt.Printf("• %s\n", principle)
	}

	fmt.Println("\n🔄 The Dependency Rule:")
	fmt.Println("Dependencies point inward. Outer layers depend on inner layers, never the reverse.")
	fmt.Println()

	fmt.Println("📊 Layer Structure (from inside out):")
	fmt.Println("1. 🏛️  Domain Layer (Entities, Business Rules)")
	fmt.Println("   - Pure business logic")
	fmt.Println("   - No external dependencies")
	fmt.Println("   - Example: User entity with validation rules")
	fmt.Println()

	fmt.Println("2. 🔧 Application Layer (Use Cases)")
	fmt.Println("   - Orchestrates domain objects")
	fmt.Println("   - Implements application-specific business rules")
	fmt.Println("   - Example: CreateUser use case")
	fmt.Println()

	fmt.Println("3. 🔌 Interface Adapters (Controllers, Gateways)")
	fmt.Println("   - Converts data between use cases and external world")
	fmt.Println("   - HTTP handlers, database repositories")
	fmt.Println("   - Example: UserController, UserRepository")
	fmt.Println()

	fmt.Println("4. 🌐 Frameworks & Drivers (Web, DB, External APIs)")
	fmt.Println("   - External tools and frameworks")
	fmt.Println("   - Gin/Echo, PostgreSQL, Redis")
	fmt.Println("   - Example: HTTP server, database connection")

	var ready string
	readyPrompt := &survey.Select{
		Message: "Ready to apply Clean Architecture to your project?",
		Options: []string{
			"Yes - Let's build with Clean Architecture",
			"Show me a real example",
			"Back to project types",
		},
	}

	if err := survey.AskOne(readyPrompt, &ready); err != nil {
		return err
	}

	if strings.HasPrefix(ready, "Show me") {
		return showCleanArchitectureExample()
	}

	if strings.HasPrefix(ready, "Back") {
		return showProjectArchitectureOverview()
	}

	return nil
}

// showCleanArchitectureExample shows a concrete example
func showCleanArchitectureExample() error {
	fmt.Println("\n💡 Clean Architecture Example: User Management API")
	fmt.Println()

	fmt.Println("📁 Project Structure:")
	fmt.Println("```")
	fmt.Println("internal/")
	fmt.Println("├── domain/              # 🏛️  Domain Layer")
	fmt.Println("│   └── user/")
	fmt.Println("│       ├── user.go      # User entity with business rules")
	fmt.Println("│       └── repository.go # Repository interface (contract)")
	fmt.Println("├── application/         # 🔧 Application Layer")
	fmt.Println("│   └── user/")
	fmt.Println("│       └── service.go   # CreateUser, UpdateUser use cases")
	fmt.Println("├── infrastructure/      # 🔌 Interface Adapters")
	fmt.Println("│   ├── http/")
	fmt.Println("│   │   └── user_handler.go # HTTP handlers")
	fmt.Println("│   └── database/")
	fmt.Println("│       └── user_repo.go    # Database implementation")
	fmt.Println("└── main.go             # 🌐 Framework & Drivers")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("🔄 Data Flow Example (Create User):")
	fmt.Println("1. HTTP Request → UserHandler (Infrastructure)")
	fmt.Println("2. Handler validates input → calls UserService (Application)")
	fmt.Println("3. Service applies business rules → calls UserRepository (Domain Interface)")
	fmt.Println("4. Repository saves to database (Infrastructure Implementation)")
	fmt.Println("5. Response flows back through the layers")
	fmt.Println()

	fmt.Println("🧪 Testing Benefits:")
	fmt.Println("• Test business logic without database (mock repository)")
	fmt.Println("• Test use cases without HTTP (call service directly)")
	fmt.Println("• Test handlers without business logic (mock service)")

	var proceed string
	proceedPrompt := &survey.Select{
		Message: "Ready to build your project with this architecture?",
		Options: []string{
			"Yes - Let's start building",
			"Back to architecture overview",
		},
	}

	return survey.AskOne(proceedPrompt, &proceed)
}

// explainProjectTypeDifferences explains the differences between project types
func explainProjectTypeDifferences() error {
	fmt.Println("\n🔍 Project Type Comparison")
	fmt.Println()

	comparisons := []struct {
		Aspect string
		API    string
		WebApp string
		CLI    string
		Micro  string
	}{
		{
			Aspect: "Primary Interface",
			API:    "HTTP REST/GraphQL endpoints",
			WebApp: "HTML pages with forms",
			CLI:    "Command-line arguments",
			Micro:  "gRPC + HTTP health checks",
		},
		{
			Aspect: "Client Interaction",
			API:    "Mobile apps, SPAs, other services",
			WebApp: "Web browsers (server-rendered)",
			CLI:    "Terminal/shell scripts",
			Micro:  "Other microservices",
		},
		{
			Aspect: "State Management",
			API:    "Stateless (database/cache)",
			WebApp: "Session-based state",
			CLI:    "Stateless (file-based config)",
			Micro:  "Stateless (distributed state)",
		},
		{
			Aspect: "Deployment",
			API:    "Containers, cloud platforms",
			WebApp: "Traditional servers, containers",
			CLI:    "Binary distribution",
			Micro:  "Kubernetes, service mesh",
		},
		{
			Aspect: "Scaling",
			API:    "Horizontal (load balancers)",
			WebApp: "Vertical + horizontal",
			CLI:    "Not applicable",
			Micro:  "Auto-scaling, service discovery",
		},
	}

	fmt.Printf("%-20s | %-25s | %-25s | %-20s | %-25s\n", "Aspect", "API", "WebApp", "CLI", "Microservice")
	fmt.Println(strings.Repeat("-", 120))

	for _, comp := range comparisons {
		fmt.Printf("%-20s | %-25s | %-25s | %-20s | %-25s\n",
			comp.Aspect, comp.API, comp.WebApp, comp.CLI, comp.Micro)
	}

	fmt.Println()
	fmt.Println("🎯 When to Choose Each Type:")
	fmt.Println()

	fmt.Println("📡 Choose API when:")
	fmt.Println("• Building backend for mobile apps or SPAs")
	fmt.Println("• Creating microservices architecture")
	fmt.Println("• Need to serve multiple client types")
	fmt.Println("• Building headless/API-first applications")
	fmt.Println()

	fmt.Println("🌐 Choose WebApp when:")
	fmt.Println("• Building traditional web applications")
	fmt.Println("• Need server-side rendering for SEO")
	fmt.Println("• Creating admin dashboards or internal tools")
	fmt.Println("• Want simpler deployment and state management")
	fmt.Println()

	fmt.Println("💻 Choose CLI when:")
	fmt.Println("• Building developer tools or utilities")
	fmt.Println("• Creating automation scripts")
	fmt.Println("• Need cross-platform binary distribution")
	fmt.Println("• Building system administration tools")
	fmt.Println()

	fmt.Println("🔧 Choose Microservice when:")
	fmt.Println("• Building distributed systems")
	fmt.Println("• Need service-to-service communication")
	fmt.Println("• Implementing domain-driven design")
	fmt.Println("• Deploying in Kubernetes/service mesh")

	var proceed string
	proceedPrompt := &survey.Select{
		Message: "Ready to choose your project type?",
		Options: []string{
			"Yes - Let's select a project type",
			"Back to architecture overview",
		},
	}

	return survey.AskOne(proceedPrompt, &proceed)
}

// selectProjectTypeWithEducation handles project type selection with educational content
func selectProjectTypeWithEducation(config *ProjectConfiguration) error {
	clearScreen()
	fmt.Println("🏗️  Step 1: Project Type Selection")
	fmt.Println("Choose the type of Go project you want to build:")
	fmt.Println()

	projectTypes := []string{
		"api - REST API with Clean Architecture (recommended for learning)",
		"webapp - Web application with server-side rendering",
		"microservice - Distributed service with gRPC support",
		"cli - Command-line tool with subcommands",
		"Quit",
	}

	var selected string
	typePrompt := &survey.Select{
		Message: "What type of project do you want to create?",
		Options: projectTypes,
		Help:    "Each type teaches different Go patterns and architectures",
	}

	if err := survey.AskOne(typePrompt, &selected); err != nil {
		return err
	}

	if selected == "Quit" {
		return nil
	}

	// Extract project type
	switch {
	case strings.HasPrefix(selected, "api"):
		config.Type = "api"
	case strings.HasPrefix(selected, "webapp"):
		config.Type = "webapp"
	case strings.HasPrefix(selected, "microservice"):
		config.Type = "microservice"
	case strings.HasPrefix(selected, "cli"):
		config.Type = "cli"
	}

	// Provide educational context for the selected type
	return explainSelectedProjectType(config.Type)
}

// explainSelectedProjectType provides detailed explanation of the selected project type
func explainSelectedProjectType(projectType string) error {
	fmt.Printf("\n🎓 You selected: %s\n", strings.ToUpper(projectType))
	fmt.Println()

	switch projectType {
	case "api":
		fmt.Println("📡 REST API Project")
		fmt.Println("Perfect choice for learning Clean Architecture!")
		fmt.Println()
		fmt.Println("What you'll learn:")
		fmt.Println("• Clean Architecture principles and layers")
		fmt.Println("• Domain-driven design patterns")
		fmt.Println("• HTTP handler patterns in Go")
		fmt.Println("• Database integration and repository pattern")
		fmt.Println("• Middleware composition and request processing")
		fmt.Println("• JWT authentication and authorization")
		fmt.Println("• API documentation and testing strategies")
		fmt.Println()
		fmt.Println("Generated structure:")
		fmt.Println("• Domain entities with business rules")
		fmt.Println("• Repository interfaces and implementations")
		fmt.Println("• Service layer for use cases")
		fmt.Println("• HTTP handlers with proper error handling")
		fmt.Println("• Database migrations and configuration")
		fmt.Println("• Comprehensive test examples")

	case "webapp":
		fmt.Println("🌐 Web Application Project")
		fmt.Println("Great for learning traditional web development patterns!")
		fmt.Println()
		fmt.Println("What you'll learn:")
		fmt.Println("• MVC (Model-View-Controller) pattern")
		fmt.Println("• HTML template rendering in Go")
		fmt.Println("• Session management and cookies")
		fmt.Println("• Form handling and validation")
		fmt.Println("• Static asset serving")
		fmt.Println("• Server-side rendering techniques")
		fmt.Println()
		fmt.Println("Generated structure:")
		fmt.Println("• Template-based views")
		fmt.Println("• Controller handlers for web pages")
		fmt.Println("• Static assets (CSS, JS, images)")
		fmt.Println("• Session management middleware")
		fmt.Println("• Form processing examples")

	case "microservice":
		fmt.Println("🔧 Microservice Project")
		fmt.Println("Advanced pattern for distributed systems!")
		fmt.Println()
		fmt.Println("What you'll learn:")
		fmt.Println("• Hexagonal architecture (ports and adapters)")
		fmt.Println("• gRPC service definitions and implementation")
		fmt.Println("• Health check patterns")
		fmt.Println("• Service discovery concepts")
		fmt.Println("• Distributed tracing and monitoring")
		fmt.Println("• Configuration management")
		fmt.Println()
		fmt.Println("Generated structure:")
		fmt.Println("• gRPC service definitions (.proto files)")
		fmt.Println("• Service implementation with business logic")
		fmt.Println("• Health check endpoints")
		fmt.Println("• Configuration and environment handling")
		fmt.Println("• Docker containerization setup")

	case "cli":
		fmt.Println("💻 CLI Tool Project")
		fmt.Println("Perfect for learning command-line application patterns!")
		fmt.Println()
		fmt.Println("What you'll learn:")
		fmt.Println("• Command pattern and subcommands")
		fmt.Println("• Flag parsing and validation")
		fmt.Println("• Configuration file handling")
		fmt.Println("• Output formatting and colors")
		fmt.Println("• Cross-platform compatibility")
		fmt.Println("• Binary distribution strategies")
		fmt.Println()
		fmt.Println("Generated structure:")
		fmt.Println("• Root command with subcommands")
		fmt.Println("• Flag definitions and parsing")
		fmt.Println("• Configuration management")
		fmt.Println("• Output formatting utilities")
		fmt.Println("• Build and release automation")
	}

	var proceed string
	proceedPrompt := &survey.Select{
		Message: "Continue with this project type?",
		Options: []string{
			"Yes - Continue with " + projectType,
			"No - Choose a different type",
			"Quit",
		},
	}

	if err := survey.AskOne(proceedPrompt, &proceed); err != nil {
		return err
	}

	if proceed == "Quit" {
		return nil
	}

	if strings.HasPrefix(proceed, "No") {
		return ErrReturnToMenu // This will restart the project type selection
	}

	return nil
}

// configureProjectBasics handles project name and path configuration
func configureProjectBasics(config *ProjectConfiguration) error {
	clearScreen()
	fmt.Println("📝 Step 2: Project Configuration")
	fmt.Println("Let's configure the basic details of your project:")
	fmt.Println()

	// Project name
	namePrompt := &survey.Input{
		Message: "What is the name of your project?",
		Help:    "This will be used as the directory name and Go module name. Use lowercase with hyphens (e.g., 'my-api', 'user-service')",
	}

	if err := survey.AskOne(namePrompt, &config.Name, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	config.Path = filepath.Join(currentDir, config.Name)

	// Path confirmation
	var confirm string
	confirmPrompt := &survey.Select{
		Message: fmt.Sprintf("Create project '%s' in %s?", config.Name, config.Path),
		Options: []string{
			"Yes - Create project here",
			"No - Choose different location",
			"Quit",
		},
	}

	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return err
	}

	if confirm == "Quit" {
		return nil
	}

	if strings.HasPrefix(confirm, "No") {
		// Ask for custom path
		var customPath string
		pathPrompt := &survey.Input{
			Message: "Enter the directory path where you want to create the project:",
			Default: currentDir,
			Help:    "The project folder will be created inside this directory",
		}

		if err := survey.AskOne(pathPrompt, &customPath, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

		config.Path = filepath.Join(customPath, config.Name)
	}

	fmt.Printf("✅ Project '%s' will be created at: %s\n", config.Name, config.Path)
	return nil
}

// selectFrameworkWithEducation handles framework selection with educational content
func selectFrameworkWithEducation(config *ProjectConfiguration) error {
	clearScreen()
	fmt.Println("🚀 Step 3: Web Framework Selection")
	fmt.Println("Choose a web framework for your API. Each has different strengths:")
	fmt.Println()

	frameworks := []struct {
		Name        string
		Description string
		Strengths   []string
		UseCase     string
		Learning    string
	}{
		{
			Name:        "gin",
			Description: "Fast HTTP web framework with martini-like API",
			Strengths:   []string{"High performance", "Simple API", "Great middleware ecosystem", "JSON binding"},
			UseCase:     "High-performance APIs, microservices, rapid prototyping",
			Learning:    "Learn performance optimization and middleware patterns",
		},
		{
			Name:        "echo",
			Description: "High performance, extensible, minimalist web framework",
			Strengths:   []string{"Minimal overhead", "Built-in middleware", "WebSocket support", "HTTP/2 support"},
			UseCase:     "Real-time applications, WebSocket APIs, modern web services",
			Learning:    "Learn modern web standards and real-time communication",
		},
		{
			Name:        "gorilla",
			Description: "Powerful HTTP toolkit with flexible routing",
			Strengths:   []string{"Flexible routing", "Rich middleware", "WebSocket support", "Session management"},
			UseCase:     "Complex routing requirements, traditional web apps, enterprise APIs",
			Learning:    "Learn advanced routing patterns and HTTP toolkit usage",
		},
	}

	// Show detailed comparison
	for i, fw := range frameworks {
		fmt.Printf("%d. %s - %s\n", i+1, strings.ToUpper(fw.Name), fw.Description)
		fmt.Printf("   💪 Strengths: %s\n", strings.Join(fw.Strengths, ", "))
		fmt.Printf("   🎯 Best for: %s\n", fw.UseCase)
		fmt.Printf("   🎓 You'll learn: %s\n\n", fw.Learning)
	}

	// Framework selection
	frameworkOptions := []string{
		"gin - Fast and simple (recommended for beginners)",
		"echo - Modern and minimal (good for real-time apps)",
		"gorilla - Flexible and powerful (best for complex routing)",
		"Compare frameworks in detail",
		"Quit",
	}

	var selected string
	frameworkPrompt := &survey.Select{
		Message: "Which web framework would you like to use?",
		Options: frameworkOptions,
		Help:    "Each framework teaches different patterns and approaches",
	}

	if err := survey.AskOne(frameworkPrompt, &selected); err != nil {
		return err
	}

	if selected == "Quit" {
		return nil
	}

	if strings.HasPrefix(selected, "Compare") {
		return showFrameworkComparison(config)
	}

	// Extract framework name
	switch {
	case strings.HasPrefix(selected, "gin"):
		config.Framework = "gin"
	case strings.HasPrefix(selected, "echo"):
		config.Framework = "echo"
	case strings.HasPrefix(selected, "gorilla"):
		config.Framework = "gorilla"
	}

	// Show what they'll learn with this framework
	return explainFrameworkChoice(config.Framework)
}

// showFrameworkComparison provides detailed framework comparison
func showFrameworkComparison(config *ProjectConfiguration) error {
	fmt.Println("\n📊 Detailed Framework Comparison")
	fmt.Println()

	fmt.Println("🏃 Performance Comparison:")
	fmt.Println("• Gin: ~40,000 req/sec (fastest)")
	fmt.Println("• Echo: ~35,000 req/sec (very fast)")
	fmt.Println("• Gorilla: ~25,000 req/sec (good performance)")
	fmt.Println()

	fmt.Println("📚 Learning Curve:")
	fmt.Println("• Gin: Easy (simple API, good docs)")
	fmt.Println("• Echo: Medium (more features, modern patterns)")
	fmt.Println("• Gorilla: Medium-Hard (flexible but complex)")
	fmt.Println()

	fmt.Println("🔧 Middleware Ecosystem:")
	fmt.Println("• Gin: Large ecosystem, many third-party packages")
	fmt.Println("• Echo: Built-in middleware, growing ecosystem")
	fmt.Println("• Gorilla: Rich toolkit, enterprise-focused")
	fmt.Println()

	fmt.Println("🎯 Code Example Comparison:")
	fmt.Println()

	fmt.Println("Gin:")
	fmt.Println("```go")
	fmt.Println("r := gin.Default()")
	fmt.Println("r.GET(\"/users/:id\", getUserHandler)")
	fmt.Println("r.Run(\":8080\")")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("Echo:")
	fmt.Println("```go")
	fmt.Println("e := echo.New()")
	fmt.Println("e.GET(\"/users/:id\", getUserHandler)")
	fmt.Println("e.Start(\":8080\")")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("Gorilla:")
	fmt.Println("```go")
	fmt.Println("r := mux.NewRouter()")
	fmt.Println("r.HandleFunc(\"/users/{id}\", getUserHandler).Methods(\"GET\")")
	fmt.Println("http.ListenAndServe(\":8080\", r)")
	fmt.Println("```")

	// Return to framework selection
	return selectFrameworkWithEducation(config)
}

// explainFrameworkChoice explains what the user will learn with their chosen framework
func explainFrameworkChoice(framework string) error {
	fmt.Printf("\n🎉 Excellent choice: %s!\n", strings.ToUpper(framework))
	fmt.Println()

	switch framework {
	case "gin":
		fmt.Println("🚀 With Gin, you'll learn:")
		fmt.Println("• High-performance HTTP handling")
		fmt.Println("• JSON binding and validation")
		fmt.Println("• Middleware composition patterns")
		fmt.Println("• Route grouping and organization")
		fmt.Println("• Custom validators and error handling")
		fmt.Println()
		fmt.Println("💡 Gin is perfect for:")
		fmt.Println("• Learning Go web development fundamentals")
		fmt.Println("• Building high-performance APIs")
		fmt.Println("• Rapid prototyping and development")

	case "echo":
		fmt.Println("🌊 With Echo, you'll learn:")
		fmt.Println("• Modern HTTP/2 and WebSocket support")
		fmt.Println("• Built-in middleware patterns")
		fmt.Println("• Context-based request handling")
		fmt.Println("• Advanced routing and grouping")
		fmt.Println("• Real-time communication patterns")
		fmt.Println()
		fmt.Println("💡 Echo is perfect for:")
		fmt.Println("• Modern web API development")
		fmt.Println("• Real-time applications")
		fmt.Println("• Learning contemporary web standards")

	case "gorilla":
		fmt.Println("🦍 With Gorilla, you'll learn:")
		fmt.Println("• Advanced routing patterns and constraints")
		fmt.Println("• Rich HTTP toolkit usage")
		fmt.Println("• Session management and cookies")
		fmt.Println("• WebSocket implementation")
		fmt.Println("• Enterprise-grade middleware patterns")
		fmt.Println()
		fmt.Println("💡 Gorilla is perfect for:")
		fmt.Println("• Complex routing requirements")
		fmt.Println("• Enterprise applications")
		fmt.Println("• Learning comprehensive HTTP handling")
	}

	var proceed string
	proceedPrompt := &survey.Select{
		Message: "Continue with " + framework + "?",
		Options: []string{
			"Yes - Continue with " + framework,
			"No - Choose different framework",
			"Quit",
		},
	}

	if err := survey.AskOne(proceedPrompt, &proceed); err != nil {
		return err
	}

	if strings.HasPrefix(proceed, "No") {
		return ErrReturnToMenu // Return to framework selection
	}

	return nil
}

// designDatabaseArchitecture handles database configuration with educational content
func designDatabaseArchitecture(config *ProjectConfiguration) error {
	clearScreen()
	fmt.Println("🗄️  Step 4: Database Architecture Design")
	fmt.Println("Let's design your data layer with best practices:")
	fmt.Println()

	fmt.Println("📚 Database Layer in Clean Architecture:")
	fmt.Println("• Repository Pattern: Abstract data access behind interfaces")
	fmt.Println("• Domain Independence: Business logic doesn't know about SQL")
	fmt.Println("• Testability: Easy to mock for unit tests")
	fmt.Println("• Flexibility: Can swap databases without changing business logic")
	fmt.Println()

	// Database selection with education
	if err := selectDatabaseWithEducation(config); err != nil {
		return err
	}

	// Configuration type selection
	if err := selectDatabaseConfigurationWithEducation(config); err != nil {
		return err
	}

	// Redis configuration
	if err := configureRedisWithEducation(config); err != nil {
		return err
	}

	return nil
}

// selectDatabaseWithEducation handles database selection with educational content
func selectDatabaseWithEducation(config *ProjectConfiguration) error {
	fmt.Println("🎯 Database Selection:")
	fmt.Println()

	databases := []struct {
		Name        string
		Type        string
		Description string
		Strengths   []string
		UseCase     string
		Learning    string
	}{
		{
			Name:        "PostgreSQL",
			Type:        "postgresql",
			Description: "Advanced open-source relational database",
			Strengths:   []string{"ACID compliance", "JSON support", "Advanced indexing", "Full-text search"},
			UseCase:     "Complex queries, data integrity, enterprise applications",
			Learning:    "Learn advanced SQL, transactions, and relational design",
		},
		{
			Name:        "MySQL",
			Type:        "mysql",
			Description: "Popular open-source relational database",
			Strengths:   []string{"High performance", "Wide adoption", "Great tooling", "Easy replication"},
			UseCase:     "Web applications, read-heavy workloads, simple schemas",
			Learning:    "Learn SQL fundamentals and web-scale database patterns",
		},
		{
			Name:        "MongoDB",
			Type:        "mongodb",
			Description: "Document-oriented NoSQL database",
			Strengths:   []string{"Flexible schema", "JSON documents", "Horizontal scaling", "Aggregation pipeline"},
			UseCase:     "Rapid development, flexible schemas, document storage",
			Learning:    "Learn NoSQL patterns and document-based design",
		},
	}

	for i, db := range databases {
		fmt.Printf("%d. %s - %s\n", i+1, db.Name, db.Description)
		fmt.Printf("   💪 Strengths: %s\n", strings.Join(db.Strengths, ", "))
		fmt.Printf("   🎯 Best for: %s\n", db.UseCase)
		fmt.Printf("   🎓 You'll learn: %s\n\n", db.Learning)
	}

	dbOptions := []string{
		"PostgreSQL - Advanced relational database (recommended for learning)",
		"MySQL - Popular and simple relational database",
		"MongoDB - Flexible document database",
		"Compare databases in detail",
		"Quit",
	}

	var selected string
	dbPrompt := &survey.Select{
		Message: "Which database would you like to use?",
		Options: dbOptions,
		Help:    "Each database teaches different data modeling approaches",
	}

	if err := survey.AskOne(dbPrompt, &selected); err != nil {
		return err
	}

	if selected == "Quit" {
		return nil
	}

	if strings.HasPrefix(selected, "Compare") {
		return showDatabaseComparison(config)
	}

	// Initialize database config
	config.DatabaseConfig = &generator.DatabaseConfig{}

	// Extract database type
	switch {
	case strings.HasPrefix(selected, "PostgreSQL"):
		config.DatabaseConfig.Type = "postgresql"
	case strings.HasPrefix(selected, "MySQL"):
		config.DatabaseConfig.Type = "mysql"
	case strings.HasPrefix(selected, "MongoDB"):
		config.DatabaseConfig.Type = "mongodb"
	}

	return explainDatabaseChoice(config.DatabaseConfig.Type)
}

// showDatabaseComparison provides detailed database comparison
func showDatabaseComparison(config *ProjectConfiguration) error {
	fmt.Println("\n📊 Detailed Database Comparison")
	fmt.Println()

	fmt.Println("🏗️  Data Model:")
	fmt.Println("• PostgreSQL: Relational (tables, rows, columns) + JSON")
	fmt.Println("• MySQL: Relational (tables, rows, columns)")
	fmt.Println("• MongoDB: Document-based (JSON-like documents)")
	fmt.Println()

	fmt.Println("🔍 Query Language:")
	fmt.Println("• PostgreSQL: Advanced SQL with window functions, CTEs")
	fmt.Println("• MySQL: Standard SQL with some extensions")
	fmt.Println("• MongoDB: MongoDB Query Language (MQL) + Aggregation Pipeline")
	fmt.Println()

	fmt.Println("📈 Scaling:")
	fmt.Println("• PostgreSQL: Vertical + read replicas + partitioning")
	fmt.Println("• MySQL: Vertical + read replicas + sharding")
	fmt.Println("• MongoDB: Built-in horizontal scaling (sharding)")
	fmt.Println()

	fmt.Println("🎓 Learning Value:")
	fmt.Println("• PostgreSQL: Advanced SQL, ACID properties, complex queries")
	fmt.Println("• MySQL: SQL fundamentals, web application patterns")
	fmt.Println("• MongoDB: NoSQL concepts, document modeling, aggregations")

	// Return to database selection
	return selectDatabaseWithEducation(config)
}

// explainDatabaseChoice explains the chosen database
func explainDatabaseChoice(dbType string) error {
	fmt.Printf("\n🎉 Great choice: %s!\n", strings.ToUpper(dbType))
	fmt.Println()

	switch dbType {
	case "postgresql":
		fmt.Println("🐘 With PostgreSQL, you'll learn:")
		fmt.Println("• Advanced SQL queries and window functions")
		fmt.Println("• ACID transactions and data consistency")
		fmt.Println("• JSON/JSONB for flexible data storage")
		fmt.Println("• Full-text search and advanced indexing")
		fmt.Println("• Database migrations and schema evolution")
		fmt.Println()
		fmt.Println("🏗️  Repository Pattern Implementation:")
		fmt.Println("• SQL query builders and prepared statements")
		fmt.Println("• Transaction management in Go")
		fmt.Println("• Connection pooling and performance optimization")

	case "mysql":
		fmt.Println("🐬 With MySQL, you'll learn:")
		fmt.Println("• SQL fundamentals and best practices")
		fmt.Println("• Database design and normalization")
		fmt.Println("• Indexing strategies for performance")
		fmt.Println("• Replication and high availability")
		fmt.Println("• Web application database patterns")
		fmt.Println()
		fmt.Println("🏗️  Repository Pattern Implementation:")
		fmt.Println("• CRUD operations with proper error handling")
		fmt.Println("• Connection management and pooling")
		fmt.Println("• Query optimization techniques")

	case "mongodb":
		fmt.Println("🍃 With MongoDB, you'll learn:")
		fmt.Println("• Document-based data modeling")
		fmt.Println("• Flexible schema design patterns")
		fmt.Println("• Aggregation pipeline for complex queries")
		fmt.Println("• Indexing strategies for documents")
		fmt.Println("• Horizontal scaling concepts")
		fmt.Println()
		fmt.Println("🏗️  Repository Pattern Implementation:")
		fmt.Println("• Document CRUD operations")
		fmt.Println("• Aggregation queries in Go")
		fmt.Println("• Schema validation and data consistency")
	}

	return nil
}

// selectDatabaseConfigurationWithEducation handles database configuration selection
func selectDatabaseConfigurationWithEducation(config *ProjectConfiguration) error {
	fmt.Println("\n⚙️  Database Configuration Pattern:")
	fmt.Println("Choose how your application will connect to the database:")
	fmt.Println()

	configTypes := []struct {
		Name        string
		Type        string
		Description string
		UseCase     string
		Learning    string
	}{
		{
			Name:        "Single Instance",
			Type:        "single",
			Description: "One database server for all operations",
			UseCase:     "Development, small applications, simple deployments",
			Learning:    "Learn basic database connectivity and connection pooling",
		},
		{
			Name:        "Read-Write Split",
			Type:        "read-write",
			Description: "Separate read and write database endpoints",
			UseCase:     "High-read applications, performance optimization",
			Learning:    "Learn read replica patterns and query routing",
		},
		{
			Name:        "Cluster",
			Type:        "cluster",
			Description: "Multiple database nodes for high availability",
			UseCase:     "Production systems, high availability requirements",
			Learning:    "Learn distributed database patterns and failover",
		},
	}

	for i, cfg := range configTypes {
		fmt.Printf("%d. %s - %s\n", i+1, cfg.Name, cfg.Description)
		fmt.Printf("   🎯 Best for: %s\n", cfg.UseCase)
		fmt.Printf("   🎓 You'll learn: %s\n\n", cfg.Learning)
	}

	configOptions := []string{
		"Single instance - Simple single database server (recommended for learning)",
		"Read-Write split - Separate read and write endpoints",
		"Cluster - Multiple database nodes",
		"Quit",
	}

	var selected string
	configPrompt := &survey.Select{
		Message: "Choose your database configuration pattern:",
		Options: configOptions,
		Help:    "Start simple and scale up as you learn more patterns",
	}

	if err := survey.AskOne(configPrompt, &selected); err != nil {
		return err
	}

	if selected == "Quit" {
		return nil
	}

	// Extract configuration type
	switch {
	case strings.HasPrefix(selected, "Single"):
		config.DatabaseConfig.ConfigType = "single"
	case strings.HasPrefix(selected, "Read-Write"):
		config.DatabaseConfig.ConfigType = "read-write"
	case strings.HasPrefix(selected, "Cluster"):
		config.DatabaseConfig.ConfigType = "cluster"
	}

	// Get database credentials (simplified for educational purposes)
	return getDatabaseCredentialsWithEducation(config.DatabaseConfig, config.Name)
}

// getDatabaseCredentialsWithEducation gets database credentials with educational context
func getDatabaseCredentialsWithEducation(dbConfig *generator.DatabaseConfig, projectName string) error {
	fmt.Println("\n🔐 Database Connection Configuration:")
	fmt.Println("Let's configure your database connection details.")
	fmt.Println()

	fmt.Println("💡 Security Note:")
	fmt.Println("In production, never hardcode credentials in your code!")
	fmt.Println("We'll generate environment variable configuration for you.")
	fmt.Println()

	// Database name
	dbNamePrompt := &survey.Input{
		Message: "Database name:",
		Default: projectName + "_db",
		Help:    "The name of the database to connect to",
	}
	if err := survey.AskOne(dbNamePrompt, &dbConfig.DatabaseName, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// Username
	usernamePrompt := &survey.Input{
		Message: "Database username:",
		Default: "admin",
		Help:    "Database user with appropriate permissions",
	}
	if err := survey.AskOne(usernamePrompt, &dbConfig.Username, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// Password
	passwordPrompt := &survey.Password{
		Message: "Database password:",
		Help:    "This will be stored in environment variables, not in code",
	}
	if err := survey.AskOne(passwordPrompt, &dbConfig.Password, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// Host and port based on configuration type
	return configureConnectionDetails(dbConfig)
}

// configureConnectionDetails configures connection details based on configuration type
func configureConnectionDetails(dbConfig *generator.DatabaseConfig) error {
	switch dbConfig.ConfigType {
	case "single":
		return configureSingleInstance(dbConfig)
	case "read-write":
		return configureReadWriteSplit(dbConfig)
	case "cluster":
		return configureCluster(dbConfig)
	}
	return nil
}

// configureSingleInstance configures single instance connection
func configureSingleInstance(dbConfig *generator.DatabaseConfig) error {
	// Host
	hostPrompt := &survey.Input{
		Message: "Database host:",
		Default: "localhost",
		Help:    "Hostname or IP address of your database server",
	}
	if err := survey.AskOne(hostPrompt, &dbConfig.Host, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// Port
	var defaultPort string
	switch dbConfig.Type {
	case "postgresql":
		defaultPort = "5432"
	case "mysql":
		defaultPort = "3306"
	case "mongodb":
		defaultPort = "27017"
	}

	portPrompt := &survey.Input{
		Message: "Database port:",
		Default: defaultPort,
		Help:    "Port number for your database server",
	}
	if err := survey.AskOne(portPrompt, &dbConfig.Port, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// SSL Mode for SQL databases
	if dbConfig.Type == "postgresql" || dbConfig.Type == "mysql" {
		var sslMode string
		sslPrompt := &survey.Select{
			Message: "SSL Mode:",
			Options: []string{"disable", "require", "verify-ca", "verify-full"},
			Default: "disable",
			Help:    "SSL connection mode (use 'require' or higher in production)",
		}
		if err := survey.AskOne(sslPrompt, &sslMode); err != nil {
			return err
		}
		dbConfig.SSLMode = sslMode
	}

	return nil
}

// configureReadWriteSplit configures read-write split
func configureReadWriteSplit(dbConfig *generator.DatabaseConfig) error {
	fmt.Println("🔄 Read-Write Split Configuration:")
	fmt.Println("This pattern uses separate endpoints for read and write operations.")
	fmt.Println()

	// Write host
	writeHostPrompt := &survey.Input{
		Message: "Write database host (master):",
		Default: "localhost",
		Help:    "Primary database server for write operations",
	}
	if err := survey.AskOne(writeHostPrompt, &dbConfig.WriteHost, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// Read host
	readHostPrompt := &survey.Input{
		Message: "Read database host (replica):",
		Default: "localhost-replica",
		Help:    "Read replica server for read operations",
	}
	if err := survey.AskOne(readHostPrompt, &dbConfig.ReadHost, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// Port (same for both)
	var defaultPort string
	switch dbConfig.Type {
	case "postgresql":
		defaultPort = "5432"
	case "mysql":
		defaultPort = "3306"
	case "mongodb":
		defaultPort = "27017"
	}

	portPrompt := &survey.Input{
		Message: "Database port:",
		Default: defaultPort,
	}
	return survey.AskOne(portPrompt, &dbConfig.Port, survey.WithValidator(survey.Required))
}

// configureCluster configures cluster setup
func configureCluster(dbConfig *generator.DatabaseConfig) error {
	fmt.Println("🏗️  Cluster Configuration:")
	fmt.Println("This pattern uses multiple database nodes for high availability.")
	fmt.Println()

	// For simplicity, configure 3 nodes
	dbConfig.ClusterNodes = make([]string, 3)
	for i := 0; i < 3; i++ {
		nodePrompt := &survey.Input{
			Message: fmt.Sprintf("Cluster node %d host:", i+1),
			Default: fmt.Sprintf("db-node-%d.cluster.local", i+1),
			Help:    "Hostname of cluster node",
		}
		if err := survey.AskOne(nodePrompt, &dbConfig.ClusterNodes[i], survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	}

	// Port
	var defaultPort string
	switch dbConfig.Type {
	case "postgresql":
		defaultPort = "5432"
	case "mysql":
		defaultPort = "3306"
	case "mongodb":
		defaultPort = "27017"
	}

	portPrompt := &survey.Input{
		Message: "Database port:",
		Default: defaultPort,
	}
	return survey.AskOne(portPrompt, &dbConfig.Port, survey.WithValidator(survey.Required))
}

// configureRedisWithEducation handles Redis configuration with educational content
func configureRedisWithEducation(config *ProjectConfiguration) error {
	fmt.Println("\n🚀 Caching Layer Configuration:")
	fmt.Println("Redis provides high-performance caching and session storage.")
	fmt.Println()

	fmt.Println("🎓 Why use Redis?")
	fmt.Println("• Caching: Store frequently accessed data in memory")
	fmt.Println("• Session Storage: Manage user sessions across multiple servers")
	fmt.Println("• Pub/Sub: Real-time messaging between services")
	fmt.Println("• Rate Limiting: Control API usage and prevent abuse")
	fmt.Println()

	var redisChoice string
	redisPrompt := &survey.Select{
		Message: "Do you want to include Redis for caching and sessions?",
		Options: []string{
			"Yes - Include Redis support (recommended for APIs)",
			"No - Skip Redis for now",
			"Tell me more about Redis patterns",
			"Quit",
		},
		Help: "Redis adds powerful caching and session management capabilities",
	}

	if err := survey.AskOne(redisPrompt, &redisChoice); err != nil {
		return err
	}

	if redisChoice == "Quit" {
		return nil
	}

	if strings.HasPrefix(redisChoice, "Tell me more") {
		return explainRedisPatterns(config)
	}

	config.RedisConfig = &generator.RedisConfig{
		Enabled: strings.HasPrefix(redisChoice, "Yes"),
	}

	if config.RedisConfig.Enabled {
		return configureRedisConnection(config.RedisConfig)
	}

	return nil
}

// explainRedisPatterns explains Redis usage patterns
func explainRedisPatterns(config *ProjectConfiguration) error {
	fmt.Println("\n🔍 Redis Patterns You'll Learn:")
	fmt.Println()

	fmt.Println("1. 🏃 Caching Pattern:")
	fmt.Println("   • Cache database query results")
	fmt.Println("   • Reduce database load and improve response times")
	fmt.Println("   • Cache invalidation strategies")
	fmt.Println()

	fmt.Println("2. 🎫 Session Management:")
	fmt.Println("   • Store user sessions in Redis")
	fmt.Println("   • Enable horizontal scaling of your API")
	fmt.Println("   • Session expiration and cleanup")
	fmt.Println()

	fmt.Println("3. 🚦 Rate Limiting:")
	fmt.Println("   • Implement API rate limiting")
	fmt.Println("   • Prevent abuse and ensure fair usage")
	fmt.Println("   • Different rate limiting algorithms")
	fmt.Println()

	fmt.Println("4. 📡 Pub/Sub Messaging:")
	fmt.Println("   • Real-time notifications")
	fmt.Println("   • Event-driven architecture")
	fmt.Println("   • Microservice communication")
	fmt.Println()

	fmt.Println("💡 Code Example - Caching:")
	fmt.Println("```go")
	fmt.Println("// Check cache first")
	fmt.Println("user, err := redis.Get(\"user:\" + userID)")
	fmt.Println("if err == nil {")
	fmt.Println("    return user // Cache hit!")
	fmt.Println("}")
	fmt.Println()
	fmt.Println("// Cache miss - get from database")
	fmt.Println("user, err = db.GetUser(userID)")
	fmt.Println("if err == nil {")
	fmt.Println("    redis.Set(\"user:\" + userID, user, 5*time.Minute)")
	fmt.Println("}")
	fmt.Println("```")

	// Return to Redis configuration
	return configureRedisWithEducation(config)
}

// configureRedisConnection configures Redis connection details
func configureRedisConnection(redisConfig *generator.RedisConfig) error {
	fmt.Println("\n⚙️  Redis Connection Configuration:")
	fmt.Println()

	// Host
	hostPrompt := &survey.Input{
		Message: "Redis host:",
		Default: "localhost",
		Help:    "Hostname or IP address of your Redis server",
	}
	if err := survey.AskOne(hostPrompt, &redisConfig.Host, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// Port
	portPrompt := &survey.Input{
		Message: "Redis port:",
		Default: "6379",
		Help:    "Port number for your Redis server",
	}
	if err := survey.AskOne(portPrompt, &redisConfig.Port, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// Password (optional)
	passwordPrompt := &survey.Password{
		Message: "Redis password (leave empty if no password):",
		Help:    "Redis AUTH password (optional)",
	}
	survey.AskOne(passwordPrompt, &redisConfig.Password)

	// Database number
	redisConfig.Database = 0 // Default to database 0

	fmt.Println("✅ Redis configuration completed!")
	return nil
}

// configureProjectFeatures handles feature selection
func configureProjectFeatures(config *ProjectConfiguration) error {
	clearScreen()
	fmt.Println("🎛️  Step 5: Feature Configuration")
	fmt.Println("Choose additional features to include in your project:")
	fmt.Println()

	features := []ProjectFeature{
		{
			Name:        "JWT Authentication",
			Description: "JSON Web Token authentication and authorization",
			Educational: "Learn modern authentication patterns and security best practices",
		},
		{
			Name:        "API Documentation",
			Description: "Automatic OpenAPI/Swagger documentation generation",
			Educational: "Learn API documentation standards and tools",
		},
		{
			Name:        "Request Validation",
			Description: "Automatic request body and parameter validation",
			Educational: "Learn input validation patterns and security",
		},
		{
			Name:        "Structured Logging",
			Description: "JSON-structured logging with different levels",
			Educational: "Learn observability and debugging best practices",
		},
		{
			Name:        "Health Checks",
			Description: "Application and dependency health monitoring",
			Educational: "Learn monitoring and operational readiness patterns",
		},
		{
			Name:        "CORS Support",
			Description: "Cross-Origin Resource Sharing configuration",
			Educational: "Learn web security and browser interaction patterns",
		},
		{
			Name:        "Rate Limiting",
			Description: "API rate limiting and abuse prevention",
			Educational: "Learn API protection and fair usage patterns",
		},
	}

	for _, feature := range features {
		var include string
		includePrompt := &survey.Select{
			Message: fmt.Sprintf("Include %s?", feature.Name),
			Options: []string{
				fmt.Sprintf("Yes - %s", feature.Description),
				"No - Skip this feature",
				"What will I learn?",
				"Quit",
			},
		}

		if err := survey.AskOne(includePrompt, &include); err != nil {
			return err
		}

		if include == "Quit" {
			return nil
		}

		if strings.HasPrefix(include, "What will") {
			fmt.Printf("\n🎓 %s:\n%s\n\n", feature.Name, feature.Educational)
			// Ask again
			includePrompt := &survey.Select{
				Message: fmt.Sprintf("Include %s?", feature.Name),
				Options: []string{
					fmt.Sprintf("Yes - %s", feature.Description),
					"No - Skip this feature",
				},
			}
			if err := survey.AskOne(includePrompt, &include); err != nil {
				return err
			}
		}

		feature.Enabled = strings.HasPrefix(include, "Yes")
		config.Features = append(config.Features, feature)
	}

	// Show selected features
	fmt.Println("\n✅ Selected Features:")
	for _, feature := range config.Features {
		if feature.Enabled {
			fmt.Printf("   • %s\n", feature.Name)
		}
	}

	return nil
}

// visualizeProjectStructure shows the project structure that will be generated
func visualizeProjectStructure(config *ProjectConfiguration) error {
	clearScreen()
	fmt.Println("🏗️  Step 6: Project Structure Visualization")
	fmt.Printf("Here's the structure that will be generated for your %s project:\n\n", config.Type)

	switch config.Type {
	case "api":
		fmt.Println("📁 Project Structure:")
		fmt.Println("```")
		fmt.Printf("%s/\n", config.Name)
		fmt.Println("├── cmd/")
		fmt.Println("│   └── api/")
		fmt.Println("│       └── main.go              # Application entry point")
		fmt.Println("├── internal/")
		fmt.Println("│   ├── api/                     # 🌐 Interface Layer")
		fmt.Println("│   │   ├── handlers/            # HTTP request handlers")
		fmt.Println("│   │   ├── middleware/          # Cross-cutting concerns")
		fmt.Println("│   │   └── routes/              # Route definitions")
		fmt.Println("│   ├── domain/                  # 🏛️  Domain Layer")
		fmt.Println("│   │   ├── user/                # User business logic")
		fmt.Println("│   │   └── post/                # Post business logic")
		fmt.Println("│   ├── infrastructure/          # 🔧 Infrastructure Layer")
		fmt.Println("│   │   ├── database/            # Database implementations")
		fmt.Println("│   │   └── auth/                # Authentication services")
		fmt.Println("│   ├── config/                  # Configuration management")
		fmt.Println("│   └── pkg/                     # Shared utilities")
		fmt.Println("├── migrations/                  # Database migrations")
		fmt.Println("├── scripts/                     # Build and deployment scripts")
		fmt.Println("├── .env.example                 # Environment variables template")
		fmt.Println("├── go.mod                       # Go module definition")
		fmt.Println("├── README.md                    # Project documentation")
		fmt.Println("└── gophex.md                    # Project metadata")
		fmt.Println("```")

		fmt.Printf("\n🚀 Framework: %s\n", strings.ToUpper(config.Framework))
		fmt.Printf("🗄️  Database: %s (%s)\n", strings.ToUpper(config.DatabaseConfig.Type), config.DatabaseConfig.ConfigType)
		if config.RedisConfig.Enabled {
			fmt.Println("🚀 Caching: Redis enabled")
		}

	case "webapp":
		fmt.Println("📁 Project Structure:")
		fmt.Println("```")
		fmt.Printf("%s/\n", config.Name)
		fmt.Println("├── cmd/")
		fmt.Println("│   └── webapp/")
		fmt.Println("│       └── main.go              # Application entry point")
		fmt.Println("├── internal/")
		fmt.Println("│   ├── handlers/                # HTTP handlers for pages")
		fmt.Println("│   ├── models/                  # Data models")
		fmt.Println("│   ├── middleware/              # Web middleware")
		fmt.Println("│   └── config/                  # Configuration")
		fmt.Println("├── web/")
		fmt.Println("│   ├── templates/               # HTML templates")
		fmt.Println("│   └── static/                  # CSS, JS, images")
		fmt.Println("├── go.mod")
		fmt.Println("├── README.md")
		fmt.Println("└── gophex.md")
		fmt.Println("```")

	case "cli":
		fmt.Println("📁 Project Structure:")
		fmt.Println("```")
		fmt.Printf("%s/\n", config.Name)
		fmt.Println("├── cmd/")
		fmt.Println("│   └── main.go                  # CLI entry point")
		fmt.Println("├── internal/")
		fmt.Println("│   └── cmd/                     # Command implementations")
		fmt.Println("│       ├── root.go              # Root command")
		fmt.Println("│       └── version.go           # Version command")
		fmt.Println("├── go.mod")
		fmt.Println("├── README.md")
		fmt.Println("└── gophex.md")
		fmt.Println("```")

	case "microservice":
		fmt.Println("📁 Project Structure:")
		fmt.Println("```")
		fmt.Printf("%s/\n", config.Name)
		fmt.Println("├── cmd/")
		fmt.Println("│   └── server/")
		fmt.Println("│       └── main.go              # Service entry point")
		fmt.Println("├── internal/")
		fmt.Println("│   ├── handlers/                # gRPC handlers")
		fmt.Println("│   ├── config/                  # Configuration")
		fmt.Println("│   └── health/                  # Health checks")
		fmt.Println("├── proto/                       # Protocol buffer definitions")
		fmt.Println("├── go.mod")
		fmt.Println("├── README.md")
		fmt.Println("└── gophex.md")
		fmt.Println("```")
	}

	fmt.Println("\n🎓 Educational Features:")
	fmt.Println("• Comprehensive code comments explaining patterns")
	fmt.Println("• Clean Architecture principles demonstrated")
	fmt.Println("• Best practices for Go development")
	fmt.Println("• Example implementations and tests")
	fmt.Println("• Step-by-step learning documentation")

	var proceed string
	proceedPrompt := &survey.Select{
		Message: "Ready to generate your project?",
		Options: []string{
			"Yes - Generate project with educational content",
			"No - Modify configuration",
			"Quit",
		},
	}

	if err := survey.AskOne(proceedPrompt, &proceed); err != nil {
		return err
	}

	if proceed == "Quit" {
		return nil
	}

	if strings.HasPrefix(proceed, "No") {
		return ErrReturnToMenu
	}

	return nil
}

// generateProjectWithExplanation generates the project and explains what was created
func generateProjectWithExplanation(config *ProjectConfiguration) error {
	clearScreen()
	fmt.Println("🚀 Step 7: Project Generation")
	fmt.Printf("Generating your %s project with educational content...\n", config.Type)
	fmt.Println()

	// Generate the project
	gen := generator.New()
	var err error
	if config.Type == "api" {
		err = gen.GenerateWithFramework(config.Type, config.Name, config.Path, config.Framework, config.DatabaseConfig, config.RedisConfig)
	} else {
		err = gen.Generate(config.Type, config.Name, config.Path)
	}

	if err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	// Create project tracking metadata
	tracker := NewProjectTracker(config.Path)
	if err := tracker.CreateInitialMetadata(config.Type, config.Name, config.Path, config.DatabaseConfig, config.RedisConfig); err != nil {
		fmt.Printf("⚠️  Warning: Failed to create project tracking metadata: %v\n", err)
	}

	fmt.Printf("✅ Successfully generated %s project '%s'!\n", config.Type, config.Name)
	fmt.Printf("📍 Location: %s\n\n", config.Path)

	// Explain what was generated
	return explainGeneratedProject(config)
}

// explainGeneratedProject explains what was generated and next steps
func explainGeneratedProject(config *ProjectConfiguration) error {
	fmt.Println("🎉 What Was Generated:")
	fmt.Println()

	switch config.Type {
	case "api":
		fmt.Println("🏛️  Clean Architecture Implementation:")
		fmt.Printf("• Domain entities (User, Post) with business rules\n")
		fmt.Printf("• Repository interfaces and %s implementations\n", config.DatabaseConfig.Type)
		fmt.Printf("• Service layer with use cases and business logic\n")
		fmt.Printf("• %s HTTP handlers with proper error handling\n", strings.ToUpper(config.Framework))
		fmt.Printf("• Middleware for logging, CORS, and authentication\n")
		fmt.Printf("• Database migrations for %s\n", config.DatabaseConfig.Type)
		if config.RedisConfig.Enabled {
			fmt.Printf("• Redis integration for caching and sessions\n")
		}

		fmt.Println("\n🎓 Educational Features:")
		fmt.Println("• Extensive code comments explaining Clean Architecture")
		fmt.Println("• Business rule examples and validation patterns")
		fmt.Println("• Dependency injection setup and explanation")
		fmt.Println("• Testing examples for each layer")
		fmt.Println("• API documentation and usage examples")

	case "webapp":
		fmt.Println("🌐 Web Application Features:")
		fmt.Println("• MVC pattern implementation")
		fmt.Println("• HTML template rendering")
		fmt.Println("• Static asset serving")
		fmt.Println("• Session management")
		fmt.Println("• Form handling examples")

	case "cli":
		fmt.Println("💻 CLI Tool Features:")
		fmt.Println("• Command pattern with subcommands")
		fmt.Println("• Flag parsing and validation")
		fmt.Println("• Configuration management")
		fmt.Println("• Cross-platform compatibility")

	case "microservice":
		fmt.Println("🔧 Microservice Features:")
		fmt.Println("• gRPC service implementation")
		fmt.Println("• Health check endpoints")
		fmt.Println("• Configuration management")
		fmt.Println("• Docker containerization")
	}

	fmt.Println("\n📚 Next Steps:")
	fmt.Println("1. 📖 Read the generated README.md for detailed instructions")
	fmt.Println("2. 🔧 Install dependencies: cd " + config.Name + " && go mod tidy")
	if config.Type == "api" {
		fmt.Println("3. 🗄️  Set up your database and run migrations")
		fmt.Println("4. ⚙️  Configure environment variables (.env file)")
		fmt.Println("5. 🚀 Start the server: go run cmd/api/main.go")
		fmt.Println("6. 🧪 Test the API endpoints")
		fmt.Println("7. 🎓 Use the Enhanced CRUD Wizard to add more entities")
	} else {
		fmt.Println("3. 🚀 Build and run: go run cmd/*/main.go")
		fmt.Println("4. 🧪 Run tests: go test ./...")
	}

	fmt.Println("\n🎯 Learning Path:")
	fmt.Println("• Study the generated code and comments")
	fmt.Println("• Experiment with modifications")
	fmt.Println("• Add new features using the same patterns")
	fmt.Println("• Run tests to understand the architecture")

	// Show post-generation menu
	opts := PostGenerationOptions{
		ProjectPath: config.Path,
		ProjectType: config.Type,
		ProjectName: config.Name,
	}

	return ShowPostGenerationMenu(opts)
}
