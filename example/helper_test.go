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
	cfgTyp           = "ini"
	cfgSection       = "example_section"
	queryTruncateTbl = `TRUNCATE TABLE users;`

	on  = 1
	off = 0
)

var (
	testDir       = string(filepath.Separator) + filepath.Join("go", "src", "work", "testdata", "example")
	cfgPath       = filepath.Join(testDir, "example.dsn")
	beforeSqlPath = filepath.Join(testDir, "before_update.sql")
	afterSqlPath  = filepath.Join(testDir, "after_update.sql")
)

func prepareDB(t *testing.T, sqlPath string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query, err := os.ReadFile(sqlPath)
	if err != nil {
		t.Fatal(err)
	}

	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)

	db, err := dbutil.NewDBContext(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

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
