package reflectsort

import (
	"reflect"

	sortdeps "{{.Module}}/sandbox/deps/sortdeps"

	"{{.Module}}/sandbox/deps"
)

// Bind fills deps.Deps.Sortdeps with an insertion sort over reflect.Swapper —
// the same contract adapters/libs/sortdeps fills with the standard library's
// sort. It is the second implementation of one contract: nothing inside the
// sandbox can tell the two apart, and which one a build binds is the
// available's decision.
//
// Insertion sort is stable, so Slice and SliceStable are one function. It is
// O(n²): this adapter is for a program that sorts short lists and would rather
// not carry `sort`, not a drop-in for large ones.
func Bind(deps *deps.Deps) {
	deps.Sortdeps = sortdeps.Sandbox{
		Strings: func(list []string) {
			insertion(len(list),
				func(i int, j int) bool { return list[i] < list[j] },
				func(i int, j int) { list[i], list[j] = list[j], list[i] })
		},
		Slice: func(slice any, less func(i int, j int) bool) {
			sortSlice(slice, less)
		},
		SliceStable: func(slice any, less func(i int, j int) bool) {
			sortSlice(slice, less)
		},
	}
}

// sortSlice sorts any slice through reflect, the one way to swap two elements
// of a slice whose element type is not known at compile time. A value that is
// not a slice is left alone, the way the standard library's own sort would
// panic on it — here the contract has no error to report it through.
func sortSlice(slice any, less func(i int, j int) bool) {
	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		return
	}
	insertion(value.Len(), less, reflect.Swapper(slice))
}

// insertion sorts length elements, comparing with less and swapping with swap.
// Both take current positions, which is what makes one pass enough: an element
// walks left only past the elements it sorts strictly before, so equal ones
// keep their original order.
func insertion(length int, less func(i int, j int) bool, swap func(i int, j int)) {
	for i := 1; i < length; i++ {
		for j := i; j > 0 && less(j, j-1); j-- {
			swap(j, j-1)
		}
	}
}
