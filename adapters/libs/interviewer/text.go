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

// editLine is the line editor itself: printable bytes are echoed and kept,
// backspace takes the last rune back, enter finishes, escape steps back to the
// previous question and ctrl-c ends the session. Raw mode does no carriage
// return of its own, so every line break here is a \r\n.
func editLine() (string, error) {
	typed := ""

	for {
		buffer := make([]byte, 1)
		if _, err := os.Stdin.Read(buffer); err != nil {
			fmt.Fprint(os.Stdout, "\r\n")
			return "", ErrNoInput
		}

		switch letter := buffer[0]; {
		case letter == '\r' || letter == '\n':
			fmt.Fprint(os.Stdout, "\r\n")
			return typed, nil

		case letter == 0x03:
			fmt.Fprint(os.Stdout, "\r\n")
			return "", ErrCancelled

		case letter == 0x04:
			fmt.Fprint(os.Stdout, "\r\n")
			return "", ErrNoInput

		case letter == 0x1b:
			// An escape the person pressed has nothing behind it; one a
			// terminal sent is the head of an arrow sequence, which is
			// read off and ignored rather than taken as a step back.
			if rest, timedOut := readPending(2); timedOut || len(rest) == 0 {
				fmt.Fprint(os.Stdout, "\r\n")
				return "", ErrBack
			}

		case letter == 0x7f || letter == 0x08:
			typed = eraseLast(typed)

		case letter >= 0x20:
			// The byte is kept and echoed as it arrived. A letter outside
			// ascii reaches us one byte at a time, and converting each of
			// them to a rune of its own would re-encode it into mojibake.
			typed += string(buffer[:1])
			os.Stdout.Write(buffer[:1])
		}
	}
}

// eraseLast drops the last rune of what has been typed and takes it off the
// screen with it. It works a rune at a time rather than a byte at a time, so
// one backspace on an accented letter removes the letter instead of half of it.
func eraseLast(typed string) string {
	if typed == "" {
		return typed
	}

	_, size := utf8.DecodeLastRuneInString(typed)
	fmt.Fprint(os.Stdout, "\b \b")
	return typed[:len(typed)-size]
}
