package example

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bxcodec/faker/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-gimei"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/exaream/go-db/dbutil"
)

const (
	// Layout of "Y-m-d H:i:s"
	layout = "2006-01-02 15:04:05"

	// SQL
	queryTruncateTbl = `TRUNCATE TABLE users`
	querySelect      = `SELECT id, name, status, created_at, updated_at FROM users WHERE id = :id AND status = :status;`
	queryInsert      = `INSERT INTO users (id, name, email, status, created_at, updated_at) 
VALUES (:id, :name, :email, :status, :created_at, :updated_at)`
	queryUpdate = `UPDATE users SET status = :afterSts, updated_at = NOW() WHERE id = :id AND status = :beforeSts;`
)

// Schema of users table
// Please use exported struct and fields because dbutil package handle these. (rows.StructScan)
type User struct {
	ID        uint64     `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Status    uint8      `db:"status"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// User's stringer.
func (u User) String() string {
	return fmt.Sprintf("%d\t%s\t%v\t%s\t%s", u.ID, u.Name, u.Status, u.CreatedAt.Format(layout), u.UpdatedAt.Format(layout))
}

// Cond has conditions to create SQL.
type Cond struct {
	id        uint64
	beforeSts uint8
	afterSts  uint8
}

// NewCond returns conditions to create SQL.
func NewCond(id uint64, beforeSts, afterSts uint8) *Cond {
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

// Executor has logger and db.
type Executor struct {
	Logger *zap.Logger
	DB     *sqlx.DB
}

// NewExecutor returns Executor.
func NewExecutor(ctx context.Context, cfg *dbutil.ConfigFile) (*Executor, error) {
	db, err := dbutil.NewDBContext(ctx, cfg)
	if err != nil {
		return nil, err
	}

	Logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return &Executor{
		Logger: Logger,
		DB:     db,
	}, nil
}

// init initialize sample data
// TODO: test
func Init(ctx context.Context, cfg *dbutil.ConfigFile, max, chunkSize uint64) (err error) {
	db, err := dbutil.NewDBContext(ctx, cfg)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := db.Close(); rerr != nil {
			err = rerr
		}
	}()

	tx := db.MustBeginTx(ctx, nil)

	if _, err := tx.ExecContext(ctx, queryTruncateTbl); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	if err := dbutil.BulkInsertTxContext(ctx, tx, users, queryInsert, max, chunkSize); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	if err := tx.Commit(); err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	return nil
}

// TODO: test
// TODO: Confirm how to use Generics to specify a type of return value.
func users(min, max uint64) (list []User) {
	now := time.Now()

	var i uint64
	for i = min; i <= max; i++ {
		list = append(list, User{i, gimei.NewName().Kanji(), faker.Email(), 0, &now, &now})
	}

	return list
}

// prepare runs SELECT clause before update.
func (ex *Executor) prepare(ctx context.Context, cond *Cond) error {
	args := map[string]any{"id": cond.id, "status": cond.beforeSts}
	rows, err := dbutil.SelectContext[User](ctx, ex.DB, querySelect, args)
	if err != nil {
		return err
	}

	if len(rows) <= 0 {
		return errors.New("there is no target rows")
	}

	return nil
}

// exec runs UPDATE and SELECT clause on the same transaction.
func (ex *Executor) exec(ctx context.Context, cond *Cond) error {
	tx := ex.DB.MustBeginTx(ctx, nil)

	args := map[string]any{"id": cond.id, "beforeSts": cond.beforeSts, "afterSts": cond.afterSts}
	num, err := dbutil.UpdateTxContext(ctx, tx, queryUpdate, args)
	if err != nil {
		return err
	}

	if num <= 0 {
		return multierr.Append(errors.New("there is no affected rows"), tx.Rollback())
	}

	args = map[string]any{"id": cond.id, "status": cond.afterSts}
	rows, err := dbutil.SelectTxContext[User](ctx, tx, querySelect, args)
	if err != nil {
		return err
	}

	if len(rows) <= 0 {
		return multierr.Append(errors.New("there is no target rows"), tx.Rollback())
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// teardown runs SELECT clause after update.
func (ex *Executor) teardown(ctx context.Context, cond *Cond) error {
	args := map[string]any{"id": cond.id, "status": cond.afterSts}
	rows, err := dbutil.SelectContext[User](ctx, ex.DB, querySelect, args)
	if err != nil {
		return err
	}

	if len(rows) <= 0 {
		return errors.New("there is no affected rows")
	}

	return nil
}
