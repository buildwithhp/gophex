package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/buildwithhp/gophex/internal/utils"
)

// CRUDTemplateData contains all data needed for CRUD template generation
type CRUDTemplateData struct {
	Entity       *CRUDEntity
	ModuleName   string
	ProjectName  string
	DatabaseType string
	Timestamp    string
}

// generateCRUDCode generates all CRUD-related files
func generateCRUDCode(projectPath string, entity *CRUDEntity) error {
	fmt.Printf("🔨 Generating CRUD operations for %s...\n", entity.Name)

	// Load project metadata to get module name and database type
	metadata, err := utils.LoadMetadata(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load project metadata: %w", err)
	}

	// Determine module name from go.mod
	moduleName, err := getModuleName(projectPath)
	if err != nil {
		return fmt.Errorf("failed to get module name: %w", err)
	}

	// Determine database type from existing config
	databaseType, err := getDatabaseType(projectPath)
	if err != nil {
		return fmt.Errorf("failed to determine database type: %w", err)
	}

	templateData := &CRUDTemplateData{
		Entity:       entity,
		ModuleName:   moduleName,
		ProjectName:  metadata.Project.Name,
		DatabaseType: databaseType,
		Timestamp:    time.Now().Format(time.RFC3339),
	}

	// Create directory structure
	if err := createCRUDDirectories(projectPath, entity); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Generate files
	if err := generateModelFile(projectPath, templateData); err != nil {
		return fmt.Errorf("failed to generate model: %w", err)
	}

	if err := generateRepositoryFile(projectPath, templateData); err != nil {
		return fmt.Errorf("failed to generate repository: %w", err)
	}

	if err := generateServiceFile(projectPath, templateData); err != nil {
		return fmt.Errorf("failed to generate service: %w", err)
	}

	if err := generateHandlerFile(projectPath, templateData); err != nil {
		return fmt.Errorf("failed to generate handler: %w", err)
	}

	if err := updateRoutesFile(projectPath, templateData); err != nil {
		return fmt.Errorf("failed to update routes: %w", err)
	}

	if err := generateMigrationFiles(projectPath, templateData); err != nil {
		return fmt.Errorf("failed to generate migrations: %w", err)
	}

	if err := generateDocumentation(projectPath, templateData); err != nil {
		return fmt.Errorf("failed to generate documentation: %w", err)
	}

	fmt.Printf("✅ Successfully generated CRUD operations for %s!\n\n", entity.Name)

	// Show next steps
	showNextSteps(entity)

	return nil
}

