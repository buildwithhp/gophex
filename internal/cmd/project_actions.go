package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

// OpenProjectDirectory opens the project directory in the system file manager
func OpenProjectDirectory(projectPath string) error {
	fmt.Printf("📁 Opening project directory: %s\n", projectPath)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", projectPath)
	case "linux":
		cmd = exec.Command("xdg-open", projectPath)
	case "windows":
		cmd = exec.Command("explorer", projectPath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open directory: %w", err)
	}

	fmt.Println("✅ Project directory opened successfully")
	return nil
}

// RunDatabaseSetup runs database migrations or initialization based on project type
func RunDatabaseSetup(projectPath, projectType string) error {
	if projectType != "api" {
		fmt.Println("ℹ️  Database setup is only available for API projects")
		return nil
	}

	fmt.Println("🗄️ Setting up database...")

	// Get the appropriate migration script for the platform
	migrateScript, err := getMigrationScript(projectPath)
	if err != nil {
		return fmt.Errorf("migration script not found: %w", err)
	}

	// Make script executable (Unix/Linux/macOS only)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(migrateScript, 0755); err != nil {
			return fmt.Errorf("failed to make migrate script executable: %w", err)
		}
	}

	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(projectPath); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Check database type from metadata
	dbType, err := getDatabaseTypeFromMetadata(projectPath)
	if err != nil {
		return fmt.Errorf("failed to determine database type: %w", err)
	}

	// For SQL databases, check if golang-migrate is installed
	if dbType != "mongodb" {
		if err := ensureGolangMigrateInstalled(dbType); err != nil {
			return fmt.Errorf("failed to ensure golang-migrate is available: %w", err)
		}
		fmt.Println("✅ Migration tool is ready")
	}

	// Run appropriate database setup command
	var cmd *exec.Cmd
	if dbType == "mongodb" {
		// Check if MongoDB shell is available
		if err := ensureMongoShellAvailable(); err != nil {
			return fmt.Errorf("MongoDB setup requires MongoDB shell: %w", err)
		}
		fmt.Println("🍃 Initializing MongoDB collections and indexes...")
		cmd = executeScript(migrateScript, "init")
	} else {
		fmt.Println("🐘 Running database migrations...")
		cmd = executeScript(migrateScript, "up")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Check if the error is related to missing golang-migrate
		if strings.Contains(err.Error(), "golang-migrate") || strings.Contains(err.Error(), "migrate") {
			fmt.Println("⚠️  Migration tool issue detected. Attempting to resolve...")

			// Try to install golang-migrate again
			if installErr := ensureGolangMigrateInstalled(dbType); installErr != nil {
				return fmt.Errorf("database setup failed and could not install migration tool: %w", err)
			}

			// Retry the migration
			fmt.Println("🔄 Retrying database setup...")
			if retryErr := cmd.Run(); retryErr != nil {
				return fmt.Errorf("database setup failed after installing migration tool: %w", retryErr)
			}
		} else {
			return fmt.Errorf("database setup failed: %w", err)
		}
	}

	fmt.Println("✅ Database setup completed successfully")
	return nil
}

// InstallDependencies runs go mod tidy to install project dependencies
func InstallDependencies(projectPath string) error {
	fmt.Println("📦 Installing dependencies...")

	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(projectPath); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}

	fmt.Println("✅ Dependencies installed successfully")
	return nil
}

