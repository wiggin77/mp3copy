package sort

import (
	"fmt"
	"sort"
)

// Song implements the Sorter interface and sorts by song name.
type Song struct {
	Descending bool
}

// Sort sorts the slice of Entry in place based on song name.
func (a Song) Sort(entries []Entry) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	sort.Slice(entries, func(i, j int) bool {
		var song1, song2 string
		var err error
		if song1, err = extractSong(entries[i]); err != nil {
			panic(err)
		}
		if song2, err = extractSong(entries[j]); err != nil {
			panic(err)
		}
		if a.Descending {
			return song1 > song2
		}
		return song1 < song2
	})
	return nil
}

func extractSong(entry Entry) (string, error) {
	if entry.Info == nil {
		info, err := ExtractInfo(entry.Filespec)
		if err != nil {
			return "", err
		}
		entry.Info = info
	}
	return entry.Info.Song, nil
}
