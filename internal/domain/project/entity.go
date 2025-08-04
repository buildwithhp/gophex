package project

import (
	"time"
)

// Project represents a Gophex project entity
type Project struct {
	Name           string
	Type           ProjectType
	Path           string
	ModuleName     string
	GeneratedAt    time.Time
	LastUpdated    time.Time
	DatabaseConfig *DatabaseConfig
	RedisConfig    *RedisConfig
	Activities     map[string]Activity
	Features       []Feature
}

// ProjectType represents the type of project
type ProjectType string

const (
	ProjectTypeAPI          ProjectType = "api"
	ProjectTypeWebApp       ProjectType = "webapp"
	ProjectTypeMicroservice ProjectType = "microservice"
	ProjectTypeCLI          ProjectType = "cli"
)

// IsValid checks if the project type is valid
func (pt ProjectType) IsValid() bool {
	switch pt {
	case ProjectTypeAPI, ProjectTypeWebApp, ProjectTypeMicroservice, ProjectTypeCLI:
		return true
	default:
		return false
	}
}

// String returns the string representation of the project type
func (pt ProjectType) String() string {
	return string(pt)
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Type         DatabaseType
	ConfigType   DatabaseConfigType
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	ReadHost     string
	WriteHost    string
	ClusterNodes []string
	SSLMode      string
	AuthSource   string
	ReplicaSet   string
}

// DatabaseType represents the type of database
type DatabaseType string

const (
	DatabaseTypeMySQL      DatabaseType = "mysql"
	DatabaseTypePostgreSQL DatabaseType = "postgresql"
	DatabaseTypeMongoDB    DatabaseType = "mongodb"
)

// DatabaseConfigType represents the database configuration type
type DatabaseConfigType string

const (
	DatabaseConfigTypeSingle    DatabaseConfigType = "single"
	DatabaseConfigTypeCluster   DatabaseConfigType = "cluster"
	DatabaseConfigTypeReadWrite DatabaseConfigType = "read-write"
)

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Password string
	Database int
}

// Activity represents a project activity
type Activity struct {
	Name      string
	Completed bool
	Timestamp *time.Time
	CanRepeat bool
}

// Feature represents a project feature
type Feature struct {
	Name        string
	Enabled     bool
	Description string
}

// Validate validates the project entity
func (p *Project) Validate() error {
	if p.Name == "" {
		return ErrInvalidProjectName
	}

	if !p.Type.IsValid() {
		return ErrInvalidProjectType
	}

	if p.Path == "" {
		return ErrInvalidProjectPath
	}

	return nil
}

// HasFeature checks if the project has a specific feature
func (p *Project) HasFeature(featureName string) bool {
	for _, feature := range p.Features {
		if feature.Name == featureName && feature.Enabled {
			return true
		}
	}
	return false
}

// IsActivityCompleted checks if an activity is completed
func (p *Project) IsActivityCompleted(activityName string) bool {
	if activity, exists := p.Activities[activityName]; exists {
		return activity.Completed
	}
	return false
}

// CompleteActivity marks an activity as completed
func (p *Project) CompleteActivity(activityName string) {
	if p.Activities == nil {
		p.Activities = make(map[string]Activity)
	}

	activity := p.Activities[activityName]
	activity.Completed = true
	now := time.Now()
	activity.Timestamp = &now
	p.Activities[activityName] = activity
	p.LastUpdated = now
}

// AddFeature adds a feature to the project
func (p *Project) AddFeature(feature Feature) {
	p.Features = append(p.Features, feature)
}

// RequiresDatabase returns true if the project type requires a database
func (p *Project) RequiresDatabase() bool {
	return p.Type == ProjectTypeAPI
}

// SupportsCRUD returns true if the project type supports CRUD generation
func (p *Project) SupportsCRUD() bool {
	return p.Type == ProjectTypeAPI
}
