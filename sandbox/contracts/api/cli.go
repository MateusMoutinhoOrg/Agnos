package api

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"

type FlagValue interface {
	String() string
	Int() int
	Float() float64
	Bool() bool
}

const (
	FlagTypeString = iota
	FlagTypeInt
	FlagTypeFloat
	FlagTypeBool
)

type CliFlag struct {
	ValidIdentifiers []string
	Values           []FlagValue
	Type             int
	Required         bool
}

type CliCommand struct {
	ValidStartIdentifiers []string
	Description           string
	Examples              []string
	Flags                 []CliFlag
	Handler               func(deps deps.Deps) int
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
