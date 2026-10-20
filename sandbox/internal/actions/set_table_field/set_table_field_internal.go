package set_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/migrationconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetTableFieldInternal parses the target database's database.yaml, rewrites the
// named field in place and writes the file back. It is add-table-field applied
// to a declaration that already exists: the keys given are written over the
// ones there, --clear takes one off, and the result goes through the same
// constructor, so an edited field and a declared one are the same bytes.
//
// An edit that changes what the stored data holds — a rename, a type some
// stored value may not read as, a field made required — opens a version whose
// operations carry the data across: the values move to the new name, convert
// to the new type, and fill a required field from --default or from the
// hand-written step --fill-by-hand writes. Every other edit leaves the version
// as it is.
func SetTableFieldInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.SetTableFieldProps) error {
	conf, err := utils.LoadDatabaseConf(sandbox, io, props.Database)
	if err != nil {
		return err
	}
	before := conf.Render()

	name := utils.DatabaseName(sandbox, props.Name)
	if name == "" {
		return sandbox.Deps.StdDeps.Errorf("set-table-field needs the name of the field to edit")
	}
	if utils.IsDatabaseFieldEditEmpty(sandbox, props) {
		return sandbox.Deps.StdDeps.Errorf("set-table-field: nothing to change (pass --rename, --type, --required, --target or --clear)")
	}
	if props.Default != "" && props.FillByHand {
		return sandbox.Deps.StdDeps.Errorf("--default and --fill-by-hand both fill the records already stored: pass one of them")
	}

	fields, where, err := utils.DatabaseFieldsAt(sandbox, conf, props.Table, props.Parent)
	if err != nil {
		return err
	}

	index := utils.FindDatabaseField(sandbox, fields, name)
	if index < 0 {
		return sandbox.Deps.StdDeps.Errorf("%s declares no field %q", where, name)
	}
	current := fields[index]

	field, err := utils.EditDatabaseField(sandbox, current, conf, props)
	if err != nil {
		return err
	}
	if field.Name != name && utils.FindDatabaseField(sandbox, fields, field.Name) >= 0 {
		return sandbox.Deps.StdDeps.Errorf("%s already declares a field named %q", where, field.Name)
	}

	ops := utils.DatabaseFieldEditOps(utils.DatabaseName(sandbox, props.Table), utils.DatabaseName(sandbox, props.Parent), current, field, props.Default)
	if err := checkFill(sandbox, props, field, ops); err != nil {
		return err
	}

	sandbox.Deps.StdDeps.Logf("set-table-field updating %s.%s in %s \n", where, name, utils.DatabaseConfPath(sandbox, props.Database))

	fields[index] = field
	if err := utils.SetDatabaseFieldsAt(sandbox, conf, props.Table, props.Parent, fields); err != nil {
		return err
	}

	if len(ops) > 0 {
		if err := utils.OpenDatabaseVersion(sandbox, io, props.Database, before, conf, ops); err != nil {
			return err
		}
		if props.FillByHand {
			if err := utils.WriteDatabaseStep(sandbox, io, props.Database, conf.Version); err != nil {
				return err
			}
		}
	}

	return utils.SaveDatabaseConf(sandbox, io, props.Database, conf)
}

// checkFill refuses a --default or a --fill-by-hand this edit has nothing to
// fill with: both belong to a field the edit makes required.
func checkFill(sandbox *api.Sandbox, props api.SetTableFieldProps, field databaseconf.Field, ops []migrationconf.Op) error {
	if props.Default == "" && !props.FillByHand {
		return nil
	}
	for _, op := range ops {
		if op.Kind == migrationconf.OpRequireField {
			if props.Default == "" {
				return nil
			}
			return utils.CheckDatabaseDefault(sandbox, field, props.Default)
		}
	}
	return sandbox.Deps.StdDeps.Errorf("--default and --fill-by-hand fill the records already stored with a field this edit makes required, and %s is not made required here", field.Name)
}
