package dbutil

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

// Timezone
const defaultTz = "Asia/Tokyo"

// DB config file
type ConfigFile struct {
	Typ     string
	Path    string
	Section string
}

// DB config
type Config struct {
	Driver   string
	Host     string
	Database string
	Username string
	Password string
	Protocol string
	Tz       string
	Port     uint16 // 1~65535
}

// NewDBContext returns DB handle.
func NewDBContext(ctx context.Context, f *ConfigFile) (*sqlx.DB, error) {
	cfg, err := ParseConfig(f.Typ, f.Path, f.Section)
	if err != nil {
		return nil, err
	}

	db, err := OpenContext(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewConfigFile returns DB config file.
func NewConfigFile(typ, path, section string) *ConfigFile {
	return &ConfigFile{
		Typ:     typ,
		Path:    path,
		Section: section,
	}
}

// ParseConfig returns DB config by DB config file.
func ParseConfig(typ, path, section string) (*Config, error) {
	v := viper.New()
	v.SetConfigType(typ)
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	sub := v.Sub(section)
	var cfg *Config
	if err := sub.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	password, err := base64.StdEncoding.DecodeString(cfg.Password)
	if err != nil {
		return nil, err
	}
	cfg.Password = string(password)

	if cfg.Tz == "" {
		cfg.Tz = defaultTz
	}

	return cfg, nil
}

func OpenMySQL(cfg *Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
		cfg.Username, cfg.Password, cfg.Protocol, cfg.Host, cfg.Port, cfg.Database)

	params := url.Values{"parseTime": {"true"}, "loc": {cfg.Tz},
		"interpolateParams": {"true"}, "collation": {"utf8mb4_bin"}}

	db, err := sqlx.Open("mysql", dsn+"?"+params.Encode())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func OpenPostgreSQL(cfg *Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", // TODO: sslmode
		cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// OpenContext returns DB handle.
// See: http://dsas.blog.klab.org/archives/52191467.html
func OpenContext(ctx context.Context, cfg *Config) (db *sqlx.DB, err error) {
	switch cfg.Driver {
	case "mysql":
		db, err = OpenMySQL(cfg)
	case "postgres":
		db, err = OpenPostgreSQL(cfg)
	}

	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

// SelectContext runs SELECT and returns the results.
func SelectContext[T any](ctx context.Context, db *sqlx.DB, query string, args map[string]any) ([]T, error) {
	rows, err := sqlx.NamedQueryContext(ctx, db, query, args)
	if err != nil {
		return nil, err
	}

	list := []T{}
	for rows.Next() {
		var row T
		if err := rows.StructScan(&row); err != nil {
			return nil, err
		}
		list = append(list, row)
		fmt.Println(row)
	}

	return list, nil
}

// SelectTxContext runs SELECT and returns the results on transaction.
func SelectTxContext[T any](ctx context.Context, tx *sqlx.Tx, query string, args map[string]any) ([]T, error) {
	rows, err := sqlx.NamedQueryContext(ctx, tx, query, args)
	if err != nil {
		return nil, err
	}

	list := []T{}
	for rows.Next() {
		var row T
		if err := rows.StructScan(&row); err != nil {
			return nil, multierr.Append(err, tx.Rollback())
		}
		list = append(list, row)
		fmt.Println(row)
	}

	return list, nil
}

// UpdateTxContext runs UPDATE on transaction.
func UpdateTxContext(ctx context.Context, tx *sqlx.Tx, query string, args map[string]any) (int64, error) {
	result, err := sqlx.NamedExecContext(ctx, tx, query, args)
	if err != nil {
		return 0, multierr.Append(err, tx.Rollback())
	}

	num, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return num, nil
}

// BulkInsertTxContext executes Bulk Insert on context and transaction.
// TODO: Too many arguments?
func BulkInsertTxContext[T any](ctx context.Context, tx *sqlx.Tx, fn func(i, j uint) []T, query string, min, max, chunkSize uint) (int64, error) {
	var i uint
	var total int64

	for i = min; i <= max; i += chunkSize {
		j := i + chunkSize - min
		if j > max {
			j = max
		}

		result, err := tx.NamedExecContext(ctx, query, fn(i, j))
		if err != nil {
			return 0, multierr.Append(err, tx.Rollback())
		}

		num, err := result.RowsAffected()
		if err != nil {
			return 0, multierr.Append(err, tx.Rollback())
		}
		total += num
	}

	return total, nil
}
