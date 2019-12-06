package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wiggin77/cfg"
)

const (
	DEFAULT_SORT = "artist:a,album:a,track:a"
)

// GetSortString returns the sort string for the specified directory.
// Sort string is determined by looking for `.mp3copy` files in the
// specified directory or parent directories. If no `.mp3copy` files
// are found then use command line options, and if none provided, use
// app default of "artist:a, album:a, track:a"
func GetSortString(opts Opts, dir string) string {
	config := &cfg.Config{}
	defer config.Shutdown()

	// walk the tree back to the original src dir and grab each
	// .mp3copy file. Add them to config in order of precedence.
	for {
		file := filepath.Join(dir, OPTS_FILE)
		src, err := cfg.NewSrcFileFromFilespec(file)
		if err == nil {
			m, err := src.GetProps()
			if err == nil {
				children, ok := m["children"]
				// don't add src if `children=false` and stop walking
				// since any opts in parent dirs are blocked.
				if ok && strings.ToLower(children) == "false" {
					break
				}
				config.AppendSource(src)
			}
		}
		if filepath.Clean(dir) == opts.src {
			break
		}
		dir = filepath.Dir(dir) // get parent dir
		if dir == "." || dir == "/" {
			break
		}
	}
	if opts.sort != "" {
		msrc := cfg.NewSrcMapFromMap(map[string]string{"sort": opts.sort})
		config.AppendSource(msrc)
	}

	s, _ := config.String("sort", DEFAULT_SORT)
	return strings.ReplaceAll(s, " ", "")
}

// ParseSortString parses a comma separated list of sort criteria of the form
// `album:a,artist:d,...`. A slice of Sorters is returned or error if the
// string is invalid.
func ParseSortString(ss string) ([]Sorter, error) {
	sorters := make([]Sorter, 0)

	// remove all whitespace
	ss = strings.ReplaceAll(ss, "\t", "")
	ss = strings.ReplaceAll(ss, " ", "")

	tokens := strings.Split(ss, ",")
	if len(tokens) == 0 {
		return sorters, nil
	}

	for _, t := range tokens {
		comp := strings.Split(t, ":")

		s := comp[0]
		a := "a"
		if len(comp) > 1 {
			a = comp[1]
		}
		sorter, err := sorterFromString(s, a)
		if err != nil {
			return sorters, err
		}
		sorters = append(sorters, sorter)
	}
	return sorters, nil
}

// sorterFromString parses the sort string components and returns a Sorter.
// ss can be one of:
//   `artist`      artist name
//   `album`       album name
//   `song`        song name
//   `track`       track id
//   `genre`       song's genre
//   `year`        year song was published
//   `file`        file name
//   `date`        last modified date
//   `random`      random order
// order can be one of:
//   `a`           ascending order
//   `d`           descending order
func sorterFromString(ss string, order string) (Sorter, error) {
	sorter := Sorter{}
	switch order {
	case "a":
	case "d":
		sorter.descending = true
	default:
		return sorter, fmt.Errorf("invalid sort order (must be ':a' or ':d'): %s:%s", ss, order)
	}

	ss = strings.ToLower(ss)
	switch ss {
	case "artist":
	case "album":
	case "song":
	case "track":
	case "genre":
	case "year":
	case "file":
	case "date":
	case "random":
	default:
		return sorter, fmt.Errorf("invalid sort string: %s:%s", ss, order)
	}
	sorter.field = ss
	return sorter, nil
}
