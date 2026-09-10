package apishape

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
)

// Violations reports every declaration of an api package that could not be
// copied into a consumer and converted.
//
// The criterion is one: the underlying type has to be identical between the
// two copies of the package. Everything follows from it. A builtin, a slice or
// map of builtins, `any`, `error`, a func of builtins and a named interface of
// builtin-only methods are already identical, and cross with a plain
// assignment or conversion. A struct whose fields name the package's own types
// is not — Go's type identity is not recursive through named types — so the
// generator writes it a converter, field by field. What is left cannot be
// written at all: a type from another package (already barred by "api imports
// nothing"), a type parameter, a channel, an anonymous struct or interface, or
// an embedded field.
//
// The check is not opt-in and has no flag. What varies between repos is the
// shape of the api, and a shape is verified, not declared; an installable
// switch defaulting to off would only move the violation to the consumer's
// install, which is where it hurts.
func Violations(deps *deps.Deps, api *Api) []string {
	var violations []string

	for _, file := range api.Files {
		for _, entry := range file.Parsed.Types {
			violations = append(violations, typeViolations(deps, api, file.Name, entry)...)
		}
		for _, entry := range file.Parsed.Functions {
			violations = append(violations, signatureViolations(deps, api, file.Name, "function "+entry.Name, entry)...)
		}
	}

	return violations
}

// typeViolations checks one type declaration, in the terms of its kind: a
// struct through its fields, an interface through its methods, anything else
// through its underlying type.
func typeViolations(deps *deps.Deps, api *Api, file string, entry goimportsdeps.Type) []string {
	var violations []string

	where := "type " + entry.Name

	switch entry.Kind {
	case "struct":
		for _, field := range entry.Fields {
			if field.Name == "" {
				violations = append(violations, violation(file, where,
					"embeds "+field.Type+"; an embedded field has no name to convert through, so give it one"))
				continue
			}
			violations = append(violations, exprViolations(deps, api, file, where+" field "+field.Name, field.Type)...)
		}
	case "interface":
		for _, method := range entry.Methods {
			violations = append(violations, signatureViolations(deps, api, file, where+" method "+method.Name, method)...)
		}
	default:
		violations = append(violations, exprViolations(deps, api, file, where, entry.Underlying)...)
	}

	return violations
}

// signatureViolations checks every parameter and result of one function or
// method declaration.
func signatureViolations(deps *deps.Deps, api *Api, file string, where string, entry goimportsdeps.Function) []string {
	var violations []string

	for _, param := range entry.Params {
		violations = append(violations, exprViolations(deps, api, file, where+" parameter", param.Type)...)
	}
	for _, result := range entry.Results {
		violations = append(violations, exprViolations(deps, api, file, where+" result", result.Type)...)
	}

	return violations
}

// exprViolations is the whole rule, applied to one type expression. It reads
// the expression as tokens rather than re-parsing it: every question it asks —
// is this identifier known, is it instantiated, is it a channel, is it an
// anonymous struct — is answered by an identifier and the character after it.
func exprViolations(deps *deps.Deps, api *Api, file string, where string, expr string) []string {
	var violations []string

	tokens := tokenize(expr)

	for index, item := range tokens {
		if !item.Ident {
			continue
		}

		if isQualified(tokens, index) {
			violations = append(violations, violation(file, where,
				"names "+expr+", a type of another package; sandbox/api/ imports nothing, so nothing there can be copied"))
			continue
		}

		if item.Text == "chan" {
			violations = append(violations, violation(file, where,
				"names "+expr+"; a channel cannot be converted between two copies of a package"))
			continue
		}

		if (item.Text == "struct" || item.Text == "interface") && item.Next == '{' {
			violations = append(violations, violation(file, where,
				"names "+expr+", an anonymous "+item.Text+"; give it a name so the generator has one to convert through"))
			continue
		}

		if isInstantiation(tokens, index) {
			violations = append(violations, violation(file, where,
				"names "+expr+", a generic instantiation; the generator has no type arguments to carry across"))
			continue
		}

		if isBinding(tokens, index) || isPredeclared(item.Text) || Declares(api, item.Text) {
			continue
		}

		violations = append(violations, violation(file, where,
			"names "+item.Text+" in "+expr+", which is neither predeclared nor declared in this package"))
	}

	return violations
}

// violation words one finding the same way for every rule.
func violation(file string, where string, reason string) string {
	return file + ": " + where + " " + reason
}
