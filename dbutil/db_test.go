package dbutil_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/exaream/go-db/dbutil"
)

const (
	cfgTyp     = "ini"
	cfgSection = "example_section"
	timeout    = 5
)

var cfgPath = string(filepath.Separator) + filepath.Join("go", "src", "work", "cmd", "example", "example.dsn")

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
