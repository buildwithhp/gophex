package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/buildwithhp/gophex/internal/utils"
)

// DiscoveredProject represents information about a discovered Gophex project
type DiscoveredProject struct {
	Name         string
	Type         string
	Path         string
	LastUpdated  string
	RelativePath string
}

// LoadExistingProject handles loading an existing Gophex project
func LoadExistingProject() error {
	fmt.Println("📁 Load Existing Gophex Project")
	fmt.Println("💡 Enter the path to a directory containing a 'gophex.md' file")
	fmt.Println()

	return browseForProject()
}

// browseForProject allows manual browsing for a project
func browseForProject() error {
	var projectPath string
	pathPrompt := &survey.Input{
		Message: "Enter the path to a Gophex project directory:",
		Help:    "The directory should contain a 'gophex.md' file",
	}

	err := survey.AskOne(pathPrompt, &projectPath)
	if err != nil {
		if isUserInterrupt(err) {
			return nil
		}
		return fmt.Errorf("path input failed: %w", err)
	}

	// Handle empty input
	if strings.TrimSpace(projectPath) == "" {
		fmt.Println("❌ No path provided")
		return Execute() // Return to main menu
	}

	// Expand ~ to home directory
	if strings.HasPrefix(projectPath, "~/") {
		home := os.Getenv("HOME")
		projectPath = filepath.Join(home, projectPath[2:])
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		fmt.Printf("❌ Invalid path: %v\n", err)
		return Execute()
	}

	// Check if directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("❌ Directory does not exist: %s\n", absPath)
		return Execute()
	}

	// Check if it's a Gophex project
	if !utils.HasGophexMetadata(absPath) {
		fmt.Printf("❌ No Gophex metadata file (gophex.md) found in: %s\n", absPath)
		fmt.Println()

		// Ask if user wants to create a new project
		var createNew bool
		createPrompt := &survey.Confirm{
			Message: "Would you like to create a new project?",
			Default: true,
		}

		err := survey.AskOne(createPrompt, &createNew)
		if err != nil {
			if isUserInterrupt(err) {
				return nil
			}
			return fmt.Errorf("create project prompt failed: %w", err)
		}

		if createNew {
			return GenerateProject()
		}

		return Execute() // Return to main menu
	}

	// Load metadata
	fmt.Println(absPath)
	metadata, err := utils.LoadMetadata(absPath)
	if err != nil {
		fmt.Printf("❌ Failed to load project metadata: %v\n", err)
		fmt.Println("💡 The gophex.md file may be corrupted or invalid.")
		return Execute()
	}

	// Create relative path for display
	cwd, _ := os.Getwd()
	relativePath, err := filepath.Rel(cwd, absPath)
	if err != nil {
		relativePath = absPath
	}

	project := &DiscoveredProject{
		Name:         metadata.Project.Name,
		Type:         metadata.Project.Type,
		Path:         absPath,
		LastUpdated:  metadata.Project.LastUpdated,
		RelativePath: relativePath,
	}

	return loadProject(project)
}

// loadProject loads a selected project and shows the post-generation menu
func loadProject(project *DiscoveredProject) error {
	fmt.Printf("📂 Loading project: %s (%s)\n", project.Name, project.Type)
	fmt.Printf("📍 Location: %s\n", project.RelativePath)

	timeAgo := formatTimeAgo(project.LastUpdated)
	fmt.Printf("🕒 Last updated: %s\n\n", timeAgo)

	// Create project options for post-generation workflow
	opts := PostGenerationOptions{
		ProjectName: project.Name,
		ProjectType: project.Type,
		ProjectPath: project.Path,
	}

	// Show post-generation menu
	return ShowPostGenerationMenu(opts)
}

// formatTimeAgo formats a timestamp into a human-readable "time ago" string
func formatTimeAgo(timestamp string) string {
	if timestamp == "" {
		return "unknown"
	}

	// Parse the timestamp
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return "unknown"
	}

	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / (24 * 7))
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := int(diff.Hours() / (24 * 365))
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}
