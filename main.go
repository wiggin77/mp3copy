package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
)

const (
	APP_VERSION = "1.01"
)

type Opts struct {
	src  string
	dest string
	sort string

	simulation bool
}

func main() {
	err := run()
	if err != nil {
		Term.Errorf("%v\n", err)
		os.Exit(-1)
	}
	Term.Printf("\n")
	os.Exit(0)
}

func run() error {
	opts := Opts{}
	var help bool
	var silent bool
	var version bool
	flag.StringVar(&opts.src, "src", "", "source directory")
	flag.StringVar(&opts.dest, "dest", "", "destination directory")
	flag.StringVar(&opts.sort, "sort", "",
		`default sort criteria, comma separated, 
		in order of precedence with optional order suffix. 
		(e.g. album:a,track:a) Used for any directories without .mp3copy file`)
	flag.BoolVar(&silent, "s", false, "silent, no progress displayed")
	flag.BoolVar(&opts.simulation, "sim", false, "simulation only; files read, nothing written")
	flag.BoolVar(&version, "version", false, "display version info")
	flag.BoolVar(&help, "h", false, "display usage help")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return nil
	}

	if version {
		fmt.Println(APP_VERSION)
		return nil
	}

	if !silent {
		Term.SetOut(os.Stdout)
		Term.SetErr(os.Stderr)
	}

	if opts.src == "" {
		flag.PrintDefaults()
		return fmt.Errorf("missing src")
	}
	if opts.dest == "" {
		flag.PrintDefaults()
		return fmt.Errorf("missing dest")
	}

	opts.src = expandTilde(opts.src)
	opts.dest = expandTilde(opts.dest)

	// check source exists and is a directory
	if fi, err := os.Stat(opts.src); err != nil || !fi.IsDir() {
		return fmt.Errorf("src %s is not a directory", opts.src)
	}

	// check dest exists and is a directory; create if needed
	if !opts.simulation {
		if fi, err := os.Stat(opts.dest); err != nil || !fi.IsDir() {
			if err := os.MkdirAll(opts.dest, fi.Mode().Perm()); err != nil {
				return fmt.Errorf("dest %s is not a directory", opts.dest)
			}
		}
	}
	return mp3copy(opts)
}

func expandTilde(s string) string {
	t, err := homedir.Expand(s)
	if err != nil {
		t = s
	}
	return t
}
