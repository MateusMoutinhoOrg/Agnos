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
		lines = append(lines, "no command runs: the empty line prints the general help and exits 2")
	default:
		lines = append(lines, "no command answers: handle_not_found.go answers it, exit 2")
	}
	return lines, nil
}
