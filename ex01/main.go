package main

import (
	"convert/imgconv"
	"flag"
	"fmt"
	"os"
)

func main() {
	// imgconv.Convert()
	flag.Parse()
	if dirname := flag.Arg(0); dirname == "" {
		fmt.Fprintf(os.Stderr, "error: invalid argument\n")
		os.Exit(0)
	} else if flag.Arg(1) != "" {
		fmt.Fprintf(os.Stderr, "error: multiple arguments\n")
		os.Exit(0)
	} else {
		imgconv.ConvertMain(dirname)
	}
}
