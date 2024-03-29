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
	timeout = 5
	active  = 1
	non     = 0

	// Config Common
	cfgType     = "ini"
	cfgSection  = "test_dbutil_section"
	cfgDatabase = "test_dbutil_db"
	cfgUsername = "exampleuser"
	cfgPassword = "examplepasswd"
	cfgProtocol = "tcp"
	cfgTz       = "Asia/Tokyo"
	cfgSSLMode  = "disable" // for PostgreSQL

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
VALUES (:name, :email, :status, :created_at, :updated_at);`
	querySelect = `SELECT id, name, status, created_at, updated_at FROM users WHERE id = :id AND status = :status;`
	queryUpdate = `UPDATE users SET status = :afterSts, updated_at = NOW() WHERE id = :id AND status = :beforeSts;`
)

var (
	// Path
	testDir       = string(filepath.Separator) + filepath.Join("go", "src", "work", "testdata")
	mysqlCfgPath  = filepath.Join(testDir, "mysql.dsn")
	pgsqlCfgPath  = filepath.Join(testDir, "pgsql.dsn")
	beforeSQLPath = filepath.Join(testDir, "before_update.sql")

	// Query
	queryTruncateTbls = map[string]string{
		mysqlDriver: `TRUNCATE TABLE users;`,
		pgsqlDriver: `TRUNCATE TABLE users RESTART IDENTITY;`,
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

	queryTruncateTbl := queryTruncateTbls[db.DriverName()]
	args := make(map[string]any)
	_, err = sqlx.NamedExecContext(ctx, db, queryTruncateTbl, args)
	if err != nil {
		t.Fatal(err)
	}

	_, err = sqlx.NamedExecContext(ctx, db, string(query), args)
	if err != nil {
		t.Fatal(err)
	}
}

// expectedConfig returns the Config that tests expect.
func expectedConfig(t *testing.T, dbType string) *dbutil.Config {
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
	case mysqlDBType:
		cfg.Host = mysqlHost
		cfg.Port = mysqlPort
		cfg.Driver = mysqlDriver
		dsn, err := dbutil.ExportDataSrcMySQL(cfg)
		if err != nil {
			t.Fatal(err)
		}
		cfg.DataSrc = dsn
		return cfg
	case pgsqlDBType:
		cfg.Host = pgsqlHost
		cfg.Port = pgsqlPort
		cfg.Driver = pgsqlDriver
		cfg.SSLMode = cfgSSLMode
		cfg.DataSrc = dbutil.ExportDataSrcPgSQL(cfg)
		return cfg
	default:
		return nil
	}
}

// fakeUsers returns fake user list.
func fakeUsers(min, max int) []*User {
	if min == 0 || max == 0 {
		return nil
	}

	users := make([]*User, 0, max)
	now := time.Now()

	for i := min; i <= max; i++ {
		users = append(users, &User{i, gimei.NewName().Kanji(), faker.Email(), 0, &now, &now})
	}

	return users
}
