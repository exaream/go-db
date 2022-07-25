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

const (
	version   = "v0.2.0"
	min       = 1
	max       = 50000
	chunkSize = 10000
)

// Arguments
var (
	app       = kingpin.New("example", "An example command made of Go to operate MySQL.")
	initFlg   = app.Flag("init", "Set true if you want to initialize data").Default("false").Bool()
	typ       = app.Flag("type", "Set a config type.").Default("ini").String()
	path      = app.Flag("path", "Set a config file path.").Default("example.dsn").String()
	section   = app.Flag("section", "Set a config section name.").Default("example_section").String()
	timeout   = app.Flag("timeout", "Set a timeout value. e.g. 5s").Default("10s").Duration()
	id        = app.Flag("id", "Set an ID.").Default("0").Int64()
	beforeSts = app.Flag("before-sts", "Set a before status.").Default("0").Int8()
	afterSts  = app.Flag("after-sts", "Set a after status.").Default("0").Int8()
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

	if *initFlg {
		if err := example.Init(ctx, cfg, min, max, chunkSize); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("Succeeded generating initial data.")
		os.Exit(0)
	}

	if errs := example.Run(ctx, cfg, cond); errs != nil {
		for _, err := range multierr.Errors(errs) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
