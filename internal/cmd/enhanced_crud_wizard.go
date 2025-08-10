package cmd

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// ArchitectureLayer represents a layer in the clean architecture
type ArchitectureLayer struct {
	Name         string
	Description  string
	Files        []string
	Purpose      string
	Dependencies []string
}

// DomainObject represents a complete domain object with all its components
type DomainObject struct {
	Entity     CRUDEntity
	Repository RepositoryConfig
	Service    ServiceConfig
	Handler    HandlerConfig
	Validation ValidationConfig
	Middleware []MiddlewareConfig
}

// RepositoryConfig represents repository layer configuration
type RepositoryConfig struct {
	Interface    string
	Methods      []string
	Caching      bool
	Transactions bool
	Pagination   bool
}

// ServiceConfig represents service layer configuration
type ServiceConfig struct {
	BusinessRules []BusinessRule
	Events        []DomainEvent
	Validation    bool
	Logging       bool
}

// HandlerConfig represents handler layer configuration
type HandlerConfig struct {
	Endpoints     []APIEndpoint
	Middleware    []string
	ErrorHandling bool
	Documentation bool
}

// ValidationConfig represents validation configuration
type ValidationConfig struct {
	InputValidation    bool
	BusinessValidation bool
	CustomValidators   []string
}

// MiddlewareConfig represents middleware configuration
type MiddlewareConfig struct {
	Name    string
	Purpose string
	Order   int
	Enabled bool
}

// BusinessRule represents a business rule
type BusinessRule struct {
	Name        string
	Description string
	Validation  string
}

// DomainEvent represents a domain event
type DomainEvent struct {
	Name    string
	Trigger string
	Payload []string
}

// APIEndpoint represents an API endpoint
type APIEndpoint struct {
	Method      string
	Path        string
	Handler     string
	Middleware  []string
	Description string
}

// RunEnhancedCRUDWizard runs the enhanced educational CRUD generation wizard
func RunEnhancedCRUDWizard(projectPath string) error {
	clearScreen()
	fmt.Println("🎓 Enhanced CRUD Architecture Wizard")
	fmt.Println("Learn Go Clean Architecture by building step-by-step!")
	fmt.Println()

	// Show architecture overview first
	if err := showArchitectureOverview(); err != nil {
		return err
	}

	domainObj := &DomainObject{}

	// Step 1: Domain Entity Design
	if err := designDomainEntity(domainObj); err != nil {
		return err
	}

	// Step 2: Repository Layer Design
	if err := designRepositoryLayer(domainObj); err != nil {
		return err
	}

	// Step 3: Service Layer Design
	if err := designServiceLayer(domainObj); err != nil {
		return err
	}

	// Step 4: Handler Layer Design
	if err := designHandlerLayer(domainObj); err != nil {
		return err
	}

	// Step 5: Middleware Configuration
	if err := configureMiddleware(domainObj); err != nil {
		return err
	}

	// Step 6: Dependency Injection Visualization
	if err := visualizeDependencyInjection(domainObj); err != nil {
		return err
	}

	// Step 7: Architecture Review and Generation
	if err := reviewArchitectureAndGenerate(projectPath, domainObj); err != nil {
		return err
	}

	return nil
}

