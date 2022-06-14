package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"go.uber.org/multierr"
	s "ops/sample"
)

const (
	version        = "0.1.0"
	defaultIniPath = "credentials/foo.ini"
	defaultSection = "sample"
	defaultTimeout = 30
)

// Arguments
var (
	app     = kingpin.New("sample", "Sample command made of Go to operate MySQL.")
	userId  = app.Flag("user-id", "Set user_id.").Int()
	status  = app.Flag("status", "Set a status.").Int()
	iniPath = app.Flag("ini-path", "Set an ini file path.").Default(defaultIniPath).String()
	section = app.Flag("section", "Set a section name.").Default(defaultSection).String()
	timeout = app.Flag("timeout", "Set seconds for timeout.").Int()
)

func init() {
	app.Version(version)

	// Parse arguments
	if _, err := app.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Set default value.
	// Because we can NOT use Int() and Default() at the same time.
	if *timeout == 0 {
		*timeout = defaultTimeout
	}
}

func main() {
	c := s.NewCond(*userId, *status, *iniPath, *section, *timeout)

	if errs := c.Run(); errs != nil {
		for _, err := range multierr.Errors(errs) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
