package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/buildwithhp/gophex/internal/generator"
	"github.com/buildwithhp/gophex/pkg/version"
)

// ProjectMetadata represents the complete project tracking information
type ProjectMetadata struct {
	Gophex GophexMetadata `json:"gophex"`
}

// GophexMetadata contains all Gophex-specific tracking data
type GophexMetadata struct {
	Version     string           `json:"version"`
	GeneratedAt string           `json:"generated_at"`
	Project     ProjectInfo      `json:"project"`
	Database    DatabaseInfo     `json:"database"`
	Redis       RedisInfo        `json:"redis"`
	Activities  ActivityTracker  `json:"activities"`
	Hierarchy   ProjectHierarchy `json:"hierarchy"`
}

// ProjectInfo contains basic project information
type ProjectInfo struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Path       string `json:"path"`
	ModuleName string `json:"module_name"`
}

// DatabaseInfo tracks database configuration (no sensitive data)
type DatabaseInfo struct {
	Configured              bool   `json:"configured"`
	Type                    string `json:"type,omitempty"`        // postgresql, mysql, mongodb
	ConfigType              string `json:"config_type,omitempty"` // single, read-write, cluster
	MigrationsExecuted      bool   `json:"migrations_executed"`
	InitializationCompleted bool   `json:"initialization_completed"`
	SupportsMigrations      bool   `json:"supports_migrations"`
}

// RedisInfo tracks Redis configuration (no sensitive data)
type RedisInfo struct {
	Configured bool `json:"configured"`
	Enabled    bool `json:"enabled"`
}

// ActivityTracker tracks completion status of various activities
type ActivityTracker struct {
	DependenciesInstalled bool `json:"dependencies_installed"`
	DatabaseSetup         bool `json:"database_setup"`
	ApplicationStarted    bool `json:"application_started"`
	TestsRun              bool `json:"tests_run"`
	HealthCheckTested     bool `json:"health_check_tested"`
	DocumentationViewed   bool `json:"documentation_viewed"`
	ChangeDetectionRun    bool `json:"change_detection_run"`
}

// ProjectHierarchy represents the file structure of the generated project
type ProjectHierarchy struct {
	Cmd        map[string]interface{} `json:"cmd,omitempty"`
	Internal   map[string]interface{} `json:"internal,omitempty"`
	Migrations []string               `json:"migrations,omitempty"`
	Scripts    []string               `json:"scripts,omitempty"`
	Web        map[string]interface{} `json:"web,omitempty"`
	RootFiles  []string               `json:"root_files"`
}

// ProjectTracker manages project metadata operations
type ProjectTracker struct {
	projectPath string
	metadata    *ProjectMetadata
}

// NewProjectTracker creates a new project tracker for the given project path
func NewProjectTracker(projectPath string) *ProjectTracker {
	return &ProjectTracker{
		projectPath: projectPath,
	}
}

// LoadMetadata loads existing project metadata from gophex.md
func (pt *ProjectTracker) LoadMetadata() error {
	metadataPath := filepath.Join(pt.projectPath, "gophex.md")

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, create default metadata
			pt.metadata = pt.createDefaultMetadata()
			return nil
		}
		return fmt.Errorf("failed to read gophex.md: %w", err)
	}

	pt.metadata = &ProjectMetadata{}
	if err := json.Unmarshal(data, pt.metadata); err != nil {
		return fmt.Errorf("failed to parse gophex.md: %w", err)
	}

	return nil
}

// SaveMetadata saves the current metadata to gophex.md
func (pt *ProjectTracker) SaveMetadata() error {
	metadataPath := filepath.Join(pt.projectPath, "gophex.md")

	data, err := json.MarshalIndent(pt.metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write gophex.md: %w", err)
	}

	return nil
}

// CreateInitialMetadata creates and saves initial project metadata
func (pt *ProjectTracker) CreateInitialMetadata(projectType, projectName, projectPath string, dbConfig *generator.DatabaseConfig, redisConfig *generator.RedisConfig) error {
	pt.metadata = &ProjectMetadata{
		Gophex: GophexMetadata{
			Version:     version.GetVersion(),
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
			Project: ProjectInfo{
				Name:       projectName,
				Type:       projectType,
				Path:       projectPath,
				ModuleName: generateModuleName(projectName),
			},
			Database:   pt.createDatabaseInfo(dbConfig),
			Redis:      pt.createRedisInfo(redisConfig),
			Activities: ActivityTracker{}, // All activities start as false
			Hierarchy:  pt.generateHierarchy(projectType, dbConfig),
		},
	}

	return pt.SaveMetadata()
}