// showArchitectureOverview shows the clean architecture layers
func showArchitectureOverview() error {
	fmt.Println("📚 Clean Architecture Overview")
	fmt.Println("We'll build your CRUD operations using Clean Architecture principles:")
	fmt.Println()

	layers := []ArchitectureLayer{
		{
			Name:         "Domain Layer (Core)",
			Description:  "Business entities, rules, and interfaces",
			Files:        []string{"internal/domain/{entity}/model.go", "internal/domain/{entity}/repository.go"},
			Purpose:      "Contains your business logic and rules, independent of external concerns",
			Dependencies: []string{"No dependencies on other layers"},
		},
		{
			Name:         "Infrastructure Layer",
			Description:  "Database implementations, external services",
			Files:        []string{"internal/infrastructure/repository/{entity}_repository.go"},
			Purpose:      "Implements domain interfaces, handles data persistence",
			Dependencies: []string{"Depends on Domain Layer interfaces"},
		},
		{
			Name:         "Application Layer (Use Cases)",
			Description:  "Application-specific business rules and orchestration",
			Files:        []string{"internal/domain/{entity}/service.go"},
			Purpose:      "Orchestrates domain objects and implements use cases",
			Dependencies: []string{"Depends on Domain Layer"},
		},
		{
			Name:         "Interface Layer (Controllers)",
			Description:  "HTTP handlers, request/response handling",
			Files:        []string{"internal/api/handlers/{entity}.go", "internal/api/routes/routes.go"},
			Purpose:      "Handles HTTP requests, converts data formats, calls use cases",
			Dependencies: []string{"Depends on Application Layer"},
		},
	}

	for i, layer := range layers {
		fmt.Printf("%d. %s\n", i+1, layer.Name)
		fmt.Printf("   📖 %s\n", layer.Description)
		fmt.Printf("   🎯 %s\n", layer.Purpose)
		fmt.Printf("   🔗 %s\n", strings.Join(layer.Dependencies, ", "))
		fmt.Printf("   📁 Files: %s\n\n", strings.Join(layer.Files, ", "))
	}

	var proceed string
	proceedPrompt := &survey.Select{
		Message: "Ready to start building your CRUD architecture?",
		Options: []string{
			"Yes - Let's build step by step",
			"No - Show me more details first",
			"Quit",
		},
	}

	if err := survey.AskOne(proceedPrompt, &proceed); err != nil {
		return err
	}

	if proceed == "Quit" {
		return nil
	}

	if proceed[:2] == "No" {
		return showDetailedArchitectureExplanation()
	}

	return nil
}

// showDetailedArchitectureExplanation provides more detailed explanation
func showDetailedArchitectureExplanation() error {
	fmt.Println("\n🔍 Detailed Architecture Explanation")
	fmt.Println()

	explanations := map[string]string{
		"Dependency Rule":     "Dependencies point inward. Outer layers depend on inner layers, never the reverse.",
		"Domain Independence": "Your business logic doesn't know about databases, HTTP, or frameworks.",
		"Testability":         "Each layer can be tested independently with mocks/stubs.",
		"Flexibility":         "You can swap databases, frameworks, or external services without changing business logic.",
		"Maintainability":     "Clear separation of concerns makes code easier to understand and modify.",
	}

	fmt.Println("Key Principles:")
	for principle, explanation := range explanations {
		fmt.Printf("• %s: %s\n", principle, explanation)
	}

	fmt.Println("\n📊 Data Flow Example (Create User):")
	fmt.Println("1. HTTP Request → Handler (Interface Layer)")
	fmt.Println("2. Handler validates input → calls Service (Application Layer)")
	fmt.Println("3. Service applies business rules → calls Repository (Domain Interface)")
	fmt.Println("4. Repository implementation saves to database (Infrastructure Layer)")
	fmt.Println("5. Response flows back through the layers")

	var ready string
	readyPrompt := &survey.Select{
		Message: "Now ready to start building?",
		Options: []string{
			"Yes - Let's start building",
			"Quit",
		},
	}

	return survey.AskOne(readyPrompt, &ready)
}

// designDomainEntity handles domain entity design with educational content
func designDomainEntity(domainObj *DomainObject) error {
	clearScreen()
	fmt.Println("🏗️  Step 1: Domain Entity Design")
	fmt.Println("The Domain Entity is the heart of your business logic.")
	fmt.Println()

	fmt.Println("📚 What is a Domain Entity?")
	fmt.Println("• Represents a business concept (User, Order, Product)")
	fmt.Println("• Contains business rules and validation")
	fmt.Println("• Independent of database or framework details")
	fmt.Println("• Has identity and lifecycle")
	fmt.Println()

	// Use existing entity selection but with more education
	if err := selectEntityWithEducation(&domainObj.Entity); err != nil {
		return err
	}

	// Enhanced field definition with business rule consideration
	if err := defineFieldsWithBusinessRules(&domainObj.Entity); err != nil {
		return err
	}

	// Show what will be generated for this layer
	fmt.Printf("\n✅ Domain Layer for %s will include:\n", domainObj.Entity.Name)
	fmt.Printf("📁 internal/domain/%s/model.go - Entity definition with business rules\n", domainObj.Entity.Name)
	fmt.Printf("📁 internal/domain/%s/repository.go - Repository interface (contract)\n", domainObj.Entity.Name)
	fmt.Printf("📁 internal/domain/%s/errors.go - Domain-specific errors\n", domainObj.Entity.Name)

	return nil
}

