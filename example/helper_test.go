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

	cfgType          = "ini"
	cfgSection       = "example_test_section"
	queryTruncateTbl = `TRUNCATE TABLE users;`

	// Config MySQL
	mysqlHost   = "go_db_mysql"
	mysqlDBType = "mysql"
	mysqlDriver = "mysql"
	mysqlPort   = 3306

	// Config PostgreSQL
	pgsqlHost   = "go_db_pgsql"
	pgsqlDBType = "pgsql"
	pgsqlDriver = "pgx"
	pgsqlPort   = 5432

	// Dummy
	dummy     = "dummy"
	dummyPort = 9999

	// Query
	queryInsert = `INSERT INTO users (name, email, status, created_at, updated_at) 
VALUES (:name, :email, :status, :created_at, :updated_at)`
	querySelect = `SELECT id, name, status, created_at, updated_at FROM users WHERE id = :id AND status = :status;`
	queryUpdate = `UPDATE users SET status = :afterSts, updated_at = NOW() WHERE id = :id AND status = :beforeSts;`
)

var (
	testDir       = string(filepath.Separator) + filepath.Join("go", "src", "work", "testdata", "example")
	mysqlCfgPath  = filepath.Join(testDir, "mysql.dsn")
	pgsqlCfgPath  = filepath.Join(testDir, "pgsql.dsn")
	beforeSqlPath = filepath.Join(testDir, "before_update.sql")
	afterSqlPath  = filepath.Join(testDir, "after_update.sql")

	// Query
	queryTruncateTbls = map[string]string{
		"mysql": `TRUNCATE TABLE users`,
		"pgsql": `TRUNCATE TABLE users RESTART IDENTITY`,
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
	_, err = sqlx.NamedExecContext(ctx, db, queryTruncateTbls[dbType], args)
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
	if _, err = sqlx.NamedExecContext(ctx, db, queryTruncateTbls[dbType], args); err != nil {
		return err
	}
	if _, err = sqlx.NamedExecContext(ctx, db, string(query), args); err != nil {
		return err
	}

	return nil
}
