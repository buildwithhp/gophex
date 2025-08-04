package project

import (
	"errors"
	"fmt"
)

// Domain errors for project operations
var (
	ErrInvalidProjectName = errors.New("project name cannot be empty")
	ErrInvalidProjectType = errors.New("invalid project type")
	ErrInvalidProjectPath = errors.New("project path cannot be empty")
	ErrProjectNotFound    = errors.New("project not found")
	ErrProjectExists      = errors.New("project already exists")
	ErrInvalidConfig      = errors.New("invalid configuration")
	ErrActivityNotFound   = errors.New("activity not found")
	ErrFeatureNotFound    = errors.New("feature not found")
)

// ValidationError represents a validation error with details
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

// Error implements the error interface
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field string, value interface{}, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// ConfigurationError represents a configuration error
type ConfigurationError struct {
	Component string
	Reason    string
	Err       error
}

// Error implements the error interface
func (e ConfigurationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("configuration error in %s: %s (%v)", e.Component, e.Reason, e.Err)
	}
	return fmt.Sprintf("configuration error in %s: %s", e.Component, e.Reason)
}

// Unwrap returns the underlying error
func (e ConfigurationError) Unwrap() error {
	return e.Err
}

// NewConfigurationError creates a new configuration error
func NewConfigurationError(component, reason string, err error) ConfigurationError {
	return ConfigurationError{
		Component: component,
		Reason:    reason,
		Err:       err,
	}
}

// GenerationError represents an error during project generation
type GenerationError struct {
	Stage   string
	Project string
	Err     error
}

// Error implements the error interface
func (e GenerationError) Error() string {
	return fmt.Sprintf("generation failed at stage '%s' for project '%s': %v", e.Stage, e.Project, e.Err)
}

// Unwrap returns the underlying error
func (e GenerationError) Unwrap() error {
	return e.Err
}

// NewGenerationError creates a new generation error
func NewGenerationError(stage, project string, err error) GenerationError {
	return GenerationError{
		Stage:   stage,
		Project: project,
		Err:     err,
	}
}
