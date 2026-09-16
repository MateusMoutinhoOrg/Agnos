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
	lineUp     = "\033[A"
)

// ErrCancelled reports that the person answering asked to stop — ctrl-c, q or
// a lone escape. It is an error because no answer can be had, which is the one
// thing the contract says an error means.
var ErrCancelled = errors.New("interview cancelled")

// ErrNoInput reports that the input ended before the question was answered.
var ErrNoInput = errors.New("no input left to answer with")

// key is one decoded keypress. Raw mode hands over bytes, and a menu only
// cares about the six meanings below.
type key int

const (
	keyUnknown key = iota
	keyUp
	keyDown
	keyEnter
	keySpace
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
func readKey() (key, error) {
	buffer := make([]byte, 1)
	if _, err := os.Stdin.Read(buffer); err != nil {
		return keyUnknown, ErrNoInput
	}

	switch buffer[0] {
	case 0x03, 0x04, 'q':
		return keyCancel, nil
	case '\r', '\n':
		return keyEnter, nil
	case ' ':
		return keySpace, nil
	case 'k':
		return keyUp, nil
	case 'j':
		return keyDown, nil
	case 0x1b:
		return readEscape()
	}

	return keyUnknown, nil
}

// readEscape decodes the two bytes that follow an escape. Anything that is not
// a recognised arrow is ignored rather than guessed at, and an input that ends
// mid-sequence is a cancel.
func readEscape() (key, error) {
	buffer := make([]byte, 2)
	if _, err := os.Stdin.Read(buffer); err != nil {
		return keyCancel, nil
	}
	if buffer[0] != '[' {
		return keyUnknown, nil
	}

	switch buffer[1] {
	case 'A':
		return keyUp, nil
	case 'B':
		return keyDown, nil
	}

	return keyUnknown, nil
}

// readLine reads one answer in line mode. An empty line is a valid answer — it
// is how a caller offering a default learns the default was taken — so only an
// input that ended with nothing on it is an error.
func readLine() (string, error) {
	line, err := lineReader.ReadString('\n')
	if err != nil && line == "" {
		return "", ErrNoInput
	}
	return strings.TrimRight(line, "\r\n"), nil
}
