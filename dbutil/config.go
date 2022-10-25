package dbutil

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
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
	SSLMode  string // for PostgreSQL
	Driver   string
	DataSrc  string
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
		dsn, err := cfg.dataSrcMySQL()
		if err != nil {
			return nil, err
		}
		cfg.DataSrc = dsn
	case pgsqlDBType:
		cfg.Driver = pgsqlDriver
		cfg.DataSrc = cfg.dataSrcPgSQL()
	default:
		return nil, errors.New("unsupported DB type")
	}

	return cfg, nil
}

// dataSrcMySQL returns data source name for MySQL.
func (cfg *Config) dataSrcMySQL() (string, error) {
	jst, err := time.LoadLocation(cfg.Tz)
	if err != nil {
		return "", err
	}
	c := mysql.Config{
		DBName:    cfg.Database,
		User:      cfg.Username,
		Passwd:    cfg.Password,
		Addr:      fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Net:       cfg.Protocol,
		ParseTime: true,
		Collation: "utf8mb4_bin",
		Loc:       jst,
	}
	return c.FormatDSN(), nil
}

// dataSrcPgSQL returns data source name for PostgreSQL.
func (cfg *Config) dataSrcPgSQL() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode)
}
