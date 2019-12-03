package sort

import (
	"fmt"
	"sort"
)

// Album implements the Sorter interface and sorts by album name.
type Album struct {
	Descending bool
}

// Sort sorts the slice of Entry in place based on album name.
func (a Album) Sort(entries []Entry) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	sort.Slice(entries, func(i, j int) bool {
		var album1, album2 string
		var err error
		if album1, err = extractAlbum(entries[i]); err != nil {
			panic(err)
		}
		if album2, err = extractAlbum(entries[j]); err != nil {
			panic(err)
		}
		if a.Descending {
			return album1 > album2
		}
		return album1 < album2
	})
	return nil
}

func extractAlbum(entry Entry) (string, error) {
	if entry.Info == nil {
		info, err := ExtractInfo(entry.Filespec)
		if err != nil {
			return "", err
		}
		entry.Info = info
	}
	return entry.Info.Album, nil
}
