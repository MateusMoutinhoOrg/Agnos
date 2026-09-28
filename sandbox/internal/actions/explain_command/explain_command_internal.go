package explain_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// ExplainCommandInternal walks the chain the way the generated dispatch does
// and says, command by command, whether it runs for the command line and why
// not when it does not. What a handler does once it runs is its own code, so a
// middleware that runs is read as handing the line on, and the first strict
// command that runs as the one that answers it. The last line is what the line
// ends on — that command, the usage error it raises before its handler runs,
// or the not-found every unmatched line ends on.
func ExplainCommandInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.ExplainCommandProps) ([]string, error) {
	if err := utils.RequireProject(sandbox, io); err != nil {
		return nil, err
	}
	if !io.IsDir(utils.CommandsDir) {
		return nil, sandbox.Deps.Std.Errorf("the project has no cli layer: run cli-init first")
	}

	chain, err := utils.LoadCommandChain(sandbox, io)
	if err != nil {
		return nil, err
	}

	binary := sandbox.Deps.Stringsdeps.ToLower(sandbox.Config.ProjectName)
	if project_conf, err := utils.LoadProjectConf(sandbox, io); err == nil && project_conf.Name != "" {
		binary = sandbox.Deps.Stringsdeps.ToLower(project_conf.Name)
	}
	lines := []string{sandbox.Deps.Stringsdeps.Join(append([]string{binary}, props.Argv...), " ")}
	consumed := make([]bool, len(props.Argv))
	answered := ""
	failure := ""

	for _, entry := range chain {
		name := utils.CommandIdentifier(sandbox, entry.Name)
		head := sandbox.Deps.Std.Sprintf("  %-4d %-24s", entry.Conf.Priority, name)

		if answered != "" {
			lines = append(lines, head+" not reached: "+answered+" answered first")
			continue
		}

		match := utils.MatchCommandArgv(sandbox, entry.Conf, props.Argv, consumed)
		if !match.Runs {
			lines = append(lines, head+" skipped: "+match.Reason)
			continue
		}
		if match.Failure != "" {
			lines = append(lines, head+" answers "+match.Failure+", before its handler runs")
			answered, failure = name, match.Failure
			continue
		}
		if !entry.Conf.Strict {
			lines = append(lines, head+" runs, and hands the line on unless it answers")
			continue
		}
		lines = append(lines, head+" runs")
		answered = name
	}

	lines = append(lines, "")
	switch {
	case failure != "":
		lines = append(lines, "the line ends on "+answered+": "+failure)
	case answered != "":
		lines = append(lines, "the line ends on "+answered)
	case len(props.Argv) == 0:
		lines = append(lines, "no command runs: the empty line prints the general help and exits 0")
	case nearCommand(sandbox, chain, props.Argv) != "":
		lines = append(lines, "the line starts with the verb of "+nearCommand(sandbox, chain, props.Argv)+
			" but does not fit its args: handle_bad_usage.go answers it, naming what is wrong, exit 2")
	default:
		lines = append(lines, "no command answers: handle_not_found.go answers it, exit 2")
	}
	return lines, nil
}

// nearCommand is the strict command whose verb — the longest of its literal
// identifiers — the line's segments start with, "" when none: the dispatch
// answers such a line with a usage error rather than a not-found.
func nearCommand(sandbox *api.Sandbox, chain []utils.CommandChainEntry, argv []string) string {
	segments, _ := utils.SplitCommandArgv(sandbox, argv)
	near := ""
	longest := 0
	for _, entry := range chain {
		if !entry.Conf.Strict {
			continue
		}
		for _, identifier := range entry.Conf.Identifiers() {
			words := sandbox.Deps.Stringsdeps.Fields(identifier)
			if len(words) == 0 || len(words) > len(segments) || len(words) <= longest {
				continue
			}
			if sandbox.Deps.Stringsdeps.Join(segments[:len(words)], " ") == identifier {
				near, longest = utils.CommandIdentifier(sandbox, entry.Name), len(words)
			}
		}
	}
	return near
}
