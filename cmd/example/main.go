package main

import (
	"context"
	"fmt"
	"os"

	"github.com/exaream/go-db/example"
	"go.uber.org/multierr"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version = "0.1.0"

// Arguments
var (
	app     = kingpin.New("example", "An example command made of Go to operate MySQL.")
	typ     = app.Flag("type", "Set an config type.").Default("ini").String()
	dir     = app.Flag("dir", "Set an config file path.").Default(".").String()
	stem    = app.Flag("stem", "Set a config stem name.").Default("example").String()
	section = app.Flag("section", "Set a config section name.").Default("example_section").String()
	timeout = app.Flag("timeout", "Set timeout. e.g. 5s").Default("30s").Duration()
	id      = app.Flag("id", "Set id.").Int()
	status  = app.Flag("status", "Set a status.").Int()
)

type Conf struct {
	Host     string
	DB       string
	Username string
	Password string
	Protocol string
	Tz       string
	Port     int
}

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

	conf := example.NewConf(*typ, *dir, *stem, *section)
	cond := example.NewCond(*id, *status)

	if errs := example.Run(ctx, conf, cond); errs != nil {
		for _, err := range multierr.Errors(errs) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
