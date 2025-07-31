package main

import (
	"fmt"
	"os"

	"github.com/buildwithhp/gophex/internal/cmd"
)

func main() {
		// Interactive mode
	fmt.Println("🚀 Welcome to Gophex!")
	fmt.Println("A CLI tool for generating Go project scaffolding")
	fmt.Println()

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
