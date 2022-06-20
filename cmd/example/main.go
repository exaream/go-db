package main

import (
	"context"
	"fmt"
	"os"

	"github.com/exaream/go-db/example"
	"go.uber.org/multierr"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	version        = "0.1.0"
	defaultIniPath = "example.ini"
	defaultSection = "example"
	defaultTimeout = "30s"
)

// Arguments
var (
	app     = kingpin.New("example", "An example command made of Go to operate MySQL.")
	iniPath = app.Flag("ini-path", "Set an ini file path.").Short('i').Default(defaultIniPath).String()
	section = app.Flag("section", "Set a section name.").Short('s').Default(defaultSection).String()
	timeout = app.Flag("timeout", "Set timeout. e.g. 5s").Short('t').Default(defaultTimeout).Duration()
	userId  = app.Flag("user-id", "Set user_id.").Int()
	status  = app.Flag("status", "Set a status.").Int()
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
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	c := example.NewCond(*iniPath, *section, *userId, *status)
	if errs := c.Run(ctx); errs != nil {
		for _, err := range multierr.Errors(errs) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
