package publicfunctions

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/keepdeps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/lib/category"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/lib/store"
)

// AddCategoryFactory fills api.Lib.AddCategory with a closure that creates
// the named category in the injected database and returns it. Creation is
// idempotent: a name already taken makes the injected database report a key
// conflict, and the stored category is returned instead of a failure.
func AddCategoryFactory(l *api.Lib) func(name string) (api.Category, bool) {
	return func(name string) (api.Category, bool) {
		if name == "" {
			return api.Category{}, false
		}
		categories, ok := store.Categories(l.Deps)
		if !ok {
			return api.Category{}, false
		}

		record, err := categories.NewItem(map[string]any{
			store.NameField:      name,
			store.CreatedAtField: l.Deps.Now().Unix(),
		})
		if err == nil {
			return category.New(l.Deps, record), true
		}
		if err.Type != keepdeps.KeyConflict {
			return api.Category{}, false
		}

		stored, found := categories.FindByKey(store.NameField, name)
		if !found {
			return api.Category{}, false
		}
		return category.New(l.Deps, stored), true
	}
}
