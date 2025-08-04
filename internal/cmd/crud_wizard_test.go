package cmd

import (
	"testing"
)

func TestIsValidEntityName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid lowercase", "user", true},
		{"valid with numbers", "user123", true},
		{"invalid uppercase", "User", false},
		{"invalid with underscore", "user_name", false},
		{"invalid with dash", "user-name", false},
		{"invalid starting with number", "123user", false},
		{"invalid empty", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isValidEntityName(test.input)
			if result != test.expected {
				t.Errorf("isValidEntityName(%q) = %v, expected %v", test.input, result, test.expected)
			}
		})
	}
}

func TestIsValidFieldName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid camelCase", "userName", true},
		{"valid PascalCase", "UserName", true},
		{"valid single letter", "a", true},
		{"valid with numbers", "user123", true},
		{"invalid with underscore", "user_name", false},
		{"invalid with dash", "user-name", false},
		{"invalid starting with number", "123user", false},
		{"invalid empty", "", false},
		{"invalid with space", "user name", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isValidFieldName(test.input)
			if result != test.expected {
				t.Errorf("isValidFieldName(%q) = %v, expected %v", test.input, result, test.expected)
			}
		})
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		singular string
		expected string
	}{
		{"user", "users"},
		{"post", "posts"},
		{"category", "categories"},
		{"box", "boxes"},
		{"dish", "dishes"},
		{"church", "churches"},
		{"book", "books"},
		{"company", "companies"},
	}

	for _, test := range tests {
		t.Run(test.singular, func(t *testing.T) {
			result := pluralize(test.singular)
			if result != test.expected {
				t.Errorf("pluralize(%q) = %q, expected %q", test.singular, result, test.expected)
			}
		})
	}
}

func TestGetCommonFields(t *testing.T) {
	tests := []struct {
		entityName    string
		expectedCount int
		hasCreatedAt  bool
		hasUpdatedAt  bool
	}{
		{"user", 5, true, true},
		{"post", 6, true, true},
		{"product", 7, true, true},
		{"task", 7, true, true},
		{"unknown", 2, true, true}, // Default fields
	}

	for _, test := range tests {
		t.Run(test.entityName, func(t *testing.T) {
			fields := getCommonFields(test.entityName)

			if len(fields) != test.expectedCount {
				t.Errorf("getCommonFields(%q) returned %d fields, expected %d",
					test.entityName, len(fields), test.expectedCount)
			}

			// Check for CreatedAt and UpdatedAt fields
			hasCreatedAt := false
			hasUpdatedAt := false

			for _, field := range fields {
				if field.Name == "CreatedAt" {
					hasCreatedAt = true
				}
				if field.Name == "UpdatedAt" {
					hasUpdatedAt = true
				}
			}

			if hasCreatedAt != test.hasCreatedAt {
				t.Errorf("getCommonFields(%q) CreatedAt field presence = %v, expected %v",
					test.entityName, hasCreatedAt, test.hasCreatedAt)
			}

			if hasUpdatedAt != test.hasUpdatedAt {
				t.Errorf("getCommonFields(%q) UpdatedAt field presence = %v, expected %v",
					test.entityName, hasUpdatedAt, test.hasUpdatedAt)
			}
		})
	}
}

func TestCRUDFieldValidation(t *testing.T) {
	field := CRUDField{
		Name:     "Email",
		Type:     "string",
		JSONTag:  "email",
		DBTag:    "email",
		Required: true,
		Unique:   true,
	}

	// Test field properties
	if field.Name != "Email" {
		t.Errorf("Expected field name 'Email', got %q", field.Name)
	}

	if field.Type != "string" {
		t.Errorf("Expected field type 'string', got %q", field.Type)
	}

	if !field.Required {
		t.Error("Expected field to be required")
	}

	if !field.Unique {
		t.Error("Expected field to be unique")
	}
}

func TestCRUDEntityValidation(t *testing.T) {
	entity := &CRUDEntity{
		Name:         "user",
		PluralName:   "users",
		UpdateMethod: "both",
		Fields: []CRUDField{
			{Name: "Name", Type: "string", Required: true},
			{Name: "Email", Type: "string", Required: true, Unique: true},
		},
	}

	// Test entity properties
	if entity.Name != "user" {
		t.Errorf("Expected entity name 'user', got %q", entity.Name)
	}

	if entity.PluralName != "users" {
		t.Errorf("Expected plural name 'users', got %q", entity.PluralName)
	}

	if entity.UpdateMethod != "both" {
		t.Errorf("Expected update method 'both', got %q", entity.UpdateMethod)
	}

	if len(entity.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(entity.Fields))
	}

	// Test required field
	nameField := entity.Fields[0]
	if !nameField.Required {
		t.Error("Expected Name field to be required")
	}

	// Test unique field
	emailField := entity.Fields[1]
	if !emailField.Unique {
		t.Error("Expected Email field to be unique")
	}
}
