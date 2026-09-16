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
	stepPrefix   = "\x00step:"
)

// InterviewInternal is the session: the steps this project has not taken and
// the areas it already has, a menu of the commands in one of them, one
// question per field that command declares, a confirm screen showing the
// command line the answers add up to, and then the command itself — over and
// over until the person asks to stop.
//
// Nothing about any command is spelled here. Cli.Commands already holds every
// declaration — its category, its help, and each flag and arg with its type,
// its bounds, whether it repeats and what it falls back to — so the questions
// are generated from the declarations and a command declared tomorrow is
// covered without this file changing. What state.go adds on top is not a
// command list but a filter: which of those declarations the project in front
// of the person can actually run.
func InterviewInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	printWelcome(sandbox, io, path)

	for {
		command, chosen, err := chooseCommand(sandbox, io)
		if err != nil {
			return endSession(sandbox, err)
		}
		if !chosen {
			break
		}

		values, asked, err := AskValues(sandbox, io, command, path)
		if err != nil {
			return endSession(sandbox, err)
		}
		if !asked {
			continue
		}

		if err := runPlan(sandbox, io, command, values, path); err != nil {
			return endSession(sandbox, err)
		}
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

// chooseCommand walks the menus that lead to one command, and reports false
// when the person chose to stop instead. The state is read again on every
// pass: a command that turned a mechanic on opens its area, and the step that
// turned it on is gone from the next menu.
func chooseCommand(sandbox *api.Sandbox, io *smartio.SmartIO) (api.Command, bool, error) {
	for {
		state := readState(sandbox, io)

		chosen, err := sandbox.Deps.Interviewer.SingleAlternativeQuestion(
			"What do you want to do?",
			topRows(sandbox, state),
		)
		if err != nil {
			// There is no question before the first one, so going back from
			// here is going out — the same answer the exit row gives, and the
			// session still ends ok.
			if sandbox.Deps.Interviewer.Back(err) {
				return api.Command{}, false, nil
			}
			return api.Command{}, false, err
		}
		if chosen == exitOptionId {
			return api.Command{}, false, nil
		}

		if verb, isStep := stepVerb(sandbox, chosen); isStep {
			if command, found := commandByVerb(sandbox, verb); found {
				return command, true, nil
			}
			continue
		}

		command, picked, err := chooseInCategory(sandbox, state, chosen)
		if err != nil {
			return api.Command{}, false, err
		}
		if picked {
			return command, true, nil
		}
	}
}

// topRows is the first menu: the steps this project has not taken, the areas
// it has, and the row that ends the session. A step is a command too — it is
// listed first, in plain words, because it is what the project needs next.
func topRows(sandbox *api.Sandbox, state projectState) []interviewer.AlternativeOption {
	rows := []interviewer.AlternativeOption{}
	suggested := false

	for _, one := range nextSteps(state) {
		if _, found := commandByVerb(sandbox, one.Verb); !found {
			continue
		}

		mark := "  "
		if len(rows) == 0 && one.Key {
			mark = "★ "
			suggested = true
		}

		rows = append(rows, interviewer.AlternativeOption{
			Id:  stepPrefix + one.Verb,
			Msg: sandbox.Deps.Std.Sprintf("%s%s  %s(%s)%s", mark, one.Msg, dim, one.Verb, reset),
		})
	}

	// A step with no ★ on it is an offer, not a recommendation, and every one
	// of them installs a whole layer. Leading the menu with one of those puts a
	// http server, or a dependency layer, under an enter pressed without
	// reading — so the way out leads instead, and the offers keep their places.
	offersOnly := len(rows) > 0 && !suggested

	rows = append(rows, categoryRows(sandbox, state)...)

	return SafeFirst(rows, interviewer.AlternativeOption{Id: exitOptionId, Msg: "· exit"}, offersOnly)
}

// stepVerb reads a chosen row back as the command a step runs, reporting false
// for every row that is not one.
func stepVerb(sandbox *api.Sandbox, chosen string) (string, bool) {
	if !sandbox.Deps.Stringsdeps.HasPrefix(chosen, stepPrefix) {
		return "", false
	}
	return sandbox.Deps.Stringsdeps.TrimPrefix(chosen, stepPrefix), true
}

// commandByVerb reads one declaration of the surface back by the name it is
// typed as.
func commandByVerb(sandbox *api.Sandbox, verb string) (api.Command, bool) {
	for _, command := range sandbox.Cli.Commands {
		if verbOf(command) == verb {
			return command, true
		}
	}
	return api.Command{}, false
}

// chooseInCategory offers the commands of one area, and reports false when the
// person asked to go back to the first menu.
func chooseInCategory(sandbox *api.Sandbox, state projectState, category string) (api.Command, bool, error) {
	commands := commandsIn(sandbox, state, category)

	rows := []interviewer.AlternativeOption{}
	for _, command := range commands {
		rows = append(rows, interviewer.AlternativeOption{
			Id:  verbOf(command),
			Msg: labelled(sandbox, verbOf(command), command.Help),
		})
	}

	// An area lists its commands in declaration order, which puts cli-purge and
	// remove-command next to add-command with nothing to tell them apart. When
	// the one that happens to come first takes something away, the way back
	// leads the menu instead.
	leads := len(commands) > 0 && DestructiveVerb(sandbox, verbOf(commands[0]))
	rows = SafeFirst(rows, interviewer.AlternativeOption{Id: backOptionId, Msg: "· back"}, leads)

	chosen, err := sandbox.Deps.Interviewer.SingleAlternativeQuestion(category, rows)
	if err != nil {
		if sandbox.Deps.Interviewer.Back(err) {
			return api.Command{}, false, nil
		}
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

// categoryRows is one row per area this project has, each saying what it is
// for. The order is the areas table's — what someone new needs first — and an
// area whose mechanic is off has no command that applies, so no row is built
// for it at all.
func categoryRows(sandbox *api.Sandbox, state projectState) []interviewer.AlternativeOption {
	offered := map[string]bool{}
	declared := []string{}

	for _, command := range sandbox.Cli.Commands {
		category := categoryOf(command)
		if !offerable(state, command) || offered[category] {
			continue
		}
		offered[category] = true
		declared = append(declared, category)
	}

	rows := []interviewer.AlternativeOption{}
	emitted := map[string]bool{}

	for _, one := range areas {
		if offered[one.Name] {
			rows = append(rows, categoryRow(sandbox, one.Name))
			emitted[one.Name] = true
		}
	}
	for _, category := range declared {
		if !emitted[category] {
			rows = append(rows, categoryRow(sandbox, category))
			emitted[category] = true
		}
	}

	return rows
}

// categoryRow is one area as a menu row, indented to line up under the steps.
func categoryRow(sandbox *api.Sandbox, category string) interviewer.AlternativeOption {
	return interviewer.AlternativeOption{
		Id:  category,
		Msg: "  " + labelled(sandbox, category, areaHelp(category)),
	}
}

// commandsIn is every command of one area this project can run, in declaration
// order.
func commandsIn(sandbox *api.Sandbox, state projectState, category string) []api.Command {
	commands := []api.Command{}
	for _, command := range sandbox.Cli.Commands {
		if offerable(state, command) && categoryOf(command) == category {
			commands = append(commands, command)
		}
	}
	return commands
}

// offerable reports whether a command belongs on a menu: a hidden one is kept
// off every listing, one with no identifier cannot be named, the interview
// does not offer itself, and the rest is what the project at --path can
// actually run — the gate in state.go.
func offerable(state projectState, command api.Command) bool {
	if command.Hidden || len(command.Identifiers) == 0 || verbOf(command) == selfVerb {
		return false
	}
	return applies(state, verbOf(command), categoryOf(command))
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

// runPlan shows the confirm screen and runs what it shows, over and over until
// the command finishes ok or the person goes back to the menu.
//
// A command that fails does not cost the answers. The screen comes back with
// every one of them still on it and the row that changes the one at fault
// already on the menu — the same mechanism a change before the first run uses.
// Losing a whole questionnaire to one rejected character was the interview's
// worst moment, and this loop is the whole of the fix.
func runPlan(sandbox *api.Sandbox, io *smartio.SmartIO, command api.Command, values map[string][]any, path string) error {
	failed := api.ExitOk

	for {
		confirmed, err := confirmPlan(sandbox, io, command, values, path, failed)
		if err != nil {
			return err
		}
		if !confirmed {
			return nil
		}

		failed = runCommand(sandbox, command, values)
		if failed == api.ExitOk {
			return nil
		}
	}
}

// confirmPlan shows the command line the answers add up to and offers to run
// it, to change one answer, or to go back to the menu without running
// anything. Changing an answer re-asks that one field and nothing else — and
// then drops the answers that one has just ruled out, so the line on screen is
// never one the command would refuse.
//
// failed is the exit code of the attempt before this one, and api.ExitOk when
// there has not been one.
func confirmPlan(sandbox *api.Sandbox, io *smartio.SmartIO, command api.Command, values map[string][]any, path string, failed int) (bool, error) {
	for {
		printPlan(sandbox, io, command, values, failed)

		chosen, err := sandbox.Deps.Interviewer.SingleAlternativeQuestion("Run it?", planRows(sandbox, command, values))
		if err != nil {
			if sandbox.Deps.Interviewer.Back(err) {
				return false, nil
			}
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

		bound, err := ResolveField(sandbox, io, command, field, values, path)
		if err != nil {
			// Going back out of a change changes nothing: the answer that was
			// there stays, and the screen that offered the change returns.
			if sandbox.Deps.Interviewer.Back(err) {
				continue
			}
			return false, err
		}
		if len(bound) > 0 {
			values[field.Id] = bound
		} else {
			delete(values, field.Id)
		}

		PruneRuledOut(sandbox, command, values)
	}
}

// planRows is the confirm menu: run it, one row per answer that can be
// changed, and a way back. The two fields the interview answers by itself are
// not offered — the path is the session's and the progress stays visible — and
// neither is one the answers have ruled out, which was never asked.
func planRows(sandbox *api.Sandbox, command api.Command, values map[string][]any) []interviewer.AlternativeOption {
	changes := []interviewer.AlternativeOption{}

	for _, field := range FieldsOf(command) {
		if AnsweredForYou(command, field) || RuledOut(sandbox, command, field, values) {
			continue
		}
		changes = append(changes, interviewer.AlternativeOption{
			Id:  field.Id,
			Msg: sandbox.Deps.Std.Sprintf("change %s  %s(%s)%s", label(field), dim, currentText(sandbox, values, field), reset),
		})
	}

	run := interviewer.AlternativeOption{Id: runOptionId, Msg: "yes, run it"}
	back := interviewer.AlternativeOption{Id: backOptionId, Msg: "· no, back to the menu"}

	// "yes, run it" leads this menu because that is what the screen is for —
	// except where running it deletes a layer, where the same enter would do
	// the one thing that cannot be undone. There, saying no is the default and
	// yes is a row that has to be moved to.
	if Destructive(sandbox, command, values) {
		return append([]interviewer.AlternativeOption{back, run}, changes...)
	}

	return append(append([]interviewer.AlternativeOption{run}, changes...), back)
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
// A command that fails is reported and its exit code handed back, so the
// caller can offer the answers again instead of losing them: a wrong answer to
// one question is no reason to lose the rest of a session.
func runCommand(sandbox *api.Sandbox, command api.Command, values map[string][]any) int {
	bound := api.BindCommand(&command)
	for id, list := range values {
		bound.Items[id] = list
	}

	if bound.Handler == nil {
		printOutcome(sandbox, command, api.ExitFailure)
		return api.ExitFailure
	}

	exit := bound.Handler(bound)
	printOutcome(sandbox, command, exit)
	return exit
}
