package project

import (
	"context"
	"testing"
	"time"
)

// MockRepository implements the Repository interface for testing
type MockRepository struct {
	projects     map[string]*Project
	activities   map[string]map[string]Activity
	existsResult bool // For controlling Exists method behavior in tests
}

func (m *MockRepository) Save(ctx context.Context, project *Project) error {
	m.projects[project.Path] = project
	return nil
}

func (m *MockRepository) FindByPath(ctx context.Context, path string) (*Project, error) {
	if project, exists := m.projects[path]; exists {
		return project, nil
	}
	return nil, ErrProjectNotFound
}

func (m *MockRepository) FindByName(ctx context.Context, name string) (*Project, error) {
	for _, project := range m.projects {
		if project.Name == name {
			return project, nil
		}
	}
	return nil, ErrProjectNotFound
}

func TestService_CreateProject_FrameworkValidation(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	metadataRepo := NewMockMetadataRepository()
	generator := NewMockGenerator()
	logger := NewMockLogger()

	service := NewService(repo, metadataRepo, generator, logger)

	// Test case 1: API project without framework should fail
	req := CreateProjectRequest{
		Name: "test-api",
		Type: ProjectTypeAPI,
		Path: "/tmp/test-api",
		// Framework is empty - should fail for API projects
	}

	ctx := context.Background()
	project, err := service.CreateProject(ctx, req)

	if err == nil {
		t.Fatal("Expected validation error for API project without framework")
	}

	if project != nil {
		t.Error("Expected no project to be created")
	}

	var validationErr ValidationError
	if !isValidationError(err, &validationErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}

	// Test case 2: API project with valid framework should succeed
	req.Framework = FrameworkTypeGin
	repo.existsResult = false // Project doesn't exist

	project, err = service.CreateProject(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if project == nil {
		t.Fatal("Expected project to be created")
	}

	if project.Framework != FrameworkTypeGin {
		t.Errorf("Expected framework to be %s, got %s", FrameworkTypeGin, project.Framework)
	}

	// Test case 3: Non-API project with framework should succeed (framework is ignored)
	req.Type = ProjectTypeCLI
	req.Framework = FrameworkTypeEcho // Should be ignored for CLI projects

	project, err = service.CreateProject(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if project == nil {
		t.Fatal("Expected project to be created")
	}
	return nil, ErrProjectNotFound
}

func (m *MockRepository) List(ctx context.Context) ([]*Project, error) {
	var projects []*Project
	for _, project := range m.projects {
		projects = append(projects, project)
	}
	return projects, nil
}

func (m *MockRepository) Delete(ctx context.Context, path string) error {
	delete(m.projects, path)
	return nil
}

func (m *MockRepository) Exists(ctx context.Context, path string) (bool, error) {
	return m.existsResult, nil
}

func (m *MockRepository) UpdateActivity(ctx context.Context, projectPath, activityName string, completed bool) error {
	if m.activities[projectPath] == nil {
		m.activities[projectPath] = make(map[string]Activity)
	}

	activity := m.activities[projectPath][activityName]
	activity.Name = activityName
	activity.Completed = completed
	if completed {
		now := time.Now()
		activity.Timestamp = &now
	}
	m.activities[projectPath][activityName] = activity

	// Update project if it exists
	if project, exists := m.projects[projectPath]; exists {
		if project.Activities == nil {
			project.Activities = make(map[string]Activity)
		}
		project.Activities[activityName] = activity
	}

	return nil
}

func (m *MockRepository) UpdateMetadata(ctx context.Context, project *Project) error {
	if _, exists := m.projects[project.Path]; exists {
		m.projects[project.Path] = project
		return nil
	}
	return ErrProjectNotFound
}

// MockMetadataRepository implements the MetadataRepository interface for testing
type MockMetadataRepository struct {
	metadata map[string]*ProjectMetadata
}

func NewMockMetadataRepository() *MockMetadataRepository {
	return &MockMetadataRepository{
		metadata: make(map[string]*ProjectMetadata),
	}
}

func (m *MockMetadataRepository) LoadMetadata(ctx context.Context, projectPath string) (*ProjectMetadata, error) {
	if metadata, exists := m.metadata[projectPath]; exists {
		return metadata, nil
	}
	return nil, ErrProjectNotFound
}

func (m *MockMetadataRepository) SaveMetadata(ctx context.Context, projectPath string, metadata *ProjectMetadata) error {
	m.metadata[projectPath] = metadata
	return nil
}

func (m *MockMetadataRepository) HasMetadata(ctx context.Context, projectPath string) (bool, error) {
	_, exists := m.metadata[projectPath]
	return exists, nil
}

// MockGenerator implements the Generator interface for testing
type MockGenerator struct {
	generateCalled bool
	generateError  error
}

func NewMockGenerator() *MockGenerator {
	return &MockGenerator{}
}

func (m *MockGenerator) Generate(ctx context.Context, project *Project) error {
	m.generateCalled = true
	return m.generateError
}

// MockLogger implements the Logger interface for testing
type MockLogger struct {
	logs []LogEntry
}

type LogEntry struct {
	Level   string
	Message string
	Fields  []interface{}
	Error   error
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		logs: make([]LogEntry, 0),
	}
}

func (m *MockLogger) Info(msg string, fields ...interface{}) {
	m.logs = append(m.logs, LogEntry{Level: "INFO", Message: msg, Fields: fields})
}

func (m *MockLogger) Error(msg string, err error, fields ...interface{}) {
	m.logs = append(m.logs, LogEntry{Level: "ERROR", Message: msg, Fields: fields, Error: err})
}

func (m *MockLogger) Debug(msg string, fields ...interface{}) {
	m.logs = append(m.logs, LogEntry{Level: "DEBUG", Message: msg, Fields: fields})
}

func (m *MockLogger) Warn(msg string, fields ...interface{}) {
	m.logs = append(m.logs, LogEntry{Level: "WARN", Message: msg, Fields: fields})
}

// Test functions

func TestService_CreateProject(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	metadataRepo := NewMockMetadataRepository()
	generator := NewMockGenerator()
	logger := NewMockLogger()

	service := NewService(repo, metadataRepo, generator, logger)

	// Test data
	req := CreateProjectRequest{
		Name: "test-project",
		Type: ProjectTypeAPI,
		Path: "/tmp/test-project",
		DatabaseConfig: &DatabaseConfig{
			Type: DatabaseTypePostgreSQL,
		},
	}

	// Execute
	ctx := context.Background()
	project, err := service.CreateProject(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if project == nil {
		t.Fatal("Expected project to be created")
	}

	if project.Name != req.Name {
		t.Errorf("Expected project name %s, got %s", req.Name, project.Name)
	}

	if project.Type != req.Type {
		t.Errorf("Expected project type %s, got %s", req.Type, project.Type)
	}

	if !generator.generateCalled {
		t.Error("Expected generator to be called")
	}

	// Check if project was saved
	savedProject, err := repo.FindByPath(ctx, req.Path)
	if err != nil {
		t.Fatalf("Expected project to be saved, got error: %v", err)
	}

	if savedProject.Name != req.Name {
		t.Errorf("Expected saved project name %s, got %s", req.Name, savedProject.Name)
	}
}

func TestService_CreateProject_ValidationError(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	metadataRepo := NewMockMetadataRepository()
	generator := NewMockGenerator()
	logger := NewMockLogger()

	service := NewService(repo, metadataRepo, generator, logger)

	// Test data with invalid name
	req := CreateProjectRequest{
		Name: "", // Invalid empty name
		Type: ProjectTypeAPI,
		Path: "/tmp/test-project",
	}

	// Execute
	ctx := context.Background()
	project, err := service.CreateProject(ctx, req)

	// Assert
	if err == nil {
		t.Fatal("Expected validation error")
	}

	if project != nil {
		t.Error("Expected no project to be created")
	}

	var validationErr ValidationError
	if !isValidationError(err, &validationErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

func TestService_CreateProject_ProjectExists(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	metadataRepo := NewMockMetadataRepository()
	generator := NewMockGenerator()
	logger := NewMockLogger()

	service := NewService(repo, metadataRepo, generator, logger)

	// Pre-create a project
	existingProject := &Project{
		Name: "existing-project",
		Type: ProjectTypeAPI,
		Path: "/tmp/existing-project",
	}
	repo.Save(context.Background(), existingProject)

	// Test data with same path
	req := CreateProjectRequest{
		Name: "new-project",
		Type: ProjectTypeAPI,
		Path: "/tmp/existing-project", // Same path as existing project
	}

	// Execute
	ctx := context.Background()
	project, err := service.CreateProject(ctx, req)

	// Assert
	if err != ErrProjectExists {
		t.Fatalf("Expected ErrProjectExists, got %v", err)
	}

	if project != nil {
		t.Error("Expected no project to be created")
	}
}

func TestService_LoadProject(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	metadataRepo := NewMockMetadataRepository()
	generator := NewMockGenerator()
	logger := NewMockLogger()

	service := NewService(repo, metadataRepo, generator, logger)

	// Pre-create a project
	existingProject := &Project{
		Name: "existing-project",
		Type: ProjectTypeAPI,
		Path: "/tmp/existing-project",
	}
	repo.Save(context.Background(), existingProject)

	// Execute
	ctx := context.Background()
	project, err := service.LoadProject(ctx, "/tmp/existing-project")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if project == nil {
		t.Fatal("Expected project to be loaded")
	}

	if project.Name != existingProject.Name {
		t.Errorf("Expected project name %s, got %s", existingProject.Name, project.Name)
	}
}

func TestService_ValidateProjectName(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	metadataRepo := NewMockMetadataRepository()
	generator := NewMockGenerator()
	logger := NewMockLogger()

	service := NewService(repo, metadataRepo, generator, logger)

	tests := []struct {
		name        string
		projectName string
		expectError bool
	}{
		{"valid name", "my-project", false},
		{"valid name with numbers", "project123", false},
		{"empty name", "", true},
		{"too short", "a", true},
		{"too long", "this-is-a-very-long-project-name-that-exceeds-the-maximum-allowed-length", true},
		{"valid name with spaces", "My Project", false},
		{"valid name with underscores", "my_project", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.ValidateProjectName(test.projectName)

			if test.expectError && err == nil {
				t.Error("Expected validation error")
			}

			if !test.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestService_CompleteActivity(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	metadataRepo := NewMockMetadataRepository()
	generator := NewMockGenerator()
	logger := NewMockLogger()

	service := NewService(repo, metadataRepo, generator, logger)

	// Pre-create a project
	projectPath := "/tmp/test-project"
	existingProject := &Project{
		Name:       "test-project",
		Type:       ProjectTypeAPI,
		Path:       projectPath,
		Activities: make(map[string]Activity),
	}
	repo.Save(context.Background(), existingProject)

	// Execute
	ctx := context.Background()
	err := service.CompleteActivity(ctx, projectPath, "test_activity")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if activity was completed
	completed, err := service.GetActivityStatus(ctx, projectPath, "test_activity")
	if err != nil {
		t.Fatalf("Expected no error getting activity status, got %v", err)
	}

	if !completed {
		t.Error("Expected activity to be completed")
	}
}

// Helper functions

func isValidationError(err error, target *ValidationError) bool {
	if ve, ok := err.(ValidationError); ok {
		*target = ve
		return true
	}
	return false
}
