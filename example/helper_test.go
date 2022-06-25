package example_test

import (
	"context"
	"path/filepath"
	"strconv"
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

var confDir = string(filepath.Separator) + filepath.Join("go", "src", "work", "cmd", "example")

// You can also use the following SQL to initialize the testing DB.
// /go/src/work/_local/mysql/setup/ddl/example_db.sql
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

	if _, err := tx.ExecContext(ctx, queryDropTbl); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	if _, err := tx.ExecContext(ctx, queryCreateTbl); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	if _, err := tx.NamedExecContext(ctx, queryInsert, testUsers()); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	return nil
}

func testUsers() []example.User {
	var users []example.User
	names := map[int]string{1: "Alice", 2: "Bobby", 3: "Chris", 4: "Daisy", 5: "Elise"}
	now := time.Now()

	for i := 1; i <= len(names); i++ {
		users = append(users, example.User{i, names[i], "example" + strconv.Itoa(i) + "@examle.com", 0, &now, &now})
	}

	return users
}
