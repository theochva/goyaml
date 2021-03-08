package main

import (
	// "github.com/jedib0t/go-pretty/text"
	"github.com/theochva/goyaml/commands"
)

var (
	version = "0.0.0"
	commit  = "none"
	date    = "unknown"
)

func main() {
	commands.SetVersion(version, commit, date)
	commands.Execute()
}
