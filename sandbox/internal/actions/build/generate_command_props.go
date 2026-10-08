package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// retiredCliFiles are generated files an older build wrote that nothing
// renders any more: the dispatch before it moved into the OpinionatedAgnosCli
// lib, the entries.yaml the build wrote for help and version before
// command.yaml, and the EntriesYaml doc CommandYaml replaced. Each one names a
// symbol the current contract dropped, so a tree still carrying it would not
// compile; it is removed on every build.
var retiredCliFiles = []string{
	utils.GeneratedDir + "/cli/main.go",
	"sandbox/internal/commands/help/entries.yaml",
	"sandbox/internal/commands/version/entries.yaml",
	"docs/EntriesYaml/doc.md",
	"docs/EntriesYaml/props.yaml",
	"docs/EntriesYaml/Index.md",
}

// GenerateCommandProps rewrites sandbox/internal/commandprops/commandprops.go
// — the CommandProps one command line's chain of commands shares, handed to
// every Handle as its first argument — as the embedding of every
// part the package declares beside it, the way GenerateRouteProps does for
// the server layer.
func GenerateCommandProps(sandbox *api.Sandbox, io *stagedfs.StagedFS, module string) error {
	for _, retired := range retiredCliFiles {
		if io.IsFile(retired) {
			io.RemoveDir(retired)
		}
	}

	return generatePropsAggregate(sandbox, io, propsAggregate{
		Dir:             utils.CommandPropsDir,
		File:            utils.CommandPropsFile,
		Type:            "CommandProps",
		Template:        "templates/commandprops.go",
		ProjectTemplate: "templates/commandprops_project.go",
	}, module)
}
