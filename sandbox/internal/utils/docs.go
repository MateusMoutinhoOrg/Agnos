package utils

import (
	"sort"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/docpropsconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/themesconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// DocsDir is the documentation tree of a project: one directory per doc, each
// holding doc.md + props.yaml and, recursively, its sub-docs.
const DocsDir = "docs"

// DocsIndexName is the reserved first-level name of the generated theme
// indexes (docs/Index/<theme-id>.md). It is never a doc.
const DocsIndexName = "Index"

// DocIndexFile is the index a doc directory gets when it has sub-docs.
const DocIndexFile = "Index.md"

// DocFile and DocPropsFile are the two files every doc directory holds.
const (
	DocFile      = "doc.md"
	DocPropsFile = "props.yaml"
)

// Doc is one node of the documentation tree: a directory holding doc.md and
// props.yaml, plus the sub-docs nested under it. Any other file in the
// directory is an asset and is ignored.
type Doc struct {
	// Dir is the doc's own directory name ("PublicApi"); Path is its
	// project-relative path ("docs/PublicApi", "docs/PublicApi/api.Actions").
	Dir  string
	Path string

	Name        string
	Description string
	Themes      []string
	Order       int
	HasOrder    bool

	Subdocs []Doc
}

// LoadThemesConf reads <ProjectName>Config/themes.yaml through the
// transaction-aware io. Like project.yaml it is written by `agnos start`, so a
// missing or unparsable file is a hard error rather than an empty fallback.
func LoadThemesConf(deps *deps.Deps, io *smartio.SmartIO) (*themesconf.ThemesConf, error) {
	rel := config.ProjectName + "Config/themes.yaml"

	content, err := io.ReadFile(rel)
	if err != nil {
		return nil, deps.Std.Errorf("could not read %s: run `agnos start` first (%w)", rel, err)
	}

	return themesconf.New(deps, string(content))
}

// LoadDocProps reads and parses one doc directory's props.yaml. A missing or
// unparsable file is an error: the doc tree is indexed from these files, so a
// doc without them cannot be listed anywhere.
func LoadDocProps(deps *deps.Deps, io *smartio.SmartIO, doc_path string) (*docpropsconf.DocPropsConf, error) {
	rel := doc_path + "/" + DocPropsFile

	content, err := io.ReadFile(rel)
	if err != nil {
		return nil, deps.Std.Errorf("%s is missing", rel)
	}

	conf, err := docpropsconf.New(deps, string(content))
	if err != nil {
		return nil, deps.Std.Errorf("%s: %w", rel, err)
	}
	return conf, nil
}

// CollectDocTree walks docs/ and returns its first-level docs, each carrying
// its sub-docs recursively. Every directory is a doc — the reserved first-level
// Index/ aside — so a missing props.yaml is an error. A project with no docs/
// directory yields an empty tree and no error.
func CollectDocTree(deps *deps.Deps, io *smartio.SmartIO) ([]Doc, error) {
	if !io.IsDir(DocsDir) {
		return nil, nil
	}
	return collectDocsIn(deps, io, DocsDir, true)
}

// collectDocsIn collects the docs directly under parent_path. first_level
// marks the docs/ directory itself, where Index/ is reserved for the generated
// theme indexes.
func collectDocsIn(deps *deps.Deps, io *smartio.SmartIO, parent_path string, first_level bool) ([]Doc, error) {
	var docs []Doc

	for _, dir := range io.ListDirs(parent_path) {
		name := LastSegment(dir)
		if name == "" {
			continue
		}
		if first_level && name == DocsIndexName {
			continue
		}

		doc_path := parent_path + "/" + name
		props, err := LoadDocProps(deps, io, doc_path)
		if err != nil {
			return nil, err
		}

		subdocs, err := collectDocsIn(deps, io, doc_path, false)
		if err != nil {
			return nil, err
		}

		doc := Doc{
			Dir:         name,
			Path:        doc_path,
			Name:        props.Name,
			Description: props.Description,
			Themes:      props.Themes,
			Order:       props.Order,
			HasOrder:    props.HasOrder,
			Subdocs:     subdocs,
		}
		if doc.Name == "" {
			doc.Name = name
		}

		docs = append(docs, doc)
	}

	SortDocs(docs)
	return docs, nil
}

// SortDocs orders docs the way every index lists them: by `order`, then by
// name. A doc with no `order` comes after every ordered one.
func SortDocs(docs []Doc) {
	sort.SliceStable(docs, func(i, j int) bool {
		left, right := docs[i], docs[j]
		if left.HasOrder != right.HasOrder {
			return left.HasOrder
		}
		if left.HasOrder && left.Order != right.Order {
			return left.Order < right.Order
		}
		return left.Name < right.Name
	})
}

// DocSegments splits a user-typed doc name into the directory segments it
// names under docs/ ("PublicApi/api.Actions" -> ["PublicApi", "api.Actions"]).
// Empty segments — a leading, trailing or doubled slash — are dropped, so the
// result is the doc's path relative to docs/.
func DocSegments(name string) []string {
	var segments []string
	for _, segment := range strings.Split(strings.TrimSpace(name), "/") {
		segment = strings.TrimSpace(segment)
		if segment != "" {
			segments = append(segments, segment)
		}
	}
	return segments
}

// ValidateDocName reports whether a user-typed doc name is usable as a
// directory path under docs/. Each segment becomes a directory name and is
// linked from a generated index, so only letters, digits, dots, dashes and
// underscores are allowed, and Index is reserved at the first level for the
// generated theme indexes.
func ValidateDocName(deps *deps.Deps, name string) error {
	segments := DocSegments(name)
	if len(segments) == 0 {
		return deps.Std.Errorf("a doc needs a name")
	}

	for index, segment := range segments {
		if index == 0 && segment == DocsIndexName {
			return deps.Std.Errorf("invalid doc name %q: %s is reserved for the generated theme indexes", name, DocsIndexName)
		}
		for _, letter := range segment {
			valid := (letter >= 'a' && letter <= 'z') ||
				(letter >= 'A' && letter <= 'Z') ||
				(letter >= '0' && letter <= '9') ||
				letter == '.' || letter == '-' || letter == '_'
			if !valid {
				return deps.Std.Errorf(
					"invalid doc name %q: only letters, digits, dots, dashes, underscores and / are allowed (it becomes the directory %s)",
					name, DocDir(name))
			}
		}
	}
	return nil
}

// DocDir is the project-relative directory holding a doc
// ("PublicApi/api.Actions" -> "docs/PublicApi/api.Actions").
func DocDir(name string) string {
	return DocsDir + "/" + strings.Join(DocSegments(name), "/")
}

// DocParentDir is the project-relative directory holding a doc's parent — the
// docs/ directory itself for a first-level doc.
func DocParentDir(name string) string {
	segments := DocSegments(name)
	if len(segments) < 2 {
		return DocsDir
	}
	return DocsDir + "/" + strings.Join(segments[:len(segments)-1], "/")
}

// DocTitle is the human-readable name a new doc is declared with: its last
// segment, with CamelCase split into words ("HandleDocuments" -> "Handle
// Documents"). A segment already carrying its own punctuation — a dot, dash,
// underscore or space, as in "api.Actions" — is kept verbatim.
func DocTitle(name string) string {
	segments := DocSegments(name)
	if len(segments) == 0 {
		return ""
	}
	title := segments[len(segments)-1]

	if strings.ContainsAny(title, ".-_ ") {
		return title
	}

	// A word starts where a lowercase letter or a digit is followed by an
	// uppercase one, so an acronym stays whole ("SmartIO" -> "Smart IO").
	var out strings.Builder
	letters := []rune(title)
	for index, letter := range letters {
		if index > 0 && letter >= 'A' && letter <= 'Z' {
			previous := letters[index-1]
			if (previous >= 'a' && previous <= 'z') || (previous >= '0' && previous <= '9') {
				out.WriteString(" ")
			}
		}
		out.WriteRune(letter)
	}
	return out.String()
}

// LastSegment is the final element of a project-relative path
// ("docs/PublicApi" -> "PublicApi").
func LastSegment(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
