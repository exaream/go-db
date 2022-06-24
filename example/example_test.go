package example_test

import (
	"fmt"
	"os"
	"testing"
)

// TODO: How to call a helper func in TestMain which does NOT have testing.T.
func TestMain(m *testing.M) {
	if err := initTable(confType, confDir, confStem, confSection); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	code := m.Run()

	os.Exit(code)
}
