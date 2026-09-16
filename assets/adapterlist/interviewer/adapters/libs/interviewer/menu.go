package interviewer

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	interviewer "{{.Module}}/sandbox/deps/interviewer"
)

// runMenu asks one alternative question and returns the indexes chosen. It is
// the whole of the difference between the two menu questions: a single-choice
// menu finishes on enter with the row under the cursor, a multiple-choice one
// toggles rows with space and finishes on enter with whatever is ticked.
//
// Escape and q answer neither: they report ErrBack, which the caller undoes by
// asking the question before this one again. Only ctrl-c and ctrl-d end a
// session.
//
// A terminal that cannot be put in raw mode answers the same question as a
// numbered list read line by line, so every question has an answer whatever it
// is being run under.
func runMenu(question string, options []interviewer.AlternativeOption, multiple bool) ([]int, error) {
	if len(options) == 0 {
		return []int{}, nil
	}

	if !interactive() || !rawOn() {
		return numberedMenu(question, options, multiple)
	}
	defer rawOff()

	fmt.Fprint(os.Stdout, cursorHide)
	defer fmt.Fprint(os.Stdout, cursorShow)

	printQuestion(question, multiple)

	cursor := 0
	ticked := make([]bool, len(options))
	painted := 0

	for {
		painted = paintMenu(options, cursor, ticked, multiple, painted)

		pressed, err := readKey()
		if err != nil {
			return nil, err
		}

		switch pressed {
		case keyUp:
			cursor = (cursor - 1 + len(options)) % len(options)
		case keyDown:
			cursor = (cursor + 1) % len(options)
		case keySpace:
			if multiple {
				ticked[cursor] = !ticked[cursor]
			}
		case keyBack:
			fmt.Fprint(os.Stdout, "\r\n")
			return nil, ErrBack
		case keyCancel:
			fmt.Fprint(os.Stdout, "\r\n")
			return nil, ErrCancelled
		case keyEnter:
			if !multiple {
				ticked[cursor] = true
			}
			fmt.Fprint(os.Stdout, "\r\n")
			return chosenIndexes(ticked), nil
		}
	}
}

// printQuestion writes the prompt above the rows. It is printed once and never
// repainted, so the rows below it are the only thing that moves.
func printQuestion(question string, multiple bool) {
	fmt.Fprintf(os.Stdout, "\r\n  %s%s%s\r\n", bold+cyan, question, reset)
	if multiple {
		fmt.Fprintf(os.Stdout, "  %sspace toggles, enter confirms, esc goes back%s\r\n", dim, reset)
		return
	}
	fmt.Fprintf(os.Stdout, "  %s↑↓ to move, enter confirms, esc goes back%s\r\n", dim, reset)
}

// paintMenu repaints the rows in place: it walks the cursor back over the rows
// it printed last time, clears from there down and writes them again. painted
// is how many rows were on screen, and the count it returns is how many are
// now. Raw mode does no carriage return of its own, so every line break here
// is a \r\n.
func paintMenu(options []interviewer.AlternativeOption, cursor int, ticked []bool, multiple bool, painted int) int {
	if painted > 0 {
		fmt.Fprint(os.Stdout, "\r"+strings.Repeat(lineUp, painted)+clearBelow)
	}

	for index, option := range options {
		pointer := "  "
		label := option.Msg
		if index == cursor {
			pointer = green + "❯" + reset + " "
			label = bold + white + option.Msg + reset
		}

		box := ""
		if multiple {
			box = dim + "[ ] " + reset
			if ticked[index] {
				box = green + "[x] " + reset
			}
		}

		fmt.Fprintf(os.Stdout, "  %s%s%s\r\n", pointer, box, label)
	}

	return len(options)
}

// numberedMenu is the same question without a terminal: the options are
// printed once, numbered, and the answer is read as a line. A multiple-choice
// question takes the numbers separated by spaces or commas, and an empty line
// chooses none.
func numberedMenu(question string, options []interviewer.AlternativeOption, multiple bool) ([]int, error) {
	for {
		fmt.Fprintf(os.Stdout, "\n  %s%s%s\n", bold+cyan, question, reset)
		for index, option := range options {
			fmt.Fprintf(os.Stdout, "    %s%d)%s %s\n", yellow, index+1, reset, option.Msg)
		}

		if multiple {
			fmt.Fprintf(os.Stdout, "  %snumbers separated by spaces, empty for none:%s ", dim, reset)
		} else {
			fmt.Fprintf(os.Stdout, "  %snumber [1]:%s ", dim, reset)
		}

		answer, err := readLine()
		if err != nil {
			return nil, err
		}

		chosen, ok := parseChoices(answer, len(options), multiple)
		if ok {
			return chosen, nil
		}

		fmt.Fprintf(os.Stdout, "  %sthat is not one of the options%s\n", yellow, reset)
	}
}

// parseChoices reads a numbered answer into indexes, reporting whether every
// number named a row. An empty answer is the first row for a single choice and
// no rows at all for a multiple one.
func parseChoices(answer string, size int, multiple bool) ([]int, bool) {
	answer = strings.TrimSpace(answer)

	if answer == "" {
		if multiple {
			return []int{}, true
		}
		return []int{0}, true
	}

	fields := strings.FieldsFunc(answer, func(letter rune) bool {
		return letter == ' ' || letter == ',' || letter == '\t'
	})

	if !multiple && len(fields) != 1 {
		return nil, false
	}

	ticked := make([]bool, size)
	for _, field := range fields {
		number, err := strconv.Atoi(field)
		if err != nil || number < 1 || number > size {
			return nil, false
		}
		ticked[number-1] = true
	}

	return chosenIndexes(ticked), true
}

// chosenIndexes reads a tick list back as the indexes it holds, in declaration
// order — so an answer never depends on the order the rows were ticked in.
func chosenIndexes(ticked []bool) []int {
	chosen := []int{}
	for index, on := range ticked {
		if on {
			chosen = append(chosen, index)
		}
	}
	return chosen
}
