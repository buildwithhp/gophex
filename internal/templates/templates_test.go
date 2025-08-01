package templates

import (
	"testing"
)

func TestGetTemplateFiles(t *testing.T) {
	files, err := GetTemplateFiles("api")
	if err != nil {
		t.Fatalf("Failed to get API template files: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No template files found for API")
	}

	t.Logf("Found %d template files:", len(files))
	for _, f := range files {
		t.Logf("- %s", f.Path)
	}
}

func TestProcessTemplate(t *testing.T) {
	content := "module {{.ModuleName}}\n\nproject: {{.ProjectName}}"
	data := TemplateData{
		ProjectName: "testapi",
		ModuleName:  "github.com/testapi",
	}

	result, err := ProcessTemplate(content, data)
	if err != nil {
		t.Fatalf("Failed to process template: %v", err)
	}

	expected := "module github.com/testapi\n\nproject: testapi"
	if result != expected {
		t.Fatalf("Template processing failed.\nExpected: %s\nGot: %s", expected, result)
	}
}
