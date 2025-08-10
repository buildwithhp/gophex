package cmd

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/buildwithhp/gophex/internal/utils"
	"github.com/buildwithhp/gophex/pkg/version"
	"github.com/dolmen-go/kittyimg"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

//go:embed assets/*
var files embed.FS

// ErrReturnToMenu is a special error that signals to return to the main menu
var ErrReturnToMenu = errors.New("return to main menu")

// ErrUserQuit is a special error that signals the user wants to quit
var ErrUserQuit = errors.New("user quit")

// isUserInterrupt checks if the error is due to user interruption (Ctrl+C, EOF, etc.)
func isUserInterrupt(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "interrupt") ||
		strings.Contains(errStr, "eof") ||
		strings.Contains(errStr, "cancelled") ||
		strings.Contains(errStr, "canceled")
}

// clearScreen clears the terminal screen for a cleaner user experience
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// askWithInterruptHandling wraps survey.AskOne with graceful interrupt handling
func askWithInterruptHandling(prompt survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	err := survey.AskOne(prompt, response, opts...)
	if err != nil && isUserInterrupt(err) {
		fmt.Println("\nOperation cancelled. Goodbye! 👋")
		os.Exit(0)
	}
	return err
}

func Execute() error {
	// Check if current directory contains gophex.md
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	var currentProject *utils.ProjectMetadata
	var options []string
	var hasCurrentProject bool

	if utils.HasGophexMetadata(cwd) {
		// Try to load current project metadata
		metadata, err := utils.LoadMetadata(cwd)
		if err == nil {
			currentProject = metadata
			hasCurrentProject = true
			fmt.Printf("📂 Current directory contains Gophex project: %s (%s)\n",
				metadata.Project.Name, metadata.Project.Type)
			fmt.Println()

			options = []string{
				"Load current project",
				"Generate a new project",
				"Enhanced CRUD Wizard - Learn Clean Architecture",
				"Load different project",
				"Show version",
				"Show help",
				"Quit",
			}
		} else {
			// gophex.md exists but is corrupted
			fmt.Printf("⚠️  Found gophex.md in current directory but failed to load: %v\n", err)
			fmt.Println()
			options = []string{
				"Generate a new project",
				"Enhanced CRUD Wizard - Learn Clean Architecture",
				"Load existing project",
				"Show version",
				"Show help",
				"Quit",
			}
		}
	} else {
		options = []string{
			"Generate a new project",
			"Enhanced CRUD Wizard - Learn Clean Architecture",
			"Load existing project",
			"Show version",
			"Print image",
			"Show help",
			"Quit",
		}
	}

	for {
		var action string
		prompt := &survey.Select{
			Message: "What would you like to do?",
			Options: options,
		}

		err = survey.AskOne(prompt, &action)
		if err != nil {
			// Handle user interruption (Ctrl+C) gracefully
			if isUserInterrupt(err) {
				fmt.Println("\nGoodbye! 👋")
				return nil
			}
			return fmt.Errorf("prompt failed: %w", err)
		}

		switch action {
		case "Load current project":
			if hasCurrentProject {
				err = loadCurrentProject(cwd, currentProject)
				if err == ErrReturnToMenu {
					continue // Return to main menu
				}
				return err
			}
			return fmt.Errorf("no current project available")
		case "Generate a new project":
			err = GenerateProject()
			if err == ErrReturnToMenu {
				continue // Return to main menu
			}
			return err
		case "Enhanced CRUD Wizard - Learn Clean Architecture":
			// First check if we're in a project directory
			if hasCurrentProject && currentProject.Project.Type == "api" {
				err = RunEnhancedCRUDWizard(cwd)
			} else {
				fmt.Println("🎓 Enhanced CRUD Wizard")
				fmt.Println("This wizard requires an existing API project.")
				fmt.Println("Let's create one first, then run the CRUD wizard.")
				fmt.Println()

				err = GenerateProject()
				if err == nil {
					// After successful project generation, run the enhanced wizard
					fmt.Println("\n🚀 Now let's create your first CRUD operations!")
					err = RunEnhancedCRUDWizard(cwd)
				}
			}
			if err == ErrReturnToMenu {
				continue // Return to main menu
			}
			return err
		case "Load existing project", "Load different project":
			err = LoadExistingProject()
			if err == ErrReturnToMenu {
				continue // Return to main menu
			}
			return err
		case "Show version":
			fmt.Printf("gophex version %s\n", version.GetVersion())
			continue // Stay in menu
		case "Show help":
			printHelp()
			continue // Stay in menu
		case "Quit":
			return GetProcessManager().HandleGracefulShutdown()
		case "Print image":

			f, err := files.Open("assets/my-image.jpeg")
			if err != nil {
				return fmt.Errorf("failed to open embedded image file: %w", err)
			}

			// kittyimg.Fprintln(os.Stdout, imageData)
			kittyimg.Transcode(os.Stdout, f)
			fmt.Println()
			f.Close()
			continue
		default:
			return fmt.Errorf("unknown action: %s", action)
		}
	}
}

// loadCurrentProject loads the project from the current directory
func loadCurrentProject(projectPath string, metadata *utils.ProjectMetadata) error {
	fmt.Printf("📂 Loading current project: %s (%s)\n", metadata.Project.Name, metadata.Project.Type)
	fmt.Printf("📍 Location: %s\n", projectPath)

	// Create project options for post-generation workflow
	opts := PostGenerationOptions{
		ProjectName: metadata.Project.Name,
		ProjectType: metadata.Project.Type,
		ProjectPath: projectPath,
	}

	// Show post-generation menu
	return ShowPostGenerationMenu(opts)
}

func printHelp() {
	fmt.Println("Gophex - Go Project Generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gophex                 Start interactive mode")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("  - Generate new projects with clean architecture")
	fmt.Println("  - Load existing Gophex projects by specifying directory path")
	fmt.Println("  - Track project activities with 're-' prefix for repeated actions")
	fmt.Println("  - Database configuration without storing sensitive data")
	fmt.Println("  - Smart validation with helpful error messages")
	fmt.Println()
	fmt.Println("Supported project types:")
	fmt.Println("  - api: REST API with clean architecture")
	fmt.Println("  - webapp: Web application with templates")
	fmt.Println("  - microservice: Microservice with gRPC support")
	fmt.Println("  - cli: Command-line tool")
}
