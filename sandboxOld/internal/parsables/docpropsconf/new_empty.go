package docpropsconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *DocPropsConf {

	// Order is left unset: a doc with no order is listed alphabetically after
	// every ordered one.
	doc_props_conf := &DocPropsConf{
		Themes: make([]string, 0),
	}
	BindMethods(sandbox, doc_props_conf)
	return doc_props_conf
}
