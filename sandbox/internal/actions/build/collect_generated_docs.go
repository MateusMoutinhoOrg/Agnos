package build

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/docpropsconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// generatedDocsGroup is the asset group whose docs every build writes, whether
// or not the project already has them.
const generatedDocsGroup = "all"

// CollectGeneratedDocs returns the first-level docs the `all` group writes on
// this build, read from that group's own props.yaml templates. They belong in
// the doc index like any other doc, but SmartIO listings read disk, so on a
// project's first build they are not there to be listed yet: without this the
// README of a freshly scaffolded project would index nothing. Any doc later
// added under assets/all/docs/ is picked up here on its own.
func CollectGeneratedDocs(deps *deps.Deps, io *smartio.SmartIO, vars map[string]interface{}) ([]utils.Doc, error) {
	files, err := deps.Embeddeps.ListFilesRecursively(generatedDocsGroup)
	if err != nil {
		return nil, err
	}

	var docs []utils.Doc

	for _, file := range files {
		dir, ok := generatedDocDir(file)
		if !ok {
			continue
		}

		src, err := deps.Embeddeps.ReadFile(generatedDocsGroup + "/" + file)
		if err != nil {
			return nil, err
		}

		rendered, err := utils.RenderTemplate(deps, io, utils.DocPropsFile, src, vars)
		if err != nil {
			return nil, err
		}

		props, err := docPropsOf(deps, string(rendered), file)
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
// declares, and whether the path is a first-level doc's props.yaml at all
// ("docs/PublicApi/props.yaml" -> "PublicApi").
func generatedDocDir(file string) (string, bool) {
	parts := strings.Split(file, "/")
	if len(parts) != 3 {
		return "", false
	}
	if parts[0] != utils.DocsDir || parts[2] != utils.DocPropsFile {
		return "", false
	}
	return parts[1], true
}

// docPropsOf parses one rendered props.yaml, naming the asset it came from
// when it does not parse.
func docPropsOf(deps *deps.Deps, content string, file string) (*docpropsconf.DocPropsConf, error) {
	conf, err := docpropsconf.New(deps, content)
	if err != nil {
		return nil, deps.Std.Errorf("assets/%s/%s: %w", generatedDocsGroup, file, err)
	}
	return conf, nil
}

// docsVars is the subset of the build's template vars a generated doc's
// props.yaml may use. The full var map cannot be handed over here: it carries
// the doc index, which is what these docs are being collected to build. A
// props.yaml that reaches for anything else renders it empty.
func docsVars(module string, name string) map[string]interface{} {
	return map[string]interface{}{
		"Module": module,
		"Name":   name,
	}
}
