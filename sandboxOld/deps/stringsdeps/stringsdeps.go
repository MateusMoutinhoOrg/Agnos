package stringsdeps

// This package is the sandbox's *copy* of the api a text library exposes —
// the same mechanic as argvdeps, embeddeps, iodeps, rundeps and std, for the
// same reason: the sandbox may import nothing but the sandbox, so `strings`
// and `strconv` may not appear inside it. The contract is restated here, and
// the adapter — which lives outside the sandbox — is what fills it.
//
// Every field is the standard library function of the same name, with the
// same semantics; the adapter is a straight delegation, so a caller can read
// the standard library documentation for the behaviour of any of them.

// Sandbox is the text library injected whole as the Deps.Stringsdeps field. The
// first group of fields is string manipulation, the second is conversion
// between strings and numbers.
type Sandbox struct {
	// TrimSpace returns s with leading and trailing white space removed.
	TrimSpace func(s string) string

	// Trim returns s with every leading and trailing character contained in
	// cutset removed.
	Trim func(s string, cutset string) string

	// TrimLeft returns s with every leading character contained in cutset
	// removed.
	TrimLeft func(s string, cutset string) string

	// TrimRight returns s with every trailing character contained in cutset
	// removed.
	TrimRight func(s string, cutset string) string

	// TrimPrefix returns s without the given leading prefix. When s does not
	// start with prefix, s is returned unchanged.
	TrimPrefix func(s string, prefix string) string

	// TrimSuffix returns s without the given trailing suffix. When s does not
	// end with suffix, s is returned unchanged.
	TrimSuffix func(s string, suffix string) string

	// HasPrefix reports whether s begins with prefix.
	HasPrefix func(s string, prefix string) bool

	// HasSuffix reports whether s ends with suffix.
	HasSuffix func(s string, suffix string) bool

	// Contains reports whether substr is within s.
	Contains func(s string, substr string) bool

	// ContainsAny reports whether any character of chars is within s.
	ContainsAny func(s string, chars string) bool

	// LastIndex returns the index of the last instance of substr in s, or -1
	// when substr is absent.
	LastIndex func(s string, substr string) int

	// Count returns the number of non-overlapping instances of substr in s.
	// When substr is empty it returns one plus the number of runes in s.
	Count func(s string, substr string) int

	// Split slices s into every substring separated by sep.
	Split func(s string, sep string) []string

	// Join concatenates elems, placing sep between consecutive elements.
	Join func(elems []string, sep string) string

	// Fields slices s around each run of white space, returning the
	// substrings between them.
	Fields func(s string) []string

	// FieldsFunc slices s at each run of runes satisfying f, returning the
	// substrings between them.
	FieldsFunc func(s string, f func(rune) bool) []string

	// Repeat returns count copies of s concatenated.
	Repeat func(s string, count int) string

	// ReplaceAll returns s with every non-overlapping instance of old
	// replaced by new.
	ReplaceAll func(s string, old string, new string) string

	// ToUpper returns s with every letter mapped to its upper case.
	ToUpper func(s string) string

	// ToLower returns s with every letter mapped to its lower case.
	ToLower func(s string) string

	// Quote returns s as a double-quoted Go string literal, escaping what
	// the Go syntax requires.
	Quote func(s string) string

	// MatchPattern reports whether s is matched by the regular expression
	// pattern, and errors when the pattern itself does not compile. It is the
	// one matching primitive the sandbox has: `regexp` lives on the adapter
	// side like every other standard package.
	MatchPattern func(pattern string, s string) (bool, error)

	// Atoi parses s as a decimal integer. The error reports a string that is
	// not one.
	Atoi func(s string) (int, error)

	// ParseInt parses s as an integer in the given base with the given bit
	// size. The error reports a string that is not one.
	ParseInt func(s string, base int, bit_size int) (int64, error)

	// ParseFloat parses s as a floating-point number of the given bit size.
	// The error reports a string that is not one.
	ParseFloat func(s string, bit_size int) (float64, error)

	// FormatInt returns the string representation of value in the given base.
	FormatInt func(value int64, base int) string

	// FormatFloat returns the string representation of value, formatted
	// according to the format byte, the precision and the bit size — the same
	// three controls the standard library takes.
	FormatFloat func(value float64, format byte, precision int, bit_size int) string
}
