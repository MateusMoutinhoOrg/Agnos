package interviewer

// This package is the sandbox's *copy* of the api an interactive terminal
// exposes — the same mechanic as argvdeps, iodeps and std, for the same
// reason: reading a keypress, putting the terminal in raw mode and repainting
// a menu are OS-bound effects, so `os`, `bufio` and `os/exec` may not appear
// inside the sandbox. The contract is restated here, and the adapter — which
// lives outside the sandbox — is what fills it.
//
// Nothing here says how a question is presented. An adapter is free to answer
// with a full-screen arrow-key menu, a numbered list read line by line, or a
// recorded script: the sandbox asks a question and is handed an answer.

// AlternativeOption is one choice offered by an alternative question: Id is
// what the answer returns, Msg is what the person reading it sees.
type AlternativeOption struct {
	// Id is the value the question returns when this option is chosen. It
	// is never displayed.
	Id string
	// Msg is the label shown for this option. It is never returned.
	Msg string
}

// Sandbox is the interview library injected whole as the Deps.Interviewer
// field. Every field asks one question and blocks until it is answered;
// an error reports that no answer can be had — the input ended, or the
// terminal could not be read — and never that the answer was invalid.
type Sandbox struct {
	// IntQuestion asks for a whole number and returns it already converted.
	IntQuestion func(question string) (int, error)

	// StrQuestion asks for a line of text. An empty answer is a valid one:
	// it is how a caller offering a default learns that the default was
	// taken.
	StrQuestion func(question string) (string, error)

	// FloatQuestion asks for a number and returns it already converted.
	FloatQuestion func(question string) (float64, error)

	// BoolQuestion asks a yes-or-no question and returns the answer.
	BoolQuestion func(question string) (bool, error)

	// SingleAlternativeQuestion offers a list of options and returns the Id
	// of the one chosen — never its Msg, so the wording of an option may
	// change without moving what the caller matches on.
	SingleAlternativeQuestion func(question string, alternatives []AlternativeOption) (string, error)

	// MultipleAlternativeQuestion offers a list of options and returns the
	// Ids of every one chosen, in the order the options were declared.
	// Choosing none is a valid answer and returns an empty slice.
	MultipleAlternativeQuestion func(question string, alternatives []AlternativeOption) ([]string, error)

	// Back reports whether an error a question returned means the person
	// asked to step back to the question before it, rather than that no
	// answer can be had at all. It is the one way a caller tells the two
	// apart, so that going back can be undone by asking again while every
	// other error ends the session.
	Back func(err error) bool
}
