package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/buildwithhp/gophex/internal/utils"
)

// PostGenerationOptions represents the available actions after project generation
type PostGenerationOptions struct {
	ProjectPath string
	ProjectType string
	ProjectName string
}

// ShowPostGenerationMenu displays the post-generation menu and handles user choices
func ShowPostGenerationMenu(opts PostGenerationOptions) error {
	// Initialize project tracker
	tracker := NewProjectTracker(opts.ProjectPath)
	if err := tracker.LoadMetadata(); err != nil {
		fmt.Printf("⚠️  Warning: Could not load project metadata: %v\n", err)
	}

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
			Options: buildMenuOptions(tracker),
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
			} else {
				// Quick start includes multiple activities - update both systems
				if utils.HasGophexMetadata(opts.ProjectPath) {
					utils.UpdateActivity(opts.ProjectPath, "dependencies_installed", true)
					utils.UpdateActivity(opts.ProjectPath, "database_migrated", true)
					utils.UpdateActivity(opts.ProjectPath, "application_started", true)
				}
				tracker.UpdateActivity("dependencies_installed", true)
				tracker.UpdateActivity("database_setup", true)
				tracker.UpdateActivity("application_started", true)
			}
		case choice[:4] == "🔄":
			if err := RunDevelopmentWorkflow(opts.ProjectPath, opts.ProjectType); err != nil {
				fmt.Printf("❌ Development workflow failed: %v\n", err)
			} else {
				// Development workflow includes all activities - update both systems
				if utils.HasGophexMetadata(opts.ProjectPath) {
					utils.UpdateActivity(opts.ProjectPath, "dependencies_installed", true)
					utils.UpdateActivity(opts.ProjectPath, "database_migrated", true)
					utils.UpdateActivity(opts.ProjectPath, "application_started", true)
					utils.UpdateActivity(opts.ProjectPath, "tests_executed", true)
				}
				tracker.UpdateActivity("dependencies_installed", true)
				tracker.UpdateActivity("database_setup", true)
				tracker.UpdateActivity("application_started", true)
				tracker.UpdateActivity("tests_run", true)
				tracker.UpdateActivity("health_check_tested", true)
			}
		case choice[:4] == "📁":
			if err := OpenProjectDirectory(opts.ProjectPath); err != nil {
				fmt.Printf("❌ Error opening directory: %v\n", err)
			}
		case choice[:4] == "🗄️":
			if err := RunDatabaseSetup(opts.ProjectPath, opts.ProjectType); err != nil {
				fmt.Printf("❌ Database setup failed: %v\n", err)
			} else {
				tracker.UpdateActivity("database_setup", true)
				tracker.UpdateDatabaseStatus(true, true)
			}
		case choice[:4] == "📦":
			if err := InstallDependencies(opts.ProjectPath); err != nil {
				fmt.Printf("❌ Dependency installation failed: %v\n", err)
			} else {
				tracker.UpdateActivity("dependencies_installed", true)
			}
		case choice[:4] == "🚀":
			if err := StartApplication(opts.ProjectPath, opts.ProjectType); err != nil {
				fmt.Printf("❌ Failed to start application: %v\n", err)
			} else {
				tracker.UpdateActivity("application_started", true)
			}
		case choice[:4] == "🧪":
			if err := RunTests(opts.ProjectPath); err != nil {
				fmt.Printf("❌ Tests failed: %v\n", err)
			} else {
				tracker.UpdateActivity("tests_run", true)
			}
		case choice[:4] == "📖":
			if err := ViewDocumentation(opts.ProjectPath); err != nil {
				fmt.Printf("❌ Error viewing documentation: %v\n", err)
			} else {
				tracker.UpdateActivity("documentation_viewed", true)
			}
		case choice[:4] == "🔍":
			if err := RunChangeDetection(opts.ProjectPath); err != nil {
				fmt.Printf("❌ Change detection failed: %v\n", err)
			} else {
				tracker.UpdateActivity("change_detection_run", true)
			}
		case choice[:4] == "🏗️":
			if err := RunCRUDWizard(opts.ProjectPath); err != nil {
				if err == ErrReturnToMenu {
					continue // Return to menu
				}
				fmt.Printf("❌ CRUD generation failed: %v\n", err)
			} else {
				tracker.UpdateActivity("crud_generated", true)
			}
		case choice[:4] == "🎓":
			if err := RunEnhancedCRUDWizard(opts.ProjectPath); err != nil {
				if err == ErrReturnToMenu {
					continue // Return to menu
				}
				fmt.Printf("❌ Enhanced CRUD generation failed: %v\n", err)
			} else {
				tracker.UpdateActivity("enhanced_crud_generated", true)
			}
		case choice[:4] == "🆕":
			// Generate another project
			return GenerateProject()
		case choice == "Quit":
			return GetProcessManager().HandleGracefulShutdown()
		}

		// Ask if user wants to continue or exit
		var continueMenu string
		continuePrompt := &survey.Select{
			Message: "What would you like to do next?",
			Options: []string{
				"Return to menu",
				"Exit Gophex",
			},
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

		if continueMenu[:4] == "Exit" {
			fmt.Println("👋 Thank you for using Gophex!")
			return nil
		}
	}
}

