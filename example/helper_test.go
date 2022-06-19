package example_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/exaream/go-rdb/dbx"
)

const (
	testIniPath = "/go/src/work/cmd/example/example.ini"
	testSection = "example"
	testTimeout = 30
)

// As types and order of of firlds are different from "user" in example package,
// I wrote the following struct in example_test package separately.
// Notice: Do NOT apply fieldalignment to the following "user".
type user struct {
	id        int
	name      string
	email     string
	status    int
	createdAt string
	updatedAt string
}

var testUsers = []user{
	{1, "Alice", "example1@example.com", 0, "2022-01-01 00:00:00", "2022-01-01 00:00:00"},
	{2, "Bob", "example2@example.com", 0, "2022-01-01 00:00:00", "2022-01-01 00:00:00"},
	{3, "Chris", "example3@example.com", 0, "2022-01-01 00:00:00", "2022-01-01 00:00:00"},
}

// We can use the following SQL to initialize DB.
// /go/src/work/local/mysql/setup/ddl/example_db.sql
// But I wrote the process using Go for learning the language.
func initTable(iniPath, section string) (rerr error) {
	ctx := context.Context(context.Background())

	db, err := dbx.OpenByIni(iniPath, section)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			rerr = err
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, dropTblStmt()); err != nil {
		return dbx.Rollback(tx, rerr, err)
	}

	if _, err := tx.ExecContext(ctx, createTblStmt()); err != nil {
		return dbx.Rollback(tx, rerr, err)
	}

	numFields := reflect.TypeOf(user{}).NumField() // the number of fields of struct
	placeHolders := make([]string, 0, len(testUsers))
	values := make([]any, 0, len(testUsers)*numFields)

	for _, u := range testUsers {
		placeHolders = append(placeHolders, " (?, ?, ?, ?, ?, ?)")

		values = append(values, u.id)
		values = append(values, u.name)
		values = append(values, u.email)
		values = append(values, u.status)
		values = append(values, u.createdAt)
		values = append(values, u.updatedAt)
	}

	query := fmt.Sprintf(insertTblStmt(), strings.Join(placeHolders, ","))

	if _, err := tx.ExecContext(ctx, query, values...); err != nil {
		return dbx.Rollback(tx, rerr, err)
	}

	return nil
}

func dropTblStmt() string {
	return "DROP TABLE IF EXISTS `example_db`.`users`"
}

func createTblStmt() string {
	// TODO: Confirm how to escape back slashes in bash slashes.
	return `CREATE TABLE example_db.users (
		id int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
		name varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
		email varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
		status int(11) UNSIGNED NOT NULL DEFAULT '0',
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
}

func insertTblStmt() string {
	return "INSERT INTO `example_db`.`users` (`id`, `name`, `email`, `status`, `created_at`, `updated_at`) VALUES %s;"
}
