package example_test

import (
	"context"
	"path/filepath"
	"time"

	"github.com/exaream/go-db/dbutil"
	"github.com/exaream/go-db/example"
	"go.uber.org/multierr"
)

const (
	// DB config
	confType    = "ini"
	confStem    = "example"
	confSection = "example_section"
	// SQL
	dropTblStmt   = "DROP TABLE IF EXISTS `example_db`.`users`"
	createTblStmt = `CREATE TABLE example_db.users (
		id int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
		name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
		email varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
		status int(11) UNSIGNED NOT NULL DEFAULT '0',
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	insertStmt = "INSERT INTO `example_db`.`users` (`id`, `name`, `email`, `status`, `created_at`, `updated_at`) VALUES (:id, :name, :email, :status, :created_at, :updated_at)"
)

var confDir = string(filepath.Separator) + filepath.Join("go", "src", "work", "cmd", "example")

// We can use the following SQL to initialize DB.
// /go/src/work/_local/mysql/setup/ddl/example_db.sql
// But I wrote bulk insert for as a sample.
func initTable(typ, dir, stem, section string) (err error) {
	ctx := context.Context(context.Background())

	db, err := dbutil.OpenWithContext(ctx, typ, dir, stem, section)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := db.Close(); rerr != nil {
			err = rerr
		}
	}()

	tx := db.MustBeginTx(ctx, nil)

	if _, err := tx.ExecContext(ctx, dropTblStmt); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	if _, err := tx.ExecContext(ctx, createTblStmt); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	defTime, _ := time.Parse(dbutil.YmdHis, "2022-01-01 00:00:00")
	var users = []example.User{
		{1, "Alice", "example1@example.com", 0, &defTime, &defTime},
		{2, "Billy", "example2@example.com", 0, &defTime, &defTime},
		{3, "Chris", "example3@example.com", 0, &defTime, &defTime},
	}

	if _, err := tx.NamedExecContext(ctx, insertStmt, users); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	return nil
}
