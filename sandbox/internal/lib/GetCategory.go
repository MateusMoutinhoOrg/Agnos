package lib

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/lib/category"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/lib/store"
)

// GetCategoryFactory fills api.Lib.GetCategory with a closure that looks a
// stored category up by its unique name. It reports false when no category
// carries that name.
func GetCategoryFactory(l *api.Lib) func(name string) (api.Category, bool) {
	return func(name string) (api.Category, bool) {
		record, ok := store.FindCategory(l.Deps, name)
		if !ok {
			return api.Category{}, false
		}
		return category.New(l.Deps, record), true
	}
}
