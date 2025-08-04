package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config represents the application configuration
type Config struct {
	// Application settings
	AppName  string
	Version  string
	LogLevel string
	Debug    bool

	// Template settings
	TemplateDir string
	OutputDir   string

	// Generation settings
	DefaultProjectType string
	DefaultModuleName  string

	// Feature flags
	EnableCRUDGeneration   bool
	EnableInteractiveMode  bool
	EnableMetadataTracking bool
}

// Provider defines the interface for configuration providers
type Provider interface {
	Get(key string) (string, bool)
	Set(key string, value string) error
	Load() error
	Save() error
}

// Manager manages application configuration
type Manager struct {
	providers []Provider
	config    *Config
}

// NewManager creates a new configuration manager
func NewManager(providers ...Provider) *Manager {
	return &Manager{
		providers: providers,
		config:    &Config{},
	}
}

// Load loads configuration from all providers
func (m *Manager) Load() error {
	// Load from all providers in order
	for _, provider := range m.providers {
		if err := provider.Load(); err != nil {
			return fmt.Errorf("failed to load from provider: %w", err)
		}
	}

	// Build final configuration
	m.config = &Config{
		AppName:                m.getString("APP_NAME", "gophex"),
		Version:                m.getString("VERSION", "1.0.0"),
		LogLevel:               m.getString("LOG_LEVEL", "info"),
		Debug:                  m.getBool("DEBUG", false),
		TemplateDir:            m.getString("TEMPLATE_DIR", "internal/templates"),
		OutputDir:              m.getString("OUTPUT_DIR", "."),
		DefaultProjectType:     m.getString("DEFAULT_PROJECT_TYPE", "api"),
		DefaultModuleName:      m.getString("DEFAULT_MODULE_NAME", "github.com/user/project"),
		EnableCRUDGeneration:   m.getBool("ENABLE_CRUD_GENERATION", true),
		EnableInteractiveMode:  m.getBool("ENABLE_INTERACTIVE_MODE", true),
		EnableMetadataTracking: m.getBool("ENABLE_METADATA_TRACKING", true),
	}

	return nil
}

// GetConfig returns the loaded configuration
func (m *Manager) GetConfig() *Config {
	return m.config
}

// getString gets a string value from providers with fallback
func (m *Manager) getString(key, defaultValue string) string {
	for _, provider := range m.providers {
		if value, exists := provider.Get(key); exists {
			return value
		}
	}
	return defaultValue
}

// getBool gets a boolean value from providers with fallback
func (m *Manager) getBool(key string, defaultValue bool) bool {
	for _, provider := range m.providers {
		if value, exists := provider.Get(key); exists {
			if parsed, err := strconv.ParseBool(value); err == nil {
				return parsed
			}
		}
	}
	return defaultValue
}

// getInt gets an integer value from providers with fallback
func (m *Manager) getInt(key string, defaultValue int) int {
	for _, provider := range m.providers {
		if value, exists := provider.Get(key); exists {
			if parsed, err := strconv.Atoi(value); err == nil {
				return parsed
			}
		}
	}
	return defaultValue
}

// EnvironmentProvider provides configuration from environment variables
type EnvironmentProvider struct{}

// NewEnvironmentProvider creates a new environment provider
func NewEnvironmentProvider() Provider {
	return &EnvironmentProvider{}
}

// Get gets a value from environment variables
func (e *EnvironmentProvider) Get(key string) (string, bool) {
	value := os.Getenv(key)
	return value, value != ""
}

// Set sets an environment variable (not persistent)
func (e *EnvironmentProvider) Set(key string, value string) error {
	return os.Setenv(key, value)
}

// Load loads environment variables (no-op for environment provider)
func (e *EnvironmentProvider) Load() error {
	return nil
}

// Save saves environment variables (no-op for environment provider)
func (e *EnvironmentProvider) Save() error {
	return nil
}

// FileProvider provides configuration from a file
type FileProvider struct {
	filePath string
	data     map[string]string
}

// NewFileProvider creates a new file provider
func NewFileProvider(filePath string) Provider {
	return &FileProvider{
		filePath: filePath,
		data:     make(map[string]string),
	}
}

// Get gets a value from the file data
func (f *FileProvider) Get(key string) (string, bool) {
	value, exists := f.data[key]
	return value, exists
}

// Set sets a value in the file data
func (f *FileProvider) Set(key string, value string) error {
	f.data[key] = value
	return nil
}

// Load loads configuration from file
func (f *FileProvider) Load() error {
	content, err := os.ReadFile(f.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, that's okay
			return nil
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse simple key=value format
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			f.data[key] = value
		}
	}

	return nil
}

// Save saves configuration to file
func (f *FileProvider) Save() error {
	var lines []string
	for key, value := range f.data {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	content := strings.Join(lines, "\n")
	return os.WriteFile(f.filePath, []byte(content), 0644)
}

// DefaultProvider provides default configuration values
type DefaultProvider struct {
	defaults map[string]string
}

// NewDefaultProvider creates a new default provider
func NewDefaultProvider(defaults map[string]string) Provider {
	return &DefaultProvider{
		defaults: defaults,
	}
}

// Get gets a value from defaults
func (d *DefaultProvider) Get(key string) (string, bool) {
	value, exists := d.defaults[key]
	return value, exists
}

// Set sets a default value
func (d *DefaultProvider) Set(key string, value string) error {
	d.defaults[key] = value
	return nil
}

// Load loads defaults (no-op)
func (d *DefaultProvider) Load() error {
	return nil
}

// Save saves defaults (no-op)
func (d *DefaultProvider) Save() error {
	return nil
}

// GetEnvWithDefault returns the value of an environment variable or a default value if not set
func GetEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
