package example

import (
	"context"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/exaream/go-db/dbutil"
)

const (
	// Layout of "Y-m-d H:i:s"
	YmdHis = "2006-01-02 15:04:05"
	// SQL
	querySelect = `SELECT id, name, status, created_at, updated_at FROM users WHERE id = :id AND status = :status;`
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
	return fmt.Sprintf("%d\t%s\t%v\t%s\t%s", u.ID, u.Name, u.Status, u.CreatedAt.Format(YmdHis), u.UpdatedAt.Format(YmdHis))
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
	ex, err := newExecutor(ctx, cfg)
	if err != nil {
		return err
	}

	defer func() {
		if err := ex.db.Close(); err != nil {
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

type executor struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func newExecutor(ctx context.Context, cfg *dbutil.ConfigFile) (*executor, error) {
	db, err := dbutil.NewDBContext(ctx, cfg)
	if err != nil {
		return nil, err
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return &executor{
		logger: logger,
		db:     db,
	}, nil
}

func (ex *executor) prepare(ctx context.Context, cond *Cond) error {
	args := map[string]any{"id": cond.id, "status": cond.beforeSts}
	rows, err := dbutil.SelectContext[User](ctx, ex.db, querySelect, args)
	if err != nil {
		return err
	}

	if len(rows) <= 0 {
		return errors.New("there is no target rows")
	}

	return nil
}

func (ex *executor) exec(ctx context.Context, cond *Cond) error {
	tx := ex.db.MustBeginTx(ctx, nil)

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

func (ex *executor) teardown(ctx context.Context, cond *Cond) error {
	args := map[string]any{"id": cond.id, "status": cond.afterSts}
	rows, err := dbutil.SelectContext[User](ctx, ex.db, querySelect, args)
	if err != nil {
		return err
	}

	if len(rows) <= 0 {
		return errors.New("there is no affected rows")
	}

	return nil
}
