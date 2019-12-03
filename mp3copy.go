package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/wiggin77/mp3copy/sort"
)

func mp3copy(opts Opts) error {

	return nil
}

func copyDir(opts Opts, dir string) error {
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return err
	}

	// enumerates all directory entries and returns sorted slice.
	items, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	// separate directories from files.
	var dirs []os.FileInfo
	var entries []sort.Entry
	for _, fi := range items {
		if fi.IsDir() {
			dirs = append(dirs, fi)
		} else {
			e := sort.Entry{Filespec: filepath.Join(dir, fi.Name()), LastModified: fi.ModTime(), Filename: fi.Name()}
			entries = append(entries, e)
		}
	}

	reldir, err := filepath.Rel(opts.src, dir)
	if err != nil {
		return err
	}
	dest := filepath.Join(opts.dest, reldir)

	// copy child directories first
	for _, d := range dirs {
		child := filepath.Join(opts.dest, reldir, d.Name())
		if err := copyDir(opts, child); err != nil {
			return err
		}
	}

	// mkdir this directory in the dest
	if reldir != "." {
		if err := os.MkdirAll(dest, dirInfo.Mode().Perm()); err != nil {
			return err
		}
	}

	// sort and copy files.
	//err = SortEntries(entries)
	return fmt.Errorf("not implemented yet")
}

func copyFiles(opts Opts, entries []sort.Entry) error {
	return fmt.Errorf("not implemented yet")
}
