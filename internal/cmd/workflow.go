package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

// RunDevelopmentWorkflow runs an automated development setup workflow
func RunDevelopmentWorkflow(projectPath, projectType string) error {
	fmt.Println("🔄 Development Workflow Automation")
	fmt.Println("   This will set up your project for development automatically.")

	// Show what will be done
	fmt.Println("\n📋 Workflow steps:")
	fmt.Println("   1. 📦 Install dependencies (go mod tidy)")
	if projectType == "api" {
		fmt.Println("   2. 🗄️ Set up database (migrations/initialization)")
	}
	fmt.Println("   3. 🧪 Run tests (go test ./...)")
	fmt.Println("   4. 🚀 Start application")

	if projectType == "api" {
		fmt.Println("\n🎯 After completion:")
		fmt.Println("   • API server will be running on http://localhost:8080")
		fmt.Println("   • Health check available at /api/v1/health")
		fmt.Println("   • Ready for development and testing")
	}
	// Confirm with user
	var choice string
	confirmPrompt := &survey.Select{
		Message: "Continue with automated workflow?",
		Options: []string{
			"Yes - Start automated workflow",
			"No - Cancel workflow",
			"Quit",
		},
	}

	if err := survey.AskOne(confirmPrompt, &choice); err != nil {
		if isUserInterrupt(err) {
			return GetProcessManager().HandleGracefulShutdown()
		}
		return err
	}

	// Handle quit option
	if choice == "Quit" {
		return GetProcessManager().HandleGracefulShutdown()
	}

	proceed := choice[:3] == "Yes"

	if !proceed {
		fmt.Println("⏹️  Workflow cancelled")
		return nil
	}

	fmt.Println("\n🚀 Starting automated workflow...")

	// Step 1: Install dependencies
	fmt.Println("📦 Step 1/4: Installing dependencies...")
	if err := InstallDependencies(projectPath); err != nil {
		return fmt.Errorf("workflow failed at dependency installation: %w", err)
	}
	time.Sleep(1 * time.Second)

	// Step 2: Database setup (for API projects)
	if projectType == "api" {
		fmt.Println("\n🗄️ Step 2/4: Setting up database...")
		if err := RunDatabaseSetup(projectPath, projectType); err != nil {
			if strings.Contains(err.Error(), "golang-migrate") {
				fmt.Printf("⚠️  Database setup requires golang-migrate tool: %v\n", err)
				fmt.Println("   The tool installation was declined or failed.")
				fmt.Println("   You can set up the database manually later using the menu option.")
			} else {
				fmt.Printf("⚠️  Database setup failed: %v\n", err)
				fmt.Println("   You can set up the database manually later using the menu option.")
			}
		}
		time.Sleep(1 * time.Second)
	}

	// Step 3: Run tests
	fmt.Println("\n🧪 Step 3/4: Running tests...")
	if err := RunTests(projectPath); err != nil {
		fmt.Printf("⚠️  Tests failed: %v\n", err)
		fmt.Println("   You can run tests manually later using the menu option.")
	}
	time.Sleep(1 * time.Second)

	// Step 4: Start application
	fmt.Println("\n🚀 Step 4/4: Starting application...")

	var startApp string
	startPrompt := &survey.Select{
		Message: "Start the application now?",
		Options: []string{
			"Yes - Start the application",
			"No - Skip for now",
			"Quit",
		},
	}

	if err := survey.AskOne(startPrompt, &startApp); err != nil {
		return err
	}

	if startApp == "Quit" {
		return nil
	}

	if startApp[:3] == "Yes" {
		if err := StartApplication(projectPath, projectType); err != nil {
			fmt.Printf("⚠️  Failed to start application: %v\n", err)
			fmt.Println("   You can start the application manually later using the menu option.")
		}
	}

	fmt.Println("\n✅ Development workflow completed!")
	fmt.Println("🎉 Your project is ready for development!")

	if projectType == "api" {
		fmt.Println("\n🌐 Quick API Test:")
		fmt.Println("   curl http://localhost:8080/api/v1/health")
	}

	return nil
}

// RunQuickStart provides a simplified quick start workflow
func RunQuickStart(projectPath, projectType string) error {
	fmt.Println("⚡ Quick Start - Setting up your project...")

	// Install dependencies
	if err := InstallDependencies(projectPath); err != nil {
		return fmt.Errorf("quick start failed: %w", err)
	}

	// For API projects, try database setup
	if projectType == "api" {
		fmt.Println("🗄️ Setting up database...")
		if err := RunDatabaseSetup(projectPath, projectType); err != nil {
			if strings.Contains(err.Error(), "golang-migrate") {
				fmt.Printf("⚠️  Database setup requires golang-migrate tool: %v\n", err)
				fmt.Println("   You can set up the database manually later.")
			} else {
				fmt.Printf("⚠️  Database setup skipped: %v\n", err)
			}
		}
	}

	// Start application
	if err := StartApplication(projectPath, projectType); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	fmt.Println("✅ Quick start completed!")
	return nil
}