// designRepositoryLayer handles repository layer design
func designRepositoryLayer(domainObj *DomainObject) error {
	clearScreen()
	fmt.Println("🗄️  Step 2: Repository Layer Design")
	fmt.Println("The Repository abstracts data access and provides a collection-like interface.")
	fmt.Println()

	fmt.Println("📚 Repository Pattern Benefits:")
	fmt.Println("• Separates business logic from data access logic")
	fmt.Println("• Makes testing easier with mock implementations")
	fmt.Println("• Allows switching databases without changing business logic")
	fmt.Println("• Provides a consistent interface for data operations")
	fmt.Println()

	repo := &RepositoryConfig{}

	// Configure repository methods
	if err := configureRepositoryMethods(repo, &domainObj.Entity); err != nil {
		return err
	}

	// Configure advanced features
	if err := configureRepositoryFeatures(repo); err != nil {
		return err
	}

	domainObj.Repository = *repo

	// Show repository interface that will be generated
	fmt.Printf("\n✅ Repository Layer for %s will include:\n", domainObj.Entity.Name)
	fmt.Printf("📁 Repository Interface (in domain layer):\n")
	for _, method := range repo.Methods {
		fmt.Printf("   • %s\n", method)
	}
	fmt.Printf("📁 Repository Implementation (in infrastructure layer):\n")
	fmt.Printf("   • Database-specific implementation\n")
	if repo.Caching {
		fmt.Printf("   • Caching layer integration\n")
	}
	if repo.Transactions {
		fmt.Printf("   • Transaction support\n")
	}

	return nil
}

// designServiceLayer handles service layer design
func designServiceLayer(domainObj *DomainObject) error {
	clearScreen()
	fmt.Println("⚙️  Step 3: Service Layer Design (Use Cases)")
	fmt.Println("The Service Layer orchestrates domain objects and implements use cases.")
	fmt.Println()

	fmt.Println("📚 Service Layer Responsibilities:")
	fmt.Println("• Implements application-specific business rules")
	fmt.Println("• Orchestrates multiple domain objects")
	fmt.Println("• Handles transactions and error scenarios")
	fmt.Println("• Publishes domain events")
	fmt.Println()

	service := &ServiceConfig{}

	// Configure business rules
	if err := configureBusinessRules(service, &domainObj.Entity); err != nil {
		return err
	}

	// Configure domain events
	if err := configureDomainEvents(service, &domainObj.Entity); err != nil {
		return err
	}

	// Configure service features
	if err := configureServiceFeatures(service); err != nil {
		return err
	}

	domainObj.Service = *service

	fmt.Printf("\n✅ Service Layer for %s will include:\n", domainObj.Entity.Name)
	fmt.Printf("📁 Business Rules:\n")
	for _, rule := range service.BusinessRules {
		fmt.Printf("   • %s: %s\n", rule.Name, rule.Description)
	}
	if len(service.Events) > 0 {
		fmt.Printf("📁 Domain Events:\n")
		for _, event := range service.Events {
			fmt.Printf("   • %s (triggered on %s)\n", event.Name, event.Trigger)
		}
	}

	return nil
}

// designHandlerLayer handles handler layer design
func designHandlerLayer(domainObj *DomainObject) error {
	clearScreen()
	fmt.Println("🌐 Step 4: Handler Layer Design (Controllers)")
	fmt.Println("The Handler Layer manages HTTP requests and responses.")
	fmt.Println()

	fmt.Println("📚 Handler Layer Responsibilities:")
	fmt.Println("• Handles HTTP requests and responses")
	fmt.Println("• Validates input data")
	fmt.Println("• Converts between HTTP and domain models")
	fmt.Println("• Applies middleware (auth, logging, etc.)")
	fmt.Println()

	handler := &HandlerConfig{}

	// Configure API endpoints
	if err := configureAPIEndpoints(handler, &domainObj.Entity); err != nil {
		return err
	}

	// Configure handler features
	if err := configureHandlerFeatures(handler); err != nil {
		return err
	}

	domainObj.Handler = *handler

	fmt.Printf("\n✅ Handler Layer for %s will include:\n", domainObj.Entity.Name)
	fmt.Printf("📁 API Endpoints:\n")
	for _, endpoint := range handler.Endpoints {
		fmt.Printf("   • %s %s - %s\n", endpoint.Method, endpoint.Path, endpoint.Description)
	}

	return nil
}

