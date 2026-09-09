package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// DefaultStaticMount is the url prefix front-init declares for the static
// route, and what pageio is rendered with when the project has no static route
// to read one off.
const DefaultStaticMount = "/static"

// CollectFrontMount returns the url prefix the static route answers under, for
// the StaticMount constant of the generated sandbox/internal/pageio.
//
// The constant is generated but the prefix is the project's: the static route
// is written once and then editable, so renaming its first segment has to move
// every link the pageio helpers build. Reading it back from the declaration is
// what makes that rename a rebuild rather than a silent breakage — every
// staticref would otherwise keep pointing at a mount nothing answers on.
//
// A missing, unparsable or capture-first static route falls back to
// DefaultStaticMount: this runs on every build, so it reports nothing and
// leaves the complaining to CollectRoutes, which parses the same file.
func CollectFrontMount(deps *deps.Deps, io *smartio.SmartIO) string {
	content, err := io.ReadFile(utils.RouteConfPath(deps, utils.StaticRouteName))
	if err != nil {
		return DefaultStaticMount
	}

	conf, err := routeconf.New(deps, string(content))
	if err != nil {
		return DefaultStaticMount
	}

	if len(conf.Paths) == 0 || conf.Paths[0].Identifier == "" {
		return DefaultStaticMount
	}

	return conf.Paths[0].Identifier
}
