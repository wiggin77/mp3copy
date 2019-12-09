package main

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	rnd    *rand.Rand
	fields []string
)

func setup(t *testing.T) {
	src := rand.NewSource(time.Now().UnixNano())
	rnd = rand.New(src)

	fields = []string{FILENAME, LASTMODIFIED, ARTIST, ALBUM, SONG, TRACK, GENRE, YEAR, RANDOM}
	require.Equal(t, FIELD_COUNT, len(fields), "field(s) added/removed, unit test needs updating")
}

func teardown(t *testing.T) {
	// do nothing
}

func TestSortEntries(t *testing.T) {
	setup(t)
	defer teardown(t)

	type args struct {
		entries []Entry
		sorters []Sorter
	}
	type test struct {
		name string
		args args
	}
	tests := []test{
		{name: "case1", args: args{entries: genEntries(500), sorters: []Sorter{}}},
		{name: "case2", args: args{entries: genEntries(500), sorters: []Sorter{
			{field: FILENAME, descending: false},
		}}},
		{name: "case3", args: args{entries: genEntries(500), sorters: []Sorter{
			{field: ARTIST, descending: false},
			{field: ALBUM, descending: false},
			{field: TRACK, descending: false},
		}}},
		{name: "case4", args: args{entries: genEntries(500), sorters: []Sorter{
			{field: RANDOM, descending: false},
		}}},
	}
	// Add some random test cases
	for i := 0; i < 500; i++ {
		tests = append(tests, test{
			name: "case_rnd" + strconv.Itoa(i),
			args: args{
				entries: genRandomEntries(500),
				sorters: genRandomSorters(3),
			},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SortEntries(tt.args.entries, tt.args.sorters)
			var b bool
			if len(tt.args.sorters) > 0 && tt.args.sorters[0].field == RANDOM {
				b = !isSorted(t, tt.args.entries, tt.args.sorters)
			} else {
				b = isSorted(t, tt.args.entries, tt.args.sorters)
			}
			require.Truef(t, b, "%s not sorted correctly", tt.name)
		})
	}
}

func isSorted(t *testing.T, entries []Entry, sorters []Sorter) bool {
	if len(entries) == 0 || len(sorters) == 0 {
		return true
	}

	sorter := sorters[0]
	if !checkSorted(t, entries, sorter) {
		return false
	}

	if !checkSubSorts(t, entries, sorters) {
		return false
	}
	return true
}

func checkSorted(t *testing.T, entries []Entry, sorter Sorter) bool {
	sortField := sorter.field
	if sortField == RANDOM {
		sortField = FILENAME
	}

	var prev, val string
	var ok bool
	for i, entry := range entries {
		if i == 0 {
			prev, ok = entry.fields[sortField]
			require.Truef(t, ok, "missing field `%s`", sortField)
			continue
		}
		val, ok = entry.fields[sortField]
		require.True(t, ok)
		if sorter.descending {
			if val > prev {
				return false
			}
		} else {
			if val < prev {
				return false
			}
		}
	}
	return true
}

func checkSubSorts(t *testing.T, entries []Entry, sorters []Sorter) bool {
	if len(sorters) < 2 {
		return true
	}

	sorter := sorters[0]

	// look for any repeating sequences and check the sub-sorts
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
				if !isSorted(t, entries[start:i], sorters[1:]) {
					return false
				}
			}
		}
		prev = val
	}
	if found && !isSorted(t, entries[start:], sorters[1:]) {
		return false
	}
	return true
}

func genEntries(num int) []Entry {
	entries := make([]Entry, 0, num)
	for artists := 0; artists < num; artists++ {
		artist := genRandomString(16)
		for albums := 0; albums < rnd.Intn(6)+1; albums++ {
			album := genRandomString(16)
			year := strconv.Itoa(rnd.Intn(45) + 1970)
			genre := genGenre()
			date := genRandomDate()
			for tracks := 1; tracks < 13; tracks++ {
				song := genRandomString(16)
				entry := Entry{
					filespec: "/home/test/file-" + genRandomString(16),
					music:    true,
					fields: map[string]string{
						FILENAME:     song + ".mp3",
						LASTMODIFIED: date,
						ARTIST:       artist,
						ALBUM:        album,
						SONG:         song,
						TRACK:        strconv.Itoa(tracks),
						GENRE:        genre,
						YEAR:         year,
					},
				}
				entries = append(entries, entry)
			}
		}
	}
	return entries
}

