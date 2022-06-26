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
	// SQL
	querySelect = `SELECT id, name, status, created_at, updated_at FROM users;`
	queryUpdate = `UPDATE users SET status = :status, updated_at = NOW() WHERE id = :id;`
)

// Struct of users table
type User struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Status    int        `db:"status"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// User's stringer
func (u User) String() string {
	return fmt.Sprintf("%d\t%s\t%v\t%s\t%s", u.ID, u.Name, u.Status, u.CreatedAt.Format(dbutil.YmdHis), u.UpdatedAt.Format(dbutil.YmdHis))
}

// Conf has configurations to create DB handle.
type Conf struct {
	typ     string
	path    string
	section string
}

// Cond has conditions to create SQL.
type Cond struct {
	id     int
	status int
}

type executor struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewConf returns configurations to create DB handle.
func NewConf(typ, path, section string) *Conf {
	return &Conf{
		typ:     typ,
		path:    path,
		section: section,
	}
}

// NewCond returns conditions to create SQL.
func NewCond(id, status int) *Cond {
	return &Cond{
		id:     id,
		status: status,
	}
}

// Run does a DB operation.
func Run(ctx context.Context, conf *Conf, cond *Cond) (rerr error) {
	var ex *executor
	ex, err := newExecutor(ctx, conf)
	if err != nil {
		return err
	}

	defer func() {
		if err := ex.db.Close(); err != nil {
			rerr = err
		}
	}()

	ex.logger.Info("before operations")
	if err := ex.check(ctx); err != nil {
		return err
	}

	ex.logger.Info("before commit")
	tx := ex.db.MustBeginTx(ctx, nil)
	if err := ex.exec(ctx, cond, tx); err != nil {
		return err
	}

	ex.logger.Info("after commit")
	return ex.check(ctx)
}

func newExecutor(ctx context.Context, conf *Conf) (*executor, error) {
	logger, err := zap.NewDevelopment() // TODO: zap options
	if err != nil {
		return nil, err
	}

	dbConf, err := dbutil.ParseConf(conf.typ, conf.path, conf.section)
	if err != nil {
		return nil, err
	}

	db, err := dbutil.OpenWithContext(ctx, dbConf)
	if err != nil {
		return nil, err
	}

	return &executor{
		logger: logger,
		db:     db,
	}, nil
}

func (ex *executor) check(ctx context.Context) error {
	users := []User{}
	if err := ex.db.SelectContext(ctx, &users, querySelect); err != nil {
		return err
	}

	for _, u := range users {
		fmt.Println(u)
	}

	return nil
}

func (ex *executor) exec(ctx context.Context, cond *Cond, tx *sqlx.Tx) (rerr error) {
	args := map[string]any{"id": cond.id, "status": cond.status}

	_, err := tx.NamedExecContext(ctx, queryUpdate, args)
	if err != nil {
		return multierr.Append(err, tx.Rollback())
	}

	users := []User{}
	err = tx.SelectContext(ctx, &users, querySelect)
	if err != nil {
		return err
	}

	for _, u := range users {
		fmt.Println(u)
	}

	if err := tx.Commit(); err != nil {
		return multierr.Append(rerr, err)
	}

	return nil
}
