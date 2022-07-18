package main

import (
	"context"
	"fmt"
	"os"

	"github.com/exaream/go-db/dbutil"
	"github.com/exaream/go-db/example"
	"go.uber.org/multierr"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version = "v0.2.0"

// Arguments
var (
	app       = kingpin.New("example", "An example command made of Go to operate MySQL.")
	typ       = app.Flag("type", "Set a config type.").Default("ini").String()
	path      = app.Flag("path", "Set a config file path.").Default("example.dsn").String()
	section   = app.Flag("section", "Set a config section name.").Default("example_section").String()
	timeout   = app.Flag("timeout", "Set a timeout value. e.g. 5s").Default("10s").Duration()
	id        = app.Flag("id", "Set an ID.").Required().Uint64()
	beforeSts = app.Flag("before-sts", "Set a before status.").Required().Uint8()
	afterSts  = app.Flag("after-sts", "Set a after status.").Required().Uint8()
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
	cfg := dbutil.NewConfigFile(*typ, *path, *section)
	cond := example.NewCond(*id, *beforeSts, *afterSts)

	if errs := example.Run(ctx, cfg, cond); errs != nil {
		for _, err := range multierr.Errors(errs) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
