package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/buildwithhp/gophex/internal/generator"
)

func GenerateProject() error {
	// Offer choice between quick generation and educational wizard
	var approach string
	approachPrompt := &survey.Select{
		Message: "How would you like to generate your project?",
		Options: []string{
			"🎓 Educational Wizard - Learn step-by-step (recommended)",
			"⚡ Quick Generation - Direct project creation",
			"Quit",
		},
		Help: "The Educational Wizard teaches Go architecture patterns while building your project",
	}

	err := survey.AskOne(approachPrompt, &approach)
	if err != nil {
		if isUserInterrupt(err) {
			fmt.Println("\nProject generation cancelled. Goodbye! 👋")
			return nil
		}
		return fmt.Errorf("approach selection failed: %w", err)
	}

	if approach == "Quit" {
		return GetProcessManager().HandleGracefulShutdown()
	}

	if strings.HasPrefix(approach, "🎓") {
		// Use the enhanced educational wizard
		return RunEnhancedProjectWizard()
	}

	// Continue with quick generation for users who want the old behavior
	return runQuickProjectGeneration()
}

// runQuickProjectGeneration provides the original quick generation experience
func runQuickProjectGeneration() error {
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
			"Quit",
		},
	}

	err := survey.AskOne(projectTypePrompt, &projectType)
	if err != nil {
		// Handle user interruption (Ctrl+C) gracefully
		if isUserInterrupt(err) {
			fmt.Println("\nProject generation cancelled. Goodbye! 👋")
			return nil
		}
		return fmt.Errorf("project type selection failed: %w", err)
	}

	// Handle quit option
	if projectType == "Quit" {
		return GetProcessManager().HandleGracefulShutdown()
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
		// Handle user interruption (Ctrl+C) gracefully
		if isUserInterrupt(err) {
			return GetProcessManager().HandleGracefulShutdown()
		}
		return fmt.Errorf("project name input failed: %w", err)
	}

	// Get framework and database configuration for API projects
	var framework string
	var dbConfig *generator.DatabaseConfig
	var redisConfig *generator.RedisConfig
	if projectType == "api" {
		framework, err = getFrameworkConfiguration()
		if err != nil {
			return fmt.Errorf("framework configuration failed: %w", err)
		}

		dbConfig, err = getDatabaseConfiguration(projectName)
		if err != nil {
			return fmt.Errorf("database configuration failed: %w", err)
		}

		redisConfig, err = getRedisConfiguration()
		if err != nil {
			return fmt.Errorf("redis configuration failed: %w", err)
		}
	}

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	projectPath := filepath.Join(currentDir, projectName)

	// Path confirmation loop
	for {
		var confirm string
		confirmPrompt := &survey.Select{
			Message: fmt.Sprintf("Generate %s project '%s' in %s?", projectType, projectName, projectPath),
			Options: []string{
				"Yes - Generate project",
				"No - Change settings",
				"Quit",
			},
		}

		err = survey.AskOne(confirmPrompt, &confirm)
		if err != nil {
			if isUserInterrupt(err) {
				return GetProcessManager().HandleGracefulShutdown()
			}
			return fmt.Errorf("confirmation failed: %w", err)
		}

		if confirm == "Quit" {
			return GetProcessManager().HandleGracefulShutdown()
		}

		if confirm[:3] == "Yes" {
			break // User confirmed, proceed with generation
		}

		// User said no, ask if they want to change the path or cancel
		var action string
		actionPrompt := &survey.Select{
			Message: "What would you like to do?",
			Options: []string{
				"Change directory path",
				"Cancel project generation",
				"Quit",
			},
		}

		err = survey.AskOne(actionPrompt, &action)
		if err != nil {
			if isUserInterrupt(err) {
				return GetProcessManager().HandleGracefulShutdown()
			}
			return fmt.Errorf("action selection failed: %w", err)
		}

		// Handle quit option
		if action == "Quit" {
			return GetProcessManager().HandleGracefulShutdown()
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
	if err := gen.GenerateWithFramework(projectType, projectName, projectPath, framework, dbConfig, redisConfig); err != nil {
		return fmt.Errorf("error generating project: %w", err)
	}

	// Create project tracking metadata
	tracker := NewProjectTracker(projectPath)
	if err := tracker.CreateInitialMetadata(projectType, projectName, projectPath, dbConfig, redisConfig); err != nil {
		fmt.Printf("⚠️  Warning: Failed to create project tracking metadata: %v\n", err)
		// Don't fail the entire generation for this
	}

	fmt.Printf("✅ Successfully generated %s project '%s' in %s\n", projectType, projectName, projectPath)

	// Show post-generation menu
	opts := PostGenerationOptions{
		ProjectPath: projectPath,
		ProjectType: projectType,
		ProjectName: projectName,
	}

	return ShowPostGenerationMenu(opts)
}

func getDatabaseConfiguration(projectName string) (*generator.DatabaseConfig, error) {
	config := &generator.DatabaseConfig{}

	// Ask for database type
	var dbType string
	dbTypePrompt := &survey.Select{
		Message: "Which database would you like to use?",
		Options: []string{
			"PostgreSQL - Advanced open-source relational database",
			"MySQL - Popular open-source relational database",
			"MongoDB - Document-oriented NoSQL database",
			"Quit",
		},
	}

	err := survey.AskOne(dbTypePrompt, &dbType)
	if err != nil {
		if isUserInterrupt(err) {
			return nil, GetProcessManager().HandleGracefulShutdown()
		}
		return nil, fmt.Errorf("database type selection failed: %w", err)
	}

	// Handle quit option
	if dbType == "Quit" {
		return nil, GetProcessManager().HandleGracefulShutdown()
	}

	// Extract database type
	switch {
	case dbType[:10] == "PostgreSQL":
		config.Type = "postgresql"
	case dbType[:5] == "MySQL":
		config.Type = "mysql"
	case dbType[:7] == "MongoDB":
		config.Type = "mongodb"
	}

	// Ask for configuration type
	var configType string
	configTypePrompt := &survey.Select{
		Message: "What type of database configuration do you need?",
		Options: []string{
			"Single instance - Simple single database server",
			"Read-Write split - Separate read and write endpoints",
			"Cluster - Multiple database nodes",
			"Quit",
		},
	}

	err = survey.AskOne(configTypePrompt, &configType)
	if err != nil {
		if isUserInterrupt(err) {
			return nil, GetProcessManager().HandleGracefulShutdown()
		}
		return nil, fmt.Errorf("configuration type selection failed: %w", err)
	}

	// Handle quit option
	if configType == "Quit" {
		return nil, GetProcessManager().HandleGracefulShutdown()
	}

	// Extract configuration type
	switch {
	case configType[:6] == "Single":
		config.ConfigType = "single"
	case configType[:10] == "Read-Write":
		config.ConfigType = "read-write"
	case configType[:7] == "Cluster":
		config.ConfigType = "cluster"
	}

	// Get database credentials and connection details
	err = getDatabaseCredentials(config, projectName)
	if err != nil {
		return nil, fmt.Errorf("failed to get database credentials: %w", err)
	}

	return config, nil
}

func getRedisConfiguration() (*generator.RedisConfig, error) {
	config := &generator.RedisConfig{}

	// Ask if user wants Redis
	var redisChoice string
	redisPrompt := &survey.Select{
		Message: "Do you want to include Redis for caching and session storage?",
		Options: []string{
			"Yes - Include Redis support",
			"No - Skip Redis",
			"Quit",
		},
		Help: "Redis provides high-performance caching, session storage, and pub/sub capabilities",
	}

	err := survey.AskOne(redisPrompt, &redisChoice)
	if err != nil {
		if isUserInterrupt(err) {
			return nil, GetProcessManager().HandleGracefulShutdown()
		}
		return nil, fmt.Errorf("redis selection failed: %w", err)
	}

	// Handle quit option
	if redisChoice == "Quit" {
		return nil, GetProcessManager().HandleGracefulShutdown()
	}

	wantsRedis := redisChoice[:3] == "Yes"

	config.Enabled = wantsRedis

	// If user wants Redis, get connection details
	if wantsRedis {
		// Redis host
		hostPrompt := &survey.Input{
			Message: "Redis host:",
			Default: "localhost",
			Help:    "The hostname or IP address of your Redis server",
		}
		err = survey.AskOne(hostPrompt, &config.Host, survey.WithValidator(survey.Required))
		if err != nil {
			return nil, err
		}

		// Redis port
		portPrompt := &survey.Input{
			Message: "Redis port:",
			Default: "6379",
			Help:    "The port number for your Redis server",
		}
		err = survey.AskOne(portPrompt, &config.Port, survey.WithValidator(survey.Required))
		if err != nil {
			return nil, err
		}

		// Redis password (optional)
		passwordPrompt := &survey.Password{
			Message: "Redis password (leave empty if no password):",
		}
		err = survey.AskOne(passwordPrompt, &config.Password)
		if err != nil {
			return nil, err
		}

		// Redis database number
		var dbNumber string
		dbPrompt := &survey.Input{
			Message: "Redis database number:",
			Default: "0",
			Help:    "Redis database number (0-15, typically use 0)",
		}
		err = survey.AskOne(dbPrompt, &dbNumber, survey.WithValidator(survey.Required))
		if err != nil {
			return nil, err
		}

		// Convert database number to int
		if dbNumber == "0" {
			config.Database = 0
		} else {
			// For simplicity, we'll just use 0 for now
			config.Database = 0
		}
	}

	return config, nil
}

func getDatabaseCredentials(config *generator.DatabaseConfig, projectName string) error {
	// Database name
	dbNamePrompt := &survey.Input{
		Message: "Database name:",
		Default: projectName,
		Help:    "The name of the database to connect to",
	}
	err := survey.AskOne(dbNamePrompt, &config.DatabaseName, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	// Username
	usernamePrompt := &survey.Input{
		Message: "Database username:",
		Default: "admin",
	}
	err = survey.AskOne(usernamePrompt, &config.Username, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	// Password
	passwordPrompt := &survey.Password{
		Message: "Database password:",
	}
	err = survey.AskOne(passwordPrompt, &config.Password, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	// Configuration-specific prompts
	switch config.ConfigType {
	case "single":
		return getSingleInstanceConfig(config)
	case "read-write":
		return getReadWriteConfig(config)
	case "cluster":
		return getClusterConfig(config)
	}

	return nil
}

func getSingleInstanceConfig(config *generator.DatabaseConfig) error {
	// Host
	hostPrompt := &survey.Input{
		Message: "Database host:",
		Default: "localhost",
	}
	err := survey.AskOne(hostPrompt, &config.Host, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	// Port
	var defaultPort string
	switch config.Type {
	case "postgresql":
		defaultPort = "5432"
	case "mysql":
		defaultPort = "3306"
	case "mongodb":
		defaultPort = "27017"
	}

	portPrompt := &survey.Input{
		Message: "Database port:",
		Default: defaultPort,
	}
	err = survey.AskOne(portPrompt, &config.Port, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	// SSL Mode for PostgreSQL/MySQL
	if config.Type == "postgresql" || config.Type == "mysql" {
		var sslMode string
		sslPrompt := &survey.Select{
			Message: "SSL Mode:",
			Options: []string{"disable", "require", "verify-ca", "verify-full"},
			Default: "disable",
		}
		err = survey.AskOne(sslPrompt, &sslMode)
		if err != nil {
			return err
		}
		config.SSLMode = sslMode
	}

	// MongoDB specific settings
	if config.Type == "mongodb" {
		authSourcePrompt := &survey.Input{
			Message: "Auth source (optional):",
			Default: "admin",
		}
		survey.AskOne(authSourcePrompt, &config.AuthSource)
	}

	return nil
}

func getReadWriteConfig(config *generator.DatabaseConfig) error {
	// Write host
	writeHostPrompt := &survey.Input{
		Message: "Write database host:",
		Default: "localhost",
	}
	err := survey.AskOne(writeHostPrompt, &config.WriteHost, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	// Ask if read host is the same as write host
	var sameHost string
	sameHostPrompt := &survey.Select{
		Message: "Use the same host for read operations?",
		Options: []string{
			"Yes - Use same host for read operations",
			"No - Use different host for read operations",
			"Quit",
		},
	}
	err = survey.AskOne(sameHostPrompt, &sameHost)
	if err != nil {
		return err
	}

	if sameHost == "Quit" {
		return GetProcessManager().HandleGracefulShutdown()
	}

	if sameHost[:3] == "Yes" {
		config.ReadHost = config.WriteHost
	} else {
		readHostPrompt := &survey.Input{
			Message: "Read database host:",
			Default: "localhost",
		}
		err = survey.AskOne(readHostPrompt, &config.ReadHost, survey.WithValidator(survey.Required))
		if err != nil {
			return err
		}
	}

	// Port
	var defaultPort string
	switch config.Type {
	case "postgresql":
		defaultPort = "5432"
	case "mysql":
		defaultPort = "3306"
	case "mongodb":
		defaultPort = "27017"
	}

	portPrompt := &survey.Input{
		Message: "Database port:",
		Default: defaultPort,
	}
	err = survey.AskOne(portPrompt, &config.Port, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	// SSL Mode for PostgreSQL/MySQL
	if config.Type == "postgresql" || config.Type == "mysql" {
		var sslMode string
		sslPrompt := &survey.Select{
			Message: "SSL Mode:",
			Options: []string{"disable", "require", "verify-ca", "verify-full"},
			Default: "disable",
		}
		err = survey.AskOne(sslPrompt, &sslMode)
		if err != nil {
			return err
		}
		config.SSLMode = sslMode
	}

	return nil
}

func getClusterConfig(config *generator.DatabaseConfig) error {
	// Number of cluster nodes
	var nodeCountStr string
	nodeCountPrompt := &survey.Input{
		Message: "Number of cluster nodes:",
		Default: "3",
		Help:    "Enter the number of database nodes in your cluster",
	}
	err := survey.AskOne(nodeCountPrompt, &nodeCountStr, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	// Get cluster node addresses
	config.ClusterNodes = make([]string, 0)
	for i := 1; i <= 3; i++ { // Default to 3 nodes, can be made dynamic
		var nodeHost string
		nodePrompt := &survey.Input{
			Message: fmt.Sprintf("Cluster node %d host:", i),
			Default: fmt.Sprintf("node%d.cluster.local", i),
		}
		err = survey.AskOne(nodePrompt, &nodeHost, survey.WithValidator(survey.Required))
		if err != nil {
			return err
		}
		config.ClusterNodes = append(config.ClusterNodes, nodeHost)
	}

	// Port
	var defaultPort string
	switch config.Type {
	case "postgresql":
		defaultPort = "5432"
	case "mysql":
		defaultPort = "3306"
	case "mongodb":
		defaultPort = "27017"
	}

	portPrompt := &survey.Input{
		Message: "Database port:",
		Default: defaultPort,
	}
	err = survey.AskOne(portPrompt, &config.Port, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	// MongoDB replica set
	if config.Type == "mongodb" {
		replicaSetPrompt := &survey.Input{
			Message: "Replica set name:",
			Default: "rs0",
		}
		survey.AskOne(replicaSetPrompt, &config.ReplicaSet)
	}

	return nil
}

func getFrameworkConfiguration() (string, error) {
	var framework string
	frameworkPrompt := &survey.Select{
		Message: "Which web framework would you like to use for your API?",
		Options: []string{
			"gin - Fast HTTP web framework with a martini-like API",
			"echo - High performance, extensible, minimalist Go web framework",
			"gorilla - A web toolkit for the Go programming language",
			"Quit",
		},
		Help: "Choose the web framework that best fits your project needs",
	}

	err := survey.AskOne(frameworkPrompt, &framework)
	if err != nil {
		if isUserInterrupt(err) {
			return "", GetProcessManager().HandleGracefulShutdown()
		}
		return "", fmt.Errorf("framework selection failed: %w", err)
	}

	// Handle quit option
	if framework == "Quit" {
		return "", GetProcessManager().HandleGracefulShutdown()
	}

	// Extract framework type from selection
	switch {
	case framework[:3] == "gin":
		return "gin", nil
	case framework[:4] == "echo":
		return "echo", nil
	case framework[:7] == "gorilla":
		return "gorilla", nil
	default:
		return "gin", nil // Default fallback
	}
}
