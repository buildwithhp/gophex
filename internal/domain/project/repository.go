package project

import (
	"context"
)

// Repository defines the interface for project data operations
type Repository interface {
	// Save saves a project to storage
	Save(ctx context.Context, project *Project) error

	// FindByPath finds a project by its path
	FindByPath(ctx context.Context, path string) (*Project, error)

	// FindByName finds a project by its name
	FindByName(ctx context.Context, name string) (*Project, error)

	// List returns all projects
	List(ctx context.Context) ([]*Project, error)

	// Delete removes a project from storage
	Delete(ctx context.Context, path string) error

	// Exists checks if a project exists at the given path
	Exists(ctx context.Context, path string) (bool, error)

	// UpdateActivity updates a specific activity for a project
	UpdateActivity(ctx context.Context, projectPath, activityName string, completed bool) error

	// UpdateMetadata updates project metadata
	UpdateMetadata(ctx context.Context, project *Project) error
}

// MetadataRepository defines the interface for metadata operations
type MetadataRepository interface {
	// LoadMetadata loads project metadata from storage
	LoadMetadata(ctx context.Context, projectPath string) (*ProjectMetadata, error)

	// SaveMetadata saves project metadata to storage
	SaveMetadata(ctx context.Context, projectPath string, metadata *ProjectMetadata) error

	// HasMetadata checks if metadata exists for a project
	HasMetadata(ctx context.Context, projectPath string) (bool, error)
}

// ProjectMetadata represents the complete metadata structure
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

// ProjectInfo represents basic project information
type ProjectInfo struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Version       string `json:"version"`
	GophexVersion string `json:"gophex_version"`
	GeneratedAt   string `json:"generated_at"`
	LastUpdated   string `json:"last_updated"`
}

// DatabaseMetadata represents database metadata
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

// RedisMetadata represents Redis metadata
type RedisMetadata struct {
	Configured bool `json:"configured"`
	Enabled    bool `json:"enabled"`
}

// ActivityInfo represents activity information
type ActivityInfo struct {
	Completed bool   `json:"completed"`
	Timestamp string `json:"timestamp,omitempty"`
	CanRepeat bool   `json:"can_repeat"`
}

// EndpointInfo represents API endpoint information
type EndpointInfo struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Protected   bool   `json:"protected"`
}

// CommandInfo represents CLI command information
type CommandInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Subcommands []string `json:"subcommands"`
}
