package tagcloud

import (
	"math"
	"sort"
)

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	tagToFreq  map[string]int
	tagsSorted []string
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() TagCloud {
	return TagCloud{make(map[string]int), make([]string, 0)}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (t *TagCloud) AddTag(tag string) {
	if _, exists := t.tagToFreq[tag]; !exists {
		t.tagToFreq[tag] = 1
		t.tagsSorted = append(t.tagsSorted, tag)
	} else {
		t.tagToFreq[tag]++
	}
	sort.SliceStable(t.tagsSorted, func(i, j int) bool {
		return t.tagToFreq[t.tagsSorted[i]] > t.tagToFreq[t.tagsSorted[j]]
	})
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (t *TagCloud) TopN(n int) []TagStat {
	reqElemCount := int(math.Min(float64(n), float64(len(t.tagsSorted))))
	reqElems := make([]TagStat, reqElemCount)
	for i := 0; i < reqElemCount; i++ {
		reqElems[i] = TagStat{t.tagsSorted[i], t.tagToFreq[t.tagsSorted[i]]}
	}
	return reqElems
}
