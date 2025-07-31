package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hemant/gophex/internal/generator"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [project-type] [project-name]",
	Short: "Generate a new Go project",
	Long: `Generate a new Go project with the specified type and name.

Available project types:
- api: REST API with clean architecture
- webapp: Web application with templates
- microservice: Microservice with gRPC support
- cli: Command-line tool`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectType := args[0]
		projectName := args[1]

		if err := validateProjectType(projectType); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		projectPath := filepath.Join(currentDir, projectName)

		gen := generator.New()
		if err := gen.Generate(projectType, projectName, projectPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully generated %s project '%s' in %s\n", projectType, projectName, projectPath)
	},
}

func validateProjectType(projectType string) error {
	validTypes := []string{"api", "webapp", "microservice", "cli"}
	for _, t := range validTypes {
		if t == projectType {
			return nil
		}
	}
	return fmt.Errorf("invalid project type '%s'. Valid types: %v", projectType, validTypes)
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
