package cli



type CliValue interface {
	String() string
	Int() int
	Float() float64
	Bool() bool
}

const (
	CliTypeString = iota
	CliTypeInt
	CliTypeFloat
	CliTypeBool
)

type Cliflag struct {
	Id               string
	ValidIdentifiers []string

	Description string
	Examples    []string

	Values   []CliValue
	Defaults []string
	Exist    bool

	Type             int
	RequiredMinSize  int
	RequiredMaxSize  int
	RequiredPresence bool
}

type CliArg struct {
	Id          string
	Description string
	Examples    []string
	Values      []CliValue
	Defaults    []string

	RequiredType    int
	RequiredMinSize int
	RequiredMaxSize int
}

type CliEntrys struct {
	GetArgById  func(id string) *CliArg
	GetFlagById func(id string) *Cliflag
}

// CliCommand describes one top-level command the CLI dispatches to.
// ValidStartIdentifiers lists every verb the user may type to reach
// it; Description is the one-line summary shown in the command list;
// LongDescription is an optional multi-line explanation shown in
// per-command help; Category groups the command under a heading in the
// general help screen; Hidden, when true, omits the command from help
// output entirely.
type CliCommand struct {
	ValidStartIdentifiers []string
	Args                  []CliArg
	Flags                 []Cliflag
	Description           string
	LongDescription       string
	Category              string
	Hidden                bool
	Examples              []string
	Handler               func(sandbox any, entries CliEntrys) int
}

const (
	ExitOk = 0
	// ExitUsage reports that the command line itself was wrong — an unknown
	// command, a missing operand, an unparsable amount — and the usage
	// screen was printed.
	ExitUsage = 1
	// ExitFailure reports that a well-formed command could not be carried
	// out, because a record was missing or could not be written.
	ExitFailure = 2
)

type CliApi struct {
	Commands []CliCommand
	CliMain  func(args []string) int
}
