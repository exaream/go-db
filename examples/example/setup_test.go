package example_test

import (
	"context"
	"testing"
	"time"

	"github.com/exaream/go-db/dbutil"
	"github.com/exaream/go-db/examples/example"
)

func TestSetup(t *testing.T) {
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
				prepareDB(t, tt.dbType, beforeSQLPath)
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

			got, err := example.Setup(ctx, cfg, tt.min, tt.max, tt.chunkSize)
			if err != nil {
				t.Error(err)
			}
			if got != tt.want {
				t.Errorf("want: %d, got: %d", tt.want, got)
			}
		})
	}
}

func TestFakeUsers(t *testing.T) {
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
