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
	initDB(mysqlDBType, beforeSqlPath)
	initDB(pgsqlDBType, beforeSqlPath)
	m.Run()
}

func TestRun(t *testing.T) {
	cases := map[string]struct {
		dbType string
		path   string
		id     uint
	}{
		"mysql": {mysqlDBType, mysqlCfgPath, 1},
		"pgsql": {pgsqlDBType, pgsqlCfgPath, 1},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			t.Cleanup(func() {
				prepareDB(t, tt.dbType, beforeSqlPath)
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
		id     uint
	}{
		"mysql": {mysqlDBType, mysqlCfgPath, 0},
		"pgsql": {pgsqlDBType, pgsqlCfgPath, 0},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			t.Cleanup(func() {
				prepareDB(t, tt.dbType, beforeSqlPath)
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
				prepareDB(t, tt.dbType, beforeSqlPath)
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
