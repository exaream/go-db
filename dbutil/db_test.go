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
	ID        uint       `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Status    uint       `db:"status"`
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

func TestParseConfig(t *testing.T) {
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

			want := expectedConfig(t, tt.dbType)
			got, err := dbutil.ParseConfig(cfgType, tt.path, cfgSection)
			if err != nil {
				t.Fatal(err)
			}

			if got.Host != want.Host {
				t.Fatalf("host want: %s, got: %s", want.Host, got.Host)
			}
			if got.Database != want.Database {
				t.Fatalf("database want: %s, got: %s", want.Database, got.Database)
			}
			if got.Username != want.Username {
				t.Fatalf("username want: %s, got: %s", want.Username, got.Username)
			}
			if got.Password != want.Password {
				t.Fatalf("password want: %s, got: %s", want.Password, got.Password)
			}
			if got.Protocol != want.Protocol {
				t.Fatalf("protocol want: %s, got: %s", want.Protocol, got.Protocol)
			}
			if got.Port != want.Port {
				t.Fatalf("port want: %d, got: %d", want.Port, got.Port)
			}
			if got.Tz != want.Tz {
				t.Fatalf("timezone want: %s, got: %s", want.Tz, got.Tz)
			}
		})
	}
}

func TestParseConfigErr(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		typ     string
		path    string
		section string
	}{
		"all":            {dummy, dummy, dummy},
		"path":           {cfgType, dummy, cfgSection},
		"type(mysql)":    {dummy, mysqlCfgPath, cfgSection},
		"type(pgsql)":    {dummy, pgsqlCfgPath, cfgSection},
		"section(mysql)": {cfgType, mysqlCfgPath, dummy},
		"section(pgsql)": {cfgType, pgsqlCfgPath, dummy},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := dbutil.ParseConfig(tt.typ, tt.path, tt.section)
			if err == nil {
				t.Error("want: error, got: nil")
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
				cfg.DataSrc = dbutil.ExportDataSrcMySQL(cfg)
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

			prepareDB(t, tt.dbType, beforeSqlPath)

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

			prepareDB(t, tt.dbType, beforeSqlPath)

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

			prepareDB(t, tt.dbType, beforeSqlPath)

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

// TODO: Confirm why int type is chosen over uint type. e.g. result.RowsAffected()
// SEE:  https://github.com/golang/go/issues/49311
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

			var min, max, chunkSize uint = 1, 5000, 1000
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
