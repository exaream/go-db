package sample

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/alecthomas/kingpin"
	_ "github.com/go-sql-driver/mysql"

	"ops/dbutil"
)

const (
	timeout = 30                    // seconds
	layout  = "2006-01-02 15:04:05" // Y-m-d H:i:s
)

type Cond struct {
	IniPath string
	Section string
	Status  int
	UserId  int
}

type user struct {
	createdAt *time.Time
	updatedAt *time.Time
	name      string
	status    int
	id        int
}

const version = "0.1.0"

var (
	// Command arguments
	app     = kingpin.New("sample", "Sample command made of Go to operate MySQL.")
	userId  = app.Flag("user-id", "Set user_id.").Int()
	status  = app.Flag("status", "Set a status.").Int()
	iniPath = app.Flag("ini-path", "Set an ini file path.").Default(".").String()
	section = app.Flag("section", "Set a section name.").String()
)

func init() {
	app.Version(version)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewCond(userId, status int, iniPath, section string) *Cond {
	return &Cond{
		IniPath: iniPath,
		Section: section,
		UserId:  userId,
		Status:  status,
	}
}

func IniPath() string {
	return path.Clean(*iniPath)
}

func Section() string {
	return *section
}

func UserId() int {
	return *userId
}

func Status() int {
	return *status
}

// TODO: Move to the internal directory. ===================================
// TODO: Get command's arguments.
func (c *Cond) Run() (rerr error) {
	// Rollback if the time limit is exceeded.
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	dbInfo, err := dbutil.ParseConf(c.IniPath, c.Section)
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
	_, err = tx.ExecContext(ctx, `UPDATE users SET status = ?, updated_at = NOW() WHERE id = ?`, c.Status, c.UserId)
	if err != nil {
		return dbutil.Rollback(tx, rerr, err)
	}

	//fmt.Println(result.RowsAffected())
	fmt.Println("------------------")

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
