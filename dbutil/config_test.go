package dbutil_test

import (
	"testing"

	"github.com/exaream/go-db/dbutil"
)

func TestParseConfig(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		dbType string
		path   string
	}{
		"mysql": {mysqlDBType, mysqlCfgPath},
		"pgsql": {pgsqlDBType, pgsqlCfgPath},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			want := expectedConfig(t, tt.dbType)
			got, err := dbutil.ParseConfig(cfgType, tt.path, cfgSection)
			if err != nil {
				t.Fatal(err)
			}

			if got.Host != want.Host {
				t.Fatalf("host want: %s, got: %s", want.Host, got.Host)
			}
			if got.Database != want.Database {
				t.Fatalf("database want: %s, got: %s", want.Database, got.Database)
			}
			if got.Username != want.Username {
				t.Fatalf("username want: %s, got: %s", want.Username, got.Username)
			}
			if got.Password != want.Password {
				t.Fatalf("password want: %s, got: %s", want.Password, got.Password)
			}
			if got.Protocol != want.Protocol {
				t.Fatalf("protocol want: %s, got: %s", want.Protocol, got.Protocol)
			}
			if got.Port != want.Port {
				t.Fatalf("port want: %d, got: %d", want.Port, got.Port)
			}
			if got.Tz != want.Tz {
				t.Fatalf("timezone want: %s, got: %s", want.Tz, got.Tz)
			}
		})
	}
}

func TestParseConfigErr(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		typ     string
		path    string
		section string
	}{
		"all":            {dummy, dummy, dummy},
		"path":           {cfgType, dummy, cfgSection},
		"type(mysql)":    {dummy, mysqlCfgPath, cfgSection},
		"type(pgsql)":    {dummy, pgsqlCfgPath, cfgSection},
		"section(mysql)": {cfgType, mysqlCfgPath, dummy},
		"section(pgsql)": {cfgType, pgsqlCfgPath, dummy},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := dbutil.ParseConfig(tt.typ, tt.path, tt.section)
			if err == nil {
				t.Error("want: error, got: nil")
			}
		})
	}
}
