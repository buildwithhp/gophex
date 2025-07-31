package cmd

import (
	"fmt"

	"github.com/buildwithhp/gophex/pkg/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gophex",
	Short: "A CLI tool for generating Go project scaffolding",
	Long: `Gophex is a command-line tool that helps you quickly generate
Go project structures following best practices and standard layouts.

It supports generating various types of Go projects including:
- REST APIs
- Web applications
- Microservices
- CLI tools`,
	Version: version.GetVersion(),
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gophex",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gophex version", version.GetVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
