package sample

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"ops/dbx"
)

const (
	selectStmt = `SELECT id, name, status, created_at, updated_at FROM users;`
	updateStmt = `UPDATE users SET status = ?, updated_at = NOW() WHERE id = ?`
)

// Cond has the fields needed to operate a DB.
type Cond struct {
	timeout time.Duration
	iniPath string
	section string
	userId  int
	status  int
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
		timeout: time.Duration(timeout) * time.Second,
		iniPath: path.Clean(iniPath),
		section: section,
		userId:  userId,
		status:  status,
	}
}

// Run does a DB operation.
func (c *Cond) Run() (rerr error) {
	// Rollback if the time limit is exceeded.
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	db, err := dbx.OpenByIni(c.iniPath, c.section)
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

	fmt.Println("-------------------------------------------------")
	fmt.Println("Before operating")

	if err := dbx.QueryWithContext(ctx, db, selectStmt, scanRows); err != nil {
		return err
	}

	// Begin transaction.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// TODO: Validate before update

	// Update
	_, err = tx.ExecContext(ctx, updateStmt, c.status, c.userId)
	if err != nil {
		return dbx.Rollback(tx, rerr, err)
	}

	//fmt.Println(result.RowsAffected())
	fmt.Println("-------------------------------------------------")
	fmt.Println("Before commit")

	if err := dbx.QueryTxWithContext(ctx, tx, selectStmt, scanRows); err != nil {
		return dbx.Rollback(tx, rerr, err)
	}

	if err := tx.Commit(); err != nil {
		return dbx.Rollback(tx, rerr, err)
	}

	fmt.Println("-------------------------------------------------")
	fmt.Println("After commit")

	if err := dbx.QueryWithContext(ctx, db, selectStmt, scanRows); err != nil {
		return err
	}

	return nil
}

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

	for rows.Next() {
		var u user
		// TODO: 動的に Scan の引数にセットする方法があるか確認
		err := rows.Scan(&u.id, &u.name, &u.status, &u.createdAt, &u.updatedAt)
		if err != nil {
			return err
		}
		fmt.Println(u.id, u.name, u.status, u.createdAt.Format(dbx.YmdHis), u.updatedAt.Format(dbx.YmdHis))
	}

	return nil
}