// StartApplication starts the generated application
func StartApplication(projectPath, projectType string) error {
	fmt.Println("🚀 Starting application...")

	// Check if dependencies are installed
	if !checkDependenciesInstalled(projectPath) {
		var installChoice string
		installPrompt := &survey.Select{
			Message: "Dependencies not installed. What would you like to do?",
			Options: []string{
				"Yes - Install dependencies now",
				"No - Skip (application may not start)",
				"Quit",
			},
		}

		if err := survey.AskOne(installPrompt, &installChoice); err != nil {
			if isUserInterrupt(err) {
				return GetProcessManager().HandleGracefulShutdown()
			}
			return err
		}

		// Handle quit option
		if installChoice == "Quit" {
			return GetProcessManager().HandleGracefulShutdown()
		}

		installDeps := installChoice[:3] == "Yes"

		if installDeps {
			if err := InstallDependencies(projectPath); err != nil {
				return fmt.Errorf("failed to install dependencies: %w", err)
			}
		} else {
			return fmt.Errorf("dependencies required to start application")
		}
	}

	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(projectPath); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Determine the main file based on project type
	var mainFile string
	switch projectType {
	case "api":
		mainFile = "cmd/api/main.go"
	case "webapp":
		mainFile = "cmd/webapp/main.go"
	case "microservice":
		mainFile = "cmd/server/main.go"
	case "cli":
		mainFile = "cmd/main.go"
	default:
		return fmt.Errorf("unsupported project type: %s", projectType)
	}

	// Check if main file exists
	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		return fmt.Errorf("main file not found: %s", mainFile)
	}

	fmt.Printf("🎯 Starting %s application...\n", projectType)
	fmt.Println("📝 Application logs:")
	fmt.Println("---")

	// Start the application
	cmd := exec.Command("go", "run", mainFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command with process tracking
	processName := fmt.Sprintf("%s-app", projectType)
	processDesc := fmt.Sprintf("%s application", strings.Title(projectType))

	pm := GetProcessManager()
	if err := pm.StartProcessWithTracking(processName, processDesc, projectPath, cmd); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	// Give the application time to start
	time.Sleep(2 * time.Second)

	// For API projects, show helpful information
	if projectType == "api" {
		fmt.Println("\n🌐 API Server Information:")
		fmt.Println("   Health Check: http://localhost:8080/api/v1/health")
		fmt.Println("   API Base URL: http://localhost:8080/api/v1")
		fmt.Println("\n📋 Available Endpoints:")
		fmt.Println("   POST /api/v1/auth/register - User registration")
		fmt.Println("   POST /api/v1/auth/login    - User login")
		fmt.Println("   GET  /api/v1/posts        - List posts")
		fmt.Println("   GET  /api/v1/health       - Health check")
	}

	fmt.Println("\n⚠️  Application is running in the background.")
	fmt.Println("   Press Ctrl+C in the terminal to stop it.")

	// Ask if user wants to test the health endpoint (for API projects)
	if projectType == "api" {
		var testChoice string
		testPrompt := &survey.Select{
			Message: "Test the health endpoint now?",
			Options: []string{
				"Yes - Test health endpoint",
				"No - Skip test",
				"Quit",
			},
		}

		if err := survey.AskOne(testPrompt, &testChoice); err != nil {
			if isUserInterrupt(err) {
				return GetProcessManager().HandleGracefulShutdown()
			}
			// Don't fail the whole function if health test prompt fails
		} else {
			// Handle quit option
			if testChoice == "Quit" {
				return GetProcessManager().HandleGracefulShutdown()
			}

			if testChoice[:3] == "Yes" {
				time.Sleep(1 * time.Second) // Give server more time
				if err := testHealthEndpoint(); err != nil {
					fmt.Printf("❌ Health check failed: %v\n", err)
				}
			}
		}
	}

	return nil
}

// RunTests runs the project tests
func RunTests(projectPath string) error {
	fmt.Println("🧪 Running tests...")

	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(projectPath); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Run tests
	cmd := exec.Command("go", "test", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	fmt.Println("✅ All tests passed")
	return nil
}

// ViewDocumentation displays the project documentation
func ViewDocumentation(projectPath string) error {
	fmt.Println("📖 Viewing project documentation...")

	readmePath := filepath.Join(projectPath, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		return fmt.Errorf("README.md not found in project")
	}

	// Read and display README content
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("failed to read README.md: %w", err)
	}

	fmt.Println("📄 Project README:")
	fmt.Println("==================")
	fmt.Println(string(content))

	return nil
}

