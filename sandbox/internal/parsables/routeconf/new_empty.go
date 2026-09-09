package routeconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *RouteConf {
	conf := &RouteConf{
		Method:   DefaultMethod,
		Paths:    []Segment{},
		Examples: []string{},
		Headers:  []Field{},
		Params:   []Field{},
		Body:     Body{Type: BodyNone, MaxBytes: DefaultMaxBytes},
	}
	BindMethods(deps, conf)
	return conf
}
