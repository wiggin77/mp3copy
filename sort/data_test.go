package sort_test

import (
	"time"

	"github.com/wiggin77/mp3copy/sort"
)

type testFields struct {
	Descending bool
}

type testArgs struct {
	entries []sort.Entry
}

type testCase struct {
	name   string
	fields testFields
	args   testArgs
}

var (
	base = time.Now().Add(-time.Hour * 8760)

	case1 = []sort.Entry{
		{Filespec: "/home/music/appley.mp3", LastModified: base.Add(time.Hour * 1)},
		{Filespec: "/home/music/banana.mp3", LastModified: base.Add(time.Hour * 2)},
		{Filespec: "/home/music/carrot.mp3", LastModified: base.Add(time.Hour * 3)},
		{Filespec: "/home/music/radish.mp3", LastModified: base.Add(time.Hour * 4)},
	}
	case2 = []sort.Entry{
		{Filespec: "/home/music/zztop.mp3", LastModified: base.Add(time.Hour * 4)},
		{Filespec: "/home/music/bart.mp3", LastModified: base.Add(time.Hour * 1)},
		{Filespec: "/home/music/stevie.mp3", LastModified: base.Add(time.Hour * 3)},
		{Filespec: "/home/music/rush.mp3", LastModified: base.Add(time.Hour * 2)},
	}
	case3 = []sort.Entry{
		{Filespec: "/home/music/rush.mp3", LastModified: base},
	}
	case4 = []sort.Entry{
		{Filespec: "/home/music/ardvark.mp3", LastModified: base.Add(time.Hour * 1)},
		{Filespec: "/home/music/bart.mp3", LastModified: base.Add(time.Hour * 2)},
		{Filespec: "/home/music/stevie.mp3", LastModified: base.Add(time.Hour * 5)},
		{Filespec: "/home/music/rush.mp3", LastModified: base.Add(time.Hour * 4)},
		{Filespec: "/home/music/monty.mp3", LastModified: base.Add(time.Hour * 3)},
	}

	testCases = []testCase{
		{name: "case 1; already sorted", fields: testFields{Descending: false}, args: testArgs{entries: case1}},
		{name: "case 2; ascending", fields: testFields{Descending: false}, args: testArgs{entries: case2}},
		{name: "case 2; descending", fields: testFields{Descending: true}, args: testArgs{entries: case2}},
		{name: "case 3; ascending", fields: testFields{Descending: false}, args: testArgs{entries: case3}},
		{name: "case 3; descending", fields: testFields{Descending: true}, args: testArgs{entries: case3}},
		{name: "case 4; ascending", fields: testFields{Descending: false}, args: testArgs{entries: case4}},
		{name: "case 4; descending", fields: testFields{Descending: true}, args: testArgs{entries: case4}},
	}
)
