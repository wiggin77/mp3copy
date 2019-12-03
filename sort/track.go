package sort

import (
	"fmt"
	"sort"
)

// Track implements the Sorter interface and sorts by track name.
type Track struct {
	Descending bool
}

// Sort sorts the slice of Entry in place based on track name.
func (a Track) Sort(entries []Entry) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	sort.Slice(entries, func(i, j int) bool {
		var track1, track2 int
		var err error
		if track1, err = extractTrack(entries[i]); err != nil {
			panic(err)
		}
		if track2, err = extractTrack(entries[j]); err != nil {
			panic(err)
		}
		if a.Descending {
			return track1 > track2
		}
		return track1 < track2
	})
	return nil
}

func extractTrack(entry Entry) (int, error) {
	if entry.Info == nil {
		info, err := ExtractInfo(entry.Filespec)
		if err != nil {
			return 0, err
		}
		entry.Info = info
	}
	return entry.Info.Track, nil
}
