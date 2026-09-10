package sortdeps

import (
	"sort"

	sortdeps "{{.Module}}/sandbox/deps/sortdeps"

	"{{.Module}}/sandbox/deps"
)

// Bind fills deps.Deps.Sortdeps with the standard library's sort. Every field
// is a straight delegation: the contract restates the standard library api so
// the sandbox can call it without importing it.
func Bind(deps *deps.Deps) {
	deps.Sortdeps = sortdeps.Sandbox{
		Strings: func(list []string) {
			sort.Strings(list)
		},
		Slice: func(slice any, less func(i int, j int) bool) {
			sort.Slice(slice, less)
		},
		SliceStable: func(slice any, less func(i int, j int) bool) {
			sort.SliceStable(slice, less)
		},
	}
}
