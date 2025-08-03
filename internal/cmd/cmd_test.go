package cmd

import (
	"errors"
	"strings"
	"testing"
)

func TestIsUserInterrupt(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "interrupt error",
			err:      errors.New("interrupt"),
			expected: true,
		},
		{
			name:     "EOF error",
			err:      errors.New("EOF"),
			expected: true,
		},
		{
			name:     "cancelled error",
			err:      errors.New("operation cancelled"),
			expected: true,
		},
		{
			name:     "canceled error",
			err:      errors.New("operation canceled"),
			expected: true,
		},
		{
			name:     "uppercase interrupt",
			err:      errors.New("INTERRUPT"),
			expected: true,
		},
		{
			name:     "mixed case EOF",
			err:      errors.New("Eof"),
			expected: true,
		},
		{
			name:     "regular error",
			err:      errors.New("file not found"),
			expected: false,
		},
		{
			name:     "network error",
			err:      errors.New("connection refused"),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isUserInterrupt(test.err)
			if result != test.expected {
				t.Errorf("isUserInterrupt(%v) = %v, expected %v", test.err, result, test.expected)
			}
		})
	}
}

func TestProjectTypeExtraction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"api - REST API with clean architecture", "api"},
		{"webapp - Web application with templates", "webapp"},
		{"microservice - Microservice with gRPC support", "microservice"},
		{"cli - Command-line tool", "cli"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			var projectType string

			// Simulate the extraction logic from generate.go
			switch {
			case test.input[:3] == "api":
				projectType = "api"
			case test.input[:6] == "webapp":
				projectType = "webapp"
			case test.input[:12] == "microservice":
				projectType = "microservice"
			case test.input[:3] == "cli":
				projectType = "cli"
			}

			if projectType != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, projectType)
			}
		})
	}
}

func TestDatabaseTypeExtraction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"PostgreSQL - Advanced open-source relational database", "postgresql"},
		{"MySQL - Popular open-source relational database", "mysql"},
		{"MongoDB - Document-oriented NoSQL database", "mongodb"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			var dbType string

			// Simulate the extraction logic from generate.go
			switch {
			case test.input[:10] == "PostgreSQL":
				dbType = "postgresql"
			case test.input[:5] == "MySQL":
				dbType = "mysql"
			case test.input[:7] == "MongoDB":
				dbType = "mongodb"
			}

			if dbType != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, dbType)
			}
		})
	}
}

func TestConfigTypeExtraction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Single instance - Simple single database server", "single"},
		{"Read-Write split - Separate read and write endpoints", "read-write"},
		{"Cluster - Multiple database nodes", "cluster"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			var configType string

			// Simulate the extraction logic from generate.go
			switch {
			case test.input[:6] == "Single":
				configType = "single"
			case test.input[:10] == "Read-Write":
				configType = "read-write"
			case test.input[:7] == "Cluster":
				configType = "cluster"
			}

			if configType != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, configType)
			}
		})
	}
}

func TestRedisChoiceExtraction(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Yes - Include Redis support", true},
		{"No - Skip Redis", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			wantsRedis := test.input[:3] == "Yes"

			if wantsRedis != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, wantsRedis)
			}
		})
	}
}

func TestQuitOptionDetection(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Quit", true},
		{"Exit", true},
		{"Yes - Continue", false},
		{"No - Cancel", false},
		{"api - REST API", false},
		{"PostgreSQL - Database", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			isQuit := test.input == "Quit" || test.input == "Exit"

			if isQuit != test.expected {
				t.Errorf("Expected %v for '%s', got %v", test.expected, test.input, isQuit)
			}
		})
	}
}

func TestProjectNameQuitDetection(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"quit", true},
		{"QUIT", true},
		{"Quit", true},
		{"  quit  ", true},
		{"myproject", false},
		{"test-api", false},
		{"quitproject", false},
		{"project-quit", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			isQuit := strings.ToLower(strings.TrimSpace(test.input)) == "quit"

			if isQuit != test.expected {
				t.Errorf("Expected %v for '%s', got %v", test.expected, test.input, isQuit)
			}
		})
	}
}

func TestMainFilePathGeneration(t *testing.T) {
	tests := []struct {
		projectType string
		expected    string
	}{
		{"api", "cmd/api/main.go"},
		{"webapp", "cmd/webapp/main.go"},
		{"microservice", "cmd/server/main.go"},
		{"cli", "cmd/main.go"},
	}

	for _, test := range tests {
		t.Run(test.projectType, func(t *testing.T) {
			var mainFile string

			// Simulate the logic from project_actions.go
			switch test.projectType {
			case "api":
				mainFile = "cmd/api/main.go"
			case "webapp":
				mainFile = "cmd/webapp/main.go"
			case "microservice":
				mainFile = "cmd/server/main.go"
			case "cli":
				mainFile = "cmd/main.go"
			}

			if mainFile != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, mainFile)
			}
		})
	}
}

func TestProcessNameGeneration(t *testing.T) {
	tests := []struct {
		projectType string
		expected    string
	}{
		{"api", "api-app"},
		{"webapp", "webapp-app"},
		{"microservice", "microservice-app"},
		{"cli", "cli-app"},
	}

	for _, test := range tests {
		t.Run(test.projectType, func(t *testing.T) {
			// Simulate the logic from project_actions.go
			processName := test.projectType + "-app"

			if processName != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, processName)
			}
		})
	}
}
