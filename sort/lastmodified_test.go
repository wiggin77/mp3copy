package sort_test

import (
	stdsort "sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wiggin77/mp3copy/sort"
)

func TestLastModified_Sort(t *testing.T) {
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			f := sort.LastModified{
				Descending: tt.fields.Descending,
			}
			err := f.Sort(tt.args.entries)
			if err != nil {
				t.Errorf("sort failed for %s; %v", tt.name, err)
			}
			require.True(t, IsSortedTime(extractLastModified(tt.args.entries), tt.fields.Descending), "sort failed for ", tt.name)
		})
	}
}

func extractLastModified(entries []sort.Entry) []time.Time {
	out := make([]time.Time, len(entries))
	for i, entry := range entries {
		out[i] = entry.LastModified
	}
	return out
}

// IsSortedTime returns true if the slice of time.Time's is sorted based
// on the ascending or descending order.
func IsSortedTime(items []time.Time, descending bool) bool {
	if descending {
		count := len(items)
		data := make([]time.Time, count)
		for i, t := range items {
			data[count-1-i] = t
		}
		items = data
	}
	return stdsort.IsSorted(TimeSlice(items))
}

// TimeSlice attaches the methods of sort.Interface to []time.Time, sorting in increasing order.
type TimeSlice []time.Time

func (t TimeSlice) Len() int           { return len(t) }
func (t TimeSlice) Less(i, j int) bool { return t[i].Before(t[j]) }
func (t TimeSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TimeSlice) Sort()              { stdsort.Sort(t) }
