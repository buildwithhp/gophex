package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// CRUDField represents a field in the entity
type CRUDField struct {
	Name        string
	Type        string
	JSONTag     string
	DBTag       string
	Required    bool
	Unique      bool
	Description string
}

// CRUDEntity represents the entity to generate CRUD for
type CRUDEntity struct {
	Name         string
	PluralName   string
	Fields       []CRUDField
	UpdateMethod string // "put", "patch", or "both"
}

// UpdateMethodChoice represents the update method selection
type UpdateMethodChoice struct {
	Value       string
	Description string
	UseCase     string
	Example     string
}

// RunCRUDWizard runs the interactive CRUD generation wizard
func RunCRUDWizard(projectPath string) error {
	fmt.Println("🚀 Interactive CRUD Generator")
	fmt.Println("Let's create your first CRUD operations with proper Go idioms!")
	fmt.Println()

	entity := &CRUDEntity{}

	// Step 1: Entity Selection
	if err := selectEntity(entity); err != nil {
		return err
	}

	// Step 2: Field Definition
	if err := defineFields(entity); err != nil {
		return err
	}

	// Step 3: Update Method Selection
	if err := selectUpdateMethod(entity); err != nil {
		return err
	}

	// Step 4: Preview and Confirm
	if err := previewAndConfirm(entity); err != nil {
		return err
	}

	// Step 5: Generate Code
	return generateCRUDCode(projectPath, entity)
}

// selectEntity handles entity name selection
func selectEntity(entity *CRUDEntity) error {
	fmt.Println("📝 Step 1: Entity Selection")
	fmt.Println("What would you like to create CRUD operations for?")
	fmt.Println()

	// Provide common examples
	commonEntities := []string{
		"user - User management (name, email, password)",
		"post - Blog posts or articles (title, content, author)",
		"product - E-commerce products (name, price, description)",
		"task - Todo or task management (title, description, status)",
		"custom - I'll define my own entity",
	}

	var selected string
	entityPrompt := &survey.Select{
		Message: "Choose an entity type:",
		Options: commonEntities,
		Help:    "Select a common entity or choose 'custom' to define your own",
	}

	if err := survey.AskOne(entityPrompt, &selected); err != nil {
		if isUserInterrupt(err) {
			return nil
		}
		return fmt.Errorf("entity selection failed: %w", err)
	}

	if strings.HasPrefix(selected, "custom") {
		// Custom entity name
		var customName string
		namePrompt := &survey.Input{
			Message: "Enter your entity name (singular, e.g., 'book', 'order'):",
			Help:    "Use lowercase, singular form. We'll generate the plural automatically.",
		}

		if err := survey.AskOne(namePrompt, &customName); err != nil {
			if isUserInterrupt(err) {
				return nil
			}
			return fmt.Errorf("custom entity name input failed: %w", err)
		}

		if !isValidEntityName(customName) {
			return fmt.Errorf("invalid entity name: must be lowercase letters only")
		}

		entity.Name = strings.TrimSpace(customName)
	} else {
		// Extract entity name from selection
		parts := strings.Split(selected, " - ")
		entity.Name = parts[0]
	}

	entity.PluralName = pluralize(entity.Name)

	fmt.Printf("✅ Selected entity: %s (plural: %s)\n\n", entity.Name, entity.PluralName)
	return nil
}

