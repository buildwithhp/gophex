package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/buildwithhp/gophex/internal/app"
	"github.com/buildwithhp/gophex/internal/infrastructure/generator"
	"github.com/buildwithhp/gophex/internal/infrastructure/repository"
	"github.com/buildwithhp/gophex/internal/shared/config"
	"github.com/buildwithhp/gophex/internal/shared/logger"
	"github.com/buildwithhp/gophex/internal/shared/template"
	"github.com/buildwithhp/gophex/internal/ui/cli"
	"github.com/buildwithhp/gophex/pkg/version"
)

func main() {
	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go handleShutdown(cancel)

	// Run the application
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create logger
	log := createLogger(cfg)
	log.Info("Starting Gophex", "version", version.Version)

	// Build application with dependencies
	application, err := buildApplication(cfg, log)
	if err != nil {
		return fmt.Errorf("failed to build application: %w", err)
	}

	// Initialize application
	if err := application.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}

	// Create and run CLI
	cliApp := cli.NewCLI(application)
	if err := cliApp.Execute(ctx); err != nil {
		return fmt.Errorf("CLI execution failed: %w", err)
	}

	// Graceful shutdown
	if err := application.Shutdown(ctx); err != nil {
		log.Error("Error during shutdown", err)
		return err
	}

	log.Info("Gophex shutdown complete")
	return nil
}

func loadConfiguration() (*config.Config, error) {
	// Create configuration manager with multiple providers
	manager := config.NewManager(
		config.NewDefaultProvider(getDefaultConfig()),
		config.NewEnvironmentProvider(),
		config.NewFileProvider(".gophex.config"),
	)

	if err := manager.Load(); err != nil {
		return nil, err
	}

	return manager.GetConfig(), nil
}

func getDefaultConfig() map[string]string {
	return map[string]string{
		"APP_NAME":                 "gophex",
		"VERSION":                  version.Version,
		"LOG_LEVEL":                "info",
		"DEBUG":                    "false",
		"TEMPLATE_DIR":             "internal/templates",
		"OUTPUT_DIR":               ".",
		"DEFAULT_PROJECT_TYPE":     "api",
		"DEFAULT_MODULE_NAME":      "github.com/user/project",
		"ENABLE_CRUD_GENERATION":   "true",
		"ENABLE_INTERACTIVE_MODE":  "true",
		"ENABLE_METADATA_TRACKING": "true",
	}
}

func createLogger(cfg *config.Config) logger.Logger {
	var level logger.Level
	switch cfg.LogLevel {
	case "debug":
		level = logger.LevelDebug
	case "info":
		level = logger.LevelInfo
	case "warn":
		level = logger.LevelWarn
	case "error":
		level = logger.LevelError
	default:
		level = logger.LevelInfo
	}

	return logger.NewWithLevel(level)
}

func buildApplication(cfg *config.Config, log logger.Logger) (*app.Application, error) {
	// Create template engine
	templateEngine := template.NewEngine()

	// Create repositories
	projectRepo := repository.NewFileRepository(cfg.OutputDir)
	metadataRepo := repository.NewMetadataRepository()

	// Create generator
	projectGenerator := generator.NewGeneratorAdapter()

	// Build application using builder pattern
	application, err := app.NewBuilder().
		WithConfig(cfg).
		WithLogger(log).
		WithTemplateEngine(templateEngine).
		WithProjectRepository(projectRepo).
		WithMetadataRepository(metadataRepo).
		WithGenerator(projectGenerator).
		Build()

	if err != nil {
		return nil, fmt.Errorf("failed to build application: %w", err)
	}

	return application, nil
}

func handleShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Printf("\nReceived signal: %v\n", sig)
	fmt.Println("Shutting down gracefully...")

	cancel()
}
