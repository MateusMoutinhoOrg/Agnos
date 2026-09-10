package sortdeps

// This package is the sandbox's *copy* of the api a sorting library exposes —
// the same mechanic as argvdeps, embeddeps, iodeps, rundeps, std and
// stringsdeps, for the same reason: the sandbox may import nothing but the
// sandbox, so `sort` may not appear inside it. The contract is restated here,
// and the adapter — which lives outside the sandbox — is what fills it.
//
// Every field sorts in place and returns nothing, exactly like the standard
// library function of the same name.

// Sandbox is the sorting library injected whole as the Deps.Sortdeps field.
type Sandbox struct {
	// Strings sorts a slice of strings into increasing order.
	Strings func(list []string)

	// Slice sorts slice given the less function, which reports whether the
	// element at index i must sort before the one at index j. The sort is not
	// guaranteed to be stable; use SliceStable when equal elements must keep
	// their original order.
	Slice func(slice any, less func(i int, j int) bool)

	// SliceStable sorts slice given the less function, keeping equal elements
	// in their original order.
	SliceStable func(slice any, less func(i int, j int) bool)
}
