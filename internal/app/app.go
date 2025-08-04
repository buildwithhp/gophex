package app

import (
	"context"
	"fmt"

	"github.com/buildwithhp/gophex/internal/domain/project"
	"github.com/buildwithhp/gophex/internal/shared/config"
	"github.com/buildwithhp/gophex/internal/shared/logger"
	"github.com/buildwithhp/gophex/internal/shared/template"
)

// Application represents the main application
type Application struct {
	config         *config.Config
	logger         logger.Logger
	templateEngine template.Engine
	projectService project.Service
	projectRepo    project.Repository
	metadataRepo   project.MetadataRepository
	generator      project.Generator
}

// Dependencies holds all application dependencies
type Dependencies struct {
	Config         *config.Config
	Logger         logger.Logger
	TemplateEngine template.Engine
	ProjectRepo    project.Repository
	MetadataRepo   project.MetadataRepository
	Generator      project.Generator
}

// New creates a new application instance
func New(deps Dependencies) *Application {
	// Create project service with dependencies
	projectService := project.NewService(
		deps.ProjectRepo,
		deps.MetadataRepo,
		deps.Generator,
		deps.Logger,
	)

	return &Application{
		config:         deps.Config,
		logger:         deps.Logger,
		templateEngine: deps.TemplateEngine,
		projectService: projectService,
		projectRepo:    deps.ProjectRepo,
		metadataRepo:   deps.MetadataRepo,
		generator:      deps.Generator,
	}
}

// Initialize initializes the application
func (a *Application) Initialize(ctx context.Context) error {
	a.logger.Info("Initializing Gophex application", "version", a.config.Version)

	// Initialize template engine
	if err := a.initializeTemplateEngine(); err != nil {
		return fmt.Errorf("failed to initialize template engine: %w", err)
	}

	a.logger.Info("Application initialized successfully")
	return nil
}

// GetProjectService returns the project service
func (a *Application) GetProjectService() project.Service {
	return a.projectService
}

// GetConfig returns the application configuration
func (a *Application) GetConfig() *config.Config {
	return a.config
}

// GetLogger returns the application logger
func (a *Application) GetLogger() logger.Logger {
	return a.logger
}

// GetTemplateEngine returns the template engine
func (a *Application) GetTemplateEngine() template.Engine {
	return a.templateEngine
}

// Shutdown gracefully shuts down the application
func (a *Application) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down Gophex application")

	// Perform cleanup operations here
	// For example, close database connections, save state, etc.

	a.logger.Info("Application shutdown complete")
	return nil
}

// initializeTemplateEngine initializes the template engine
func (a *Application) initializeTemplateEngine() error {
	// Load templates from the configured template directory
	// This would be implemented based on your template loading strategy
	a.logger.Debug("Template engine initialized", "templateDir", a.config.TemplateDir)
	return nil
}

// Builder provides a fluent interface for building the application
type Builder struct {
	config         *config.Config
	logger         logger.Logger
	templateEngine template.Engine
	projectRepo    project.Repository
	metadataRepo   project.MetadataRepository
	generator      project.Generator
}

// NewBuilder creates a new application builder
func NewBuilder() *Builder {
	return &Builder{}
}

// WithConfig sets the configuration
func (b *Builder) WithConfig(cfg *config.Config) *Builder {
	b.config = cfg
	return b
}

// WithLogger sets the logger
func (b *Builder) WithLogger(log logger.Logger) *Builder {
	b.logger = log
	return b
}

// WithTemplateEngine sets the template engine
func (b *Builder) WithTemplateEngine(engine template.Engine) *Builder {
	b.templateEngine = engine
	return b
}

// WithProjectRepository sets the project repository
func (b *Builder) WithProjectRepository(repo project.Repository) *Builder {
	b.projectRepo = repo
	return b
}

// WithMetadataRepository sets the metadata repository
func (b *Builder) WithMetadataRepository(repo project.MetadataRepository) *Builder {
	b.metadataRepo = repo
	return b
}

// WithGenerator sets the project generator
func (b *Builder) WithGenerator(gen project.Generator) *Builder {
	b.generator = gen
	return b
}

// Build builds the application with all dependencies
func (b *Builder) Build() (*Application, error) {
	// Validate required dependencies
	if b.config == nil {
		return nil, fmt.Errorf("configuration is required")
	}

	if b.logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if b.templateEngine == nil {
		return nil, fmt.Errorf("template engine is required")
	}

	if b.projectRepo == nil {
		return nil, fmt.Errorf("project repository is required")
	}

	if b.metadataRepo == nil {
		return nil, fmt.Errorf("metadata repository is required")
	}

	if b.generator == nil {
		return nil, fmt.Errorf("generator is required")
	}

	deps := Dependencies{
		Config:         b.config,
		Logger:         b.logger,
		TemplateEngine: b.templateEngine,
		ProjectRepo:    b.projectRepo,
		MetadataRepo:   b.metadataRepo,
		Generator:      b.generator,
	}

	return New(deps), nil
}

// DefaultBuilder creates a builder with default implementations
func DefaultBuilder() *Builder {
	return NewBuilder().
		WithLogger(logger.New()).
		WithTemplateEngine(template.NewEngine())
}
