package interviewer

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	interviewer "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/interviewer"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

// prompt writes one question and leaves the cursor after it. The line break is
// the caller's because a question read key by key is printed in raw mode, where
// nothing returns the carriage on its own.
func prompt(question string, newline string) {
	fmt.Fprintf(os.Stdout, "%s  %s%s%s%s  %s>%s ", newline, bold+cyan, question, reset, newline, gray, reset)
}

// retry reports one unusable answer and is followed by another attempt. A
// value that does not convert is not an error the contract reports: the
// question is simply asked again.
func retry(reason string) {
	fmt.Fprintf(os.Stdout, "  %s%s%s\n", yellow, reason, reset)
}

// strQuestion fills interviewer.Sandbox.StrQuestion, reading one line. An
// empty line is returned as an empty string, which is how a caller offering a
// default learns the default was taken.
func strQuestion(question string) (string, error) {
	return askText(question)
}

// intQuestion fills interviewer.Sandbox.IntQuestion, asking again until the
// answer converts.
func intQuestion(question string) (int, error) {
	for {
		answer, err := askText(question)
		if err != nil {
			return 0, err
		}

		value, convert_error := strconv.Atoi(strings.TrimSpace(answer))
		if convert_error != nil {
			retry(fmt.Sprintf("%q is not a whole number", answer))
			continue
		}

		return value, nil
	}
}

// floatQuestion fills interviewer.Sandbox.FloatQuestion, asking again until
// the answer converts.
func floatQuestion(question string) (float64, error) {
	for {
		answer, err := askText(question)
		if err != nil {
			return 0, err
		}

		value, convert_error := strconv.ParseFloat(strings.TrimSpace(answer), 64)
		if convert_error != nil {
			retry(fmt.Sprintf("%q is not a number", answer))
			continue
		}

		return value, nil
	}
}

// boolQuestion fills interviewer.Sandbox.BoolQuestion. It is asked as a
// two-row menu rather than as a y/n line, so a yes-or-no question is answered
// with the same keys every other question is.
func boolQuestion(question string) (bool, error) {
	chosen, err := runMenu(question, []interviewer.AlternativeOption{
		{Id: "no", Msg: "No"},
		{Id: "yes", Msg: "Yes"},
	}, false)
	if err != nil {
		return false, err
	}

	return len(chosen) == 1 && chosen[0] == 1, nil
}

// singleAlternativeQuestion fills
// interviewer.Sandbox.SingleAlternativeQuestion, returning the Id of the row
// that was chosen.
func singleAlternativeQuestion(question string, alternatives []interviewer.AlternativeOption) (string, error) {
	chosen, err := runMenu(question, alternatives, false)
	if err != nil {
		return "", err
	}
	if len(chosen) == 0 {
		return "", nil
	}

	return alternatives[chosen[0]].Id, nil
}

// multipleAlternativeQuestion fills
// interviewer.Sandbox.MultipleAlternativeQuestion, returning the Ids of every
// row ticked, in declaration order.
func multipleAlternativeQuestion(question string, alternatives []interviewer.AlternativeOption) ([]string, error) {
	chosen, err := runMenu(question, alternatives, true)
	if err != nil {
		return nil, err
	}

	ids := []string{}
	for _, index := range chosen {
		ids = append(ids, alternatives[index].Id)
	}

	return ids, nil
}

// Bind fills deps.Deps.Interviewer with the terminal interview below: an
// arrow-key menu over stdin in raw mode where the terminal allows it, and the
// same questions as numbered lists where it does not.
func Bind(deps *deps.Deps) {
	deps.Interviewer = interviewer.Sandbox{
		IntQuestion: func(question string) (int, error) {
			return intQuestion(question)
		},
		StrQuestion: func(question string) (string, error) {
			return strQuestion(question)
		},
		FloatQuestion: func(question string) (float64, error) {
			return floatQuestion(question)
		},
		BoolQuestion: func(question string) (bool, error) {
			return boolQuestion(question)
		},
		SingleAlternativeQuestion: func(question string, alternatives []interviewer.AlternativeOption) (string, error) {
			return singleAlternativeQuestion(question, alternatives)
		},
		MultipleAlternativeQuestion: func(question string, alternatives []interviewer.AlternativeOption) ([]string, error) {
			return multipleAlternativeQuestion(question, alternatives)
		},
		Back: func(err error) bool {
			return isBack(err)
		},
	}
}
