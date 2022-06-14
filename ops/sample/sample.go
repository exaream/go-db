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

const layout = "2006-01-02 15:04:05" // Y-m-d H:i:s

type Cond struct {
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

func NewCond(userId, status int, iniPath, section string, timeout int) *Cond {
	return &Cond{
		iniPath: path.Clean(iniPath),
		section: section,
		userId:  userId,
		status:  status,
		timeout: time.Duration(timeout) * time.Second,
	}
}

func (c *Cond) Run() (rerr error) {
	// Rollback if the time limit is exceeded.
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	dbInfo, err := dbutil.ParseConf(c.iniPath, c.section)
	if err != nil {
		return err
	}

	db, err := dbInfo.Open()
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

	// Begin transaction.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Check DB values before execution.
	if err := selectUsers(ctx, tx); err != nil {
		return err
	}

	// TODO: Validate before update

	// Update
	_, err = tx.ExecContext(ctx, `UPDATE users SET status = ?, updated_at = NOW() WHERE id = ?`, c.status, c.userId)
	if err != nil {
		return dbutil.Rollback(tx, rerr, err)
	}

	//fmt.Println(result.RowsAffected())
	fmt.Println("-------------------------------------------------")

	// Before commit
	if err := selectUsers(ctx, tx); err != nil {
		return dbutil.Rollback(tx, rerr, err)
	}

	if err := tx.Commit(); err != nil {
		return dbutil.Rollback(tx, rerr, err)
	}

	return nil
}

func selectUsers(ctx context.Context, tx *sql.Tx) (err error) {
	rows, err := tx.QueryContext(ctx, `SELECT id, name, status, created_at, updated_at FROM users;`)
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
		err := rows.Scan(&u.id, &u.name, &u.status, &u.createdAt, &u.updatedAt)
		if err != nil {
			return err
		}
		fmt.Println(u.id, u.name, u.status, u.createdAt.Format(layout), u.updatedAt.Format(layout))
	}

	return nil
}
