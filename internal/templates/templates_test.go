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
		ModuleName:  "testapi",
	}

	result, err := ProcessTemplate(content, data)
	if err != nil {
		t.Fatalf("Failed to process template: %v", err)
	}

	expected := "module testapi\n\nproject: testapi"
	if result != expected {
		t.Fatalf("Template processing failed.\nExpected: %s\nGot: %s", expected, result)
	}
}

func TestGenerateModuleName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"MyProject", "myproject"},
		{"test-api", "test-api"},
		{"TestAPI", "testapi"},
		{"my_service", "my_service"},
	}

	for _, test := range tests {
		result := GenerateModuleName(test.input)
		if result != test.expected {
			t.Errorf("GenerateModuleName(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}
