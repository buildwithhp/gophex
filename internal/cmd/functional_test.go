package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/buildwithhp/gophex/internal/generator"
	"github.com/buildwithhp/gophex/internal/utils"
)

// TestProjectGeneration tests the complete project generation workflow
func TestProjectGeneration(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gophex-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		projectType string
		projectName string
		framework   string
		expectFiles []string
	}{
		{
			name:        "API Project Generation - Gin",
			projectType: "api",
			projectName: "test-api-gin",
			framework:   "gin",
			expectFiles: []string{
				"go.mod",
				"cmd/api/main.go",
				"internal/config/config.go",
				"internal/api/handlers/health.go",
				"internal/api/handlers/auth.go",
				"internal/api/handlers/users.go",
				"internal/api/handlers/posts.go",
				"internal/api/middleware/cors.go",
				"internal/api/middleware/logging.go",
				"internal/api/routes/routes.go",
				"internal/database/database.go",
				"migrations",
				"README.md",
				"gophex.md",
			},
		},
		{
			name:        "API Project Generation - Echo",
			projectType: "api",
			projectName: "test-api-echo",
			framework:   "echo",
			expectFiles: []string{
				"go.mod",
				"cmd/api/main.go",
				"internal/config/config.go",
				"internal/api/handlers/health.go",
				"internal/api/handlers/auth.go",
				"internal/api/handlers/users.go",
				"internal/api/handlers/posts.go",
				"internal/api/middleware/cors.go",
				"internal/api/middleware/logging.go",
				"internal/api/routes/routes.go",
				"internal/database/database.go",
				"migrations",
				"README.md",
				"gophex.md",
			},
		},
		{
			name:        "API Project Generation - Gorilla",
			projectType: "api",
			projectName: "test-api-gorilla",
			framework:   "gorilla",
			expectFiles: []string{
				"go.mod",
				"cmd/api/main.go",
				"internal/config/config.go",
				"internal/api/handlers/health.go",
				"internal/api/handlers/auth.go",
				"internal/api/handlers/users.go",
				"internal/api/handlers/posts.go",
				"internal/api/middleware/cors.go",
				"internal/api/middleware/logging.go",
				"internal/api/routes/routes.go",
				"internal/database/database.go",
				"migrations",
				"README.md",
				"gophex.md",
			},
		},
		{
			name:        "CLI Project Generation",
			projectType: "cli",
			projectName: "test-cli",
			framework:   "", // No framework for CLI
			expectFiles: []string{
				"go.mod",
				"cmd/main.go",
				"internal/cmd/root.go",
				"README.md",
				"gophex.md",
			},
		},
		{
			name:        "WebApp Project Generation",
			projectType: "webapp",
			projectName: "test-webapp",
			framework:   "", // No framework for WebApp
			expectFiles: []string{
				"go.mod",
				"cmd/webapp/main.go",
				"web/static/css/style.css",
				"web/templates/index.html",
				"README.md",
				"gophex.md",
			},
		},
		{
			name:        "Microservice Project Generation",
			projectType: "microservice",
			projectName: "test-microservice",
			framework:   "", // No framework for Microservice
			expectFiles: []string{
				"go.mod",
				"cmd/server/main.go",
				"internal/handlers/handlers.go",
				"README.md",
				"gophex.md",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			projectPath := filepath.Join(tempDir, test.projectName)

			// Generate project
			gen := generator.New()
			var err error
			if test.framework != "" && test.projectType == "api" {
				err = gen.GenerateWithFramework(test.projectType, test.projectName, projectPath, test.framework, nil, nil)
			} else {
				err = gen.Generate(test.projectType, test.projectName, projectPath)
			}
			if err != nil {
				t.Fatalf("Failed to generate %s project: %v", test.projectType, err)
			}

			// Verify expected files exist
			for _, expectedFile := range test.expectFiles {
				fullPath := filepath.Join(projectPath, expectedFile)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					t.Errorf("Expected file %s does not exist", expectedFile)
				}
			}

			// Verify metadata file exists and is valid
			metadataPath := filepath.Join(projectPath, "gophex.md")
			if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
				t.Error("Metadata file gophex.md does not exist")
			} else {
				// Test metadata loading
				metadata, err := utils.LoadMetadata(projectPath)
				if err != nil {
					t.Errorf("Failed to load metadata: %v", err)
				} else {
					if metadata.Project.Name != test.projectName {
						t.Errorf("Expected project name %s, got %s", test.projectName, metadata.Project.Name)
					}
					if metadata.Project.Type != test.projectType {
						t.Errorf("Expected project type %s, got %s", test.projectType, metadata.Project.Type)
					}
				}
			}

			// Verify go.mod file has correct module name
			goModPath := filepath.Join(projectPath, "go.mod")
			if content, err := os.ReadFile(goModPath); err == nil {
				if !strings.Contains(string(content), test.projectName) {
					t.Error("go.mod does not contain project name")
				}
			}
		})
	}
}

