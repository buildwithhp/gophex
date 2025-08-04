package project

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// Service defines the business logic interface for project operations
type Service interface {
	// CreateProject creates a new project
	CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error)

	// LoadProject loads an existing project
	LoadProject(ctx context.Context, projectPath string) (*Project, error)

	// ListProjects lists all projects
	ListProjects(ctx context.Context) ([]*Project, error)

	// ValidateProjectName validates a project name
	ValidateProjectName(name string) error

	// ValidateProjectPath validates a project path
	ValidateProjectPath(path string) error

	// CompleteActivity marks an activity as completed
	CompleteActivity(ctx context.Context, projectPath, activityName string) error

	// GetActivityStatus gets the status of an activity
	GetActivityStatus(ctx context.Context, projectPath, activityName string) (bool, error)

	// UpdateProjectMetadata updates project metadata
	UpdateProjectMetadata(ctx context.Context, project *Project) error
}

// CreateProjectRequest represents a request to create a project
type CreateProjectRequest struct {
	Name           string
	Type           ProjectType
	Path           string
	DatabaseConfig *DatabaseConfig
	RedisConfig    *RedisConfig
}

// Validate validates the create project request
func (r CreateProjectRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return NewValidationError("name", r.Name, "project name cannot be empty")
	}

	if !r.Type.IsValid() {
		return NewValidationError("type", r.Type, "invalid project type")
	}

	if strings.TrimSpace(r.Path) == "" {
		return NewValidationError("path", r.Path, "project path cannot be empty")
	}

	return nil
}

// service implements the Service interface
type service struct {
	repo         Repository
	metadataRepo MetadataRepository
	generator    Generator
	logger       Logger
}

// Generator defines the interface for project generation
type Generator interface {
	Generate(ctx context.Context, project *Project) error
}

// Logger defines the interface for logging
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

// NewService creates a new project service
func NewService(repo Repository, metadataRepo MetadataRepository, generator Generator, logger Logger) Service {
	return &service{
		repo:         repo,
		metadataRepo: metadataRepo,
		generator:    generator,
		logger:       logger,
	}
}

// CreateProject creates a new project
func (s *service) CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error) {
	s.logger.Info("Creating new project", "name", req.Name, "type", req.Type, "path", req.Path)

	// Validate request
	if err := req.Validate(); err != nil {
		s.logger.Error("Invalid create project request", err, "request", req)
		return nil, err
	}

	// Check if project already exists
	exists, err := s.repo.Exists(ctx, req.Path)
	if err != nil {
		s.logger.Error("Failed to check if project exists", err, "path", req.Path)
		return nil, fmt.Errorf("failed to check project existence: %w", err)
	}

	if exists {
		s.logger.Warn("Project already exists", "path", req.Path)
		return nil, ErrProjectExists
	}

	// Create project entity
	now := time.Now()
	project := &Project{
		Name:           req.Name,
		Type:           req.Type,
		Path:           req.Path,
		ModuleName:     s.generateModuleName(req.Name),
		GeneratedAt:    now,
		LastUpdated:    now,
		DatabaseConfig: req.DatabaseConfig,
		RedisConfig:    req.RedisConfig,
		Activities:     s.createDefaultActivities(req.Type),
		Features:       s.createDefaultFeatures(req.Type),
	}

	// Generate project files
	if err := s.generator.Generate(ctx, project); err != nil {
		s.logger.Error("Failed to generate project", err, "project", project.Name)
		return nil, NewGenerationError("file_generation", project.Name, err)
	}

	// Save project to repository
	if err := s.repo.Save(ctx, project); err != nil {
		s.logger.Error("Failed to save project", err, "project", project.Name)
		return nil, fmt.Errorf("failed to save project: %w", err)
	}

	s.logger.Info("Project created successfully", "name", project.Name, "path", project.Path)
	return project, nil
}

// LoadProject loads an existing project
func (s *service) LoadProject(ctx context.Context, projectPath string) (*Project, error) {
	s.logger.Debug("Loading project", "path", projectPath)

	if strings.TrimSpace(projectPath) == "" {
		return nil, NewValidationError("path", projectPath, "project path cannot be empty")
	}

	project, err := s.repo.FindByPath(ctx, projectPath)
	if err != nil {
		s.logger.Error("Failed to load project", err, "path", projectPath)
		return nil, fmt.Errorf("failed to load project: %w", err)
	}

	s.logger.Debug("Project loaded successfully", "name", project.Name, "path", project.Path)
	return project, nil
}

// ListProjects lists all projects
func (s *service) ListProjects(ctx context.Context) ([]*Project, error) {
	s.logger.Debug("Listing all projects")

	projects, err := s.repo.List(ctx)
	if err != nil {
		s.logger.Error("Failed to list projects", err)
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	s.logger.Debug("Projects listed successfully", "count", len(projects))
	return projects, nil
}

// ValidateProjectName validates a project name
func (s *service) ValidateProjectName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return NewValidationError("name", name, "project name cannot be empty")
	}

	if len(name) < 2 {
		return NewValidationError("name", name, "project name must be at least 2 characters long")
	}

	if len(name) > 50 {
		return NewValidationError("name", name, "project name cannot exceed 50 characters")
	}

	// Check for invalid characters
	for _, char := range name {
		if !isValidNameChar(char) {
			return NewValidationError("name", name, "project name contains invalid characters")
		}
	}

	return nil
}

