package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/buildwithhp/gophex/internal/types"
)

// ProjectMetadata represents the complete metadata structure for a Gophex project
type ProjectMetadata struct {
	Project    ProjectInfo             `json:"project"`
	Hierarchy  map[string]interface{}  `json:"hierarchy"`
	Database   DatabaseMetadata        `json:"database"`
	Redis      RedisMetadata           `json:"redis"`
	Activities map[string]ActivityInfo `json:"activities"`
	Features   map[string]bool         `json:"features"`
	Endpoints  []EndpointInfo          `json:"endpoints,omitempty"`
	Commands   []CommandInfo           `json:"commands,omitempty"`
}

type ProjectInfo struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Version       string `json:"version"`
	GophexVersion string `json:"gophex_version"`
	GeneratedAt   string `json:"generated_at"`
	LastUpdated   string `json:"last_updated"`
}

type DatabaseMetadata struct {
	Configured         bool   `json:"configured"`
	Type               string `json:"type,omitempty"`
	ConfigType         string `json:"config_type,omitempty"`
	IsClustered        bool   `json:"is_clustered"`
	HasReadWriteSplit  bool   `json:"has_read_write_split"`
	SSLEnabled         bool   `json:"ssl_enabled"`
	MigrationsExecuted bool   `json:"migrations_executed"`
	SchemaInitialized  bool   `json:"schema_initialized"`
}

type RedisMetadata struct {
	Configured bool `json:"configured"`
	Enabled    bool `json:"enabled"`
}

type ActivityInfo struct {
	Completed bool   `json:"completed"`
	Timestamp string `json:"timestamp,omitempty"`
	CanRepeat bool   `json:"can_repeat"`
}

type EndpointInfo struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Protected   bool   `json:"protected"`
}

type CommandInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Subcommands []string `json:"subcommands"`
}

// MetadataGenerator handles dynamic generation of project metadata
type MetadataGenerator struct {
	projectPath string
	projectType string
}

// NewMetadataGenerator creates a new metadata generator
func NewMetadataGenerator(projectPath, projectType string) *MetadataGenerator {
	return &MetadataGenerator{
		projectPath: projectPath,
		projectType: projectType,
	}
}

// GenerateMetadata creates metadata by scanning the actual generated project
func (mg *MetadataGenerator) GenerateMetadata(projectName string, dbConfig *types.DatabaseConfig, redisConfig *types.RedisConfig, gophexVersion string) (*ProjectMetadata, error) {
	now := time.Now().Format(time.RFC3339)

	metadata := &ProjectMetadata{
		Project: ProjectInfo{
			Name:          projectName,
			Type:          mg.projectType,
			Version:       "1.0.0",
			GophexVersion: gophexVersion,
			GeneratedAt:   now,
			LastUpdated:   now,
		},
		Hierarchy:  make(map[string]interface{}),
		Activities: mg.generateDefaultActivities(now),
	}

	// Scan project structure
	hierarchy, err := mg.scanProjectHierarchy()
	if err != nil {
		return nil, fmt.Errorf("failed to scan project hierarchy: %w", err)
	}
	metadata.Hierarchy = hierarchy

	// Generate database metadata
	metadata.Database = mg.generateDatabaseMetadata(dbConfig)

	// Generate Redis metadata
	metadata.Redis = mg.generateRedisMetadata(redisConfig)

	// Generate features based on what's actually present
	metadata.Features = mg.scanFeatures()

	// Generate endpoints/commands based on project type
	switch mg.projectType {
	case "api":
		endpoints, err := mg.scanAPIEndpoints()
		if err != nil {
			return nil, fmt.Errorf("failed to scan API endpoints: %w", err)
		}
		metadata.Endpoints = endpoints
	case "cli":
		commands, err := mg.scanCLICommands()
		if err != nil {
			return nil, fmt.Errorf("failed to scan CLI commands: %w", err)
		}
		metadata.Commands = commands
	}

	return metadata, nil
}

// scanProjectHierarchy dynamically scans the project directory structure
func (mg *MetadataGenerator) scanProjectHierarchy() (map[string]interface{}, error) {
	hierarchy := make(map[string]interface{})

	err := filepath.Walk(mg.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories (except .env files)
		if strings.HasPrefix(info.Name(), ".") && !strings.HasPrefix(info.Name(), ".env") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path from project root
		relPath, err := filepath.Rel(mg.projectPath, path)
		if err != nil {
			return err
		}

		// Skip the root directory
		if relPath == "." {
			return nil
		}

		// Build nested structure
		mg.addToHierarchy(hierarchy, relPath, info.IsDir())
		return nil
	})

	return hierarchy, err
}

