package add_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/migrationconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddTableFieldInternal parses the target database's database.yaml, appends the
// new field to the table — or to the nested collection --parent names — and
// writes the file back. Where the field lands is what decides which methods it
// generates: a `key` brings a Find, a `link` a Get, a `database` the Add/List
// pair, and every plain field an Update and a place in the filter.
//
// An optional field reads as its zero value on every record stored before it,
// so the data needs nothing. A required one opens a version: the records
// already stored hold no value for it, so the migration fills them with
// --default, or the hand-written step --fill-by-hand writes does, and a record
// left without one fails the migration rather than the data.
func AddTableFieldInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.AddTableFieldProps) error {
	conf, err := utils.LoadDatabaseConf(sandbox, io, props.Database)
	if err != nil {
		return err
	}
	before := conf.Render()

	fields, where, err := utils.DatabaseFieldsAt(sandbox, conf, props.Table, props.Parent)
	if err != nil {
		return err
	}

	field, err := utils.NewDatabaseField(sandbox, conf, props)
	if err != nil {
		return err
	}

	if utils.FindDatabaseField(sandbox, fields, field.Name) >= 0 {
		return sandbox.Deps.StdDeps.Errorf("%s already declares a field named %q", where, field.Name)
	}

	if err := checkFill(sandbox, props); err != nil {
		return err
	}
	if props.Default != "" {
		if err := utils.CheckDatabaseDefault(sandbox, field, props.Default); err != nil {
			return err
		}
	}

	sandbox.Deps.StdDeps.Logf("add-table-field adding %s.%s to %s \n", where, field.Name, utils.DatabaseConfPath(sandbox, props.Database))

	if err := utils.SetDatabaseFieldsAt(sandbox, conf, props.Table, props.Parent, append(fields, field)); err != nil {
		return err
	}

	if field.Required {
		op := utils.DatabaseRequireOp(utils.DatabaseName(sandbox, props.Table), utils.DatabaseName(sandbox, props.Parent), field, props.Default)
		if err := utils.OpenDatabaseVersion(sandbox, io, props.Database, before, conf, []migrationconf.Op{op}); err != nil {
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

// checkFill refuses a --default or a --fill-by-hand with nothing to fill:
// both are how the records already stored get a value for a required field.
func checkFill(sandbox *api.Sandbox, props api.AddTableFieldProps) error {
	if props.Default != "" && props.FillByHand {
		return sandbox.Deps.StdDeps.Errorf("--default and --fill-by-hand both fill the records already stored: pass one of them")
	}
	if (props.Default != "" || props.FillByHand) && !props.Required {
		return sandbox.Deps.StdDeps.Errorf("--default and --fill-by-hand fill the records already stored with a required field: pass --required, or leave them out — an optional field reads as empty on a record stored before it")
	}
	return nil
}