// createCRUDDirectories creates necessary directory structure
func createCRUDDirectories(projectPath string, entity *CRUDEntity) error {
	dirs := []string{
		filepath.Join(projectPath, "internal", "domain", entity.Name),
		filepath.Join(projectPath, "internal", "api", "handlers"),
		filepath.Join(projectPath, "migrations"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateModelFile generates the model file
func generateModelFile(projectPath string, data *CRUDTemplateData) error {
	tmpl := `package {{.Entity.Name}}

import (
	"time"
{{if hasTimeFields .Entity.Fields}}	"database/sql/driver"
	"fmt"{{end}}
)

// {{title .Entity.Name}} represents a {{.Entity.Name}} entity
type {{title .Entity.Name}} struct {
	ID {{if eq .DatabaseType "mongodb"}}primitive.ObjectID ` + "`json:\"id\" bson:\"_id,omitempty\"`" + `{{else}}int64 ` + "`json:\"id\" db:\"id\"`" + `{{end}}
{{range .Entity.Fields}}	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONTag}}\" db:\"{{.DBTag}}\"`" + `
{{end}}
}

// Create{{title .Entity.Name}}Request represents the request payload for creating a {{.Entity.Name}}
type Create{{title .Entity.Name}}Request struct {
{{range .Entity.Fields}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONTag}}\"{{if .Required}} validate:\"required\"{{end}}`" + `
{{end}}{{end}}{{end}}
}

{{if or (eq .Entity.UpdateMethod "put") (eq .Entity.UpdateMethod "both")}}
// Update{{title .Entity.Name}}Request represents the request payload for updating a {{.Entity.Name}} (PUT - complete replacement)
// All fields are required as this replaces the entire resource
type Update{{title .Entity.Name}}Request struct {
{{range .Entity.Fields}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONTag}}\"{{if .Required}} validate:\"required\"{{end}}`" + `
{{end}}{{end}}{{end}}
}
{{end}}

{{if or (eq .Entity.UpdateMethod "patch") (eq .Entity.UpdateMethod "both")}}
// Patch{{title .Entity.Name}}Request represents the request payload for patching a {{.Entity.Name}} (PATCH - partial update)
// All fields are optional as this only updates provided fields
type Patch{{title .Entity.Name}}Request struct {
{{range .Entity.Fields}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}	{{.Name}} *{{.Type}} ` + "`json:\"{{.JSONTag}},omitempty\"`" + `
{{end}}{{end}}{{end}}
}
{{end}}

// {{title .Entity.Name}}Response represents the response payload for a {{.Entity.Name}}
type {{title .Entity.Name}}Response struct {
	ID {{if eq .DatabaseType "mongodb"}}string{{else}}int64{{end}} ` + "`json:\"id\"`" + `
{{range .Entity.Fields}}	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONTag}}\"`" + `
{{end}}
}

// List{{title .Entity.PluralName}}Response represents the response payload for listing {{.Entity.PluralName}}
type List{{title .Entity.PluralName}}Response struct {
	{{title .Entity.PluralName}} []{{title .Entity.Name}}Response ` + "`json:\"{{.Entity.PluralName}}\"`" + `
	Total    int64                     ` + "`json:\"total\"`" + `
	Page     int                       ` + "`json:\"page\"`" + `
	PageSize int                       ` + "`json:\"page_size\"`" + `
}

// ToResponse converts a {{title .Entity.Name}} to {{title .Entity.Name}}Response
func ({{lower .Entity.Name}} *{{title .Entity.Name}}) ToResponse() {{title .Entity.Name}}Response {
	return {{title .Entity.Name}}Response{
		ID: {{if eq .DatabaseType "mongodb"}}{{lower .Entity.Name}}.ID.Hex(){{else}}{{lower .Entity.Name}}.ID{{end}},
{{range .Entity.Fields}}		{{.Name}}: {{lower $.Entity.Name}}.{{.Name}},
{{end}}
	}
}

// Validate validates the {{title .Entity.Name}} fields
func ({{lower .Entity.Name}} *{{title .Entity.Name}}) Validate() error {
{{range .Entity.Fields}}{{if .Required}}	if {{if eq .Type "string"}}{{lower $.Entity.Name}}.{{.Name}} == ""{{else if eq .Type "int"}}{{lower $.Entity.Name}}.{{.Name}} == 0{{else if eq .Type "int64"}}{{lower $.Entity.Name}}.{{.Name}} == 0{{else}}{{lower $.Entity.Name}}.{{.Name}} == nil{{end}} {
		return fmt.Errorf("{{.Name}} is required")
	}
{{end}}{{end}}
	return nil
}
`

	filePath := filepath.Join(projectPath, "internal", "domain", data.Entity.Name, "model.go")
	return executeTemplate(tmpl, filePath, data)
}

// generateRepositoryFile generates the repository file
func generateRepositoryFile(projectPath string, data *CRUDTemplateData) error {
	tmpl := `package {{.Entity.Name}}

import (
	"context"
	"fmt"
{{if eq .DatabaseType "mongodb"}}	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"{{else}}	"database/sql"
	"strings"{{end}}
)

// Repository defines the interface for {{.Entity.Name}} data operations
type Repository interface {
	Create(ctx context.Context, {{.Entity.Name}} *{{title .Entity.Name}}) error
	GetByID(ctx context.Context, id {{if eq .DatabaseType "mongodb"}}string{{else}}int64{{end}}) (*{{title .Entity.Name}}, error)
	List(ctx context.Context, page, pageSize int) ([]{{title .Entity.Name}}, int64, error)
{{if or (eq .Entity.UpdateMethod "put") (eq .Entity.UpdateMethod "both")}}	Update(ctx context.Context, {{.Entity.Name}} *{{title .Entity.Name}}) error{{end}}
{{if or (eq .Entity.UpdateMethod "patch") (eq .Entity.UpdateMethod "both")}}	Patch(ctx context.Context, id {{if eq .DatabaseType "mongodb"}}string{{else}}int64{{end}}, updates map[string]interface{}) error{{end}}
	Delete(ctx context.Context, id {{if eq .DatabaseType "mongodb"}}string{{else}}int64{{end}}) error
}

{{if eq .DatabaseType "mongodb"}}
// mongoRepository implements Repository for MongoDB
type mongoRepository struct {
	collection *mongo.Collection
}

// NewRepository creates a new MongoDB repository
func NewRepository(db *mongo.Database) Repository {
	return &mongoRepository{
		collection: db.Collection("{{.Entity.PluralName}}"),
	}
}

func (r *mongoRepository) Create(ctx context.Context, {{.Entity.Name}} *{{title .Entity.Name}}) error {
	{{.Entity.Name}}.ID = primitive.NewObjectID()
	result, err := r.collection.InsertOne(ctx, {{.Entity.Name}})
	if err != nil {
		return fmt.Errorf("failed to create {{.Entity.Name}}: %w", err)
	}
	{{.Entity.Name}}.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *mongoRepository) GetByID(ctx context.Context, id string) (*{{title .Entity.Name}}, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	var {{.Entity.Name}} {{title .Entity.Name}}
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&{{.Entity.Name}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("{{.Entity.Name}} not found")
		}
		return nil, fmt.Errorf("failed to get {{.Entity.Name}}: %w", err)
	}
	return &{{.Entity.Name}}, nil
}

func (r *mongoRepository) List(ctx context.Context, page, pageSize int) ([]{{title .Entity.Name}}, int64, error) {
	skip := (page - 1) * pageSize
	
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list {{.Entity.PluralName}}: %w", err)
	}
	defer cursor.Close(ctx)

	var {{.Entity.PluralName}} []{{title .Entity.Name}}
	if err = cursor.All(ctx, &{{.Entity.PluralName}}); err != nil {
		return nil, 0, fmt.Errorf("failed to decode {{.Entity.PluralName}}: %w", err)
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count {{.Entity.PluralName}}: %w", err)
	}

	return {{.Entity.PluralName}}, total, nil
}

{{if or (eq .Entity.UpdateMethod "put") (eq .Entity.UpdateMethod "both")}}
func (r *mongoRepository) Update(ctx context.Context, {{.Entity.Name}} *{{title .Entity.Name}}) error {
	filter := bson.M{"_id": {{.Entity.Name}}.ID}
	update := bson.M{"$set": {{.Entity.Name}}}
	
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update {{.Entity.Name}}: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("{{.Entity.Name}} not found")
	}
	return nil
}
{{end}}

{{if or (eq .Entity.UpdateMethod "patch") (eq .Entity.UpdateMethod "both")}}
func (r *mongoRepository) Patch(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updates}
	
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to patch {{.Entity.Name}}: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("{{.Entity.Name}} not found")
	}
	return nil
}
{{end}}

func (r *mongoRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete {{.Entity.Name}}: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("{{.Entity.Name}} not found")
	}
	return nil
}

{{else}}
// sqlRepository implements Repository for SQL databases
type sqlRepository struct {
	db *sql.DB
}

// NewRepository creates a new SQL repository
func NewRepository(db *sql.DB) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) Create(ctx context.Context, {{.Entity.Name}} *{{title .Entity.Name}}) error {
	query := ` + "`INSERT INTO {{.Entity.PluralName}} ({{range $i, $field := .Entity.Fields}}{{if $i}}, {{end}}{{.DBTag}}{{end}}) VALUES ({{range $i, $field := .Entity.Fields}}{{if $i}}, {{end}}${{add $i 1}}{{end}}) RETURNING id`" + `
	
	err := r.db.QueryRowContext(ctx, query{{range .Entity.Fields}}, {{$.Entity.Name}}.{{.Name}}{{end}}).Scan(&{{.Entity.Name}}.ID)
	if err != nil {
		return fmt.Errorf("failed to create {{.Entity.Name}}: %w", err)
	}
	return nil
}

func (r *sqlRepository) GetByID(ctx context.Context, id int64) (*{{title .Entity.Name}}, error) {
	query := ` + "`SELECT id{{range .Entity.Fields}}, {{.DBTag}}{{end}} FROM {{.Entity.PluralName}} WHERE id = $1`" + `
	
	var {{.Entity.Name}} {{title .Entity.Name}}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&{{.Entity.Name}}.ID{{range .Entity.Fields}}, &{{$.Entity.Name}}.{{.Name}}{{end}})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("{{.Entity.Name}} not found")
		}
		return nil, fmt.Errorf("failed to get {{.Entity.Name}}: %w", err)
	}
	return &{{.Entity.Name}}, nil
}

func (r *sqlRepository) List(ctx context.Context, page, pageSize int) ([]{{title .Entity.Name}}, int64, error) {
	offset := (page - 1) * pageSize
	
	query := ` + "`SELECT id{{range .Entity.Fields}}, {{.DBTag}}{{end}} FROM {{.Entity.PluralName}} ORDER BY id LIMIT $1 OFFSET $2`" + `
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list {{.Entity.PluralName}}: %w", err)
	}
	defer rows.Close()

	var {{.Entity.PluralName}} []{{title .Entity.Name}}
	for rows.Next() {
		var {{.Entity.Name}} {{title .Entity.Name}}
		err := rows.Scan(&{{.Entity.Name}}.ID{{range .Entity.Fields}}, &{{$.Entity.Name}}.{{.Name}}{{end}})
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan {{.Entity.Name}}: %w", err)
		}
		{{.Entity.PluralName}} = append({{.Entity.PluralName}}, {{.Entity.Name}})
	}

	// Get total count
	var total int64
	countQuery := ` + "`SELECT COUNT(*) FROM {{.Entity.PluralName}}`" + `
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count {{.Entity.PluralName}}: %w", err)
	}

	return {{.Entity.PluralName}}, total, nil
}

{{if or (eq .Entity.UpdateMethod "put") (eq .Entity.UpdateMethod "both")}}
func (r *sqlRepository) Update(ctx context.Context, {{.Entity.Name}} *{{title .Entity.Name}}) error {
	query := ` + "`UPDATE {{.Entity.PluralName}} SET {{range $i, $field := .Entity.Fields}}{{if $i}}, {{end}}{{.DBTag}} = ${{add $i 2}}{{end}} WHERE id = $1`" + `
	
	result, err := r.db.ExecContext(ctx, query, {{.Entity.Name}}.ID{{range .Entity.Fields}}, {{$.Entity.Name}}.{{.Name}}{{end}})
	if err != nil {
		return fmt.Errorf("failed to update {{.Entity.Name}}: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("{{.Entity.Name}} not found")
	}
	return nil
}
{{end}}

{{if or (eq .Entity.UpdateMethod "patch") (eq .Entity.UpdateMethod "both")}}
func (r *sqlRepository) Patch(ctx context.Context, id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setParts := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+1)
	args = append(args, id) // First argument is always the ID
	
	argIndex := 2 // Start from $2 since $1 is the ID
	for field, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf("UPDATE {{.Entity.PluralName}} SET %s WHERE id = $1", strings.Join(setParts, ", "))
	
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to patch {{.Entity.Name}}: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("{{.Entity.Name}} not found")
	}
	return nil
}
{{end}}

func (r *sqlRepository) Delete(ctx context.Context, id int64) error {
	query := ` + "`DELETE FROM {{.Entity.PluralName}} WHERE id = $1`" + `
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete {{.Entity.Name}}: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("{{.Entity.Name}} not found")
	}
	return nil
}
{{end}}
`

	filePath := filepath.Join(projectPath, "internal", "domain", data.Entity.Name, "repository.go")
	return executeTemplate(tmpl, filePath, data)
}

// Helper functions for template execution
func executeTemplate(tmplStr, filePath string, data interface{}) error {
	funcMap := template.FuncMap{
		"title": strings.Title,
		"lower": strings.ToLower,
		"add": func(a, b int) int {
			return a + b
		},
		"hasTimeFields": func(fields []CRUDField) bool {
			for _, field := range fields {
				if field.Type == "time.Time" {
					return true
				}
			}
			return false
		},
		"hasField": func(fields []CRUDField, fieldName string) bool {
			for _, field := range fields {
				if field.Name == fieldName {
					return true
				}
			}
			return false
		},
	}

	tmpl, err := template.New("crud").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// Helper functions to get project information
func getModuleName(projectPath string) (string, error) {
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

func getDatabaseType(projectPath string) (string, error) {
	// Check if MongoDB files exist
	mongoPath := filepath.Join(projectPath, "internal", "infrastructure", "database", "mongodb")
	if _, err := os.Stat(mongoPath); err == nil {
		return "mongodb", nil
	}

	// Default to PostgreSQL for SQL databases
	return "postgresql", nil
}

// generateServiceFile generates the service file
func generateServiceFile(projectPath string, data *CRUDTemplateData) error {
	tmpl := `package {{.Entity.Name}}

import (
	"context"
	"fmt"
	"time"
)

// Service defines the business logic interface for {{.Entity.Name}}
type Service interface {
	Create(ctx context.Context, req Create{{title .Entity.Name}}Request) (*{{title .Entity.Name}}Response, error)
	GetByID(ctx context.Context, id {{if eq .DatabaseType "mongodb"}}string{{else}}int64{{end}}) (*{{title .Entity.Name}}Response, error)
	List(ctx context.Context, page, pageSize int) (*List{{title .Entity.PluralName}}Response, error)
{{if or (eq .Entity.UpdateMethod "put") (eq .Entity.UpdateMethod "both")}}	Update(ctx context.Context, id {{if eq .DatabaseType "mongodb"}}string{{else}}int64{{end}}, req Update{{title .Entity.Name}}Request) (*{{title .Entity.Name}}Response, error){{end}}
{{if or (eq .Entity.UpdateMethod "patch") (eq .Entity.UpdateMethod "both")}}	Patch(ctx context.Context, id {{if eq .DatabaseType "mongodb"}}string{{else}}int64{{end}}, req Patch{{title .Entity.Name}}Request) (*{{title .Entity.Name}}Response, error){{end}}
	Delete(ctx context.Context, id {{if eq .DatabaseType "mongodb"}}string{{else}}int64{{end}}) error
}

// service implements Service interface
type service struct {
	repo Repository
}

// NewService creates a new {{.Entity.Name}} service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Create creates a new {{.Entity.Name}}
func (s *service) Create(ctx context.Context, req Create{{title .Entity.Name}}Request) (*{{title .Entity.Name}}Response, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create entity
	{{.Entity.Name}} := &{{title .Entity.Name}}{
{{range .Entity.Fields}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}		{{.Name}}: req.{{.Name}},
{{end}}{{end}}{{end}}{{if hasField .Entity.Fields "CreatedAt"}}		CreatedAt: time.Now(),{{end}}
{{if hasField .Entity.Fields "UpdatedAt"}}		UpdatedAt: time.Now(),{{end}}
	}

	if err := s.repo.Create(ctx, {{.Entity.Name}}); err != nil {
		return nil, fmt.Errorf("failed to create {{.Entity.Name}}: %w", err)
	}

	response := {{.Entity.Name}}.ToResponse()
	return &response, nil
}

// GetByID retrieves a {{.Entity.Name}} by ID
func (s *service) GetByID(ctx context.Context, id {{if eq .DatabaseType "mongodb"}}string{{else}}int64{{end}}) (*{{title .Entity.Name}}Response, error) {
	{{.Entity.Name}}, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get {{.Entity.Name}}: %w", err)
	}

	response := {{.Entity.Name}}.ToResponse()
	return &response, nil
}

// List retrieves a paginated list of {{.Entity.PluralName}}
func (s *service) List(ctx context.Context, page, pageSize int) (*List{{title .Entity.PluralName}}Response, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	{{.Entity.PluralName}}, total, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list {{.Entity.PluralName}}: %w", err)
	}

	// Convert to response format
	responses := make([]{{title .Entity.Name}}Response, len({{.Entity.PluralName}}))
	for i, {{.Entity.Name}} := range {{.Entity.PluralName}} {
		responses[i] = {{.Entity.Name}}.ToResponse()
	}

	return &List{{title .Entity.PluralName}}Response{
		{{title .Entity.PluralName}}: responses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

{{if or (eq .Entity.UpdateMethod "put") (eq .Entity.UpdateMethod "both")}}
// Update performs a complete update of a {{.Entity.Name}} (PUT - replaces entire resource)
// All required fields must be provided as this replaces the entire resource
func (s *service) Update(ctx context.Context, id {{if eq .DatabaseType "mongodb"}}string{{else}}int64{{end}}, req Update{{title .Entity.Name}}Request) (*{{title .Entity.Name}}Response, error) {
	// Validate request - all required fields must be present for PUT
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if {{.Entity.Name}} exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get {{.Entity.Name}}: %w", err)
	}

	// Update all fields (complete replacement)
	updated := &{{title .Entity.Name}}{
		ID: existing.ID,
{{range .Entity.Fields}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}		{{.Name}}: req.{{.Name}},
{{else}}		{{.Name}}: time.Now(),{{end}}{{else}}		{{.Name}}: existing.{{.Name}},{{end}}{{end}}
	}

	if err := s.repo.Update(ctx, updated); err != nil {
		return nil, fmt.Errorf("failed to update {{.Entity.Name}}: %w", err)
	}

	response := updated.ToResponse()
	return &response, nil
}
{{end}}

{{if or (eq .Entity.UpdateMethod "patch") (eq .Entity.UpdateMethod "both")}}
// Patch performs a partial update of a {{.Entity.Name}} (PATCH - updates only provided fields)
// Only the fields provided in the request will be updated
func (s *service) Patch(ctx context.Context, id {{if eq .DatabaseType "mongodb"}}string{{else}}int64{{end}}, req Patch{{title .Entity.Name}}Request) (*{{title .Entity.Name}}Response, error) {
	// Check if {{.Entity.Name}} exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get {{.Entity.Name}}: %w", err)
	}

	// Build updates map with only provided fields
	updates := make(map[string]interface{})
{{range .Entity.Fields}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}	if req.{{.Name}} != nil {
		updates["{{.DBTag}}"] = *req.{{.Name}}
	}
{{end}}{{end}}{{end}}
{{if hasField .Entity.Fields "UpdatedAt"}}	// Always update the UpdatedAt timestamp for PATCH operations
	updates["updated_at"] = time.Now()
{{end}}

	if len(updates) == 0 {
		// No fields to update, return existing {{.Entity.Name}}
		response := existing.ToResponse()
		return &response, nil
	}

	// Validate the fields being updated
	if err := s.validatePatchRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.repo.Patch(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to patch {{.Entity.Name}}: %w", err)
	}

	// Get updated {{.Entity.Name}} to return
	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated {{.Entity.Name}}: %w", err)
	}

	response := updated.ToResponse()
	return &response, nil
}
{{end}}

// Delete deletes a {{.Entity.Name}} by ID
func (s *service) Delete(ctx context.Context, id {{if eq .DatabaseType "mongodb"}}string{{else}}int64{{end}}) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete {{.Entity.Name}}: %w", err)
	}
	return nil
}

// Validation methods

func (s *service) validateCreateRequest(req Create{{title .Entity.Name}}Request) error {
{{range .Entity.Fields}}{{if .Required}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}	if {{if eq .Type "string"}}req.{{.Name}} == ""{{else if eq .Type "int"}}req.{{.Name}} == 0{{else if eq .Type "int64"}}req.{{.Name}} == 0{{else}}req.{{.Name}} == nil{{end}} {
		return fmt.Errorf("{{.Name}} is required")
	}
{{end}}{{end}}{{end}}{{end}}
	return nil
}

{{if or (eq .Entity.UpdateMethod "put") (eq .Entity.UpdateMethod "both")}}
func (s *service) validateUpdateRequest(req Update{{title .Entity.Name}}Request) error {
	// For PUT requests, all required fields must be provided (complete replacement)
{{range .Entity.Fields}}{{if .Required}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}	if {{if eq .Type "string"}}req.{{.Name}} == ""{{else if eq .Type "int"}}req.{{.Name}} == 0{{else if eq .Type "int64"}}req.{{.Name}} == 0{{else}}req.{{.Name}} == nil{{end}} {
		return fmt.Errorf("{{.Name}} is required for complete update")
	}
{{end}}{{end}}{{end}}{{end}}
	return nil
}
{{end}}

{{if or (eq .Entity.UpdateMethod "patch") (eq .Entity.UpdateMethod "both")}}
func (s *service) validatePatchRequest(req Patch{{title .Entity.Name}}Request) error {
	// For PATCH requests, only validate the fields that are being updated
{{range .Entity.Fields}}{{if .Required}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}	if req.{{.Name}} != nil && {{if eq .Type "string"}}*req.{{.Name}} == ""{{else if eq .Type "int"}}*req.{{.Name}} == 0{{else if eq .Type "int64"}}*req.{{.Name}} == 0{{else}}*req.{{.Name}} == nil{{end}} {
		return fmt.Errorf("{{.Name}} cannot be empty when provided")
	}
{{end}}{{end}}{{end}}{{end}}
	return nil
}
{{end}}
`

	filePath := filepath.Join(projectPath, "internal", "domain", data.Entity.Name, "service.go")
	return executeTemplate(tmpl, filePath, data)
}

// generateHandlerFile generates the HTTP handler file
func generateHandlerFile(projectPath string, data *CRUDTemplateData) error {
	tmpl := `package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"{{.ModuleName}}/internal/domain/{{.Entity.Name}}"
	"{{.ModuleName}}/internal/api/responses"
)

// {{title .Entity.Name}}Handler handles HTTP requests for {{.Entity.Name}} operations
type {{title .Entity.Name}}Handler struct {
	service {{.Entity.Name}}.Service
}

// New{{title .Entity.Name}}Handler creates a new {{.Entity.Name}} handler
func New{{title .Entity.Name}}Handler(service {{.Entity.Name}}.Service) *{{title .Entity.Name}}Handler {
	return &{{title .Entity.Name}}Handler{service: service}
}

// Create{{title .Entity.Name}} handles POST /api/{{.Entity.PluralName}}
func (h *{{title .Entity.Name}}Handler) Create{{title .Entity.Name}}(w http.ResponseWriter, r *http.Request) {
	var req {{.Entity.Name}}.Create{{title .Entity.Name}}Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	{{.Entity.Name}}Response, err := h.service.Create(r.Context(), req)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to create {{.Entity.Name}}", err)
		return
	}

	responses.Success(w, http.StatusCreated, "{{title .Entity.Name}} created successfully", {{.Entity.Name}}Response)
}

// Get{{title .Entity.Name}} handles GET /api/{{.Entity.PluralName}}/{id}
func (h *{{title .Entity.Name}}Handler) Get{{title .Entity.Name}}(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

{{if eq .DatabaseType "mongodb"}}	{{.Entity.Name}}Response, err := h.service.GetByID(r.Context(), idStr){{else}}	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	{{.Entity.Name}}Response, err := h.service.GetByID(r.Context(), id){{end}}
	if err != nil {
		responses.Error(w, http.StatusNotFound, "{{title .Entity.Name}} not found", err)
		return
	}

	responses.Success(w, http.StatusOK, "{{title .Entity.Name}} retrieved successfully", {{.Entity.Name}}Response)
}

// List{{title .Entity.PluralName}} handles GET /api/{{.Entity.PluralName}}
func (h *{{title .Entity.Name}}Handler) List{{title .Entity.PluralName}}(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page := 1
	pageSize := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	{{.Entity.PluralName}}Response, err := h.service.List(r.Context(), page, pageSize)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to list {{.Entity.PluralName}}", err)
		return
	}

	responses.Success(w, http.StatusOK, "{{title .Entity.PluralName}} retrieved successfully", {{.Entity.PluralName}}Response)
}

{{if or (eq .Entity.UpdateMethod "put") (eq .Entity.UpdateMethod "both")}}
// Update{{title .Entity.Name}} handles PUT /api/{{.Entity.PluralName}}/{id}
// PUT performs a complete replacement of the resource - all fields must be provided
func (h *{{title .Entity.Name}}Handler) Update{{title .Entity.Name}}(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

{{if eq .DatabaseType "mongodb"}}	var req {{.Entity.Name}}.Update{{title .Entity.Name}}Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	{{.Entity.Name}}Response, err := h.service.Update(r.Context(), idStr, req){{else}}	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	var req {{.Entity.Name}}.Update{{title .Entity.Name}}Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	{{.Entity.Name}}Response, err := h.service.Update(r.Context(), id, req){{end}}
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to update {{.Entity.Name}}", err)
		return
	}

	responses.Success(w, http.StatusOK, "{{title .Entity.Name}} updated successfully", {{.Entity.Name}}Response)
}
{{end}}

{{if or (eq .Entity.UpdateMethod "patch") (eq .Entity.UpdateMethod "both")}}
// Patch{{title .Entity.Name}} handles PATCH /api/{{.Entity.PluralName}}/{id}
// PATCH performs a partial update - only provided fields will be updated
func (h *{{title .Entity.Name}}Handler) Patch{{title .Entity.Name}}(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

{{if eq .DatabaseType "mongodb"}}	var req {{.Entity.Name}}.Patch{{title .Entity.Name}}Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	{{.Entity.Name}}Response, err := h.service.Patch(r.Context(), idStr, req){{else}}	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	var req {{.Entity.Name}}.Patch{{title .Entity.Name}}Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	{{.Entity.Name}}Response, err := h.service.Patch(r.Context(), id, req){{end}}
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to patch {{.Entity.Name}}", err)
		return
	}

	responses.Success(w, http.StatusOK, "{{title .Entity.Name}} patched successfully", {{.Entity.Name}}Response)
}
{{end}}

// Delete{{title .Entity.Name}} handles DELETE /api/{{.Entity.PluralName}}/{id}
func (h *{{title .Entity.Name}}Handler) Delete{{title .Entity.Name}}(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

{{if eq .DatabaseType "mongodb"}}	err := h.service.Delete(r.Context(), idStr){{else}}	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	err = h.service.Delete(r.Context(), id){{end}}
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to delete {{.Entity.Name}}", err)
		return
	}

	responses.Success(w, http.StatusOK, "{{title .Entity.Name}} deleted successfully", nil)
}
`

	filePath := filepath.Join(projectPath, "internal", "api", "handlers", data.Entity.Name+".go")
	return executeTemplate(tmpl, filePath, data)
}

func updateRoutesFile(projectPath string, data *CRUDTemplateData) error {
	routesPath := filepath.Join(projectPath, "internal", "api", "routes", "routes.go")

	// Check if routes file exists
	if _, err := os.Stat(routesPath); os.IsNotExist(err) {
		// Create a basic routes file if it doesn't exist
		return createRoutesFile(projectPath, data)
	}

	// For now, just create a comment about manual route addition
	// In a full implementation, this would parse and modify the existing routes.go file
	fmt.Printf("📝 Please add the following routes to your routes.go file:\n")
	fmt.Printf("   router.HandleFunc(\"/api/%s\", %sHandler.Create%s).Methods(\"POST\")\n",
		data.Entity.PluralName, data.Entity.Name, strings.Title(data.Entity.Name))
	fmt.Printf("   router.HandleFunc(\"/api/%s\", %sHandler.List%s).Methods(\"GET\")\n",
		data.Entity.PluralName, data.Entity.Name, strings.Title(data.Entity.PluralName))
	fmt.Printf("   router.HandleFunc(\"/api/%s/{id}\", %sHandler.Get%s).Methods(\"GET\")\n",
		data.Entity.PluralName, data.Entity.Name, strings.Title(data.Entity.Name))

	switch data.Entity.UpdateMethod {
	case "put":
		fmt.Printf("   router.HandleFunc(\"/api/%s/{id}\", %sHandler.Update%s).Methods(\"PUT\")\n",
			data.Entity.PluralName, data.Entity.Name, strings.Title(data.Entity.Name))
	case "patch":
		fmt.Printf("   router.HandleFunc(\"/api/%s/{id}\", %sHandler.Patch%s).Methods(\"PATCH\")\n",
			data.Entity.PluralName, data.Entity.Name, strings.Title(data.Entity.Name))
	case "both":
		fmt.Printf("   router.HandleFunc(\"/api/%s/{id}\", %sHandler.Update%s).Methods(\"PUT\")\n",
			data.Entity.PluralName, data.Entity.Name, strings.Title(data.Entity.Name))
		fmt.Printf("   router.HandleFunc(\"/api/%s/{id}\", %sHandler.Patch%s).Methods(\"PATCH\")\n",
			data.Entity.PluralName, data.Entity.Name, strings.Title(data.Entity.Name))
	}

	fmt.Printf("   router.HandleFunc(\"/api/%s/{id}\", %sHandler.Delete%s).Methods(\"DELETE\")\n",
		data.Entity.PluralName, data.Entity.Name, strings.Title(data.Entity.Name))
	fmt.Println()

	return nil
}

func createRoutesFile(projectPath string, data *CRUDTemplateData) error {
	tmpl := `package routes

import (
	"github.com/gorilla/mux"
	"{{.ModuleName}}/internal/api/handlers"
	"{{.ModuleName}}/internal/domain/{{.Entity.Name}}"
)

// SetupRoutes configures all API routes
func SetupRoutes({{.Entity.Name}}Service {{.Entity.Name}}.Service) *mux.Router {
	router := mux.NewRouter()
	
	// Initialize handlers
	{{.Entity.Name}}Handler := handlers.New{{title .Entity.Name}}Handler({{.Entity.Name}}Service)
	
	// {{title .Entity.Name}} routes
	router.HandleFunc("/api/{{.Entity.PluralName}}", {{.Entity.Name}}Handler.Create{{title .Entity.Name}}).Methods("POST")
	router.HandleFunc("/api/{{.Entity.PluralName}}", {{.Entity.Name}}Handler.List{{title .Entity.PluralName}}).Methods("GET")
	router.HandleFunc("/api/{{.Entity.PluralName}}/{id}", {{.Entity.Name}}Handler.Get{{title .Entity.Name}}).Methods("GET")
{{if or (eq .Entity.UpdateMethod "put") (eq .Entity.UpdateMethod "both")}}	router.HandleFunc("/api/{{.Entity.PluralName}}/{id}", {{.Entity.Name}}Handler.Update{{title .Entity.Name}}).Methods("PUT"){{end}}
{{if or (eq .Entity.UpdateMethod "patch") (eq .Entity.UpdateMethod "both")}}	router.HandleFunc("/api/{{.Entity.PluralName}}/{id}", {{.Entity.Name}}Handler.Patch{{title .Entity.Name}}).Methods("PATCH"){{end}}
	router.HandleFunc("/api/{{.Entity.PluralName}}/{id}", {{.Entity.Name}}Handler.Delete{{title .Entity.Name}}).Methods("DELETE")
	
	return router
}
`

	routesDir := filepath.Join(projectPath, "internal", "api", "routes")
	if err := os.MkdirAll(routesDir, 0755); err != nil {
		return fmt.Errorf("failed to create routes directory: %w", err)
	}

	filePath := filepath.Join(routesDir, "routes.go")
	return executeTemplate(tmpl, filePath, data)
}

func generateMigrationFiles(projectPath string, data *CRUDTemplateData) error {
	if data.DatabaseType == "mongodb" {
		// MongoDB doesn't need migrations, but we can create an initialization script
		return generateMongoInitScript(projectPath, data)
	}

	// Generate SQL migration files
	return generateSQLMigration(projectPath, data)
}

func generateSQLMigration(projectPath string, data *CRUDTemplateData) error {
	timestamp := time.Now().Format("20060102150405")

	// Up migration
	upTmpl := `-- Create {{.Entity.PluralName}} table
CREATE TABLE {{.Entity.PluralName}} (
    id SERIAL PRIMARY KEY,
{{range .Entity.Fields}}    {{.DBTag}} {{getSQLType .Type}}{{if .Required}} NOT NULL{{end}}{{if .Unique}} UNIQUE{{end}},
{{end}}    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
{{range .Entity.Fields}}{{if .Unique}}CREATE UNIQUE INDEX idx_{{$.Entity.PluralName}}_{{.DBTag}} ON {{$.Entity.PluralName}}({{.DBTag}});
{{end}}{{end}}

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_{{.Entity.PluralName}}_updated_at 
    BEFORE UPDATE ON {{.Entity.PluralName}} 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
`

	// Down migration
	downTmpl := `-- Drop {{.Entity.PluralName}} table
DROP TRIGGER IF EXISTS update_{{.Entity.PluralName}}_updated_at ON {{.Entity.PluralName}};
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS {{.Entity.PluralName}};
`

	// Create migration files
	migrationDir := filepath.Join(projectPath, "migrations")
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	upFile := filepath.Join(migrationDir, fmt.Sprintf("%s_create_%s_table.up.sql", timestamp, data.Entity.PluralName))
	downFile := filepath.Join(migrationDir, fmt.Sprintf("%s_create_%s_table.down.sql", timestamp, data.Entity.PluralName))

	// Add SQL type mapping function
	funcMap := template.FuncMap{
		"title": strings.Title,
		"getSQLType": func(goType string) string {
			switch goType {
			case "string":
				return "VARCHAR(255)"
			case "int", "int32":
				return "INTEGER"
			case "int64":
				return "BIGINT"
			case "float64":
				return "DECIMAL(10,2)"
			case "bool":
				return "BOOLEAN"
			case "time.Time":
				return "TIMESTAMP"
			case "[]string":
				return "TEXT[]"
			default:
				return "TEXT"
			}
		},
	}

	// Execute up migration template
	upTemplate, err := template.New("up").Funcs(funcMap).Parse(upTmpl)
	if err != nil {
		return fmt.Errorf("failed to parse up migration template: %w", err)
	}

	upFileHandle, err := os.Create(upFile)
	if err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}
	defer upFileHandle.Close()

	if err := upTemplate.Execute(upFileHandle, data); err != nil {
		return fmt.Errorf("failed to execute up migration template: %w", err)
	}

	// Execute down migration template
	downTemplate, err := template.New("down").Funcs(funcMap).Parse(downTmpl)
	if err != nil {
		return fmt.Errorf("failed to parse down migration template: %w", err)
	}

	downFileHandle, err := os.Create(downFile)
	if err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}
	defer downFileHandle.Close()

	if err := downTemplate.Execute(downFileHandle, data); err != nil {
		return fmt.Errorf("failed to execute down migration template: %w", err)
	}

	return nil
}

func generateMongoInitScript(projectPath string, data *CRUDTemplateData) error {
	tmpl := `// MongoDB initialization script for {{.Entity.PluralName}} collection
// Run this script in MongoDB shell or use it as reference

use {{.ProjectName}};

// Create {{.Entity.PluralName}} collection with validation
db.createCollection("{{.Entity.PluralName}}", {
   validator: {
      $jsonSchema: {
         bsonType: "object",
         required: [{{range .Entity.Fields}}{{if .Required}}"{{.JSONTag}}", {{end}}{{end}}],
         properties: {
{{range .Entity.Fields}}            {{.JSONTag}}: {
               bsonType: "{{getMongoType .Type}}",
               description: "{{.Name}} field"
            },
{{end}}         }
      }
   }
});

// Create indexes
{{range .Entity.Fields}}{{if .Unique}}db.{{$.Entity.PluralName}}.createIndex({ "{{.JSONTag}}": 1 }, { unique: true });
{{end}}{{end}}

// Create compound indexes if needed
// db.{{.Entity.PluralName}}.createIndex({ "field1": 1, "field2": 1 });

console.log("{{title .Entity.PluralName}} collection initialized successfully");
`

	funcMap := template.FuncMap{
		"title": strings.Title,
		"getMongoType": func(goType string) string {
			switch goType {
			case "string":
				return "string"
			case "int", "int32", "int64":
				return "int"
			case "float64":
				return "double"
			case "bool":
				return "bool"
			case "time.Time":
				return "date"
			case "[]string":
				return "array"
			default:
				return "string"
			}
		},
	}

	template, err := template.New("mongo").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse mongo template: %w", err)
	}

	migrationDir := filepath.Join(projectPath, "migrations")
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	filePath := filepath.Join(migrationDir, fmt.Sprintf("mongodb_init_%s.js", data.Entity.PluralName))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create mongo init file: %w", err)
	}
	defer file.Close()

	if err := template.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute mongo template: %w", err)
	}

	return nil
}

