package example_test

import (
	"context"
	"testing"
	"time"

	"github.com/exaream/go-db/dbutil"
	"github.com/exaream/go-db/example"
	"go.uber.org/multierr"
)

func TestMain(m *testing.M) {
	var t *testing.T
	defer prepareDB(t, beforeSqlPath) // TODO: Confirm that using defer in TestMain is OK
	prepareDB(t, beforeSqlPath)
	m.Run()
}

func TestRun(t *testing.T) {
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath) // TODO: Confirm better way. Is it a good practice to revert the DB for each test?
	})

	var id uint = 1
	cond := example.NewCond(id, off, on)
	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // TODO: Is it better to use t.Cleanup?

	if errs := example.Run(ctx, cfg, cond); errs != nil {
		for _, err := range multierr.Errors(errs) {
			t.Error(err)
		}
	}
}

// TODO: TestRun と TestRunErr をまとめたほうが良いか確認
func TestRunErr(t *testing.T) {
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	var id uint = 0
	cond := example.NewCond(id, off, on)
	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := example.Run(ctx, cfg, cond); err == nil {
		t.Error("want: error, got: nil")
	}
}

func TestNewExecutor(t *testing.T) {
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := ex.DB.PingContext(ctx); err != nil {
		t.Error(err)
	}
}

func TestInit(t *testing.T) {
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	cases := map[string]struct {
		min       uint
		max       uint
		chunkSize uint

		want int64
	}{
		"one":               {1, 1, 1, 1},
		"divisible chunk":   {1, 100, 10, 100},
		"indivisible chunk": {1, 100, 11, 100},
	}

	for name, tt := range cases {
		tt := tt

		t.Run(name, func(t *testing.T) {
			cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			db, err := dbutil.NewDBContext(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}

			got, err := example.Init(ctx, cfg, tt.min, tt.max, tt.chunkSize)
			if err != nil {
				t.Error(err)
			}
			if got != tt.want {
				t.Errorf("want: %d, got: %d", tt.want, got)
			}

			t.Cleanup(func() {
				if rerr := db.Close(); rerr != nil {
					t.Fatal(rerr)
				}
			})
		})
	}
}

func TestUsers(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		min uint
		max uint

		want int
	}{
		"min = 0":  {0, 1, 0},
		"max = 0":  {1, 0, 0},
		"both = 0": {0, 0, 0},
		"max = 1":  {1, 1, 1},
		"max > 1":  {1, 2, 2},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			users := example.ExportFakeUsers(tt.min, tt.max)
			if got := len(users); got != tt.want {
				t.Errorf("len(users) want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestPrepare(t *testing.T) {
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	var id uint = 1
	cond := example.NewCond(id, off, on)
	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := example.ExportPrepare(ex, ctx, cond); err != nil {
		t.Error(err)
	}
}

func TestPrepareErr(t *testing.T) {
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	var id uint = 0
	cond := example.NewCond(id, off, on)
	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := example.ExportPrepare(ex, ctx, cond); err == nil {
		t.Error("want: error, got: nil")
	}
}

func TestExec(t *testing.T) {
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	var id uint = 1
	cond := example.NewCond(id, off, on)
	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := example.ExportExec(ex, ctx, cond); err != nil {
		t.Error(err)
	}
}

func TestExecErr(t *testing.T) {
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	var id uint = 0
	cond := example.NewCond(id, off, on)
	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := example.ExportExec(ex, ctx, cond); err == nil {
		t.Error("want: error, got: nil")
	}
}

func TestTeardown(t *testing.T) {
	prepareDB(t, afterSqlPath)
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	var id uint = 1
	cond := example.NewCond(id, off, on)
	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := example.ExportTeardown(ex, ctx, cond); err != nil {
		t.Error(err)
	}
}

func TestTeardownErr(t *testing.T) {
	prepareDB(t, afterSqlPath)
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	var id uint = 0
	cond := example.NewCond(id, off, on)
	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ex, err := example.NewExecutor(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := example.ExportTeardown(ex, ctx, cond); err == nil {
		t.Error("want: error, got: nil")
	}
}
