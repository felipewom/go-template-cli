package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "template-cli",
	Short: "Template CLI tool to clone a base project from a public GitHub repo",
}

func Execute() error {
	return rootCmd.Execute()
}
