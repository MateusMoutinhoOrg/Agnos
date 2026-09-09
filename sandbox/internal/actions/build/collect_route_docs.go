package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RouteDocField is one captured segment, header or query parameter as
// docs/Routes prints it: where it is read from, a single label carrying type,
// arity and bounds, the default, and the declared description.
type RouteDocField struct {
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

// CollectRouteDocs renders every sandbox/internal/routes/<name>/route.yaml into
// the sections docs/Routes prints, grouped by category in first-seen order —
// the server layer's CollectCommandDocs. Hidden routes are skipped.
//
// The declaration is the only source: a route, a field or an example reaches
// the page by being declared with `add-route`, `add-field` or `set-route`,
// never by the page being edited.
func CollectRouteDocs(deps *deps.Deps, io *smartio.SmartIO) ([]RouteDocGroup, error) {
	var groups []RouteDocGroup
	index := map[string]int{}

	for _, dir := range io.ListDirs(routesDir) {
		name := lastSegmentOf(deps, dir)
		if name == "" {
			continue
		}

		content, err := io.ReadFile(routesDir + "/" + name + "/route.yaml")
		if err != nil {
			continue
		}

		conf, err := routeconf.New(deps, string(content))
		if err != nil {
			return nil, deps.Std.Errorf("routes/%s/route.yaml: %w", name, err)
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

		groups[position].Routes = append(groups[position].Routes, routeDoc(deps, name, conf))
	}

	return groups, nil
}

// routeDoc turns one parsed declaration into its page section.
func routeDoc(deps *deps.Deps, name string, conf *routeconf.RouteConf) RouteDoc {
	doc := RouteDoc{
		Name:            name,
		Method:          conf.Method,
		Pattern:         conf.Pattern(),
		Help:            docCell(deps, conf.Help),
		LongDescription: docText(deps, conf.LongDescription),
		Body:            routeDocBody(deps, conf.Body),
		Examples:        conf.Examples,
	}

	for _, segment := range conf.Paths {
		if segment.Field == nil {
			continue
		}
		doc.Fields = append(doc.Fields, routeDocField(deps, *segment.Field, "path"))
	}
	for _, field := range conf.Headers {
		doc.Fields = append(doc.Fields, routeDocField(deps, field, "header"))
	}
	for _, field := range conf.Params {
		doc.Fields = append(doc.Fields, routeDocField(deps, field, "query"))
	}

	return doc
}

// routeDocField renders one field as its table row.
func routeDocField(deps *deps.Deps, field routeconf.Field, in string) RouteDocField {
	value := ""
	if field.HasDefault {
		value = "`" + field.Default + "`"
	}

	return RouteDocField{
		Key:         field.Key,
		In:          in,
		Type:        routeFieldTypeLabel(deps, field),
		Default:     value,
		Description: docCell(deps, field.Description),
	}
}

// routeFieldTypeLabel is the one cell carrying everything the type of a field
// implies: its kind, whether it repeats, whether it must be given, and the
// bounds a numeric field declares.
func routeFieldTypeLabel(deps *deps.Deps, field routeconf.Field) string {
	label := field.Type
	if label == "" {
		label = "string"
	}
	if field.Array {
		label += ", repeatable"
	}
	if field.Required {
		label += ", required"
	}
	if bounds := routeFieldBounds(deps, field); bounds != "" {
		label += ", " + bounds
	}
	return label
}

// routeFieldBounds spells the min/max a numeric field declares, "" when it
// declares neither.
func routeFieldBounds(deps *deps.Deps, field routeconf.Field) string {
	switch {
	case field.HasMin && field.HasMax:
		return routeNumberLabel(deps, field.Type, field.Min, true) + ".." + routeNumberLabel(deps, field.Type, field.Max, true)
	case field.HasMin:
		return ">= " + routeNumberLabel(deps, field.Type, field.Min, true)
	case field.HasMax:
		return "<= " + routeNumberLabel(deps, field.Type, field.Max, true)
	}
	return ""
}

// routeDocBody is the one line describing a route's body: its kind, whether it
// is required, the content type it accepts and the size it stops at.
func routeDocBody(deps *deps.Deps, body routeconf.Body) string {
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
	return label + ", up to " + deps.Stringsdeps.FormatInt(int64(body.MaxBytes), 10) + " bytes"
}