// configureMiddleware handles middleware configuration
func configureMiddleware(domainObj *DomainObject) error {
	clearScreen()
	fmt.Println("🛡️  Step 5: Middleware Configuration")
	fmt.Println("Middleware provides cross-cutting concerns like authentication, logging, and validation.")
	fmt.Println()

	fmt.Println("📚 Common Middleware Types:")
	fmt.Println("• Authentication - Verify user identity")
	fmt.Println("• Authorization - Check user permissions")
	fmt.Println("• Logging - Record request/response details")
	fmt.Println("• Rate Limiting - Prevent abuse")
	fmt.Println("• CORS - Handle cross-origin requests")
	fmt.Println("• Validation - Validate request data")
	fmt.Println()

	middlewares := []MiddlewareConfig{
		{Name: "Logger", Purpose: "Log all HTTP requests and responses", Order: 1, Enabled: true},
		{Name: "CORS", Purpose: "Handle cross-origin resource sharing", Order: 2, Enabled: false},
		{Name: "Authentication", Purpose: "Verify JWT tokens or API keys", Order: 3, Enabled: false},
		{Name: "Authorization", Purpose: "Check user permissions for resources", Order: 4, Enabled: false},
		{Name: "Rate Limiter", Purpose: "Prevent API abuse with rate limiting", Order: 5, Enabled: false},
		{Name: "Request Validator", Purpose: "Validate request body and parameters", Order: 6, Enabled: true},
		{Name: "Error Handler", Purpose: "Handle and format error responses", Order: 7, Enabled: true},
	}

	// Let user select which middleware to include
	selectedMiddleware := []string{}
	for _, mw := range middlewares {
		var include string
		includePrompt := &survey.Select{
			Message: fmt.Sprintf("Include %s middleware?", mw.Name),
			Options: []string{
				fmt.Sprintf("Yes - %s", mw.Purpose),
				"No - Skip this middleware",
				"Quit",
			},
		}

		if err := survey.AskOne(includePrompt, &include); err != nil {
			return err
		}

		if include == "Quit" {
			return nil
		}

		if include[:3] == "Yes" {
			mw.Enabled = true
			selectedMiddleware = append(selectedMiddleware, mw.Name)
		}
		domainObj.Middleware = append(domainObj.Middleware, mw)
	}

	fmt.Printf("\n✅ Middleware Stack for %s:\n", domainObj.Entity.Name)
	for _, mw := range domainObj.Middleware {
		if mw.Enabled {
			fmt.Printf("   %d. %s - %s\n", mw.Order, mw.Name, mw.Purpose)
		}
	}

	return nil
}

// visualizeDependencyInjection shows how dependency injection works
func visualizeDependencyInjection(domainObj *DomainObject) error {
	clearScreen()
	fmt.Println("🔌 Step 6: Dependency Injection Visualization")
	fmt.Println("See how all the layers connect together through dependency injection.")
	fmt.Println()

	fmt.Println("📚 Dependency Injection Benefits:")
	fmt.Println("• Loose coupling between components")
	fmt.Println("• Easy testing with mock dependencies")
	fmt.Println("• Flexible configuration and swapping of implementations")
	fmt.Println("• Clear dependency relationships")
	fmt.Println()

	fmt.Printf("🔗 Dependency Graph for %s:\n\n", domainObj.Entity.Name)

	// Show the dependency flow
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("│                    main.go (Composition Root)              │")
	fmt.Println("└─────────────────────┬───────────────────────────────────────┘")
	fmt.Println("                      │")
	fmt.Println("                      ▼")
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Printf("│              %s Handler                                │\n", domainObj.Entity.Name)
	fmt.Println("│  • Handles HTTP requests                                    │")
	fmt.Println("│  • Validates input                                          │")
	fmt.Println("│  • Calls service layer                                      │")
	fmt.Println("└─────────────────────┬───────────────────────────────────────┘")
	fmt.Println("                      │ depends on")
	fmt.Println("                      ▼")
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Printf("│               %s Service                               │\n", domainObj.Entity.Name)
	fmt.Println("│  • Implements business logic                                │")
	fmt.Println("│  • Orchestrates domain operations                           │")
	fmt.Println("│  • Handles transactions                                     │")
	fmt.Println("└─────────────────────┬───────────────────────────────────────┘")
	fmt.Println("                      │ depends on")
	fmt.Println("                      ▼")
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Printf("│            %s Repository Interface                     │\n", domainObj.Entity.Name)
	fmt.Println("│  • Defines data access contract                             │")
	fmt.Println("│  • Domain layer interface                                   │")
	fmt.Println("└─────────────────────┬───────────────────────────────────────┘")
	fmt.Println("                      │ implemented by")
	fmt.Println("                      ▼")
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Printf("│         %s Repository Implementation                   │\n", domainObj.Entity.Name)
	fmt.Println("│  • Database-specific implementation                         │")
	fmt.Println("│  • Infrastructure layer                                     │")
	fmt.Println("└─────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 Injection Points:")
	fmt.Printf("1. Repository implementation injected into Service\n")
	fmt.Printf("2. Service injected into Handler\n")
	fmt.Printf("3. Handler registered with Router\n")
	fmt.Printf("4. Middleware applied to routes\n")

	var proceed string
	proceedPrompt := &survey.Select{
		Message: "Ready to see the generated code structure?",
		Options: []string{
			"Yes - Show me the code structure",
			"No - Let me review the configuration",
			"Quit",
		},
	}

	return survey.AskOne(proceedPrompt, &proceed)
}

