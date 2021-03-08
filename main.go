package main

import (
	"github.com/theochva/goyaml/commands"
)

var (
	version = "0.0.0"
	commit  = "none"
	date    = "unknown"
)

func main() {
	commands.NewGoyamlApp(version, commit, date).Execute()
}