// GetMetadata returns the current metadata
func (pt *ProjectTracker) GetMetadata() *ProjectMetadata {
	if pt.metadata == nil {
		pt.metadata = pt.createDefaultMetadata()
	}
	return pt.metadata
}

// UpdateActivity updates the status of a specific activity
func (pt *ProjectTracker) UpdateActivity(activity string, completed bool) error {
	if pt.metadata == nil {
		if err := pt.LoadMetadata(); err != nil {
			return err
		}
	}

	switch activity {
	case "dependencies_installed":
		pt.metadata.Gophex.Activities.DependenciesInstalled = completed
	case "database_setup":
		pt.metadata.Gophex.Activities.DatabaseSetup = completed
	case "application_started":
		pt.metadata.Gophex.Activities.ApplicationStarted = completed
	case "tests_run":
		pt.metadata.Gophex.Activities.TestsRun = completed
	case "health_check_tested":
		pt.metadata.Gophex.Activities.HealthCheckTested = completed
	case "documentation_viewed":
		pt.metadata.Gophex.Activities.DocumentationViewed = completed
	case "change_detection_run":
		pt.metadata.Gophex.Activities.ChangeDetectionRun = completed
	default:
		return fmt.Errorf("unknown activity: %s", activity)
	}

	return pt.SaveMetadata()
}

// UpdateDatabaseStatus updates database-related status
func (pt *ProjectTracker) UpdateDatabaseStatus(migrationsExecuted, initializationCompleted bool) error {
	if pt.metadata == nil {
		if err := pt.LoadMetadata(); err != nil {
			return err
		}
	}

	pt.metadata.Gophex.Database.MigrationsExecuted = migrationsExecuted
	pt.metadata.Gophex.Database.InitializationCompleted = initializationCompleted

	return pt.SaveMetadata()
}

// IsActivityCompleted checks if a specific activity has been completed
func (pt *ProjectTracker) IsActivityCompleted(activity string) bool {
	if pt.metadata == nil {
		return false
	}

	switch activity {
	case "dependencies_installed":
		return pt.metadata.Gophex.Activities.DependenciesInstalled
	case "database_setup":
		return pt.metadata.Gophex.Activities.DatabaseSetup
	case "application_started":
		return pt.metadata.Gophex.Activities.ApplicationStarted
	case "tests_run":
		return pt.metadata.Gophex.Activities.TestsRun
	case "health_check_tested":
		return pt.metadata.Gophex.Activities.HealthCheckTested
	case "documentation_viewed":
		return pt.metadata.Gophex.Activities.DocumentationViewed
	case "change_detection_run":
		return pt.metadata.Gophex.Activities.ChangeDetectionRun
	default:
		return false
	}
}

// GetActivityPrefix returns "Re-" if activity is completed, empty string otherwise
func (pt *ProjectTracker) GetActivityPrefix(activity string) string {
	if pt.IsActivityCompleted(activity) {
		return "Re-"
	}
	return ""
}

// createDefaultMetadata creates default metadata when none exists
func (pt *ProjectTracker) createDefaultMetadata() *ProjectMetadata {
	return &ProjectMetadata{
		Gophex: GophexMetadata{
			Version:     version.GetVersion(),
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
			Project:     ProjectInfo{},
			Database:    DatabaseInfo{},
			Redis:       RedisInfo{},
			Activities:  ActivityTracker{},
			Hierarchy:   ProjectHierarchy{RootFiles: []string{}},
		},
	}
}

// createDatabaseInfo creates database info from config (no sensitive data)
func (pt *ProjectTracker) createDatabaseInfo(dbConfig *generator.DatabaseConfig) DatabaseInfo {
	if dbConfig == nil {
		return DatabaseInfo{
			Configured: false,
		}
	}

	supportsMigrations := dbConfig.Type == "postgresql" || dbConfig.Type == "mysql"

	return DatabaseInfo{
		Configured:         true,
		Type:               dbConfig.Type,
		ConfigType:         dbConfig.ConfigType,
		SupportsMigrations: supportsMigrations,
	}
}

