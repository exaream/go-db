package example_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/exaream/go-db/dbutil"
	"github.com/exaream/go-db/example"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-gimei"
	"go.uber.org/multierr"
)

const (
	// SQL
	queryDropTbl   = `DROP TABLE IF EXISTS example_db.users`
	queryCreateTbl = `CREATE TABLE example_db.users (
		id bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
		name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
		email varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
		status tinyint(11) UNSIGNED NOT NULL DEFAULT 0,
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	queryInsert = `INSERT INTO example_db.users (id, name, email, status, created_at, updated_at) 
	    VALUES (:id, :name, :email, :status, :created_at, :updated_at)`
)

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

	if _, err := tx.ExecContext(ctx, queryDropTbl); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	if _, err := tx.ExecContext(ctx, queryCreateTbl); err != nil {
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
func testUsers(min, max uint64) (users []example.User) {
	now := time.Now()
	var i uint64
	for i = min; i <= max; i++ {
		users = append(users, example.User{i, gimei.NewName().Kanji(), faker.Email(), 0, &now, &now})
	}

	return users
}
