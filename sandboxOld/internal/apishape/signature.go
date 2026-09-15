package apishape

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// Param is one parameter or result of a func type expression, read back apart.
// Name is "" for an unnamed one, which is what most results are.
type Param struct {
	Name string
	Type string
}

// Signature is one func type expression read back apart. The Go parser hands a
// struct field's type back as source ("func(props StartProps) error"), so a
// generator that has to write a closure over that field has to read it again.
type Signature struct {
	Params  []Param
	Results []Param
}

// ParseSignature reads a func type expression into its parameters and results,
// and reports whether expr is a func type at all.
func ParseSignature(sandbox *api.Sandbox, expr string) (Signature, bool) {
	rest := sandbox.Deps.Stringsdeps.TrimSpace(expr)

	if !sandbox.Deps.Stringsdeps.HasPrefix(rest, "func") {
		return Signature{}, false
	}
	rest = sandbox.Deps.Stringsdeps.TrimSpace(rest[len("func"):])

	if len(rest) == 0 || rest[0] != '(' {
		return Signature{}, false
	}

	close_index := matching(rest, 0)
	if close_index < 0 {
		return Signature{}, false
	}

	signature := Signature{Params: splitParams(sandbox, rest[1:close_index])}

	results := sandbox.Deps.Stringsdeps.TrimSpace(rest[close_index+1:])
	if results == "" {
		return signature, true
	}

	if results[0] == '(' && matching(results, 0) == len(results)-1 {
		signature.Results = splitParams(sandbox, results[1:len(results)-1])
		return signature, true
	}

	signature.Results = []Param{{Type: results}}
	return signature, true
}

// splitParams cuts one parameter or result list at its top-level commas and
// reads each item.
func splitParams(sandbox *api.Sandbox, list string) []Param {
	var params []Param

	for _, chunk := range splitTop(list, ',') {
		chunk = sandbox.Deps.Stringsdeps.TrimSpace(chunk)
		if chunk == "" {
			continue
		}
		params = append(params, readParam(sandbox, chunk))
	}

	return params
}

// readParam tells `props StartProps` from `[]string`: a name is a bare
// identifier followed by the type it binds, and a type on its own never is.
func readParam(sandbox *api.Sandbox, chunk string) Param {
	space := indexTop(chunk, ' ')
	if space < 0 {
		return Param{Type: chunk}
	}

	name := chunk[:space]
	rest := sandbox.Deps.Stringsdeps.TrimSpace(chunk[space:])

	if rest == "" || !isIdentifier(name) {
		return Param{Type: chunk}
	}

	return Param{Name: name, Type: rest}
}

// isIdentifier reports whether text is a plain Go identifier.
func isIdentifier(text string) bool {
	if text == "" || (text[0] >= '0' && text[0] <= '9') {
		return false
	}
	for index := 0; index < len(text); index++ {
		if !isIdentChar(text[index]) {
			return false
		}
	}
	return true
}

// matching returns the index of the bracket closing the one at open, or -1
// when the expression is unbalanced.
func matching(expr string, open int) int {
	depth := 0
	for index := open; index < len(expr); index++ {
		switch expr[index] {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
			if depth == 0 {
				return index
			}
		}
	}
	return -1
}

// splitTop cuts expr at every occurrence of separator that is not inside
// brackets.
func splitTop(expr string, separator byte) []string {
	var parts []string

	depth := 0
	start := 0
	for index := 0; index < len(expr); index++ {
		switch expr[index] {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
		case separator:
			if depth == 0 {
				parts = append(parts, expr[start:index])
				start = index + 1
			}
		}
	}

	return append(parts, expr[start:])
}

// indexTop returns the first index of separator outside brackets, or -1.
func indexTop(expr string, separator byte) int {
	depth := 0
	for index := 0; index < len(expr); index++ {
		switch expr[index] {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
		case separator:
			if depth == 0 {
				return index
			}
		}
	}
	return -1
}
