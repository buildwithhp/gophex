package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/buildwithhp/gophex/internal/templates"
)

type DatabaseConfig struct {
	Type         string // mysql, postgresql, mongodb
	ConfigType   string // cluster, multi-cluster, read-write
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	ReadHost     string   // for read-write setup
	WriteHost    string   // for read-write setup
	ClusterNodes []string // for multi-cluster
	SSLMode      string
	AuthSource   string // for MongoDB
	ReplicaSet   string // for MongoDB
}

type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Password string
	Database int
}

type Generator struct{}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(projectType, projectName, projectPath string) error {
	return g.GenerateWithConfig(projectType, projectName, projectPath, nil)
}

func (g *Generator) GenerateWithConfig(projectType, projectName, projectPath string, dbConfig *DatabaseConfig) error {
	return g.GenerateWithFullConfig(projectType, projectName, projectPath, dbConfig, nil)
}

func (g *Generator) GenerateWithFullConfig(projectType, projectName, projectPath string, dbConfig *DatabaseConfig, redisConfig *RedisConfig) error {
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	switch projectType {
	case "api":
		return g.generateAPI(projectName, projectPath, dbConfig, redisConfig)
	case "webapp":
		return g.generateWebApp(projectName, projectPath)
	case "microservice":
		return g.generateMicroservice(projectName, projectPath)
	case "cli":
		return g.generateCLI(projectName, projectPath)
	default:
		return fmt.Errorf("unsupported project type: %s", projectType)
	}
}

func (g *Generator) generateAPI(projectName, projectPath string, dbConfig *DatabaseConfig, redisConfig *RedisConfig) error {
	return g.createFromTemplate("api", projectName, projectPath, dbConfig, redisConfig)
}

func (g *Generator) generateWebApp(projectName, projectPath string) error {
	return g.createFromTemplate("webapp", projectName, projectPath, nil, nil)
}

func (g *Generator) generateMicroservice(projectName, projectPath string) error {
	return g.createFromTemplate("microservice", projectName, projectPath, nil, nil)
}

func (g *Generator) generateCLI(projectName, projectPath string) error {
	return g.createFromTemplate("cli", projectName, projectPath, nil, nil)
}

func (g *Generator) createFromTemplate(templateType, projectName, projectPath string, dbConfig *DatabaseConfig, redisConfig *RedisConfig) error {
	// Get template files from embedded filesystem
	templateFiles, err := templates.GetTemplateFiles(templateType)
	if err != nil {
		return fmt.Errorf("failed to get template files for %s: %w", templateType, err)
	}

	// Prepare template data
	data := templates.TemplateData{
		ProjectName:   projectName,
		ModuleName:    templates.GenerateModuleName(projectName),
		GeneratedAt:   time.Now().Format(time.RFC3339),
		GophexVersion: "1.0.0", // TODO: Get from version package
		Checksums:     make(map[string]string),
	}

	// Add database configuration if provided
	if dbConfig != nil {
		data.DatabaseConfig = templates.DatabaseConfig{
			Type:         dbConfig.Type,
			ConfigType:   dbConfig.ConfigType,
			Host:         dbConfig.Host,
			Port:         dbConfig.Port,
			Username:     dbConfig.Username,
			Password:     dbConfig.Password,
			DatabaseName: dbConfig.DatabaseName,
			ReadHost:     dbConfig.ReadHost,
			WriteHost:    dbConfig.WriteHost,
			ClusterNodes: dbConfig.ClusterNodes,
			SSLMode:      dbConfig.SSLMode,
			AuthSource:   dbConfig.AuthSource,
			ReplicaSet:   dbConfig.ReplicaSet,
		}
	}

	// Add Redis configuration if provided
	if redisConfig != nil {
		data.RedisConfig = templates.RedisConfig{
			Enabled:  redisConfig.Enabled,
			Host:     redisConfig.Host,
			Port:     redisConfig.Port,
			Password: redisConfig.Password,
			Database: redisConfig.Database,
		}
	}

	for _, file := range templateFiles {
		// Skip Redis-related files if Redis is disabled
		if redisConfig != nil && !redisConfig.Enabled &&
			(strings.Contains(file.Path, "/redis/") || strings.Contains(file.Path, "redis")) {
			continue
		}

		filePath := filepath.Join(projectPath, file.Path)

		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
		}

		// Process template with proper template engine
		content, err := templates.ProcessTemplate(file.Content, data)
		if err != nil {
			return fmt.Errorf("failed to process template for %s: %w", file.Path, err)
		}

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}
