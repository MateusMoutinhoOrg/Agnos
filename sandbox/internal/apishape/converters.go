package apishape

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

// The two directions a conversion runs in. A field carries values out of the
// remote package and into the local one; a func field's parameters run the
// other way, because the closure the generator writes is called with local
// values and has to hand remote ones to the remote func.
const (
	// Out converts a remote value into the local copy of its type.
	Out = "conv"
	// Back converts a local value into the remote copy of its type.
	Back = "rev"
)

// Converter is one function of the generated shim: a whole rendered body, so
// the template that writes the file only has to join them.
type Converter struct {
	Name string
	From string
	To   string
	Body string
}

// Plan is every converter one shim needs, in generation order, plus the entry
// point Bind calls.
type Plan struct {
	// Entry is the converter that turns the remote Sandbox into the local
	// contract — what Bind assigns to the Deps field.
	Entry string

	// Converters are the functions to write, in the order they were reached.
	Converters []Converter
}

// planner carries what every step of the walk needs: the api being converted
// and the two package qualifiers its types are spelled with on each side.
type planner struct {
	deps   *deps.Deps
	api    *Api
	local  string
	remote string
	names  map[string]bool
	order  []Converter
	err    error
}

// Converters plans the conversion of one api package between two copies of
// itself: `remote` is the qualifier the module's own sandbox/api is imported
// under, `local` the one the copied contract is imported under.
//
// It walks from the Sandbox type outwards, generating a converter for every
// named struct it reaches and in whichever direction it reaches it. Everything
// else — builtins, slices and maps of builtins, `any`, `error`, a named
// interface of builtin-only methods — has an identical underlying type in both
// copies and crosses with an assignment or a plain conversion.
func Converters(deps *deps.Deps, api *Api, local string, remote string) (*Plan, error) {
	if !IsStruct(api, "Sandbox") {
		return nil, deps.Std.Errorf("the remote api declares no Sandbox struct: there is nothing to convert")
	}

	plan := &planner{deps: deps, api: api, local: local, remote: remote, names: map[string]bool{}}

	entry := converterFor(plan, "Sandbox", Out)
	if plan.err != nil {
		return nil, plan.err
	}

	return &Plan{Entry: entry, Converters: plan.order}, nil
}

// qualifiers returns the package a value of the given direction is read from
// and the one it is written to.
func qualifiers(plan *planner, direction string) (string, string) {
	if direction == Out {
		return plan.remote, plan.local
	}
	return plan.local, plan.remote
}

// flip is the direction a func parameter runs in, the reverse of the one its
// field runs in.
func flip(direction string) string {
	if direction == Out {
		return Back
	}
	return Out
}

// converterFor queues one converter and returns its name. A converter already
// planned is returned as it is, which is also what stops a recursive type from
// walking forever.
func converterFor(plan *planner, expr string, direction string) string {
	name := direction + expressionName(plan.deps, expr)

	if plan.names[name] {
		return name
	}
	plan.names[name] = true

	from, to := qualifiers(plan, direction)
	converter := Converter{
		Name: name,
		From: Qualify(plan.api, expr, from),
		To:   Qualify(plan.api, expr, to),
	}

	// The entry is reserved before the body is built: the body may reach the
	// same type again, and it has to find the name already taken.
	index := len(plan.order)
	plan.order = append(plan.order, converter)

	plan.order[index].Body = converterBody(plan, expr, direction)
	return name
}

// converterBody writes the body of one converter, by the shape of the type it
// converts.
func converterBody(plan *planner, expr string, direction string) string {
	_, to := qualifiers(plan, direction)
	destination := Qualify(plan.api, expr, to)

	if Declares(plan.api, expr) {
		return structBody(plan, expr, direction, destination)
	}

	if inner, ok := sliceElement(expr); ok {
		return sliceBody(plan, inner, direction, destination)
	}

	if inner, ok := pointerElement(expr); ok {
		return pointerBody(plan, inner, direction, destination)
	}

	if key, inner, ok := mapElement(expr); ok {
		return mapBody(plan, key, inner, direction, destination)
	}

	plan.err = plan.deps.Std.Errorf("cannot convert %s: only a named struct, a slice, a map, a pointer and a func of convertible types can cross", expr)
	return ""
}

// structBody writes one struct converter, field by field: the only shape whose
// two copies are never identical types, because Go's type identity does not
// reach through the named types a struct's fields are declared with.
func structBody(plan *planner, name string, direction string, destination string) string {
	entry := plan.api.ByName[name]

	body := "\treturn " + destination + "{\n"
	for _, field := range entry.Fields {
		body += "\t\t" + field.Name + ": " + convert(plan, field.Type, direction, "v."+field.Name) + ",\n"
	}
	return body + "\t}"
}

// sliceBody writes the loop a slice of a convertible element needs. A nil
// slice stays nil: the two copies mean the same thing by it.
func sliceBody(plan *planner, element string, direction string, destination string) string {
	return "\tif v == nil {\n\t\treturn nil\n\t}\n\n" +
		"\tout := make(" + destination + ", len(v))\n" +
		"\tfor index := range v {\n" +
		"\t\tout[index] = " + convert(plan, element, direction, "v[index]") + "\n" +
		"\t}\n\n\treturn out"
}

// pointerBody writes the conversion of a pointer to a convertible type.
func pointerBody(plan *planner, element string, direction string, destination string) string {
	return "\tif v == nil {\n\t\treturn nil\n\t}\n\n" +
		"\tout := " + convert(plan, element, direction, "*v") + "\n" +
		"\treturn &out"
}

// mapBody writes the loop a map with a convertible value needs. The key is
// always a plain type: a key that named one of the package's own types would
// have to be converted to be looked up, and the shape rule keeps it out.
func mapBody(plan *planner, key string, element string, direction string, destination string) string {
	return "\tif v == nil {\n\t\treturn nil\n\t}\n\n" +
		"\tout := make(" + destination + ", len(v))\n" +
		"\tfor key, item := range v {\n" +
		"\t\tout[key] = " + convert(plan, element, direction, "item") + "\n" +
		"\t}\n\n\treturn out"
}
