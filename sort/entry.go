package sort

import (
	"os"
	"time"

	"github.com/dhowden/tag"
)

type Info struct {
	Album  string
	Artist string
	Song   string
	Track  int
}

type Entry struct {
	Filespec     string
	Filename     string
	LastModified time.Time
	Info         *Info
}

// ExtractInfo reads the specified music file and extracts
// the metadata, such as ID3 fields.
func ExtractInfo(filespec string) (*Info, error) {
	file, err := os.Open(filespec)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	info := Info{}
	info.Artist = m.Artist()
	info.Album = m.Album()
	info.Song = m.Title()
	info.Track, _ = m.Track()
	return &info, nil
}
