package cli

import (
	"context"

	"github.com/buildwithhp/gophex/internal/app"
	"github.com/buildwithhp/gophex/internal/cmd"
)

// CLI represents the command line interface
type CLI struct {
	app *app.Application
}

// NewCLI creates a new CLI instance
func NewCLI(application *app.Application) *CLI {
	return &CLI{
		app: application,
	}
}

// Execute runs the CLI application
func (c *CLI) Execute(ctx context.Context) error {
	// Use the existing Execute function from the cmd package
	return cmd.Execute()
}
