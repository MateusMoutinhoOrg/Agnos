package docpropsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *DocPropsConf {

	// Order is left unset: a doc with no order is listed alphabetically after
	// every ordered one.
	doc_props_conf := &DocPropsConf{
		Themes: make([]string, 0),
	}
	BindMethods(deps, doc_props_conf)
	return doc_props_conf
}