// RunChangeDetection runs the change detection script
func RunChangeDetection(projectPath string) error {
	fmt.Println("🔍 Running change detection...")

	// Get the appropriate change detection script for the platform
	scriptPath, err := getChangeDetectionScript(projectPath)
	if err != nil {
		fmt.Println("ℹ️  Change detection script not available for this project type")
		return nil
	}

	// Make script executable (Unix/Linux/macOS only)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(scriptPath, 0755); err != nil {
			return fmt.Errorf("failed to make script executable: %w", err)
		}
	}

	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(projectPath); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Run change detection
	cmd := executeScript(scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("change detection failed: %w", err)
	}

	return nil
}

// Helper functions

// getDatabaseTypeFromMetadata reads the database type from the project metadata
func getDatabaseTypeFromMetadata(projectPath string) (string, error) {
	metadataPath := filepath.Join(projectPath, ".gophex-generated")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return "postgresql", nil // Default fallback
	}

	file, err := os.Open(metadataPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "database_type=") {
			return strings.TrimPrefix(line, "database_type="), nil
		}
	}

	return "postgresql", nil // Default fallback
}

// ensureGolangMigrateInstalled checks if golang-migrate is installed and offers to install it
func ensureGolangMigrateInstalled(dbType string) error {
	// Check if golang-migrate is already installed
	if isGolangMigrateInstalled() {
		return nil
	}

	fmt.Println("⚠️  golang-migrate tool is not installed")
	fmt.Printf("   This tool is required for %s database migrations\n", dbType)

	var installMigrate string
	installPrompt := &survey.Select{
		Message: "Would you like Gophex to install golang-migrate for you?",
		Options: []string{
			"Yes - Install golang-migrate tool",
			"No - Skip installation",
			"Quit",
		},
		Help: "This will install the golang-migrate tool using 'go install'",
	}

	if err := survey.AskOne(installPrompt, &installMigrate); err != nil {
		return err
	}

	if installMigrate == "Quit" {
		return nil
	}

	if installMigrate[:2] == "No" {
		fmt.Println("❌ Database migrations require golang-migrate tool")
		fmt.Printf("   You can install it manually with: go install -tags '%s' github.com/golang-migrate/migrate/v4/cmd/migrate@latest\n", dbType)
		return fmt.Errorf("golang-migrate tool is required but not installed")
	}

	return installGolangMigrate(dbType)
}

// isGolangMigrateInstalled checks if golang-migrate is available in PATH
func isGolangMigrateInstalled() bool {
	_, err := exec.LookPath("migrate")
	return err == nil
}

// installGolangMigrate installs golang-migrate using go install
func installGolangMigrate(dbType string) error {
	fmt.Println("📦 Installing golang-migrate tool...")
	fmt.Println("   This may take a few moments depending on your internet connection...")

	// Check if Go is available
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("Go is not installed or not available in PATH. Please install Go first")
	}

	// Determine the appropriate tags for the database type
	var tags string
	switch dbType {
	case "postgresql":
		tags = "postgres"
	case "mysql":
		tags = "mysql"
	default:
		tags = "postgres" // Default fallback
	}

	// Install golang-migrate with appropriate database tags
	installCmd := fmt.Sprintf("go install -tags '%s' github.com/golang-migrate/migrate/v4/cmd/migrate@latest", tags)

	fmt.Printf("   Running: %s\n", installCmd)
	fmt.Println("   📡 Downloading and compiling...")

	cmd := exec.Command("go", "install", "-tags", tags, "github.com/golang-migrate/migrate/v4/cmd/migrate@latest")

	// Capture both stdout and stderr for better error reporting
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("   ❌ Installation failed with output:\n%s\n", string(output))
		return fmt.Errorf("failed to install golang-migrate: %w", err)
	}

	// Verify installation
	if !isGolangMigrateInstalled() {
		fmt.Println("   ⚠️  Installation completed but tool is not available in PATH")
		fmt.Println("   💡 Try running the command in a new terminal or check your GOPATH/GOBIN settings")
		return fmt.Errorf("golang-migrate installation completed but tool is not available in PATH")
	}

	fmt.Println("✅ golang-migrate installed successfully!")

	// Show version information
	versionCmd := exec.Command("migrate", "-version")
	if output, err := versionCmd.Output(); err == nil {
		fmt.Printf("   📋 Version: %s\n", strings.TrimSpace(string(output)))
	}

	fmt.Printf("   🎯 Ready for %s database migrations\n", dbType)
	fmt.Println("   🚀 Continuing with database setup...")
	return nil
}

