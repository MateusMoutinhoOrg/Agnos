package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializables "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// databaseIdField is the one field name a table may not declare: every
// <T>Item already carries the record's permanent Id, so a field of that name
// would generate a struct with the same field twice.
const databaseIdField = "id"

// CheckDatabases enforces the shape the database layer's generators read by
// convention: the four generated files of a database package, a parsable
// declaration, links that point at a table, nested collections that hold
// fields and nest no further, and a hand-written methods_custom.go that does
// not redeclare anything the build writes.
//
// A project with no sandbox/internal/databases has no database layer and
// nothing to check.
func CheckDatabases(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	var violations []string

	if !io.IsDir(utils.DatabasesDir) {
		return violations
	}

	for _, dir := range io.ListDirs(utils.DatabasesDir) {
		name := databaseNameOf(sandbox, dir)
		if name == "" {
			continue
		}

		violations = append(violations, checkDatabaseFiles(sandbox, io, name)...)

		content, err := io.ReadFile(utils.DatabasesDir + "/" + name + "/" + utils.DatabaseSpecsFile)
		if err != nil {
			continue
		}

		conf, err := databaseconf.New(sandbox, string(content))
		if err != nil {
			violations = append(violations, databaseViolation(name, utils.DatabaseSpecsFile+" does not parse: "+err.Error()))
			continue
		}

		violations = append(violations, checkDatabaseDeclaration(sandbox, name, conf)...)
		violations = append(violations, checkDatabaseRawFields(sandbox, name, string(content))...)
		violations = append(violations, checkDatabaseCustom(sandbox, io, name, conf)...)
	}

	return violations
}

// checkDatabaseFiles reports a database package missing the declaration or any
// of the three files a build writes from it.
func checkDatabaseFiles(sandbox *api.Sandbox, io *smartio.SmartIO, name string) []string {
	var violations []string

	for _, file := range []string{utils.DatabaseSpecsFile, "api.go", "new.go", "methods.go"} {
		if !io.IsFile(utils.DatabasesDir + "/" + name + "/" + file) {
			violations = append(violations, databaseViolation(name, "has no "+file))
		}
	}

	return violations
}

// checkDatabaseDeclaration enforces the rules that survive parsing: the names
// a database spells once, the tables a link may point at, and the one level of
// nesting a collection is generated for.
func checkDatabaseDeclaration(sandbox *api.Sandbox, name string, conf *databaseconf.DatabaseConf) []string {
	var violations []string

	if conf.Name == "" {
		violations = append(violations, databaseViolation(name, utils.DatabaseSpecsFile+" declares no name"))
	}
	if len(conf.Tables) == 0 {
		violations = append(violations, databaseViolation(name, "declares no `tables`; a database needs at least one"))
	}

	// One Go record is declared per table and per nested collection, so the
	// two families share one namespace: two of them spelling the same name
	// would generate the same struct twice.
	records := map[string]string{}
	tables := map[string]bool{}

	for _, table := range conf.Tables {
		if tables[table.Name] {
			violations = append(violations, databaseViolation(name, "declares the table "+table.Name+" twice"))
		}
		tables[table.Name] = true
		violations = append(violations, claimRecord(sandbox, name, records, table.Name, "the table "+table.Name)...)

		fields := map[string]bool{}
		for _, field := range table.Fields {
			if fields[field.Name] {
				violations = append(violations, databaseViolation(name, "declares the field "+table.Name+"."+field.Name+" twice"))
			}
			fields[field.Name] = true

			if field.Name == databaseIdField {
				violations = append(violations, databaseViolation(name,
					"declares the field "+table.Name+".id; every record already carries its permanent Id"))
			}

			violations = append(violations, checkDatabaseField(sandbox, name, conf, records, table.Name, field)...)
		}
	}

	return violations
}

// checkDatabaseField enforces what one field may declare: a link names a table
// of this database, a nested collection carries fields of its own and nests no
// further, and no other kind of field claims a target.
func checkDatabaseField(sandbox *api.Sandbox, name string, conf *databaseconf.DatabaseConf,
	records map[string]string, table string, field databaseconf.Field) []string {

	var violations []string
	where := table + "." + field.Name

	switch field.Type {
	case databaseconf.FieldLink:
		if field.Target == "" {
			violations = append(violations, databaseViolation(name,
				"declares the link "+where+" with no `target`; a link names the table it points at"))
			break
		}
		if utils.FindDatabaseTable(sandbox, conf.Tables, field.Target) < 0 {
			violations = append(violations, databaseViolation(name,
				"declares the link "+where+" targeting "+field.Target+", which is not a table of this database"))
		}
	case databaseconf.FieldDatabase:
		if len(field.Fields) == 0 {
			violations = append(violations, databaseViolation(name,
				"declares the nested collection "+where+" with no `fields`; it would hold records of nothing"))
		}
		violations = append(violations, claimRecord(sandbox, name, records, field.Name, "the nested collection "+where)...)

		nested_names := map[string]bool{}
		for _, nested := range field.Fields {
			if nested.Type == databaseconf.FieldDatabase {
				violations = append(violations, databaseViolation(name,
					"nests a collection inside "+where+"; only one level is generated"))
			}
			if nested.Name == databaseIdField {
				violations = append(violations, databaseViolation(name,
					"declares the field "+where+".id; every record already carries its permanent Id"))
			}
			if nested_names[nested.Name] {
				violations = append(violations, databaseViolation(name, "declares the field "+where+"."+nested.Name+" twice"))
			}
			nested_names[nested.Name] = true
		}
	default:
		if field.Target != "" {
			violations = append(violations, databaseViolation(name,
				"declares a `target` on the "+field.Type+" field "+where+"; only a link points at a table"))
		}
	}

	return violations
}

