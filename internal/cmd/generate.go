package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/buildwithhp/gophex/internal/generator"
)

func GenerateProject() error {
	var projectType string
	var projectName string

	// Ask for project type
	projectTypePrompt := &survey.Select{
		Message: "What type of Go project would you like to generate?",
		Options: []string{
			"api - REST API with clean architecture",
			"webapp - Web application with templates",
			"microservice - Microservice with gRPC support",
			"cli - Command-line tool",
		},
	}

	err := survey.AskOne(projectTypePrompt, &projectType)
	if err != nil {
		return fmt.Errorf("project type selection failed: %w", err)
	}

	// Extract the actual type from the selection (before the " - " description)
	switch {
	case projectType[:3] == "api":
		projectType = "api"
	case projectType[:6] == "webapp":
		projectType = "webapp"
	case projectType[:12] == "microservice":
		projectType = "microservice"
	case projectType[:3] == "cli":
		projectType = "cli"
	}

	// Ask for project name
	projectNamePrompt := &survey.Input{
		Message: "What is the name of your project?",
		Help:    "This will be used as the directory name and module name",
	}

	err = survey.AskOne(projectNamePrompt, &projectName, survey.WithValidator(survey.Required))
	if err != nil {
		return fmt.Errorf("project name input failed: %w", err)
	}

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	projectPath := filepath.Join(currentDir, projectName)

	// Path confirmation loop
	for {
		var confirm bool
		confirmPrompt := &survey.Confirm{
			Message: fmt.Sprintf("Generate %s project '%s' in %s?", projectType, projectName, projectPath),
			Default: true,
		}

		err = survey.AskOne(confirmPrompt, &confirm)
		if err != nil {
			return fmt.Errorf("confirmation failed: %w", err)
		}

		if confirm {
			break // User confirmed, proceed with generation
		}

		// User said no, ask if they want to change the path or cancel
		var action string
		actionPrompt := &survey.Select{
			Message: "What would you like to do?",
			Options: []string{
				"Change directory path",
				"Cancel project generation",
			},
		}

		err = survey.AskOne(actionPrompt, &action)
		if err != nil {
			return fmt.Errorf("action selection failed: %w", err)
		}

		if action == "Cancel project generation" {
			fmt.Println("Project generation cancelled.")
			return nil
		}

		// Ask for new directory path
		var newPath string
		pathPrompt := &survey.Input{
			Message: "Enter the directory path where you want to create the project:",
			Default: currentDir,
			Help:    "Enter the full path or relative path. The project folder will be created inside this directory.",
		}

		err = survey.AskOne(pathPrompt, &newPath, survey.WithValidator(survey.Required))
		if err != nil {
			return fmt.Errorf("path input failed: %w", err)
		}

		// Update project path
		projectPath = filepath.Join(newPath, projectName)
	}

	// Generate the project
	gen := generator.New()
	if err := gen.Generate(projectType, projectName, projectPath); err != nil {
		return fmt.Errorf("error generating project: %w", err)
	}

	fmt.Printf("✅ Successfully generated %s project '%s' in %s\n", projectType, projectName, projectPath)
	return nil
}
