package example_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/exaream/go-db/dbutil"
)

const (
	confTyp     = "ini"
	confSection = "example_section"
)

var confPath = string(filepath.Separator) + filepath.Join("go", "src", "work", "cmd", "example", "example.dsn")

// TODO: How to call a helper func in TestMain which does NOT have testing.T.
func TestMain(m *testing.M) {
	ctx := context.Context(context.Background())

	if err := setup(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	code := m.Run()

	os.Exit(code)
}

func setup(ctx context.Context) error {
	conf, err := dbutil.ParseConf(confTyp, confPath, confSection)
	if err != nil {
		return err
	}

	if err := initTableContext(ctx, conf); err != nil {
		return err
	}

	return nil
}
