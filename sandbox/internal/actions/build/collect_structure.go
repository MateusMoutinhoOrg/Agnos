package build

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// structureIndent is the run of spaces one level of nesting adds to a line of
// the rendered tree.
const structureIndent = "  "

// structureGap is the minimum run of spaces between the longest label of the
// tree and the descriptions column.
const structureGap = 2

// StructureEntry is one rendered line of docs/Structure's tree: the item's
// label, indented by its depth and padded so every description in the block
// starts at the same column.
type StructureEntry struct {
	Line string
}

// CollectStructure renders <ProjectName>Config/structure.yaml into the lines of
// the tree docs/Structure prints. The declaration is the only source: nothing
// is read off disk here, so an item that names a pattern renders like any
// other. `verify` is what keeps the declaration and the disk in step.
//
// A project with no structure.yaml collects no lines, and its Structure doc
// renders an empty tree.
func CollectStructure(deps *deps.Deps, io *smartio.SmartIO) ([]StructureEntry, error) {
	structure_conf, err := utils.LoadStructureConf(deps, io)
	if err != nil {
		return nil, err
	}

	nodes := utils.FlattenStructure(structure_conf.Items)

	labels := make([]string, len(nodes))
	width := 0
	for index, node := range nodes {
		labels[index] = strings.Repeat(structureIndent, node.Depth) + node.Label
		if len(labels[index]) > width {
			width = len(labels[index])
		}
	}

	entries := make([]StructureEntry, 0, len(nodes))
	for index, node := range nodes {
		entries = append(entries, StructureEntry{
			Line: structureLine(labels[index], width, node),
		})
	}

	return entries, nil
}

// structureLine pads one label out to the descriptions column and appends the
// item's description, marking a generated file with the "(gen)" prefix every
// page of the docs uses for one. An item with no description is just its
// label, with no trailing spaces to diff against.
func structureLine(label string, width int, node utils.StructureNode) string {
	description := node.Description
	if node.Gen {
		description = strings.TrimSpace("(gen) " + description)
	}

	if description == "" {
		return label
	}

	padding := width - len(label) + structureGap
	return label + strings.Repeat(" ", padding) + description
}
