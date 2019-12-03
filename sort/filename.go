package sort

import (
	"path/filepath"
	"sort"
)

// Filename implements the Sorter interface and sorts by filename.
type Filename struct {
	Descending bool
}

// Sort sorts the slice of Entry in place based on filename.
func (f Filename) Sort(entries []Entry) error {
	sort.Slice(entries, func(i, j int) bool {
		fn1 := extractFilename(entries[i])
		fn2 := extractFilename(entries[j])
		if f.Descending {
			return fn1 > fn2
		}
		return fn1 < fn2
	})
	return nil
}

func extractFilename(entry Entry) string {
	if entry.Filename == "" {
		entry.Filename = filepath.Base(entry.Filespec)
	}
	return entry.Filename
}
