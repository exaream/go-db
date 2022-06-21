package example

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/multierr"

	"github.com/exaream/go-db/dbx"
	"github.com/go-logr/logr"
)

const (
	stmtQuery   = `SELECT id, name, status, created_at, updated_at FROM users;`
	stmtCommand = `UPDATE users SET status = ?, updated_at = NOW() WHERE id = ?;`
)

type (
	S any
	R any
	T any
)
type Executor[S, R, T any] struct {
	DB     *sql.DB
	logger *logr.Logger
}

type Action[S, R, T any] struct {
	//Setup    func(ctx context.Context, tx *sql.Tx) (S, error)
	//Run      func(ctx context.Context, tx *sql.Tx, s S) (R, error)
	//Teardown func(ctx context.Context, db *sql.DB, r R) (T, error)
}

func Do(ctx context.Context, iniPath, section string) (err error) {
	db, err := dbx.OpenByIniWithContext(ctx, iniPath, section)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := db.Close(); err != nil {
			err = rerr
		}
	}()

	var e *Executor[S, R, T]
	var act *Action[S, R, T]
	e.DB = db

	_, err = e.Do(ctx, act)

	if err != nil {
		return err
	}

	return nil
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

func (act *Action[S, R, T]) Setup(ctx context.Context, tx *sql.Tx) (S, error) {
	return act.setup(ctx, tx)
}

func (act *Action[S, R, T]) setup(ctx context.Context, tx *sql.Tx) (zeroS S, err error) {
	return zeroS, err
}

func (act *Action[S, R, T]) Run(ctx context.Context, tx *sql.Tx, s S) (R, error) {
	return act.run(ctx, tx, s)
}

func (act *Action[S, R, T]) run(ctx context.Context, tx *sql.Tx, s S) (zeroR R, err error) {
	return zeroR, err
}

func (act *Action[S, R, T]) Teardown(ctx context.Context, db *sql.DB, r R) (T, error) {
	return act.teardown(ctx, db, r)
}

func (act *Action[S, R, T]) teardown(ctx context.Context, db *sql.DB, r R) (zeroT T, err error) {
	return zeroT, err
}

// Cond has the fields needed to operate a DB.
type Cond struct {
	userId int
	status int
}

type user struct {
	createdAt *time.Time
	updatedAt *time.Time
	name      string
	id        int
	status    int
}

// NewCond returns the info needed to operate a DB.
func NewCond(userId, status int) *Cond {
	return &Cond{
		status: status,
		userId: userId,
	}
}

// Run does a DB operation.
//TODO: How to shorten this function
func (c *Cond) Run(ctx context.Context, iniPath, section string) (rerr error) {
	// Get DB handle.
	db, err := dbx.OpenByIniWithContext(ctx, iniPath, section)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			rerr = err
		}
	}()

	fmt.Println(dbx.LF + "Before operating")

	records, err := dbx.QueryWithContext(ctx, db, stmtQuery, scanRows)
	if err != nil {
		return err
	}
	fmt.Println(records)

	// Begin transaction.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, stmtCommand, c.status, c.userId)
	if err != nil {
		return multierr.Append(rerr, err)
	}

	// fmt.Println(result.RowsAffected())
	fmt.Println(dbx.LF + "Before commit")

	records, err = dbx.QueryTxWithContext(ctx, tx, stmtQuery, scanRows)
	if err != nil {
		return err
	}
	fmt.Println(records)

	if err := tx.Commit(); err != nil {
		return multierr.Append(rerr, err)
	}

	fmt.Println(dbx.LF + "After operating")

	records, err = dbx.QueryWithContext(ctx, db, stmtQuery, scanRows)
	if err != nil {
		return err
	}
	fmt.Println(records)

	fmt.Println("")

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
	_, err := fmt.Println(msg)
	return err
}