// ValidateProjectPath validates a project path
func (s *service) ValidateProjectPath(path string) error {
	path = strings.TrimSpace(path)

	if path == "" {
		return NewValidationError("path", path, "project path cannot be empty")
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return NewValidationError("path", path, "invalid path format")
	}

	// Check if path is valid
	if !filepath.IsAbs(absPath) {
		return NewValidationError("path", path, "path must be absolute")
	}

	return nil
}

// CompleteActivity marks an activity as completed
func (s *service) CompleteActivity(ctx context.Context, projectPath, activityName string) error {
	s.logger.Debug("Completing activity", "path", projectPath, "activity", activityName)

	if err := s.repo.UpdateActivity(ctx, projectPath, activityName, true); err != nil {
		s.logger.Error("Failed to complete activity", err, "path", projectPath, "activity", activityName)
		return fmt.Errorf("failed to complete activity: %w", err)
	}

	s.logger.Info("Activity completed", "path", projectPath, "activity", activityName)
	return nil
}

// GetActivityStatus gets the status of an activity
func (s *service) GetActivityStatus(ctx context.Context, projectPath, activityName string) (bool, error) {
	project, err := s.repo.FindByPath(ctx, projectPath)
	if err != nil {
		return false, fmt.Errorf("failed to load project: %w", err)
	}

	return project.IsActivityCompleted(activityName), nil
}

// UpdateProjectMetadata updates project metadata
func (s *service) UpdateProjectMetadata(ctx context.Context, project *Project) error {
	s.logger.Debug("Updating project metadata", "name", project.Name)

	if err := s.repo.UpdateMetadata(ctx, project); err != nil {
		s.logger.Error("Failed to update project metadata", err, "name", project.Name)
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	s.logger.Debug("Project metadata updated successfully", "name", project.Name)
	return nil
}

// Helper methods

// generateModuleName generates a Go module name from project name
func (s *service) generateModuleName(projectName string) string {
	// Simple implementation - in real world, this might be more sophisticated
	name := strings.ToLower(projectName)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	return fmt.Sprintf("github.com/user/%s", name)
}

// createDefaultActivities creates default activities for a project type
func (s *service) createDefaultActivities(projectType ProjectType) map[string]Activity {
	activities := map[string]Activity{
		"project_generated": {
			Name:      "project_generated",
			Completed: true,
			Timestamp: timePtr(time.Now()),
			CanRepeat: false,
		},
		"dependencies_installed": {
			Name:      "dependencies_installed",
			Completed: false,
			CanRepeat: true,
		},
		"tests_executed": {
			Name:      "tests_executed",
			Completed: false,
			CanRepeat: true,
		},
		"documentation_viewed": {
			Name:      "documentation_viewed",
			Completed: false,
			CanRepeat: true,
		},
	}

	// Add project-type specific activities
	switch projectType {
	case ProjectTypeAPI:
		activities["database_migrated"] = Activity{
			Name:      "database_migrated",
			Completed: false,
			CanRepeat: true,
		}
		activities["application_started"] = Activity{
			Name:      "application_started",
			Completed: false,
			CanRepeat: true,
		}
		activities["crud_generated"] = Activity{
			Name:      "crud_generated",
			Completed: false,
			CanRepeat: true,
		}
	case ProjectTypeWebApp, ProjectTypeMicroservice:
		activities["application_started"] = Activity{
			Name:      "application_started",
			Completed: false,
			CanRepeat: true,
		}
	case ProjectTypeCLI:
		activities["application_built"] = Activity{
			Name:      "application_built",
			Completed: false,
			CanRepeat: true,
		}
	}

	return activities
}

// createDefaultFeatures creates default features for a project type
func (s *service) createDefaultFeatures(projectType ProjectType) []Feature {
	var features []Feature

	switch projectType {
	case ProjectTypeAPI:
		features = []Feature{
			{Name: "authentication", Enabled: true, Description: "JWT-based authentication"},
			{Name: "user_management", Enabled: true, Description: "User CRUD operations"},
			{Name: "health_checks", Enabled: true, Description: "Health check endpoints"},
			{Name: "cors_enabled", Enabled: true, Description: "CORS middleware"},
			{Name: "rate_limiting", Enabled: true, Description: "Rate limiting middleware"},
			{Name: "request_logging", Enabled: true, Description: "Request logging middleware"},
			{Name: "clean_architecture", Enabled: true, Description: "Clean architecture pattern"},
		}
	case ProjectTypeWebApp:
		features = []Feature{
			{Name: "web_server", Enabled: true, Description: "HTTP web server"},
			{Name: "static_files", Enabled: true, Description: "Static file serving"},
			{Name: "html_templates", Enabled: true, Description: "HTML template rendering"},
		}
	case ProjectTypeMicroservice:
		features = []Feature{
			{Name: "grpc_support", Enabled: true, Description: "gRPC server support"},
			{Name: "health_checks", Enabled: true, Description: "Health check endpoints"},
			{Name: "lightweight_design", Enabled: true, Description: "Lightweight microservice design"},
		}
	case ProjectTypeCLI:
		features = []Feature{
			{Name: "cobra_framework", Enabled: true, Description: "Cobra CLI framework"},
			{Name: "subcommands", Enabled: true, Description: "Subcommand support"},
		}
	}

	return features
}

// isValidNameChar checks if a character is valid for project names
func isValidNameChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == ' '
}

// timePtr returns a pointer to a time value
func timePtr(t time.Time) *time.Time {
	return &t
}
