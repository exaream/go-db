package dbutil_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/exaream/go-db/dbutil"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-gimei"
)

const (
	timeout      = 5
	active  uint = 1
	non     uint = 0

	// Config Common
	cfgType     = "ini"
	cfgSection  = "dbutil_test_section"
	cfgDatabase = "example_db_dbutil_pkg_test"
	cfgUsername = "exampleuser"
	cfgPassword = "examplepasswd"
	cfgProtocol = "tcp"
	cfgTz       = "Asia/Tokyo"

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
	// Path
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

// prepareDB prepares a DB.
func prepareDB(t *testing.T, dbType, sqlPath string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
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

// wantConfig returns the Config that tests expect.
func wantedConfig(t *testing.T, dbType string) *dbutil.Config {
	t.Helper()

	cfg := &dbutil.Config{
		Type:     dbType,
		Database: cfgDatabase,
		Username: cfgUsername,
		Password: cfgPassword,
		Protocol: cfgProtocol,
		Tz:       cfgTz,
	}

	switch dbType {
	case "mysql":
		cfg.Host = mysqlHost
		cfg.Port = mysqlPort
		cfg.Driver = mysqlDriver
		cfg.Src = dbutil.ExportDataSrcMySQL(cfg)
		return cfg
	case "pgsql":
		cfg.Host = pgsqlHost
		cfg.Port = pgsqlPort
		cfg.Driver = pgsqlDriver
		cfg.Src = dbutil.ExportDataSrcPgSQL(cfg)
		return cfg
	default:
		return nil
	}
}

// fakeUsers returns fake user list.
// TODO: Confirm how to pass *testing.T as an argument.
func fakeUsers(min, max uint) (list []userWithoutID) {
	if min == 0 || max == 0 {
		return list
	}

	now := time.Now()
	for i := min; i <= max; i++ {
		list = append(list, userWithoutID{gimei.NewName().Kanji(), faker.Email(), 0, &now, &now})
	}

	return list
}
