package dbutil_test

import (
	"context"
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
)

// setup initializes a DB for testing.
func setup(ctx context.Context) error {
	cfg, err := dbutil.ParseConfig(cfgTyp, cfgPath, cfgSection)
	if err != nil {
		return err
	}

	if err := initTblContext(ctx, cfg, testDataNum, chunkSize); err != nil {
		return err
	}

	return nil
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
