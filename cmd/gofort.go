package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kdevo/gofort/pkg/fortune"
)

// Common program vars overwritten by linker:
var (
	NAME    = "gofort"
	VERSION = "v0.1.0"
	COMMIT  = "N/A"
	OS      = "os"
	ARCH    = "arch"
	DATE    = "2021-08-22T22:00:00"
)

func main() {
	var (
		version = flag.Bool("version", false, "Show version information and exit.")
	)
	flag.Parse()

	if *version {
		fmt.Printf("%s %s %s/%s, built at %s (%s)\n", NAME, VERSION, OS, ARCH, DATE, COMMIT)
		os.Exit(0)
	}

	fmt.Println(fortune.New().Fortune())
}
