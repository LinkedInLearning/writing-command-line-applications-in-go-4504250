package main

import (
	"bytes"
	"fmt"
)

func ExampleBanner() {
	var buf bytes.Buffer
	Banner(&buf, "Go", 6)
	fmt.Println(buf.String())

	// Output:
	//   Go
	// ======
}
