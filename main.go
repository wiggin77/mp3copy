package main

import (
	"flag"
	"fmt"
	"os"
)

type Opts struct {
	src  string
	dest string
	sort string
}

func main() {
	err := run()
	if err != nil {
		Term.Errorf("%v\n", err)
		os.Exit(-1)
	}
	os.Exit(0)
}

func run() error {
	opts := Opts{}
	var help bool
	flag.StringVar(&opts.src, "src", "", "source directory")
	flag.StringVar(&opts.dest, "dest", "", "destination directory")
	flag.StringVar(&opts.sort, "sort", "",
		`default sort criteria, comma separated, 
		in order of precedence with optional order suffix. 
		(e.g. album:a,track:a) Used for any directories without .mp3copy file`)
	flag.BoolVar(&help, "h", false, "display usage help")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return nil
	}

	if opts.src == "" {
		flag.PrintDefaults()
		return fmt.Errorf("missing src")
	}
	if opts.dest == "" {
		flag.PrintDefaults()
		return fmt.Errorf("missing dest")
	}

	// check source exists and is a directory
	if fi, err := os.Stat(opts.src); err != nil || !fi.IsDir() {
		return fmt.Errorf("src %s is not a directory", opts.src)
	}

	// check dest exists and is a directory
	if fi, err := os.Stat(opts.src); err != nil || !fi.IsDir() {
		return fmt.Errorf("dest %s is not a directory", opts.dest)
	}
	return mp3copy(opts)
}
