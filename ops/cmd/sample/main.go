package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"go.uber.org/multierr"

	s "ops/sample"
)

const version = "0.1.0"

// Arguments
var (
	app     = kingpin.New("sample", "Sample command made of Go to operate MySQL.")
	userId  = app.Flag("user-id", "Set user_id.").Int()
	status  = app.Flag("status", "Set a status.").Int()
	iniPath = app.Flag("ini-path", "Set an ini file path.").Default(".").String()
	section = app.Flag("section", "Set a section name.").String()
)

func init() {
	app.Version(version)

	// Parse arguments
	if _, err := app.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	c := s.NewCond(*userId, *status, *iniPath, *section)

	if errs := c.Run(); errs != nil {
		for _, err := range multierr.Errors(errs) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
