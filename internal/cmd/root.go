package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/buildwithhp/gophex/pkg/version"
)

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
	var action string
	prompt := &survey.Select{
		Message: "What would you like to do?",
		Options: []string{"Generate a new project", "Load existing project", "Show version", "Show help", "Quit"},
	}

	err := survey.AskOne(prompt, &action)
	if err != nil {
		// Handle user interruption (Ctrl+C) gracefully
		if isUserInterrupt(err) {
			fmt.Println("\nGoodbye! 👋")
			return nil
		}
		return fmt.Errorf("prompt failed: %w", err)
	}

	switch action {
	case "Generate a new project":
		return GenerateProject()
	case "Load existing project":
		return LoadExistingProject()
	case "Show version":
		fmt.Printf("gophex version %s\n", version.GetVersion())
		return Execute()
	case "Show help":
		printHelp()
		return Execute()
	case "Quit":
		return GetProcessManager().HandleGracefulShutdown()
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
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