// createRedisInfo creates Redis info from config (no sensitive data)
func (pt *ProjectTracker) createRedisInfo(redisConfig *generator.RedisConfig) RedisInfo {
	if redisConfig == nil {
		return RedisInfo{
			Configured: false,
		}
	}

	return RedisInfo{
		Configured: true,
		Enabled:    redisConfig.Enabled,
	}
}

// generateHierarchy creates the project hierarchy based on project type and database
func (pt *ProjectTracker) generateHierarchy(projectType string, dbConfig *generator.DatabaseConfig) ProjectHierarchy {
	hierarchy := ProjectHierarchy{
		RootFiles: []string{".env", ".env.example", ".gophex-generated", "go.mod", "README.md"},
	}

	switch projectType {
	case "api":
		hierarchy.Cmd = map[string]interface{}{
			"api/": []string{"main.go"},
		}

		hierarchy.Internal = map[string]interface{}{
			"api/": map[string]interface{}{
				"handlers/":   []string{"auth.go", "health.go", "posts.go", "users.go"},
				"middleware/": []string{"auth.go", "cors.go", "logging.go", "ratelimit.go"},
				"responses/":  []string{"error.go", "success.go"},
				"routes/":     []string{"routes.go"},
			},
			"config/":   []string{"config.go"},
			"database/": []string{"config.go", "database.go", "factory.go"},
			"domain/": map[string]interface{}{
				"post/": []string{"model.go", "repository.go", "service.go"},
				"user/": []string{"model.go", "repository.go", "service.go"},
			},
			"infrastructure/": map[string]interface{}{
				"auth/":     []string{"jwt.go", "password.go"},
				"database/": pt.getDatabaseInfrastructure(dbConfig),
			},
			"pkg/": map[string]interface{}{
				"errors/":    []string{"errors.go"},
				"logger/":    []string{"logger.go"},
				"validator/": []string{"validator.go"},
			},
		}

		if dbConfig != nil && (dbConfig.Type == "postgresql" || dbConfig.Type == "mysql") {
			hierarchy.Migrations = []string{
				"000001_create_users_table.up.sql",
				"000001_create_users_table.down.sql",
				"000002_create_posts_table.up.sql",
				"000002_create_posts_table.down.sql",
				"README.md",
			}
		} else if dbConfig != nil && dbConfig.Type == "mongodb" {
			hierarchy.Migrations = []string{
				"mongodb_init.js",
				"README.md",
			}
		}

		hierarchy.Scripts = []string{"detect-changes.sh", "migrate.sh"}

	case "webapp":
		hierarchy.Cmd = map[string]interface{}{
			"webapp/": []string{"main.go"},
		}
		hierarchy.Web = map[string]interface{}{
			"static/": map[string]interface{}{
				"css/": []string{"style.css"},
			},
			"templates/": []string{"index.html"},
		}

	case "microservice":
		hierarchy.Cmd = map[string]interface{}{
			"server/": []string{"main.go"},
		}
		hierarchy.Internal = map[string]interface{}{
			"handlers/": []string{"handlers.go"},
		}

	case "cli":
		hierarchy.Cmd = map[string]interface{}{
			"main.go": nil,
		}
		hierarchy.Internal = map[string]interface{}{
			"cmd/": []string{"root.go"},
		}
	}

	return hierarchy
}

// getDatabaseInfrastructure returns database-specific infrastructure files
func (pt *ProjectTracker) getDatabaseInfrastructure(dbConfig *generator.DatabaseConfig) map[string]interface{} {
	infrastructure := make(map[string]interface{})

	if dbConfig == nil {
		return infrastructure
	}

	switch dbConfig.Type {
	case "postgresql":
		infrastructure["postgres/"] = []string{"connection.go", "post_repo.go", "user_repo.go"}
	case "mysql":
		infrastructure["mysql/"] = []string{"connection.go", "post_repo.go", "user_repo.go"}
	case "mongodb":
		infrastructure["mongodb/"] = []string{"connection.go", "post_repo.go", "user_repo.go"}
	}

	infrastructure["redis/"] = []string{"client.go"}

	return infrastructure
}

// generateModuleName generates a Go module name from project name
func generateModuleName(projectName string) string {
	// Simple module name generation - can be enhanced later
	return fmt.Sprintf("github.com/user/%s", projectName)
}
