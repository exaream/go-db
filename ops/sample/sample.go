package sample

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"ops/dbx"
)

// Cond has the fields needed to operate a DB.
type Cond struct {
	Writer  io.Writer
	ini     ini
	stmt    stmt
	timeout time.Duration
	where   where
	set     set
}
type ini struct {
	path    string
	section string
}
type stmt struct {
	query   string // Do NOT use the word "select" because it is a reserved word in Go.
	command string
}
type set struct {
	status int
}
type where struct {
	userId int
}

type user struct {
	createdAt *time.Time
	updatedAt *time.Time
	name      string
	id        int
	status    int
}

// NewCond returns the info needed to operate a DB.
func NewCond(iniPath, section string, timeout, userId, status int) *Cond {
	return &Cond{
		Writer:  os.Stdout,
		timeout: time.Duration(timeout) * time.Second,
		ini: ini{
			path:    iniPath,
			section: section,
		},
		// Please change the following when creating your own package.
		stmt: stmt{
			query:   `SELECT id, name, status, created_at, updated_at FROM users;`,
			command: `UPDATE users SET status = ?, updated_at = NOW() WHERE id = ?;`,
		},
		set:   set{status: status},
		where: where{userId: userId},
	}
}

// Run does a DB operation.
func (c *Cond) Run() (rerr error) {
	// Rollback if the time limit is exceeded.
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Get DB handle.
	db, err := dbx.OpenByIni(c.ini.path, c.ini.section)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			rerr = err
		}
	}()

	// Check DB connection.
	if err := db.PingContext(ctx); err != nil {
		return err
	}

	if err := c.log(dbx.LF + "Before operating"); err != nil {
		return err
	}

	if err := dbx.QueryWithContext(ctx, db, c.stmt.query, scanRows); err != nil {
		return err
	}

	// Begin transaction.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Please add validation here when creating your own package.

	// Please change the following when creating your own package.
	// TODO: Confirm how to abstruct and inject the following arguments.
	_, err = tx.ExecContext(ctx, c.stmt.command, c.set.status, c.where.userId)
	if err != nil {
		return dbx.Rollback(tx, rerr, err)
	}

	// fmt.Println(result.RowsAffected())
	if err := c.log(dbx.LF + "Before commit"); err != nil {
		return err
	}

	if err := dbx.QueryTxWithContext(ctx, tx, c.stmt.query, scanRows); err != nil {
		return dbx.Rollback(tx, rerr, err)
	}

	if err := tx.Commit(); err != nil {
		return dbx.Rollback(tx, rerr, err)
	}

	if err := c.log(dbx.LF + "After operating"); err != nil {
		return err
	}

	if err := dbx.QueryWithContext(ctx, db, c.stmt.query, scanRows); err != nil {
		return err
	}

	if err := c.log(""); err != nil {
		return err
	}

	return nil
}

// TODO: Generics で型を指定し ([]*user, error) を返却するためには
// インスタンス化してからしか dbx.QueryWithContext 等の引数として
// scanRows を渡せないため、対応方法をレビュー時に確認
// 今回は *sql.Rows の値を []*user 型で返却できてからテストを書くこととする
func scanRows(ctx context.Context, rows *sql.Rows) (err error) {
	defer func() {
		if rerr := rows.Close(); err == nil && rerr != nil {
			err = rerr
		}
		// Check errors other than EOL error
		if rerr := rows.Err(); err == nil && rerr != nil {
			err = rerr
		}
	}()

	// Please change the following when creating your own package.
	for rows.Next() {
		var u user
		// TODO: Confirm how to abstruct and inject the following arguments.
		err := rows.Scan(&u.id, &u.name, &u.status, &u.createdAt, &u.updatedAt)
		if err != nil {
			return err
		}
		fmt.Println(u.id, u.name, u.status, u.createdAt.Format(dbx.YmdHis), u.updatedAt.Format(dbx.YmdHis))
	}

	return nil
}

func (c *Cond) log(msg string) error {
	_, err := fmt.Fprintln(c.Writer, msg)
	return err
}
