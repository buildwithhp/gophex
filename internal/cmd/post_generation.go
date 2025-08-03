package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// PostGenerationOptions represents the available actions after project generation
type PostGenerationOptions struct {
	ProjectPath string
	ProjectType string
	ProjectName string
}

// ShowPostGenerationMenu displays the post-generation menu and handles user choices
func ShowPostGenerationMenu(opts PostGenerationOptions) error {
	for {
		fmt.Printf("\n✅ Project '%s' is ready at %s\n", opts.ProjectName, opts.ProjectPath)

		// Show helpful information based on project type
		if opts.ProjectType == "api" {
			fmt.Println("🌐 API project with database integration")
			fmt.Println("   • Database migrations ready")
			fmt.Println("   • JWT authentication included")
			fmt.Println("   • Health check endpoint available")
		} else if opts.ProjectType == "webapp" {
			fmt.Println("🌍 Web application with templates")
		} else if opts.ProjectType == "microservice" {
			fmt.Println("🔧 Microservice with health checks")
		} else if opts.ProjectType == "cli" {
			fmt.Println("💻 Command-line application")
		}
		fmt.Println()

		var choice string
		menuPrompt := &survey.Select{
			Message: "🚀 What would you like to do next?",
			Options: []string{
				"⚡ Quick start (install deps + start app)",
				"🔄 Development workflow (full auto-setup)",
				"📁 Open project directory",
				"🗄️ Run database migrations/initialization",
				"📦 Install dependencies (go mod tidy)",
				"🚀 Start the application",
				"🧪 Run tests",
				"📖 View project documentation",
				"🔍 Run change detection",
				"🆕 Generate another project",
				"Quit",
			},
		}

		err := survey.AskOne(menuPrompt, &choice)
		if err != nil {
			// Handle user interruption (Ctrl+C) gracefully
			if isUserInterrupt(err) {
				fmt.Println("\n👋 Thank you for using Gophex!")
				return nil
			}
			return fmt.Errorf("menu selection failed: %w", err)
		}

		// Handle the selected option
		switch {
		case choice[:2] == "⚡":
			if err := RunQuickStart(opts.ProjectPath, opts.ProjectType); err != nil {
				fmt.Printf("❌ Quick start failed: %v\n", err)
			}
		case choice[:4] == "🔄":
			if err := RunDevelopmentWorkflow(opts.ProjectPath, opts.ProjectType); err != nil {
				fmt.Printf("❌ Development workflow failed: %v\n", err)
			}
		case choice[:4] == "📁":
			if err := OpenProjectDirectory(opts.ProjectPath); err != nil {
				fmt.Printf("❌ Error opening directory: %v\n", err)
			}
		case choice[:4] == "🗄️":
			if err := RunDatabaseSetup(opts.ProjectPath, opts.ProjectType); err != nil {
				fmt.Printf("❌ Database setup failed: %v\n", err)
			}
		case choice[:4] == "📦":
			if err := InstallDependencies(opts.ProjectPath); err != nil {
				fmt.Printf("❌ Dependency installation failed: %v\n", err)
			}
		case choice[:4] == "🚀":
			if err := StartApplication(opts.ProjectPath, opts.ProjectType); err != nil {
				fmt.Printf("❌ Failed to start application: %v\n", err)
			}
		case choice[:4] == "🧪":
			if err := RunTests(opts.ProjectPath); err != nil {
				fmt.Printf("❌ Tests failed: %v\n", err)
			}
		case choice[:4] == "📖":
			if err := ViewDocumentation(opts.ProjectPath); err != nil {
				fmt.Printf("❌ Error viewing documentation: %v\n", err)
			}
		case choice[:4] == "🔍":
			if err := RunChangeDetection(opts.ProjectPath); err != nil {
				fmt.Printf("❌ Change detection failed: %v\n", err)
			}
		case choice[:4] == "🆕":
			// Generate another project
			return GenerateProject()
		case choice == "Quit":
			return GetProcessManager().HandleGracefulShutdown()
		}

		// Ask if user wants to continue or exit
		var continueMenu bool
		continuePrompt := &survey.Confirm{
			Message: "Return to menu?",
			Default: true,
		}

		err = survey.AskOne(continuePrompt, &continueMenu)
		if err != nil {
			// Handle user interruption (Ctrl+C) gracefully
			if isUserInterrupt(err) {
				fmt.Println("\n👋 Thank you for using Gophex!")
				return nil
			}
			return fmt.Errorf("continue prompt failed: %w", err)
		}

		if !continueMenu {
			fmt.Println("👋 Thank you for using Gophex!")
			return nil
		}
	}
}
