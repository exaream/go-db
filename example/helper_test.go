package example_test

import (
	"context"
	"strconv"
	"time"

	"github.com/exaream/go-db/dbutil"
	"github.com/exaream/go-db/example"
	gimei "github.com/mattn/go-gimei"
	"go.uber.org/multierr"
)

const (
	// SQL
	queryDropTbl   = `DROP TABLE IF EXISTS example_db.users`
	queryCreateTbl = `CREATE TABLE example_db.users (
		id int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
		name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
		email varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
		status int(11) UNSIGNED NOT NULL DEFAULT '0',
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	queryInsert = `INSERT INTO example_db.users (id, name, email, status, created_at, updated_at) 
	    VALUES (:id, :name, :email, :status, :created_at, :updated_at)`
)

// initTblContext initializes table(s) for testing.
// You can also use the following SQL to initialize the testing DB.
// /go/src/work/_local/mysql/setup/ddl/example_db.sql
func initTblContext(ctx context.Context, cfg *dbutil.Config) (err error) {
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

	// Bulk insert
	if _, err := tx.NamedExecContext(ctx, queryInsert, testUsers()); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	if err := tx.Commit(); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	return nil
}

// testUsers returns user data for testing.
func testUsers() (users []example.User) {
	now := time.Now()

	for i := 1; i <= 10; i++ {
		users = append(users, example.User{i, gimei.NewName().Kanji(), "a" + strconv.Itoa(i) + "@examle.com", 0, &now, &now})
	}

	return users
}
