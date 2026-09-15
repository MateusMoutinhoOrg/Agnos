package apishape

// predeclared is every identifier a type expression may use without the
// package declaring it: Go's predeclared types plus the four type keywords.
// Anything else in an expression is either a type of the package or something
// that cannot cross a package boundary — an import, or a type parameter.
var predeclared = []string{
	"any", "bool", "byte", "chan", "comparable", "complex64", "complex128",
	"error", "float32", "float64", "func", "int", "int8", "int16", "int32",
	"int64", "interface", "map", "rune", "string", "struct", "uint", "uint8",
	"uint16", "uint32", "uint64", "uintptr",
}

// token is one lexical piece of a type expression: an identifier, or the
// punctuation between identifiers. Splitting an expression this way is enough
// for every question asked of it, because a type expression has no literals
// and no operators beyond the ones that build types.
type token struct {
	Text  string
	Ident bool
	// Next is the first non-space character after the token, 0 at the end of
	// the expression. It is what tells `Box[T]` from `[]T`, `struct {` from a
	// named struct, and — the case a type expression is full of — a parameter
	// name from a type: `path` in `func(path string) error` is followed by
	// another identifier, a type never is.
	Next byte
	// After is the first non-space character past Next, which is the only
	// thing that tells `args []string` from `Box[T]`: a slice opens with an
	// empty bracket, an instantiation with a type argument.
	After byte
}

// tokenize splits a type expression into identifiers and the punctuation
// between them.
func tokenize(expr string) []token {
	var tokens []token

	start := 0
	for index := 0; index <= len(expr); index++ {
		if index < len(expr) && isIdentChar(expr[index]) {
			continue
		}
		if index > start {
			tokens = append(tokens, token{Text: expr[start:index], Ident: true, Next: nextRune(expr, index), After: afterRune(expr, index)})
		}
		if index < len(expr) {
			tokens = append(tokens, token{Text: string(expr[index]), Next: nextRune(expr, index+1), After: afterRune(expr, index+1)})
		}
		start = index + 1
	}

	return tokens
}

// isIdentChar reports whether a byte may appear in a Go identifier. Every
// identifier a contract may name is ASCII: the shape rule rejects anything the
// generated shim could not spell.
func isIdentChar(letter byte) bool {
	return letter == '_' ||
		(letter >= 'a' && letter <= 'z') ||
		(letter >= 'A' && letter <= 'Z') ||
		(letter >= '0' && letter <= '9')
}

// nextRune returns the first non-space byte at or after index, or 0 past the
// end of the expression.
func nextRune(expr string, index int) byte {
	for ; index < len(expr); index++ {
		if expr[index] != ' ' && expr[index] != '\t' {
			return expr[index]
		}
	}
	return 0
}

// afterRune returns the first non-space byte past the one nextRune finds.
func afterRune(expr string, index int) byte {
	for ; index < len(expr); index++ {
		if expr[index] != ' ' && expr[index] != '\t' {
			return nextRune(expr, index+1)
		}
	}
	return 0
}

// isInstantiation reports whether the identifier at index is a generic type
// being instantiated. `map` is the one predeclared identifier that takes a
// bracket of its own, and a slice or array type after a parameter name opens
// with a bracket that is not a type argument.
func isInstantiation(tokens []token, index int) bool {
	item := tokens[index]
	if item.Next != '[' || item.Text == "map" {
		return false
	}
	return item.After != ']' && !(item.After >= '0' && item.After <= '9')
}

// isBinding reports whether the identifier at index is a name being declared
// rather than a type being named — a parameter or result name. A name is
// always followed by the type it binds, and a type never is.
func isBinding(tokens []token, index int) bool {
	item := tokens[index]
	if item.Next == '[' {
		return !isInstantiation(tokens, index)
	}
	return isIdentChar(item.Next) || item.Next == '*' || item.Next == '(' || item.Next == '.'
}

// isPredeclared reports whether name is one of Go's own.
func isPredeclared(name string) bool {
	for _, candidate := range predeclared {
		if candidate == name {
			return true
		}
	}
	return false
}

// isQualified reports whether the identifier at index is the selector half of
// a qualified name (`time.Time`), which is what an import looks like from
// inside a type expression.
func isQualified(tokens []token, index int) bool {
	return index > 0 && tokens[index-1].Text == "."
}

// References reports whether expr names any type this package declares. An
// expression that does not is identical in both copies of the package, so the
// generated shim assigns it across with no conversion at all.
func References(shape *Api, expr string) bool {
	tokens := tokenize(expr)
	for index, item := range tokens {
		if !item.Ident || isQualified(tokens, index) {
			continue
		}
		if Declares(shape, item.Text) {
			return true
		}
	}
	return false
}

// Qualify rewrites every identifier of expr that this package declares into
// qualifier.Identifier, turning one type expression into the same expression
// as read from another package. It is what makes `[]Item` become
// `[]remoteapi.Item` on one side and `[]mathlib.Item` on the other.
func Qualify(shape *Api, expr string, qualifier string) string {
	tokens := tokenize(expr)

	out := ""
	for index, item := range tokens {
		if item.Ident && !isQualified(tokens, index) && Declares(shape, item.Text) {
			out += qualifier + "." + item.Text
			continue
		}
		out += item.Text
	}

	return out
}
