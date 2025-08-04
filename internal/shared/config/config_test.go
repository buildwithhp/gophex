package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvironmentProvider(t *testing.T) {
	provider := NewEnvironmentProvider()

	// Test setting and getting environment variable
	key := "TEST_ENV_VAR"
	value := "test_value"

	err := provider.Set(key, value)
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	defer os.Unsetenv(key)

	retrievedValue, exists := provider.Get(key)
	if !exists {
		t.Error("Environment variable should exist")
	}

	if retrievedValue != value {
		t.Errorf("Expected %s, got %s", value, retrievedValue)
	}

	// Test non-existent variable
	_, exists = provider.Get("NON_EXISTENT_VAR")
	if exists {
		t.Error("Non-existent variable should not exist")
	}

	// Test Load and Save (should be no-ops)
	if err := provider.Load(); err != nil {
		t.Errorf("Load should not return error: %v", err)
	}

	if err := provider.Save(); err != nil {
		t.Errorf("Save should not return error: %v", err)
	}
}

func TestFileProvider(t *testing.T) {
	// Create temporary file
	tempDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configFile := filepath.Join(tempDir, "test.config")
	provider := NewFileProvider(configFile)

	// Test setting values
	err = provider.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Failed to set config value: %v", err)
	}

	err = provider.Set("key2", "value2")
	if err != nil {
		t.Fatalf("Failed to set config value: %v", err)
	}

	// Test getting values before save
	value, exists := provider.Get("key1")
	if !exists || value != "value1" {
		t.Error("Should be able to get value before save")
	}

	// Test save
	err = provider.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Create new provider and load
	newProvider := NewFileProvider(configFile)
	err = newProvider.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test getting values after load
	value, exists = newProvider.Get("key1")
	if !exists || value != "value1" {
		t.Error("Should be able to get value after load")
	}

	value, exists = newProvider.Get("key2")
	if !exists || value != "value2" {
		t.Error("Should be able to get value after load")
	}

	// Test non-existent key
	_, exists = newProvider.Get("non_existent")
	if exists {
		t.Error("Non-existent key should not exist")
	}
}

func TestFileProvider_LoadNonExistentFile(t *testing.T) {
	provider := NewFileProvider("/non/existent/file.config")

	// Should not return error for non-existent file
	err := provider.Load()
	if err != nil {
		t.Errorf("Load should not return error for non-existent file: %v", err)
	}
}

func TestDefaultProvider(t *testing.T) {
	defaults := map[string]string{
		"default_key1": "default_value1",
		"default_key2": "default_value2",
	}

	provider := NewDefaultProvider(defaults)

	// Test getting default values
	value, exists := provider.Get("default_key1")
	if !exists || value != "default_value1" {
		t.Error("Should get default value")
	}

	// Test non-existent key
	_, exists = provider.Get("non_existent")
	if exists {
		t.Error("Non-existent key should not exist")
	}

	// Test setting new default
	err := provider.Set("new_key", "new_value")
	if err != nil {
		t.Fatalf("Failed to set default value: %v", err)
	}

	value, exists = provider.Get("new_key")
	if !exists || value != "new_value" {
		t.Error("Should get newly set default value")
	}

	// Test Load and Save (should be no-ops)
	if err := provider.Load(); err != nil {
		t.Errorf("Load should not return error: %v", err)
	}

	if err := provider.Save(); err != nil {
		t.Errorf("Save should not return error: %v", err)
	}
}

func TestConfigManager(t *testing.T) {
	// Create providers
	defaults := map[string]string{
		"APP_NAME":  "test-app",
		"LOG_LEVEL": "info",
		"DEBUG":     "false",
	}
	defaultProvider := NewDefaultProvider(defaults)
	envProvider := NewEnvironmentProvider()

	// Set environment variable to override default
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Unsetenv("LOG_LEVEL")

	// Create manager with providers (env should override defaults)
	manager := NewManager(envProvider, defaultProvider)

	err := manager.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	config := manager.GetConfig()
	if config == nil {
		t.Fatal("Config should not be nil")
	}

	// Test that environment variable overrides default
	if config.LogLevel != "debug" {
		t.Errorf("Expected LOG_LEVEL to be 'debug', got '%s'", config.LogLevel)
	}

	// Test that default is used when env var is not set
	if config.AppName != "test-app" {
		t.Errorf("Expected APP_NAME to be 'test-app', got '%s'", config.AppName)
	}

	// Test boolean parsing
	if config.Debug != false {
		t.Errorf("Expected DEBUG to be false, got %v", config.Debug)
	}
}