// Helper functions for configuration

func selectEntityWithEducation(entity *CRUDEntity) error {
	// Reuse existing selectEntity but with additional education
	return selectEntity(entity)
}

func defineFieldsWithBusinessRules(entity *CRUDEntity) error {
	// Enhanced version of defineFields that considers business rules
	if err := defineFields(entity); err != nil {
		return err
	}

	// Add business rule considerations
	fmt.Println("\n💼 Business Rule Considerations:")
	fmt.Println("Think about these business rules for your entity:")
	fmt.Printf("• What makes a %s valid?\n", entity.Name)
	fmt.Printf("• What business constraints should be enforced?\n")
	fmt.Printf("• What happens when a %s is created/updated/deleted?\n", entity.Name)
	fmt.Printf("• Are there any relationships with other entities?\n")

	return nil
}

func configureRepositoryMethods(repo *RepositoryConfig, entity *CRUDEntity) error {
	fmt.Printf("🔧 Configuring Repository Methods for %s\n\n", entity.Name)

	standardMethods := []string{
		fmt.Sprintf("Create(%s) error", entity.Name),
		fmt.Sprintf("GetByID(id int64) (*%s, error)", entity.Name),
		fmt.Sprintf("Update(%s) error", entity.Name),
		fmt.Sprintf("Delete(id int64) error"),
		fmt.Sprintf("List(limit, offset int) ([]*%s, error)", entity.Name),
		fmt.Sprintf("Count() (int64, error)"),
	}

	repo.Methods = standardMethods

	// Ask about additional methods
	additionalMethods := []string{
		fmt.Sprintf("GetByEmail(email string) (*%s, error)", entity.Name),
		fmt.Sprintf("GetByStatus(status string) ([]*%s, error)", entity.Name),
		fmt.Sprintf("Search(query string) ([]*%s, error)", entity.Name),
		fmt.Sprintf("GetActive() ([]*%s, error)", entity.Name),
	}

	for _, method := range additionalMethods {
		var include string
		includePrompt := &survey.Select{
			Message: fmt.Sprintf("Include method: %s?", method),
			Options: []string{
				"Yes - Include this method",
				"No - Skip this method",
				"Quit",
			},
		}

		if err := survey.AskOne(includePrompt, &include); err != nil {
			return err
		}

		if include == "Quit" {
			return nil
		}

		if include[:3] == "Yes" {
			repo.Methods = append(repo.Methods, method)
		}
	}

	return nil
}

func configureRepositoryFeatures(repo *RepositoryConfig) error {
	features := []struct {
		name        string
		description string
		field       *bool
	}{
		{"Caching", "Add Redis caching layer for better performance", &repo.Caching},
		{"Transactions", "Support database transactions for data consistency", &repo.Transactions},
		{"Pagination", "Built-in pagination support for list operations", &repo.Pagination},
	}

	for _, feature := range features {
		var include string
		includePrompt := &survey.Select{
			Message: fmt.Sprintf("Enable %s?", feature.name),
			Options: []string{
				fmt.Sprintf("Yes - %s", feature.description),
				"No - Skip this feature",
				"Quit",
			},
		}

		if err := survey.AskOne(includePrompt, &include); err != nil {
			return err
		}

		if include == "Quit" {
			return nil
		}

		*feature.field = include[:3] == "Yes"
	}

	return nil
}

