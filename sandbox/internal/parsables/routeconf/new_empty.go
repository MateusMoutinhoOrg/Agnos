package routeconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *RouteConf {
	conf := &RouteConf{
		Methods:      []string{DefaultMethod},
		HasPriority:  true,
		ResponseType: DefaultResponseType,
		Phase:        PhaseBefore,
		Paths:        []Path{},
		Parameters:   []Parameter{},
		Examples:     []string{},
		Legacy:       []string{},
		Body:         Body{Type: BodyNone, MaxBytes: DefaultMaxBytes},
	}
	BindMethods(sandbox, conf)
	return conf
}