// addToHierarchy adds a file/directory to the hierarchy map
func (mg *MetadataGenerator) addToHierarchy(hierarchy map[string]interface{}, path string, isDir bool) {
	parts := strings.Split(filepath.ToSlash(path), "/")
	current := hierarchy

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - add the file/directory
			if isDir {
				if current[part] == nil {
					current[part] = make(map[string]interface{})
				}
			} else {
				current[part] = mg.getFileDescription(path)
			}
		} else {
			// Intermediate directory
			if current[part] == nil {
				current[part] = make(map[string]interface{})
			}
			if nested, ok := current[part].(map[string]interface{}); ok {
				current = nested
			}
		}
	}
}

// getFileDescription returns a description for a file based on its path and content
func (mg *MetadataGenerator) getFileDescription(filePath string) string {
	fileName := filepath.Base(filePath)
	dir := filepath.Dir(filePath)

	// Common file descriptions based on patterns
	descriptions := map[string]string{
		"main.go":           "application_entry_point",
		"go.mod":            "go_module_definition",
		"README.md":         "project_documentation",
		"gophex.md":         "project_metadata",
		".env":              "environment_variables",
		".env.example":      "environment_template",
		".gophex-generated": "generation_metadata",
		"config.go":         "configuration_management",
		"database.go":       "database_interface",
		"factory.go":        "database_factory",
		"connection.go":     "database_connection",
		"auth.go":           "authentication_logic",
		"health.go":         "health_check_handlers",
		"users.go":          "user_management",
		"posts.go":          "post_management",
		"routes.go":         "route_definitions",
		"middleware.go":     "middleware_functions",
		"cors.go":           "cors_middleware",
		"logging.go":        "logging_middleware",
		"ratelimit.go":      "rate_limiting_middleware",
		"error.go":          "error_handling",
		"success.go":        "success_responses",
		"jwt.go":            "jwt_implementation",
		"password.go":       "password_utilities",
		"validator.go":      "input_validation",
		"logger.go":         "logging_utilities",
		"errors.go":         "custom_error_types",
		"client.go":         "client_implementation",
		"handlers.go":       "request_handlers",
		"root.go":           "cli_root_command",
		"style.css":         "application_styles",
		"index.html":        "main_page_template",
	}

	if desc, exists := descriptions[fileName]; exists {
		return desc
	}

	// Pattern-based descriptions
	if strings.HasSuffix(fileName, "_repo.go") {
		return "repository_implementation"
	}
	if strings.HasSuffix(fileName, "_test.go") {
		return "test_file"
	}
	if strings.HasSuffix(fileName, ".sql") {
		if strings.Contains(fileName, ".up.") {
			return "database_migration"
		}
		if strings.Contains(fileName, ".down.") {
			return "migration_rollback"
		}
		return "sql_script"
	}
	if strings.HasSuffix(fileName, ".js") && strings.Contains(dir, "migrations") {
		return "mongodb_initialization"
	}
	if strings.HasSuffix(fileName, ".sh") {
		return "shell_script"
	}
	if strings.HasSuffix(fileName, ".bat") {
		return "batch_script"
	}

	// Directory-based descriptions
	if strings.Contains(dir, "handlers") {
		return "request_handler"
	}
	if strings.Contains(dir, "middleware") {
		return "middleware_component"
	}
	if strings.Contains(dir, "repository") || strings.Contains(dir, "repo") {
		return "repository_implementation"
	}
	if strings.Contains(dir, "service") {
		return "business_logic"
	}
	if strings.Contains(dir, "model") {
		return "data_model"
	}

	return "source_file"
}

// generateDatabaseMetadata creates database metadata based on configuration
func (mg *MetadataGenerator) generateDatabaseMetadata(dbConfig *types.DatabaseConfig) DatabaseMetadata {
	if dbConfig == nil {
		return DatabaseMetadata{
			Configured:         false,
			IsClustered:        false,
			HasReadWriteSplit:  false,
			SSLEnabled:         false,
			MigrationsExecuted: false,
			SchemaInitialized:  false,
		}
	}

	return DatabaseMetadata{
		Configured:         true,
		Type:               dbConfig.Type,
		ConfigType:         dbConfig.ConfigType,
		IsClustered:        dbConfig.ConfigType == "cluster",
		HasReadWriteSplit:  dbConfig.ConfigType == "read-write",
		SSLEnabled:         dbConfig.SSLMode != "" && dbConfig.SSLMode != "disable",
		MigrationsExecuted: false,
		SchemaInitialized:  false,
	}
}

