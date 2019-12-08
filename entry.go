package main

import (
	"os"
	"strconv"
	"time"

	"github.com/dhowden/tag"
)

const (
	FILENAME     = "file"
	LASTMODIFIED = "date"
	ARTIST       = "artist"
	ALBUM        = "album"
	SONG         = "song"
	TRACK        = "track"
	GENRE        = "genre"
	YEAR         = "year"
	RANDOM       = "random"
	FIELD_COUNT  = 9

	TIME_FORMAT = "2006-01-02T15:04:05.000"
)

type Entry struct {
	filespec string
	mode     os.FileMode
	music    bool
	fields   map[string]string
}

func NewEntry(filespec string) (Entry, error) {
	entry := Entry{filespec: filespec, fields: make(map[string]string, FIELD_COUNT)}
	fi, err := os.Stat(filespec)
	if err != nil {
		return Entry{}, err
	}

	entry.mode = fi.Mode()
	entry.fields[FILENAME] = fi.Name()
	entry.fields[LASTMODIFIED] = timeToString(fi.ModTime())
	err = extractMeta(entry)
	if err == nil {
		entry.music = true
	}
	return entry, nil
}

// extractInfo reads the specified music file and extracts
// the metadata, such as ID3 fields.
func extractMeta(entry Entry) error {
	file, err := os.Open(entry.filespec)
	if err != nil {
		return err
	}
	defer file.Close()

	m, err := tag.ReadFrom(file)
	if err != nil {
		return err
	}

	entry.fields[ARTIST] = m.Artist()
	entry.fields[ALBUM] = m.Album()
	entry.fields[SONG] = m.Title()
	track, _ := m.Track()
	entry.fields[TRACK] = trackToString(track)
	entry.fields[GENRE] = m.Genre()
	entry.fields[YEAR] = strconv.Itoa(m.Year())
	return nil
}

// timeToString converts a time.Time to a string that can be lexically sorted.
func timeToString(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

// trackToString converts a track index to a string that can be lexically sorted.
func trackToString(track int) string {
	s := "00000" + strconv.Itoa(track)
	return s[len(s)-5:]
}
