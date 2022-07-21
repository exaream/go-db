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
	"go.uber.org/multierr"
)

const (
	queryTruncateTbl = `TRUNCATE TABLE users`
	queryInsert      = `INSERT INTO users (id, name, email, status, created_at, updated_at) 
	    VALUES (:id, :name, :email, :status, :created_at, :updated_at)`

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

// initTblContext initializes a table for testing.
func initTblContext(ctx context.Context, cfg *dbutil.Config, max, size uint64) (err error) {
	db, err := dbutil.OpenContext(ctx, cfg)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := db.Close(); rerr != nil {
			err = rerr
		}
	}()

	tx := db.MustBeginTx(ctx, nil)

	if _, err := tx.ExecContext(ctx, queryTruncateTbl); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	if err := bulkInsertTxContext(ctx, tx, max, size); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	if err := tx.Commit(); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	return nil
}

// bulkInsertTxContext executes Bulk Insert on context and transaction.
func bulkInsertTxContext(ctx context.Context, tx *sqlx.Tx, max, size uint64) error {
	var i uint64
	for i = 1; i < max; i += size {
		j := i + size - 1
		if j > max {
			j = max
		}

		if _, err := tx.NamedExecContext(ctx, queryInsert, testUsers(i, j)); err != nil {
			return multierr.Append(err, tx.Rollback())
		}
	}

	return nil
}

// testUsers returns user data for testing.
func testUsers(min, max uint64) (users []user) {
	now := time.Now()
	var i uint64
	for i = min; i <= max; i++ {
		users = append(users, user{i, gimei.NewName().Kanji(), faker.Email(), 0, &now, &now})
	}

	return users
}
