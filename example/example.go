package example

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/multierr"

	"github.com/exaream/go-db/dbx"
)

const (
	stmtQuery   = `SELECT id, name, status, created_at, updated_at FROM users;`
	stmtCommand = `UPDATE users SET status = ?, updated_at = NOW() WHERE id = ?;`
)

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
