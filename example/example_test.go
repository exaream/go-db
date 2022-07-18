package example_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/exaream/go-db/dbutil"
	"github.com/exaream/go-db/example"
	"go.uber.org/multierr"
)

const (
	cfgTyp     = "ini"
	cfgSection = "example_section"
)

var cfgPath = string(filepath.Separator) + filepath.Join("go", "src", "work", "testdata", "example", "example.dsn")

func TestRun(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	cond := example.NewCond(1, 0, 1)

	if errs := example.Run(ctx, cfg, cond); errs != nil {
		for _, err := range multierr.Errors(errs) {
			t.Error(err)
		}
	}

	cond = example.NewCond(1, 1, 0)
	if errs := example.Run(ctx, cfg, cond); errs != nil {
		for _, err := range multierr.Errors(errs) {
			t.Error(err)
		}
	}
}

func TestNewExecutor(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := ex.DB.PingContext(ctx); err != nil {
		t.Error(err)
	}
}

func TestPrepare(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	cond := example.NewCond(1, 0, 1)

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := ex.ExportPrepare(ctx, cond); err != nil {
		t.Error(err)
	}
}

func TestExec(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	cond := example.NewCond(1, 0, 1)

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := ex.ExportExec(ctx, cond); err != nil {
		t.Error(err)
	}
}

func TestTeardown(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	cond := example.NewCond(1, 0, 1)

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := ex.ExportTeardown(ctx, cond); err != nil {
		t.Error(err)
	}
}
