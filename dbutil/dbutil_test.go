package dbutil_test

import (
	"context"
	"testing"
	"time"

	"github.com/exaream/go-db/dbutil"
)

// Schema of users table
// Please use exported struct and fields because dbutil package handle these. (rows.StructScan)
type User struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Status    int        `db:"status"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func TestNewDBContext(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		path string
	}{
		"mysql": {mysqlCfgPath},
		"pgsql": {pgsqlCfgPath},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
			t.Cleanup(cancel)

			file := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)
			db, err := dbutil.NewDBContext(ctx, file)
			if err != nil {
				t.Fatal(err)
			}

			cfg, err := dbutil.ParseConfig(cfgType, tt.path, cfgSection)
			if err != nil {
				t.Fatal(err)
			}

			if got := db.DriverName(); got != cfg.Driver {
				t.Errorf("want: %s, got: %s", cfg.Driver, got)
			}

			if err := db.PingContext(ctx); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestOpenContext(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		path string
	}{
		"mysql": {mysqlCfgPath},
		"pgsql": {pgsqlCfgPath},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
			t.Cleanup(cancel)

			cfg, err := dbutil.ParseConfig(cfgType, tt.path, cfgSection)
			if err != nil {
				t.Fatal(err)
			}

			db, err := dbutil.OpenContext(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}

			want := cfg.Driver
			if got := db.DriverName(); got != want {
				t.Errorf("want: %s, got: %s", want, got)
			}

			if err := db.PingContext(ctx); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestOpenContextErr(t *testing.T) {
	t.Parallel()

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

			ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
			t.Cleanup(cancel)

			cfg, err := dbutil.ParseConfig(cfgType, tt.path, cfgSection)
			if err != nil {
				t.Fatal(err)
			}
			cfg.Port = dummyPort

			switch tt.dbType {
			case mysqlDBType:
				dsn, err := dbutil.ExportDataSrcMySQL(cfg)
				if err != nil {
					t.Fatal(err)
				}
				cfg.DataSrc = dsn
			case pgsqlDBType:
				cfg.DataSrc = dbutil.ExportDataSrcPgSQL(cfg)
			}

			if _, err := dbutil.OpenContext(ctx, cfg); err == nil {
				t.Error("want: error, got: nil")
			}
		})
	}
}

func TestSelectContext(t *testing.T) {
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

			prepareDB(t, tt.dbType, beforeSQLPath)

			ctx := context.Background()
			f := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)

			db, err := dbutil.NewDBContext(ctx, f)
			if err != nil {
				t.Fatal(err)
			}

			want := 1 // record
			args := map[string]any{"id": 1, "status": non}
			list, err := dbutil.SelectContext[User](ctx, db, querySelect, args)
			if err != nil {
				t.Error(err)
			}

			if len(list) != want {
				t.Errorf("len(list) want: %d, got: %d", want, len(list))
			}
		})
	}
}

func TestSelectTxContext(t *testing.T) {
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

			prepareDB(t, tt.dbType, beforeSQLPath)

			ctx := context.Background()
			f := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)

			db, err := dbutil.NewDBContext(ctx, f)
			if err != nil {
				t.Fatal(err)
			}

			tx := db.MustBeginTx(ctx, nil)
			t.Cleanup(func() {
				if err := tx.Rollback(); err != nil {
					t.Fatal(err)
				}
			})

			want := 1 // record
			args := map[string]any{"id": 1, "status": non}
			list, err := dbutil.SelectTxContext[User](ctx, tx, querySelect, args)
			if err != nil {
				t.Error(err)
			}

			if len(list) != want {
				t.Errorf("len(list) want: %d, got: %d", want, len(list))
			}
		})
	}
}

func TestUpdateTxContext(t *testing.T) {
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

			prepareDB(t, tt.dbType, beforeSQLPath)

			ctx := context.Background()
			f := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)

			db, err := dbutil.NewDBContext(ctx, f)
			if err != nil {
				t.Fatal(err)
			}

			tx := db.MustBeginTx(ctx, nil)
			t.Cleanup(func() {
				if err := tx.Rollback(); err != nil {
					t.Fatal(err)
				}
			})

			var want int64 = 1
			args := map[string]any{"id": 1, "beforeSts": non, "afterSts": active}
			got, err := dbutil.UpdateTxContext(ctx, tx, queryUpdate, args)
			if err != nil {
				t.Error(err)
			}

			if got != want {
				t.Errorf("num want: %d, got: %d", want, got)
			}
		})
	}
}

func TestBulkInsertTxContext(t *testing.T) {
	cases := map[string]struct {
		path string
	}{
		"mysql": {mysqlCfgPath},
		"pgsql": {pgsqlCfgPath},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			f := dbutil.NewConfigFile(cfgType, tt.path, cfgSection)

			db, err := dbutil.NewDBContext(ctx, f)
			if err != nil {
				t.Fatal(err)
			}

			var min, max, chunkSize = 1, 5000, 1000
			tx := db.MustBeginTx(ctx, nil)

			num, err := dbutil.BulkInsertTxContext(ctx, tx, fakeUsers, queryInsert, min, max, chunkSize)
			if err != nil {
				t.Error(err)
			}

			if num != int64(max) {
				t.Errorf("num want: %d, got: %d", max, num)
			}

			if err := tx.Rollback(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
