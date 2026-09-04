package utils

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/structureconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// StructureConfFile is the declaration docs/Structure is rendered from, held
// in the project's config directory beside project.yaml and themes.yaml.
const StructureConfFile = "structure.yaml"

// StructureNode is one item of the structure declaration, flattened out of the
// tree in the order it renders: depth-first, siblings by `order` then name.
type StructureNode struct {
	// Path is the item's project-relative path, every ancestor's name joined
	// by "/" ("sandbox/internal/config.go"). Label is the last segment as the
	// tree prints it, with a trailing "/" when the item is a directory.
	Path  string
	Label string
	Depth int

	Description string
	Dir         bool
	Gen         bool
}

// IsStructurePattern reports whether a declared path stands for a family of
// paths rather than one ("libs/<lib>/<lib>.go", "collect_*.go"). Such an item
// has nothing to exist on disk under its own name, so `verify` checks the
// literal part of the path instead (see StructureParentPath).
func IsStructurePattern(path string) bool {
	return strings.ContainsAny(path, "<>*?")
}

// StructureConfPath is the project-relative path of the structure declaration.
func StructureConfPath() string {
	return config.ProjectName + "Config/" + StructureConfFile
}

// LoadStructureConf reads <ProjectName>Config/structure.yaml through the
// transaction-aware io. Unlike project.yaml and themes.yaml it may be absent —
// a project scaffolded by an older agnos has none — and then describes nothing
// rather than failing the build. A file that is there but does not parse is a
// hard error.
func LoadStructureConf(deps *deps.Deps, io *smartio.SmartIO) (*structureconf.StructureConf, error) {
	rel := StructureConfPath()

	content, err := io.ReadFile(rel)
	if err != nil {
		return structureconf.NewEmpty(deps), nil
	}

	conf, err := structureconf.New(deps, string(content))
	if err != nil {
		return nil, deps.Std.Errorf("%s: %w", rel, err)
	}
	return conf, nil
}

// FlattenStructure walks the declared tree and returns one node per item, in
// render order. It is the single walk both readers of the declaration share:
// `build` turns the nodes into the lines of docs/Structure, and `verify` turns
// them into the paths a ghost spec would name.
func FlattenStructure(items []structureconf.Item) []StructureNode {
	return flattenStructureIn(items, "", 0)
}

// flattenStructureIn appends the items directly under parent, then everything
// nested beneath each of them.
func flattenStructureIn(items []structureconf.Item, parent string, depth int) []StructureNode {
	var nodes []StructureNode

	for _, item := range items {
		path := item.Name
		if parent != "" {
			path = parent + "/" + item.Name
		}

		label := item.Name
		if item.Dir {
			label += "/"
		}

		nodes = append(nodes, StructureNode{
			Path:        path,
			Label:       label,
			Depth:       depth,
			Description: item.Description,
			Dir:         item.Dir,
			Gen:         item.Gen,
		})

		nodes = append(nodes, flattenStructureIn(item.Children, path, depth+1)...)
	}

	return nodes
}

// StructureParentPath is the path a pattern node is checked against: the
// longest leading run of segments that holds no pattern character. It is empty
// when the very first segment is already a pattern.
func StructureParentPath(path string) string {
	var literal []string

	for _, segment := range strings.Split(path, "/") {
		if strings.ContainsAny(segment, "<>*?") {
			break
		}
		literal = append(literal, segment)
	}

	return strings.Join(literal, "/")
}