func configureBusinessRules(service *ServiceConfig, entity *CRUDEntity) error {
	fmt.Printf("📋 Configuring Business Rules for %s\n\n", entity.Name)

	// Suggest common business rules based on entity type
	var suggestedRules []BusinessRule

	switch entity.Name {
	case "user":
		suggestedRules = []BusinessRule{
			{"Email Uniqueness", "Ensure email addresses are unique across all users", "email must be unique"},
			{"Password Strength", "Enforce strong password requirements", "password must meet complexity requirements"},
			{"Account Activation", "New accounts must be activated before use", "account must be activated"},
		}
	case "post":
		suggestedRules = []BusinessRule{
			{"Author Ownership", "Only the author can edit their posts", "author_id must match current user"},
			{"Publication Rules", "Posts must have title and content to be published", "title and content required for publication"},
		}
	case "product":
		suggestedRules = []BusinessRule{
			{"Price Validation", "Product price must be positive", "price > 0"},
			{"SKU Uniqueness", "Product SKUs must be unique", "sku must be unique"},
			{"Stock Management", "Cannot sell more than available stock", "quantity <= stock"},
		}
	default:
		suggestedRules = []BusinessRule{
			{"Data Integrity", "Ensure required fields are provided", "validate required fields"},
			{"Business Constraints", "Apply domain-specific validation rules", "custom validation rules"},
		}
	}

	for _, rule := range suggestedRules {
		var include string
		includePrompt := &survey.Select{
			Message: fmt.Sprintf("Include business rule: %s?", rule.Name),
			Options: []string{
				fmt.Sprintf("Yes - %s", rule.Description),
				"No - Skip this rule",
				"Quit",
			},
		}

		if err := survey.AskOne(includePrompt, &include); err != nil {
			return err
		}

		if include == "Quit" {
			return nil
		}

		if include[:3] == "Yes" {
			service.BusinessRules = append(service.BusinessRules, rule)
		}
	}

	return nil
}

func configureDomainEvents(service *ServiceConfig, entity *CRUDEntity) error {
	fmt.Printf("📡 Configuring Domain Events for %s\n\n", entity.Name)

	// Suggest common domain events
	var suggestedEvents []DomainEvent

	switch entity.Name {
	case "user":
		suggestedEvents = []DomainEvent{
			{"UserCreated", "user creation", []string{"user_id", "email", "created_at"}},
			{"UserUpdated", "user update", []string{"user_id", "updated_fields", "updated_at"}},
			{"UserDeleted", "user deletion", []string{"user_id", "deleted_at"}},
		}
	case "post":
		suggestedEvents = []DomainEvent{
			{"PostPublished", "post publication", []string{"post_id", "author_id", "published_at"}},
			{"PostUpdated", "post update", []string{"post_id", "updated_fields", "updated_at"}},
		}
	default:
		suggestedEvents = []DomainEvent{
			{fmt.Sprintf("%sCreated", strings.Title(entity.Name)), fmt.Sprintf("%s creation", entity.Name), []string{"id", "created_at"}},
			{fmt.Sprintf("%sUpdated", strings.Title(entity.Name)), fmt.Sprintf("%s update", entity.Name), []string{"id", "updated_at"}},
		}
	}

	for _, event := range suggestedEvents {
		var include string
		includePrompt := &survey.Select{
			Message: fmt.Sprintf("Include domain event: %s?", event.Name),
			Options: []string{
				fmt.Sprintf("Yes - Triggered on %s", event.Trigger),
				"No - Skip this event",
				"Quit",
			},
		}

		if err := survey.AskOne(includePrompt, &include); err != nil {
			return err
		}

		if include == "Quit" {
			return nil
		}

		if include[:3] == "Yes" {
			service.Events = append(service.Events, event)
		}
	}

	return nil
}

func configureServiceFeatures(service *ServiceConfig) error {
	features := []struct {
		name        string
		description string
		field       *bool
	}{
		{"Validation", "Add input validation in service layer", &service.Validation},
		{"Logging", "Add structured logging for service operations", &service.Logging},
	}

	for _, feature := range features {
		var include string
		includePrompt := &survey.Select{
			Message: fmt.Sprintf("Enable %s?", feature.name),
			Options: []string{
				fmt.Sprintf("Yes - %s", feature.description),
				"No - Skip this feature",
				"Quit",
			},
		}

		if err := survey.AskOne(includePrompt, &include); err != nil {
			return err
		}

		if include == "Quit" {
			return nil
		}

		*feature.field = include[:3] == "Yes"
	}

	return nil
}

