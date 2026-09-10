package apishape

// convert returns the expression that carries value across, for one type of
// the api read in one direction. It is the whole of the convertibility rule in
// code: an expression that names none of the package's own types is identical
// on both sides and crosses untouched, and everything else is either a plain
// Go conversion or a generated function.
func convert(plan *planner, expr string, direction string, value string) string {
	if !References(plan.api, expr) {
		return value
	}

	if Declares(plan.api, expr) {
		if IsStruct(plan.api, expr) {
			return converterFor(plan, expr, direction) + "(" + value + ")"
		}
		_, to := qualifiers(plan, direction)
		return Qualify(plan.api, expr, to) + "(" + value + ")"
	}

	if signature, ok := ParseSignature(plan.deps, expr); ok {
		return closure(plan, signature, direction, value)
	}

	if _, ok := sliceElement(expr); ok {
		return converterFor(plan, expr, direction) + "(" + value + ")"
	}
	if _, ok := pointerElement(expr); ok {
		return converterFor(plan, expr, direction) + "(" + value + ")"
	}
	if _, _, ok := mapElement(expr); ok {
		return converterFor(plan, expr, direction) + "(" + value + ")"
	}

	plan.err = plan.deps.Std.Errorf("cannot convert %s: only a named type, a slice, a map, a pointer and a func of convertible types can cross", expr)
	return value
}

// closure writes the func value that stands in for one func field. Its
// parameters run the other way: the closure is called with local values and
// has to hand remote ones to the func it wraps, so every parameter is
// converted in the flipped direction and every result in this one.
func closure(plan *planner, signature Signature, direction string, value string) string {
	flipped := flip(direction)
	_, to := qualifiers(plan, direction)
	from, _ := qualifiers(plan, flipped)

	params := ""
	arguments := ""
	for index, param := range signature.Params {
		name := "p" + plan.deps.Stringsdeps.FormatInt(int64(index), 10)

		if index > 0 {
			params += ", "
			arguments += ", "
		}

		element, variadic := variadicElement(param.Type)
		if variadic {
			if References(plan.api, element) {
				plan.err = plan.deps.Std.Errorf("cannot convert %s: a variadic parameter of a type of the api has no element to convert one by one", param.Type)
				return value
			}
			params += name + " ..." + element
			arguments += name + "..."
			continue
		}

		params += name + " " + Qualify(plan.api, param.Type, from)
		arguments += convert(plan, param.Type, flipped, name)
	}

	call := value + "(" + arguments + ")"

	switch len(signature.Results) {
	case 0:
		return "func(" + params + ") { " + call + " }"
	case 1:
		result := Qualify(plan.api, signature.Results[0].Type, to)
		return "func(" + params + ") " + result + " { return " + convert(plan, signature.Results[0].Type, direction, call) + " }"
	}

	results := ""
	names := ""
	returns := ""
	for index, result := range signature.Results {
		name := "r" + plan.deps.Stringsdeps.FormatInt(int64(index), 10)
		if index > 0 {
			results += ", "
			names += ", "
			returns += ", "
		}
		results += Qualify(plan.api, result.Type, to)
		names += name
		returns += convert(plan, result.Type, direction, name)
	}

	return "func(" + params + ") (" + results + ") {\n" +
		names + " := " + call + "\n" +
		"return " + returns + "\n}"
}

// variadicElement reports a variadic parameter and the element type it
// collects.
func variadicElement(expr string) (string, bool) {
	if len(expr) > 3 && expr[:3] == "..." {
		return expr[3:], true
	}
	return "", false
}
