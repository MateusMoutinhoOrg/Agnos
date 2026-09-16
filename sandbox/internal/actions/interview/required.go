package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// moduleFieldId is the go module path a scaffold writes into go.mod. It is the
// one field of the surface whose declaration cannot say whether it is needed —
// that depends on the folder, not on the command.
const moduleFieldId = "module"

// RequiredHere reports a field that the declaration calls optional and the
// project in front of the person makes mandatory. A declaration is written
// once and read everywhere, so it can only say "optional"; the folder being
// worked on is what decides, and the session has already read it.
//
// It is the third of the tables this package keys by field id, beside
// SuggestFor — where a field's answers come from — and RuledOut — which fields
// a previous answer removes. Every rule here is one a handler already enforces
// with an error, so making the question required only moves that error from
// after the confirm screen to the question itself.
//
// A field with no entry here is asked exactly as its declaration words it.
func RequiredHere(sandbox *api.Sandbox, command api.Command, field Field, values map[string][]any, session string) bool {
	if verbOf(command) != scaffoldVerb || field.Id != moduleFieldId {
		return false
	}
	return !sandbox.Deps.Iodeps.Exist(targetOf(sandbox, values, session) + "/go.mod")
}

// targetOf is the folder a scaffold is about to write into: the --path already
// answered, or the folder the session was opened on while that answer is still
// to come. start declares --path before --module, so the first is what this
// reads in practice.
func targetOf(sandbox *api.Sandbox, values map[string][]any, session string) string {
	if answered := answeredText(sandbox, values, pathFieldId); answered != "" {
		return answered
	}
	return session
}