// defineFields handles field definition
func defineFields(entity *CRUDEntity) error {
	fmt.Println("🏗️  Step 2: Field Definition")
	fmt.Printf("Let's define the fields for your %s entity.\n", entity.Name)
	fmt.Println()

	// Add common fields based on entity type
	entity.Fields = getCommonFields(entity.Name)

	if len(entity.Fields) > 0 {
		fmt.Printf("I've added some common fields for %s:\n", entity.Name)
		for _, field := range entity.Fields {
			fmt.Printf("  - %s (%s) %s\n", field.Name, field.Type, field.Description)
		}
		fmt.Println()

		var useCommon bool
		commonPrompt := &survey.Confirm{
			Message: "Would you like to use these common fields?",
			Default: true,
			Help:    "You can modify or add more fields in the next step",
		}

		if err := survey.AskOne(commonPrompt, &useCommon); err != nil {
			if isUserInterrupt(err) {
				return nil
			}
			return fmt.Errorf("common fields prompt failed: %w", err)
		}

		if !useCommon {
			entity.Fields = []CRUDField{}
		}
	}

	// Interactive field addition
	for {
		var addMore bool
		if len(entity.Fields) == 0 {
			addMore = true
		} else {
			addPrompt := &survey.Confirm{
				Message: "Would you like to add more fields?",
				Default: false,
			}

			if err := survey.AskOne(addPrompt, &addMore); err != nil {
				if isUserInterrupt(err) {
					return nil
				}
				return fmt.Errorf("add field prompt failed: %w", err)
			}
		}

		if !addMore {
			break
		}

		field, err := defineField()
		if err != nil {
			return err
		}

		entity.Fields = append(entity.Fields, field)
		fmt.Printf("✅ Added field: %s (%s)\n", field.Name, field.Type)
	}

	if len(entity.Fields) == 0 {
		return fmt.Errorf("at least one field is required")
	}

	fmt.Printf("\n📋 Final fields for %s:\n", entity.Name)
	for _, field := range entity.Fields {
		required := ""
		if field.Required {
			required = " (required)"
		}
		unique := ""
		if field.Unique {
			unique = " (unique)"
		}
		fmt.Printf("  - %s: %s%s%s\n", field.Name, field.Type, required, unique)
	}
	fmt.Println()

	return nil
}

// selectUpdateMethod handles update method selection with education
func selectUpdateMethod(entity *CRUDEntity) error {
	fmt.Println("🔄 Step 3: Update Method Selection")
	fmt.Println("How would you like to handle updates? Let me explain the differences:")
	fmt.Println()

	choices := []UpdateMethodChoice{
		{
			Value:       "put",
			Description: "PUT - Complete Replacement",
			UseCase:     "Replaces the entire resource. All fields must be provided.",
			Example:     "Updating a user profile where all fields are required for consistency",
		},
		{
			Value:       "patch",
			Description: "PATCH - Partial Update",
			UseCase:     "Updates only the fields provided. Missing fields remain unchanged.",
			Example:     "Updating just a user's email or status without affecting other fields",
		},
		{
			Value:       "both",
			Description: "Both PUT and PATCH",
			UseCase:     "Provides both options for maximum flexibility in your API.",
			Example:     "Different clients can choose the most appropriate method for their use case",
		},
	}

	// Display detailed explanations
	for i, choice := range choices {
		fmt.Printf("%d. %s\n", i+1, choice.Description)
		fmt.Printf("   📖 %s\n", choice.UseCase)
		fmt.Printf("   💡 Example: %s\n\n", choice.Example)
	}

	options := make([]string, len(choices))
	for i, choice := range choices {
		options[i] = choice.Description
	}

	var selected string
	methodPrompt := &survey.Select{
		Message: "Choose your update method:",
		Options: options,
		Help:    "This affects how your API will handle resource updates",
	}

	if err := survey.AskOne(methodPrompt, &selected); err != nil {
		if isUserInterrupt(err) {
			return nil
		}
		return fmt.Errorf("update method selection failed: %w", err)
	}

	// Find the selected choice
	for _, choice := range choices {
		if choice.Description == selected {
			entity.UpdateMethod = choice.Value
			break
		}
	}

	fmt.Printf("✅ Selected: %s\n\n", selected)
	return nil
}

