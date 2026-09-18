package databaseconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
)

// New parses one specs.yaml body into a DatabaseConf.
func New(sandbox *api.Sandbox, content string) (*DatabaseConf, error) {

	if content == "" {
		return nil, sandbox.Deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := sandbox.Deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !specs.IsObject() {
		return nil, sandbox.Deps.Std.Errorf("specs.yaml is not an object")
	}

	conf := &DatabaseConf{
		Name:   readString(specs, "name"),
		Prefix: readString(specs, "prefix"),
		Tables: []Table{},
	}

	tables_item, _ := specs.GetObjectItem("tables")
	if tables_item != nil {
		tables, err := readTables(sandbox, tables_item)
		if err != nil {
			return nil, err
		}
		conf.Tables = tables
	}

	BindMethods(sandbox, conf)
	return conf, nil
}

// readTables parses the `tables` sequence. A table with no name declares a
// collection nothing can be spelled after, so it is refused here rather than
// generating a method called Add.
func readTables(sandbox *api.Sandbox, item *serializibles.SerializibleObject) ([]Table, error) {
	if item.IsNull() {
		return []Table{}, nil
	}
	if !item.IsArray() {
		return nil, sandbox.Deps.Std.Errorf("`tables` must be a sequence of tables")
	}

	size, err := item.GetArraySize()
	if err != nil {
		return nil, err
	}

	tables := make([]Table, 0, size)
	for i := 0; i < size; i++ {
		entry := item.GetArrayItem(i)
		if entry == nil || !entry.IsObject() {
			return nil, sandbox.Deps.Std.Errorf("`tables` entry #%d is not an object", i)
		}

		table := Table{Name: readString(entry, "name"), Fields: []Field{}}
		if table.Name == "" {
			return nil, sandbox.Deps.Std.Errorf("`tables` entry #%d needs a name", i)
		}

		fields_item, _ := entry.GetObjectItem("fields")
		if fields_item != nil {
			fields, err := readFields(sandbox, fields_item, table.Name)
			if err != nil {
				return nil, err
			}
			table.Fields = fields
		}

		tables = append(tables, table)
	}

	return tables, nil
}

// readFields parses the `fields` sequence of one table or of one nested
// database field. The `where` label is the path the error message names, so a
// bad entry is reported at the table or the nested collection it sits in.
func readFields(sandbox *api.Sandbox, item *serializibles.SerializibleObject, where string) ([]Field, error) {
	if item.IsNull() {
		return []Field{}, nil
	}
	if !item.IsArray() {
		return nil, sandbox.Deps.Std.Errorf("`fields` of %s must be a sequence of fields", where)
	}

	size, err := item.GetArraySize()
	if err != nil {
		return nil, err
	}

	fields := make([]Field, 0, size)
	for i := 0; i < size; i++ {
		entry := item.GetArrayItem(i)
		if entry == nil || !entry.IsObject() {
			return nil, sandbox.Deps.Std.Errorf("`fields` entry #%d of %s is not an object", i, where)
		}

		field := Field{
			Name:     readString(entry, "name"),
			Type:     readString(entry, "type"),
			Required: readBool(entry, "required"),
			Target:   readString(entry, "target"),
			Fields:   []Field{},
		}
		if field.Name == "" {
			return nil, sandbox.Deps.Std.Errorf("`fields` entry #%d of %s needs a name", i, where)
		}
		if field.Type == "" {
			field.Type = FieldString
		}
		if !IsFieldType(field.Type) {
			return nil, sandbox.Deps.Std.Errorf("field %s of %s declares the unknown type %q", field.Name, where, field.Type)
		}

		// A nested collection is never written by an insert, so `required`
		// on one says nothing the generator could honour. It is dropped
		// here, and `verify` reports the file that still carries it.
		if field.Type == FieldDatabase {
			field.Required = false
			nested_item, _ := entry.GetObjectItem("fields")
			if nested_item != nil {
				nested, err := readFields(sandbox, nested_item, where+"."+field.Name)
				if err != nil {
					return nil, err
				}
				field.Fields = nested
			}
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// IsFieldType reports whether kind is one of the declared field types.
func IsFieldType(kind string) bool {
	for _, known := range FieldTypes {
		if known == kind {
			return true
		}
	}
	return false
}

func readString(item *serializibles.SerializibleObject, key string) string {
	value, _ := item.GetObjectItem(key)
	if value == nil || !value.IsString() {
		return ""
	}
	text, err := value.GetString()
	if err != nil {
		return ""
	}
	return text
}

func readBool(item *serializibles.SerializibleObject, key string) bool {
	value, _ := item.GetObjectItem(key)
	if value == nil || !value.IsBool() {
		return false
	}
	truth, err := value.GetBool()
	if err != nil {
		return false
	}
	return truth
}
