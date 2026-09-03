//go:build ignore

package api

// This file is an illustrative sample, not part of the build.

// The runtimes `build` can hand a rendered project to.
const (
	RuntimeGo   = "go"
	RuntimeNone = "none"
)

// BuildProps describes one (re)render of a project: the directory holding it
// and the runtime that then checks the result.
type BuildProps struct {
	Path    string
	Runtime string
}

// Actions is the contract of every operation the sandbox performs. Each field
// is filled by ActionsBind in sandbox/binds/actions.go.
type Actions struct {
	Build   func(props BuildProps) error
	Verify  func(path string) error
	DepList func(path string) ([]string, error)
}