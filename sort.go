package main

import (
	"math/rand"
	"sort"
	"time"
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
	if sorter.field == RANDOM {
		ShuffleEntries(entries)
		return
	}

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

// ShuffleEntries randomizes the order of entries in place.
func ShuffleEntries(entries []Entry) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	size := len(entries)
	var tmp Entry
	for i := range entries {
		j := r.Intn(size)
		tmp = entries[i]
		entries[i] = entries[j]
		entries[j] = tmp
	}
}