// ensureMongoShellAvailable checks if MongoDB shell is available
func ensureMongoShellAvailable() error {
	// Check for mongosh (MongoDB 5.0+)
	if _, err := exec.LookPath("mongosh"); err == nil {
		return nil
	}

	// Check for legacy mongo shell
	if _, err := exec.LookPath("mongo"); err == nil {
		return nil
	}

	fmt.Println("⚠️  MongoDB shell (mongosh or mongo) is not installed")
	fmt.Println("   MongoDB initialization requires a MongoDB shell to run scripts")
	fmt.Println()
	fmt.Println("📋 Installation options:")
	fmt.Println("   • Install MongoDB Community Edition (includes shell)")
	fmt.Println("   • Install MongoDB Shell separately: https://docs.mongodb.com/mongodb-shell/install/")
	fmt.Println("   • Use package manager:")

	switch runtime.GOOS {
	case "darwin":
		fmt.Println("     brew install mongosh")
	case "linux":
		fmt.Println("     # Ubuntu/Debian: apt install mongodb-mongosh")
		fmt.Println("     # CentOS/RHEL: yum install mongodb-mongosh")
	case "windows":
		fmt.Println("     # Download from: https://www.mongodb.com/try/download/shell")
	}

	fmt.Println()
	fmt.Println("💡 After installation, you can run database setup from the menu")

	return fmt.Errorf("MongoDB shell not available")
}

// checkDependenciesInstalled checks if go.mod and go.sum exist and are up to date
func checkDependenciesInstalled(projectPath string) bool {
	goModPath := filepath.Join(projectPath, "go.mod")
	goSumPath := filepath.Join(projectPath, "go.sum")

	// Check if go.mod exists
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return false
	}

	// Check if go.sum exists (indicates dependencies have been downloaded)
	if _, err := os.Stat(goSumPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// getMigrationScript returns the appropriate migration script path for the current platform
func getMigrationScript(projectPath string) (string, error) {
	var scriptName string
	if runtime.GOOS == "windows" {
		scriptName = "migrate.bat"
	} else {
		scriptName = "migrate.sh"
	}

	scriptPath := filepath.Join(projectPath, "scripts", scriptName)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("migration script not found: %s", scriptPath)
	}

	return scriptPath, nil
}

// getChangeDetectionScript returns the appropriate change detection script path for the current platform
func getChangeDetectionScript(projectPath string) (string, error) {
	var scriptName string
	if runtime.GOOS == "windows" {
		scriptName = "detect-changes.bat"
	} else {
		scriptName = "detect-changes.sh"
	}

	scriptPath := filepath.Join(projectPath, "scripts", scriptName)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("change detection script not found: %s", scriptPath)
	}

	return scriptPath, nil
}

// executeScript runs a script with the appropriate command for the platform
func executeScript(scriptPath string, args ...string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		// For Windows batch files
		cmdArgs := append([]string{"/c", scriptPath}, args...)
		return exec.Command("cmd", cmdArgs...)
	} else {
		// For Unix shell scripts
		cmdArgs := append([]string{scriptPath}, args...)
		return exec.Command("bash", cmdArgs...)
	}
}

// testHealthEndpoint tests the API health endpoint using HTTP client (cross-platform)
func testHealthEndpoint() error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("http://localhost:8080/api/v1/health")
	if err != nil {
		return fmt.Errorf("failed to connect to health endpoint: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("✅ Health check response (%d): %s\n", resp.StatusCode, string(body))
	return nil
}
