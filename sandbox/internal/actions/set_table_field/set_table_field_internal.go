package set_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetTableFieldInternal parses the target database's specs.yaml, rewrites the
// named field in place and writes the file back. It is add-table-field applied
// to a declaration that already exists: the keys given are written over the
// ones there, --clear takes one off, and the result goes through the same
// constructor, so an edited field and a declared one are the same bytes.
func SetTableFieldInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.DatabaseFieldEditProps) error {
	conf, err := utils.LoadDatabaseConf(sandbox, io, props.Database)
	if err != nil {
		return err
	}

	name := utils.DatabaseName(sandbox, props.Name)
	if name == "" {
		return sandbox.Deps.Std.Errorf("set-table-field needs the name of the field to edit")
	}
	if utils.DatabaseFieldEditEmpty(sandbox, props) {
		return sandbox.Deps.Std.Errorf("set-table-field: nothing to change (pass --rename, --type, --required, --target or --clear)")
	}

	fields, where, err := utils.DatabaseFieldsAt(sandbox, conf, props.Table, props.Parent)
	if err != nil {
		return err
	}

	index := utils.FindDatabaseField(sandbox, fields, name)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("%s declares no field %q", where, name)
	}

	field, err := utils.DatabaseFieldEdited(sandbox, fields[index], conf, props)
	if err != nil {
		return err
	}
	if field.Name != name && utils.FindDatabaseField(sandbox, fields, field.Name) >= 0 {
		return sandbox.Deps.Std.Errorf("%s already declares a field named %q", where, field.Name)
	}

	sandbox.Deps.Std.Log("set-table-field updating %s.%s in %s \n", where, name, utils.DatabaseConfPath(sandbox, props.Database))

	fields[index] = field
	if err := utils.SetDatabaseFieldsAt(sandbox, conf, props.Table, props.Parent, fields); err != nil {
		return err
	}
	return utils.SaveDatabaseConf(sandbox, io, props.Database, conf)
}
