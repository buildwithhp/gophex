package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerator_Generate(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "gophex-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	gen := New()
	projectPath := filepath.Join(tempDir, "testproject")

	err = gen.Generate("cli", "testproject", projectPath)
	if err != nil {
		t.Fatalf("Failed to generate CLI project: %v", err)
	}

	// Check if project directory was created
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Fatal("Project directory was not created")
	}

	// Check if main.go was created
	mainFile := filepath.Join(projectPath, "cmd", "main.go")
	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		t.Fatal("Main file was not created")
	}
}

func TestGenerator_GenerateWithDatabaseConfig(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "gophex-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	gen := New()
	projectPath := filepath.Join(tempDir, "testapi")

	dbConfig := &DatabaseConfig{
		Type:         "postgresql",
		ConfigType:   "single",
		Host:         "localhost",
		Port:         "5432",
		Username:     "testuser",
		Password:     "testpass",
		DatabaseName: "testapi",
		SSLMode:      "disable",
	}

	err = gen.GenerateWithConfig("api", "testapi", projectPath, dbConfig)
	if err != nil {
		t.Fatalf("Failed to generate API project with database config: %v", err)
	}

	// Check if project directory was created
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Fatal("Project directory was not created")
	}

	// Check if API main file was created
	mainFile := filepath.Join(projectPath, "cmd", "api", "main.go")
	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		t.Fatal("API main file was not created")
	}

	// Check if database config file was created
	configFile := filepath.Join(projectPath, "internal", "config", "config.go")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}
}

func TestGenerator_GenerateWithFullConfig(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "gophex-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	gen := New()
	projectPath := filepath.Join(tempDir, "testapi")

	dbConfig := &DatabaseConfig{
		Type:         "postgresql",
		ConfigType:   "single",
		Host:         "localhost",
		Port:         "5432",
		Username:     "testuser",
		Password:     "testpass",
		DatabaseName: "testapi",
		SSLMode:      "disable",
	}

	redisConfig := &RedisConfig{
		Enabled:  true,
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		Database: 0,
	}

	err = gen.GenerateWithFullConfig("api", "testapi", projectPath, dbConfig, redisConfig)
	if err != nil {
		t.Fatalf("Failed to generate API project with full config: %v", err)
	}

	// Check if project directory was created
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Fatal("Project directory was not created")
	}

	// Check if Redis client file was created
	redisFile := filepath.Join(projectPath, "internal", "infrastructure", "database", "redis", "client.go")
	if _, err := os.Stat(redisFile); os.IsNotExist(err) {
		t.Fatal("Redis client file was not created")
	}

	// Read go.mod to check if Redis dependency is included
	goModFile := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModFile)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	goModContent := string(content)
	if !contains(goModContent, "github.com/go-redis/redis/v8") {
		t.Error("Redis dependency not found in go.mod")
	}
}

func TestGenerator_GenerateWithoutRedis(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "gophex-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	gen := New()
	projectPath := filepath.Join(tempDir, "testapi")

	dbConfig := &DatabaseConfig{
		Type:         "postgresql",
		ConfigType:   "single",
		Host:         "localhost",
		Port:         "5432",
		Username:     "testuser",
		Password:     "testpass",
		DatabaseName: "testapi",
		SSLMode:      "disable",
	}

	redisConfig := &RedisConfig{
		Enabled: false,
	}

	err = gen.GenerateWithFullConfig("api", "testapi", projectPath, dbConfig, redisConfig)
	if err != nil {
		t.Fatalf("Failed to generate API project without Redis: %v", err)
	}

	// Check if Redis client file was NOT created
	redisFile := filepath.Join(projectPath, "internal", "infrastructure", "database", "redis", "client.go")
	if _, err := os.Stat(redisFile); !os.IsNotExist(err) {
		t.Error("Redis client file should not be created when Redis is disabled")
	}

	// Read go.mod to check if Redis dependency is NOT included
	goModFile := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModFile)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	goModContent := string(content)
	if contains(goModContent, "github.com/go-redis/redis/v8") {
		t.Error("Redis dependency should not be in go.mod when Redis is disabled")
	}
}

func TestGenerator_UnsupportedProjectType(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gophex-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	gen := New()
	projectPath := filepath.Join(tempDir, "testproject")

	err = gen.Generate("unsupported", "testproject", projectPath)
	if err == nil {
		t.Fatal("Expected error for unsupported project type")
	}

	expectedError := "unsupported project type: unsupported"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDatabaseConfig_AllTypes(t *testing.T) {
	tests := []struct {
		dbType     string
		configType string
	}{
		{"postgresql", "single"},
		{"postgresql", "read-write"},
		{"postgresql", "cluster"},
		{"mysql", "single"},
		{"mysql", "read-write"},
		{"mysql", "cluster"},
		{"mongodb", "single"},
		{"mongodb", "read-write"},
		{"mongodb", "cluster"},
	}

	for _, test := range tests {
		t.Run(test.dbType+"_"+test.configType, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "gophex-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tempDir)

			gen := New()
			projectPath := filepath.Join(tempDir, "testapi")

			dbConfig := &DatabaseConfig{
				Type:         test.dbType,
				ConfigType:   test.configType,
				Host:         "localhost",
				Port:         "5432",
				Username:     "testuser",
				Password:     "testpass",
				DatabaseName: "testapi",
				SSLMode:      "disable",
			}

			if test.configType == "read-write" {
				dbConfig.ReadHost = "read.localhost"
				dbConfig.WriteHost = "write.localhost"
			}

			if test.configType == "cluster" {
				dbConfig.ClusterNodes = []string{"node1", "node2", "node3"}
			}

			err = gen.GenerateWithConfig("api", "testapi", projectPath, dbConfig)
			if err != nil {
				t.Fatalf("Failed to generate project with %s %s: %v", test.dbType, test.configType, err)
			}

			// Verify project was created
			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				t.Fatal("Project directory was not created")
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsAt(s, substr))))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
