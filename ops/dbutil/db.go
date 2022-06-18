package dbutil

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"go.uber.org/multierr"
	"gopkg.in/ini.v1"
)

const YmdHis = "2006-01-02 15:04:05" // layout of "Y-m-d H:i:s"

type Conf struct {
	Host     string
	DB       string
	Username string
	Password string
	Protocol string
	Tz       string
	Port     int
}

// OpenByConf returns a DB handle by a config file's path.
func OpenByConf(iniPath, section string) (*sql.DB, error) {
	conf, err := ParseConf(iniPath, section)
	if err != nil {
		return nil, err
	}

	db, err := Open(conf)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// ParseConf returns DB's config info.
func ParseConf(iniPath, section string) (*Conf, error) {
	iniFile, err := ini.Load(iniPath)
	if err != nil {
		return nil, errors.New("faild to load a DSN file")
	}

	sec := iniFile.Section(section)

	encodedPwd := sec.Key("password").String()
	decodedPwd, err := base64.StdEncoding.DecodeString(encodedPwd)
	if err != nil {
		return nil, errors.New("failed to decode DB password")
	}

	port, err := sec.Key("port").Int()
	if err != nil {
		return nil, errors.New("failed to get port")
	}

	return &Conf{
		Host:     sec.Key("host").String(),
		DB:       sec.Key("database").String(),
		Username: sec.Key("username").String(),
		Password: string(decodedPwd),
		Protocol: sec.Key("protocol").String(),
		Port:     port,
		Tz:       "Asia/Tokyo",
	}, nil
}

// Open returns a DB handle.
// TODO: Is it better to use a receiver?
func Open(c *Conf) (*sql.DB, error) {
	srcName := fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
		c.Username, c.Password, c.Protocol, c.Host, strconv.Itoa(c.Port), c.DB)

	params := url.Values{"parseTime": {"true"}, "loc": {c.Tz}}

	return sql.Open("mysql", srcName+"?"+params.Encode())
}

// Rollback rollbacks using transaction.
// It can return multiple errors.
func Rollback(tx *sql.Tx, rerr, err error) error {
	rerr = multierr.Append(rerr, err)
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return multierr.Append(rerr, rollbackErr)
	}
	return rerr
}
