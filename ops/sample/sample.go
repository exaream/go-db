package sample

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"ops/dbutil"
)

const (
	selectStmt = `SELECT id, name, status, created_at, updated_at FROM users;`
	updateStmt = `UPDATE users SET status = ?, updated_at = NOW() WHERE id = ?`
)

// Cond has the fields needed to operate a DB.
// TODO: 様々な項目が混在しているため要精査(何をレシーバーとして渡すべきか)
type Cond struct {
	db      *sql.DB
	tx      *sql.Tx
	iniPath string
	section string
	userId  int
	status  int
	timeout time.Duration
}

type user struct {
	createdAt *time.Time
	updatedAt *time.Time
	name      string
	id        int
	status    int
}

// NewCond returns the info needed to operate a DB.
func NewCond(userId, status int, iniPath, section string, timeout int) *Cond {
	return &Cond{
		iniPath: path.Clean(iniPath),
		section: section,
		userId:  userId,
		status:  status,
		timeout: time.Duration(timeout) * time.Second,
	}
}

// Run does a DB operation.
func (c *Cond) Run() (rerr error) {
	// Rollback if the time limit is exceeded.
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	db, err := dbutil.OpenByConf(c.iniPath, c.section)
	if err != nil {
		return err
	}
	c.db = db // TODO: Is it OK?

	defer func() {
		if err := c.db.Close(); err != nil {
			rerr = err
		}
	}()

	// Check DB connection.
	if err := c.db.PingContext(ctx); err != nil {
		return err
	}

	fmt.Println("-------------------------------------------------")
	fmt.Println("Before operating")

	if err := c.queryContext(ctx); err != nil {
		return err
	}

	// Begin transaction.
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	c.tx = tx // TODO: Is it OK?

	// TODO: Validate before update

	// Update
	_, err = c.tx.ExecContext(ctx, updateStmt, c.status, c.userId)
	if err != nil {
		return c.rollback(rerr, err)
	}

	//fmt.Println(result.RowsAffected())
	fmt.Println("-------------------------------------------------")
	fmt.Println("Before commit")

	if err := c.queryTxContext(ctx); err != nil {
		return c.rollback(rerr, err)
	}

	if err := c.tx.Commit(); err != nil {
		return c.rollback(rerr, err)
	}

	fmt.Println("-------------------------------------------------")
	fmt.Println("After commit")

	if err := c.queryContext(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Cond) queryTxContext(ctx context.Context) (err error) {
	rows, err := c.tx.QueryContext(ctx, selectStmt)
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
		// TODO: 動的に Scan の引数にセットする方法があるか確認
		err := rows.Scan(&u.id, &u.name, &u.status, &u.createdAt, &u.updatedAt)
		if err != nil {
			return err
		}
		fmt.Println(u.id, u.name, u.status, u.createdAt.Format(dbutil.YmdHis), u.updatedAt.Format(dbutil.YmdHis))
	}

	return nil
}

func (c *Cond) queryContext(ctx context.Context) (err error) {
	rows, err := c.db.QueryContext(ctx, selectStmt) // TODO: c や c.db を使う必要があるのはここだけ
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
		// TODO: 動的に Scan の引数にセットする方法があるか確認
		err := rows.Scan(&u.id, &u.name, &u.status, &u.createdAt, &u.updatedAt)
		if err != nil {
			return err
		}
		fmt.Println(u.id, u.name, u.status, u.createdAt.Format(dbutil.YmdHis), u.updatedAt.Format(dbutil.YmdHis))
	}

	return nil
}

// Rollback rollbacks using transaction.
// It can return multiple errors.
func (c *Cond) rollback(rerr, err error) error {
	return dbutil.Rollback(c.tx, rerr, err) // TODO: c や c.tx を使う必要があるのはここだけ
}