func TestConfigManager_BooleanParsing(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{"true string", "true", true},
		{"false string", "false", false},
		{"1 string", "1", true},
		{"0 string", "0", false},
		{"invalid string", "invalid", false}, // should use default
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defaults := map[string]string{
				"TEST_BOOL": "false", // default to false
			}
			defaultProvider := NewDefaultProvider(defaults)
			envProvider := NewEnvironmentProvider()

			// Set environment variable
			os.Setenv("TEST_BOOL", test.envValue)
			defer os.Unsetenv("TEST_BOOL")

			manager := NewManager(envProvider, defaultProvider)
			err := manager.Load()
			if err != nil {
				t.Fatalf("Failed to load configuration: %v", err)
			}

			// Test getBool method indirectly through config loading
			result := manager.getBool("TEST_BOOL", false)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestConfigManager_MultipleProviders(t *testing.T) {
	// Create temporary file
	tempDir, err := os.MkdirTemp("", "config-multi-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configFile := filepath.Join(tempDir, "test.config")

	// Create file provider and set some values
	fileProvider := NewFileProvider(configFile)
	fileProvider.Set("FILE_KEY", "file_value")
	fileProvider.Set("OVERRIDE_KEY", "file_override")
	fileProvider.Save()

	// Create environment provider
	envProvider := NewEnvironmentProvider()
	os.Setenv("ENV_KEY", "env_value")
	os.Setenv("OVERRIDE_KEY", "env_override")
	defer func() {
		os.Unsetenv("ENV_KEY")
		os.Unsetenv("OVERRIDE_KEY")
	}()

	// Create defaults
	defaults := map[string]string{
		"DEFAULT_KEY":  "default_value",
		"OVERRIDE_KEY": "default_override",
	}
	defaultProvider := NewDefaultProvider(defaults)

	// Create manager with providers in order: env, file, defaults
	// Environment should have highest priority
	manager := NewManager(envProvider, fileProvider, defaultProvider)
	err = manager.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Test that each provider's values are accessible
	if value := manager.getString("DEFAULT_KEY", ""); value != "default_value" {
		t.Errorf("Expected default_value, got %s", value)
	}

	if value := manager.getString("FILE_KEY", ""); value != "file_value" {
		t.Errorf("Expected file_value, got %s", value)
	}

	if value := manager.getString("ENV_KEY", ""); value != "env_value" {
		t.Errorf("Expected env_value, got %s", value)
	}

	// Test that environment overrides file and defaults
	if value := manager.getString("OVERRIDE_KEY", ""); value != "env_override" {
		t.Errorf("Expected env_override, got %s", value)
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "Environment variable exists",
			envKey:       "TEST_EXISTING",
			envValue:     "existing_value",
			defaultValue: "default_value",
			expected:     "existing_value",
		},
		{
			name:         "Environment variable does not exist",
			envKey:       "TEST_NON_EXISTING",
			envValue:     "",
			defaultValue: "default_value",
			expected:     "default_value",
		},
		{
			name:         "Environment variable is empty",
			envKey:       "TEST_EMPTY",
			envValue:     "",
			defaultValue: "default_value",
			expected:     "default_value",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Clean up any existing env var
			os.Unsetenv(test.envKey)

			// Set env var if value is provided
			if test.envValue != "" {
				os.Setenv(test.envKey, test.envValue)
				defer os.Unsetenv(test.envKey)
			}

			result := GetEnvWithDefault(test.envKey, test.defaultValue)
			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestConfig_DefaultValues(t *testing.T) {
	manager := NewManager(NewEnvironmentProvider())
	err := manager.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	config := manager.GetConfig()

	// Test default values are set correctly
	if config.AppName != "gophex" {
		t.Errorf("Expected default AppName 'gophex', got '%s'", config.AppName)
	}

	if config.Version != "1.0.0" {
		t.Errorf("Expected default Version '1.0.0', got '%s'", config.Version)
	}

	if config.LogLevel != "info" {
		t.Errorf("Expected default LogLevel 'info', got '%s'", config.LogLevel)
	}

	if config.Debug != false {
		t.Errorf("Expected default Debug false, got %v", config.Debug)
	}

	if config.EnableCRUDGeneration != true {
		t.Errorf("Expected default EnableCRUDGeneration true, got %v", config.EnableCRUDGeneration)
	}
}

func BenchmarkConfigManager_Load(b *testing.B) {
	defaults := map[string]string{
		"APP_NAME":  "test-app",
		"LOG_LEVEL": "info",
		"DEBUG":     "false",
	}

	manager := NewManager(
		NewDefaultProvider(defaults),
		NewEnvironmentProvider(),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Load()
	}
}

func BenchmarkGetEnvWithDefault(b *testing.B) {
	os.Setenv("BENCH_TEST", "test_value")
	defer os.Unsetenv("BENCH_TEST")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetEnvWithDefault("BENCH_TEST", "default")
	}
}