func genRandomEntries(num int) []Entry {
	entries := make([]Entry, 0, num)
	for i := 0; i < num; i++ {
		entry := Entry{
			filespec: "/home/test/file-" + genRandomString(16),
			music:    true,
			fields:   make(map[string]string, len(fields)),
		}
		for j := 0; j < len(fields); j++ {
			entry.fields[fields[j]] = genRandomString(10)
		}
		entries = append(entries, entry)
	}
	return entries
}

func genRandomString(size int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!$ "
	const charCount = len(chars)
	sb := strings.Builder{}
	for i := 0; i < size; i++ {
		idx := rnd.Intn(charCount)
		sb.WriteByte(chars[idx])
	}
	return sb.String()
}

func genRandomDate() string {
	hours := rnd.Intn(1000)
	return time.Now().Add(time.Hour * time.Duration(hours)).Format(TIME_FORMAT)
}

func genGenre() string {
	var genres = []string{"rock", "blues", "punk", "death metal", "gangsta rap", "polka", "disco", "dubstep", "trance", "ambient"}
	var size = len(genres)
	return genres[rnd.Intn(size)]
}

var booleans = []bool{true, false}

func genRandomSorters(num int) []Sorter {
	sorters := make([]Sorter, 0, num)
	for i := 0; i < num; i++ {
		sorter := Sorter{
			field:      fields[rnd.Intn(len(fields))],
			descending: booleans[rnd.Intn(len(booleans))],
		}
		sorters = append(sorters, sorter)
	}
	return sorters
}

// TestIsSorted ensures that the testing helper method `isSorted` is correct.
func TestIsSorted(t *testing.T) {
	setup(t)
	defer teardown(t)

	type args struct {
		entries []Entry
		sorters []Sorter
	}
	type test struct {
		name       string
		wantSorted bool
		args       args
	}
	tests := []test{
		{name: "case1- artist unsorted, album unsorted, track unsorted", wantSorted: false, args: args{entries: []Entry{
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "3"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "1"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "4"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "2"}},

			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "4"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "1"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "3"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "2"}},

			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "3"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "1"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "4"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "2"}},
		}, sorters: []Sorter{
			{field: ARTIST, descending: false},
			{field: ALBUM, descending: true},
			{field: TRACK, descending: false},
		}}},
		{name: "case2 - artist sorted, album unsorted, track unsorted", wantSorted: false, args: args{entries: []Entry{
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "3"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "1"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "4"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "2"}},

			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "3"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "1"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "4"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "2"}},

			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "4"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "1"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "3"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "2"}},
		}, sorters: []Sorter{
			{field: ARTIST, descending: false},
			{field: ALBUM, descending: true},
			{field: TRACK, descending: false},
		}}},
		{name: "case3 - artist sorted, album sorted, track unsorted", wantSorted: false, args: args{entries: []Entry{
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "3"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "1"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "4"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "2"}},

			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "4"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "1"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "3"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "2"}},

			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "3"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "1"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "4"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "2"}},
		}, sorters: []Sorter{
			{field: ARTIST, descending: false},
			{field: ALBUM, descending: true},
			{field: TRACK, descending: false},
		}}},
		{name: "case4 - artist sorted, album sorted, track sorted", wantSorted: true, args: args{entries: []Entry{
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "4"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "3"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "2"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Led Zeppelin", ALBUM: "IV", TRACK: "1"}},

			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "4"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "3"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "2"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Moving Pictures", TRACK: "1"}},

			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "4"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "3"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "2"}},
			{filespec: genRandomString(10), music: true, fields: map[string]string{ARTIST: "Rush", ALBUM: "Permanent Waves", TRACK: "1"}},
		}, sorters: []Sorter{
			{field: ARTIST, descending: false},
			{field: ALBUM, descending: false},
			{field: TRACK, descending: true},
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted := isSorted(t, tt.args.entries, tt.args.sorters)
			require.Equalf(t, tt.wantSorted, sorted, "%s - want %t, got %t", tt.name, tt.wantSorted, sorted)
		})
	}
}
