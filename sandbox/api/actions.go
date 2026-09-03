package api

type StartProps struct {
	Path        string
	ProjectName string
	Module      *string
	Force       bool
}

// FieldProps describes one flag or positional arg to add to a command's
// entries.yaml. Default, Min and Max are the raw literals typed on the
// command line ("" means unset) so the action can tell "not given" from a
// zero value; Position is the index to insert at (< 0 appends).
type FieldProps struct {
	Path        string
	Command     string
	Name        string
	Identifiers []string
	Description string
	Examples    []string
	Type        string
	Default     string
	Required    bool
	Array       bool
	Min         string
	Max         string
	Position    int
}

// CommandProps carries the command-level keys of entries.yaml that
// set-command may rewrite. Empty strings leave the current value alone;
// Identifiers / Examples are appended (deduplicated).
type CommandProps struct {
	Path            string
	Command         string
	Help            string
	Category        string
	LongDescription string
	Hidden          bool
	Visible         bool
	Identifiers     []string
	Examples        []string
}

type Actions struct {
	Build         func(path string) error
	Verify        func(path string) error
	Start         func(props StartProps) error
	DepsInit      func(path string) error
	DepsPurge     func(path string) error
	DepInstall    func(path string, dep string) error
	DepRemove     func(path string, dep string) error
	DepList       func(path string) ([]string, error)
	CliInit       func(path string) error
	CliPurge      func(path string) error
	AddCommand    func(path string, name string, help string, category string) error
	RemoveCommand func(path string, name string) error
	SetCommand    func(props CommandProps) error
	AddFlag       func(props FieldProps) error
	RemoveFlag    func(path string, command string, name string) error
	AddArg        func(props FieldProps) error
	RemoveArg     func(path string, command string, name string) error
}
