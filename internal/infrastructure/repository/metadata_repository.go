package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildwithhp/gophex/internal/domain/project"
)

// metadataRepository implements the domain MetadataRepository interface using file system
type metadataRepository struct{}

// NewMetadataRepository creates a new file-based metadata repository
func NewMetadataRepository() project.MetadataRepository {
	return &metadataRepository{}
}

// LoadMetadata loads project metadata from storage
func (r *metadataRepository) LoadMetadata(ctx context.Context, projectPath string) (*project.ProjectMetadata, error) {
	metadataPath := filepath.Join(projectPath, "gophex.md")

	// Check if metadata file exists
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return nil, project.ErrProjectNotFound
	}

	// For now, we'll create a basic metadata structure
	// In a full implementation, we would parse the gophex.md file
	metadata := &project.ProjectMetadata{
		Project: project.ProjectInfo{
			Name:    filepath.Base(projectPath),
			Type:    "unknown", // Would be parsed from file
			Version: "1.0.0",
		},
		Hierarchy: make(map[string]interface{}),
		Database: project.DatabaseMetadata{
			Configured: false,
		},
		Redis: project.RedisMetadata{
			Configured: false,
		},
		Activities: make(map[string]project.ActivityInfo),
		Features:   make(map[string]bool),
	}

	return metadata, nil
}

// SaveMetadata saves project metadata to storage
func (r *metadataRepository) SaveMetadata(ctx context.Context, projectPath string, metadata *project.ProjectMetadata) error {
	metadataPath := filepath.Join(projectPath, "project_metadata.json")

	// Marshal metadata to JSON
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Write to file
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// HasMetadata checks if metadata exists for a project
func (r *metadataRepository) HasMetadata(ctx context.Context, projectPath string) (bool, error) {
	metadataPath := filepath.Join(projectPath, "gophex.md")

	_, err := os.Stat(metadataPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check metadata existence: %w", err)
	}

	return true, nil
}
