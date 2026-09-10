package apishape

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

// sliceElement returns the element type of a slice expression.
func sliceElement(expr string) (string, bool) {
	if len(expr) > 2 && expr[0] == '[' && expr[1] == ']' {
		return expr[2:], true
	}
	return "", false
}

// pointerElement returns the pointed-to type of a pointer expression.
func pointerElement(expr string) (string, bool) {
	if len(expr) > 1 && expr[0] == '*' {
		return expr[1:], true
	}
	return "", false
}

// mapElement returns the key and value types of a map expression.
func mapElement(expr string) (string, string, bool) {
	if len(expr) < len("map[x]y") || expr[:4] != "map[" {
		return "", "", false
	}

	close_index := matching(expr, 3)
	if close_index < 0 || close_index+1 >= len(expr) {
		return "", "", false
	}

	return expr[4:close_index], expr[close_index+1:], true
}

// expressionName is the name a helper converter for one type expression is
// filed under: the expression spelled as an identifier, so `[]Item` becomes
// SliceItem and `map[string]Item` MapStringItem. It is derived rather than
// counted, so the same expression always gets the same function.
func expressionName(deps *deps.Deps, expr string) string {
	tokens := tokenize(expr)

	name := ""
	for index, item := range tokens {
		switch {
		case item.Ident:
			name += title(deps, item.Text)
		case item.Text == "*":
			name += "Ptr"
		case item.Text == "[" && index+1 < len(tokens) && tokens[index+1].Text == "]":
			name += "Slice"
		}
	}

	return name
}

// title upper-cases the first letter of a name, the one spelling change a
// generated identifier ever makes.
func title(deps *deps.Deps, name string) string {
	if name == "" {
		return name
	}
	return deps.Stringsdeps.ToUpper(name[:1]) + name[1:]
}
