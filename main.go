package main

import (
	"os"

	"github.com/chippy-ao/redmine-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
