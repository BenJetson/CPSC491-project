package demo

import (
	"fmt"
	"strconv"
)

type Bar struct {
	ID   int    `json:"id" db:"bar_id"`
	Name string `json:"name" db:"name"`
}

// Bad is a sentinel function that shall intentionally trigger golangci-lint.
// FIXME delete this once testing is done.
func Bad(s string) {
	x, _ := strconv.Atoi(s)
	fmt.Println(x)
}
