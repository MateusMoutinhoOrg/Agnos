package update_tests

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

// The package is update_tests, not update_test, for the same reason exec_tests
// is plural: Go reserves every file whose name ends in _test.go for the
// testing toolchain, so an action directory named after the `update-test`
// command could not hold its own <name>.go. The command, the api field and the
// docs all keep the name.

// UpdateTest rewrites the goldens of one example with what it produces now,
// printing what each write changed. It is `exec-test --update` narrowed to a
// single name, which is what makes an update reviewable: a suite-wide rewrite
// hides the one golden that moved for a reason nobody meant.
//
// The name is required. Without it the command is `exec-test --update` under
// another spelling, and one golden at a time is the whole of it.
func UpdateTest(deps *deps.Deps, path string, name string) error {
	if deps.Stringsdeps.TrimSpace(name) == "" {
		return deps.Std.Errorf("update-test: an example name is required (rewrite every golden with `exec-test --update`)")
	}

	deps.Std.Log("update-test started with path %s \n", path)

	return UpdateTestInternal(deps, path, name)
}
