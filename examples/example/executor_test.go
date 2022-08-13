package example_test

import (
	"context"
	"testing"
	"time"

	"github.com/exaream/go-db/dbutil"
	"github.com/exaream/go-db/examples/example"
)

func TestPrepare(t *testing.T) {
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
				prepareDB(t, tt.dbType, beforeSQLPath)
			})

			cond := example.NewCond(tt.id, non, active)
			cfg := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			t.Cleanup(cancel)

			ex, err := example.NewExecutor(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}

			if err := example.ExportPrepare(ex, ctx, cond); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestPrepareErr(t *testing.T) {
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
				prepareDB(t, tt.dbType, beforeSQLPath)
			})

			cond := example.NewCond(tt.id, non, active)
			cfg := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			t.Cleanup(cancel)

			ex, err := example.NewExecutor(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}

			if err := example.ExportPrepare(ex, ctx, cond); err == nil {
				t.Error("want: error, got: nil")
			}
		})
	}
}

func TestExec(t *testing.T) {
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
				prepareDB(t, tt.dbType, beforeSQLPath)
			})

			cond := example.NewCond(tt.id, non, active)
			cfg := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)
			ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
			t.Cleanup(cancel)

			ex, err := example.NewExecutor(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}

			if err := example.ExportExec(ex, ctx, cond); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestExecErr(t *testing.T) {
	t.Parallel()

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
				prepareDB(t, tt.dbType, beforeSQLPath)
			})

			cond := example.NewCond(tt.id, non, active)
			cfg := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)
			ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
			t.Cleanup(cancel)

			ex, err := example.NewExecutor(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}

			if err := example.ExportExec(ex, ctx, cond); err == nil {
				t.Error("want: error, got: nil")
			}
		})
	}
}

func TestTeardown(t *testing.T) {
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
			prepareDB(t, tt.dbType, afterSQLPath)
			t.Cleanup(func() {
				prepareDB(t, tt.dbType, beforeSQLPath)
			})

			cond := example.NewCond(tt.id, non, active)
			cfg := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)
			ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
			t.Cleanup(cancel)

			ex, err := example.NewExecutor(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}

			if err := example.ExportTeardown(ex, ctx, cond); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestTeardownErr(t *testing.T) {
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
			prepareDB(t, tt.dbType, afterSQLPath)
			t.Cleanup(func() {
				prepareDB(t, tt.dbType, beforeSQLPath)
			})

			cond := example.NewCond(tt.id, non, active)
			cfg := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)
			ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
			t.Cleanup(cancel)

			ex, err := example.NewExecutor(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}

			if err := example.ExportTeardown(ex, ctx, cond); err == nil {
				t.Error("want: error, got: nil")
			}
		})
	}
}
