package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// generatedMarker is the phrase every file the build rewrites carries in its
// doc comment. A props file without it is one an older build wrote once and
// the project then typed by hand.
const generatedMarker = "do not edit by hand"

// propsAggregate describes one props package: the struct a chain of handlers
// shares, generated as the embedding of every part declared beside it.
type propsAggregate struct {
	// Dir is the package's project-relative directory.
	Dir string
	// File is the aggregate's file name inside Dir, before
	// utils.GeneratedPrefix.
	File string
	// Type is the aggregate's type name: RouteProps, CommandProps.
	Type string
	// Template renders the aggregate.
	Template string
	// ProjectTemplate renders the project's own part, project.go.
	ProjectTemplate string
}

// projectPropsFile is the part a props package starts with: the project's own,
// written once, where its middlewares declare what they hand on.
const projectPropsFile = "project.go"

// projectPropsType is the struct projectPropsFile declares.
const projectPropsType = "Project"

// generatePropsAggregate rewrites <Dir>/generated.<File> as the struct embedding every
// exported struct of every other file of the package, sorted by name.
//
// A part is a file of its own so that a mechanic — backoffice-init writing
// backoffice.go — adds to what a handler is handed without editing a file the
// project wrote in. The aggregate is the one file no part owns, so it is the
// one the build rewrites.
//
// Two steps run first. A <File> an older build wrote once, and the project
// may have typed since, is moved to project.go (one the build wrote, marked,
// is dropped when generated.<File> is written) with its struct renamed
// Project: every field it declared is still promoted, so props.User reads
// the same. And a package left with no part at all is given an empty
// project.go, the place the project declares its own fields.
func generatePropsAggregate(sandbox *api.Sandbox, io *stagedfs.StagedFS, props propsAggregate, module string) error {
	dest := props.Dir + "/" + utils.GeneratedFile(sandbox, props.File)
	project := props.Dir + "/" + projectPropsFile

	if err := migrateHandWrittenProps(sandbox, io, props, props.Dir+"/"+props.File, project); err != nil {
		return err
	}

	vars := map[string]any{
		"Module":        module,
		"GeneratorName": generatorName(sandbox),
	}

	structs, err := utils.CollectEmbeddedStructs(sandbox, io, props.Dir, utils.AllBut(props.File), []string{project})
	if err != nil {
		return err
	}
	if len(structs) == 0 {
		if err := utils.RenderTemplateToDest(sandbox, io, props.ProjectTemplate, vars, project); err != nil {
			return err
		}
		structs = []utils.EmbeddedStruct{{Name: projectPropsType, File: project}}
	}

	vars["Structs"] = structs
	return utils.RenderTemplateToDest(sandbox, io, props.Template, vars, dest)
}
