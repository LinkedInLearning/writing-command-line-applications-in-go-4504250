package main

import (
	"fmt"
	"strconv"
)

// Count is a positive number flag.Value.
type Count struct {
	val *int
}

// String returns the string representation of the current value.
func (c *Count) String() string {
	if c.val == nil {
		return ""
	}

	return fmt.Sprintf("%d", *c.val)
}

// Set sets the value from command line.
func (c *Count) Set(val string) error {
	n, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("bad number: %s", err)
	}

	if n < 0 {
		return fmt.Errorf("expected positive number, got %d", n)
	}

	*c.val = n
	return nil
}
