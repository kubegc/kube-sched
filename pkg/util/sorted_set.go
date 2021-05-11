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
func (sm *SortedSet) Contains(key string) bool {
	if _, ok := sm.m[key]; ok{
		return true
	}
	return false
}
func (sm *SortedSet) Add(key string) {
	sm.m[key] = true
}

func (sm *SortedSet) Delete(key string) {
	delete(sm.m, key)
}

func(sm *SortedSet) SortedKeys() []string{
	keys := make([]string, 0)

	for k, _ := range sm.m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}


