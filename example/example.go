package example

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/multierr"

	"github.com/exaream/go-db/dbx"
	"github.com/go-logr/logr"
)

type Executor[S, R, T any] struct {
	DB     *sql.DB
	logger *logr.Logger
}

type Action[S, R, T any] struct {
	Setup    func(ctx context.Context, tx *sql.Tx) (S, error)
	Run      func(ctx context.Context, tx *sql.Tx, s S) (R, error)
	Teardown func(ctx context.Context, db *sql.DB, r R) (T, error)
}

// NOTE: We wrap `do()` with `Do()` to avoid named arguments appering in the document of pkg.go.dev.
func (e *Executor[S, R, T]) Do(ctx context.Context, act *Action[S, R, T]) (T, error) {
	return e.do(ctx, act)
}

// NOTE: We use `zeroT` for returning the zero value of `T` type when an error occurs.
func (e *Executor[S, R, T]) do(ctx context.Context, act *Action[S, R, T]) (zeroT T, _ error) {
	tx, err := e.DB.BeginTx(ctx, nil)
	if err != nil {
		return zeroT, err
	}

	// TODO: Check current DB values. Return an error if the values does not meet preconditions.
	s, err := act.Setup(ctx, tx)
	if err != nil {
		return zeroT, multierr.Append(err, tx.Rollback())
	}

	// TODO: Insert or update DB values. Return an error if the values does not meet postconditions.
	r, err := act.Run(ctx, tx, s)
	if err != nil {
		return zeroT, multierr.Append(err, tx.Rollback())
	}

	if err := tx.Commit(); err != nil {
		return zeroT, multierr.Append(err, tx.Rollback())
	}

	// TODO: Check DB values after commit whether they meet postconditions.
	t, err := act.Teardown(ctx, e.DB, r)
	if err != nil {
		return zeroT, err
	}

	return t, nil
}

// Cond has the fields needed to operate a DB.
type Cond struct {
	Writer io.Writer
	ini    ini
	stmt   stmt
	where  where
	set    set
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
func NewCond(iniPath, section string, userId, status int) *Cond {
	return &Cond{
		Writer: os.Stdout,
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
//TODO: How to shorten this function
func (c *Cond) Run(ctx context.Context) (rerr error) {
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

	records, err := dbx.QueryWithContext(ctx, db, c.stmt.query, scanRows)
	if err != nil {
		return err
	}
	fmt.Println(records)

	// Begin transaction.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Please add validation here when creating your own package.

	// Please change the following when creating your own package.
	// TODO: How to abstruct and inject the following arguments.
	_, err = tx.ExecContext(ctx, c.stmt.command, c.set.status, c.where.userId)
	if err != nil {
		return dbx.Rollback(tx, rerr, err)
	}

	// fmt.Println(result.RowsAffected())
	if err := c.log(dbx.LF + "Before commit"); err != nil {
		return err
	}

	records, err = dbx.QueryTxWithContext(ctx, tx, c.stmt.query, scanRows)
	if err != nil {
		return err
	}
	fmt.Println(records)

	if err := tx.Commit(); err != nil {
		return dbx.Rollback(tx, rerr, err)
	}

	if err := c.log(dbx.LF + "After operating"); err != nil {
		return err
	}

	records, err = dbx.QueryWithContext(ctx, db, c.stmt.query, scanRows)
	if err != nil {
		return err
	}
	fmt.Println(records)

	if err := c.log(""); err != nil {
		return err
	}

	return nil
}

// TODO: How to apply `user` type to `records` using generics.
func scanRows(ctx context.Context, rows *sql.Rows) (_ dbx.Records, err error) {
	defer func() {
		if rerr := rows.Close(); err == nil && rerr != nil {
			err = rerr
		}
		// Check errors other than EOL error
		if rerr := rows.Err(); err == nil && rerr != nil {
			err = rerr
		}
	}()

	records := make(dbx.Records) // 並行時の競合を避けるため初期化

	// Please change the following when creating your own package.
	for rows.Next() {
		var u user
		// TODO: How to abstruct and inject the following arguments.
		err := rows.Scan(&u.id, &u.name, &u.status, &u.createdAt, &u.updatedAt)
		if err != nil {
			return nil, err
		}

		records[u.id] = map[string]any{
			"id":        u.id,
			"name":      u.name,
			"status":    u.status,
			"createdAt": u.createdAt.Format(dbx.YmdHis),
			"updatedAt": u.updatedAt.Format(dbx.YmdHis),
		}

		fmt.Println(u.id, u.name, u.status, u.createdAt.Format(dbx.YmdHis), u.updatedAt.Format(dbx.YmdHis))
	}

	return records, nil
}

func (c *Cond) log(msg string) error {
	_, err := fmt.Fprintln(c.Writer, msg)
	return err
}
