package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetEnvWithDefault returns the value of an environment variable or a default value if not set
func GetEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// ActivityInfo represents information about a project activity
type ActivityInfo struct {
	Completed bool   `json:"completed"`
	Timestamp string `json:"timestamp,omitempty"`
	CanRepeat bool   `json:"can_repeat"`
}

// ProjectMetadata represents basic project metadata structure
type ProjectMetadata struct {
	Project struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		LastUpdated string `json:"last_updated"`
	} `json:"project"`
	Database struct {
		MigrationsExecuted bool `json:"migrations_executed"`
		SchemaInitialized  bool `json:"schema_initialized"`
	} `json:"database"`
	Activities map[string]ActivityInfo `json:"activities"`
}

// LegacyMetadata represents the old gophex.md format
type LegacyMetadata struct {
	Gophex struct {
		Version     string `json:"version"`
		GeneratedAt string `json:"generated_at"`
		Project     struct {
			Name   string `json:"name"`
			Type   string `json:"type"`
			Path   string `json:"path"`
			Module string `json:"module"`
		} `json:"project"`
	} `json:"gophex"`
}

// UpdateActivity updates the status of a specific activity in the project metadata
func UpdateActivity(projectPath, activityName string, completed bool) error {
	metadata, err := LoadMetadata(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
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

// UpdateDatabaseStatus updates database-related status in the project metadata
func UpdateDatabaseStatus(projectPath string, migrationsExecuted, schemaInitialized bool) error {
	metadata, err := LoadMetadata(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	metadata.Database.MigrationsExecuted = migrationsExecuted
	metadata.Database.SchemaInitialized = schemaInitialized
	metadata.Project.LastUpdated = time.Now().Format(time.RFC3339)

	return SaveMetadata(projectPath, metadata)
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

// LoadMetadata loads metadata from gophex.md file
func LoadMetadata(projectPath string) (*ProjectMetadata, error) {
	filePath := filepath.Join(projectPath, "gophex.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	// Extract JSON from markdown or handle legacy JSON format
	contentStr := string(content)
	var jsonContent string

	// First, try to extract from markdown format
	startMarkers := []string{
		"```json\n",
		"```json\r\n",
		"``` json\n",
		"``` json\r\n",
	}

	endMarkers := []string{
		"\n```",
		"\r\n```",
		"\n```\n",
		"\r\n```\r\n",
	}

	var found bool

	for _, startMarker := range startMarkers {
		startIdx := strings.Index(contentStr, startMarker)
		if startIdx == -1 {
			continue
		}

		jsonStart := startIdx + len(startMarker)

		for _, endMarker := range endMarkers {
			endIdx := strings.Index(contentStr[jsonStart:], endMarker)
			if endIdx != -1 {
				jsonContent = contentStr[jsonStart : jsonStart+endIdx]
				found = true
				break
			}
		}

		if found {
			break
		}
	}

	// If no markdown markers found, try to parse as legacy JSON format
	if !found {
		// Check if it looks like JSON (starts with { and contains "gophex")
		trimmed := strings.TrimSpace(contentStr)
		if strings.HasPrefix(trimmed, "{") && strings.Contains(trimmed, "\"gophex\"") {
			jsonContent = trimmed
			found = true
		}
	}

	if !found {
		preview := contentStr
		if len(preview) > 200 {
			preview = preview[:200]
		}
		return nil, fmt.Errorf("no JSON markers found and content doesn't appear to be legacy JSON format. File content preview: %q", preview)
	}

	// Try to parse as new format first
	var metadata ProjectMetadata
	err = json.Unmarshal([]byte(jsonContent), &metadata)
	if err != nil || metadata.Project.Name == "" {
		// If new format fails or results in empty project name, try legacy format
		var legacyMetadata LegacyMetadata
		err = json.Unmarshal([]byte(jsonContent), &legacyMetadata)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata in both new and legacy formats: %w", err)
		}

		// Validate that we actually got legacy data
		if legacyMetadata.Gophex.Project.Name == "" {
			return nil, fmt.Errorf("metadata file appears to be corrupted - no project name found in either format")
		}

		// Convert legacy format to new format
		metadata = ProjectMetadata{
			Project: struct {
				Name        string `json:"name"`
				Type        string `json:"type"`
				LastUpdated string `json:"last_updated"`
			}{
				Name:        legacyMetadata.Gophex.Project.Name,
				Type:        legacyMetadata.Gophex.Project.Type,
				LastUpdated: legacyMetadata.Gophex.GeneratedAt,
			},
			Database: struct {
				MigrationsExecuted bool `json:"migrations_executed"`
				SchemaInitialized  bool `json:"schema_initialized"`
			}{
				MigrationsExecuted: false,
				SchemaInitialized:  false,
			},
			Activities: make(map[string]ActivityInfo),
		}
	}
	return &metadata, nil
}

// SaveMetadata saves metadata to gophex.md file
func SaveMetadata(projectPath string, metadata *ProjectMetadata) error {
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
	filePath := filepath.Join(projectPath, "gophex.md")
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// HasGophexMetadata checks if a directory contains a gophex.md file
func HasGophexMetadata(projectPath string) bool {
	filePath := filepath.Join(projectPath, "gophex.md")
	_, err := os.Stat(filePath)
	return err == nil
}
