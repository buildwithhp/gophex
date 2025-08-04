package template

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

// Engine defines the interface for template processing
type Engine interface {
	// LoadTemplates loads templates from a filesystem
	LoadTemplates(fsys fs.FS, pattern string) error

	// Execute executes a template with the given data
	Execute(templateName string, data interface{}) (string, error)

	// ExecuteToFile executes a template and writes to a file
	ExecuteToFile(templateName string, data interface{}, outputPath string) error

	// AddFunction adds a custom function to the template engine
	AddFunction(name string, fn interface{}) error

	// ListTemplates returns a list of loaded template names
	ListTemplates() []string
}

// engine implements the Engine interface
type engine struct {
	templates *template.Template
	funcMap   template.FuncMap
}

// NewEngine creates a new template engine
func NewEngine() Engine {
	return &engine{
		templates: template.New("gophex"),
		funcMap:   getDefaultFuncMap(),
	}
}

// LoadTemplates loads templates from a filesystem
func (e *engine) LoadTemplates(fsys fs.FS, pattern string) error {
	// Reset templates with function map
	e.templates = template.New("gophex").Funcs(e.funcMap)

	// Walk through the filesystem and load templates
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Check if file matches pattern
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return err
		}

		if !matched {
			return nil
		}

		// Read template content
		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}

		// Parse template
		templateName := strings.TrimSuffix(path, ".tmpl")
		_, err = e.templates.New(templateName).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		return nil
	})

	return err
}

// Execute executes a template with the given data
func (e *engine) Execute(templateName string, data interface{}) (string, error) {
	tmpl := e.templates.Lookup(templateName)
	if tmpl == nil {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	var buf strings.Builder
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// ExecuteToFile executes a template and writes to a file
func (e *engine) ExecuteToFile(templateName string, data interface{}, outputPath string) error {
	content, err := e.Execute(templateName, data)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := createDirIfNotExists(dir); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write to file
	return writeFile(outputPath, []byte(content))
}

// AddFunction adds a custom function to the template engine
func (e *engine) AddFunction(name string, fn interface{}) error {
	if e.funcMap == nil {
		e.funcMap = make(template.FuncMap)
	}

	e.funcMap[name] = fn

	// Recreate templates with updated function map
	if e.templates != nil {
		e.templates = e.templates.Funcs(e.funcMap)
	}

	return nil
}

// ListTemplates returns a list of loaded template names
func (e *engine) ListTemplates() []string {
	if e.templates == nil {
		return nil
	}

	var names []string
	for _, tmpl := range e.templates.Templates() {
		if tmpl.Name() != "gophex" { // Skip root template
			names = append(names, tmpl.Name())
		}
	}

	return names
}

// getDefaultFuncMap returns the default function map for templates
func getDefaultFuncMap() template.FuncMap {
	return template.FuncMap{
		// String functions
		"title":     strings.Title,
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
		"trim":      strings.TrimSpace,
		"replace":   strings.ReplaceAll,
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,

		// Utility functions
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},

		// Collection functions
		"len": func(v interface{}) int {
			switch val := v.(type) {
			case []interface{}:
				return len(val)
			case []string:
				return len(val)
			case string:
				return len(val)
			default:
				return 0
			}
		},

		// Conditional functions
		"eq":  func(a, b interface{}) bool { return a == b },
		"ne":  func(a, b interface{}) bool { return a != b },
		"not": func(a bool) bool { return !a },
		"and": func(a, b bool) bool { return a && b },
		"or":  func(a, b bool) bool { return a || b },

		// Type checking functions
		"isString": func(v interface{}) bool {
			_, ok := v.(string)
			return ok
		},
		"isInt": func(v interface{}) bool {
			switch v.(type) {
			case int, int8, int16, int32, int64:
				return true
			default:
				return false
			}
		},
		"isBool": func(v interface{}) bool {
			_, ok := v.(bool)
			return ok
		},

		// Custom Gophex functions
		"pluralize":  pluralize,
		"camelCase":  camelCase,
		"snakeCase":  snakeCase,
		"kebabCase":  kebabCase,
		"pascalCase": pascalCase,
	}
}

// Helper functions for template processing

// pluralize converts a singular word to plural
func pluralize(word string) string {
	if strings.HasSuffix(word, "y") && len(word) > 1 {
		beforeY := word[len(word)-2]
		if beforeY != 'a' && beforeY != 'e' && beforeY != 'i' && beforeY != 'o' && beforeY != 'u' {
			return strings.TrimSuffix(word, "y") + "ies"
		}
	}

	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") ||
		strings.HasSuffix(word, "z") || strings.HasSuffix(word, "sh") ||
		strings.HasSuffix(word, "ch") {
		return word + "es"
	}

	if strings.HasSuffix(word, "f") {
		return strings.TrimSuffix(word, "f") + "ves"
	}

	if strings.HasSuffix(word, "fe") {
		return strings.TrimSuffix(word, "fe") + "ves"
	}

	return word + "s"
}

// camelCase converts a string to camelCase
func camelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	if len(words) == 0 {
		return s
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		result += strings.Title(strings.ToLower(words[i]))
	}

	return result
}

// snakeCase converts a string to snake_case
func snakeCase(s string) string {
	var result strings.Builder

	for i, r := range s {
		if i > 0 && (r >= 'A' && r <= 'Z') {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}

	return strings.ToLower(result.String())
}

// kebabCase converts a string to kebab-case
func kebabCase(s string) string {
	return strings.ReplaceAll(snakeCase(s), "_", "-")
}

// pascalCase converts a string to PascalCase
func pascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	var result strings.Builder
	for _, word := range words {
		result.WriteString(strings.Title(strings.ToLower(word)))
	}

	return result.String()
}

// File system helper functions (these would be implemented based on your needs)
func createDirIfNotExists(dir string) error {
	// Implementation would go here
	return nil
}

func writeFile(path string, content []byte) error {
	// Implementation would go here
	return nil
}
