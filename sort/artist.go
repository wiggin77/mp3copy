package sort

import (
	"fmt"
	"sort"
)

// Artist implements the Sorter interface and sorts by artist name.
type Artist struct {
	Descending bool
}

// Sort sorts the slice of Entry in place based on artist name.
func (a Artist) Sort(entries []Entry) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	sort.Slice(entries, func(i, j int) bool {
		var artist1, artist2 string
		var err error
		if artist1, err = extractArtist(entries[i]); err != nil {
			panic(err)
		}
		if artist2, err = extractArtist(entries[j]); err != nil {
			panic(err)
		}
		if a.Descending {
			return artist1 > artist2
		}
		return artist1 < artist2
	})
	return nil
}

func extractArtist(entry Entry) (string, error) {
	if entry.Info == nil {
		info, err := ExtractInfo(entry.Filespec)
		if err != nil {
			return "", err
		}
		entry.Info = info
	}
	return entry.Info.Artist, nil
}
