package main

import (
	"go-template-cli/cmd"
	_ "go-template-cli/internal/config"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
