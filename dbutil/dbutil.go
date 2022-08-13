// Package dbutil provides general purpose functions for operating MySQL and PostgreSQL.
package dbutil

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
)

type stringConstant string

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
func SelectContext[T any](ctx context.Context, db *sqlx.DB, query stringConstant, args map[string]any) ([]*T, error) {
	rows, err := sqlx.NamedQueryContext(ctx, db, string(query), args)
	if err != nil {
		return nil, err
	}

	var list []*T
	for rows.Next() {
		var row T
		if err := rows.StructScan(&row); err != nil {
			return nil, err
		}
		list = append(list, &row)
		fmt.Println(row)
	}

	return list, nil
}

// SelectTxContext runs SELECT and returns the results on transaction.
func SelectTxContext[T any](ctx context.Context, tx *sqlx.Tx, query stringConstant, args map[string]any) ([]*T, error) {
	rows, err := sqlx.NamedQueryContext(ctx, tx, string(query), args)
	if err != nil {
		return nil, err
	}

	var list []*T
	for rows.Next() {
		var row T
		if err := rows.StructScan(&row); err != nil {
			return nil, multierr.Append(err, tx.Rollback())
		}
		list = append(list, &row)
		fmt.Println(row)
	}

	return list, nil
}

// UpdateTxContext runs UPDATE on transaction.
func UpdateTxContext(ctx context.Context, tx *sqlx.Tx, query stringConstant, args map[string]any) (int64, error) {
	result, err := sqlx.NamedExecContext(ctx, tx, string(query), args)
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
func BulkInsertTxContext[T any](ctx context.Context, tx *sqlx.Tx,
	fn func(i, j uint) []*T, query stringConstant, min, max, chunkSize uint) (int64, error) {
	var i uint
	var total int64

	queryStr := string(query)
	for i = min; i <= max; i += chunkSize {
		j := i + chunkSize - min
		if j > max {
			j = max
		}

		result, err := tx.NamedExecContext(ctx, queryStr, fn(i, j))
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
