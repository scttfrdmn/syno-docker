package main

import (
	"os"

	"github.com/scttfrdmn/synodeploy/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