// claimRecord reserves the Go record name one table or one nested collection
// generates, and reports the second claimant of a name.
func claimRecord(sandbox *api.Sandbox, name string, records map[string]string, declared string, what string) []string {
	kind := utils.ExportedName(sandbox, declared)
	if taken, seen := records[kind]; seen {
		return []string{databaseViolation(name,
			what+" and "+taken+" both generate "+kind+"Item; two of them cannot share one record type")}
	}
	records[kind] = what
	return nil
}

// checkDatabaseCustom enforces that the hand-written half of a database
// package redeclares nothing the build writes. The two files are one Go
// package, so a collision is a compile error — this reports it against the
// declaration, before the compiler reports it against a generated file nobody
// wrote.
func checkDatabaseCustom(sandbox *api.Sandbox, io *smartio.SmartIO, name string, conf *databaseconf.DatabaseConf) []string {
	custom := utils.DatabasesDir + "/" + name + "/" + utils.DatabaseCustomFile
	if !io.IsFile(custom) {
		return nil
	}

	content, err := io.ReadFile(custom)
	if err != nil {
		return nil
	}

	parsed, err := sandbox.Deps.Goimportsdeps.Parse(string(content))
	if err != nil {
		return []string{databaseViolation(name, utils.DatabaseCustomFile+" is not parsable Go: "+err.Error())}
	}

	generated := generatedDatabaseNames(sandbox, conf)

	var violations []string
	for _, function := range parsed.Functions {
		if function.Receiver != "" {
			continue
		}
		if generated[function.Name] {
			violations = append(violations, databaseViolation(name,
				utils.DatabaseCustomFile+" declares "+function.Name+", which the build generates in methods.go"))
		}
	}
	for _, declared := range parsed.Types {
		if generated[declared.Name] {
			violations = append(violations, databaseViolation(name,
				utils.DatabaseCustomFile+" declares the type "+declared.Name+", which the build generates in api.go"))
		}
	}

	return violations
}

// generatedDatabaseNames is every top-level name a build writes into a
// database package: the methods, the records, and the two helpers each record
// carries.
func generatedDatabaseNames(sandbox *api.Sandbox, conf *databaseconf.DatabaseConf) map[string]bool {
	names := map[string]bool{"New": true}

	claim := func(record string, filtrage bool) {
		names[record+"Item"] = true
		names[record+"New"] = true
		names["new"+record+"Fields"] = true
		names["build"+record+"Item"] = true
		if filtrage {
			names[record+"Filtrage"] = true
			names["match"+record] = true
		}
	}

	for _, table := range conf.Tables {
		kind := utils.ExportedName(sandbox, table.Name)
		claim(kind, true)
		for _, method := range utils.DatabaseMethodNames(sandbox, table) {
			names[method] = true
		}
		for _, field := range table.Fields {
			if field.Type == databaseconf.FieldDatabase {
				claim(utils.ExportedName(sandbox, field.Name), false)
			}
		}
	}

	names[utils.ExportedName(sandbox, conf.Name)] = true
	return names
}

// checkDatabaseRawFields reads the declaration a second time, unparsed, for
// the rule the parser resolves away: it drops `required` from a nested
// collection, which is never written by an insert, so by the time a
// DatabaseConf exists the contradiction is gone. The file is what has to be
// right.
func checkDatabaseRawFields(sandbox *api.Sandbox, name string, content string) []string {
	specs, err := sandbox.Deps.Serializables.ParseYaml(content)
	if err != nil || !specs.IsObject() {
		return nil
	}

	tables, _ := specs.GetObjectItem("tables")
	if tables == nil || !tables.IsArray() {
		return nil
	}
	size, err := tables.GetArraySize()
	if err != nil {
		return nil
	}

	var violations []string
	for i := 0; i < size; i++ {
		entry := tables.GetArrayItem(i)
		if entry == nil || !entry.IsObject() {
			continue
		}
		table := rawString(entry, "name")

		fields, _ := entry.GetObjectItem("fields")
		if fields == nil || !fields.IsArray() {
			continue
		}
		count, err := fields.GetArraySize()
		if err != nil {
			continue
		}
		for at := 0; at < count; at++ {
			field := fields.GetArrayItem(at)
			if field == nil || !field.IsObject() {
				continue
			}
			violations = append(violations, checkRawDatabaseField(sandbox, name, table, field)...)
		}
	}

	return violations
}

// checkRawDatabaseField reports the contradiction one unparsed field entry may
// carry.
func checkRawDatabaseField(sandbox *api.Sandbox, name string, table string, entry *serializables.SerializibleObject) []string {
	if rawString(entry, "type") != databaseconf.FieldDatabase || !rawBool(entry, "required") {
		return nil
	}
	return []string{databaseViolation(name,
		"declares the nested collection "+table+"."+rawString(entry, "name")+
			" as required; a nested collection is never written by an insert")}
}

// databaseNameOf is the last segment of a listed database directory.
func databaseNameOf(sandbox *api.Sandbox, path string) string {
	parts := sandbox.Deps.Stringsdeps.Split(path, "/")
	return parts[len(parts)-1]
}

// databaseViolation words one violation the same way for every rule.
func databaseViolation(name string, reason string) string {
	return utils.DatabasesDir + "/" + name + " " + reason
}
