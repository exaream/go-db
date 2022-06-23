package example

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/exaream/go-db/dbx"
)

const (
	stmtQuery   = `SELECT id, name, status, created_at, updated_at FROM users;`
	stmtCommand = `UPDATE users SET status = ?, updated_at = NOW() WHERE id = ?;`
)

// Conf has configurations to create DB handle.
type Conf struct {
	typ     string
	dir     string
	stem    string
	section string
}

// Cond has conditions to create SQL.
type Cond struct {
	id     int
	status int
}

type user struct {
	createdAt *time.Time
	updatedAt *time.Time
	name      string
	id        int
	status    int
}

// NewConf returns configurations to create DB handle.
func NewConf(typ, dir, stem, section string) *Conf {
	return &Conf{
		typ:     typ,
		dir:     dir,
		stem:    stem,
		section: section,
	}
}

// NewCond returns conditions to create SQL.
func NewCond(id, status int) *Cond {
	return &Cond{
		status: status,
		id:     id,
	}
}

// Run does a DB operation.
func Run(ctx context.Context, conf *Conf, cond *Cond) (rerr error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	db, err := dbx.OpenWithContext(ctx, conf.typ, conf.dir, conf.stem, conf.section)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			rerr = err
		}
	}()

	logger.Info("Before operation")

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

	_, err = tx.ExecContext(ctx, stmtCommand, cond.status, cond.id)
	if err != nil {
		return multierr.Append(rerr, err)
	}

	// fmt.Println(result.RowsAffected())
	logger.Info("Before commit")

	records, err = dbx.QueryTxWithContext(ctx, tx, stmtQuery, scanRows)
	if err != nil {
		return err
	}
	fmt.Println(records)

	if err := tx.Commit(); err != nil {
		return multierr.Append(rerr, err)
	}

	logger.Info("After commit")

	records, err = dbx.QueryWithContext(ctx, db, stmtQuery, scanRows)
	if err != nil {
		return err
	}
	fmt.Println(records)

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
