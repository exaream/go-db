package dbutil_test

import (
	"testing"

	"github.com/exaream/go-db/dbutil"
	"github.com/google/go-cmp/cmp"
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

			if diff := cmp.Diff(want, got); diff != "" {
				t.Error(diff)
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
