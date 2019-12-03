package sort_test

import (
	stdsort "sort"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wiggin77/mp3copy/sort"
)

func TestFilename_Sort(t *testing.T) {
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			f := sort.Filename{
				Descending: tt.fields.Descending,
			}
			err := f.Sort(tt.args.entries)
			if err != nil {
				t.Errorf("sort failed for %s; %v", tt.name, err)
			}
			require.True(t, IsSortedStrings(extractFilenames(tt.args.entries), tt.fields.Descending), "sort failed for ", tt.name)
		})
	}
}

func extractFilenames(entries []sort.Entry) []string {
	out := make([]string, len(entries))
	for i, entry := range entries {
		out[i] = entry.Filespec
	}
	return out
}

// IsSortedStrings returns true if the string slice is sorted based
// on the ascending or descending order.
func IsSortedStrings(items []string, descending bool) bool {
	if descending {
		count := len(items)
		data := make([]string, count)
		for i, s := range items {
			data[count-1-i] = s
		}
		items = data
	}
	return stdsort.IsSorted(stdsort.StringSlice(items))
}
