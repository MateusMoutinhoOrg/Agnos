package resultconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func BindMethods(deps *deps.Deps, conf *ResultConf) {
	conf.AddTreeEntry = func(file string, sha string) {
		conf.Tree = append(conf.Tree, TreeEntry{File: file, Sha: sha})
	}

	conf.Render = func() string {
		return Render(deps, conf)
	}
}
