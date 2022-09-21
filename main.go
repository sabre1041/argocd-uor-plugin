package main

import (
	"os"

	"github.com/uor-community/argocd-uor-plugin/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