// TestCRUDGenerationWorkflow tests the complete CRUD generation workflow
func TestCRUDGenerationWorkflow(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gophex-crud-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	projectPath := filepath.Join(tempDir, "test-api")

	// First generate an API project
	gen := generator.New()
	err = gen.Generate("api", "test-api", projectPath)
	if err != nil {
		t.Fatalf("Failed to generate API project: %v", err)
	}

	// Test CRUD entity creation
	entity := &CRUDEntity{
		Name:         "user",
		PluralName:   "users",
		UpdateMethod: "both",
		Fields: []CRUDField{
			{Name: "Name", Type: "string", JSONTag: "name", DBTag: "name", Required: true},
			{Name: "Email", Type: "string", JSONTag: "email", DBTag: "email", Required: true, Unique: true},
			{Name: "CreatedAt", Type: "time.Time", JSONTag: "created_at", DBTag: "created_at"},
			{Name: "UpdatedAt", Type: "time.Time", JSONTag: "updated_at", DBTag: "updated_at"},
		},
	}

	// Test CRUD code generation
	err = generateCRUDCode(projectPath, entity)
	if err != nil {
		t.Fatalf("Failed to generate CRUD code: %v", err)
	}

	// Verify CRUD files were created
	expectedFiles := []string{
		"internal/domain/user/model.go",
		"internal/domain/user/repository.go",
		"internal/domain/user/service.go",
		"internal/api/handlers/user.go",
		"README_user.md",
	}

	for _, expectedFile := range expectedFiles {
		fullPath := filepath.Join(projectPath, expectedFile)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected CRUD file %s does not exist", expectedFile)
		}
	}

	// Verify model file contains expected structs
	modelPath := filepath.Join(projectPath, "internal/domain/user/model.go")
	if content, err := os.ReadFile(modelPath); err == nil {
		contentStr := string(content)
		expectedStructs := []string{
			"type User struct",
			"type CreateUserRequest struct",
			"type UpdateUserRequest struct",
			"type PatchUserRequest struct",
			"type UserResponse struct",
		}

		for _, expectedStruct := range expectedStructs {
			if !strings.Contains(contentStr, expectedStruct) {
				t.Errorf("Model file does not contain expected struct: %s", expectedStruct)
			}
		}
	}

	// Verify repository file contains expected methods
	repoPath := filepath.Join(projectPath, "internal/domain/user/repository.go")
	if content, err := os.ReadFile(repoPath); err == nil {
		contentStr := string(content)
		expectedMethods := []string{
			"Create(ctx context.Context",
			"GetByID(ctx context.Context",
			"List(ctx context.Context",
			"Update(ctx context.Context",
			"Patch(ctx context.Context",
			"Delete(ctx context.Context",
		}

		for _, expectedMethod := range expectedMethods {
			if !strings.Contains(contentStr, expectedMethod) {
				t.Errorf("Repository file does not contain expected method: %s", expectedMethod)
			}
		}
	}
}