// buildMenuOptions creates dynamic menu options based on project state
func buildMenuOptions(tracker *ProjectTracker) []string {
	projectPath := tracker.projectPath

	options := []string{
		"⚡ Quick start (install deps + start app)",
		"🔄 Development workflow (full auto-setup)",
		"📁 Open project directory",
	}

	// Check if we have gophex metadata to use the new system
	if utils.HasGophexMetadata(projectPath) {
		// Use new metadata system for activity prefixes

		// Add database-specific options if database is configured
		metadata := tracker.GetMetadata()
		if metadata.Gophex.Database.Configured {
			prefix := utils.GetActivityPrefix(projectPath, "database_migrated")
			options = append(options, fmt.Sprintf("🗄️  %sRun database migrations/initialization", prefix))
		}

		// Add dependency installation option
		prefix := utils.GetActivityPrefix(projectPath, "dependencies_installed")
		options = append(options, fmt.Sprintf("📦 %sInstall dependencies (go mod tidy)", prefix))

		// Add application start option
		prefix = utils.GetActivityPrefix(projectPath, "application_started")
		options = append(options, fmt.Sprintf("🚀 %sStart the application", prefix))

		// Add test option
		prefix = utils.GetActivityPrefix(projectPath, "tests_executed")
		options = append(options, fmt.Sprintf("🧪 %sRun tests", prefix))

		// Add documentation option
		prefix = utils.GetActivityPrefix(projectPath, "documentation_viewed")
		options = append(options, fmt.Sprintf("📖 %sView project documentation", prefix))

		// Add change detection option
		prefix = utils.GetActivityPrefix(projectPath, "change_detection_run")
		options = append(options, fmt.Sprintf("🔍 %sRun change detection", prefix))

		// Add CRUD generation option (only for API projects)
		if projectMetadata, err := utils.LoadMetadata(projectPath); err == nil && projectMetadata.Project.Type == "api" {
			prefix = utils.GetActivityPrefix(projectPath, "crud_generated")
			options = append(options, fmt.Sprintf("🏗️  %sGenerate CRUD operations", prefix))

			// Add enhanced CRUD wizard option
			prefix = utils.GetActivityPrefix(projectPath, "enhanced_crud_generated")
			options = append(options, fmt.Sprintf("🎓 %sEnhanced CRUD Wizard - Learn Clean Architecture", prefix))
		}
	} else {
		// Fallback to old system
		metadata := tracker.GetMetadata()
		if metadata.Gophex.Database.Configured {
			prefix := tracker.GetActivityPrefix("database_setup")
			options = append(options, fmt.Sprintf("🗄️  %sRun database migrations/initialization", prefix))
		}

		// Add dependency installation option
		prefix := tracker.GetActivityPrefix("dependencies_installed")
		options = append(options, fmt.Sprintf("📦 %sInstall dependencies (go mod tidy)", prefix))

		// Add application start option
		prefix = tracker.GetActivityPrefix("application_started")
		options = append(options, fmt.Sprintf("🚀 %sStart the application", prefix))

		// Add test option
		prefix = tracker.GetActivityPrefix("tests_run")
		options = append(options, fmt.Sprintf("🧪 %sRun tests", prefix))

		// Add documentation option
		prefix = tracker.GetActivityPrefix("documentation_viewed")
		options = append(options, fmt.Sprintf("📖 %sView project documentation", prefix))

		// Add change detection option
		prefix = tracker.GetActivityPrefix("change_detection_run")
		options = append(options, fmt.Sprintf("🔍 %sRun change detection", prefix))

		// Add CRUD generation option (only for API projects)
		trackerMetadata := tracker.GetMetadata()
		if trackerMetadata.Gophex.Project.Type == "api" {
			prefix = tracker.GetActivityPrefix("crud_generated")
			options = append(options, fmt.Sprintf("🏗️  %sGenerate CRUD operations", prefix))

			// Add enhanced CRUD wizard option
			prefix = tracker.GetActivityPrefix("enhanced_crud_generated")
			options = append(options, fmt.Sprintf("🎓 %sEnhanced CRUD Wizard - Learn Clean Architecture", prefix))
		}
	}

	// Add static options
	options = append(options,
		"🆕 Generate another project",
		"Quit",
	)

	return options
}
