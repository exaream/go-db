package example_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/exaream/go-db/dbutil"
	"github.com/jmoiron/sqlx"
)

const (
	timeout = 5
	active  = 1
	non     = 0

	cfgType    = "ini"
	cfgSection = "test_example_section"

	// Config MySQL
	mysqlDBType = "mysql"
	mysqlDriver = "mysql"

	// Config PostgreSQL
	pgsqlDBType = "pgsql"
	pgsqlDriver = "pgx"
)

var (
	testDir       = string(filepath.Separator) + filepath.Join("go", "src", "work", "examples", "testdata", "example")
	mysqlCfgPath  = filepath.Join(testDir, "mysql.dsn")
	pgsqlCfgPath  = filepath.Join(testDir, "pgsql.dsn")
	beforeSQLPath = filepath.Join(testDir, "before_update.sql")
	afterSQLPath  = filepath.Join(testDir, "after_update.sql")

	// Query
	queryTruncateTbls = map[string]string{
		mysqlDriver: `TRUNCATE TABLE users;`,
		pgsqlDriver: `TRUNCATE TABLE users RESTART IDENTITY;`,
	}
)

// prepareDB prepares initial data
func prepareDB(t *testing.T, dbType, sqlPath string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	query, err := os.ReadFile(sqlPath)
	if err != nil {
		t.Fatal(err)
	}

	cfgPath := filepath.Join(testDir, dbType+".dsn")
	cfg := dbutil.NewConfigFile(cfgType, cfgPath, cfgSection)

	db, err := dbutil.NewDBContext(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	args := make(map[string]any)
	queryTruncateTbl := queryTruncateTbls[db.DriverName()]
	_, err = sqlx.NamedExecContext(ctx, db, queryTruncateTbl, args)
	if err != nil {
		t.Fatal(err)
	}

	_, err = sqlx.NamedExecContext(ctx, db, string(query), args)
	if err != nil {
		t.Fatal(err)
	}
}

func initDB(dbType, sqlPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	query, err := os.ReadFile(sqlPath)
	if err != nil {
		return err
	}

	cfgPath := filepath.Join(testDir, dbType+".dsn")
	cfg := dbutil.NewConfigFile(cfgType, cfgPath, cfgSection)

	db, err := dbutil.NewDBContext(ctx, cfg)
	if err != nil {
		return err
	}

	args := make(map[string]any)
	queryTruncateTbl := queryTruncateTbls[db.DriverName()]
	if _, err = sqlx.NamedExecContext(ctx, db, queryTruncateTbl, args); err != nil {
		return err
	}
	if _, err = sqlx.NamedExecContext(ctx, db, string(query), args); err != nil {
		return err
	}

	return nil
}
