package main

import (
	"project-scaffolder-cli/cmd"
	_ "project-scaffolder-cli/internal/config"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
