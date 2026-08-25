package api

type MemoryCliFlag struct {
	Id               string
	ValidIdentifiers []string

	Description string
	Examples    []string

	Type             int
	RequiredMinSize  int
	RequiredMaxSize  int
	RequiredPresence bool
}

type MemoryCliArg struct {
	Id           string
	Description  string
	Examples     []string
	RequiredType int
	RequiredSize int
}

type MemoryCliCommand struct {
	ValidStartIdentifiers []string
	Args                  []MemoryCliArg
	Flags                 []MemoryCliFlag
	Description           string
	Examples              []string
}

type Memory struct {
	ProjectName string
	Version     func() string
	Commands    []MemoryCliCommand

	Persist func()
}

type MemoryApi struct {
	NewMemory func(path string) Memory
}
