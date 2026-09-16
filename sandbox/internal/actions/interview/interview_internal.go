package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	interviewer "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/interviewer"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// selfVerb is the one command the interview never offers: answering questions
// in order to start another interview is a loop with no way out of it.
const selfVerb = "interview"

// The rows the session adds to a menu of its own making. They carry the same
// leading NUL as the field rows, so no declared value collides with them.
const (
	exitOptionId = "\x00exit"
	backOptionId = "\x00back"
	runOptionId  = "\x00run"
)

// InterviewInternal is the session: a menu of the categories the command
// surface declares, a menu of the commands in one of them, one question per
// field that command declares, a confirm screen showing the command line the
// answers add up to, and then the command itself — over and over until the
// person asks to stop.
//
// Nothing about any command is spelled here. Cli.Commands already holds every
// declaration — its category, its help, and each flag and arg with its type,
// its bounds, whether it repeats and what it falls back to — so the questions
// are generated from the declarations and a command declared tomorrow is
// covered without this file changing.
func InterviewInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	printWelcome(sandbox, path)

	for {
		command, chosen, err := chooseCommand(sandbox)
		if err != nil {
			return endSession(sandbox, err)
		}
		if !chosen {
			break
		}

		values, err := AskValues(sandbox, io, command, path)
		if err != nil {
			return endSession(sandbox, err)
		}

		confirmed, err := confirmPlan(sandbox, io, command, values, path)
		if err != nil {
			return endSession(sandbox, err)
		}
		if !confirmed {
			continue
		}

		runCommand(sandbox, command, values)
	}

	sandbox.Deps.Std.Printf("  %sNothing else then. Run %s%s interview%s%s whenever you want me back.%s\n\n",
		dim, reset+cyan, binaryName(sandbox), reset, dim, reset)
	return nil
}

// endSession closes an interview whose question could not be answered — the
// person pressed ctrl-c, or the input ran out. Neither is a failure of the
// command, so the session ends quietly and the process exits ok.
func endSession(sandbox *api.Sandbox, reason error) error {
	sandbox.Deps.Std.Printf("\n")
	sandbox.Deps.Std.Log("interview ended: %s \n", reason.Error())
	return nil
}

// ─── Choosing a command ─────────────────────────────────────────────────────

// chooseCommand walks the two menus that lead to one command, and reports
// false when the person chose to stop instead.
func chooseCommand(sandbox *api.Sandbox) (api.Command, bool, error) {
	for {
		category, err := sandbox.Deps.Interviewer.SingleAlternativeQuestion(
			"What do you want to do?",
			categoryRows(sandbox),
		)
		if err != nil {
			return api.Command{}, false, err
		}
		if category == exitOptionId {
			return api.Command{}, false, nil
		}

		command, chosen, err := chooseInCategory(sandbox, category)
		if err != nil {
			return api.Command{}, false, err
		}
		if chosen {
			return command, true, nil
		}
	}
}

// chooseInCategory offers the commands of one category, and reports false when
// the person asked to go back to the categories.
func chooseInCategory(sandbox *api.Sandbox, category string) (api.Command, bool, error) {
	commands := commandsIn(sandbox, category)

	rows := []interviewer.AlternativeOption{}
	for _, command := range commands {
		rows = append(rows, interviewer.AlternativeOption{
			Id:  verbOf(command),
			Msg: labelled(sandbox, verbOf(command), command.Help),
		})
	}
	rows = append(rows, interviewer.AlternativeOption{Id: backOptionId, Msg: "· back"})

	chosen, err := sandbox.Deps.Interviewer.SingleAlternativeQuestion(category, rows)
	if err != nil {
		return api.Command{}, false, err
	}
	if chosen == backOptionId {
		return api.Command{}, false, nil
	}

	for _, command := range commands {
		if verbOf(command) == chosen {
			return command, true, nil
		}
	}

	return api.Command{}, false, nil
}

// categoryRows is one row per category the command surface declares, in the
// order the help screen groups them, plus the row that ends the session.
func categoryRows(sandbox *api.Sandbox) []interviewer.AlternativeOption {
	rows := []interviewer.AlternativeOption{}
	seen := map[string]bool{}

	for _, command := range sandbox.Cli.Commands {
		category := categoryOf(command)
		if !offerable(command) || seen[category] {
			continue
		}
		seen[category] = true
		rows = append(rows, interviewer.AlternativeOption{Id: category, Msg: category})
	}

	return append(rows, interviewer.AlternativeOption{Id: exitOptionId, Msg: "· exit"})
}

