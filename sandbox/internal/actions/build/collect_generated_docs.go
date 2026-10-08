package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/docconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CollectGeneratedDocs returns the first-level docs the given groups write on
// this build, read from those groups' own doc.yaml templates. They belong in
// the doc index like any other doc, but StagedFS listings read disk, so on a
// project's first build they are not there to be listed yet: without this the
// README of a freshly scaffolded project would index nothing. Any doc later
// added under assets/<group>/docs/ is picked up here on its own.
func CollectGeneratedDocs(sandbox *api.Sandbox, io *stagedfs.StagedFS, vars map[string]interface{}, groups []string) ([]utils.Doc, error) {
	var docs []utils.Doc

	for _, group := range groups {
		group_docs, err := collectGroupDocs(sandbox, io, vars, group)
		if err != nil {
			return nil, err
		}
		docs = append(docs, group_docs...)
	}

	return docs, nil
}

// collectGroupDocs is CollectGeneratedDocs over one asset group.
func collectGroupDocs(sandbox *api.Sandbox, io *stagedfs.StagedFS, vars map[string]interface{}, group string) ([]utils.Doc, error) {
	files, err := sandbox.Deps.EmbedDeps.ListFilesRecursively(group)
	if err != nil {
		return nil, err
	}

	var docs []utils.Doc

	for _, file := range files {
		dir, ok := generatedDocDir(sandbox, file)
		if !ok {
			continue
		}

		src, err := sandbox.Deps.EmbedDeps.ReadFile(group + "/" + file)
		if err != nil {
			return nil, err
		}

		rendered, err := utils.RenderTemplate(sandbox, io, utils.DocConfFile, src, vars)
		if err != nil {
			return nil, err
		}

		props, err := docConfOf(sandbox, string(rendered), group, file)
		if err != nil {
			return nil, err
		}

		doc := utils.Doc{
			Dir:         dir,
			Path:        utils.DocsDir + "/" + dir,
			Name:        props.Name,
			Description: props.Description,
			Themes:      props.Themes,
			Order:       props.Order,
			HasOrder:    props.HasOrder,
		}
		if doc.Name == "" {
			doc.Name = dir
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

// generatedDocDir reports the doc directory a group-relative asset path
// declares, and whether the path is a first-level doc's doc.yaml at all
// ("docs/PublicApi/doc.yaml" -> "PublicApi").
func generatedDocDir(sandbox *api.Sandbox, file string) (string, bool) {
	parts := sandbox.Deps.StringsDeps.Split(file, "/")
	if len(parts) != 3 {
		return "", false
	}
	if parts[0] != utils.DocsDir || parts[2] != utils.DocConfFile {
		return "", false
	}
	return parts[1], true
}

// docConfOf parses one rendered doc.yaml, naming the asset it came from
// when it does not parse.
func docConfOf(sandbox *api.Sandbox, content string, group string, file string) (*docconf.DocConf, error) {
	conf, err := docconf.New(sandbox, content)
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("assets/%s/%s: %w", group, file, err)
	}
	return conf, nil
}

// docsVars is the subset of the build's template vars a generated doc's
// doc.yaml may use. The full var map cannot be handed over here: it carries
// the doc index, which is what these docs are being collected to build. A
// doc.yaml that reaches for anything else renders it empty.
func docsVars(module string, name string, generator_name string) map[string]interface{} {
	return map[string]interface{}{
		"Module":        module,
		"ProjectName":   name,
		"GeneratorName": generator_name,
	}
}
