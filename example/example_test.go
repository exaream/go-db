package example_test

import (
	"context"
	"testing"
	"time"

	"github.com/exaream/go-db/dbutil"
	"github.com/exaream/go-db/example"
	"go.uber.org/multierr"
)

func TestRun(t *testing.T) {
	prepareDB(t, beforeSqlPath)
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	cond := example.NewCond(1, off, on)
	if errs := example.Run(ctx, cfg, cond); errs != nil {
		for _, err := range multierr.Errors(errs) {
			t.Error(err)
		}
	}
}

func TestNewExecutor(t *testing.T) {
	prepareDB(t, beforeSqlPath)

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
	prepareDB(t, beforeSqlPath)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	cond := example.NewCond(1, off, on)

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := example.ExportPrepare(ex, ctx, cond); err != nil {
		t.Error(err)
	}
}

func TestExec(t *testing.T) {
	prepareDB(t, beforeSqlPath)
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	cond := example.NewCond(1, off, on)

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := example.ExportExec(ex, ctx, cond); err != nil {
		t.Error(err)
	}
}

func TestTeardown(t *testing.T) {
	prepareDB(t, afterSqlPath)
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	cond := example.NewCond(1, off, on)

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := example.ExportTeardown(ex, ctx, cond); err != nil {
		t.Error(err)
	}
}
