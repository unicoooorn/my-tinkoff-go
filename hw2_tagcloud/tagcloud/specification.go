package tagcloud

import (
	"sort"
)

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	toFreq map[string]int
	sorted []string
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag     string
	Counter int
}

// New should create a valid TagCloud instance
func New() TagCloud {
	return TagCloud{make(map[string]int), make([]string, 0)}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (t *TagCloud) AddTag(tag string) {
	if _, ok := t.toFreq[tag]; !ok {
		t.toFreq[tag] = 1
		t.sorted = append(t.sorted, tag)
	} else {
		t.toFreq[tag]++
	}
	sort.SliceStable(t.sorted, func(i, j int) bool {
		return t.toFreq[t.sorted[i]] > t.toFreq[t.sorted[j]]
	})
}

// TopN returns top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (t *TagCloud) TopN(n int) []TagStat {
	if len(t.sorted) < n {
		n = len(t.sorted)
	}
	elems := make([]TagStat, n)
	for i := 0; i < n; i++ {
		elems[i] = TagStat{t.sorted[i], t.toFreq[t.sorted[i]]}
	}
	return elems
}
