package example

import (
	"context"
	"errors"

	"github.com/exaream/go-db/dbutil"
	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

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
