package opinionatedagnosdatabase

import (
	"fmt"
	"strconv"
	"strings"

	opinionatedagnosdatabase "{{.Module}}/sandbox/deps/OpinionatedAgnosDatabase"
	"{{.Module}}/sandbox/deps/databasedeps"
)

// The operations a version runs on a migration's copy. Each one opens the
// copy through a lens — the one table it reaches, holding only the fields it
// reads and writes, every one optional so it can be cleared — and goes
// through the store, so an index follows every value it moves. A whole
// collection, a table or one nested under a record, is a folder of its own,
// and is moved or removed as one.

// applyOp runs one operation on the copy.
func applyOp(props opinionatedagnosdatabase.MigrateProps, at folders, op opinionatedagnosdatabase.Op) error {
	switch op.Kind {
	case opinionatedagnosdatabase.OpRemoveTable:
		props.IoDeps.RemoveDir(props.IoDeps.Join(at.work, op.Table))
		return nil
	case opinionatedagnosdatabase.OpRenameField:
		if op.Type == opinionatedagnosdatabase.TypeObject {
			return eachNestedFolder(props, at, op, func(folder string, owner string) error {
				return props.IoDeps.Rename(folder, props.IoDeps.Join(owner, op.To))
			})
		}
		return renameValues(props.Databases, at.workPath, op)
	case opinionatedagnosdatabase.OpRetypeField:
		return retypeValues(props.Databases, at.workPath, op)
	case opinionatedagnosdatabase.OpRequireField:
		return fillValues(props.Databases, at.workPath, op)
	case opinionatedagnosdatabase.OpRemoveField:
		if op.Type == opinionatedagnosdatabase.TypeObject {
			return eachNestedFolder(props, at, op, func(folder string, owner string) error {
				props.IoDeps.RemoveDir(folder)
				return nil
			})
		}
		return clearValues(props.Databases, at.workPath, op)
	}
	return fmt.Errorf("unknown operation %q", op.Kind)
}

// renameValues moves every stored value of op.Field to op.To.
func renameValues(databases databasedeps.Databases, path string, op opinionatedagnosdatabase.Op) error {
	kind := lensType(op.Type)
	database, err := lens(databases, path, op.Table, op.Parent,
		databasedeps.Field{Name: op.Field, Type: kind},
		databasedeps.Field{Name: op.To, Type: kind})
	if err != nil {
		return err
	}
	return walk(database, op.Table, op.Parent, func(owner databasedeps.Record, record databasedeps.Record) error {
		value, found, err := valueOf(record, op.Field)
		if err != nil || !found {
			return err
		}
		if failure := record.Update(op.To, value); failure != nil {
			return recordFailure(record, failure)
		}
		return recordFailure(record, record.Update(op.Field, nil))
	})
}

// retypeValues converts every stored value of op.Field from op.Type to
// op.To: read through a lens of the old type, cleared there — which drops
// its index entry when it was a key — and written through a lens of the new
// one, which indexes it when it now is.
func retypeValues(databases databasedeps.Databases, path string, op opinionatedagnosdatabase.Op) error {
	from, err := lens(databases, path, op.Table, op.Parent, databasedeps.Field{Name: op.Field, Type: lensType(op.Type)})
	if err != nil {
		return err
	}
	to, err := lens(databases, path, op.Table, op.Parent, databasedeps.Field{Name: op.Field, Type: lensType(op.To)})
	if err != nil {
		return err
	}
	return walk(from, op.Table, op.Parent, func(owner databasedeps.Record, record databasedeps.Record) error {
		value, found, err := valueOf(record, op.Field)
		if err != nil || !found {
			return err
		}
		converted, err := convert(value, op.To)
		if err != nil {
			return fmt.Errorf("record %d: %w", record.ID, err)
		}
		if failure := record.Update(op.Field, nil); failure != nil {
			return recordFailure(record, failure)
		}
		target, err := recordAt(to, op.Table, op.Parent, owner.ID, record.ID)
		if err != nil {
			return err
		}
		return recordFailure(record, target.Update(op.Field, converted))
	})
}

// fillValues gives op.Default to every record holding no value for
// op.Field. With no default it fills nothing: a record still holding none
// is what the check after the last step reports.
func fillValues(databases databasedeps.Databases, path string, op opinionatedagnosdatabase.Op) error {
	if op.Default == "" {
		return nil
	}
	value, err := convert(op.Default, op.Type)
	if err != nil {
		return fmt.Errorf("the default: %w", err)
	}
	database, err := lens(databases, path, op.Table, op.Parent, databasedeps.Field{Name: op.Field, Type: lensType(op.Type)})
	if err != nil {
		return err
	}
	return walk(database, op.Table, op.Parent, func(owner databasedeps.Record, record databasedeps.Record) error {
		has, failure := record.HasValues([]string{op.Field})
		if failure != nil {
			return recordFailure(record, failure)
		}
		if has {
			return nil
		}
		return recordFailure(record, record.Update(op.Field, value))
	})
}

