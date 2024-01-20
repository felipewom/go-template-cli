package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scaffolder",
	Short: "Project Scaffolder CLI tool to clone a base project from a public GitHub repo",
}

func Execute() error {
	return rootCmd.Execute()
}
