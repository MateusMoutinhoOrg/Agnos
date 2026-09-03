package frozen

// This file is an illustrative sample, not part of the build.

import (
	"time"

	clockdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/clockdeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/iodeps"
	std "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/std"
	verb "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/verb"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

// New assembles a Deps whose clock is frozen at the given time: the real
// libs for everything the sandbox cannot do without, then one override.
func New(now time.Time) deps.Deps {
	deps := deps.Deps{}
	clockdeps.Bind(&deps)
	iodeps.Bind(&deps)
	std.Bind(&deps)
	verb.Bind(&deps)

	// Override after the binds, so the lib's other fields stay real.
	deps.Clockdeps.Now = func() time.Time { return now }
	return deps
}
