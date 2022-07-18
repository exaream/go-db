package dbutil_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/exaream/go-db/dbutil"
)

const (
	cfgTyp      = "ini"
	cfgSection  = "example_section"
	cfgHost     = "go_db_mysql"
	cfgDatabase = "example_test"
	cfgUsername = "exampleuser"
	cfgPassword = "examplepasswd"
	cfgProtocol = "tcp"
	driver      = "mysql"
	cfgPort     = 3306
	timeout     = 5
	querySelect = `SELECT id, name, status, created_at, updated_at FROM users WHERE id = :id AND status = :status;`
	queryUpdate = `UPDATE users SET status = :afterSts, updated_at = NOW() WHERE id = :id AND status = :beforeSts;`

	testDataNum = 10 // 50000
	chunkSize   = 10 // 10000
)

var cfgPath = string(filepath.Separator) + filepath.Join("go", "src", "work", "testdata", "example", "example.dsn")

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

func TestMain(m *testing.M) {
	//ctx := context.Context(context.Background())
	// if errs := setup(ctx); errs != nil {
	// 	for _, err := range multierr.Errors(errs) {
	// 		fmt.Fprintln(os.Stderr, err)
	// 	}
	// 	os.Exit(1)
	// }

	code := m.Run()

	os.Exit(code)
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
	ctx := context.Background()
	cfg := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)

	db, err := dbutil.NewDBContext(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	want := 1
	args := map[string]any{"id": 1, "status": 0}
	list, err := dbutil.SelectContext[user](ctx, db, querySelect, args)
	if err != nil {
		t.Error(err)
	}

	if len(list) != want {
		t.Errorf("len(list) want: %d, got: %d", want, len(list))
	}
}

func TestSelectTxContext(t *testing.T) {
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
	args := map[string]any{"id": 1, "status": 0}
	list, err := dbutil.SelectTxContext[user](ctx, tx, querySelect, args)
	if err != nil {
		t.Error(err)
	}

	if len(list) != want {
		t.Errorf("len(list) want: %d, got: %d", want, len(list))
	}
}

func TestUpdateTxContext(t *testing.T) {
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
	args := map[string]any{"id": 1, "beforeSts": 0, "afterSts": 1}
	got, err := dbutil.UpdateTxContext(ctx, tx, queryUpdate, args)
	if err != nil {
		t.Error(err)
	}

	if got != want {
		t.Errorf("num want: %d, got: %d", want, got)
	}
}