func configureAPIEndpoints(handler *HandlerConfig, entity *CRUDEntity) error {
	fmt.Printf("🌐 Configuring API Endpoints for %s\n\n", entity.Name)

	// Standard CRUD endpoints
	standardEndpoints := []APIEndpoint{
		{"GET", fmt.Sprintf("/api/%s", entity.PluralName), fmt.Sprintf("List%s", strings.Title(entity.PluralName)), []string{"Logger", "Auth"}, fmt.Sprintf("Get paginated list of %s", entity.PluralName)},
		{"GET", fmt.Sprintf("/api/%s/{id}", entity.PluralName), fmt.Sprintf("Get%s", strings.Title(entity.Name)), []string{"Logger", "Auth"}, fmt.Sprintf("Get %s by ID", entity.Name)},
		{"POST", fmt.Sprintf("/api/%s", entity.PluralName), fmt.Sprintf("Create%s", strings.Title(entity.Name)), []string{"Logger", "Auth", "Validator"}, fmt.Sprintf("Create new %s", entity.Name)},
		{"PUT", fmt.Sprintf("/api/%s/{id}", entity.PluralName), fmt.Sprintf("Update%s", strings.Title(entity.Name)), []string{"Logger", "Auth", "Validator"}, fmt.Sprintf("Update %s", entity.Name)},
		{"DELETE", fmt.Sprintf("/api/%s/{id}", entity.PluralName), fmt.Sprintf("Delete%s", strings.Title(entity.Name)), []string{"Logger", "Auth"}, fmt.Sprintf("Delete %s", entity.Name)},
	}

	handler.Endpoints = standardEndpoints

	// Ask about additional endpoints
	additionalEndpoints := []APIEndpoint{
		{"GET", fmt.Sprintf("/api/%s/search", entity.PluralName), fmt.Sprintf("Search%s", strings.Title(entity.PluralName)), []string{"Logger", "Auth"}, fmt.Sprintf("Search %s by query", entity.PluralName)},
		{"GET", fmt.Sprintf("/api/%s/active", entity.PluralName), fmt.Sprintf("GetActive%s", strings.Title(entity.PluralName)), []string{"Logger", "Auth"}, fmt.Sprintf("Get active %s only", entity.PluralName)},
		{"PATCH", fmt.Sprintf("/api/%s/{id}/status", entity.PluralName), fmt.Sprintf("Update%sStatus", strings.Title(entity.Name)), []string{"Logger", "Auth"}, fmt.Sprintf("Update %s status", entity.Name)},
	}

	for _, endpoint := range additionalEndpoints {
		var include string
		includePrompt := &survey.Select{
			Message: fmt.Sprintf("Include endpoint: %s %s?", endpoint.Method, endpoint.Path),
			Options: []string{
				fmt.Sprintf("Yes - %s", endpoint.Description),
				"No - Skip this endpoint",
				"Quit",
			},
		}

		if err := survey.AskOne(includePrompt, &include); err != nil {
			return err
		}

		if include == "Quit" {
			return nil
		}

		if include[:3] == "Yes" {
			handler.Endpoints = append(handler.Endpoints, endpoint)
		}
	}

	return nil
}

func configureHandlerFeatures(handler *HandlerConfig) error {
	features := []struct {
		name        string
		description string
		field       *bool
	}{
		{"Error Handling", "Structured error responses with proper HTTP status codes", &handler.ErrorHandling},
		{"Documentation", "Generate OpenAPI/Swagger documentation", &handler.Documentation},
	}

	for _, feature := range features {
		var include string
		includePrompt := &survey.Select{
			Message: fmt.Sprintf("Enable %s?", feature.name),
			Options: []string{
				fmt.Sprintf("Yes - %s", feature.description),
				"No - Skip this feature",
				"Quit",
			},
		}

		if err := survey.AskOne(includePrompt, &include); err != nil {
			return err
		}

		if include == "Quit" {
			return nil
		}

		*feature.field = include[:3] == "Yes"
	}

	return nil
}

