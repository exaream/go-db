package example_test

import (
	"fmt"
	"os"
	"testing"
)

// TODO: Confirm how to call a helper func in TestMain which does NOT have testing.T.
func TestMain(m *testing.M) {
	if err := initTable(testIniPath, testSection); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	code := m.Run()

	os.Exit(code)
}

/*
func TestRun(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		id     int
		name   string
		status int
	}{
		"change status": {1, "Alice", 1},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := example.NewCond(tt.id, tt.status, testIniPath, testSection, testTimeout)
			if err := c.Run(); err != nil {
				t.Error(err)
			}

		})
	}
}
*/