// commandsIn is every offerable command of one category, in declaration order.
func commandsIn(sandbox *api.Sandbox, category string) []api.Command {
	commands := []api.Command{}
	for _, command := range sandbox.Cli.Commands {
		if offerable(command) && categoryOf(command) == category {
			commands = append(commands, command)
		}
	}
	return commands
}

// offerable reports whether a command belongs on a menu: a hidden one is kept
// off every listing, one with no identifier cannot be named, and the interview
// does not offer itself.
func offerable(command api.Command) bool {
	return !command.Hidden && len(command.Identifiers) > 0 && verbOf(command) != selfVerb
}

// categoryOf is the heading a command is listed under, with the same fallback
// the help screen uses for one that declares none.
func categoryOf(command api.Command) string {
	if command.Category == "" {
		return "Other"
	}
	return command.Category
}

// ─── Confirming and running ─────────────────────────────────────────────────

// confirmPlan shows the command line the answers add up to and offers to run
// it, to change one answer, or to go back to the menu without running
// anything. Changing an answer re-asks that one field and nothing else.
func confirmPlan(sandbox *api.Sandbox, io *smartio.SmartIO, command api.Command, values map[string][]any, path string) (bool, error) {
	for {
		printPlan(sandbox, command, values)

		chosen, err := sandbox.Deps.Interviewer.SingleAlternativeQuestion("Run it?", planRows(sandbox, command, values))
		if err != nil {
			return false, err
		}

		switch chosen {
		case runOptionId:
			return true, nil
		case backOptionId:
			return false, nil
		}

		field, found := findField(command, chosen)
		if !found {
			continue
		}

		bound, err := ResolveField(sandbox, io, command, field, path)
		if err != nil {
			return false, err
		}
		if len(bound) > 0 {
			values[field.Id] = bound
			continue
		}
		delete(values, field.Id)
	}
}

// planRows is the confirm menu: run it, one row per answer that can be
// changed, and a way back. The two fields the interview answers by itself are
// not offered — the path is the session's and the progress stays visible.
func planRows(sandbox *api.Sandbox, command api.Command, values map[string][]any) []interviewer.AlternativeOption {
	rows := []interviewer.AlternativeOption{{Id: runOptionId, Msg: "run it"}}

	for _, field := range FieldsOf(command) {
		if AnsweredForYou(command, field) {
			continue
		}
		rows = append(rows, interviewer.AlternativeOption{
			Id:  field.Id,
			Msg: sandbox.Deps.Std.Sprintf("change %s  %s(%s)%s", label(field), dim, currentText(sandbox, values, field), reset),
		})
	}

	return append(rows, interviewer.AlternativeOption{Id: backOptionId, Msg: "· back, run nothing"})
}

// currentText is what one answer holds now, for the row that offers to change
// it.
func currentText(sandbox *api.Sandbox, values map[string][]any, field Field) string {
	bound := values[field.Id]
	if len(bound) == 0 {
		return "unset"
	}

	texts := []string{}
	for _, value := range bound {
		texts = append(texts, valueText(sandbox, value))
	}

	return sandbox.Deps.Stringsdeps.Join(texts, ", ")
}

// findField reads one field of a command back by its id.
func findField(command api.Command, id string) (Field, bool) {
	for _, field := range FieldsOf(command) {
		if field.Id == id {
			return field, true
		}
	}
	return Field{}, false
}

// runCommand binds the answers onto a copy of the declaration and hands it to
// the command's own handler — the same call the cli dispatch makes, with the
// values coming from questions instead of from a command line. The handler
// runs its own action, which persists and builds; the interview writes
// nothing itself.
//
// A command that fails is reported and the session goes on: a wrong answer to
// one question is no reason to lose the rest of a session.
func runCommand(sandbox *api.Sandbox, command api.Command, values map[string][]any) {
	bound := api.BindCommand(&command)
	for id, list := range values {
		bound.Items[id] = list
	}

	if bound.Handler == nil {
		printOutcome(sandbox, command, api.ExitFailure)
		return
	}

	printOutcome(sandbox, command, bound.Handler(bound))
}
