package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hemant/gophex/internal/templates"
)

type Generator struct{}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(projectType, projectName, projectPath string) error {
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	switch projectType {
	case "api":
		return g.generateAPI(projectName, projectPath)
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

func (g *Generator) generateAPI(projectName, projectPath string) error {
	return g.createFromTemplate("api", projectName, projectPath)
}

func (g *Generator) generateWebApp(projectName, projectPath string) error {
	return g.createFromTemplate("webapp", projectName, projectPath)
}

func (g *Generator) generateMicroservice(projectName, projectPath string) error {
	return g.createFromTemplate("microservice", projectName, projectPath)
}

func (g *Generator) generateCLI(projectName, projectPath string) error {
	return g.createFromTemplate("cli", projectName, projectPath)
}

func (g *Generator) createFromTemplate(templateType, projectName, projectPath string) error {
	template := templates.GetTemplate(templateType)

	for _, file := range template.Files {
		filePath := filepath.Join(projectPath, file.Path)

		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
		}

		content := templates.ProcessTemplate(file.Content, projectName)

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}
