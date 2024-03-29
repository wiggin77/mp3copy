package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	OPTS_FILE = ".mp3copy"
)

func mp3copy(opts Opts) error {
	return copyDir(opts, opts.src)
}

func copyDir(opts Opts, dir string) error {
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}

	var dirs []os.FileInfo
	var entries []Entry
	dirs, entries, err = getDirectoriesAndFiles(dir)
	if err != nil {
		return err
	}

	// calculate this directory relative to original source dir
	reldir, err := filepath.Rel(opts.src, dir)
	if err != nil {
		return err
	}
	// calculate what this directory would be in the target tree
	dest := filepath.Join(opts.dest, reldir)

	// mkdir this directory in the dest
	if !opts.simulation && reldir != "." {
		if err := os.MkdirAll(dest, dirInfo.Mode().Perm()); err != nil {
			return err
		}
	}

	// copy child directories first
	for _, d := range dirs {
		child := filepath.Join(opts.src, reldir, d.Name())
		if err := copyDir(opts, child); err != nil {
			return err
		}
	}

	sorters, err := getSorters(opts, dir)
	if err != nil {
		return err
	}

	SortEntries(entries, sorters)
	return copyFiles(opts, entries)
}

// getDirectoriesAndFiles returns all the directories and files within the specified directory.
func getDirectoriesAndFiles(dir string) (dirs []os.FileInfo, entries []Entry, err error) {
	items, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}

	for _, fi := range items {
		if fi.IsDir() {
			dirs = append(dirs, fi)
		} else {
			filespec := filepath.Join(dir, fi.Name())
			entry, err := NewEntry(filespec)
			if err != nil {
				return nil, nil, fmt.Errorf("error opening %s: %v", filespec, err)
			}
			entries = append(entries, entry)
		}
	}
	return dirs, entries, nil
}

func getSorters(opts Opts, dir string) ([]Sorter, error) {
	ss := GetSortString(opts, dir)
	return ParseSortString(ss)
}

func copyFiles(opts Opts, entries []Entry) error {
	buf := make([]byte, 64*1024)
	for _, entry := range entries {
		// calculate this file's directory relative to original source dir
		reldir, err := filepath.Rel(opts.src, entry.filespec)
		if err != nil {
			return err
		}
		dest := filepath.Join(opts.dest, reldir)
		err = copyFile(opts, entry, dest, buf)
		if err != nil {
			Term.Errorf("%v", err)
		}
	}
	return nil
}

func copyFile(opts Opts, entry Entry, dest string, buf []byte) error {
	srcFile, err := os.Open(entry.filespec)
	if err != nil {
		return fmt.Errorf("cannot open %s: %v", entry.filespec, err)
	}
	defer srcFile.Close()

	var destFile *os.File
	if !opts.simulation {
		destFile, err = os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, entry.mode)
		if err != nil {
			return fmt.Errorf("cannot create %s: %v", dest, err)
		}
		defer destFile.Close()
	}

	Term.Printf("%s\n", entry.filespec)

	var count int64
	var n int
	for {
		n, err = srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if !opts.simulation {
			if n, err = destFile.Write(buf[:n]); err != nil {
				return err
			}
		}
		count = count + int64(n)
		Term.Progress(count, entry.size)
	}
	Term.Progress(0, 0)
	return nil
}
