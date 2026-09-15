package routeconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *RouteConf {
	conf := &RouteConf{
		Method:   DefaultMethod,
		Paths:    []Segment{},
		Examples: []string{},
		Headers:  []Field{},
		Params:   []Field{},
		Body:     Body{Type: BodyNone, MaxBytes: DefaultMaxBytes},
	}
	BindMethods(sandbox, conf)
	return conf
}
