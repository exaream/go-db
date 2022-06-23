package dbx

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const (
	DefaultTz = "Asia/Tokyo"
	LF        = "\n"
	YmdHis    = "2006-01-02 15:04:05" // layout of "Y-m-d H:i:s"
)

type Conf struct {
	Host     string
	Database string
	Username string
	Password string
	Protocol string
	Tz       string
	Port     int
}

func OpenWithContext(ctx context.Context, typ, dir, stem, section string) (*sqlx.DB, error) {
	c, err := ParseConf(typ, dir, stem, section)
	if err != nil {
		return nil, err
	}

	srcName := fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
		c.Username, c.Password, c.Protocol, c.Host, strconv.Itoa(c.Port), c.Database)

	params := url.Values{"parseTime": {"true"}, "loc": {c.Tz}}

	db, err := sqlx.Open("mysql", srcName+"?"+params.Encode())
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func ParseConf(typ, confPath, stem, section string) (*Conf, error) {
	v := viper.New()
	v.SetConfigType(typ)
	v.AddConfigPath(confPath)
	v.SetConfigName(stem) // the file stem (= the file name without the extension)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	sub := v.Sub(section)
	var c *Conf
	if err := sub.Unmarshal(&c); err != nil {
		return nil, err
	}

	password, err := base64.StdEncoding.DecodeString(c.Password)
	if err != nil {
		return nil, err
	}
	c.Password = string(password)

	if c.Tz == "" {
		c.Tz = DefaultTz
	}

	return c, nil
}

type Records map[int]map[string]any

func QueryTxWithContext(ctx context.Context, tx *sql.Tx, stmt string, fn func(context.Context, *sql.Rows) (Records, error)) (Records, error) {
	rows, err := tx.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	return fn(ctx, rows)
}

func QueryWithContext(ctx context.Context, db *sqlx.DB, stmt string, fn func(context.Context, *sql.Rows) (Records, error)) (Records, error) {
	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	return fn(ctx, rows)
}
