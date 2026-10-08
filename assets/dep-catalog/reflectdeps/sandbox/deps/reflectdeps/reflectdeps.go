package reflectdeps

// This package is the sandbox's *copy* of the api a reflection library exposes
// — the same mechanic as sortdeps, stringsdeps and std, for the same reason:
// the sandbox may import nothing but the sandbox, so `reflect` may not appear
// inside it. The contract is restated here, and the adapter — which lives
// outside the sandbox — is what fills it.
//
// Every field is one reflection primitive and nothing more: it builds, reads,
// fills or calls a value whose type is known only at run time. What a caller
// does with a struct tag, or which function it calls, is the caller's business.

// Sandbox is the reflection library injected whole as the Deps.Reflectdeps
// field.
type Sandbox struct {
	// NumIn returns how many parameters the function fn takes, -1 when fn
	// is not a function.
	NumIn func(fn any) int

	// NewIn builds a fresh value for the parameter at index of the function
	// fn: a pointer to a new zero value when that parameter is a pointer, the
	// zero value itself otherwise. It returns nil when fn is not a function
	// or index is out of range.
	NewIn func(fn any, index int) any

	// Call calls the function fn with args, one per parameter, and returns
	// what it returned. A nil arg is passed as the zero value of its
	// parameter. It returns nil when fn is not a function or the argument
	// count differs.
	Call func(fn any, args []any) []any

	// NumField returns how many fields the struct target points at, -1 when
	// target is not a pointer to a struct.
	NumField func(target any) int

	// FieldName returns the Go name of the field at index of the struct
	// target points at, "" when there is none.
	FieldName func(target any, index int) string

	// FieldTag returns the value of key in the tag of the field at index of
	// the struct target points at, "" when the tag or the field is absent.
	FieldTag func(target any, index int, key string) string

	// FieldKind returns the kind of the field at index of the struct target
	// points at, spelled as the standard library spells it ("string",
	// "float64", "bool", "slice", "struct"), "" when there is none.
	FieldKind func(target any, index int) string

	// SetField stores value in the field at index of the struct target
	// points at, and reports an error when the field is absent, not
	// settable, or of a type value cannot be assigned to.
	SetField func(target any, index int, value any) error
}
