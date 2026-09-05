package utils

import (
	"sort"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// ExamplesDir is the example tree of a project: one directory per side, each
// holding one directory per example. Every example is both documentation and
// a test — `exec-test` runs it and compares what it produced with its golden.
const ExamplesDir = "examples"

// ExampleCliSide and ExampleLibSide are the two sides of ExamplesDir: the one
// exercised through the project's cli and the one exercised through its lib.
const (
	ExampleCliSide = "cli"
	ExampleLibSide = "lib"
)

// ExampleCliFile and ExampleLibFile are the one file each side's example
// directory holds — the example itself, run with the directory as its cwd.
const (
	ExampleCliFile = "example.sh"
	ExampleLibFile = "example.go"
)

// ExampleResultFile is the golden one example run is compared against, and
// ExampleTestDir the only directory an example may write into.
const (
	ExampleResultFile = "result.yaml"
	ExampleTestDir    = "TestDir"
)

// ExampleSides is the fixed order the two sides are walked in: an example
// present on both sides is run through the cli first.
var ExampleSides = []string{ExampleCliSide, ExampleLibSide}

// ExampleSideDir is the project-relative directory holding one side's
// examples ("cli" -> "examples/cli").
func ExampleSideDir(side string) string {
	return ExamplesDir + "/" + side
}

// ExampleDir is the project-relative directory of one example
// ("cli", "start" -> "examples/cli/start").
func ExampleDir(side string, name string) string {
	return ExampleSideDir(side) + "/" + name
}

// ExampleFile is the name of the one file an example of that side holds.
func ExampleFile(side string) string {
	if side == ExampleCliSide {
		return ExampleCliFile
	}
	return ExampleLibFile
}

// CollectExamples returns the names of the examples declared on one side of
// examples/, sorted. A project with no examples/ directory — every project
// before its first add-cli-example / add-lib-example — yields none, which is
// what the generated listing then documents.
func CollectExamples(io *smartio.SmartIO, side string) []string {
	dir := ExampleSideDir(side)
	if !io.IsDir(dir) {
		return nil
	}

	var names []string
	for _, entry := range io.ListDirs(dir) {
		names = append(names, LastSegment(entry))
	}
	sort.Strings(names)
	return names
}

// ValidateExampleName reports whether a user-typed example name is usable as
// a directory name under examples/<side>/. It becomes a directory and is
// linked from a generated listing, so only letters, digits, dots, dashes and
// underscores are allowed, and it is one segment — an example is never nested.
func ValidateExampleName(deps *deps.Deps, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return deps.Std.Errorf("an example needs a name")
	}

	for _, letter := range name {
		valid := (letter >= 'a' && letter <= 'z') ||
			(letter >= 'A' && letter <= 'Z') ||
			(letter >= '0' && letter <= '9') ||
			letter == '.' || letter == '-' || letter == '_'
		if !valid {
			return deps.Std.Errorf(
				"invalid example name %q: only letters, digits, dots, dashes and underscores are allowed (it becomes one directory under %s)",
				name, ExamplesDir)
		}
	}
	return nil
}

// RemoveExample deletes one example directory and everything in it — the
// example file, its golden result.yaml and any TestDir the last run left
// behind. It is the whole of both remove-*-example actions: the two differ
// only in the side they name.
func RemoveExample(deps *deps.Deps, io *smartio.SmartIO, side string, name string) error {
	if err := ValidateExampleName(deps, name); err != nil {
		return err
	}

	dir := ExampleDir(side, name)
	if !io.IsDir(dir) {
		return deps.Std.Errorf("example %s not found", dir)
	}

	deps.Std.Log("remove-%s-example removing %s \n", side, dir)

	for _, entry := range io.ListAllRecursively(dir) {
		io.RemoveDir(entry)
	}
	io.RemoveDir(dir)
	return nil
}
