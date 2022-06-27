package example

import (
	"context"
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
	querySelect = `SELECT id, name, status, created_at, updated_at FROM users WHERE id = :id;`
	queryUpdate = `UPDATE users SET status = :status, updated_at = NOW() WHERE id = :id;`
)

// Struct of users table
type User struct {
	ID        uint64     `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Status    uint8      `db:"status"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// User's stringer
func (u User) String() string {
	return fmt.Sprintf("%d\t%s\t%v\t%s\t%s", u.ID, u.Name, u.Status, u.CreatedAt.Format(YmdHis), u.UpdatedAt.Format(YmdHis))
}

// Cond has conditions to create SQL.
type Cond struct {
	id        uint64
	beforeSts uint8
	afterSts  uint8
}

// Config has configurations to create DB handle.
type Config struct {
	typ     string
	path    string
	section string
}

type executor struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewConfig returns configurations to create DB handle.
func NewConfig(typ, path, section string) *Config {
	return &Config{
		typ:     typ,
		path:    path,
		section: section,
	}
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
func Run(ctx context.Context, cfg *Config, cond *Cond) (rerr error) {
	var ex *executor
	ex, err := newExecutor(ctx, cfg)
	if err != nil {
		return err
	}

	ex.logger.Info("Start")
	defer func() {
		if err := ex.db.Close(); err != nil {
			rerr = err
		}
		ex.logger.Info("End")
	}()

	if _, err := ex.selectContext(ctx, cond); err != nil {
		return err
	}

	tx := ex.db.MustBeginTx(ctx, nil)
	if _, err := ex.updateTxContext(ctx, tx, cond); err != nil {
		return err
	}

	if _, err := ex.selectTxContext(ctx, tx, cond); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if _, err := ex.selectContext(ctx, cond); err != nil {
		return err
	}
	return nil
}

func newExecutor(ctx context.Context, cfg *Config) (*executor, error) {
	logger, err := zap.NewDevelopment() // TODO: zap options
	if err != nil {
		return nil, err
	}

	dbCfg, err := dbutil.ParseConfig(cfg.typ, cfg.path, cfg.section)
	if err != nil {
		return nil, err
	}

	db, err := dbutil.OpenContext(ctx, dbCfg)
	if err != nil {
		return nil, err
	}

	return &executor{
		logger: logger,
		db:     db,
	}, nil
}

func (ex *executor) selectContext(ctx context.Context, cond *Cond) ([]User, error) {
	args := map[string]any{"id": cond.id}
	rows, err := ex.db.NamedQueryContext(ctx, querySelect, args)
	if err != nil {
		return nil, err
	}

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.StructScan(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
		fmt.Println(u)
	}

	return users, nil
}

func (ex *executor) updateTxContext(ctx context.Context, tx *sqlx.Tx, cond *Cond) (int64, error) {
	args := map[string]any{"id": cond.id, "status": cond.afterSts}
	result, err := tx.NamedExecContext(ctx, queryUpdate, args)
	if err != nil {
		return 0, multierr.Append(err, tx.Rollback())
	}

	num, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (ex *executor) selectTxContext(ctx context.Context, tx *sqlx.Tx, cond *Cond) ([]User, error) {
	args := map[string]any{"id": cond.id}
	rows, err := sqlx.NamedQueryContext(ctx, tx, querySelect, args)
	if err != nil {
		return nil, err
	}

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.StructScan(&u); err != nil {
			return nil, multierr.Append(err, tx.Rollback())
		}
		users = append(users, u)
		fmt.Println(u)
	}

	return users, nil
}
