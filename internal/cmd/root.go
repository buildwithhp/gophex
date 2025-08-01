package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/buildwithhp/gophex/pkg/version"
)

func Execute() error {
	// Handle command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version", "-v":
			fmt.Printf("gophex version %s\n", version.GetVersion())
			return nil
		case "generate", "gen":
			return GenerateProject()
		case "help", "--help", "-h":
			printHelp()
			return nil
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			printHelp()
			return nil
		}
	}

	// Interactive mode
	fmt.Println("🚀 Welcome to Gophex!")
	fmt.Println("A CLI tool for generating Go project scaffolding")
	fmt.Println()

	var action string
	prompt := &survey.Select{
		Message: "What would you like to do?",
		Options: []string{"Generate a new project", "Show version", "Exit"},
	}

	err := survey.AskOne(prompt, &action)
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}

	switch action {
	case "Generate a new project":
		return GenerateProject()
	case "Show version":
		fmt.Printf("gophex version %s\n", version.GetVersion())
		return Execute()
	case "Exit":
		fmt.Println("Goodbye! 👋")
		return nil
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

func printHelp() {
	fmt.Println("Gophex - Go Project Generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gophex                 Start interactive mode")
	fmt.Println("  gophex generate        Generate a new project interactively")
	fmt.Println("  gophex version         Show version")
	fmt.Println("  gophex help            Show this help")
	fmt.Println()
	fmt.Println("Supported project types:")
	fmt.Println("  - api: REST API with clean architecture")
	fmt.Println("  - webapp: Web application with templates")
	fmt.Println("  - microservice: Microservice with gRPC support")
	fmt.Println("  - cli: Command-line tool")
}
