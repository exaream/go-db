package dbutil

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const (
	Tz     = "Asia/Tokyo"          // Timezone
	YmdHis = "2006-01-02 15:04:05" // Layout of "Y-m-d H:i:s"
)

// DB config
type Conf struct {
	Host     string
	Database string
	Username string
	Password string
	Protocol string
	Tz       string
	Port     int
}

// ParseConf returns DB config by a config file.
func ParseConf(typ, path, section string) (*Conf, error) {
	v := viper.New()
	v.SetConfigType(typ)
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	sub := v.Sub(section)
	var conf *Conf
	if err := sub.Unmarshal(&conf); err != nil {
		return nil, err
	}

	password, err := base64.StdEncoding.DecodeString(conf.Password)
	if err != nil {
		return nil, err
	}
	conf.Password = string(password)

	if conf.Tz == "" {
		conf.Tz = Tz
	}

	return conf, nil
}

// OpenWithContext returns DB handle.
func OpenWithContext(ctx context.Context, conf *Conf) (*sqlx.DB, error) {
	srcName := fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
		conf.Username, conf.Password, conf.Protocol, conf.Host, strconv.Itoa(conf.Port), conf.Database)

	params := url.Values{"parseTime": {"true"}, "loc": {conf.Tz}}

	db, err := sqlx.Open("mysql", srcName+"?"+params.Encode())
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
