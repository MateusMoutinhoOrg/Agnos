package api

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"

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

	Values []CliValue
	Exist  bool

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

	RequiredType int
	RequiredSize int
}

type CliEntrys struct {
	GetValue func(id string) CliValue
	Exist    func(id string) bool
	Size     func(id string) int
}

type CliCommand struct {
	ValidStartIdentifiers []string
	Args                  []CliArg
	Flags                 []Cliflag
	Description           string
	Examples              []string
	Handler               func(deps deps.Deps, entries CliEntrys) int
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
