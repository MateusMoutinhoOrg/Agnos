package interviewer

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"
)

// ─── ANSI escape sequences ──────────────────────────────────────────────────
//
// The same palette sandbox/internal/commands/help/handler.go prints its
// screens with, so an interview and a help screen read as one interface.

const (
	bold       = "\033[1m"
	dim        = "\033[2m"
	reset      = "\033[0m"
	cyan       = "\033[36m"
	green      = "\033[32m"
	yellow     = "\033[33m"
	white      = "\033[97m"
	gray       = "\033[90m"
	cursorHide = "\033[?25l"
	cursorShow = "\033[?25h"
	clearBelow = "\033[J"
	clearLine  = "\033[K"
	cursorBack = "\033[%dD"
	lineUp     = "\033[A"
)

// ErrCancelled reports that the person answering asked to stop the whole
// session — ctrl-c or ctrl-d. It is an error because no answer can be had,
// which is the one thing the contract says an error means.
var ErrCancelled = errors.New("interview cancelled")

// ErrBack reports that the person asked to step back to the question before
// this one. It is an error for the same reason ErrCancelled is — this question
// has no answer — and the caller tells the two apart through the contract's
// Back, which is what makes going back undoable instead of fatal.
var ErrBack = errors.New("going back")

// ErrNoInput reports that the input ended before the question was answered.
var ErrNoInput = errors.New("no input left to answer with")

// backToken is how a session with no terminal asks to go back. A piped stdin
// has no escape key to press — every byte of it is the answer — so the one
// line that means "back" is spelled out instead.
const backToken = ":back"

// isBack fills interviewer.Sandbox.Back: of the errors a question returns,
// only this one is undone by asking again.
func isBack(err error) bool {
	return errors.Is(err, ErrBack)
}

// key is one decoded keypress. Raw mode hands over bytes, and a menu only
// cares about the seven meanings below.
type key int

const (
	keyUnknown key = iota
	keyUp
	keyDown
	keyEnter
	keySpace
	keyBack
	keyCancel
)

// lineReader is the one buffered reader every line-mode question reads
// through. It is a package-level value because a second reader over os.Stdin
// would buffer ahead of the first and swallow an answer meant for it.
var lineReader = bufio.NewReader(os.Stdin)

// interactive reports whether stdin is a terminal a menu can be painted on.
// A pipe, a file and /dev/null all answer false, and every question falls back
// to the numbered form.
func interactive() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// rawOn puts the terminal in raw mode, so a keypress arrives as it is typed
// rather than at the end of a line, and is not echoed. It reports whether it
// worked: stty is absent on Windows and in a stripped container, and a false
// here is what sends a question to its fallback.
func rawOn() bool {
	return stty("raw", "-echo")
}

// rawOff restores the terminal. Every path that called rawOn calls this one,
// including the error paths: a process that exits in raw mode leaves the shell
// behind it unusable.
func rawOff() {
	stty("-raw", "echo")
}

// stty runs the terminal driver with stdin inherited, which is what makes it
// act on this process's terminal rather than on a pipe of its own.
func stty(args ...string) bool {
	command := exec.Command("stty", args...)
	command.Stdin = os.Stdin
	return command.Run() == nil
}

// readKey reads one keypress in raw mode. Arrow keys arrive as an escape
// sequence (0x1b 0x5b 0x41 is up), and j/k are accepted beside them so the
// menu answers to the hands that expect either.
//
// The two ways out are kept apart on purpose: ctrl-c and ctrl-d end the
// session, while escape and q step back one question. They used to be the same
// key, which made "go back" delete every answer already given.
func readKey() (key, error) {
	buffer := make([]byte, 1)
	if _, err := os.Stdin.Read(buffer); err != nil {
		return keyUnknown, ErrNoInput
	}

	switch buffer[0] {
	case 0x03, 0x04:
		return keyCancel, nil
	case 'q':
		return keyBack, nil
	case '\r', '\n':
		return keyEnter, nil
	case ' ':
		return keySpace, nil
	case 'k':
		return keyUp, nil
	case 'j':
		return keyDown, nil
	case 0x1b:
		return readEscape(), nil
	}

	return keyUnknown, nil
}

// readEscape decodes what follows an escape byte. An escape on its own is the
// escape key, and an escape followed by a bracket is an arrow — the two are
// told apart by how long the rest takes to arrive, because nothing follows a
// key the person pressed themselves.
func readEscape() key {
	rest, timedOut := readPending(2)
	if timedOut || len(rest) < 2 || rest[0] != '[' {
		return keyBack
	}

	switch rest[1] {
	case 'A':
		return keyUp
	case 'B':
		return keyDown
	}

	return keyUnknown
}

// readPending reads up to size bytes that are already on their way, reporting
// true when none arrived. VMIN 0 with VTIME 1 is a tenth of a second of
// patience — long enough for a terminal to finish an escape sequence it is
// already sending, short enough that a pressed escape key does not hang.
func readPending(size int) ([]byte, bool) {
	if !stty("min", "0", "time", "1") {
		return nil, true
	}
	defer stty("min", "1", "time", "0")

	buffer := make([]byte, size)
	read, err := os.Stdin.Read(buffer)
	if err != nil || read == 0 {
		return nil, true
	}
	return buffer[:read], false
}

// readLine reads one answer in line mode. An empty line is a valid answer — it
// is how a caller offering a default learns the default was taken — so only an
// input that ended with nothing on it is an error. The one line that is not an
// answer at all is the back token.
func readLine() (string, error) {
	line, err := lineReader.ReadString('\n')
	if err != nil && line == "" {
		return "", ErrNoInput
	}

	answer := strings.TrimRight(line, "\r\n")
	if strings.TrimSpace(answer) == backToken {
		return "", ErrBack
	}
	return answer, nil
}
