// A utility package for DB operations.
package dbutil

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

const (
	defaultTz = "Asia/Tokyo"
	// MySQL
	mysqlDBType = "mysql"
	mysqlDriver = "mysql"
	// PostgreSQL
	pgsqlDBType = "pgsql"
	pgsqlDriver = "pgx"
)

// DB config file
type ConfigFile struct {
	Type    string
	Path    string
	Section string
}

// DB config
type Config struct {
	Type     string
	Host     string
	Database string
	Username string
	Password string
	Port     uint16 // 1~65535
	Protocol string
	Tz       string
	Driver   string
	DataSrc  string
}

// NewDBContext returns DB handle.
func NewDBContext(ctx context.Context, f *ConfigFile) (*sqlx.DB, error) {
	cfg, err := ParseConfig(f.Type, f.Path, f.Section)
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
		Type:    typ,
		Path:    path,
		Section: section,
	}
}

// ParseConfig returns DB config by DB config file.
func ParseConfig(typ, path, section string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	v := viper.New()
	v.SetConfigType(typ)
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	sub := v.Sub(section)
	if sub == nil {
		return nil, errors.New("failed to parse config by section")
	}

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

	switch cfg.Type {
	case mysqlDBType:
		cfg.Driver = mysqlDriver
		cfg.DataSrc = cfg.dataSrcMySQL()
	case pgsqlDBType:
		cfg.Driver = pgsqlDriver
		cfg.DataSrc = cfg.dataSrcPgSQL()
	default:
		return nil, errors.New("Unsupported DB type")
	}

	return cfg, nil
}

// dataSrcMySQL returns data source name for MySQL.
func (cfg *Config) dataSrcMySQL() string {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
		cfg.Username, cfg.Password, cfg.Protocol, cfg.Host, cfg.Port, cfg.Database)

	params := url.Values{"parseTime": {"true"}, "loc": {cfg.Tz},
		"interpolateParams": {"true"}, "collation": {"utf8mb4_bin"}}

	return dsn + "?" + params.Encode()
}

// dataSrcPgSQL returns data source name for PostgreSQL.
func (cfg *Config) dataSrcPgSQL() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", // TODO: sslmode
		cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password)
}

// OpenContext returns DB handle.
// See: http://dsas.blog.klab.org/archives/52191467.html
func OpenContext(ctx context.Context, cfg *Config) (db *sqlx.DB, err error) {
	db, err = sqlx.Open(cfg.Driver, cfg.DataSrc)

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