// previewAndConfirm shows what will be generated
func previewAndConfirm(entity *CRUDEntity) error {
	fmt.Println("👀 Step 4: Preview")
	fmt.Printf("Here's what will be generated for your %s entity:\n\n", entity.Name)

	// Show endpoints
	fmt.Println("📡 API Endpoints:")
	fmt.Printf("  GET    /api/%s     - List %s with pagination\n", entity.PluralName, entity.PluralName)
	fmt.Printf("  GET    /api/%s/{id} - Get %s by ID\n", entity.PluralName, entity.Name)
	fmt.Printf("  POST   /api/%s     - Create new %s\n", entity.PluralName, entity.Name)

	switch entity.UpdateMethod {
	case "put":
		fmt.Printf("  PUT    /api/%s/{id} - Complete %s replacement (all fields required)\n", entity.PluralName, entity.Name)
	case "patch":
		fmt.Printf("  PATCH  /api/%s/{id} - Partial %s update (only provided fields)\n", entity.PluralName, entity.Name)
	case "both":
		fmt.Printf("  PUT    /api/%s/{id} - Complete %s replacement (all fields required)\n", entity.PluralName, entity.Name)
		fmt.Printf("  PATCH  /api/%s/{id} - Partial %s update (only provided fields)\n", entity.PluralName, entity.Name)
	}

	fmt.Printf("  DELETE /api/%s/{id} - Delete %s\n\n", entity.PluralName, entity.Name)

	// Show files that will be created
	fmt.Println("📁 Files to be created/updated:")
	fmt.Printf("  internal/domain/%s/model.go       - Data model\n", entity.Name)
	fmt.Printf("  internal/domain/%s/repository.go  - Database operations\n", entity.Name)
	fmt.Printf("  internal/domain/%s/service.go     - Business logic\n", entity.Name)
	fmt.Printf("  internal/api/handlers/%s.go       - HTTP handlers\n", entity.Name)
	fmt.Printf("  internal/api/routes/routes.go     - Route registration (updated)\n")
	fmt.Printf("  migrations/                       - Database migration files\n")
	fmt.Printf("  README_%s.md                      - Documentation and examples\n\n", entity.Name)

	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: "Generate these CRUD operations?",
		Default: true,
	}

	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		if isUserInterrupt(err) {
			return nil
		}
		return fmt.Errorf("confirmation prompt failed: %w", err)
	}

	if !confirm {
		fmt.Println("❌ CRUD generation cancelled")
		return ErrReturnToMenu
	}

	return nil
}

// defineField handles individual field definition
func defineField() (CRUDField, error) {
	field := CRUDField{}

	// Field name
	namePrompt := &survey.Input{
		Message: "Field name (e.g., 'email', 'title', 'price'):",
		Help:    "Use camelCase for Go conventions",
	}

	if err := survey.AskOne(namePrompt, &field.Name); err != nil {
		return field, fmt.Errorf("field name input failed: %w", err)
	}

	if !isValidFieldName(field.Name) {
		return field, fmt.Errorf("invalid field name: must start with letter and contain only letters/numbers")
	}

	// Field type
	fieldTypes := []string{
		"string - Text data",
		"int - Integer numbers",
		"int64 - Large integer numbers",
		"float64 - Decimal numbers",
		"bool - True/false values",
		"time.Time - Date and time",
		"[]string - Array of strings",
	}

	var selectedType string
	typePrompt := &survey.Select{
		Message: "Field type:",
		Options: fieldTypes,
	}

	if err := survey.AskOne(typePrompt, &selectedType); err != nil {
		return field, fmt.Errorf("field type selection failed: %w", err)
	}

	field.Type = strings.Split(selectedType, " - ")[0]

	// Field properties
	requiredPrompt := &survey.Confirm{
		Message: "Is this field required?",
		Default: false,
	}

	if err := survey.AskOne(requiredPrompt, &field.Required); err != nil {
		return field, fmt.Errorf("required prompt failed: %w", err)
	}

	uniquePrompt := &survey.Confirm{
		Message: "Should this field be unique?",
		Default: false,
	}

	if err := survey.AskOne(uniquePrompt, &field.Unique); err != nil {
		return field, fmt.Errorf("unique prompt failed: %w", err)
	}

	// Generate tags
	field.JSONTag = strings.ToLower(field.Name)
	field.DBTag = strings.ToLower(field.Name)

	return field, nil
}

// Helper functions

func isValidEntityName(name string) bool {
	matched, _ := regexp.MatchString("^[a-z][a-z0-9]*$", name)
	return matched
}

func isValidFieldName(name string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9]*$", name)
	return matched
}

func pluralize(singular string) string {
	// Simple pluralization rules
	if strings.HasSuffix(singular, "y") && len(singular) > 1 {
		// Check if the letter before 'y' is a consonant
		beforeY := singular[len(singular)-2]
		if beforeY != 'a' && beforeY != 'e' && beforeY != 'i' && beforeY != 'o' && beforeY != 'u' {
			return strings.TrimSuffix(singular, "y") + "ies"
		}
	}
	if strings.HasSuffix(singular, "s") || strings.HasSuffix(singular, "x") ||
		strings.HasSuffix(singular, "z") || strings.HasSuffix(singular, "sh") ||
		strings.HasSuffix(singular, "ch") {
		return singular + "es"
	}
	if strings.HasSuffix(singular, "f") {
		return strings.TrimSuffix(singular, "f") + "ves"
	}
	if strings.HasSuffix(singular, "fe") {
		return strings.TrimSuffix(singular, "fe") + "ves"
	}
	return singular + "s"
}

