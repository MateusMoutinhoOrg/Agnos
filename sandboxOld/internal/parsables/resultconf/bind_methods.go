package resultconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, conf *ResultConf) {
	conf.AddTreeEntry = func(file string, sha string) {
		conf.Tree = append(conf.Tree, TreeEntry{File: file, Sha: sha})
	}

	conf.Render = func() string {
		return Render(sandbox, conf)
	}
}
