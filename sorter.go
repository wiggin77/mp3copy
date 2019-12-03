package main

import (
	"fmt"
	"strings"

	"github.com/wiggin77/mp3copy/sort"
)

// Sorter can sort a slice of Entry based on a field of Entry (e.g. Entry.album)
type Sorter interface {
	// Sort sorts the slice of filespecs in place.
	Sort(entries []sort.Entry) error
}

// SortEntries sorts the entries using the specified sort criteria.
func SortEntries(entries []sort.Entry, sort string) error {
	sorters, err := parseSort(sort)
	if err != nil {
		return err
	}
	return sortEntries(entries, sorters)
}

// sortEntries sorts the entries using the specified Sorter's.
func sortEntries(entries []sort.Entry, sorters []Sorter) error {
	if len(sorters) == 0 {
		return nil
	}

	err := sorters[0].Sort(entries)
	if err != nil {
		return err
	}

	// check for any repeating
	return fmt.Errorf("not implemented yet")
}

// parseSort parses a comma separated list of sort criteria of the form
// `album:a,artist:d,...`. A slice of Sorters is returned or error if the
// string is invalid.
func parseSort(sort string) ([]Sorter, error) {
	// remove all whitespace
	sort = strings.ReplaceAll(sort, " ", "")

	tokens := strings.Split(sort, ",")
	if len(tokens) == 0 {
		return defaultSorters(), nil
	}
	return nil, fmt.Errorf("not implemented yet")
}

// DefaultSorters returns the default slice of sorters
// (artist, album, track) all using ascending order.
func defaultSorters() []Sorter {
	return []Sorter{
		sort.Artist{}, sort.Album{}, sort.Track{},
	}
}
