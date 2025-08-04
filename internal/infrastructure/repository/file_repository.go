package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildwithhp/gophex/internal/domain/project"
)

// fileRepository implements the domain Repository interface using file system
type fileRepository struct {
	basePath string
}

// NewFileRepository creates a new file-based repository
func NewFileRepository(basePath string) project.Repository {
	return &fileRepository{
		basePath: basePath,
	}
}

// Save saves a project to the file system
func (r *fileRepository) Save(ctx context.Context, proj *project.Project) error {
	projectDir := filepath.Join(r.basePath, proj.Name)

	// Create project directory if it doesn't exist
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Save project metadata as JSON
	metadataPath := filepath.Join(projectDir, "project.json")
	data, err := json.MarshalIndent(proj, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project data: %w", err)
	}

	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write project metadata: %w", err)
	}

	return nil
}

// FindByName finds a project by its name
func (r *fileRepository) FindByName(ctx context.Context, name string) (*project.Project, error) {
	projectDir := filepath.Join(r.basePath, name)
	metadataPath := filepath.Join(projectDir, "project.json")

	// Check if project exists
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return nil, project.ErrProjectNotFound
	}

	// Read project metadata
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read project metadata: %w", err)
	}

	var proj project.Project
	if err := json.Unmarshal(data, &proj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project data: %w", err)
	}

	return &proj, nil
}

// FindByPath finds a project by its path
func (r *fileRepository) FindByPath(ctx context.Context, path string) (*project.Project, error) {
	// For file repository, we'll look for a project.json file in the given path
	metadataPath := filepath.Join(path, "project.json")

	// Check if project exists
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return nil, project.ErrProjectNotFound
	}

	// Read project metadata
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read project metadata: %w", err)
	}

	var proj project.Project
	if err := json.Unmarshal(data, &proj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project data: %w", err)
	}

	return &proj, nil
}

// Exists checks if a project exists at the given path
func (r *fileRepository) Exists(ctx context.Context, path string) (bool, error) {
	// Check if the path exists and is a directory
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check path existence: %w", err)
	}

	return info.IsDir(), nil
}

// Delete removes a project from the file system
func (r *fileRepository) Delete(ctx context.Context, name string) error {
	projectDir := filepath.Join(r.basePath, name)

	if err := os.RemoveAll(projectDir); err != nil {
		return fmt.Errorf("failed to delete project directory: %w", err)
	}

	return nil
}

// List returns all projects in the repository
func (r *fileRepository) List(ctx context.Context) ([]*project.Project, error) {
	entries, err := os.ReadDir(r.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*project.Project{}, nil
		}
		return nil, fmt.Errorf("failed to read repository directory: %w", err)
	}

	var projects []*project.Project
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		proj, err := r.FindByName(ctx, entry.Name())
		if err != nil {
			// Skip projects that can't be loaded
			continue
		}

		projects = append(projects, proj)
	}

	return projects, nil
}

// UpdateActivity updates a specific activity for a project
func (r *fileRepository) UpdateActivity(ctx context.Context, projectPath, activityName string, completed bool) error {
	// Find the project by path
	proj, err := r.FindByPath(ctx, projectPath)
	if err != nil {
		return err
	}

	// Update the activity
	if completed {
		proj.CompleteActivity(activityName)
	} else {
		// Mark as incomplete
		if proj.Activities == nil {
			proj.Activities = make(map[string]project.Activity)
		}
		activity := proj.Activities[activityName]
		activity.Completed = false
		activity.Timestamp = nil
		proj.Activities[activityName] = activity
	}

	// Save the updated project
	return r.Save(ctx, proj)
}

// UpdateMetadata updates project metadata
func (r *fileRepository) UpdateMetadata(ctx context.Context, proj *project.Project) error {
	return r.Save(ctx, proj)
}