// clearValues deletes every stored value of op.Field, and its index entry.
func clearValues(databases databasedeps.Databases, path string, op opinionatedagnosdatabase.Op) error {
	database, err := lens(databases, path, op.Table, op.Parent, databasedeps.Field{Name: op.Field, Type: lensType(op.Type)})
	if err != nil {
		return err
	}
	return walk(database, op.Table, op.Parent, func(owner databasedeps.Record, record databasedeps.Record) error {
		return recordFailure(record, record.Update(op.Field, nil))
	})
}

// eachNestedFolder hands visit the folder of the nested collection op.Field
// under every live record of op.Table, and the record's own folder, skipping
// a record that never had one.
func eachNestedFolder(props opinionatedagnosdatabase.MigrateProps, at folders, op opinionatedagnosdatabase.Op, visit func(folder string, owner string) error) error {
	database, err := lens(props.Databases, at.workPath, op.Table, "")
	if err != nil {
		return err
	}
	return walk(database, op.Table, "", func(_ databasedeps.Record, record databasedeps.Record) error {
		owner := props.IoDeps.Join(at.work, op.Table, strconv.FormatInt(record.ID, 10))
		folder := props.IoDeps.Join(owner, op.Field)
		if !props.IoDeps.IsDir(folder) {
			return nil
		}
		return visit(folder, owner)
	})
}

// checkData is the last check before the swap: every table of the last
// schema is repaired — which indexes a field that became a key, and fails on
// two records holding one value of it — and every record holds every field
// that schema requires.
func checkData(databases databasedeps.Databases, path string, last opinionatedagnosdatabase.Version) error {
	database, err := open(databases, path, last.Schemas)
	if err != nil {
		return err
	}
	for _, schema := range last.Schemas {
		table, err := collection(database, schema.Name)
		if err != nil {
			return err
		}
		if failure := table.Repair(); failure != nil {
			return fmt.Errorf("repairing %s: %w", schema.Name, fail(failure))
		}
		if err := checkRequired(table, schema.Name, schema.Fields, last.Version); err != nil {
			return err
		}
	}
	return nil
}

// checkRequired reports the first record of a collection holding no value
// for a field the schema requires, and recurses into every nested one.
func checkRequired(table databasedeps.Collection, where string, fields []databasedeps.Field, version int) error {
	required := []string{}
	nested := []databasedeps.Field{}
	for _, field := range fields {
		if field.Type == databasedeps.Nested {
			nested = append(nested, field)
			continue
		}
		if field.Required {
			required = append(required, field.Name)
		}
	}

	records, failure := table.ListAll()
	if failure != nil {
		return fail(failure)
	}
	for _, record := range records {
		for _, name := range required {
			has, failure := record.HasValues([]string{name})
			if failure != nil {
				return recordFailure(record, failure)
			}
			if !has {
				return fmt.Errorf("record %d of %s holds no %s, which version %d requires: give it a value in a hand-written step of that version, or a default", record.ID, where, name, version)
			}
		}
		for _, field := range nested {
			children, failure := record.Nested(field.Name)
			if failure != nil {
				return recordFailure(record, failure)
			}
			if err := checkRequired(children, where+"."+field.Name, field.Fields, version); err != nil {
				return err
			}
		}
	}
	return nil
}

// lens opens the copy with one table holding the given fields alone — on the
// table, or inside parent — every one optional, so an operation can clear what
// it moves.
func lens(databases databasedeps.Databases, path string, table string, parent string, fields ...databasedeps.Field) (databasedeps.Database, error) {
	if parent != "" {
		nested := databasedeps.Field{Name: parent, Type: databasedeps.Nested, Fields: fields}
		fields = []databasedeps.Field{nested}
	}
	schema := databasedeps.Schema{Name: table, Fields: fields}
	return open(databases, path, []databasedeps.Schema{schema})
}

