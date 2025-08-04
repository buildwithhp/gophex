package generator

import (
	"context"

	"github.com/buildwithhp/gophex/internal/domain/project"
	"github.com/buildwithhp/gophex/internal/generator"
)

// generatorAdapter adapts the existing generator to implement the domain Generator interface
type generatorAdapter struct {
	generator *generator.Generator
}

// NewGeneratorAdapter creates a new generator adapter
func NewGeneratorAdapter() project.Generator {
	return &generatorAdapter{
		generator: generator.New(),
	}
}

// Generate implements the domain Generator interface
func (g *generatorAdapter) Generate(ctx context.Context, proj *project.Project) error {
	// Convert domain types to generator types
	var dbConfig *generator.DatabaseConfig
	if proj.DatabaseConfig != nil {
		dbConfig = &generator.DatabaseConfig{
			Type:         string(proj.DatabaseConfig.Type),
			ConfigType:   string(proj.DatabaseConfig.ConfigType),
			Host:         proj.DatabaseConfig.Host,
			Port:         proj.DatabaseConfig.Port,
			Username:     proj.DatabaseConfig.Username,
			Password:     proj.DatabaseConfig.Password,
			DatabaseName: proj.DatabaseConfig.DatabaseName,
			ReadHost:     proj.DatabaseConfig.ReadHost,
			WriteHost:    proj.DatabaseConfig.WriteHost,
			ClusterNodes: proj.DatabaseConfig.ClusterNodes,
			SSLMode:      proj.DatabaseConfig.SSLMode,
			AuthSource:   proj.DatabaseConfig.AuthSource,
			ReplicaSet:   proj.DatabaseConfig.ReplicaSet,
		}
	}

	var redisConfig *generator.RedisConfig
	if proj.RedisConfig != nil {
		redisConfig = &generator.RedisConfig{
			Enabled:  proj.RedisConfig.Enabled,
			Host:     proj.RedisConfig.Host,
			Port:     proj.RedisConfig.Port,
			Password: proj.RedisConfig.Password,
			Database: proj.RedisConfig.Database,
		}
	}

	// Call the existing generator
	return g.generator.GenerateWithFullConfig(
		string(proj.Type),
		proj.Name,
		proj.Path,
		dbConfig,
		redisConfig,
	)
}