// generateRedisMetadata creates Redis metadata based on configuration
func (mg *MetadataGenerator) generateRedisMetadata(redisConfig *types.RedisConfig) RedisMetadata {
	if redisConfig == nil {
		return RedisMetadata{
			Configured: false,
			Enabled:    false,
		}
	}

	return RedisMetadata{
		Configured: true,
		Enabled:    redisConfig.Enabled,
	}
}

// generateDefaultActivities creates the default activity tracking structure
func (mg *MetadataGenerator) generateDefaultActivities(timestamp string) map[string]ActivityInfo {
	activities := map[string]ActivityInfo{
		"project_generated": {
			Completed: true,
			Timestamp: timestamp,
			CanRepeat: false,
		},
		"dependencies_installed": {
			Completed: false,
			CanRepeat: true,
		},
		"tests_executed": {
			Completed: false,
			CanRepeat: true,
		},
		"project_opened": {
			Completed: false,
			CanRepeat: true,
		},
		"documentation_viewed": {
			Completed: false,
			CanRepeat: true,
		},
	}

	// Add project-type specific activities
	switch mg.projectType {
	case "api":
		activities["database_migrated"] = ActivityInfo{
			Completed: false,
			CanRepeat: true,
		}
		activities["application_started"] = ActivityInfo{
			Completed: false,
			CanRepeat: true,
		}
		activities["change_detection_run"] = ActivityInfo{
			Completed: false,
			CanRepeat: true,
		}
	case "webapp", "microservice":
		activities["application_started"] = ActivityInfo{
			Completed: false,
			CanRepeat: true,
		}
	case "cli":
		activities["application_built"] = ActivityInfo{
			Completed: false,
			CanRepeat: true,
		}
	}

	return activities
}

// scanFeatures detects which features are present in the generated project
func (mg *MetadataGenerator) scanFeatures() map[string]bool {
	features := make(map[string]bool)

	// Check for common features based on file existence
	featureFiles := map[string]string{
		"authentication":     "internal/api/handlers/auth.go",
		"user_management":    "internal/api/handlers/users.go",
		"post_management":    "internal/api/handlers/posts.go",
		"health_checks":      "internal/api/handlers/health.go",
		"cors_enabled":       "internal/api/middleware/cors.go",
		"rate_limiting":      "internal/api/middleware/ratelimit.go",
		"request_logging":    "internal/api/middleware/logging.go",
		"input_validation":   "internal/pkg/validator/validator.go",
		"structured_logging": "internal/pkg/logger/logger.go",
		"clean_architecture": "internal/domain",
		"graceful_shutdown":  "cmd",
		"web_server":         "web",
		"static_files":       "web/static",
		"html_templates":     "web/templates",
		"grpc_support":       "internal/handlers",
		"cobra_framework":    "internal/cmd/root.go",
	}

	for feature, path := range featureFiles {
		fullPath := filepath.Join(mg.projectPath, path)
		if _, err := os.Stat(fullPath); err == nil {
			features[feature] = true
		}
	}

	return features
}

// scanAPIEndpoints scans for API endpoints in handler files
func (mg *MetadataGenerator) scanAPIEndpoints() ([]EndpointInfo, error) {
	var endpoints []EndpointInfo

	// Default endpoints that are typically generated
	defaultEndpoints := []EndpointInfo{
		{Method: "GET", Path: "/api/v1/health", Description: "Health check endpoint", Protected: false},
	}

	// Check if auth handlers exist
	authPath := filepath.Join(mg.projectPath, "internal/api/handlers/auth.go")
	if _, err := os.Stat(authPath); err == nil {
		defaultEndpoints = append(defaultEndpoints,
			EndpointInfo{Method: "POST", Path: "/api/v1/auth/register", Description: "User registration", Protected: false},
			EndpointInfo{Method: "POST", Path: "/api/v1/auth/login", Description: "User login", Protected: false},
		)
	}

	// Check if user handlers exist
	userPath := filepath.Join(mg.projectPath, "internal/api/handlers/users.go")
	if _, err := os.Stat(userPath); err == nil {
		defaultEndpoints = append(defaultEndpoints,
			EndpointInfo{Method: "GET", Path: "/api/v1/users", Description: "List users", Protected: true},
			EndpointInfo{Method: "GET", Path: "/api/v1/users/{id}", Description: "Get user by ID", Protected: true},
			EndpointInfo{Method: "PUT", Path: "/api/v1/users/{id}", Description: "Update user", Protected: true},
			EndpointInfo{Method: "DELETE", Path: "/api/v1/users/{id}", Description: "Delete user", Protected: true},
		)
	}

	// Check if post handlers exist
	postPath := filepath.Join(mg.projectPath, "internal/api/handlers/posts.go")
	if _, err := os.Stat(postPath); err == nil {
		defaultEndpoints = append(defaultEndpoints,
			EndpointInfo{Method: "GET", Path: "/api/v1/posts", Description: "List posts", Protected: false},
			EndpointInfo{Method: "GET", Path: "/api/v1/posts/{id}", Description: "Get post by ID", Protected: false},
			EndpointInfo{Method: "POST", Path: "/api/v1/posts", Description: "Create post", Protected: true},
			EndpointInfo{Method: "PUT", Path: "/api/v1/posts/{id}", Description: "Update post", Protected: true},
			EndpointInfo{Method: "DELETE", Path: "/api/v1/posts/{id}", Description: "Delete post", Protected: true},
		)
	}

	endpoints = append(endpoints, defaultEndpoints...)

	// TODO: In the future, we could parse the actual route files to extract custom endpoints
	return endpoints, nil
}

