package dbutil_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/exaream/go-db/dbutil"
	"github.com/exaream/go-db/example"
)

const (
	cfgTyp      = "ini"
	cfgSection  = "example_section"
	cfgHost     = "go_db_mysql"
	cfgDatabase = "example_test"
	cfgUsername = "exampleuser"
	cfgPassword = "examplepasswd"
	cfgProtocol = "tcp"
	cfgPort     = 3306
	timeout     = 5
	querySelect = `SELECT id, name, status, created_at, updated_at FROM users WHERE id = :id AND status = :status;`
	queryUpdate = `UPDATE users SET status = :afterSts, updated_at = NOW() WHERE id = :id AND status = :beforeSts;`
)

var cfgPath = string(filepath.Separator) + filepath.Join("go", "src", "work", "testdata", "example", "example.dsn")

func TestNewDBContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	f := dbutil.NewConfigFile(cfgTyp, cfgPath, cfgSection)
	db, err := dbutil.NewDBContext(ctx, f)

	if err != nil {
		t.Error(err)
	}

	if err := db.PingContext(ctx); err != nil {
		t.Error(err)
	}
}

func TestParseConfig(t *testing.T) {
	cfg, err := dbutil.ParseConfig(cfgTyp, cfgPath, cfgSection)
	if err != nil {
		t.Error(err)
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
		t.Error(err)
	}

	db, err := dbutil.OpenContext(ctx, cfg)
	if err != nil {
		t.Error(err)
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
		t.Error(err)
	}
	args := map[string]any{"id": 1, "status": 0}
	list, err := dbutil.SelectContext[example.User](ctx, db, querySelect, args)
	if err != nil {
		t.Error(err)
	}

	if len(list) != 1 {
		t.Errorf("want: 1, got: %d", len(list))
	}
	if list[0].Status != 0 {
		t.Errorf("want: 1, got: %d", list[0].Status)
	}
}
