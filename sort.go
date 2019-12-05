package main

import (
	"sort"
)

type Sorter struct {
	field      string
	descending bool
}

// SortEntries sorts the entries using the specified Sorter's.
// Will sub-sort the entries based on the supplied slice of
// sorters such that the first sort is done with sorter[0],
// the sub-sort is done with sorter[1], ...
func SortEntries(entries []Entry, sorters []Sorter) {
	if len(sorters) == 0 {
		return
	}

	sorter := sorters[0]

	sort.Slice(entries, func(i, j int) bool {
		field1 := entries[i].fields[sorter.field]
		field2 := entries[j].fields[sorter.field]
		if sorter.descending {
			return field1 > field2
		}
		return field1 < field2
	})

	// if no subsorters then nothing left to do
	if len(sorters) < 2 {
		return
	}

	// look for any repeating sequences and sub-sort them.
	var prev string
	var found bool
	var start int
	for i, entry := range entries {
		val := entry.fields[sorter.field]
		if i == 0 {
			prev = val
			continue
		}

		if val == prev {
			if !found {
				found = true
				start = i
			}
		} else {
			if found {
				found = false
				SortEntries(entries[start:i], sorters[1:])
			}
		}
		prev = val
	}
}