// lensType is the store type one declared type is read and written as in a
// lens. A link is read as the integer it is stored as: a lens holds no table
// for it to point at.
func lensType(kind string) int {
	switch kind {
	case opinionatedagnosdatabase.TypeKey:
		return databasedeps.Key
	case opinionatedagnosdatabase.TypeInteger, opinionatedagnosdatabase.TypeLink:
		return databasedeps.Int
	case opinionatedagnosdatabase.TypeNumber:
		return databasedeps.Float
	case opinionatedagnosdatabase.TypeBytes:
		return databasedeps.Bytes
	}
	return databasedeps.String
}

// walk hands visit every live record of table — or, with a parent, every
// record nested under each of them, beside the record owning it.
func walk(database databasedeps.Database, table string, parent string, visit func(owner databasedeps.Record, record databasedeps.Record) error) error {
	records, err := listAll(database, table)
	if err != nil {
		return err
	}
	for _, record := range records {
		if parent == "" {
			if err := visit(databasedeps.Record{}, record); err != nil {
				return err
			}
			continue
		}
		children, failure := record.ListNested(parent)
		if failure != nil {
			return recordFailure(record, failure)
		}
		for _, child := range children {
			if err := visit(record, child); err != nil {
				return err
			}
		}
	}
	return nil
}

// listAll is every live record of one table.
func listAll(database databasedeps.Database, table string) ([]databasedeps.Record, error) {
	resolved, err := collection(database, table)
	if err != nil {
		return nil, err
	}
	records, failure := resolved.ListAll()
	if failure != nil {
		return nil, fail(failure)
	}
	return records, nil
}

// recordAt is the record of another lens that one walked record is: the
// same id, under the same owner.
func recordAt(database databasedeps.Database, table string, parent string, owner int64, id int64) (databasedeps.Record, error) {
	resolved, err := collection(database, table)
	if err != nil {
		return databasedeps.Record{}, err
	}
	if parent != "" {
		record, found, failure := resolved.FindByID(owner)
		if failure != nil || !found {
			return databasedeps.Record{}, fmt.Errorf("record %d of %s is gone", owner, table)
		}
		resolved, failure = record.Nested(parent)
		if failure != nil {
			return databasedeps.Record{}, fail(failure)
		}
	}
	record, found, failure := resolved.FindByID(id)
	if failure != nil || !found {
		return databasedeps.Record{}, fmt.Errorf("record %d of %s is gone", id, table)
	}
	return record, nil
}

// valueOf is the stored value of one field, and false when there is none.
func valueOf(record databasedeps.Record, field string) (any, bool, error) {
	raw, failure := record.Get(field)
	if failure == nil {
		return raw, true, nil
	}
	if failure.Type == databasedeps.NoValue {
		return nil, false, nil
	}
	return nil, false, recordFailure(record, failure)
}

// convert turns one stored value — or a default, given as text — into the
// value a field of type kind stores. A value that does not read as that
// type is an error, never a zero.
func convert(value any, kind string) (any, error) {
	text := ""
	switch typed := value.(type) {
	case string:
		text = typed
	case []byte:
		text = string(typed)
	case int64:
		text = strconv.FormatInt(typed, 10)
	case float64:
		text = strconv.FormatFloat(typed, 'g', -1, 64)
	default:
		return nil, fmt.Errorf("a %T cannot be converted", value)
	}

	switch kind {
	case opinionatedagnosdatabase.TypeBytes:
		return []byte(text), nil
	case opinionatedagnosdatabase.TypeInteger, opinionatedagnosdatabase.TypeLink:
		if whole, ok := value.(int64); ok {
			return whole, nil
		}
		if number, ok := value.(float64); ok {
			if number != float64(int64(number)) {
				return nil, fmt.Errorf("%s is not a whole number", text)
			}
			return int64(number), nil
		}
		whole, err := strconv.ParseInt(strings.TrimSpace(text), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%q is not a whole number", text)
		}
		return whole, nil
	case opinionatedagnosdatabase.TypeNumber:
		if whole, ok := value.(int64); ok {
			return float64(whole), nil
		}
		if number, ok := value.(float64); ok {
			return number, nil
		}
		number, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
		if err != nil {
			return nil, fmt.Errorf("%q is not a number", text)
		}
		return number, nil
	}
	return text, nil
}

// recordFailure is a store failure, named after the record it happened on.
func recordFailure(record databasedeps.Record, failure *databasedeps.Error) error {
	if failure == nil {
		return nil
	}
	return fmt.Errorf("record %d: %w", record.ID, fail(failure))
}

// opPlace is where an operation reaches, for the message naming it.
func opPlace(op opinionatedagnosdatabase.Op) string {
	place := op.Table
	if op.Parent != "" {
		place += "." + op.Parent
	}
	if op.Field != "" {
		place += "." + op.Field
	}
	return place
}
