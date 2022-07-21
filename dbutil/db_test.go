package dbutil_test

import (
	"context"
	"testing"
	"time"

	"github.com/exaream/go-db/dbutil"
)

const (
	timeout     = 5
	driver      = "mysql"
	cfgTyp      = "ini"
	cfgSection  = "example_section"
	cfgHost     = "go_db_mysql"
	cfgDatabase = "example_test"
	cfgUsername = "exampleuser"
	cfgPassword = "examplepasswd"
	cfgProtocol = "tcp"
	cfgPort     = 3306
	querySelect = `SELECT id, name, status, created_at, updated_at FROM users WHERE id = :id AND status = :status;`
	queryUpdate = `UPDATE users SET status = :afterSts, updated_at = NOW() WHERE id = :id AND status = :beforeSts;`
	testDataNum = 10 // 50000
	chunkSize   = 10 // 10000
)

// Schema of users table
// Please use exported struct and fields because dbutil package handle these. (rows.StructScan)
type user struct {
	ID        uint64     `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Status    uint8      `db:"status"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func TestNewDBContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	file := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	db, err := dbutil.NewDBContext(ctx, file)
	if err != nil {
		t.Fatal(err)
	}

	if got := db.DriverName(); got != driver {
		t.Errorf("want: %s, got: %s", driver, got)
	}

	if err := db.PingContext(ctx); err != nil {
		t.Error(err)
	}
}

func TestParseConfig(t *testing.T) {
	cfg, err := dbutil.ParseConfig(cfgTyp, cfgPath, cfgSection)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Host != cfgHost {
		t.Error("failed to get host")
	}
	if cfg.Database != cfgDatabase {
		t.Error("failed to get database")
	}
	if cfg.Username != cfgUsername {
		t.Error("failed to get username")
	}
	if cfg.Password != cfgPassword {
		t.Error("failed to get password")
	}
	if cfg.Protocol != cfgProtocol {
		t.Error("failed to get protocol")
	}
	if cfg.Port != cfgPort {
		t.Error("failed to get protocol")
	}
}

func TestOpenContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	cfg, err := dbutil.ParseConfig(cfgTyp, cfgPath, cfgSection)
	if err != nil {
		t.Fatal(err)
	}

	db, err := dbutil.OpenContext(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if got := db.DriverName(); got != driver {
		t.Errorf("want: %s, got: %s", driver, got)
	}

	if err := db.PingContext(ctx); err != nil {
		t.Error(err)
	}
}

func TestSelectContext(t *testing.T) {
	prepareDB(t, beforeSqlPath)

	ctx := context.Background()
	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)

	db, err := dbutil.NewDBContext(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	want := 1
	args := map[string]any{"id": 1, "status": off}
	list, err := dbutil.SelectContext[user](ctx, db, querySelect, args)
	if err != nil {
		t.Error(err)
	}

	if len(list) != want {
		t.Errorf("len(list) want: %d, got: %d", want, len(list))
	}
}

func TestSelectTxContext(t *testing.T) {
	prepareDB(t, beforeSqlPath)

	ctx := context.Background()
	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)

	db, err := dbutil.NewDBContext(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	tx := db.MustBeginTx(ctx, nil)
	defer func() {
		if err := tx.Rollback(); err != nil {
			t.Fatal(err)
		}
	}()

	want := 1
	args := map[string]any{"id": 1, "status": off}
	list, err := dbutil.SelectTxContext[user](ctx, tx, querySelect, args)
	if err != nil {
		t.Error(err)
	}

	if len(list) != want {
		t.Errorf("len(list) want: %d, got: %d", want, len(list))
	}
}

func TestUpdateTxContext(t *testing.T) {
	prepareDB(t, beforeSqlPath)
	t.Cleanup(func() {
		prepareDB(t, beforeSqlPath)
	})

	ctx := context.Background()
	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)

	db, err := dbutil.NewDBContext(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	tx := db.MustBeginTx(ctx, nil)
	defer func() {
		if err := tx.Rollback(); err != nil {
			t.Fatal(err)
		}
	}()

	var want int64 = 1
	args := map[string]any{"id": 1, "beforeSts": off, "afterSts": on}
	got, err := dbutil.UpdateTxContext(ctx, tx, queryUpdate, args)
	if err != nil {
		t.Error(err)
	}

	if got != want {
		t.Errorf("num want: %d, got: %d", want, got)
	}
}
