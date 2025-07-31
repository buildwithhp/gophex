package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/buildwithhp/gophex/pkg/version"
)

func Execute() error {
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
