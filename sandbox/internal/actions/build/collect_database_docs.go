package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// DatabaseDocField is one field of a table as docs/Databases prints it.
type DatabaseDocField struct {
	Name     string
	Type     string
	Required string
	Target   string
}

// DatabaseDocTable is one table's section of a database's page: its fields and
// the nested collections it holds.
type DatabaseDocTable struct {
	Name   string
	Type   string
	Fields []DatabaseDocField
}

// DatabaseDocMethod is one generated method as the page prints it: the whole
// signature, and the line saying what it does.
type DatabaseDocMethod struct {
	Name      string
	Signature string
	Help      string
}

// DatabaseDoc is one database's page of docs/Databases, rendered from its
// specs.yaml alone: nothing here is written by hand on the page.
type DatabaseDoc struct {
	Name    string
	Package string
	Type    string
	Prefix  string
	Tables  []DatabaseDocTable
	Methods []DatabaseDocMethod
}

// CollectDatabaseDocs renders every sandbox/internal/databases/<db>/specs.yaml
// into the page docs/Databases prints for it — the database layer's
// CollectRouteDocs. The declaration is the only source: a table, a field or a
// method reaches the page by being declared, never by the page being edited.
func CollectDatabaseDocs(sandbox *api.Sandbox, io *smartio.SmartIO) ([]DatabaseDoc, error) {
	var docs []DatabaseDoc

	for _, dir := range io.ListDirs(utils.DatabasesDir) {
		name := lastSegmentOf(sandbox, dir)
		if name == "" {
			continue
		}

		content, err := io.ReadFile(utils.DatabasesDir + "/" + name + "/" + utils.DatabaseSpecsFile)
		if err != nil {
			continue
		}

		conf, err := databaseconf.New(sandbox, string(content))
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("databases/%s/%s: %w", name, utils.DatabaseSpecsFile, err)
		}

		docs = append(docs, databaseDoc(sandbox, name, conf))
	}

	return docs, nil
}

// databaseDoc turns one parsed declaration into its page.
func databaseDoc(sandbox *api.Sandbox, name string, conf *databaseconf.DatabaseConf) DatabaseDoc {
	doc := DatabaseDoc{
		Name:    utils.DatabaseIdentifier(sandbox, name),
		Package: name,
		Type:    utils.ExportedName(sandbox, name),
		Prefix:  conf.Prefix,
	}

	for _, table := range conf.Tables {
		section := DatabaseDocTable{Name: table.Name, Type: utils.ExportedName(sandbox, table.Name)}
		for _, field := range table.Fields {
			section.Fields = append(section.Fields, databaseDocField(field, ""))
			for _, nested := range field.Fields {
				section.Fields = append(section.Fields, databaseDocField(nested, field.Name+"."))
			}
		}
		doc.Tables = append(doc.Tables, section)

		for _, method := range utils.DatabaseMethods(sandbox, table) {
			doc.Methods = append(doc.Methods, databaseDocMethod(method))
		}
	}

	return doc
}

// databaseDocField renders one field as its table row. A field of a nested
// collection is listed under the collection it belongs to, by the dotted path
// the declaration gives it.
func databaseDocField(field databaseconf.Field, prefix string) DatabaseDocField {
	required := ""
	if field.Required {
		required = "yes"
	}

	target := ""
	if field.Target != "" {
		target = "`" + field.Target + "`"
	}

	return DatabaseDocField{
		Name:     prefix + field.Name,
		Type:     "`" + field.Type + "`",
		Required: required,
		Target:   target,
	}
}

// databaseDocMethod renders one generated method as the signature the page
// prints. It is the same derivation the three templates are rendered from, so
// the page and the code can never disagree about a signature.
func databaseDocMethod(method utils.DatabaseMethod) DatabaseDocMethod {
	return DatabaseDocMethod{
		Name:      method.Name,
		Signature: method.Name + "(" + method.Params + ") " + method.Results,
		Help:      method.Help,
	}
}
