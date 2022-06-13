package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"ops/dbutil"
)

const (
	timeout = 30                    // seconds
	layout  = "2006-01-02 15:04:05" // Y-m-d H:i:s
)

type user struct {
	createdAt *time.Time
	updatedAt *time.Time
	name      string
	email     string
	id        int
}

func main() {
	// Rollback if the time limit is exceeded.
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	iniPath := filepath.Join("credentials", "foo.ini")
	conf, err := dbutil.ParseConf(iniPath, "sample")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	db, err := conf.Open()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer func() {
		if err := db.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	// Check DB connection.
	if err := db.PingContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Begin transaction.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Check DB values before execution.
	if err := selectUsers(ctx, tx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Update
	email := sampleEmail("alice", "sample.com")
	result, err := tx.ExecContext(ctx, `UPDATE users SET email = ?, updated_at = NOW() WHERE id = ?`, email, 1)

	// Affected rows
	fmt.Println(result.RowsAffected())

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if rerr := tx.Rollback(); rerr != nil {
			fmt.Fprintln(os.Stderr, rerr)
			os.Exit(1)
		}
	}

	// Before commit
	if err := selectUsers(ctx, tx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if rerr := tx.Rollback(); rerr != nil {
			fmt.Fprintln(os.Stderr, rerr)
			os.Exit(1)
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if rerr := tx.Rollback(); rerr != nil {
			fmt.Fprintln(os.Stderr, rerr)
			os.Exit(1)
		}
	}
}

func selectUsers(ctx context.Context, tx *sql.Tx) (err error) {
	rows, err := tx.QueryContext(ctx, `SELECT id, name, email, created_at, updated_at FROM users;`)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := rows.Close(); err == nil && rerr != nil {
			err = rerr
		}
		// Check errors other than EOL error
		if rerr := rows.Err(); err == nil && rerr != nil {
			err = rerr
		}
	}()

	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.email, &u.createdAt, &u.updatedAt)
		if err != nil {
			return err
		}
		fmt.Println(u.id, u.name, u.email, u.createdAt.Format(layout), u.updatedAt.Format(layout))
	}

	return nil
}

func sampleEmail(prefix, domain string) string {
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))
	return prefix + timestamp + "@" + domain
}
