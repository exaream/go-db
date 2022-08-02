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

func TestInit(t *testing.T) {
	cases := map[string]struct {
		dbType    string
		path      string
		min       uint
		max       uint
		chunkSize uint

		want int64
	}{
		// MySQL
		"mysql 1":                 {mysqlDBType, mysqlCfgPath, 1, 1, 1, 1},
		"mysql divisible chunk":   {mysqlDBType, mysqlCfgPath, 1, 100, 10, 100},
		"mysql indivisible chunk": {mysqlDBType, mysqlCfgPath, 1, 100, 11, 100},
		// PostgreSQL
		"pgsql 1":                 {pgsqlDBType, pgsqlCfgPath, 1, 1, 1, 1},
		"pgsql divisible chunk":   {pgsqlDBType, pgsqlCfgPath, 1, 100, 10, 100},
		"pgsql indivisible chunk": {pgsqlDBType, pgsqlCfgPath, 1, 100, 11, 100},
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

			db, err := dbutil.NewDBContext(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := db.Close(); err != nil {
					t.Fatal(err)
				}
			})

			got, err := example.Init(ctx, cfg, tt.min, tt.max, tt.chunkSize)
			if err != nil {
				t.Error(err)
			}
			if got != tt.want {
				t.Errorf("want: %d, got: %d", tt.want, got)
			}

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
				prepareDB(t, tt.dbType, beforeSqlPath)
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
				prepareDB(t, tt.dbType, beforeSqlPath)
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
				prepareDB(t, tt.dbType, beforeSqlPath)
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
			prepareDB(t, tt.dbType, afterSqlPath)
			t.Cleanup(func() {
				prepareDB(t, tt.dbType, beforeSqlPath)
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
			prepareDB(t, tt.dbType, afterSqlPath)
			t.Cleanup(func() {
				prepareDB(t, tt.dbType, beforeSqlPath)
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
