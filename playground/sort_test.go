package playground

import (
	"testing"
	"fmt"
	"sort"
)

type StringID string

func TestStringCompare(t *testing.T) {
	s1 := "abc"
	s2 := "bbc"
	fmt.Println(s1 < s2) // true
}

func TestStringIDCompare(t *testing.T) {
	i1 := StringID("a")
	i2 := StringID("b")
	i3 := StringID("a1")
	fmt.Println(i1 < i2) // true
	fmt.Println(i1 < i3) // true
}

func TestStringIDSearch(t *testing.T) {
	s := []StringID{"a", "b", "c"}
	i := sort.Search(len(s), func(i int) bool { return s[i] >= "b" })
	fmt.Println(i)
}

func TestStringIDSort(t *testing.T) {
	//s := []StringID{"a", "b", "c"}
	//sort.Strings(s) // cannot use s (type []StringID) as type []string in argument to sort.Strings
}
