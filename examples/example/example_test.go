package example_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/exaream/go-db/dbutil"
	"github.com/exaream/go-db/examples/example"
	"go.uber.org/multierr"
)

func TestMain(m *testing.M) {
	if err := initDB(mysqlDBType, beforeSQLPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := initDB(pgsqlDBType, beforeSQLPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	m.Run()
}

func TestRun(t *testing.T) {
	cases := map[string]struct {
		dbType string
		path   string
		id     int
	}{
		"mysql": {mysqlDBType, mysqlCfgPath, 1},
		"pgsql": {pgsqlDBType, pgsqlCfgPath, 1},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			t.Cleanup(func() {
				prepareDB(t, tt.dbType, beforeSQLPath)
			})

			cond := example.NewCond(tt.id, non, active)
			cfg := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			t.Cleanup(cancel)

			if errs := example.Run(ctx, cfg, cond); errs != nil {
				for _, err := range multierr.Errors(errs) {
					t.Error(err)
				}
			}
		})
	}
}

func TestRunErr(t *testing.T) {
	cases := map[string]struct {
		dbType string
		path   string
		id     int
	}{
		"mysql": {mysqlDBType, mysqlCfgPath, 0},
		"pgsql": {pgsqlDBType, pgsqlCfgPath, 0},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			t.Cleanup(func() {
				prepareDB(t, tt.dbType, beforeSQLPath)
			})

			cond := example.NewCond(tt.id, non, active)
			cfg := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			t.Cleanup(cancel)

			if err := example.Run(ctx, cfg, cond); err == nil {
				t.Error("want: error, got: nil")
			}
		})
	}
}

func TestNewExecutor(t *testing.T) {
	cases := map[string]struct {
		dbType string
		path   string
	}{
		"mysql": {mysqlDBType, mysqlCfgPath},
		"pgsql": {pgsqlDBType, pgsqlCfgPath},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			t.Cleanup(func() {
				prepareDB(t, tt.dbType, beforeSQLPath)
			})

			cfg := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			t.Cleanup(cancel)

			ex, err := example.NewExecutor(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}

			if err := ex.DB.PingContext(ctx); err != nil {
				t.Error(err)
			}
		})
	}
}