// TestMetadataManagement tests metadata creation and management
func TestMetadataManagement(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gophex-metadata-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	projectPath := filepath.Join(tempDir, "test-project")

	// Generate project
	gen := generator.New()
	err = gen.Generate("api", "test-project", projectPath)
	if err != nil {
		t.Fatalf("Failed to generate project: %v", err)
	}

	// Test metadata loading
	metadata, err := utils.LoadMetadata(projectPath)
	if err != nil {
		t.Fatalf("Failed to load metadata: %v", err)
	}

	// Verify metadata structure
	if metadata.Project.Name != "test-project" {
		t.Errorf("Expected project name 'test-project', got %s", metadata.Project.Name)
	}

	if metadata.Project.Type != "api" {
		t.Errorf("Expected project type 'api', got %s", metadata.Project.Type)
	}

	// Test activity updates
	err = utils.UpdateActivity(projectPath, "test_activity", true)
	if err != nil {
		t.Fatalf("Failed to update activity: %v", err)
	}

	// Verify activity was updated
	if !utils.IsActivityCompleted(projectPath, "test_activity") {
		t.Error("Activity should be marked as completed")
	}

	// Test database status updates
	err = utils.UpdateDatabaseStatus(projectPath, true, true)
	if err != nil {
		t.Fatalf("Failed to update database status: %v", err)
	}

	// Reload metadata and verify updates
	updatedMetadata, err := utils.LoadMetadata(projectPath)
	if err != nil {
		t.Fatalf("Failed to reload metadata: %v", err)
	}

	if !updatedMetadata.Database.MigrationsExecuted {
		t.Error("Database migrations should be marked as executed")
	}

	if !updatedMetadata.Database.SchemaInitialized {
		t.Error("Database schema should be marked as initialized")
	}
}

// TestTemplateProcessing tests template processing functionality
func TestTemplateProcessing(t *testing.T) {

	// Test template functions
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Pluralize user", "user", "users"},
		{"Pluralize category", "category", "categories"},
		{"Pluralize box", "box", "boxes"},
		{"Title case", "hello world", "Hello World"},
		{"Lower case", "HELLO WORLD", "hello world"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result string
			switch test.name {
			case "Pluralize user", "Pluralize category", "Pluralize box":
				result = pluralize(test.input)
			case "Title case":
				result = strings.Title(test.input)
			case "Lower case":
				result = strings.ToLower(test.input)
			}

			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}

	// Test field validation
	validField := CRUDField{
		Name:     "TestField",
		Type:     "string",
		JSONTag:  "test_field",
		DBTag:    "test_field",
		Required: true,
	}

	if validField.Name == "" {
		t.Error("Field name should not be empty")
	}

	if validField.Type == "" {
		t.Error("Field type should not be empty")
	}
}

// TestProjectValidation tests project validation functionality
func TestProjectValidation(t *testing.T) {
	tests := []struct {
		name            string
		projectName     string
		projectType     string
		expectNameError bool
		expectTypeError bool
	}{
		{"Valid API project", "my-api", "api", false, false},
		{"Valid CLI project", "my-cli", "cli", false, false},
		{"Valid WebApp project", "my-webapp", "webapp", false, false},
		{"Valid Microservice project", "my-service", "microservice", false, false},
		{"Empty project name", "", "api", true, false},
		{"Invalid project type", "my-project", "invalid", false, true},
		{"Project name with spaces", "my project", "api", false, false},
		{"Project name with numbers", "project123", "api", false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Test project name validation
			if test.projectName != "" {
				nameValid := isValidProjectName(test.projectName)
				if test.expectNameError && nameValid {
					t.Error("Expected project name validation to fail")
				}
				if !test.expectNameError && !nameValid {
					t.Error("Expected project name validation to pass")
				}
			}

			// Test project type validation
			typeValid := isValidProjectType(test.projectType)
			if test.expectTypeError && typeValid {
				t.Error("Expected project type validation to fail")
			}
			if !test.expectTypeError && !typeValid {
				t.Error("Expected project type validation to pass")
			}
		})
	}
}

