package example

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/exaream/go-db/dbutil"
)

const (
	stmtQuery   = `SELECT id, name, status, created_at, updated_at FROM users;`
	stmtCommand = `UPDATE users SET status = :status, updated_at = NOW() WHERE id = :id;`
)

type User struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Status    int        `db:"status"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

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

	db, err := dbutil.OpenWithContext(ctx, conf.typ, conf.dir, conf.stem, conf.section)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			rerr = err
		}
	}()

	logger.Info("Before operation")

	var users []User
	err = db.SelectContext(ctx, &users, stmtQuery)
	if err != nil {
		return err
	}
	for _, u := range users {
		fmt.Println(u.ID, u.Name, u.Status, u.CreatedAt.Format(dbutil.YmdHis), u.UpdatedAt.Format(dbutil.YmdHis))
	}

	// Begin transaction.
	tx := db.MustBeginTx(ctx, nil)

	args := map[string]any{"id": cond.id, "status": cond.status}
	_, err = tx.NamedExecContext(ctx, stmtCommand, args)
	if err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	//fmt.Println(result.RowsAffected())
	logger.Info("Before commit")

	users = []User{}
	err = tx.SelectContext(ctx, &users, stmtQuery)
	if err != nil {
		return err
	}
	for _, u := range users {
		fmt.Println(u.ID, u.Name, u.Status, u.CreatedAt.Format(dbutil.YmdHis), u.UpdatedAt.Format(dbutil.YmdHis))
	}
	if err := tx.Commit(); err != nil {
		return multierr.Append(rerr, err)
	}

	logger.Info("After commit")
	users = []User{}
	err = db.SelectContext(ctx, &users, stmtQuery)
	if err != nil {
		return err
	}
	for _, u := range users {
		fmt.Println(u.ID, u.Name, u.Status, u.CreatedAt.Format(dbutil.YmdHis), u.UpdatedAt.Format(dbutil.YmdHis))
	}

	return nil
}
