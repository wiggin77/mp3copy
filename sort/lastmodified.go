package sort

import (
	"fmt"
	"os"
	"sort"
	"time"
)

// LastModified implements the Sorter interface and sorts by last modified date.
type LastModified struct {
	Descending bool
}

// Sort sorts the slice of Entry in place based on last modified date.
func (l LastModified) Sort(entries []Entry) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	sort.Slice(entries, func(i, j int) bool {
		var lm1, lm2 time.Time
		var err error
		if lm1, err = extractLastModified(entries[i]); err != nil {
			panic(err)
		}
		if lm2, err = extractLastModified(entries[j]); err != nil {
			panic(err)
		}
		if l.Descending {
			return lm1.After(lm2)
		}
		return lm1.Before(lm2)
	})
	return nil
}

func extractLastModified(entry Entry) (time.Time, error) {
	if entry.LastModified.IsZero() {
		fi, err := os.Stat(entry.Filespec)
		if err != nil {
			return time.Time{}, fmt.Errorf("cannot get last modified time for %s", entry.Filespec)
		}
		entry.LastModified = fi.ModTime()
	}
	return entry.LastModified, nil
}
