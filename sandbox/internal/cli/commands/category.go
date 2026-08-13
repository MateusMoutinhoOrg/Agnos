package commands

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// Category runs the `category` command group — add, list, remove — over the
// arguments still unread in the injected parser.
func Category(l *api.Lib, quiet bool) int {
	verb := l.Deps.VerbLib
	action, err := verb.GetNextStringArg()
	if err != nil {
		return UsageError(l, config.CategoryActionMissing)
	}

	switch action {
	case "list":
		categories := l.ListCategories()
		if len(categories) == 0 {
			l.Deps.Printf("%s\n", config.NoCategories)
			return api.ExitOk
		}
		for _, stored := range categories {
			l.Deps.Printf("%s\n", stored.String())
		}
		return api.ExitOk

	case "add":
		name, err := verb.GetNextStringArg()
		if err != nil {
			return UsageError(l, config.CategoryAddNameMissing)
		}
		created, ok := l.AddCategory(name)
		if !ok {
			return Failure(l, config.CategoryNotCreated, name)
		}
		if !quiet {
			l.Deps.Printf(config.CategoryAdded+"\n", created.String())
		}
		return api.ExitOk

	case "remove":
		name, err := verb.GetNextStringArg()
		if err != nil {
			return UsageError(l, config.CategoryRemoveNameMissing)
		}
		stored, found := l.GetCategory(name)
		if !found {
			return Failure(l, config.CategoryNotFound, name)
		}
		if !stored.Remove() {
			return Failure(l, config.CategoryNotRemoved, name)
		}
		if !quiet {
			l.Deps.Printf(config.CategoryRemoved+"\n", name)
		}
		return api.ExitOk
	}
	return UsageError(l, config.CategoryActionUnknown, action)
}
