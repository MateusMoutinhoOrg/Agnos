package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RouteDocField is one path slice or parameter as docs/Routes prints it: the
// Entries field it binds, where it is read from, a single label carrying its
// type and its conditions, the default, and the declared description.
type RouteDocField struct {
	Id          string
	Key         string
	In          string
	Type        string
	Default     string
	Description string
}

// RouteDoc is one route's section of docs/Routes, rendered from its route.yaml
// alone: nothing here is written by hand on the page.
type RouteDoc struct {
	Name            string
	Method          string
	Pattern         string
	Help            string
	LongDescription string
	Fields          []RouteDocField
	Body            string
	Examples        []string
}

// RouteDocGroup is one category section of docs/Routes, holding the routes
// that declare that category.
type RouteDocGroup struct {
	Category string
	Routes   []RouteDoc
}

// routeDocOther is the category a route with no declared one falls into.
const routeDocOther = "Other"

// CollectRouteDocs renders every sandbox/internal/routeslist/<name>/route.yaml into
// the sections docs/Routes prints, grouped by category in first-seen order —
// the server layer's CollectCommandDocs. Hidden routes are skipped.
//
// The declaration is the only source: a route, a field or an example reaches
// the page by being declared with `add-route`, `add-path`, `add-parameter` or
// `set-route`,
// never by the page being edited.
func CollectRouteDocs(sandbox *api.Sandbox, io *smartio.SmartIO) ([]RouteDocGroup, error) {
	var groups []RouteDocGroup
	index := map[string]int{}

	for _, dir := range io.ListDirs(routesDir) {
		name := lastSegmentOf(sandbox, dir)
		if name == "" {
			continue
		}

		content, err := io.ReadFile(routesDir + "/" + name + "/route.yaml")
		if err != nil {
			continue
		}

		conf, err := routeconf.New(sandbox, string(content))
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("routeslist/%s/route.yaml: %w", name, err)
		}

		if conf.Hidden {
			continue
		}

		category := conf.Category
		if category == "" {
			category = routeDocOther
		}

		position, seen := index[category]
		if !seen {
			position = len(groups)
			index[category] = position
			groups = append(groups, RouteDocGroup{Category: category})
		}

		groups[position].Routes = append(groups[position].Routes, routeDoc(sandbox, name, conf))
	}

	return groups, nil
}

// routeDoc turns one parsed declaration into its page section.
func routeDoc(sandbox *api.Sandbox, name string, conf *routeconf.RouteConf) RouteDoc {
	doc := RouteDoc{
		Name:            name,
		Method:          sandbox.Deps.Stringsdeps.Join(conf.Methods, ", "),
		Pattern:         conf.Pattern(),
		Help:            docCell(sandbox, conf.Help),
		LongDescription: docText(sandbox, conf.LongDescription),
		Body:            routeDocBody(sandbox, conf.Body),
		Examples:        conf.Examples,
	}

	for _, path := range conf.Paths {
		doc.Fields = append(doc.Fields, RouteDocField{
			Id:          path.Id,
			Key:         routeDocSlice(sandbox, path),
			In:          "path",
			Type:        path.Type + routeDocTrigger(path.Trigger),
			Description: docCell(sandbox, path.Description),
		})
	}
	for _, parameter := range conf.Parameters {
		value := ""
		if parameter.HasDefault {
			value = "`" + parameter.Default + "`"
		}
		label := parameter.Type
		if parameter.Required {
			label += ", required"
		}
		doc.Fields = append(doc.Fields, RouteDocField{
			Id:          parameter.Id,
			Key:         "`" + parameter.Key + "`",
			In:          sandbox.Deps.Stringsdeps.Join(parameter.Fonts, ", "),
			Type:        label + routeDocTrigger(parameter.Trigger),
			Default:     value,
			Description: docCell(sandbox, parameter.Description),
		})
	}

	return doc
}

// routeDocSlice spells the segments a path reads: "0..-1" is the whole path.
func routeDocSlice(sandbox *api.Sandbox, path routeconf.Path) string {
	return "segments " + sandbox.Deps.Stringsdeps.FormatInt(int64(path.Start), 10) + ".." +
		sandbox.Deps.Stringsdeps.FormatInt(int64(path.End), 10)
}

// routeDocTrigger is the part of a type label a trigger adds, "" when there is
// none.
func routeDocTrigger(trigger routeconf.Trigger) string {
	if !trigger.Exists {
		return ""
	}
	text := ", " + trigger.Type + " `" + trigger.Value + "`"
	if trigger.IgnoreCase {
		text += " ignoring case"
	}
	if trigger.Negate {
		text = ", not" + text[1:]
	}
	return text
}

// routeDocBody is the one line describing a route's body: its kind, whether it
// is required, the content type it accepts and the size it stops at.
func routeDocBody(sandbox *api.Sandbox, body routeconf.Body) string {
	if body.Type == routeconf.BodyNone {
		return ""
	}

	label := "`" + body.Type + "`"
	if body.Required {
		label += ", required"
	}
	if body.ContentType != "" {
		label += ", `" + body.ContentType + "`"
	}
	if body.HasSchema {
		label += ", json-schema"
	}
	return label + ", up to " + sandbox.Deps.Stringsdeps.FormatInt(int64(body.MaxBytes), 10) + " bytes"
}
