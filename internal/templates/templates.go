package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"text/template"
)

//go:embed api webapp microservice cli
var templateFS embed.FS

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

type TemplateData struct {
	ProjectName    string
	ModuleName     string
	DatabaseConfig DatabaseConfig
	RedisConfig    RedisConfig
	GeneratedAt    string
	GophexVersion  string
	Checksums      map[string]string
}

type FileTemplate struct {
	Path    string
	Content string
}

func GetTemplateFiles(templateType string) ([]FileTemplate, error) {
	var files []FileTemplate

	templateDir := templateType
	err := fs.WalkDir(templateFS, templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Read the template file
		content, err := templateFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", path, err)
		}

		// Remove the template type prefix and .tmpl suffix from path
		relativePath := strings.TrimPrefix(path, templateType+"/")
		relativePath = strings.TrimSuffix(relativePath, ".tmpl")

		// Handle special cases for hidden files
		if relativePath == "env.example" {
			relativePath = ".env.example"
		}
		if relativePath == "env" {
			relativePath = ".env"
		}

		files = append(files, FileTemplate{
			Path:    relativePath,
			Content: string(content),
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk template directory %s: %w", templateType, err)
	}

	return files, nil
}

func ProcessTemplate(content string, data TemplateData) (string, error) {
	tmpl, err := template.New("template").Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func GenerateModuleName(projectName string) string {
	// Use a simple local module name that will work out of the box
	// Users can change this later if they want to publish to a specific repository
	return strings.ToLower(projectName)
}
