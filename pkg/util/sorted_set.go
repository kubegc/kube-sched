package util

import "sort"

type SortedSet struct {
	m map[string]bool
}
func NewSortedSet() *SortedSet {
	return &SortedSet{
		m:make(map[string]bool),
	}
}
func (ss *SortedSet) Contains(key string) bool {
	if _, ok := ss.m[key]; ok{
		return true
	}
	return false
}
func (ss *SortedSet) Add(key string) {
	ss.m[key] = true
}

func (ss *SortedSet) Delete(key string) {
	delete(ss.m, key)
}

func(ss *SortedSet) SortedKeys() []string{

	keys := make([]string, 0)
	if ss.m == nil {
		return keys
	}
	for k, _ := range ss.m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (ss *SortedSet) Size() int {
	return len(ss.m)
}