func generateDocumentation(projectPath string, data *CRUDTemplateData) error {
	tmpl := `# {{title .Entity.Name}} CRUD Operations

This document provides comprehensive information about the {{title .Entity.Name}} CRUD operations generated by Gophex.

## Overview

The {{title .Entity.Name}} entity includes the following operations:
- **Create**: Add new {{.Entity.PluralName}}
- **Read**: Retrieve {{.Entity.PluralName}} (single or list)
{{if or (eq .Entity.UpdateMethod "put") (eq .Entity.UpdateMethod "both")}}
- **Update (PUT)**: Complete replacement of {{.Entity.Name}} resource
{{end}}
{{if or (eq .Entity.UpdateMethod "patch") (eq .Entity.UpdateMethod "both")}}
- **Update (PATCH)**: Partial update of {{.Entity.Name}} resource
{{end}}
- **Delete**: Remove {{.Entity.PluralName}}

## API Endpoints

### Create {{title .Entity.Name}}
` + "```" + `
POST /api/{{.Entity.PluralName}}
Content-Type: application/json

{
{{range .Entity.Fields}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}  "{{.JSONTag}}": {{getExampleValue .Type}}{{if .Required}} // Required{{end}},
{{end}}{{end}}{{end}}
}
` + "```" + `

**Response (201 Created):**
` + "```json" + `
{
  "success": true,
  "message": "{{title .Entity.Name}} created successfully",
  "data": {
    "id": {{if eq .DatabaseType "mongodb"}}"507f1f77bcf86cd799439011"{{else}}1{{end}},
{{range .Entity.Fields}}    "{{.JSONTag}}": {{getExampleValue .Type}},
{{end}}
  }
}
` + "```" + `

### Get {{title .Entity.Name}} by ID
` + "```" + `
GET /api/{{.Entity.PluralName}}/{id}
` + "```" + `

**Response (200 OK):**
` + "```json" + `
{
  "success": true,
  "message": "{{title .Entity.Name}} retrieved successfully",
  "data": {
    "id": {{if eq .DatabaseType "mongodb"}}"507f1f77bcf86cd799439011"{{else}}1{{end}},
{{range .Entity.Fields}}    "{{.JSONTag}}": {{getExampleValue .Type}},
{{end}}
  }
}
` + "```" + `

### List {{title .Entity.PluralName}}
` + "```" + `
GET /api/{{.Entity.PluralName}}?page=1&page_size=10
` + "```" + `

**Query Parameters:**
- ` + "`page`" + `: Page number (default: 1)
- ` + "`page_size`" + `: Items per page (default: 10, max: 100)

**Response (200 OK):**
` + "```json" + `
{
  "success": true,
  "message": "{{title .Entity.PluralName}} retrieved successfully",
  "data": {
    "{{.Entity.PluralName}}": [
      {
        "id": {{if eq .DatabaseType "mongodb"}}"507f1f77bcf86cd799439011"{{else}}1{{end}},
{{range .Entity.Fields}}        "{{.JSONTag}}": {{getExampleValue .Type}},
{{end}}
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
` + "```" + `

{{if or (eq .Entity.UpdateMethod "put") (eq .Entity.UpdateMethod "both")}}
### Update {{title .Entity.Name}} (PUT - Complete Replacement)

**Important**: PUT replaces the entire resource. All required fields must be provided.

` + "```" + `
PUT /api/{{.Entity.PluralName}}/{id}
Content-Type: application/json

{
{{range .Entity.Fields}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}  "{{.JSONTag}}": {{getExampleValue .Type}}{{if .Required}} // Required for PUT{{end}},
{{end}}{{end}}{{end}}
}
` + "```" + `

**Response (200 OK):**
` + "```json" + `
{
  "success": true,
  "message": "{{title .Entity.Name}} updated successfully",
  "data": {
    "id": {{if eq .DatabaseType "mongodb"}}"507f1f77bcf86cd799439011"{{else}}1{{end}},
{{range .Entity.Fields}}    "{{.JSONTag}}": {{getExampleValue .Type}},
{{end}}
  }
}
` + "```" + `
{{end}}

{{if or (eq .Entity.UpdateMethod "patch") (eq .Entity.UpdateMethod "both")}}
### Update {{title .Entity.Name}} (PATCH - Partial Update)

**Important**: PATCH updates only the provided fields. Missing fields remain unchanged.

` + "```" + `
PATCH /api/{{.Entity.PluralName}}/{id}
Content-Type: application/json

{
  // Only include fields you want to update
{{range .Entity.Fields}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}  "{{.JSONTag}}": {{getExampleValue .Type}}, // Optional
{{end}}{{end}}{{end}}
}
` + "```" + `

**Example - Update only email:**
` + "```json" + `
{
  "email": "newemail@example.com"
}
` + "```" + `

**Response (200 OK):**
` + "```json" + `
{
  "success": true,
  "message": "{{title .Entity.Name}} patched successfully",
  "data": {
    "id": {{if eq .DatabaseType "mongodb"}}"507f1f77bcf86cd799439011"{{else}}1{{end}},
{{range .Entity.Fields}}    "{{.JSONTag}}": {{getExampleValue .Type}},
{{end}}
  }
}
` + "```" + `
{{end}}

### Delete {{title .Entity.Name}}
` + "```" + `
DELETE /api/{{.Entity.PluralName}}/{id}
` + "```" + `

**Response (200 OK):**
` + "```json" + `
{
  "success": true,
  "message": "{{title .Entity.Name}} deleted successfully",
  "data": null
}
` + "```" + `

{{if or (eq .Entity.UpdateMethod "both")}}
## PUT vs PATCH: When to Use Which?

### Use PUT when:
- You have the complete resource data
- You want to ensure data consistency
- You're implementing a "save" or "replace" operation
- You want to reset fields to default values

### Use PATCH when:
- You only want to update specific fields
- You're implementing incremental updates
- You want to preserve existing field values
- You're building forms that update individual fields

### Example Scenarios:

**PUT Example - User Profile Update:**
` + "```json" + `
// Client sends complete user profile
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "address": "123 Main St"
}
// All fields are replaced, missing fields would be set to null/default
` + "```" + `

**PATCH Example - Status Update:**
` + "```json" + `
// Client only wants to update status
{
  "status": "active"
}
// Only status is updated, all other fields remain unchanged
` + "```" + `
{{end}}

## Error Responses

All endpoints return consistent error responses:

**400 Bad Request:**
` + "```json" + `
{
  "success": false,
  "message": "Validation failed",
  "error": "Name is required"
}
` + "```" + `

**404 Not Found:**
` + "```json" + `
{
  "success": false,
  "message": "{{title .Entity.Name}} not found",
  "error": "{{.Entity.Name}} with ID 123 not found"
}
` + "```" + `

**500 Internal Server Error:**
` + "```json" + `
{
  "success": false,
  "message": "Internal server error",
  "error": "Database connection failed"
}
` + "```" + `

## Testing with curl

### Create a new {{.Entity.Name}}:
` + "```bash" + `
curl -X POST http://localhost:8080/api/{{.Entity.PluralName}} \
  -H "Content-Type: application/json" \
  -d '{
{{range .Entity.Fields}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}    "{{.JSONTag}}": {{getExampleValue .Type}}{{if not (isLast $.Entity.Fields .)}},{{end}}
{{end}}{{end}}{{end}}  }'
` + "```" + `

### Get all {{.Entity.PluralName}}:
` + "```bash" + `
curl http://localhost:8080/api/{{.Entity.PluralName}}
` + "```" + `

### Get a specific {{.Entity.Name}}:
` + "```bash" + `
curl http://localhost:8080/api/{{.Entity.PluralName}}/1
` + "```" + `

{{if or (eq .Entity.UpdateMethod "put") (eq .Entity.UpdateMethod "both")}}
### Update a {{.Entity.Name}} (PUT):
` + "```bash" + `
curl -X PUT http://localhost:8080/api/{{.Entity.PluralName}}/1 \
  -H "Content-Type: application/json" \
  -d '{
{{range .Entity.Fields}}{{if not (eq .Name "CreatedAt")}}{{if not (eq .Name "UpdatedAt")}}    "{{.JSONTag}}": {{getExampleValue .Type}}{{if not (isLast $.Entity.Fields .)}},{{end}}
{{end}}{{end}}{{end}}  }'
` + "```" + `
{{end}}

{{if or (eq .Entity.UpdateMethod "patch") (eq .Entity.UpdateMethod "both")}}
### Partially update a {{.Entity.Name}} (PATCH):
` + "```bash" + `
curl -X PATCH http://localhost:8080/api/{{.Entity.PluralName}}/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "updated"}'
` + "```" + `
{{end}}

### Delete a {{.Entity.Name}}:
` + "```bash" + `
curl -X DELETE http://localhost:8080/api/{{.Entity.PluralName}}/1
` + "```" + `

## Database Schema

{{if eq .DatabaseType "mongodb"}}
### MongoDB Collection: {{.Entity.PluralName}}

` + "```javascript" + `
{
  "_id": ObjectId("507f1f77bcf86cd799439011"),
{{range .Entity.Fields}}  "{{.JSONTag}}": {{getExampleValue .Type}},
{{end}}
}
` + "```" + `

### Indexes:
{{range .Entity.Fields}}{{if .Unique}}
- Unique index on ` + "`{{.JSONTag}}`" + `
{{end}}{{end}}

{{else}}
### PostgreSQL Table: {{.Entity.PluralName}}

` + "```sql" + `
CREATE TABLE {{.Entity.PluralName}} (
    id SERIAL PRIMARY KEY,
{{range .Entity.Fields}}    {{.DBTag}} {{getSQLType .Type}}{{if .Required}} NOT NULL{{end}}{{if .Unique}} UNIQUE{{end}},
{{end}}    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
` + "```" + `

### Indexes:
{{range .Entity.Fields}}{{if .Unique}}
- Unique index on ` + "`{{.DBTag}}`" + `
{{end}}{{end}}
{{end}}

## Next Steps

1. **Run Database Migrations**: Execute the generated migration files
2. **Start Your Server**: Run your API server
3. **Test the Endpoints**: Use the curl examples above
4. **Customize Business Logic**: Modify the service layer for your specific needs
5. **Add Authentication**: Integrate with your authentication middleware
6. **Add Validation**: Enhance field validation in the service layer

## File Structure

The CRUD generation created the following files:

` + "```" + `
internal/domain/{{.Entity.Name}}/
├── model.go       # Data models and request/response structs
├── repository.go  # Database operations
└── service.go     # Business logic

internal/api/handlers/
└── {{.Entity.Name}}.go  # HTTP handlers

migrations/
{{if eq .DatabaseType "mongodb"}}└── mongodb_init_{{.Entity.PluralName}}.js  # MongoDB initialization{{else}}├── [timestamp]_create_{{.Entity.PluralName}}_table.up.sql
└── [timestamp]_create_{{.Entity.PluralName}}_table.down.sql{{end}}
` + "```" + `

Generated on: {{.Timestamp}}
`

	funcMap := template.FuncMap{
		"title": strings.Title,
		"getExampleValue": func(goType string) string {
			switch goType {
			case "string":
				return `"example"`
			case "int", "int32":
				return "123"
			case "int64":
				return "123"
			case "float64":
				return "99.99"
			case "bool":
				return "true"
			case "time.Time":
				return `"2023-01-01T00:00:00Z"`
			case "[]string":
				return `["item1", "item2"]`
			default:
				return `"example"`
			}
		},
		"getSQLType": func(goType string) string {
			switch goType {
			case "string":
				return "VARCHAR(255)"
			case "int", "int32":
				return "INTEGER"
			case "int64":
				return "BIGINT"
			case "float64":
				return "DECIMAL(10,2)"
			case "bool":
				return "BOOLEAN"
			case "time.Time":
				return "TIMESTAMP"
			case "[]string":
				return "TEXT[]"
			default:
				return "TEXT"
			}
		},
		"isLast": func(fields []CRUDField, current CRUDField) bool {
			for i, field := range fields {
				if field.Name == current.Name {
					return i == len(fields)-1
				}
			}
			return false
		},
	}

	template, err := template.New("docs").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse documentation template: %w", err)
	}

	filePath := filepath.Join(projectPath, fmt.Sprintf("README_%s.md", data.Entity.Name))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create documentation file: %w", err)
	}
	defer file.Close()

	if err := template.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute documentation template: %w", err)
	}

	return nil
}