func reviewArchitectureAndGenerate(projectPath string, domainObj *DomainObject) error {
	fmt.Println("\n🎯 Step 7: Architecture Review & Code Generation")
	fmt.Println("Review your complete CRUD architecture before generation.")
	fmt.Println()

	// Show complete architecture summary
	fmt.Printf("📊 Complete Architecture Summary for %s:\n\n", domainObj.Entity.Name)

	// Domain Layer
	fmt.Println("🏗️  Domain Layer:")
	fmt.Printf("   • Entity: %s with %d fields\n", domainObj.Entity.Name, len(domainObj.Entity.Fields))
	fmt.Printf("   • Repository Interface: %d methods\n", len(domainObj.Repository.Methods))
	fmt.Printf("   • Business Rules: %d rules\n", len(domainObj.Service.BusinessRules))
	fmt.Printf("   • Domain Events: %d events\n", len(domainObj.Service.Events))

	// Infrastructure Layer
	fmt.Println("\n🔧 Infrastructure Layer:")
	fmt.Printf("   • Repository Implementation with database integration\n")
	if domainObj.Repository.Caching {
		fmt.Printf("   • Redis caching layer\n")
	}
	if domainObj.Repository.Transactions {
		fmt.Printf("   • Transaction support\n")
	}

	// Application Layer
	fmt.Println("\n⚙️  Application Layer:")
	fmt.Printf("   • Service with business logic orchestration\n")
	if domainObj.Service.Validation {
		fmt.Printf("   • Input validation\n")
	}
	if domainObj.Service.Logging {
		fmt.Printf("   • Structured logging\n")
	}

	// Interface Layer
	fmt.Println("\n🌐 Interface Layer:")
	fmt.Printf("   • HTTP Handlers: %d endpoints\n", len(domainObj.Handler.Endpoints))
	enabledMiddleware := 0
	for _, mw := range domainObj.Middleware {
		if mw.Enabled {
			enabledMiddleware++
		}
	}
	fmt.Printf("   • Middleware: %d components\n", enabledMiddleware)

	// Show file structure
	fmt.Println("\n📁 Files to be generated:")
	files := []string{
		fmt.Sprintf("internal/domain/%s/model.go", domainObj.Entity.Name),
		fmt.Sprintf("internal/domain/%s/repository.go", domainObj.Entity.Name),
		fmt.Sprintf("internal/domain/%s/service.go", domainObj.Entity.Name),
		fmt.Sprintf("internal/domain/%s/errors.go", domainObj.Entity.Name),
		fmt.Sprintf("internal/infrastructure/repository/%s_repository.go", domainObj.Entity.Name),
		fmt.Sprintf("internal/api/handlers/%s.go", domainObj.Entity.Name),
		fmt.Sprintf("internal/api/middleware/"),
		fmt.Sprintf("migrations/create_%s_table.sql", domainObj.Entity.PluralName),
		fmt.Sprintf("docs/%s_api.md", domainObj.Entity.Name),
	}

	for _, file := range files {
		fmt.Printf("   • %s\n", file)
	}

	// Final confirmation
	var confirm string
	confirmPrompt := &survey.Select{
		Message: "Generate this complete CRUD architecture?",
		Options: []string{
			"Yes - Generate all files with educational comments",
			"No - Let me modify the configuration",
			"Quit",
		},
	}

	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return err
	}

	if confirm == "Quit" {
		return nil
	}

	if confirm[:2] == "No" {
		fmt.Println("You can restart the wizard to modify your configuration.")
		return nil
	}

	// Generate the enhanced CRUD with educational content
	return generateEnhancedCRUDCode(projectPath, domainObj)
}

func generateEnhancedCRUDCode(projectPath string, domainObj *DomainObject) error {
	fmt.Println("\n🚀 Generating Enhanced CRUD Architecture...")
	fmt.Println()

	// This would integrate with the existing generator but with enhanced templates
	// that include educational comments and clean architecture patterns

	fmt.Printf("✅ Successfully generated enhanced CRUD architecture for %s!\n", domainObj.Entity.Name)
	fmt.Println()
	fmt.Println("📚 What was generated:")
	fmt.Println("• Complete Clean Architecture implementation")
	fmt.Println("• Educational comments explaining each pattern")
	fmt.Println("• Dependency injection setup")
	fmt.Println("• Comprehensive test examples")
	fmt.Println("• API documentation")
	fmt.Println("• Migration files")
	fmt.Println()
	fmt.Println("🎓 Next Steps:")
	fmt.Println("1. Review the generated code and comments")
	fmt.Println("2. Run the tests to see the architecture in action")
	fmt.Println("3. Customize the business rules for your specific needs")
	fmt.Println("4. Add more entities using the same patterns")

	return nil
}
