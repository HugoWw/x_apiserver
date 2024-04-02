package sets

import "sort"

type Set[T comparable] map[T]struct{}

func KeySet[T comparable, V any](theMap map[T]V) Set[T] {
	ss := Set[T]{}
	for item := range theMap {
		ss.Inset(item)
	}

	return ss
}

func NewSet[T comparable](items ...T) Set[T] {
	s := make(Set[T], len(items))
	s.Inset(items...)
	return s
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Inset(items ...T) Set[T] {
	for _, item := range items {
		s[item] = struct{}{}
	}

	return s
}

func (s Set[T]) Delete(items ...T) Set[T] {
	for _, item := range items {
		delete(s, item)
	}

	return s
}

func (s Set[T]) Has(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Set[T]) HasAll(items ...T) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}

	return true
}

func (s Set[T]) HasAny(items ...T) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}

	return false
}

func (s Set[T]) Clone() Set[T] {
	newSet := make(Set[T], len(s))
	for item := range s {
		newSet.Inset(item)
	}

	return newSet
}

// Difference returns a set of objects that are not in s2.
func (s Set[T]) Difference(s2 Set[T]) Set[T] {
	diffSet := NewSet[T]()

	for item := range s {
		if !s2.Has(item) {
			diffSet.Inset(item)
		}
	}
	return diffSet
}

func (s Set[T]) SymmetricDifference(s2 Set[T]) Set[T] {
	return s.Difference(s2).Union(s2.Difference(s))
}

// Union returns a new set which includes items in either s1 or s2.
func (s Set[T]) Union(s2 Set[T]) Set[T] {
	newSet := s.Clone()
	for item := range s2 {
		newSet.Inset(item)
	}

	return newSet
}

// Intersection return a new set which includes items in both s1 and s2.
func (s Set[T]) Intersection(s2 Set[T]) Set[T] {
	var walk, other Set[T]

	result := NewSet[T]()

	if s.Len() < s2.Len() {
		walk = s
		other = s2
	} else {
		walk = s2
		other = s
	}

	for item := range walk {
		if !other.Has(item) {
			result.Inset(item)
		}
	}

	return result
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s Set[T]) IsSuperset(s2 Set[T]) bool {
	for item := range s2 {
		if !s.Has(item) {
			return false
		}
	}

	return true
}

func (s Set[T]) Equal(s2 Set[T]) bool {
	return len(s) == len(s2) && s.IsSuperset(s2)
}

// PopAny returns a single element from the set.
func (s Set[T]) PopAny() (T, bool) {
	for item := range s {
		s.Delete(item)
		return item, true
	}

	var zeroValueT T
	return zeroValueT, false
}

// UnSortedList convert set to slice,return unsorted slice
func (s Set[T]) UnSortedList() []T {
	setList := make([]T, len(s))
	for item := range s {
		setList = append(setList, item)
	}

	return setList
}

type sortableSliceOfGeneric[T ordered] []T

func (g sortableSliceOfGeneric[T]) Len() int           { return len(g) }
func (g sortableSliceOfGeneric[T]) Less(i, j int) bool { return less[T](g[i], g[j]) }
func (g sortableSliceOfGeneric[T]) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }

func less[T ordered](lhs, rhs T) bool {
	return lhs < rhs
}

func List[T ordered](s Set[T]) []T {
	res := make(sortableSliceOfGeneric[T], 0, len(s))
	for key := range s {
		res = append(res, key)
	}

	sort.Sort(res)
	return res
}
