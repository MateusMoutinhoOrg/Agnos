package interviewer

import (
	"fmt"
	"os"
	"unicode/utf8"
)

// askText reads one line of free text. Where the terminal allows it the line is
// read key by key, which is the whole point: a line read in cooked mode cannot
// see the escape key, so "go back" had to be spelled with a letter and any
// question expecting a word swallowed it as the answer. Reading the keys
// ourselves keeps every printable byte an answer and leaves escape meaning
// escape.
//
// A pipe, a file and a terminal that will not go raw all fall back to the line
// reader, where the back token is the way back instead.
func askText(question string) (string, error) {
	if !interactive() || !rawOn() {
		prompt(question, "\n")
		return readLine()
	}
	defer rawOff()

	prompt(question, "\r\n")
	return editLine()
}

// editor is what has been typed so far together with where the cursor sits in
// it. The cursor is a byte index rather than a rune one because raw mode hands
// over one byte at a time: a letter outside ascii arrives in pieces, and each
// piece is written where the one before it left off, so the bytes of a rune
// stay whole and in order.
type editor struct {
	typed  string
	cursor int
}

// insert writes one byte where the cursor is and leaves the cursor after it.
// The byte is put back as a byte and never as a rune: a letter outside ascii
// arrives in pieces, and turning each piece into a rune of its own would
// re-encode it into mojibake.
func (state *editor) insert(letter []byte) {
	state.typed = state.typed[:state.cursor] + string(letter) + state.typed[state.cursor:]
	state.cursor += len(letter)
}

// left and right move a whole rune at a time, so one press of an arrow steps
// over an accented letter instead of landing inside it.
func (state *editor) left() {
	if state.cursor == 0 {
		return
	}

	_, size := utf8.DecodeLastRuneInString(state.typed[:state.cursor])
	state.cursor -= size
}

func (state *editor) right() {
	if state.cursor == len(state.typed) {
		return
	}

	_, size := utf8.DecodeRuneInString(state.typed[state.cursor:])
	state.cursor += size
}

// home and end are the two ends of the line — the arrows' own keys, and
// ctrl-a / ctrl-e beside them.
func (state *editor) home() {
	state.cursor = 0
}

func (state *editor) end() {
	state.cursor = len(state.typed)
}

// backspace drops the rune before the cursor, delete the one after it. Both
// work a rune at a time for the same reason the arrows do.
func (state *editor) backspace() {
	if state.cursor == 0 {
		return
	}

	_, size := utf8.DecodeLastRuneInString(state.typed[:state.cursor])
	state.typed = state.typed[:state.cursor-size] + state.typed[state.cursor:]
	state.cursor -= size
}

func (state *editor) delete() {
	if state.cursor == len(state.typed) {
		return
	}

	_, size := utf8.DecodeRuneInString(state.typed[state.cursor:])
	state.typed = state.typed[:state.cursor] + state.typed[state.cursor+size:]
}

// render paints the whole line again from the left margin and then walks the
// cursor back to where it belongs. Printing the answer as it arrived was
// enough while the only edit was at the end; a cursor that moves means what is
// on screen is redrawn, because an insert in the middle changes every column
// after it.
func (state *editor) render() {
	fmt.Fprintf(os.Stdout, "\r  %s>%s %s%s", gray, reset, state.typed, clearLine)

	if back := utf8.RuneCountInString(state.typed[state.cursor:]); back > 0 {
		fmt.Fprintf(os.Stdout, cursorBack, back)
	}
}

// editLine is the line editor itself: printable bytes are inserted at the
// cursor, the arrows move inside the line, backspace and delete take a rune
// off either side of it, enter finishes, escape steps back to the previous
// question and ctrl-c ends the session. Raw mode does no carriage return of
// its own, so every line break here is a \r\n.
func editLine() (string, error) {
	state := &editor{}

	for {
		buffer := make([]byte, 1)
		if _, err := os.Stdin.Read(buffer); err != nil {
			fmt.Fprint(os.Stdout, "\r\n")
			return "", ErrNoInput
		}

		switch letter := buffer[0]; {
		case letter == '\r' || letter == '\n':
			fmt.Fprint(os.Stdout, "\r\n")
			return state.typed, nil

		case letter == 0x03:
			fmt.Fprint(os.Stdout, "\r\n")
			return "", ErrCancelled

		case letter == 0x04:
			fmt.Fprint(os.Stdout, "\r\n")
			return "", ErrNoInput

		case letter == 0x1b:
			// An escape the person pressed has nothing behind it; one a
			// terminal sent is the head of a sequence, which is read off
			// and applied rather than taken as a step back.
			rest, timedOut := readPending(2)
			if timedOut || len(rest) == 0 {
				fmt.Fprint(os.Stdout, "\r\n")
				return "", ErrBack
			}

			state.escape(readSequence(rest))

		case letter == 0x7f || letter == 0x08:
			state.backspace()

		case letter == 0x01:
			state.home()

		case letter == 0x05:
			state.end()

		case letter >= 0x20:
			state.insert(buffer[:1])
		}

		state.render()
	}
}

// escape applies a sequence the terminal sent. Only the keys that move inside
// the line or take a rune out of it mean anything here; the ones that would
// walk a history there is none of are dropped, so nothing a terminal sends
// ends up typed into the answer.
func (state *editor) escape(sequence string) {
	switch sequence {
	case "[D", "OD":
		state.left()
	case "[C", "OC":
		state.right()
	case "[H", "OH", "[1~", "[7~":
		state.home()
	case "[F", "OF", "[4~", "[8~":
		state.end()
	case "[3~":
		state.delete()
	}
}

// readSequence finishes reading a control sequence whose first bytes are
// already in hand. A sequence ends on a byte between @ and ~, and the keys
// that carry a parameter — delete is esc [ 3 ~, and a terminal may spell the
// arrows with one too — are longer than the two bytes an arrow takes. Reading
// to the end is what keeps the tail of a long sequence out of the answer.
func readSequence(rest []byte) string {
	sequence := string(rest)

	for sequence[0] == '[' && (len(sequence) < 2 || !sequenceEnd(sequence[len(sequence)-1])) {
		more, timedOut := readPending(1)
		if timedOut {
			break
		}

		sequence += string(more)
	}

	return sequence
}

// sequenceEnd reports whether a byte is the final one of a control sequence.
func sequenceEnd(letter byte) bool {
	return letter >= 0x40 && letter <= 0x7e
}