// scanCLICommands scans for CLI commands in the project
func (mg *MetadataGenerator) scanCLICommands() ([]CommandInfo, error) {
	commands := []CommandInfo{
		{
			Name:        "root",
			Description: "Root command",
			Subcommands: []string{},
		},
	}

	// TODO: Parse actual command files to extract subcommands
	return commands, nil
}

// WriteMetadataFile writes the metadata to gophex.md file
func (mg *MetadataGenerator) WriteMetadataFile(metadata *ProjectMetadata) error {
	// Convert to JSON with proper formatting
	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Create the markdown content
	content := fmt.Sprintf(`# Gophex Project Metadata

This file contains project metadata and progress tracking for Gophex-generated projects.
**Do not edit this file manually** - it is automatically maintained by Gophex.

`+"```json\n%s\n```\n", string(jsonData))

	// Write to file
	filePath := filepath.Join(mg.projectPath, "gophex.md")
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// UpdateActivity updates the status of a specific activity
func UpdateActivity(projectPath, activityName string, completed bool) error {
	metadata, err := LoadMetadata(projectPath)
	if err != nil {
		return err
	}

	if metadata.Activities == nil {
		metadata.Activities = make(map[string]ActivityInfo)
	}

	activity := metadata.Activities[activityName]
	activity.Completed = completed
	if completed {
		activity.Timestamp = time.Now().Format(time.RFC3339)
	}
	metadata.Activities[activityName] = activity
	metadata.Project.LastUpdated = time.Now().Format(time.RFC3339)

	return SaveMetadata(projectPath, metadata)
}

// UpdateDatabaseStatus updates database-related status
func UpdateDatabaseStatus(projectPath string, migrationsExecuted, schemaInitialized bool) error {
	metadata, err := LoadMetadata(projectPath)
	if err != nil {
		return err
	}

	metadata.Database.MigrationsExecuted = migrationsExecuted
	metadata.Database.SchemaInitialized = schemaInitialized
	metadata.Project.LastUpdated = time.Now().Format(time.RFC3339)

	return SaveMetadata(projectPath, metadata)
}

// LoadMetadata loads metadata from gophex.md file
func LoadMetadata(projectPath string) (*ProjectMetadata, error) {
	filePath := filepath.Join(projectPath, "gophex.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	// Extract JSON from markdown
	jsonRegex := regexp.MustCompile(`(?s)` + "```json\n(.*?)\n```")
	matches := jsonRegex.FindSubmatch(content)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no JSON found in metadata file")
	}

	var metadata ProjectMetadata
	err = json.Unmarshal(matches[1], &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &metadata, nil
}

// SaveMetadata saves metadata to gophex.md file
func SaveMetadata(projectPath string, metadata *ProjectMetadata) error {
	mg := NewMetadataGenerator(projectPath, metadata.Project.Type)
	return mg.WriteMetadataFile(metadata)
}

// IsActivityCompleted checks if an activity has been completed
func IsActivityCompleted(projectPath, activityName string) bool {
	metadata, err := LoadMetadata(projectPath)
	if err != nil {
		return false
	}

	if activity, exists := metadata.Activities[activityName]; exists {
		return activity.Completed
	}
	return false
}

// GetActivityPrefix returns "re-" prefix if activity was already completed
func GetActivityPrefix(projectPath, activityName string) string {
	if IsActivityCompleted(projectPath, activityName) {
		return "re-"
	}
	return ""
}
