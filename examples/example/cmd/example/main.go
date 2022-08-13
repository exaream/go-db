// An example of a tool to operate MySQL and PostgreSQL.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/exaream/go-db/dbutil"
	"github.com/exaream/go-db/examples/example"
	"go.uber.org/multierr"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "v0.3.0"

	// Initial data
	min       = 1
	max       = 50000
	chunkSize = 10000
)

// Arguments
var (
	app       = kingpin.New("example", "An example command made of Go to operate MySQL and PostgreSQL.")
	typ       = app.Flag("type", "Set a config type.").Default("ini").String()
	path      = app.Flag("path", "Set a config file path.").Default("mysql.dsn").String()
	section   = app.Flag("section", "Set a config section name.").Default("example_section").String()
	timeout   = app.Flag("timeout", "Set a timeout value. e.g. 5s").Default("10s").Duration()
	id        = app.Flag("id", "Set an ID.").Default("0").Uint()
	beforeSts = app.Flag("before-sts", "Set a before status.").Default("0").Uint()
	afterSts  = app.Flag("after-sts", "Set a after status.").Default("0").Uint()
	setupFlg  = app.Flag("setup", "Set true if you want to initialize data.").Default("false").Bool()
)

func init() {
	app.Version(version)

	// Parse command's arguments.
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

	// Generate initial data.
	if *setupFlg {
		total, err := example.Setup(ctx, cfg, min, max, chunkSize)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("succeeded to generate %d records as initial data.\n", total)
		os.Exit(0)
	}

	// Run DB update.
	if errs := example.Run(ctx, cfg, cond); errs != nil {
		for _, err := range multierr.Errors(errs) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
