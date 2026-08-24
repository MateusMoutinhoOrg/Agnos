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

type Cliflag struct {
	Name        string
	Description string
	Examples    []string

	ValidIdentifiers []string
	Values           []FlagValue
	Type             int
	MinSize          int
	MaxSize          int
	Required         bool
}
type CliArg struct {
	Name        string
	Description string
	Required    bool
}

// it will never have errors here, since the validation steps will ensure
// the flags its correct.
type FlagsRetriver struct {
	GetIntFlag    func(name string, index int) int
	GetStringFlag func(name string, index int) string
	GetFloatFlag  func(name string, index int) float64
	GetBoolFlag   func(name string, index int) bool
}
type ArgsRetriver struct {
	GetIntArg    func(name string, index int) int
	GetStringArg func(name string, index int) string
	GetFloatArg  func(name string, index int) float64
	GetBoolArg   func(name string, index int) bool
}

type CliCommand struct {
	ValidStartIdentifiers []string
	ArgsList              []CliArg
	FlagsList             []Cliflag
	Description           string
	Examples              []string
	Handler               func(deps deps.Deps, argsRetriver *ArgsRetriver, flagsRetriver *FlagsRetriver) int
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
