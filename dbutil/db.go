package dbutil

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const (
	Tz     = "Asia/Tokyo"          // Timezone
	YmdHis = "2006-01-02 15:04:05" // Layout of "Y-m-d H:i:s"
)

// DB config
type Config struct {
	Host     string
	Database string
	Username string
	Password string
	Protocol string
	Tz       string
	Port     uint16 // 1~65535
}

// ParseConfig returns DB config by a config file.
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
		cfg.Tz = Tz
	}

	return cfg, nil
}

// OpenContext returns DB handle.
func OpenContext(ctx context.Context, cfg *Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
		cfg.Username, cfg.Password, cfg.Protocol, cfg.Host, cfg.Port, cfg.Database)

	params := url.Values{"parseTime": {"true"}, "loc": {cfg.Tz},
		// See: http://dsas.blog.klab.org/archives/52191467.html
		"interpolateParams": {"true"}, "collation": {"utf8mb4_bin"}}

	db, err := sqlx.Open("mysql", dsn+"?"+params.Encode())
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
