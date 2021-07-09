package main

import (
	"fmt"
	"os"

	chesseract "github.com/thijzert/chesseract/pkg"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Printf("Package version: %s. Hello, world!\n", chesseract.PackageVersion)
	return nil
}