// TestEnvironmentVariableHandling tests environment variable functionality
func TestEnvironmentVariableHandling(t *testing.T) {
	tests := []struct {
		name         string
		envVar       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{"Environment variable exists", "TEST_VAR", "test_value", "default", "test_value"},
		{"Environment variable missing", "MISSING_VAR", "", "default", "default"},
		{"Empty environment variable", "EMPTY_VAR", "", "default", "default"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set environment variable if provided
			if test.envValue != "" {
				os.Setenv(test.envVar, test.envValue)
				defer os.Unsetenv(test.envVar)
			}

			result := utils.GetEnvWithDefault(test.envVar, test.defaultValue)
			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}

// Helper functions for testing

func isValidProjectName(name string) bool {
	return strings.TrimSpace(name) != ""
}

func isValidProjectType(projectType string) bool {
	validTypes := []string{"api", "webapp", "microservice", "cli"}
	for _, validType := range validTypes {
		if projectType == validType {
			return true
		}
	}
	return false
}

// TestBackwardCompatibility tests that existing functionality still works
func TestBackwardCompatibility(t *testing.T) {
	// Test that existing command structures still work
	tests := []struct {
		name     string
		function func() error
	}{
		{"Project generation workflow", func() error {
			tempDir, err := os.MkdirTemp("", "gophex-compat-test-*")
			if err != nil {
				return err
			}
			defer os.RemoveAll(tempDir)

			gen := generator.New()
			return gen.Generate("api", "compat-test", filepath.Join(tempDir, "compat-test"))
		}},
		{"Metadata utilities", func() error {
			tempDir, err := os.MkdirTemp("", "gophex-metadata-compat-*")
			if err != nil {
				return err
			}
			defer os.RemoveAll(tempDir)

			projectPath := filepath.Join(tempDir, "metadata-test")
			gen := generator.New()
			if err := gen.Generate("api", "metadata-test", projectPath); err != nil {
				return err
			}

			return utils.UpdateActivity(projectPath, "test_activity", true)
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.function(); err != nil {
				t.Errorf("Backward compatibility test failed: %v", err)
			}
		})
	}
}

// TestConcurrentOperations tests thread safety of operations
func TestConcurrentOperations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gophex-concurrent-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test concurrent project generation
	const numProjects = 5
	done := make(chan error, numProjects)

	for i := 0; i < numProjects; i++ {
		go func(index int) {
			projectName := fmt.Sprintf("concurrent-test-%d", index)
			projectPath := filepath.Join(tempDir, projectName)

			gen := generator.New()
			err := gen.Generate("api", projectName, projectPath)
			done <- err
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numProjects; i++ {
		if err := <-done; err != nil {
			t.Errorf("Concurrent project generation failed: %v", err)
		}
	}

	// Verify all projects were created
	for i := 0; i < numProjects; i++ {
		projectName := fmt.Sprintf("concurrent-test-%d", i)
		projectPath := filepath.Join(tempDir, projectName)

		if _, err := os.Stat(filepath.Join(projectPath, "gophex.md")); os.IsNotExist(err) {
			t.Errorf("Project %s was not created properly", projectName)
		}
	}
}

// TestErrorHandling tests error handling in various scenarios
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		operation   func() error
		expectError bool
	}{
		{
			"Generate project in non-existent directory",
			func() error {
				gen := generator.New()
				return gen.Generate("api", "test", "/non/existent/path/project")
			},
			true,
		},
		{
			"Load metadata from non-existent project",
			func() error {
				_, err := utils.LoadMetadata("/non/existent/project")
				return err
			},
			true,
		},
		{
			"Update activity in non-existent project",
			func() error {
				return utils.UpdateActivity("/non/existent/project", "test", true)
			},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.operation()
			if test.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !test.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
