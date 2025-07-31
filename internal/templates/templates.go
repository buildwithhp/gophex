package templates

import (
	"strings"
)

type File struct {
	Path    string
	Content string
}

type Template struct {
	Files []File
}

func GetTemplate(templateType string) Template {
	switch templateType {
	case "api":
		return getAPITemplate()
	case "webapp":
		return getWebAppTemplate()
	case "microservice":
		return getMicroserviceTemplate()
	case "cli":
		return getCLITemplate()
	default:
		return Template{}
	}
}

func ProcessTemplate(content, projectName string) string {
	content = strings.ReplaceAll(content, "{{.ProjectName}}", projectName)
	content = strings.ReplaceAll(content, "{{.ModuleName}}", "github.com/"+strings.ToLower(projectName))
	return content
}

func getAPITemplate() Template {
	return Template{
		Files: []File{
			{
				Path: "go.mod",
				Content: `module {{.ModuleName}}

go 1.21

require (
	github.com/gorilla/mux v1.8.0
)
`,
			},
			{
				Path: "cmd/api/main.go",
				Content: `package main

import (
	"log"
	"net/http"

	"{{.ModuleName}}/internal/api/routes"
)

func main() {
	router := routes.Setup()
	
	log.Println("API server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
`,
			},
			{
				Path: "internal/api/routes/routes.go",
				Content: `package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"{{.ModuleName}}/internal/api/handlers"
)

func Setup() *mux.Router {
	r := mux.NewRouter()
	
	r.HandleFunc("/health", handlers.Health).Methods("GET")
	r.HandleFunc("/api/v1/users", handlers.GetUsers).Methods("GET")
	
	return r
}
`,
			},
			{
				Path: "internal/api/handlers/health.go",
				Content: `package handlers

import (
	"encoding/json"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
`,
			},
			{
				Path: "internal/api/handlers/users.go",
				Content: `package handlers

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Jane Smith"},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
`,
			},
			{
				Path: "README.md",
				Content: `# {{.ProjectName}}

A REST API built with Go.

## Getting Started

1. Install dependencies:
   ` + "```bash" + `
   go mod tidy
   ` + "```" + `

2. Run the server:
   ` + "```bash" + `
   go run cmd/api/main.go
   ` + "```" + `

3. Test the API:
   ` + "```bash" + `
   curl http://localhost:8080/health
   ` + "```" + `
`,
			},
		},
	}
}

func getWebAppTemplate() Template {
	return Template{
		Files: []File{
			{
				Path: "go.mod",
				Content: `module {{.ModuleName}}

go 1.21

require (
	github.com/gorilla/mux v1.8.0
)
`,
			},
			{
				Path: "cmd/webapp/main.go",
				Content: `package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	
	log.Println("Web server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	data := struct {
		Title string
	}{
		Title: "{{.ProjectName}}",
	}
	tmpl.Execute(w, data)
}
`,
			},
			{
				Path: "web/templates/index.html",
				Content: `<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <h1>Welcome to {{.Title}}</h1>
    <p>Your Go web application is running!</p>
</body>
</html>
`,
			},
			{
				Path: "web/static/css/style.css",
				Content: `body {
    font-family: Arial, sans-serif;
    margin: 40px;
    background-color: #f4f4f4;
}

h1 {
    color: #333;
}
`,
			},
			{
				Path: "README.md",
				Content: `# {{.ProjectName}}

A web application built with Go.

## Getting Started

1. Install dependencies:
   ` + "```bash" + `
   go mod tidy
   ` + "```" + `

2. Run the server:
   ` + "```bash" + `
   go run cmd/webapp/main.go
   ` + "```" + `

3. Open your browser to http://localhost:8080
`,
			},
		},
	}
}

func getMicroserviceTemplate() Template {
	return Template{
		Files: []File{
			{
				Path: "go.mod",
				Content: `module {{.ModuleName}}

go 1.21

require (
	github.com/gorilla/mux v1.8.0
)
`,
			},
			{
				Path: "cmd/server/main.go",
				Content: `package main

import (
	"log"
	"net/http"

	"{{.ModuleName}}/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/health", handlers.Health).Methods("GET")
	r.HandleFunc("/api/{{.ProjectName}}", handlers.Service).Methods("GET")
	
	log.Printf("{{.ProjectName}} microservice starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
`,
			},
			{
				Path: "internal/handlers/handlers.go",
				Content: `package handlers

import (
	"encoding/json"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"service": "{{.ProjectName}}",
		"status":  "healthy",
	})
}

func Service(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "{{.ProjectName}} microservice is running",
	})
}
`,
			},
			{
				Path: "README.md",
				Content: `# {{.ProjectName}} Microservice

A microservice built with Go.

## Getting Started

1. Install dependencies:
   ` + "```bash" + `
   go mod tidy
   ` + "```" + `

2. Run the service:
   ` + "```bash" + `
   go run cmd/server/main.go
   ` + "```" + `

3. Test the service:
   ` + "```bash" + `
   curl http://localhost:8080/health
   ` + "```" + `
`,
			},
		},
	}
}

func getCLITemplate() Template {
	return Template{
		Files: []File{
			{
				Path: "go.mod",
				Content: `module {{.ModuleName}}

go 1.21

require (
	github.com/spf13/cobra v1.8.0
)
`,
			},
			{
				Path: "cmd/{{.ProjectName}}/main.go",
				Content: `package main

import (
	"fmt"
	"os"

	"{{.ModuleName}}/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
`,
			},
			{
				Path: "internal/cmd/root.go",
				Content: `package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "{{.ProjectName}}",
	Short: "A CLI tool built with Go",
	Long:  "{{.ProjectName}} is a command-line tool built with Go and Cobra.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from {{.ProjectName}}!")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
`,
			},
			{
				Path: "README.md",
				Content: `# {{.ProjectName}}

A CLI tool built with Go.

## Installation

` + "```bash" + `
go install {{.ModuleName}}/cmd/{{.ProjectName}}@latest
` + "```" + `

## Usage

` + "```bash" + `
{{.ProjectName}}
` + "```" + `
`,
			},
		},
	}
}
