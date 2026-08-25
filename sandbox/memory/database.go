package memory

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/keepdeps"
)

var Schemas = []keepdeps.Schema{
	{
		Name: "config",
		Itens: []keepdeps.Item{
			{Name: "name", Type: keepdeps.Key, Required: true},
			{Name: "value", Type: keepdeps, Required: true},

		},
	},
	{
		Name: "states",
		Itens: []keepdeps.Item{
			{Name: "name", Type: keepdeps.Key, Required: true},
			{}
		},
	},
}

var Props = keeptypes.Props{
	Path:    "testDatabase/",
	Schemas: Schemas,
}

func NewDatabase(path string) keepdeps.KeepDatabase {

}
