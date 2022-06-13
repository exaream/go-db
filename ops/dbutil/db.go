package dbutil

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"gopkg.in/ini.v1"
)

type Conf struct {
	Host     string
	DB       string
	Username string
	Password string
	Protocol string
	Port     int
	Tz       string
}

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

func (c *Conf) Open() (*sql.DB, error) {
	srcName := fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
		c.Username, c.Password, c.Protocol, c.Host, strconv.Itoa(c.Port), c.DB)

	params := url.Values{"parseTime": {"true"}, "loc": {c.Tz}}

	return sql.Open("mysql", srcName+"?"+params.Encode())
}
