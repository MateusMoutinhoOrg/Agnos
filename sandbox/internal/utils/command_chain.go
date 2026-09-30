package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CommandChainEntry is one declared command of the chain: its package name,
// the directory it sits in and its parsed declaration.
type CommandChainEntry struct {
	Name string
	Dir  string
	Conf *commandconf.CommandConf
}

// LoadCommandChain reads every command.yaml under sandbox/internal/commands
// and returns them in the order the dispatch runs them: by priority, lowest
// first, then by name — the order the build collector lays Cli.Commands down
// in.
func LoadCommandChain(sandbox *api.Sandbox, io *smartio.SmartIO) ([]CommandChainEntry, error) {
	chain := []CommandChainEntry{}

	for _, unit := range CommandDirs(sandbox, io) {
		content, err := io.ReadFile(unit.Dir + "/" + CommandConfFile)
		if err != nil {
			continue
		}
		conf, err := commandconf.New(sandbox, string(content))
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("%s/%s: %w", unit.Dir, CommandConfFile, err)
		}
		chain = append(chain, CommandChainEntry{Name: unit.Name, Dir: unit.Dir, Conf: conf})
	}

	SortCommandChain(sandbox, chain)
	return chain, nil
}

// SortCommandChain puts a chain in run order, in place.
func SortCommandChain(sandbox *api.Sandbox, chain []CommandChainEntry) {
	sandbox.Deps.Sortdeps.SliceStable(chain, func(i int, j int) bool {
		left, right := chain[i].Conf, chain[j].Conf
		if left.Priority != right.Priority {
			return left.Priority < right.Priority
		}
		return chain[i].Name < chain[j].Name
	})
}

// CommandRelativePriority is the rung one rung below (before) or above (after)
// the command named, for --before and --after. At most one of the two is
// given; it reports false when neither is.
func CommandRelativePriority(sandbox *api.Sandbox, io *smartio.SmartIO, before string, after string) (int, bool, error) {
	before = sandbox.Deps.Stringsdeps.TrimSpace(before)
	after = sandbox.Deps.Stringsdeps.TrimSpace(after)

	if before != "" && after != "" {
		return 0, false, sandbox.Deps.Std.Errorf("--before and --after exclude each other")
	}
	if before == "" && after == "" {
		return 0, false, nil
	}

	if before != "" {
		other, err := LoadCommandConf(sandbox, io, before)
		if err != nil {
			return 0, false, err
		}
		if other.Priority == 0 {
			return 0, false, sandbox.Deps.Std.Errorf(
				"--before %s: it runs on rung 0, and nothing runs below it — spread the chain with rebalance-commands first", before)
		}
		return other.Priority - 1, true, nil
	}

	other, err := LoadCommandConf(sandbox, io, after)
	if err != nil {
		return 0, false, err
	}
	return other.Priority + 1, true, nil
}