func getCommonFields(entityName string) []CRUDField {
	switch entityName {
	case "user":
		return []CRUDField{
			{Name: "Name", Type: "string", JSONTag: "name", DBTag: "name", Required: true, Description: "- User's full name"},
			{Name: "Email", Type: "string", JSONTag: "email", DBTag: "email", Required: true, Unique: true, Description: "- User's email address"},
			{Name: "Password", Type: "string", JSONTag: "password", DBTag: "password", Required: true, Description: "- User's password (will be hashed)"},
			{Name: "CreatedAt", Type: "time.Time", JSONTag: "created_at", DBTag: "created_at", Description: "- Account creation timestamp"},
			{Name: "UpdatedAt", Type: "time.Time", JSONTag: "updated_at", DBTag: "updated_at", Description: "- Last update timestamp"},
		}
	case "post":
		return []CRUDField{
			{Name: "Title", Type: "string", JSONTag: "title", DBTag: "title", Required: true, Description: "- Post title"},
			{Name: "Content", Type: "string", JSONTag: "content", DBTag: "content", Required: true, Description: "- Post content"},
			{Name: "AuthorID", Type: "int64", JSONTag: "author_id", DBTag: "author_id", Required: true, Description: "- Author's user ID"},
			{Name: "Published", Type: "bool", JSONTag: "published", DBTag: "published", Description: "- Publication status"},
			{Name: "CreatedAt", Type: "time.Time", JSONTag: "created_at", DBTag: "created_at", Description: "- Creation timestamp"},
			{Name: "UpdatedAt", Type: "time.Time", JSONTag: "updated_at", DBTag: "updated_at", Description: "- Last update timestamp"},
		}
	case "product":
		return []CRUDField{
			{Name: "Name", Type: "string", JSONTag: "name", DBTag: "name", Required: true, Description: "- Product name"},
			{Name: "Description", Type: "string", JSONTag: "description", DBTag: "description", Description: "- Product description"},
			{Name: "Price", Type: "float64", JSONTag: "price", DBTag: "price", Required: true, Description: "- Product price"},
			{Name: "SKU", Type: "string", JSONTag: "sku", DBTag: "sku", Unique: true, Description: "- Stock keeping unit"},
			{Name: "InStock", Type: "bool", JSONTag: "in_stock", DBTag: "in_stock", Description: "- Availability status"},
			{Name: "CreatedAt", Type: "time.Time", JSONTag: "created_at", DBTag: "created_at", Description: "- Creation timestamp"},
			{Name: "UpdatedAt", Type: "time.Time", JSONTag: "updated_at", DBTag: "updated_at", Description: "- Last update timestamp"},
		}
	case "task":
		return []CRUDField{
			{Name: "Title", Type: "string", JSONTag: "title", DBTag: "title", Required: true, Description: "- Task title"},
			{Name: "Description", Type: "string", JSONTag: "description", DBTag: "description", Description: "- Task description"},
			{Name: "Status", Type: "string", JSONTag: "status", DBTag: "status", Required: true, Description: "- Task status (pending, in_progress, completed)"},
			{Name: "Priority", Type: "string", JSONTag: "priority", DBTag: "priority", Description: "- Task priority (low, medium, high)"},
			{Name: "DueDate", Type: "time.Time", JSONTag: "due_date", DBTag: "due_date", Description: "- Task due date"},
			{Name: "CreatedAt", Type: "time.Time", JSONTag: "created_at", DBTag: "created_at", Description: "- Creation timestamp"},
			{Name: "UpdatedAt", Type: "time.Time", JSONTag: "updated_at", DBTag: "updated_at", Description: "- Last update timestamp"},
		}
	default:
		return []CRUDField{
			{Name: "CreatedAt", Type: "time.Time", JSONTag: "created_at", DBTag: "created_at", Description: "- Creation timestamp"},
			{Name: "UpdatedAt", Type: "time.Time", JSONTag: "updated_at", DBTag: "updated_at", Description: "- Last update timestamp"},
		}
	}
}
