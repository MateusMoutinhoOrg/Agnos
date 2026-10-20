package migrationconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializabledeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializabledeps"
)

// New parses one v<N>.ops.yaml body into a MigrationConf.
func New(sandbox *api.Sandbox, content string) (*MigrationConf, error) {

	if content == "" {
		return nil, sandbox.Deps.StdDeps.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := sandbox.Deps.SerializableDeps.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}

	if !specs.IsObject() {
		return nil, sandbox.Deps.StdDeps.Errorf("the ops file is not an object")
	}

	conf := &MigrationConf{Ops: []Op{}}

	ops_item, _ := specs.GetObjectItem("ops")
	if ops_item != nil && !ops_item.IsNull() {
		if !ops_item.IsArray() {
			return nil, sandbox.Deps.StdDeps.Errorf("`ops` must be a sequence of operations")
		}
		size, err := ops_item.GetArraySize()
		if err != nil {
			return nil, err
		}
		for i := 0; i < size; i++ {
			op, err := readOp(sandbox, ops_item.GetArrayItem(i), i)
			if err != nil {
				return nil, err
			}
			conf.Ops = append(conf.Ops, op)
		}
	}

	BindMethods(sandbox, conf)
	return conf, nil
}

// readOp parses one entry of `ops`. An entry with no kind, or one the engine
// does not run, is refused here rather than skipped by the migration.
func readOp(sandbox *api.Sandbox, entry *serializabledeps.SerializableObject, index int) (Op, error) {
	if entry == nil || !entry.IsObject() {
		return Op{}, sandbox.Deps.StdDeps.Errorf("`ops` entry #%d is not an object", index)
	}

	op := Op{
		Kind:    readString(entry, "kind"),
		Table:   readString(entry, "table"),
		Parent:  readString(entry, "parent"),
		Field:   readString(entry, "field"),
		Type:    readString(entry, "type"),
		To:      readString(entry, "to"),
		Default: readString(entry, "default"),
	}
	if !IsOpKind(op.Kind) {
		return Op{}, sandbox.Deps.StdDeps.Errorf("`ops` entry #%d declares the unknown kind %q", index, op.Kind)
	}
	if op.Table == "" {
		return Op{}, sandbox.Deps.StdDeps.Errorf("`ops` entry #%d (%s) names no table", index, op.Kind)
	}
	return op, nil
}

// IsOpKind reports whether kind is one of the declared op kinds.
func IsOpKind(kind string) bool {
	for _, known := range OpKinds {
		if known == kind {
			return true
		}
	}
	return false
}

func readString(item *serializabledeps.SerializableObject, key string) string {
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
