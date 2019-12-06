package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSortString(t *testing.T) {
	src, err := ioutil.TempDir("", "src_")
	require.NoError(t, err, err)
	defer os.RemoveAll(src)

	dest, err := ioutil.TempDir("", "dest_")
	require.NoError(t, err, err)
	defer os.RemoveAll(dest)

	type args struct {
		opts  Opts
		dir   string
		fsort string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case1", args: args{opts: Opts{src: src, dest: dest, sort: ""}, dir: src, fsort: "track:a,album:d,file:a"}, want: "track:a,album:d,file:a"},
		{name: "case2", args: args{opts: Opts{src: src, dest: dest, sort: "random:a,title:d,date:a"}, dir: src, fsort: "track:a,album:d,file:a"}, want: "track:a,album:d,file:a"},
		{name: "case3", args: args{opts: Opts{src: src, dest: dest, sort: "random:a,title:d,date:a"}, dir: src, fsort: ""}, want: "random:a,title:d,date:a"},
		{name: "case4", args: args{opts: Opts{src: src, dest: dest, sort: ""}, dir: src, fsort: ""}, want: DEFAULT_SORT},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var deleter func() error
			var err error
			if tt.args.fsort != "" {
				deleter, err = createOptsFile(tt.args.opts.src, tt.args.fsort, true)
				require.NoErrorf(t, err, "%s cannot create opts file - %v", tt.name, err)
			}
			if got := GetSortString(tt.args.opts, tt.args.dir); got != tt.want {
				t.Errorf("GetSortString() = %v, want %v", got, tt.want)
			}
			if deleter != nil {
				err = deleter()
				require.NoErrorf(t, err, "%s cannot delete opts file - %v", tt.name, err)
			}
		})
	}
}

func createOptsFile(srcDir string, sort string, children bool) (func() error, error) {
	fn := filepath.Join(srcDir, ".mp3copy")
	file, err := os.Create(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	out := fmt.Sprintf("sort = %s\nchildren = %t", sort, children)
	_, err = file.WriteString(out)
	if err != nil {
		return nil, err
	}
	return func() error { return os.Remove(fn) }, nil
}

func TestParseSortString(t *testing.T) {
	type args struct {
		ss string
	}
	tests := []struct {
		name    string
		args    args
		want    []Sorter
		wantErr bool
	}{
		{name: "case1", args: args{ss: "artist, album, track"}, wantErr: false,
			want: []Sorter{
				{field: "artist", descending: false},
				{field: "album", descending: false},
				{field: "track", descending: false},
			},
		},
		{name: "case2", args: args{ss: "file:d, date:d, random"}, wantErr: false,
			want: []Sorter{
				{field: "file", descending: true},
				{field: "date", descending: true},
				{field: "random", descending: false},
			},
		},
		{name: "case3", args: args{ss: "year:a,song:d,genre"}, wantErr: false,
			want: []Sorter{
				{field: "year", descending: false},
				{field: "song", descending: true},
				{field: "genre", descending: false},
			},
		},
		{name: "case4", args: args{ss: "song"}, wantErr: false,
			want: []Sorter{
				{field: "song", descending: false},
			},
		},
		{name: "case5", args: args{ss: ""}, wantErr: true, want: []Sorter{}},
		{name: "case6", args: args{ss: "blap:d, bloop:d, bleep"}, wantErr: true, want: []Sorter{}},
		{name: "case7", args: args{ss: " , ,"}, wantErr: true, want: []Sorter{}},
		{name: "case8", args: args{ss: "artist:x,song:a"}, wantErr: true, want: []Sorter{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSortString(tt.args.ss)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSortString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSortString() = %v, want %v", got, tt.want)
			}
		})
	}
}
