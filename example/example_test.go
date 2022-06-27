package example_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/exaream/go-db/dbutil"
	"go.uber.org/multierr"
)

const (
	cfgTyp      = "ini"
	cfgSection  = "example_section"
	testDataNum = 50000
	chunkSize   = 10000
)

var cfgPath = string(filepath.Separator) + filepath.Join("go", "src", "work", "cmd", "example", "example.dsn")

// TODO: How to call a helper func in TestMain which does NOT have testing.T.
func TestMain(m *testing.M) {
	ctx := context.Context(context.Background())

	if errs := setup(ctx); errs != nil {
		for _, err := range multierr.Errors(errs) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}

	code := m.Run()

	os.Exit(code)
}

func setup(ctx context.Context) error {
	cfg, err := dbutil.ParseConfig(cfgTyp, cfgPath, cfgSection)
	if err != nil {
		return err
	}

	if err := initTblContext(ctx, cfg, testDataNum, chunkSize); err != nil {
		return err
	}

	return nil
}
