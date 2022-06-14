package main

import (
	"fmt"
	"os"

	"go.uber.org/multierr"

	s "ops/sample"
)

func main() {
	c := s.NewCond(s.UserId(), s.Status(), s.IniPath(), s.Section())

	if rerr := c.Run(); rerr != nil {
		for _, err := range multierr.Errors(rerr) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