func showNextSteps(entity *CRUDEntity) {
	fmt.Println("🎉 Next Steps:")
	fmt.Printf("1. Run database migrations: `make migrate` or `./scripts/migrate.sh`\n")
	fmt.Printf("2. Start your server: `go run cmd/api/main.go`\n")
	fmt.Printf("3. Test your API endpoints:\n")
	fmt.Printf("   - POST   /api/%s     (Create)\n", entity.PluralName)
	fmt.Printf("   - GET    /api/%s     (List)\n", entity.PluralName)
	fmt.Printf("   - GET    /api/%s/{id} (Get by ID)\n", entity.PluralName)

	switch entity.UpdateMethod {
	case "put":
		fmt.Printf("   - PUT    /api/%s/{id} (Complete update)\n", entity.PluralName)
	case "patch":
		fmt.Printf("   - PATCH  /api/%s/{id} (Partial update)\n", entity.PluralName)
	case "both":
		fmt.Printf("   - PUT    /api/%s/{id} (Complete update)\n", entity.PluralName)
		fmt.Printf("   - PATCH  /api/%s/{id} (Partial update)\n", entity.PluralName)
	}

	fmt.Printf("   - DELETE /api/%s/{id} (Delete)\n", entity.PluralName)
	fmt.Printf("4. Check README_%s.md for detailed examples and documentation\n\n", entity.Name)
}
