// Package example is a simple tool for updating column status of table users.
package example

import (
	"context"
	"fmt"
	"time"

	"github.com/exaream/go-db/dbutil"
)

const (
	// Layout of "Y-m-d H:i:s"
	layout = "2006-01-02 15:04:05"

	// Driver
	mysqlDriver = "mysql"
	pgsqlDriver = "pgx"

	// SQL
	querySelect = `SELECT id, name, status, created_at, updated_at FROM users WHERE id = :id AND status = :status;`
	queryInsert = `INSERT INTO users (name, email, status, created_at, updated_at) 
VALUES (:name, :email, :status, :created_at, :updated_at);`
	queryUpdate = `UPDATE users SET status = :afterSts, updated_at = NOW() WHERE id = :id AND status = :beforeSts;`
)

var queryTruncateTbls = map[string]string{
	mysqlDriver: `TRUNCATE TABLE users;`,
	pgsqlDriver: `TRUNCATE TABLE users RESTART IDENTITY;`,
}

// Schema of users table
// Please use exported struct and fields because dbutil package handle these. (rows.StructScan)
type User struct {
	ID        uint       `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Status    int        `db:"status"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// User's stringer.
func (u User) String() string {
	return fmt.Sprintf("%d\t%s\t%v\t%s\t%s", u.ID, u.Name, u.Status, u.CreatedAt.Format(layout), u.UpdatedAt.Format(layout))
}

// Cond has conditions to create SQL.
type Cond struct {
	id        uint
	beforeSts uint
	afterSts  uint
}

// NewCond returns conditions to create SQL.
func NewCond(id, beforeSts, afterSts uint) *Cond {
	return &Cond{
		id:        id,
		beforeSts: beforeSts,
		afterSts:  afterSts,
	}
}

// Run does a DB operation.
func Run(ctx context.Context, cfg *dbutil.ConfigFile, cond *Cond) (rerr error) {
	ex, err := NewExecutor(ctx, cfg)
	if err != nil {
		return err
	}

	defer func() {
		if err := ex.DB.Close(); err != nil {
			rerr = err
		}
	}()

	if err := ex.prepare(ctx, cond); err != nil {
		return err
	}

	if err := ex.exec(ctx, cond); err != nil {
		return err
	}

	if err := ex.teardown(ctx, cond); err != nil {
		return err
	}

	return nil
}
